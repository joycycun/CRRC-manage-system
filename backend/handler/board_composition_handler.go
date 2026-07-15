package handler

import (
	"crrc_pm_backend/config"
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
)

type BoardComposition struct {
	ID            int64  `json:"id"`
	ProjectID     int64  `json:"projectId"`
	ProjectName   string `json:"projectName"`
	ProductName   string `json:"productName"`
	InboundModel  string `json:"inboundModel"`
	OutboundModel string `json:"outboundModel"`
	CreatedBy     string `json:"createdBy"`
	CreatedAt     string `json:"createdAt"`
	UpdatedAt     string `json:"updatedAt"`
}

func ensureBoardCompositionTables() {
	_, _ = config.DB.Exec(`
		CREATE TABLE IF NOT EXISTS board_compositions (
			id BIGINT PRIMARY KEY AUTO_INCREMENT,
			project_id BIGINT NOT NULL DEFAULT 0,
			project_name VARCHAR(128) NOT NULL DEFAULT '',
			product_name VARCHAR(128) NOT NULL DEFAULT '',
			inbound_model VARCHAR(128) NOT NULL DEFAULT '',
			outbound_model VARCHAR(128) NOT NULL DEFAULT '',
			created_by VARCHAR(64) NOT NULL DEFAULT '',
			is_deleted TINYINT NOT NULL DEFAULT 0,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
			KEY idx_bc_project (project_id, project_name),
			KEY idx_bc_outbound (outbound_model),
			KEY idx_bc_inbound (product_name, inbound_model)
		) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4
	`)
	_, _ = config.DB.Exec(`
		CREATE TABLE IF NOT EXISTS board_composition_deductions (
			id BIGINT PRIMARY KEY AUTO_INCREMENT,
			burn_record_id BIGINT NOT NULL DEFAULT 0,
			composition_id BIGINT NOT NULL DEFAULT 0,
			inventory_device_id BIGINT NOT NULL DEFAULT 0,
			quantity INT NOT NULL DEFAULT 0,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			KEY idx_bcd_burn_record (burn_record_id),
			KEY idx_bcd_inventory (inventory_device_id)
		) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4
	`)
}

func BoardCompositionsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	ensureBoardCompositionTables()
	switch r.Method {
	case http.MethodGet:
		GetBoardCompositionsHandler(w, r)
	case http.MethodPost:
		CreateBoardCompositionHandler(w, r)
	default:
		http.Error(w, "不支持该请求方法", http.StatusMethodNotAllowed)
	}
}

func BoardCompositionActionHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	ensureBoardCompositionTables()

	path := strings.TrimPrefix(r.URL.Path, "/api/board-compositions/")
	path = strings.Trim(path, "/")
	id, err := strconv.ParseInt(path, 10, 64)
	if err != nil || id == 0 {
		http.Error(w, "板卡组成ID错误", http.StatusBadRequest)
		return
	}

	if r.Method == http.MethodDelete {
		DeleteBoardCompositionHandler(w, r, id)
		return
	}

	http.Error(w, "不支持该请求方法", http.StatusMethodNotAllowed)
}

func canManageBoardComposition(r *http.Request) bool {
	return hasRequestRole(r, "hardware_owner") ||
		hasRequestRole(r, "system_admin") ||
		hasRequestPermission(r, "board-composition:upload") ||
		hasRequestPermission(r, "board-composition:create")
}

func canViewBoardComposition(r *http.Request) bool {
	return canManageBoardComposition(r) ||
		hasRequestPermission(r, "board-composition:view") ||
		hasRequestRole(r, "leader")
}

func GetBoardCompositionsHandler(w http.ResponseWriter, r *http.Request) {
	if !canViewBoardComposition(r) {
		http.Error(w, "无板卡组成查看权限", http.StatusForbidden)
		return
	}

	rows, err := config.DB.Query(`
		SELECT
			id,
			IFNULL(project_id, 0),
			IFNULL(project_name, ''),
			IFNULL(product_name, ''),
			IFNULL(inbound_model, ''),
			IFNULL(outbound_model, ''),
			IFNULL(created_by, ''),
			IFNULL(DATE_FORMAT(created_at, '%Y-%m-%d %H:%i:%s'), ''),
			IFNULL(DATE_FORMAT(updated_at, '%Y-%m-%d %H:%i:%s'), '')
		FROM board_compositions
		WHERE IFNULL(is_deleted, 0) = 0
		ORDER BY project_name ASC, outbound_model ASC, product_name ASC, id DESC
	`)
	if err != nil {
		http.Error(w, "查询板卡组成失败: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	list := make([]BoardComposition, 0)
	for rows.Next() {
		var item BoardComposition
		if err := rows.Scan(
			&item.ID,
			&item.ProjectID,
			&item.ProjectName,
			&item.ProductName,
			&item.InboundModel,
			&item.OutboundModel,
			&item.CreatedBy,
			&item.CreatedAt,
			&item.UpdatedAt,
		); err != nil {
			http.Error(w, "板卡组成数据解析失败: "+err.Error(), http.StatusInternalServerError)
			return
		}
		list = append(list, item)
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"code": 200,
		"msg":  "查询成功",
		"data": list,
	})
}

func CreateBoardCompositionHandler(w http.ResponseWriter, r *http.Request) {
	if !canManageBoardComposition(r) {
		http.Error(w, "无板卡组成维护权限", http.StatusForbidden)
		return
	}

	var req struct {
		Records []BoardComposition `json:"records"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "参数解析失败: "+err.Error(), http.StatusBadRequest)
		return
	}
	if len(req.Records) == 0 {
		http.Error(w, "板卡组成数据不能为空", http.StatusBadRequest)
		return
	}

	_, currentUserName := currentRequestUser(r)
	tx, err := config.DB.Begin()
	if err != nil {
		http.Error(w, "开启事务失败: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer tx.Rollback()

	count := 0
	for index, item := range req.Records {
		item.ProjectName = strings.TrimSpace(item.ProjectName)
		item.ProductName = strings.TrimSpace(item.ProductName)
		item.InboundModel = strings.TrimSpace(item.InboundModel)
		item.OutboundModel = strings.TrimSpace(item.OutboundModel)
		if item.ProjectName == "" || item.ProductName == "" || item.InboundModel == "" || item.OutboundModel == "" {
			http.Error(w, "第 "+strconv.Itoa(index+1)+" 行项目名称、产品名称、入库型号、出库型号不能为空", http.StatusBadRequest)
			return
		}
		if item.ProjectID == 0 {
			item.ProjectID = findProjectIDByName(tx, item.ProjectName)
		}
		createdBy := strings.TrimSpace(item.CreatedBy)
		if createdBy == "" {
			createdBy = currentUserName
		}

		_, err := tx.Exec(`
			INSERT INTO board_compositions (
				project_id,
				project_name,
				product_name,
				inbound_model,
				outbound_model,
				created_by,
				is_deleted,
				created_at,
				updated_at
			) VALUES (?, ?, ?, ?, ?, ?, 0, NOW(), NOW())
		`, item.ProjectID, item.ProjectName, item.ProductName, item.InboundModel, item.OutboundModel, createdBy)
		if err != nil {
			http.Error(w, "保存板卡组成失败: "+err.Error(), http.StatusInternalServerError)
			return
		}
		count++
	}

	if err := tx.Commit(); err != nil {
		http.Error(w, "提交事务失败: "+err.Error(), http.StatusInternalServerError)
		return
	}

	go BackfillBoardCompositionDeductionsForFactoryInventory()

	json.NewEncoder(w).Encode(map[string]interface{}{
		"code": 200,
		"msg":  "保存成功",
		"data": map[string]interface{}{"count": count},
	})
}

func DeleteBoardCompositionHandler(w http.ResponseWriter, r *http.Request, id int64) {
	if !canManageBoardComposition(r) && !hasRequestPermission(r, "board-composition:delete") {
		http.Error(w, "无板卡组成删除权限", http.StatusForbidden)
		return
	}

	result, err := config.DB.Exec(`
		UPDATE board_compositions
		SET is_deleted = 1,
			updated_at = NOW()
		WHERE id = ?
		  AND IFNULL(is_deleted, 0) = 0
	`, id)
	if err != nil {
		http.Error(w, "删除板卡组成失败: "+err.Error(), http.StatusInternalServerError)
		return
	}

	affected, _ := result.RowsAffected()
	if affected == 0 {
		http.Error(w, "板卡组成不存在或已删除", http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"code": 200,
		"msg":  "删除成功",
	})
}

func findProjectIDByName(tx *sql.Tx, projectName string) int64 {
	var id int64
	_ = tx.QueryRow(`
		SELECT id
		FROM projects
		WHERE project_name = ?
		  AND IFNULL(is_deleted, 0) = 0
		ORDER BY id DESC
		LIMIT 1
	`, projectName).Scan(&id)
	return id
}
