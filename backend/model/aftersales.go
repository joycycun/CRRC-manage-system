package model

import (
	"database/sql"
	"time"
)

type RepairRecord struct {
	ID                int64        `json:"id"`
	ProjectID         int64        `json:"projectId"`
	DeviceType        string       `json:"deviceType"`
	InventoryDeviceID int64        `json:"inventoryDeviceId"`
	SN                string       `json:"sn"`
	MacAddress        string       `json:"macAddress"`
	FaultDesc         string       `json:"faultDesc"`
	RepairUserID      int64        `json:"repairUserId"`
	RepairUserName    string       `json:"repairUserName"`
	RepairTime        sql.NullTime `json:"repairTime"`
	RepairFinishTime  sql.NullTime `json:"repairFinishTime"`
	RepairMethod      string       `json:"repairMethod"`
	RepairProcess     string       `json:"repairProcess"`
	ConfirmStatus     string       `json:"confirmStatus"`
	Remark            string       `json:"remark"`
	IsDeleted         bool         `json:"isDeleted"`
	CreatedAt         time.Time    `json:"createdAt"`
	UpdatedAt         time.Time    `json:"updatedAt"`
}

type FaultAnalysis struct {
	ID             int64        `json:"id"`
	ProjectID      int64        `json:"projectId"`
	ProjectName    string       `json:"projectName"`
	IssueID        int64        `json:"issueId"`
	RepairID       int64        `json:"repairId"`
	BoardType      string       `json:"boardType"`
	AnalysisName   string       `json:"analysisName"`
	FileID         int64        `json:"fileId"`
	FileName       string       `json:"fileName"`
	FileURL        string       `json:"fileUrl"`
	SubmitUserID   int64        `json:"submitUserId"`
	SubmitUserName string       `json:"submitUserName"`
	SubmitTime     sql.NullTime `json:"submitTime"`
	AuditStatus    string       `json:"auditStatus"`
	AuditorID      int64        `json:"auditorId"`
	AuditorName    string       `json:"auditorName"`
	AuditTime      sql.NullTime `json:"auditTime"`
	RejectReason   string       `json:"rejectReason"`
	AnalysisDesc   string       `json:"analysisDesc"`
	IsDeleted      bool         `json:"isDeleted"`
	CreatedAt      time.Time    `json:"createdAt"`
	UpdatedAt      time.Time    `json:"updatedAt"`
}
