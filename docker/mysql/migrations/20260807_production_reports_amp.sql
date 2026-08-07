CREATE TABLE IF NOT EXISTS schema_migrations (
    version VARCHAR(64) PRIMARY KEY,
    applied_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

ALTER TABLE burn_records
    MODIFY COLUMN sn VARCHAR(128) NULL COMMENT 'SN序列号';

ALTER TABLE burn_records
    ADD COLUMN source_row_no INT NOT NULL DEFAULT 0 AFTER source_file_id;

ALTER TABLE inventory_devices
    MODIFY COLUMN sn VARCHAR(128) NULL COMMENT 'SN';

ALTER TABLE factory_tests
    ADD COLUMN report_file_id BIGINT NOT NULL DEFAULT 0 AFTER file_id,
    ADD COLUMN report_uploader_id BIGINT NOT NULL DEFAULT 0 AFTER report_file_id,
    ADD COLUMN report_uploader_name VARCHAR(64) NOT NULL DEFAULT '' AFTER report_uploader_id,
    ADD COLUMN report_upload_time DATETIME NULL AFTER report_uploader_name;

INSERT INTO schema_migrations (version, applied_at)
VALUES ('20260807_production_reports_amp', NOW())
ON DUPLICATE KEY UPDATE applied_at = VALUES(applied_at);
