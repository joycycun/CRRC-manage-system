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

CALL crrc_add_column_if_missing('fault_analysis', 'reason', '`reason` varchar(128) DEFAULT '''' COMMENT ''原因'' AFTER `board_type`');

DROP PROCEDURE IF EXISTS crrc_add_column_if_missing;
