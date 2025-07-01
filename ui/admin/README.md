# 数字惠农OA管理系统前端

## 项目简介

数字惠农OA管理系统前端是一个基于Vue 3 + TypeScript构建的现代化后台管理系统，主要用于农业贷款审批和系统管理。

## 技术栈

- **框架**: Vue 3.3+ (Composition API)
- **构建工具**: Vite 4+
- **语言**: TypeScript 5+
- **UI库**: Element Plus 2.4+
- **路由**: Vue Router 4+
- **状态管理**: Pinia 2+
- **HTTP客户端**: Axios 1.6+
- **日期处理**: Day.js 1.11+
- **包管理器**: pnpm 8+

## 快速开始

### 环境要求

- Node.js >= 18.0.0
- pnpm >= 8.0.0

### 安装依赖

```bash
pnpm install
```

### 开发环境

```bash
pnpm dev
```

开发服务器将在 `http://localhost:5173` 启动

### 生产构建

```bash
pnpm build
```

### 本地预览

```bash
pnpm preview
```

## 功能特性

### 🔐 认证与权限
- JWT Token认证
- 基于角色的权限控制（管理员/审批员）
- 路由守卫和页面权限

### 📊 工作台
- 数据统计概览
- 待办事项管理
- AI审批状态监控
- 快捷操作面板
- 最近活动时间线

### 📋 审批管理
- 贷款申请列表查看
- 多维度筛选和搜索
- 快速审批功能
- 详细申请信息查看
- AI风险评估结果
- 文档预览和下载
- 审批历史记录

### 👥 用户管理
- OA用户创建和编辑
- 用户状态管理（启用/禁用）
- 角色权限分配
- 用户信息维护

### 📝 操作日志
- 系统操作记录查看
- 多维度日志筛选
- 操作详情查看
- 日志导出功能

### ⚙️ 系统设置
- AI审批功能配置
- 系统参数管理
- 配置项在线编辑
- 系统状态监控
- 快捷运维操作

## 页面结构

```
src/views/
├── LoginView.vue          # 登录页面
├── LayoutView.vue         # 主布局
├── DashboardView.vue      # 工作台
├── ApprovalView.vue       # 审批看板
├── ApprovalDetailView.vue # 审批详情
├── UsersView.vue          # 用户管理
├── LogsView.vue           # 操作日志
├── SystemView.vue         # 系统设置
└── NotFoundView.vue       # 404页面
```

## API接口

系统通过RESTful API与后端通信，默认后端地址：`http://localhost:8080/api/v1`

主要接口模块：
- 认证接口 (`/admin/login`)
- 工作台接口 (`/admin/dashboard`)
- 审批接口 (`/admin/loans/applications/*`)
- 用户管理接口 (`/admin/users`)
- 日志接口 (`/admin/logs`)
- 系统配置接口 (`/admin/configs`)

## 测试账号

- **管理员**: `admin` / `admin123`
- **审批员**: `reviewer` / `reviewer123`

## 设计特色

### 🎨 现代化UI
- Element Plus组件库
- 响应式设计
- 暗色/亮色主题
- 流畅的交互动画

### 📱 移动端适配
- 响应式布局
- 触摸友好的交互
- 移动端优化

### 🚀 性能优化
- 路由懒加载
- 组件按需加载
- 代码分割优化
- 资源压缩和缓存

### 🔒 安全特性
- XSS防护
- CSRF保护
- 敏感信息脱敏
- 安全的Token管理

## 项目架构

### 组件化设计
- 可复用的通用组件
- 业务组件封装
- 合理的组件划分

### 状态管理
- Pinia状态管理
- 模块化Store设计
- 持久化存储

### 类型安全
- 完整的TypeScript类型定义
- API接口类型约束
- 组件Props类型检查

## 浏览器支持

- Chrome >= 87
- Firefox >= 78
- Safari >= 14
- Edge >= 88

## 开发规范

项目遵循严格的代码规范：
- ESLint代码检查
- Prettier格式化
- Vue官方风格指南
- TypeScript最佳实践

## 贡献指南

1. Fork项目
2. 创建功能分支
3. 提交代码变更
4. 创建Pull Request

## 许可证

本项目使用MIT许可证

## 联系方式

如有问题或建议，请联系开发团队。
