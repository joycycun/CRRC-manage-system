package router

import (
	"crrc_pm_backend/handler"
	"crrc_pm_backend/middleware"
	"net/http"
)

func InitRouter() http.Handler {
	mux := http.NewServeMux()

	// 登录
	mux.HandleFunc("/api/login", handler.LoginHandler)
	mux.HandleFunc("/api/users/software-owners", handler.SoftwareOwnersHandler)

	// 项目管理
	mux.HandleFunc("/api/projects", handler.ProjectsHandler)
	mux.HandleFunc("/api/projects/", handler.ProjectActionHandler)

	// 需求书管理
	mux.HandleFunc("/api/requirement-books", handler.RequirementBooksHandler)
	mux.HandleFunc("/api/requirement-books/", handler.RequirementBookActionHandler)
	//需求书更改
	mux.HandleFunc("/api/requirement-changes", handler.RequirementChangesHandler)
	mux.HandleFunc("/api/requirement-changes/", handler.RequirementChangeActionHandler)
	// 客供资料管理
	mux.HandleFunc("/api/customer-supplied-files", handler.CustomerSuppliedFilesHandler)
	mux.HandleFunc("/api/customer-supplied-files/", handler.CustomerSuppliedFileActionHandler)
	//硬件版本和测试
	mux.HandleFunc("/api/hardware-versions", handler.HardwareVersionsHandler)
	mux.HandleFunc("/api/hardware-versions/", handler.HardwareVersionActionHandler)

	mux.HandleFunc("/api/hardware-tests", handler.HardwareTestsHandler)
	mux.HandleFunc("/api/hardware-tests/", handler.HardwareTestActionHandler)

	// 软件版本
	mux.HandleFunc("/api/software-versions", handler.SoftwareVersionsHandler)
	mux.HandleFunc("/api/software-versions/", handler.SoftwareVersionActionHandler)

	// 项目分支
	mux.HandleFunc("/api/branches", handler.BranchesHandler)
	mux.HandleFunc("/api/branches/", handler.BranchActionHandler)
	// 测试用例,报告
	mux.HandleFunc("/api/test-cases", handler.TestCasesHandler)
	mux.HandleFunc("/api/test-cases/", handler.TestCaseActionHandler)

	// // 测试报告
	// mux.HandleFunc("/api/test-reports", handler.TestReportsHandler)
	// mux.HandleFunc("/api/test-reports/", handler.TestReportActionHandler)

	// 问题闭环
	mux.HandleFunc("/api/issues", handler.IssuesHandler)
	mux.HandleFunc("/api/issues/", handler.IssueActionHandler)
	// 生产工单
	mux.HandleFunc("/api/production-orders", handler.ProductionOrdersHandler)
	mux.HandleFunc("/api/production-orders/", handler.ProductionOrderActionHandler)

	// 烧录记录
	mux.HandleFunc("/api/burn-records", handler.BurnRecordsHandler)
	mux.HandleFunc("/api/burn-records/", handler.BurnRecordActionHandler)

	// 出厂测试
	mux.HandleFunc("/api/factory-tests", handler.FactoryTestsHandler)
	// mux.HandleFunc("/api/factory-tests/import", handler.ImportFactoryTestsHandler)
	// mux.HandleFunc("/api/factory-tests/delete", handler.DeleteFactoryTestsHandler)
	// mux.HandleFunc("/api/factory-tests/submit", handler.SubmitFactoryTestsHandler)
	// mux.HandleFunc("/api/factory-tests/audit", handler.AuditFactoryTestsHandler)
	mux.HandleFunc("/api/factory-tests/", handler.FactoryTestActionHandler)

	// 库存
	mux.HandleFunc("/api/inventory", handler.InventoryHandler)
	mux.HandleFunc("/api/inventory/", handler.InventoryActionHandler)

	// 发货批次
	mux.HandleFunc("/api/shipping-batches", handler.ShippingBatchesHandler)
	mux.HandleFunc("/api/shipping-batches/", handler.ShippingBatchActionHandler)
	mux.HandleFunc("/api/production-requests", handler.ProductionRequestsHandler)
	mux.HandleFunc("/api/production-requests/", handler.ProductionRequestActionHandler)

	// 发货批次设备明细
	mux.HandleFunc("/api/shipping-batch-devices", handler.ShippingBatchDevicesHandler)

	// 出库记录
	mux.HandleFunc("/api/outbound-records", handler.OutboundRecordsHandler)
	mux.HandleFunc("/api/outbound-records/", handler.OutboundRecordActionHandler)

	// 维修记录
	mux.HandleFunc("/api/repair-records", handler.RepairRecordsHandler)
	mux.HandleFunc("/api/repair-records/", handler.RepairRecordActionHandler)

	// 故障分析
	mux.HandleFunc("/api/fault-analysis", handler.FaultAnalysisHandler)
	mux.HandleFunc("/api/fault-analysis/", handler.FaultAnalysisActionHandler)

	// 首页统计
	mux.HandleFunc("/api/dashboard/summary", handler.DashboardSummaryHandler)
	mux.HandleFunc("/api/notifications/read", handler.MarkNotificationReadHandler)
	mux.HandleFunc("/api/global-search", handler.GlobalSearchHandler)

	// 报表
	mux.HandleFunc("/api/reports/project-progress", handler.ProjectProgressReportHandler)
	mux.HandleFunc("/api/reports/version-matrix", handler.VersionMatrixReportHandler)
	mux.HandleFunc("/api/reports/issue-statistics", handler.IssueStatisticsReportHandler)
	mux.HandleFunc("/api/reports/trace", handler.TraceReportHandler)

	return middleware.CORS(mux)
}
