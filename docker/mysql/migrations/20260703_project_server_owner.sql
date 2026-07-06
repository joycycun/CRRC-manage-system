DELIMITER $$

DROP PROCEDURE IF EXISTS add_column_if_missing$$
CREATE PROCEDURE add_column_if_missing(
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
END$$

CALL add_column_if_missing('projects', 'has_server', '`has_server` TINYINT NOT NULL DEFAULT 0 COMMENT ''是否有服务器'' AFTER `owner_name`')$$
CALL add_column_if_missing('projects', 'server_owner_id', '`server_owner_id` BIGINT NOT NULL DEFAULT 0 COMMENT ''服务器负责人ID'' AFTER `has_server`')$$
CALL add_column_if_missing('projects', 'server_owner_name', '`server_owner_name` VARCHAR(64) NOT NULL DEFAULT '''' COMMENT ''服务器负责人'' AFTER `server_owner_id`')$$

DROP PROCEDURE IF EXISTS add_column_if_missing$$

DELIMITER ;

INSERT INTO users (
  username,
  password_hash,
  real_name,
  department,
  status,
  created_at,
  updated_at
)
SELECT
  '都俊成',
  'pbkdf2_sha256$120000$JZ4CH/SiZq0yEp1+ycbZPg$e4cp/lhdLxseDQg9EnGKWAnmnklNQ2bPDgp7YmQ88wQ',
  '都俊成',
  '软件研发',
  '启用',
  NOW(),
  NOW()
WHERE NOT EXISTS (
  SELECT 1 FROM users WHERE username = '都俊成'
);

INSERT IGNORE INTO user_roles (user_id, role_id, created_at)
SELECT u.id, r.id, NOW()
FROM users u
JOIN roles r ON r.role_code = 'software_owner'
WHERE u.username = '都俊成';

UPDATE projects p
JOIN users u ON u.username = '都俊成'
SET
  p.server_owner_id = u.id,
  p.server_owner_name = IFNULL(NULLIF(u.real_name, ''), u.username)
WHERE IFNULL(p.has_server, 0) = 1
  AND IFNULL(p.server_owner_id, 0) = 0;
