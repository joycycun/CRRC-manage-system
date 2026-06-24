package model

import (
	"database/sql"
	"time"
)

type HardwareVersion struct {
	ID                   int64     `json:"id"`
	HardwareVersion      string    `json:"hardwareVersion"`
	ProjectID            int64     `json:"projectId"`
	DeviceType           string    `json:"deviceType"`
	Status               string    `json:"status"`
	OwnerID              int64     `json:"ownerId"`
	OwnerName            string    `json:"ownerName"`
	ChangeDocFileID      int64     `json:"changeDocFileId"`
	ChangeDocFileName    string    `json:"changeDocFileName"`
	ChangeDocFileURL     string    `json:"changeDocFileUrl"`
	ChangeDocDownloadURL string    `json:"changeDocDownloadUrl"`
	Description          string    `json:"description"`
	CreatedAt            time.Time `json:"createdAt"`
	UpdatedAt            time.Time `json:"updatedAt"`
}

type HardwareTest struct {
	ID              int64        `json:"id"`
	ProjectID       int64        `json:"projectId"`
	ProjectName     string       `json:"projectName"`
	HardwareVersion string       `json:"hardwareVersion"`
	HardwareID      int64        `json:"hardwareId"`
	RecordName      string       `json:"recordName"`
	DeviceType      string       `json:"deviceType"`
	FileID          int64        `json:"fileId"`
	FileName        string       `json:"fileName"`
	FileURL         string       `json:"fileUrl"`
	DownloadURL     string       `json:"downloadUrl"`
	AuditStatus     string       `json:"auditStatus"`
	AuditorID       int64        `json:"auditorId"`
	AuditorName     string       `json:"auditorName"`
	AuditTime       sql.NullTime `json:"auditTime"`
	RejectReason    string       `json:"rejectReason"`
	Remark          string       `json:"remark"`
	CreatedAt       time.Time    `json:"createdAt"`
	UpdatedAt       time.Time    `json:"updatedAt"`
	IsDeleted       bool         `json:"isDeleted"`
	UploaderID      int64        `json:"uploaderId"`
	UploaderName    string       `json:"uploaderName"`
	UploadTime      sql.NullTime `json:"uploadTime"`
}
