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
  `application_id` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '申请编号',
  `user_id` bigint UNSIGNED NOT NULL COMMENT '用户ID',
  `applicant_name` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '申请人姓名',
  `product_id` bigint UNSIGNED NOT NULL COMMENT '贷款产品ID',
  `name` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '申请名称',
  `type` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '贷款类型',
  `amount` decimal(15,2) NOT NULL COMMENT '申请金额',
  `duration` int UNSIGNED NOT NULL COMMENT '贷款期限(月)',
  `purpose` text CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci COMMENT '贷款用途',
  `status` varchar(20) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT 'pending' COMMENT '状态 pending/approved/rejected/cancelled',
  `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_application_id` (`application_id`),
  KEY `idx_user_id` (`user_id`),
  KEY `idx_product_id` (`product_id`),
  KEY `idx_status` (`status`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='贷款申请表';

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
  `status` varchar(20) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT 'active' COMMENT '状态 active/completed/cancelled',
  `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_contract_no` (`contract_no`),
  KEY `idx_application_id` (`application_id`),
  KEY `idx_user_id` (`user_id`),
  KEY `idx_status` (`status`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='贷款合同表';

-- ----------------------------
-- 初始化数据
-- ----------------------------
INSERT INTO `loan_applications` (`application_id`, `user_id`, `applicant_name`, `product_id`, `name`, `type`, `amount`, `duration`, `purpose`, `status`) VALUES
('LN202506010001', 1, '张三', 1, '春耕资金贷款', '农业贷', 50000.00, 12, '春季农业生产资金', 'approved'),
('LN202506010002', 2, '李四', 2, '农机购买贷款', '经营贷', 200000.00, 36, '购买农机设备', 'pending'),
('LN202506020001', 3, '王五', 3, '日常消费贷款', '消费贷', 30000.00, 24, '家庭日常开支', 'approved');

INSERT INTO `loan_contracts` (`application_id`, `contract_no`, `user_id`, `amount`, `duration`, `interest_rate`, `monthly_payment`, `start_date`, `end_date`, `status`) VALUES
(1, 'HN202506010001', 1, 50000.00, 12, 5.20, 4350.00, '2025-01-01', '2025-12-31', 'active'),
(3, 'HN202506020001', 3, 30000.00, 24, 7.80, 1380.00, '2025-02-01', '2026-01-31', 'active');

SET FOREIGN_KEY_CHECKS = 1;
