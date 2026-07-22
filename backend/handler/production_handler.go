package handler

import (
	"crrc_pm_backend/config"
	"crrc_pm_backend/model"
	"database/sql"
	"encoding/json"
	"log"
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
	ensureUploadedFilesTable()

	rows, err := config.DB.Query(`
		SELECT
			br.id,
			IFNULL(br.batch_no, ''),
			IFNULL(br.project_id, 0),
			IFNULL(br.production_order_id, 0),
			IFNULL(br.product_name, ''),
			IFNULL(br.product_model, ''),
			IFNULL(br.product_code, ''),
			IFNULL(br.device_type, ''),
			IFNULL(br.sn, ''),
			IFNULL(br.mac_address, ''),
			IFNULL(br.hardware_id, 0),
			IFNULL(br.hardware_version, ''),
			IFNULL(br.software_id, 0),
			IFNULL(br.software_version, ''),
			IFNULL(br.pcb_qr_code, ''),
			IFNULL(br.note, ''),
			IFNULL(br.source_file_id, 0),
			IFNULL(uf.file_name, ''),
			IFNULL(br.uploader_id, 0),
			IFNULL(br.uploader_name, ''),
			IFNULL(DATE_FORMAT(br.upload_time, '%Y-%m-%d'), ''),
			IFNULL(br.burn_desc, ''),
			IFNULL(DATE_FORMAT(br.created_at, '%Y-%m-%d %H:%i:%s'), ''),
			IFNULL(DATE_FORMAT(br.updated_at, '%Y-%m-%d %H:%i:%s'), '')
		FROM burn_records br
		LEFT JOIN uploaded_files uf ON uf.id = br.source_file_id
		WHERE IFNULL(br.is_deleted, 0) = 0
		ORDER BY br.id DESC
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
			&item.FileName,
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

		if item.FileName == "" && item.SourceFileID > 0 {
			item.FileName = "文件ID-" + strconv.FormatInt(item.SourceFileID, 10)
		}
		item.FileURL = filePreviewURL(item.SourceFileID)

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
	ensureUploadedFilesTable()

	visibilitySQL := reviewVisibilitySQL(r, "ft.audit_status", "ft.uploader_id", "ft.uploader_name", "production_staff")
	rows, err := config.DB.Query(`
		SELECT
			ft.id,
			IFNULL(ft.burn_record_id, 0),
			IFNULL(ft.project_id, 0),
			IFNULL(ft.product_model, ''),
			IFNULL(ft.device_type, ''),
			IFNULL(ft.mac_address, ''),
			IFNULL(ft.sn, ''),
			IFNULL(ft.file_id, 0),
			IFNULL(uf.file_name, ''),
			IFNULL(ft.uploader_id, 0),
			IFNULL(ft.uploader_name, ''),
			IFNULL(DATE_FORMAT(ft.upload_time, '%Y-%m-%d'), ''),
			CASE
				WHEN ft.audit_status IN ('草稿', 'draft') THEN 'draft'
				WHEN ft.audit_status IN ('待审核', 'submitted') THEN 'submitted'
				WHEN ft.audit_status IN ('审核通过', '已通过', 'approved') THEN 'approved'
				WHEN ft.audit_status IN ('审核驳回', '已驳回', 'rejected') THEN 'rejected'
				ELSE IFNULL(ft.audit_status, 'draft')
			END,
			IFNULL(ft.reject_reason, ''),
			IFNULL(ft.auditor_id, 0),
			IFNULL(ft.auditor_name, ''),
			IFNULL(DATE_FORMAT(ft.audit_time, '%Y-%m-%d'), ''),
			IFNULL(ft.remark, '')
		FROM factory_tests ft
		LEFT JOIN uploaded_files uf ON uf.id = ft.file_id
		WHERE IFNULL(ft.is_deleted, 0) = 0 ` + visibilitySQL + `
		ORDER BY ft.id DESC
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
			&item.FileName,
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

		item.RecordName = "出厂测试文档"
		if item.FileName == "" && item.FileID > 0 {
			item.FileName = "出厂测试文档"
		}
		item.FileURL = filePreviewURL(item.FileID)

		list = append(list, item)
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"code": 200,
		"msg":  "查询成功",
		"data": list,
	})
}

func CreateFactoryTestHandler(w http.ResponseWriter, r *http.Request) {
	var req struct {
		model.FactoryTest
		FileName        string `json:"fileName"`
		FileContentType string `json:"fileContentType"`
		FileData        string `json:"fileData"`
	}

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "参数解析失败: "+err.Error(), http.StatusBadRequest)
		return
	}

	item := req.FactoryTest
	if item.BurnRecordID == 0 {
		http.Error(w, "烧录记录ID不能为空", http.StatusBadRequest)
		return
	}
	var burnProductName string
	var burnProductModel string
	if err := config.DB.QueryRow(`
		SELECT IFNULL(product_name, ''), IFNULL(product_model, '')
		FROM burn_records
		WHERE id = ? AND IFNULL(is_deleted, 0) = 0
	`, item.BurnRecordID).Scan(&burnProductName, &burnProductModel); err != nil {
		http.Error(w, "烧录记录不存在", http.StatusBadRequest)
		return
	}
	if isHandsetProduct(burnProductName, burnProductModel) {
		http.Error(w, "联络电话手持话柄烧录后已直接入库，无需创建出厂测试", http.StatusBadRequest)
		return
	}

	if item.FileID == 0 {
		http.Error(w, "测试文件ID不能为空", http.StatusBadRequest)
		return
	}
	if err := saveUploadedFile(UploadedFilePayload{
		FileID:          item.FileID,
		FileName:        req.FileName,
		FileContentType: req.FileContentType,
		FileData:        req.FileData,
	}); err != nil {
		http.Error(w, "保存出厂测试文件失败: "+err.Error(), http.StatusBadRequest)
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
	if !requireFactoryQualityAuditorPermission(w, r) {
		return
	}

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
	req.AuditorID, req.AuditorName = normalizeAuditUser(r, req.AuditorID, req.AuditorName)

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
				IFNULL(br.product_code, ''),
				IFNULL(br.device_type, ''),
				br.sn,
				br.mac_address,
				IFNULL(br.pcb_qr_code, ''),
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
			&burn.ProductCode,
			&burn.DeviceType,
			&burn.SN,
			&burn.MacAddress,
			&burn.PCBQRCode,
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
				product_code,
				sn,
				mac_address,
				pcb_qr_code,
				hardware_id,
				hardware_version,
				software_id,
				software_version,
				inventory_status,
				source_burn_record_id,
				factory_test_id,
				inbound_type,
				in_time,
				update_time,
				remark,
				is_deleted
			) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, '在库', ?, ?, 'factory', NOW(), NOW(), '', 0)
			ON DUPLICATE KEY UPDATE
				project_id = VALUES(project_id),
				device_type = VALUES(device_type),
				product_name = VALUES(product_name),
				product_model = VALUES(product_model),
				product_code = VALUES(product_code),
				pcb_qr_code = VALUES(pcb_qr_code),
				hardware_id = VALUES(hardware_id),
				hardware_version = VALUES(hardware_version),
				software_id = VALUES(software_id),
				software_version = VALUES(software_version),
				inventory_status = '在库',
				source_burn_record_id = VALUES(source_burn_record_id),
				factory_test_id = VALUES(factory_test_id),
				inbound_type = 'factory',
				update_time = NOW()
		`,
			burn.ProjectID,
			burn.DeviceType,
			burn.ProductName,
			burn.ProductModel,
			burn.ProductCode,
			burn.SN,
			burn.MacAddress,
			burn.PCBQRCode,
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

		if err = ensureBoardCompositionDeductedForFactoryTest(tx, id); err != nil {
			if err == sql.ErrNoRows {
				http.Error(w, "自动入库失败：板卡组成对应的板卡入库库存不足", http.StatusBadRequest)
				return
			}
			http.Error(w, "自动入库失败，扣减板卡组成库存失败: "+err.Error(), http.StatusInternalServerError)
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

func BoardInboundHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	switch r.Method {
	case http.MethodGet:
		GetBoardInboundHandler(w, r)
	case http.MethodPost:
		ImportBoardInboundHandler(w, r)
	case http.MethodDelete:
		DeleteBoardInboundHandler(w, r)
	default:
		http.Error(w, "不支持该请求方法", http.StatusMethodNotAllowed)
	}
}

func ensureInventoryBoardColumns() {
	_, _ = config.DB.Exec(`ALTER TABLE inventory_devices ADD COLUMN product_code VARCHAR(128) DEFAULT '' AFTER product_model`)
	_, _ = config.DB.Exec(`ALTER TABLE inventory_devices ADD COLUMN quantity INT NOT NULL DEFAULT 1 AFTER product_code`)
	_, _ = config.DB.Exec(`ALTER TABLE inventory_devices ADD COLUMN pcb_qr_code VARCHAR(255) DEFAULT '' AFTER mac_address`)
	_, _ = config.DB.Exec(`ALTER TABLE inventory_devices ADD COLUMN source_file_id BIGINT DEFAULT 0 AFTER factory_test_id`)
	_, _ = config.DB.Exec(`ALTER TABLE inventory_devices ADD COLUMN source_file_name VARCHAR(255) DEFAULT '' AFTER source_file_id`)
	_, _ = config.DB.Exec(`ALTER TABLE inventory_devices ADD COLUMN inbound_type VARCHAR(32) DEFAULT '' AFTER source_file_name`)
	_, _ = config.DB.Exec(`ALTER TABLE inventory_devices ADD COLUMN scrap_audit_status VARCHAR(32) DEFAULT '' AFTER inbound_type`)
	_, _ = config.DB.Exec(`ALTER TABLE inventory_devices ADD COLUMN scrap_request_user_id BIGINT DEFAULT 0 AFTER scrap_audit_status`)
	_, _ = config.DB.Exec(`ALTER TABLE inventory_devices ADD COLUMN scrap_request_user_name VARCHAR(64) DEFAULT '' AFTER scrap_request_user_id`)
	_, _ = config.DB.Exec(`ALTER TABLE inventory_devices ADD COLUMN scrap_request_time DATETIME NULL AFTER scrap_request_user_name`)
	_, _ = config.DB.Exec(`ALTER TABLE inventory_devices ADD COLUMN scrap_audit_user_id BIGINT DEFAULT 0 AFTER scrap_request_time`)
	_, _ = config.DB.Exec(`ALTER TABLE inventory_devices ADD COLUMN scrap_audit_user_name VARCHAR(64) DEFAULT '' AFTER scrap_audit_user_id`)
	_, _ = config.DB.Exec(`ALTER TABLE inventory_devices ADD COLUMN scrap_audit_time DATETIME NULL AFTER scrap_audit_user_name`)
	_, _ = config.DB.Exec(`ALTER TABLE inventory_devices ADD COLUMN scrap_reject_reason VARCHAR(255) DEFAULT '' AFTER scrap_audit_time`)
	_, _ = config.DB.Exec(`ALTER TABLE inventory_devices MODIFY COLUMN mac_address VARCHAR(64) NULL`)
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

	if len(parts) == 2 && r.Method == http.MethodPost && parts[1] == "scrap-submit" {
		SubmitInventoryScrapHandler(w, r, id)
		return
	}

	if len(parts) == 2 && r.Method == http.MethodPost && parts[1] == "scrap-audit" {
		AuditInventoryScrapHandler(w, r, id)
		return
	}

	http.Error(w, "接口不存在", http.StatusNotFound)
}

func GetInventoryHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if !hasRequestRole(r, "production_staff") && !hasRequestRole(r, "system_admin") && !hasRequestRole(r, "leader") && !hasRequestPermission(r, "inventory:view") && !hasRequestPermission(r, "production:view") {
		http.Error(w, "无库存情况查看权限", http.StatusForbidden)
		return
	}

	ensureInventoryBoardColumns()

	rows, err := config.DB.Query(`
		SELECT
			id,
			IFNULL(project_id, 0),
			IFNULL(device_type, ''),
			IFNULL(product_name, ''),
			IFNULL(product_model, ''),
			IFNULL(product_code, ''),
			IFNULL(quantity, 1),
			IFNULL(sn, ''),
			IFNULL(mac_address, ''),
			IFNULL(pcb_qr_code, ''),
			IFNULL(hardware_id, 0),
			IFNULL(hardware_version, ''),
			IFNULL(software_id, 0),
			IFNULL(software_version, ''),
			IFNULL(inventory_status, ''),
			IFNULL(source_burn_record_id, 0),
			IFNULL(factory_test_id, 0),
			IFNULL(source_file_id, 0),
			IFNULL(source_file_name, ''),
			IFNULL(inbound_type, ''),
			IFNULL(scrap_audit_status, ''),
			IFNULL(scrap_request_user_id, 0),
			IFNULL(scrap_request_user_name, ''),
			IFNULL(DATE_FORMAT(scrap_request_time, '%Y-%m-%d %H:%i:%s'), ''),
			IFNULL(scrap_audit_user_id, 0),
			IFNULL(scrap_audit_user_name, ''),
			IFNULL(DATE_FORMAT(scrap_audit_time, '%Y-%m-%d %H:%i:%s'), ''),
			IFNULL(scrap_reject_reason, ''),
			IFNULL(DATE_FORMAT(in_time, '%Y-%m-%d %H:%i:%s'), ''),
			IFNULL(DATE_FORMAT(update_time, '%Y-%m-%d %H:%i:%s'), ''),
			IFNULL(remark, ''),
			IFNULL(is_deleted, 0)
		FROM inventory_devices
		WHERE IFNULL(is_deleted, 0) = 0
		  AND IFNULL(inventory_status, '') NOT IN ('板卡入库', '已烧录', '已出库')
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
		ProductCode        string `json:"productCode"`
		Quantity           int    `json:"quantity"`
		SN                 string `json:"sn"`
		MacAddress         string `json:"macAddress"`
		PcbQrCode          string `json:"pcbQrCode"`
		HardwareID         int64  `json:"hardwareId"`
		HardwareVersion    string `json:"hardwareVersion"`
		SoftwareID         int64  `json:"softwareId"`
		SoftwareVersion    string `json:"softwareVersion"`
		InventoryStatus    string `json:"inventoryStatus"`
		SourceBurnRecordID int64  `json:"sourceBurnRecordId"`
		FactoryTestID      int64  `json:"factoryTestId"`
		SourceFileID       int64  `json:"sourceFileId"`
		SourceFileName     string `json:"sourceFileName"`
		InboundType        string `json:"inboundType"`
		ScrapAuditStatus   string `json:"scrapAuditStatus"`
		ScrapRequestUserID int64  `json:"scrapRequestUserId"`
		ScrapRequestUser   string `json:"scrapRequestUserName"`
		ScrapRequestTime   string `json:"scrapRequestTime"`
		ScrapAuditUserID   int64  `json:"scrapAuditUserId"`
		ScrapAuditUser     string `json:"scrapAuditUserName"`
		ScrapAuditTime     string `json:"scrapAuditTime"`
		ScrapRejectReason  string `json:"scrapRejectReason"`
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
			&item.ProductCode,
			&item.Quantity,
			&item.SN,
			&item.MacAddress,
			&item.PcbQrCode,
			&item.HardwareID,
			&item.HardwareVersion,
			&item.SoftwareID,
			&item.SoftwareVersion,
			&item.InventoryStatus,
			&item.SourceBurnRecordID,
			&item.FactoryTestID,
			&item.SourceFileID,
			&item.SourceFileName,
			&item.InboundType,
			&item.ScrapAuditStatus,
			&item.ScrapRequestUserID,
			&item.ScrapRequestUser,
			&item.ScrapRequestTime,
			&item.ScrapAuditUserID,
			&item.ScrapAuditUser,
			&item.ScrapAuditTime,
			&item.ScrapRejectReason,
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

func GetBoardInboundHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if !hasRequestRole(r, "production_staff") && !hasRequestRole(r, "shipping_staff") && !hasRequestRole(r, "system_admin") && !hasRequestRole(r, "leader") && !hasRequestPermission(r, "board-inbound:view") && !hasRequestPermission(r, "production:view") {
		http.Error(w, "无板卡入库查看权限", http.StatusForbidden)
		return
	}

	ensureInventoryBoardColumns()

	rows, err := config.DB.Query(`
		SELECT
			id,
			IFNULL(product_name, ''),
			IFNULL(product_model, ''),
			IFNULL(product_code, ''),
			IFNULL(quantity, 1),
			IFNULL(sn, ''),
			IFNULL(mac_address, ''),
			IFNULL(pcb_qr_code, ''),
			IFNULL(remark, ''),
			IFNULL(source_file_id, 0),
			IFNULL(source_file_name, ''),
			IFNULL(inventory_status, ''),
			IFNULL(DATE_FORMAT(in_time, '%Y-%m-%d %H:%i:%s'), ''),
			IFNULL(DATE_FORMAT(update_time, '%Y-%m-%d %H:%i:%s'), '')
		FROM inventory_devices
		WHERE IFNULL(is_deleted, 0) = 0
		  AND IFNULL(inbound_type, '') = 'board'
		  AND IFNULL(inventory_status, '') IN ('板卡入库', '已烧录')
		ORDER BY id DESC
	`)
	if err != nil {
		http.Error(w, "查询板卡入库失败: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	type BoardInboundVO struct {
		ID              int64  `json:"id"`
		ProductName     string `json:"productName"`
		ProductModel    string `json:"productModel"`
		ProductCode     string `json:"productCode"`
		Quantity        int    `json:"quantity"`
		SN              string `json:"sn"`
		MacAddress      string `json:"macAddress"`
		PcbQrCode       string `json:"pcbQrCode"`
		Remark          string `json:"remark"`
		SourceFileID    int64  `json:"sourceFileId"`
		SourceFileName  string `json:"sourceFileName"`
		InventoryStatus string `json:"inventoryStatus"`
		InTime          string `json:"inTime"`
		UpdateTime      string `json:"updateTime"`
	}

	list := make([]BoardInboundVO, 0)
	for rows.Next() {
		var item BoardInboundVO
		if err := rows.Scan(
			&item.ID,
			&item.ProductName,
			&item.ProductModel,
			&item.ProductCode,
			&item.Quantity,
			&item.SN,
			&item.MacAddress,
			&item.PcbQrCode,
			&item.Remark,
			&item.SourceFileID,
			&item.SourceFileName,
			&item.InventoryStatus,
			&item.InTime,
			&item.UpdateTime,
		); err != nil {
			http.Error(w, "板卡入库数据解析失败: "+err.Error(), http.StatusInternalServerError)
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

func ImportBoardInboundHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if !hasRequestRole(r, "production_staff") && !hasRequestRole(r, "system_admin") && !hasRequestPermission(r, "board-inbound:import") {
		http.Error(w, "无板卡入库导入权限", http.StatusForbidden)
		return
	}

	var req struct {
		SourceFileID    int64  `json:"sourceFileId"`
		FileName        string `json:"fileName"`
		FileContentType string `json:"fileContentType"`
		FileData        string `json:"fileData"`
		Records         []struct {
			ProductName  string `json:"productName"`
			ProductModel string `json:"productModel"`
			ProductCode  string `json:"productCode"`
			Quantity     int    `json:"quantity"`
			RowNumber    int    `json:"rowNumber"`
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

	if err := saveUploadedFile(UploadedFilePayload{
		FileID:          req.SourceFileID,
		FileName:        req.FileName,
		FileContentType: req.FileContentType,
		FileData:        req.FileData,
	}); err != nil {
		http.Error(w, "保存板卡入库源文件失败: "+err.Error(), http.StatusBadRequest)
		return
	}

	ensureInventoryBoardColumns()

	tx, err := config.DB.Begin()
	if err != nil {
		http.Error(w, "开启事务失败: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer tx.Rollback()

	insertCount := 0
	for index, item := range req.Records {
		productName := strings.TrimSpace(item.ProductName)
		productModel := strings.TrimSpace(item.ProductModel)
		productCode := strings.TrimSpace(item.ProductCode)
		quantity := item.Quantity
		if item.RowNumber > 0 && quantity <= 0 {
			http.Error(w, "导入失败：第 "+strconv.Itoa(item.RowNumber)+" 行数量必须大于 0", http.StatusBadRequest)
			return
		}
		if productName == "" || productModel == "" || productCode == "" || quantity <= 0 {
			rowLabel := strconv.Itoa(index + 1)
			if item.RowNumber > 0 {
				rowLabel = strconv.Itoa(item.RowNumber)
			}
			http.Error(w, "导入失败：第 "+rowLabel+" 行产品名称、产品型号、生产编码、数量不能为空", http.StatusBadRequest)
			return
		}

		rowNumber := item.RowNumber
		if rowNumber <= 0 {
			rowNumber = index + 4
		}

		syntheticCode := strconv.FormatInt(req.SourceFileID, 10) + "-" + strconv.Itoa(rowNumber)
		sn := "BOARD-IN-SN-" + syntheticCode
		macAddress := "BOARD-IN-MAC-" + syntheticCode

		_, err = tx.Exec(`
			INSERT INTO inventory_devices (
				project_id,
				device_type,
				product_name,
				product_model,
				product_code,
				quantity,
				sn,
				mac_address,
				pcb_qr_code,
				hardware_id,
				hardware_version,
				software_id,
				software_version,
				inventory_status,
				source_burn_record_id,
				factory_test_id,
				source_file_id,
				source_file_name,
				inbound_type,
				in_time,
				update_time,
				remark,
				is_deleted
			) VALUES (0, ?, ?, ?, ?, ?, ?, ?, '', 0, '', 0, '', '板卡入库', 0, 0, ?, ?, 'board', NOW(), NOW(), '', 0)
			ON DUPLICATE KEY UPDATE
				device_type = VALUES(device_type),
				product_name = VALUES(product_name),
				product_model = VALUES(product_model),
				product_code = VALUES(product_code),
				quantity = VALUES(quantity),
				inventory_status = '板卡入库',
				source_file_id = VALUES(source_file_id),
				source_file_name = VALUES(source_file_name),
				inbound_type = 'board',
				update_time = NOW(),
				is_deleted = 0
		`,
			productName,
			productName,
			productModel,
			productCode,
			quantity,
			sn,
			macAddress,
			req.SourceFileID,
			req.FileName,
		)

		if err != nil {
			http.Error(w, "导入板卡入库失败: "+err.Error(), http.StatusInternalServerError)
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

func DeleteBoardInboundHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if !hasRequestRole(r, "production_staff") && !hasRequestRole(r, "system_admin") && !hasRequestPermission(r, "board-inbound:delete") {
		http.Error(w, "无板卡入库删除权限", http.StatusForbidden)
		return
	}

	id, err := strconv.ParseInt(strings.TrimSpace(r.URL.Query().Get("id")), 10, 64)
	if err != nil || id <= 0 {
		http.Error(w, "板卡入库ID错误", http.StatusBadRequest)
		return
	}

	ensureInventoryBoardColumns()

	result, err := config.DB.Exec(`
		UPDATE inventory_devices
		SET is_deleted = 1,
			update_time = NOW()
		WHERE id = ?
		  AND IFNULL(is_deleted, 0) = 0
		  AND IFNULL(inbound_type, '') = 'board'
		  AND IFNULL(inventory_status, '') = '板卡入库'
	`, id)
	if err != nil {
		http.Error(w, "删除板卡入库失败: "+err.Error(), http.StatusInternalServerError)
		return
	}

	affected, _ := result.RowsAffected()
	if affected == 0 {
		http.Error(w, "删除失败：记录不存在或已经进入后续流程", http.StatusBadRequest)
		return
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"code": 200,
		"msg":  "删除成功",
	})
}

type boardCompositionForDeduct struct {
	ID            int64
	ProjectName   string
	ProductName   string
	InboundModel  string
	OutboundModel string
}

func getBoardCompositionsForBurn(tx *sql.Tx, projectID int64, outboundModel string) ([]boardCompositionForDeduct, error) {
	ensureBoardCompositionTables()
	outboundModel = strings.TrimSpace(outboundModel)
	if outboundModel == "" || outboundModel == "-" {
		return nil, nil
	}

	rows, err := tx.Query(`
		SELECT
			id,
			IFNULL(project_name, ''),
			IFNULL(product_name, ''),
			IFNULL(inbound_model, ''),
			IFNULL(outbound_model, '')
		FROM board_compositions
		WHERE IFNULL(is_deleted, 0) = 0
		  AND outbound_model = ?
		  AND (? = 0 OR project_id = 0 OR project_id = ?)
		ORDER BY id ASC
	`, outboundModel, projectID, projectID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	list := make([]boardCompositionForDeduct, 0)
	for rows.Next() {
		var item boardCompositionForDeduct
		if err := rows.Scan(&item.ID, &item.ProjectName, &item.ProductName, &item.InboundModel, &item.OutboundModel); err != nil {
			return nil, err
		}
		list = append(list, item)
	}
	return list, rows.Err()
}

func deductBoardInboundByComposition(tx *sql.Tx, burnRecordID int64, projectID int64, outboundModel string) (int, error) {
	compositions, err := getBoardCompositionsForBurn(tx, projectID, outboundModel)
	if err != nil {
		return 0, err
	}
	if len(compositions) == 0 {
		return 0, nil
	}

	for _, composition := range compositions {
		if err := deductOneBoardComponent(tx, burnRecordID, composition); err != nil {
			return 0, err
		}
	}
	return len(compositions), nil
}

func ensureBoardCompositionDeductedForFactoryTest(tx *sql.Tx, factoryTestID int64) error {
	ensureBoardCompositionTables()

	var burnRecordID int64
	var projectID int64
	var productModel string
	err := tx.QueryRow(`
		SELECT
			IFNULL(br.id, 0),
			IFNULL(br.project_id, 0),
			IFNULL(br.product_model, '')
		FROM factory_tests ft
		INNER JOIN burn_records br ON br.id = ft.burn_record_id
		WHERE ft.id = ?
		  AND IFNULL(ft.is_deleted, 0) = 0
		  AND IFNULL(br.is_deleted, 0) = 0
		LIMIT 1
	`, factoryTestID).Scan(&burnRecordID, &projectID, &productModel)
	if err != nil {
		return err
	}
	if burnRecordID <= 0 {
		return nil
	}

	var existing int
	err = tx.QueryRow(`
		SELECT COUNT(1)
		FROM board_composition_deductions
		WHERE burn_record_id = ?
	`, burnRecordID).Scan(&existing)
	if err != nil {
		return err
	}
	if existing > 0 {
		return nil
	}

	_, err = deductBoardInboundByComposition(tx, burnRecordID, projectID, productModel)
	return err
}

func BackfillBoardCompositionDeductionsForFactoryInventory() {
	ensureBoardCompositionTables()
	ensureInventoryBoardColumns()
	backfillBoardCompositionDeductionsByInventoryModels()

	rows, err := config.DB.Query(`
		SELECT DISTINCT ft.id
		FROM factory_tests ft
		INNER JOIN burn_records br ON br.id = ft.burn_record_id
		INNER JOIN inventory_devices inv
			ON (
				inv.factory_test_id = ft.id
				OR inv.source_burn_record_id = br.id
				OR (IFNULL(inv.sn, '') <> '' AND inv.sn = br.sn)
				OR (IFNULL(inv.mac_address, '') <> '' AND inv.mac_address = br.mac_address)
			)
		WHERE IFNULL(ft.is_deleted, 0) = 0
		  AND IFNULL(br.is_deleted, 0) = 0
		  AND IFNULL(inv.is_deleted, 0) = 0
		  AND IFNULL(inv.inbound_type, '') = 'factory'
		  AND IFNULL(inv.inventory_status, '') NOT IN ('已出库', '板卡入库', '已烧录')
		  AND ft.audit_status IN ('approved', '审核通过', '已通过')
		  AND NOT EXISTS (
			  SELECT 1
			  FROM board_composition_deductions d
			  WHERE d.burn_record_id = br.id
		  )
		ORDER BY ft.id ASC
	`)
	if err != nil {
		log.Printf("板卡组成历史补扣查询失败: %v", err)
		return
	}
	defer rows.Close()

	factoryTestIDs := make([]int64, 0)
	for rows.Next() {
		var id int64
		if err := rows.Scan(&id); err != nil {
			log.Printf("板卡组成历史补扣数据解析失败: %v", err)
			return
		}
		factoryTestIDs = append(factoryTestIDs, id)
	}
	if err := rows.Err(); err != nil {
		log.Printf("板卡组成历史补扣读取失败: %v", err)
		return
	}
	if len(factoryTestIDs) == 0 {
		return
	}

	successCount := 0
	skipCount := 0
	for _, factoryTestID := range factoryTestIDs {
		tx, err := config.DB.Begin()
		if err != nil {
			log.Printf("板卡组成历史补扣开启事务失败，出厂测试ID=%d: %v", factoryTestID, err)
			skipCount++
			continue
		}

		err = ensureBoardCompositionDeductedForFactoryTest(tx, factoryTestID)
		if err != nil {
			tx.Rollback()
			if err == sql.ErrNoRows {
				log.Printf("板卡组成历史补扣跳过，板卡入库库存不足，出厂测试ID=%d", factoryTestID)
			} else {
				log.Printf("板卡组成历史补扣失败，出厂测试ID=%d: %v", factoryTestID, err)
			}
			skipCount++
			continue
		}
		if err := tx.Commit(); err != nil {
			log.Printf("板卡组成历史补扣提交失败，出厂测试ID=%d: %v", factoryTestID, err)
			skipCount++
			continue
		}
		successCount++
	}

	log.Printf("板卡组成历史补扣完成：成功 %d 条，跳过 %d 条", successCount, skipCount)
}

func backfillBoardCompositionDeductionsByInventoryModels() {
	rows, err := config.DB.Query(`
		SELECT
			CASE
				WHEN IFNULL(inv.source_burn_record_id, 0) > 0 THEN inv.source_burn_record_id
				WHEN IFNULL(br.id, 0) > 0 THEN br.id
				ELSE -inv.id
			END AS ledger_burn_id,
			IFNULL(inv.project_id, IFNULL(br.project_id, 0)) AS project_id,
			IFNULL(inv.product_model, '') AS outbound_model
		FROM inventory_devices inv
		LEFT JOIN burn_records br
			ON IFNULL(br.is_deleted, 0) = 0
		   AND (
				br.id = inv.source_burn_record_id
				OR (IFNULL(inv.sn, '') <> '' AND br.sn = inv.sn)
				OR (IFNULL(inv.mac_address, '') <> '' AND br.mac_address = inv.mac_address)
		   )
		WHERE IFNULL(inv.is_deleted, 0) = 0
		  AND IFNULL(inv.inbound_type, '') <> 'board'
		  AND IFNULL(inv.inventory_status, '') NOT IN ('已出库', '板卡入库', '已烧录', '已废弃', '已报废')
		  AND IFNULL(inv.product_model, '') <> ''
		  AND EXISTS (
			  SELECT 1
			  FROM board_compositions bc
			  WHERE IFNULL(bc.is_deleted, 0) = 0
			    AND CONVERT(bc.outbound_model USING utf8mb4) COLLATE utf8mb4_general_ci
			        = CONVERT(inv.product_model USING utf8mb4) COLLATE utf8mb4_general_ci
		  )
		  AND NOT EXISTS (
			  SELECT 1
			  FROM board_composition_deductions d
			  WHERE d.burn_record_id = CASE
				  WHEN IFNULL(inv.source_burn_record_id, 0) > 0 THEN inv.source_burn_record_id
				  WHEN IFNULL(br.id, 0) > 0 THEN br.id
				  ELSE -inv.id
			  END
		  )
		GROUP BY ledger_burn_id, project_id, outbound_model
		ORDER BY ledger_burn_id ASC
	`)
	if err != nil {
		log.Printf("库存型号历史补扣查询失败: %v", err)
		return
	}
	defer rows.Close()

	type backfillItem struct {
		LedgerBurnID  int64
		ProjectID     int64
		OutboundModel string
	}
	items := make([]backfillItem, 0)
	for rows.Next() {
		var item backfillItem
		if err := rows.Scan(&item.LedgerBurnID, &item.ProjectID, &item.OutboundModel); err != nil {
			log.Printf("库存型号历史补扣数据解析失败: %v", err)
			return
		}
		items = append(items, item)
	}
	if err := rows.Err(); err != nil {
		log.Printf("库存型号历史补扣读取失败: %v", err)
		return
	}
	if len(items) == 0 {
		return
	}

	successCount := 0
	skipCount := 0
	for _, item := range items {
		tx, err := config.DB.Begin()
		if err != nil {
			log.Printf("库存型号历史补扣开启事务失败，出库型号=%s: %v", item.OutboundModel, err)
			skipCount++
			continue
		}

		_, err = deductBoardInboundByComposition(tx, item.LedgerBurnID, item.ProjectID, item.OutboundModel)
		if err != nil {
			tx.Rollback()
			if err == sql.ErrNoRows {
				log.Printf("库存型号历史补扣跳过，板卡入库库存不足，出库型号=%s", item.OutboundModel)
			} else {
				log.Printf("库存型号历史补扣失败，出库型号=%s: %v", item.OutboundModel, err)
			}
			skipCount++
			continue
		}
		if err := tx.Commit(); err != nil {
			log.Printf("库存型号历史补扣提交失败，出库型号=%s: %v", item.OutboundModel, err)
			skipCount++
			continue
		}
		successCount++
	}

	log.Printf("库存型号历史补扣完成：成功 %d 条，跳过 %d 条", successCount, skipCount)
}

func deductOneBoardComponent(tx *sql.Tx, burnRecordID int64, composition boardCompositionForDeduct) error {
	rows, err := tx.Query(`
		SELECT id, IFNULL(quantity, 1)
		FROM inventory_devices
		WHERE IFNULL(is_deleted, 0) = 0
		  AND IFNULL(inbound_type, '') = 'board'
		  AND IFNULL(inventory_status, '') = '板卡入库'
		  AND product_model = ?
		  AND IFNULL(quantity, 0) > 0
		ORDER BY in_time ASC, id ASC
	`, composition.InboundModel)
	if err != nil {
		return err
	}

	type inventoryQty struct {
		ID       int64
		Quantity int
	}
	candidates := make([]inventoryQty, 0)
	for rows.Next() {
		var item inventoryQty
		if err := rows.Scan(&item.ID, &item.Quantity); err != nil {
			rows.Close()
			return err
		}
		candidates = append(candidates, item)
	}
	if err := rows.Err(); err != nil {
		rows.Close()
		return err
	}
	rows.Close()

	remaining := 1
	for _, candidate := range candidates {
		if remaining <= 0 {
			break
		}
		if candidate.Quantity <= 0 {
			continue
		}
		deductQty := 1
		if candidate.Quantity < deductQty {
			deductQty = candidate.Quantity
		}

		_, err = tx.Exec(`
			UPDATE inventory_devices
			SET quantity = quantity - ?,
				inventory_status = CASE WHEN quantity - ? <= 0 THEN '已烧录' ELSE inventory_status END,
				update_time = NOW(),
				remark = CASE
					WHEN IFNULL(remark, '') = '' THEN ?
					ELSE CONCAT(remark, '；', ?)
				END
			WHERE id = ?
		`, deductQty, deductQty, "板卡组成扣减："+composition.OutboundModel, "板卡组成扣减："+composition.OutboundModel, candidate.ID)
		if err != nil {
			return err
		}
		_, err = tx.Exec(`
			INSERT INTO board_composition_deductions (
				burn_record_id,
				composition_id,
				inventory_device_id,
				quantity,
				created_at
			) VALUES (?, ?, ?, ?, NOW())
		`, burnRecordID, composition.ID, candidate.ID, deductQty)
		if err != nil {
			return err
		}
		remaining -= deductQty
	}
	if remaining > 0 {
		return sql.ErrNoRows
	}
	return nil
}

func restoreBoardCompositionDeductions(tx *sql.Tx, burnWhereSQL string, args ...interface{}) error {
	ensureBoardCompositionTables()
	query := `
		SELECT d.inventory_device_id, SUM(d.quantity)
		FROM board_composition_deductions d
		JOIN burn_records br ON br.id = d.burn_record_id
		WHERE ` + burnWhereSQL + `
		GROUP BY d.inventory_device_id
	`
	rows, err := tx.Query(query, args...)
	if err != nil {
		return err
	}
	defer rows.Close()

	type restoreItem struct {
		InventoryID int64
		Quantity    int
	}
	items := make([]restoreItem, 0)
	for rows.Next() {
		var item restoreItem
		if err := rows.Scan(&item.InventoryID, &item.Quantity); err != nil {
			return err
		}
		items = append(items, item)
	}
	if err := rows.Err(); err != nil {
		return err
	}

	for _, item := range items {
		_, err := tx.Exec(`
			UPDATE inventory_devices
			SET quantity = quantity + ?,
				inventory_status = '板卡入库',
				update_time = NOW(),
				remark = CASE
					WHEN IFNULL(remark, '') = '' THEN '烧录删除后恢复板卡组成扣减'
					ELSE CONCAT(remark, '；烧录删除后恢复板卡组成扣减')
				END
			WHERE id = ?
		`, item.Quantity, item.InventoryID)
		if err != nil {
			return err
		}
	}

	deleteSQL := `
		DELETE d
		FROM board_composition_deductions d
		JOIN burn_records br ON br.id = d.burn_record_id
		WHERE ` + burnWhereSQL
	_, err = tx.Exec(deleteSQL, args...)
	return err
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
		"返厂":  true,
		"更换":  true,
		"已废弃": true,
		"已报废": true,

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

func SubmitInventoryScrapHandler(w http.ResponseWriter, r *http.Request, id int64) {
	w.Header().Set("Content-Type", "application/json")

	if !hasRequestRole(r, "production_staff") && !hasRequestRole(r, "system_admin") && !hasRequestPermission(r, "production:update") {
		http.Error(w, "无废弃申请权限：只有生产人员可以提交", http.StatusForbidden)
		return
	}

	ensureInventoryBoardColumns()
	userID, userName := currentRequestUser(r)
	if strings.TrimSpace(userName) == "" {
		userName = "生产人员"
	}

	result, err := config.DB.Exec(`
		UPDATE inventory_devices
		SET
			scrap_audit_status = '待审核',
			scrap_request_user_id = ?,
			scrap_request_user_name = ?,
			scrap_request_time = NOW(),
			scrap_audit_user_id = 0,
			scrap_audit_user_name = '',
			scrap_audit_time = NULL,
			scrap_reject_reason = '',
			update_time = NOW()
		WHERE id = ?
		  AND IFNULL(is_deleted, 0) = 0
		  AND IFNULL(inventory_status, '') NOT IN ('已废弃', '已报废', '已出库')
	`,
		userID,
		userName,
		id,
	)
	if err != nil {
		http.Error(w, "提交废弃申请失败: "+err.Error(), http.StatusInternalServerError)
		return
	}

	affected, _ := result.RowsAffected()
	if affected == 0 {
		http.Error(w, "库存设备不存在，或当前状态不允许废弃", http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"code": 200,
		"msg":  "废弃申请已提交",
	})
}

func AuditInventoryScrapHandler(w http.ResponseWriter, r *http.Request, id int64) {
	w.Header().Set("Content-Type", "application/json")

	if !hasLeaderPermission(r) {
		http.Error(w, "无废弃审核权限：只有领导可以审核", http.StatusForbidden)
		return
	}

	var req struct {
		AuditStatus  string `json:"auditStatus"`
		Status       string `json:"status"`
		RejectReason string `json:"rejectReason"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "参数解析失败: "+err.Error(), http.StatusBadRequest)
		return
	}

	status := strings.TrimSpace(req.AuditStatus)
	if status == "" {
		status = strings.TrimSpace(req.Status)
	}
	approved := status == "approved" || status == "已通过" || status == "审核通过"
	rejected := status == "rejected" || status == "已驳回" || status == "审核驳回"
	if !approved && !rejected {
		http.Error(w, "审核状态不合法", http.StatusBadRequest)
		return
	}

	ensureInventoryBoardColumns()
	auditUserID, auditUserName := currentRequestUser(r)
	if auditUserName == "" {
		auditUserName = "领导"
	}

	auditStatus := "已通过"
	inventoryStatusSQL := "inventory_status"
	if approved {
		inventoryStatusSQL = "'已废弃'"
	} else {
		auditStatus = "已驳回"
	}

	result, err := config.DB.Exec(`
		UPDATE inventory_devices
		SET
			inventory_status = `+inventoryStatusSQL+`,
			scrap_audit_status = ?,
			scrap_audit_user_id = ?,
			scrap_audit_user_name = ?,
			scrap_audit_time = NOW(),
			scrap_reject_reason = ?,
			update_time = NOW()
		WHERE id = ?
		  AND IFNULL(is_deleted, 0) = 0
		  AND IFNULL(scrap_audit_status, '') = '待审核'
	`,
		auditStatus,
		auditUserID,
		auditUserName,
		strings.TrimSpace(req.RejectReason),
		id,
	)
	if err != nil {
		http.Error(w, "废弃审核失败: "+err.Error(), http.StatusInternalServerError)
		return
	}

	affected, _ := result.RowsAffected()
	if affected == 0 {
		http.Error(w, "库存设备不存在，或没有待审核的废弃申请", http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"code": 200,
		"msg":  "审核完成",
	})
}

func ImportBurnRecordsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var req struct {
		SourceFileID    int64  `json:"sourceFileId"`
		FileName        string `json:"fileName"`
		FileContentType string `json:"fileContentType"`
		FileData        string `json:"fileData"`
		Records         []struct {
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

	if err := saveUploadedFile(UploadedFilePayload{
		FileID:          req.SourceFileID,
		FileName:        req.FileName,
		FileContentType: req.FileContentType,
		FileData:        req.FileData,
	}); err != nil {
		http.Error(w, "保存烧录源文件失败: "+err.Error(), http.StatusBadRequest)
		return
	}
	ensureInventoryBoardColumns()
	ensureBoardCompositionTables()

	tx, err := config.DB.Begin()
	if err != nil {
		http.Error(w, "开启事务失败: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer tx.Rollback()

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

		sourceFileID := item.SourceFileID
		if sourceFileID == 0 {
			sourceFileID = req.SourceFileID
		}

		productName := strings.TrimSpace(item.ProductName)
		productModel := strings.TrimSpace(item.ProductModel)
		productCode := strings.TrimSpace(item.ProductCode)
		handset := isHandsetProduct(productName, productModel)
		if handset {
			if productName == "" || productModel == "" || productCode == "" {
				http.Error(w, "手持话柄烧录记录必须包含产品名称、产品型号、产品编码和序列号", http.StatusBadRequest)
				return
			}
			macAddress = ""
			item.HardwareID = 0
			item.HardwareVersion = ""
			item.SoftwareID = 0
			item.SoftwareVersion = ""
			pcbQrCode = ""
		}

		result, err := tx.Exec(`
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
			productName,
			productModel,
			productCode,
			deviceType,
			sn,
			macAddress,
			item.HardwareID,
			item.HardwareVersion,
			item.SoftwareID,
			item.SoftwareVersion,
			pcbQrCode,
			item.Note,
			sourceFileID,
			item.UploaderID,
			uploaderName,
			burnDesc,
		)

		if err != nil {
			http.Error(w, "导入失败，SN或MAC可能重复: "+err.Error(), http.StatusInternalServerError)
			return
		}

		if handset {
			burnRecordID, err := result.LastInsertId()
			if err != nil {
				http.Error(w, "读取手持话柄烧录记录ID失败: "+err.Error(), http.StatusInternalServerError)
				return
			}
			if err := deductHandsetBoardInbound(tx, burnRecordID, productName, productModel); err != nil {
				if err == sql.ErrNoRows {
					http.Error(w, "手持话柄板卡入库数量不足，无法完成烧录入库", http.StatusBadRequest)
					return
				}
				http.Error(w, "扣减手持话柄板卡库存失败: "+err.Error(), http.StatusInternalServerError)
				return
			}
			if err := createHandsetInventory(tx, burnRecordID, item.ProjectID, productName, productModel, productCode, sn, sourceFileID, req.FileName); err != nil {
				http.Error(w, "手持话柄烧录后入库失败: "+err.Error(), http.StatusInternalServerError)
				return
			}
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

func isHandsetProduct(productName string, productModel string) bool {
	name := strings.ReplaceAll(strings.TrimSpace(productName), " ", "")
	model := strings.ToLower(strings.TrimSpace(productModel))
	return strings.Contains(name, "手持话柄") || model == "handheld mic-zycoo"
}

func deductHandsetBoardInbound(tx *sql.Tx, burnRecordID int64, productName string, productModel string) error {
	var inventoryID int64
	var quantity int
	err := tx.QueryRow(`
		SELECT id, IFNULL(quantity, 1)
		FROM inventory_devices
		WHERE IFNULL(is_deleted, 0) = 0
		  AND IFNULL(inbound_type, '') = 'board'
		  AND IFNULL(inventory_status, '') = '板卡入库'
		  AND IFNULL(quantity, 0) > 0
		  AND (
			LOWER(TRIM(IFNULL(product_model, ''))) = LOWER(TRIM(?))
			OR REPLACE(IFNULL(product_name, ''), ' ', '') = REPLACE(?, ' ', '')
		  )
		ORDER BY in_time ASC, id ASC
		LIMIT 1
		FOR UPDATE
	`, productModel, productName).Scan(&inventoryID, &quantity)
	if err != nil {
		return err
	}
	_, err = tx.Exec(`
		UPDATE inventory_devices
		SET quantity = quantity - 1,
			inventory_status = CASE WHEN quantity - 1 <= 0 THEN '已烧录' ELSE '板卡入库' END,
			update_time = NOW(),
			remark = CASE
				WHEN IFNULL(remark, '') = '' THEN '手持话柄烧录扣减'
				ELSE CONCAT(remark, '；手持话柄烧录扣减')
			END
		WHERE id = ? AND quantity = ?
	`, inventoryID, quantity)
	if err != nil {
		return err
	}
	_, err = tx.Exec(`
		INSERT INTO board_composition_deductions (
			burn_record_id, composition_id, inventory_device_id, quantity, created_at
		) VALUES (?, 0, ?, 1, NOW())
	`, burnRecordID, inventoryID)
	return err
}

func createHandsetInventory(tx *sql.Tx, burnRecordID int64, projectID int64, productName string, productModel string, productCode string, sn string, sourceFileID int64, sourceFileName string) error {
	_, err := tx.Exec(`
		INSERT INTO inventory_devices (
			project_id, device_type, product_name, product_model, product_code, quantity,
			sn, mac_address, pcb_qr_code, hardware_id, hardware_version, software_id,
			software_version, inventory_status, source_burn_record_id, factory_test_id,
			source_file_id, source_file_name, inbound_type, in_time, update_time, remark, is_deleted
		) VALUES (?, ?, ?, ?, ?, 1, ?, NULL, '', 0, '', 0, '', '在库', ?, 0, ?, ?, 'handset_burn', NOW(), NOW(), '手持话柄烧录完成，无需出厂测试，直接入库', 0)
		ON DUPLICATE KEY UPDATE
			project_id = VALUES(project_id), device_type = VALUES(device_type),
			product_name = VALUES(product_name), product_model = VALUES(product_model),
			product_code = VALUES(product_code), quantity = 1, mac_address = NULL,
			hardware_id = 0, hardware_version = '', software_id = 0, software_version = '',
			inventory_status = '在库', source_burn_record_id = VALUES(source_burn_record_id),
			factory_test_id = 0, source_file_id = VALUES(source_file_id),
			source_file_name = VALUES(source_file_name), inbound_type = 'handset_burn',
			update_time = NOW(), remark = VALUES(remark), is_deleted = 0
	`, projectID, productName, productName, productModel, productCode, sn, burnRecordID, sourceFileID, sourceFileName)
	return err
}
func DeleteBurnBatchHandler(w http.ResponseWriter, r *http.Request, batchNo string) {
	w.Header().Set("Content-Type", "application/json")

	if !hasRequestRole(r, "production_staff") && !hasRequestRole(r, "system_admin") && !hasRequestPermission(r, "burn:deleteBatch") && !hasRequestPermission(r, "production:delete") {
		http.Error(w, "无删除权限：生产烧录批次仅生产人员或管理员可删除", http.StatusForbidden)
		return
	}

	batchNo = strings.TrimSpace(batchNo)

	if batchNo == "" {
		http.Error(w, "生产批次号不能为空", http.StatusBadRequest)
		return
	}

	tx, err := config.DB.Begin()
	if err != nil {
		http.Error(w, "开启删除事务失败: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer tx.Rollback()

	if err := restoreBoardCompositionDeductions(tx, "br.batch_no = ?", batchNo); err != nil {
		http.Error(w, "恢复板卡组成扣减失败: "+err.Error(), http.StatusInternalServerError)
		return
	}

	_, err = tx.Exec(`
		UPDATE inventory_devices inv
		JOIN burn_records br ON (
			inv.source_burn_record_id = br.id
			OR (inv.mac_address <> '' AND br.mac_address <> '' AND inv.mac_address = br.mac_address)
			OR (inv.sn <> '' AND br.sn <> '' AND inv.sn = br.sn)
		)
		SET
			inv.inventory_status = '板卡入库',
			inv.source_burn_record_id = 0,
			inv.update_time = NOW(),
			inv.remark = CASE
				WHEN IFNULL(inv.remark, '') = '' THEN '烧录删除后恢复板卡入库'
				ELSE CONCAT(inv.remark, '；烧录删除后恢复板卡入库')
			END
		WHERE br.batch_no = ?
		  AND IFNULL(inv.inbound_type, '') = 'board'
		  AND IFNULL(inv.inventory_status, '') = '已烧录'
	`, batchNo)
	if err != nil {
		http.Error(w, "恢复板卡入库失败: "+err.Error(), http.StatusInternalServerError)
		return
	}

	inventoryResult, err := tx.Exec(`
		DELETE inv
		FROM inventory_devices inv
		JOIN burn_records br ON (
			inv.source_burn_record_id = br.id
			OR (inv.mac_address <> '' AND br.mac_address <> '' AND inv.mac_address = br.mac_address)
			OR (inv.sn <> '' AND br.sn <> '' AND inv.sn = br.sn)
		)
		WHERE br.batch_no = ?
		  AND NOT (
			IFNULL(inv.inbound_type, '') = 'board'
			AND IFNULL(inv.inventory_status, '') = '板卡入库'
		  )
	`, batchNo)
	if err != nil {
		http.Error(w, "删除关联库存信息失败: "+err.Error(), http.StatusInternalServerError)
		return
	}

	factoryResult, err := tx.Exec(`
		DELETE ft
		FROM factory_tests ft
		JOIN burn_records br ON (
			ft.burn_record_id = br.id
			OR (ft.mac_address <> '' AND br.mac_address <> '' AND ft.mac_address = br.mac_address)
			OR (ft.sn <> '' AND br.sn <> '' AND ft.sn = br.sn)
		)
		WHERE br.batch_no = ?
	`, batchNo)
	if err != nil {
		http.Error(w, "删除关联出厂测试失败: "+err.Error(), http.StatusInternalServerError)
		return
	}

	result, err := tx.Exec(`
		DELETE FROM burn_records
		WHERE batch_no = ?
	`, batchNo)
	if err != nil {
		http.Error(w, "删除批次失败: "+err.Error(), http.StatusInternalServerError)
		return
	}

	affected, _ := result.RowsAffected()
	factoryAffected, _ := factoryResult.RowsAffected()
	inventoryAffected, _ := inventoryResult.RowsAffected()

	if err := tx.Commit(); err != nil {
		http.Error(w, "提交删除事务失败: "+err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"code": 200,
		"msg":  "删除批次成功",
		"data": map[string]interface{}{
			"count":                affected,
			"factoryTestCount":     factoryAffected,
			"inventoryDeviceCount": inventoryAffected,
		},
	})
}

func DeleteBurnRecordHandler(w http.ResponseWriter, r *http.Request, id int64) {
	w.Header().Set("Content-Type", "application/json")

	if !hasRequestRole(r, "production_staff") && !hasRequestRole(r, "system_admin") && !hasRequestPermission(r, "burn:deleteBatch") && !hasRequestPermission(r, "production:delete") {
		http.Error(w, "无删除权限：生产烧录记录仅生产人员或管理员可删除", http.StatusForbidden)
		return
	}

	tx, err := config.DB.Begin()
	if err != nil {
		http.Error(w, "开启删除事务失败: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer tx.Rollback()

	if err := restoreBoardCompositionDeductions(tx, "br.id = ?", id); err != nil {
		http.Error(w, "恢复板卡组成扣减失败: "+err.Error(), http.StatusInternalServerError)
		return
	}

	_, err = tx.Exec(`
		UPDATE inventory_devices inv
		JOIN burn_records br ON (
			inv.source_burn_record_id = br.id
			OR (inv.mac_address <> '' AND br.mac_address <> '' AND inv.mac_address = br.mac_address)
			OR (inv.sn <> '' AND br.sn <> '' AND inv.sn = br.sn)
		)
		SET
			inv.inventory_status = '板卡入库',
			inv.source_burn_record_id = 0,
			inv.update_time = NOW(),
			inv.remark = CASE
				WHEN IFNULL(inv.remark, '') = '' THEN '烧录删除后恢复板卡入库'
				ELSE CONCAT(inv.remark, '；烧录删除后恢复板卡入库')
			END
		WHERE br.id = ?
		  AND IFNULL(inv.inbound_type, '') = 'board'
		  AND IFNULL(inv.inventory_status, '') = '已烧录'
	`, id)
	if err != nil {
		http.Error(w, "恢复板卡入库失败: "+err.Error(), http.StatusInternalServerError)
		return
	}

	_, err = tx.Exec(`
		DELETE inv
		FROM inventory_devices inv
		JOIN burn_records br ON (
			inv.source_burn_record_id = br.id
			OR (inv.mac_address <> '' AND br.mac_address <> '' AND inv.mac_address = br.mac_address)
			OR (inv.sn <> '' AND br.sn <> '' AND inv.sn = br.sn)
		)
		WHERE br.id = ?
		  AND NOT (
			IFNULL(inv.inbound_type, '') = 'board'
			AND IFNULL(inv.inventory_status, '') = '板卡入库'
		  )
	`, id)
	if err != nil {
		http.Error(w, "删除关联库存信息失败: "+err.Error(), http.StatusInternalServerError)
		return
	}

	_, err = tx.Exec(`
		DELETE ft
		FROM factory_tests ft
		JOIN burn_records br ON (
			ft.burn_record_id = br.id
			OR (ft.mac_address <> '' AND br.mac_address <> '' AND ft.mac_address = br.mac_address)
			OR (ft.sn <> '' AND br.sn <> '' AND ft.sn = br.sn)
		)
		WHERE br.id = ?
	`, id)
	if err != nil {
		http.Error(w, "删除关联出厂测试失败: "+err.Error(), http.StatusInternalServerError)
		return
	}

	result, err := tx.Exec(`
		DELETE FROM burn_records
		WHERE id = ?
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

	if err := tx.Commit(); err != nil {
		http.Error(w, "提交删除事务失败: "+err.Error(), http.StatusInternalServerError)
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
		FileID          int64  `json:"fileId"`
		FileName        string `json:"fileName"`
		FileContentType string `json:"fileContentType"`
		FileData        string `json:"fileData"`
		Records         []struct {
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

	if err := saveUploadedFile(UploadedFilePayload{
		FileID:          req.FileID,
		FileName:        req.FileName,
		FileContentType: req.FileContentType,
		FileData:        req.FileData,
	}); err != nil {
		http.Error(w, "保存出厂测试文件失败: "+err.Error(), http.StatusBadRequest)
		return
	}

	tx, err := config.DB.Begin()
	if err != nil {
		http.Error(w, "开启事务失败: "+err.Error(), http.StatusInternalServerError)
		return
	}

	count := 0
	skipCount := 0

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

		fileID := item.FileID
		if fileID == 0 {
			fileID = req.FileID
		}

		var existingID int64
		var existingStatus string
		err = tx.QueryRow(`
			SELECT id, IFNULL(audit_status, 'draft')
			FROM factory_tests
			WHERE IFNULL(is_deleted, 0) = 0
			  AND (
				burn_record_id = ?
				OR (? <> '' AND mac_address = ?)
				OR (? <> '' AND sn = ?)
			  )
			ORDER BY id DESC
			LIMIT 1
		`, burnRecordID, macAddress, macAddress, sn, sn).Scan(&existingID, &existingStatus)
		if err != nil && err != sql.ErrNoRows {
			tx.Rollback()
			http.Error(w, "检查重复出厂测试失败: "+err.Error(), http.StatusInternalServerError)
			return
		}

		if existingID > 0 {
			if existingStatus == "draft" || existingStatus == "rejected" || existingStatus == "草稿" || existingStatus == "审核驳回" || existingStatus == "已驳回" {
				_, err = tx.Exec(`
					UPDATE factory_tests
					SET product_model = ?,
						device_type = ?,
						file_id = ?,
						uploader_id = ?,
						uploader_name = ?,
						upload_time = NOW(),
						audit_status = ?,
						reject_reason = '',
						auditor_id = 0,
						auditor_name = '',
						audit_time = NULL,
						remark = ?,
						updated_at = NOW()
					WHERE id = ?
				`, productModel, deviceType, fileID, item.UploaderID, item.UploaderName, auditStatus, item.Remark, existingID)
				if err != nil {
					tx.Rollback()
					http.Error(w, "更新重复出厂测试失败: "+err.Error(), http.StatusInternalServerError)
					return
				}
				count++
			} else {
				skipCount++
			}
			continue
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
			fileID,
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
			"count":     count,
			"skipCount": skipCount,
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

	if !requireFactoryQualityAuditorPermission(w, r) {
		return
	}

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

	req.AuditorID, req.AuditorName = normalizeAuditUser(r, req.AuditorID, req.AuditorName)

	tx, err := config.DB.Begin()
	if err != nil {
		http.Error(w, "开启事务失败: "+err.Error(), http.StatusInternalServerError)
		return
	}

	auditCount := int64(0)
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
		auditCount += affected

		if req.Status == "approved" {
			err = SyncFactoryTestToInventoryTx(tx, id)
			if err != nil {
				tx.Rollback()
				if err == sql.ErrNoRows {
					http.Error(w, "审核通过，但板卡组成对应的板卡入库库存不足", http.StatusBadRequest)
					return
				}
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
		"data": map[string]interface{}{
			"count": auditCount,
		},
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
		  AND NOT (
			REPLACE(IFNULL(product_name, ''), ' ', '') LIKE '%手持话柄%'
			OR LOWER(TRIM(IFNULL(product_model, ''))) = 'handheld mic-zycoo'
		  )
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
			product_code,
			sn,
			mac_address,
			pcb_qr_code,
			hardware_id,
			hardware_version,
			software_id,
			software_version,
			inventory_status,
			source_burn_record_id,
			factory_test_id,
			inbound_type,
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
			IFNULL(br.product_code, ''),
			IFNULL(br.sn, ''),
			IFNULL(br.mac_address, ''),
			IFNULL(br.pcb_qr_code, ''),
			IFNULL(br.hardware_id, 0),
			IFNULL(br.hardware_version, ''),
			IFNULL(br.software_id, 0),
			IFNULL(br.software_version, ''),
			'在库',
			br.id,
			ft.id,
			'factory',
			NOW(),
			NOW(),
			'出厂测试审核通过，自动入库',
			0
		FROM factory_tests ft
		INNER JOIN burn_records br ON br.id = ft.burn_record_id
		WHERE ft.id = ?
		  AND IFNULL(ft.is_deleted, 0) = 0
		  AND IFNULL(br.is_deleted, 0) = 0
		  AND ft.audit_status IN ('approved', '审核通过', '已通过')
		ON DUPLICATE KEY UPDATE
			project_id = VALUES(project_id),
			device_type = VALUES(device_type),
			product_name = VALUES(product_name),
			product_model = VALUES(product_model),
			product_code = VALUES(product_code),
			mac_address = VALUES(mac_address),
			pcb_qr_code = VALUES(pcb_qr_code),
			hardware_id = VALUES(hardware_id),
			hardware_version = VALUES(hardware_version),
			software_id = VALUES(software_id),
			software_version = VALUES(software_version),
			inventory_status = '在库',
			source_burn_record_id = VALUES(source_burn_record_id),
			factory_test_id = VALUES(factory_test_id),
			inbound_type = 'factory',
			update_time = NOW(),
			remark = VALUES(remark),
			is_deleted = 0
	`, factoryTestID)

	if err != nil {
		return err
	}
	return ensureBoardCompositionDeductedForFactoryTest(tx, factoryTestID)
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
	ProductName       string `json:"productName"`
	ProductModel      string `json:"productModel"`
	ProductCode       string `json:"productCode"`
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
			sbd.id,
			sbd.batch_id,
			sbd.inventory_device_id,
			IFNULL(sbd.sn, ''),
			IFNULL(sbd.mac_address, ''),
			IFNULL(sbd.device_type, ''),
			IFNULL(inv.product_name, ''),
			IFNULL(inv.product_model, ''),
			IFNULL(inv.product_code, ''),
			IFNULL(sbd.hardware_version, ''),
			IFNULL(sbd.software_version, ''),
			IFNULL(DATE_FORMAT(sbd.created_at, '%Y-%m-%d %H:%i:%s'), '')
		FROM shipping_batch_devices sbd
		LEFT JOIN inventory_devices inv ON inv.id = sbd.inventory_device_id
		WHERE sbd.batch_id = ?
		  AND IFNULL(sbd.is_deleted, 0) = 0
		ORDER BY sbd.id ASC
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
			&item.ProductName,
			&item.ProductModel,
			&item.ProductCode,
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
