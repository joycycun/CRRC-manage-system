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

// 总入口：处理 GET /api/requirement-changes 和 POST /api/requirement-changes
func RequirementChangesHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	switch r.Method {
	case http.MethodGet:
		GetRequirementChangesHandler(w, r)
	case http.MethodPost:
		CreateRequirementChangeHandler(w, r)
	default:
		http.Error(w, "不支持该请求方法", http.StatusMethodNotAllowed)
	}
}

// 带 ID 的入口：处理提交、审核、关闭、删除
func RequirementChangeActionHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	path := strings.TrimPrefix(r.URL.Path, "/api/requirement-changes/")
	path = strings.Trim(path, "/")

	if path == "" {
		http.Error(w, "缺少需求变更ID", http.StatusBadRequest)
		return
	}

	parts := strings.Split(path, "/")

	id, err := strconv.ParseInt(parts[0], 10, 64)
	if err != nil {
		http.Error(w, "需求变更ID错误", http.StatusBadRequest)
		return
	}

	// DELETE /api/requirement-changes/1
	if len(parts) == 1 && r.Method == http.MethodDelete {
		DeleteRequirementChangeHandler(w, r, id)
		return
	}

	// POST /api/requirement-changes/1
	// 默认按提交处理，方便测试
	if len(parts) == 1 && r.Method == http.MethodPost {
		SubmitRequirementChangeHandler(w, r, id)
		return
	}

	if len(parts) == 2 && r.Method == http.MethodPost {
		action := parts[1]

		switch action {
		case "submit":
			SubmitRequirementChangeHandler(w, r, id)
			return
		case "audit":
			AuditRequirementChangeHandler(w, r, id)
			return
		case "close":
			CloseRequirementChangeHandler(w, r, id)
			return
		default:
			http.Error(w, "不支持的操作", http.StatusNotFound)
			return
		}
	}

	http.Error(w, "接口不存在", http.StatusNotFound)
}

// GET /api/requirement-changes
func GetRequirementChangesHandler(w http.ResponseWriter, r *http.Request) {
	rows, err := config.DB.Query(`
		SELECT
			id,
			project_id,
			change_title,
			IFNULL(change_type, ''),
			file_id,
			IFNULL(status, ''),
			IFNULL(close_status, ''),
			IFNULL(submit_user_id, 0),
			IFNULL(submit_user_name, ''),
			IFNULL(audit_user_id, 0),
			IFNULL(audit_user_name, ''),
			submit_time,
			audit_time,
			IFNULL(close_user_id, 0),
			IFNULL(close_user_name, ''),
			close_time,
			IFNULL(reject_reason, ''),
			IFNULL(remark, ''),
			created_at,
			updated_at,
			is_deleted
		FROM requirement_changes
		WHERE is_deleted = 0
		ORDER BY id DESC
	`)
	if err != nil {
		http.Error(w, "查询失败: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	list := make([]model.RequirementChange, 0)

	for rows.Next() {
		var item model.RequirementChange

		err := rows.Scan(
			&item.ID,
			&item.ProjectID,
			&item.ChangeTitle,
			&item.ChangeType,
			&item.FileID,
			&item.Status,
			&item.CloseStatus,
			&item.SubmitUserID,
			&item.SubmitUserName,
			&item.AuditUserID,
			&item.AuditUserName,
			&item.SubmitTime,
			&item.AuditTime,
			&item.CloseUserID,
			&item.CloseUserName,
			&item.CloseTime,
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

// POST /api/requirement-changes
func CreateRequirementChangeHandler(w http.ResponseWriter, r *http.Request) {
	var item model.RequirementChange

	err := json.NewDecoder(r.Body).Decode(&item)
	if err != nil {
		http.Error(w, "参数解析失败: "+err.Error(), http.StatusBadRequest)
		return
	}

	if item.ProjectID == 0 {
		http.Error(w, "项目ID不能为空", http.StatusBadRequest)
		return
	}

	if item.ChangeTitle == "" {
		http.Error(w, "变更标题不能为空", http.StatusBadRequest)
		return
	}

	if item.FileID == 0 {
		http.Error(w, "变更文件ID不能为空", http.StatusBadRequest)
		return
	}

	if item.ChangeType == "" {
		item.ChangeType = "普通变更"
	}

	now := time.Now()

	result, err := config.DB.Exec(`
		INSERT INTO requirement_changes (
			project_id,
			change_title,
			change_type,
			file_id,
			status,
			close_status,
			submit_user_id,
			submit_user_name,
			reject_reason,
			remark,
			created_at,
			updated_at,
			is_deleted
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, 0)
	`,
		item.ProjectID,
		item.ChangeTitle,
		item.ChangeType,
		item.FileID,
		"草稿",
		"未关闭",
		item.SubmitUserID,
		item.SubmitUserName,
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

	item.ID = id
	item.Status = "草稿"
	item.CloseStatus = "未关闭"
	item.CreatedAt = now
	item.UpdatedAt = now
	item.IsDeleted = false

	json.NewEncoder(w).Encode(map[string]interface{}{
		"code": 200,
		"msg":  "新增成功",
		"data": item,
	})
}

// POST /api/requirement-changes/{id}/submit
func SubmitRequirementChangeHandler(w http.ResponseWriter, r *http.Request, id int64) {
	result, err := config.DB.Exec(`
		UPDATE requirement_changes
		SET
			status = '待审核',
			submit_time = NOW(),
			updated_at = NOW()
		WHERE id = ? AND is_deleted = 0
	`, id)

	if err != nil {
		http.Error(w, "提交失败: "+err.Error(), http.StatusInternalServerError)
		return
	}

	affected, _ := result.RowsAffected()
	if affected == 0 {
		http.Error(w, "需求变更不存在或已删除", http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"code": 200,
		"msg":  "提交成功",
	})
}

type RequirementChangeAuditRequest struct {
	AuditUserID   int64  `json:"auditUserId"`
	AuditUserName string `json:"auditUserName"`
	AuditStatus   string `json:"auditStatus"`
	RejectReason  string `json:"rejectReason"`
}

// POST /api/requirement-changes/{id}/audit
func AuditRequirementChangeHandler(w http.ResponseWriter, r *http.Request, id int64) {
	if !requireLeaderPermission(w, r) {
		return
	}

	var req RequirementChangeAuditRequest

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
		UPDATE requirement_changes
		SET
			status = ?,
			audit_user_id = ?,
			audit_user_name = ?,
			audit_time = NOW(),
			reject_reason = ?,
			updated_at = NOW()
		WHERE id = ? AND is_deleted = 0
	`,
		req.AuditStatus,
		req.AuditUserID,
		req.AuditUserName,
		req.RejectReason,
		id,
	)

	if err != nil {
		http.Error(w, "审核失败: "+err.Error(), http.StatusInternalServerError)
		return
	}

	affected, _ := result.RowsAffected()
	if affected == 0 {
		http.Error(w, "需求变更不存在或已删除", http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"code": 200,
		"msg":  "审核成功",
	})
}

type RequirementChangeCloseRequest struct {
	CloseUserID   int64  `json:"closeUserId"`
	CloseUserName string `json:"closeUserName"`
}

// POST /api/requirement-changes/{id}/close
func CloseRequirementChangeHandler(w http.ResponseWriter, r *http.Request, id int64) {
	if !hasRequestRole(r, "project_assistant") {
		http.Error(w, "无关闭权限：需求变更仅项目助理可关闭", http.StatusForbidden)
		return
	}

	var req RequirementChangeCloseRequest

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "参数解析失败: "+err.Error(), http.StatusBadRequest)
		return
	}

	result, err := config.DB.Exec(`
		UPDATE requirement_changes
		SET
			close_status = '已关闭',
			close_user_id = ?,
			close_user_name = ?,
			close_time = NOW(),
			updated_at = NOW()
		WHERE id = ? AND is_deleted = 0
	`,
		req.CloseUserID,
		req.CloseUserName,
		id,
	)

	if err != nil {
		http.Error(w, "关闭失败: "+err.Error(), http.StatusInternalServerError)
		return
	}

	affected, _ := result.RowsAffected()
	if affected == 0 {
		http.Error(w, "需求变更不存在或已删除", http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"code": 200,
		"msg":  "关闭成功",
	})
}

// DELETE /api/requirement-changes/{id}
func DeleteRequirementChangeHandler(w http.ResponseWriter, r *http.Request, id int64) {
	result, err := config.DB.Exec(`
		UPDATE requirement_changes
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
		http.Error(w, "需求变更不存在", http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"code": 200,
		"msg":  "删除成功",
	})
}
