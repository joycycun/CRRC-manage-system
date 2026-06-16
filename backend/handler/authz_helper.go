package handler

import (
	"net/http"
	"strings"
)

func hasLeaderPermission(r *http.Request) bool {
	return hasRequestRole(r, "leader") || hasRequestRole(r, "system_admin")
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

func requireLeaderPermission(w http.ResponseWriter, r *http.Request) bool {
	if hasLeaderPermission(r) {
		return true
	}
	http.Error(w, "无审核权限：只有领导角色可以审核", http.StatusForbidden)
	return false
}
