package model

import (
	"database/sql"
	"time"
)

type Project struct {
	ID            int64        `json:"id"`
	ProjectName   string       `json:"projectName"`
	ProjectCode   string       `json:"projectCode"`
	OwnerID       int64        `json:"ownerId"`
	OwnerName     string       `json:"owner"`
	Stage         string       `json:"stage"`
	Status        string       `json:"status"`
	SubmitTime    sql.NullTime `json:"-"`
	AuditStatus   string       `json:"auditStatus"`
	AuditUserID   int64        `json:"auditUserId"`
	AuditUserName string       `json:"auditUserName"`
	AuditTime     sql.NullTime `json:"-"`
	ArchiveTime   sql.NullTime `json:"-"`
	Remark        string       `json:"remark"`
	CreatedAt     time.Time    `json:"createdAt"`
	UpdatedAt     time.Time    `json:"updatedAt"`
	CloseTime     sql.NullTime `json:"-"`
	IsDeleted     bool         `json:"isDeleted"`
}

type ProjectResponse struct {
	ID            int64      `json:"id"`
	ProjectName   string     `json:"projectName"`
	ProjectCode   string     `json:"projectCode"`
	OwnerID       int64      `json:"ownerId"`
	OwnerName     string     `json:"owner"`
	Stage         string     `json:"stage"`
	Status        string     `json:"status"`
	SubmitTime    *time.Time `json:"submitTime"`
	AuditStatus   string     `json:"auditStatus"`
	AuditUserID   int64      `json:"auditUserId"`
	AuditUserName string     `json:"auditUserName"`
	AuditTime     *time.Time `json:"auditTime"`
	ArchiveTime   *time.Time `json:"archiveTime"`
	Remark        string     `json:"remark"`
	CreatedAt     time.Time  `json:"createdAt"`
	UpdatedAt     time.Time  `json:"updatedAt"`
	CloseTime     *time.Time `json:"closeTime"`
	IsDeleted     bool       `json:"isDeleted"`
}

func ProjectNullTimeToPtr(t sql.NullTime) *time.Time {
	if t.Valid {
		return &t.Time
	}
	return nil
}

func (p Project) ToResponse() ProjectResponse {
	return ProjectResponse{
		ID:            p.ID,
		ProjectName:   p.ProjectName,
		ProjectCode:   p.ProjectCode,
		OwnerID:       p.OwnerID,
		OwnerName:     p.OwnerName,
		Stage:         p.Stage,
		Status:        p.Status,
		SubmitTime:    ProjectNullTimeToPtr(p.SubmitTime),
		AuditStatus:   p.AuditStatus,
		AuditUserID:   p.AuditUserID,
		AuditUserName: p.AuditUserName,
		AuditTime:     ProjectNullTimeToPtr(p.AuditTime),
		ArchiveTime:   ProjectNullTimeToPtr(p.ArchiveTime),
		Remark:        p.Remark,
		CreatedAt:     p.CreatedAt,
		UpdatedAt:     p.UpdatedAt,
		CloseTime:     ProjectNullTimeToPtr(p.CloseTime),
		IsDeleted:     p.IsDeleted,
	}
}
