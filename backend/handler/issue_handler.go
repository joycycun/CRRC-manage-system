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

func IssuesHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	switch r.Method {
	case http.MethodGet:
		GetIssuesHandler(w, r)
	case http.MethodPost:
		CreateIssueHandler(w, r)
	default:
		http.Error(w, "不支持该请求方法", http.StatusMethodNotAllowed)
	}
}

func IssueActionHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	path := strings.TrimPrefix(r.URL.Path, "/api/issues/")
	path = strings.Trim(path, "/")
	parts := strings.Split(path, "/")

	if len(parts) < 1 || parts[0] == "" {
		http.Error(w, "缺少问题ID", http.StatusBadRequest)
		return
	}

	id, err := strconv.ParseInt(parts[0], 10, 64)
	if err != nil {
		http.Error(w, "问题ID错误", http.StatusBadRequest)
		return
	}

	if len(parts) == 1 && r.Method == http.MethodPut {
		UpdateIssueHandler(w, r, id)
		return
	}

	if len(parts) == 1 && r.Method == http.MethodDelete {
		DeleteIssueHandler(w, r, id)
		return
	}

	if len(parts) == 2 && r.Method == http.MethodPost {
		switch parts[1] {
		case "reply":
			ReplyIssueHandler(w, r, id)
			return
		case "close":
			CloseIssueHandler(w, r, id)
			return
		case "reopen":
			ReopenIssueHandler(w, r, id)
			return
		default:
			http.Error(w, "不支持的操作", http.StatusNotFound)
			return
		}
	}

	http.Error(w, "接口不存在", http.StatusNotFound)
}

func GetIssuesHandler(w http.ResponseWriter, r *http.Request) {
	rows, err := config.DB.Query(`
		SELECT
			id,
			project_id,
			IFNULL(device_type, ''),
			IFNULL(issue_source, ''),
			IFNULL(level, ''),
			issue_title,
			IFNULL(issue_desc, ''),
			IFNULL(owner_id, 0),
			IFNULL(owner_name, ''),
			IFNULL(creator_id, 0),
			IFNULL(creator_name, ''),
			create_time,
			plan_close_time,
			real_close_time,
			IFNULL(close_status, ''),
			IFNULL(close_user_id, 0),
			IFNULL(close_user_name, ''),
			IFNULL(reopen_reason, ''),
			updated_at,
			is_deleted
		FROM issues
		WHERE is_deleted = 0
		ORDER BY id DESC
	`)
	if err != nil {
		http.Error(w, "查询失败: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	list := make([]model.Issue, 0)

	for rows.Next() {
		var item model.Issue

		err := rows.Scan(
			&item.ID,
			&item.ProjectID,
			&item.DeviceType,
			&item.IssueSource,
			&item.Level,
			&item.IssueTitle,
			&item.IssueDesc,
			&item.OwnerID,
			&item.OwnerName,
			&item.CreatorID,
			&item.CreatorName,
			&item.CreateTime,
			&item.PlanCloseTime,
			&item.RealCloseTime,
			&item.CloseStatus,
			&item.CloseUserID,
			&item.CloseUserName,
			&item.ReopenReason,
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

func CreateIssueHandler(w http.ResponseWriter, r *http.Request) {
	var item model.Issue
	err := json.NewDecoder(r.Body).Decode(&item)
	if err != nil {
		http.Error(w, "参数解析失败: "+err.Error(), http.StatusBadRequest)
		return
	}

	if item.ProjectID == 0 {
		http.Error(w, "项目ID不能为空", http.StatusBadRequest)
		return
	}
	if item.IssueTitle == "" {
		http.Error(w, "问题标题不能为空", http.StatusBadRequest)
		return
	}

	if item.CloseStatus == "" {
		item.CloseStatus = "打开"
	}
	if item.Level == "" {
		item.Level = "一般"
	}

	now := time.Now()

	result, err := config.DB.Exec(`
		INSERT INTO issues (
			project_id,
			device_type,
			issue_source,
			level,
			issue_title,
			issue_desc,
			owner_id,
			owner_name,
			creator_id,
			creator_name,
			create_time,
			close_status,
			reopen_reason,
			updated_at,
			is_deleted
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, 0)
	`,
		item.ProjectID,
		item.DeviceType,
		item.IssueSource,
		item.Level,
		item.IssueTitle,
		item.IssueDesc,
		item.OwnerID,
		item.OwnerName,
		item.CreatorID,
		item.CreatorName,
		now,
		item.CloseStatus,
		item.ReopenReason,
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

func UpdateIssueHandler(w http.ResponseWriter, r *http.Request, id int64) {
	var item model.Issue
	err := json.NewDecoder(r.Body).Decode(&item)
	if err != nil {
		http.Error(w, "参数解析失败: "+err.Error(), http.StatusBadRequest)
		return
	}

	if item.IssueTitle == "" {
		http.Error(w, "问题标题不能为空", http.StatusBadRequest)
		return
	}

	result, err := config.DB.Exec(`
		UPDATE issues
		SET
			device_type = ?,
			issue_source = ?,
			level = ?,
			issue_title = ?,
			issue_desc = ?,
			owner_id = ?,
			owner_name = ?,
			updated_at = NOW()
		WHERE id = ? AND is_deleted = 0
	`,
		item.DeviceType,
		item.IssueSource,
		item.Level,
		item.IssueTitle,
		item.IssueDesc,
		item.OwnerID,
		item.OwnerName,
		id,
	)

	if err != nil {
		http.Error(w, "修改失败: "+err.Error(), http.StatusInternalServerError)
		return
	}

	affected, _ := result.RowsAffected()
	if affected == 0 {
		http.Error(w, "问题不存在或已删除", http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"code": 200,
		"msg":  "修改成功",
	})
}

type IssueReplyRequest struct {
	ReplyUserID   int64  `json:"replyUserId"`
	ReplyUserName string `json:"replyUserName"`
	Content       string `json:"content"`
}

func ReplyIssueHandler(w http.ResponseWriter, r *http.Request, id int64) {
	var req IssueReplyRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "参数解析失败: "+err.Error(), http.StatusBadRequest)
		return
	}

	if req.Content == "" {
		http.Error(w, "回复内容不能为空", http.StatusBadRequest)
		return
	}

	result, err := config.DB.Exec(`
		INSERT INTO issue_replies (
			issue_id, reply_user_id, reply_user_name, reply_time, content
		) VALUES (?, ?, ?, NOW(), ?)
	`, id, req.ReplyUserID, req.ReplyUserName, req.Content)

	if err != nil {
		http.Error(w, "回复失败: "+err.Error(), http.StatusInternalServerError)
		return
	}

	replyID, _ := result.LastInsertId()

	_, _ = config.DB.Exec(`
		UPDATE issues
		SET close_status = '处理中', updated_at = NOW()
		WHERE id = ? AND is_deleted = 0
	`, id)

	json.NewEncoder(w).Encode(map[string]interface{}{
		"code": 200,
		"msg":  "回复成功",
		"data": map[string]interface{}{"id": replyID},
	})
}

type IssueCloseRequest struct {
	CloseUserID   int64  `json:"closeUserId"`
	CloseUserName string `json:"closeUserName"`
}

func CloseIssueHandler(w http.ResponseWriter, r *http.Request, id int64) {
	var req IssueCloseRequest
	_ = json.NewDecoder(r.Body).Decode(&req)

	result, err := config.DB.Exec(`
		UPDATE issues
		SET
			close_status = '关闭',
			close_user_id = ?,
			close_user_name = ?,
			real_close_time = NOW(),
			updated_at = NOW()
		WHERE id = ? AND is_deleted = 0
	`, req.CloseUserID, req.CloseUserName, id)

	if err != nil {
		http.Error(w, "关闭失败: "+err.Error(), http.StatusInternalServerError)
		return
	}

	affected, _ := result.RowsAffected()
	if affected == 0 {
		http.Error(w, "问题不存在或已删除", http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"code": 200,
		"msg":  "关闭成功",
	})
}

type IssueReopenRequest struct {
	ReopenReason string `json:"reopenReason"`
}

func ReopenIssueHandler(w http.ResponseWriter, r *http.Request, id int64) {
	var req IssueReopenRequest
	_ = json.NewDecoder(r.Body).Decode(&req)

	result, err := config.DB.Exec(`
		UPDATE issues
		SET
			close_status = '重开',
			reopen_reason = ?,
			real_close_time = NULL,
			updated_at = NOW()
		WHERE id = ? AND is_deleted = 0
	`, req.ReopenReason, id)

	if err != nil {
		http.Error(w, "重开失败: "+err.Error(), http.StatusInternalServerError)
		return
	}

	affected, _ := result.RowsAffected()
	if affected == 0 {
		http.Error(w, "问题不存在或已删除", http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"code": 200,
		"msg":  "重开成功",
	})
}

func DeleteIssueHandler(w http.ResponseWriter, r *http.Request, id int64) {
	result, err := config.DB.Exec(`
		UPDATE issues
		SET is_deleted = 1, updated_at = NOW()
		WHERE id = ?
	`, id)

	if err != nil {
		http.Error(w, "删除失败: "+err.Error(), http.StatusInternalServerError)
		return
	}

	affected, _ := result.RowsAffected()
	if affected == 0 {
		http.Error(w, "问题不存在", http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"code": 200,
		"msg":  "删除成功",
	})
}
