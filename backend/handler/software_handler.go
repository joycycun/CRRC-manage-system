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
// 软件版本入口
// GET  /api/software-versions
// POST /api/software-versions
// ============================================================

func SoftwareVersionsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	switch r.Method {
	case http.MethodGet:
		GetSoftwareVersionsHandler(w, r)
	case http.MethodPost:
		CreateSoftwareVersionHandler(w, r)
	default:
		http.Error(w, "不支持该请求方法", http.StatusMethodNotAllowed)
	}
}

// ============================================================
// 软件版本带 ID 操作入口
// PUT  /api/software-versions/{id}
// POST /api/software-versions/{id}/release
// POST /api/software-versions/{id}/discard
// DELETE /api/software-versions/{id}
// ============================================================

func SoftwareVersionActionHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	path := strings.TrimPrefix(r.URL.Path, "/api/software-versions/")
	path = strings.Trim(path, "/")

	if path == "" {
		http.Error(w, "缺少软件版本ID", http.StatusBadRequest)
		return
	}

	parts := strings.Split(path, "/")

	id, err := strconv.ParseInt(parts[0], 10, 64)
	if err != nil {
		http.Error(w, "软件版本ID错误", http.StatusBadRequest)
		return
	}

	if len(parts) == 1 && r.Method == http.MethodPut {
		UpdateSoftwareVersionHandler(w, r, id)
		return
	}

	if len(parts) == 1 && r.Method == http.MethodDelete {
		DeleteSoftwareVersionHandler(w, r, id)
		return
	}

	if len(parts) == 2 && r.Method == http.MethodPost {
		action := parts[1]

		switch action {
		case "release":
			ReleaseSoftwareVersionHandler(w, r, id)
			return
		case "discard":
			DiscardSoftwareVersionHandler(w, r, id)
			return
		default:
			http.Error(w, "不支持的操作", http.StatusNotFound)
			return
		}
	}

	http.Error(w, "接口不存在", http.StatusNotFound)
}

// ============================================================
// GET /api/software-versions
// ============================================================

func GetSoftwareVersionsHandler(w http.ResponseWriter, r *http.Request) {
	rows, err := config.DB.Query(`
		SELECT
			id,
			IFNULL(project_id, 0),
			software_version,
			IFNULL(device_type, ''),
			IFNULL(hardware_id, 0),
			IFNULL(hardware_version, ''),
			IFNULL(owner_id, 0),
			IFNULL(owner_name, ''),
			release_date,
			IFNULL(download_url, ''),
			IFNULL(business_desc, ''),
			IFNULL(description, ''),
			IFNULL(software_status, ''),
			created_at,
			updated_at,
			is_deleted
		FROM software_versions
		WHERE is_deleted = 0
		ORDER BY id DESC
	`)
	if err != nil {
		http.Error(w, "查询失败: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	list := make([]model.SoftwareVersionResponse, 0)

	for rows.Next() {
		var item model.SoftwareVersion

		err := rows.Scan(
			&item.ID,
			&item.ProjectID,
			&item.SoftwareVersion,
			&item.DeviceType,
			&item.HardwareID,
			&item.HardwareVersion,
			&item.OwnerID,
			&item.OwnerName,
			&item.ReleaseDate,
			&item.DownloadURL,
			&item.BusinessDesc,
			&item.Description,
			&item.SoftwareStatus,
			&item.CreatedAt,
			&item.UpdatedAt,
			&item.IsDeleted,
		)
		if err != nil {
			http.Error(w, "数据解析失败: "+err.Error(), http.StatusInternalServerError)
			return
		}

		list = append(list, item.ToResponse())
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"code": 200,
		"msg":  "查询成功",
		"data": list,
	})
}

// ============================================================
// POST /api/software-versions
// ============================================================

func CreateSoftwareVersionHandler(w http.ResponseWriter, r *http.Request) {
	var item model.SoftwareVersion

	err := json.NewDecoder(r.Body).Decode(&item)
	if err != nil {
		http.Error(w, "参数解析失败: "+err.Error(), http.StatusBadRequest)
		return
	}

	if item.SoftwareVersion == "" {
		http.Error(w, "软件版本号不能为空", http.StatusBadRequest)
		return
	}

	if item.DeviceType == "" {
		http.Error(w, "终端类型不能为空", http.StatusBadRequest)
		return
	}

	if item.SoftwareStatus == "" {
		item.SoftwareStatus = "草稿"
	}

	now := time.Now()

	result, err := config.DB.Exec(`
		INSERT INTO software_versions (
			project_id,
			software_version,
			device_type,
			hardware_id,
			hardware_version,
			owner_id,
			owner_name,
			download_url,
			business_desc,
			description,
			software_status,
			created_at,
			updated_at,
			is_deleted
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, 0)
	`,
		item.ProjectID,
		item.SoftwareVersion,
		item.DeviceType,
		item.HardwareID,
		item.HardwareVersion,
		item.OwnerID,
		item.OwnerName,
		item.DownloadURL,
		item.BusinessDesc,
		item.Description,
		item.SoftwareStatus,
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

// ============================================================
// PUT /api/software-versions/{id}
// ============================================================

func UpdateSoftwareVersionHandler(w http.ResponseWriter, r *http.Request, id int64) {
	var item model.SoftwareVersion

	err := json.NewDecoder(r.Body).Decode(&item)
	if err != nil {
		http.Error(w, "参数解析失败: "+err.Error(), http.StatusBadRequest)
		return
	}

	if item.SoftwareVersion == "" {
		http.Error(w, "软件版本号不能为空", http.StatusBadRequest)
		return
	}

	if item.DeviceType == "" {
		http.Error(w, "终端类型不能为空", http.StatusBadRequest)
		return
	}

	result, err := config.DB.Exec(`
		UPDATE software_versions
		SET
			project_id = ?,
			software_version = ?,
			device_type = ?,
			hardware_id = ?,
			hardware_version = ?,
			owner_id = ?,
			owner_name = ?,
			download_url = ?,
			business_desc = ?,
			description = ?,
			software_status = ?,
			updated_at = NOW()
		WHERE id = ? AND is_deleted = 0
	`,
		item.ProjectID,
		item.SoftwareVersion,
		item.DeviceType,
		item.HardwareID,
		item.HardwareVersion,
		item.OwnerID,
		item.OwnerName,
		item.DownloadURL,
		item.BusinessDesc,
		item.Description,
		item.SoftwareStatus,
		id,
	)

	if err != nil {
		http.Error(w, "修改失败: "+err.Error(), http.StatusInternalServerError)
		return
	}

	affected, _ := result.RowsAffected()
	if affected == 0 {
		http.Error(w, "软件版本不存在或已删除", http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"code": 200,
		"msg":  "修改成功",
	})
}

type SoftwareReleaseRequest struct {
	ReleaseDate string `json:"releaseDate"`
	DownloadURL string `json:"downloadUrl"`
}

// ============================================================
// POST /api/software-versions/{id}/release
// ============================================================

func ReleaseSoftwareVersionHandler(w http.ResponseWriter, r *http.Request, id int64) {
	var req SoftwareReleaseRequest
	_ = json.NewDecoder(r.Body).Decode(&req)

	if req.ReleaseDate == "" {
		req.ReleaseDate = time.Now().Format("2006-01-02")
	}

	result, err := config.DB.Exec(`
		UPDATE software_versions
		SET
			software_status = '已发布',
			release_date = ?,
			download_url = IF(? = '', download_url, ?),
			updated_at = NOW()
		WHERE id = ? AND is_deleted = 0
	`,
		req.ReleaseDate,
		req.DownloadURL,
		req.DownloadURL,
		id,
	)

	if err != nil {
		http.Error(w, "发布失败: "+err.Error(), http.StatusInternalServerError)
		return
	}

	affected, _ := result.RowsAffected()
	if affected == 0 {
		http.Error(w, "软件版本不存在或已删除", http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"code": 200,
		"msg":  "发布成功",
	})
}

// ============================================================
// POST /api/software-versions/{id}/discard
// ============================================================

func DiscardSoftwareVersionHandler(w http.ResponseWriter, r *http.Request, id int64) {
	result, err := config.DB.Exec(`
		UPDATE software_versions
		SET
			software_status = '已废弃',
			updated_at = NOW()
		WHERE id = ? AND is_deleted = 0
	`, id)

	if err != nil {
		http.Error(w, "废弃失败: "+err.Error(), http.StatusInternalServerError)
		return
	}

	affected, _ := result.RowsAffected()
	if affected == 0 {
		http.Error(w, "软件版本不存在或已删除", http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"code": 200,
		"msg":  "废弃成功",
	})
}

// ============================================================
// DELETE /api/software-versions/{id}
// ============================================================

func DeleteSoftwareVersionHandler(w http.ResponseWriter, r *http.Request, id int64) {
	result, err := config.DB.Exec(`
		UPDATE software_versions
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
		http.Error(w, "软件版本不存在", http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"code": 200,
		"msg":  "删除成功",
	})
}

// ============================================================
// 项目分支入口
// GET  /api/project-branches
// POST /api/project-branches
// ============================================================

func ProjectBranchesHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	switch r.Method {
	case http.MethodGet:
		GetProjectBranchesHandler(w, r)
	case http.MethodPost:
		CreateProjectBranchHandler(w, r)
	default:
		http.Error(w, "不支持该请求方法", http.StatusMethodNotAllowed)
	}
}

// ============================================================
// 项目分支带 ID 操作入口
// DELETE /api/project-branches/{id}
// ============================================================

func ProjectBranchActionHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	path := strings.TrimPrefix(r.URL.Path, "/api/project-branches/")
	path = strings.Trim(path, "/")

	if path == "" {
		http.Error(w, "缺少项目分支ID", http.StatusBadRequest)
		return
	}

	parts := strings.Split(path, "/")

	id, err := strconv.ParseInt(parts[0], 10, 64)
	if err != nil {
		http.Error(w, "项目分支ID错误", http.StatusBadRequest)
		return
	}

	if len(parts) == 1 && r.Method == http.MethodDelete {
		DeleteProjectBranchHandler(w, r, id)
		return
	}

	http.Error(w, "接口不存在", http.StatusNotFound)
}

// ============================================================
// GET /api/project-branches
// ============================================================

func GetProjectBranchesHandler(w http.ResponseWriter, r *http.Request) {
	rows, err := config.DB.Query(`
		SELECT
			id,
			project_id,
			IFNULL(device_type, ''),
			IFNULL(repo_name, ''),
			IFNULL(repo_url, ''),
			branch_name,
			IFNULL(clone_url, ''),
			IFNULL(owner_id, 0),
			IFNULL(owner_name, ''),
			IFNULL(remark, ''),
			created_at,
			updated_at,
			is_deleted
		FROM project_branches
		WHERE is_deleted = 0
		ORDER BY id DESC
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
			&item.DeviceType,
			&item.RepoName,
			&item.RepoURL,
			&item.BranchName,
			&item.CloneURL,
			&item.OwnerID,
			&item.OwnerName,
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

// ============================================================
// POST /api/project-branches
// ============================================================

func CreateProjectBranchHandler(w http.ResponseWriter, r *http.Request) {
	var item model.ProjectBranch

	err := json.NewDecoder(r.Body).Decode(&item)
	if err != nil {
		http.Error(w, "参数解析失败: "+err.Error(), http.StatusBadRequest)
		return
	}

	if item.ProjectID == 0 {
		http.Error(w, "项目ID不能为空", http.StatusBadRequest)
		return
	}

	if item.BranchName == "" {
		http.Error(w, "分支名称不能为空", http.StatusBadRequest)
		return
	}

	now := time.Now()

	result, err := config.DB.Exec(`
		INSERT INTO project_branches (
			project_id,
			device_type,
			repo_name,
			repo_url,
			branch_name,
			clone_url,
			owner_id,
			owner_name,
			remark,
			created_at,
			updated_at,
			is_deleted
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, 0)
	`,
		item.ProjectID,
		item.DeviceType,
		item.RepoName,
		item.RepoURL,
		item.BranchName,
		item.CloneURL,
		item.OwnerID,
		item.OwnerName,
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
		"data": map[string]interface{}{
			"id": id,
		},
	})
}

// ============================================================
// DELETE /api/project-branches/{id}
// ============================================================

func DeleteProjectBranchHandler(w http.ResponseWriter, r *http.Request, id int64) {
	result, err := config.DB.Exec(`
		UPDATE project_branches
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
		http.Error(w, "项目分支不存在", http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"code": 200,
		"msg":  "删除成功",
	})
}
