package handler

import (
	"crrc_pm_backend/config"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

func hasLeaderPermission(r *http.Request) bool {
	return hasRequestRole(r, "leader") || hasRequestRole(r, "system_admin")
}

func hasProductionAuditPermission(r *http.Request) bool {
	return hasLeaderPermission(r) || hasRequestRole(r, "quality_staff")
}

func hasFactoryQualityAuditorPermission(r *http.Request) bool {
	return hasRequestRole(r, "quality_staff") || hasRequestRole(r, "system_admin")
}

func hasRequestRole(r *http.Request, roleCode string) bool {
	roles := strings.Split(r.Header.Get("X-User-Roles"), ",")
	for _, role := range roles {
		role = strings.TrimSpace(role)
		if role == roleCode {
			return true
		}
	}
	return false
}

func hasRequestPermission(r *http.Request, permissionCode string) bool {
	if hasRequestRole(r, "system_admin") {
		return true
	}
	userID, _ := currentRequestUser(r)
	if userID == 0 || strings.TrimSpace(permissionCode) == "" {
		return false
	}
	ensureUserPermissionsTable()

	var exists int
	err := config.DB.QueryRow(`
		SELECT 1
		FROM user_permissions
		WHERE user_id = ?
		  AND permission_code = ?
		LIMIT 1
	`, userID, permissionCode).Scan(&exists)
	return err == nil && exists == 1
}

func requireLeaderPermission(w http.ResponseWriter, r *http.Request) bool {
	if hasLeaderPermission(r) {
		return true
	}
	http.Error(w, "无审核权限：只有领导角色可以审核", http.StatusForbidden)
	return false
}

func requireShippingAuditPermission(w http.ResponseWriter, r *http.Request) bool {
	if hasLeaderPermission(r) || hasRequestRole(r, "shipping_auditor") {
		return true
	}
	http.Error(w, "无发货审核权限：只有领导或发货审核角色可以审核", http.StatusForbidden)
	return false
}

func requireShippingManagePermission(w http.ResponseWriter, r *http.Request) bool {
	if hasLeaderPermission(r) || hasRequestRole(r, "shipping_staff") {
		return true
	}
	http.Error(w, "无发货管理权限", http.StatusForbidden)
	return false
}

func requireProjectAssistantPermission(w http.ResponseWriter, r *http.Request) bool {
	if hasRequestRole(r, "project_assistant") || hasRequestRole(r, "system_admin") {
		return true
	}
	http.Error(w, "无项目助理权限", http.StatusForbidden)
	return false
}

func requireHardwareOwnerPermission(w http.ResponseWriter, r *http.Request) bool {
	if hasRequestRole(r, "hardware_owner") || hasRequestRole(r, "system_admin") {
		return true
	}
	http.Error(w, "无硬件负责人权限", http.StatusForbidden)
	return false
}

func requireProductionAuditPermission(w http.ResponseWriter, r *http.Request) bool {
	if hasProductionAuditPermission(r) {
		return true
	}
	http.Error(w, "无生产测试审核权限：只有领导或质量检查人员可以审核", http.StatusForbidden)
	return false
}

func requireFactoryQualityAuditorPermission(w http.ResponseWriter, r *http.Request) bool {
	if hasFactoryQualityAuditorPermission(r) {
		return true
	}
	http.Error(w, "无出厂测试审核权限：只有质量检查人员可以审核", http.StatusForbidden)
	return false
}

func currentRequestUser(r *http.Request) (int64, string) {
	var userID int64
	if headerUserID := strings.TrimSpace(r.Header.Get("X-User-Id")); headerUserID != "" {
		if id, err := strconv.ParseInt(headerUserID, 10, 64); err == nil {
			userID = id
		}
	}

	userName := strings.TrimSpace(r.Header.Get("X-User-Name"))
	if decoded, err := url.QueryUnescape(userName); err == nil {
		userName = decoded
	}
	return userID, strings.TrimSpace(userName)
}

func normalizeAuditUser(r *http.Request, userID int64, userName string) (int64, string) {
	headerUserID, headerUserName := currentRequestUser(r)
	if userID == 0 && headerUserID > 0 {
		userID = headerUserID
	}

	userName = strings.TrimSpace(userName)
	if userName == "" || userName == "领导" || userName == "当前领导" || userName == "当前用户" || userName == "质量检查人员" {
		userName = headerUserName
	}
	if userName == "" {
		userName = "审核人"
	}
	return userID, userName
}

func reviewVisibilitySQL(r *http.Request, statusExpr string, creatorIDExpr string, creatorNameExpr string, creatorRole string) string {
	if hasRequestRole(r, "system_admin") {
		return ""
	}

	if hasRequestRole(r, "leader") || hasRequestRole(r, "shipping_auditor") || hasRequestRole(r, "quality_staff") {
		return " AND IFNULL(" + statusExpr + ", '草稿') IN ('待审核', 'submitted', '已提交', '已通过', '审核通过', 'approved')"
	}

	userID, userName := currentRequestUser(r)
	ownSQL := ""
	if userID > 0 {
		ownSQL += " OR IFNULL(" + creatorIDExpr + ", 0) = " + strconv.FormatInt(userID, 10)
	}
	if userName != "" {
		ownSQL += " OR IFNULL(" + creatorNameExpr + ", '') = '" + strings.ReplaceAll(userName, "'", "''") + "'"
	}

	return " AND (IFNULL(" + statusExpr + ", '草稿') IN ('已通过', '审核通过', 'approved')" + ownSQL + ")"
}
