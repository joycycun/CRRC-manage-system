-- Migration: 20260714_production_outline_board_inbound
-- Purpose: bind production test outlines to projects and support board inbound quantity rows.

SET NAMES utf8mb4;

CREATE TABLE IF NOT EXISTS schema_migrations (
  version varchar(64) NOT NULL PRIMARY KEY,
  applied_at datetime NOT NULL DEFAULT CURRENT_TIMESTAMP
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

DELIMITER $$

DROP PROCEDURE IF EXISTS crrc_add_column_if_missing $$
CREATE PROCEDURE crrc_add_column_if_missing(
  IN table_name_param varchar(64),
  IN column_name_param varchar(64),
  IN column_definition_param text
)
BEGIN
  IF NOT EXISTS (
    SELECT 1
    FROM information_schema.columns
    WHERE table_schema = DATABASE()
      AND table_name = table_name_param
      AND column_name = column_name_param
  ) THEN
    SET @alter_sql = CONCAT('ALTER TABLE `', table_name_param, '` ADD COLUMN ', column_definition_param);
    PREPARE alter_stmt FROM @alter_sql;
    EXECUTE alter_stmt;
    DEALLOCATE PREPARE alter_stmt;
  END IF;
END $$

DROP PROCEDURE IF EXISTS crrc_add_index_if_missing $$
CREATE PROCEDURE crrc_add_index_if_missing(
  IN table_name_param varchar(64),
  IN index_name_param varchar(64),
  IN index_definition_param text
)
BEGIN
  IF NOT EXISTS (
    SELECT 1
    FROM information_schema.statistics
    WHERE table_schema = DATABASE()
      AND table_name = table_name_param
      AND index_name = index_name_param
  ) THEN
    SET @index_sql = CONCAT('ALTER TABLE `', table_name_param, '` ADD ', index_definition_param);
    PREPARE index_stmt FROM @index_sql;
    EXECUTE index_stmt;
    DEALLOCATE PREPARE index_stmt;
  END IF;
END $$

DELIMITER ;

CALL crrc_add_column_if_missing('production_test_outlines', 'project_id', '`project_id` bigint NOT NULL DEFAULT 0 AFTER `id`');
CALL crrc_add_index_if_missing('production_test_outlines', 'idx_pto_project_id', 'KEY `idx_pto_project_id` (`project_id`)');

CALL crrc_add_column_if_missing('inventory_devices', 'quantity', '`quantity` int NOT NULL DEFAULT 1 COMMENT ''数量'' AFTER `product_code`');

INSERT INTO schema_migrations (version, applied_at)
VALUES ('20260714_production_outline_board_inbound', NOW())
ON DUPLICATE KEY UPDATE applied_at = VALUES(applied_at);

DROP PROCEDURE IF EXISTS crrc_add_column_if_missing;
DROP PROCEDURE IF EXISTS crrc_add_index_if_missing;
