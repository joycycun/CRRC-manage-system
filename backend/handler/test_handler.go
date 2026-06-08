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

// ================= 测试用例 =================

func TestCasesHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	switch r.Method {
	case http.MethodGet:
		GetTestCasesHandler(w, r)
	case http.MethodPost:
		CreateTestCaseHandler(w, r)
	default:
		http.Error(w, "不支持该请求方法", http.StatusMethodNotAllowed)
	}
}

func TestCaseActionHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	path := strings.TrimPrefix(r.URL.Path, "/api/test-cases/")
	path = strings.Trim(path, "/")
	parts := strings.Split(path, "/")

	if len(parts) < 1 || parts[0] == "" {
		http.Error(w, "缺少测试用例ID", http.StatusBadRequest)
		return
	}

	id, err := strconv.ParseInt(parts[0], 10, 64)
	if err != nil {
		http.Error(w, "测试用例ID错误", http.StatusBadRequest)
		return
	}

	if len(parts) == 1 && r.Method == http.MethodDelete {
		DeleteTestCaseHandler(w, r, id)
		return
	}

	if len(parts) == 2 && r.Method == http.MethodPost && parts[1] == "audit" {
		AuditTestCaseHandler(w, r, id)
		return
	}

	http.Error(w, "接口不存在", http.StatusNotFound)
}

func GetTestCasesHandler(w http.ResponseWriter, r *http.Request) {
	rows, err := config.DB.Query(`
		SELECT
			id, project_id, case_name, case_file_id,
			IFNULL(report_file_id, 0),
			IFNULL(uploader_id, 0),
			IFNULL(uploader_name, ''),
			upload_time,
			IFNULL(report_uploader_id, 0),
			IFNULL(report_uploader_name, ''),
			report_upload_time,
			IFNULL(audit_status, ''),
			IFNULL(auditor_id, 0),
			IFNULL(auditor_name, ''),
			audit_time,
			IFNULL(reject_reason, ''),
			IFNULL(remark, ''),
			created_at,
			updated_at,
			is_deleted
		FROM test_cases
		WHERE is_deleted = 0
		ORDER BY id DESC
	`)
	if err != nil {
		http.Error(w, "查询失败: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	list := make([]model.TestCase, 0)

	for rows.Next() {
		var item model.TestCase
		err := rows.Scan(
			&item.ID,
			&item.ProjectID,
			&item.CaseName,
			&item.CaseFileID,
			&item.ReportFileID,
			&item.UploaderID,
			&item.UploaderName,
			&item.UploadTime,
			&item.ReportUploaderID,
			&item.ReportUploaderName,
			&item.ReportUploadTime,
			&item.AuditStatus,
			&item.AuditorID,
			&item.AuditorName,
			&item.AuditTime,
			&item.RejectReason,
			&item.Remark,
			&item.CreatedAt,
			&item.UpdatedAt,
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

func CreateTestCaseHandler(w http.ResponseWriter, r *http.Request) {
	var item model.TestCase
	err := json.NewDecoder(r.Body).Decode(&item)
	if err != nil {
		http.Error(w, "参数解析失败: "+err.Error(), http.StatusBadRequest)
		return
	}

	if item.ProjectID == 0 {
		http.Error(w, "项目ID不能为空", http.StatusBadRequest)
		return
	}
	if item.CaseName == "" {
		http.Error(w, "测试用例名称不能为空", http.StatusBadRequest)
		return
	}
	if item.CaseFileID == 0 {
		http.Error(w, "测试用例文件ID不能为空", http.StatusBadRequest)
		return
	}

	now := time.Now()

	result, err := config.DB.Exec(`
		INSERT INTO test_cases (
			project_id, case_name, case_file_id, report_file_id,
			uploader_id, uploader_name, upload_time,
			audit_status, reject_reason, remark,
			created_at, updated_at, is_deleted
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, 0)
	`,
		item.ProjectID,
		item.CaseName,
		item.CaseFileID,
		item.ReportFileID,
		item.UploaderID,
		item.UploaderName,
		now,
		"待审核",
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
		"data": map[string]interface{}{"id": id},
	})
}

type TestAuditRequest struct {
	AuditorID    int64  `json:"auditorId"`
	AuditorName  string `json:"auditorName"`
	AuditStatus  string `json:"auditStatus"`
	RejectReason string `json:"rejectReason"`
}

func AuditTestCaseHandler(w http.ResponseWriter, r *http.Request, id int64) {
	var req TestAuditRequest
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
		UPDATE test_cases
		SET audit_status = ?, auditor_id = ?, auditor_name = ?, audit_time = NOW(),
			reject_reason = ?, updated_at = NOW()
		WHERE id = ? AND is_deleted = 0
	`, req.AuditStatus, req.AuditorID, req.AuditorName, req.RejectReason, id)

	if err != nil {
		http.Error(w, "审核失败: "+err.Error(), http.StatusInternalServerError)
		return
	}

	affected, _ := result.RowsAffected()
	if affected == 0 {
		http.Error(w, "测试用例不存在或已删除", http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"code": 200,
		"msg":  "审核成功",
	})
}

func DeleteTestCaseHandler(w http.ResponseWriter, r *http.Request, id int64) {
	result, err := config.DB.Exec(`
		UPDATE test_cases
		SET is_deleted = 1, updated_at = NOW()
		WHERE id = ?
	`, id)

	if err != nil {
		http.Error(w, "删除失败: "+err.Error(), http.StatusInternalServerError)
		return
	}

	affected, _ := result.RowsAffected()
	if affected == 0 {
		http.Error(w, "测试用例不存在", http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"code": 200,
		"msg":  "删除成功",
	})
}

// ================= 测试报告 =================

func TestReportsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	switch r.Method {
	case http.MethodGet:
		GetTestReportsHandler(w, r)
	case http.MethodPost:
		CreateTestReportHandler(w, r)
	default:
		http.Error(w, "不支持该请求方法", http.StatusMethodNotAllowed)
	}
}

func TestReportActionHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	path := strings.TrimPrefix(r.URL.Path, "/api/test-reports/")
	path = strings.Trim(path, "/")
	parts := strings.Split(path, "/")

	if len(parts) < 1 || parts[0] == "" {
		http.Error(w, "缺少测试报告ID", http.StatusBadRequest)
		return
	}

	id, err := strconv.ParseInt(parts[0], 10, 64)
	if err != nil {
		http.Error(w, "测试报告ID错误", http.StatusBadRequest)
		return
	}

	if len(parts) == 1 && r.Method == http.MethodDelete {
		DeleteTestReportHandler(w, r, id)
		return
	}

	if len(parts) == 2 && r.Method == http.MethodPost && parts[1] == "audit" {
		AuditTestReportHandler(w, r, id)
		return
	}

	http.Error(w, "接口不存在", http.StatusNotFound)
}

func GetTestReportsHandler(w http.ResponseWriter, r *http.Request) {
	rows, err := config.DB.Query(`
		SELECT
			id, project_id,
			IFNULL(test_case_id, 0),
			report_name,
			report_file_id,
			IFNULL(uploader_id, 0),
			IFNULL(uploader_name, ''),
			upload_time,
			IFNULL(audit_status, ''),
			IFNULL(auditor_id, 0),
			IFNULL(auditor_name, ''),
			audit_time,
			IFNULL(reject_reason, ''),
			IFNULL(remark, ''),
			created_at,
			updated_at,
			is_deleted
		FROM test_reports
		WHERE is_deleted = 0
		ORDER BY id DESC
	`)
	if err != nil {
		http.Error(w, "查询失败: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	list := make([]model.TestReport, 0)

	for rows.Next() {
		var item model.TestReport
		err := rows.Scan(
			&item.ID,
			&item.ProjectID,
			&item.TestCaseID,
			&item.ReportName,
			&item.ReportFileID,
			&item.UploaderID,
			&item.UploaderName,
			&item.UploadTime,
			&item.AuditStatus,
			&item.AuditorID,
			&item.AuditorName,
			&item.AuditTime,
			&item.RejectReason,
			&item.Remark,
			&item.CreatedAt,
			&item.UpdatedAt,
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

func CreateTestReportHandler(w http.ResponseWriter, r *http.Request) {
	var item model.TestReport
	err := json.NewDecoder(r.Body).Decode(&item)
	if err != nil {
		http.Error(w, "参数解析失败: "+err.Error(), http.StatusBadRequest)
		return
	}

	if item.ProjectID == 0 {
		http.Error(w, "项目ID不能为空", http.StatusBadRequest)
		return
	}
	if item.ReportName == "" {
		http.Error(w, "测试报告名称不能为空", http.StatusBadRequest)
		return
	}
	if item.ReportFileID == 0 {
		http.Error(w, "测试报告文件ID不能为空", http.StatusBadRequest)
		return
	}

	now := time.Now()

	result, err := config.DB.Exec(`
		INSERT INTO test_reports (
			project_id, test_case_id, report_name, report_file_id,
			uploader_id, uploader_name, upload_time,
			audit_status, reject_reason, remark,
			created_at, updated_at, is_deleted
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, 0)
	`,
		item.ProjectID,
		item.TestCaseID,
		item.ReportName,
		item.ReportFileID,
		item.UploaderID,
		item.UploaderName,
		now,
		"待审核",
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
		"data": map[string]interface{}{"id": id},
	})
}

func AuditTestReportHandler(w http.ResponseWriter, r *http.Request, id int64) {
	var req TestAuditRequest
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
		UPDATE test_reports
		SET audit_status = ?, auditor_id = ?, auditor_name = ?, audit_time = NOW(),
			reject_reason = ?, updated_at = NOW()
		WHERE id = ? AND is_deleted = 0
	`, req.AuditStatus, req.AuditorID, req.AuditorName, req.RejectReason, id)

	if err != nil {
		http.Error(w, "审核失败: "+err.Error(), http.StatusInternalServerError)
		return
	}

	affected, _ := result.RowsAffected()
	if affected == 0 {
		http.Error(w, "测试报告不存在或已删除", http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"code": 200,
		"msg":  "审核成功",
	})
}

func DeleteTestReportHandler(w http.ResponseWriter, r *http.Request, id int64) {
	result, err := config.DB.Exec(`
		UPDATE test_reports
		SET is_deleted = 1, updated_at = NOW()
		WHERE id = ?
	`, id)

	if err != nil {
		http.Error(w, "删除失败: "+err.Error(), http.StatusInternalServerError)
		return
	}

	affected, _ := result.RowsAffected()
	if affected == 0 {
		http.Error(w, "测试报告不存在", http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"code": 200,
		"msg":  "删除成功",
	})
}