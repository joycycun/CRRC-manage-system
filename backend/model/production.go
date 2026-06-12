package model

import (
	"database/sql"
	"time"
)

type ProductionOrder struct {
	ID              int64        `json:"id"`
	OrderNo         string       `json:"orderNo"`
	ProjectID       int64        `json:"projectId"`
	DeviceType      string       `json:"deviceType"`
	ProductName     string       `json:"productName"`
	ProductModel    string       `json:"productModel"`
	PlanQty         int          `json:"planQty"`
	HardwareID      int64        `json:"hardwareId"`
	HardwareVersion string       `json:"hardwareVersion"`
	SoftwareID      int64        `json:"softwareId"`
	SoftwareVersion string       `json:"softwareVersion"`
	Status          string       `json:"status"`
	CreateUserID    int64        `json:"createUserId"`
	CreateUserName  string       `json:"createUserName"`
	CreateTime      sql.NullTime `json:"createTime"`
	Remark          string       `json:"remark"`
	UpdatedAt       time.Time    `json:"updatedAt"`
	IsDeleted       bool         `json:"isDeleted"`
}

type BurnRecord struct {
	ID                int64        `json:"id"`
	BatchNo           string       `json:"batchNo"`
	ProjectID         int64        `json:"projectId"`
	ProductionOrderID int64        `json:"productionOrderId"`
	ProductName       string       `json:"productName"`
	ProductModel      string       `json:"productModel"`
	ProductCode       string       `json:"productCode"`
	DeviceType        string       `json:"deviceType"`
	SN                string       `json:"sn"`
	MacAddress        string       `json:"macAddress"`
	HardwareID        int64        `json:"hardwareId"`
	HardwareVersion   string       `json:"hardwareVersion"`
	SoftwareID        int64        `json:"softwareId"`
	SoftwareVersion   string       `json:"softwareVersion"`
	PCBQRCode         string       `json:"pcbQrCode"`
	Note              string       `json:"note"`
	SourceFileID      int64        `json:"sourceFileId"`
	UploaderID        int64        `json:"uploaderId"`
	UploaderName      string       `json:"uploaderName"`
	UploadTime        sql.NullTime `json:"uploadTime"`
	BurnDesc          string       `json:"burnDesc"`
	IsDeleted         bool         `json:"isDeleted"`
	CreatedAt         time.Time    `json:"createdAt"`
	UpdatedAt         time.Time    `json:"updatedAt"`
}

type FactoryTest struct {
	ID           int64        `json:"id"`
	BurnRecordID int64        `json:"burnRecordId"`
	ProjectID    int64        `json:"projectId"`
	ProductModel string       `json:"productModel"`
	DeviceType   string       `json:"deviceType"`
	MacAddress   string       `json:"macAddress"`
	SN           string       `json:"sn"`
	FileID       int64        `json:"fileId"`
	UploaderID   int64        `json:"uploaderId"`
	UploaderName string       `json:"uploaderName"`
	UploadTime   sql.NullTime `json:"uploadTime"`
	AuditStatus  string       `json:"auditStatus"`
	RejectReason string       `json:"rejectReason"`
	AuditorID    int64        `json:"auditorId"`
	AuditorName  string       `json:"auditorName"`
	AuditTime    sql.NullTime `json:"auditTime"`
	Remark       string       `json:"remark"`
	IsDeleted    bool         `json:"isDeleted"`
	CreatedAt    time.Time    `json:"createdAt"`
	UpdatedAt    time.Time    `json:"updatedAt"`
}

type InventoryDevice struct {
	ID                 int64        `json:"id"`
	ProjectID          int64        `json:"projectId"`
	DeviceType         string       `json:"deviceType"`
	ProductName        string       `json:"productName"`
	ProductModel       string       `json:"productModel"`
	SN                 string       `json:"sn"`
	MacAddress         string       `json:"macAddress"`
	HardwareID         int64        `json:"hardwareId"`
	HardwareVersion    string       `json:"hardwareVersion"`
	SoftwareID         int64        `json:"softwareId"`
	SoftwareVersion    string       `json:"softwareVersion"`
	InventoryStatus    string       `json:"inventoryStatus"`
	SourceBurnRecordID int64        `json:"sourceBurnRecordId"`
	FactoryTestID      int64        `json:"factoryTestId"`
	InTime             sql.NullTime `json:"inTime"`
	UpdateTime         sql.NullTime `json:"updateTime"`
	Remark             string       `json:"remark"`
	IsDeleted          bool         `json:"isDeleted"`
}
