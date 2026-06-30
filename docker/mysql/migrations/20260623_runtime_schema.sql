-- Migration: 20260623_runtime_schema
-- Purpose: upgrade an existing crrc_pm database to the schema expected by the
-- current backend/frontend code. This script is safe to run repeatedly.

SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;

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

CREATE TABLE IF NOT EXISTS uploaded_files (
  id bigint NOT NULL,
  file_name varchar(255) NOT NULL DEFAULT '',
  content_type varchar(128) NOT NULL DEFAULT '',
  file_category varchar(64) NOT NULL DEFAULT '',
  file_path varchar(512) NOT NULL DEFAULT '',
  file_size bigint NOT NULL DEFAULT 0,
  file_data longblob,
  created_at datetime DEFAULT CURRENT_TIMESTAMP,
  updated_at datetime DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

CREATE TABLE IF NOT EXISTS notification_reads (
  id bigint NOT NULL AUTO_INCREMENT,
  user_id bigint NOT NULL DEFAULT 0,
  username varchar(64) NOT NULL DEFAULT '',
  notification_id varchar(128) NOT NULL DEFAULT '',
  read_at datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (id),
  UNIQUE KEY uk_notification_read (user_id, username, notification_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

CREATE TABLE IF NOT EXISTS issue_confirmations (
  id bigint NOT NULL AUTO_INCREMENT,
  issue_id bigint NOT NULL,
  confirm_user_id bigint NOT NULL DEFAULT 0,
  confirm_user_name varchar(64) NOT NULL DEFAULT '',
  confirm_time datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (id),
  UNIQUE KEY uk_issue_confirm_user (issue_id, confirm_user_id, confirm_user_name)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

CREATE TABLE IF NOT EXISTS requirement_change_confirmations (
  id bigint NOT NULL AUTO_INCREMENT,
  requirement_change_id bigint NOT NULL,
  confirm_user_id bigint NOT NULL DEFAULT 0,
  confirm_user_name varchar(64) NOT NULL DEFAULT '',
  confirm_time datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (id),
  UNIQUE KEY uk_requirement_change_confirm_user (requirement_change_id, confirm_user_id, confirm_user_name)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

CREATE TABLE IF NOT EXISTS production_requests (
  id bigint NOT NULL AUTO_INCREMENT,
  requester_id bigint NOT NULL DEFAULT 0,
  requester_name varchar(64) NOT NULL DEFAULT '',
  product_model varchar(128) NOT NULL DEFAULT '',
  device_type varchar(64) NOT NULL DEFAULT '',
  quantity int NOT NULL DEFAULT 0,
  detail text,
  status varchar(32) NOT NULL DEFAULT 'pending',
  confirmer_id bigint NOT NULL DEFAULT 0,
  confirmer_name varchar(64) NOT NULL DEFAULT '',
  confirm_time datetime DEFAULT NULL,
  created_at datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  is_deleted tinyint NOT NULL DEFAULT 0,
  PRIMARY KEY (id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

CREATE TABLE IF NOT EXISTS production_test_outlines (
  id bigint NOT NULL AUTO_INCREMENT,
  hardware_id bigint NOT NULL DEFAULT 0,
  hardware_version varchar(128) NOT NULL DEFAULT '',
  board_models varchar(512) NOT NULL DEFAULT '',
  device_type varchar(128) NOT NULL DEFAULT '',
  file_id bigint NOT NULL DEFAULT 0,
  file_name varchar(255) NOT NULL DEFAULT '',
  uploader_id bigint NOT NULL DEFAULT 0,
  uploader_name varchar(64) NOT NULL DEFAULT '',
  upload_time datetime DEFAULT CURRENT_TIMESTAMP,
  remark text,
  is_deleted tinyint NOT NULL DEFAULT 0,
  created_at datetime DEFAULT CURRENT_TIMESTAMP,
  updated_at datetime DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (id),
  KEY idx_pto_hardware_id (hardware_id),
  KEY idx_pto_upload_time (upload_time)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

CALL crrc_add_column_if_missing('projects', 'proposal_file_name', '`proposal_file_name` varchar(255) DEFAULT ''''');
CALL crrc_add_column_if_missing('projects', 'proposal_content_type', '`proposal_content_type` varchar(128) DEFAULT ''''');
CALL crrc_add_column_if_missing('projects', 'proposal_file_data', '`proposal_file_data` longblob NULL');

CALL crrc_add_column_if_missing('users', 'avatar_url', '`avatar_url` blob NULL');

CALL crrc_add_column_if_missing('production_test_outlines', 'board_models', '`board_models` varchar(512) NOT NULL DEFAULT '''' AFTER `hardware_version`');

CALL crrc_add_column_if_missing('shipping_batches', 'file_id', '`file_id` bigint DEFAULT NULL COMMENT ''发货单文件ID''');

CALL crrc_add_column_if_missing('burn_records', 'source_file_id', '`source_file_id` bigint DEFAULT 0 COMMENT ''来源文件ID''');
CALL crrc_add_column_if_missing('burn_records', 'source_file_name', '`source_file_name` varchar(255) DEFAULT '''' COMMENT ''来源文件名称''');

CALL crrc_add_column_if_missing('factory_tests', 'file_id', '`file_id` bigint DEFAULT 0 COMMENT ''测试文件ID''');

CALL crrc_add_index_if_missing('production_requests', 'idx_pr_status', 'KEY `idx_pr_status` (`status`)');
CALL crrc_add_index_if_missing('production_requests', 'idx_pr_requester', 'KEY `idx_pr_requester` (`requester_id`, `requester_name`)');

INSERT INTO schema_migrations (version, applied_at)
VALUES ('20260623_runtime_schema', NOW())
ON DUPLICATE KEY UPDATE applied_at = VALUES(applied_at);

DROP PROCEDURE IF EXISTS crrc_add_column_if_missing;
DROP PROCEDURE IF EXISTS crrc_add_index_if_missing;

SET FOREIGN_KEY_CHECKS = 1;
