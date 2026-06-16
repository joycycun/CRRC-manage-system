package model

type ProjectProgressReport struct {
	ProjectID       int64   `json:"projectId"`
	ProjectName     string  `json:"projectName"`
	ProjectCode     string  `json:"projectCode"`
	Stage           string  `json:"stage"`
	Status          string  `json:"status"`
	RequirementDone bool    `json:"requirementDone"`
	HardwareDone    bool    `json:"hardwareDone"`
	SoftwareDone    bool    `json:"softwareDone"`
	TestDone        bool    `json:"testDone"`
	IssueDone       bool    `json:"issueDone"`
	ProductionDone  bool    `json:"productionDone"`
	ShippingDone    bool    `json:"shippingDone"`
	Progress        float64 `json:"progress"`
}

type VersionMatrixReport struct {
	ProjectID       int64  `json:"projectId"`
	ProjectName     string `json:"projectName"`
	DeviceType      string `json:"deviceType"`
	HardwareVersion string `json:"hardwareVersion"`
	SoftwareVersion string `json:"softwareVersion"`
	SN              string `json:"sn"`
	MacAddress      string `json:"macAddress"`
	InventoryStatus string `json:"inventoryStatus"`
	ShippingBatchNo string `json:"shippingBatchNo"`
	OutboundStatus  string `json:"outboundStatus"`
}

type IssueStatisticsReport struct {
	ProjectID   int64  `json:"projectId"`
	ProjectName string `json:"projectName"`
	IssueSource string `json:"issueSource"`
	Level       string `json:"level"`
	CloseStatus string `json:"closeStatus"`
	OwnerName   string `json:"ownerName"`
	IssueCount  int64  `json:"issueCount"`
}

type DashboardSummary struct {
	OngoingProjects      int64                  `json:"ongoingProjects"`
	OngoingProjectDelta  int64                  `json:"ongoingProjectDelta"`
	OngoingProjectMonth  int64                  `json:"ongoingProjectMonth"`
	OngoingProjectPrev   int64                  `json:"ongoingProjectPrev"`
	PendingIssues        int64                  `json:"pendingIssues"`
	CurrentMonthReleases int64                  `json:"currentMonthReleases"`
	ReleaseDelta         int64                  `json:"releaseDelta"`
	ReleasePrevMonth     int64                  `json:"releasePrevMonth"`
	InventoryInStock     int64                  `json:"inventoryInStock"`
	InventoryOutbound    int64                  `json:"inventoryOutbound"`
	RepairPending        int64                  `json:"repairPending"`
	Kpis                 []DashboardKpiItem     `json:"kpis"`
	Todos                []DashboardTodoItem    `json:"todos"`
	Notifications        []DashboardTodoItem    `json:"notifications"`
	RecentReleases       []DashboardReleaseItem `json:"recentReleases"`
}

type DashboardKpiItem struct {
	Title string `json:"title"`
	Value int64  `json:"value"`
	Desc  string `json:"desc"`
	Icon  string `json:"icon"`
}

type DashboardTodoItem struct {
	ID       string `json:"id"`
	Title    string `json:"title"`
	Deadline string `json:"deadline"`
	Level    string `json:"level"`
	Type     string `json:"type"`
	Link     string `json:"link"`
}

type DashboardReleaseItem struct {
	ID          int64  `json:"id"`
	Version     string `json:"version"`
	ReleaseTime string `json:"releaseTime"`
	Status      string `json:"status"`
	StatusType  string `json:"statusType"`
}
