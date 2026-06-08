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
	mux.HandleFunc("/api/project-branches", handler.ProjectBranchesHandler)
	mux.HandleFunc("/api/project-branches/", handler.ProjectBranchActionHandler)
	// 测试用例
	mux.HandleFunc("/api/test-cases", handler.TestCasesHandler)
	mux.HandleFunc("/api/test-cases/", handler.TestCaseActionHandler)

	// 测试报告
	mux.HandleFunc("/api/test-reports", handler.TestReportsHandler)
	mux.HandleFunc("/api/test-reports/", handler.TestReportActionHandler)

	// 问题闭环
	mux.HandleFunc("/api/issues", handler.IssuesHandler)
	mux.HandleFunc("/api/issues/", handler.IssueActionHandler)

	return middleware.CORS(mux)
}
