<template>
  <div class="logs-view">
    <div class="page-header">
      <h2 class="page-title">操作日志</h2>
      <div class="header-actions">
        <el-button @click="refreshData" :loading="loading">
          <el-icon><Refresh /></el-icon>
          刷新
        </el-button>
        <el-button @click="exportLogs">
          <el-icon><Download /></el-icon>
          导出
        </el-button>
      </div>
    </div>

    <!-- 筛选条件 -->
    <el-card class="filter-card" shadow="never">
      <el-form :model="filterForm" inline size="default">
        <el-form-item label="操作人">
          <el-input
            v-model="filterForm.operator_id"
            placeholder="请输入操作人ID"
            clearable
          />
        </el-form-item>
        
        <el-form-item label="操作类型">
          <el-select v-model="filterForm.action" placeholder="请选择操作类型" clearable>
            <el-option label="全部" value="" />
            <el-option label="登录" value="登录" />
            <el-option label="登出" value="登出" />
            <el-option label="审批申请" value="审批申请" />
            <el-option label="创建用户" value="创建用户" />
            <el-option label="更新用户状态" value="更新用户状态" />
            <el-option label="切换AI审批" value="切换AI审批" />
            <el-option label="查看申请详情" value="查看申请详情" />
          </el-select>
        </el-form-item>
        
        <el-form-item label="时间范围">
          <el-date-picker
            v-model="dateRange"
            type="datetimerange"
            range-separator="至"
            start-placeholder="开始时间"
            end-placeholder="结束时间"
            format="YYYY-MM-DD HH:mm"
            value-format="YYYY-MM-DD HH:mm:ss"
            @change="handleDateChange"
          />
        </el-form-item>
        
        <el-form-item>
          <el-button type="primary" @click="handleSearch">
            <el-icon><Search /></el-icon>
            搜索
          </el-button>
          <el-button @click="handleReset">
            <el-icon><RefreshLeft /></el-icon>
            重置
          </el-button>
        </el-form-item>
      </el-form>
    </el-card>

    <!-- 统计信息 -->
    <div class="stats-row">
      <el-card class="stat-item" shadow="hover">
        <div class="stat-content">
          <div class="stat-value">{{ statistics.total }}</div>
          <div class="stat-label">总操作数</div>
        </div>
      </el-card>
      <el-card class="stat-item" shadow="hover">
        <div class="stat-content">
          <div class="stat-value today">{{ statistics.today }}</div>
          <div class="stat-label">今日操作</div>
        </div>
      </el-card>
      <el-card class="stat-item" shadow="hover">
        <div class="stat-content">
          <div class="stat-value success">{{ statistics.success }}</div>
          <div class="stat-label">成功操作</div>
        </div>
      </el-card>
      <el-card class="stat-item" shadow="hover">
        <div class="stat-content">
          <div class="stat-value failed">{{ statistics.failed }}</div>
          <div class="stat-label">失败操作</div>
        </div>
      </el-card>
    </div>

    <!-- 日志列表 -->
    <el-card class="table-card" shadow="never">
      <template #header>
        <div class="card-header">
          <span>操作日志</span>
          <div class="header-extra">
            <el-tag type="info">
              共 {{ pagination.total }} 条记录
            </el-tag>
          </div>
        </div>
      </template>
      
      <el-table
        v-loading="loading"
        :data="logs"
        stripe
        style="width: 100%"
      >
        <el-table-column prop="id" label="ID" width="80" />
        
        <el-table-column prop="operator_name" label="操作人" width="120">
          <template #default="{ row }">
            <el-tag size="small" type="info">
              {{ row.operator_name }}
            </el-tag>
          </template>
        </el-table-column>
        
        <el-table-column prop="action" label="操作类型" width="120">
          <template #default="{ row }">
            <el-tag :type="getActionType(row.action)" size="small">
              {{ row.action }}
            </el-tag>
          </template>
        </el-table-column>
        
        <el-table-column prop="target" label="操作对象" min-width="200">
          <template #default="{ row }">
            <el-tooltip :content="row.target" placement="top">
              <span class="target-text">{{ row.target }}</span>
            </el-tooltip>
          </template>
        </el-table-column>
        
        <el-table-column prop="result" label="结果" width="100">
          <template #default="{ row }">
            <el-tag :type="getResultType(row.result)" size="small">
              {{ row.result }}
            </el-tag>
          </template>
        </el-table-column>
        
        <el-table-column prop="ip_address" label="IP地址" width="140">
          <template #default="{ row }">
            <code class="ip-address">{{ row.ip_address }}</code>
          </template>
        </el-table-column>
        
        <el-table-column prop="occurred_at" label="操作时间" width="160">
          <template #default="{ row }">
            {{ formatDateTime(row.occurred_at) }}
          </template>
        </el-table-column>
        
        <el-table-column label="操作" width="120" fixed="right">
          <template #default="{ row }">
            <el-button
              type="primary"
              size="small"
              @click="viewLogDetail(row)"
            >
              详情
            </el-button>
          </template>
        </el-table-column>
      </el-table>
      
      <!-- 分页 -->
      <div class="pagination-wrapper">
        <el-pagination
          v-model:current-page="pagination.page"
          v-model:page-size="pagination.limit"
          :page-sizes="[10, 20, 50, 100]"
          :total="pagination.total"
          layout="total, sizes, prev, pager, next, jumper"
          @size-change="handleSizeChange"
          @current-change="handleCurrentChange"
        />
      </div>
    </el-card>

    <!-- 日志详情对话框 -->
    <el-dialog
      v-model="detailDialogVisible"
      title="操作日志详情"
      width="600px"
    >
      <el-descriptions v-if="selectedLog" :column="1" border>
        <el-descriptions-item label="日志ID">
          {{ selectedLog.id }}
        </el-descriptions-item>
        <el-descriptions-item label="操作人ID">
          {{ selectedLog.operator_id }}
        </el-descriptions-item>
        <el-descriptions-item label="操作人名称">
          {{ selectedLog.operator_name }}
        </el-descriptions-item>
        <el-descriptions-item label="操作类型">
          <el-tag :type="getActionType(selectedLog.action)" size="small">
            {{ selectedLog.action }}
          </el-tag>
        </el-descriptions-item>
        <el-descriptions-item label="操作对象">
          {{ selectedLog.target }}
        </el-descriptions-item>
        <el-descriptions-item label="操作结果">
          <el-tag :type="getResultType(selectedLog.result)" size="small">
            {{ selectedLog.result }}
          </el-tag>
        </el-descriptions-item>
        <el-descriptions-item label="IP地址">
          <code>{{ selectedLog.ip_address }}</code>
        </el-descriptions-item>
        <el-descriptions-item label="用户代理">
          <div class="user-agent">{{ selectedLog.user_agent }}</div>
        </el-descriptions-item>
        <el-descriptions-item label="操作时间">
          {{ formatDateTime(selectedLog.occurred_at) }}
        </el-descriptions-item>
      </el-descriptions>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, computed, onMounted } from 'vue'
import { ElMessage } from 'element-plus'
import {
  Refresh,
  Download,
  Search,
  RefreshLeft
} from '@element-plus/icons-vue'
import { getOperationLogs } from '@/api/admin'
import type { OperationLog, PaginationResponse } from '@/types'
import dayjs from 'dayjs'

const loading = ref(false)
const detailDialogVisible = ref(false)
const selectedLog = ref<OperationLog | null>(null)

// 筛选表单
const filterForm = reactive({
  operator_id: '',
  action: '',
  start_date: '',
  end_date: ''
})

const dateRange = ref<[string, string] | null>(null)

// 分页信息
const pagination = reactive({
  page: 1,
  limit: 20,
  total: 0
})

// 日志列表
const logs = ref<OperationLog[]>([])

// 统计信息
const statistics = computed(() => {
  const total = logs.value.length
  const today = logs.value.filter(log => 
    dayjs(log.occurred_at).format('YYYY-MM-DD') === dayjs().format('YYYY-MM-DD')
  ).length
  const success = logs.value.filter(log => log.result === '成功').length
  const failed = logs.value.filter(log => log.result === '失败').length
  
  return { total, today, success, failed }
})

// 方法
const fetchLogs = async () => {
  try {
    loading.value = true
    const params = {
      operator_id: filterForm.operator_id,
      action: filterForm.action,
      start_date: filterForm.start_date,
      end_date: filterForm.end_date,
      page: pagination.page,
      limit: pagination.limit
    }
    
    const data: PaginationResponse<OperationLog> = await getOperationLogs(params)
    logs.value = data.data || []
    pagination.total = data.total || 0
  } catch (error) {
    ElMessage.error('获取操作日志失败')
  } finally {
    loading.value = false
  }
}

const refreshData = () => {
  fetchLogs()
}

const handleDateChange = (dates: [string, string] | null) => {
  if (dates) {
    filterForm.start_date = dates[0]
    filterForm.end_date = dates[1]
  } else {
    filterForm.start_date = ''
    filterForm.end_date = ''
  }
}

const handleSearch = () => {
  pagination.page = 1
  fetchLogs()
}

const handleReset = () => {
  Object.assign(filterForm, {
    operator_id: '',
    action: '',
    start_date: '',
    end_date: ''
  })
  dateRange.value = null
  pagination.page = 1
  fetchLogs()
}

const handleSizeChange = () => {
  pagination.page = 1
  fetchLogs()
}

const handleCurrentChange = () => {
  fetchLogs()
}

const viewLogDetail = (log: OperationLog) => {
  selectedLog.value = log
  detailDialogVisible.value = true
}

const exportLogs = () => {
  // 导出日志功能
  ElMessage.info('导出功能开发中...')
}

const getActionType = (action: string) => {
  const actionMap: Record<string, string> = {
    '登录': 'success',
    '登出': 'info',
    '审批申请': 'primary',
    '创建用户': 'warning',
    '更新用户状态': 'warning',
    '切换AI审批': 'danger',
    '查看申请详情': 'info'
  }
  return actionMap[action] || 'info'
}

const getResultType = (result: string) => {
  return result === '成功' ? 'success' : 'danger'
}

const formatDateTime = (datetime: string) => {
  return dayjs(datetime).format('YYYY-MM-DD HH:mm:ss')
}

onMounted(() => {
  fetchLogs()
})
</script>

<style scoped>
.logs-view {
  padding: 0;
}

.page-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 20px;
}

.page-title {
  margin: 0;
  color: #333;
  font-size: 24px;
  font-weight: 600;
}

.header-actions {
  display: flex;
  gap: 12px;
}

.filter-card {
  margin-bottom: 20px;
  border-radius: 8px;
}

.stats-row {
  display: grid;
  grid-template-columns: repeat(4, 1fr);
  gap: 16px;
  margin-bottom: 20px;
}

.stat-item {
  border-radius: 8px;
  text-align: center;
}

.stat-content {
  padding: 20px;
}

.stat-value {
  font-size: 28px;
  font-weight: 700;
  color: #333;
  margin-bottom: 8px;
}

.stat-value.today {
  color: #409eff;
}

.stat-value.success {
  color: #67c23a;
}

.stat-value.failed {
  color: #f56c6c;
}

.stat-label {
  color: #666;
  font-size: 14px;
}

.table-card {
  border-radius: 8px;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  font-weight: 600;
}

.target-text {
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
  text-overflow: ellipsis;
  cursor: help;
}

.ip-address {
  background: #f1f2f6;
  padding: 2px 6px;
  border-radius: 4px;
  font-size: 12px;
  color: #333;
}

.user-agent {
  word-break: break-all;
  line-height: 1.5;
  color: #666;
  font-size: 13px;
}

.pagination-wrapper {
  display: flex;
  justify-content: center;
  margin-top: 20px;
}

:deep(.el-descriptions__label) {
  font-weight: 500;
}

:deep(.el-date-editor) {
  width: 100%;
}
</style> 