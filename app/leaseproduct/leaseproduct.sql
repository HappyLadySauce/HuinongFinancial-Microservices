-- LeaseProduct服务数据库初始化脚本
-- 数据库: leaseproduct
-- 连接: leaseproduct:leaseproduct@tcp(10.10.10.7:3306)/leaseproduct?charset=utf8mb4&parseTime=True&loc=Local
-- 职责：专门管理租赁产品信息

SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;

-- 创建数据库
CREATE DATABASE IF NOT EXISTS `leaseproduct` DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
USE `leaseproduct`;

-- ----------------------------
-- 租赁产品表
-- ----------------------------
DROP TABLE IF EXISTS `lease_products`;
CREATE TABLE `lease_products` (
  `id` bigint UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '产品ID',
  `product_code` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '产品编码',
  `name` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '产品名称',
  `type` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '租赁类型',
  `machinery` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT '' COMMENT '设备名称',
  `brand` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT '' COMMENT '品牌',
  `model` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT '' COMMENT '型号',
  `daily_rate` decimal(10,2) NOT NULL COMMENT '日租金',
  `deposit` decimal(10,2) DEFAULT 0.00 COMMENT '押金',
  `max_duration` int UNSIGNED DEFAULT 365 COMMENT '最大租期(天)',
  `min_duration` int UNSIGNED DEFAULT 1 COMMENT '最小租期(天)',
  `description` text CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '产品描述',
  `inventory_count` int UNSIGNED DEFAULT 0 COMMENT '库存数量',
  `available_count` int UNSIGNED DEFAULT 0 COMMENT '可用数量',
  `status` tinyint UNSIGNED DEFAULT 1 COMMENT '状态 1:可租 2:维护中 3:已租完',
  `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_product_code` (`product_code`),
  KEY `idx_type` (`type`),
  KEY `idx_brand` (`brand`),
  KEY `idx_status` (`status`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='租赁产品表';

-- ----------------------------
-- 初始化数据
-- ----------------------------
INSERT INTO `lease_products` (`product_code`, `name`, `type`, `machinery`, `brand`, `model`, `daily_rate`, `deposit`, `max_duration`, `min_duration`, `description`, `inventory_count`, `available_count`, `status`) VALUES
('LEASE001', '大型拖拉机租赁', '农机租赁', '约翰迪尔1204拖拉机', '约翰迪尔', '1204', 500.00, 5000.00, 90, 1, '大型农业拖拉机租赁，适用于耕地、播种等大面积作业', 5, 3, 1),
('LEASE002', '联合收割机租赁', '农机租赁', '久保田4LZ-4J收割机', '久保田', '4LZ-4J', 800.00, 10000.00, 60, 1, '联合收割机租赁，提高粮食收获效率', 3, 2, 1),
('LEASE003', '精密播种机租赁', '农机租赁', '大华宝来2BYF-6播种机', '大华宝来', '2BYF-6', 300.00, 2000.00, 30, 1, '精密播种机租赁，提高播种质量和效率', 4, 4, 1),
('LEASE004', '农产品运输车租赁', '车辆租赁', '福田货车', '福田', 'BJ1043V9JDA-F3', 400.00, 3000.00, 180, 1, '农产品专用运输车辆租赁服务', 8, 6, 1);

SET FOREIGN_KEY_CHECKS = 1; 