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
  `phone` varchar(20) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '手机号',
  `password_hash` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '密码哈希',
  `name` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT '' COMMENT '姓名',
  `nickname` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT '' COMMENT '昵称',
  `age` tinyint UNSIGNED DEFAULT 0 COMMENT '年龄',
  `gender` tinyint UNSIGNED DEFAULT 0 COMMENT '性别 0:未知 1:男 2:女',
  `role` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT '' COMMENT '管理员(admin)/普通操作员(operator)',
  `status` tinyint UNSIGNED DEFAULT 1 COMMENT '状态 1:正常 2:禁用',
  `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_phone` (`phone`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='后台用户表';

-- ----------------------------
-- 初始化测试数据
-- 密码哈希对应的明文密码均为: 123456
-- 生成方式: bcrypt.GenerateFromPassword([]byte("123456"), bcrypt.DefaultCost)
-- ----------------------------
INSERT INTO `oa_users` (`id`, `phone`, `password_hash`, `name`, `nickname`, `age`, `gender`, `role`, `status`) VALUES
(2001, '13900139000', '$2a$10$N.zmdr9k7uOCQb96looxdOm8DtNkdtK67ZNWEaMMEPYANxJAPJV6C', 'B端测试管理员', 'admin', 30, 1, 'admin', 1),
(2002, '13900139001', '$2a$10$N.zmdr9k7uOCQb96looxdOm8DtNkdtK67ZNWEaMMEPYANxJAPJV6C', '普通操作员1', 'operator1', 30, 1, 'operator', 1),
(2003, '13900139002', '$2a$10$N.zmdr9k7uOCQb96looxdOm8DtNkdtK67ZNWEaMMEPYANxJAPJV6C', '普通操作员2', 'operator2', 30, 1, 'operator', 1),
(2004, '13900139003', '$2a$10$N.zmdr9k7uOCQb96looxdOm8DtNkdtK67ZNWEaMMEPYANxJAPJV6C', '普通操作员3', 'operator3', 30, 1, 'operator', 1);

SET FOREIGN_KEY_CHECKS = 1;