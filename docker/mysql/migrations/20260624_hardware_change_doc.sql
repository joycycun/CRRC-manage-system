USE `crrc_pm`;

DELIMITER $$

DROP PROCEDURE IF EXISTS migrate_hardware_change_doc $$
CREATE PROCEDURE migrate_hardware_change_doc()
BEGIN
  IF EXISTS (
    SELECT 1
    FROM information_schema.columns
    WHERE table_schema = DATABASE()
      AND table_name = 'hardware_versions'
      AND column_name = 'zip_file_id'
  ) AND NOT EXISTS (
    SELECT 1
    FROM information_schema.columns
    WHERE table_schema = DATABASE()
      AND table_name = 'hardware_versions'
      AND column_name = 'change_doc_file_id'
  ) THEN
    ALTER TABLE hardware_versions
      CHANGE COLUMN zip_file_id change_doc_file_id BIGINT DEFAULT NULL COMMENT '硬件更改文档文件ID';
  END IF;

  IF NOT EXISTS (
    SELECT 1
    FROM information_schema.columns
    WHERE table_schema = DATABASE()
      AND table_name = 'hardware_versions'
      AND column_name = 'change_doc_file_id'
  ) THEN
    ALTER TABLE hardware_versions
      ADD COLUMN change_doc_file_id BIGINT DEFAULT NULL COMMENT '硬件更改文档文件ID' AFTER owner_name;
  END IF;
END $$

CALL migrate_hardware_change_doc() $$
DROP PROCEDURE IF EXISTS migrate_hardware_change_doc $$

DELIMITER ;
