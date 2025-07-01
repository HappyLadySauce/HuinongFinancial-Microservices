/*
 Navicat Premium Dump SQL

 Source Server         : 172.18.122.19
 Source Server Type    : MySQL
 Source Server Version : 80405 (8.4.5)
 Source Host           : 172.18.122.19:3306
 Source Schema         : lease

 Target Server Type    : MySQL
 Target Server Version : 80405 (8.4.5)
 File Encoding         : 65001

 Date: 01/07/2025 11:20:01
*/

SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;

-- ----------------------------
-- Table structure for lease_applications
-- ----------------------------
DROP TABLE IF EXISTS `lease_applications`;
CREATE TABLE `lease_applications`  (
  `id` bigint UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '申请ID',
  `application_id` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '申请编号',
  `user_id` bigint UNSIGNED NOT NULL COMMENT '用户ID',
  `applicant_name` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '申请人姓名',
  `product_id` bigint UNSIGNED NOT NULL COMMENT '租赁产品ID',
  `product_code` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '产品编码',
  `name` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '申请名称',
  `type` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '租赁类型',
  `machinery` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL DEFAULT '' COMMENT '设备名称',
  `start_date` date NOT NULL COMMENT '开始日期',
  `end_date` date NOT NULL COMMENT '结束日期',
  `duration` int UNSIGNED NOT NULL COMMENT '租期(天)',
  `daily_rate` decimal(10, 2) NOT NULL COMMENT '日租金',
  `total_amount` decimal(10, 2) NOT NULL COMMENT '总金额',
  `deposit` decimal(10, 2) NULL DEFAULT 0.00 COMMENT '押金',
  `delivery_address` varchar(500) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL DEFAULT '' COMMENT '交付地址',
  `contact_phone` varchar(20) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL DEFAULT '' COMMENT '联系电话',
  `purpose` text CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL COMMENT '使用目的',
  `status` varchar(20) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL DEFAULT 'pending' COMMENT '状态 pending/approved/rejected/cancelled',
  `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE INDEX `uk_application_id`(`application_id` ASC) USING BTREE,
  INDEX `idx_user_id`(`user_id` ASC) USING BTREE,
  INDEX `idx_product_id`(`product_id` ASC) USING BTREE,
  INDEX `idx_status`(`status` ASC) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 7 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_unicode_ci COMMENT = '租赁申请表' ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of lease_applications
-- ----------------------------
INSERT INTO `lease_applications` VALUES (1, 'LS202506010001', 1, '张三', 1, 'LEASE001', '春耕拖拉机租赁', '农机租赁', '约翰迪尔1204拖拉机', '2025-03-01', '2025-04-30', 60, 500.00, 30000.00, 5000.00, '江苏省南京市浦口区永宁街道农业园区', '13800138000', '春季农田耕作', 'approved', '2025-06-30 13:27:39', '2025-06-30 13:27:39');
INSERT INTO `lease_applications` VALUES (2, 'LS202506010002', 2, '李四', 2, 'LEASE002', '水稻收割机租赁', '农机租赁', '久保田4LZ-4J收割机', '2025-09-15', '2025-10-15', 30, 800.00, 24000.00, 10000.00, '江苏省扬州市高邮市汤庄镇', '13800138001', '水稻收获', 'pending', '2025-06-30 13:27:39', '2025-06-30 13:27:39');
INSERT INTO `lease_applications` VALUES (3, 'LS202506020001', 3, '王五', 3, 'LEASE003', '播种机租赁', '农机租赁', '大华宝来2BYF-6播种机', '2025-04-01', '2025-04-15', 14, 300.00, 4200.00, 2000.00, '江苏省盐城市东台市安丰镇', '13800138002', '大豆播种', 'approved', '2025-06-30 13:27:39', '2025-06-30 13:27:39');
INSERT INTO `lease_applications` VALUES (4, 'LA20250701vBu0LY', 1007, '测试用户更新', 1, 'LP001', '测试租赁申请', '挖掘机', '大型挖掘机', '2025-07-01', '2025-07-10', 10, 800.00, 8000.00, 10000.00, '北京市朝阳区建国路1号', '13800138000', '建筑工程施工', 'pending', '2025-07-01 03:08:49', '2025-07-01 03:08:49');
INSERT INTO `lease_applications` VALUES (5, 'LA20250701j2hvoN', 1007, '测试用户更新', 1, 'LP001', '测试租赁申请', '挖掘机', '大型挖掘机', '2025-07-01', '2025-07-10', 10, 800.00, 8000.00, 10000.00, '更新后的地址', '13900139000', '更新后的使用目的', 'approved', '2025-07-01 03:10:37', '2025-07-01 03:10:42');
INSERT INTO `lease_applications` VALUES (6, 'LA20250701Qb0aTm', 1007, '测试用户更新', 1, 'LP001', '测试租赁申请', '挖掘机', '大型挖掘机', '2025-07-01', '2025-07-10', 10, 800.00, 8000.00, 10000.00, '更新后的地址', '13900139000', '更新后的使用目的', 'approved', '2025-07-01 03:18:42', '2025-07-01 03:18:47');

-- ----------------------------
-- Table structure for lease_approvals
-- ----------------------------
DROP TABLE IF EXISTS `lease_approvals`;
CREATE TABLE `lease_approvals`  (
  `id` bigint UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '审批ID',
  `application_id` bigint UNSIGNED NOT NULL COMMENT '申请ID',
  `auditor_id` bigint UNSIGNED NOT NULL COMMENT '审核员ID',
  `auditor_name` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '审核员姓名',
  `action` varchar(20) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '审批动作 approve/reject',
  `suggestions` text CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL COMMENT '审批意见',
  `approved_duration` int UNSIGNED NULL DEFAULT NULL COMMENT '批准租期(天)',
  `approved_amount` decimal(10, 2) NULL DEFAULT NULL COMMENT '批准金额',
  `approved_deposit` decimal(10, 2) NULL DEFAULT NULL COMMENT '批准押金',
  `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  PRIMARY KEY (`id`) USING BTREE,
  INDEX `idx_application_id`(`application_id` ASC) USING BTREE,
  INDEX `idx_auditor_id`(`auditor_id` ASC) USING BTREE,
  INDEX `idx_action`(`action` ASC) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 5 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_unicode_ci COMMENT = '租赁审批记录表' ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of lease_approvals
-- ----------------------------
INSERT INTO `lease_approvals` VALUES (1, 1, 2001, 'B端测试管理员', 'approve', '设备需求合理，客户信用良好，同意租赁', 60, 30000.00, 5000.00, '2025-06-30 13:27:39');
INSERT INTO `lease_approvals` VALUES (2, 3, 2002, '普通操作员1', 'approve', '短期租赁，风险可控，设备适用，建议批准', 14, 4200.00, 2000.00, '2025-06-30 13:27:39');
INSERT INTO `lease_approvals` VALUES (3, 5, 2006, '13452552490', 'approve', '申请已通过审核', 10, 8000.00, 10000.00, '2025-07-01 03:10:42');
INSERT INTO `lease_approvals` VALUES (4, 6, 2006, '13452552490', 'approve', '申请已通过审核', 10, 8000.00, 10000.00, '2025-07-01 03:18:47');

SET FOREIGN_KEY_CHECKS = 1;
