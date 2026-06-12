package handler

import (
	"crrc_pm_backend/config"
	"crrc_pm_backend/model"
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"time"
)

// ============================================================
// 生产工单
// ============================================================

func ProductionOrdersHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	switch r.Method {
	case http.MethodGet:
		GetProductionOrdersHandler(w, r)
	case http.MethodPost:
		CreateProductionOrderHandler(w, r)
	default:
		http.Error(w, "不支持该请求方法", http.StatusMethodNotAllowed)
	}
}

func ProductionOrderActionHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	path := strings.TrimPrefix(r.URL.Path, "/api/production-orders/")
	path = strings.Trim(path, "/")
	parts := strings.Split(path, "/")

	if len(parts) < 1 || parts[0] == "" {
		http.Error(w, "缺少生产工单ID", http.StatusBadRequest)
		return
	}

	id, err := strconv.ParseInt(parts[0], 10, 64)
	if err != nil {
		http.Error(w, "生产工单ID错误", http.StatusBadRequest)
		return
	}

	if len(parts) == 1 && r.Method == http.MethodPut {
		UpdateProductionOrderHandler(w, r, id)
		return
	}

	if len(parts) == 1 && r.Method == http.MethodDelete {
		DeleteProductionOrderHandler(w, r, id)
		return
	}

	if len(parts) == 2 && r.Method == http.MethodPost && parts[1] == "close" {
		CloseProductionOrderHandler(w, r, id)
		return
	}

	http.Error(w, "接口不存在", http.StatusNotFound)
}

func GetProductionOrdersHandler(w http.ResponseWriter, r *http.Request) {
	rows, err := config.DB.Query(`
		SELECT
			id,
			order_no,
			IFNULL(project_id, 0),
			IFNULL(device_type, ''),
			IFNULL(product_name, ''),
			IFNULL(product_model, ''),
			IFNULL(plan_qty, 0),
			IFNULL(hardware_id, 0),
			IFNULL(hardware_version, ''),
			IFNULL(software_id, 0),
			IFNULL(software_version, ''),
			IFNULL(status, ''),
			IFNULL(create_user_id, 0),
			IFNULL(create_user_name, ''),
			create_time,
			IFNULL(remark, ''),
			updated_at,
			is_deleted
		FROM production_orders
		WHERE is_deleted = 0
		ORDER BY id DESC
	`)
	if err != nil {
		http.Error(w, "查询失败: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	list := make([]model.ProductionOrder, 0)

	for rows.Next() {
		var item model.ProductionOrder

		err := rows.Scan(
			&item.ID,
			&item.OrderNo,
			&item.ProjectID,
			&item.DeviceType,
			&item.ProductName,
			&item.ProductModel,
			&item.PlanQty,
			&item.HardwareID,
			&item.HardwareVersion,
			&item.SoftwareID,
			&item.SoftwareVersion,
			&item.Status,
			&item.CreateUserID,
			&item.CreateUserName,
			&item.CreateTime,
			&item.Remark,
			&item.UpdatedAt,
			&item.IsDeleted,
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

func CreateProductionOrderHandler(w http.ResponseWriter, r *http.Request) {
	var item model.ProductionOrder

	err := json.NewDecoder(r.Body).Decode(&item)
	if err != nil {
		http.Error(w, "参数解析失败: "+err.Error(), http.StatusBadRequest)
		return
	}

	if item.OrderNo == "" {
		http.Error(w, "生产工单号不能为空", http.StatusBadRequest)
		return
	}

	if item.DeviceType == "" {
		http.Error(w, "终端类型不能为空", http.StatusBadRequest)
		return
	}

	if item.Status == "" {
		item.Status = "待生产"
	}

	now := time.Now()

	result, err := config.DB.Exec(`
		INSERT INTO production_orders (
			order_no,
			project_id,
			device_type,
			product_name,
			product_model,
			plan_qty,
			hardware_id,
			hardware_version,
			software_id,
			software_version,
			status,
			create_user_id,
			create_user_name,
			create_time,
			remark,
			updated_at,
			is_deleted
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, 0)
	`,
		item.OrderNo,
		item.ProjectID,
		item.DeviceType,
		item.ProductName,
		item.ProductModel,
		item.PlanQty,
		item.HardwareID,
		item.HardwareVersion,
		item.SoftwareID,
		item.SoftwareVersion,
		item.Status,
		item.CreateUserID,
		item.CreateUserName,
		now,
		item.Remark,
		now,
	)

	if err != nil {
		http.Error(w, "新增失败: "+err.Error(), http.StatusInternalServerError)
		return
	}

	id, _ := result.LastInsertId()

	json.NewEncoder(w).Encode(map[string]interface{}{
		"code": 200,
		"msg":  "新增成功",
		"data": map[string]interface{}{"id": id},
	})
}

func UpdateProductionOrderHandler(w http.ResponseWriter, r *http.Request, id int64) {
	var item model.ProductionOrder

	err := json.NewDecoder(r.Body).Decode(&item)
	if err != nil {
		http.Error(w, "参数解析失败: "+err.Error(), http.StatusBadRequest)
		return
	}

	if item.OrderNo == "" {
		http.Error(w, "生产工单号不能为空", http.StatusBadRequest)
		return
	}

	result, err := config.DB.Exec(`
		UPDATE production_orders
		SET
			order_no = ?,
			project_id = ?,
			device_type = ?,
			product_name = ?,
			product_model = ?,
			plan_qty = ?,
			hardware_id = ?,
			hardware_version = ?,
			software_id = ?,
			software_version = ?,
			status = ?,
			remark = ?,
			updated_at = NOW()
		WHERE id = ? AND is_deleted = 0
	`,
		item.OrderNo,
		item.ProjectID,
		item.DeviceType,
		item.ProductName,
		item.ProductModel,
		item.PlanQty,
		item.HardwareID,
		item.HardwareVersion,
		item.SoftwareID,
		item.SoftwareVersion,
		item.Status,
		item.Remark,
		id,
	)

	if err != nil {
		http.Error(w, "修改失败: "+err.Error(), http.StatusInternalServerError)
		return
	}

	affected, _ := result.RowsAffected()
	if affected == 0 {
		http.Error(w, "生产工单不存在或已删除", http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"code": 200,
		"msg":  "修改成功",
	})
}

func CloseProductionOrderHandler(w http.ResponseWriter, r *http.Request, id int64) {
	result, err := config.DB.Exec(`
		UPDATE production_orders
		SET status = '已关闭', updated_at = NOW()
		WHERE id = ? AND is_deleted = 0
	`, id)

	if err != nil {
		http.Error(w, "关闭失败: "+err.Error(), http.StatusInternalServerError)
		return
	}

	affected, _ := result.RowsAffected()
	if affected == 0 {
		http.Error(w, "生产工单不存在或已删除", http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"code": 200,
		"msg":  "关闭成功",
	})
}

func DeleteProductionOrderHandler(w http.ResponseWriter, r *http.Request, id int64) {
	result, err := config.DB.Exec(`
		UPDATE production_orders
		SET is_deleted = 1, updated_at = NOW()
		WHERE id = ?
	`, id)

	if err != nil {
		http.Error(w, "删除失败: "+err.Error(), http.StatusInternalServerError)
		return
	}

	affected, _ := result.RowsAffected()
	if affected == 0 {
		http.Error(w, "生产工单不存在", http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"code": 200,
		"msg":  "删除成功",
	})
}

// ============================================================
// 烧录记录
// ============================================================

func BurnRecordsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	switch r.Method {
	case http.MethodGet:
		GetBurnRecordsHandler(w, r)
	case http.MethodPost:
		CreateBurnRecordHandler(w, r)
	default:
		http.Error(w, "不支持该请求方法", http.StatusMethodNotAllowed)
	}
}

func BurnRecordActionHandler(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimPrefix(r.URL.Path, "/api/burn-records/")
	parts := strings.Split(path, "/")

	if len(parts) < 1 || parts[0] == "" {
		http.NotFound(w, r)
		return
	}
	if len(parts) == 1 && parts[0] == "options" && r.Method == http.MethodGet {
		GetBurnRecordOptionsHandler(w, r)
		return
	}
	if len(parts) == 1 && parts[0] == "import" && r.Method == http.MethodPost {
		ImportBurnRecordsHandler(w, r)
		return
	}

	if len(parts) == 2 && parts[0] == "batch" && r.Method == http.MethodDelete {
		DeleteBurnBatchHandler(w, r, parts[1])
		return
	}

	if len(parts) == 1 && r.Method == http.MethodDelete {
		id, err := strconv.ParseInt(parts[0], 10, 64)
		if err != nil {
			http.Error(w, "烧录记录ID错误", http.StatusBadRequest)
			return
		}

		DeleteBurnRecordHandler(w, r, id)
		return
	}

	http.NotFound(w, r)
}
func GetBurnRecordsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	rows, err := config.DB.Query(`
		SELECT
			id,
			IFNULL(batch_no, ''),
			IFNULL(project_id, 0),
			IFNULL(production_order_id, 0),
			IFNULL(product_name, ''),
			IFNULL(product_model, ''),
			IFNULL(product_code, ''),
			IFNULL(device_type, ''),
			IFNULL(sn, ''),
			IFNULL(mac_address, ''),
			IFNULL(hardware_id, 0),
			IFNULL(hardware_version, ''),
			IFNULL(software_id, 0),
			IFNULL(software_version, ''),
			IFNULL(pcb_qr_code, ''),
			IFNULL(note, ''),
			IFNULL(source_file_id, 0),
			IFNULL(uploader_id, 0),
			IFNULL(uploader_name, ''),
			IFNULL(DATE_FORMAT(upload_time, '%Y-%m-%d'), ''),
			IFNULL(burn_desc, ''),
			IFNULL(DATE_FORMAT(created_at, '%Y-%m-%d %H:%i:%s'), ''),
			IFNULL(DATE_FORMAT(updated_at, '%Y-%m-%d %H:%i:%s'), '')
		FROM burn_records
		WHERE IFNULL(is_deleted, 0) = 0
		ORDER BY id DESC
	`)
	if err != nil {
		http.Error(w, "查询失败: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	type BurnRecordVO struct {
		ID int64 `json:"id"`

		BatchNo string `json:"batchNo"`

		ProjectID         int64 `json:"projectId"`
		ProductionOrderID int64 `json:"productionOrderId"`

		ProductName  string `json:"productName"`
		ProductModel string `json:"productModel"`
		ProductCode  string `json:"productCode"`
		DeviceType   string `json:"deviceType"`

		SerialNumber string `json:"serialNumber"`
		SN           string `json:"sn"`

		MacAddress string `json:"macAddress"`

		HardwareID      int64  `json:"hardwareId"`
		HardwareVersion string `json:"hardwareVersion"`

		SoftwareID      int64  `json:"softwareId"`
		SoftwareVersion string `json:"softwareVersion"`

		PcbQrCode string `json:"pcbQrCode"`
		PcbQRCode string `json:"pcbQRCode"`

		Note string `json:"note"`

		SourceFileID int64  `json:"sourceFileId"`
		FileName     string `json:"fileName"`
		FileURL      string `json:"fileUrl"`

		UploaderID   int64  `json:"uploaderId"`
		UploaderName string `json:"uploaderName"`
		Uploader     string `json:"uploader"`

		UploadTime string `json:"uploadTime"`

		BurnDesc     string `json:"burnDesc"`
		ImportRemark string `json:"importRemark"`

		CreatedAt string `json:"createdAt"`
		UpdatedAt string `json:"updatedAt"`
	}

	list := make([]BurnRecordVO, 0)

	for rows.Next() {
		var item BurnRecordVO

		err := rows.Scan(
			&item.ID,
			&item.BatchNo,
			&item.ProjectID,
			&item.ProductionOrderID,
			&item.ProductName,
			&item.ProductModel,
			&item.ProductCode,
			&item.DeviceType,
			&item.SerialNumber,
			&item.MacAddress,
			&item.HardwareID,
			&item.HardwareVersion,
			&item.SoftwareID,
			&item.SoftwareVersion,
			&item.PcbQrCode,
			&item.Note,
			&item.SourceFileID,
			&item.UploaderID,
			&item.UploaderName,
			&item.UploadTime,
			&item.BurnDesc,
			&item.CreatedAt,
			&item.UpdatedAt,
		)
		if err != nil {
			http.Error(w, "数据解析失败: "+err.Error(), http.StatusInternalServerError)
			return
		}

		item.SN = item.SerialNumber
		item.PcbQRCode = item.PcbQrCode
		item.Uploader = item.UploaderName
		item.ImportRemark = item.BurnDesc

		// 现在还没有真实文件上传表，所以这里先给前端占位
		item.FileName = "暂无源文件名称"
		item.FileURL = ""

		if item.MacAddress == "" {
			item.MacAddress = "-"
		}

		if item.HardwareVersion == "" {
			item.HardwareVersion = "-"
		}

		if item.SoftwareVersion == "" {
			item.SoftwareVersion = "-"
		}

		if item.PcbQrCode == "" {
			item.PcbQrCode = "-"
			item.PcbQRCode = "-"
		}

		list = append(list, item)
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"code": 200,
		"msg":  "查询成功",
		"data": list,
	})
}

func CreateBurnRecordHandler(w http.ResponseWriter, r *http.Request) {
	var item model.BurnRecord

	err := json.NewDecoder(r.Body).Decode(&item)
	if err != nil {
		http.Error(w, "参数解析失败: "+err.Error(), http.StatusBadRequest)
		return
	}

	if item.SN == "" {
		http.Error(w, "SN不能为空", http.StatusBadRequest)
		return
	}

	if item.MacAddress == "" {
		http.Error(w, "MAC地址不能为空", http.StatusBadRequest)
		return
	}

	now := time.Now()

	result, err := config.DB.Exec(`
		INSERT INTO burn_records (
			batch_no,
			project_id,
			production_order_id,
			product_name,
			product_model,
			product_code,
			device_type,
			sn,
			mac_address,
			hardware_id,
			hardware_version,
			software_id,
			software_version,
			pcb_qr_code,
			note,
			source_file_id,
			uploader_id,
			uploader_name,
			upload_time,
			burn_desc,
			is_deleted,
			created_at,
			updated_at
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, 0, ?, ?)
	`,
		item.BatchNo,
		item.ProjectID,
		item.ProductionOrderID,
		item.ProductName,
		item.ProductModel,
		item.ProductCode,
		item.DeviceType,
		item.SN,
		item.MacAddress,
		item.HardwareID,
		item.HardwareVersion,
		item.SoftwareID,
		item.SoftwareVersion,
		item.PCBQRCode,
		item.Note,
		item.SourceFileID,
		item.UploaderID,
		item.UploaderName,
		now,
		item.BurnDesc,
		now,
		now,
	)

	if err != nil {
		http.Error(w, "新增失败: "+err.Error(), http.StatusInternalServerError)
		return
	}

	id, _ := result.LastInsertId()

	json.NewEncoder(w).Encode(map[string]interface{}{
		"code": 200,
		"msg":  "新增成功",
		"data": map[string]interface{}{"id": id},
	})
}

// func (w http.ResponseWriter, r *http.Request, id int64) {
// 	result, err := config.DB.Exec(`
// 		UPDATE burn_records
// 		SET is_deleted = 1, updated_at = NOW()
// 		WHERE id = ?
// 	`, id)

// 	if err != nil {
// 		http.Error(w, "删除失败: "+err.Error(), http.StatusInternalServerError)
// 		return
// 	}

// 	affected, _ := result.RowsAffected()
// 	if affected == 0 {
// 		http.Error(w, "烧录记录不存在", http.StatusNotFound)
// 		return
// 	}

// 	json.NewEncoder(w).Encode(map[string]interface{}{
// 		"code": 200,
// 		"msg":  "删除成功",
// 	})
// }

// ============================================================
// 出厂测试
// ============================================================

func FactoryTestsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	switch r.Method {
	case http.MethodGet:
		GetFactoryTestsHandler(w, r)
	case http.MethodPost:
		CreateFactoryTestHandler(w, r)
	default:
		http.Error(w, "不支持该请求方法", http.StatusMethodNotAllowed)
	}
}

func FactoryTestActionHandler(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimPrefix(r.URL.Path, "/api/factory-tests/")
	parts := strings.Split(path, "/")

	if len(parts) < 1 || parts[0] == "" {
		http.NotFound(w, r)
		return
	}

	// 保险：即使 router.go 没有精确匹配，这里也先拦截 import/delete/submit/audit
	if len(parts) == 1 && parts[0] == "import" && r.Method == http.MethodPost {
		ImportFactoryTestsHandler(w, r)
		return
	}

	if len(parts) == 1 && parts[0] == "delete" && r.Method == http.MethodPost {
		DeleteFactoryTestsHandler(w, r)
		return
	}

	if len(parts) == 1 && parts[0] == "submit" && r.Method == http.MethodPost {
		SubmitFactoryTestsHandler(w, r)
		return
	}

	if len(parts) == 1 && parts[0] == "audit" && r.Method == http.MethodPost {
		AuditFactoryTestsHandler(w, r)
		return
	}

	// 下面才允许把路径当 ID
	if len(parts) == 1 && r.Method == http.MethodDelete {
		id, err := strconv.ParseInt(parts[0], 10, 64)
		if err != nil {
			http.Error(w, "出厂测试ID错误", http.StatusBadRequest)
			return
		}

		DeleteFactoryTestByIDHandler(w, r, id)
		return
	}

	http.NotFound(w, r)
}
func GetFactoryTestsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	rows, err := config.DB.Query(`
		SELECT
			id,
			IFNULL(burn_record_id, 0),
			IFNULL(project_id, 0),
			IFNULL(product_model, ''),
			IFNULL(device_type, ''),
			IFNULL(mac_address, ''),
			IFNULL(sn, ''),
			IFNULL(file_id, 0),
			IFNULL(uploader_id, 0),
			IFNULL(uploader_name, ''),
			IFNULL(DATE_FORMAT(upload_time, '%Y-%m-%d'), ''),
			CASE
				WHEN audit_status IN ('草稿', 'draft') THEN 'draft'
				WHEN audit_status IN ('待审核', 'submitted') THEN 'submitted'
				WHEN audit_status IN ('审核通过', '已通过', 'approved') THEN 'approved'
				WHEN audit_status IN ('审核驳回', '已驳回', 'rejected') THEN 'rejected'
				ELSE IFNULL(audit_status, 'draft')
			END,
			IFNULL(reject_reason, ''),
			IFNULL(auditor_id, 0),
			IFNULL(auditor_name, ''),
			IFNULL(DATE_FORMAT(audit_time, '%Y-%m-%d'), ''),
			IFNULL(remark, '')
		FROM factory_tests
		WHERE IFNULL(is_deleted, 0) = 0
		ORDER BY id DESC
	`)

	if err != nil {
		http.Error(w, "查询出厂测试记录失败: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	type FactoryTestVO struct {
		ID           int64 `json:"id"`
		BurnRecordID int64 `json:"burnRecordId"`
		ProjectID    int64 `json:"projectId"`

		ProductModel string `json:"productModel"`
		DeviceType   string `json:"deviceType"`
		MacAddress   string `json:"macAddress"`
		SN           string `json:"sn"`

		FileID int64 `json:"fileId"`

		// 前端页面需要这些字段，但数据库没有 file_name/file_url/record_name，
		// 所以这里给占位，避免前端报 undefined。
		RecordName string `json:"recordName"`
		FileName   string `json:"fileName"`
		FileURL    string `json:"fileUrl"`

		UploaderID   int64  `json:"uploaderId"`
		Uploader     string `json:"uploader"`
		UploaderName string `json:"uploaderName"`
		UploadTime   string `json:"uploadTime"`

		AuditStatus  string `json:"auditStatus"`
		RejectReason string `json:"rejectReason"`

		AuditorID   int64  `json:"auditorId"`
		Auditor     string `json:"auditor"`
		AuditorName string `json:"auditorName"`
		AuditTime   string `json:"auditTime"`

		Remark string `json:"remark"`
	}

	list := make([]FactoryTestVO, 0)

	for rows.Next() {
		var item FactoryTestVO

		err := rows.Scan(
			&item.ID,
			&item.BurnRecordID,
			&item.ProjectID,
			&item.ProductModel,
			&item.DeviceType,
			&item.MacAddress,
			&item.SN,
			&item.FileID,
			&item.UploaderID,
			&item.UploaderName,
			&item.UploadTime,
			&item.AuditStatus,
			&item.RejectReason,
			&item.AuditorID,
			&item.AuditorName,
			&item.AuditTime,
			&item.Remark,
		)

		if err != nil {
			http.Error(w, "数据解析失败: "+err.Error(), http.StatusInternalServerError)
			return
		}

		item.Uploader = item.UploaderName
		item.Auditor = item.AuditorName

		// 当前 factory_tests 表没有 file_name/file_url 字段，只能给前端占位
		item.RecordName = "出厂测试文档"
		item.FileName = "出厂测试文档"
		item.FileURL = ""

		list = append(list, item)
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"code": 200,
		"msg":  "查询成功",
		"data": list,
	})
}

func CreateFactoryTestHandler(w http.ResponseWriter, r *http.Request) {
	var item model.FactoryTest

	err := json.NewDecoder(r.Body).Decode(&item)
	if err != nil {
		http.Error(w, "参数解析失败: "+err.Error(), http.StatusBadRequest)
		return
	}

	if item.BurnRecordID == 0 {
		http.Error(w, "烧录记录ID不能为空", http.StatusBadRequest)
		return
	}

	if item.FileID == 0 {
		http.Error(w, "测试文件ID不能为空", http.StatusBadRequest)
		return
	}

	if item.AuditStatus == "" {
		item.AuditStatus = "待审核"
	}

	now := time.Now()

	if item.SN == "" || item.MacAddress == "" {
		_ = config.DB.QueryRow(`
			SELECT
				IFNULL(project_id, 0),
				IFNULL(product_model, ''),
				IFNULL(device_type, ''),
				sn,
				mac_address
			FROM burn_records
			WHERE id = ? AND is_deleted = 0
		`, item.BurnRecordID).Scan(
			&item.ProjectID,
			&item.ProductModel,
			&item.DeviceType,
			&item.SN,
			&item.MacAddress,
		)
	}

	result, err := config.DB.Exec(`
		INSERT INTO factory_tests (
			burn_record_id,
			project_id,
			product_model,
			device_type,
			mac_address,
			sn,
			file_id,
			uploader_id,
			uploader_name,
			upload_time,
			audit_status,
			reject_reason,
			remark,
			is_deleted,
			created_at,
			updated_at
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, 0, ?, ?)
	`,
		item.BurnRecordID,
		item.ProjectID,
		item.ProductModel,
		item.DeviceType,
		item.MacAddress,
		item.SN,
		item.FileID,
		item.UploaderID,
		item.UploaderName,
		now,
		item.AuditStatus,
		item.RejectReason,
		item.Remark,
		now,
		now,
	)

	if err != nil {
		http.Error(w, "新增失败: "+err.Error(), http.StatusInternalServerError)
		return
	}

	id, _ := result.LastInsertId()

	json.NewEncoder(w).Encode(map[string]interface{}{
		"code": 200,
		"msg":  "新增成功",
		"data": map[string]interface{}{"id": id},
	})
}

type FactoryTestAuditRequest struct {
	AuditorID    int64  `json:"auditorId"`
	AuditorName  string `json:"auditorName"`
	AuditStatus  string `json:"auditStatus"`
	RejectReason string `json:"rejectReason"`
}

func AuditFactoryTestHandler(w http.ResponseWriter, r *http.Request, id int64) {
	var req FactoryTestAuditRequest

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "参数解析失败: "+err.Error(), http.StatusBadRequest)
		return
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
		UPDATE factory_tests
		SET
			audit_status = ?,
			reject_reason = ?,
			auditor_id = ?,
			auditor_name = ?,
			audit_time = NOW(),
			updated_at = NOW()
		WHERE id = ? AND is_deleted = 0
	`,
		req.AuditStatus,
		req.RejectReason,
		req.AuditorID,
		req.AuditorName,
		id,
	)

	if err != nil {
		http.Error(w, "审核失败: "+err.Error(), http.StatusInternalServerError)
		return
	}

	affected, _ := result.RowsAffected()
	if affected == 0 {
		http.Error(w, "出厂测试不存在或已删除", http.StatusNotFound)
		return
	}

	// 出厂测试通过后，自动入库
	if req.AuditStatus == "已通过" {
		var burn model.BurnRecord
		var factoryTestID int64 = id

		err = tx.QueryRow(`
			SELECT
				br.id,
				IFNULL(br.project_id, 0),
				IFNULL(br.product_name, ''),
				IFNULL(br.product_model, ''),
				IFNULL(br.device_type, ''),
				br.sn,
				br.mac_address,
				IFNULL(br.hardware_id, 0),
				IFNULL(br.hardware_version, ''),
				IFNULL(br.software_id, 0),
				IFNULL(br.software_version, '')
			FROM factory_tests ft
			JOIN burn_records br ON ft.burn_record_id = br.id
			WHERE ft.id = ? AND ft.is_deleted = 0 AND br.is_deleted = 0
		`, id).Scan(
			&burn.ID,
			&burn.ProjectID,
			&burn.ProductName,
			&burn.ProductModel,
			&burn.DeviceType,
			&burn.SN,
			&burn.MacAddress,
			&burn.HardwareID,
			&burn.HardwareVersion,
			&burn.SoftwareID,
			&burn.SoftwareVersion,
		)

		if err != nil {
			http.Error(w, "读取烧录记录失败: "+err.Error(), http.StatusInternalServerError)
			return
		}

		_, err = tx.Exec(`
			INSERT INTO inventory_devices (
				project_id,
				device_type,
				product_name,
				product_model,
				sn,
				mac_address,
				hardware_id,
				hardware_version,
				software_id,
				software_version,
				inventory_status,
				source_burn_record_id,
				factory_test_id,
				in_time,
				update_time,
				remark,
				is_deleted
			) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, '在库', ?, ?, NOW(), NOW(), '', 0)
			ON DUPLICATE KEY UPDATE
				inventory_status = '在库',
				factory_test_id = VALUES(factory_test_id),
				update_time = NOW()
		`,
			burn.ProjectID,
			burn.DeviceType,
			burn.ProductName,
			burn.ProductModel,
			burn.SN,
			burn.MacAddress,
			burn.HardwareID,
			burn.HardwareVersion,
			burn.SoftwareID,
			burn.SoftwareVersion,
			burn.ID,
			factoryTestID,
		)

		if err != nil {
			http.Error(w, "自动入库失败: "+err.Error(), http.StatusInternalServerError)
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

func DeleteFactoryTestHandler(w http.ResponseWriter, r *http.Request, id int64) {
	result, err := config.DB.Exec(`
		UPDATE factory_tests
		SET is_deleted = 1, updated_at = NOW()
		WHERE id = ?
	`, id)

	if err != nil {
		http.Error(w, "删除失败: "+err.Error(), http.StatusInternalServerError)
		return
	}

	affected, _ := result.RowsAffected()
	if affected == 0 {
		http.Error(w, "出厂测试不存在", http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"code": 200,
		"msg":  "删除成功",
	})
}

// ============================================================
// 库存
// ============================================================

func InventoryHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	switch r.Method {
	case http.MethodGet:
		GetInventoryHandler(w, r)
	default:
		http.Error(w, "不支持该请求方法", http.StatusMethodNotAllowed)
	}
}

func InventoryActionHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	path := strings.TrimPrefix(r.URL.Path, "/api/inventory/")
	path = strings.Trim(path, "/")
	parts := strings.Split(path, "/")

	if len(parts) < 1 || parts[0] == "" {
		http.Error(w, "缺少库存ID", http.StatusBadRequest)
		return
	}

	id, err := strconv.ParseInt(parts[0], 10, 64)
	if err != nil {
		http.Error(w, "库存ID错误", http.StatusBadRequest)
		return
	}

	if len(parts) == 1 && r.Method == http.MethodPut {
		UpdateInventoryHandler(w, r, id)
		return
	}

	http.Error(w, "接口不存在", http.StatusNotFound)
}

func GetInventoryHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	rows, err := config.DB.Query(`
		SELECT
			id,
			IFNULL(project_id, 0),
			IFNULL(device_type, ''),
			IFNULL(product_name, ''),
			IFNULL(product_model, ''),
			IFNULL(sn, ''),
			IFNULL(mac_address, ''),
			IFNULL(hardware_id, 0),
			IFNULL(hardware_version, ''),
			IFNULL(software_id, 0),
			IFNULL(software_version, ''),
			IFNULL(inventory_status, ''),
			IFNULL(source_burn_record_id, 0),
			IFNULL(factory_test_id, 0),
			IFNULL(DATE_FORMAT(in_time, '%Y-%m-%d %H:%i:%s'), ''),
			IFNULL(DATE_FORMAT(update_time, '%Y-%m-%d %H:%i:%s'), ''),
			IFNULL(remark, ''),
			IFNULL(is_deleted, 0)
		FROM inventory_devices
		WHERE IFNULL(is_deleted, 0) = 0
		ORDER BY id DESC
	`)
	if err != nil {
		http.Error(w, "查询库存失败: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	type InventoryVO struct {
		ID                 int64  `json:"id"`
		ProjectID          int64  `json:"projectId"`
		DeviceType         string `json:"deviceType"`
		ProductName        string `json:"productName"`
		ProductModel       string `json:"productModel"`
		SN                 string `json:"sn"`
		MacAddress         string `json:"macAddress"`
		HardwareID         int64  `json:"hardwareId"`
		HardwareVersion    string `json:"hardwareVersion"`
		SoftwareID         int64  `json:"softwareId"`
		SoftwareVersion    string `json:"softwareVersion"`
		InventoryStatus    string `json:"inventoryStatus"`
		SourceBurnRecordID int64  `json:"sourceBurnRecordId"`
		FactoryTestID      int64  `json:"factoryTestId"`
		InTime             string `json:"inTime"`
		UpdateTime         string `json:"updateTime"`
		Remark             string `json:"remark"`
		IsDeleted          int    `json:"isDeleted"`
	}

	list := make([]InventoryVO, 0)

	for rows.Next() {
		var item InventoryVO

		err := rows.Scan(
			&item.ID,
			&item.ProjectID,
			&item.DeviceType,
			&item.ProductName,
			&item.ProductModel,
			&item.SN,
			&item.MacAddress,
			&item.HardwareID,
			&item.HardwareVersion,
			&item.SoftwareID,
			&item.SoftwareVersion,
			&item.InventoryStatus,
			&item.SourceBurnRecordID,
			&item.FactoryTestID,
			&item.InTime,
			&item.UpdateTime,
			&item.Remark,
			&item.IsDeleted,
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

func UpdateInventoryHandler(w http.ResponseWriter, r *http.Request, id int64) {
	w.Header().Set("Content-Type", "application/json")

	var req struct {
		InventoryStatus string `json:"inventoryStatus"`
		Remark          string `json:"remark"`
	}

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "参数解析失败: "+err.Error(), http.StatusBadRequest)
		return
	}

	req.InventoryStatus = strings.TrimSpace(req.InventoryStatus)

	if req.InventoryStatus == "" {
		http.Error(w, "库存状态不能为空", http.StatusBadRequest)
		return
	}

	allowStatus := map[string]bool{
		"在库":  true,
		"锁定":  true,
		"已出库": true,
		"返修":  true,

		// 兼容前端旧状态
		"ready":        true,
		"locked":       true,
		"outbound":     true,
		"repair":       true,
		"testing":      true,
		"waiting_burn": true,
	}

	if !allowStatus[req.InventoryStatus] {
		http.Error(w, "库存状态不合法", http.StatusBadRequest)
		return
	}

	result, err := config.DB.Exec(`
		UPDATE inventory_devices
		SET
			inventory_status = ?,
			remark = ?,
			update_time = NOW()
		WHERE id = ?
		  AND IFNULL(is_deleted, 0) = 0
	`,
		req.InventoryStatus,
		req.Remark,
		id,
	)

	if err != nil {
		http.Error(w, "修改库存失败: "+err.Error(), http.StatusInternalServerError)
		return
	}

	affected, _ := result.RowsAffected()
	if affected == 0 {
		http.Error(w, "库存设备不存在或已删除", http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"code": 200,
		"msg":  "修改成功",
	})
}

func ImportBurnRecordsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var req struct {
		Records []struct {
			BatchNo string `json:"batchNo"`

			ProjectID         int64 `json:"projectId"`
			ProductionOrderID int64 `json:"productionOrderId"`

			ProductName  string `json:"productName"`
			ProductModel string `json:"productModel"`
			ProductCode  string `json:"productCode"`
			DeviceType   string `json:"deviceType"`

			SerialNumber string `json:"serialNumber"`
			SN           string `json:"sn"`

			MacAddress string `json:"macAddress"`

			HardwareID      int64  `json:"hardwareId"`
			HardwareVersion string `json:"hardwareVersion"`

			SoftwareID      int64  `json:"softwareId"`
			SoftwareVersion string `json:"softwareVersion"`

			PcbQrCode string `json:"pcbQrCode"`
			PcbQRCode string `json:"pcbQRCode"`

			Note string `json:"note"`

			SourceFileID int64 `json:"sourceFileId"`

			UploaderID   int64  `json:"uploaderId"`
			UploaderName string `json:"uploaderName"`
			Uploader     string `json:"uploader"`

			BurnDesc     string `json:"burnDesc"`
			ImportRemark string `json:"importRemark"`
		} `json:"records"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "参数解析失败: "+err.Error(), http.StatusBadRequest)
		return
	}

	if len(req.Records) == 0 {
		http.Error(w, "导入数据不能为空", http.StatusBadRequest)
		return
	}

	tx, err := config.DB.Begin()
	if err != nil {
		http.Error(w, "开启事务失败: "+err.Error(), http.StatusInternalServerError)
		return
	}

	insertCount := 0

	for _, item := range req.Records {
		sn := strings.TrimSpace(item.SerialNumber)
		if sn == "" {
			sn = strings.TrimSpace(item.SN)
		}

		pcbQrCode := strings.TrimSpace(item.PcbQrCode)
		if pcbQrCode == "" {
			pcbQrCode = strings.TrimSpace(item.PcbQRCode)
		}

		uploaderName := strings.TrimSpace(item.UploaderName)
		if uploaderName == "" {
			uploaderName = strings.TrimSpace(item.Uploader)
		}

		burnDesc := strings.TrimSpace(item.BurnDesc)
		if burnDesc == "" {
			burnDesc = strings.TrimSpace(item.ImportRemark)
		}

		deviceType := strings.TrimSpace(item.DeviceType)
		if deviceType == "" {
			deviceType = strings.TrimSpace(item.ProductName)
		}

		batchNo := strings.TrimSpace(item.BatchNo)
		if batchNo == "" {
			tx.Rollback()
			http.Error(w, "导入失败：生产批次号不能为空", http.StatusBadRequest)
			return
		}

		if sn == "" || sn == "-" {
			tx.Rollback()
			http.Error(w, "导入失败：序列号不能为空，且不能为 -", http.StatusBadRequest)
			return
		}

		macAddress := strings.TrimSpace(item.MacAddress)
		if macAddress == "-" {
			macAddress = ""
		}

		_, err := tx.Exec(`
			INSERT INTO burn_records (
				batch_no,
				project_id,
				production_order_id,
				product_name,
				product_model,
				product_code,
				device_type,
				sn,
				mac_address,
				hardware_id,
				hardware_version,
				software_id,
				software_version,
				pcb_qr_code,
				note,
				source_file_id,
				uploader_id,
				uploader_name,
				upload_time,
				burn_desc,
				is_deleted,
				created_at,
				updated_at
			) VALUES (?, ?, ?, ?, ?, ?, ?, ?, NULLIF(?, ''), ?, ?, ?, ?, ?, ?, ?, ?, ?, NOW(), ?, 0, NOW(), NOW())
		`,
			batchNo,
			item.ProjectID,
			item.ProductionOrderID,
			item.ProductName,
			item.ProductModel,
			item.ProductCode,
			deviceType,
			sn,
			macAddress,
			item.HardwareID,
			item.HardwareVersion,
			item.SoftwareID,
			item.SoftwareVersion,
			pcbQrCode,
			item.Note,
			item.SourceFileID,
			item.UploaderID,
			uploaderName,
			burnDesc,
		)

		if err != nil {
			tx.Rollback()
			http.Error(w, "导入失败，SN或MAC可能重复: "+err.Error(), http.StatusInternalServerError)
			return
		}

		insertCount++
	}

	if err := tx.Commit(); err != nil {
		http.Error(w, "提交事务失败: "+err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"code": 200,
		"msg":  "导入成功",
		"data": map[string]interface{}{
			"count": insertCount,
		},
	})
}
func DeleteBurnBatchHandler(w http.ResponseWriter, r *http.Request, batchNo string) {
	w.Header().Set("Content-Type", "application/json")

	batchNo = strings.TrimSpace(batchNo)

	if batchNo == "" {
		http.Error(w, "生产批次号不能为空", http.StatusBadRequest)
		return
	}

	result, err := config.DB.Exec(`
		UPDATE burn_records
		SET is_deleted = 1
		WHERE batch_no = ? AND IFNULL(is_deleted, 0) = 0
	`, batchNo)
	if err != nil {
		http.Error(w, "删除批次失败: "+err.Error(), http.StatusInternalServerError)
		return
	}

	affected, _ := result.RowsAffected()

	json.NewEncoder(w).Encode(map[string]interface{}{
		"code": 200,
		"msg":  "删除批次成功",
		"data": map[string]interface{}{
			"count": affected,
		},
	})
}

func DeleteBurnRecordHandler(w http.ResponseWriter, r *http.Request, id int64) {
	w.Header().Set("Content-Type", "application/json")

	result, err := config.DB.Exec(`
		UPDATE burn_records
		SET is_deleted = 1
		WHERE id = ? AND IFNULL(is_deleted, 0) = 0
	`, id)
	if err != nil {
		http.Error(w, "删除失败: "+err.Error(), http.StatusInternalServerError)
		return
	}

	affected, _ := result.RowsAffected()
	if affected == 0 {
		http.Error(w, "烧录记录不存在或已删除", http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"code": 200,
		"msg":  "删除成功",
	})
}

func ImportFactoryTestsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var req struct {
		Records []struct {
			ProductModel string `json:"productModel"`
			MacAddress   string `json:"macAddress"`
			SN           string `json:"sn"`

			FileID int64 `json:"fileId"`

			UploaderID   int64  `json:"uploaderId"`
			UploaderName string `json:"uploaderName"`

			AuditStatus string `json:"auditStatus"`
			Remark      string `json:"remark"`
		} `json:"records"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "参数解析失败: "+err.Error(), http.StatusBadRequest)
		return
	}

	if len(req.Records) == 0 {
		http.Error(w, "出厂测试记录不能为空", http.StatusBadRequest)
		return
	}

	tx, err := config.DB.Begin()
	if err != nil {
		http.Error(w, "开启事务失败: "+err.Error(), http.StatusInternalServerError)
		return
	}

	count := 0

	for _, item := range req.Records {
		productModel := strings.TrimSpace(item.ProductModel)
		macAddress := strings.TrimSpace(item.MacAddress)
		sn := strings.TrimSpace(item.SN)

		if productModel == "" {
			tx.Rollback()
			http.Error(w, "产品型号不能为空", http.StatusBadRequest)
			return
		}

		if macAddress == "" && sn == "" {
			tx.Rollback()
			http.Error(w, "MAC地址和SN不能同时为空", http.StatusBadRequest)
			return
		}

		var burnRecordID int64
		var projectID int64
		var deviceType string

		err := tx.QueryRow(`
			SELECT
				id,
				IFNULL(project_id, 0),
				IFNULL(device_type, '')
			FROM burn_records
			WHERE IFNULL(is_deleted, 0) = 0
			  AND (
					(? <> '' AND mac_address = ?)
					OR
					(? <> '' AND sn = ?)
			  )
			LIMIT 1
		`,
			macAddress, macAddress,
			sn, sn,
		).Scan(&burnRecordID, &projectID, &deviceType)

		if err == sql.ErrNoRows {
			tx.Rollback()
			http.Error(w, "找不到对应的烧录记录，MAC: "+macAddress+"，SN: "+sn, http.StatusBadRequest)
			return
		}

		if err != nil {
			tx.Rollback()
			http.Error(w, "查询烧录记录失败: "+err.Error(), http.StatusInternalServerError)
			return
		}

		auditStatus := strings.TrimSpace(item.AuditStatus)
		if auditStatus == "" {
			auditStatus = "draft"
		}

		_, err = tx.Exec(`
			INSERT INTO factory_tests (
				burn_record_id,
				project_id,
				product_model,
				device_type,
				mac_address,
				sn,
				file_id,
				uploader_id,
				uploader_name,
				upload_time,
				audit_status,
				reject_reason,
				auditor_id,
				auditor_name,
				audit_time,
				remark,
				is_deleted,
				created_at,
				updated_at
			) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, NOW(), ?, '', 0, '', NULL, ?, 0, NOW(), NOW())
		`,
			burnRecordID,
			projectID,
			productModel,
			deviceType,
			macAddress,
			sn,
			item.FileID,
			item.UploaderID,
			item.UploaderName,
			auditStatus,
			item.Remark,
		)

		if err != nil {
			tx.Rollback()
			http.Error(w, "保存出厂测试记录失败: "+err.Error(), http.StatusInternalServerError)
			return
		}

		count++
	}

	if err := tx.Commit(); err != nil {
		http.Error(w, "提交事务失败: "+err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"code": 200,
		"msg":  "保存成功",
		"data": map[string]interface{}{
			"count": count,
		},
	})
}

func DeleteFactoryTestsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var req struct {
		IDs []int64 `json:"ids"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "参数解析失败: "+err.Error(), http.StatusBadRequest)
		return
	}

	if len(req.IDs) == 0 {
		http.Error(w, "请选择要删除的出厂测试记录", http.StatusBadRequest)
		return
	}

	tx, err := config.DB.Begin()
	if err != nil {
		http.Error(w, "开启事务失败: "+err.Error(), http.StatusInternalServerError)
		return
	}

	count := int64(0)

	for _, id := range req.IDs {
		result, err := tx.Exec(`
			UPDATE factory_tests
			SET is_deleted = 1,
			    updated_at = NOW()
			WHERE id = ?
			  AND IFNULL(is_deleted, 0) = 0
			  AND IFNULL(audit_status, 'draft') <> 'approved'
		`, id)

		if err != nil {
			tx.Rollback()
			http.Error(w, "删除失败: "+err.Error(), http.StatusInternalServerError)
			return
		}

		affected, _ := result.RowsAffected()
		count += affected
	}

	if err := tx.Commit(); err != nil {
		http.Error(w, "提交事务失败: "+err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"code": 200,
		"msg":  "删除成功",
		"data": map[string]interface{}{
			"count": count,
		},
	})
}

func SubmitFactoryTestsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var req struct {
		IDs []int64 `json:"ids"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "参数解析失败: "+err.Error(), http.StatusBadRequest)
		return
	}

	if len(req.IDs) == 0 {
		http.Error(w, "请选择要提交的记录", http.StatusBadRequest)
		return
	}

	tx, err := config.DB.Begin()
	if err != nil {
		http.Error(w, "开启事务失败: "+err.Error(), http.StatusInternalServerError)
		return
	}

	for _, id := range req.IDs {
		_, err := tx.Exec(`
			UPDATE factory_tests
			SET audit_status = 'submitted',
			    auditor_id = 0,
			    auditor_name = '',
			    audit_time = NULL,
			    updated_at = NOW()
			WHERE id = ?
			  AND IFNULL(is_deleted, 0) = 0
			  AND IFNULL(audit_status, 'draft') IN ('draft', 'rejected')
		`, id)

		if err != nil {
			tx.Rollback()
			http.Error(w, "提交失败: "+err.Error(), http.StatusInternalServerError)
			return
		}
	}

	if err := tx.Commit(); err != nil {
		http.Error(w, "提交事务失败: "+err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"code": 200,
		"msg":  "提交成功",
	})
}

func AuditFactoryTestsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var req struct {
		IDs          []int64 `json:"ids"`
		Status       string  `json:"status"`
		AuditorID    int64   `json:"auditorId"`
		AuditorName  string  `json:"auditorName"`
		RejectReason string  `json:"rejectReason"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "参数解析失败: "+err.Error(), http.StatusBadRequest)
		return
	}

	if len(req.IDs) == 0 {
		http.Error(w, "请选择要审核的记录", http.StatusBadRequest)
		return
	}

	if req.Status != "approved" && req.Status != "rejected" {
		http.Error(w, "审核状态错误", http.StatusBadRequest)
		return
	}

	if strings.TrimSpace(req.AuditorName) == "" {
		req.AuditorName = "领导"
	}

	tx, err := config.DB.Begin()
	if err != nil {
		http.Error(w, "开启事务失败: "+err.Error(), http.StatusInternalServerError)
		return
	}

	for _, id := range req.IDs {
		result, err := tx.Exec(`
			UPDATE factory_tests
			SET audit_status = ?,
			    reject_reason = ?,
			    auditor_id = ?,
			    auditor_name = ?,
			    audit_time = NOW(),
			    updated_at = NOW()
			WHERE id = ?
			  AND IFNULL(is_deleted, 0) = 0
			  AND audit_status IN ('submitted', '待审核')
		`,
			req.Status,
			req.RejectReason,
			req.AuditorID,
			req.AuditorName,
			id,
		)

		if err != nil {
			tx.Rollback()
			http.Error(w, "审核失败: "+err.Error(), http.StatusInternalServerError)
			return
		}

		affected, _ := result.RowsAffected()
		if affected == 0 {
			tx.Rollback()
			http.Error(w, "审核失败：记录不存在，或当前状态不是待审核", http.StatusBadRequest)
			return
		}

		if req.Status == "approved" {
			err = SyncFactoryTestToInventoryTx(tx, id)
			if err != nil {
				tx.Rollback()
				http.Error(w, "审核通过，但自动入库失败: "+err.Error(), http.StatusInternalServerError)
				return
			}
		}
	}

	if err := tx.Commit(); err != nil {
		http.Error(w, "提交事务失败: "+err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"code": 200,
		"msg":  "审核成功",
	})
}
func DeleteFactoryTestByIDHandler(w http.ResponseWriter, r *http.Request, id int64) {
	w.Header().Set("Content-Type", "application/json")

	result, err := config.DB.Exec(`
		UPDATE factory_tests
		SET is_deleted = 1,
		    updated_at = NOW()
		WHERE id = ?
		  AND IFNULL(is_deleted, 0) = 0
		  AND IFNULL(audit_status, 'draft') <> 'approved'
	`, id)

	if err != nil {
		http.Error(w, "删除失败: "+err.Error(), http.StatusInternalServerError)
		return
	}

	affected, _ := result.RowsAffected()
	if affected == 0 {
		http.Error(w, "记录不存在，或审核通过后不允许删除", http.StatusBadRequest)
		return
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"code": 200,
		"msg":  "删除成功",
	})
}
func GetBurnRecordOptionsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	rows, err := config.DB.Query(`
		SELECT
			IFNULL(batch_no, ''),
			IFNULL(product_name, ''),
			IFNULL(product_model, ''),
			IFNULL(product_code, ''),
			IFNULL(device_type, ''),
			COUNT(*) AS total
		FROM burn_records
		WHERE IFNULL(is_deleted, 0) = 0
		GROUP BY
			batch_no,
			product_name,
			product_model,
			product_code,
			device_type
		ORDER BY batch_no DESC, product_model ASC
	`)
	if err != nil {
		http.Error(w, "查询烧录型号失败: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	type BurnOption struct {
		BatchNo      string `json:"batchNo"`
		ProductName  string `json:"productName"`
		ProductModel string `json:"productModel"`
		ProductCode  string `json:"productCode"`
		DeviceType   string `json:"deviceType"`
		Count        int64  `json:"count"`
	}

	list := make([]BurnOption, 0)

	for rows.Next() {
		var item BurnOption

		err := rows.Scan(
			&item.BatchNo,
			&item.ProductName,
			&item.ProductModel,
			&item.ProductCode,
			&item.DeviceType,
			&item.Count,
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

func SyncFactoryTestToInventoryTx(tx *sql.Tx, factoryTestID int64) error {
	_, err := tx.Exec(`
		INSERT INTO inventory_devices (
			project_id,
			device_type,
			product_name,
			product_model,
			sn,
			mac_address,
			hardware_id,
			hardware_version,
			software_id,
			software_version,
			inventory_status,
			source_burn_record_id,
			factory_test_id,
			in_time,
			update_time,
			remark,
			is_deleted
		)
		SELECT
			IFNULL(br.project_id, 0),
			IFNULL(br.device_type, ''),
			IFNULL(br.product_name, ''),
			IFNULL(br.product_model, ''),
			IFNULL(br.sn, ''),
			IFNULL(br.mac_address, ''),
			IFNULL(br.hardware_id, 0),
			IFNULL(br.hardware_version, ''),
			IFNULL(br.software_id, 0),
			IFNULL(br.software_version, ''),
			'在库',
			br.id,
			ft.id,
			NOW(),
			NOW(),
			'出厂测试审核通过，自动入库',
			0
		FROM factory_tests ft
		INNER JOIN burn_records br ON br.id = ft.burn_record_id
		WHERE ft.id = ?
		  AND IFNULL(ft.is_deleted, 0) = 0
		  AND IFNULL(br.is_deleted, 0) = 0
		  AND ft.audit_status IN ('approved', '审核通过')
		ON DUPLICATE KEY UPDATE
			project_id = VALUES(project_id),
			device_type = VALUES(device_type),
			product_name = VALUES(product_name),
			product_model = VALUES(product_model),
			mac_address = VALUES(mac_address),
			hardware_id = VALUES(hardware_id),
			hardware_version = VALUES(hardware_version),
			software_id = VALUES(software_id),
			software_version = VALUES(software_version),
			inventory_status = '在库',
			source_burn_record_id = VALUES(source_burn_record_id),
			factory_test_id = VALUES(factory_test_id),
			update_time = NOW(),
			remark = VALUES(remark),
			is_deleted = 0
	`, factoryTestID)

	return err
}

// ============================================================
// 发货批次 / 出库记录
// 说明：出库记录不单独落表，直接查询“已审核通过的发货批次设备明细”。
// ============================================================

type ShippingBatchDeviceVO struct {
	ID                int64  `json:"id"`
	BatchID           int64  `json:"batchId"`
	InventoryDeviceID int64  `json:"inventoryDeviceId"`
	SN                string `json:"sn"`
	MacAddress        string `json:"macAddress"`
	DeviceType        string `json:"deviceType"`
	HardwareVersion   string `json:"hardwareVersion"`
	SoftwareVersion   string `json:"softwareVersion"`
	CreatedAt         string `json:"createdAt"`
}

type ShippingBatchVO struct {
	ID           int64  `json:"id"`
	ProjectID    int64  `json:"projectId"`
	BatchNo      string `json:"batchNo"`
	ExpressNo    string `json:"expressNo"`
	UploaderID   int64  `json:"uploaderId"`
	UploaderName string `json:"uploaderName"`

	AuditStatus string `json:"auditStatus"`
	AuditorID   int64  `json:"auditorId"`
	AuditorName string `json:"auditorName"`
	AuditTime   string `json:"auditTime"`

	Remark    string                  `json:"remark"`
	CreatedAt string                  `json:"createdAt"`
	UpdatedAt string                  `json:"updatedAt"`
	Devices   []ShippingBatchDeviceVO `json:"devices"`
}

type AuditShippingBatchRequest struct {
	AuditorID    int64  `json:"auditorId"`
	AuditorName  string `json:"auditorName"`
	AuditStatus  string `json:"auditStatus"`
	RejectReason string `json:"rejectReason"`
}

type InventoryDeviceOptionVO struct {
	ID              int64  `json:"id"`
	ProjectID       int64  `json:"projectId"`
	DeviceType      string `json:"deviceType"`
	ProductName     string `json:"productName"`
	ProductModel    string `json:"productModel"`
	SN              string `json:"sn"`
	MacAddress      string `json:"macAddress"`
	HardwareVersion string `json:"hardwareVersion"`
	SoftwareVersion string `json:"softwareVersion"`
	InventoryStatus string `json:"inventoryStatus"`
	InTime          string `json:"inTime"`
}

type OutboundRecordVO struct {
	ID                int64 `json:"id"`
	BatchID           int64 `json:"batchId"`
	InventoryDeviceID int64 `json:"inventoryDeviceId"`

	DeviceType      string `json:"deviceType"`
	SN              string `json:"sn"`
	MacAddress      string `json:"macAddress"`
	SoftwareVersion string `json:"softwareVersion"`
	HardwareVersion string `json:"hardwareVersion"`

	InTime       string `json:"inTime"`
	OutboundTime string `json:"outboundTime"`

	UploaderName string `json:"uploaderName"`
	BatchNo      string `json:"batchNo"`
	ExpressNo    string `json:"expressNo"`
	Remark       string `json:"remark"`
}

func shippingWriteJSON(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(data)
}

func shippingSuccess(w http.ResponseWriter, msg string, data interface{}) {
	shippingWriteJSON(w, map[string]interface{}{
		"code": 200,
		"msg":  msg,
		"data": data,
	})
}

func shippingError(w http.ResponseWriter, status int, msg string) {
	http.Error(w, msg, status)
}

func getShippingBatchDevicesByBatchID(batchID int64) ([]ShippingBatchDeviceVO, error) {
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
			IFNULL(DATE_FORMAT(created_at, '%Y-%m-%d %H:%i:%s'), '')
		FROM shipping_batch_devices
		WHERE batch_id = ?
		  AND IFNULL(is_deleted, 0) = 0
		ORDER BY id ASC
	`, batchID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	list := make([]ShippingBatchDeviceVO, 0)

	for rows.Next() {
		var item ShippingBatchDeviceVO

		err := rows.Scan(
			&item.ID,
			&item.BatchID,
			&item.InventoryDeviceID,
			&item.SN,
			&item.MacAddress,
			&item.DeviceType,
			&item.HardwareVersion,
			&item.SoftwareVersion,
			&item.CreatedAt,
		)
		if err != nil {
			return nil, err
		}

		list = append(list, item)
	}

	return list, rows.Err()
}

func OutboundRecordActionHandler(w http.ResponseWriter, r *http.Request) {
	shippingError(w, http.StatusNotFound, "出库记录由发货批次审核通过后自动生成，不支持单独操作")
}
