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

// GET /api/customer-supplied-files
// POST /api/customer-supplied-files
func CustomerSuppliedFilesHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	switch r.Method {
	case http.MethodGet:
		GetCustomerSuppliedFilesHandler(w, r)
	case http.MethodPost:
		CreateCustomerSuppliedFileHandler(w, r)
	default:
		http.Error(w, "不支持该请求方法", http.StatusMethodNotAllowed)
	}
}

// DELETE /api/customer-supplied-files/{id}
func CustomerSuppliedFileActionHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	path := strings.TrimPrefix(r.URL.Path, "/api/customer-supplied-files/")
	path = strings.Trim(path, "/")

	if path == "" {
		http.Error(w, "缺少客供资料ID", http.StatusBadRequest)
		return
	}

	parts := strings.Split(path, "/")

	id, err := strconv.ParseInt(parts[0], 10, 64)
	if err != nil {
		http.Error(w, "客供资料ID错误", http.StatusBadRequest)
		return
	}

	if len(parts) == 1 && r.Method == http.MethodDelete {
		DeleteCustomerSuppliedFileHandler(w, r, id)
		return
	}

	http.Error(w, "接口不存在", http.StatusNotFound)
}

func GetCustomerSuppliedFilesHandler(w http.ResponseWriter, r *http.Request) {
	rows, err := config.DB.Query(`
		SELECT
			id,
			project_id,
			file_id,
			material_name,
			IFNULL(file_display_name, ''),
			IFNULL(material_desc, ''),
			IFNULL(upload_user_id, 0),
			IFNULL(upload_user_name, ''),
			upload_time,
			IFNULL(remark, ''),
			created_at,
			updated_at,
			is_deleted
		FROM customer_supplied_files
		WHERE is_deleted = 0
		ORDER BY id DESC
	`)
	if err != nil {
		http.Error(w, "查询失败: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	list := make([]model.CustomerSuppliedFile, 0)

	for rows.Next() {
		var item model.CustomerSuppliedFile

		err := rows.Scan(
			&item.ID,
			&item.ProjectID,
			&item.FileID,
			&item.MaterialName,
			&item.FileDisplayName,
			&item.MaterialDesc,
			&item.UploadUserID,
			&item.UploadUserName,
			&item.UploadTime,
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

func CreateCustomerSuppliedFileHandler(w http.ResponseWriter, r *http.Request) {
	var item model.CustomerSuppliedFile

	err := json.NewDecoder(r.Body).Decode(&item)
	if err != nil {
		http.Error(w, "参数解析失败: "+err.Error(), http.StatusBadRequest)
		return
	}

	if item.ProjectID == 0 {
		http.Error(w, "项目ID不能为空", http.StatusBadRequest)
		return
	}

	if item.FileID == 0 {
		http.Error(w, "文件ID不能为空", http.StatusBadRequest)
		return
	}

	if item.MaterialName == "" {
		http.Error(w, "客供资料名称不能为空", http.StatusBadRequest)
		return
	}

	if item.FileDisplayName == "" {
		item.FileDisplayName = item.MaterialName
	}

	now := time.Now()

	result, err := config.DB.Exec(`
		INSERT INTO customer_supplied_files (
			project_id,
			file_id,
			material_name,
			file_display_name,
			material_desc,
			upload_user_id,
			upload_user_name,
			upload_time,
			remark,
			created_at,
			updated_at,
			is_deleted
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, 0)
	`,
		item.ProjectID,
		item.FileID,
		item.MaterialName,
		item.FileDisplayName,
		item.MaterialDesc,
		item.UploadUserID,
		item.UploadUserName,
		now,
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
	item.CreatedAt = now
	item.UpdatedAt = now
	item.IsDeleted = false

	json.NewEncoder(w).Encode(map[string]interface{}{
		"code": 200,
		"msg":  "新增成功",
		"data": map[string]interface{}{
			"id": id,
		},
	})
}

func DeleteCustomerSuppliedFileHandler(w http.ResponseWriter, r *http.Request, id int64) {
	if !hasRequestRole(r, "project_assistant") {
		http.Error(w, "无删除权限：客供资料仅项目助理可删除", http.StatusForbidden)
		return
	}

	result, err := config.DB.Exec(`
		UPDATE customer_supplied_files
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
		http.Error(w, "客供资料不存在", http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"code": 200,
		"msg":  "删除成功",
	})
}
