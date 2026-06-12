package handler

import (
	"crrc_pm_backend/config"
	"crrc_pm_backend/model"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"time"
)

// ============================================================
// 发货批次请求结构
// ============================================================

type CreateShippingBatchRequest struct {
	BatchNo            string  `json:"batchNo"`
	ProjectID          int64   `json:"projectId"`
	ExpressNo          string  `json:"expressNo"`
	FileID             int64   `json:"fileId"`
	UploaderID         int64   `json:"uploaderId"`
	UploaderName       string  `json:"uploaderName"`
	Remark             string  `json:"remark"`
	ShippingDesc       string  `json:"shippingDesc"`
	InventoryDeviceIDs []int64 `json:"inventoryDeviceIds"`
	DeviceIDs          []int64 `json:"deviceIds"`
}

type ShippingAuditRequest struct {
	AuditorID    int64  `json:"auditorId"`
	AuditorName  string `json:"auditorName"`
	AuditStatus  string `json:"auditStatus"`
	Status       string `json:"status"`
	RejectReason string `json:"rejectReason"`
}

// ============================================================
// 发货批次入口
// GET  /api/shipping-batches
// POST /api/shipping-batches
// ============================================================

func ShippingBatchesHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	switch r.Method {
	case http.MethodGet:
		GetShippingBatchesHandler(w, r)
	case http.MethodPost:
		CreateShippingBatchHandler(w, r)
	default:
		http.Error(w, "不支持该请求方法", http.StatusMethodNotAllowed)
	}
}

// ============================================================
// 发货批次操作入口
// POST   /api/shipping-batches/{id}/submit
// POST   /api/shipping-batches/{id}/audit
// DELETE /api/shipping-batches/{id}
// ============================================================

func ShippingBatchActionHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	path := strings.TrimPrefix(r.URL.Path, "/api/shipping-batches/")
	path = strings.Trim(path, "/")
	parts := strings.Split(path, "/")

	if len(parts) < 1 || parts[0] == "" {
		http.Error(w, "缺少发货批次ID", http.StatusBadRequest)
		return
	}

	id, err := strconv.ParseInt(parts[0], 10, 64)
	if err != nil {
		http.Error(w, "发货批次ID错误", http.StatusBadRequest)
		return
	}

	if len(parts) == 1 && r.Method == http.MethodDelete {
		DeleteShippingBatchHandler(w, r, id)
		return
	}

	if len(parts) == 2 && r.Method == http.MethodPost {
		switch parts[1] {
		case "submit":
			SubmitShippingBatchHandler(w, r, id)
			return
		case "audit":
			AuditShippingBatchHandler(w, r, id)
			return
		default:
			http.Error(w, "不支持的操作", http.StatusNotFound)
			return
		}
	}

	http.Error(w, "接口不存在", http.StatusNotFound)
}

// ============================================================
// GET /api/shipping-batches
// ============================================================

func GetShippingBatchesHandler(w http.ResponseWriter, r *http.Request) {
	rows, err := config.DB.Query(`
		SELECT
			id,
			batch_no,
			IFNULL(project_id, 0),
			IFNULL(express_no, ''),
			IFNULL(device_count, 0),
			IFNULL(file_id, 0),
			IFNULL(uploader_id, 0),
			IFNULL(uploader_name, ''),
			upload_time,
			IFNULL(audit_status, ''),
			IFNULL(auditor_id, 0),
			IFNULL(auditor_name, ''),
			audit_time,
			IFNULL(remark, ''),
			IFNULL(reject_reason, ''),
			IFNULL(shipping_desc, ''),
			is_deleted,
			created_at,
			updated_at
		FROM shipping_batches
		WHERE is_deleted = 0
		ORDER BY id DESC
	`)
	if err != nil {
		http.Error(w, "查询失败: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	list := make([]map[string]interface{}, 0)

	for rows.Next() {
		var item model.ShippingBatch

		err := rows.Scan(
			&item.ID,
			&item.BatchNo,
			&item.ProjectID,
			&item.ExpressNo,
			&item.DeviceCount,
			&item.FileID,
			&item.UploaderID,
			&item.UploaderName,
			&item.UploadTime,
			&item.AuditStatus,
			&item.AuditorID,
			&item.AuditorName,
			&item.AuditTime,
			&item.Remark,
			&item.RejectReason,
			&item.ShippingDesc,
			&item.IsDeleted,
			&item.CreatedAt,
			&item.UpdatedAt,
		)
		if err != nil {
			http.Error(w, "数据解析失败: "+err.Error(), http.StatusInternalServerError)
			return
		}

		devices, err := getShippingBatchDevicesByBatchID(item.ID)
		if err != nil {
			http.Error(w, "查询批次设备明细失败: "+err.Error(), http.StatusInternalServerError)
			return
		}

		list = append(list, map[string]interface{}{
			"id":           item.ID,
			"batchNo":      item.BatchNo,
			"projectId":    item.ProjectID,
			"expressNo":    item.ExpressNo,
			"deviceCount":  item.DeviceCount,
			"fileId":       item.FileID,
			"uploaderId":   item.UploaderID,
			"uploaderName": item.UploaderName,
			"uploadTime":   item.UploadTime,
			"auditStatus":  item.AuditStatus,
			"auditorId":    item.AuditorID,
			"auditorName":  item.AuditorName,
			"auditTime":    item.AuditTime,
			"remark":       item.Remark,
			"rejectReason": item.RejectReason,
			"shippingDesc": item.ShippingDesc,
			"isDeleted":    item.IsDeleted,
			"createdAt":    item.CreatedAt,
			"updatedAt":    item.UpdatedAt,
			"deviceList":   devices,
			"devices":      devices,
		})
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"code": 200,
		"msg":  "查询成功",
		"data": list,
	})
}

// ============================================================
// POST /api/shipping-batches
// 新增批次，并锁定库存设备
// ============================================================

func CreateShippingBatchHandler(w http.ResponseWriter, r *http.Request) {
	var req CreateShippingBatchRequest

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "参数解析失败: "+err.Error(), http.StatusBadRequest)
		return
	}

	if len(req.InventoryDeviceIDs) == 0 && len(req.DeviceIDs) > 0 {
		req.InventoryDeviceIDs = req.DeviceIDs
	}

	if req.BatchNo == "" {
		http.Error(w, "发货批次号不能为空", http.StatusBadRequest)
		return
	}

	if len(req.InventoryDeviceIDs) == 0 {
		http.Error(w, "必须选择至少一台库存设备", http.StatusBadRequest)
		return
	}

	tx, err := config.DB.Begin()
	if err != nil {
		http.Error(w, "开启事务失败: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer tx.Rollback()

	now := time.Now()

	result, err := tx.Exec(`
		INSERT INTO shipping_batches (
			batch_no,
			project_id,
			express_no,
			device_count,
			file_id,
			uploader_id,
			uploader_name,
			upload_time,
			audit_status,
			remark,
			shipping_desc,
			is_deleted,
			created_at,
			updated_at
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, '草稿', ?, ?, 0, ?, ?)
	`,
		req.BatchNo,
		req.ProjectID,
		req.ExpressNo,
		len(req.InventoryDeviceIDs),
		req.FileID,
		req.UploaderID,
		req.UploaderName,
		now,
		req.Remark,
		req.ShippingDesc,
		now,
		now,
	)

	if err != nil {
		http.Error(w, "新增发货批次失败: "+err.Error(), http.StatusInternalServerError)
		return
	}

	batchID, _ := result.LastInsertId()

	for _, inventoryID := range req.InventoryDeviceIDs {
		var sn string
		var mac string
		var deviceType string
		var hardwareVersion string
		var softwareVersion string
		var inventoryStatus string

		err = tx.QueryRow(`
			SELECT
				sn,
				mac_address,
				IFNULL(device_type, ''),
				IFNULL(hardware_version, ''),
				IFNULL(software_version, ''),
				IFNULL(inventory_status, '')
			FROM inventory_devices
			WHERE id = ? AND is_deleted = 0
		`, inventoryID).Scan(
			&sn,
			&mac,
			&deviceType,
			&hardwareVersion,
			&softwareVersion,
			&inventoryStatus,
		)

		if err != nil {
			http.Error(w, "读取库存设备失败: "+err.Error(), http.StatusInternalServerError)
			return
		}

		if inventoryStatus != "在库" {
			http.Error(w, "库存设备不是在库状态，不能加入发货批次", http.StatusBadRequest)
			return
		}

		_, err = tx.Exec(`
			INSERT INTO shipping_batch_devices (
				batch_id,
				inventory_device_id,
				sn,
				mac_address,
				device_type,
				hardware_version,
				software_version,
				is_deleted,
				created_at
			) VALUES (?, ?, ?, ?, ?, ?, ?, 0, ?)
		`,
			batchID,
			inventoryID,
			sn,
			mac,
			deviceType,
			hardwareVersion,
			softwareVersion,
			now,
		)

		if err != nil {
			http.Error(w, "保存发货设备明细失败: "+err.Error(), http.StatusInternalServerError)
			return
		}

		_, err = tx.Exec(`
			UPDATE inventory_devices
			SET inventory_status = '已锁定', update_time = NOW()
			WHERE id = ? AND is_deleted = 0
		`, inventoryID)

		if err != nil {
			http.Error(w, "锁定库存设备失败: "+err.Error(), http.StatusInternalServerError)
			return
		}
	}

	if err = tx.Commit(); err != nil {
		http.Error(w, "提交事务失败: "+err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"code": 200,
		"msg":  "新增成功",
		"data": map[string]interface{}{
			"id": batchID,
		},
	})
}

// ============================================================
// POST /api/shipping-batches/{id}/submit
// ============================================================

func SubmitShippingBatchHandler(w http.ResponseWriter, r *http.Request, id int64) {
	result, err := config.DB.Exec(`
		UPDATE shipping_batches
		SET audit_status = '待审核', updated_at = NOW()
		WHERE id = ? AND is_deleted = 0
	`, id)

	if err != nil {
		http.Error(w, "提交失败: "+err.Error(), http.StatusInternalServerError)
		return
	}

	affected, _ := result.RowsAffected()
	if affected == 0 {
		http.Error(w, "发货批次不存在或已删除", http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"code": 200,
		"msg":  "提交成功",
	})
}

// ============================================================
// POST /api/shipping-batches/{id}/audit
// 审核通过：自动出库
// 审核驳回：释放库存为在库
// ============================================================

func AuditShippingBatchHandler(w http.ResponseWriter, r *http.Request, id int64) {
	var req ShippingAuditRequest

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "参数解析失败: "+err.Error(), http.StatusBadRequest)
		return
	}

	if req.AuditStatus == "" {
		req.AuditStatus = req.Status
	}
	switch req.AuditStatus {
	case "approved":
		req.AuditStatus = "已通过"
	case "rejected":
		req.AuditStatus = "已驳回"
	}

	if req.AuditStatus != "已通过" && req.AuditStatus != "已驳回" {
		http.Error(w, "审核状态只能是 已通过 或 已驳回", http.StatusBadRequest)
		return
	}

	tx, err := config.DB.Begin()
	if err != nil {
		http.Error(w, "开启事务失败: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer tx.Rollback()

	result, err := tx.Exec(`
		UPDATE shipping_batches
		SET
			audit_status = ?,
			auditor_id = ?,
			auditor_name = ?,
			audit_time = NOW(),
			reject_reason = ?,
			updated_at = NOW()
		WHERE id = ? AND is_deleted = 0
	`,
		req.AuditStatus,
		req.AuditorID,
		req.AuditorName,
		req.RejectReason,
		id,
	)

	if err != nil {
		http.Error(w, "审核失败: "+err.Error(), http.StatusInternalServerError)
		return
	}

	affected, _ := result.RowsAffected()
	if affected == 0 {
		http.Error(w, "发货批次不存在或已删除", http.StatusNotFound)
		return
	}

	if req.AuditStatus == "已通过" {
		rows, err := tx.Query(`
			SELECT
				inventory_device_id,
				IFNULL(sn, ''),
				IFNULL(mac_address, '')
			FROM shipping_batch_devices
			WHERE batch_id = ? AND is_deleted = 0
		`, id)

		if err != nil {
			http.Error(w, "读取发货设备失败: "+err.Error(), http.StatusInternalServerError)
			return
		}

		type shippingDevice struct {
			inventoryDeviceID int64
			sn                string
			mac               string
		}
		devices := make([]shippingDevice, 0)
		for rows.Next() {
			var device shippingDevice

			err = rows.Scan(&device.inventoryDeviceID, &device.sn, &device.mac)
			if err != nil {
				rows.Close()
				http.Error(w, "解析发货设备失败: "+err.Error(), http.StatusInternalServerError)
				return
			}
			devices = append(devices, device)
		}
		if err = rows.Err(); err != nil {
			rows.Close()
			http.Error(w, "读取发货设备失败: "+err.Error(), http.StatusInternalServerError)
			return
		}
		rows.Close()

		if len(devices) == 0 {
			http.Error(w, "该发货批次没有设备明细，不能出库", http.StatusBadRequest)
			return
		}

		for _, device := range devices {
			var outboundExists int
			err = tx.QueryRow(`
				SELECT COUNT(1)
				FROM outbound_records
				WHERE batch_id = ?
				  AND inventory_device_id = ?
				  AND IFNULL(is_deleted, 0) = 0
			`, id, device.inventoryDeviceID).Scan(&outboundExists)

			if err != nil {
				http.Error(w, "检查出库记录失败: "+err.Error(), http.StatusInternalServerError)
				return
			}

			if outboundExists == 0 {
				_, err = tx.Exec(`
					INSERT INTO outbound_records (
						batch_id,
						inventory_device_id,
						sn,
						mac_address,
						outbound_time,
						outbound_user_id,
						outbound_user_name,
						status,
						remark,
						is_deleted,
						created_at
					) VALUES (?, ?, ?, ?, NOW(), ?, ?, '已出库', '', 0, NOW())
				`,
					id,
					device.inventoryDeviceID,
					device.sn,
					device.mac,
					req.AuditorID,
					req.AuditorName,
				)

				if err != nil {
					http.Error(w, "生成出库记录失败: "+err.Error(), http.StatusInternalServerError)
					return
				}
			}

			_, err = tx.Exec(`
				UPDATE inventory_devices
				SET inventory_status = '已出库', update_time = NOW()
				WHERE id = ? AND is_deleted = 0
			`, device.inventoryDeviceID)

			if err != nil {
				http.Error(w, "更新库存出库状态失败: "+err.Error(), http.StatusInternalServerError)
				return
			}
		}
	}

	if req.AuditStatus == "已驳回" {
		_, err = tx.Exec(`
			UPDATE inventory_devices inv
			JOIN shipping_batch_devices sbd ON inv.id = sbd.inventory_device_id
			SET inv.inventory_status = '在库', inv.update_time = NOW()
			WHERE sbd.batch_id = ? AND sbd.is_deleted = 0
		`, id)

		if err != nil {
			http.Error(w, "释放库存失败: "+err.Error(), http.StatusInternalServerError)
			return
		}
	}

	if err = tx.Commit(); err != nil {
		http.Error(w, "提交事务失败: "+err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"code": 200,
		"msg":  "审核成功",
	})
}

// ============================================================
// DELETE /api/shipping-batches/{id}
// 删除草稿/未通过批次时释放已锁定库存
// ============================================================

func DeleteShippingBatchHandler(w http.ResponseWriter, r *http.Request, id int64) {
	tx, err := config.DB.Begin()
	if err != nil {
		http.Error(w, "开启事务失败: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer tx.Rollback()

	_, err = tx.Exec(`
		UPDATE inventory_devices inv
		JOIN shipping_batch_devices sbd ON inv.id = sbd.inventory_device_id
		SET inv.inventory_status = '在库', inv.update_time = NOW()
		WHERE sbd.batch_id = ?
		  AND sbd.is_deleted = 0
		  AND inv.inventory_status = '已锁定'
	`, id)

	if err != nil {
		http.Error(w, "释放库存失败: "+err.Error(), http.StatusInternalServerError)
		return
	}

	_, err = tx.Exec(`
		UPDATE shipping_batch_devices
		SET is_deleted = 1
		WHERE batch_id = ?
	`, id)

	if err != nil {
		http.Error(w, "删除批次设备失败: "+err.Error(), http.StatusInternalServerError)
		return
	}

	result, err := tx.Exec(`
		UPDATE shipping_batches
		SET is_deleted = 1, updated_at = NOW()
		WHERE id = ?
	`, id)

	if err != nil {
		http.Error(w, "删除发货批次失败: "+err.Error(), http.StatusInternalServerError)
		return
	}

	affected, _ := result.RowsAffected()
	if affected == 0 {
		http.Error(w, "发货批次不存在", http.StatusNotFound)
		return
	}

	if err = tx.Commit(); err != nil {
		http.Error(w, "提交事务失败: "+err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"code": 200,
		"msg":  "删除成功",
	})
}

// ============================================================
// 发货批次设备明细
// GET /api/shipping-batch-devices?batchId=1
// ============================================================

func ShippingBatchDevicesHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	batchIDStr := r.URL.Query().Get("batchId")
	if batchIDStr == "" {
		http.Error(w, "batchId不能为空", http.StatusBadRequest)
		return
	}

	batchID, err := strconv.ParseInt(batchIDStr, 10, 64)
	if err != nil {
		http.Error(w, "batchId错误", http.StatusBadRequest)
		return
	}

	rows, err := config.DB.Query(`
		SELECT
			id,
			batch_id,
			inventory_device_id,
			IFNULL(sn, ''),
			IFNULL(mac_address, ''),
			IFNULL(device_type, ''),
			IFNULL(hardware_version, ''),
			IFNULL(software_version, ''),
			is_deleted,
			created_at
		FROM shipping_batch_devices
		WHERE batch_id = ? AND is_deleted = 0
		ORDER BY id ASC
	`, batchID)

	if err != nil {
		http.Error(w, "查询失败: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	list := make([]model.ShippingBatchDevice, 0)

	for rows.Next() {
		var item model.ShippingBatchDevice

		err := rows.Scan(
			&item.ID,
			&item.BatchID,
			&item.InventoryDeviceID,
			&item.SN,
			&item.MacAddress,
			&item.DeviceType,
			&item.HardwareVersion,
			&item.SoftwareVersion,
			&item.IsDeleted,
			&item.CreatedAt,
		)

		if err != nil {
			http.Error(w, "数据解析失败: "+err.Error(), http.StatusInternalServerError)
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

// ============================================================
// 出库记录
// GET /api/outbound-records
// ============================================================

func OutboundRecordsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodGet {
		http.Error(w, "不支持该请求方法", http.StatusMethodNotAllowed)
		return
	}

	rows, err := config.DB.Query(`
		SELECT
			o.id,
			o.batch_id,
			o.inventory_device_id,
			IFNULL(sbd.device_type, IFNULL(inv.device_type, '')),
			IFNULL(o.sn, ''),
			IFNULL(o.mac_address, ''),
			IFNULL(sbd.software_version, IFNULL(inv.software_version, '')),
			IFNULL(sbd.hardware_version, IFNULL(inv.hardware_version, '')),
			IFNULL(DATE_FORMAT(inv.in_time, '%Y-%m-%d %H:%i:%s'), ''),
			IFNULL(DATE_FORMAT(o.outbound_time, '%Y-%m-%d %H:%i:%s'), ''),
			IFNULL(o.outbound_user_name, ''),
			IFNULL(sb.batch_no, ''),
			IFNULL(sb.express_no, ''),
			IFNULL(o.status, ''),
			IFNULL(o.remark, '')
		FROM outbound_records o
		LEFT JOIN shipping_batches sb ON sb.id = o.batch_id
		LEFT JOIN shipping_batch_devices sbd
			ON sbd.batch_id = o.batch_id
			AND sbd.inventory_device_id = o.inventory_device_id
			AND IFNULL(sbd.is_deleted, 0) = 0
		LEFT JOIN inventory_devices inv ON inv.id = o.inventory_device_id
		WHERE IFNULL(o.is_deleted, 0) = 0
		ORDER BY o.id DESC
	`)

	if err != nil {
		http.Error(w, "查询失败: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	type OutboundRecordVO struct {
		ID                int64  `json:"id"`
		OutboundID        int64  `json:"outboundId"`
		BatchID           int64  `json:"batchId"`
		InventoryDeviceID int64  `json:"inventoryDeviceId"`
		DeviceType        string `json:"deviceType"`
		SN                string `json:"sn"`
		MacAddress        string `json:"macAddress"`
		SoftwareVersion   string `json:"softwareVersion"`
		HardwareVersion   string `json:"hardwareVersion"`
		InTime            string `json:"inTime"`
		OutboundTime      string `json:"outboundTime"`
		Operator          string `json:"operator"`
		OutboundUserName  string `json:"outboundUserName"`
		ShippingBatchNo   string `json:"shippingBatchNo"`
		BatchNo           string `json:"batchNo"`
		ExpressNo         string `json:"expressNo"`
		OutboundStatus    string `json:"outboundStatus"`
		Status            string `json:"status"`
		Remark            string `json:"remark"`
	}

	list := make([]OutboundRecordVO, 0)

	for rows.Next() {
		var item OutboundRecordVO

		err := rows.Scan(
			&item.ID,
			&item.BatchID,
			&item.InventoryDeviceID,
			&item.DeviceType,
			&item.SN,
			&item.MacAddress,
			&item.SoftwareVersion,
			&item.HardwareVersion,
			&item.InTime,
			&item.OutboundTime,
			&item.OutboundUserName,
			&item.ShippingBatchNo,
			&item.ExpressNo,
			&item.Status,
			&item.Remark,
		)

		if err != nil {
			http.Error(w, "数据解析失败: "+err.Error(), http.StatusInternalServerError)
			return
		}

		item.OutboundID = item.ID
		item.Operator = item.OutboundUserName
		item.BatchNo = item.ShippingBatchNo
		item.OutboundStatus = item.Status
		list = append(list, item)
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"code": 200,
		"msg":  "查询成功",
		"data": list,
	})
}
