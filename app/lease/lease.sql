-- Lease服务数据库初始化脚本
-- 数据库: lease
-- 连接: lease:lease@tcp(10.10.10.6:3306)/lease?charset=utf8mb4&parseTime=True&loc=Local
-- 职责：管理租赁申请、审批等业务流程

SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;

-- 创建数据库
CREATE DATABASE IF NOT EXISTS `lease` DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
USE `lease`;

-- ----------------------------
-- 租赁申请表
-- ----------------------------
DROP TABLE IF EXISTS `lease_applications`;
CREATE TABLE `lease_applications` (
  `id` bigint UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '申请ID',
  `application_id` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '申请编号',
  `user_id` bigint UNSIGNED NOT NULL COMMENT '用户ID',
  `applicant_name` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '申请人姓名',
  `product_id` bigint UNSIGNED NOT NULL COMMENT '租赁产品ID',
  `product_code` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '产品编码',
  `name` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '申请名称',
  `type` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '租赁类型',
  `machinery` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT '' COMMENT '设备名称',
  `start_date` date NOT NULL COMMENT '开始日期',
  `end_date` date NOT NULL COMMENT '结束日期',
  `duration` int UNSIGNED NOT NULL COMMENT '租期(天)',
  `daily_rate` decimal(10,2) NOT NULL COMMENT '日租金',
  `total_amount` decimal(10,2) NOT NULL COMMENT '总金额',
  `deposit` decimal(10,2) DEFAULT 0.00 COMMENT '押金',
  `delivery_address` varchar(500) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT '' COMMENT '交付地址',
  `contact_phone` varchar(20) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT '' COMMENT '联系电话',
  `purpose` text CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci COMMENT '使用目的',
  `status` varchar(20) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT 'pending' COMMENT '状态 pending/approved/rejected/cancelled',
  `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_application_id` (`application_id`),
  KEY `idx_user_id` (`user_id`),
  KEY `idx_product_id` (`product_id`),
  KEY `idx_status` (`status`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='租赁申请表';

-- ----------------------------
-- 租赁审批记录表
-- ----------------------------
DROP TABLE IF EXISTS `lease_approvals`;
CREATE TABLE `lease_approvals` (
  `id` bigint UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '审批ID',
  `application_id` bigint UNSIGNED NOT NULL COMMENT '申请ID',
  `auditor_id` bigint UNSIGNED NOT NULL COMMENT '审核员ID',
  `auditor_name` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '审核员姓名',
  `action` varchar(20) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '审批动作 approve/reject',
  `suggestions` text CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci COMMENT '审批意见',
  `approved_duration` int UNSIGNED DEFAULT NULL COMMENT '批准租期(天)',
  `approved_amount` decimal(10,2) DEFAULT NULL COMMENT '批准金额',
  `approved_deposit` decimal(10,2) DEFAULT NULL COMMENT '批准押金',
  `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  PRIMARY KEY (`id`),
  KEY `idx_application_id` (`application_id`),
  KEY `idx_auditor_id` (`auditor_id`),
  KEY `idx_action` (`action`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='租赁审批记录表';

-- ----------------------------
-- 初始化数据
-- ----------------------------
INSERT INTO `lease_applications` (`application_id`, `user_id`, `applicant_name`, `product_id`, `product_code`, `name`, `type`, `machinery`, `start_date`, `end_date`, `duration`, `daily_rate`, `total_amount`, `deposit`, `delivery_address`, `contact_phone`, `purpose`, `status`) VALUES
('LS202506010001', 1, '张三', 1, 'LEASE001', '春耕拖拉机租赁', '农机租赁', '约翰迪尔1204拖拉机', '2025-03-01', '2025-04-30', 60, 500.00, 30000.00, 5000.00, '江苏省南京市浦口区永宁街道农业园区', '13800138000', '春季农田耕作', 'approved'),
('LS202506010002', 2, '李四', 2, 'LEASE002', '水稻收割机租赁', '农机租赁', '久保田4LZ-4J收割机', '2025-09-15', '2025-10-15', 30, 800.00, 24000.00, 10000.00, '江苏省扬州市高邮市汤庄镇', '13800138001', '水稻收获', 'pending'),
('LS202506020001', 3, '王五', 3, 'LEASE003', '播种机租赁', '农机租赁', '大华宝来2BYF-6播种机', '2025-04-01', '2025-04-15', 14, 300.00, 4200.00, 2000.00, '江苏省盐城市东台市安丰镇', '13800138002', '大豆播种', 'approved');

INSERT INTO `lease_approvals` (`application_id`, `auditor_id`, `auditor_name`, `action`, `suggestions`, `approved_duration`, `approved_amount`, `approved_deposit`) VALUES
(1, 2001, 'B端测试管理员', 'approve', '设备需求合理，客户信用良好，同意租赁', 60, 30000.00, 5000.00),
(3, 2002, '普通操作员1', 'approve', '短期租赁，风险可控，设备适用，建议批准', 14, 4200.00, 2000.00);

SET FOREIGN_KEY_CHECKS = 1; 