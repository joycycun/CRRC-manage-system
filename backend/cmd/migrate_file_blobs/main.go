package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	db, err := openDB()
	if err != nil {
		log.Fatalf("连接数据库失败: %v", err)
	}
	defer db.Close()

	if err := ensureColumns(db); err != nil {
		log.Fatalf("补齐文件路径字段失败: %v", err)
	}

	uploadedCount, err := migrateUploadedFiles(db)
	if err != nil {
		log.Fatalf("迁移 uploaded_files 失败: %v", err)
	}

	projectCount, err := migrateProjectProposalFiles(db)
	if err != nil {
		log.Fatalf("迁移项目立项书失败: %v", err)
	}

	log.Printf("文件迁移完成：uploaded_files=%d，project_proposals=%d", uploadedCount, projectCount)
}

func openDB() (*sql.DB, error) {
	user := getEnv("DB_USER", "crrc_user")
	password := getEnv("DB_PASSWORD", "123456")
	host := getEnv("DB_HOST", "127.0.0.1")
	port := getEnv("DB_PORT", "3306")
	dbName := getEnv("DB_NAME", "crrc_pm")

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", user, password, host, port, dbName)
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}
	return db, db.Ping()
}

func ensureColumns(db *sql.DB) error {
	statements := []string{
		"ALTER TABLE uploaded_files ADD COLUMN file_category VARCHAR(64) NOT NULL DEFAULT '' AFTER content_type",
		"ALTER TABLE uploaded_files ADD COLUMN file_path VARCHAR(512) NOT NULL DEFAULT '' AFTER file_category",
		"ALTER TABLE uploaded_files ADD COLUMN file_size BIGINT NOT NULL DEFAULT 0 AFTER file_path",
		"ALTER TABLE projects ADD COLUMN proposal_file_path VARCHAR(512) DEFAULT '' AFTER proposal_content_type",
		"ALTER TABLE projects ADD COLUMN proposal_file_size BIGINT NOT NULL DEFAULT 0 AFTER proposal_file_path",
	}

	for _, statement := range statements {
		if _, err := db.Exec(statement); err != nil && !strings.Contains(err.Error(), "Duplicate column name") {
			return err
		}
	}
	return nil
}

func migrateUploadedFiles(db *sql.DB) (int, error) {
	rows, err := db.Query(`
		SELECT
			id,
			IFNULL(file_name, ''),
			IFNULL(content_type, ''),
			file_data
		FROM uploaded_files
		WHERE file_data IS NOT NULL
		  AND OCTET_LENGTH(file_data) > 0
	`)
	if err != nil {
		return 0, err
	}
	defer rows.Close()

	count := 0
	for rows.Next() {
		var id int64
		var fileName string
		var contentType string
		var data []byte
		if err := rows.Scan(&id, &fileName, &contentType, &data); err != nil {
			return count, err
		}

		filePath, category, size, err := writeFile(id, fileName, contentType, data, "")
		if err != nil {
			return count, err
		}

		if _, err := db.Exec(`
			UPDATE uploaded_files
			SET
				file_category = ?,
				file_path = ?,
				file_size = ?,
				file_data = NULL,
				updated_at = NOW()
			WHERE id = ?
		`, category, filePath, size, id); err != nil {
			return count, err
		}
		count++
	}
	return count, rows.Err()
}

func migrateProjectProposalFiles(db *sql.DB) (int, error) {
	rows, err := db.Query(`
		SELECT
			id,
			IFNULL(proposal_file_name, ''),
			IFNULL(proposal_content_type, ''),
			proposal_file_data
		FROM projects
		WHERE proposal_file_data IS NOT NULL
		  AND OCTET_LENGTH(proposal_file_data) > 0
	`)
	if err != nil {
		return 0, err
	}
	defer rows.Close()

	count := 0
	for rows.Next() {
		var id int64
		var fileName string
		var contentType string
		var data []byte
		if err := rows.Scan(&id, &fileName, &contentType, &data); err != nil {
			return count, err
		}

		filePath, _, size, err := writeFile(id, fileName, contentType, data, "projects")
		if err != nil {
			return count, err
		}

		if _, err := db.Exec(`
			UPDATE projects
			SET
				proposal_file_path = ?,
				proposal_file_size = ?,
				proposal_file_data = NULL,
				updated_at = NOW()
			WHERE id = ?
		`, filePath, size, id); err != nil {
			return count, err
		}
		count++
	}
	return count, rows.Err()
}

func writeFile(id int64, fileName string, contentType string, data []byte, subdir string) (string, string, int64, error) {
	category := fileCategory(fileName, contentType)
	dirParts := []string{uploadRoot(), category}
	if strings.TrimSpace(subdir) != "" {
		dirParts = append(dirParts, strings.Trim(strings.ReplaceAll(subdir, "\\", "/"), "/"))
	}
	dir := filepath.Join(dirParts...)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return "", "", 0, err
	}

	path := filepath.Join(dir, strconv.FormatInt(id, 10)+"_"+safeFileName(fileName))
	if err := os.WriteFile(path, data, 0644); err != nil {
		return "", "", 0, err
	}
	return path, category, int64(len(data)), nil
}

func fileCategory(fileName string, contentType string) string {
	name := strings.ToLower(fileName)
	contentType = strings.ToLower(contentType)

	switch {
	case strings.HasSuffix(name, ".doc") || strings.HasSuffix(name, ".docx") || strings.Contains(contentType, "word"):
		return "documents"
	case strings.HasSuffix(name, ".txt") || strings.Contains(contentType, "text/plain"):
		return "texts"
	case strings.HasSuffix(name, ".pdf") || strings.Contains(contentType, "pdf"):
		return "pdf"
	case strings.HasSuffix(name, ".xls") || strings.HasSuffix(name, ".xlsx") || strings.Contains(contentType, "spreadsheet") || strings.Contains(contentType, "excel"):
		return "spreadsheets"
	case strings.HasSuffix(name, ".zip") || strings.HasSuffix(name, ".rar") || strings.HasSuffix(name, ".7z") || strings.Contains(contentType, "zip") || strings.Contains(contentType, "rar"):
		return "archives"
	case strings.HasSuffix(name, ".png") || strings.HasSuffix(name, ".jpg") || strings.HasSuffix(name, ".jpeg") || strings.HasSuffix(name, ".gif") || strings.HasSuffix(name, ".webp") || strings.HasPrefix(contentType, "image/"):
		return "images"
	default:
		return "others"
	}
}

func safeFileName(fileName string) string {
	fileName = strings.TrimSpace(filepath.Base(fileName))
	if fileName == "" || fileName == "." {
		fileName = "file"
	}
	re := regexp.MustCompile(`[^A-Za-z0-9._\-\p{Han}]`)
	return re.ReplaceAllString(fileName, "_")
}

func uploadRoot() string {
	return getEnv("UPLOAD_ROOT", "/app/uploads")
}

func getEnv(key string, fallback string) string {
	if value := strings.TrimSpace(os.Getenv(key)); value != "" {
		return value
	}
	return fallback
}
