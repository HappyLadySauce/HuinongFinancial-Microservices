/*
 Navicat Premium Dump SQL

 Source Server         : 172.18.122.19
 Source Server Type    : MySQL
 Source Server Version : 80405 (8.4.5)
 Source Host           : 172.18.122.19:3306
 Source Schema         : loan

 Target Server Type    : MySQL
 Target Server Version : 80405 (8.4.5)
 File Encoding         : 65001

 Date: 01/07/2025 11:20:15
*/

SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;

-- ----------------------------
-- Table structure for loan_applications
-- ----------------------------
DROP TABLE IF EXISTS `loan_applications`;
CREATE TABLE `loan_applications`  (
  `id` bigint UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '申请ID',
  `application_id` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '申请编号',
  `user_id` bigint UNSIGNED NOT NULL COMMENT '用户ID',
  `applicant_name` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '申请人姓名',
  `product_id` bigint UNSIGNED NOT NULL COMMENT '贷款产品ID',
  `name` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '申请名称',
  `type` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '贷款类型',
  `amount` decimal(15, 2) NOT NULL COMMENT '申请金额',
  `duration` int UNSIGNED NOT NULL COMMENT '贷款期限(月)',
  `purpose` text CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL COMMENT '贷款用途',
  `status` varchar(20) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL DEFAULT 'pending' COMMENT '状态 pending/approved/rejected/cancelled',
  `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE INDEX `uk_application_id`(`application_id` ASC) USING BTREE,
  INDEX `idx_user_id`(`user_id` ASC) USING BTREE,
  INDEX `idx_product_id`(`product_id` ASC) USING BTREE,
  INDEX `idx_status`(`status` ASC) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 8 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_unicode_ci COMMENT = '贷款申请表' ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of loan_applications
-- ----------------------------
INSERT INTO `loan_applications` VALUES (1, 'LN202506010001', 1, '张三', 1, '春耕资金贷款', '农业贷', 50000.00, 12, '春季农业生产资金', 'approved', '2025-06-30 13:27:24', '2025-06-30 13:27:24');
INSERT INTO `loan_applications` VALUES (2, 'LN202506010002', 2, '李四', 2, '农机购买贷款', '经营贷', 200000.00, 36, '购买农机设备', 'pending', '2025-06-30 13:27:24', '2025-06-30 13:27:24');
INSERT INTO `loan_applications` VALUES (3, 'LN202506020001', 3, '王五', 3, '日常消费贷款', '消费贷', 30000.00, 24, '家庭日常开支', 'approved', '2025-06-30 13:27:24', '2025-06-30 13:27:24');
INSERT INTO `loan_applications` VALUES (4, 'LN20250701NJdY3q', 1007, '测试用户更新', 1, '测试贷款申请', '信用贷款', 120000.00, 18, '更新后的贷款用途', 'approved', '2025-07-01 02:21:02', '2025-07-01 02:21:07');
INSERT INTO `loan_applications` VALUES (5, 'LN20250701LD3Sz2', 1007, '测试用户更新', 1, '测试贷款申请', '信用贷款', 120000.00, 18, '更新后的贷款用途', 'approved', '2025-07-01 02:58:08', '2025-07-01 02:58:13');
INSERT INTO `loan_applications` VALUES (6, 'LN20250701TYY9cw', 1007, '测试用户更新', 1, '测试贷款申请', '信用贷款', 120000.00, 18, '更新后的贷款用途', 'approved', '2025-07-01 03:10:45', '2025-07-01 03:10:51');
INSERT INTO `loan_applications` VALUES (7, 'LN20250701dKqwE5', 1007, '测试用户更新', 1, '测试贷款申请', '信用贷款', 120000.00, 18, '更新后的贷款用途', 'approved', '2025-07-01 03:18:50', '2025-07-01 03:18:56');

-- ----------------------------
-- Table structure for loan_approvals
-- ----------------------------
DROP TABLE IF EXISTS `loan_approvals`;
CREATE TABLE `loan_approvals`  (
  `id` bigint UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '审批ID',
  `application_id` bigint UNSIGNED NOT NULL COMMENT '申请ID',
  `auditor_id` bigint UNSIGNED NOT NULL COMMENT '审核员ID',
  `auditor_name` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '审核员姓名',
  `action` varchar(20) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '审批动作 approve/reject',
  `suggestions` text CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL COMMENT '审批意见',
  `approved_amount` decimal(15, 2) NULL DEFAULT NULL COMMENT '批准金额',
  `approved_duration` int UNSIGNED NULL DEFAULT NULL COMMENT '批准期限(月)',
  `interest_rate` decimal(5, 2) NULL DEFAULT NULL COMMENT '利率(%)',
  `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  PRIMARY KEY (`id`) USING BTREE,
  INDEX `idx_application_id`(`application_id` ASC) USING BTREE,
  INDEX `idx_auditor_id`(`auditor_id` ASC) USING BTREE,
  INDEX `idx_action`(`action` ASC) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 7 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_unicode_ci COMMENT = '贷款审批记录表' ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of loan_approvals
-- ----------------------------
INSERT INTO `loan_approvals` VALUES (1, 1, 2001, 'B端测试管理员', 'approve', '申请材料完整，收入稳定，同意批准', 50000.00, 12, 5.20, '2025-06-30 13:27:24');
INSERT INTO `loan_approvals` VALUES (2, 3, 2002, '普通操作员1', 'approve', '消费用途合理，风险可控，建议批准', 30000.00, 24, 7.80, '2025-06-30 13:27:24');
INSERT INTO `loan_approvals` VALUES (3, 4, 2006, '13452552490', 'approve', '申请已通过审核', 100000.00, 12, 0.08, '2025-07-01 02:21:07');
INSERT INTO `loan_approvals` VALUES (4, 5, 2006, '13452552490', 'approve', '申请已通过审核', 100000.00, 12, 0.08, '2025-07-01 02:58:13');
INSERT INTO `loan_approvals` VALUES (5, 6, 2006, '13452552490', 'approve', '申请已通过审核', 100000.00, 12, 0.08, '2025-07-01 03:10:51');
INSERT INTO `loan_approvals` VALUES (6, 7, 2006, '13452552490', 'approve', '申请已通过审核', 100000.00, 12, 0.08, '2025-07-01 03:18:56');

SET FOREIGN_KEY_CHECKS = 1;
