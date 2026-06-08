package model

import (
	"database/sql"
	"time"
)

type TestCase struct {
	ID                 int64        `json:"id"`
	ProjectID          int64        `json:"projectId"`
	CaseName           string       `json:"caseName"`
	CaseFileID         int64        `json:"caseFileId"`
	ReportFileID       int64        `json:"reportFileId"`
	UploaderID         int64        `json:"uploaderId"`
	UploaderName       string       `json:"uploaderName"`
	UploadTime         sql.NullTime `json:"uploadTime"`
	ReportUploaderID   int64        `json:"reportUploaderId"`
	ReportUploaderName string       `json:"reportUploaderName"`
	ReportUploadTime   sql.NullTime `json:"reportUploadTime"`
	AuditStatus        string       `json:"auditStatus"`
	AuditorID          int64        `json:"auditorId"`
	AuditorName        string       `json:"auditorName"`
	AuditTime          sql.NullTime `json:"auditTime"`
	RejectReason       string       `json:"rejectReason"`
	Remark             string       `json:"remark"`
	CreatedAt          time.Time    `json:"createdAt"`
	UpdatedAt          time.Time    `json:"updatedAt"`
	IsDeleted          bool         `json:"isDeleted"`
}

type TestReport struct {
	ID           int64        `json:"id"`
	ProjectID    int64        `json:"projectId"`
	TestCaseID   int64        `json:"testCaseId"`
	ReportName   string       `json:"reportName"`
	ReportFileID int64        `json:"reportFileId"`
	UploaderID   int64        `json:"uploaderId"`
	UploaderName string       `json:"uploaderName"`
	UploadTime   sql.NullTime `json:"uploadTime"`
	AuditStatus  string       `json:"auditStatus"`
	AuditorID    int64        `json:"auditorId"`
	AuditorName  string       `json:"auditorName"`
	AuditTime    sql.NullTime `json:"auditTime"`
	RejectReason string       `json:"rejectReason"`
	Remark       string       `json:"remark"`
	CreatedAt    time.Time    `json:"createdAt"`
	UpdatedAt    time.Time    `json:"updatedAt"`
	IsDeleted    bool         `json:"isDeleted"`
}
