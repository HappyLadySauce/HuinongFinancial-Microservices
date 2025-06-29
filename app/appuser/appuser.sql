-- AppUser服务数据库初始化脚本
-- 数据库: appuser  
-- 描述: 本库仅用于存储C端(App)用户的基本信息
-- 连接: appuser:appuser@tcp(10.10.10.6:3306)/appuser?charset=utf8mb4&parseTime=True&loc=Local

SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;

-- 创建数据库
CREATE DATABASE IF NOT EXISTS `appuser` DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
USE `appuser`;

-- ----------------------------
-- App用户表
-- ----------------------------
DROP TABLE IF EXISTS `app_users`;
CREATE TABLE `app_users` (
  `id` bigint UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '用户ID',
  `phone` varchar(20) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '手机号',
  `password` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '密码哈希',
  `name` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '用户姓名',
  `nickname` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT '' COMMENT '昵称',
  `age` tinyint UNSIGNED DEFAULT 0 COMMENT '年龄',
  `gender` tinyint UNSIGNED DEFAULT 0 COMMENT '性别 0:未知 1:男 2:女',
  `occupation` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT '' COMMENT '职业',
  `address` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT '' COMMENT '联系地址',
  `income` decimal(10,2) DEFAULT 0.00 COMMENT '月收入',
  `status` tinyint UNSIGNED DEFAULT 1 COMMENT '状态 1:正常 2:冻结 3:禁用',
  `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_phone` (`phone`),
  KEY `idx_status` (`status`),
  KEY `idx_created_at` (`created_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='App用户表';

-- ----------------------------
-- 初始化测试数据
-- 密码哈希对应的明文密码均为: 123456
-- 生成方式: bcrypt.GenerateFromPassword([]byte("123456"), bcrypt.DefaultCost)
-- ----------------------------
INSERT INTO `app_users` (`id`, `phone`, `password`, `name`, `nickname`, `age`, `gender`, `occupation`, `address`, `income`, `status`) VALUES
(1001, '13800138000', '$2a$10$N.zmdr9k7uOCQb96looxdOm8DtNkdtK67ZNWEaMMEPYANxJAPJV6C', '王小明', 'C端测试用户', 30, 1, '软件工程师', '北京市朝阳区', 15000.00, 1),
(1002, '13800138001', '$2a$10$N.zmdr9k7uOCQb96looxdOm8DtNkdtK67ZNWEaMMEPYANxJAPJV6C', '张小农', '农业小张', 32, 1, '农民', '北京市朝阳区农业园区', 8500.00, 1),
(1003, '13800138002', '$2a$10$N.zmdr9k7uOCQb96looxdOm8DtNkdtK67ZNWEaMMEPYANxJAPJV6C', '李惠农', '惠农小李', 29, 2, '农村基层干部', '北京市昌平区新农村', 6200.00, 1),
(1004, '13800138003', '$2a$10$N.zmdr9k7uOCQb96looxdOm8DtNkdtK67ZNWEaMMEPYANxJAPJV6C', '赵机达', '农机达人', 35, 1, '农机师', '北京市房山区机械村', 12000.00, 1),
(1005, '13800138004', '$2a$10$N.zmdr9k7uOCQb96looxdOm8DtNkdtK67ZNWEaMMEPYANxJAPJV6C', '赵养殖', '养殖老赵', 40, 1, '畜牧员', '北京市大兴区养殖基地', 7500.00, 1),
(1006, '13800138005', '$2a$10$N.zmdr9k7uOCQb96looxdOm8DtNkdtK67ZNWEaMMEPYANxJAPJV6C', '刘有机', '有机农夫', 31, 1, '农业推广员', '北京市顺义区有机农场', 9000.00, 1);

SET FOREIGN_KEY_CHECKS = 1;
