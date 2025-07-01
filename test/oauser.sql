/*
 Navicat Premium Dump SQL

 Source Server         : 172.18.122.19
 Source Server Type    : MySQL
 Source Server Version : 80405 (8.4.5)
 Source Host           : 172.18.122.19:3306
 Source Schema         : oauser

 Target Server Type    : MySQL
 Target Server Version : 80405 (8.4.5)
 File Encoding         : 65001

 Date: 01/07/2025 11:20:29
*/

SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;

-- ----------------------------
-- Table structure for oa_users
-- ----------------------------
DROP TABLE IF EXISTS `oa_users`;
CREATE TABLE `oa_users`  (
  `id` bigint UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '用户ID',
  `phone` varchar(20) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '手机号',
  `password_hash` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '密码哈希',
  `name` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL DEFAULT '' COMMENT '姓名',
  `nickname` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL DEFAULT '' COMMENT '昵称',
  `age` tinyint UNSIGNED NULL DEFAULT 0 COMMENT '年龄',
  `gender` tinyint UNSIGNED NULL DEFAULT 0 COMMENT '性别 0:未知 1:男 2:女',
  `role` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL DEFAULT '' COMMENT '管理员(admin)/普通操作员(operator)',
  `status` tinyint UNSIGNED NULL DEFAULT 1 COMMENT '状态 1:正常 2:禁用',
  `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE INDEX `uk_phone`(`phone` ASC) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 2007 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_unicode_ci COMMENT = '后台用户表' ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of oa_users
-- ----------------------------
INSERT INTO `oa_users` VALUES (2001, '13900139000', '$2a$10$N.zmdr9k7uOCQb96looxdOm8DtNkdtK67ZNWEaMMEPYANxJAPJV6C', 'B端测试管理员', 'admin', 30, 1, 'admin', 1, '2025-06-30 13:27:20', '2025-06-30 13:27:20');
INSERT INTO `oa_users` VALUES (2002, '13900139001', '$2a$10$N.zmdr9k7uOCQb96looxdOm8DtNkdtK67ZNWEaMMEPYANxJAPJV6C', '普通操作员1', 'operator1', 30, 1, 'operator', 1, '2025-06-30 13:27:20', '2025-06-30 13:27:20');
INSERT INTO `oa_users` VALUES (2003, '13900139002', '$2a$10$N.zmdr9k7uOCQb96looxdOm8DtNkdtK67ZNWEaMMEPYANxJAPJV6C', '普通操作员2', 'operator2', 30, 1, 'operator', 1, '2025-06-30 13:27:20', '2025-06-30 13:27:20');
INSERT INTO `oa_users` VALUES (2004, '13900139003', '$2a$10$N.zmdr9k7uOCQb96looxdOm8DtNkdtK67ZNWEaMMEPYANxJAPJV6C', '普通操作员3', 'operator3', 30, 1, 'operator', 1, '2025-06-30 13:27:20', '2025-06-30 13:27:20');
INSERT INTO `oa_users` VALUES (2005, '13452552491', '$2a$10$.5.yRts3P3be20UuEEUrBuvcjL2YTeuFwZryUycKoqpyBNcKD44L.', '', '', 0, 0, 'operator', 1, '2025-06-30 16:29:27', '2025-06-30 16:29:27');
INSERT INTO `oa_users` VALUES (2006, '13452552490', '$2a$10$F4qwKzU5dUm.hsLzY9Jgxu.in1RgAGWaaIDE7asafnYg70jyOn3du', '', '', 0, 0, 'admin', 1, '2025-06-30 16:29:33', '2025-06-30 16:29:33');

SET FOREIGN_KEY_CHECKS = 1;
