-- Migration: 20260715_board_composition
-- Purpose: add board composition mapping and burn deduction ledger.

SET NAMES utf8mb4;

CREATE TABLE IF NOT EXISTS schema_migrations (
  version varchar(64) NOT NULL PRIMARY KEY,
  applied_at datetime NOT NULL DEFAULT CURRENT_TIMESTAMP
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

CREATE TABLE IF NOT EXISTS board_compositions (
  id bigint NOT NULL AUTO_INCREMENT,
  project_id bigint NOT NULL DEFAULT 0,
  project_name varchar(128) NOT NULL DEFAULT '',
  product_name varchar(128) NOT NULL DEFAULT '',
  inbound_model varchar(128) NOT NULL DEFAULT '',
  outbound_model varchar(128) NOT NULL DEFAULT '',
  created_by varchar(64) NOT NULL DEFAULT '',
  is_deleted tinyint NOT NULL DEFAULT 0,
  created_at datetime DEFAULT CURRENT_TIMESTAMP,
  updated_at datetime DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (id),
  KEY idx_bc_project (project_id, project_name),
  KEY idx_bc_outbound (outbound_model),
  KEY idx_bc_inbound (product_name, inbound_model)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

CREATE TABLE IF NOT EXISTS board_composition_deductions (
  id bigint NOT NULL AUTO_INCREMENT,
  burn_record_id bigint NOT NULL DEFAULT 0,
  composition_id bigint NOT NULL DEFAULT 0,
  inventory_device_id bigint NOT NULL DEFAULT 0,
  quantity int NOT NULL DEFAULT 0,
  created_at datetime DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (id),
  KEY idx_bcd_burn_record (burn_record_id),
  KEY idx_bcd_inventory (inventory_device_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

INSERT INTO schema_migrations (version, applied_at)
VALUES ('20260715_board_composition', NOW())
ON DUPLICATE KEY UPDATE applied_at = VALUES(applied_at);
