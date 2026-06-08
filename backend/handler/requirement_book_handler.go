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

// 总入口，处理 GET/POST
func RequirementBooksHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	switch r.Method {
	case http.MethodGet:
		GetRequirementBooksHandler(w, r)
	case http.MethodPost:
		CreateRequirementBookHandler(w, r)
	default:
		http.Error(w, "不支持该请求方法", http.StatusMethodNotAllowed)
	}
}

// 带 ID 的操作入口，处理提交、删除
func RequirementBookActionHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	path := strings.TrimPrefix(r.URL.Path, "/api/requirement-books/")
	path = strings.Trim(path, "/")

	if path == "" {
		http.Error(w, "缺少需求书ID", http.StatusBadRequest)
		return
	}

	parts := strings.Split(path, "/")

	id, err := strconv.ParseInt(parts[0], 10, 64)
	if err != nil {
		http.Error(w, "需求书ID错误", http.StatusBadRequest)
		return
	}

	// POST /api/requirement-books/1 默认提交
	if len(parts) == 1 && r.Method == http.MethodPost {
		SubmitRequirementBookHandler(w, r, id)
		return
	}

	// DELETE /api/requirement-books/1
	if len(parts) == 1 && r.Method == http.MethodDelete {
		DeleteRequirementBookHandler(w, r, id)
		return
	}

	// POST /api/requirement-books/1/submit
	if len(parts) == 2 && r.Method == http.MethodPost && parts[1] == "submit" {
		SubmitRequirementBookHandler(w, r, id)
		return
	}

	http.Error(w, "接口不存在", http.StatusNotFound)
}

// ---------------- 内部函数 -----------------

func GetRequirementBooksHandler(w http.ResponseWriter, r *http.Request) {
	rows, err := config.DB.Query(`
		SELECT id, project_id, book_name, file_id, status,
		       submit_user_id, submit_user_name, submit_time,
		       audit_user_id, audit_user_name, audit_time,
		       reject_reason, remark, created_at, updated_at, is_deleted
		FROM requirement_books
		WHERE is_deleted=0
		ORDER BY id DESC
	`)
	if err != nil {
		http.Error(w, "查询失败: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	list := make([]model.RequirementBook, 0)

	for rows.Next() {
		var rbook model.RequirementBook
		err := rows.Scan(
			&rbook.ID,
			&rbook.ProjectID,
			&rbook.BookName,
			&rbook.FileID,
			&rbook.Status,
			&rbook.SubmitUserID,
			&rbook.SubmitUserName,
			&rbook.SubmitTime,
			&rbook.AuditUserID,
			&rbook.AuditUserName,
			&rbook.AuditTime,
			&rbook.RejectReason,
			&rbook.Remark,
			&rbook.CreatedAt,
			&rbook.UpdatedAt,
			&rbook.IsDeleted,
		)
		if err != nil {
			http.Error(w, "数据解析失败: "+err.Error(), http.StatusInternalServerError)
			return
		}
		list = append(list, rbook)
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"code": 200,
		"msg":  "查询成功",
		"data": list,
	})
}

func CreateRequirementBookHandler(w http.ResponseWriter, r *http.Request) {
	var rbook model.RequirementBook
	err := json.NewDecoder(r.Body).Decode(&rbook)
	if err != nil {
		http.Error(w, "参数解析失败: "+err.Error(), http.StatusBadRequest)
		return
	}

	if rbook.ProjectID == 0 || rbook.BookName == "" || rbook.FileID == 0 {
		http.Error(w, "缺少必要参数", http.StatusBadRequest)
		return
	}

	now := time.Now()
	result, err := config.DB.Exec(`
		INSERT INTO requirement_books
		(project_id, book_name, file_id, status, submit_user_id, submit_user_name, created_at, updated_at, is_deleted)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, 0)
	`, rbook.ProjectID, rbook.BookName, rbook.FileID, "草稿", rbook.SubmitUserID, rbook.SubmitUserName, now, now)

	if err != nil {
		http.Error(w, "新增失败: "+err.Error(), http.StatusInternalServerError)
		return
	}

	id, _ := result.LastInsertId()
	rbook.ID = id
	rbook.CreatedAt = now
	rbook.UpdatedAt = now
	rbook.Status = "草稿"

	json.NewEncoder(w).Encode(map[string]interface{}{
		"code": 200,
		"msg":  "新增成功",
		"data": rbook,
	})
}

func SubmitRequirementBookHandler(w http.ResponseWriter, r *http.Request, id int64) {
	result, err := config.DB.Exec(`
		UPDATE requirement_books
		SET status='待审核', submit_time=NOW(), updated_at=NOW()
		WHERE id=? AND is_deleted=0
	`, id)

	if err != nil {
		http.Error(w, "提交失败: "+err.Error(), http.StatusInternalServerError)
		return
	}

	affected, _ := result.RowsAffected()
	if affected == 0 {
		http.Error(w, "需求书不存在或已删除", http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"code": 200,
		"msg":  "提交成功",
	})
}

func DeleteRequirementBookHandler(w http.ResponseWriter, r *http.Request, id int64) {
	result, err := config.DB.Exec(`
		UPDATE requirement_books
		SET is_deleted=1, updated_at=NOW()
		WHERE id=?
	`, id)

	if err != nil {
		http.Error(w, "删除失败: "+err.Error(), http.StatusInternalServerError)
		return
	}

	affected, _ := result.RowsAffected()
	if affected == 0 {
		http.Error(w, "需求书不存在", http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"code": 200,
		"msg":  "删除成功",
	})
}