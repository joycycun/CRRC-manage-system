package handler

import (
	"crrc_pm_backend/config"
	"crrc_pm_backend/model"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
)

func IssuesHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		GetIssuesHandler(w, r)
	case http.MethodPost:
		CreateIssueHandler(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func IssueActionHandler(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimPrefix(r.URL.Path, "/api/issues/")
	parts := strings.Split(path, "/")

	if len(parts) < 1 || parts[0] == "" {
		http.NotFound(w, r)
		return
	}

	id, err := strconv.ParseInt(parts[0], 10, 64)
	if err != nil {
		http.Error(w, "无效的问题ID", http.StatusBadRequest)
		return
	}

	if len(parts) == 1 && r.Method == http.MethodPut {
		UpdateIssueHandler(w, r, id)
		return
	}

	if len(parts) == 2 && r.Method == http.MethodPost && parts[1] == "reply" {
		ReplyIssueHandler(w, r, id)
		return
	}

	if len(parts) == 2 && r.Method == http.MethodPost && parts[1] == "close" {
		CloseIssueHandler(w, r, id)
		return
	}

	if len(parts) == 2 && r.Method == http.MethodPost && parts[1] == "reopen" {
		ReopenIssueHandler(w, r, id)
		return
	}

	if len(parts) == 1 && r.Method == http.MethodDelete {
		DeleteIssueHandler(w, r, id)
		return
	}

	http.NotFound(w, r)
}

func GetIssuesHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	rows, err := config.DB.Query(`
		SELECT
			i.id,
			i.project_id,
			IFNULL(p.project_name, ''),
			IFNULL(i.device_type, ''),
			IFNULL(i.issue_source, ''),
			IFNULL(i.level, ''),
			IFNULL(i.issue_title, ''),
			IFNULL(i.issue_desc, ''),
			IFNULL(i.owner_id, 0),
			IFNULL(i.owner_name, ''),
			IFNULL(i.creator_id, 0),
			IFNULL(i.creator_name, ''),
			IFNULL(DATE_FORMAT(i.create_time, '%Y-%m-%d'), ''),
			IFNULL(DATE_FORMAT(i.plan_close_time, '%Y-%m-%d'), ''),
			IFNULL(DATE_FORMAT(i.real_close_time, '%Y-%m-%d'), ''),
			CASE
					WHEN i.close_status IN ('已关闭', '关闭', 'closed') THEN '已关闭'
					ELSE '打开'
			END,
			IFNULL(i.close_user_id, 0),
			IFNULL(i.close_user_name, ''),
			IFNULL(i.reopen_reason, ''),
			IFNULL(DATE_FORMAT(i.updated_at, '%Y-%m-%d %H:%i:%s'), '')
		FROM issues i
		LEFT JOIN projects p ON i.project_id = p.id
		WHERE IFNULL(i.is_deleted, 0) = 0
		ORDER BY i.id DESC
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
			&item.ProjectName,
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
		)
		if err != nil {
			http.Error(w, "数据解析失败: "+err.Error(), http.StatusInternalServerError)
			return
		}

		item.Source = item.IssueSource
		item.Severity = item.Level
		item.Title = item.IssueTitle
		item.Owner = item.OwnerName
		item.Creator = item.CreatorName
		item.CloseUser = item.CloseUserName

		replies, err := getIssueReplies(item.ID)
		if err != nil {
			http.Error(w, "查询回复失败: "+err.Error(), http.StatusInternalServerError)
			return
		}

		item.Replies = replies
		list = append(list, item)
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"code": 200,
		"msg":  "查询成功",
		"data": list,
	})
}

func getIssueReplies(issueID int64) ([]model.IssueReply, error) {
	rows, err := config.DB.Query(`
		SELECT
			id,
			issue_id,
			IFNULL(reply_user_id, 0),
			IFNULL(reply_user_name, ''),
			IFNULL(DATE_FORMAT(reply_time, '%Y-%m-%d'), ''),
			IFNULL(content, '')
		FROM issue_replies
		WHERE issue_id = ?
		ORDER BY id ASC
	`, issueID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	list := make([]model.IssueReply, 0)

	for rows.Next() {
		var item model.IssueReply

		err := rows.Scan(
			&item.ID,
			&item.IssueID,
			&item.ReplyUserID,
			&item.ReplyUserName,
			&item.ReplyTime,
			&item.Content,
		)
		if err != nil {
			return nil, err
		}

		item.ReplyUser = item.ReplyUserName
		item.ReplyContent = item.Content
		list = append(list, item)
	}

	return list, nil
}

func CreateIssueHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var req model.Issue
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "参数解析失败: "+err.Error(), http.StatusBadRequest)
		return
	}

	if req.ProjectID == 0 {
		http.Error(w, "项目ID不能为空", http.StatusBadRequest)
		return
	}

	if req.IssueTitle == "" {
		req.IssueTitle = req.Title
	}

	if req.IssueTitle == "" {
		http.Error(w, "问题名称不能为空", http.StatusBadRequest)
		return
	}

	if req.IssueSource == "" {
		req.IssueSource = req.Source
	}

	if req.Level == "" {
		req.Level = req.Severity
	}

	if req.OwnerName == "" {
		req.OwnerName = req.Owner
	}

	if req.CreatorName == "" {
		req.CreatorName = req.Creator
	}

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
			plan_close_time,
			close_status,
			is_deleted
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, NOW(), NULLIF(?, ''), '未关闭', 0)
	`,
		req.ProjectID,
		req.DeviceType,
		req.IssueSource,
		req.Level,
		req.IssueTitle,
		req.IssueDesc,
		req.OwnerID,
		req.OwnerName,
		req.CreatorID,
		req.CreatorName,
		req.PlanCloseTime,
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

func UpdateIssueHandler(w http.ResponseWriter, r *http.Request, id int64) {
	w.Header().Set("Content-Type", "application/json")

	var req model.Issue
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "参数解析失败: "+err.Error(), http.StatusBadRequest)
		return
	}

	if req.IssueTitle == "" {
		req.IssueTitle = req.Title
	}

	if req.IssueSource == "" {
		req.IssueSource = req.Source
	}

	if req.Level == "" {
		req.Level = req.Severity
	}

	if req.OwnerName == "" {
		req.OwnerName = req.Owner
	}

	result, err := config.DB.Exec(`
		UPDATE issues
		SET
			project_id = ?,
			device_type = ?,
			issue_source = ?,
			level = ?,
			issue_title = ?,
			issue_desc = ?,
			owner_id = ?,
			owner_name = ?,
			plan_close_time = NULLIF(?, '')
		WHERE id = ? AND IFNULL(is_deleted, 0) = 0
	`,
		req.ProjectID,
		req.DeviceType,
		req.IssueSource,
		req.Level,
		req.IssueTitle,
		req.IssueDesc,
		req.OwnerID,
		req.OwnerName,
		req.PlanCloseTime,
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

func ReplyIssueHandler(w http.ResponseWriter, r *http.Request, id int64) {
	w.Header().Set("Content-Type", "application/json")

	var req model.IssueReply
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "参数解析失败: "+err.Error(), http.StatusBadRequest)
		return
	}

	if req.Content == "" {
		req.Content = req.ReplyContent
	}

	if req.ReplyUserName == "" {
		req.ReplyUserName = req.ReplyUser
	}

	if req.Content == "" {
		http.Error(w, "回复内容不能为空", http.StatusBadRequest)
		return
	}

	_, err := config.DB.Exec(`
		INSERT INTO issue_replies (
			issue_id,
			reply_user_id,
			reply_user_name,
			reply_time,
			content
		) VALUES (?, ?, ?, NOW(), ?)
	`,
		id,
		req.ReplyUserID,
		req.ReplyUserName,
		req.Content,
	)
	if err != nil {
		http.Error(w, "回复失败: "+err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"code": 200,
		"msg":  "回复成功",
	})
}

func CloseIssueHandler(w http.ResponseWriter, r *http.Request, id int64) {
	w.Header().Set("Content-Type", "application/json")

	var req model.Issue
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "参数解析失败: "+err.Error(), http.StatusBadRequest)
		return
	}

	if req.CloseUserName == "" {
		req.CloseUserName = req.CloseUser
	}

	result, err := config.DB.Exec(`
		UPDATE issues
		SET
			close_status = '已关闭',
			close_user_id = ?,
			close_user_name = ?,
			real_close_time = NOW()
		WHERE id = ? AND IFNULL(is_deleted, 0) = 0
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
		http.Error(w, "问题不存在或已删除", http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"code": 200,
		"msg":  "关闭成功",
	})
}

func ReopenIssueHandler(w http.ResponseWriter, r *http.Request, id int64) {
	w.Header().Set("Content-Type", "application/json")

	var req model.Issue
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "参数解析失败: "+err.Error(), http.StatusBadRequest)
		return
	}

	result, err := config.DB.Exec(`
		UPDATE issues
		SET
			close_status = '未关闭',
			close_user_id = 0,
			close_user_name = '',
			real_close_time = NULL,
			reopen_reason = ?,
			reopen_count = IFNULL(reopen_count, 0) + 1
		WHERE id = ? AND IFNULL(is_deleted, 0) = 0
	`,
		req.ReopenReason,
		id,
	)
	if err != nil {
		http.Error(w, "重新打开失败: "+err.Error(), http.StatusInternalServerError)
		return
	}

	affected, _ := result.RowsAffected()
	if affected == 0 {
		http.Error(w, "问题不存在或已删除", http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"code": 200,
		"msg":  "重新打开成功",
	})
}

func DeleteIssueHandler(w http.ResponseWriter, r *http.Request, id int64) {
	w.Header().Set("Content-Type", "application/json")

	result, err := config.DB.Exec(`
		UPDATE issues
		SET is_deleted = 1
		WHERE id = ? AND IFNULL(is_deleted, 0) = 0
	`, id)
	if err != nil {
		http.Error(w, "删除失败: "+err.Error(), http.StatusInternalServerError)
		return
	}

	affected, _ := result.RowsAffected()
	if affected == 0 {
		http.Error(w, "问题不存在或已删除", http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"code": 200,
		"msg":  "删除成功",
	})
}
