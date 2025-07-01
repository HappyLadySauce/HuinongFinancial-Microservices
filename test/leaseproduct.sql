/*
 Navicat Premium Dump SQL

 Source Server         : 172.18.122.19
 Source Server Type    : MySQL
 Source Server Version : 80405 (8.4.5)
 Source Host           : 172.18.122.19:3306
 Source Schema         : leaseproduct

 Target Server Type    : MySQL
 Target Server Version : 80405 (8.4.5)
 File Encoding         : 65001

 Date: 01/07/2025 11:20:09
*/

SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;

-- ----------------------------
-- Table structure for lease_products
-- ----------------------------
DROP TABLE IF EXISTS `lease_products`;
CREATE TABLE `lease_products`  (
  `id` bigint UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '产品ID',
  `product_code` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '产品编码',
  `name` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '产品名称',
  `type` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '租赁类型',
  `machinery` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL DEFAULT '' COMMENT '设备名称',
  `brand` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL DEFAULT '' COMMENT '品牌',
  `model` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL DEFAULT '' COMMENT '型号',
  `daily_rate` decimal(10, 2) NOT NULL COMMENT '日租金',
  `deposit` decimal(10, 2) NULL DEFAULT 0.00 COMMENT '押金',
  `max_duration` int UNSIGNED NULL DEFAULT 365 COMMENT '最大租期(天)',
  `min_duration` int UNSIGNED NULL DEFAULT 1 COMMENT '最小租期(天)',
  `description` text CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '产品描述',
  `inventory_count` int UNSIGNED NULL DEFAULT 0 COMMENT '库存数量',
  `available_count` int UNSIGNED NULL DEFAULT 0 COMMENT '可用数量',
  `status` tinyint UNSIGNED NULL DEFAULT 1 COMMENT '状态 1:上架 2:下架',
  `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE INDEX `uk_product_code`(`product_code` ASC) USING BTREE,
  INDEX `idx_type`(`type` ASC) USING BTREE,
  INDEX `idx_brand`(`brand` ASC) USING BTREE,
  INDEX `idx_status`(`status` ASC) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 20 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_unicode_ci COMMENT = '租赁产品表' ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of lease_products
-- ----------------------------
INSERT INTO `lease_products` VALUES (1, 'LEASE001', '大型拖拉机租赁', '农机租赁', '约翰迪尔1204拖拉机', '约翰迪尔', '1204', 500.00, 5000.00, 90, 1, '大型农业拖拉机租赁，适用于耕地、播种等大面积作业', 5, 3, 1, '2025-06-30 13:27:34', '2025-06-30 13:27:34');
INSERT INTO `lease_products` VALUES (2, 'LEASE002', '联合收割机租赁', '农机租赁', '久保田4LZ-4J收割机', '久保田', '4LZ-4J', 800.00, 10000.00, 60, 1, '联合收割机租赁，提高粮食收获效率', 3, 2, 1, '2025-06-30 13:27:34', '2025-06-30 13:27:34');
INSERT INTO `lease_products` VALUES (3, 'LEASE003', '精密播种机租赁', '农机租赁', '大华宝来2BYF-6播种机', '大华宝来', '2BYF-6', 300.00, 2000.00, 30, 1, '精密播种机租赁，提高播种质量和效率', 4, 4, 1, '2025-06-30 13:27:34', '2025-06-30 13:27:34');
INSERT INTO `lease_products` VALUES (4, 'LEASE004', '农产品运输车租赁', '车辆租赁', '福田货车', '福田', 'BJ1043V9JDA-F3', 400.00, 3000.00, 180, 1, '农产品专用运输车辆租赁服务', 8, 6, 1, '2025-06-30 13:27:34', '2025-06-30 13:27:34');
INSERT INTO `lease_products` VALUES (19, 'LP001', '更新后的挖掘机', '挖掘机', '大型挖掘机', '卡特彼勒', 'CAT320D', 850.00, 12000.00, 300, 3, '更新后的产品描述', 5, 5, 1, '2025-07-01 03:08:32', '2025-07-01 03:08:41');

SET FOREIGN_KEY_CHECKS = 1;
