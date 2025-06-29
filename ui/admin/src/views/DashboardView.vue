<template>
  <div class="dashboard">
    <!-- 欢迎区域 -->
    <div class="welcome-section">
      <el-card class="welcome-card">
        <div class="welcome-content">
          <div class="welcome-info">
            <h2 class="welcome-title">{{ greeting }}，{{ currentUser?.username || '管理员' }}！</h2>
            <p class="welcome-desc">今天是 {{ currentDate }}，{{ weekDay }}</p>
            <div class="quick-stats">
              <div class="stat-item">
                <span class="stat-label">今日审批</span>
                <span class="stat-value">{{ todayStats.approvals }}</span>
              </div>
              <div class="stat-item">
                <span class="stat-label">待处理</span>
                <span class="stat-value">{{ todayStats.pending }}</span>
              </div>
            </div>
          </div>
          <div class="welcome-actions">
            <el-button type="primary" size="large" @click="goToApproval">
              <el-icon><DocumentChecked /></el-icon>
              进入审批
            </el-button>
            <el-button size="large" @click="goToStatistics">
              <el-icon><DataLine /></el-icon>
              查看报表
            </el-button>
          </div>
        </div>
      </el-card>
    </div>

    <!-- 统计卡片 -->
    <div class="stats-section">
      <el-row :gutter="24">
        <el-col :xs="24" :sm="12" :md="6" v-for="(stat, index) in statsCards" :key="index">
          <el-card class="stat-card" :class="`stat-card-${stat.type}`">
            <div class="stat-card-content">
              <div class="stat-icon">
                <el-icon :class="stat.iconClass">
                  <component :is="stat.icon" />
                </el-icon>
              </div>
              <div class="stat-info">
                <div class="stat-number">{{ stat.value }}</div>
                <div class="stat-title">{{ stat.title }}</div>
                <div class="stat-change" :class="stat.trend">
                  <el-icon>
                    <ArrowUp v-if="stat.trend === 'up'" />
                    <ArrowDown v-if="stat.trend === 'down'" />
                    <Minus v-if="stat.trend === 'same'" />
                  </el-icon>
                  {{ stat.change }}
                </div>
              </div>
            </div>
          </el-card>
        </el-col>
      </el-row>
    </div>

    <!-- 图表和待办事项 -->
    <div class="charts-section">
      <el-row :gutter="24">
        <!-- 审批趋势图 -->
        <el-col :xs="24" :lg="16">
          <el-card class="chart-card" header="审批趋势统计">
            <template #header>
              <div class="card-header">
                <span>审批趋势统计</span>
                <div class="header-actions">
                  <el-radio-group v-model="chartPeriod" size="small">
                    <el-radio-button label="week">近7天</el-radio-button>
                    <el-radio-button label="month">近30天</el-radio-button>
                    <el-radio-button label="year">本年</el-radio-button>
                  </el-radio-group>
                </div>
              </div>
            </template>
            <div class="chart-container">
              <div ref="trendChart" style="width: 100%; height: 300px;"></div>
            </div>
          </el-card>
        </el-col>

        <!-- 待办事项 -->
        <el-col :xs="24" :lg="8">
          <el-card class="todo-card" header="待办事项">
            <template #header>
              <div class="card-header">
                <span>待办事项</span>
                <el-badge :value="todoList.length" :max="99">
                  <el-button link size="small" @click="refreshTodos">
                    <el-icon><Refresh /></el-icon>
                  </el-button>
                </el-badge>
              </div>
            </template>
            <div class="todo-list">
              <div 
                v-for="item in todoList" 
                :key="item.id" 
                class="todo-item"
                @click="handleTodoClick(item)"
              >
                <div class="todo-info">
                  <div class="todo-title">{{ item.title }}</div>
                  <div class="todo-desc">{{ item.description }}</div>
                  <div class="todo-time">{{ formatTime(item.createdAt) }}</div>
                </div>
                <div class="todo-action">
                  <el-tag :type="getTagType(item.type)" size="small">
                    {{ item.type }}
                  </el-tag>
                </div>
              </div>
              <div v-if="todoList.length === 0" class="todo-empty">
                <el-empty description="暂无待办事项" />
              </div>
            </div>
          </el-card>
        </el-col>
      </el-row>
    </div>

    <!-- 审批类型分布和最近活动 -->
    <div class="bottom-section">
      <el-row :gutter="24">
        <!-- 审批类型分布 -->
        <el-col :xs="24" :lg="12">
          <el-card class="chart-card" header="审批类型分布">
            <div class="chart-container">
              <div ref="pieChart" style="width: 100%; height: 300px;"></div>
            </div>
          </el-card>
        </el-col>

        <!-- 最近活动 -->
        <el-col :xs="24" :lg="12">
          <el-card class="activity-card" header="最近活动">
            <div class="activity-list">
              <div 
                v-for="activity in recentActivities" 
                :key="activity.id" 
                class="activity-item"
              >
                <div class="activity-avatar">
                  <el-avatar :size="32">
                    <el-icon><User /></el-icon>
                  </el-avatar>
                </div>
                <div class="activity-content">
                  <div class="activity-text">
                    <span class="activity-user">{{ activity.user }}</span>
                    <span class="activity-action">{{ activity.action }}</span>
                    <span class="activity-target">{{ activity.target }}</span>
                  </div>
                  <div class="activity-time">{{ formatTime(activity.createdAt) }}</div>
                </div>
              </div>
              <div v-if="recentActivities.length === 0" class="activity-empty">
                <el-empty description="暂无活动记录" />
              </div>
            </div>
          </el-card>
        </el-col>
      </el-row>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, computed, nextTick } from 'vue'
import { useRouter } from 'vue-router'
import { useUserStore } from '@/stores/user'
import { 
  DocumentChecked, DataLine, ArrowUp, ArrowDown, Minus, Refresh, User,
  CreditCard, Van, TrendCharts, Clock
} from '@element-plus/icons-vue'
import dayjs from 'dayjs'
import relativeTime from 'dayjs/plugin/relativeTime'
import 'dayjs/locale/zh-cn'

dayjs.extend(relativeTime)
dayjs.locale('zh-cn')
import * as echarts from 'echarts'

const router = useRouter()
const userStore = useUserStore()

// 响应式数据
const chartPeriod = ref('month')
const trendChart = ref<HTMLElement>()
const pieChart = ref<HTMLElement>()
const todoList = ref([])
const recentActivities = ref([])
const todayStats = ref({
  approvals: 0,
  pending: 0
})

// 计算属性
const currentUser = computed(() => userStore.userInfo)

const currentDate = computed(() => {
  return new Date().toLocaleDateString('zh-CN', {
    year: 'numeric',
    month: 'long',
    day: 'numeric'
  })
})

const weekDay = computed(() => {
  const days = ['周日', '周一', '周二', '周三', '周四', '周五', '周六']
  return days[new Date().getDay()]
})

const greeting = computed(() => {
  const hour = new Date().getHours()
  if (hour < 6) return '凌晨好'
  if (hour < 9) return '早上好'
  if (hour < 12) return '上午好'
  if (hour < 14) return '中午好'
  if (hour < 17) return '下午好'
  if (hour < 19) return '傍晚好'
  return '晚上好'
})

const statsCards = ref([
  {
    title: '贷款审批',
    value: '114',
    change: '+12.5%',
    trend: 'up',
    type: 'loan',
    icon: CreditCard,
    iconClass: 'loan-icon'
  },
  {
    title: '租赁审批',
    value: '64',
    change: '+8.2%',
    trend: 'up',
    type: 'lease',
    icon: Van,
    iconClass: 'lease-icon'
  },
  {
    title: '待处理',
    value: '23',
    change: '-5.1%',
    trend: 'down',
    type: 'pending',
    icon: Clock,
    iconClass: 'pending-icon'
  },
  {
    title: '本月完成',
    value: '342',
    change: '+15.3%',
    trend: 'up',
    type: 'completed',
    icon: TrendCharts,
    iconClass: 'completed-icon'
  }
])

// 方法
const goToApproval = () => {
  router.push('/approval/dashboard')
}

const goToStatistics = () => {
  router.push('/operation/statistics')
}

const formatTime = (time: string | Date) => {
  return dayjs(time).fromNow()
}

const getTagType = (type: string) => {
  const typeMap = {
    '贷款审批': 'primary',
    '租赁审批': 'success',
    '系统通知': 'info'
  }
  return typeMap[type] || 'default'
}

const handleTodoClick = (item: any) => {
  if (item.link) {
    router.push(item.link)
  }
}

const refreshTodos = () => {
  loadTodoList()
}

const initTrendChart = () => {
  if (!trendChart.value) return
  
  const chart = echarts.init(trendChart.value)
  const option = {
    title: {
      show: false
    },
    tooltip: {
      trigger: 'axis',
      axisPointer: {
        type: 'cross',
        label: {
          backgroundColor: '#6a7985'
        }
      }
    },
    legend: {
      data: ['贷款审批', '租赁审批', '智能审批']
    },
    grid: {
      left: '3%',
      right: '4%',
      bottom: '3%',
      containLabel: true
    },
    xAxis: [
      {
        type: 'category',
        boundaryGap: false,
        data: ['周一', '周二', '周三', '周四', '周五', '周六', '周日']
      }
    ],
    yAxis: [
      {
        type: 'value'
      }
    ],
    series: [
      {
        name: '贷款审批',
        type: 'line',
        stack: 'Total',
        smooth: true,
        lineStyle: {
          width: 3
        },
        areaStyle: {
          opacity: 0.3
        },
        emphasis: {
          focus: 'series'
        },
        data: [12, 13, 10, 13, 9, 23, 21]
      },
      {
        name: '租赁审批',
        type: 'line',
        stack: 'Total',
        smooth: true,
        lineStyle: {
          width: 3
        },
        areaStyle: {
          opacity: 0.3
        },
        emphasis: {
          focus: 'series'
        },
        data: [22, 18, 19, 23, 29, 33, 31]
      },
      {
        name: '智能审批',
        type: 'line',
        stack: 'Total',
        smooth: true,
        lineStyle: {
          width: 3
        },
        areaStyle: {
          opacity: 0.3
        },
        emphasis: {
          focus: 'series'
        },
        data: [15, 23, 20, 15, 20, 24, 18]
      }
    ]
  }
  
  chart.setOption(option)
  
  // 响应式
  window.addEventListener('resize', () => {
    chart.resize()
  })
}

const initPieChart = () => {
  if (!pieChart.value) return
  
  const chart = echarts.init(pieChart.value)
  const option = {
    title: {
      show: false
    },
    tooltip: {
      trigger: 'item'
    },
    legend: {
      bottom: '5%',
      left: 'center'
    },
    series: [
      {
        type: 'pie',
        radius: ['40%', '70%'],
        center: ['50%', '45%'],
        avoidLabelOverlap: false,
        itemStyle: {
          borderRadius: 10,
          borderColor: '#fff',
          borderWidth: 2
        },
        label: {
          show: false,
          position: 'center'
        },
        emphasis: {
          label: {
            show: true,
            fontSize: 40,
            fontWeight: 'bold'
          }
        },
        labelLine: {
          show: false
        },
        data: [
          { value: 1048, name: '农业贷' },
          { value: 735, name: '创业贷' },
          { value: 580, name: '消费贷' },
          { value: 484, name: '经营贷' },
          { value: 300, name: '助学贷' }
        ]
      }
    ]
  }
  
  chart.setOption(option)
  
  // 响应式
  window.addEventListener('resize', () => {
    chart.resize()
  })
}

const loadTodoList = async () => {
  // 模拟数据
  todoList.value = [
    {
      id: 1,
      title: '审批农业贷款申请',
      description: '张三的春耕资金贷款申请待审批',
      type: '贷款审批',
      createdAt: new Date(),
      link: '/approval/loan'
    },
    {
      id: 2,
      title: '拖拉机租赁审批',
      description: '李四申请租赁大型拖拉机',
      type: '租赁审批',
      createdAt: new Date(Date.now() - 3600000),
      link: '/approval/lease'
    },
    {
      id: 3,
      title: '系统更新通知',
      description: '系统将于今晚进行维护更新',
      type: '系统通知',
      createdAt: new Date(Date.now() - 7200000),
      link: null
    }
  ]
}

const loadRecentActivities = async () => {
  // 模拟数据
  recentActivities.value = [
    {
      id: 1,
      user: '王审批员',
      action: '通过了',
      target: '张三的农业贷款申请',
      createdAt: new Date(Date.now() - 1800000)
    },
    {
      id: 2,
      user: '李审批员',
      action: '拒绝了',
      target: '王五的消费贷款申请',
      createdAt: new Date(Date.now() - 3600000)
    },
    {
      id: 3,
      user: '赵审批员',
      action: '审批了',
      target: '陈六的拖拉机租赁申请',
      createdAt: new Date(Date.now() - 5400000)
    },
    {
      id: 4,
      user: '孙审批员',
      action: '通过了',
      target: '刘七的经营贷款申请',
      createdAt: new Date(Date.now() - 7200000)
    }
  ]
}

const loadTodayStats = async () => {
  // 模拟数据
  todayStats.value = {
    approvals: 12,
    pending: 8
  }
}

// 生命周期
onMounted(async () => {
  await Promise.all([
    loadTodoList(),
    loadRecentActivities(),
    loadTodayStats()
  ])
  
  nextTick(() => {
    initTrendChart()
    initPieChart()
  })
})
</script>

<style scoped>
.dashboard {
  padding: 0;
}

/* 欢迎区域 */
.welcome-section {
  margin-bottom: 24px;
}

.welcome-card {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  border: none;
  color: white;
}

.welcome-card :deep(.el-card__body) {
  padding: 32px;
}

.welcome-content {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.welcome-info {
  flex: 1;
}

.welcome-title {
  font-size: 28px;
  font-weight: 600;
  margin: 0 0 8px 0;
  color: white;
}

.welcome-desc {
  font-size: 16px;
  margin: 0 0 20px 0;
  opacity: 0.9;
}

.quick-stats {
  display: flex;
  gap: 32px;
}

.stat-item {
  display: flex;
  flex-direction: column;
  align-items: flex-start;
}

.stat-label {
  font-size: 14px;
  opacity: 0.8;
  margin-bottom: 4px;
}

.stat-value {
  font-size: 24px;
  font-weight: 600;
}

.welcome-actions {
  display: flex;
  gap: 16px;
}

.welcome-actions .el-button {
  height: 40px;
  font-size: 16px;
}

/* 统计卡片 */
.stats-section {
  margin-bottom: 24px;
}

.stat-card {
  border: none;
  border-radius: 12px;
  overflow: hidden;
  height: 120px;
  position: relative;
  transition: all 0.3s ease;
}

.stat-card:hover {
  transform: translateY(-4px);
  box-shadow: 0 8px 25px rgba(0, 0, 0, 0.1);
}

.stat-card-content {
  display: flex;
  align-items: center;
  height: 100%;
  padding: 20px;
}

.stat-icon {
  width: 60px;
  height: 60px;
  border-radius: 12px;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 24px;
  margin-right: 16px;
}

.stat-card-loan .stat-icon {
  background: linear-gradient(135deg, #4facfe 0%, #00f2fe 100%);
  color: white;
}

.stat-card-lease .stat-icon {
  background: linear-gradient(135deg, #43e97b 0%, #38f9d7 100%);
  color: white;
}

.stat-card-pending .stat-icon {
  background: linear-gradient(135deg, #fa709a 0%, #fee140 100%);
  color: white;
}

.stat-card-completed .stat-icon {
  background: linear-gradient(135deg, #a8edea 0%, #fed6e3 100%);
  color: #333;
}

.stat-info {
  flex: 1;
}

.stat-number {
  font-size: 32px;
  font-weight: 700;
  color: #333;
  line-height: 1;
  margin-bottom: 4px;
}

.stat-title {
  font-size: 14px;
  color: #666;
  margin-bottom: 8px;
}

.stat-change {
  display: flex;
  align-items: center;
  gap: 4px;
  font-size: 12px;
  font-weight: 500;
}

.stat-change.up {
  color: #67c23a;
}

.stat-change.down {
  color: #f56c6c;
}

.stat-change.same {
  color: #909399;
}

/* 图表区域 */
.charts-section {
  margin-bottom: 24px;
}

.chart-card,
.todo-card,
.activity-card {
  border: none;
  border-radius: 12px;
  height: 400px;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.header-actions {
  display: flex;
  align-items: center;
  gap: 12px;
}

.chart-container {
  height: 300px;
}

/* 待办事项 */
.todo-list {
  height: 300px;
  overflow-y: auto;
}

.todo-item {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  padding: 16px 0;
  border-bottom: 1px solid #f0f0f0;
  cursor: pointer;
  transition: background-color 0.3s;
}

.todo-item:hover {
  background-color: #f8f9fa;
  border-radius: 8px;
  margin: 0 -12px;
  padding: 16px 12px;
}

.todo-item:last-child {
  border-bottom: none;
}

.todo-info {
  flex: 1;
}

.todo-title {
  font-weight: 500;
  color: #333;
  margin-bottom: 4px;
}

.todo-desc {
  font-size: 14px;
  color: #666;
  margin-bottom: 8px;
  line-height: 1.4;
}

.todo-time {
  font-size: 12px;
  color: #999;
}

.todo-action {
  margin-left: 12px;
}

.todo-empty {
  height: 200px;
  display: flex;
  align-items: center;
  justify-content: center;
}

/* 底部区域 */
.bottom-section {
  margin-bottom: 24px;
}

/* 活动列表 */
.activity-list {
  height: 300px;
  overflow-y: auto;
}

.activity-item {
  display: flex;
  align-items: flex-start;
  padding: 12px 0;
  border-bottom: 1px solid #f0f0f0;
}

.activity-item:last-child {
  border-bottom: none;
}

.activity-avatar {
  margin-right: 12px;
  flex-shrink: 0;
}

.activity-content {
  flex: 1;
}

.activity-text {
  margin-bottom: 4px;
}

.activity-user {
  font-weight: 500;
  color: #333;
}

.activity-action {
  color: #666;
  margin: 0 4px;
}

.activity-target {
  color: #409eff;
}

.activity-time {
  font-size: 12px;
  color: #999;
}

.activity-empty {
  height: 200px;
  display: flex;
  align-items: center;
  justify-content: center;
}

/* 响应式设计 */
@media (max-width: 768px) {
  .welcome-content {
    flex-direction: column;
    text-align: center;
    gap: 24px;
  }
  
  .welcome-actions {
    width: 100%;
    justify-content: center;
  }
  
  .quick-stats {
    justify-content: center;
  }
  
  .stat-card-content {
    padding: 16px;
  }
  
  .stat-icon {
    width: 48px;
    height: 48px;
    font-size: 20px;
  }
  
  .stat-number {
    font-size: 24px;
  }
}
</style> 