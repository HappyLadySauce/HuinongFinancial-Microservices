<template>
  <div class="logs-view">
    <!-- 页面头部 -->
    <PageHeader
      title="操作日志"
      subtitle="系统操作记录和审计日志"
      :icon="Document"
      :breadcrumbs="breadcrumbs"
    />

    <!-- 日志查询 -->
    <TableActions
      @search="handleSearch"
      @refresh="handleRefresh"
      :search-placeholder="'搜索用户名、操作内容...'"
    >
      <template #filters>
        <el-form-item label="日志类型">
          <el-select v-model="filters.type" placeholder="全部类型" clearable>
            <el-option label="登录日志" value="login" />
            <el-option label="操作日志" value="operation" />
            <el-option label="审批日志" value="approval" />
            <el-option label="系统日志" value="system" />
            <el-option label="错误日志" value="error" />
          </el-select>
        </el-form-item>

        <el-form-item label="操作结果">
          <el-select v-model="filters.result" placeholder="全部结果" clearable>
            <el-option label="成功" value="success" />
            <el-option label="失败" value="failure" />
            <el-option label="警告" value="warning" />
          </el-select>
        </el-form-item>

        <el-form-item label="时间范围">
          <el-date-picker
            v-model="filters.dateRange"
            type="datetimerange"
            range-separator="至"
            start-placeholder="开始时间"
            end-placeholder="结束时间"
            format="YYYY-MM-DD HH:mm:ss"
            value-format="YYYY-MM-DD HH:mm:ss"
          />
        </el-form-item>
      </template>

      <template #stats>
        <div class="log-stats">
          <StatCard
            title="今日日志"
            :value="stats.todayLogs"
            format="number"
            type="primary"
            size="small"
          />
          <StatCard
            title="异常日志"
            :value="stats.errorLogs"
            format="number"
            type="danger"
            size="small"
          />
          <StatCard
            title="活跃用户"
            :value="stats.activeUsers"
            format="number"
            type="success"
            size="small"
          />
        </div>
      </template>

      <template #actions>
        <el-button @click="handleExport">
          <el-icon><Download /></el-icon>
          导出日志
        </el-button>
        <el-button @click="handleCleanup" type="danger">
          <el-icon><Delete /></el-icon>
          清理日志
        </el-button>
      </template>
    </TableActions>

    <!-- 日志表格 -->
    <el-card class="table-card" shadow="never">
      <el-table 
        :data="logs" 
        v-loading="loading"
        border
        stripe
        @sort-change="handleSortChange"
      >
        <el-table-column prop="id" label="ID" width="80" />
        
        <el-table-column prop="type" label="类型" width="100">
          <template #default="{ row }">
            <StatusTag :status="row.type" />
          </template>
        </el-table-column>

        <el-table-column prop="username" label="用户" width="120">
          <template #default="{ row }">
            <div class="user-info">
              <el-avatar :size="24" :src="row.userAvatar">
                {{ row.username?.charAt(0) }}
              </el-avatar>
              <span>{{ row.username }}</span>
            </div>
          </template>
        </el-table-column>

        <el-table-column prop="operation" label="操作内容" min-width="200">
          <template #default="{ row }">
            <div class="operation-content">
              <div class="action">{{ row.operation }}</div>
              <div class="target" v-if="row.target">目标: {{ row.target }}</div>
            </div>
          </template>
        </el-table-column>

        <el-table-column prop="result" label="结果" width="100">
          <template #default="{ row }">
            <StatusTag :status="row.result" :show-icon="true" />
          </template>
        </el-table-column>

        <el-table-column prop="ip" label="IP地址" width="120" />

        <el-table-column prop="userAgent" label="设备信息" width="150" show-overflow-tooltip>
          <template #default="{ row }">
            <div class="device-info">
              <el-icon v-if="row.userAgent.includes('Mobile')"><Iphone /></el-icon>
              <el-icon v-else><Monitor /></el-icon>
              <span>{{ getDeviceInfo(row.userAgent) }}</span>
            </div>
          </template>
        </el-table-column>

        <el-table-column prop="createdAt" label="操作时间" width="160" sortable="custom">
          <template #default="{ row }">
            {{ formatDateTime(row.createdAt) }}
          </template>
        </el-table-column>

        <el-table-column label="操作" width="120" fixed="right">
          <template #default="{ row }">
            <el-button 
              size="small" 
              @click="handleViewDetail(row)"
              :disabled="!row.details"
            >
              <el-icon><View /></el-icon>
              详情
            </el-button>
          </template>
        </el-table-column>
      </el-table>

      <!-- 分页 -->
      <div class="pagination-wrapper">
        <el-pagination
          v-model:current-page="pagination.page"
          v-model:page-size="pagination.size"
          :total="pagination.total"
          :page-sizes="[10, 20, 50, 100]"
          layout="total, sizes, prev, pager, next, jumper"
          @size-change="loadLogs"
          @current-change="loadLogs"
        />
      </div>
    </el-card>

    <!-- 详情对话框 -->
    <el-dialog v-model="detailDialog.visible" title="日志详情" width="800px">
      <div v-if="detailDialog.data" class="log-detail">
        <el-descriptions :column="2" border>
          <el-descriptions-item label="日志ID">
            {{ detailDialog.data.id }}
          </el-descriptions-item>
          <el-descriptions-item label="操作类型">
            <StatusTag :status="detailDialog.data.type" />
          </el-descriptions-item>
          <el-descriptions-item label="操作用户">
            {{ detailDialog.data.username }}
          </el-descriptions-item>
          <el-descriptions-item label="操作时间">
            {{ formatDateTime(detailDialog.data.createdAt) }}
          </el-descriptions-item>
          <el-descriptions-item label="IP地址">
            {{ detailDialog.data.ip }}
          </el-descriptions-item>
          <el-descriptions-item label="操作结果">
            <StatusTag :status="detailDialog.data.result" :show-icon="true" />
          </el-descriptions-item>
          <el-descriptions-item label="操作内容" :span="2">
            {{ detailDialog.data.operation }}
          </el-descriptions-item>
          <el-descriptions-item label="User Agent" :span="2">
            {{ detailDialog.data.userAgent }}
          </el-descriptions-item>
        </el-descriptions>

        <div v-if="detailDialog.data.details" class="detail-content">
          <h4>详细信息</h4>
          <el-input
            v-model="detailDialog.data.details"
            type="textarea"
            :rows="8"
            readonly
          />
        </div>

        <div v-if="detailDialog.data.error" class="error-content">
          <h4>错误信息</h4>
          <pre class="error-trace">{{ detailDialog.data.error }}</pre>
        </div>
      </div>

      <template #footer>
        <el-button @click="detailDialog.visible = false">关闭</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, computed, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import dayjs from 'dayjs'
import {
  Document,
  Download,
  Delete,
  View,
  Iphone,
  Monitor
} from '@element-plus/icons-vue'

// 导入通用组件
import { PageHeader, TableActions, StatCard, StatusTag } from '@/components/common'

const loading = ref(false)

const filters = reactive({
  keyword: '',
  type: '',
  result: '',
  dateRange: null
})

const pagination = reactive({
  page: 1,
  size: 20,
  total: 0
})

const stats = ref({
  todayLogs: 1520,
  errorLogs: 12,
  activeUsers: 89
})

const detailDialog = reactive({
  visible: false,
  data: null
})

// 日志数据
const logs = ref([
  {
    id: 'LOG001',
    type: 'login',
    username: 'admin',
    userAvatar: '',
    operation: '用户登录',
    target: null,
    result: 'success',
    ip: '192.168.1.100',
    userAgent: 'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36',
    createdAt: '2024-12-29 18:15:32',
    details: '用户从管理端登录系统，登录成功'
  },
  {
    id: 'LOG002',
    type: 'approval',
    username: 'reviewer1',
    userAvatar: '',
    operation: '审批申请',
    target: '贷款申请 #12345',
    result: 'success',
    ip: '192.168.1.105',
    userAgent: 'Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36',
    createdAt: '2024-12-29 17:45:18',
    details: '审批通过农户贷款申请，金额: 50000元'
  },
  {
    id: 'LOG003',
    type: 'operation',
    username: 'admin',
    userAvatar: '',
    operation: '修改用户权限',
    target: '用户 reviewer2',
    result: 'success',
    ip: '192.168.1.100',
    userAgent: 'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36',
    createdAt: '2024-12-29 16:30:45',
    details: '将用户 reviewer2 权限从普通用户升级为审批员'
  },
  {
    id: 'LOG004',
    type: 'error',
    username: 'system',
    userAvatar: '',
    operation: '系统错误',
    target: '数据库连接',
    result: 'failure',
    ip: '127.0.0.1',
    userAgent: 'System',
    createdAt: '2024-12-29 15:22:10',
    details: '数据库连接超时',
    error: 'Connection timeout after 30 seconds\nat DatabaseConnection.connect()\nat UserService.findById()'
  },
  {
    id: 'LOG005',
    type: 'system',
    username: 'system',
    userAvatar: '',
    operation: '系统备份',
    target: '数据库备份',
    result: 'success',
    ip: '127.0.0.1',
    userAgent: 'System',
    createdAt: '2024-12-29 02:00:00',
    details: '定时备份任务执行成功，备份文件大小: 2.5GB'
  }
])

// 计算属性
const breadcrumbs = computed(() => [
  { title: '首页', to: '/' },
  { title: '系统管理', to: '/system' },
  { title: '操作日志' }
])

// 方法
const handleSearch = (keyword: string) => {
  filters.keyword = keyword
  pagination.page = 1
  loadLogs()
}

const handleRefresh = () => {
  loadLogs()
  ElMessage.success('日志数据已刷新')
}

const handleSortChange = ({ prop, order }: any) => {
  console.log('排序变化:', prop, order)
  // 实现排序逻辑
  loadLogs()
}

const handleViewDetail = (log: any) => {
  detailDialog.data = log
  detailDialog.visible = true
}

const handleExport = async () => {
  try {
    await ElMessageBox.confirm('确定要导出当前筛选条件下的日志吗？', '确认导出', {
      type: 'info'
    })
    
    console.log('导出日志:', filters)
    ElMessage.success('日志导出功能待实现')
  } catch {
    // 用户取消
  }
}

const handleCleanup = async () => {
  try {
    await ElMessageBox.confirm(
      '清理操作将删除30天前的日志，此操作不可恢复，确定要继续吗？',
      '确认清理',
      {
        type: 'warning',
        confirmButtonText: '确定清理',
        cancelButtonText: '取消'
      }
    )
    
    console.log('清理日志')
    ElMessage.success('日志清理功能待实现')
  } catch {
    // 用户取消
  }
}

const loadLogs = () => {
  loading.value = true
  
  console.log('加载日志:', {
    filters,
    pagination: {
      page: pagination.page,
      size: pagination.size
    }
  })
  
  // 模拟加载
  setTimeout(() => {
    pagination.total = 1250 // 模拟总数
    loading.value = false
  }, 500)
}

const getDeviceInfo = (userAgent: string) => {
  if (userAgent.includes('Windows')) return 'Windows'
  if (userAgent.includes('Macintosh')) return 'macOS'
  if (userAgent.includes('Linux')) return 'Linux'
  if (userAgent.includes('Mobile')) return 'Mobile'
  return '未知设备'
}

const formatDateTime = (datetime: string) => {
  return dayjs(datetime).format('YYYY-MM-DD HH:mm:ss')
}

onMounted(() => {
  loadLogs()
})
</script>

<style scoped>
.logs-view {
  height: 100vh;
  display: flex;
  flex-direction: column;
  gap: 16px;
  padding: 0 20px;
  overflow: hidden;
}

.log-stats {
  display: flex;
  gap: 16px;
}

.table-card {
  flex: 1;
  border: none;
  border-radius: 12px;
  overflow: hidden;
}

.user-info {
  display: flex;
  align-items: center;
  gap: 8px;
}

.operation-content .action {
  font-weight: 500;
  margin-bottom: 4px;
}

.operation-content .target {
  font-size: 12px;
  color: #666;
}

.device-info {
  display: flex;
  align-items: center;
  gap: 4px;
}

.pagination-wrapper {
  display: flex;
  justify-content: center;
  margin-top: 20px;
}

.log-detail {
  margin-bottom: 20px;
}

.detail-content,
.error-content {
  margin-top: 20px;
}

.detail-content h4,
.error-content h4 {
  margin-bottom: 12px;
  color: #333;
  font-weight: 500;
}

.error-trace {
  background: #f5f5f5;
  padding: 12px;
  border-radius: 4px;
  font-family: 'Courier New', monospace;
  font-size: 12px;
  color: #e74c3c;
  overflow-x: auto;
  white-space: pre-wrap;
}

:deep(.el-card__body) {
  padding: 0;
}

:deep(.el-table) {
  border-radius: 0;
}

:deep(.el-table__header) {
  background-color: #fafafa;
}
</style>
