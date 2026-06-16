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

// ============================================================
// 项目入口
// GET  /api/projects
// POST /api/projects
// ============================================================

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

// ============================================================
// 项目带 ID 操作入口
// PUT    /api/projects/{id}
// POST   /api/projects/{id}/submit
// POST   /api/projects/{id}/audit
// POST   /api/projects/{id}/archive
// POST   /api/projects/{id}/close
// DELETE /api/projects/{id}
// ============================================================

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

	if len(parts) == 1 && r.Method == http.MethodGet {
		GetProjectDetail(w, r, id)
		return
	}

	if len(parts) == 1 && r.Method == http.MethodPut {
		UpdateProject(w, r, id)
		return
	}

	if len(parts) == 1 && r.Method == http.MethodDelete {
		DeleteProject(w, r, id)
		return
	}

	if len(parts) == 2 && r.Method == http.MethodPost {
		switch parts[1] {
		case "submit":
			SubmitProject(w, r, id)
			return
		case "audit":
			AuditProject(w, r, id)
			return
		case "archive":
			ArchiveProject(w, r, id)
			return
		case "close":
			CloseProject(w, r, id)
			return
		default:
			http.Error(w, "不支持的操作", http.StatusNotFound)
			return
		}
	}

	http.Error(w, "接口不存在", http.StatusNotFound)
}

// ============================================================
// GET /api/projects
// ============================================================

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
			close_time,
			IFNULL(is_deleted, 0)
		FROM projects
		WHERE IFNULL(is_deleted, 0) = 0
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
			&p.IsDeleted,
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

// ============================================================
// GET /api/projects/{id}
// ============================================================

func GetProjectDetail(w http.ResponseWriter, r *http.Request, id int64) {
	var p model.Project

	err := config.DB.QueryRow(`
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
			close_time,
			IFNULL(is_deleted, 0)
		FROM projects
		WHERE id = ? AND IFNULL(is_deleted, 0) = 0
		LIMIT 1
	`, id).Scan(
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
		&p.IsDeleted,
	)

	if err != nil {
		http.Error(w, "项目不存在: "+err.Error(), http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"code": 200,
		"msg":  "查询成功",
		"data": p.ToResponse(),
	})
}

// ============================================================
// POST /api/projects
// ============================================================

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

	if p.OwnerID == 0 || !isSoftwareOwner(p.OwnerID) {
		http.Error(w, "请选择已注册的软件负责人", http.StatusBadRequest)
		return
	}

	var ownerName string
	if err := config.DB.QueryRow(`
		SELECT real_name
		FROM users
		WHERE id = ?
		  AND IFNULL(status, '启用') = '启用'
	`, p.OwnerID).Scan(&ownerName); err == nil && ownerName != "" {
		p.OwnerName = ownerName
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
			updated_at,
			is_deleted
		) VALUES (?, ?, ?, ?, ?, ?, '未提交', ?, ?, ?, 0)
	`,
		p.ProjectName,
		p.ProjectCode,
		p.OwnerID,
		p.OwnerName,
		p.Stage,
		p.Status,
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

func isSoftwareOwner(userID int64) bool {
	var count int
	err := config.DB.QueryRow(`
		SELECT COUNT(*)
		FROM users u
		JOIN user_roles ur ON ur.user_id = u.id
		JOIN roles r ON r.id = ur.role_id
		WHERE u.id = ?
		  AND r.role_code = 'software_owner'
		  AND IFNULL(u.status, '启用') = '启用'
	`, userID).Scan(&count)
	return err == nil && count > 0
}

// ============================================================
// PUT /api/projects/{id}
// ============================================================

func UpdateProject(w http.ResponseWriter, r *http.Request, id int64) {
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

	result, err := config.DB.Exec(`
		UPDATE projects
		SET
			project_name = ?,
			project_code = ?,
			owner_id = ?,
			owner_name = ?,
			stage = ?,
			status = ?,
			remark = ?,
			updated_at = NOW()
		WHERE id = ? AND IFNULL(is_deleted, 0) = 0
	`,
		p.ProjectName,
		p.ProjectCode,
		p.OwnerID,
		p.OwnerName,
		p.Stage,
		p.Status,
		p.Remark,
		id,
	)

	if err != nil {
		http.Error(w, "修改失败: "+err.Error(), http.StatusInternalServerError)
		return
	}

	affected, _ := result.RowsAffected()
	if affected == 0 {
		http.Error(w, "项目不存在或已删除", http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"code": 200,
		"msg":  "修改成功",
	})
}

// ============================================================
// POST /api/projects/{id}/submit
// ============================================================

func SubmitProject(w http.ResponseWriter, r *http.Request, id int64) {
	result, err := config.DB.Exec(`
		UPDATE projects
		SET 
			status = '立项中',
			audit_status = '待审核',
			submit_time = NOW(),
			updated_at = NOW()
		WHERE id = ? AND IFNULL(is_deleted, 0) = 0
	`, id)

	if err != nil {
		http.Error(w, "提交失败: "+err.Error(), http.StatusInternalServerError)
		return
	}

	affected, _ := result.RowsAffected()
	if affected == 0 {
		http.Error(w, "项目不存在或已删除", http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"code": 200,
		"msg":  "提交成功",
	})
}

// ============================================================
// POST /api/projects/{id}/audit
// ============================================================

type ProjectAuditRequest struct {
	AuditUserID   int64  `json:"auditUserId"`
	AuditUserName string `json:"auditUserName"`
	AuditStatus   string `json:"auditStatus"`
	RejectReason  string `json:"rejectReason"`
}

func AuditProject(w http.ResponseWriter, r *http.Request, id int64) {
	if !requireLeaderPermission(w, r) {
		return
	}

	var req ProjectAuditRequest

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "参数解析失败: "+err.Error(), http.StatusBadRequest)
		return
	}

	if req.AuditStatus != "已通过" && req.AuditStatus != "已驳回" {
		http.Error(w, "审核状态只能是 已通过 或 已驳回", http.StatusBadRequest)
		return
	}

	projectStatus := "立项中"
	stage := "立项阶段"

	if req.AuditStatus == "已通过" {
		projectStatus = "进行中"
		stage = "需求阶段"
	}

	result, err := config.DB.Exec(`
		UPDATE projects
		SET
			audit_status = ?,
			audit_user_id = ?,
			audit_user_name = ?,
			audit_time = NOW(),
			status = ?,
			stage = ?,
			remark = IF(? = '', remark, ?),
			updated_at = NOW()
		WHERE id = ? AND IFNULL(is_deleted, 0) = 0
	`,
		req.AuditStatus,
		req.AuditUserID,
		req.AuditUserName,
		projectStatus,
		stage,
		req.RejectReason,
		req.RejectReason,
		id,
	)

	if err != nil {
		http.Error(w, "审核失败: "+err.Error(), http.StatusInternalServerError)
		return
	}

	affected, _ := result.RowsAffected()
	if affected == 0 {
		http.Error(w, "项目不存在或已删除", http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"code": 200,
		"msg":  "审核成功",
	})
}

// ============================================================
// POST /api/projects/{id}/archive
// ============================================================

func ArchiveProject(w http.ResponseWriter, r *http.Request, id int64) {
	result, err := config.DB.Exec(`
		UPDATE projects
		SET
			status = '归档',
			stage = '已归档',
			archive_time = NOW(),
			updated_at = NOW()
		WHERE id = ? AND IFNULL(is_deleted, 0) = 0
	`, id)

	if err != nil {
		http.Error(w, "归档失败: "+err.Error(), http.StatusInternalServerError)
		return
	}

	affected, _ := result.RowsAffected()
	if affected == 0 {
		http.Error(w, "项目不存在或已删除", http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"code": 200,
		"msg":  "归档成功",
	})
}

// ============================================================
// POST /api/projects/{id}/close
// ============================================================

func CloseProject(w http.ResponseWriter, r *http.Request, id int64) {
	result, err := config.DB.Exec(`
		UPDATE projects
		SET
			status = '已关闭',
			stage = '已关闭',
			close_time = NOW(),
			updated_at = NOW()
		WHERE id = ? AND IFNULL(is_deleted, 0) = 0
	`, id)

	if err != nil {
		http.Error(w, "关闭失败: "+err.Error(), http.StatusInternalServerError)
		return
	}

	affected, _ := result.RowsAffected()
	if affected == 0 {
		http.Error(w, "项目不存在或已删除", http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"code": 200,
		"msg":  "关闭成功",
	})
}

// ============================================================
// DELETE /api/projects/{id}
// ============================================================

func DeleteProject(w http.ResponseWriter, r *http.Request, id int64) {
	result, err := config.DB.Exec(`
		UPDATE projects
		SET
			is_deleted = 1,
			updated_at = NOW()
		WHERE id = ?
		  AND IFNULL(is_deleted, 0) = 0
		  AND IFNULL(audit_status, '未提交') <> '已通过'
		  AND IFNULL(status, '') NOT IN ('进行中', '已关闭', '归档')
	`, id)

	if err != nil {
		http.Error(w, "删除失败: "+err.Error(), http.StatusInternalServerError)
		return
	}

	affected, _ := result.RowsAffected()
	if affected == 0 {
		http.Error(w, "项目不存在，或审核通过后的项目不允许删除", http.StatusBadRequest)
		return
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"code": 200,
		"msg":  "删除成功",
	})
}
