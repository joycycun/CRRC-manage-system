package handler

import (
	"crrc_pm_backend/config"
	"crrc_pm_backend/model"
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
// PUT  /api/hardware-versions/{id}
// POST /api/hardware-versions/{id}/upload-zip
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

	// POST /api/hardware-versions/1/upload-zip
	if len(parts) == 2 && r.Method == http.MethodPost && parts[1] == "upload-zip" {
		UploadHardwareZipHandler(w, r, id)
		return
	}

	http.Error(w, "接口不存在", http.StatusNotFound)
}

// ==========================
// GET /api/hardware-versions
// ==========================

func GetHardwareVersionsHandler(w http.ResponseWriter, r *http.Request) {
	rows, err := config.DB.Query(`
		SELECT
			id,
			hardware_version,
			IFNULL(project_id, 0),
			IFNULL(device_type, ''),
			IFNULL(status, ''),
			IFNULL(owner_id, 0),
			IFNULL(owner_name, ''),
			IFNULL(zip_file_id, 0),
			IFNULL(description, ''),
			created_at,
			updated_at
		FROM hardware_versions
		ORDER BY id DESC
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
			&item.ZipFileID,
			&item.Description,
			&item.CreatedAt,
			&item.UpdatedAt,
		)
		if err != nil {
			http.Error(w, "数据解析失败: "+err.Error(), http.StatusInternalServerError)
			return
		}

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
	var item model.HardwareVersion

	err := json.NewDecoder(r.Body).Decode(&item)
	if err != nil {
		http.Error(w, "参数解析失败: "+err.Error(), http.StatusBadRequest)
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

	if item.Status == "" {
		item.Status = "样品"
	}

	now := time.Now()

	result, err := config.DB.Exec(`
		INSERT INTO hardware_versions (
			hardware_version,
			project_id,
			device_type,
			status,
			owner_id,
			owner_name,
			uploader_id,
			uploader_name,
			zip_file_id,
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
		item.ZipFileID,
		item.Description,
		now,
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

// ==========================
// PUT /api/hardware-versions/{id}
// ==========================

func UpdateHardwareVersionHandler(w http.ResponseWriter, r *http.Request, id int64) {
	var item model.HardwareVersion

	err := json.NewDecoder(r.Body).Decode(&item)
	if err != nil {
		http.Error(w, "参数解析失败: "+err.Error(), http.StatusBadRequest)
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

	result, err := config.DB.Exec(`
		UPDATE hardware_versions
		SET
			hardware_version = ?,
			project_id = ?,
			device_type = ?,
			status = ?,
			owner_id = ?,
			owner_name = ?,
			zip_file_id = ?,
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
		item.ZipFileID,
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

	json.NewEncoder(w).Encode(map[string]interface{}{
		"code": 200,
		"msg":  "修改成功",
	})
}

// ==========================
// POST /api/hardware-versions/{id}/upload-zip
// 当前先接收 JSON 里的 zipFileId
// 后面统一文件上传完成后，再改成 multipart 上传
// ==========================

type HardwareZipRequest struct {
	ZipFileID int64 `json:"zipFileId"`
}

func UploadHardwareZipHandler(w http.ResponseWriter, r *http.Request, id int64) {
	var req HardwareZipRequest

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "参数解析失败: "+err.Error(), http.StatusBadRequest)
		return
	}

	if req.ZipFileID == 0 {
		http.Error(w, "zipFileId 不能为空", http.StatusBadRequest)
		return
	}

	result, err := config.DB.Exec(`
		UPDATE hardware_versions
		SET
			zip_file_id = ?,
			updated_at = NOW()
		WHERE id = ?
	`, req.ZipFileID, id)

	if err != nil {
		http.Error(w, "上传硬件资料失败: "+err.Error(), http.StatusInternalServerError)
		return
	}

	affected, _ := result.RowsAffected()
	if affected == 0 {
		http.Error(w, "硬件版本不存在", http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"code": 200,
		"msg":  "硬件资料绑定成功",
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
		IFNULL(ht.uploader_id, 0),
		IFNULL(ht.uploader_name, ''),
		IFNULL(ht.audit_status, ''),
		IFNULL(ht.auditor_id, 0),
		IFNULL(ht.auditor_name, ''),
		ht.audit_time,
		IFNULL(ht.reject_reason, ''),
		IFNULL(ht.remark, ''),
		ht.created_at,
		IFNULL(ht.is_deleted, 0)
	FROM hardware_tests ht
	LEFT JOIN projects p ON ht.project_id = p.id
	LEFT JOIN hardware_versions hv ON ht.hardware_id = hv.id
	WHERE ht.is_deleted = 0
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
			&item.UploaderID,
			&item.UploaderName,
			&item.AuditStatus,
			&item.AuditorID,
			&item.AuditorName,
			&item.AuditTime,
			&item.RejectReason,
			&item.Remark,
			&item.CreatedAt,
			&item.IsDeleted,
		)
		if err != nil {
			http.Error(w, "数据解析失败: "+err.Error(), http.StatusInternalServerError)
			return
		}

		list = append(list, item)
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"code": 200,
		"msg":  "查询成功",
		"data": list,
	})
}

// ============================================================
// POST /api/hardware-tests
// 当前先接收 JSON 里的 fileId
// 后面统一文件上传完成后，再改成 multipart 上传
// ============================================================

func CreateHardwareTestHandler(w http.ResponseWriter, r *http.Request) {
	var item model.HardwareTest

	err := json.NewDecoder(r.Body).Decode(&item)
	if err != nil {
		http.Error(w, "参数解析失败: "+err.Error(), http.StatusBadRequest)
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

	if item.RecordName == "" {
		item.RecordName = "硬件测试记录"
	}

	if item.AuditStatus == "" {
		item.AuditStatus = "待审核"
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
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, 0)
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
	result, err := config.DB.Exec(`
		UPDATE hardware_tests
		SET
			is_deleted = 1,
			updated_at = NOW()
		WHERE id = ?
	`, id)

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

	result, err := config.DB.Exec(`
		UPDATE hardware_tests
		SET
			audit_status = '待审核',
			reject_reason = '',
			updated_at = NOW()
		WHERE id = ? AND is_deleted = 0
	`, id)

	if err != nil {
		http.Error(w, "提交失败: "+err.Error(), http.StatusInternalServerError)
		return
	}

	affected, _ := result.RowsAffected()
	if affected == 0 {
		http.Error(w, "硬件测试记录不存在或已删除", http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"code": 200,
		"msg":  "提交成功",
	})
}
