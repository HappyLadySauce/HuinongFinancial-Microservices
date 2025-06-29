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
  `category` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT 'equipment' COMMENT '产品分类',
  `machinery` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT '' COMMENT '设备名称',
  `brand` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT '' COMMENT '品牌',
  `model` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT '' COMMENT '型号',
  `serial_number` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT '' COMMENT '设备序列号',
  `daily_rate` decimal(10,2) NOT NULL COMMENT '日租金',
  `weekly_rate` decimal(10,2) DEFAULT NULL COMMENT '周租金',
  `monthly_rate` decimal(10,2) DEFAULT NULL COMMENT '月租金',
  `yearly_rate` decimal(10,2) DEFAULT NULL COMMENT '年租金',
  `deposit` decimal(10,2) DEFAULT 0.00 COMMENT '押金',
  `max_duration` int UNSIGNED DEFAULT 365 COMMENT '最大租期(天)',
  `min_duration` int UNSIGNED DEFAULT 1 COMMENT '最小租期(天)',
  `description` text CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '产品描述',
  `specifications` json DEFAULT NULL COMMENT '设备规格',
  `features` json DEFAULT NULL COMMENT '产品特点',
  `requirements` json DEFAULT NULL COMMENT '租赁要求',
  `risk_factors` json DEFAULT NULL COMMENT '风险要素',
  `maintenance_info` json DEFAULT NULL COMMENT '维护信息',
  `images` json DEFAULT NULL COMMENT '产品图片',
  `tags` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT '' COMMENT '产品标签',
  `inventory_count` int UNSIGNED DEFAULT 0 COMMENT '库存数量',
  `available_count` int UNSIGNED DEFAULT 0 COMMENT '可用数量',
  `maintenance_count` int UNSIGNED DEFAULT 0 COMMENT '维护中数量',
  `rented_count` int UNSIGNED DEFAULT 0 COMMENT '已租出数量',
  `status` tinyint UNSIGNED DEFAULT 1 COMMENT '状态 1:可租 2:维护中 3:已租完 4:停用',
  `sort_order` int DEFAULT 0 COMMENT '排序',
  `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_product_code` (`product_code`),
  KEY `idx_type` (`type`),
  KEY `idx_category` (`category`),
  KEY `idx_brand` (`brand`),
  KEY `idx_status` (`status`),
  KEY `idx_sort_order` (`sort_order`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='租赁产品表';

-- ----------------------------
-- 租赁产品分类表
-- ----------------------------
DROP TABLE IF EXISTS `lease_product_categories`;
CREATE TABLE `lease_product_categories` (
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
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='租赁产品分类表';

-- ----------------------------
-- 设备库存记录表
-- ----------------------------
DROP TABLE IF EXISTS `equipment_inventory`;
CREATE TABLE `equipment_inventory` (
  `id` bigint UNSIGNED NOT NULL AUTO_INCREMENT COMMENT 'ID',
  `product_id` bigint UNSIGNED NOT NULL COMMENT '产品ID',
  `product_code` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '产品编码',
  `equipment_id` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '设备编号',
  `serial_number` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT '' COMMENT '设备序列号',
  `purchase_date` date DEFAULT NULL COMMENT '采购日期',
  `purchase_price` decimal(12,2) DEFAULT 0.00 COMMENT '采购价格',
  `depreciation_rate` decimal(5,2) DEFAULT 0.00 COMMENT '折旧率(%/年)',
  `current_value` decimal(12,2) DEFAULT 0.00 COMMENT '当前价值',
  `location` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT '' COMMENT '存放位置',
  `condition_level` varchar(20) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT 'good' COMMENT '设备状况 excellent/good/fair/poor',
  `last_maintenance` date DEFAULT NULL COMMENT '上次维护日期',
  `next_maintenance` date DEFAULT NULL COMMENT '下次维护日期',
  `status` tinyint UNSIGNED DEFAULT 1 COMMENT '状态 1:可用 2:租出 3:维护中 4:报废 5:丢失',
  `notes` text CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci COMMENT '备注信息',
  `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_equipment_id` (`equipment_id`),
  KEY `idx_product_id` (`product_id`),
  KEY `idx_product_code` (`product_code`),
  KEY `idx_status` (`status`),
  KEY `idx_location` (`location`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='设备库存记录表';

-- ----------------------------
-- 租赁价格策略表
-- ----------------------------
DROP TABLE IF EXISTS `lease_pricing_strategies`;
CREATE TABLE `lease_pricing_strategies` (
  `id` bigint UNSIGNED NOT NULL AUTO_INCREMENT COMMENT 'ID',
  `product_id` bigint UNSIGNED NOT NULL COMMENT '产品ID',
  `product_code` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '产品编码',
  `strategy_name` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '策略名称',
  `duration_min` int UNSIGNED NOT NULL COMMENT '最小租期(天)',
  `duration_max` int UNSIGNED NOT NULL COMMENT '最大租期(天)',
  `base_daily_rate` decimal(10,2) NOT NULL COMMENT '基础日租金',
  `weekly_discount` decimal(5,2) DEFAULT 0.00 COMMENT '周租折扣(%)',
  `monthly_discount` decimal(5,2) DEFAULT 0.00 COMMENT '月租折扣(%)',
  `yearly_discount` decimal(5,2) DEFAULT 0.00 COMMENT '年租折扣(%)',
  `quantity_discount` json DEFAULT NULL COMMENT '数量折扣配置',
  `seasonal_adjustment` json DEFAULT NULL COMMENT '季节性调价配置',
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
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='租赁价格策略表';

-- ----------------------------
-- 设备维护记录表
-- ----------------------------
DROP TABLE IF EXISTS `equipment_maintenance`;
CREATE TABLE `equipment_maintenance` (
  `id` bigint UNSIGNED NOT NULL AUTO_INCREMENT COMMENT 'ID',
  `equipment_id` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '设备编号',
  `product_id` bigint UNSIGNED NOT NULL COMMENT '产品ID',
  `maintenance_type` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '维护类型',
  `maintenance_date` date NOT NULL COMMENT '维护日期',
  `maintenance_cost` decimal(10,2) DEFAULT 0.00 COMMENT '维护费用',
  `maintenance_content` text CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci COMMENT '维护内容',
  `maintenance_staff` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT '' COMMENT '维护人员',
  `parts_replaced` json DEFAULT NULL COMMENT '更换部件记录',
  `before_condition` varchar(20) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT '' COMMENT '维护前状况',
  `after_condition` varchar(20) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT '' COMMENT '维护后状况',
  `next_maintenance_date` date DEFAULT NULL COMMENT '下次维护建议日期',
  `notes` text CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci COMMENT '备注',
  `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  KEY `idx_equipment_id` (`equipment_id`),
  KEY `idx_product_id` (`product_id`),
  KEY `idx_maintenance_date` (`maintenance_date`),
  KEY `idx_maintenance_type` (`maintenance_type`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='设备维护记录表';

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
  `module` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT 'leaseproduct' COMMENT '所属模块',
  `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_key` (`key`),
  KEY `idx_module` (`module`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='系统配置表';

-- ----------------------------
-- 初始化数据
-- ----------------------------

-- 租赁产品分类数据
INSERT INTO `lease_product_categories` (`code`, `name`, `parent_id`, `level`, `description`, `icon`, `sort_order`) VALUES
('lease', '租赁产品', 0, 1, '所有租赁产品的根分类', 'tools', 1),
('lease_machinery', '农机租赁', 1, 2, '农业机械设备租赁', 'tractor', 1),
('lease_vehicle', '车辆租赁', 1, 2, '运输车辆租赁', 'truck', 2),
('lease_irrigation', '灌溉设备', 1, 2, '农田灌溉设备租赁', 'water', 3),
('lease_processing', '加工设备', 1, 2, '农产品加工设备租赁', 'factory', 4),
('lease_harvesting', '收获设备', 1, 2, '收获相关设备租赁', 'harvest', 5),
('lease_planting', '种植设备', 1, 2, '种植相关设备租赁', 'seed', 6);

-- 租赁产品测试数据
INSERT INTO `lease_products` (`product_code`, `name`, `type`, `category`, `machinery`, `brand`, `model`, `serial_number`, `daily_rate`, `weekly_rate`, `monthly_rate`, `yearly_rate`, `deposit`, `max_duration`, `min_duration`, `description`, `specifications`, `features`, `requirements`, `risk_factors`, `maintenance_info`, `tags`, `inventory_count`, `available_count`, `maintenance_count`, `rented_count`, `sort_order`) VALUES
('LEASE001', '大型拖拉机租赁', '农机租赁', 'lease_machinery', '约翰迪尔1204拖拉机', '约翰迪尔', '1204', 'JD1204-2024', 500.00, 3200.00, 12000.00, 140000.00, 5000.00, 90, 1, '大型农业拖拉机租赁，适用于耕地、播种等大面积作业',
 '{"功率": "120马力", "重量": "4.5吨", "作业宽度": "2.5米", "燃油类型": "柴油", "变速箱": "16档", "轮胎规格": "18.4-34"}',
 '["动力强劲", "操作简便", "燃油经济", "维护便利", "适应性强"]',
 '["有驾驶经验", "提供押金", "购买保险", "提供身份证明"]',
 '["机械故障风险", "操作不当风险", "燃油价格波动", "天气影响"]',
 '{"保养周期": "100小时", "主要检查项目": ["机油", "液压油", "齿轮油", "滤清器"], "年检要求": "每年一次"}',
 '拖拉机,耕地,农机,大型', 5, 3, 1, 1, 1),

('LEASE002', '联合收割机租赁', '农机租赁', 'lease_harvesting', '久保田4LZ-4J收割机', '久保田', '4LZ-4J', 'KB4LZ-2024', 800.00, 5200.00, 20000.00, 230000.00, 10000.00, 60, 1, '联合收割机租赁，提高粮食收获效率',
 '{"适用作物": "水稻、小麦", "收获宽度": "2.0米", "粮仓容量": "1.7立方米", "作业效率": "0.8-1.2公顷/小时", "发动机功率": "74马力"}',
 '["效率高", "损失少", "适应性强", "自动化程度高", "操作简单"]',
 '["有操作证书", "提供押金", "收获季节预约", "购买保险"]',
 '["季节性需求风险", "设备损坏风险", "天气影响风险", "作物适应性"]',
 '{"保养周期": "50小时", "重点检查": ["切割器", "脱粒滚筒", "清选筛", "传动带"], "季节保养": "作业前后全面检查"}',
 '收割机,收获,粮食,高效', 3, 2, 0, 1, 2),

('LEASE003', '精密播种机租赁', '农机租赁', 'lease_planting', '大华宝来2BYF-6播种机', '大华宝来', '2BYF-6', 'DH2BYF-2024', 300.00, 1800.00, 7500.00, 85000.00, 2000.00, 30, 1, '精密播种机租赁，提高播种质量和效率',
 '{"播种行数": "6行", "行距": "可调15-25cm", "播种深度": "可调1-5cm", "作业宽度": "1.8米", "种箱容量": "120L"}',
 '["播种精准", "节约种子", "提高出苗率", "作业效率高", "深度一致"]',
 '["有拖拉机配套", "提供押金", "春播季节预约", "种子适配确认"]',
 '["配套设备风险", "种子适应性风险", "土壤条件限制", "操作技术要求"]',
 '{"保养周期": "每季使用前", "关键部件": ["播种盘", "输种管", "开沟器", "覆土器"], "存储要求": "干燥通风处存放"}',
 '播种机,种植,精密,高效', 4, 4, 0, 0, 3),

('LEASE004', '农产品运输车租赁', '车辆租赁', 'lease_vehicle', '福田货车', '福田', 'BJ1043V9JDA-F3', 'FT1043-2024', 400.00, 2500.00, 10000.00, 115000.00, 3000.00, 180, 1, '农产品专用运输车辆租赁服务',
 '{"载重": "3吨", "车型": "厢式货车", "冷藏": "可选", "长宽高": "5.2m×2.0m×2.3m", "发动机": "130马力柴油机"}',
 '["载重大", "保鲜好", "运输稳定", "成本低", "维护简单"]',
 '["有驾驶证", "提供押金", "购买保险", "定期检查"]',
 '["交通事故风险", "货物损坏风险", "道路条件风险", "天气影响"]',
 '{"保养周期": "5000公里", "检查项目": ["刹车系统", "轮胎状况", "冷藏设备", "货厢密封"], "年检": "按法规要求"}',
 '运输,货车,物流,冷链', 8, 6, 1, 1, 4),

('LEASE005', '灌溉水泵租赁', '灌溉设备', 'lease_irrigation', '南方水泵', '南方', 'CDL4-19', 'NF-CDL4-2024', 150.00, 900.00, 3500.00, 40000.00, 500.00, 120, 1, '农田灌溉专用水泵设备租赁',
 '{"流量": "100立方米/时", "扬程": "50米", "功率": "15KW", "进出口径": "150mm", "材质": "不锈钢"}',
 '["流量大", "扬程高", "节能环保", "安装简便", "耐腐蚀"]',
 '["有用电条件", "提供押金", "专业安装", "定期维护"]',
 '["电力供应风险", "水源条件风险", "设备维护风险", "季节性需求"]',
 '{"保养周期": "200小时", "检查内容": ["泵体密封", "轴承润滑", "电机绝缘", "管路连接"], "冬季保养": "防冻保护"}',
 '水泵,灌溉,农田,节能', 10, 8, 1, 1, 5),

('LEASE006', '粮食烘干设备租赁', '加工设备', 'lease_processing', '三久烘干机', '三久', 'SJ-120', 'SJ120-2024', 600.00, 3800.00, 15000.00, 170000.00, 8000.00, 60, 7, '粮食烘干加工设备租赁，保证粮食品质',
 '{"处理量": "20吨/批", "烘干时间": "8-12小时", "适用作物": "水稻、玉米、小麦", "热源": "燃煤/燃气可选", "水分降幅": "30%-14%"}',
 '["处理量大", "烘干均匀", "节能高效", "自动控制", "安全可靠"]',
 '["有场地条件", "提供押金", "专业操作", "安全培训"]',
 '["季节性需求风险", "能源价格风险", "操作技术风险", "安全操作风险"]',
 '{"保养周期": "每批次后", "重点检查": ["燃烧器", "热交换器", "提升机", "温控系统"], "安全检查": "燃气管路和安全装置"}',
 '烘干机,粮食,加工,自动化', 2, 1, 0, 1, 6);

-- 设备库存记录测试数据
INSERT INTO `equipment_inventory` (`product_id`, `product_code`, `equipment_id`, `serial_number`, `purchase_date`, `purchase_price`, `depreciation_rate`, `current_value`, `location`, `condition_level`, `last_maintenance`, `next_maintenance`, `status`, `notes`) VALUES
(1, 'LEASE001', 'JD1204-001', 'JD1204-001-2024', '2024-01-15', 280000.00, 10.00, 252000.00, '设备库A区', 'excellent', '2024-01-01', '2024-04-01', 1, '新购设备，状况良好'),
(1, 'LEASE001', 'JD1204-002', 'JD1204-002-2024', '2024-01-15', 280000.00, 10.00, 252000.00, '设备库A区', 'good', '2024-02-15', '2024-05-15', 2, '当前租出给张三农场'),
(2, 'LEASE002', 'KB4LZ-001', 'KB4LZ-001-2024', '2024-02-01', 350000.00, 12.00, 315000.00, '设备库B区', 'good', '2024-02-01', '2024-06-01', 1, '收割季节前检查完毕'),
(3, 'LEASE003', 'DH2BYF-001', 'DH2BYF-001-2024', '2024-03-01', 85000.00, 8.00, 79900.00, '设备库C区', 'excellent', '2024-03-01', '2024-03-15', 1, '春播设备，使用频率高');

-- 租赁价格策略测试数据
INSERT INTO `lease_pricing_strategies` (`product_id`, `product_code`, `strategy_name`, `duration_min`, `duration_max`, `base_daily_rate`, `weekly_discount`, `monthly_discount`, `yearly_discount`, `quantity_discount`, `seasonal_adjustment`, `effective_date`) VALUES
(1, 'LEASE001', '拖拉机标准价格', 1, 7, 500.00, 0.00, 0.00, 0.00, '{"2台以上": "5%", "5台以上": "10%"}', '{"春季": "20%", "秋季": "15%"}', '2024-01-01'),
(1, 'LEASE001', '拖拉机长租优惠', 8, 90, 500.00, 8.00, 15.00, 25.00, '{"2台以上": "5%", "5台以上": "10%"}', '{"春季": "20%", "秋季": "15%"}', '2024-01-01'),
(2, 'LEASE002', '收割机季节价格', 1, 60, 800.00, 5.00, 12.00, 0.00, '{"3台以上": "8%"}', '{"收获季": "30%"}', '2024-06-01'),
(3, 'LEASE003', '播种机春季价格', 1, 30, 300.00, 10.00, 20.00, 0.00, '{"2台以上": "5%"}', '{"春播季": "25%"}', '2024-03-01');

-- 设备维护记录测试数据
INSERT INTO `equipment_maintenance` (`equipment_id`, `product_id`, `maintenance_type`, `maintenance_date`, `maintenance_cost`, `maintenance_content`, `maintenance_staff`, `parts_replaced`, `before_condition`, `after_condition`, `next_maintenance_date`, `notes`) VALUES
('JD1204-001', 1, '定期保养', '2024-01-01', 800.00, '更换机油、液压油，检查各系统', '李师傅', '["机油滤清器", "液压滤清器"]', 'good', 'excellent', '2024-04-01', '保养完成，设备状态良好'),
('KB4LZ-001', 2, '季前检查', '2024-02-01', 1200.00, '收割前全面检查，更换磨损部件', '王师傅', '["切割刀片", "传动带"]', 'fair', 'good', '2024-06-01', '已为收割季做好准备'),
('DH2BYF-001', 3, '新机调试', '2024-03-01', 300.00, '新设备调试，播种精度校准', '赵师傅', '[]', 'new', 'excellent', '2024-03-15', '新设备调试完成');

-- 系统配置测试数据
INSERT INTO `system_config` (`key`, `value`, `description`, `type`, `module`) VALUES
('lease.product.auto_publish', 'false', '租赁产品自动上架', 'bool', 'leaseproduct'),
('lease.product.approval_required', 'true', '产品上架需要审批', 'bool', 'leaseproduct'),
('lease.inventory.low_stock_threshold', '2', '库存预警阈值', 'int', 'leaseproduct'),
('lease.maintenance.reminder_days', '7', '维护提醒提前天数', 'int', 'leaseproduct'),
('lease.deposit.auto_calculate', 'true', '押金自动计算', 'bool', 'leaseproduct'),
('lease.pricing.seasonal_enabled', 'true', '启用季节性定价', 'bool', 'leaseproduct'),
('lease.equipment.depreciation_method', 'straight_line', '设备折旧方法', 'string', 'leaseproduct'),
('system.maintenance_mode', 'false', '系统维护模式', 'bool', 'system');

SET FOREIGN_KEY_CHECKS = 1; 