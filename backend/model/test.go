package model

type TestCase struct {
	ID          int64  `json:"id"`
	ProjectID   int64  `json:"projectId"`
	ProjectName string `json:"projectName"`

	CaseName string `json:"caseName"`

	FileID   int64  `json:"fileId"`
	FileName string `json:"fileName"`
	FileURL  string `json:"fileUrl"`

	UploaderID   int64  `json:"uploaderId"`
	UploaderName string `json:"uploaderName"`
	Uploader     string `json:"uploader"`
	UploadTime   string `json:"uploadTime"`

	AuditStatus string `json:"auditStatus"`
	AuditorID   int64  `json:"auditorId"`
	AuditorName string `json:"auditorName"`
	Auditor     string `json:"auditor"`
	AuditTime   string `json:"auditTime"`

	RejectReason string `json:"rejectReason"`
	Remark       string `json:"remark"`

	ReportName         string `json:"reportName"`
	ReportFileID       int64  `json:"reportFileId"`
	ReportFileName     string `json:"reportFileName"`
	ReportFileURL      string `json:"reportFileUrl"`
	ReportUploaderID   int64  `json:"reportUploaderId"`
	ReportUploaderName string `json:"reportUploaderName"`
	ReportUploader     string `json:"reportUploader"`
	ReportUploadTime   string `json:"reportUploadTime"`
	ReportRemark       string `json:"reportRemark"`

	CreatedAt string `json:"createdAt"`
	UpdatedAt string `json:"updatedAt"`
}
