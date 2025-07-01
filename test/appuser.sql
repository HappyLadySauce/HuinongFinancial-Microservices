/*
 Navicat Premium Dump SQL

 Source Server         : 172.18.122.19
 Source Server Type    : MySQL
 Source Server Version : 80405 (8.4.5)
 Source Host           : 172.18.122.19:3306
 Source Schema         : appuser

 Target Server Type    : MySQL
 Target Server Version : 80405 (8.4.5)
 File Encoding         : 65001

 Date: 01/07/2025 11:19:48
*/

SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;

-- ----------------------------
-- Table structure for app_users
-- ----------------------------
DROP TABLE IF EXISTS `app_users`;
CREATE TABLE `app_users`  (
  `id` bigint UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '用户ID',
  `phone` varchar(20) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '手机号',
  `password` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '密码哈希',
  `name` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '用户姓名',
  `nickname` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL DEFAULT '' COMMENT '昵称',
  `age` tinyint UNSIGNED NULL DEFAULT 0 COMMENT '年龄',
  `gender` tinyint UNSIGNED NULL DEFAULT 0 COMMENT '性别 0:未知 1:男 2:女',
  `occupation` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL DEFAULT '' COMMENT '职业',
  `address` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL DEFAULT '' COMMENT '联系地址',
  `income` decimal(10, 2) NULL DEFAULT 0.00 COMMENT '月收入',
  `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE INDEX `uk_phone`(`phone` ASC) USING BTREE,
  INDEX `idx_created_at`(`created_at` ASC) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 1008 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_unicode_ci COMMENT = 'App用户表' ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of app_users
-- ----------------------------
INSERT INTO `app_users` VALUES (1001, '13800138000', '$2a$10$N.zmdr9k7uOCQb96looxdOm8DtNkdtK67ZNWEaMMEPYANxJAPJV6C', '王小明', 'C端测试用户', 30, 1, '软件工程师', '北京市朝阳区', 15000.00, '2025-06-30 13:27:15', '2025-06-30 13:27:15');
INSERT INTO `app_users` VALUES (1002, '13800138001', '$2a$10$N.zmdr9k7uOCQb96looxdOm8DtNkdtK67ZNWEaMMEPYANxJAPJV6C', '张小农', '农业小张', 32, 1, '农民', '北京市朝阳区农业园区', 8500.00, '2025-06-30 13:27:15', '2025-06-30 13:27:15');
INSERT INTO `app_users` VALUES (1003, '13800138002', '$2a$10$N.zmdr9k7uOCQb96looxdOm8DtNkdtK67ZNWEaMMEPYANxJAPJV6C', '李惠农', '惠农小李', 29, 2, '农村基层干部', '北京市昌平区新农村', 6200.00, '2025-06-30 13:27:15', '2025-06-30 13:27:15');
INSERT INTO `app_users` VALUES (1004, '13800138003', '$2a$10$N.zmdr9k7uOCQb96looxdOm8DtNkdtK67ZNWEaMMEPYANxJAPJV6C', '赵机达', '农机达人', 35, 1, '农机师', '北京市房山区机械村', 12000.00, '2025-06-30 13:27:15', '2025-06-30 13:27:15');
INSERT INTO `app_users` VALUES (1005, '13800138004', '$2a$10$N.zmdr9k7uOCQb96looxdOm8DtNkdtK67ZNWEaMMEPYANxJAPJV6C', '赵养殖', '养殖老赵', 40, 1, '畜牧员', '北京市大兴区养殖基地', 7500.00, '2025-06-30 13:27:15', '2025-06-30 13:27:15');
INSERT INTO `app_users` VALUES (1006, '13800138005', '$2a$10$N.zmdr9k7uOCQb96looxdOm8DtNkdtK67ZNWEaMMEPYANxJAPJV6C', '刘有机', '有机农夫', 31, 1, '农业推广员', '北京市顺义区有机农场', 9000.00, '2025-06-30 13:27:15', '2025-06-30 13:27:15');
INSERT INTO `app_users` VALUES (1007, '13452552490', '$2a$10$kuZyB33jRLY.HYOXoJRl8Oek1WkHeb1fGtdvGlI8BBaHHFPSJBG1i', '测试用户更新', '测试昵称', 25, 1, '软件工程师', '北京市朝阳区', 15000.00, '2025-06-30 16:29:43', '2025-06-30 17:07:29');

SET FOREIGN_KEY_CHECKS = 1;
