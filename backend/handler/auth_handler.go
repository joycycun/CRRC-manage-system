package handler

import (
	"crrc_pm_backend/config"
	"crrc_pm_backend/model"
	"database/sql"
	"encoding/json"
	"net/http"
	"time"
)

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodPost {
		json.NewEncoder(w).Encode(LoginResponse{
			Code: 405,
			Msg:  "请求方法错误",
			Data: nil,
		})
		return
	}

	var req LoginRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		json.NewEncoder(w).Encode(LoginResponse{
			Code: 400,
			Msg:  "请求参数错误",
			Data: nil,
		})
		return
	}

	var user model.User

	sqlStr := `
	SELECT id, username, password_hash, real_name, 
	       IFNULL(email, ''), IFNULL(phone, ''), 
	       IFNULL(department, ''), status
	FROM users
	WHERE username = ?
	LIMIT 1
`

	err = config.DB.QueryRow(sqlStr, req.Username).Scan(
		&user.ID,
		&user.Username,
		&user.PasswordHash,
		&user.RealName,
		&user.Email,
		&user.Phone,
		&user.Department,
		&user.Status,
	)

	if err == sql.ErrNoRows {
		json.NewEncoder(w).Encode(LoginResponse{
			Code: 401,
			Msg:  "账号不存在",
			Data: nil,
		})
		return
	}

	if err != nil {
		json.NewEncoder(w).Encode(LoginResponse{
			Code: 500,
			Msg:  "数据库查询失败: " + err.Error(),
			Data: nil,
		})
		return
	}

	if user.Status != "启用" {
		json.NewEncoder(w).Encode(LoginResponse{
			Code: 403,
			Msg:  "账号已禁用",
			Data: nil,
		})
		return
	}

	if user.PasswordHash != req.Password {
		json.NewEncoder(w).Encode(LoginResponse{
			Code: 401,
			Msg:  "密码错误",
			Data: nil,
		})
		return
	}

	_, _ = config.DB.Exec(
		"UPDATE users SET last_login_time = ? WHERE id = ?",
		time.Now(),
		user.ID,
	)

	roles, _ := queryUserRoles(user.ID)
	permissions, _ := queryUserPermissions(user.ID)

	if len(permissions) == 0 {
		permissions = []string{
			"project:view",
			"project:create",
			"project:audit",
		}
	}

	json.NewEncoder(w).Encode(LoginResponse{
		Code: 200,
		Msg:  "登录成功",
		Data: map[string]interface{}{
			"token":       "test-token",
			"user":        user,
			"roles":       roles,
			"permissions": permissions,
		},
	})
}

func SoftwareOwnersHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodGet {
		http.Error(w, "请求方法错误", http.StatusMethodNotAllowed)
		return
	}

	if !authTableExists("roles") || !authTableExists("user_roles") {
		http.Error(w, "角色表未初始化", http.StatusInternalServerError)
		return
	}

	rows, err := config.DB.Query(`
		SELECT
			u.id,
			u.username,
			u.real_name
		FROM users u
		JOIN user_roles ur ON ur.user_id = u.id
		JOIN roles r ON r.id = ur.role_id
		WHERE r.role_code = 'software_owner'
		  AND IFNULL(u.status, '启用') = '启用'
		ORDER BY u.id ASC
	`)
	if err != nil {
		http.Error(w, "查询软件负责人失败: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	list := make([]map[string]interface{}, 0)
	for rows.Next() {
		var id int64
		var username string
		var realName string
		if err := rows.Scan(&id, &username, &realName); err != nil {
			http.Error(w, "解析软件负责人失败: "+err.Error(), http.StatusInternalServerError)
			return
		}
		list = append(list, map[string]interface{}{
			"id":       id,
			"username": username,
			"name":     realName,
			"realName": realName,
		})
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"code": 200,
		"msg":  "查询成功",
		"data": list,
	})
}

func authTableExists(tableName string) bool {
	var count int
	err := config.DB.QueryRow(`
		SELECT COUNT(*)
		FROM information_schema.tables
		WHERE table_schema = DATABASE()
		  AND table_name = ?
	`, tableName).Scan(&count)
	return err == nil && count > 0
}

func queryUserRoles(userID int64) ([]map[string]interface{}, error) {
	if !authTableExists("roles") || !authTableExists("user_roles") {
		return []map[string]interface{}{}, nil
	}

	rows, err := config.DB.Query(`
		SELECT
			r.id,
			r.role_code,
			r.role_name,
			IFNULL(r.description, '')
		FROM user_roles ur
		JOIN roles r ON ur.role_id = r.id
		WHERE ur.user_id = ?
		ORDER BY r.id
	`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	list := make([]map[string]interface{}, 0)
	for rows.Next() {
		var id int64
		var code string
		var name string
		var description string
		if err := rows.Scan(&id, &code, &name, &description); err != nil {
			return nil, err
		}
		list = append(list, map[string]interface{}{
			"id":          id,
			"roleCode":    code,
			"roleName":    name,
			"description": description,
		})
	}

	return list, rows.Err()
}

func queryUserPermissions(userID int64) ([]string, error) {
	if !authTableExists("permissions") || !authTableExists("user_roles") || !authTableExists("role_permissions") {
		return []string{}, nil
	}

	rows, err := config.DB.Query(`
		SELECT DISTINCT p.permission_code
		FROM user_roles ur
		JOIN role_permissions rp ON ur.role_id = rp.role_id
		JOIN permissions p ON rp.permission_id = p.id
		WHERE ur.user_id = ?
		ORDER BY p.permission_code
	`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	list := make([]string, 0)
	for rows.Next() {
		var code string
		if err := rows.Scan(&code); err != nil {
			return nil, err
		}
		list = append(list, code)
	}

	return list, rows.Err()
}
