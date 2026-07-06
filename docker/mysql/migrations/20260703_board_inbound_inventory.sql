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

CALL crrc_add_column_if_missing('inventory_devices', 'product_code', '`product_code` varchar(128) DEFAULT '''' AFTER `product_model`');
CALL crrc_add_column_if_missing('inventory_devices', 'pcb_qr_code', '`pcb_qr_code` varchar(255) DEFAULT '''' AFTER `mac_address`');
CALL crrc_add_column_if_missing('inventory_devices', 'source_file_id', '`source_file_id` bigint DEFAULT 0 AFTER `factory_test_id`');
CALL crrc_add_column_if_missing('inventory_devices', 'source_file_name', '`source_file_name` varchar(255) DEFAULT '''' AFTER `source_file_id`');
CALL crrc_add_column_if_missing('inventory_devices', 'inbound_type', '`inbound_type` varchar(32) DEFAULT '''' AFTER `source_file_name`');

DROP PROCEDURE IF EXISTS crrc_add_column_if_missing;
