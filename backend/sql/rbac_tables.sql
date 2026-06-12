CREATE TABLE IF NOT EXISTS roles (
  id BIGINT PRIMARY KEY AUTO_INCREMENT,
  role_code VARCHAR(64) NOT NULL UNIQUE,
  role_name VARCHAR(64) NOT NULL,
  description VARCHAR(255) DEFAULT '',
  created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
  updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS permissions (
  id BIGINT PRIMARY KEY AUTO_INCREMENT,
  permission_code VARCHAR(128) NOT NULL UNIQUE,
  permission_name VARCHAR(128) NOT NULL,
  module VARCHAR(64) DEFAULT '',
  description VARCHAR(255) DEFAULT '',
  created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS user_roles (
  id BIGINT PRIMARY KEY AUTO_INCREMENT,
  user_id BIGINT NOT NULL,
  role_id BIGINT NOT NULL,
  created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
  UNIQUE KEY uk_user_role (user_id, role_id),
  KEY idx_user_roles_user (user_id),
  KEY idx_user_roles_role (role_id)
);

CREATE TABLE IF NOT EXISTS role_permissions (
  id BIGINT PRIMARY KEY AUTO_INCREMENT,
  role_id BIGINT NOT NULL,
  permission_id BIGINT NOT NULL,
  created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
  UNIQUE KEY uk_role_permission (role_id, permission_id),
  KEY idx_role_permissions_role (role_id),
  KEY idx_role_permissions_permission (permission_id)
);

CREATE TABLE IF NOT EXISTS operation_logs (
  id BIGINT PRIMARY KEY AUTO_INCREMENT,
  user_id BIGINT DEFAULT NULL,
  real_name VARCHAR(64) DEFAULT '',
  module VARCHAR(64) DEFAULT '',
  action VARCHAR(64) DEFAULT '',
  business_id BIGINT DEFAULT NULL,
  project_id BIGINT DEFAULT NULL,
  before_data TEXT,
  after_data TEXT,
  ip_address VARCHAR(64) DEFAULT '',
  result VARCHAR(32) DEFAULT '',
  created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
  KEY idx_operation_logs_user (user_id),
  KEY idx_operation_logs_project (project_id),
  KEY idx_operation_logs_module (module),
  KEY idx_operation_logs_created_at (created_at)
);

ALTER TABLE users
  ADD COLUMN IF NOT EXISTS avatar_url BLOB NULL;

INSERT INTO roles (role_code, role_name, description)
VALUES
  ('system_admin', '系统管理员', '拥有系统全部模块权限'),
  ('project_assistant', '项目助理', '维护项目立项和项目资料'),
  ('software_owner', '软件负责人', '维护软件版本和软件问题闭环'),
  ('hardware_owner', '硬件负责人', '维护硬件版本和硬件检测'),
  ('production_staff', '生产人员', '维护烧录记录和出厂测试'),
  ('shipping_staff', '发货人员', '维护库存、发货批次和出库'),
  ('aftersales_staff', '售后人员', '维护维修记录和故障分析'),
  ('leader', '领导', '审核项目和查看统计看板')
ON DUPLICATE KEY UPDATE
  role_name = VALUES(role_name),
  description = VALUES(description);

INSERT INTO permissions (permission_code, permission_name, module, description)
VALUES
  ('project:view', '查看项目', '项目管理', '查看项目立项数据'),
  ('project:create', '新建项目', '项目管理', '新增项目立项'),
  ('project:audit', '审核项目', '项目管理', '审核项目立项'),
  ('issue:view', '查看问题', '测试管理', '查看问题闭环'),
  ('issue:handle', '处理问题', '测试管理', '处理和关闭问题'),
  ('software:view', '查看软件版本', '软件管理', '查看软件版本'),
  ('software:release', '发布软件版本', '软件管理', '发布软件版本'),
  ('hardware:view', '查看硬件版本', '硬件管理', '查看硬件版本'),
  ('production:view', '查看生产记录', '生产管理', '查看生产数据'),
  ('production:test', '出厂测试', '生产管理', '处理烧录后的出厂测试'),
  ('shipping:view', '查看发货数据', '发货管理', '查看发货和出库数据'),
  ('shipping:audit', '审核发货批次', '发货管理', '审核发货批次'),
  ('aftersales:view', '查看售后数据', '售后管理', '查看维修和故障分析'),
  ('report:view', '查看统计报表', '统计报表', '查看项目进度、版本矩阵和问题统计')
ON DUPLICATE KEY UPDATE
  permission_name = VALUES(permission_name),
  module = VALUES(module),
  description = VALUES(description);

INSERT IGNORE INTO role_permissions (role_id, permission_id)
SELECT r.id, p.id
FROM roles r
JOIN permissions p
WHERE r.role_code = 'system_admin';

INSERT IGNORE INTO user_roles (user_id, role_id)
SELECT u.id, r.id
FROM users u
JOIN roles r ON r.role_code = 'system_admin'
WHERE u.username = 'admin';
