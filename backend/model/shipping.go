package model

import (
	"database/sql"
	"time"
)

type ShippingBatch struct {
	ID           int64        `json:"id"`
	BatchNo      string       `json:"batchNo"`
	ProjectID    int64        `json:"projectId"`
	ExpressNo    string       `json:"expressNo"`
	DeviceCount  int          `json:"deviceCount"`
	FileID       int64        `json:"fileId"`
	UploaderID   int64        `json:"uploaderId"`
	UploaderName string       `json:"uploaderName"`
	UploadTime   sql.NullTime `json:"uploadTime"`
	AuditStatus  string       `json:"auditStatus"`
	AuditorID    int64        `json:"auditorId"`
	AuditorName  string       `json:"auditorName"`
	AuditTime    sql.NullTime `json:"auditTime"`
	Remark       string       `json:"remark"`
	RejectReason string       `json:"rejectReason"`
	ShippingDesc string       `json:"shippingDesc"`
	IsDeleted    bool         `json:"isDeleted"`
	CreatedAt    time.Time    `json:"createdAt"`
	UpdatedAt    time.Time    `json:"updatedAt"`
}

type ShippingBatchDevice struct {
	ID                int64     `json:"id"`
	BatchID           int64     `json:"batchId"`
	InventoryDeviceID int64     `json:"inventoryDeviceId"`
	SN                string    `json:"sn"`
	MacAddress        string    `json:"macAddress"`
	DeviceType        string    `json:"deviceType"`
	HardwareVersion   string    `json:"hardwareVersion"`
	SoftwareVersion   string    `json:"softwareVersion"`
	IsDeleted         bool      `json:"isDeleted"`
	CreatedAt         time.Time `json:"createdAt"`
}

type OutboundRecord struct {
	ID                int64        `json:"id"`
	BatchID           int64        `json:"batchId"`
	InventoryDeviceID int64        `json:"inventoryDeviceId"`
	SN                string       `json:"sn"`
	MacAddress        string       `json:"macAddress"`
	OutboundTime      sql.NullTime `json:"outboundTime"`
	OutboundUserID    int64        `json:"outboundUserId"`
	OutboundUserName  string       `json:"outboundUserName"`
	Status            string       `json:"status"`
	Remark            string       `json:"remark"`
	InboundTime       sql.NullTime `json:"inboundTime"`
	IsDeleted         bool         `json:"isDeleted"`
	CreatedAt         time.Time    `json:"createdAt"`
}