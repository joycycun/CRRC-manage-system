package handler

import (
	"crrc_pm_backend/config"
	"crrc_pm_backend/model"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
)

func TestCasesHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		GetTestCasesHandler(w, r)
	case http.MethodPost:
		CreateTestCaseHandler(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func TestCaseActionHandler(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimPrefix(r.URL.Path, "/api/test-cases/")
	parts := strings.Split(path, "/")

	if len(parts) < 1 || parts[0] == "" {
		http.NotFound(w, r)
		return
	}

	id, err := strconv.ParseInt(parts[0], 10, 64)
	if err != nil {
		http.Error(w, "无效的测试用例ID", http.StatusBadRequest)
		return
	}

	if len(parts) == 1 && r.Method == http.MethodDelete {
		DeleteTestCaseHandler(w, r, id)
		return
	}

	if len(parts) == 2 && r.Method == http.MethodPost && parts[1] == "submit" {
		SubmitTestCaseHandler(w, r, id)
		return
	}

	if len(parts) == 2 && r.Method == http.MethodPost && parts[1] == "audit" {
		AuditTestCaseHandler(w, r, id)
		return
	}

	if len(parts) == 2 && r.Method == http.MethodPost && parts[1] == "report" {
		UploadTestReportHandler(w, r, id)
		return
	}

	http.NotFound(w, r)
}

func GetTestCasesHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	rows, err := config.DB.Query(`
		SELECT
			tc.id,
			IFNULL(tc.project_id, 0),
			IFNULL(p.project_name, ''),
			IFNULL(tc.case_name, ''),
			IFNULL(tc.file_id, 0),
			IFNULL(tc.file_name, ''),
			IFNULL(tc.uploader_id, 0),
			IFNULL(tc.uploader_name, ''),
			IFNULL(DATE_FORMAT(tc.upload_time, '%Y-%m-%d'), ''),
			IFNULL(tc.audit_status, ''),
			IFNULL(tc.auditor_id, 0),
			IFNULL(tc.auditor_name, ''),
			IFNULL(DATE_FORMAT(tc.audit_time, '%Y-%m-%d'), ''),
			IFNULL(tc.reject_reason, ''),
			IFNULL(tc.remark, ''),
			IFNULL(tc.report_name, ''),
			IFNULL(tc.report_file_id, 0),
			IFNULL(tc.report_file_name, ''),
			IFNULL(tc.report_uploader_id, 0),
			IFNULL(tc.report_uploader_name, ''),
			IFNULL(DATE_FORMAT(tc.report_upload_time, '%Y-%m-%d'), ''),
			IFNULL(tc.report_remark, ''),
			IFNULL(DATE_FORMAT(tc.created_at, '%Y-%m-%d %H:%i:%s'), ''),
			IFNULL(DATE_FORMAT(tc.updated_at, '%Y-%m-%d %H:%i:%s'), '')
		FROM test_cases tc
		LEFT JOIN projects p ON tc.project_id = p.id
		WHERE tc.is_deleted = 0
		ORDER BY tc.id DESC
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
			&item.ProjectName,
			&item.CaseName,
			&item.FileID,
			&item.FileName,
			&item.UploaderID,
			&item.UploaderName,
			&item.UploadTime,
			&item.AuditStatus,
			&item.AuditorID,
			&item.AuditorName,
			&item.AuditTime,
			&item.RejectReason,
			&item.Remark,
			&item.ReportName,
			&item.ReportFileID,
			&item.ReportFileName,
			&item.ReportUploaderID,
			&item.ReportUploaderName,
			&item.ReportUploadTime,
			&item.ReportRemark,
			&item.CreatedAt,
			&item.UpdatedAt,
		)
		if err != nil {
			http.Error(w, "数据解析失败: "+err.Error(), http.StatusInternalServerError)
			return
		}

		item.Uploader = item.UploaderName
		item.Auditor = item.AuditorName
		item.ReportUploader = item.ReportUploaderName

		list = append(list, item)
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"code": 200,
		"msg":  "查询成功",
		"data": list,
	})
}

func CreateTestCaseHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var req model.TestCase
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "参数解析失败: "+err.Error(), http.StatusBadRequest)
		return
	}

	if req.ProjectID == 0 {
		http.Error(w, "项目ID不能为空", http.StatusBadRequest)
		return
	}

	if req.CaseName == "" {
		http.Error(w, "测试用例名称不能为空", http.StatusBadRequest)
		return
	}

	if req.UploaderName == "" {
		req.UploaderName = req.Uploader
	}

	if req.AuditStatus == "" {
		req.AuditStatus = "草稿"
	}

	result, err := config.DB.Exec(`
		INSERT INTO test_cases (
			project_id,
			case_name,
			file_id,
			file_name,
			uploader_id,
			uploader_name,
			upload_time,
			audit_status,
			remark,
			is_deleted,
			created_at,
			updated_at
		) VALUES (?, ?, ?, ?, ?, ?, NOW(), ?, ?, 0, NOW(), NOW())
	`,
		req.ProjectID,
		req.CaseName,
		req.FileID,
		req.FileName,
		req.UploaderID,
		req.UploaderName,
		req.AuditStatus,
		req.Remark,
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

func SubmitTestCaseHandler(w http.ResponseWriter, r *http.Request, id int64) {
	w.Header().Set("Content-Type", "application/json")

	result, err := config.DB.Exec(`
		UPDATE test_cases
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
		http.Error(w, "测试用例不存在或已删除", http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"code": 200,
		"msg":  "提交成功",
	})
}

func AuditTestCaseHandler(w http.ResponseWriter, r *http.Request, id int64) {
	w.Header().Set("Content-Type", "application/json")

	if !requireLeaderPermission(w, r) {
		return
	}

	var req model.TestCase
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "参数解析失败: "+err.Error(), http.StatusBadRequest)
		return
	}

	if req.AuditorName == "" {
		req.AuditorName = req.Auditor
	}

	if req.AuditStatus == "" {
		http.Error(w, "审核状态不能为空", http.StatusBadRequest)
		return
	}

	result, err := config.DB.Exec(`
		UPDATE test_cases
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
		http.Error(w, "测试用例不存在或已删除", http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"code": 200,
		"msg":  "审核成功",
	})
}

func UploadTestReportHandler(w http.ResponseWriter, r *http.Request, id int64) {
	w.Header().Set("Content-Type", "application/json")

	var req model.TestCase
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "参数解析失败: "+err.Error(), http.StatusBadRequest)
		return
	}

	if req.ReportName == "" {
		http.Error(w, "测试报告名称不能为空", http.StatusBadRequest)
		return
	}

	if req.ReportFileName == "" {
		req.ReportFileName = req.FileName
	}

	if req.ReportUploaderName == "" {
		req.ReportUploaderName = req.UploaderName
	}

	result, err := config.DB.Exec(`
		UPDATE test_cases
		SET
			report_name = ?,
			report_file_id = ?,
			report_file_name = ?,
			report_uploader_id = ?,
			report_uploader_name = ?,
			report_upload_time = NOW(),
			report_remark = ?,
			updated_at = NOW()
		WHERE id = ? AND is_deleted = 0
	`,
		req.ReportName,
		req.ReportFileID,
		req.ReportFileName,
		req.ReportUploaderID,
		req.ReportUploaderName,
		req.ReportRemark,
		id,
	)
	if err != nil {
		http.Error(w, "上传测试报告失败: "+err.Error(), http.StatusInternalServerError)
		return
	}

	affected, _ := result.RowsAffected()
	if affected == 0 {
		http.Error(w, "测试用例不存在或已删除", http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"code": 200,
		"msg":  "上传测试报告成功",
	})
}

func DeleteTestCaseHandler(w http.ResponseWriter, r *http.Request, id int64) {
	w.Header().Set("Content-Type", "application/json")

	result, err := config.DB.Exec(`
		UPDATE test_cases
		SET is_deleted = 1, updated_at = NOW()
		WHERE id = ? AND is_deleted = 0
	`, id)
	if err != nil {
		http.Error(w, "删除失败: "+err.Error(), http.StatusInternalServerError)
		return
	}

	affected, _ := result.RowsAffected()
	if affected == 0 {
		http.Error(w, "测试用例不存在或已删除", http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"code": 200,
		"msg":  "删除成功",
	})
}
