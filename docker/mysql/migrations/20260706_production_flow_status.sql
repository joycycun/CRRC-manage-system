UPDATE inventory_devices inv
JOIN burn_records br ON (
  (inv.mac_address <> '' AND br.mac_address <> '' AND inv.mac_address = br.mac_address)
  OR (inv.sn <> '' AND br.sn <> '' AND inv.sn = br.sn)
)
SET
  inv.inventory_status = '已烧录',
  inv.source_burn_record_id = br.id,
  inv.update_time = NOW(),
  inv.remark = CASE
    WHEN IFNULL(inv.remark, '') = '' THEN '已进入生产烧录'
    WHEN inv.remark LIKE '%已进入生产烧录%' THEN inv.remark
    ELSE CONCAT(inv.remark, '；已进入生产烧录')
  END
WHERE IFNULL(inv.is_deleted, 0) = 0
  AND IFNULL(inv.inbound_type, '') = 'board'
  AND IFNULL(inv.inventory_status, '') = '板卡入库'
  AND IFNULL(br.is_deleted, 0) = 0;

UPDATE inventory_devices inv
JOIN burn_records br ON (
  inv.source_burn_record_id = br.id
  OR (inv.mac_address <> '' AND br.mac_address <> '' AND inv.mac_address = br.mac_address)
  OR (inv.sn <> '' AND br.sn <> '' AND inv.sn = br.sn)
)
JOIN factory_tests ft ON ft.burn_record_id = br.id
SET
  inv.project_id = IFNULL(br.project_id, 0),
  inv.device_type = IFNULL(br.device_type, ''),
  inv.product_name = IFNULL(br.product_name, ''),
  inv.product_model = IFNULL(br.product_model, ''),
  inv.product_code = IFNULL(br.product_code, ''),
  inv.pcb_qr_code = IFNULL(br.pcb_qr_code, ''),
  inv.hardware_id = IFNULL(br.hardware_id, 0),
  inv.hardware_version = IFNULL(br.hardware_version, ''),
  inv.software_id = IFNULL(br.software_id, 0),
  inv.software_version = IFNULL(br.software_version, ''),
  inv.inventory_status = '在库',
  inv.source_burn_record_id = br.id,
  inv.factory_test_id = ft.id,
  inv.inbound_type = 'factory',
  inv.update_time = NOW(),
  inv.remark = '出厂测试审核通过，自动入库',
  inv.is_deleted = 0
WHERE IFNULL(inv.is_deleted, 0) = 0
  AND IFNULL(br.is_deleted, 0) = 0
  AND IFNULL(ft.is_deleted, 0) = 0
  AND ft.audit_status IN ('approved', '审核通过', '已通过');
