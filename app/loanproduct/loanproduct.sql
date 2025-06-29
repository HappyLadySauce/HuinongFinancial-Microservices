-- LoanProduct服务数据库初始化脚本
-- 数据库: loanproduct
-- 连接: loanproduct:loanproduct@tcp(10.10.10.6:3306)/loanproduct?charset=utf8mb4&parseTime=True&loc=Local
-- 职责：专门管理贷款产品信息

SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;

-- 创建数据库
CREATE DATABASE IF NOT EXISTS `loanproduct` DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
USE `loanproduct`;

-- ----------------------------
-- 贷款产品表
-- ----------------------------
DROP TABLE IF EXISTS `loan_products`;
CREATE TABLE `loan_products` (
  `id` bigint UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '产品ID',
  `product_code` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '产品编码',
  `name` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '产品名称',
  `type` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '产品类型',
  `max_amount` decimal(15,2) NOT NULL COMMENT '最大金额',
  `min_amount` decimal(15,2) DEFAULT 1000.00 COMMENT '最小金额',
  `max_duration` int UNSIGNED DEFAULT 60 COMMENT '最大期限(月)',
  `min_duration` int UNSIGNED DEFAULT 1 COMMENT '最小期限(月)',
  `interest_rate` decimal(5,2) NOT NULL COMMENT '年利率(%)',
  `description` text CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '产品描述',
  `status` tinyint UNSIGNED DEFAULT 1 COMMENT '状态 1:上架 2:下架',
  `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_product_code` (`product_code`),
  KEY `idx_type` (`type`),
  KEY `idx_status` (`status`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='贷款产品表';

-- ----------------------------
-- 初始化数据
-- ----------------------------
INSERT INTO `loan_products` (`product_code`, `name`, `type`, `max_amount`, `min_amount`, `max_duration`, `min_duration`, `interest_rate`, `description`, `status`) VALUES
('LOAN001', '农业生产贷', '农业贷', 500000.00, 5000.00, 36, 3, 5.20, '专为农业生产提供的资金支持，支持种植、养殖等农业项目', 1),
('LOAN002', '农村创业贷', '创业贷', 300000.00, 10000.00, 60, 6, 6.50, '支持农村创业项目的专项贷款，助力乡村振兴发展', 1),
('LOAN003', '农村消费贷', '消费贷', 100000.00, 1000.00, 24, 1, 7.80, '满足农村居民日常消费需求的个人贷款产品', 1),
('LOAN004', '农机设备贷', '经营贷', 800000.00, 20000.00, 60, 12, 5.80, '专门用于购买农机设备的经营性贷款', 1);

SET FOREIGN_KEY_CHECKS = 1;
