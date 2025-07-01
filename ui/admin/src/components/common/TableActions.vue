<template>
  <div class="table-actions">
    <!-- 搜索和过滤区域 -->
    <el-card v-if="showFilters" shadow="never" class="filter-card">
      <el-form 
        ref="filterFormRef"
        :model="filterForm" 
        :inline="true"
        @submit.prevent="handleSearch"
      >
        <slot name="filters" :form="filterForm" />
        
        <el-form-item>
          <el-button type="primary" @click="handleSearch">
            <el-icon><Search /></el-icon>
            {{ searchText }}
          </el-button>
          <el-button @click="handleReset">
            <el-icon><RefreshLeft /></el-icon>
            {{ resetText }}
          </el-button>
        </el-form-item>
      </el-form>
    </el-card>

    <!-- 操作栏 -->
    <el-card shadow="never" class="actions-card">
      <div class="actions-header">
        <!-- 左侧操作按钮 -->
        <div class="left-actions">
          <slot name="left-actions" />
          
          <!-- 批量操作 -->
          <template v-if="batchActions.length > 0 && selectedItems.length > 0">
            <el-dropdown 
              trigger="click" 
              @command="handleBatchAction"
              placement="bottom-start"
            >
              <el-button type="primary">
                批量操作 ({{ selectedItems.length }})
                <el-icon class="el-icon--right"><ArrowDown /></el-icon>
              </el-button>
              <template #dropdown>
                <el-dropdown-menu>
                  <el-dropdown-item 
                    v-for="action in batchActions" 
                    :key="action.key"
                    :command="action.key"
                    :disabled="action.disabled"
                  >
                    <el-icon v-if="action.icon">
                      <component :is="action.icon" />
                    </el-icon>
                    {{ action.label }}
                  </el-dropdown-item>
                </el-dropdown-menu>
              </template>
            </el-dropdown>
          </template>
        </div>

        <!-- 右侧操作 -->
        <div class="right-actions">
          <slot name="right-actions" />
          
          <!-- 表格设置 -->
          <el-tooltip v-if="showTableSettings" content="表格设置" placement="top">
            <el-button circle @click="showColumnSettings = true">
              <el-icon><Setting /></el-icon>
            </el-button>
          </el-tooltip>
          
          <!-- 刷新按钮 -->
          <el-tooltip v-if="showRefresh" content="刷新" placement="top">
            <el-button circle @click="handleRefresh">
              <el-icon><Refresh /></el-icon>
            </el-button>
          </el-tooltip>
        </div>
      </div>

      <!-- 统计信息 -->
      <div v-if="showStats" class="stats-info">
        <slot name="stats" :total="total" :selected="selectedItems.length" />
      </div>
    </el-card>

    <!-- 列设置对话框 -->
    <el-dialog
      v-model="showColumnSettings"
      title="列设置"
      width="400px"
    >
      <div class="column-settings">
        <el-checkbox-group v-model="visibleColumns">
          <div 
            v-for="column in columns" 
            :key="column.key"
            class="column-item"
          >
            <el-checkbox :label="column.key" :disabled="column.required">
              {{ column.label }}
            </el-checkbox>
          </div>
        </el-checkbox-group>
      </div>
      
      <template #footer>
        <el-button @click="showColumnSettings = false">取消</el-button>
        <el-button type="primary" @click="saveColumnSettings">确定</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, watch } from 'vue'
import { Search, RefreshLeft, ArrowDown, Setting, Refresh } from '@element-plus/icons-vue'
import type { FormInstance } from 'element-plus'

interface BatchAction {
  key: string
  label: string
  icon?: any
  disabled?: boolean
  danger?: boolean
}

interface Column {
  key: string
  label: string
  required?: boolean
}

interface Props {
  showFilters?: boolean
  showStats?: boolean
  showTableSettings?: boolean
  showRefresh?: boolean
  batchActions?: BatchAction[]
  selectedItems?: any[]
  total?: number
  columns?: Column[]
  searchText?: string
  resetText?: string
  initialFilters?: Record<string, any>
}

const props = withDefaults(defineProps<Props>(), {
  showFilters: true,
  showStats: true,
  showTableSettings: true,
  showRefresh: true,
  batchActions: () => [],
  selectedItems: () => [],
  total: 0,
  columns: () => [],
  searchText: '搜索',
  resetText: '重置',
  initialFilters: () => ({})
})

const emit = defineEmits<{
  search: [filters: Record<string, any>]
  reset: []
  refresh: []
  batchAction: [action: string, items: any[]]
  'update:visibleColumns': [columns: string[]]
}>()

const filterFormRef = ref<FormInstance>()
const filterForm = reactive({ ...props.initialFilters })
const showColumnSettings = ref(false)
const visibleColumns = ref<string[]>(
  props.columns.filter(col => col.required !== false).map(col => col.key)
)

// 监听初始过滤器变化
watch(() => props.initialFilters, (newFilters) => {
  Object.assign(filterForm, newFilters)
}, { deep: true })

const handleSearch = () => {
  emit('search', { ...filterForm })
}

const handleReset = () => {
  if (filterFormRef.value) {
    filterFormRef.value.resetFields()
  }
  Object.assign(filterForm, props.initialFilters)
  emit('reset')
}

const handleRefresh = () => {
  emit('refresh')
}

const handleBatchAction = (action: string) => {
  emit('batchAction', action, props.selectedItems)
}

const saveColumnSettings = () => {
  emit('update:visibleColumns', visibleColumns.value)
  showColumnSettings.value = false
}
</script>

<style scoped>
.table-actions {
  margin-bottom: 16px;
}

.filter-card,
.actions-card {
  border: none;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.1);
}

.filter-card {
  margin-bottom: 16px;
}

.actions-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 12px;
}

.left-actions,
.right-actions {
  display: flex;
  align-items: center;
  gap: 12px;
}

.stats-info {
  padding: 12px 0;
  border-top: 1px solid #f0f0f0;
  color: #666;
  font-size: 14px;
}

.column-settings {
  max-height: 300px;
  overflow-y: auto;
}

.column-item {
  margin-bottom: 12px;
  padding: 8px 0;
  border-bottom: 1px solid #f0f0f0;
}

.column-item:last-child {
  border-bottom: none;
}

:deep(.el-form--inline .el-form-item) {
  margin-right: 16px;
  margin-bottom: 12px;
}
</style> 