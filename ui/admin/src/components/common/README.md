# 通用组件使用指南

本目录包含了一系列可复用的Vue组件，旨在减少代码重复，统一项目风格，提高开发效率。

## 组件列表

### 1. StatCard - 统计卡片组件

用于显示各种统计数据的卡片组件。

**Props:**
- `value: string | number` - 显示的数值
- `label: string` - 标签文本
- `icon?: any` - 图标组件
- `type?: 'primary' | 'success' | 'warning' | 'danger' | 'info'` - 卡片类型
- `change?: string` - 变化文本
- `changeType?: 'up' | 'down' | 'same'` - 变化趋势
- `showIcon?: boolean` - 是否显示图标
- `clickable?: boolean` - 是否可点击
- `format?: 'number' | 'currency' | 'percent'` - 数值格式化类型

**Events:**
- `click` - 点击事件

**使用示例:**
```vue
<StatCard
  :value="1234"
  label="总用户数"
  :icon="User"
  type="primary"
  change="+12%"
  changeType="up"
  :clickable="true"
  @click="handleCardClick"
/>
```

### 2. StatusTag - 状态标签组件

通用的状态标签组件，支持多种预定义状态。

**Props:**
- `status: string` - 状态值
- `statusMap?: Record<string, { type: string; text?: string; icon?: any }>` - 自定义状态映射
- `showIcon?: boolean` - 是否显示图标
- `size?: 'large' | 'default' | 'small'` - 标签大小
- `effect?: 'dark' | 'light' | 'plain'` - 标签效果
- `customText?: string` - 自定义文本

**使用示例:**
```vue
<StatusTag status="pending" :show-icon="true" />
<StatusTag status="approved" size="small" />
<StatusTag 
  status="custom_status" 
  :status-map="{ custom_status: { type: 'success', text: '自定义状态' } }"
/>
```

### 3. TableActions - 表格操作组件

封装了搜索、过滤、批量操作等表格常用功能。

**Props:**
- `showFilters?: boolean` - 是否显示过滤器
- `showStats?: boolean` - 是否显示统计信息
- `showTableSettings?: boolean` - 是否显示表格设置
- `showRefresh?: boolean` - 是否显示刷新按钮
- `batchActions?: BatchAction[]` - 批量操作配置
- `selectedItems?: any[]` - 选中项
- `total?: number` - 总数据量

**Events:**
- `search` - 搜索事件
- `reset` - 重置事件
- `refresh` - 刷新事件
- `batchAction` - 批量操作事件

**使用示例:**
```vue
<TableActions
  :selected-items="selectedItems"
  :total="total"
  :batch-actions="batchActions"
  @search="handleSearch"
  @refresh="handleRefresh"
>
  <template #filters="{ form }">
    <el-form-item label="状态">
      <el-select v-model="form.status">
        <el-option label="全部" value="" />
        <el-option label="激活" value="active" />
      </el-select>
    </el-form-item>
  </template>
</TableActions>
```

### 4. FormDialog - 表单对话框组件

可复用的表单对话框，支持新增和编辑模式。

**Props:**
- `modelValue: boolean` - 对话框可见性
- `title?: string` - 对话框标题
- `mode?: 'create' | 'edit'` - 模式
- `width?: string` - 对话框宽度
- `initialData?: Record<string, any>` - 初始数据
- `rules?: FormRules` - 表单验证规则

**Events:**
- `submit` - 提交事件
- `cancel` - 取消事件

**使用示例:**
```vue
<FormDialog
  v-model="dialogVisible"
  :mode="editMode"
  :initial-data="formData"
  :rules="formRules"
  @submit="handleSubmit"
>
  <template #form="{ form }">
    <el-form-item label="名称" prop="name">
      <el-input v-model="form.name" />
    </el-form-item>
  </template>
</FormDialog>
```

### 5. ChartCard - 图表卡片组件

封装了ECharts的图表卡片组件。

**Props:**
- `title: string` - 图表标题
- `chartHeight?: string` - 图表高度
- `loading?: boolean` - 加载状态
- `showRefresh?: boolean` - 是否显示刷新按钮
- `showFullscreen?: boolean` - 是否显示全屏按钮
- `chartOptions?: any` - ECharts配置

**Events:**
- `refresh` - 刷新事件

**使用示例:**
```vue
<ChartCard
  title="销售趋势"
  :chart-options="chartOptions"
  :loading="chartLoading"
  @refresh="refreshChart"
/>
```

### 6. PageHeader - 页面头部组件

统一的页面头部组件，包含标题、面包屑、统计等。

**Props:**
- `title: string` - 页面标题
- `subtitle?: string` - 副标题
- `icon?: any` - 页面图标
- `showBack?: boolean` - 是否显示返回按钮
- `showBreadcrumb?: boolean` - 是否显示面包屑
- `breadcrumbs?: BreadcrumbItem[]` - 面包屑数据
- `stats?: StatItem[]` - 统计数据

**使用示例:**
```vue
<PageHeader
  title="用户管理"
  subtitle="管理系统用户"
  :icon="User"
  :breadcrumbs="breadcrumbs"
  :stats="pageStats"
>
  <template #actions>
    <el-button type="primary">新增用户</el-button>
  </template>
</PageHeader>
```

## 使用建议

### 1. 导入方式

```typescript
// 单个组件导入
import { StatCard, StatusTag } from '@/components/common'

// 全部导入
import * as CommonComponents from '@/components/common'
```

### 2. 类型定义

组件已包含完整的TypeScript类型定义，建议使用TypeScript开发以获得更好的开发体验。

### 3. 样式定制

所有组件都支持CSS变量定制，可以通过修改CSS变量来调整组件样式：

```css
:root {
  --stat-card-border-radius: 12px;
  --status-tag-font-size: 12px;
}
```

### 4. 扩展组件

如果需要特殊功能，建议继承现有组件进行扩展，而不是重新创建：

```vue
<template>
  <StatCard v-bind="$attrs" :class="customClass">
    <!-- 自定义内容 -->
  </StatCard>
</template>
```

## 代码重构指南

### 重构前后对比

**重构前:**
```vue
<!-- 原始的统计卡片代码 -->
<el-card class="stat-card" shadow="hover">
  <div class="stat-content">
    <div class="stat-number">{{ stats.pending }}</div>
    <div class="stat-label">待审批</div>
  </div>
</el-card>
```

**重构后:**
```vue
<!-- 使用通用组件 -->
<StatCard
  :value="stats.pending"
  label="待审批"
  type="warning"
  :icon="Clock"
  :clickable="true"
  @click="handleFilter('pending')"
/>
```

### 重构收益

1. **代码减少**: 单个页面代码量减少约30-50%
2. **维护性提升**: 统一的组件逻辑，易于维护和升级
3. **一致性**: 统一的UI风格和交互体验
4. **开发效率**: 快速搭建新页面，减少重复开发
5. **类型安全**: 完整的TypeScript支持

### 逐步迁移策略

1. 新页面直接使用通用组件
2. 存量页面在修改时逐步重构
3. 先重构相对独立的组件部分
4. 最后重构复杂的交互逻辑

## 最佳实践

1. **组件职责单一**: 每个组件只负责一个特定功能
2. **Props设计合理**: 提供合理的默认值和可选配置
3. **事件命名规范**: 使用清晰的事件名称
4. **插槽灵活运用**: 通过插槽提供自定义能力
5. **文档完善**: 为每个组件提供详细的使用文档

通过使用这些通用组件，可以显著提高开发效率，减少代码重复，并确保整个应用的UI一致性。 