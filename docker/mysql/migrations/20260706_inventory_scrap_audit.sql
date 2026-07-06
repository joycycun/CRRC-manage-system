USE `crrc_pm`;

DELIMITER $$

DROP PROCEDURE IF EXISTS crrc_add_column_if_missing $$
CREATE PROCEDURE crrc_add_column_if_missing(
  IN p_table_name VARCHAR(64),
  IN p_column_name VARCHAR(64),
  IN p_column_definition TEXT
)
BEGIN
  IF NOT EXISTS (
    SELECT 1
    FROM information_schema.COLUMNS
    WHERE TABLE_SCHEMA = DATABASE()
      AND TABLE_NAME = p_table_name
      AND COLUMN_NAME = p_column_name
  ) THEN
    SET @ddl = CONCAT('ALTER TABLE `', p_table_name, '` ADD COLUMN ', p_column_definition);
    PREPARE stmt FROM @ddl;
    EXECUTE stmt;
    DEALLOCATE PREPARE stmt;
  END IF;
END $$

DELIMITER ;

CALL crrc_add_column_if_missing('inventory_devices', 'scrap_audit_status', '`scrap_audit_status` varchar(32) DEFAULT '''' COMMENT ''废弃审核状态'' AFTER `inbound_type`');
CALL crrc_add_column_if_missing('inventory_devices', 'scrap_request_user_id', '`scrap_request_user_id` bigint DEFAULT 0 COMMENT ''废弃申请人ID'' AFTER `scrap_audit_status`');
CALL crrc_add_column_if_missing('inventory_devices', 'scrap_request_user_name', '`scrap_request_user_name` varchar(64) DEFAULT '''' COMMENT ''废弃申请人'' AFTER `scrap_request_user_id`');
CALL crrc_add_column_if_missing('inventory_devices', 'scrap_request_time', '`scrap_request_time` datetime DEFAULT NULL COMMENT ''废弃申请时间'' AFTER `scrap_request_user_name`');
CALL crrc_add_column_if_missing('inventory_devices', 'scrap_audit_user_id', '`scrap_audit_user_id` bigint DEFAULT 0 COMMENT ''废弃审核人ID'' AFTER `scrap_request_time`');
CALL crrc_add_column_if_missing('inventory_devices', 'scrap_audit_user_name', '`scrap_audit_user_name` varchar(64) DEFAULT '''' COMMENT ''废弃审核人'' AFTER `scrap_audit_user_id`');
CALL crrc_add_column_if_missing('inventory_devices', 'scrap_audit_time', '`scrap_audit_time` datetime DEFAULT NULL COMMENT ''废弃审核时间'' AFTER `scrap_audit_user_name`');
CALL crrc_add_column_if_missing('inventory_devices', 'scrap_reject_reason', '`scrap_reject_reason` varchar(255) DEFAULT '''' COMMENT ''废弃驳回原因'' AFTER `scrap_audit_time`');

DROP PROCEDURE IF EXISTS crrc_add_column_if_missing;
