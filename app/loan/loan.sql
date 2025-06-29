-- Loan服务数据库初始化脚本
-- 数据库: loan
-- 连接: loan:loan@tcp(10.10.10.6:3306)/loan?charset=utf8mb4&parseTime=True&loc=Local

SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;

-- 创建数据库
CREATE DATABASE IF NOT EXISTS `loan` DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
USE `loan`;

-- ----------------------------
-- 贷款申请表
-- ----------------------------
DROP TABLE IF EXISTS `loan_applications`;
CREATE TABLE `loan_applications` (
  `id` bigint UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '申请ID',
  `application_id` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '申请编号(业务编号)',
  `user_id` bigint UNSIGNED NOT NULL COMMENT '用户ID',
  `applicant_name` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '申请人姓名',
  `product_id` bigint UNSIGNED NOT NULL COMMENT '贷款产品ID',
  `name` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '申请名称',
  `type` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '贷款类型',
  `amount` decimal(15,2) NOT NULL COMMENT '申请金额',
  `duration` int UNSIGNED NOT NULL COMMENT '贷款期限(月)',
  `purpose` text CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci COMMENT '贷款用途',
  `description` text CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci COMMENT '申请描述',
  `status` varchar(20) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT 'pending' COMMENT '状态 pending/reviewing/approved/rejected/cancelled',
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
  KEY `idx_status` (`status`),
  KEY `idx_risk_level` (`risk_level`),
  KEY `idx_submission_time` (`submission_time`),
  KEY `idx_created_at` (`created_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='贷款申请表';

-- ----------------------------
-- 贷款审批记录表
-- ----------------------------
DROP TABLE IF EXISTS `loan_approvals`;
CREATE TABLE `loan_approvals` (
  `id` bigint UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '审批ID',
  `application_id` bigint UNSIGNED NOT NULL COMMENT '申请ID',
  `auditor_id` bigint UNSIGNED NOT NULL COMMENT '审核员ID',
  `auditor_name` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '审核员姓名',
  `action` varchar(20) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '审批动作 approve/reject/return',
  `suggestions` text CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci COMMENT '审批意见',
  `approved_amount` decimal(15,2) DEFAULT NULL COMMENT '批准金额',
  `approved_duration` int UNSIGNED DEFAULT NULL COMMENT '批准期限(月)',
  `interest_rate` decimal(5,2) DEFAULT NULL COMMENT '利率(%)',
  `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  PRIMARY KEY (`id`),
  KEY `idx_application_id` (`application_id`),
  KEY `idx_auditor_id` (`auditor_id`),
  KEY `idx_action` (`action`),
  KEY `idx_created_at` (`created_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='贷款审批记录表';

-- ----------------------------
-- 贷款合同表
-- ----------------------------
DROP TABLE IF EXISTS `loan_contracts`;
CREATE TABLE `loan_contracts` (
  `id` bigint UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '合同ID',
  `application_id` bigint UNSIGNED NOT NULL COMMENT '申请ID',
  `contract_no` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '合同编号',
  `user_id` bigint UNSIGNED NOT NULL COMMENT '用户ID',
  `amount` decimal(15,2) NOT NULL COMMENT '贷款金额',
  `duration` int UNSIGNED NOT NULL COMMENT '贷款期限(月)',
  `interest_rate` decimal(5,2) NOT NULL COMMENT '年利率(%)',
  `monthly_payment` decimal(10,2) NOT NULL COMMENT '月还款额',
  `start_date` date NOT NULL COMMENT '放款日期',
  `end_date` date NOT NULL COMMENT '到期日期',
  `status` varchar(20) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT 'active' COMMENT '状态 active/completed/overdue/cancelled',
  `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_contract_no` (`contract_no`),
  KEY `idx_application_id` (`application_id`),
  KEY `idx_user_id` (`user_id`),
  KEY `idx_status` (`status`),
  KEY `idx_start_date` (`start_date`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='贷款合同表';

-- ----------------------------
-- 还款计划表
-- ----------------------------
DROP TABLE IF EXISTS `repayment_schedules`;
CREATE TABLE `repayment_schedules` (
  `id` bigint UNSIGNED NOT NULL AUTO_INCREMENT COMMENT 'ID',
  `contract_id` bigint UNSIGNED NOT NULL COMMENT '合同ID',
  `period` int UNSIGNED NOT NULL COMMENT '期数',
  `due_date` date NOT NULL COMMENT '还款日期',
  `principal` decimal(10,2) NOT NULL COMMENT '本金',
  `interest` decimal(10,2) NOT NULL COMMENT '利息',
  `total_amount` decimal(10,2) NOT NULL COMMENT '还款总额',
  `status` varchar(20) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT 'pending' COMMENT '状态 pending/paid/overdue',
  `paid_at` timestamp NULL DEFAULT NULL COMMENT '实际还款时间',
  `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  KEY `idx_contract_id` (`contract_id`),
  KEY `idx_due_date` (`due_date`),
  KEY `idx_status` (`status`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='还款计划表';

-- ----------------------------
-- 贷款申请材料表
-- ----------------------------
DROP TABLE IF EXISTS `loan_documents`;
CREATE TABLE `loan_documents` (
  `id` bigint UNSIGNED NOT NULL AUTO_INCREMENT COMMENT 'ID',
  `application_id` bigint UNSIGNED NOT NULL COMMENT '申请ID',
  `type` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '文件类型 identity/income/asset/other',
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
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='贷款申请材料表';

-- ----------------------------
-- 初始化数据
-- ----------------------------

-- 贷款申请测试数据
INSERT INTO `loan_applications` (`id`, `application_id`, `user_id`, `applicant_name`, `product_id`, `name`, `type`, `amount`, `duration`, `purpose`, `description`, `status`, `risk_level`, `ai_risk_score`, `ai_suggestion`, `submission_time`) VALUES
(1, 'LN202506010001', 1, '张三', 1, '春耕资金贷款', '农业贷', 32567.00, 2, '春季农业生产', '申请春耕生产资金，用于购买种子、化肥等农资', 'pending', 'medium', 75, 'AI建议人工审核：申请人收入稳定，但申请金额较大，建议核实农业生产计划', '2024-06-01 09:30:00'),
(2, 'LN202506010002', 2, '沈女士', 1, '种植基地扩建', '农业贷', 318752.00, 3, '扩大种植规模', '计划扩建有机蔬菜种植基地，提高产能', 'pending', 'high', 45, 'AI建议详细评估：申请金额较大，需要详细核实项目可行性和还款能力', '2024-06-01 14:20:00'),
(3, 'LN202506020001', 3, '许先生', 2, '农机购买贷款', '经营贷', 229461.00, 5, '购买农机设备', '购买联合收割机等现代化农机设备，提高作业效率', 'pending', 'low', 85, 'AI建议通过：申请人信用良好，收入稳定，农机设备可作为抵押物', '2024-06-02 10:15:00'),
(4, 'LN202506020002', 4, '柴大三', 3, '养殖场建设', '创业贷', 361760.00, 4, '建设现代化养殖场', '计划建设标准化养殖场，发展生态养殖业', 'rejected', 'high', 35, 'AI建议拒绝：收入证明不充分，项目风险较高，建议完善材料后重新申请', '2024-06-02 16:45:00'),
(5, 'LN202506030001', 5, '高先生', 1, '有机蔬菜种植', '农业贷', 135316.00, 3, '发展有机农业', '发展有机蔬菜种植，打造绿色农产品品牌', 'approved', 'low', 92, 'AI建议通过：申请材料完善，信用评分高，项目前景良好，风险评估通过', '2024-06-03 08:30:00');

-- 贷款审批记录测试数据
INSERT INTO `loan_approvals` (`application_id`, `auditor_id`, `auditor_name`, `action`, `suggestions`, `approved_amount`, `approved_duration`, `interest_rate`) VALUES
(4, 2, '柴大三', 'reject', '收入证明不充分，风险较高', NULL, NULL, NULL),
(5, 3, '王大兰', 'approve', '申请材料完善，风险评估通过', 135316.00, 3, 5.20);

-- 贷款合同测试数据 (仅批准的申请)
INSERT INTO `loan_contracts` (`application_id`, `contract_no`, `user_id`, `amount`, `duration`, `interest_rate`, `monthly_payment`, `start_date`, `end_date`, `status`) VALUES
(5, 'HN202506050001', 5, 135316.00, 3, 5.20, 45789.50, '2025-06-01', '2025-09-01', 'active');

-- 还款计划测试数据
INSERT INTO `repayment_schedules` (`contract_id`, `period`, `due_date`, `principal`, `interest`, `total_amount`, `status`) VALUES
(1, 1, '2025-07-01', 44549.33, 1240.17, 45789.50, 'pending'),
(1, 2, '2025-08-01', 44549.33, 1240.17, 45789.50, 'pending'),
(1, 3, '2025-09-01', 46217.34, 1572.16, 47789.50, 'pending');

SET FOREIGN_KEY_CHECKS = 1;
