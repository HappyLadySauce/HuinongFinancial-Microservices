<template>
  <div class="approval-dashboard">
    <!-- 头部统计卡片 -->
    <div class="stats-cards">
      <el-card class="stat-card" shadow="hover">
        <div class="stat-content">
          <div class="stat-icon pending">
            <el-icon><Clock /></el-icon>
          </div>
          <div class="stat-info">
            <div class="stat-number">{{ pendingCount }}</div>
            <div class="stat-label">待审批</div>
          </div>
        </div>
      </el-card>

      <el-card class="stat-card" shadow="hover">
        <div class="stat-content">
          <div class="stat-icon approved">
            <el-icon><Check /></el-icon>
          </div>
          <div class="stat-info">
            <div class="stat-number">{{ approvedCount }}</div>
            <div class="stat-label">已通过</div>
          </div>
        </div>
      </el-card>

      <el-card class="stat-card" shadow="hover">
        <div class="stat-content">
          <div class="stat-icon rejected">
            <el-icon><Close /></el-icon>
          </div>
          <div class="stat-info">
            <div class="stat-number">{{ rejectedCount }}</div>
            <div class="stat-label">已拒绝</div>
          </div>
        </div>
      </el-card>

      <el-card class="stat-card" shadow="hover">
        <div class="stat-content">
          <div class="stat-icon total">
            <el-icon><Document /></el-icon>
          </div>
          <div class="stat-info">
            <div class="stat-number">{{ totalCount }}</div>
            <div class="stat-label">总数量</div>
          </div>
        </div>
      </el-card>
    </div>

    <!-- 主要内容区域 -->
    <div class="main-content">
      <!-- 左侧：待办事项 -->
      <div class="left-panel">
        <el-card class="panel-card" shadow="never">
          <template #header>
            <div class="card-header">
              <span>待审批事项</span>
              <el-button type="primary" size="small" @click="handleRefresh">
                <el-icon><Refresh /></el-icon>
                刷新
              </el-button>
            </div>
          </template>

          <div class="todo-list">
            <div v-for="item in pendingApprovals" :key="item.id" class="todo-item">
              <div class="todo-info">
                <div class="todo-title">
                  <el-tag :type="item.type === 'loan' ? 'primary' : 'success'" size="small">
                    {{ item.type === 'loan' ? '贷款' : '租赁' }}
                  </el-tag>
                  <span class="title">{{ item.title }}</span>
                </div>
                <div class="todo-meta">
                  <span class="amount">{{ item.displayAmount || formatAmount(item.amount) }}</span>
                  <span class="time">{{ formatTime(item.createdAt) }}</span>
                </div>
              </div>
              <div class="todo-actions">
                <el-button type="primary" size="small" @click="handleApprove(item)">
                  审批
                </el-button>
                <el-button size="small" @click="handleViewDetail(item)">
                  详情
                </el-button>
              </div>
            </div>

            <div v-if="pendingApprovals.length === 0" class="empty-state">
              <el-empty description="暂无待审批事项" />
            </div>
          </div>
        </el-card>
      </div>

      <!-- 右侧：图表统计 -->
      <div class="right-panel">
        <el-card class="panel-card" shadow="never">
          <template #header>
            <div class="card-header">
              <span>审批统计</span>
              <el-radio-group v-model="chartPeriod" size="small">
                <el-radio-button value="week">本周</el-radio-button>
                <el-radio-button value="month">本月</el-radio-button>
                <el-radio-button value="quarter">本季度</el-radio-button>
              </el-radio-group>
            </div>
          </template>

          <div class="chart-container">
            <div ref="approvalChart" class="chart"></div>
          </div>
        </el-card>
      </div>
    </div>

    <!-- 底部：最近活动 -->
    <el-card class="activity-card" shadow="never">
      <template #header>
        <div class="card-header">
          <span>最近活动</span>
          <el-button type="text" size="small" @click="viewAllActivities">
            查看全部
          </el-button>
        </div>
      </template>

      <div class="activity-list">
        <div v-for="activity in recentActivities" :key="activity.id" class="activity-item">
          <div class="activity-icon">
            <el-icon v-if="activity.action === 'approve'" class="approve"><Check /></el-icon>
            <el-icon v-else-if="activity.action === 'reject'" class="reject"><Close /></el-icon>
            <el-icon v-else class="default"><Document /></el-icon>
          </div>
          <div class="activity-content">
            <div class="activity-title">{{ activity.title }}</div>
            <div class="activity-meta">
              <span>{{ activity.user }}</span>
              <span class="time">{{ formatTime(activity.createdAt) }}</span>
            </div>
          </div>
          <div class="activity-status">
            <el-tag 
              :type="getStatusType(activity.status)" 
              size="small"
            >
              {{ activity.status }}
            </el-tag>
          </div>
        </div>
      </div>
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, nextTick } from 'vue'
import { useRouter } from 'vue-router'
import { Clock, Check, Close, Document, Refresh } from '@element-plus/icons-vue'
import * as echarts from 'echarts'
import { ElMessage } from 'element-plus'
import { adminApi } from '@/services/api'

const router = useRouter()

// 响应式数据
const loading = ref(false)
const pendingCount = ref(0)
const approvedCount = ref(0)
const rejectedCount = ref(0)
const totalCount = ref(0)
const chartPeriod = ref('month')

// 图表实例
const approvalChart = ref<HTMLElement>()
let chartInstance: echarts.ECharts | null = null

// 待审批事项
const pendingApprovals = ref<any[]>([])

// 最近活动
const recentActivities = ref<any[]>([])

// 加载统计数据
const loadStats = async () => {
  try {
    loading.value = true
    
    // 分别获取贷款和租赁的统计数据
    const [loanPending, loanApproved, loanRejected, leasePending, leaseApproved, leaseRejected] = await Promise.all([
      adminApi.getAllLoanApprovals({ status: 'pending', page_size: 1 }),
      adminApi.getAllLoanApprovals({ status: 'approved', page_size: 1 }),
      adminApi.getAllLoanApprovals({ status: 'rejected', page_size: 1 }),
      adminApi.getAllLeaseApprovals({ status: 'pending', page_size: 1 }),
      adminApi.getAllLeaseApprovals({ status: 'approved', page_size: 1 }),
      adminApi.getAllLeaseApprovals({ status: 'rejected', page_size: 1 })
    ])
    
    pendingCount.value = loanPending.total + leasePending.total
    approvedCount.value = loanApproved.total + leaseApproved.total  
    rejectedCount.value = loanRejected.total + leaseRejected.total
    totalCount.value = pendingCount.value + approvedCount.value + rejectedCount.value
  } catch (error) {
    console.error('加载统计数据失败:', error)
    ElMessage.error('加载统计数据失败')
  } finally {
    loading.value = false
  }
}

// 加载待审批事项
const loadPendingApprovals = async () => {
  try {
    // 获取待审批的贷款和租赁申请
    const [loanPending, leasePending] = await Promise.all([
      adminApi.getAllLoanApprovals({ status: 'pending', page_size: 3 }),
      adminApi.getAllLeaseApprovals({ status: 'pending', page_size: 2 })
    ])
    
    // 合并并格式化数据
    const allPending = [
      ...loanPending.list.map(item => ({ 
        ...item, 
        type: 'loan',
        title: item.name || '贷款申请',
        createdAt: item.created_at,
        amount: item.amount || 0,
        displayAmount: `¥${((item.amount || 0) / 10000).toFixed(1)}万`
      })),
      ...leasePending.list.map(item => ({ 
        ...item, 
        type: 'lease',
        title: item.name || '租赁申请',
        createdAt: item.created_at,
        amount: 0, // 租赁没有金额字段
        displayAmount: `${item.start_at || ''} 至 ${item.end_at || ''}`,
        startAt: item.start_at,
        endAt: item.end_at
      }))
    ]
    
    // 按创建时间排序，取最新的5个
    pendingApprovals.value = allPending
      .filter(item => item.createdAt) // 过滤掉没有时间的记录
      .sort((a, b) => new Date(b.createdAt).getTime() - new Date(a.createdAt).getTime())
      .slice(0, 5)
  } catch (error) {
    console.error('加载待审批事项失败:', error)
    ElMessage.error('加载待审批事项失败')
  }
}

// 加载最近活动
const loadRecentActivities = async () => {
  try {
    // 获取最近的已审批申请
    const [loanApproved, loanRejected, leaseApproved, leaseRejected] = await Promise.all([
      adminApi.getAllLoanApprovals({ status: 'approved', page_size: 3 }),
      adminApi.getAllLoanApprovals({ status: 'rejected', page_size: 2 }),
      adminApi.getAllLeaseApprovals({ status: 'approved', page_size: 3 }),
      adminApi.getAllLeaseApprovals({ status: 'rejected', page_size: 2 })
    ])
    
    // 合并并格式化数据
    const allActivities = [
      ...loanApproved.list.map(item => ({ 
        ...item, 
        type: 'loan',
        title: `贷款申请 - ${item.name}`,
        user: item.auditor || '系统',
        action: 'approve',
        createdAt: item.updated_at || item.created_at,
        amount: item.amount || 0
      })),
      ...loanRejected.list.map(item => ({ 
        ...item, 
        type: 'loan',
        title: `贷款申请 - ${item.name}`,
        user: item.auditor || '系统',
        action: 'reject',
        createdAt: item.updated_at || item.created_at,
        amount: item.amount || 0
      })),
      ...leaseApproved.list.map(item => ({ 
        ...item, 
        type: 'lease',
        title: `租赁申请 - ${item.name}`,
        user: item.auditor || '系统',
        action: 'approve',
        createdAt: item.updated_at || item.created_at,
        amount: 0, // 租赁没有金额字段
        startAt: item.start_at,
        endAt: item.end_at
      })),
      ...leaseRejected.list.map(item => ({ 
        ...item, 
        type: 'lease',
        title: `租赁申请 - ${item.name}`,
        user: item.auditor || '系统',
        action: 'reject',
        createdAt: item.updated_at || item.created_at,
        amount: 0, // 租赁没有金额字段
        startAt: item.start_at,
        endAt: item.end_at
      }))
    ]
    
    // 按更新时间排序，取最新的10个
    recentActivities.value = allActivities
      .filter(item => item.createdAt) // 过滤掉没有时间的记录
      .sort((a, b) => new Date(b.createdAt).getTime() - new Date(a.createdAt).getTime())
      .slice(0, 10)
  } catch (error) {
    console.error('加载最近活动失败:', error)
    ElMessage.error('加载最近活动失败')
  }
}

// 格式化金额
const formatAmount = (amount: number) => {
  if (amount === 0) return '租赁申请'
  return `¥${(amount / 10000).toFixed(1)}万`
}

// 格式化时间
const formatTime = (dateStr: string | Date | null | undefined) => {
  if (!dateStr) return '未知时间'
  
  const date = typeof dateStr === 'string' ? new Date(dateStr) : dateStr
  
  // 检查日期是否有效
  if (isNaN(date.getTime())) return '无效时间'
  
  const now = new Date()
  const diff = now.getTime() - date.getTime()
  
  if (diff < 0) return '刚刚'
  if (diff < 60 * 60 * 1000) {
    return `${Math.floor(diff / (60 * 1000))}分钟前`
  } else if (diff < 24 * 60 * 60 * 1000) {
    return `${Math.floor(diff / (60 * 60 * 1000))}小时前`
  } else {
    return `${Math.floor(diff / (24 * 60 * 60 * 1000))}天前`
  }
}

// 获取状态类型
const getStatusType = (status: string) => {
  switch (status) {
    case 'approved': return 'success'
    case 'rejected': return 'danger'
    case 'pending': return 'warning'
    default: return 'info'
  }
}

// 处理审批
const handleApprove = (item: any) => {
  const route = item.type === 'loan' ? `/approval/loan/${item.id}` : `/approval/lease/${item.id}`
  router.push(route)
}

// 查看详情
const handleViewDetail = (item: any) => {
  const route = item.type === 'loan' ? `/approval/loan/${item.id}` : `/approval/lease/${item.id}`
  router.push(route)
}

// 刷新数据
const handleRefresh = async () => {
  await Promise.all([
    loadStats(),
    loadPendingApprovals(),
    loadRecentActivities()
  ])
  ElMessage.success('数据已刷新')
}

// 查看全部活动
const viewAllActivities = () => {
  router.push('/approval/activities')
}

// 初始化图表
const initChart = () => {
  if (!approvalChart.value) return
  
  chartInstance = echarts.init(approvalChart.value)
  
  const option = {
    title: {
      text: '审批趋势',
      textStyle: {
        fontSize: 14,
        fontWeight: 'normal'
      }
    },
    tooltip: {
      trigger: 'axis'
    },
    legend: {
      data: ['贷款审批', '租赁审批'],
      bottom: 0
    },
    xAxis: {
      type: 'category',
      data: ['周一', '周二', '周三', '周四', '周五', '周六', '周日']
    },
    yAxis: {
      type: 'value'
    },
    series: [
      {
        name: '贷款审批',
        type: 'line',
        data: [12, 19, 15, 23, 18, 25, 20],
        smooth: true,
        itemStyle: {
          color: '#409EFF'
        }
      },
      {
        name: '租赁审批',
        type: 'line',
        data: [8, 12, 10, 15, 12, 18, 14],
        smooth: true,
        itemStyle: {
          color: '#67C23A'
        }
      }
    ]
  }
  
  chartInstance.setOption(option)
}

onMounted(() => {
  // 加载数据
  Promise.all([
    loadStats(),
    loadPendingApprovals(),
    loadRecentActivities()
  ])
  
  // 初始化图表
  nextTick(() => {
    initChart()
  })
})
</script>

<style scoped>
.approval-dashboard {
  padding: 20px;
  background: #f5f5f5;
  min-height: calc(100vh - 60px);
}

.stats-cards {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(200px, 1fr));
  gap: 16px;
  margin-bottom: 20px;
}

.stat-card {
  border: none;
}

.stat-content {
  display: flex;
  align-items: center;
  gap: 16px;
}

.stat-icon {
  width: 48px;
  height: 48px;
  border-radius: 8px;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 24px;
  color: white;
}

.stat-icon.pending { background: #E6A23C; }
.stat-icon.approved { background: #67C23A; }
.stat-icon.rejected { background: #F56C6C; }
.stat-icon.total { background: #409EFF; }

.stat-info {
  flex: 1;
}

.stat-number {
  font-size: 24px;
  font-weight: bold;
  color: #303133;
  margin-bottom: 4px;
}

.stat-label {
  font-size: 14px;
  color: #909399;
}

.main-content {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 20px;
  margin-bottom: 20px;
}

.panel-card {
  border: none;
  height: 500px;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  font-weight: 500;
}

.todo-list {
  max-height: 400px;
  overflow-y: auto;
}

.todo-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 12px 0;
  border-bottom: 1px solid #EBEEF5;
}

.todo-item:last-child {
  border-bottom: none;
}

.todo-info {
  flex: 1;
}

.todo-title {
  display: flex;
  align-items: center;
  gap: 8px;
  margin-bottom: 4px;
}

.todo-title .title {
  font-weight: 500;
  color: #303133;
}

.todo-meta {
  display: flex;
  gap: 16px;
  font-size: 12px;
  color: #909399;
}

.todo-meta .amount {
  font-weight: 500;
  color: #E6A23C;
}

.todo-actions {
  display: flex;
  gap: 8px;
}

.chart-container {
  height: 400px;
  padding: 20px 0;
}

.chart {
  width: 100%;
  height: 100%;
}

.activity-card {
  border: none;
}

.activity-list {
  max-height: 300px;
  overflow-y: auto;
}

.activity-item {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 12px 0;
  border-bottom: 1px solid #EBEEF5;
}

.activity-item:last-child {
  border-bottom: none;
}

.activity-icon {
  width: 32px;
  height: 32px;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 16px;
  color: white;
}

.activity-icon.approve { background: #67C23A; }
.activity-icon.reject { background: #F56C6C; }
.activity-icon.default { background: #909399; }

.activity-content {
  flex: 1;
}

.activity-title {
  font-weight: 500;
  color: #303133;
  margin-bottom: 4px;
}

.activity-meta {
  display: flex;
  gap: 16px;
  font-size: 12px;
  color: #909399;
}

.empty-state {
  text-align: center;
  padding: 40px 0;
}

@media (max-width: 768px) {
  .stats-cards {
    grid-template-columns: repeat(2, 1fr);
  }
  
  .main-content {
    grid-template-columns: 1fr;
  }
  
  .todo-item {
    flex-direction: column;
    align-items: flex-start;
    gap: 8px;
  }
  
  .todo-actions {
    width: 100%;
    justify-content: flex-end;
  }
}
</style> 