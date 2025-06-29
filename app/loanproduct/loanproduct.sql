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
  `category` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT 'general' COMMENT '产品分类',
  `max_amount` decimal(15,2) NOT NULL COMMENT '最大金额',
  `min_amount` decimal(15,2) DEFAULT 1000.00 COMMENT '最小金额',
  `max_duration` int UNSIGNED DEFAULT 60 COMMENT '最大期限(月)',
  `min_duration` int UNSIGNED DEFAULT 1 COMMENT '最小期限(月)',
  `base_interest_rate` decimal(5,2) NOT NULL COMMENT '基准利率(%)',
  `min_interest_rate` decimal(5,2) NOT NULL COMMENT '最低利率(%)',
  `max_interest_rate` decimal(5,2) NOT NULL COMMENT '最高利率(%)',
  `description` text CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '产品描述',
  `features` json DEFAULT NULL COMMENT '产品特点',
  `requirements` json DEFAULT NULL COMMENT '申请要求',
  `risk_factors` json DEFAULT NULL COMMENT '风险要素',
  `target_audience` json DEFAULT NULL COMMENT '目标客户群体',
  `approval_process` json DEFAULT NULL COMMENT '审批流程配置',
  `tags` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT '' COMMENT '产品标签',
  `status` tinyint UNSIGNED DEFAULT 1 COMMENT '状态 1:上架 2:下架 3:停用',
  `sort_order` int DEFAULT 0 COMMENT '排序',
  `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_product_code` (`product_code`),
  KEY `idx_type` (`type`),
  KEY `idx_category` (`category`),
  KEY `idx_status` (`status`),
  KEY `idx_sort_order` (`sort_order`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='贷款产品表';

-- ----------------------------
-- 贷款产品分类表
-- ----------------------------
DROP TABLE IF EXISTS `loan_product_categories`;
CREATE TABLE `loan_product_categories` (
  `id` bigint UNSIGNED NOT NULL AUTO_INCREMENT COMMENT 'ID',
  `code` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '分类编码',
  `name` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '分类名称',
  `parent_id` bigint UNSIGNED DEFAULT 0 COMMENT '父分类ID',
  `level` int UNSIGNED DEFAULT 1 COMMENT '分类级别',
  `description` text CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci COMMENT '分类描述',
  `icon` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT '' COMMENT '分类图标',
  `sort_order` int DEFAULT 0 COMMENT '排序',
  `status` tinyint UNSIGNED DEFAULT 1 COMMENT '状态 1:启用 0:禁用',
  `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_code` (`code`),
  KEY `idx_parent_id` (`parent_id`),
  KEY `idx_sort_order` (`sort_order`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='贷款产品分类表';

-- ----------------------------
-- 贷款产品利率配置表
-- ----------------------------
DROP TABLE IF EXISTS `loan_interest_rates`;
CREATE TABLE `loan_interest_rates` (
  `id` bigint UNSIGNED NOT NULL AUTO_INCREMENT COMMENT 'ID',
  `product_id` bigint UNSIGNED NOT NULL COMMENT '产品ID',
  `product_code` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '产品编码',
  `duration_min` int UNSIGNED NOT NULL COMMENT '最小期限(月)',
  `duration_max` int UNSIGNED NOT NULL COMMENT '最大期限(月)',
  `amount_min` decimal(15,2) NOT NULL COMMENT '最小金额',
  `amount_max` decimal(15,2) NOT NULL COMMENT '最大金额',
  `base_rate` decimal(5,2) NOT NULL COMMENT '基准利率(%)',
  `discount_rate` decimal(5,2) DEFAULT 0.00 COMMENT '优惠利率(%)',
  `risk_adjustment` decimal(5,2) DEFAULT 0.00 COMMENT '风险调整率(%)',
  `final_rate` decimal(5,2) NOT NULL COMMENT '最终利率(%)',
  `effective_date` date NOT NULL COMMENT '生效日期',
  `expire_date` date DEFAULT NULL COMMENT '失效日期',
  `status` tinyint UNSIGNED DEFAULT 1 COMMENT '状态 1:启用 0:禁用',
  `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  KEY `idx_product_id` (`product_id`),
  KEY `idx_product_code` (`product_code`),
  KEY `idx_effective_date` (`effective_date`),
  KEY `idx_status` (`status`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='贷款产品利率配置表';

-- ----------------------------
-- 贷款产品额度配置表
-- ----------------------------
DROP TABLE IF EXISTS `loan_amount_limits`;
CREATE TABLE `loan_amount_limits` (
  `id` bigint UNSIGNED NOT NULL AUTO_INCREMENT COMMENT 'ID',
  `product_id` bigint UNSIGNED NOT NULL COMMENT '产品ID',
  `product_code` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '产品编码',
  `customer_type` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '客户类型',
  `credit_level` varchar(20) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '信用等级',
  `min_amount` decimal(15,2) NOT NULL COMMENT '最小金额',
  `max_amount` decimal(15,2) NOT NULL COMMENT '最大金额',
  `min_duration` int UNSIGNED NOT NULL COMMENT '最小期限(月)',
  `max_duration` int UNSIGNED NOT NULL COMMENT '最大期限(月)',
  `conditions` json DEFAULT NULL COMMENT '附加条件',
  `status` tinyint UNSIGNED DEFAULT 1 COMMENT '状态 1:启用 0:禁用',
  `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  KEY `idx_product_id` (`product_id`),
  KEY `idx_product_code` (`product_code`),
  KEY `idx_customer_type` (`customer_type`),
  KEY `idx_credit_level` (`credit_level`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='贷款产品额度配置表';

-- ----------------------------
-- 系统配置表
-- ----------------------------
DROP TABLE IF EXISTS `system_config`;
CREATE TABLE `system_config` (
  `id` bigint UNSIGNED NOT NULL AUTO_INCREMENT COMMENT 'ID',
  `key` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '配置键',
  `value` text CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '配置值',
  `description` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT '' COMMENT '配置描述',
  `type` varchar(20) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT 'string' COMMENT '值类型 string/int/float/json/bool',
  `module` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT 'loanproduct' COMMENT '所属模块',
  `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_key` (`key`),
  KEY `idx_module` (`module`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='系统配置表';

-- ----------------------------
-- 初始化数据
-- ----------------------------

-- 贷款产品分类数据
INSERT INTO `loan_product_categories` (`code`, `name`, `parent_id`, `level`, `description`, `icon`, `sort_order`) VALUES
('loan', '贷款产品', 0, 1, '所有贷款产品的根分类', 'money', 1),
('loan_agriculture', '农业贷款', 1, 2, '农业生产相关贷款', 'plant', 1),
('loan_business', '创业贷款', 1, 2, '创业经营相关贷款', 'business', 2),
('loan_consumer', '消费贷款', 1, 2, '个人消费相关贷款', 'shopping', 3),
('loan_education', '教育贷款', 1, 2, '教育助学相关贷款', 'education', 4),
('loan_housing', '房屋贷款', 1, 2, '房屋购买相关贷款', 'home', 5),
('loan_vehicle', '车辆贷款', 1, 2, '车辆购买相关贷款', 'car', 6);

-- 贷款产品测试数据
INSERT INTO `loan_products` (`product_code`, `name`, `type`, `category`, `max_amount`, `min_amount`, `max_duration`, `min_duration`, `base_interest_rate`, `min_interest_rate`, `max_interest_rate`, `description`, `features`, `requirements`, `risk_factors`, `target_audience`, `approval_process`, `tags`, `sort_order`) VALUES
('LOAN001', '农业生产贷', '农业贷', 'loan_agriculture', 500000.00, 5000.00, 36, 3, 5.20, 4.80, 6.50, '专为农业生产提供的资金支持，支持种植、养殖等农业项目', 
 '["利率优惠", "审批快速", "还款灵活", "政策支持"]', 
 '["年满18周岁", "有稳定收入", "有农业项目", "信用良好"]',
 '["天气风险", "市场价格波动", "农产品滞销风险"]',
 '["农民", "农业合作社", "农业企业"]',
 '["初审", "实地调研", "风险评估", "审批", "放款"]',
 '农业,种植,养殖,优惠', 1),

('LOAN002', '农村创业贷', '创业贷', 'loan_business', 300000.00, 10000.00, 60, 6, 6.50, 5.80, 8.00, '支持农村创业项目的专项贷款，助力乡村振兴发展', 
 '["支持创业", "期限较长", "额度灵活", "创业指导"]', 
 '["有创业计划", "年满20周岁", "有一定启动资金", "通过创业评估"]',
 '["创业失败风险", "市场竞争风险", "经营管理风险"]',
 '["农村创业者", "返乡青年", "大学生村官"]',
 '["创业计划审核", "可行性分析", "风险评估", "审批", "放款", "跟踪指导"]',
 '创业,乡村振兴,小微企业', 2),

('LOAN003', '农村消费贷', '消费贷', 'loan_consumer', 100000.00, 1000.00, 24, 1, 7.80, 6.50, 9.20, '满足农村居民日常消费需求的个人贷款产品', 
 '["用途广泛", "申请简便", "快速放款", "随借随还"]', 
 '["年满18周岁", "有稳定收入", "信用记录良好"]',
 '["收入不稳定风险", "过度消费风险"]',
 '["农村居民", "农民工", "个体工商户"]',
 '["身份验证", "收入证明", "信用评估", "快速审批", "放款"]',
 '消费,个人,便民', 3),

('LOAN004', '农机设备贷', '经营贷', 'loan_business', 800000.00, 20000.00, 60, 12, 5.80, 5.20, 7.50, '专门用于购买农机设备的经营性贷款', 
 '["专款专用", "利率优惠", "期限灵活", "设备抵押"]', 
 '["有农机需求", "有还款能力", "提供设备清单", "抵押设备"]',
 '["设备贬值风险", "技术淘汰风险", "市场需求风险"]',
 '["农机合作社", "种植大户", "农业公司"]',
 '["设备评估", "抵押登记", "风险评估", "审批", "放款"]',
 '农机,设备,经营,抵押', 4),

('LOAN005', '教育助学贷', '助学贷', 'loan_education', 50000.00, 2000.00, 96, 12, 4.50, 4.00, 5.50, '支持农村学生完成学业的教育贷款', 
 '["利率最低", "期限最长", "宽限期长", "政府贴息"]', 
 '["在校学生", "农村户口", "家庭困难", "学习成绩良好"]',
 '["学业中断风险", "就业不确定风险"]',
 '["在校大学生", "高中生", "职业学校学生"]',
 '["学籍验证", "家庭状况调查", "学习成绩审核", "审批", "分期放款"]',
 '教育,助学,学生,公益', 5);

-- 贷款产品利率配置测试数据
INSERT INTO `loan_interest_rates` (`product_id`, `product_code`, `duration_min`, `duration_max`, `amount_min`, `amount_max`, `base_rate`, `discount_rate`, `risk_adjustment`, `final_rate`, `effective_date`) VALUES
(1, 'LOAN001', 3, 12, 5000.00, 100000.00, 5.20, 0.20, 0.00, 5.00, '2024-01-01'),
(1, 'LOAN001', 13, 36, 5000.00, 500000.00, 5.20, 0.00, 0.30, 5.50, '2024-01-01'),
(2, 'LOAN002', 6, 24, 10000.00, 150000.00, 6.50, 0.30, 0.00, 6.20, '2024-01-01'),
(2, 'LOAN002', 25, 60, 10000.00, 300000.00, 6.50, 0.00, 0.50, 7.00, '2024-01-01'),
(3, 'LOAN003', 1, 12, 1000.00, 50000.00, 7.80, 0.50, 0.00, 7.30, '2024-01-01'),
(3, 'LOAN003', 13, 24, 1000.00, 100000.00, 7.80, 0.00, 0.40, 8.20, '2024-01-01');

-- 贷款产品额度配置测试数据
INSERT INTO `loan_amount_limits` (`product_id`, `product_code`, `customer_type`, `credit_level`, `min_amount`, `max_amount`, `min_duration`, `max_duration`, `conditions`) VALUES
(1, 'LOAN001', '个人农户', 'AAA', 5000.00, 200000.00, 3, 36, '["有土地承包权", "近三年无逾期记录"]'),
(1, 'LOAN001', '农业合作社', 'AA', 20000.00, 500000.00, 6, 36, '["注册满1年", "有固定经营场所", "提供担保"]'),
(2, 'LOAN002', '返乡创业', 'A', 10000.00, 150000.00, 6, 48, '["有创业计划书", "提供担保人"]'),
(2, 'LOAN002', '大学生创业', 'AA', 20000.00, 300000.00, 12, 60, '["大学毕业证", "创业项目评估通过"]'),
(3, 'LOAN003', '普通农户', 'A', 1000.00, 50000.00, 1, 24, '["有稳定收入来源"]'),
(3, 'LOAN003', '农民工', 'AA', 5000.00, 100000.00, 3, 24, '["有工作证明", "连续缴纳社保6个月以上"]');

-- 系统配置测试数据
INSERT INTO `system_config` (`key`, `value`, `description`, `type`, `module`) VALUES
('loan.product.auto_publish', 'false', '贷款产品自动上架', 'bool', 'loanproduct'),
('loan.product.approval_required', 'true', '产品上架需要审批', 'bool', 'loanproduct'),
('loan.interest.update_frequency', '30', '利率更新频率(天)', 'int', 'loanproduct'),
('loan.amount.max_limit', '1000000', '单笔贷款最大金额限制', 'float', 'loanproduct'),
('loan.duration.max_limit', '120', '贷款最大期限限制(月)', 'int', 'loanproduct'),
('system.maintenance_mode', 'false', '系统维护模式', 'bool', 'system');

SET FOREIGN_KEY_CHECKS = 1;
