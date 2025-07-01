<template>
  <div class="page-header">
    <!-- 面包屑导航 -->
    <el-breadcrumb v-if="showBreadcrumb && breadcrumbs.length > 0" class="page-breadcrumb">
      <el-breadcrumb-item 
        v-for="(item, index) in breadcrumbs" 
        :key="index"
        :to="item.to"
      >
        <el-icon v-if="item.icon">
          <component :is="item.icon" />
        </el-icon>
        {{ item.title }}
      </el-breadcrumb-item>
    </el-breadcrumb>

    <!-- 页面标题区域 -->
    <div class="page-title-section">
      <div class="title-left">
        <!-- 返回按钮 -->
        <el-button 
          v-if="showBack"
          size="small"
          circle
          @click="handleBack"
          class="back-button"
        >
          <el-icon><ArrowLeft /></el-icon>
        </el-button>

        <!-- 页面图标和标题 -->
        <div class="title-content">
          <div class="title-main">
            <el-icon v-if="icon" class="title-icon">
              <component :is="icon" />
            </el-icon>
            <h1 class="page-title">{{ title }}</h1>
            <el-tag v-if="badge" :type="badgeType" size="small" class="title-badge">
              {{ badge }}
            </el-tag>
          </div>
          
          <div v-if="subtitle" class="page-subtitle">
            {{ subtitle }}
          </div>

          <!-- 页面描述 -->
          <div v-if="description" class="page-description">
            {{ description }}
          </div>
        </div>
      </div>

      <!-- 右侧操作区域 -->
      <div class="title-actions">
        <slot name="actions" />
      </div>
    </div>

    <!-- 页面状态信息 -->
    <div v-if="showStatus && (lastUpdate || currentUser)" class="page-status">
      <div class="status-info">
        <span v-if="lastUpdate" class="status-item">
          <el-icon><Clock /></el-icon>
          最后更新: {{ formatTime(lastUpdate) }}
        </span>
        <span v-if="currentUser" class="status-item">
          <el-icon><User /></el-icon>
          操作人: {{ currentUser }}
        </span>
      </div>
    </div>

    <!-- 页面统计信息 -->
    <div v-if="showStats && stats.length > 0" class="page-stats">
      <div class="stats-list">
        <div 
          v-for="(stat, index) in stats"
          :key="index"
          class="stat-item"
          :class="{ 'stat-clickable': stat.clickable }"
          @click="stat.clickable && handleStatClick(stat)"
        >
          <div class="stat-value" :style="{ color: stat.color }">
            {{ stat.value }}
          </div>
          <div class="stat-label">{{ stat.label }}</div>
        </div>
      </div>
    </div>

    <!-- 页面标签 -->
    <div v-if="showTabs && tabs.length > 0" class="page-tabs">
      <el-tabs 
        v-model="activeTab" 
        class="header-tabs"
        @tab-change="handleTabChange"
      >
        <el-tab-pane 
          v-for="tab in tabs"
          :key="tab.name"
          :label="tab.label"
          :name="tab.name"
          :disabled="tab.disabled"
        >
          <template #label>
            <el-icon v-if="tab.icon">
              <component :is="tab.icon" />
            </el-icon>
            {{ tab.label }}
            <el-badge 
              v-if="tab.badge" 
              :value="tab.badge" 
              :type="tab.badgeType || 'danger'"
              class="tab-badge"
            />
          </template>
        </el-tab-pane>
      </el-tabs>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue'
import { useRouter } from 'vue-router'
import { ArrowLeft, Clock, User } from '@element-plus/icons-vue'
import dayjs from 'dayjs'

interface BreadcrumbItem {
  title: string
  to?: string | { name: string; params?: any }
  icon?: any
}

interface TabItem {
  name: string
  label: string
  icon?: any
  disabled?: boolean
  badge?: string | number
  badgeType?: 'primary' | 'success' | 'warning' | 'danger' | 'info'
}

interface StatItem {
  label: string
  value: string | number
  color?: string
  clickable?: boolean
  key?: string
}

interface Props {
  title: string
  subtitle?: string
  description?: string
  icon?: any
  badge?: string
  badgeType?: 'primary' | 'success' | 'warning' | 'danger' | 'info'
  showBack?: boolean
  showBreadcrumb?: boolean
  showStatus?: boolean
  showStats?: boolean
  showTabs?: boolean
  breadcrumbs?: BreadcrumbItem[]
  tabs?: TabItem[]
  stats?: StatItem[]
  lastUpdate?: string | Date
  currentUser?: string
  defaultTab?: string
}

const props = withDefaults(defineProps<Props>(), {
  badgeType: 'primary',
  showBack: false,
  showBreadcrumb: true,
  showStatus: false,
  showStats: false,
  showTabs: false,
  breadcrumbs: () => [],
  tabs: () => [],
  stats: () => []
})

const emit = defineEmits<{
  back: []
  tabChange: [tabName: string]
  statClick: [stat: StatItem]
}>()

const router = useRouter()
const activeTab = ref(props.defaultTab || (props.tabs.length > 0 ? props.tabs[0].name : ''))

const handleBack = () => {
  emit('back')
  if (!emit('back')) {
    router.back()
  }
}

const handleTabChange = (tabName: string) => {
  emit('tabChange', tabName)
}

const handleStatClick = (stat: StatItem) => {
  emit('statClick', stat)
}

const formatTime = (time: string | Date) => {
  return dayjs(time).format('YYYY-MM-DD HH:mm:ss')
}
</script>

<style scoped>
.page-header {
  margin-bottom: 24px;
}

.page-breadcrumb {
  margin-bottom: 16px;
}

.page-title-section {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  margin-bottom: 16px;
}

.title-left {
  flex: 1;
  display: flex;
  align-items: flex-start;
  gap: 12px;
}

.back-button {
  margin-top: 8px;
  flex-shrink: 0;
}

.title-content {
  flex: 1;
}

.title-main {
  display: flex;
  align-items: center;
  gap: 12px;
  margin-bottom: 8px;
}

.title-icon {
  font-size: 24px;
  color: #409eff;
}

.page-title {
  font-size: 24px;
  font-weight: 600;
  color: #333;
  margin: 0;
}

.title-badge {
  margin-left: 8px;
}

.page-subtitle {
  font-size: 16px;
  color: #666;
  margin-bottom: 8px;
}

.page-description {
  font-size: 14px;
  color: #999;
  line-height: 1.5;
}

.title-actions {
  flex-shrink: 0;
  display: flex;
  align-items: center;
  gap: 12px;
}

.page-status {
  background: #f8f9fa;
  border-radius: 8px;
  padding: 12px 16px;
  margin-bottom: 16px;
}

.status-info {
  display: flex;
  gap: 24px;
}

.status-item {
  display: flex;
  align-items: center;
  gap: 6px;
  font-size: 13px;
  color: #666;
}

.page-stats {
  background: linear-gradient(135deg, #f5f7fa 0%, #c3cfe2 100%);
  border-radius: 12px;
  padding: 20px;
  margin-bottom: 16px;
}

.stats-list {
  display: flex;
  gap: 32px;
}

.stat-item {
  text-align: center;
  transition: all 0.3s ease;
}

.stat-clickable {
  cursor: pointer;
}

.stat-clickable:hover {
  transform: translateY(-2px);
}

.stat-value {
  font-size: 28px;
  font-weight: bold;
  margin-bottom: 4px;
}

.stat-label {
  font-size: 14px;
  color: #666;
}

.page-tabs {
  margin-top: 16px;
}

.header-tabs {
  border-bottom: 1px solid #e4e7ed;
}

.tab-badge {
  margin-left: 6px;
}

:deep(.el-breadcrumb__item:last-child .el-breadcrumb__inner) {
  color: #409eff;
}

:deep(.el-tabs__item) {
  height: 48px;
  line-height: 48px;
  font-size: 14px;
}

:deep(.el-tabs__nav-wrap::after) {
  background-color: #e4e7ed;
}
</style> 