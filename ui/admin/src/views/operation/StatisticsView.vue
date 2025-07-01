<template>
  <div class="statistics-view">
    <!-- 页面头部 -->
    <PageHeader
      title="数据统计"
      subtitle="业务数据统计和分析报告"
      :icon="DataAnalysis"
      :breadcrumbs="breadcrumbs"
    >
      <template #actions>
        <el-button @click="handleExport">
          <el-icon><Download /></el-icon>
          导出报告
        </el-button>
        <el-button @click="handleRefresh">
          <el-icon><Refresh /></el-icon>
          刷新数据
        </el-button>
      </template>
    </PageHeader>

    <!-- 统计概览 -->
    <el-row :gutter="24" class="stats-overview">
      <el-col :span="6">
        <StatCard
          title="总申请量"
          :value="overviewData.totalApplications"
          format="number"
          type="primary"
          :change="overviewData.applicationChange"
          icon="Document"
          @click="handleStatClick('applications')"
        />
      </el-col>
      <el-col :span="6">
        <StatCard
          title="审批通过率"
          :value="overviewData.approvalRate"
          format="percent"
          type="success"
          :change="overviewData.approvalChange"
          icon="SuccessFilled"
          @click="handleStatClick('approval')"
        />
      </el-col>
      <el-col :span="6">
        <StatCard
          title="平均审批时长"
          :value="overviewData.avgProcessTime"
          suffix="小时"
          type="warning"
          :change="overviewData.timeChange"
          icon="Timer"
          @click="handleStatClick('time')"
        />
      </el-col>
      <el-col :span="6">
        <StatCard
          title="累计放款金额"
          :value="overviewData.totalAmount"
          format="currency"
          type="danger"
          :change="overviewData.amountChange"
          icon="Money"
          @click="handleStatClick('amount')"
        />
      </el-col>
    </el-row>

    <el-row :gutter="24">
      <!-- 左侧图表 -->
      <el-col :span="16">
        <!-- 趋势分析 -->
        <ChartCard
          title="业务趋势分析"
          :chart-data="trendChartData"
          chart-type="line"
          height="400px"
          :time-range="timeRange"
          @time-change="handleTimeChange"
          @refresh="loadTrendData"
        />

        <!-- 业务分布 -->
        <ChartCard
          title="业务类型分布"
          :chart-data="distributionChartData"
          chart-type="pie"
          height="350px"
          style="margin-top: 20px;"
          @refresh="loadDistributionData"
        />
      </el-col>

      <!-- 右侧详情 -->
      <el-col :span="8">
        <!-- 实时监控 -->
        <el-card class="monitor-card" shadow="never">
          <template #header>
            <div class="card-header">
              <span>
                <el-icon><Monitor /></el-icon>
                实时监控
              </span>
              <el-button size="small" @click="loadRealtimeData">
                <el-icon><Refresh /></el-icon>
              </el-button>
            </div>
          </template>

          <div class="realtime-stats">
            <div class="stat-item">
              <span class="label">今日新增申请</span>
              <span class="value">{{ realtimeData.todayApplications }}</span>
            </div>
            <div class="stat-item">
              <span class="label">待审批数量</span>
              <span class="value pending">{{ realtimeData.pendingCount }}</span>
            </div>
            <div class="stat-item">
              <span class="label">在线用户</span>
              <span class="value online">{{ realtimeData.onlineUsers }}</span>
            </div>
            <div class="stat-item">
              <span class="label">系统负载</span>
              <span class="value" :class="getLoadClass(realtimeData.systemLoad)">
                {{ realtimeData.systemLoad }}%
              </span>
            </div>
          </div>
        </el-card>

        <!-- Top榜单 -->
        <el-card class="ranking-card" shadow="never">
          <template #header>
            <div class="card-header">
              <span>
                <el-icon><Trophy /></el-icon>
                业务排行榜
              </span>
              <el-select v-model="rankingType" size="small" style="width: 120px;">
                <el-option label="申请量" value="applications" />
                <el-option label="通过率" value="approval" />
                <el-option label="放款额" value="amount" />
              </el-select>
            </div>
          </template>

          <div class="ranking-list">
            <div
              v-for="(item, index) in rankingData"
              :key="item.name"
              class="ranking-item"
            >
              <div class="rank">
                <el-icon v-if="index === 0" class="gold"><Trophy /></el-icon>
                <el-icon v-else-if="index === 1" class="silver"><Trophy /></el-icon>
                <el-icon v-else-if="index === 2" class="bronze"><Trophy /></el-icon>
                <span v-else>{{ index + 1 }}</span>
              </div>
              <div class="info">
                <div class="name">{{ item.name }}</div>
                <div class="value">{{ formatRankingValue(item.value) }}</div>
              </div>
              <div class="progress">
                <el-progress
                  :percentage="(item.value / rankingData[0].value) * 100"
                  :show-text="false"
                  :stroke-width="6"
                />
              </div>
            </div>
          </div>
        </el-card>
      </el-col>
    </el-row>

    <!-- 详细数据表格 -->
    <el-card class="table-card" shadow="never">
      <template #header>
        <div class="card-header">
          <span>
            <el-icon><List /></el-icon>
            详细数据
          </span>
          <div class="header-actions">
            <el-select v-model="tableFilter" size="small" style="width: 120px; margin-right: 12px;">
              <el-option label="全部" value="all" />
              <el-option label="本月" value="month" />
              <el-option label="本周" value="week" />
              <el-option label="今日" value="today" />
            </el-select>
            <el-button size="small" @click="loadTableData">
              <el-icon><Search /></el-icon>
              查询
            </el-button>
          </div>
        </div>
      </template>

      <el-table :data="tableData" border stripe>
        <el-table-column prop="date" label="日期" width="120" />
        <el-table-column prop="applications" label="申请数量" width="100" />
        <el-table-column prop="approved" label="通过数量" width="100" />
        <el-table-column prop="rejected" label="拒绝数量" width="100" />
        <el-table-column prop="approvalRate" label="通过率" width="100">
          <template #default="{ row }">
            {{ (row.approvalRate * 100).toFixed(1) }}%
          </template>
        </el-table-column>
        <el-table-column prop="amount" label="放款金额" width="120">
          <template #default="{ row }">
            ¥{{ row.amount.toLocaleString() }}
          </template>
        </el-table-column>
        <el-table-column prop="avgTime" label="平均审批时长" width="120">
          <template #default="{ row }">
            {{ row.avgTime }}小时
          </template>
        </el-table-column>
        <el-table-column prop="status" label="状态" width="80">
          <template #default="{ row }">
            <StatusTag :status="row.status" />
          </template>
        </el-table-column>
      </el-table>

      <div class="pagination-wrapper">
        <el-pagination
          v-model:current-page="pagination.page"
          v-model:page-size="pagination.size"
          :total="pagination.total"
          :page-sizes="[10, 20, 50, 100]"
          layout="total, sizes, prev, pager, next, jumper"
          @size-change="loadTableData"
          @current-change="loadTableData"
        />
      </div>
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, computed, onMounted, watch } from 'vue'
import { ElMessage } from 'element-plus'
import dayjs from 'dayjs'
import {
  DataAnalysis,
  Download,
  Refresh,
  Monitor,
  Trophy,
  List,
  Search
} from '@element-plus/icons-vue'

// 导入通用组件
import { PageHeader, StatCard, ChartCard, StatusTag } from '@/components/common'

const timeRange = ref('7days')
const rankingType = ref('applications')
const tableFilter = ref('all')

const pagination = reactive({
  page: 1,
  size: 20,
  total: 0
})

// 概览数据
const overviewData = ref({
  totalApplications: 12580,
  applicationChange: 8.5,
  approvalRate: 0.762,
  approvalChange: -2.1,
  avgProcessTime: 24.5,
  timeChange: -12.3,
  totalAmount: 158760000,
  amountChange: 15.2
})

// 实时数据
const realtimeData = ref({
  todayApplications: 45,
  pendingCount: 128,
  onlineUsers: 89,
  systemLoad: 68
})

// 趋势图表数据
const trendChartData = ref({
  xAxis: ['12-23', '12-24', '12-25', '12-26', '12-27', '12-28', '12-29'],
  series: [
    {
      name: '申请量',
      data: [120, 135, 145, 158, 162, 175, 180]
    },
    {
      name: '通过量',
      data: [98, 108, 115, 125, 128, 135, 142]
    }
  ]
})

// 分布图表数据
const distributionChartData = ref({
  series: [
    { name: '贷款申请', value: 6800 },
    { name: '租赁申请', value: 3200 },
    { name: '保险申请', value: 1580 },
    { name: '其他', value: 1000 }
  ]
})

// 排行榜数据
const rankingData = ref([
  { name: '农业贷款', value: 2580 },
  { name: '设备租赁', value: 1890 },
  { name: '土地抵押', value: 1650 },
  { name: '种植保险', value: 1200 },
  { name: '养殖贷款', value: 980 }
])

// 表格数据
const tableData = ref([
  {
    date: '2024-12-29',
    applications: 180,
    approved: 142,
    rejected: 38,
    approvalRate: 0.789,
    amount: 5680000,
    avgTime: 22.5,
    status: 'normal'
  },
  {
    date: '2024-12-28',
    applications: 175,
    approved: 135,
    rejected: 40,
    approvalRate: 0.771,
    amount: 5240000,
    avgTime: 24.8,
    status: 'normal'
  }
])

// 计算属性
const breadcrumbs = computed(() => [
  { title: '首页', to: '/' },
  { title: '运营管理', to: '/operation' },
  { title: '数据统计' }
])

// 监听排行榜类型变化
watch(rankingType, (newType) => {
  loadRankingData(newType)
})

// 方法
const handleStatClick = (type: string) => {
  console.log('点击统计卡片:', type)
  // 可以跳转到详细页面或显示详细信息
}

const handleTimeChange = (range: string) => {
  timeRange.value = range
  loadTrendData()
}

const handleExport = () => {
  ElMessage.success('报告导出功能待实现')
}

const handleRefresh = () => {
  loadAllData()
  ElMessage.success('数据已刷新')
}

const loadTrendData = () => {
  console.log('加载趋势数据:', timeRange.value)
  // 根据时间范围加载数据
}

const loadDistributionData = () => {
  console.log('加载分布数据')
  // 实现加载逻辑
}

const loadRealtimeData = () => {
  console.log('加载实时数据')
  // 实现实时数据加载
}

const loadRankingData = (type: string) => {
  console.log('加载排行榜数据:', type)
  // 根据类型加载不同的排行榜数据
}

const loadTableData = () => {
  console.log('加载表格数据:', {
    filter: tableFilter.value,
    page: pagination.page,
    size: pagination.size
  })
  // 实现表格数据加载
}

const loadAllData = () => {
  loadTrendData()
  loadDistributionData()
  loadRealtimeData()
  loadRankingData(rankingType.value)
  loadTableData()
}

const getLoadClass = (load: number) => {
  if (load > 80) return 'danger'
  if (load > 60) return 'warning'
  return 'success'
}

const formatRankingValue = (value: number) => {
  if (rankingType.value === 'amount') {
    return `¥${(value / 10000).toFixed(1)}万`
  } else if (rankingType.value === 'approval') {
    return `${(value * 100).toFixed(1)}%`
  }
  return value.toLocaleString()
}

onMounted(() => {
  loadAllData()
})
</script>

<style scoped>
.statistics-view {
  height: 100vh;
  display: flex;
  flex-direction: column;
  gap: 16px;
  padding: 0 20px;
  overflow: hidden;
}

.stats-overview {
  margin-bottom: 20px;
}

.monitor-card,
.ranking-card,
.table-card {
  margin-bottom: 20px;
  border: none;
  border-radius: 12px;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  font-weight: 500;
}

.header-actions {
  display: flex;
  align-items: center;
}

.realtime-stats {
  padding: 16px 0;
}

.stat-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 16px;
  padding: 12px 0;
  border-bottom: 1px solid #f0f0f0;
}

.stat-item:last-child {
  border-bottom: none;
}

.stat-item .label {
  color: #666;
  font-size: 14px;
}

.stat-item .value {
  font-weight: 600;
  font-size: 16px;
}

.stat-item .value.pending {
  color: #f56c6c;
}

.stat-item .value.online {
  color: #67c23a;
}

.stat-item .value.success {
  color: #67c23a;
}

.stat-item .value.warning {
  color: #e6a23c;
}

.stat-item .value.danger {
  color: #f56c6c;
}

.ranking-list {
  padding: 16px 0;
}

.ranking-item {
  display: flex;
  align-items: center;
  margin-bottom: 16px;
  padding: 12px 0;
  border-bottom: 1px solid #f0f0f0;
}

.ranking-item:last-child {
  border-bottom: none;
}

.ranking-item .rank {
  width: 40px;
  text-align: center;
  font-weight: 600;
}

.ranking-item .rank .gold {
  color: #ffd700;
}

.ranking-item .rank .silver {
  color: #c0c0c0;
}

.ranking-item .rank .bronze {
  color: #cd7f32;
}

.ranking-item .info {
  flex: 1;
  margin-left: 12px;
}

.ranking-item .info .name {
  font-weight: 500;
  margin-bottom: 4px;
}

.ranking-item .info .value {
  font-size: 12px;
  color: #666;
}

.ranking-item .progress {
  width: 80px;
  margin-left: 12px;
}

.pagination-wrapper {
  display: flex;
  justify-content: center;
  margin-top: 20px;
}

:deep(.el-card__header) {
  padding: 16px 20px;
  border-bottom: 1px solid #f0f0f0;
}

:deep(.el-card__body) {
  padding: 20px;
}
</style>
