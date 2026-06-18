package handler

import (
	"crrc_pm_backend/config"
	"encoding/base64"
	"net/http"
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
			file_data LONGBLOB NULL,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
		) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4
	`)
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

	_, err = config.DB.Exec(`
		INSERT INTO uploaded_files (
			id,
			file_name,
			content_type,
			file_data,
			created_at,
			updated_at
		) VALUES (?, ?, ?, ?, ?, ?)
		ON DUPLICATE KEY UPDATE
			file_name = VALUES(file_name),
			content_type = VALUES(content_type),
			file_data = VALUES(file_data),
			updated_at = VALUES(updated_at)
	`, payload.FileID, payload.FileName, payload.FileContentType, data, time.Now(), time.Now())
	return err
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
	var data []byte

	err := config.DB.QueryRow(`
		SELECT
			IFNULL(file_name, ''),
			IFNULL(content_type, ''),
			file_data
		FROM uploaded_files
		WHERE id = ?
		LIMIT 1
	`, id).Scan(&fileName, &contentType, &data)
	if err != nil || len(data) == 0 {
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
