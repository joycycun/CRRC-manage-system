USE `crrc_pm`;

UPDATE users
SET
  username = '丁宇',
  real_name = '丁宇',
  updated_at = NOW()
WHERE username IN ('丁sir', '丁Sir', '丁SIr')
   OR real_name IN ('丁sir', '丁Sir', '丁SIr');
