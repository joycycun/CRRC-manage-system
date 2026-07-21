-- Migration: document attachments for usage guides and direct inventory model compatibility.

SET NAMES utf8mb4;

CREATE TABLE IF NOT EXISTS schema_migrations (
  version varchar(64) NOT NULL PRIMARY KEY,
  applied_at datetime NOT NULL DEFAULT CURRENT_TIMESTAMP
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

CREATE TABLE IF NOT EXISTS usage_guide_files (
  id bigint NOT NULL AUTO_INCREMENT,
  guide_id bigint NOT NULL,
  file_id bigint NOT NULL,
  sort_order int NOT NULL DEFAULT 0,
  created_at datetime DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (id),
  UNIQUE KEY uk_usage_guide_file (guide_id, file_id),
  KEY idx_usage_guide_files_guide (guide_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

INSERT IGNORE INTO usage_guide_files (guide_id, file_id, sort_order, created_at)
SELECT guide_id, file_id, sort_order, created_at
FROM usage_guide_images;

INSERT INTO schema_migrations (version, applied_at)
VALUES ('20260722_usage_guide_files', NOW())
ON DUPLICATE KEY UPDATE applied_at = VALUES(applied_at);
