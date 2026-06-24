/*M!999999\- enable the sandbox mode */ 

/*!40101 SET @OLD_CHARACTER_SET_CLIENT=@@CHARACTER_SET_CLIENT */;
/*!40101 SET @OLD_CHARACTER_SET_RESULTS=@@CHARACTER_SET_RESULTS */;
/*!40101 SET @OLD_COLLATION_CONNECTION=@@COLLATION_CONNECTION */;
/*!40101 SET NAMES utf8mb4 */;
/*!40103 SET @OLD_TIME_ZONE=@@TIME_ZONE */;
/*!40103 SET TIME_ZONE='+00:00' */;
/*!40014 SET @OLD_UNIQUE_CHECKS=@@UNIQUE_CHECKS, UNIQUE_CHECKS=0 */;
/*!40014 SET @OLD_FOREIGN_KEY_CHECKS=@@FOREIGN_KEY_CHECKS, FOREIGN_KEY_CHECKS=0 */;
/*!40101 SET @OLD_SQL_MODE=@@SQL_MODE, SQL_MODE='NO_AUTO_VALUE_ON_ZERO' */;
/*!40111 SET @OLD_SQL_NOTES=@@SQL_NOTES, SQL_NOTES=0 */;

CREATE DATABASE /*!32312 IF NOT EXISTS*/ `crrc_pm` /*!40100 DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci */ /*!80016 DEFAULT ENCRYPTION='N' */;

USE `crrc_pm`;
DROP TABLE IF EXISTS `burn_records`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8mb4 */;
CREATE TABLE `burn_records` (
  `id` bigint NOT NULL AUTO_INCREMENT COMMENT '烧录记录ID',
  `batch_no` varchar(64) COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT '生产/烧录批次',
  `project_id` bigint DEFAULT NULL COMMENT '关联项目',
  `production_order_id` bigint DEFAULT NULL COMMENT '生产工单ID',
  `product_name` varchar(128) COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT '产品名称',
  `product_model` varchar(128) COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT '产品型号',
  `product_code` varchar(64) COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT '产品编码',
  `device_type` varchar(64) COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT '终端类型',
  `sn` varchar(128) COLLATE utf8mb4_general_ci NOT NULL COMMENT 'SN序列号',
  `mac_address` varchar(64) COLLATE utf8mb4_general_ci DEFAULT NULL,
  `hardware_id` bigint DEFAULT NULL COMMENT '硬件版本ID',
  `hardware_version` varchar(64) COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT '硬件版本',
  `software_id` bigint DEFAULT NULL COMMENT '软件版本ID',
  `software_version` varchar(64) COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT '软件版本',
  `pcb_qr_code` varchar(128) COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT 'PCB二维码',
  `note` text COLLATE utf8mb4_general_ci COMMENT '备注',
  `source_file_id` bigint DEFAULT NULL COMMENT '导入文件ID',
  `source_file_name` varchar(255) COLLATE utf8mb4_general_ci DEFAULT '' COMMENT '来源文件名称',
  `uploader_id` bigint DEFAULT NULL COMMENT '上传人ID',
  `uploader_name` varchar(64) COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT '上传人姓名',
  `upload_time` datetime DEFAULT NULL COMMENT '上传时间',
  `burn_desc` text COLLATE utf8mb4_general_ci COMMENT '烧录说明',
  `is_deleted` tinyint DEFAULT '0' COMMENT '是否删除',
  `created_at` datetime DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` datetime DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_burn_sn` (`sn`),
  UNIQUE KEY `uk_burn_mac` (`mac_address`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci COMMENT='烧录记录表';
/*!40101 SET character_set_client = @saved_cs_client */;
DROP TABLE IF EXISTS `customer_supplied_files`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8mb4 */;
CREATE TABLE `customer_supplied_files` (
  `id` bigint NOT NULL AUTO_INCREMENT COMMENT '客供资料ID',
  `project_id` bigint NOT NULL COMMENT '关联项目',
  `file_id` bigint NOT NULL COMMENT '客供资料文件',
  `material_name` varchar(128) COLLATE utf8mb4_general_ci NOT NULL COMMENT '客供资料名称',
  `file_display_name` varchar(128) COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT '页面显示文件名称',
  `material_desc` text COLLATE utf8mb4_general_ci COMMENT '资料说明',
  `upload_user_id` bigint DEFAULT NULL COMMENT '上传人ID',
  `upload_user_name` varchar(64) COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT '上传人姓名',
  `upload_time` datetime DEFAULT NULL COMMENT '上传时间',
  `remark` text COLLATE utf8mb4_general_ci COMMENT '备注',
  `created_at` datetime DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` datetime DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  `is_deleted` tinyint DEFAULT '0' COMMENT '是否删除，0否，1是',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci COMMENT='客供资料业务表';
/*!40101 SET character_set_client = @saved_cs_client */;
DROP TABLE IF EXISTS `factory_tests`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8mb4 */;
CREATE TABLE `factory_tests` (
  `id` bigint NOT NULL AUTO_INCREMENT COMMENT '出厂测试ID',
  `burn_record_id` bigint NOT NULL COMMENT '关联烧录记录',
  `project_id` bigint DEFAULT NULL COMMENT '关联项目',
  `product_model` varchar(128) COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT '产品型号',
  `device_type` varchar(64) COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT '终端类型',
  `mac_address` varchar(64) COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT 'MAC地址',
  `sn` varchar(128) COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT 'SN序列号',
  `file_id` bigint DEFAULT NULL COMMENT '测试文件ID',
  `uploader_id` bigint DEFAULT NULL COMMENT '上传人ID',
  `uploader_name` varchar(64) COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT '上传人姓名',
  `upload_time` datetime DEFAULT NULL COMMENT '上传时间',
  `audit_status` varchar(32) COLLATE utf8mb4_general_ci DEFAULT '待审核' COMMENT '审核状态',
  `reject_reason` text COLLATE utf8mb4_general_ci COMMENT '驳回原因',
  `auditor_id` bigint DEFAULT NULL COMMENT '审核人ID',
  `auditor_name` varchar(64) COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT '审核人姓名',
  `audit_time` datetime DEFAULT NULL COMMENT '审核时间',
  `remark` text COLLATE utf8mb4_general_ci COMMENT '备注',
  `is_deleted` tinyint DEFAULT '0' COMMENT '是否删除',
  `created_at` datetime DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` datetime DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci COMMENT='出厂测试表';
/*!40101 SET character_set_client = @saved_cs_client */;
DROP TABLE IF EXISTS `fault_analysis`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8mb4 */;
CREATE TABLE `fault_analysis` (
  `id` bigint NOT NULL AUTO_INCREMENT COMMENT '故障分析ID',
  `project_id` bigint DEFAULT NULL COMMENT '关联项目',
  `issue_id` bigint DEFAULT NULL COMMENT '关联问题',
  `repair_id` bigint DEFAULT NULL COMMENT '关联维修记录',
  `board_type` varchar(64) COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT '板卡/终端类型',
  `analysis_name` varchar(128) COLLATE utf8mb4_general_ci NOT NULL COMMENT '分析方案名称',
  `file_id` bigint DEFAULT NULL COMMENT '故障分析文件ID',
  `submit_user_id` bigint DEFAULT NULL COMMENT '提交人ID',
  `submit_user_name` varchar(64) COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT '提交人姓名',
  `submit_time` datetime DEFAULT NULL COMMENT '提交时间',
  `audit_status` varchar(32) COLLATE utf8mb4_general_ci DEFAULT '待审核' COMMENT '审核状态',
  `auditor_id` bigint DEFAULT NULL COMMENT '审核人ID',
  `auditor_name` varchar(64) COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT '审核人姓名',
  `audit_time` datetime DEFAULT NULL COMMENT '审核时间',
  `reject_reason` text COLLATE utf8mb4_general_ci COMMENT '驳回原因',
  `analysis_desc` text COLLATE utf8mb4_general_ci COMMENT '故障分析说明',
  `is_deleted` tinyint DEFAULT '0' COMMENT '是否删除',
  `created_at` datetime DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` datetime DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  `file_name` varchar(255) COLLATE utf8mb4_general_ci DEFAULT '',
  `file_url` text COLLATE utf8mb4_general_ci,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci COMMENT='故障分析方案表';
/*!40101 SET character_set_client = @saved_cs_client */;
DROP TABLE IF EXISTS `hardware_tests`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8mb4 */;
CREATE TABLE `hardware_tests` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `project_id` bigint NOT NULL,
  `hardware_id` bigint NOT NULL,
  `record_name` varchar(128) COLLATE utf8mb4_general_ci DEFAULT NULL,
  `device_type` varchar(64) COLLATE utf8mb4_general_ci DEFAULT NULL,
  `file_id` bigint NOT NULL,
  `uploader_id` bigint DEFAULT '0' COMMENT '上传人ID',
  `uploader_name` varchar(64) COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT '上传人姓名',
  `upload_time` datetime DEFAULT NULL COMMENT '上传时间',
  `audit_status` varchar(32) COLLATE utf8mb4_general_ci DEFAULT '待审核',
  `auditor_id` bigint DEFAULT NULL,
  `auditor_name` varchar(64) COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT '审核人姓名',
  `audit_time` datetime DEFAULT NULL,
  `reject_reason` text COLLATE utf8mb4_general_ci,
  `remark` text COLLATE utf8mb4_general_ci,
  `created_at` datetime DEFAULT CURRENT_TIMESTAMP,
  `updated_at` datetime DEFAULT CURRENT_TIMESTAMP,
  `is_deleted` tinyint DEFAULT '0',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;
/*!40101 SET character_set_client = @saved_cs_client */;
DROP TABLE IF EXISTS `hardware_versions`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8mb4 */;
CREATE TABLE `hardware_versions` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `hardware_version` varchar(64) COLLATE utf8mb4_general_ci NOT NULL,
  `project_id` bigint DEFAULT NULL,
  `device_type` varchar(64) COLLATE utf8mb4_general_ci DEFAULT NULL,
  `status` varchar(32) COLLATE utf8mb4_general_ci DEFAULT '样品',
  `owner_id` bigint DEFAULT NULL,
  `owner_name` varchar(64) COLLATE utf8mb4_general_ci DEFAULT NULL,
  `change_doc_file_id` bigint DEFAULT NULL COMMENT '硬件更改文档文件ID',
  `description` text COLLATE utf8mb4_general_ci,
  `created_at` datetime DEFAULT CURRENT_TIMESTAMP,
  `updated_at` datetime DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;
/*!40101 SET character_set_client = @saved_cs_client */;
DROP TABLE IF EXISTS `inventory_devices`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8mb4 */;
CREATE TABLE `inventory_devices` (
  `id` bigint NOT NULL AUTO_INCREMENT COMMENT '库存ID',
  `project_id` bigint DEFAULT NULL COMMENT '关联项目ID',
  `device_type` varchar(64) COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT '终端类型',
  `product_name` varchar(128) COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT '产品名称',
  `product_model` varchar(128) COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT '产品型号',
  `sn` varchar(128) COLLATE utf8mb4_general_ci NOT NULL COMMENT 'SN',
  `mac_address` varchar(64) COLLATE utf8mb4_general_ci NOT NULL COMMENT 'MAC',
  `hardware_id` bigint DEFAULT NULL COMMENT '硬件版本ID',
  `hardware_version` varchar(64) COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT '硬件版本',
  `software_id` bigint DEFAULT NULL COMMENT '软件版本ID',
  `software_version` varchar(64) COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT '软件版本',
  `inventory_status` varchar(32) COLLATE utf8mb4_general_ci DEFAULT '在库' COMMENT '在库/已锁定/已出库/返修',
  `source_burn_record_id` bigint DEFAULT NULL COMMENT '来源烧录记录',
  `factory_test_id` bigint DEFAULT NULL COMMENT '出厂测试记录',
  `in_time` datetime DEFAULT NULL COMMENT '入库时间',
  `update_time` datetime DEFAULT NULL COMMENT '更新时间',
  `remark` text COLLATE utf8mb4_general_ci COMMENT '备注',
  `is_deleted` tinyint DEFAULT '0' COMMENT '是否删除',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_inventory_sn` (`sn`),
  UNIQUE KEY `uk_inventory_mac` (`mac_address`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci COMMENT='库存设备表';
/*!40101 SET character_set_client = @saved_cs_client */;
DROP TABLE IF EXISTS `issue_confirmations`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8mb4 */;
CREATE TABLE `issue_confirmations` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `issue_id` bigint NOT NULL,
  `confirm_user_id` bigint NOT NULL DEFAULT '0',
  `confirm_user_name` varchar(64) COLLATE utf8mb4_general_ci NOT NULL DEFAULT '',
  `confirm_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_issue_confirm_user` (`issue_id`,`confirm_user_id`,`confirm_user_name`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;
/*!40101 SET character_set_client = @saved_cs_client */;
DROP TABLE IF EXISTS `issue_replies`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8mb4 */;
CREATE TABLE `issue_replies` (
  `id` bigint NOT NULL AUTO_INCREMENT COMMENT '回复ID',
  `issue_id` bigint NOT NULL COMMENT '关联问题ID',
  `reply_user_id` bigint DEFAULT NULL COMMENT '回复人ID',
  `reply_user_name` varchar(64) COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT '回复人',
  `reply_time` datetime DEFAULT NULL COMMENT '回复时间',
  `content` text COLLATE utf8mb4_general_ci COMMENT '回复内容',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci COMMENT='问题回复表';
/*!40101 SET character_set_client = @saved_cs_client */;
DROP TABLE IF EXISTS `issues`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8mb4 */;
CREATE TABLE `issues` (
  `id` bigint NOT NULL AUTO_INCREMENT COMMENT '问题ID',
  `project_id` bigint NOT NULL COMMENT '关联项目',
  `device_type` varchar(64) COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT '终端类型',
  `issue_source` varchar(64) COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT '问题来源',
  `level` varchar(32) COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT '严重等级',
  `issue_title` varchar(255) COLLATE utf8mb4_general_ci NOT NULL COMMENT '问题名称',
  `issue_desc` text COLLATE utf8mb4_general_ci COMMENT '问题描述',
  `owner_id` bigint DEFAULT NULL COMMENT '负责人ID',
  `owner_name` varchar(64) COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT '负责人',
  `creator_id` bigint DEFAULT NULL COMMENT '创建人ID',
  `creator_name` varchar(64) COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT '创建人',
  `create_time` datetime DEFAULT NULL COMMENT '提出时间',
  `plan_close_time` datetime DEFAULT NULL COMMENT '计划关闭时间',
  `real_close_time` datetime DEFAULT NULL COMMENT '实际关闭时间',
  `close_status` varchar(32) COLLATE utf8mb4_general_ci DEFAULT '打开' COMMENT '打开/关闭/重开',
  `close_user_id` bigint DEFAULT NULL COMMENT '关闭人ID',
  `close_user_name` varchar(64) COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT '关闭人姓名',
  `reopen_reason` text COLLATE utf8mb4_general_ci COMMENT '重新打开原因',
  `updated_at` datetime DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  `is_deleted` tinyint DEFAULT '0' COMMENT '是否删除',
  `reopen_count` int NOT NULL DEFAULT '0',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci COMMENT='问题闭环表';
/*!40101 SET character_set_client = @saved_cs_client */;
DROP TABLE IF EXISTS `notification_reads`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8mb4 */;
CREATE TABLE `notification_reads` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `user_id` bigint NOT NULL DEFAULT '0',
  `username` varchar(64) COLLATE utf8mb4_general_ci NOT NULL DEFAULT '',
  `notification_id` varchar(160) COLLATE utf8mb4_general_ci NOT NULL,
  `read_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_notification_read` (`user_id`,`username`,`notification_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;
/*!40101 SET character_set_client = @saved_cs_client */;
DROP TABLE IF EXISTS `operation_logs`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8mb4 */;
CREATE TABLE `operation_logs` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `user_id` bigint DEFAULT NULL,
  `real_name` varchar(64) COLLATE utf8mb4_general_ci DEFAULT '',
  `module` varchar(64) COLLATE utf8mb4_general_ci DEFAULT '',
  `action` varchar(64) COLLATE utf8mb4_general_ci DEFAULT '',
  `business_id` bigint DEFAULT NULL,
  `project_id` bigint DEFAULT NULL,
  `before_data` text COLLATE utf8mb4_general_ci,
  `after_data` text COLLATE utf8mb4_general_ci,
  `ip_address` varchar(64) COLLATE utf8mb4_general_ci DEFAULT '',
  `result` varchar(32) COLLATE utf8mb4_general_ci DEFAULT '',
  `created_at` datetime DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  KEY `idx_operation_logs_user` (`user_id`),
  KEY `idx_operation_logs_project` (`project_id`),
  KEY `idx_operation_logs_module` (`module`),
  KEY `idx_operation_logs_created_at` (`created_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;
/*!40101 SET character_set_client = @saved_cs_client */;
DROP TABLE IF EXISTS `outbound_records`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8mb4 */;
CREATE TABLE `outbound_records` (
  `id` bigint NOT NULL AUTO_INCREMENT COMMENT '出库ID',
  `batch_id` bigint NOT NULL COMMENT '发货批次ID',
  `inventory_device_id` bigint NOT NULL COMMENT '库存设备ID',
  `sn` varchar(128) COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT 'SN',
  `mac_address` varchar(64) COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT 'MAC',
  `outbound_time` datetime DEFAULT NULL COMMENT '出库时间',
  `outbound_user_id` bigint DEFAULT NULL COMMENT '出库操作人ID',
  `outbound_user_name` varchar(64) COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT '出库操作人姓名',
  `status` varchar(32) COLLATE utf8mb4_general_ci DEFAULT '已出库' COMMENT '已出库/取消',
  `remark` text COLLATE utf8mb4_general_ci COMMENT '备注',
  `inbound_time` datetime DEFAULT NULL COMMENT '入库时间',
  `is_deleted` tinyint DEFAULT '0' COMMENT '是否删除',
  `created_at` datetime DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci COMMENT='出库记录表';
/*!40101 SET character_set_client = @saved_cs_client */;
DROP TABLE IF EXISTS `permissions`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8mb4 */;
CREATE TABLE `permissions` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `permission_code` varchar(128) COLLATE utf8mb4_general_ci NOT NULL,
  `permission_name` varchar(128) COLLATE utf8mb4_general_ci NOT NULL,
  `module` varchar(64) COLLATE utf8mb4_general_ci DEFAULT '',
  `description` varchar(255) COLLATE utf8mb4_general_ci DEFAULT '',
  `created_at` datetime DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  UNIQUE KEY `permission_code` (`permission_code`)
) ENGINE=InnoDB AUTO_INCREMENT=35 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;
/*!40101 SET character_set_client = @saved_cs_client */;
DROP TABLE IF EXISTS `production_orders`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8mb4 */;
CREATE TABLE `production_orders` (
  `id` bigint NOT NULL AUTO_INCREMENT COMMENT '生产工单ID',
  `order_no` varchar(64) COLLATE utf8mb4_general_ci NOT NULL COMMENT '生产工单号',
  `project_id` bigint DEFAULT NULL COMMENT '关联项目ID',
  `device_type` varchar(64) COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT '终端类型',
  `product_name` varchar(128) COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT '产品名称',
  `product_model` varchar(128) COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT '产品型号',
  `plan_qty` int DEFAULT '0' COMMENT '计划生产数量',
  `hardware_id` bigint DEFAULT NULL COMMENT '硬件版本ID',
  `hardware_version` varchar(64) COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT '硬件版本号',
  `software_id` bigint DEFAULT NULL COMMENT '软件版本ID',
  `software_version` varchar(64) COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT '软件版本号',
  `create_user_id` bigint DEFAULT NULL COMMENT '创建人ID',
  `create_user_name` varchar(64) COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT '创建人姓名',
  `create_time` datetime DEFAULT NULL COMMENT '创建时间',
  `remark` text COLLATE utf8mb4_general_ci COMMENT '备注',
  `updated_at` datetime DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  `is_deleted` tinyint DEFAULT '0' COMMENT '是否删除',
  `status` varchar(32) COLLATE utf8mb4_general_ci DEFAULT '待生产' COMMENT '待生产/生产中/已完成/已关闭',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci COMMENT='生产工单表';
/*!40101 SET character_set_client = @saved_cs_client */;
DROP TABLE IF EXISTS `production_test_outlines`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8mb4 */;
CREATE TABLE `production_test_outlines` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `hardware_id` bigint NOT NULL DEFAULT '0',
  `hardware_version` varchar(128) COLLATE utf8mb4_general_ci NOT NULL DEFAULT '',
  `board_models` varchar(512) COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '板卡型号，多个用逗号分隔',
  `device_type` varchar(128) COLLATE utf8mb4_general_ci NOT NULL DEFAULT '',
  `file_id` bigint NOT NULL DEFAULT '0',
  `file_name` varchar(255) COLLATE utf8mb4_general_ci NOT NULL DEFAULT '',
  `uploader_id` bigint NOT NULL DEFAULT '0',
  `uploader_name` varchar(64) COLLATE utf8mb4_general_ci NOT NULL DEFAULT '',
  `upload_time` datetime DEFAULT CURRENT_TIMESTAMP,
  `remark` text COLLATE utf8mb4_general_ci,
  `is_deleted` tinyint NOT NULL DEFAULT '0',
  `created_at` datetime DEFAULT CURRENT_TIMESTAMP,
  `updated_at` datetime DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  KEY `idx_pto_hardware_id` (`hardware_id`),
  KEY `idx_pto_upload_time` (`upload_time`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci COMMENT='生产测试大纲表';
/*!40101 SET character_set_client = @saved_cs_client */;
DROP TABLE IF EXISTS `production_requests`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8mb4 */;
CREATE TABLE `production_requests` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `requester_id` bigint NOT NULL DEFAULT '0',
  `requester_name` varchar(64) COLLATE utf8mb4_general_ci NOT NULL DEFAULT '',
  `product_model` varchar(128) COLLATE utf8mb4_general_ci NOT NULL DEFAULT '',
  `device_type` varchar(64) COLLATE utf8mb4_general_ci NOT NULL DEFAULT '',
  `quantity` int NOT NULL DEFAULT '0',
  `detail` text COLLATE utf8mb4_general_ci,
  `status` varchar(32) COLLATE utf8mb4_general_ci NOT NULL DEFAULT 'pending',
  `confirmer_id` bigint NOT NULL DEFAULT '0',
  `confirmer_name` varchar(64) COLLATE utf8mb4_general_ci NOT NULL DEFAULT '',
  `confirm_time` datetime DEFAULT NULL,
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  `is_deleted` tinyint NOT NULL DEFAULT '0',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;
/*!40101 SET character_set_client = @saved_cs_client */;
DROP TABLE IF EXISTS `project_branches`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8mb4 */;
CREATE TABLE `project_branches` (
  `id` bigint NOT NULL AUTO_INCREMENT COMMENT '主键',
  `project_id` bigint NOT NULL COMMENT '关联项目',
  `device_type` varchar(64) COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT '终端类型',
  `repo_name` varchar(128) COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT '代码库名称',
  `repo_url` varchar(512) COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT '代码库地址',
  `branch_name` varchar(128) COLLATE utf8mb4_general_ci NOT NULL COMMENT '分支名称',
  `clone_url` varchar(512) COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT '克隆地址',
  `owner_id` bigint DEFAULT NULL COMMENT '软件负责人ID',
  `owner_name` varchar(64) COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT '软件负责人',
  `remark` text COLLATE utf8mb4_general_ci COMMENT '备注',
  `created_at` datetime DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` datetime DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  `is_deleted` tinyint DEFAULT '0' COMMENT '是否删除，0否，1是',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci COMMENT='项目分支绑定表';
/*!40101 SET character_set_client = @saved_cs_client */;
DROP TABLE IF EXISTS `projects`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8mb4 */;
CREATE TABLE `projects` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `project_name` varchar(128) COLLATE utf8mb4_general_ci NOT NULL,
  `project_code` varchar(64) COLLATE utf8mb4_general_ci NOT NULL,
  `owner_id` bigint NOT NULL,
  `owner_name` varchar(64) COLLATE utf8mb4_general_ci NOT NULL,
  `stage` varchar(64) COLLATE utf8mb4_general_ci DEFAULT NULL,
  `status` varchar(32) COLLATE utf8mb4_general_ci DEFAULT NULL,
  `submit_time` datetime DEFAULT NULL,
  `audit_status` varchar(32) COLLATE utf8mb4_general_ci DEFAULT NULL,
  `audit_user_id` bigint DEFAULT NULL,
  `audit_user_name` varchar(64) COLLATE utf8mb4_general_ci DEFAULT NULL,
  `audit_time` datetime DEFAULT NULL,
  `archive_time` datetime DEFAULT NULL,
  `remark` text COLLATE utf8mb4_general_ci,
  `created_at` datetime DEFAULT CURRENT_TIMESTAMP,
  `updated_at` datetime DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  `close_time` datetime DEFAULT NULL,
  `is_deleted` tinyint DEFAULT '0' COMMENT '是否删除',
  `proposal_file_name` varchar(255) COLLATE utf8mb4_general_ci DEFAULT '',
  `proposal_content_type` varchar(128) COLLATE utf8mb4_general_ci DEFAULT '',
  `proposal_file_data` longblob,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=5 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;
/*!40101 SET character_set_client = @saved_cs_client */;
DROP TABLE IF EXISTS `repair_records`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8mb4 */;
CREATE TABLE `repair_records` (
  `id` bigint NOT NULL AUTO_INCREMENT COMMENT '维修记录ID',
  `project_id` bigint DEFAULT NULL COMMENT '关联项目',
  `device_type` varchar(64) COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT '终端类型',
  `inventory_device_id` bigint DEFAULT NULL COMMENT '关联库存设备',
  `sn` varchar(128) COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT 'SN',
  `mac_address` varchar(64) COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT 'MAC',
  `fault_desc` text COLLATE utf8mb4_general_ci COMMENT '故障现象',
  `repair_user_id` bigint DEFAULT NULL COMMENT '维修人ID',
  `repair_user_name` varchar(64) COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT '维修人',
  `repair_time` datetime DEFAULT NULL COMMENT '维修时间',
  `repair_finish_time` datetime DEFAULT NULL COMMENT '返修完成时间',
  `repair_method` varchar(128) COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT '维修方式',
  `repair_process` text COLLATE utf8mb4_general_ci COMMENT '维修处理过程',
  `confirm_status` varchar(32) COLLATE utf8mb4_general_ci DEFAULT '待确认' COMMENT '待确认/已完成',
  `remark` text COLLATE utf8mb4_general_ci COMMENT '备注',
  `is_deleted` tinyint DEFAULT '0' COMMENT '是否删除',
  `created_at` datetime DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` datetime DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci COMMENT='维修记录表';
/*!40101 SET character_set_client = @saved_cs_client */;
DROP TABLE IF EXISTS `requirement_books`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8mb4 */;
CREATE TABLE `requirement_books` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `project_id` bigint NOT NULL,
  `book_name` varchar(128) COLLATE utf8mb4_general_ci NOT NULL,
  `file_id` bigint NOT NULL,
  `status` varchar(32) COLLATE utf8mb4_general_ci DEFAULT '草稿',
  `submit_user_id` bigint DEFAULT NULL,
  `submit_user_name` varchar(64) COLLATE utf8mb4_general_ci DEFAULT NULL,
  `submit_time` datetime DEFAULT NULL,
  `audit_user_id` bigint DEFAULT NULL,
  `audit_user_name` varchar(64) COLLATE utf8mb4_general_ci DEFAULT NULL,
  `audit_time` datetime DEFAULT NULL,
  `reject_reason` text COLLATE utf8mb4_general_ci,
  `remark` text COLLATE utf8mb4_general_ci,
  `created_at` datetime DEFAULT CURRENT_TIMESTAMP,
  `updated_at` datetime DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  `is_deleted` tinyint DEFAULT '0',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;
/*!40101 SET character_set_client = @saved_cs_client */;
DROP TABLE IF EXISTS `requirement_changes`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8mb4 */;
CREATE TABLE `requirement_changes` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `project_id` bigint NOT NULL,
  `change_title` varchar(128) COLLATE utf8mb4_general_ci NOT NULL,
  `change_type` varchar(64) COLLATE utf8mb4_general_ci DEFAULT NULL,
  `file_id` bigint NOT NULL,
  `status` varchar(32) COLLATE utf8mb4_general_ci DEFAULT '草稿',
  `close_status` varchar(32) COLLATE utf8mb4_general_ci DEFAULT '未关闭',
  `submit_user_id` bigint DEFAULT NULL,
  `submit_user_name` varchar(64) COLLATE utf8mb4_general_ci DEFAULT NULL,
  `audit_user_id` bigint DEFAULT NULL,
  `audit_user_name` varchar(64) COLLATE utf8mb4_general_ci DEFAULT NULL,
  `submit_time` datetime DEFAULT NULL,
  `audit_time` datetime DEFAULT NULL,
  `close_user_id` bigint DEFAULT NULL,
  `close_user_name` varchar(64) COLLATE utf8mb4_general_ci DEFAULT NULL,
  `close_time` datetime DEFAULT NULL,
  `reject_reason` text COLLATE utf8mb4_general_ci,
  `remark` text COLLATE utf8mb4_general_ci,
  `created_at` datetime DEFAULT CURRENT_TIMESTAMP,
  `updated_at` datetime DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  `is_deleted` tinyint DEFAULT '0',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;
/*!40101 SET character_set_client = @saved_cs_client */;
DROP TABLE IF EXISTS `role_permissions`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8mb4 */;
CREATE TABLE `role_permissions` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `role_id` bigint NOT NULL,
  `permission_id` bigint NOT NULL,
  `created_at` datetime DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_role_permission` (`role_id`,`permission_id`),
  KEY `idx_role_permissions_role` (`role_id`),
  KEY `idx_role_permissions_permission` (`permission_id`)
) ENGINE=InnoDB AUTO_INCREMENT=146 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;
/*!40101 SET character_set_client = @saved_cs_client */;
DROP TABLE IF EXISTS `roles`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8mb4 */;
CREATE TABLE `roles` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `role_code` varchar(64) COLLATE utf8mb4_general_ci NOT NULL,
  `role_name` varchar(64) COLLATE utf8mb4_general_ci NOT NULL,
  `description` varchar(255) COLLATE utf8mb4_general_ci DEFAULT '',
  `created_at` datetime DEFAULT CURRENT_TIMESTAMP,
  `updated_at` datetime DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  UNIQUE KEY `role_code` (`role_code`)
) ENGINE=InnoDB AUTO_INCREMENT=10 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;
/*!40101 SET character_set_client = @saved_cs_client */;
DROP TABLE IF EXISTS `shipping_batch_devices`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8mb4 */;
CREATE TABLE `shipping_batch_devices` (
  `id` bigint NOT NULL AUTO_INCREMENT COMMENT '主键',
  `batch_id` bigint NOT NULL COMMENT '发货批次ID',
  `inventory_device_id` bigint NOT NULL COMMENT '库存设备ID',
  `sn` varchar(128) COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT 'SN',
  `mac_address` varchar(64) COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT 'MAC',
  `device_type` varchar(64) COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT '终端类型',
  `hardware_version` varchar(64) COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT '硬件版本',
  `software_version` varchar(64) COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT '软件版本',
  `is_deleted` tinyint DEFAULT '0' COMMENT '是否删除',
  `created_at` datetime DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci COMMENT='发货批次设备明细表';
/*!40101 SET character_set_client = @saved_cs_client */;
DROP TABLE IF EXISTS `shipping_batches`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8mb4 */;
CREATE TABLE `shipping_batches` (
  `id` bigint NOT NULL AUTO_INCREMENT COMMENT '发货批次ID',
  `batch_no` varchar(64) COLLATE utf8mb4_general_ci NOT NULL COMMENT '发货批次号',
  `project_id` bigint DEFAULT NULL COMMENT '关联项目',
  `express_no` varchar(128) COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT '物流单号',
  `device_count` int DEFAULT '0' COMMENT '设备数量',
  `file_id` bigint DEFAULT NULL COMMENT '发货附件或SN清单文件',
  `uploader_id` bigint DEFAULT NULL COMMENT '上传人ID',
  `uploader_name` varchar(64) COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT '上传人姓名',
  `upload_time` datetime DEFAULT NULL COMMENT '上传时间',
  `audit_status` varchar(32) COLLATE utf8mb4_general_ci DEFAULT '草稿' COMMENT '草稿/待审核/已通过/已驳回',
  `auditor_id` bigint DEFAULT NULL COMMENT '审核人ID',
  `auditor_name` varchar(64) COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT '审核人姓名',
  `audit_time` datetime DEFAULT NULL COMMENT '审核时间',
  `remark` text COLLATE utf8mb4_general_ci COMMENT '备注',
  `reject_reason` text COLLATE utf8mb4_general_ci COMMENT '驳回原因',
  `shipping_desc` text COLLATE utf8mb4_general_ci COMMENT '发货说明',
  `is_deleted` tinyint DEFAULT '0' COMMENT '是否删除',
  `created_at` datetime DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` datetime DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci COMMENT='发货批次主表';
/*!40101 SET character_set_client = @saved_cs_client */;
DROP TABLE IF EXISTS `software_versions`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8mb4 */;
CREATE TABLE `software_versions` (
  `id` bigint NOT NULL AUTO_INCREMENT COMMENT '软件版本ID',
  `project_id` bigint DEFAULT NULL COMMENT '所属项目ID',
  `software_version` varchar(64) COLLATE utf8mb4_general_ci NOT NULL COMMENT '软件版本号',
  `device_type` varchar(64) COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT '终端类型',
  `hardware_id` bigint DEFAULT NULL COMMENT '适配硬件版本ID',
  `hardware_version` varchar(64) COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT '硬件版本号冗余',
  `owner_id` bigint DEFAULT NULL COMMENT '软件负责人ID',
  `owner_name` varchar(64) COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT '软件负责人',
  `release_date` date DEFAULT NULL COMMENT '发布日期',
  `download_url` varchar(512) COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT '软件下载地址',
  `business_desc` text COLLATE utf8mb4_general_ci COMMENT '业务说明',
  `description` text COLLATE utf8mb4_general_ci COMMENT '版本说明',
  `software_status` varchar(32) COLLATE utf8mb4_general_ci DEFAULT '草稿' COMMENT '草稿/已发布/已废弃',
  `created_at` datetime DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` datetime DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  `is_deleted` tinyint DEFAULT '0' COMMENT '是否删除，0否，1是',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci COMMENT='软件版本表';
/*!40101 SET character_set_client = @saved_cs_client */;
DROP TABLE IF EXISTS `test_cases`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8mb4 */;
CREATE TABLE `test_cases` (
  `id` bigint NOT NULL AUTO_INCREMENT COMMENT '测试用例ID',
  `project_id` bigint DEFAULT '0' COMMENT '项目ID',
  `case_name` varchar(128) COLLATE utf8mb4_general_ci DEFAULT '' COMMENT '测试用例名称',
  `file_id` bigint DEFAULT '0' COMMENT '测试用例文件ID',
  `file_name` varchar(255) COLLATE utf8mb4_general_ci DEFAULT '' COMMENT '测试用例文件名',
  `uploader_id` bigint DEFAULT '0' COMMENT '上传人ID',
  `uploader_name` varchar(64) COLLATE utf8mb4_general_ci DEFAULT '' COMMENT '上传人姓名',
  `upload_time` datetime DEFAULT NULL COMMENT '上传时间',
  `audit_status` varchar(32) COLLATE utf8mb4_general_ci DEFAULT '草稿' COMMENT '审核状态',
  `auditor_id` bigint DEFAULT '0' COMMENT '审核人ID',
  `auditor_name` varchar(64) COLLATE utf8mb4_general_ci DEFAULT '' COMMENT '审核人姓名',
  `audit_time` datetime DEFAULT NULL COMMENT '审核时间',
  `reject_reason` text COLLATE utf8mb4_general_ci COMMENT '驳回原因',
  `remark` text COLLATE utf8mb4_general_ci COMMENT '备注',
  `report_name` varchar(128) COLLATE utf8mb4_general_ci DEFAULT '' COMMENT '测试报告名称',
  `report_file_id` bigint DEFAULT '0' COMMENT '测试报告文件ID',
  `report_file_name` varchar(255) COLLATE utf8mb4_general_ci DEFAULT '' COMMENT '测试报告文件名',
  `report_uploader_id` bigint DEFAULT '0' COMMENT '测试报告上传人ID',
  `report_uploader_name` varchar(64) COLLATE utf8mb4_general_ci DEFAULT '' COMMENT '测试报告上传人姓名',
  `report_upload_time` datetime DEFAULT NULL COMMENT '测试报告上传时间',
  `report_remark` text COLLATE utf8mb4_general_ci COMMENT '测试报告说明',
  `is_deleted` tinyint DEFAULT '0' COMMENT '是否删除',
  `created_at` datetime DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` datetime DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci COMMENT='测试用例表';
/*!40101 SET character_set_client = @saved_cs_client */;
DROP TABLE IF EXISTS `uploaded_files`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8mb4 */;
CREATE TABLE `uploaded_files` (
  `id` bigint NOT NULL,
  `file_name` varchar(255) NOT NULL DEFAULT '',
  `content_type` varchar(128) NOT NULL DEFAULT '',
  `file_data` longblob,
  `created_at` datetime DEFAULT CURRENT_TIMESTAMP,
  `updated_at` datetime DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;
DROP TABLE IF EXISTS `user_roles`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8mb4 */;
CREATE TABLE `user_roles` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `user_id` bigint NOT NULL,
  `role_id` bigint NOT NULL,
  `created_at` datetime DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_user_role` (`user_id`,`role_id`),
  KEY `idx_user_roles_user` (`user_id`),
  KEY `idx_user_roles_role` (`role_id`)
) ENGINE=InnoDB AUTO_INCREMENT=16 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;
/*!40101 SET character_set_client = @saved_cs_client */;
DROP TABLE IF EXISTS `users`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8mb4 */;
CREATE TABLE `users` (
  `id` bigint NOT NULL AUTO_INCREMENT COMMENT '用户ID',
  `username` varchar(64) COLLATE utf8mb4_general_ci NOT NULL COMMENT '登录账号',
  `password_hash` varchar(255) COLLATE utf8mb4_general_ci NOT NULL COMMENT '密码',
  `real_name` varchar(64) COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT '真实姓名',
  `email` varchar(128) COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT '邮箱',
  `phone` varchar(32) COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT '手机号',
  `department` varchar(64) COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT '部门',
  `status` varchar(32) COLLATE utf8mb4_general_ci DEFAULT '启用' COMMENT '启用/禁用',
  `last_login_time` datetime DEFAULT NULL COMMENT '最后登录时间',
  `created_at` datetime DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` datetime DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  `avatar_url` blob,
  PRIMARY KEY (`id`),
  UNIQUE KEY `username` (`username`)
) ENGINE=InnoDB AUTO_INCREMENT=14 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci COMMENT='用户表';
/*!40101 SET character_set_client = @saved_cs_client */;
/*!40103 SET TIME_ZONE=@OLD_TIME_ZONE */;

/*!40101 SET SQL_MODE=@OLD_SQL_MODE */;
/*!40014 SET FOREIGN_KEY_CHECKS=@OLD_FOREIGN_KEY_CHECKS */;
/*!40014 SET UNIQUE_CHECKS=@OLD_UNIQUE_CHECKS */;
/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40101 SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */;
/*!40101 SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION */;
/*!40111 SET SQL_NOTES=@OLD_SQL_NOTES */;
