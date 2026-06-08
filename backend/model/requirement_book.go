package model

import (
	"database/sql"
	"time"
)

type RequirementBook struct {
	ID             int64        `json:"id"`
	ProjectID      int64        `json:"projectId"`
	BookName       string       `json:"bookName"`
	FileID         int64        `json:"fileId"`
	Status         string       `json:"status"`
	SubmitUserID   int64        `json:"submitUserId"`
	SubmitUserName string       `json:"submitUserName"`
	SubmitTime     sql.NullTime `json:"submitTime"`
	AuditUserID    int64        `json:"auditUserId"`
	AuditUserName  string       `json:"auditUserName"`
	AuditTime      sql.NullTime `json:"auditTime"`
	RejectReason   string       `json:"rejectReason"`
	Remark         string       `json:"remark"`
	CreatedAt      time.Time    `json:"createdAt"`
	UpdatedAt      time.Time    `json:"updatedAt"`
	IsDeleted      bool         `json:"isDeleted"`
}
