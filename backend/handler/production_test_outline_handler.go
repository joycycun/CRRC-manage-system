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

type ProductionTestOutlineVO struct {
	ID              int64  `json:"id"`
	ProjectID       int64  `json:"projectId"`
	ProjectName     string `json:"projectName"`
	HardwareID      int64  `json:"hardwareId"`
	HardwareVersion string `json:"hardwareVersion"`
	BoardModels     string `json:"boardModels"`
	DeviceType      string `json:"deviceType"`
	FileID          int64  `json:"fileId"`
	FileName        string `json:"fileName"`
	FileURL         string `json:"fileUrl"`
	DownloadURL     string `json:"downloadUrl"`
	UploaderID      int64  `json:"uploaderId"`
	UploaderName    string `json:"uploaderName"`
	UploadTime      string `json:"uploadTime"`
	Remark          string `json:"remark"`
	CreatedAt       string `json:"createdAt"`
	UpdatedAt       string `json:"updatedAt"`
}

func ProductionTestOutlinesHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	switch r.Method {
	case http.MethodGet:
		GetProductionTestOutlinesHandler(w, r)
	case http.MethodPost:
		CreateProductionTestOutlineHandler(w, r)
	default:
		http.Error(w, "不支持该请求方法", http.StatusMethodNotAllowed)
	}
}

func ProductionTestOutlineActionHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	path := strings.TrimPrefix(r.URL.Path, "/api/production-test-outlines/")
	path = strings.Trim(path, "/")
	if path == "" {
		http.Error(w, "缺少生产测试大纲ID", http.StatusBadRequest)
		return
	}

	parts := strings.Split(path, "/")
	if len(parts) != 1 {
		http.Error(w, "接口不存在", http.StatusNotFound)
		return
	}

	id, err := strconv.ParseInt(parts[0], 10, 64)
	if err != nil || id == 0 {
		http.Error(w, "生产测试大纲ID错误", http.StatusBadRequest)
		return
	}

	if r.Method == http.MethodDelete {
		DeleteProductionTestOutlineHandler(w, r, id)
		return
	}

	http.Error(w, "不支持该请求方法", http.StatusMethodNotAllowed)
}

func ensureProductionTestOutlinesTable() {
	ensureUploadedFilesTable()
	_, _ = config.DB.Exec(`
		CREATE TABLE IF NOT EXISTS production_test_outlines (
			id BIGINT PRIMARY KEY AUTO_INCREMENT,
			project_id BIGINT NOT NULL DEFAULT 0,
			hardware_id BIGINT NOT NULL DEFAULT 0,
			hardware_version VARCHAR(128) NOT NULL DEFAULT '',
			board_models VARCHAR(512) NOT NULL DEFAULT '',
			device_type VARCHAR(128) NOT NULL DEFAULT '',
			file_id BIGINT NOT NULL DEFAULT 0,
			file_name VARCHAR(255) NOT NULL DEFAULT '',
			uploader_id BIGINT NOT NULL DEFAULT 0,
			uploader_name VARCHAR(64) NOT NULL DEFAULT '',
			upload_time DATETIME DEFAULT CURRENT_TIMESTAMP,
			remark TEXT NULL,
			is_deleted TINYINT NOT NULL DEFAULT 0,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
			KEY idx_pto_project_id (project_id),
			KEY idx_pto_hardware_id (hardware_id),
			KEY idx_pto_upload_time (upload_time)
		) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4
	`)
	_, _ = config.DB.Exec(`
		ALTER TABLE production_test_outlines
		ADD COLUMN project_id BIGINT NOT NULL DEFAULT 0 AFTER id
	`)
	_, _ = config.DB.Exec(`
		ALTER TABLE production_test_outlines
		ADD COLUMN board_models VARCHAR(512) NOT NULL DEFAULT '' AFTER hardware_version
	`)
	_, _ = config.DB.Exec(`ALTER TABLE production_test_outlines ADD KEY idx_pto_project_id (project_id)`)
}

func splitBoardModels(value string) []string {
	parts := strings.FieldsFunc(value, func(r rune) bool {
		return r == ',' || r == '，'
	})

	list := make([]string, 0, len(parts))
	seen := map[string]bool{}
	for _, part := range parts {
		model := strings.TrimSpace(part)
		if model == "" || seen[model] {
			continue
		}
		seen[model] = true
		list = append(list, model)
	}
	return list
}

func canViewProductionTestOutline(r *http.Request) bool {
	return hasRequestRole(r, "hardware_owner") ||
		hasRequestRole(r, "production_staff") ||
		hasRequestRole(r, "quality_staff") ||
		hasRequestRole(r, "system_admin") ||
		hasRequestRole(r, "leader")
}

func GetProductionTestOutlinesHandler(w http.ResponseWriter, r *http.Request) {
	if !canViewProductionTestOutline(r) {
		http.Error(w, "无生产测试大纲查看权限", http.StatusForbidden)
		return
	}

	ensureProductionTestOutlinesTable()

	query := `
		SELECT
			pto.id,
			IFNULL(pto.project_id, 0),
			IFNULL(p.project_name, '未绑定项目'),
			pto.hardware_id,
			pto.hardware_version,
			IFNULL(pto.board_models, ''),
			pto.device_type,
			pto.file_id,
			IFNULL(uf.file_name, pto.file_name),
			pto.uploader_id,
			pto.uploader_name,
			pto.upload_time,
			IFNULL(pto.remark, ''),
			pto.created_at,
			pto.updated_at
		FROM production_test_outlines pto
		LEFT JOIN projects p ON p.id = pto.project_id
		LEFT JOIN uploaded_files uf ON uf.id = pto.file_id
		WHERE pto.is_deleted = 0
		ORDER BY IFNULL(p.project_name, '未绑定项目') ASC, pto.board_models ASC, pto.id DESC
	`

	if hasRequestRole(r, "production_staff") && !hasRequestRole(r, "hardware_owner") && !hasRequestRole(r, "system_admin") && !hasRequestRole(r, "leader") {
		query = `
			SELECT
				pto.id,
				IFNULL(pto.project_id, 0),
				IFNULL(p.project_name, '未绑定项目'),
				pto.hardware_id,
				pto.hardware_version,
				IFNULL(pto.board_models, ''),
				pto.device_type,
				pto.file_id,
				IFNULL(uf.file_name, pto.file_name),
				pto.uploader_id,
				pto.uploader_name,
				pto.upload_time,
				IFNULL(pto.remark, ''),
				pto.created_at,
				pto.updated_at
			FROM production_test_outlines pto
			INNER JOIN (
				SELECT IFNULL(project_id, 0) AS project_id, IFNULL(board_models, '') AS board_models, MAX(id) AS latest_id
				FROM production_test_outlines
				WHERE is_deleted = 0
				GROUP BY IFNULL(project_id, 0), IFNULL(board_models, '')
			) latest ON latest.latest_id = pto.id
			LEFT JOIN projects p ON p.id = pto.project_id
			LEFT JOIN uploaded_files uf ON uf.id = pto.file_id
			WHERE pto.is_deleted = 0
			ORDER BY IFNULL(p.project_name, '未绑定项目') ASC, pto.board_models ASC, pto.id DESC
		`
	}

	rows, err := config.DB.Query(query)
	if err != nil {
		http.Error(w, "查询生产测试大纲失败: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	list := make([]ProductionTestOutlineVO, 0)
	for rows.Next() {
		var item ProductionTestOutlineVO
		var uploadTime, createdAt, updatedAt sql.NullTime

		if err := rows.Scan(
			&item.ID,
			&item.ProjectID,
			&item.ProjectName,
			&item.HardwareID,
			&item.HardwareVersion,
			&item.BoardModels,
			&item.DeviceType,
			&item.FileID,
			&item.FileName,
			&item.UploaderID,
			&item.UploaderName,
			&uploadTime,
			&item.Remark,
			&createdAt,
			&updatedAt,
		); err != nil {
			http.Error(w, "数据解析失败: "+err.Error(), http.StatusInternalServerError)
			return
		}

		item.FileURL = filePreviewURL(item.FileID)
		item.DownloadURL = fileDownloadURL(item.FileID)
		item.UploadTime = nullTime(uploadTime)
		item.CreatedAt = nullTime(createdAt)
		item.UpdatedAt = nullTime(updatedAt)
		list = append(list, item)
	}

	if err := rows.Err(); err != nil {
		http.Error(w, "数据读取失败: "+err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"code": 200,
		"msg":  "查询成功",
		"data": list,
	})
}

func CreateProductionTestOutlineHandler(w http.ResponseWriter, r *http.Request) {
	if !requireHardwareOwnerPermission(w, r) {
		return
	}

	ensureProductionTestOutlinesTable()

	var req struct {
		ProjectID       int64  `json:"projectId"`
		HardwareID      int64  `json:"hardwareId"`
		BoardModels     string `json:"boardModels"`
		FileID          int64  `json:"fileId"`
		FileName        string `json:"fileName"`
		FileContentType string `json:"fileContentType"`
		FileData        string `json:"fileData"`
		UploaderID      int64  `json:"uploaderId"`
		UploaderName    string `json:"uploaderName"`
		Remark          string `json:"remark"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "参数解析失败: "+err.Error(), http.StatusBadRequest)
		return
	}
	if req.FileID == 0 || strings.TrimSpace(req.FileName) == "" {
		http.Error(w, "请上传测试大纲文件", http.StatusBadRequest)
		return
	}
	req.BoardModels = strings.Join(splitBoardModels(req.BoardModels), ", ")
	if req.BoardModels == "" {
		http.Error(w, "请填写板卡型号", http.StatusBadRequest)
		return
	}
	if req.ProjectID == 0 {
		http.Error(w, "请选择绑定项目", http.StatusBadRequest)
		return
	}

	var projectName string
	err := config.DB.QueryRow(`
		SELECT IFNULL(project_name, '')
		FROM projects
		WHERE id = ?
		  AND IFNULL(is_deleted, 0) = 0
		  AND IFNULL(audit_status, '') IN ('approved', '已通过', '通过')
		LIMIT 1
	`, req.ProjectID).Scan(&projectName)
	if err != nil {
		http.Error(w, "绑定项目不存在或未审核通过", http.StatusBadRequest)
		return
	}

	var hardwareVersion, deviceType string
	if req.HardwareID > 0 {
		err := config.DB.QueryRow(`
			SELECT IFNULL(hardware_version, ''), IFNULL(device_type, '')
			FROM hardware_versions
			WHERE id = ?
			LIMIT 1
		`, req.HardwareID).Scan(&hardwareVersion, &deviceType)
		if err != nil {
			http.Error(w, "产品型号不存在", http.StatusBadRequest)
			return
		}
	}

	if err := saveUploadedFile(UploadedFilePayload{
		FileID:          req.FileID,
		FileName:        req.FileName,
		FileContentType: req.FileContentType,
		FileData:        req.FileData,
	}); err != nil {
		http.Error(w, "保存测试大纲文件失败: "+err.Error(), http.StatusBadRequest)
		return
	}

	now := time.Now()
	result, err := config.DB.Exec(`
		INSERT INTO production_test_outlines (
			project_id,
			hardware_id,
			hardware_version,
			board_models,
			device_type,
			file_id,
			file_name,
			uploader_id,
			uploader_name,
			upload_time,
			remark,
			created_at,
			updated_at,
			is_deleted
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, 0)
	`, req.ProjectID, req.HardwareID, hardwareVersion, req.BoardModels, deviceType, req.FileID, req.FileName, req.UploaderID, req.UploaderName, now, req.Remark, now, now)
	if err != nil {
		http.Error(w, "新增生产测试大纲失败: "+err.Error(), http.StatusInternalServerError)
		return
	}

	id, _ := result.LastInsertId()
	json.NewEncoder(w).Encode(map[string]interface{}{
		"code": 200,
		"msg":  "保存成功",
		"data": map[string]interface{}{
			"id": id,
		},
	})
}

func DeleteProductionTestOutlineHandler(w http.ResponseWriter, r *http.Request, id int64) {
	if !hasRequestRole(r, "system_admin") && !hasRequestPermission(r, "production:outline:delete") {
		http.Error(w, "无删除生产测试大纲权限：只有管理员可以删除", http.StatusForbidden)
		return
	}

	ensureProductionTestOutlinesTable()

	result, err := config.DB.Exec(`
		UPDATE production_test_outlines
		SET is_deleted = 1,
			updated_at = NOW()
		WHERE id = ?
		  AND is_deleted = 0
	`, id)
	if err != nil {
		http.Error(w, "删除生产测试大纲失败: "+err.Error(), http.StatusInternalServerError)
		return
	}

	affected, _ := result.RowsAffected()
	if affected == 0 {
		http.Error(w, "生产测试大纲不存在或已删除", http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"code": 200,
		"msg":  "删除成功",
	})
}
