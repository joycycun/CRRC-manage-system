package model

import (
	"database/sql"
	"time"
)

type RequirementChange struct {
	ID             int64        `json:"id"`
	ProjectID      int64        `json:"projectId"`
	ChangeTitle    string       `json:"changeTitle"`
	ChangeType     string       `json:"changeType"`
	FileID         int64        `json:"fileId"`
	Status         string       `json:"status"`
	CloseStatus    string       `json:"closeStatus"`
	SubmitUserID   int64        `json:"submitUserId"`
	SubmitUserName string       `json:"submitUserName"`
	AuditUserID    int64        `json:"auditUserId"`
	AuditUserName  string       `json:"auditUserName"`
	SubmitTime     sql.NullTime `json:"submitTime"`
	AuditTime      sql.NullTime `json:"auditTime"`
	CloseUserID    int64        `json:"closeUserId"`
	CloseUserName  string       `json:"closeUserName"`
	CloseTime      sql.NullTime `json:"closeTime"`
	RejectReason   string       `json:"rejectReason"`
	Remark         string       `json:"remark"`
	CreatedAt      time.Time    `json:"createdAt"`
	UpdatedAt      time.Time    `json:"updatedAt"`
	IsDeleted      bool         `json:"isDeleted"`
}
