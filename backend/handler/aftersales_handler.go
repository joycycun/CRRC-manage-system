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

// ============================================================
// 维修记录
// GET    /api/repair-records
// POST   /api/repair-records
// PUT    /api/repair-records/{id}
// POST   /api/repair-records/{id}/confirm
// DELETE /api/repair-records/{id}
// ============================================================

func RepairRecordsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	switch r.Method {
	case http.MethodGet:
		GetRepairRecordsHandler(w, r)
	case http.MethodPost:
		CreateRepairRecordHandler(w, r)
	default:
		http.Error(w, "不支持该请求方法", http.StatusMethodNotAllowed)
	}
}

func RepairRecordActionHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	path := strings.TrimPrefix(r.URL.Path, "/api/repair-records/")
	path = strings.Trim(path, "/")
	parts := strings.Split(path, "/")

	if len(parts) < 1 || parts[0] == "" {
		http.Error(w, "缺少维修记录ID", http.StatusBadRequest)
		return
	}

	id, err := strconv.ParseInt(parts[0], 10, 64)
	if err != nil {
		http.Error(w, "维修记录ID错误", http.StatusBadRequest)
		return
	}

	if len(parts) == 1 && r.Method == http.MethodPut {
		UpdateRepairRecordHandler(w, r, id)
		return
	}

	if len(parts) == 1 && r.Method == http.MethodDelete {
		DeleteRepairRecordHandler(w, r, id)
		return
	}

	if len(parts) == 2 && r.Method == http.MethodPost && parts[1] == "confirm" {
		ConfirmRepairRecordHandler(w, r, id)
		return
	}

	http.Error(w, "接口不存在", http.StatusNotFound)
}

type RepairRecordRequest struct {
	ProjectID          int64  `json:"projectId"`
	ProjectName        string `json:"projectName"`
	DeviceType         string `json:"deviceType"`
	InventoryDeviceID  int64  `json:"inventoryDeviceId"`
	SN                 string `json:"sn"`
	MacAddress         string `json:"macAddress"`
	FaultDesc          string `json:"faultDesc"`
	RepairUserID       int64  `json:"repairUserId"`
	RepairUserName     string `json:"repairUserName"`
	RepairUser         string `json:"repairUser"`
	RepairTime         string `json:"repairTime"`
	ReturnCompleteTime string `json:"returnCompleteTime"`
	RepairFinishTime   string `json:"repairFinishTime"`
	RepairMethod       string `json:"repairMethod"`
	RepairProcess      string `json:"repairProcess"`
	ConfirmStatus      string `json:"confirmStatus"`
	Remark             string `json:"remark"`
}

type RepairRecordVO struct {
	ID                 int64  `json:"id"`
	ProjectID          int64  `json:"projectId"`
	ProjectName        string `json:"projectName"`
	DeviceType         string `json:"deviceType"`
	InventoryDeviceID  int64  `json:"inventoryDeviceId"`
	SN                 string `json:"sn"`
	MacAddress         string `json:"macAddress"`
	FaultDesc          string `json:"faultDesc"`
	RepairUserID       int64  `json:"repairUserId"`
	RepairUserName     string `json:"repairUserName"`
	RepairUser         string `json:"repairUser"`
	RepairTime         string `json:"repairTime"`
	ReturnCompleteTime string `json:"returnCompleteTime"`
	RepairFinishTime   string `json:"repairFinishTime"`
	RepairMethod       string `json:"repairMethod"`
	RepairProcess      string `json:"repairProcess"`
	ConfirmStatus      string `json:"confirmStatus"`
	Remark             string `json:"remark"`
	UpdateTime         string `json:"updateTime"`
	CreatedAt          string `json:"createdAt"`
}

func resolveAfterSalesProjectID(projectID int64, projectName string) (int64, error) {
	if projectID > 0 {
		var exists int
		err := config.DB.QueryRow(`
			SELECT COUNT(1)
			FROM projects
			WHERE id = ?
			  AND IFNULL(is_deleted, 0) = 0
			  AND status = '已关闭'
		`, projectID).Scan(&exists)
		if err != nil {
			return 0, err
		}
		if exists == 0 {
			return 0, sql.ErrNoRows
		}
		return projectID, nil
	}

	projectName = strings.TrimSpace(projectName)
	if projectName == "" {
		return 0, sql.ErrNoRows
	}

	var id int64
	err := config.DB.QueryRow(`
		SELECT id
		FROM projects
		WHERE project_name = ?
		  AND IFNULL(is_deleted, 0) = 0
		  AND status = '已关闭'
		LIMIT 1
	`, projectName).Scan(&id)
	if err != nil {
		return 0, err
	}

	return id, nil
}

func parseRepairTime(value string, fallback time.Time) (time.Time, error) {
	value = strings.TrimSpace(value)
	if value == "" {
		return fallback, nil
	}

	for _, layout := range []string{"2006-01-02", "2006-01-02 15:04:05", time.RFC3339} {
		parsed, err := time.ParseInLocation(layout, value, time.Local)
		if err == nil {
			return parsed, nil
		}
	}

	return time.Time{}, strconv.ErrSyntax
}

func parseNullableRepairTime(value string) (interface{}, error) {
	value = strings.TrimSpace(value)
	if value == "" {
		return nil, nil
	}

	parsed, err := parseRepairTime(value, time.Now())
	if err != nil {
		return nil, err
	}

	return parsed, nil
}

func GetRepairRecordsHandler(w http.ResponseWriter, r *http.Request) {
	rows, err := config.DB.Query(`
		SELECT
			rr.id,
			IFNULL(rr.project_id, 0),
			IFNULL(p.project_name, ''),
			IFNULL(rr.device_type, ''),
			IFNULL(rr.inventory_device_id, 0),
			IFNULL(rr.sn, ''),
			IFNULL(rr.mac_address, ''),
			IFNULL(rr.fault_desc, ''),
			IFNULL(rr.repair_user_id, 0),
			IFNULL(rr.repair_user_name, ''),
			IFNULL(DATE_FORMAT(rr.repair_time, '%Y-%m-%d'), ''),
			IFNULL(DATE_FORMAT(rr.repair_finish_time, '%Y-%m-%d'), ''),
			IFNULL(rr.repair_method, ''),
			IFNULL(rr.repair_process, ''),
			IFNULL(rr.confirm_status, ''),
			IFNULL(rr.remark, ''),
			IFNULL(DATE_FORMAT(rr.created_at, '%Y-%m-%d %H:%i:%s'), ''),
			IFNULL(DATE_FORMAT(rr.updated_at, '%Y-%m-%d %H:%i:%s'), '')
		FROM repair_records rr
		INNER JOIN projects p
			ON rr.project_id = p.id
			AND IFNULL(p.is_deleted, 0) = 0
			AND p.status = '已关闭'
		WHERE IFNULL(rr.is_deleted, 0) = 0
		ORDER BY rr.id DESC
	`)
	if err != nil {
		http.Error(w, "查询失败: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	list := make([]RepairRecordVO, 0)

	for rows.Next() {
		var item RepairRecordVO

		err := rows.Scan(
			&item.ID,
			&item.ProjectID,
			&item.ProjectName,
			&item.DeviceType,
			&item.InventoryDeviceID,
			&item.SN,
			&item.MacAddress,
			&item.FaultDesc,
			&item.RepairUserID,
			&item.RepairUserName,
			&item.RepairTime,
			&item.ReturnCompleteTime,
			&item.RepairMethod,
			&item.RepairProcess,
			&item.ConfirmStatus,
			&item.Remark,
			&item.CreatedAt,
			&item.UpdateTime,
		)
		if err != nil {
			http.Error(w, "数据解析失败: "+err.Error(), http.StatusInternalServerError)
			return
		}

		item.RepairUser = item.RepairUserName
		item.RepairFinishTime = item.ReturnCompleteTime
		list = append(list, item)
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"code": 200,
		"msg":  "查询成功",
		"data": list,
	})
}

func CreateRepairRecordHandler(w http.ResponseWriter, r *http.Request) {
	var req RepairRecordRequest

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "参数解析失败: "+err.Error(), http.StatusBadRequest)
		return
	}

	projectID, err := resolveAfterSalesProjectID(req.ProjectID, req.ProjectName)
	if err != nil {
		http.Error(w, "请选择项目立项中已关闭的项目", http.StatusBadRequest)
		return
	}

	req.SN = strings.TrimSpace(req.SN)
	req.MacAddress = strings.TrimSpace(req.MacAddress)
	req.FaultDesc = strings.TrimSpace(req.FaultDesc)

	if req.SN == "" && req.MacAddress == "" {
		http.Error(w, "SN 或 MAC 至少填写一个", http.StatusBadRequest)
		return
	}

	if req.FaultDesc == "" {
		http.Error(w, "故障现象不能为空", http.StatusBadRequest)
		return
	}

	if req.ConfirmStatus == "" {
		req.ConfirmStatus = "待确认"
	}

	if req.RepairMethod == "" {
		req.RepairMethod = "返厂维修"
	}

	repairUserName := strings.TrimSpace(req.RepairUserName)
	if repairUserName == "" {
		repairUserName = strings.TrimSpace(req.RepairUser)
	}
	if repairUserName == "" {
		repairUserName = "当前用户"
	}

	now := time.Now()
	repairTime, err := parseRepairTime(req.RepairTime, now)
	if err != nil {
		http.Error(w, "维修时间格式错误", http.StatusBadRequest)
		return
	}
	finishTime, err := parseNullableRepairTime(req.ReturnCompleteTime)
	if err != nil {
		http.Error(w, "返修完成时间格式错误", http.StatusBadRequest)
		return
	}

	result, err := config.DB.Exec(`
		INSERT INTO repair_records (
			project_id,
			device_type,
			inventory_device_id,
			sn,
			mac_address,
			fault_desc,
			repair_user_id,
			repair_user_name,
			repair_time,
			repair_method,
			repair_process,
			confirm_status,
			remark,
			is_deleted,
			created_at,
			updated_at
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, 0, ?, ?)
	`,
		projectID,
		req.DeviceType,
		req.InventoryDeviceID,
		req.SN,
		req.MacAddress,
		req.FaultDesc,
		req.RepairUserID,
		repairUserName,
		repairTime,
		req.RepairMethod,
		req.RepairProcess,
		req.ConfirmStatus,
		req.Remark,
		now,
		now,
	)

	if err != nil {
		http.Error(w, "新增失败: "+err.Error(), http.StatusInternalServerError)
		return
	}

	id, _ := result.LastInsertId()

	if finishTime != nil {
		_, err = config.DB.Exec(`
			UPDATE repair_records
			SET repair_finish_time = ?, updated_at = NOW()
			WHERE id = ?
		`, finishTime, id)
		if err != nil {
			http.Error(w, "保存返修完成时间失败: "+err.Error(), http.StatusInternalServerError)
			return
		}
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"code": 200,
		"msg":  "新增成功",
		"data": map[string]interface{}{
			"id": id,
		},
	})
}

func UpdateRepairRecordHandler(w http.ResponseWriter, r *http.Request, id int64) {
	var req RepairRecordRequest

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "参数解析失败: "+err.Error(), http.StatusBadRequest)
		return
	}

	projectID, err := resolveAfterSalesProjectID(req.ProjectID, req.ProjectName)
	if err != nil {
		http.Error(w, "请选择项目立项中已关闭的项目", http.StatusBadRequest)
		return
	}

	req.SN = strings.TrimSpace(req.SN)
	req.MacAddress = strings.TrimSpace(req.MacAddress)
	req.FaultDesc = strings.TrimSpace(req.FaultDesc)

	if req.SN == "" && req.MacAddress == "" {
		http.Error(w, "SN 或 MAC 至少填写一个", http.StatusBadRequest)
		return
	}

	if req.FaultDesc == "" {
		http.Error(w, "故障现象不能为空", http.StatusBadRequest)
		return
	}

	if req.ConfirmStatus == "" {
		req.ConfirmStatus = "待确认"
	}

	if req.RepairMethod == "" {
		req.RepairMethod = "返厂维修"
	}

	repairUserName := strings.TrimSpace(req.RepairUserName)
	if repairUserName == "" {
		repairUserName = strings.TrimSpace(req.RepairUser)
	}
	if repairUserName == "" {
		repairUserName = "当前用户"
	}

	repairTime, err := parseRepairTime(req.RepairTime, time.Now())
	if err != nil {
		http.Error(w, "维修时间格式错误", http.StatusBadRequest)
		return
	}
	finishTime, err := parseNullableRepairTime(req.ReturnCompleteTime)
	if err != nil {
		http.Error(w, "返修完成时间格式错误", http.StatusBadRequest)
		return
	}

	result, err := config.DB.Exec(`
		UPDATE repair_records
		SET
			project_id = ?,
			device_type = ?,
			inventory_device_id = ?,
			sn = ?,
			mac_address = ?,
			fault_desc = ?,
			repair_user_id = ?,
			repair_user_name = ?,
			repair_time = ?,
			repair_finish_time = ?,
			repair_method = ?,
			repair_process = ?,
			confirm_status = ?,
			remark = ?,
			updated_at = NOW()
		WHERE id = ? AND is_deleted = 0
	`,
		projectID,
		req.DeviceType,
		req.InventoryDeviceID,
		req.SN,
		req.MacAddress,
		req.FaultDesc,
		req.RepairUserID,
		repairUserName,
		repairTime,
		finishTime,
		req.RepairMethod,
		req.RepairProcess,
		req.ConfirmStatus,
		req.Remark,
		id,
	)

	if err != nil {
		http.Error(w, "修改失败: "+err.Error(), http.StatusInternalServerError)
		return
	}

	affected, _ := result.RowsAffected()
	if affected == 0 {
		http.Error(w, "维修记录不存在或已删除", http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"code": 200,
		"msg":  "修改成功",
	})
}

type RepairConfirmRequest struct {
	Remark string `json:"remark"`
}

func ConfirmRepairRecordHandler(w http.ResponseWriter, r *http.Request, id int64) {
	var req RepairConfirmRequest
	_ = json.NewDecoder(r.Body).Decode(&req)

	result, err := config.DB.Exec(`
		UPDATE repair_records
		SET
			confirm_status = '已完成',
			repair_finish_time = NOW(),
			remark = ?,
			updated_at = NOW()
		WHERE id = ? AND is_deleted = 0
	`, req.Remark, id)

	if err != nil {
		http.Error(w, "确认失败: "+err.Error(), http.StatusInternalServerError)
		return
	}

	affected, _ := result.RowsAffected()
	if affected == 0 {
		http.Error(w, "维修记录不存在或已删除", http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"code": 200,
		"msg":  "确认完成",
	})
}

func DeleteRepairRecordHandler(w http.ResponseWriter, r *http.Request, id int64) {
	result, err := config.DB.Exec(`
		UPDATE repair_records
		SET is_deleted = 1, updated_at = NOW()
		WHERE id = ?
	`, id)

	if err != nil {
		http.Error(w, "删除失败: "+err.Error(), http.StatusInternalServerError)
		return
	}

	affected, _ := result.RowsAffected()
	if affected == 0 {
		http.Error(w, "维修记录不存在", http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"code": 200,
		"msg":  "删除成功",
	})
}

// ============================================================
// 故障分析
// GET    /api/fault-analysis
// POST   /api/fault-analysis
// POST   /api/fault-analysis/{id}/audit
// DELETE /api/fault-analysis/{id}
// ============================================================

func FaultAnalysisHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	switch r.Method {
	case http.MethodGet:
		GetFaultAnalysisHandler(w, r)
	case http.MethodPost:
		CreateFaultAnalysisHandler(w, r)
	default:
		http.Error(w, "不支持该请求方法", http.StatusMethodNotAllowed)
	}
}

func FaultAnalysisActionHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	path := strings.TrimPrefix(r.URL.Path, "/api/fault-analysis/")
	path = strings.Trim(path, "/")
	parts := strings.Split(path, "/")

	if len(parts) < 1 || parts[0] == "" {
		http.Error(w, "缺少故障分析ID", http.StatusBadRequest)
		return
	}

	id, err := strconv.ParseInt(parts[0], 10, 64)
	if err != nil {
		http.Error(w, "故障分析ID错误", http.StatusBadRequest)
		return
	}

	if len(parts) == 1 && r.Method == http.MethodDelete {
		DeleteFaultAnalysisHandler(w, r, id)
		return
	}

	if len(parts) == 2 && r.Method == http.MethodPost && parts[1] == "audit" {
		AuditFaultAnalysisHandler(w, r, id)
		return
	}

	http.Error(w, "接口不存在", http.StatusNotFound)
}

func GetFaultAnalysisHandler(w http.ResponseWriter, r *http.Request) {
	rows, err := config.DB.Query(`
		SELECT
			fa.id,
			IFNULL(fa.project_id, 0),
			IFNULL(fa.issue_id, 0),
			IFNULL(fa.repair_id, 0),
			IFNULL(fa.board_type, ''),
			fa.analysis_name,
			IFNULL(fa.file_id, 0),
			IFNULL(fa.submit_user_id, 0),
			IFNULL(fa.submit_user_name, ''),
			fa.submit_time,
			IFNULL(fa.audit_status, ''),
			IFNULL(fa.auditor_id, 0),
			IFNULL(fa.auditor_name, ''),
			fa.audit_time,
			IFNULL(fa.reject_reason, ''),
			IFNULL(fa.analysis_desc, ''),
			fa.is_deleted,
			fa.created_at,
			fa.updated_at
		FROM fault_analysis fa
		INNER JOIN projects p
			ON fa.project_id = p.id
			AND IFNULL(p.is_deleted, 0) = 0
			AND p.status = '已关闭'
		WHERE fa.is_deleted = 0
		ORDER BY fa.id DESC
	`)
	if err != nil {
		http.Error(w, "查询失败: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	list := make([]model.FaultAnalysis, 0)

	for rows.Next() {
		var item model.FaultAnalysis

		err := rows.Scan(
			&item.ID,
			&item.ProjectID,
			&item.IssueID,
			&item.RepairID,
			&item.BoardType,
			&item.AnalysisName,
			&item.FileID,
			&item.SubmitUserID,
			&item.SubmitUserName,
			&item.SubmitTime,
			&item.AuditStatus,
			&item.AuditorID,
			&item.AuditorName,
			&item.AuditTime,
			&item.RejectReason,
			&item.AnalysisDesc,
			&item.IsDeleted,
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

func CreateFaultAnalysisHandler(w http.ResponseWriter, r *http.Request) {
	var item model.FaultAnalysis

	err := json.NewDecoder(r.Body).Decode(&item)
	if err != nil {
		http.Error(w, "参数解析失败: "+err.Error(), http.StatusBadRequest)
		return
	}

	if item.AnalysisName == "" {
		http.Error(w, "分析方案名称不能为空", http.StatusBadRequest)
		return
	}

	projectID, err := resolveAfterSalesProjectID(item.ProjectID, "")
	if err != nil {
		http.Error(w, "请选择项目立项中已关闭的项目", http.StatusBadRequest)
		return
	}

	if item.AuditStatus == "" {
		item.AuditStatus = "待审核"
	}

	now := time.Now()

	result, err := config.DB.Exec(`
		INSERT INTO fault_analysis (
			project_id,
			issue_id,
			repair_id,
			board_type,
			analysis_name,
			file_id,
			submit_user_id,
			submit_user_name,
			submit_time,
			audit_status,
			reject_reason,
			analysis_desc,
			is_deleted,
			created_at,
			updated_at
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, 0, ?, ?)
	`,
		projectID,
		item.IssueID,
		item.RepairID,
		item.BoardType,
		item.AnalysisName,
		item.FileID,
		item.SubmitUserID,
		item.SubmitUserName,
		now,
		item.AuditStatus,
		item.RejectReason,
		item.AnalysisDesc,
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

type FaultAnalysisAuditRequest struct {
	AuditorID    int64  `json:"auditorId"`
	AuditorName  string `json:"auditorName"`
	AuditStatus  string `json:"auditStatus"`
	RejectReason string `json:"rejectReason"`
}

func AuditFaultAnalysisHandler(w http.ResponseWriter, r *http.Request, id int64) {
	var req FaultAnalysisAuditRequest

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "参数解析失败: "+err.Error(), http.StatusBadRequest)
		return
	}

	if req.AuditStatus != "已通过" && req.AuditStatus != "已驳回" {
		http.Error(w, "审核状态只能是 已通过 或 已驳回", http.StatusBadRequest)
		return
	}

	result, err := config.DB.Exec(`
		UPDATE fault_analysis
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
		http.Error(w, "故障分析不存在或已删除", http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"code": 200,
		"msg":  "审核成功",
	})
}

func DeleteFaultAnalysisHandler(w http.ResponseWriter, r *http.Request, id int64) {
	result, err := config.DB.Exec(`
		UPDATE fault_analysis
		SET is_deleted = 1, updated_at = NOW()
		WHERE id = ?
	`, id)

	if err != nil {
		http.Error(w, "删除失败: "+err.Error(), http.StatusInternalServerError)
		return
	}

	affected, _ := result.RowsAffected()
	if affected == 0 {
		http.Error(w, "故障分析不存在", http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"code": 200,
		"msg":  "删除成功",
	})
}
