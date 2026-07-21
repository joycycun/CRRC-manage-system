-- Migration: hardware version multi-project binding and image-based usage guides.

SET NAMES utf8mb4;

CREATE TABLE IF NOT EXISTS schema_migrations (
  version varchar(64) NOT NULL PRIMARY KEY,
  applied_at datetime NOT NULL DEFAULT CURRENT_TIMESTAMP
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

INSERT INTO permissions (permission_code, permission_name, module, description, created_at)
VALUES ('inventory:view', '库存情况查看', '生产管理', '查看生产管理中的库存情况', NOW())
ON DUPLICATE KEY UPDATE
  permission_name = VALUES(permission_name),
  module = VALUES(module),
  description = VALUES(description);

ALTER TABLE user_permissions
  MODIFY COLUMN permission_code varchar(128)
  CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL;

ALTER TABLE user_permissions
  DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci;

CREATE TABLE IF NOT EXISTS hardware_version_projects (
  hardware_version_id bigint NOT NULL,
  project_id bigint NOT NULL,
  created_at datetime DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (hardware_version_id, project_id),
  KEY idx_hvp_project_id (project_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

INSERT IGNORE INTO hardware_version_projects (hardware_version_id, project_id, created_at)
SELECT id, project_id, NOW()
FROM hardware_versions
WHERE IFNULL(project_id, 0) > 0;

CREATE TABLE IF NOT EXISTS usage_guides (
  id bigint NOT NULL AUTO_INCREMENT,
  title varchar(128) NOT NULL DEFAULT '',
  description text,
  created_by bigint NOT NULL DEFAULT 0,
  created_by_name varchar(64) NOT NULL DEFAULT '',
  is_deleted tinyint NOT NULL DEFAULT 0,
  created_at datetime DEFAULT CURRENT_TIMESTAMP,
  updated_at datetime DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

CREATE TABLE IF NOT EXISTS usage_guide_images (
  id bigint NOT NULL AUTO_INCREMENT,
  guide_id bigint NOT NULL,
  file_id bigint NOT NULL,
  sort_order int NOT NULL DEFAULT 0,
  created_at datetime DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (id),
  KEY idx_usage_guide_images_guide (guide_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

INSERT INTO schema_migrations (version, applied_at)
VALUES ('20260721_hardware_projects_usage_guides', NOW())
ON DUPLICATE KEY UPDATE applied_at = VALUES(applied_at);
