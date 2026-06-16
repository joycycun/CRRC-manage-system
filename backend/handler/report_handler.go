package handler

import (
	"crrc_pm_backend/config"
	"crrc_pm_backend/model"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
)

// ============================================================
// GET /api/dashboard/summary
// 首页统计
// ============================================================

func DashboardSummaryHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodGet {
		http.Error(w, "请求方法错误", http.StatusMethodNotAllowed)
		return
	}

	var data model.DashboardSummary
	userID, _ := strconv.ParseInt(r.URL.Query().Get("userId"), 10, 64)
	username := strings.TrimSpace(r.URL.Query().Get("username"))
	realName := strings.TrimSpace(r.URL.Query().Get("realName"))
	department := strings.TrimSpace(r.URL.Query().Get("department"))

	err := config.DB.QueryRow(`
		SELECT
			(SELECT COUNT(*) FROM projects WHERE IFNULL(is_deleted, 0) = 0 AND status = '进行中') AS ongoing_projects,
			(SELECT COUNT(*) FROM projects
			 WHERE IFNULL(is_deleted, 0) = 0
			   AND status = '进行中'
			   AND COALESCE(submit_time, created_at) >= DATE_FORMAT(CURDATE(), '%Y-%m-01')
			   AND COALESCE(submit_time, created_at) < DATE_ADD(DATE_FORMAT(CURDATE(), '%Y-%m-01'), INTERVAL 1 MONTH)) AS current_month_projects,
			(SELECT COUNT(*) FROM projects
			 WHERE IFNULL(is_deleted, 0) = 0
			   AND status = '进行中'
			   AND COALESCE(submit_time, created_at) >= DATE_SUB(DATE_FORMAT(CURDATE(), '%Y-%m-01'), INTERVAL 1 MONTH)
			   AND COALESCE(submit_time, created_at) < DATE_FORMAT(CURDATE(), '%Y-%m-01')) AS previous_month_projects,
			(SELECT COUNT(*) FROM issues
			 WHERE is_deleted = 0
			   AND IFNULL(close_status,'打开') NOT IN ('关闭', '已关闭', 'closed')
			   AND (
				 ? = 0
				 OR owner_id = ?
				 OR owner_name IN (?, ?)
			   )) AS pending_issues,
			(SELECT COUNT(*) FROM software_versions 
			 WHERE is_deleted = 0 
			   AND software_status = '已发布'
			   AND release_date IS NOT NULL
			   AND DATE_FORMAT(release_date, '%Y-%m') = DATE_FORMAT(CURDATE(), '%Y-%m')) AS current_month_releases,
			(SELECT COUNT(*) FROM software_versions 
			 WHERE is_deleted = 0 
			   AND software_status = '已发布'
			   AND release_date IS NOT NULL
			   AND release_date >= DATE_SUB(DATE_FORMAT(CURDATE(), '%Y-%m-01'), INTERVAL 1 MONTH)
			   AND release_date < DATE_FORMAT(CURDATE(), '%Y-%m-01')) AS previous_month_releases,
			(SELECT COUNT(*) FROM inventory_devices WHERE is_deleted = 0 AND inventory_status = '在库') AS inventory_in_stock,
			(SELECT COUNT(*) FROM inventory_devices WHERE is_deleted = 0 AND inventory_status = '已出库') AS inventory_outbound,
			(SELECT COUNT(*) FROM repair_records WHERE is_deleted = 0 AND IFNULL(confirm_status,'待确认') <> '已完成') AS repair_pending
	`, userID, userID, realName, username).Scan(
		&data.OngoingProjects,
		&data.OngoingProjectMonth,
		&data.OngoingProjectPrev,
		&data.PendingIssues,
		&data.CurrentMonthReleases,
		&data.ReleasePrevMonth,
		&data.InventoryInStock,
		&data.InventoryOutbound,
		&data.RepairPending,
	)

	if err != nil {
		http.Error(w, "首页统计失败: "+err.Error(), http.StatusInternalServerError)
		return
	}

	data.OngoingProjectDelta = data.OngoingProjectMonth - data.OngoingProjectPrev
	data.ReleaseDelta = data.CurrentMonthReleases - data.ReleasePrevMonth
	data.Kpis = []model.DashboardKpiItem{
		{
			Title: "进行中项目",
			Value: data.OngoingProjects,
			Desc:  formatDeltaDesc(data.OngoingProjectDelta, "较上月"),
			Icon:  "项",
		},
		{
			Title: "待处理问题",
			Value: data.PendingIssues,
			Desc:  "当前负责人未闭环问题",
			Icon:  "问",
		},
		{
			Title: "本月发布版本",
			Value: data.CurrentMonthReleases,
			Desc:  formatDeltaDesc(data.ReleaseDelta, "较上月"),
			Icon:  "版",
		},
	}

	data.Todos, err = queryDashboardTodos(userID, username, realName, department)
	if err != nil {
		http.Error(w, "查询待办事项失败: "+err.Error(), http.StatusInternalServerError)
		return
	}

	data.Notifications, err = queryAuditResultNotifications(userID, username, realName)
	if err != nil {
		http.Error(w, "查询消息提醒失败: "+err.Error(), http.StatusInternalServerError)
		return
	}
	productionNotifications, err := queryProductionRequestNotifications(userID, username, realName)
	if err != nil {
		http.Error(w, "查询生产请求通知失败: "+err.Error(), http.StatusInternalServerError)
		return
	}
	data.Notifications = append(data.Notifications, productionNotifications...)
	issueNotifications, err := queryIssueConfirmedNotifications(userID, username, realName)
	if err != nil {
		http.Error(w, "查询问题确认通知失败: "+err.Error(), http.StatusInternalServerError)
		return
	}
	data.Notifications = append(data.Notifications, issueNotifications...)

	data.RecentReleases, err = queryRecentSoftwareReleases()
	if err != nil {
		http.Error(w, "查询近期软件发布失败: "+err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"code": 200,
		"msg":  "查询成功",
		"data": data,
	})
}

func ensureNotificationReadTable() error {
	_, err := config.DB.Exec(`
		CREATE TABLE IF NOT EXISTS notification_reads (
			id BIGINT PRIMARY KEY AUTO_INCREMENT,
			user_id BIGINT NOT NULL DEFAULT 0,
			username VARCHAR(64) NOT NULL DEFAULT '',
			notification_id VARCHAR(160) NOT NULL,
			read_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
			UNIQUE KEY uk_notification_read (user_id, username, notification_id)
		) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4
	`)
	return err
}

func MarkNotificationReadHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodPost {
		http.Error(w, "请求方法错误", http.StatusMethodNotAllowed)
		return
	}

	if err := ensureNotificationReadTable(); err != nil {
		http.Error(w, "初始化消息已读表失败: "+err.Error(), http.StatusInternalServerError)
		return
	}

	var req struct {
		UserID         int64  `json:"userId"`
		Username       string `json:"username"`
		NotificationID string `json:"notificationId"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "参数解析失败: "+err.Error(), http.StatusBadRequest)
		return
	}

	req.Username = strings.TrimSpace(req.Username)
	req.NotificationID = strings.TrimSpace(req.NotificationID)

	if req.NotificationID == "" {
		http.Error(w, "消息ID不能为空", http.StatusBadRequest)
		return
	}

	if req.UserID == 0 && req.Username == "" {
		http.Error(w, "用户信息不能为空", http.StatusBadRequest)
		return
	}

	_, err := config.DB.Exec(`
		INSERT INTO notification_reads (user_id, username, notification_id, read_at)
		VALUES (?, ?, ?, NOW())
		ON DUPLICATE KEY UPDATE read_at = NOW()
	`, req.UserID, req.Username, req.NotificationID)
	if err != nil {
		http.Error(w, "标记消息已读失败: "+err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"code": 200,
		"msg":  "已读成功",
	})
}

func formatDeltaDesc(delta int64, prefix string) string {
	if delta > 0 {
		return prefix + "增加 " + strconv.FormatInt(delta, 10) + " 个"
	}
	if delta < 0 {
		return prefix + "减少 " + strconv.FormatInt(-delta, 10) + " 个"
	}
	return prefix + "持平"
}

func queryDashboardTodos(userID int64, username string, realName string, department string) ([]model.DashboardTodoItem, error) {
	todos := make([]model.DashboardTodoItem, 0)
	identity := username + " " + realName + " " + department
	isLeader := userID == 0 || username == "admin" ||
		strings.Contains(identity, "领导") ||
		strings.Contains(identity, "管理") ||
		strings.Contains(identity, "经理") ||
		strings.Contains(identity, "主管")
	isProduction := userID == 0 || username == "admin" ||
		strings.Contains(identity, "生产")

	if isLeader {
		items, err := queryAuditTodos()
		if err != nil {
			return nil, err
		}
		todos = append(todos, items...)
	}

	issueTodos, err := queryIssueTodos(userID, username, realName)
	if err != nil {
		return nil, err
	}
	todos = append(todos, issueTodos...)

	if isProduction {
		items, err := queryBurnTestTodos()
		if err != nil {
			return nil, err
		}
		todos = append(todos, items...)

		requestItems, err := queryProductionRequestTodos()
		if err != nil {
			return nil, err
		}
		todos = append(todos, requestItems...)
	}

	if len(todos) > 8 {
		return todos[:8], nil
	}
	return todos, nil
}

func queryAuditTodos() ([]model.DashboardTodoItem, error) {
	rows, err := config.DB.Query(`
		SELECT id, title, deadline, level, todo_type, link
		FROM (
			SELECT
				CONCAT('project-audit-', id) AS id,
				CONCAT('审核项目立项：', project_name) AS title,
				IFNULL(DATE_FORMAT(submit_time, '%m-%d'), DATE_FORMAT(updated_at, '%m-%d')) AS deadline,
				'高' AS level,
				'projectAudit' AS todo_type,
				'/project/manage' AS link,
				COALESCE(submit_time, updated_at) AS sort_time
			FROM projects
			WHERE IFNULL(is_deleted, 0) = 0
			  AND audit_status = '待审核'

			UNION ALL

			SELECT
				CONCAT('requirement-book-audit-', rb.id),
				CONCAT('审核需求书：', rb.book_name),
				IFNULL(DATE_FORMAT(rb.submit_time, '%m-%d'), DATE_FORMAT(rb.updated_at, '%m-%d')),
				'中',
				'requirementBookAudit',
				'/requirement/book',
				COALESCE(rb.submit_time, rb.updated_at)
			FROM requirement_books rb
			WHERE IFNULL(rb.is_deleted, 0) = 0
			  AND rb.status IN ('待审核', 'submitted', '已提交')

			UNION ALL

			SELECT
				CONCAT('requirement-change-audit-', rc.id),
				CONCAT('审核需求变更：', rc.change_title),
				IFNULL(DATE_FORMAT(rc.submit_time, '%m-%d'), DATE_FORMAT(rc.updated_at, '%m-%d')),
				'中',
				'requirementChangeAudit',
				'/requirement/change',
				COALESCE(rc.submit_time, rc.updated_at)
			FROM requirement_changes rc
			WHERE IFNULL(rc.is_deleted, 0) = 0
			  AND rc.status IN ('待审核', 'submitted', '已提交')

			UNION ALL

			SELECT
				CONCAT('test-case-audit-', tc.id),
				CONCAT('审核测试用例：', IFNULL(NULLIF(tc.case_name, ''), IFNULL(NULLIF(tc.file_name, ''), '未命名测试用例'))),
				IFNULL(DATE_FORMAT(tc.upload_time, '%m-%d'), DATE_FORMAT(tc.updated_at, '%m-%d')),
				'中',
				'testCaseAudit',
				'/test/case',
				COALESCE(tc.upload_time, tc.updated_at)
			FROM test_cases tc
			WHERE IFNULL(tc.is_deleted, 0) = 0
			  AND tc.audit_status IN ('待审核', 'submitted', '已提交')

			UNION ALL

			SELECT
				CONCAT('hardware-test-audit-', ht.id),
				CONCAT('审核硬件测试：', IFNULL(NULLIF(ht.record_name, ''), '未命名硬件测试')),
				IFNULL(DATE_FORMAT(ht.upload_time, '%m-%d'), DATE_FORMAT(ht.updated_at, '%m-%d')),
				'中',
				'hardwareTestAudit',
				'/hardware/test',
				COALESCE(ht.upload_time, ht.updated_at)
			FROM hardware_tests ht
			WHERE IFNULL(ht.is_deleted, 0) = 0
			  AND ht.audit_status IN ('待审核', 'submitted', '已提交')

			UNION ALL

			SELECT
				CONCAT('factory-test-audit-', ft.id),
				CONCAT('审核出厂测试：', IFNULL(NULLIF(ft.product_model, ''), IFNULL(NULLIF(ft.sn, ''), '未命名出厂测试'))),
				IFNULL(DATE_FORMAT(ft.upload_time, '%m-%d'), DATE_FORMAT(ft.updated_at, '%m-%d')),
				'高',
				'factoryTestAudit',
				'/production/factory-test',
				COALESCE(ft.upload_time, ft.updated_at)
			FROM factory_tests ft
			WHERE IFNULL(ft.is_deleted, 0) = 0
			  AND ft.audit_status IN ('待审核', 'submitted', '已提交')

			UNION ALL

			SELECT
				CONCAT('shipping-batch-audit-', sb.id),
				CONCAT('审核发货批次：', sb.batch_no),
				IFNULL(DATE_FORMAT(sb.upload_time, '%m-%d'), DATE_FORMAT(sb.updated_at, '%m-%d')),
				'高',
				'shippingBatchAudit',
				'/shipping/batch',
				COALESCE(sb.upload_time, sb.updated_at)
			FROM shipping_batches sb
			WHERE IFNULL(sb.is_deleted, 0) = 0
			  AND sb.audit_status IN ('待审核', 'submitted', '已提交')

			UNION ALL

			SELECT
				CONCAT('fault-analysis-audit-', fa.id),
				CONCAT('审核故障分析：', fa.analysis_name),
				IFNULL(DATE_FORMAT(fa.submit_time, '%m-%d'), DATE_FORMAT(fa.updated_at, '%m-%d')),
				'中',
				'faultAnalysisAudit',
				'/aftersales/fault-analysis',
				COALESCE(fa.submit_time, fa.updated_at)
			FROM fault_analysis fa
			WHERE IFNULL(fa.is_deleted, 0) = 0
			  AND fa.audit_status IN ('待审核', 'submitted', '已提交')
		) pending_audits
		ORDER BY sort_time ASC
		LIMIT 20
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	list := make([]model.DashboardTodoItem, 0)
	for rows.Next() {
		var id string
		var title string
		var deadline string
		var level string
		var todoType string
		var link string
		if err := rows.Scan(&id, &title, &deadline, &level, &todoType, &link); err != nil {
			return nil, err
		}
		list = append(list, model.DashboardTodoItem{
			ID:       id,
			Title:    title,
			Deadline: deadline,
			Level:    level,
			Type:     todoType,
			Link:     link,
		})
	}
	return list, rows.Err()
}

func queryAuditResultNotifications(userID int64, username string, realName string) ([]model.DashboardTodoItem, error) {
	if err := ensureNotificationReadTable(); err != nil {
		return nil, err
	}

	rows, err := config.DB.Query(`
		SELECT id, title, deadline, level, todo_type, link
		FROM (
			SELECT
				CONCAT('project-audit-result-', p.id) AS id,
				CONCAT('项目立项', IF(p.audit_status IN ('已通过', '审核通过', 'approved'), '审核通过：', '审核驳回：'), p.project_name) AS title,
				IFNULL(DATE_FORMAT(p.audit_time, '%m-%d'), DATE_FORMAT(p.updated_at, '%m-%d')) AS deadline,
				IF(p.audit_status IN ('已驳回', '审核驳回', 'rejected'), '高', '中') AS level,
				'projectAuditResult' AS todo_type,
				'/project/manage' AS link,
				COALESCE(p.audit_time, p.updated_at) AS sort_time
			FROM projects p
			WHERE IFNULL(p.is_deleted, 0) = 0
			  AND p.audit_status IN ('已通过', '审核通过', 'approved', '已驳回', '审核驳回', 'rejected')
			  AND (? = 0 OR p.owner_id = ? OR p.owner_name IN (?, ?))

			UNION ALL

			SELECT
				CONCAT('requirement-book-audit-result-', rb.id),
				CONCAT('需求书', IF(rb.status IN ('已通过', '审核通过', 'approved'), '审核通过：', '审核驳回：'), rb.book_name),
				IFNULL(DATE_FORMAT(rb.audit_time, '%m-%d'), DATE_FORMAT(rb.updated_at, '%m-%d')),
				IF(rb.status IN ('已驳回', '审核驳回', 'rejected'), '高', '中'),
				'requirementBookAuditResult',
				'/requirement/book',
				COALESCE(rb.audit_time, rb.updated_at)
			FROM requirement_books rb
			WHERE IFNULL(rb.is_deleted, 0) = 0
			  AND rb.status IN ('已通过', '审核通过', 'approved', '已驳回', '审核驳回', 'rejected')
			  AND (? = 0 OR rb.submit_user_id = ? OR rb.submit_user_name IN (?, ?))

			UNION ALL

			SELECT
				CONCAT('requirement-change-audit-result-', rc.id),
				CONCAT('需求变更', IF(rc.status IN ('已通过', '审核通过', 'approved'), '审核通过：', '审核驳回：'), rc.change_title),
				IFNULL(DATE_FORMAT(rc.audit_time, '%m-%d'), DATE_FORMAT(rc.updated_at, '%m-%d')),
				IF(rc.status IN ('已驳回', '审核驳回', 'rejected'), '高', '中'),
				'requirementChangeAuditResult',
				'/requirement/change',
				COALESCE(rc.audit_time, rc.updated_at)
			FROM requirement_changes rc
			WHERE IFNULL(rc.is_deleted, 0) = 0
			  AND rc.status IN ('已通过', '审核通过', 'approved', '已驳回', '审核驳回', 'rejected')
			  AND (? = 0 OR rc.submit_user_id = ? OR rc.submit_user_name IN (?, ?))

			UNION ALL

			SELECT
				CONCAT('test-case-audit-result-', tc.id),
				CONCAT('测试用例', IF(tc.audit_status IN ('已通过', '审核通过', 'approved'), '审核通过：', '审核驳回：'), IFNULL(NULLIF(tc.case_name, ''), IFNULL(NULLIF(tc.file_name, ''), '未命名测试用例'))),
				IFNULL(DATE_FORMAT(tc.audit_time, '%m-%d'), DATE_FORMAT(tc.updated_at, '%m-%d')),
				IF(tc.audit_status IN ('已驳回', '审核驳回', 'rejected'), '高', '中'),
				'testCaseAuditResult',
				'/test/case',
				COALESCE(tc.audit_time, tc.updated_at)
			FROM test_cases tc
			WHERE IFNULL(tc.is_deleted, 0) = 0
			  AND tc.audit_status IN ('已通过', '审核通过', 'approved', '已驳回', '审核驳回', 'rejected')
			  AND (? = 0 OR tc.uploader_id = ? OR tc.uploader_name IN (?, ?))

			UNION ALL

			SELECT
				CONCAT('hardware-test-audit-result-', ht.id),
				CONCAT('硬件测试', IF(ht.audit_status IN ('已通过', '审核通过', 'approved'), '审核通过：', '审核驳回：'), IFNULL(NULLIF(ht.record_name, ''), '未命名硬件测试')),
				IFNULL(DATE_FORMAT(ht.audit_time, '%m-%d'), DATE_FORMAT(ht.updated_at, '%m-%d')),
				IF(ht.audit_status IN ('已驳回', '审核驳回', 'rejected'), '高', '中'),
				'hardwareTestAuditResult',
				'/hardware/test',
				COALESCE(ht.audit_time, ht.updated_at)
			FROM hardware_tests ht
			WHERE IFNULL(ht.is_deleted, 0) = 0
			  AND ht.audit_status IN ('已通过', '审核通过', 'approved', '已驳回', '审核驳回', 'rejected')
			  AND (? = 0 OR ht.uploader_id = ? OR ht.uploader_name IN (?, ?))

			UNION ALL

			SELECT
				CONCAT(
					'factory-test-audit-result-',
					IFNULL(ft.uploader_id, 0),
					'-',
					REPLACE(IFNULL(NULLIF(ft.product_model, ''), 'unknown'), ' ', '_'),
					'-',
					IF(ft.audit_status IN ('已通过', '审核通过', 'approved'), 'approved', 'rejected')
				),
				CONCAT(
					'出厂测试',
					IF(ft.audit_status IN ('已通过', '审核通过', 'approved'), '审核通过：', '审核驳回：'),
					IFNULL(NULLIF(ft.product_model, ''), '未命名型号'),
					'（',
					COUNT(*),
					' 台）'
				),
				IFNULL(DATE_FORMAT(MAX(ft.audit_time), '%m-%d'), DATE_FORMAT(MAX(ft.updated_at), '%m-%d')),
				IF(ft.audit_status IN ('已驳回', '审核驳回', 'rejected'), '高', '中'),
				'factoryTestAuditResult',
				'/production/factory-test',
				MAX(COALESCE(ft.audit_time, ft.updated_at))
			FROM factory_tests ft
			WHERE IFNULL(ft.is_deleted, 0) = 0
			  AND ft.audit_status IN ('已通过', '审核通过', 'approved', '已驳回', '审核驳回', 'rejected')
			  AND (? = 0 OR ft.uploader_id = ? OR ft.uploader_name IN (?, ?))
			GROUP BY
				IFNULL(ft.uploader_id, 0),
				IFNULL(ft.uploader_name, ''),
				IFNULL(ft.product_model, ''),
				IF(ft.audit_status IN ('已通过', '审核通过', 'approved'), 'approved', 'rejected')

			UNION ALL

			SELECT
				CONCAT('shipping-batch-audit-result-', sb.id),
				CONCAT('发货批次', IF(sb.audit_status IN ('已通过', '审核通过', 'approved'), '审核通过：', '审核驳回：'), sb.batch_no),
				IFNULL(DATE_FORMAT(sb.audit_time, '%m-%d'), DATE_FORMAT(sb.updated_at, '%m-%d')),
				IF(sb.audit_status IN ('已驳回', '审核驳回', 'rejected'), '高', '中'),
				'shippingBatchAuditResult',
				'/shipping/batch',
				COALESCE(sb.audit_time, sb.updated_at)
			FROM shipping_batches sb
			WHERE IFNULL(sb.is_deleted, 0) = 0
			  AND sb.audit_status IN ('已通过', '审核通过', 'approved', '已驳回', '审核驳回', 'rejected')
			  AND (? = 0 OR sb.uploader_id = ? OR sb.uploader_name IN (?, ?))

			UNION ALL

			SELECT
				CONCAT('fault-analysis-audit-result-', fa.id),
				CONCAT('故障分析', IF(fa.audit_status IN ('已通过', '审核通过', 'approved'), '审核通过：', '审核驳回：'), fa.analysis_name),
				IFNULL(DATE_FORMAT(fa.audit_time, '%m-%d'), DATE_FORMAT(fa.updated_at, '%m-%d')),
				IF(fa.audit_status IN ('已驳回', '审核驳回', 'rejected'), '高', '中'),
				'faultAnalysisAuditResult',
				'/aftersales/fault-analysis',
				COALESCE(fa.audit_time, fa.updated_at)
			FROM fault_analysis fa
			WHERE IFNULL(fa.is_deleted, 0) = 0
			  AND fa.audit_status IN ('已通过', '审核通过', 'approved', '已驳回', '审核驳回', 'rejected')
			  AND (? = 0 OR fa.submit_user_id = ? OR fa.submit_user_name IN (?, ?))
		) audit_results
		WHERE NOT EXISTS (
			SELECT 1
			FROM notification_reads nr
			WHERE nr.notification_id = audit_results.id
			  AND (
				(? <> 0 AND nr.user_id = ?)
				OR (? <> '' AND nr.username = ?)
			  )
		)
		ORDER BY sort_time DESC
		LIMIT 12
	`,
		userID, userID, realName, username,
		userID, userID, realName, username,
		userID, userID, realName, username,
		userID, userID, realName, username,
		userID, userID, realName, username,
		userID, userID, realName, username,
		userID, userID, realName, username,
		userID, userID, realName, username,
		userID, userID, username, username,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	list := make([]model.DashboardTodoItem, 0)
	for rows.Next() {
		var item model.DashboardTodoItem
		if err := rows.Scan(&item.ID, &item.Title, &item.Deadline, &item.Level, &item.Type, &item.Link); err != nil {
			return nil, err
		}
		list = append(list, item)
	}
	return list, rows.Err()
}

func queryProductionRequestNotifications(userID int64, username string, realName string) ([]model.DashboardTodoItem, error) {
	if err := ensureProductionRequestTable(); err != nil {
		return nil, err
	}
	if err := ensureNotificationReadTable(); err != nil {
		return nil, err
	}

	rows, err := config.DB.Query(`
		SELECT id, title, deadline, level, todo_type, link
		FROM (
			SELECT
				CONCAT('production-request-confirmed-', pr.id) AS id,
				CONCAT('生产已确认：', pr.product_model, '（', pr.quantity, ' 台）') AS title,
				IFNULL(DATE_FORMAT(pr.confirm_time, '%m-%d'), DATE_FORMAT(pr.updated_at, '%m-%d')) AS deadline,
				'中' AS level,
				'productionRequestConfirmed' AS todo_type,
				'/shipping/batch' AS link,
				COALESCE(pr.confirm_time, pr.updated_at) AS sort_time
			FROM production_requests pr
			WHERE IFNULL(pr.is_deleted, 0) = 0
			  AND pr.status = 'confirmed'
			  AND (? = 0 OR pr.requester_id = ? OR pr.requester_name IN (?, ?))
		) production_results
		WHERE NOT EXISTS (
			SELECT 1
			FROM notification_reads nr
			WHERE nr.notification_id = production_results.id
			  AND (
				(? <> 0 AND nr.user_id = ?)
				OR (? <> '' AND nr.username = ?)
			  )
		)
		ORDER BY sort_time DESC
		LIMIT 8
	`, userID, userID, realName, username, userID, userID, username, username)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	list := make([]model.DashboardTodoItem, 0)
	for rows.Next() {
		var item model.DashboardTodoItem
		if err := rows.Scan(&item.ID, &item.Title, &item.Deadline, &item.Level, &item.Type, &item.Link); err != nil {
			return nil, err
		}
		list = append(list, item)
	}
	return list, rows.Err()
}

func queryIssueTodos(userID int64, username string, realName string) ([]model.DashboardTodoItem, error) {
	if err := ensureIssueConfirmationTable(); err != nil {
		return nil, err
	}

	rows, err := config.DB.Query(`
		SELECT
			id,
			issue_title,
			IFNULL(DATE_FORMAT(plan_close_time, '%m-%d'), DATE_FORMAT(updated_at, '%m-%d')),
			IFNULL(level, '')
		FROM issues
		WHERE IFNULL(is_deleted, 0) = 0
		  AND IFNULL(close_status, '打开') NOT IN ('关闭', '已关闭', 'closed')
		  AND (
			 ? = 0
			 OR owner_id = ?
			 OR owner_name IN (?, ?)
		  )
		  AND NOT EXISTS (
			SELECT 1
			FROM issue_confirmations ic
			WHERE ic.issue_id = issues.id
			  AND (
				(? <> 0 AND ic.confirm_user_id = ?)
				OR (? <> '' AND ic.confirm_user_name = ?)
			  )
		  )
		ORDER BY
			CASE IFNULL(level, '')
				WHEN '高' THEN 1
				WHEN '中' THEN 2
				WHEN '低' THEN 3
				ELSE 4
			END,
			plan_close_time ASC,
			updated_at DESC
		LIMIT 4
	`, userID, userID, realName, username, userID, userID, realName, realName)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	list := make([]model.DashboardTodoItem, 0)
	for rows.Next() {
		var id int64
		var title string
		var deadline string
		var level string
		if err := rows.Scan(&id, &title, &deadline, &level); err != nil {
			return nil, err
		}
		list = append(list, model.DashboardTodoItem{
			ID:       "issue-" + strconv.FormatInt(id, 10),
			Title:    "处理问题闭环：" + title,
			Deadline: deadline,
			Level:    normalizeTodoLevel(level),
			Type:     "issue",
			Link:     "/test/issue",
		})
	}
	return list, rows.Err()
}

func queryIssueConfirmedNotifications(userID int64, username string, realName string) ([]model.DashboardTodoItem, error) {
	if err := ensureIssueConfirmationTable(); err != nil {
		return nil, err
	}
	if err := ensureNotificationReadTable(); err != nil {
		return nil, err
	}

	rows, err := config.DB.Query(`
		SELECT id, title, deadline, level, todo_type, link
		FROM (
			SELECT
				CONCAT('issue-confirmed-', i.id, '-', ic.id) AS id,
				CONCAT('问题已确认：', i.issue_title, '（', ic.confirm_user_name, '）') AS title,
				IFNULL(DATE_FORMAT(ic.confirm_time, '%m-%d'), '') AS deadline,
				'中' AS level,
				'issueConfirmed' AS todo_type,
				'/test/issue' AS link,
				ic.confirm_time AS sort_time
			FROM issue_confirmations ic
			JOIN issues i ON i.id = ic.issue_id
			WHERE IFNULL(i.is_deleted, 0) = 0
			  AND (
				? = 0
				OR i.creator_id = ?
				OR i.creator_name IN (?, ?)
			  )
		) issue_results
		WHERE NOT EXISTS (
			SELECT 1
			FROM notification_reads nr
			WHERE nr.notification_id = issue_results.id
			  AND (
				(? <> 0 AND nr.user_id = ?)
				OR (? <> '' AND nr.username = ?)
			  )
		)
		ORDER BY sort_time DESC
		LIMIT 8
	`, userID, userID, realName, username, userID, userID, username, username)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	list := make([]model.DashboardTodoItem, 0)
	for rows.Next() {
		var item model.DashboardTodoItem
		if err := rows.Scan(&item.ID, &item.Title, &item.Deadline, &item.Level, &item.Type, &item.Link); err != nil {
			return nil, err
		}
		list = append(list, item)
	}
	return list, rows.Err()
}

func queryBurnTestTodos() ([]model.DashboardTodoItem, error) {
	rows, err := config.DB.Query(`
		SELECT
			IFNULL(NULLIF(br.batch_no, ''), CONCAT('未分批-', DATE_FORMAT(MIN(br.created_at), '%Y%m%d'))) AS batch_no,
			COUNT(DISTINCT br.id) AS burn_count,
			COUNT(DISTINCT ft.burn_record_id) AS tested_count,
			IFNULL(DATE_FORMAT(MAX(br.upload_time), '%m-%d'), DATE_FORMAT(MAX(br.created_at), '%m-%d')) AS deadline
		FROM burn_records br
		LEFT JOIN factory_tests ft
		  ON ft.burn_record_id = br.id
		 AND IFNULL(ft.is_deleted, 0) = 0
		WHERE IFNULL(br.is_deleted, 0) = 0
		GROUP BY IFNULL(NULLIF(br.batch_no, ''), CONCAT('未分批-', DATE_FORMAT(br.created_at, '%Y%m%d')))
		HAVING burn_count > tested_count
		ORDER BY MAX(br.created_at) DESC
		LIMIT 4
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	list := make([]model.DashboardTodoItem, 0)
	for rows.Next() {
		var batchNo string
		var burnCount int64
		var testedCount int64
		var deadline string
		if err := rows.Scan(&batchNo, &burnCount, &testedCount, &deadline); err != nil {
			return nil, err
		}
		list = append(list, model.DashboardTodoItem{
			ID:       "burn-test-" + batchNo,
			Title:    "烧录批次待测试：" + batchNo + "（" + strconv.FormatInt(burnCount-testedCount, 10) + " 台）",
			Deadline: deadline,
			Level:    "中",
			Type:     "burnTest",
			Link:     "/production/factory-test",
		})
	}
	return list, rows.Err()
}

func queryProductionRequestTodos() ([]model.DashboardTodoItem, error) {
	if err := ensureProductionRequestTable(); err != nil {
		return nil, err
	}

	rows, err := config.DB.Query(`
		SELECT
			id,
			product_model,
			IFNULL(device_type, ''),
			quantity,
			IFNULL(DATE_FORMAT(created_at, '%m-%d'), '')
		FROM production_requests
		WHERE IFNULL(is_deleted, 0) = 0
		  AND status = 'pending'
		ORDER BY created_at ASC
		LIMIT 8
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	list := make([]model.DashboardTodoItem, 0)
	for rows.Next() {
		var id int64
		var productModel string
		var deviceType string
		var quantity int64
		var deadline string
		if err := rows.Scan(&id, &productModel, &deviceType, &quantity, &deadline); err != nil {
			return nil, err
		}
		title := "生产请求：" + productModel + "（" + strconv.FormatInt(quantity, 10) + " 台）"
		if deviceType != "" {
			title = "生产请求：" + productModel + " / " + deviceType + "（" + strconv.FormatInt(quantity, 10) + " 台）"
		}
		list = append(list, model.DashboardTodoItem{
			ID:       "production-request-" + strconv.FormatInt(id, 10),
			Title:    title,
			Deadline: deadline,
			Level:    "高",
			Type:     "productionRequest",
			Link:     "/production/burn",
		})
	}
	return list, rows.Err()
}

func normalizeTodoLevel(level string) string {
	switch level {
	case "紧急", "高", "中", "低":
		return level
	case "严重":
		return "紧急"
	case "一般":
		return "中"
	case "轻微", "建议":
		return "低"
	default:
		return "中"
	}
}

func queryRecentSoftwareReleases() ([]model.DashboardReleaseItem, error) {
	rows, err := config.DB.Query(`
		SELECT
			id,
			software_version,
			IFNULL(DATE_FORMAT(COALESCE(release_date, updated_at), '%Y-%m-%d'), ''),
			IFNULL(software_status, '')
		FROM software_versions
		WHERE IFNULL(is_deleted, 0) = 0
		  AND software_status = '已发布'
		ORDER BY COALESCE(release_date, updated_at) DESC, id DESC
		LIMIT 5
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	list := make([]model.DashboardReleaseItem, 0)
	for rows.Next() {
		var item model.DashboardReleaseItem
		if err := rows.Scan(&item.ID, &item.Version, &item.ReleaseTime, &item.Status); err != nil {
			return nil, err
		}
		item.StatusType = releaseStatusType(item.Status)
		list = append(list, item)
	}
	return list, rows.Err()
}

func releaseStatusType(status string) string {
	switch status {
	case "已发布", "发布成功":
		return "success"
	case "开发中", "测试中", "进行中":
		return "running"
	case "已废弃", "已回滚":
		return "rollback"
	default:
		return "running"
	}
}

// ============================================================
// GET /api/reports/project-progress
// 项目进度报表
// ============================================================

func ProjectProgressReportHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	rows, err := config.DB.Query(`
		SELECT
			p.id,
			p.project_name,
			p.project_code,
			IFNULL(p.owner_name, ''),
			IFNULL(p.stage, ''),
			IFNULL(p.status, ''),
			IFNULL(DATE_FORMAT(p.updated_at, '%Y-%m-%d %H:%i:%s'), '')
		FROM projects p
		WHERE IFNULL(p.is_deleted, 0) = 0
		ORDER BY p.id DESC
	`)
	if err != nil {
		http.Error(w, "查询项目进度失败: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	type ProjectProgressVO struct {
		ID            int64   `json:"id"`
		ProjectID     int64   `json:"projectId"`
		Name          string  `json:"name"`
		ProjectName   string  `json:"projectName"`
		ProjectCode   string  `json:"projectCode"`
		Owner         string  `json:"owner"`
		CurrentStage  string  `json:"currentStage"`
		Stage         string  `json:"stage"`
		ProjectStatus string  `json:"projectStatus"`
		Status        string  `json:"status"`
		UpdateTime    string  `json:"updateTime"`
		Progress      float64 `json:"progress"`
	}

	list := make([]ProjectProgressVO, 0)

	for rows.Next() {
		var item ProjectProgressVO

		err := rows.Scan(
			&item.ProjectID,
			&item.ProjectName,
			&item.ProjectCode,
			&item.Owner,
			&item.Stage,
			&item.Status,
			&item.UpdateTime,
		)
		if err != nil {
			http.Error(w, "解析项目进度失败: "+err.Error(), http.StatusInternalServerError)
			return
		}

		item.ID = item.ProjectID
		item.Name = item.ProjectName
		item.CurrentStage = item.Stage
		item.ProjectStatus = "running"
		if item.Status == "已关闭" || item.Stage == "已关闭" {
			item.ProjectStatus = "closed"
		}
		item.Progress = projectStageProgress(item.Stage, item.Status)

		list = append(list, item)
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"code": 200,
		"msg":  "查询成功",
		"data": list,
	})
}

func projectStageProgress(stage string, status string) float64 {
	if status == "已关闭" || stage == "已关闭" {
		return 100
	}

	progressMap := map[string]float64{
		"立项":           10,
		"立项阶段":         10,
		"需求阶段":         20,
		"需求书首次确认":      20,
		"硬件检测":         30,
		"硬件开发":         30,
		"软件研发":         40,
		"内部初始测试":       50,
		"内部初始测试是否问题闭环": 60,
		"内部初始测试问题闭环":   60,
		"样机联调":         70,
		"出厂测试":         80,
		"收货":           90,
		"项目关闭":         100,
	}

	if progress, ok := progressMap[stage]; ok {
		return progress
	}

	return 0
}

// ============================================================
// GET /api/reports/version-matrix
// 版本矩阵
// ============================================================

func VersionMatrixReportHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	rows, err := config.DB.Query(`
		SELECT
			CONCAT('sw-', sv.id) AS id,
			IFNULL(p.id, 0) AS project_id,
			IFNULL(p.project_name, '') AS project_name,
			IFNULL(sv.device_type, '') AS device_type,
			IFNULL(COALESCE(NULLIF(sv.hardware_version, ''), hv.hardware_version), '') AS hardware_version,
			IFNULL(sv.software_version, '') AS software_version,
			IFNULL(sv.software_status, '') AS status,
			IFNULL(sv.owner_name, '') AS owner_name,
			IFNULL(DATE_FORMAT(sv.updated_at, '%Y-%m-%d %H:%i:%s'), '') AS update_time,
			IFNULL(sv.description, '') AS remark
		FROM software_versions sv
		LEFT JOIN projects p ON sv.project_id = p.id AND IFNULL(p.is_deleted, 0) = 0
		LEFT JOIN hardware_versions hv ON sv.hardware_id = hv.id
		WHERE IFNULL(sv.is_deleted, 0) = 0

		UNION ALL

		SELECT
			CONCAT('hw-', hv.id) AS id,
			IFNULL(p.id, 0) AS project_id,
			IFNULL(p.project_name, '') AS project_name,
			IFNULL(hv.device_type, '') AS device_type,
			IFNULL(hv.hardware_version, '') AS hardware_version,
			'' AS software_version,
			IFNULL(hv.status, '') AS status,
			IFNULL(hv.owner_name, '') AS owner_name,
			IFNULL(DATE_FORMAT(hv.updated_at, '%Y-%m-%d %H:%i:%s'), '') AS update_time,
			IFNULL(hv.description, '') AS remark
		FROM hardware_versions hv
		LEFT JOIN projects p ON hv.project_id = p.id AND IFNULL(p.is_deleted, 0) = 0
		WHERE NOT EXISTS (
			SELECT 1
			FROM software_versions sv
			WHERE IFNULL(sv.is_deleted, 0) = 0
			  AND sv.project_id = hv.project_id
			  AND IFNULL(sv.device_type, '') = IFNULL(hv.device_type, '')
			  AND (
				sv.hardware_id = hv.id
				OR IFNULL(sv.hardware_version, '') = IFNULL(hv.hardware_version, '')
			  )
		)
		ORDER BY project_id DESC, device_type, update_time DESC
	`)
	if err != nil {
		http.Error(w, "查询版本矩阵失败: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	type VersionMatrixVO struct {
		ID              string `json:"id"`
		ProjectID       int64  `json:"projectId"`
		ProjectName     string `json:"projectName"`
		DeviceType      string `json:"deviceType"`
		HardwareVersion string `json:"hardwareVersion"`
		SoftwareVersion string `json:"softwareVersion"`
		Status          string `json:"status"`
		StatusType      string `json:"statusType"`
		Owner           string `json:"owner"`
		OwnerName       string `json:"ownerName"`
		UpdateTime      string `json:"updateTime"`
		Remark          string `json:"remark"`
	}

	list := make([]VersionMatrixVO, 0)

	for rows.Next() {
		var item VersionMatrixVO

		err := rows.Scan(
			&item.ID,
			&item.ProjectID,
			&item.ProjectName,
			&item.DeviceType,
			&item.HardwareVersion,
			&item.SoftwareVersion,
			&item.Status,
			&item.OwnerName,
			&item.UpdateTime,
			&item.Remark,
		)
		if err != nil {
			http.Error(w, "解析版本矩阵失败: "+err.Error(), http.StatusInternalServerError)
			return
		}

		item.Owner = item.OwnerName
		item.StatusType = versionStatusType(item.Status)
		list = append(list, item)
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"code": 200,
		"msg":  "查询成功",
		"data": list,
	})
}

func versionStatusType(status string) string {
	switch status {
	case "草稿", "样品":
		return "developing"
	case "测试中", "集成测试", "Beta":
		return "testing"
	case "已发布", "正式发布":
		return "production"
	case "冻结", "版本冻结":
		return "frozen"
	default:
		return "developing"
	}
}

// ============================================================
// GET /api/reports/issue-statistics
// 问题统计
// ============================================================

func IssueStatisticsReportHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	rows, err := config.DB.Query(`
		SELECT
			CONCAT('issue-', i.id) AS id,
			IFNULL(p.id, 0) AS project_id,
			IFNULL(p.project_name, '') AS project_name,
			IFNULL(i.issue_title, '') AS issue_title,
			CASE
				WHEN IFNULL(i.reopen_count, 0) > 2 THEN 'serious'
				WHEN i.level IN ('严重', 'serious') THEN 'serious'
				WHEN i.level IN ('轻微', 'minor') THEN 'minor'
				WHEN i.level IN ('建议', 'suggestion') THEN 'suggestion'
				ELSE 'normal'
			END AS issue_level,
			CASE
				WHEN i.close_status IN ('关闭', '已关闭', '已闭环', 'closed') THEN 'closed'
				WHEN i.close_status IN ('处理中', 'processing') THEN 'processing'
				ELSE 'open'
			END AS issue_status,
			IFNULL(i.owner_name, '') AS owner_name,
			1 + IFNULL(i.reopen_count, 0) AS frequency,
			IFNULL(DATE_FORMAT(i.updated_at, '%Y-%m-%d %H:%i:%s'), '') AS update_time,
			IFNULL(i.issue_desc, '') AS remark,
			IFNULL(i.issue_source, '测试问题') AS issue_source,
			IFNULL(i.reopen_count, 0) AS reopen_count
		FROM issues i
		LEFT JOIN projects p ON i.project_id = p.id AND IFNULL(p.is_deleted, 0) = 0
		WHERE IFNULL(i.is_deleted, 0) = 0

		UNION ALL

		SELECT
			CONCAT('repair-', rr.id) AS id,
			IFNULL(p.id, 0) AS project_id,
			IFNULL(p.project_name, '') AS project_name,
			IFNULL(rr.fault_desc, '') AS issue_title,
			'normal' AS issue_level,
			'closed' AS issue_status,
			IFNULL(rr.repair_user_name, '') AS owner_name,
			1 AS frequency,
			IFNULL(DATE_FORMAT(rr.updated_at, '%Y-%m-%d %H:%i:%s'), '') AS update_time,
			IFNULL(rr.repair_process, IFNULL(rr.remark, '')) AS remark,
			'维修记录' AS issue_source,
			0 AS reopen_count
		FROM repair_records rr
		LEFT JOIN projects p ON rr.project_id = p.id AND IFNULL(p.is_deleted, 0) = 0
		WHERE IFNULL(rr.is_deleted, 0) = 0
		ORDER BY update_time DESC
	`)
	if err != nil {
		http.Error(w, "查询问题统计失败: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	type IssueStatisticsVO struct {
		ID          string `json:"id"`
		ProjectID   int64  `json:"projectId"`
		ProjectName string `json:"projectName"`
		IssueTitle  string `json:"issueTitle"`
		IssueLevel  string `json:"issueLevel"`
		IssueStatus string `json:"issueStatus"`
		Owner       string `json:"owner"`
		OwnerName   string `json:"ownerName"`
		Frequency   int64  `json:"frequency"`
		UpdateTime  string `json:"updateTime"`
		Remark      string `json:"remark"`
		IssueSource string `json:"issueSource"`
		ReopenCount int64  `json:"reopenCount"`
	}

	list := make([]IssueStatisticsVO, 0)

	for rows.Next() {
		var item IssueStatisticsVO

		err := rows.Scan(
			&item.ID,
			&item.ProjectID,
			&item.ProjectName,
			&item.IssueTitle,
			&item.IssueLevel,
			&item.IssueStatus,
			&item.OwnerName,
			&item.Frequency,
			&item.UpdateTime,
			&item.Remark,
			&item.IssueSource,
			&item.ReopenCount,
		)
		if err != nil {
			http.Error(w, "解析问题统计失败: "+err.Error(), http.StatusInternalServerError)
			return
		}

		item.Owner = item.OwnerName
		list = append(list, item)
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"code": 200,
		"msg":  "查询成功",
		"data": list,
	})
}

// ============================================================
// GET /api/global-search?keyword=xxx
// 全局搜索：项目、版本、SN/MAC库存与出库定位
// ============================================================

func GlobalSearchHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodGet {
		http.Error(w, "请求方法错误", http.StatusMethodNotAllowed)
		return
	}

	keyword := strings.TrimSpace(r.URL.Query().Get("keyword"))
	if keyword == "" {
		json.NewEncoder(w).Encode(map[string]interface{}{
			"code": 200,
			"msg":  "查询成功",
			"data": map[string]interface{}{
				"projects": []map[string]interface{}{},
				"versions": []map[string]interface{}{},
				"devices":  []map[string]interface{}{},
			},
		})
		return
	}

	like := "%" + keyword + "%"
	data := make(map[string]interface{})

	projects, err := queryTraceList(`
		SELECT
			id,
			project_name,
			project_code,
			IFNULL(owner_name, '') AS owner_name,
			IFNULL(stage, '') AS stage,
			IFNULL(status, '') AS status
		FROM projects
		WHERE IFNULL(is_deleted, 0) = 0
		  AND (project_name LIKE ? OR project_code LIKE ? OR owner_name LIKE ?)
		ORDER BY updated_at DESC
		LIMIT 8
	`, like, like, like)
	if err != nil {
		http.Error(w, "搜索项目失败: "+err.Error(), http.StatusInternalServerError)
		return
	}
	data["projects"] = projects

	versions, err := queryTraceList(`
		SELECT *
		FROM (
			SELECT
				'software' AS version_type,
				sv.id,
				IFNULL(p.project_name, '') AS project_name,
				IFNULL(sv.device_type, '') AS device_type,
				IFNULL(sv.software_version, '') AS software_version,
				IFNULL(sv.hardware_version, '') AS hardware_version,
				IFNULL(sv.software_status, '') AS version_status,
				sv.updated_at AS sort_time
			FROM software_versions sv
			LEFT JOIN projects p ON p.id = sv.project_id
			WHERE IFNULL(sv.is_deleted, 0) = 0
			  AND (
				sv.software_version LIKE ?
				OR sv.hardware_version LIKE ?
				OR sv.device_type LIKE ?
				OR p.project_name LIKE ?
			  )

			UNION ALL

			SELECT
				'hardware' AS version_type,
				hv.id,
				IFNULL(p.project_name, '') AS project_name,
				IFNULL(hv.device_type, '') AS device_type,
				'' AS software_version,
				IFNULL(hv.hardware_version, '') AS hardware_version,
				IFNULL(hv.status, '') AS version_status,
				hv.updated_at AS sort_time
			FROM hardware_versions hv
			LEFT JOIN projects p ON p.id = hv.project_id
			WHERE (
				hv.hardware_version LIKE ?
				OR hv.device_type LIKE ?
				OR p.project_name LIKE ?
			  )
		) version_matches
		ORDER BY sort_time DESC
		LIMIT 12
	`, like, like, like, like, like, like, like)
	if err != nil {
		http.Error(w, "搜索版本失败: "+err.Error(), http.StatusInternalServerError)
		return
	}
	data["versions"] = versions

	devices, err := queryTraceList(`
		SELECT *
		FROM (
			SELECT
				'inventory' AS location,
				inv.id,
				inv.id AS inventory_device_id,
				IFNULL(inv.sn, '') AS sn,
				IFNULL(inv.mac_address, '') AS mac_address,
				IFNULL(inv.device_type, '') AS device_type,
				IFNULL(inv.software_version, '') AS software_version,
				IFNULL(inv.hardware_version, '') AS hardware_version,
				IFNULL(inv.inventory_status, '') AS status,
				IFNULL(DATE_FORMAT(inv.update_time, '%Y-%m-%d %H:%i:%s'), '') AS event_time
			FROM inventory_devices inv
			WHERE IFNULL(inv.is_deleted, 0) = 0
			  AND IFNULL(inv.inventory_status, '') <> '已出库'
			  AND (inv.sn LIKE ? OR inv.mac_address LIKE ?)

			UNION ALL

			SELECT
				'outbound' AS location,
				o.id,
				o.inventory_device_id,
				IFNULL(o.sn, '') AS sn,
				IFNULL(o.mac_address, '') AS mac_address,
				IFNULL(sbd.device_type, IFNULL(inv.device_type, '')) AS device_type,
				IFNULL(sbd.software_version, IFNULL(inv.software_version, '')) AS software_version,
				IFNULL(sbd.hardware_version, IFNULL(inv.hardware_version, '')) AS hardware_version,
				IFNULL(o.status, '已出库') AS status,
				IFNULL(DATE_FORMAT(o.outbound_time, '%Y-%m-%d %H:%i:%s'), '') AS event_time
			FROM outbound_records o
			LEFT JOIN shipping_batch_devices sbd
				ON sbd.batch_id = o.batch_id
				AND sbd.inventory_device_id = o.inventory_device_id
				AND IFNULL(sbd.is_deleted, 0) = 0
			LEFT JOIN inventory_devices inv ON inv.id = o.inventory_device_id
			WHERE IFNULL(o.is_deleted, 0) = 0
			  AND (o.sn LIKE ? OR o.mac_address LIKE ?)
		) device_matches
		ORDER BY event_time DESC
		LIMIT 12
	`, like, like, like, like)
	if err != nil {
		http.Error(w, "搜索终端失败: "+err.Error(), http.StatusInternalServerError)
		return
	}
	data["devices"] = devices

	json.NewEncoder(w).Encode(map[string]interface{}{
		"code": 200,
		"msg":  "查询成功",
		"data": data,
	})
}

// ============================================================
// GET /api/reports/trace?sn=xxx&mac=xxx
// SN/MAC 全链路追溯
// ============================================================

func TraceReportHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	sn := r.URL.Query().Get("sn")
	mac := r.URL.Query().Get("mac")

	if sn == "" && mac == "" {
		http.Error(w, "sn 或 mac 至少传一个", http.StatusBadRequest)
		return
	}

	data := map[string]interface{}{}

	burnRecords, err := queryTraceList(`
		SELECT
			id,
			IFNULL(batch_no, ''),
			IFNULL(project_id, 0),
			IFNULL(production_order_id, 0),
			IFNULL(product_name, ''),
			IFNULL(product_model, ''),
			IFNULL(device_type, ''),
			sn,
			mac_address,
			IFNULL(hardware_version, ''),
			IFNULL(software_version, ''),
			upload_time,
			created_at
		FROM burn_records
		WHERE is_deleted = 0
		  AND (? = '' OR sn = ?)
		  AND (? = '' OR mac_address = ?)
		ORDER BY id DESC
	`, sn, sn, mac, mac)
	if err != nil {
		http.Error(w, "查询烧录记录失败: "+err.Error(), http.StatusInternalServerError)
		return
	}
	data["burnRecords"] = burnRecords

	factoryTests, err := queryTraceList(`
		SELECT
			id,
			burn_record_id,
			IFNULL(project_id, 0),
			IFNULL(product_model, ''),
			IFNULL(device_type, ''),
			IFNULL(sn, ''),
			IFNULL(mac_address, ''),
			IFNULL(file_id, 0),
			IFNULL(audit_status, ''),
			IFNULL(auditor_name, ''),
			audit_time,
			created_at
		FROM factory_tests
		WHERE is_deleted = 0
		  AND (? = '' OR sn = ?)
		  AND (? = '' OR mac_address = ?)
		ORDER BY id DESC
	`, sn, sn, mac, mac)
	if err != nil {
		http.Error(w, "查询出厂测试失败: "+err.Error(), http.StatusInternalServerError)
		return
	}
	data["factoryTests"] = factoryTests

	inventoryDevices, err := queryTraceList(`
		SELECT
			id,
			IFNULL(project_id, 0),
			IFNULL(device_type, ''),
			IFNULL(product_name, ''),
			IFNULL(product_model, ''),
			sn,
			mac_address,
			IFNULL(hardware_version, ''),
			IFNULL(software_version, ''),
			IFNULL(inventory_status, ''),
			in_time,
			update_time
		FROM inventory_devices
		WHERE is_deleted = 0
		  AND (? = '' OR sn = ?)
		  AND (? = '' OR mac_address = ?)
		ORDER BY id DESC
	`, sn, sn, mac, mac)
	if err != nil {
		http.Error(w, "查询库存记录失败: "+err.Error(), http.StatusInternalServerError)
		return
	}
	data["inventoryDevices"] = inventoryDevices

	shippingDevices, err := queryTraceList(`
		SELECT
			sbd.id,
			sbd.batch_id,
			IFNULL(sb.batch_no, ''),
			sbd.inventory_device_id,
			IFNULL(sbd.sn, ''),
			IFNULL(sbd.mac_address, ''),
			IFNULL(sbd.device_type, ''),
			IFNULL(sbd.hardware_version, ''),
			IFNULL(sbd.software_version, ''),
			IFNULL(sb.audit_status, ''),
			sb.audit_time,
			sbd.created_at
		FROM shipping_batch_devices sbd
		LEFT JOIN shipping_batches sb ON sbd.batch_id = sb.id
		WHERE sbd.is_deleted = 0
		  AND (? = '' OR sbd.sn = ?)
		  AND (? = '' OR sbd.mac_address = ?)
		ORDER BY sbd.id DESC
	`, sn, sn, mac, mac)
	if err != nil {
		http.Error(w, "查询发货批次设备失败: "+err.Error(), http.StatusInternalServerError)
		return
	}
	data["shippingBatchDevices"] = shippingDevices

	outboundRecords, err := queryTraceList(`
		SELECT
			id,
			batch_id,
			inventory_device_id,
			IFNULL(sn, ''),
			IFNULL(mac_address, ''),
			outbound_time,
			IFNULL(outbound_user_name, ''),
			IFNULL(status, ''),
			created_at
		FROM outbound_records
		WHERE is_deleted = 0
		  AND (? = '' OR sn = ?)
		  AND (? = '' OR mac_address = ?)
		ORDER BY id DESC
	`, sn, sn, mac, mac)
	if err != nil {
		http.Error(w, "查询出库记录失败: "+err.Error(), http.StatusInternalServerError)
		return
	}
	data["outboundRecords"] = outboundRecords

	repairRecords, err := queryTraceList(`
		SELECT
			id,
			IFNULL(project_id, 0),
			IFNULL(device_type, ''),
			IFNULL(inventory_device_id, 0),
			IFNULL(sn, ''),
			IFNULL(mac_address, ''),
			IFNULL(fault_desc, ''),
			IFNULL(repair_user_name, ''),
			repair_time,
			repair_finish_time,
			IFNULL(repair_method, ''),
			IFNULL(confirm_status, ''),
			created_at
		FROM repair_records
		WHERE is_deleted = 0
		  AND (? = '' OR sn = ?)
		  AND (? = '' OR mac_address = ?)
		ORDER BY id DESC
	`, sn, sn, mac, mac)
	if err != nil {
		http.Error(w, "查询维修记录失败: "+err.Error(), http.StatusInternalServerError)
		return
	}
	data["repairRecords"] = repairRecords

	json.NewEncoder(w).Encode(map[string]interface{}{
		"code": 200,
		"msg":  "查询成功",
		"data": data,
	})
}

// ============================================================
// 通用查询工具：把 rows 转成 []map[string]interface{}
// ============================================================

func queryTraceList(query string, args ...interface{}) ([]map[string]interface{}, error) {
	rows, err := config.DB.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	columns, err := rows.Columns()
	if err != nil {
		return nil, err
	}

	result := make([]map[string]interface{}, 0)

	for rows.Next() {
		values := make([]interface{}, len(columns))
		valuePtrs := make([]interface{}, len(columns))

		for i := range values {
			valuePtrs[i] = &values[i]
		}

		if err := rows.Scan(valuePtrs...); err != nil {
			return nil, err
		}

		rowMap := make(map[string]interface{})

		for i, col := range columns {
			val := values[i]

			switch v := val.(type) {
			case []byte:
				rowMap[col] = string(v)
			default:
				rowMap[col] = v
			}
		}

		result = append(result, rowMap)
	}

	return result, nil
}
