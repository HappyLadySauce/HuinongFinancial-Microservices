-- OAUser服务数据库初始化脚本
-- 数据库: oauser
-- 描述: 本库用于存储B端(后台)用户的基本信息
-- 连接: oauser:oauser@tcp(10.10.10.6:3306)/oauser?charset=utf8mb4&parseTime=True&loc=Local

SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;

-- 创建数据库
CREATE DATABASE IF NOT EXISTS `oauser` DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
USE `oauser`;

-- ----------------------------
-- 后台用户表
-- ----------------------------
DROP TABLE IF EXISTS `oa_users`;
CREATE TABLE `oa_users` (
  `id` bigint UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '用户ID',
  `username` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '用户名',
  `password_hash` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '密码哈希',
  `name` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT '' COMMENT '姓名',
  `email` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT '' COMMENT '邮箱',
  `mobile` varchar(20) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT '' COMMENT '手机号',
  `roles` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT '' COMMENT '角色列表(逗号分隔)',
  `status` tinyint UNSIGNED DEFAULT 1 COMMENT '状态 1:正常 2:禁用',
  `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_username` (`username`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='后台用户表';

-- ----------------------------
-- 初始化测试数据
-- 密码哈希对应的明文密码均为: 123456
-- 生成方式: bcrypt.GenerateFromPassword([]byte("123456"), bcrypt.DefaultCost)
-- ----------------------------
INSERT INTO `oa_users` (`id`, `username`, `password_hash`, `name`, `email`, `mobile`, `roles`, `status`) VALUES
(2001, 'admin', '$2a$10$N.zmdr9k7uOCQb96looxdOm8DtNkdtK67ZNWEaMMEPYANxJAPJV6C', 'B端测试管理员', 'admin@huinong.com', '13900139000', 'admin,manager', 1),
(2002, 'manager', '$2a$10$N.zmdr9k7uOCQb96looxdOm8DtNkdtK67ZNWEaMMEPYANxJAPJV6C', '风控经理', 'manager@huinong.com', '13900139001', 'manager', 1),
(2003, 'operator', '$2a$10$N.zmdr9k7uOCQb96looxdOm8DtNkdtK67ZNWEaMMEPYANxJAPJV6C', '普通操作员', 'operator@huinong.com', '13900139002', 'operator', 1),
(2004, 'auditor', '$2a$10$N.zmdr9k7uOCQb96looxdOm8DtNkdtK67ZNWEaMMEPYANxJAPJV6C', '审核员', 'auditor@huinong.com', '13900139003', 'auditor', 1),
(2005, 'financial', '$2a$10$N.zmdr9k7uOCQb96looxdOm8DtNkdtK67ZNWEaMMEPYANxJAPJV6C', '财务专员', 'financial@huinong.com', '13900139004', 'financial', 1);

SET FOREIGN_KEY_CHECKS = 1;