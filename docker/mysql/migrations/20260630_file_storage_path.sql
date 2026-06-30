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

CREATE TABLE IF NOT EXISTS uploaded_files (
  id BIGINT PRIMARY KEY,
  file_name VARCHAR(255) NOT NULL DEFAULT '',
  content_type VARCHAR(128) NOT NULL DEFAULT '',
  file_category VARCHAR(64) NOT NULL DEFAULT '',
  file_path VARCHAR(512) NOT NULL DEFAULT '',
  file_size BIGINT NOT NULL DEFAULT 0,
  file_data LONGBLOB NULL,
  created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
  updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CALL crrc_add_column_if_missing('uploaded_files', 'file_category', '`file_category` varchar(64) NOT NULL DEFAULT '''' AFTER `content_type`');
CALL crrc_add_column_if_missing('uploaded_files', 'file_path', '`file_path` varchar(512) NOT NULL DEFAULT '''' AFTER `file_category`');
CALL crrc_add_column_if_missing('uploaded_files', 'file_size', '`file_size` bigint NOT NULL DEFAULT 0 AFTER `file_path`');

CALL crrc_add_column_if_missing('projects', 'proposal_file_path', '`proposal_file_path` varchar(512) DEFAULT '''' AFTER `proposal_content_type`');
CALL crrc_add_column_if_missing('projects', 'proposal_file_size', '`proposal_file_size` bigint NOT NULL DEFAULT 0 AFTER `proposal_file_path`');

DROP PROCEDURE IF EXISTS crrc_add_column_if_missing;
