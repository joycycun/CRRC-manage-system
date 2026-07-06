-- 硬件版本已改为只保存硬件更改文档，删除旧的压缩包字段。
SET @zip_file_id_exists := (
  SELECT COUNT(*)
  FROM information_schema.columns
  WHERE table_schema = DATABASE()
    AND table_name = 'hardware_versions'
    AND column_name = 'zip_file_id'
);

SET @drop_zip_file_id_sql := IF(
  @zip_file_id_exists > 0,
  'ALTER TABLE hardware_versions DROP COLUMN zip_file_id',
  'SELECT 1'
);

PREPARE stmt FROM @drop_zip_file_id_sql;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;
