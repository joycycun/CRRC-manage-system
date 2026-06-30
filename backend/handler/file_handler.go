package handler

import (
	"crrc_pm_backend/config"
	"encoding/base64"
	"fmt"
	"io/fs"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type UploadedFilePayload struct {
	FileID          int64  `json:"fileId"`
	FileName        string `json:"fileName"`
	FileURL         string `json:"fileUrl"`
	FileContentType string `json:"fileContentType"`
	FileData        string `json:"fileData"`
}

func ensureUploadedFilesTable() {
	_, _ = config.DB.Exec(`
		CREATE TABLE IF NOT EXISTS uploaded_files (
			id BIGINT PRIMARY KEY,
			file_name VARCHAR(255) NOT NULL DEFAULT '',
			content_type VARCHAR(128) NOT NULL DEFAULT '',
			file_category VARCHAR(64) NOT NULL DEFAULT '',
			file_path VARCHAR(512) NOT NULL DEFAULT '',
			file_size BIGINT NOT NULL DEFAULT 0,
			file_data LONGBLOB NULL,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
		) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4
	`)
	_, _ = config.DB.Exec(`ALTER TABLE uploaded_files ADD COLUMN file_category VARCHAR(64) NOT NULL DEFAULT '' AFTER content_type`)
	_, _ = config.DB.Exec(`ALTER TABLE uploaded_files ADD COLUMN file_path VARCHAR(512) NOT NULL DEFAULT '' AFTER file_category`)
	_, _ = config.DB.Exec(`ALTER TABLE uploaded_files ADD COLUMN file_size BIGINT NOT NULL DEFAULT 0 AFTER file_path`)
}

func decodeDataURL(value string) ([]byte, error) {
	value = strings.TrimSpace(value)
	if idx := strings.Index(value, ","); idx >= 0 {
		value = value[idx+1:]
	}
	return base64.StdEncoding.DecodeString(value)
}

func saveUploadedFile(payload UploadedFilePayload) error {
	if payload.FileID == 0 || strings.TrimSpace(payload.FileData) == "" {
		return nil
	}

	ensureUploadedFilesTable()

	data, err := decodeDataURL(payload.FileData)
	if err != nil {
		return err
	}

	if payload.FileName == "" {
		payload.FileName = "未命名文件"
	}
	if payload.FileContentType == "" {
		payload.FileContentType = "application/octet-stream"
	}

	filePath, category, fileSize, err := saveFileBytesToUploads(payload.FileID, payload.FileName, payload.FileContentType, data)
	if err != nil {
		return err
	}

	_, err = config.DB.Exec(`
		INSERT INTO uploaded_files (
			id,
			file_name,
			content_type,
			file_category,
			file_path,
			file_size,
			file_data,
			created_at,
			updated_at
		) VALUES (?, ?, ?, ?, ?, ?, NULL, ?, ?)
		ON DUPLICATE KEY UPDATE
			file_name = VALUES(file_name),
			content_type = VALUES(content_type),
			file_category = VALUES(file_category),
			file_path = VALUES(file_path),
			file_size = VALUES(file_size),
			file_data = NULL,
			updated_at = VALUES(updated_at)
	`, payload.FileID, payload.FileName, payload.FileContentType, category, filePath, fileSize, time.Now(), time.Now())
	return err
}

func uploadRoot() string {
	root := strings.TrimSpace(os.Getenv("UPLOAD_ROOT"))
	if root == "" {
		root = "/app/uploads"
	}
	return root
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
	if fileName == "" || fileName == "." || fileName == string(filepath.Separator) {
		fileName = "file"
	}
	re := regexp.MustCompile(`[^A-Za-z0-9._\-\p{Han}]`)
	fileName = re.ReplaceAllString(fileName, "_")
	return fileName
}

func saveFileBytesToUploads(fileID int64, fileName string, contentType string, data []byte) (string, string, int64, error) {
	category := fileCategory(fileName, contentType)
	dir := filepath.Join(uploadRoot(), category)
	if err := os.MkdirAll(dir, fs.FileMode(0755)); err != nil {
		return "", "", 0, err
	}

	name := fmt.Sprintf("%d_%s", fileID, safeFileName(fileName))
	path := filepath.Join(dir, name)
	if err := os.WriteFile(path, data, 0644); err != nil {
		return "", "", 0, err
	}

	return path, category, int64(len(data)), nil
}

func saveNamedFileBytesToUploads(subdir string, fileName string, contentType string, data []byte, prefix string) (string, string, int64, error) {
	category := fileCategory(fileName, contentType)
	parts := []string{uploadRoot(), category}
	if subdir != "" {
		parts = append(parts, strings.Trim(strings.ReplaceAll(subdir, "\\", "/"), "/"))
	}
	dir := filepath.Join(parts...)
	if err := os.MkdirAll(dir, fs.FileMode(0755)); err != nil {
		return "", "", 0, err
	}

	if prefix == "" {
		prefix = strconv.FormatInt(time.Now().UnixNano(), 10)
	}
	path := filepath.Join(dir, prefix+"_"+safeFileName(fileName))
	if err := os.WriteFile(path, data, 0644); err != nil {
		return "", "", 0, err
	}

	return path, category, int64(len(data)), nil
}

func filePreviewURL(fileID int64) string {
	if fileID == 0 {
		return ""
	}
	return "/api/files/" + strconv.FormatInt(fileID, 10) + "/preview"
}

func fileDownloadURL(fileID int64) string {
	if fileID == 0 {
		return ""
	}
	return "/api/files/" + strconv.FormatInt(fileID, 10) + "/download"
}

func UploadedFileActionHandler(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimPrefix(r.URL.Path, "/api/files/")
	path = strings.Trim(path, "/")
	parts := strings.Split(path, "/")

	if len(parts) != 2 || parts[0] == "" {
		http.Error(w, "缺少文件ID", http.StatusBadRequest)
		return
	}

	id, err := strconv.ParseInt(parts[0], 10, 64)
	if err != nil || id == 0 {
		http.Error(w, "文件ID错误", http.StatusBadRequest)
		return
	}

	switch parts[1] {
	case "preview":
		serveUploadedFile(w, id, true)
	case "download":
		serveUploadedFile(w, id, false)
	default:
		http.Error(w, "接口不存在", http.StatusNotFound)
	}
}

func serveUploadedFile(w http.ResponseWriter, id int64, inline bool) {
	ensureUploadedFilesTable()

	var fileName string
	var contentType string
	var filePath string
	var data []byte

	err := config.DB.QueryRow(`
		SELECT
			IFNULL(file_name, ''),
			IFNULL(content_type, ''),
			IFNULL(file_path, ''),
			file_data
		FROM uploaded_files
		WHERE id = ?
		LIMIT 1
	`, id).Scan(&fileName, &contentType, &filePath, &data)
	if err != nil {
		http.Error(w, "文件不存在或尚未保存内容", http.StatusNotFound)
		return
	}

	if strings.TrimSpace(filePath) != "" {
		diskData, readErr := os.ReadFile(filePath)
		if readErr == nil && len(diskData) > 0 {
			data = diskData
		}
	}

	if len(data) == 0 {
		http.Error(w, "文件不存在或尚未保存内容", http.StatusNotFound)
		return
	}

	if contentType == "" {
		contentType = "application/octet-stream"
	}
	if fileName == "" {
		fileName = "file"
	}

	disposition := "attachment"
	if inline {
		disposition = "inline"
	}

	w.Header().Set("Content-Type", contentType)
	w.Header().Set("Content-Disposition", disposition+`; filename="`+strings.ReplaceAll(fileName, `"`, "")+`"`)
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(data)
}
