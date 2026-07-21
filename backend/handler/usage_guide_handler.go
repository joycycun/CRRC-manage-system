package handler

import (
	"crrc_pm_backend/config"
	"encoding/json"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"
)

type usageGuideFilePayload struct {
	FileID          int64  `json:"fileId"`
	FileName        string `json:"fileName"`
	FileContentType string `json:"fileContentType"`
	FileData        string `json:"fileData"`
}

func ensureUsageGuideTables() {
	_, _ = config.DB.Exec(`
		CREATE TABLE IF NOT EXISTS usage_guides (
			id BIGINT PRIMARY KEY AUTO_INCREMENT,
			title VARCHAR(128) NOT NULL DEFAULT '',
			description TEXT,
			created_by BIGINT NOT NULL DEFAULT 0,
			created_by_name VARCHAR(64) NOT NULL DEFAULT '',
			is_deleted TINYINT NOT NULL DEFAULT 0,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
		) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4
	`)
	_, _ = config.DB.Exec(`
		CREATE TABLE IF NOT EXISTS usage_guide_images (
			id BIGINT PRIMARY KEY AUTO_INCREMENT,
			guide_id BIGINT NOT NULL,
			file_id BIGINT NOT NULL,
			sort_order INT NOT NULL DEFAULT 0,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			KEY idx_usage_guide_images_guide (guide_id)
		) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4
	`)
	_, _ = config.DB.Exec(`
		CREATE TABLE IF NOT EXISTS usage_guide_files (
			id BIGINT PRIMARY KEY AUTO_INCREMENT,
			guide_id BIGINT NOT NULL,
			file_id BIGINT NOT NULL,
			sort_order INT NOT NULL DEFAULT 0,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			UNIQUE KEY uk_usage_guide_file (guide_id, file_id),
			KEY idx_usage_guide_files_guide (guide_id)
		) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4
	`)
	_, _ = config.DB.Exec(`
		INSERT IGNORE INTO usage_guide_files (guide_id, file_id, sort_order, created_at)
		SELECT guide_id, file_id, sort_order, created_at FROM usage_guide_images
	`)
}

func UsageGuidesHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	ensureUsageGuideTables()

	switch r.Method {
	case http.MethodGet:
		getUsageGuides(w)
	case http.MethodPost:
		if !hasRequestRole(r, "system_admin") {
			http.Error(w, "只有系统管理员可以新增使用说明", http.StatusForbidden)
			return
		}
		createUsageGuide(w, r)
	default:
		http.Error(w, "不支持该请求方法", http.StatusMethodNotAllowed)
	}
}

func UsageGuideActionHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if !hasRequestRole(r, "system_admin") {
		http.Error(w, "只有系统管理员可以使用使用说明管理", http.StatusForbidden)
		return
	}
	if r.Method != http.MethodDelete {
		http.Error(w, "不支持该请求方法", http.StatusMethodNotAllowed)
		return
	}
	idText := strings.Trim(strings.TrimPrefix(r.URL.Path, "/api/usage-guides/"), "/")
	id, err := strconv.ParseInt(idText, 10, 64)
	if err != nil || id <= 0 {
		http.Error(w, "使用说明ID错误", http.StatusBadRequest)
		return
	}
	ensureUsageGuideTables()
	result, err := config.DB.Exec("UPDATE usage_guides SET is_deleted = 1, updated_at = NOW() WHERE id = ? AND is_deleted = 0", id)
	if err != nil {
		http.Error(w, "删除使用说明失败: "+err.Error(), http.StatusInternalServerError)
		return
	}
	affected, _ := result.RowsAffected()
	if affected == 0 {
		http.Error(w, "使用说明不存在", http.StatusNotFound)
		return
	}
	json.NewEncoder(w).Encode(map[string]interface{}{"code": 200, "msg": "删除成功"})
}

func getUsageGuides(w http.ResponseWriter) {
	rows, err := config.DB.Query(`
		SELECT id, title, IFNULL(description, ''), IFNULL(created_by_name, ''),
		       DATE_FORMAT(created_at, '%Y-%m-%d %H:%i:%s')
		FROM usage_guides
		WHERE is_deleted = 0
		ORDER BY updated_at DESC, id DESC
	`)
	if err != nil {
		http.Error(w, "查询使用说明失败: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	type guideFile struct {
		ID          int64  `json:"id"`
		FileID      int64  `json:"fileId"`
		FileName    string `json:"fileName"`
		ContentType string `json:"contentType"`
		FileURL     string `json:"fileUrl"`
		DownloadURL string `json:"downloadUrl"`
		SortOrder   int    `json:"sortOrder"`
	}
	type guide struct {
		ID          int64       `json:"id"`
		Title       string      `json:"title"`
		Description string      `json:"description"`
		CreatedBy   string      `json:"createdBy"`
		CreatedAt   string      `json:"createdAt"`
		Files       []guideFile `json:"files"`
	}

	list := make([]guide, 0)
	for rows.Next() {
		var item guide
		item.Files = make([]guideFile, 0)
		if err := rows.Scan(&item.ID, &item.Title, &item.Description, &item.CreatedBy, &item.CreatedAt); err != nil {
			http.Error(w, "解析使用说明失败: "+err.Error(), http.StatusInternalServerError)
			return
		}
		fileRows, fileErr := config.DB.Query(`
			SELECT ugf.id, ugf.file_id, IFNULL(uf.file_name, ''), IFNULL(uf.content_type, ''), ugf.sort_order
			FROM usage_guide_files ugf
			LEFT JOIN uploaded_files uf ON uf.id = ugf.file_id
			WHERE ugf.guide_id = ?
			ORDER BY ugf.sort_order, ugf.id
		`, item.ID)
		if fileErr != nil {
			http.Error(w, "查询说明附件失败: "+fileErr.Error(), http.StatusInternalServerError)
			return
		}
		for fileRows.Next() {
			var file guideFile
			if err := fileRows.Scan(&file.ID, &file.FileID, &file.FileName, &file.ContentType, &file.SortOrder); err != nil {
				fileRows.Close()
				http.Error(w, "解析说明附件失败: "+err.Error(), http.StatusInternalServerError)
				return
			}
			file.FileURL = filePreviewURL(file.FileID)
			file.DownloadURL = fileDownloadURL(file.FileID)
			item.Files = append(item.Files, file)
		}
		fileRows.Close()
		list = append(list, item)
	}
	json.NewEncoder(w).Encode(map[string]interface{}{"code": 200, "msg": "查询成功", "data": list})
}

func createUsageGuide(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Title       string                  `json:"title"`
		Description string                  `json:"description"`
		Files       []usageGuideFilePayload `json:"files"`
		Images      []usageGuideFilePayload `json:"images"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "参数解析失败: "+err.Error(), http.StatusBadRequest)
		return
	}
	req.Title = strings.TrimSpace(req.Title)
	if req.Title == "" {
		http.Error(w, "说明名称不能为空", http.StatusBadRequest)
		return
	}
	if len(req.Files) == 0 {
		req.Files = req.Images
	}
	if len(req.Files) == 0 {
		http.Error(w, "请至少上传一个说明文件", http.StatusBadRequest)
		return
	}
	for _, file := range req.Files {
		if file.FileID == 0 || !isUsageGuideFileAllowed(file.FileName) {
			http.Error(w, "使用说明支持 Word、PDF、TXT、Excel、PPT 和常用图片文件", http.StatusBadRequest)
			return
		}
		if err := saveUploadedFile(UploadedFilePayload{
			FileID: file.FileID, FileName: file.FileName,
			FileContentType: file.FileContentType, FileData: file.FileData,
		}); err != nil {
			http.Error(w, "保存说明文件失败: "+err.Error(), http.StatusBadRequest)
			return
		}
	}

	userID, userName := currentRequestUser(r)
	tx, err := config.DB.Begin()
	if err != nil {
		http.Error(w, "开启事务失败: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer tx.Rollback()
	result, err := tx.Exec(`
		INSERT INTO usage_guides (title, description, created_by, created_by_name, created_at, updated_at)
		VALUES (?, ?, ?, ?, NOW(), NOW())
	`, req.Title, strings.TrimSpace(req.Description), userID, userName)
	if err != nil {
		http.Error(w, "新增使用说明失败: "+err.Error(), http.StatusInternalServerError)
		return
	}
	guideID, _ := result.LastInsertId()
	for index, file := range req.Files {
		if _, err := tx.Exec(`
			INSERT INTO usage_guide_files (guide_id, file_id, sort_order, created_at)
			VALUES (?, ?, ?, NOW())
		`, guideID, file.FileID, index+1); err != nil {
			http.Error(w, "保存说明文件关系失败: "+err.Error(), http.StatusInternalServerError)
			return
		}
	}
	if err := tx.Commit(); err != nil {
		http.Error(w, "提交使用说明失败: "+err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(map[string]interface{}{"code": 200, "msg": "新增成功", "data": map[string]int64{"id": guideID}})
}

func isUsageGuideFileAllowed(fileName string) bool {
	switch strings.ToLower(filepath.Ext(strings.TrimSpace(fileName))) {
	case ".doc", ".docx", ".pdf", ".txt", ".xls", ".xlsx", ".ppt", ".pptx", ".png", ".jpg", ".jpeg", ".gif", ".webp":
		return true
	default:
		return false
	}
}
