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

func ProjectsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	switch r.Method {
	case http.MethodGet:
		GetProjects(w, r)
	case http.MethodPost:
		CreateProject(w, r)
	default:
		http.Error(w, "不支持该请求方法", http.StatusMethodNotAllowed)
	}
}

func ProjectActionHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	path := strings.TrimPrefix(r.URL.Path, "/api/projects/")
	path = strings.Trim(path, "/")

	if path == "" {
		http.Error(w, "缺少项目ID", http.StatusBadRequest)
		return
	}

	parts := strings.Split(path, "/")

	id, err := strconv.ParseInt(parts[0], 10, 64)
	if err != nil {
		http.Error(w, "项目ID错误", http.StatusBadRequest)
		return
	}

	// 现在先兼容你的测试方式：POST /api/projects/1
	// 也就是默认表示提交立项
	if len(parts) == 1 && r.Method == http.MethodPost {
		SubmitProject(w, r, id)
		return
	}

	// 后面我们会支持：
	// POST /api/projects/1/submit
	// POST /api/projects/1/audit
	// POST /api/projects/1/archive
	if len(parts) == 2 {
		action := parts[1]

		if r.Method == http.MethodPost && action == "submit" {
			SubmitProject(w, r, id)
			return
		}
	}

	http.Error(w, "接口不存在", http.StatusNotFound)
}

func GetProjects(w http.ResponseWriter, r *http.Request) {
	rows, err := config.DB.Query(`
		SELECT 
			id,
			project_name,
			project_code,
			IFNULL(owner_id, 0),
			IFNULL(owner_name, ''),
			IFNULL(stage, ''),
			IFNULL(status, ''),
			submit_time,
			IFNULL(audit_status, ''),
			IFNULL(audit_user_id, 0),
			IFNULL(audit_user_name, ''),
			audit_time,
			archive_time,
			IFNULL(remark, ''),
			created_at,
			updated_at,
			close_time
		FROM projects
		ORDER BY id DESC
	`)
	if err != nil {
		http.Error(w, "查询失败: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	list := make([]model.ProjectResponse, 0)

	for rows.Next() {
		var p model.Project

		err := rows.Scan(
			&p.ID,
			&p.ProjectName,
			&p.ProjectCode,
			&p.OwnerID,
			&p.OwnerName,
			&p.Stage,
			&p.Status,
			&p.SubmitTime,
			&p.AuditStatus,
			&p.AuditUserID,
			&p.AuditUserName,
			&p.AuditTime,
			&p.ArchiveTime,
			&p.Remark,
			&p.CreatedAt,
			&p.UpdatedAt,
			&p.CloseTime,
		)
		if err != nil {
			http.Error(w, "数据解析失败: "+err.Error(), http.StatusInternalServerError)
			return
		}

		list = append(list, p.ToResponse())
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"code": 200,
		"msg":  "查询成功",
		"data": list,
	})
}

func CreateProject(w http.ResponseWriter, r *http.Request) {
	var p model.Project

	err := json.NewDecoder(r.Body).Decode(&p)
	if err != nil {
		http.Error(w, "参数解析失败: "+err.Error(), http.StatusBadRequest)
		return
	}

	if p.ProjectName == "" {
		http.Error(w, "项目名称不能为空", http.StatusBadRequest)
		return
	}

	if p.ProjectCode == "" {
		http.Error(w, "项目编号不能为空", http.StatusBadRequest)
		return
	}

	if p.Status == "" {
		p.Status = "立项中"
	}

	if p.Stage == "" {
		p.Stage = "立项阶段"
	}

	now := time.Now()

	result, err := config.DB.Exec(`
		INSERT INTO projects (
			project_name,
			project_code,
			owner_id,
			owner_name,
			stage,
			status,
			audit_status,
			remark,
			created_at,
			updated_at
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`,
		p.ProjectName,
		p.ProjectCode,
		p.OwnerID,
		p.OwnerName,
		p.Stage,
		p.Status,
		"未提交",
		p.Remark,
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

func SubmitProject(w http.ResponseWriter, r *http.Request, id int64) {
	result, err := config.DB.Exec(`
		UPDATE projects
		SET 
			status = '立项中',
			audit_status = '待审核',
			submit_time = NOW(),
			updated_at = NOW()
		WHERE id = ?
	`, id)

	if err != nil {
		http.Error(w, "提交失败: "+err.Error(), http.StatusInternalServerError)
		return
	}

	affected, _ := result.RowsAffected()
	if affected == 0 {
		http.Error(w, "项目不存在", http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"code": 200,
		"msg":  "提交成功",
	})
}
