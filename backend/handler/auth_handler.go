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

	json.NewEncoder(w).Encode(LoginResponse{
		Code: 200,
		Msg:  "登录成功",
		Data: map[string]interface{}{
			"token": "test-token",
			"user":  user,
			"permissions": []string{
				"project:view",
				"project:create",
				"project:audit",
			},
		},
	})
}
