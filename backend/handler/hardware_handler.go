package handler

import (
	"crrc_pm_backend/config"
	"crrc_pm_backend/model"
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"time"
)

// ==========================
// 硬件版本入口
// GET  /api/hardware-versions
// POST /api/hardware-versions
// ==========================

func HardwareVersionsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	switch r.Method {
	case http.MethodGet:
		GetHardwareVersionsHandler(w, r)
	case http.MethodPost:
		CreateHardwareVersionHandler(w, r)
	default:
		http.Error(w, "不支持该请求方法", http.StatusMethodNotAllowed)
	}
}

// ==========================
// 硬件版本带 ID 操作入口
// PUT    /api/hardware-versions/{id}
// DELETE /api/hardware-versions/{id}
// POST /api/hardware-versions/{id}/upload-document
// ==========================

func HardwareVersionActionHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	path := strings.TrimPrefix(r.URL.Path, "/api/hardware-versions/")
	path = strings.Trim(path, "/")

	if path == "" {
		http.Error(w, "缺少硬件版本ID", http.StatusBadRequest)
		return
	}

	parts := strings.Split(path, "/")

	id, err := strconv.ParseInt(parts[0], 10, 64)
	if err != nil {
		http.Error(w, "硬件版本ID错误", http.StatusBadRequest)
		return
	}

	// PUT /api/hardware-versions/1
	if len(parts) == 1 && r.Method == http.MethodPut {
		UpdateHardwareVersionHandler(w, r, id)
		return
	}

	// DELETE /api/hardware-versions/1
	if len(parts) == 1 && r.Method == http.MethodDelete {
		DeleteHardwareVersionHandler(w, r, id)
		return
	}

	// POST /api/hardware-versions/1/upload-document
	if len(parts) == 2 && r.Method == http.MethodPost && parts[1] == "upload-document" {
		UploadHardwareDocumentHandler(w, r, id)
		return
	}

	http.Error(w, "接口不存在", http.StatusNotFound)
}

func DeleteHardwareVersionHandler(w http.ResponseWriter, r *http.Request, id int64) {
	if !hasRequestRole(r, "system_admin") && !hasRequestPermission(r, "hardware:delete") {
		http.Error(w, "只有系统管理员可以删除硬件版本", http.StatusForbidden)
		return
	}
	ensureHardwareVersionDocumentColumn()

	tx, err := config.DB.Begin()
	if err != nil {
		http.Error(w, "开启事务失败: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer tx.Rollback()
	if _, err := tx.Exec("DELETE FROM hardware_version_projects WHERE hardware_version_id = ?", id); err != nil {
		http.Error(w, "删除硬件版本项目关联失败: "+err.Error(), http.StatusInternalServerError)
		return
	}
	result, err := tx.Exec(`
		DELETE FROM hardware_versions
		WHERE id = ?
	`, id)
	if err != nil {
		http.Error(w, "删除硬件版本失败: "+err.Error(), http.StatusInternalServerError)
		return
	}

	affected, _ := result.RowsAffected()
	if affected == 0 {
		http.Error(w, "硬件版本不存在", http.StatusNotFound)
		return
	}
	if err := tx.Commit(); err != nil {
		http.Error(w, "提交删除失败: "+err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"code": 200,
		"msg":  "删除成功",
	})
}

func ensureHardwareVersionDocumentColumn() {
	_, _ = config.DB.Exec(`
		ALTER TABLE hardware_versions
		ADD COLUMN change_doc_file_id BIGINT DEFAULT NULL AFTER owner_name
	`)
	_, _ = config.DB.Exec(`
		CREATE TABLE IF NOT EXISTS hardware_version_projects (
			hardware_version_id BIGINT NOT NULL,
			project_id BIGINT NOT NULL,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			PRIMARY KEY (hardware_version_id, project_id),
			KEY idx_hvp_project_id (project_id)
		) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4
	`)
	_, _ = config.DB.Exec(`
		INSERT IGNORE INTO hardware_version_projects (hardware_version_id, project_id, created_at)
		SELECT id, project_id, NOW()
		FROM hardware_versions
		WHERE IFNULL(project_id, 0) > 0
	`)
}

func normalizeHardwareProjectIDs(item model.HardwareVersion) []int64 {
	seen := make(map[int64]bool)
	result := make([]int64, 0, len(item.ProjectIDs)+1)
	for _, projectID := range item.ProjectIDs {
		if projectID > 0 && !seen[projectID] {
			seen[projectID] = true
			result = append(result, projectID)
		}
	}
	if item.ProjectID > 0 && !seen[item.ProjectID] {
		result = append(result, item.ProjectID)
	}
	return result
}

func syncHardwareVersionProjects(tx interface {
	Exec(query string, args ...interface{}) (sql.Result, error)
}, hardwareVersionID int64, projectIDs []int64) error {
	if _, err := tx.Exec("DELETE FROM hardware_version_projects WHERE hardware_version_id = ?", hardwareVersionID); err != nil {
		return err
	}
	for _, projectID := range projectIDs {
		if _, err := tx.Exec(`
			INSERT INTO hardware_version_projects (hardware_version_id, project_id, created_at)
			VALUES (?, ?, NOW())
		`, hardwareVersionID, projectID); err != nil {
			return err
		}
	}
	return nil
}

func isAllowedHardwareDeviceType(deviceType string) bool {
	deviceType = strings.TrimSpace(deviceType)
	allowed := map[string]bool{
		"控制盒（主）":  true,
		"控制盒（辅）":  true,
		"报警器":     true,
		"解码板":     true,
		"编码板":     true,
		"解编码板":    true,
		"司机提醒单元":  true,
		"紧急联络电话":  true,
		"功放板":     true,
		"噪声检测":    true,
		"SIP广播终端": true,
		"SIP对讲终端": true,
		"SIP网关":   true,
	}
	return allowed[deviceType]
}

// ==========================
// GET /api/hardware-versions
// ==========================

func GetHardwareVersionsHandler(w http.ResponseWriter, r *http.Request) {
	ensureUploadedFilesTable()
	ensureHardwareVersionDocumentColumn()

	rows, err := config.DB.Query(`
		SELECT
			hv.id,
			hv.hardware_version,
			IFNULL(hv.project_id, 0),
			IFNULL(hv.device_type, ''),
			IFNULL(hv.status, ''),
			IFNULL(hv.owner_id, 0),
			IFNULL(hv.owner_name, ''),
			IFNULL(hv.change_doc_file_id, 0),
			IFNULL(uf.file_name, ''),
			IFNULL(hv.description, ''),
			hv.created_at,
			hv.updated_at
		FROM hardware_versions hv
		LEFT JOIN uploaded_files uf ON uf.id = hv.change_doc_file_id
		ORDER BY hv.id DESC
	`)
	if err != nil {
		http.Error(w, "查询失败: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	list := make([]model.HardwareVersion, 0)

	for rows.Next() {
		var item model.HardwareVersion

		err := rows.Scan(
			&item.ID,
			&item.HardwareVersion,
			&item.ProjectID,
			&item.DeviceType,
			&item.Status,
			&item.OwnerID,
			&item.OwnerName,
			&item.ChangeDocFileID,
			&item.ChangeDocFileName,
			&item.Description,
			&item.CreatedAt,
			&item.UpdatedAt,
		)
		if err != nil {
			http.Error(w, "数据解析失败: "+err.Error(), http.StatusInternalServerError)
			return
		}

		item.ChangeDocFileURL = filePreviewURL(item.ChangeDocFileID)
		item.ChangeDocDownloadURL = fileDownloadURL(item.ChangeDocFileID)
		projectRows, projectErr := config.DB.Query(`
			SELECT hvp.project_id, IFNULL(p.project_name, '')
			FROM hardware_version_projects hvp
			LEFT JOIN projects p ON p.id = hvp.project_id
			WHERE hvp.hardware_version_id = ?
			ORDER BY hvp.created_at, hvp.project_id
		`, item.ID)
		if projectErr != nil {
			http.Error(w, "查询硬件版本绑定项目失败: "+projectErr.Error(), http.StatusInternalServerError)
			return
		}
		for projectRows.Next() {
			var projectID int64
			var projectName string
			if err := projectRows.Scan(&projectID, &projectName); err != nil {
				projectRows.Close()
				http.Error(w, "解析硬件版本绑定项目失败: "+err.Error(), http.StatusInternalServerError)
				return
			}
			item.ProjectIDs = append(item.ProjectIDs, projectID)
			if projectName != "" {
				item.ProjectNames = append(item.ProjectNames, projectName)
			}
		}
		projectRows.Close()
		list = append(list, item)
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"code": 200,
		"msg":  "查询成功",
		"data": list,
	})
}

// ==========================
// POST /api/hardware-versions
// ==========================

func CreateHardwareVersionHandler(w http.ResponseWriter, r *http.Request) {
	ensureHardwareVersionDocumentColumn()

	var req struct {
		model.HardwareVersion
		FileContentType string `json:"fileContentType"`
		FileData        string `json:"fileData"`
	}

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "参数解析失败: "+err.Error(), http.StatusBadRequest)
		return
	}

	item := req.HardwareVersion
	if err := saveUploadedFile(UploadedFilePayload{
		FileID:          item.ChangeDocFileID,
		FileName:        item.ChangeDocFileName,
		FileContentType: req.FileContentType,
		FileData:        req.FileData,
	}); err != nil {
		http.Error(w, "保存硬件更改文档失败: "+err.Error(), http.StatusBadRequest)
		return
	}

	if item.HardwareVersion == "" {
		http.Error(w, "硬件版本号不能为空", http.StatusBadRequest)
		return
	}

	if item.DeviceType == "" {
		http.Error(w, "终端类型不能为空", http.StatusBadRequest)
		return
	}
	if !isAllowedHardwareDeviceType(item.DeviceType) {
		http.Error(w, "终端类型不在允许范围内", http.StatusBadRequest)
		return
	}

	if item.Status == "" {
		item.Status = "样品"
	}

	now := time.Now()

	projectIDs := normalizeHardwareProjectIDs(item)
	if len(projectIDs) == 0 {
		http.Error(w, "至少需要绑定一个项目", http.StatusBadRequest)
		return
	}
	item.ProjectID = projectIDs[0]

	tx, err := config.DB.Begin()
	if err != nil {
		http.Error(w, "开启事务失败: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer tx.Rollback()

	result, err := tx.Exec(`
		INSERT INTO hardware_versions (
			hardware_version,
			project_id,
			device_type,
			status,
			owner_id,
			owner_name,
			change_doc_file_id,
			description,
			created_at,
			updated_at
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`,
		item.HardwareVersion,
		item.ProjectID,
		item.DeviceType,
		item.Status,
		item.OwnerID,
		item.OwnerName,
		item.ChangeDocFileID,
		item.Description,
		now,
		now,
	)

	if err != nil {
		http.Error(w, "新增失败: "+err.Error(), http.StatusInternalServerError)
		return
	}

	id, _ := result.LastInsertId()
	if err := syncHardwareVersionProjects(tx, id, projectIDs); err != nil {
		http.Error(w, "保存绑定项目失败: "+err.Error(), http.StatusInternalServerError)
		return
	}
	if err := tx.Commit(); err != nil {
		http.Error(w, "提交硬件版本失败: "+err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"code": 200,
		"msg":  "新增成功",
		"data": map[string]interface{}{
			"id": id,
		},
	})
}

// ==========================
// PUT /api/hardware-versions/{id}
// ==========================

func UpdateHardwareVersionHandler(w http.ResponseWriter, r *http.Request, id int64) {
	ensureHardwareVersionDocumentColumn()

	var req struct {
		model.HardwareVersion
		FileContentType string `json:"fileContentType"`
		FileData        string `json:"fileData"`
	}

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "参数解析失败: "+err.Error(), http.StatusBadRequest)
		return
	}

	item := req.HardwareVersion
	if err := saveUploadedFile(UploadedFilePayload{
		FileID:          item.ChangeDocFileID,
		FileName:        item.ChangeDocFileName,
		FileContentType: req.FileContentType,
		FileData:        req.FileData,
	}); err != nil {
		http.Error(w, "保存硬件更改文档失败: "+err.Error(), http.StatusBadRequest)
		return
	}

	if item.HardwareVersion == "" {
		http.Error(w, "硬件版本号不能为空", http.StatusBadRequest)
		return
	}

	if item.DeviceType == "" {
		http.Error(w, "终端类型不能为空", http.StatusBadRequest)
		return
	}
	if !isAllowedHardwareDeviceType(item.DeviceType) {
		http.Error(w, "终端类型不在允许范围内", http.StatusBadRequest)
		return
	}

	projectIDs := normalizeHardwareProjectIDs(item)
	if len(projectIDs) == 0 {
		http.Error(w, "至少需要绑定一个项目", http.StatusBadRequest)
		return
	}
	item.ProjectID = projectIDs[0]

	tx, err := config.DB.Begin()
	if err != nil {
		http.Error(w, "开启事务失败: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer tx.Rollback()

	result, err := tx.Exec(`
		UPDATE hardware_versions
		SET
			hardware_version = ?,
			project_id = ?,
			device_type = ?,
			status = ?,
			owner_id = ?,
			owner_name = ?,
			change_doc_file_id = ?,
			description = ?,
			updated_at = NOW()
		WHERE id = ?
	`,
		item.HardwareVersion,
		item.ProjectID,
		item.DeviceType,
		item.Status,
		item.OwnerID,
		item.OwnerName,
		item.ChangeDocFileID,
		item.Description,
		id,
	)

	if err != nil {
		http.Error(w, "修改失败: "+err.Error(), http.StatusInternalServerError)
		return
	}

	affected, _ := result.RowsAffected()
	if affected == 0 {
		http.Error(w, "硬件版本不存在", http.StatusNotFound)
		return
	}
	if err := syncHardwareVersionProjects(tx, id, projectIDs); err != nil {
		http.Error(w, "保存绑定项目失败: "+err.Error(), http.StatusInternalServerError)
		return
	}
	if err := tx.Commit(); err != nil {
		http.Error(w, "提交硬件版本修改失败: "+err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"code": 200,
		"msg":  "修改成功",
	})
}

// ==========================
// POST /api/hardware-versions/{id}/upload-document
// ==========================

type HardwareDocumentRequest struct {
	ChangeDocFileID   int64  `json:"changeDocFileId"`
	ChangeDocFileName string `json:"changeDocFileName"`
	FileContentType   string `json:"fileContentType"`
	FileData          string `json:"fileData"`
}

func UploadHardwareDocumentHandler(w http.ResponseWriter, r *http.Request, id int64) {
	ensureHardwareVersionDocumentColumn()

	var req HardwareDocumentRequest

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "参数解析失败: "+err.Error(), http.StatusBadRequest)
		return
	}

	if req.ChangeDocFileID == 0 {
		http.Error(w, "文件ID不能为空", http.StatusBadRequest)
		return
	}

	if err := saveUploadedFile(UploadedFilePayload{
		FileID:          req.ChangeDocFileID,
		FileName:        req.ChangeDocFileName,
		FileContentType: req.FileContentType,
		FileData:        req.FileData,
	}); err != nil {
		http.Error(w, "保存硬件更改文档失败: "+err.Error(), http.StatusBadRequest)
		return
	}

	result, err := config.DB.Exec(`
		UPDATE hardware_versions
		SET
			change_doc_file_id = ?,
			updated_at = NOW()
		WHERE id = ?
	`, req.ChangeDocFileID, id)

	if err != nil {
		http.Error(w, "上传硬件更改文档失败: "+err.Error(), http.StatusInternalServerError)
		return
	}

	affected, _ := result.RowsAffected()
	if affected == 0 {
		http.Error(w, "硬件版本不存在", http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"code": 200,
		"msg":  "硬件更改文档绑定成功",
	})
}

// ============================================================
// 硬件测试入口
// GET  /api/hardware-tests
// POST /api/hardware-tests
// ============================================================

func HardwareTestsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	switch r.Method {
	case http.MethodGet:
		GetHardwareTestsHandler(w, r)
	case http.MethodPost:
		CreateHardwareTestHandler(w, r)
	default:
		http.Error(w, "不支持该请求方法", http.StatusMethodNotAllowed)
	}
}

// ============================================================
// 硬件测试带 ID 操作入口
// POST   /api/hardware-tests/{id}/audit
// DELETE /api/hardware-tests/{id}
// ============================================================

func HardwareTestActionHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	path := strings.TrimPrefix(r.URL.Path, "/api/hardware-tests/")
	path = strings.Trim(path, "/")

	if path == "" {
		http.Error(w, "缺少硬件测试ID", http.StatusBadRequest)
		return
	}

	parts := strings.Split(path, "/")

	id, err := strconv.ParseInt(parts[0], 10, 64)
	if err != nil {
		http.Error(w, "硬件测试ID错误", http.StatusBadRequest)
		return
	}

	// DELETE /api/hardware-tests/1
	if len(parts) == 1 && r.Method == http.MethodDelete {
		DeleteHardwareTestHandler(w, r, id)
		return
	}

	// POST /api/hardware-tests/1/audit
	if len(parts) == 2 && r.Method == http.MethodPost && parts[1] == "audit" {
		AuditHardwareTestHandler(w, r, id)
		return
	}
	if len(parts) == 2 && r.Method == http.MethodPost && parts[1] == "submit" {
		SubmitHardwareTestHandler(w, r, id)
		return
	}

	http.Error(w, "接口不存在", http.StatusNotFound)
}

// ============================================================
// GET /api/hardware-tests
// ============================================================

func GetHardwareTestsHandler(w http.ResponseWriter, r *http.Request) {
	ensureUploadedFilesTable()

	visibilitySQL := hardwareTestVisibilitySQL(r)
	rows, err := config.DB.Query(`
	SELECT
		ht.id,
		ht.project_id,
		IFNULL(p.project_name, ''),
		ht.hardware_id,
		IFNULL(hv.hardware_version, ''),
		ht.record_name,
		ht.device_type,
		IFNULL(ht.file_id, 0),
		IFNULL(uf.file_name, ''),
		IFNULL(ht.uploader_id, 0),
		IFNULL(ht.uploader_name, ''),
		IFNULL(ht.audit_status, ''),
		IFNULL(ht.auditor_id, 0),
		IFNULL(ht.auditor_name, ''),
		ht.audit_time,
		IFNULL(ht.reject_reason, ''),
		IFNULL(ht.remark, ''),
		ht.upload_time,
		ht.created_at,
		IFNULL(ht.is_deleted, 0)
	FROM hardware_tests ht
	INNER JOIN projects p
		ON ht.project_id = p.id
		AND IFNULL(p.is_deleted, 0) = 0
		AND IFNULL(p.audit_status, '未提交') = '已通过'
	LEFT JOIN hardware_versions hv ON ht.hardware_id = hv.id
	LEFT JOIN uploaded_files uf ON uf.id = ht.file_id
	WHERE ht.is_deleted = 0 ` + visibilitySQL + `
	ORDER BY ht.id DESC
`)
	if err != nil {
		http.Error(w, "查询失败: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	list := make([]model.HardwareTest, 0)

	for rows.Next() {
		var item model.HardwareTest

		err := rows.Scan(
			&item.ID,
			&item.ProjectID,
			&item.ProjectName,
			&item.HardwareID,
			&item.HardwareVersion,
			&item.RecordName,
			&item.DeviceType,
			&item.FileID,
			&item.FileName,
			&item.UploaderID,
			&item.UploaderName,
			&item.AuditStatus,
			&item.AuditorID,
			&item.AuditorName,
			&item.AuditTime,
			&item.RejectReason,
			&item.Remark,
			&item.UploadTime,
			&item.CreatedAt,
			&item.IsDeleted,
		)
		if err != nil {
			http.Error(w, "数据解析失败: "+err.Error(), http.StatusInternalServerError)
			return
		}

		item.FileURL = filePreviewURL(item.FileID)
		item.DownloadURL = fileDownloadURL(item.FileID)
		list = append(list, item)
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"code": 200,
		"msg":  "查询成功",
		"data": list,
	})
}

func hardwareTestVisibilitySQL(r *http.Request) string {
	userID, userName := currentRequestUser(r)
	ownSQL := ""
	if userID > 0 {
		ownSQL += " OR IFNULL(ht.uploader_id, 0) = " + strconv.FormatInt(userID, 10)
	}
	if userName != "" {
		ownSQL += " OR IFNULL(ht.uploader_name, '') = '" + strings.ReplaceAll(userName, "'", "''") + "'"
	}

	if hasRequestRole(r, "system_admin") || hasRequestRole(r, "leader") || hasRequestPermission(r, "hardware:audit") {
		return " AND (IFNULL(ht.audit_status, '草稿') IN ('待审核', 'submitted', '已提交', '已通过', '审核通过', 'approved', '已驳回', '审核驳回', 'rejected')" + ownSQL + ")"
	}

	return " AND (IFNULL(ht.audit_status, '草稿') IN ('已通过', '审核通过', 'approved')" + ownSQL + ")"
}

// ============================================================
// POST /api/hardware-tests
// 当前先接收 JSON 里的 fileId
// 后面统一文件上传完成后，再改成 multipart 上传
// ============================================================

func CreateHardwareTestHandler(w http.ResponseWriter, r *http.Request) {
	var req struct {
		model.HardwareTest
		FileContentType string `json:"fileContentType"`
		FileData        string `json:"fileData"`
	}

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "参数解析失败: "+err.Error(), http.StatusBadRequest)
		return
	}

	item := req.HardwareTest
	if err := saveUploadedFile(UploadedFilePayload{
		FileID:          item.FileID,
		FileName:        item.FileName,
		FileContentType: req.FileContentType,
		FileData:        req.FileData,
	}); err != nil {
		http.Error(w, "保存硬件测试文件失败: "+err.Error(), http.StatusBadRequest)
		return
	}

	if item.ProjectID == 0 {
		http.Error(w, "项目ID不能为空", http.StatusBadRequest)
		return
	}

	if item.HardwareID == 0 {
		http.Error(w, "硬件版本ID不能为空", http.StatusBadRequest)
		return
	}

	if item.FileID == 0 {
		http.Error(w, "测试文件ID不能为空", http.StatusBadRequest)
		return
	}
	if item.DeviceType != "" && !isAllowedHardwareDeviceType(item.DeviceType) {
		http.Error(w, "终端类型不在允许范围内", http.StatusBadRequest)
		return
	}

	if item.RecordName == "" {
		item.RecordName = item.FileName
	}
	if item.RecordName == "" {
		item.RecordName = "硬件测试记录"
	}

	if item.AuditStatus == "" {
		item.AuditStatus = "草稿"
	}

	now := time.Now()

	result, err := config.DB.Exec(`
		INSERT INTO hardware_tests (
			project_id,
			hardware_id,
			record_name,
			device_type,
			file_id,
			audit_status,
			reject_reason,
			remark,
			created_at,
			updated_at,
			is_deleted,
			uploader_id,
			uploader_name,
			upload_time
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, 0, ?, ?, ?)
	`,
		item.ProjectID,
		item.HardwareID,
		item.RecordName,
		item.DeviceType,
		item.FileID,
		item.AuditStatus,
		item.RejectReason,
		item.Remark,
		now,
		now,
		item.UploaderID,
		item.UploaderName,
		now,
	)

	if err != nil {
		http.Error(w, "新增失败: "+err.Error(), http.StatusInternalServerError)
		return
	}

	id, _ := result.LastInsertId()

	json.NewEncoder(w).Encode(map[string]interface{}{
		"code": 200,
		"msg":  "新增成功",
		"data": map[string]interface{}{
			"id": id,
		},
	})
}

// ============================================================
// POST /api/hardware-tests/{id}/audit
// ============================================================

type HardwareTestAuditRequest struct {
	AuditorID    int64  `json:"auditorId"`
	AuditorName  string `json:"auditorName"`
	AuditStatus  string `json:"auditStatus"`
	RejectReason string `json:"rejectReason"`
}

func AuditHardwareTestHandler(w http.ResponseWriter, r *http.Request, id int64) {
	if !requireLeaderPermission(w, r) {
		return
	}

	var req HardwareTestAuditRequest

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "参数解析失败: "+err.Error(), http.StatusBadRequest)
		return
	}

	if req.AuditStatus == "" {
		http.Error(w, "审核状态不能为空", http.StatusBadRequest)
		return
	}

	if req.AuditStatus != "已通过" && req.AuditStatus != "已驳回" {
		http.Error(w, "审核状态只能是 已通过 或 已驳回", http.StatusBadRequest)
		return
	}
	req.AuditorID, req.AuditorName = normalizeAuditUser(r, req.AuditorID, req.AuditorName)

	result, err := config.DB.Exec(`
		UPDATE hardware_tests
		SET
			audit_status = ?,
			auditor_id = ?,
			auditor_name = ?,
			audit_time = NOW(),
			reject_reason = ?,
			updated_at = NOW()
		WHERE id = ? AND is_deleted = 0
	`,
		req.AuditStatus,
		req.AuditorID,
		req.AuditorName,
		req.RejectReason,
		id,
	)

	if err != nil {
		http.Error(w, "审核失败: "+err.Error(), http.StatusInternalServerError)
		return
	}

	affected, _ := result.RowsAffected()
	if affected == 0 {
		http.Error(w, "硬件测试记录不存在或已删除", http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"code": 200,
		"msg":  "审核成功",
	})
}

// ============================================================
// DELETE /api/hardware-tests/{id}
// ============================================================

func DeleteHardwareTestHandler(w http.ResponseWriter, r *http.Request, id int64) {
	canDelete := hasRequestRole(r, "system_admin") ||
		hasRequestRole(r, "hardware_owner") ||
		hasRequestPermission(r, "hardware:delete")
	if !canDelete {
		http.Error(w, "无删除硬件测试记录权限", http.StatusForbidden)
		return
	}

	statusSQL := ""
	args := []interface{}{id}
	if !hasRequestRole(r, "system_admin") {
		userID, userName := currentRequestUser(r)
		statusSQL = " AND IFNULL(audit_status, '草稿') IN ('草稿', 'draft', '已驳回', '审核驳回', 'rejected')" +
			" AND (IFNULL(uploader_id, 0) = ? OR IFNULL(uploader_name, '') = ?)"
		args = append(args, userID, userName)
	}

	result, err := config.DB.Exec(`
		UPDATE hardware_tests
		SET
			is_deleted = 1,
			updated_at = NOW()
		WHERE id = ?
		  AND is_deleted = 0
	`+statusSQL, args...)

	if err != nil {
		http.Error(w, "删除失败: "+err.Error(), http.StatusInternalServerError)
		return
	}

	affected, _ := result.RowsAffected()
	if affected == 0 {
		http.Error(w, "硬件测试记录不存在", http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"code": 200,
		"msg":  "删除成功",
	})
}
func SubmitHardwareTestHandler(w http.ResponseWriter, r *http.Request, id int64) {
	w.Header().Set("Content-Type", "application/json")

	if !canSubmitHardwareTest(r, id) {
		http.Error(w, "无提交权限：只有上传人可以提交当前硬件测试记录", http.StatusForbidden)
		return
	}

	result, err := config.DB.Exec(`
		UPDATE hardware_tests
		SET
			audit_status = '待审核',
			reject_reason = '',
			updated_at = NOW()
		WHERE id = ? AND is_deleted = 0
		  AND IFNULL(audit_status, '草稿') IN ('草稿', 'draft', '未提交', '已驳回', '审核驳回', 'rejected')
	`, id)

	if err != nil {
		http.Error(w, "提交失败: "+err.Error(), http.StatusInternalServerError)
		return
	}

	affected, _ := result.RowsAffected()
	if affected == 0 {
		http.Error(w, "硬件测试记录不存在、已删除或当前状态不可提交", http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"code": 200,
		"msg":  "提交成功",
	})
}

func canSubmitHardwareTest(r *http.Request, id int64) bool {
	if hasRequestRole(r, "system_admin") {
		return true
	}

	var uploaderID int64
	var uploaderName string
	err := config.DB.QueryRow(`
		SELECT IFNULL(uploader_id, 0), IFNULL(uploader_name, '')
		FROM hardware_tests
		WHERE id = ? AND is_deleted = 0
		LIMIT 1
	`, id).Scan(&uploaderID, &uploaderName)
	if err != nil {
		return false
	}

	userID, userName := currentRequestUser(r)
	if userID > 0 && uploaderID > 0 && userID == uploaderID {
		return true
	}
	return userName != "" && uploaderName != "" && userName == uploaderName
}
