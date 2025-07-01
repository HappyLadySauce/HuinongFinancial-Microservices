/*
 Navicat Premium Dump SQL

 Source Server         : 172.18.122.19
 Source Server Type    : MySQL
 Source Server Version : 80405 (8.4.5)
 Source Host           : 172.18.122.19:3306
 Source Schema         : loanproduct

 Target Server Type    : MySQL
 Target Server Version : 80405 (8.4.5)
 File Encoding         : 65001

 Date: 01/07/2025 11:20:21
*/

SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;

-- ----------------------------
-- Table structure for loan_products
-- ----------------------------
DROP TABLE IF EXISTS `loan_products`;
CREATE TABLE `loan_products`  (
  `id` bigint UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '产品ID',
  `product_code` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '产品编码',
  `name` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '产品名称',
  `type` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '产品类型',
  `max_amount` decimal(15, 2) NOT NULL COMMENT '最大金额',
  `min_amount` decimal(15, 2) NULL DEFAULT 1000.00 COMMENT '最小金额',
  `max_duration` int UNSIGNED NULL DEFAULT 60 COMMENT '最大期限(月)',
  `min_duration` int UNSIGNED NULL DEFAULT 1 COMMENT '最小期限(月)',
  `interest_rate` decimal(5, 2) NOT NULL COMMENT '年利率(%)',
  `description` text CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '产品描述',
  `status` tinyint UNSIGNED NULL DEFAULT 1 COMMENT '状态 1:上架 2:下架',
  `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE INDEX `uk_product_code`(`product_code` ASC) USING BTREE,
  INDEX `idx_type`(`type` ASC) USING BTREE,
  INDEX `idx_status`(`status` ASC) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 21 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_unicode_ci COMMENT = '贷款产品表' ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of loan_products
-- ----------------------------
INSERT INTO `loan_products` VALUES (1, 'LOAN001', '农业生产贷', '农业贷', 500000.00, 5000.00, 36, 3, 5.20, '专为农业生产提供的资金支持，支持种植、养殖等农业项目', 1, '2025-06-30 13:27:29', '2025-06-30 13:27:29');
INSERT INTO `loan_products` VALUES (2, 'LOAN002', '农村创业贷', '创业贷', 300000.00, 10000.00, 60, 6, 6.50, '支持农村创业项目的专项贷款，助力乡村振兴发展', 1, '2025-06-30 13:27:29', '2025-06-30 13:27:29');
INSERT INTO `loan_products` VALUES (3, 'LOAN003', '农村消费贷', '消费贷', 100000.00, 1000.00, 24, 1, 7.80, '满足农村居民日常消费需求的个人贷款产品', 1, '2025-06-30 13:27:29', '2025-06-30 13:27:29');
INSERT INTO `loan_products` VALUES (4, 'LOAN004', '农机设备贷', '经营贷', 800000.00, 20000.00, 60, 12, 5.80, '专门用于购买农机设备的经营性贷款', 1, '2025-06-30 13:27:29', '2025-06-30 13:27:29');
INSERT INTO `loan_products` VALUES (20, 'LN001', '更新后的信用贷款', '信用贷款', 600000.00, 5000.00, 48, 3, 0.09, '更新后的产品描述', 1, '2025-07-01 03:08:33', '2025-07-01 03:08:45');

SET FOREIGN_KEY_CHECKS = 1;
