package handler

import (
	"crrc_pm_backend/config"
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type HardwareDevDocument struct {
	ID             int64        `json:"id"`
	ProjectID      int64        `json:"projectId"`
	ProjectName    string       `json:"projectName"`
	FileID         int64        `json:"fileId"`
	FileName       string       `json:"fileName"`
	FileSize       int64        `json:"fileSize"`
	FileURL        string       `json:"fileUrl"`
	DownloadURL    string       `json:"downloadUrl"`
	UploadUserID   int64        `json:"uploadUserId"`
	UploadUserName string       `json:"uploadUserName"`
	UploadTime     sql.NullTime `json:"uploadTime"`
	Remark         string       `json:"remark"`
	CreatedAt      time.Time    `json:"createdAt"`
	UpdatedAt      time.Time    `json:"updatedAt"`
	IsDeleted      bool         `json:"isDeleted"`
}

func HardwareDevDocumentsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	switch r.Method {
	case http.MethodGet:
		GetHardwareDevDocumentsHandler(w, r)
	case http.MethodPost:
		CreateHardwareDevDocumentHandler(w, r)
	default:
		http.Error(w, "不支持该请求方法", http.StatusMethodNotAllowed)
	}
}

func HardwareDevDocumentActionHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	path := strings.TrimPrefix(r.URL.Path, "/api/hardware-dev-documents/")
	path = strings.Trim(path, "/")
	id, err := strconv.ParseInt(path, 10, 64)
	if err != nil || id == 0 {
		http.Error(w, "开发文档ID错误", http.StatusBadRequest)
		return
	}

	if r.Method == http.MethodDelete {
		DeleteHardwareDevDocumentHandler(w, r, id)
		return
	}

	http.Error(w, "接口不存在", http.StatusNotFound)
}

func ensureHardwareDevDocumentsTable() {
	ensureUploadedFilesTable()
	_, _ = config.DB.Exec(`
		CREATE TABLE IF NOT EXISTS hardware_dev_documents (
			id BIGINT PRIMARY KEY AUTO_INCREMENT,
			project_id BIGINT NOT NULL DEFAULT 0,
			file_id BIGINT NOT NULL DEFAULT 0,
			upload_user_id BIGINT NOT NULL DEFAULT 0,
			upload_user_name VARCHAR(64) NOT NULL DEFAULT '',
			upload_time DATETIME DEFAULT CURRENT_TIMESTAMP,
			remark TEXT NULL,
			is_deleted TINYINT NOT NULL DEFAULT 0,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
			KEY idx_hdd_project_id (project_id),
			KEY idx_hdd_file_id (file_id)
		) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4
	`)
}

func GetHardwareDevDocumentsHandler(w http.ResponseWriter, r *http.Request) {
	ensureHardwareDevDocumentsTable()

	rows, err := config.DB.Query(`
		SELECT
			hdd.id,
			hdd.project_id,
			IFNULL(p.project_name, ''),
			hdd.file_id,
			IFNULL(uf.file_name, ''),
			IFNULL(uf.file_size, 0),
			IFNULL(hdd.upload_user_id, 0),
			IFNULL(hdd.upload_user_name, ''),
			hdd.upload_time,
			IFNULL(hdd.remark, ''),
			hdd.created_at,
			hdd.updated_at,
			hdd.is_deleted
		FROM hardware_dev_documents hdd
		LEFT JOIN projects p ON p.id = hdd.project_id
		LEFT JOIN uploaded_files uf ON uf.id = hdd.file_id
		WHERE hdd.is_deleted = 0
		ORDER BY hdd.id DESC
	`)
	if err != nil {
		http.Error(w, "查询开发文档失败: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	list := make([]HardwareDevDocument, 0)
	for rows.Next() {
		var item HardwareDevDocument
		if err := rows.Scan(
			&item.ID,
			&item.ProjectID,
			&item.ProjectName,
			&item.FileID,
			&item.FileName,
			&item.FileSize,
			&item.UploadUserID,
			&item.UploadUserName,
			&item.UploadTime,
			&item.Remark,
			&item.CreatedAt,
			&item.UpdatedAt,
			&item.IsDeleted,
		); err != nil {
			http.Error(w, "解析开发文档失败: "+err.Error(), http.StatusInternalServerError)
			return
		}
		item.FileURL = filePreviewURL(item.FileID)
		item.DownloadURL = fileDownloadURL(item.FileID)
		list = append(list, item)
	}

	json.NewEncoder(w).Encode(map[string]interface{}{"code": 200, "msg": "查询成功", "data": list})
}

func CreateHardwareDevDocumentHandler(w http.ResponseWriter, r *http.Request) {
	if !hasRequestRole(r, "hardware_owner") && !hasRequestPermission(r, "hardware-dev-doc:upload") {
		http.Error(w, "无上传开发文档权限", http.StatusForbidden)
		return
	}

	ensureHardwareDevDocumentsTable()

	var req struct {
		ProjectID       int64  `json:"projectId"`
		FileID          int64  `json:"fileId"`
		FileName        string `json:"fileName"`
		FileContentType string `json:"fileContentType"`
		FileData        string `json:"fileData"`
		UploadUserID    int64  `json:"uploadUserId"`
		UploadUserName  string `json:"uploadUserName"`
		Remark          string `json:"remark"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "参数解析失败: "+err.Error(), http.StatusBadRequest)
		return
	}
	if req.ProjectID == 0 {
		http.Error(w, "项目ID不能为空", http.StatusBadRequest)
		return
	}
	if req.FileID == 0 || strings.TrimSpace(req.FileName) == "" {
		http.Error(w, "请上传开发文档文件", http.StatusBadRequest)
		return
	}

	if err := saveUploadedFile(UploadedFilePayload{
		FileID:          req.FileID,
		FileName:        req.FileName,
		FileContentType: req.FileContentType,
		FileData:        req.FileData,
	}); err != nil {
		http.Error(w, "保存开发文档文件失败: "+err.Error(), http.StatusBadRequest)
		return
	}

	now := time.Now()
	result, err := config.DB.Exec(`
		INSERT INTO hardware_dev_documents (
			project_id,
			file_id,
			upload_user_id,
			upload_user_name,
			upload_time,
			remark,
			created_at,
			updated_at,
			is_deleted
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, 0)
	`, req.ProjectID, req.FileID, req.UploadUserID, req.UploadUserName, now, req.Remark, now, now)
	if err != nil {
		http.Error(w, "新增开发文档失败: "+err.Error(), http.StatusInternalServerError)
		return
	}

	id, _ := result.LastInsertId()
	json.NewEncoder(w).Encode(map[string]interface{}{"code": 200, "msg": "新增成功", "data": map[string]interface{}{"id": id}})
}

func DeleteHardwareDevDocumentHandler(w http.ResponseWriter, r *http.Request, id int64) {
	if !hasRequestRole(r, "hardware_owner") && !hasRequestPermission(r, "hardware-dev-doc:delete") {
		http.Error(w, "无删除开发文档权限", http.StatusForbidden)
		return
	}

	result, err := config.DB.Exec(`
		UPDATE hardware_dev_documents
		SET is_deleted = 1,
			updated_at = NOW()
		WHERE id = ?
	`, id)
	if err != nil {
		http.Error(w, "删除开发文档失败: "+err.Error(), http.StatusInternalServerError)
		return
	}
	if affected, _ := result.RowsAffected(); affected == 0 {
		http.Error(w, "开发文档不存在", http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(map[string]interface{}{"code": 200, "msg": "删除成功"})
}
