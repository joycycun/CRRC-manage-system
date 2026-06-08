package model

import (
	"database/sql"
	"time"
)

type Issue struct {
	ID            int64        `json:"id"`
	ProjectID     int64        `json:"projectId"`
	DeviceType    string       `json:"deviceType"`
	IssueSource   string       `json:"issueSource"`
	Level         string       `json:"level"`
	IssueTitle    string       `json:"issueTitle"`
	IssueDesc     string       `json:"issueDesc"`
	OwnerID       int64        `json:"ownerId"`
	OwnerName     string       `json:"ownerName"`
	CreatorID     int64        `json:"creatorId"`
	CreatorName   string       `json:"creatorName"`
	CreateTime    sql.NullTime `json:"createTime"`
	PlanCloseTime sql.NullTime `json:"planCloseTime"`
	RealCloseTime sql.NullTime `json:"realCloseTime"`
	CloseStatus   string       `json:"closeStatus"`
	CloseUserID   int64        `json:"closeUserId"`
	CloseUserName string       `json:"closeUserName"`
	ReopenReason  string       `json:"reopenReason"`
	UpdatedAt     time.Time    `json:"updatedAt"`
	IsDeleted     bool         `json:"isDeleted"`
}

type IssueReply struct {
	ID            int64        `json:"id"`
	IssueID       int64        `json:"issueId"`
	ReplyUserID   int64        `json:"replyUserId"`
	ReplyUserName string       `json:"replyUserName"`
	ReplyTime     sql.NullTime `json:"replyTime"`
	Content       string       `json:"content"`
}
