package model

import (
	"database/sql"
	"time"
)

type CustomerSuppliedFile struct {
	ID              int64        `json:"id"`
	ProjectID       int64        `json:"projectId"`
	FileID          int64        `json:"fileId"`
	MaterialName    string       `json:"materialName"`
	FileDisplayName string       `json:"fileDisplayName"`
	MaterialDesc    string       `json:"materialDesc"`
	UploadUserID    int64        `json:"uploadUserId"`
	UploadUserName  string       `json:"uploadUserName"`
	UploadTime      sql.NullTime `json:"uploadTime"`
	Remark          string       `json:"remark"`
	CreatedAt       time.Time    `json:"createdAt"`
	UpdatedAt       time.Time    `json:"updatedAt"`
	IsDeleted       bool         `json:"isDeleted"`
}
