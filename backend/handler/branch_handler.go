package handler

import (
	"crrc_pm_backend/config"
	"crrc_pm_backend/model"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
)

func BranchesHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		GetBranchesHandler(w, r)
	case http.MethodPost:
		CreateBranchHandler(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func BranchActionHandler(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimPrefix(r.URL.Path, "/api/branches/")
	parts := strings.Split(path, "/")

	if len(parts) < 1 || parts[0] == "" {
		http.NotFound(w, r)
		return
	}

	id, err := strconv.ParseInt(parts[0], 10, 64)
	if err != nil {
		http.Error(w, "无效的分支ID", http.StatusBadRequest)
		return
	}

	if len(parts) == 1 && r.Method == http.MethodPut {
		UpdateBranchHandler(w, r, id)
		return
	}

	if len(parts) == 1 && r.Method == http.MethodDelete {
		DeleteBranchHandler(w, r, id)
		return
	}

	http.NotFound(w, r)
}

func GetBranchesHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	rows, err := config.DB.Query(`
		SELECT
			pb.id,
			IFNULL(pb.project_id, 0),
			IFNULL(p.project_name, ''),
			IFNULL(pb.device_type, ''),
			IFNULL(pb.branch_name, ''),
			IFNULL(pb.owner_id, 0),
			IFNULL(pb.owner_name, ''),
			IFNULL(DATE_FORMAT(pb.created_at, '%Y-%m-%d'), ''),
			IFNULL(pb.clone_url, ''),
			IFNULL(DATE_FORMAT(pb.created_at, '%Y-%m-%d %H:%i:%s'), ''),
			IFNULL(DATE_FORMAT(pb.updated_at, '%Y-%m-%d %H:%i:%s'), '')
		FROM project_branches pb
		LEFT JOIN projects p ON pb.project_id = p.id
		WHERE pb.is_deleted = 0
		ORDER BY pb.id DESC
	`)
	if err != nil {
		http.Error(w, "查询失败: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	list := make([]model.ProjectBranch, 0)

	for rows.Next() {
		var item model.ProjectBranch

		err := rows.Scan(
			&item.ID,
			&item.ProjectID,
			&item.ProjectName,
			&item.DeviceType,
			&item.BranchName,
			&item.OwnerID,
			&item.OwnerName,
			&item.CreateTime,
			&item.CloneURL,
			&item.CreatedAt,
			&item.UpdatedAt,
		)
		if err != nil {
			http.Error(w, "数据解析失败: "+err.Error(), http.StatusInternalServerError)
			return
		}

		item.Owner = item.OwnerName
		list = append(list, item)
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"code": 200,
		"msg":  "查询成功",
		"data": list,
	})
}

func CreateBranchHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var req model.ProjectBranch
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "参数解析失败: "+err.Error(), http.StatusBadRequest)
		return
	}

	if req.ProjectID == 0 {
		http.Error(w, "项目ID不能为空", http.StatusBadRequest)
		return
	}

	if req.BranchName == "" {
		http.Error(w, "分支名称不能为空", http.StatusBadRequest)
		return
	}

	if req.OwnerName == "" {
		req.OwnerName = req.Owner
	}

	if req.RepoName == "" {
		req.RepoName = req.BranchName
	}

	if req.RepoURL == "" {
		req.RepoURL = req.CloneURL
	}

	if req.OwnerName == "" {
		req.OwnerName = req.Owner
	}

	result, err := config.DB.Exec(`
	INSERT INTO project_branches (
		project_id,
		device_type,
		branch_name,
		repo_name,
		repo_url,
		clone_url,
		owner_id,
		owner_name,
		remark,
		is_deleted,
		created_at,
		updated_at
	) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, 0, NOW(), NOW())
`,
		req.ProjectID,
		req.DeviceType,
		req.BranchName,
		req.RepoName,
		req.RepoURL,
		req.CloneURL,
		req.OwnerID,
		req.OwnerName,
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

func UpdateBranchHandler(w http.ResponseWriter, r *http.Request, id int64) {
	w.Header().Set("Content-Type", "application/json")

	var req model.ProjectBranch
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "参数解析失败: "+err.Error(), http.StatusBadRequest)
		return
	}

	if req.OwnerName == "" {
		req.OwnerName = req.Owner
	}

	result, err := config.DB.Exec(`
		UPDATE project_branches
		SET
			project_id = ?,
			device_type = ?,
			branch_name = ?,
			owner_id = ?,
			owner_name = ?,
			create_time = ?,
			clone_url = ?,
			updated_at = NOW()
		WHERE id = ? AND is_deleted = 0
	`,
		req.ProjectID,
		req.DeviceType,
		req.BranchName,
		req.OwnerID,
		req.OwnerName,
		req.CreateTime,
		req.CloneURL,
		id,
	)
	if err != nil {
		http.Error(w, "修改失败: "+err.Error(), http.StatusInternalServerError)
		return
	}

	affected, _ := result.RowsAffected()
	if affected == 0 {
		http.Error(w, "项目分支不存在或已删除", http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"code": 200,
		"msg":  "修改成功",
	})
}

func DeleteBranchHandler(w http.ResponseWriter, r *http.Request, id int64) {
	w.Header().Set("Content-Type", "application/json")

	result, err := config.DB.Exec(`
		UPDATE project_branches
		SET is_deleted = 1, updated_at = NOW()
		WHERE id = ? AND is_deleted = 0
	`, id)
	if err != nil {
		http.Error(w, "删除失败: "+err.Error(), http.StatusInternalServerError)
		return
	}

	affected, _ := result.RowsAffected()
	if affected == 0 {
		http.Error(w, "项目分支不存在或已删除", http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"code": 200,
		"msg":  "删除成功",
	})
}
