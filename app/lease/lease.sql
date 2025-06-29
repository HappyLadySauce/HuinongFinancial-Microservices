-- Lease服务数据库初始化脚本
-- 数据库: lease
-- 连接: lease:lease@tcp(10.10.10.6:3306)/lease?charset=utf8mb4&parseTime=True&loc=Local
-- 职责：管理租赁申请、审批、合同等业务流程

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
  `application_id` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '申请编号(业务编号)',
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
  `description` text CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci COMMENT '申请描述',
  `special_requirements` text CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci COMMENT '特殊要求',
  `status` varchar(20) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT 'pending' COMMENT '状态 pending/approved/rejected/cancelled/completed',
  `risk_level` varchar(20) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT 'unknown' COMMENT '风险等级 low/medium/high/unknown',
  `ai_risk_score` int UNSIGNED DEFAULT NULL COMMENT 'AI风险评分(0-100)',
  `ai_suggestion` text CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci COMMENT 'AI评估建议',
  `submission_time` timestamp NULL DEFAULT CURRENT_TIMESTAMP COMMENT '提交时间',
  `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_application_id` (`application_id`),
  KEY `idx_user_id` (`user_id`),
  KEY `idx_product_id` (`product_id`),
  KEY `idx_product_code` (`product_code`),
  KEY `idx_status` (`status`),
  KEY `idx_risk_level` (`risk_level`),
  KEY `idx_submission_time` (`submission_time`),
  KEY `idx_start_date` (`start_date`),
  KEY `idx_created_at` (`created_at`)
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
  `action` varchar(20) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '审批动作 approve/reject/return',
  `suggestions` text CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci COMMENT '审批意见',
  `approved_duration` int UNSIGNED DEFAULT NULL COMMENT '批准租期(天)',
  `approved_amount` decimal(10,2) DEFAULT NULL COMMENT '批准金额',
  `approved_deposit` decimal(10,2) DEFAULT NULL COMMENT '批准押金',
  `conditions` text CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci COMMENT '附加条件',
  `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  PRIMARY KEY (`id`),
  KEY `idx_application_id` (`application_id`),
  KEY `idx_auditor_id` (`auditor_id`),
  KEY `idx_action` (`action`),
  KEY `idx_created_at` (`created_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='租赁审批记录表';

-- ----------------------------
-- 租赁合同表
-- ----------------------------
DROP TABLE IF EXISTS `lease_contracts`;
CREATE TABLE `lease_contracts` (
  `id` bigint UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '合同ID',
  `application_id` bigint UNSIGNED NOT NULL COMMENT '申请ID',
  `contract_no` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '合同编号',
  `user_id` bigint UNSIGNED NOT NULL COMMENT '用户ID',
  `product_id` bigint UNSIGNED NOT NULL COMMENT '产品ID',
  `product_code` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '产品编码',
  `machinery` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '设备名称',
  `start_date` date NOT NULL COMMENT '租赁开始日期',
  `end_date` date NOT NULL COMMENT '租赁结束日期',
  `duration` int UNSIGNED NOT NULL COMMENT '租期(天)',
  `daily_rate` decimal(10,2) NOT NULL COMMENT '日租金',
  `total_amount` decimal(10,2) NOT NULL COMMENT '租金总额',
  `deposit` decimal(10,2) NOT NULL COMMENT '押金',
  `delivery_address` varchar(500) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '交付地址',
  `delivery_date` date DEFAULT NULL COMMENT '实际交付日期',
  `return_date` date DEFAULT NULL COMMENT '实际归还日期',
  `status` varchar(20) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT 'active' COMMENT '状态 active/completed/terminated/overdue',
  `payment_status` varchar(20) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT 'unpaid' COMMENT '付款状态 unpaid/partial/paid/overdue',
  `equipment_status` varchar(20) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT 'normal' COMMENT '设备状态 normal/damaged/lost',
  `remarks` text CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci COMMENT '备注',
  `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_contract_no` (`contract_no`),
  KEY `idx_application_id` (`application_id`),
  KEY `idx_user_id` (`user_id`),
  KEY `idx_product_id` (`product_id`),
  KEY `idx_product_code` (`product_code`),
  KEY `idx_status` (`status`),
  KEY `idx_payment_status` (`payment_status`),
  KEY `idx_start_date` (`start_date`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='租赁合同表';

-- ----------------------------
-- 租金支付计划表
-- ----------------------------
DROP TABLE IF EXISTS `lease_payment_schedules`;
CREATE TABLE `lease_payment_schedules` (
  `id` bigint UNSIGNED NOT NULL AUTO_INCREMENT COMMENT 'ID',
  `contract_id` bigint UNSIGNED NOT NULL COMMENT '合同ID',
  `period` int UNSIGNED NOT NULL COMMENT '期数',
  `payment_type` varchar(20) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '付款类型 rent/deposit/penalty',
  `due_date` date NOT NULL COMMENT '应付日期',
  `amount` decimal(10,2) NOT NULL COMMENT '应付金额',
  `paid_amount` decimal(10,2) DEFAULT 0.00 COMMENT '已付金额',
  `status` varchar(20) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT 'pending' COMMENT '状态 pending/paid/overdue',
  `paid_at` timestamp NULL DEFAULT NULL COMMENT '实际付款时间',
  `payment_method` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT '' COMMENT '付款方式',
  `payment_reference` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT '' COMMENT '付款凭证号',
  `remarks` text CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci COMMENT '备注',
  `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  KEY `idx_contract_id` (`contract_id`),
  KEY `idx_due_date` (`due_date`),
  KEY `idx_status` (`status`),
  KEY `idx_payment_type` (`payment_type`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='租金支付计划表';

-- ----------------------------
-- 设备交付记录表
-- ----------------------------
DROP TABLE IF EXISTS `lease_delivery_records`;
CREATE TABLE `lease_delivery_records` (
  `id` bigint UNSIGNED NOT NULL AUTO_INCREMENT COMMENT 'ID',
  `contract_id` bigint UNSIGNED NOT NULL COMMENT '合同ID',
  `delivery_type` varchar(20) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '交付类型 delivery/return',
  `scheduled_date` date NOT NULL COMMENT '计划日期',
  `actual_date` date DEFAULT NULL COMMENT '实际日期',
  `address` varchar(500) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '交付地址',
  `contact_person` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '联系人',
  `contact_phone` varchar(20) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '联系电话',
  `equipment_condition` text CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci COMMENT '设备状况描述',
  `photos` json DEFAULT NULL COMMENT '照片记录',
  `status` varchar(20) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT 'scheduled' COMMENT '状态 scheduled/in_transit/completed/failed',
  `driver_name` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT '' COMMENT '司机姓名',
  `driver_phone` varchar(20) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT '' COMMENT '司机电话',
  `vehicle_info` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT '' COMMENT '运输车辆信息',
  `remarks` text CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci COMMENT '备注',
  `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  KEY `idx_contract_id` (`contract_id`),
  KEY `idx_delivery_type` (`delivery_type`),
  KEY `idx_scheduled_date` (`scheduled_date`),
  KEY `idx_status` (`status`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='设备交付记录表';

-- ----------------------------
-- 租赁申请材料表
-- ----------------------------
DROP TABLE IF EXISTS `lease_documents`;
CREATE TABLE `lease_documents` (
  `id` bigint UNSIGNED NOT NULL AUTO_INCREMENT COMMENT 'ID',
  `application_id` bigint UNSIGNED NOT NULL COMMENT '申请ID',
  `type` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '文件类型 identity/business/certificate/other',
  `name` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '文件名',
  `url` varchar(500) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '文件URL',
  `size` bigint DEFAULT 0 COMMENT '文件大小',
  `status` varchar(20) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT 'pending' COMMENT '审核状态 pending/approved/rejected',
  `remark` text CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci COMMENT '备注',
  `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  PRIMARY KEY (`id`),
  KEY `idx_application_id` (`application_id`),
  KEY `idx_type` (`type`),
  KEY `idx_status` (`status`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='租赁申请材料表';

-- ----------------------------
-- 初始化数据
-- ----------------------------

-- 租赁申请测试数据
INSERT INTO `lease_applications` (`application_id`, `user_id`, `applicant_name`, `product_id`, `product_code`, `name`, `type`, `machinery`, `start_date`, `end_date`, `duration`, `daily_rate`, `total_amount`, `deposit`, `delivery_address`, `contact_phone`, `purpose`, `description`, `status`, `risk_level`, `ai_risk_score`, `ai_suggestion`, `submission_time`) VALUES
('LS202506010001', 1, '张三', 1, 'LEASE001', '春耕拖拉机租赁', '农机租赁', '约翰迪尔1204拖拉机', '2025-03-01', '2025-04-30', 60, 500.00, 30000.00, 5000.00, '江苏省南京市浦口区永宁街道农业园区', '13800138000', '春季农田耕作', '春季播种期间需要租赁大型拖拉机进行土地耕作，作业面积约200亩', 'pending', 'medium', 75, 'AI建议通过：申请人有农机操作经验，租赁目的明确，建议批准', '2024-06-01 09:30:00'),

('LS202506010002', 2, '李四', 2, 'LEASE002', '水稻收割机租赁', '农机租赁', '久保田4LZ-4J收割机', '2025-09-15', '2025-10-15', 30, 800.00, 24000.00, 10000.00, '江苏省扬州市高邮市汤庄镇', '13800138001', '水稻收获', '水稻收获季节租赁联合收割机，预计作业面积150亩', 'pending', 'low', 85, 'AI建议通过：收获季节租赁需求合理，申请人信用良好', '2024-06-01 14:20:00'),

('LS202506020001', 3, '王五', 3, 'LEASE003', '播种机租赁', '农机租赁', '大华宝来2BYF-6播种机', '2025-04-01', '2025-04-15', 14, 300.00, 4200.00, 2000.00, '江苏省盐城市东台市安丰镇', '13800138002', '大豆播种', '大豆播种期间设备租赁，配合拖拉机作业', 'approved', 'low', 92, 'AI建议通过：短期租赁，风险可控，设备需求明确', '2024-06-02 10:15:00'),

('LS202506020002', 4, '赵六', 6, 'LEASE006', '粮食烘干设备租赁', '加工设备', '三久烘干机', '2025-10-01', '2025-11-30', 60, 600.00, 36000.00, 8000.00, '江苏省泰州市兴化市戴南镇', '13800138003', '玉米烘干加工', '玉米收获后烘干处理，预计处理量800吨', 'approved', 'medium', 78, 'AI建议通过：烘干需求合理，申请人有加工经验', '2024-06-02 16:45:00'),

('LS202506030001', 5, '孙七', 5, 'LEASE005', '灌溉水泵租赁', '灌溉设备', '南方水泵', '2025-06-01', '2025-08-31', 91, 150.00, 13650.00, 500.00, '江苏省连云港市灌云县伊山镇', '13800138004', '夏季农田灌溉', '夏季干旱期间农田灌溉，覆盖面积300亩', 'pending', 'low', 88, 'AI建议通过：灌溉需求合理，设备适用性强', '2024-06-03 08:30:00');

-- 租赁审批记录测试数据
INSERT INTO `lease_approvals` (`application_id`, `auditor_id`, `auditor_name`, `action`, `suggestions`, `approved_duration`, `approved_amount`, `approved_deposit`, `conditions`) VALUES
(3, 3, '许先生', 'approve', '设备状态良好，申请人有相关经验，同意租赁', 14, 4200.00, 2000.00, '需要提供拖拉机配套使用证明'),
(4, 2, '柴大三', 'approve', '烘干需求合理，申请人有加工场地，符合租赁条件', 60, 36000.00, 8000.00, '需要提供场地安全检查报告，配备专业操作人员');

-- 租赁合同测试数据 (仅批准的申请)
INSERT INTO `lease_contracts` (`application_id`, `contract_no`, `user_id`, `product_id`, `product_code`, `machinery`, `start_date`, `end_date`, `duration`, `daily_rate`, `total_amount`, `deposit`, `delivery_address`, `delivery_date`, `status`) VALUES
(3, 'HL202506050001', 3, 3, 'LEASE003', '大华宝来2BYF-6播种机', '2025-04-01', '2025-04-15', 14, 300.00, 4200.00, 2000.00, '江苏省盐城市东台市安丰镇', '2025-03-31', 'completed'),
(4, 'HL202506050002', 4, 6, 'LEASE006', '三久烘干机', '2025-10-01', '2025-11-30', 60, 600.00, 36000.00, 8000.00, '江苏省泰州市兴化市戴南镇', '2025-09-30', 'active');

-- 租金支付计划测试数据
INSERT INTO `lease_payment_schedules` (`contract_id`, `period`, `payment_type`, `due_date`, `amount`, `paid_amount`, `status`, `paid_at`, `payment_method`) VALUES
-- 播种机合同支付记录（已完成）
(1, 1, 'deposit', '2025-03-30', 2000.00, 2000.00, 'paid', '2025-03-30 10:00:00', '银行转账'),
(1, 2, 'rent', '2025-04-01', 4200.00, 4200.00, 'paid', '2025-04-01 14:30:00', '银行转账'),

-- 烘干机合同支付记录
(2, 1, 'deposit', '2025-09-25', 8000.00, 8000.00, 'paid', '2025-09-25 09:15:00', '银行转账'),
(2, 2, 'rent', '2025-10-01', 18000.00, 18000.00, 'paid', '2025-10-01 11:20:00', '银行转账'),
(2, 3, 'rent', '2025-11-01', 18000.00, 0.00, 'pending', NULL, '');

-- 设备交付记录测试数据
INSERT INTO `lease_delivery_records` (`contract_id`, `delivery_type`, `scheduled_date`, `actual_date`, `address`, `contact_person`, `contact_phone`, `equipment_condition`, `status`, `driver_name`, `driver_phone`, `vehicle_info`) VALUES
-- 播种机交付记录
(1, 'delivery', '2025-03-31', '2025-03-31', '江苏省盐城市东台市安丰镇', '王五', '13800138002', '设备状况良好，已完成保养检查', 'completed', '张师傅', '13900139000', '苏B12345 平板运输车'),
(1, 'return', '2025-04-15', '2025-04-15', '江苏省盐城市东台市安丰镇', '王五', '13800138002', '设备正常使用，无明显损坏', 'completed', '张师傅', '13900139000', '苏B12345 平板运输车'),

-- 烘干机交付记录
(2, 'delivery', '2025-09-30', '2025-09-30', '江苏省泰州市兴化市戴南镇', '赵六', '13800138003', '设备状况良好，已完成调试', 'completed', '李师傅', '13900139001', '苏C67890 专用运输车');

SET FOREIGN_KEY_CHECKS = 1; 