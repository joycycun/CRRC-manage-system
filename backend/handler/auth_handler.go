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
