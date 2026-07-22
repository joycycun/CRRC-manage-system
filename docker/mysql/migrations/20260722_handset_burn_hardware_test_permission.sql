-- Migration: handset burn-to-inventory flow and independent hardware-test view permission.

SET NAMES utf8mb4;

CREATE TABLE IF NOT EXISTS schema_migrations (
  version varchar(64) NOT NULL PRIMARY KEY,
  applied_at datetime NOT NULL DEFAULT CURRENT_TIMESTAMP
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

INSERT INTO permissions (permission_code, permission_name, module, description, created_at)
VALUES ('hardware-test:view', '硬件测试查看', '硬件管理', '查看全部硬件测试记录', NOW())
ON DUPLICATE KEY UPDATE
  permission_name = VALUES(permission_name),
  module = VALUES(module),
  description = VALUES(description);

ALTER TABLE inventory_devices
  MODIFY COLUMN mac_address varchar(64) NULL COMMENT 'MAC，可选';

UPDATE inventory_devices
SET inventory_status = '板卡入库',
    inbound_type = 'board',
    quantity = 1,
    update_time = NOW(),
    remark = '手持话柄需完成烧录后进入库存'
WHERE IFNULL(is_deleted, 0) = 0
  AND IFNULL(inbound_type, '') = 'direct'
  AND IFNULL(source_burn_record_id, 0) = 0
  AND IFNULL(inventory_status, '') = '在库'
  AND REPLACE(IFNULL(product_name, ''), ' ', '') LIKE '%手持话柄%';

INSERT INTO schema_migrations (version, applied_at)
VALUES ('20260722_handset_burn_hardware_test_permission', NOW())
ON DUPLICATE KEY UPDATE applied_at = VALUES(applied_at);
