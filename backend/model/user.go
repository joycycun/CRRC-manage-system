package model

import "time"

type User struct {
	ID            int64      `json:"id"`
	Username      string     `json:"username"`
	PasswordHash  string     `json:"-"`
	RealName      string     `json:"realName"`
	Email         string     `json:"email"`
	Phone         string     `json:"phone"`
	Department    string     `json:"department"`
	Status        string     `json:"status"`
	LastLoginTime *time.Time `json:"lastLoginTime"`
	CreatedAt     time.Time  `json:"createdAt"`
	UpdatedAt     time.Time  `json:"updatedAt"`
	AvatarURL     string     `json:"avatarUrl"`
}