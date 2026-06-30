package handler

import (
	"crrc_pm_backend/config"
	"crrc_pm_backend/model"
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"crypto/subtle"
	"database/sql"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"
)

const (
	passwordHashPrefix     = "pbkdf2_sha256"
	passwordHashIterations = 120000
	passwordSaltSize       = 16
	passwordKeySize        = 32
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

	passwordOK, needsUpgrade := verifyPassword(user.PasswordHash, req.Password)
	if !passwordOK {
		json.NewEncoder(w).Encode(LoginResponse{
			Code: 401,
			Msg:  "密码错误",
			Data: nil,
		})
		return
	}

	if needsUpgrade {
		if encryptedPassword, err := hashPassword(req.Password); err == nil {
			_, _ = config.DB.Exec(
				"UPDATE users SET password_hash = ?, updated_at = NOW() WHERE id = ?",
				encryptedPassword,
				user.ID,
			)
			user.PasswordHash = encryptedPassword
		}
	}

	_, _ = config.DB.Exec(
		"UPDATE users SET last_login_time = ? WHERE id = ?",
		time.Now(),
		user.ID,
	)

	ensureUserRoleBinding("丁宇", "software_owner")
	ensureUserRoleBinding("王宇", "hardware_owner")

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

type ChangePasswordRequest struct {
	UserID      int64  `json:"userId"`
	Username    string `json:"username"`
	OldPassword string `json:"oldPassword"`
	NewPassword string `json:"newPassword"`
}

func ChangePasswordHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodPost {
		json.NewEncoder(w).Encode(LoginResponse{Code: 405, Msg: "请求方法错误"})
		return
	}

	var req ChangePasswordRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		json.NewEncoder(w).Encode(LoginResponse{Code: 400, Msg: "请求参数错误"})
		return
	}

	if headerUserID := strings.TrimSpace(r.Header.Get("X-User-Id")); headerUserID != "" {
		if id, err := strconv.ParseInt(headerUserID, 10, 64); err == nil {
			req.UserID = id
		}
	}

	req.Username = strings.TrimSpace(req.Username)
	if req.UserID == 0 && req.Username == "" {
		json.NewEncoder(w).Encode(LoginResponse{Code: 400, Msg: "无法识别当前用户"})
		return
	}

	if req.OldPassword == "" {
		json.NewEncoder(w).Encode(LoginResponse{Code: 400, Msg: "请输入原密码"})
		return
	}

	if len(req.NewPassword) < 6 {
		json.NewEncoder(w).Encode(LoginResponse{Code: 400, Msg: "新密码至少需要6位"})
		return
	}

	var user model.User
	var err error
	if req.UserID > 0 {
		err = config.DB.QueryRow(`
			SELECT id, username, password_hash, real_name, status
			FROM users
			WHERE id = ?
			LIMIT 1
		`, req.UserID).Scan(&user.ID, &user.Username, &user.PasswordHash, &user.RealName, &user.Status)
	} else {
		err = config.DB.QueryRow(`
			SELECT id, username, password_hash, real_name, status
			FROM users
			WHERE username = ?
			LIMIT 1
		`, req.Username).Scan(&user.ID, &user.Username, &user.PasswordHash, &user.RealName, &user.Status)
	}

	if err == sql.ErrNoRows {
		json.NewEncoder(w).Encode(LoginResponse{Code: 404, Msg: "账号不存在"})
		return
	}

	if err != nil {
		json.NewEncoder(w).Encode(LoginResponse{Code: 500, Msg: "数据库查询失败: " + err.Error()})
		return
	}

	if user.Status != "启用" {
		json.NewEncoder(w).Encode(LoginResponse{Code: 403, Msg: "账号已禁用"})
		return
	}

	passwordOK, _ := verifyPassword(user.PasswordHash, req.OldPassword)
	if !passwordOK {
		json.NewEncoder(w).Encode(LoginResponse{Code: 401, Msg: "原密码错误"})
		return
	}

	encryptedPassword, err := hashPassword(req.NewPassword)
	if err != nil {
		json.NewEncoder(w).Encode(LoginResponse{Code: 500, Msg: "密码加密失败"})
		return
	}

	_, err = config.DB.Exec(
		"UPDATE users SET password_hash = ?, updated_at = NOW() WHERE id = ?",
		encryptedPassword,
		user.ID,
	)
	if err != nil {
		json.NewEncoder(w).Encode(LoginResponse{Code: 500, Msg: "修改密码失败: " + err.Error()})
		return
	}

	json.NewEncoder(w).Encode(LoginResponse{
		Code: 200,
		Msg:  "密码修改成功，请重新登录",
	})
}

type CreateUserRequest struct {
	Username   string   `json:"username"`
	Password   string   `json:"password"`
	RealName   string   `json:"realName"`
	Email      string   `json:"email"`
	Phone      string   `json:"phone"`
	Department string   `json:"department"`
	Status     string   `json:"status"`
	RoleCodes  []string `json:"roleCodes"`
}

func UsersHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	switch r.Method {
	case http.MethodGet:
		GetUsersHandler(w, r)
		return
	case http.MethodPost:
		CreateUserHandler(w, r)
		return
	default:
		json.NewEncoder(w).Encode(LoginResponse{Code: 405, Msg: "请求方法错误"})
		return
	}
}

func CreateUserHandler(w http.ResponseWriter, r *http.Request) {
	if !hasRequestRole(r, "system_admin") {
		json.NewEncoder(w).Encode(LoginResponse{Code: 403, Msg: "只有系统管理员可以新增用户"})
		return
	}

	if !authTableExists("roles") || !authTableExists("user_roles") {
		json.NewEncoder(w).Encode(LoginResponse{Code: 500, Msg: "角色表未初始化"})
		return
	}
	ensureQualityStaffRole()
	ensureShippingAuditorRole()

	var req CreateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		json.NewEncoder(w).Encode(LoginResponse{Code: 400, Msg: "请求参数错误"})
		return
	}

	req.Username = strings.TrimSpace(req.Username)
	req.RealName = strings.TrimSpace(req.RealName)
	req.Email = strings.TrimSpace(req.Email)
	req.Phone = strings.TrimSpace(req.Phone)
	req.Department = strings.TrimSpace(req.Department)
	req.Status = strings.TrimSpace(req.Status)
	if req.Status == "" {
		req.Status = "启用"
	}
	if req.Password == "" {
		req.Password = "123456"
	}

	if req.Username == "" {
		json.NewEncoder(w).Encode(LoginResponse{Code: 400, Msg: "登录账号不能为空"})
		return
	}
	if req.RealName == "" {
		req.RealName = req.Username
	}
	if len(req.Password) < 6 {
		json.NewEncoder(w).Encode(LoginResponse{Code: 400, Msg: "密码至少需要6位"})
		return
	}
	roleCodes := uniqueRoleCodes(req.RoleCodes)
	if len(roleCodes) == 0 {
		json.NewEncoder(w).Encode(LoginResponse{Code: 400, Msg: "请至少选择一个角色"})
		return
	}

	passwordHash, err := hashPassword(req.Password)
	if err != nil {
		json.NewEncoder(w).Encode(LoginResponse{Code: 500, Msg: "密码加密失败"})
		return
	}

	tx, err := config.DB.Begin()
	if err != nil {
		json.NewEncoder(w).Encode(LoginResponse{Code: 500, Msg: "开启事务失败: " + err.Error()})
		return
	}
	defer tx.Rollback()

	result, err := tx.Exec(`
		INSERT INTO users (
			username,
			password_hash,
			real_name,
			email,
			phone,
			department,
			status,
			created_at,
			updated_at
		) VALUES (?, ?, ?, ?, ?, ?, ?, NOW(), NOW())
	`,
		req.Username,
		passwordHash,
		req.RealName,
		req.Email,
		req.Phone,
		req.Department,
		req.Status,
	)
	if err != nil {
		if strings.Contains(strings.ToLower(err.Error()), "duplicate") {
			json.NewEncoder(w).Encode(LoginResponse{Code: 409, Msg: "登录账号已存在"})
			return
		}
		json.NewEncoder(w).Encode(LoginResponse{Code: 500, Msg: "新增用户失败: " + err.Error()})
		return
	}

	userID, _ := result.LastInsertId()
	for _, roleCode := range roleCodes {
		roleResult, err := tx.Exec(`
			INSERT IGNORE INTO user_roles (user_id, role_id, created_at)
			SELECT ?, id, NOW()
			FROM roles
			WHERE role_code = ?
		`, userID, roleCode)
		if err != nil {
			json.NewEncoder(w).Encode(LoginResponse{Code: 500, Msg: "绑定角色失败: " + err.Error()})
			return
		}
		affected, _ := roleResult.RowsAffected()
		if affected == 0 {
			var exists int
			_ = tx.QueryRow("SELECT COUNT(1) FROM roles WHERE role_code = ?", roleCode).Scan(&exists)
			if exists == 0 {
				json.NewEncoder(w).Encode(LoginResponse{Code: 400, Msg: "角色不存在: " + roleCode})
				return
			}
		}
	}

	if err := tx.Commit(); err != nil {
		json.NewEncoder(w).Encode(LoginResponse{Code: 500, Msg: "提交事务失败: " + err.Error()})
		return
	}

	json.NewEncoder(w).Encode(LoginResponse{
		Code: 200,
		Msg:  "新增用户成功",
		Data: map[string]interface{}{
			"id": userID,
		},
	})
}

func GetUsersHandler(w http.ResponseWriter, r *http.Request) {
	if !hasRequestRole(r, "system_admin") {
		json.NewEncoder(w).Encode(LoginResponse{Code: 403, Msg: "只有系统管理员可以查看用户"})
		return
	}

	rows, err := config.DB.Query(`
		SELECT
			u.id,
			u.username,
			IFNULL(u.real_name, ''),
			IFNULL(u.department, ''),
			IFNULL(u.status, ''),
			IFNULL(GROUP_CONCAT(r.role_name ORDER BY r.id SEPARATOR '、'), '')
		FROM users u
		LEFT JOIN user_roles ur ON ur.user_id = u.id
		LEFT JOIN roles r ON r.id = ur.role_id
		GROUP BY u.id, u.username, u.real_name, u.department, u.status
		ORDER BY u.id ASC
	`)
	if err != nil {
		json.NewEncoder(w).Encode(LoginResponse{Code: 500, Msg: "查询用户失败: " + err.Error()})
		return
	}
	defer rows.Close()

	list := make([]map[string]interface{}, 0)
	for rows.Next() {
		var id int64
		var username string
		var realName string
		var department string
		var status string
		var roles string
		if err := rows.Scan(&id, &username, &realName, &department, &status, &roles); err != nil {
			json.NewEncoder(w).Encode(LoginResponse{Code: 500, Msg: "解析用户失败: " + err.Error()})
			return
		}
		list = append(list, map[string]interface{}{
			"id":         id,
			"username":   username,
			"realName":   realName,
			"department": department,
			"status":     status,
			"roles":      roles,
		})
	}

	json.NewEncoder(w).Encode(LoginResponse{Code: 200, Msg: "查询成功", Data: list})
}

func UserActionHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	path := strings.TrimPrefix(r.URL.Path, "/api/users/")
	path = strings.Trim(path, "/")
	if path == "" {
		json.NewEncoder(w).Encode(LoginResponse{Code: 400, Msg: "缺少用户ID"})
		return
	}

	id, err := strconv.ParseInt(path, 10, 64)
	if err != nil || id <= 0 {
		json.NewEncoder(w).Encode(LoginResponse{Code: 400, Msg: "用户ID错误"})
		return
	}

	if r.Method != http.MethodDelete {
		json.NewEncoder(w).Encode(LoginResponse{Code: 405, Msg: "请求方法错误"})
		return
	}

	DeleteUserHandler(w, r, id)
}

func DeleteUserHandler(w http.ResponseWriter, r *http.Request, id int64) {
	if !hasRequestRole(r, "system_admin") {
		json.NewEncoder(w).Encode(LoginResponse{Code: 403, Msg: "只有系统管理员可以删除用户"})
		return
	}

	var username string
	if err := config.DB.QueryRow("SELECT username FROM users WHERE id = ?", id).Scan(&username); err != nil {
		json.NewEncoder(w).Encode(LoginResponse{Code: 404, Msg: "用户不存在"})
		return
	}
	if username == "admin" {
		json.NewEncoder(w).Encode(LoginResponse{Code: 400, Msg: "admin 账户不能删除"})
		return
	}

	tx, err := config.DB.Begin()
	if err != nil {
		json.NewEncoder(w).Encode(LoginResponse{Code: 500, Msg: "开启事务失败: " + err.Error()})
		return
	}
	defer tx.Rollback()

	if _, err := tx.Exec("DELETE FROM user_roles WHERE user_id = ?", id); err != nil {
		json.NewEncoder(w).Encode(LoginResponse{Code: 500, Msg: "删除用户角色失败: " + err.Error()})
		return
	}
	result, err := tx.Exec("DELETE FROM users WHERE id = ?", id)
	if err != nil {
		json.NewEncoder(w).Encode(LoginResponse{Code: 500, Msg: "删除用户失败: " + err.Error()})
		return
	}
	affected, _ := result.RowsAffected()
	if affected == 0 {
		json.NewEncoder(w).Encode(LoginResponse{Code: 404, Msg: "用户不存在"})
		return
	}
	if err := tx.Commit(); err != nil {
		json.NewEncoder(w).Encode(LoginResponse{Code: 500, Msg: "提交事务失败: " + err.Error()})
		return
	}

	json.NewEncoder(w).Encode(LoginResponse{Code: 200, Msg: "删除用户成功"})
}

func uniqueRoleCodes(values []string) []string {
	seen := make(map[string]bool)
	result := make([]string, 0, len(values))
	for _, value := range values {
		value = strings.TrimSpace(value)
		if value == "" || seen[value] {
			continue
		}
		seen[value] = true
		result = append(result, value)
	}
	return result
}

func UserRolesOptionsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodGet {
		json.NewEncoder(w).Encode(LoginResponse{Code: 405, Msg: "请求方法错误"})
		return
	}

	if !hasRequestRole(r, "system_admin") {
		json.NewEncoder(w).Encode(LoginResponse{Code: 403, Msg: "只有系统管理员可以查看角色"})
		return
	}

	ensureQualityStaffRole()
	ensureShippingAuditorRole()

	rows, err := config.DB.Query(`
		SELECT role_code, role_name, IFNULL(description, '')
		FROM roles
		ORDER BY id ASC
	`)
	if err != nil {
		json.NewEncoder(w).Encode(LoginResponse{Code: 500, Msg: "查询角色失败: " + err.Error()})
		return
	}
	defer rows.Close()

	list := make([]map[string]interface{}, 0)
	for rows.Next() {
		var roleCode string
		var roleName string
		var description string
		if err := rows.Scan(&roleCode, &roleName, &description); err != nil {
			json.NewEncoder(w).Encode(LoginResponse{Code: 500, Msg: "解析角色失败: " + err.Error()})
			return
		}
		list = append(list, map[string]interface{}{
			"roleCode":    roleCode,
			"roleName":    roleName,
			"description": description,
		})
	}

	json.NewEncoder(w).Encode(LoginResponse{
		Code: 200,
		Msg:  "查询成功",
		Data: list,
	})
}

func hashPassword(password string) (string, error) {
	salt := make([]byte, passwordSaltSize)
	if _, err := io.ReadFull(rand.Reader, salt); err != nil {
		return "", err
	}

	key := pbkdf2Key([]byte(password), salt, passwordHashIterations, passwordKeySize)
	return fmt.Sprintf(
		"%s$%d$%s$%s",
		passwordHashPrefix,
		passwordHashIterations,
		base64.RawStdEncoding.EncodeToString(salt),
		base64.RawStdEncoding.EncodeToString(key),
	), nil
}

func verifyPassword(storedPassword, plainPassword string) (bool, bool) {
	if strings.HasPrefix(storedPassword, passwordHashPrefix+"$") {
		return verifyHashedPassword(storedPassword, plainPassword), false
	}

	return subtle.ConstantTimeCompare([]byte(storedPassword), []byte(plainPassword)) == 1, true
}

func verifyHashedPassword(storedPassword, plainPassword string) bool {
	parts := strings.Split(storedPassword, "$")
	if len(parts) != 4 || parts[0] != passwordHashPrefix {
		return false
	}

	iterations, err := strconv.Atoi(parts[1])
	if err != nil || iterations <= 0 {
		return false
	}

	salt, err := base64.RawStdEncoding.DecodeString(parts[2])
	if err != nil {
		return false
	}

	expected, err := base64.RawStdEncoding.DecodeString(parts[3])
	if err != nil || len(expected) == 0 {
		return false
	}

	actual := pbkdf2Key([]byte(plainPassword), salt, iterations, len(expected))
	return subtle.ConstantTimeCompare(actual, expected) == 1
}

func pbkdf2Key(password, salt []byte, iterations, keyLen int) []byte {
	hashLen := sha256.Size
	numBlocks := (keyLen + hashLen - 1) / hashLen
	result := make([]byte, 0, numBlocks*hashLen)

	for block := 1; block <= numBlocks; block++ {
		mac := hmac.New(sha256.New, password)
		mac.Write(salt)
		mac.Write([]byte{
			byte(block >> 24),
			byte(block >> 16),
			byte(block >> 8),
			byte(block),
		})

		u := mac.Sum(nil)
		t := make([]byte, len(u))
		copy(t, u)

		for i := 1; i < iterations; i++ {
			mac = hmac.New(sha256.New, password)
			mac.Write(u)
			u = mac.Sum(nil)

			for j := range t {
				t[j] ^= u[j]
			}
		}

		result = append(result, t...)
	}

	return result[:keyLen]
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

	ensureUserRoleBinding("丁宇", "software_owner")

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

func ensureQualityStaffRole() {
	_, _ = config.DB.Exec(`
		INSERT INTO roles (role_code, role_name, description, created_at, updated_at)
		VALUES ('quality_staff', '质量检查人员', '负责生产管理里面生产测试审查', NOW(), NOW())
		ON DUPLICATE KEY UPDATE
			role_name = VALUES(role_name),
			description = VALUES(description),
			updated_at = NOW()
	`)

	_, _ = config.DB.Exec(`
		INSERT INTO permissions (permission_code, permission_name, module, description, created_at)
		VALUES
			('production:view', '查看生产记录', '生产管理', '查看生产数据', NOW()),
			('production:test', '出厂测试', '生产管理', '处理烧录后的出厂测试', NOW()),
			('production:audit', '审核生产测试', '生产管理', '质量检查人员审核出厂测试', NOW())
		ON DUPLICATE KEY UPDATE
			permission_name = VALUES(permission_name),
			module = VALUES(module),
			description = VALUES(description)
	`)

	_, _ = config.DB.Exec(`
		INSERT IGNORE INTO role_permissions (role_id, permission_id, created_at)
		SELECT r.id, p.id, NOW()
		FROM roles r
		JOIN permissions p
		WHERE r.role_code = 'quality_staff'
		  AND p.permission_code IN ('project:view', 'production:view', 'production:test', 'production:audit', 'report:view', 'report:manage')
	`)
}

func ensureShippingAuditorRole() {
	_, _ = config.DB.Exec(`
		INSERT INTO roles (role_code, role_name, description, created_at, updated_at)
		VALUES ('shipping_auditor', '发货审核', '发货人员权限基础上增加发货批次审核', NOW(), NOW())
		ON DUPLICATE KEY UPDATE
			role_name = VALUES(role_name),
			description = VALUES(description),
			updated_at = NOW()
	`)

	_, _ = config.DB.Exec(`
		INSERT INTO permissions (permission_code, permission_name, module, description, created_at)
		VALUES
			('project:view', '查看项目', '项目管理', '查看项目立项数据', NOW()),
			('shipping:view', '查看发货数据', '发货管理', '查看发货和出库数据', NOW()),
			('shipping:audit', '审核发货批次', '发货管理', '审核发货批次', NOW()),
			('report:view', '查看统计报表', '统计报表', '查看项目进度、版本矩阵和问题统计', NOW()),
			('report:manage', '查看统计报表', '统计报表', '查看统计报表', NOW())
		ON DUPLICATE KEY UPDATE
			permission_name = VALUES(permission_name),
			module = VALUES(module),
			description = VALUES(description)
	`)

	_, _ = config.DB.Exec(`
		INSERT IGNORE INTO role_permissions (role_id, permission_id, created_at)
		SELECT r.id, p.id, NOW()
		FROM roles r
		JOIN permissions p
		WHERE r.role_code = 'shipping_auditor'
		  AND p.permission_code IN ('project:view', 'shipping:view', 'shipping:audit', 'report:view', 'report:manage')
	`)

	_, _ = config.DB.Exec(`
		DELETE rp
		FROM role_permissions rp
		JOIN roles r ON r.id = rp.role_id
		JOIN permissions p ON p.id = rp.permission_id
		WHERE r.role_code = 'shipping_auditor'
		  AND p.permission_code = 'shipping:manage'
	`)
}

func ensureUserRoleBinding(username string, roleCode string) {
	username = strings.TrimSpace(username)
	roleCode = strings.TrimSpace(roleCode)
	if username == "" || roleCode == "" {
		return
	}

	_, _ = config.DB.Exec(`
		INSERT IGNORE INTO user_roles (user_id, role_id, created_at)
		SELECT u.id, r.id, NOW()
		FROM users u
		JOIN roles r ON r.role_code = ?
		WHERE (u.username = ? OR u.real_name = ?)
		  AND IFNULL(u.status, '启用') = '启用'
	`, roleCode, username, username)
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
