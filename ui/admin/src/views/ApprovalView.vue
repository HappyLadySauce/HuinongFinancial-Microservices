<template>
  <div class="approval-view">
    <div class="page-header">
      <h2 class="page-title">审批看板</h2>
      <div class="header-actions">
        <el-button type="primary" @click="refreshData" :loading="loading">
          <el-icon><Refresh /></el-icon>
          刷新
        </el-button>
      </div>
    </div>
    
    <!-- 筛选条件 -->
    <el-card class="filter-card" shadow="never">
      <el-form :model="filterForm" inline size="default">
        <el-form-item label="申请状态">
          <el-select v-model="filterForm.status_filter" placeholder="请选择状态" clearable>
            <el-option label="全部" value="" />
            <el-option label="AI审批中" value="AI_审批中" />
            <el-option label="待人工复核" value="待人工复核" />
            <el-option label="已批准" value="已批准" />
            <el-option label="已拒绝" value="已拒绝" />
          </el-select>
        </el-form-item>
        
        <el-form-item label="申请人姓名">
          <el-input
            v-model="filterForm.applicant_name"
            placeholder="请输入申请人姓名"
            clearable
          />
        </el-form-item>
        
        <el-form-item label="申请编号">
          <el-input
            v-model="filterForm.application_id"
            placeholder="请输入申请编号"
            clearable
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
          <div class="stat-label">总申请数</div>
          <div class="chart-container">
            <div ref="barChartRef" class="mini-chart"></div>
          </div>
        </div>
      </el-card>
      <el-card class="stat-item" shadow="hover">
        <div class="stat-content">
          <div class="stat-value pending">{{ statistics.pending }}</div>
          <div class="stat-label">待处理</div>
          <div class="chart-container chart-container-large">
            <div ref="pieChartRef" class="mini-chart"></div>
          </div>
        </div>
      </el-card>
      <el-card class="stat-item" shadow="hover">
        <div class="stat-content">
          <div class="stat-value approved">{{ statistics.approved }}</div>
          <div class="stat-label">已批准</div>
          <div class="chart-container">
            <div ref="lineChartRef" class="mini-chart"></div>
          </div>
        </div>
      </el-card>
      <el-card class="stat-item" shadow="hover">
        <div class="stat-content">
          <div class="stat-value rejected">{{ statistics.rejected }}</div>
          <div class="stat-label">已拒绝</div>
          <div class="chart-container chart-container-large">
            <div ref="doughnutChartRef" class="mini-chart"></div>
          </div>
        </div>
      </el-card>
    </div>
    
    <!-- 申请列表 -->
    <el-card class="table-card" shadow="never">
      <template #header>
        <div class="card-header">
          <span>申请列表</span>
          <div class="header-extra">
            <el-tag v-if="selectedRows.length > 0" type="info">
              已选择 {{ selectedRows.length }} 项
            </el-tag>
          </div>
        </div>
      </template>
      
      <el-table
        v-loading="loading"
        :data="applications"
        @selection-change="handleSelectionChange"
        @row-click="handleRowClick"
        stripe
        style="width: 100%"
      >
        <el-table-column type="selection" width="55" />
        
        <el-table-column prop="application_id" label="申请编号" width="140">
          <template #default="{ row }">
            <el-link type="primary" @click.stop="viewDetail(row.application_id)">
              {{ row.application_id }}
            </el-link>
          </template>
        </el-table-column>
        
        <el-table-column prop="applicant_name" label="申请人" width="100" />
        
        <el-table-column prop="amount" label="申请金额" width="120">
          <template #default="{ row }">
            <span class="amount">¥{{ formatAmount(row.amount) }}</span>
          </template>
        </el-table-column>
        
        <el-table-column prop="status" label="状态" width="120">
          <template #default="{ row }">
            <el-tag :type="getStatusType(row.status)" size="small">
              {{ row.status }}
            </el-tag>
          </template>
        </el-table-column>
        
        <el-table-column prop="ai_risk_score" label="AI风险评分" width="120">
          <template #default="{ row }">
            <div v-if="row.ai_risk_score !== undefined" class="risk-score">
              <el-progress
                :percentage="row.ai_risk_score"
                :color="getRiskColor(row.ai_risk_score)"
                :stroke-width="6"
                :show-text="false"
              />
              <span class="score-text">{{ row.ai_risk_score }}分</span>
            </div>
            <span v-else class="no-score">-</span>
          </template>
        </el-table-column>
        
        <el-table-column prop="ai_suggestion" label="AI建议" min-width="200">
          <template #default="{ row }">
            <el-tooltip
              v-if="row.ai_suggestion"
              :content="row.ai_suggestion"
              placement="top"
              :show-after="500"
            >
              <span class="suggestion-text">{{ row.ai_suggestion }}</span>
            </el-tooltip>
            <span v-else class="no-suggestion">-</span>
          </template>
        </el-table-column>
        
        <el-table-column prop="submission_time" label="提交时间" width="160">
          <template #default="{ row }">
            {{ formatDateTime(row.submission_time) }}
          </template>
        </el-table-column>
        
        <el-table-column label="操作" width="200" fixed="right">
          <template #default="{ row }">
            <el-button
              type="primary"
              size="small"
              @click.stop="viewDetail(row.application_id)"
            >
              查看详情
            </el-button>
            <el-button
              v-if="canReview(row.status)"
              type="success"
              size="small"
              @click.stop="quickApprove(row)"
            >
              快速审批
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
    
    <!-- 快速审批对话框 -->
    <el-dialog
      v-model="quickApprovalVisible"
      title="快速审批"
      width="500px"
      @close="resetQuickApprovalForm"
    >
      <el-form
        ref="quickApprovalFormRef"
        :model="quickApprovalForm"
        :rules="quickApprovalRules"
        label-width="100px"
      >
        <el-form-item label="申请编号">
          <el-input :value="selectedApplication?.application_id" disabled />
        </el-form-item>
        
        <el-form-item label="申请人">
          <el-input :value="selectedApplication?.applicant_name" disabled />
        </el-form-item>
        
        <el-form-item label="申请金额">
          <el-input :value="formatAmount(selectedApplication?.amount || 0)" disabled />
        </el-form-item>
        
        <el-form-item label="审批决策" prop="decision" required>
          <el-radio-group v-model="quickApprovalForm.decision">
            <el-radio value="approved">批准</el-radio>
            <el-radio value="rejected">拒绝</el-radio>
            <el-radio value="require_more_info">要求补充材料</el-radio>
          </el-radio-group>
        </el-form-item>
        
        <el-form-item
          v-if="quickApprovalForm.decision === 'approved'"
          label="批准金额"
          prop="approved_amount"
        >
          <el-input-number
            v-model="quickApprovalForm.approved_amount"
            :min="1"
            :max="selectedApplication?.amount || 0"
            :step="1000"
            style="width: 100%"
          />
        </el-form-item>
        
        <el-form-item label="审批意见" prop="comments" required>
          <el-input
            v-model="quickApprovalForm.comments"
            type="textarea"
            :rows="4"
            placeholder="请输入审批意见"
          />
        </el-form-item>
        
        <el-form-item
          v-if="quickApprovalForm.decision === 'require_more_info'"
          label="补充说明"
          prop="required_info_details"
        >
          <el-input
            v-model="quickApprovalForm.required_info_details"
            type="textarea"
            :rows="3"
            placeholder="请说明需要补充的材料或信息"
          />
        </el-form-item>
      </el-form>
      
      <template #footer>
        <el-button @click="quickApprovalVisible = false">取消</el-button>
        <el-button
          type="primary"
          @click="submitQuickApproval"
          :loading="submitting"
        >
          提交审批
        </el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, computed, onMounted, onBeforeUnmount } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage, ElMessageBox, type FormInstance, type FormRules } from 'element-plus'
import {
  Refresh,
  Search,
  RefreshLeft
} from '@element-plus/icons-vue'
import { adminLoanApprovalApi } from '@/services/api'
import type { LoanApplication } from '@/types'
import dayjs from 'dayjs'
import * as echarts from 'echarts'

const router = useRouter()

const loading = ref(false)
const submitting = ref(false)
const quickApprovalVisible = ref(false)
const selectedRows = ref<LoanApplication[]>([])
const selectedApplication = ref<LoanApplication | null>(null)
const quickApprovalFormRef = ref<FormInstance>()

// 图表DOM引用
const barChartRef = ref<HTMLElement | null>(null)
const pieChartRef = ref<HTMLElement | null>(null)
const lineChartRef = ref<HTMLElement | null>(null)
const doughnutChartRef = ref<HTMLElement | null>(null)

// 图表实例
let barChart: echarts.ECharts | null = null
let pieChart: echarts.ECharts | null = null
let lineChart: echarts.ECharts | null = null
let doughnutChart: echarts.ECharts | null = null

// 筛选表单
const filterForm = reactive({
  status_filter: '',
  applicant_name: '',
  application_id: ''
})

// 分页信息
const pagination = reactive({
  page: 1,
  limit: 20,
  total: 0
})

// 申请列表
const applications = ref<LoanApplication[]>([])

// 统计信息 - 固定数值
const statistics = computed(() => {
  return {
    total: 278,
    pending: 38,
    approved: 195,
    rejected: 45
  }
})

// 快速审批表单
const quickApprovalForm = reactive({
  decision: '',
  approved_amount: 0,
  comments: '',
  required_info_details: ''
})

const quickApprovalRules: FormRules = {
  decision: [
    { required: true, message: '请选择审批决策', trigger: 'change' }
  ],
  comments: [
    { required: true, message: '请输入审批意见', trigger: 'blur' }
  ],
  approved_amount: [
    { required: true, message: '请输入批准金额', trigger: 'blur' }
  ]
}

// 方法
const fetchApplications = async () => {
  try {
    loading.value = true
    const params = {
      ...filterForm,
      page: pagination.page,
      size: pagination.limit
    }
    
    const data = await adminLoanApprovalApi.getApprovals(params)
    applications.value = data.list || []
    pagination.total = data.total || 0
  } catch (error) {
    ElMessage.error('获取申请列表失败')
  } finally {
    loading.value = false
  }
}

const refreshData = () => {
  fetchApplications()
}

const handleSearch = () => {
  pagination.page = 1
  fetchApplications()
}

const handleReset = () => {
  Object.assign(filterForm, {
    status_filter: '',
    applicant_name: '',
    application_id: ''
  })
  pagination.page = 1
  fetchApplications()
}

const handleSelectionChange = (selection: LoanApplication[]) => {
  selectedRows.value = selection
}

const handleRowClick = (row: LoanApplication) => {
  viewDetail(row.application_id)
}

const handleSizeChange = () => {
  pagination.page = 1
  fetchApplications()
}

const handleCurrentChange = () => {
  fetchApplications()
}

const viewDetail = (applicationId: string) => {
  router.push(`/approval/${applicationId}`)
}

const canReview = (status: string) => {
  return status === '待人工复核'
}

const quickApprove = (application: LoanApplication) => {
  selectedApplication.value = application
  quickApprovalForm.approved_amount = application.amount
  quickApprovalVisible.value = true
}

const resetQuickApprovalForm = () => {
  Object.assign(quickApprovalForm, {
    decision: '',
    approved_amount: 0,
    comments: '',
    required_info_details: ''
  })
  selectedApplication.value = null
}

const submitQuickApproval = async () => {
  if (!quickApprovalFormRef.value || !selectedApplication.value) return
  
  try {
    await quickApprovalFormRef.value.validate()
    submitting.value = true
    
    const submitData = {
      decision: quickApprovalForm.decision as 'approved' | 'rejected' | 'require_more_info',
      approved_amount: quickApprovalForm.decision === 'approved' ? quickApprovalForm.approved_amount : undefined,
      comments: quickApprovalForm.comments,
      required_info_details: quickApprovalForm.decision === 'require_more_info' ? quickApprovalForm.required_info_details : undefined
    }
    
    await adminLoanApprovalApi.reviewApproval(selectedApplication.value.application_id, submitData)
    
    ElMessage.success('审批提交成功')
    quickApprovalVisible.value = false
    resetQuickApprovalForm()
    fetchApplications()
  } catch (error) {
    ElMessage.error('审批提交失败')
  } finally {
    submitting.value = false
  }
}

const getStatusType = (status: string) => {
  const statusMap: Record<string, string> = {
    'AI_审批中': 'warning',
    '待人工复核': 'info',
    '已批准': 'success',
    '已拒绝': 'danger'
  }
  return statusMap[status] || 'info'
}

const getRiskColor = (score: number) => {
  if (score <= 30) return '#67c23a'
  if (score <= 70) return '#e6a23c'
  return '#f56c6c'
}

const formatAmount = (amount: number) => {
  return amount.toLocaleString()
}

const formatDateTime = (datetime: string) => {
  return dayjs(datetime).format('YYYY-MM-DD HH:mm')
}

// 初始化图表
const initCharts = () => {
  // 柱状图 - 总申请数
  if (barChartRef.value) {
    barChart = echarts.init(barChartRef.value)
    const option = {
      tooltip: {
        trigger: 'axis'
      },
      grid: {
        left: '5%',
        right: '5%',
        bottom: '5%',
        top: '5%',
        containLabel: false
      },
      xAxis: {
        type: 'category',
        data: ['1月', '2月', '3月', '4月', '5月', '6月'],
        axisLine: { show: false },
        axisTick: { show: false },
        axisLabel: { show: false }
      },
      yAxis: {
        type: 'value',
        show: false
      },
      series: [{
        data: [220, 240, 260, 250, 270, 278],
        type: 'bar',
        itemStyle: {
          color: new echarts.graphic.LinearGradient(0, 0, 0, 1, [
            { offset: 0, color: '#667eea' },
            { offset: 1, color: '#764ba2' }
          ])
        },
        barWidth: '60%'
      }]
    }
    barChart.setOption(option)
  }

  // 饼图 - 待处理
  if (pieChartRef.value) {
    pieChart = echarts.init(pieChartRef.value)
    const option = {
      tooltip: {
        trigger: 'item'
      },
      series: [{
        type: 'pie',
        radius: ['40%', '70%'],
        data: [
          { value: 38, name: '待处理', itemStyle: { color: '#e6a23c' } },
          { value: 240, name: '其他', itemStyle: { color: '#f0f2f5' } }
        ],
        label: { show: false },
        emphasis: { label: { show: false } }
      }]
    }
    pieChart.setOption(option)
  }

  // 走势图 - 已批准
  if (lineChartRef.value) {
    lineChart = echarts.init(lineChartRef.value)
    const option = {
      tooltip: {
        trigger: 'axis'
      },
      grid: {
        left: '5%',
        right: '5%',
        bottom: '5%',
        top: '5%',
        containLabel: false
      },
      xAxis: {
        type: 'category',
        data: ['1月', '2月', '3月', '4月', '5月', '6月'],
        axisLine: { show: false },
        axisTick: { show: false },
        axisLabel: { show: false }
      },
      yAxis: {
        type: 'value',
        show: false
      },
      series: [{
        data: [150, 160, 170, 180, 190, 195],
        type: 'line',
        smooth: true,
        symbol: 'none',
        lineStyle: {
          width: 3,
          color: '#67c23a'
        },
        areaStyle: {
          color: new echarts.graphic.LinearGradient(0, 0, 0, 1, [
            { offset: 0, color: 'rgba(103, 194, 58, 0.3)' },
            { offset: 1, color: 'rgba(103, 194, 58, 0.05)' }
          ])
        }
      }]
    }
    lineChart.setOption(option)
  }

  // 扇形图 - 已拒绝
  if (doughnutChartRef.value) {
    doughnutChart = echarts.init(doughnutChartRef.value)
    const option = {
      tooltip: {
        trigger: 'item'
      },
      series: [{
        type: 'pie',
        radius: ['50%', '80%'],
        data: [
          { value: 45, name: '已拒绝', itemStyle: { color: '#f56c6c' } },
          { value: 233, name: '其他', itemStyle: { color: '#f0f2f5' } }
        ],
        label: { show: false },
        emphasis: { label: { show: false } },
        startAngle: 90
      }]
    }
    doughnutChart.setOption(option)
  }
}

const resizeCharts = () => {
  barChart?.resize()
  pieChart?.resize()
  lineChart?.resize()
  doughnutChart?.resize()
}

onMounted(() => {
  fetchApplications()
  // 等待DOM更新后初始化图表
  setTimeout(() => {
    initCharts()
    window.addEventListener('resize', resizeCharts)
  }, 200)
})

onBeforeUnmount(() => {
  window.removeEventListener('resize', resizeCharts)
  barChart?.dispose()
  pieChart?.dispose()
  lineChart?.dispose()
  doughnutChart?.dispose()
})
</script>

<style scoped>
.approval-view {
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
  transition: transform 0.2s, box-shadow 0.2s;
}

.stat-item:hover {
  transform: translateY(-2px);
  box-shadow: 0 8px 25px rgba(0, 0, 0, 0.1);
}

.stat-content {
  padding: 30px;
}

.stat-value {
  font-size: 28px;
  font-weight: 700;
  color: #333;
  margin-bottom: 8px;
}

.stat-value.pending {
  color: #e6a23c;
}

.stat-value.approved {
  color: #67c23a;
}

.stat-value.rejected {
  color: #f56c6c;
}

.stat-label {
  color: #666;
  font-size: 14px;
  margin-bottom: 12px;
}

.chart-container {
  height: 80px;
  width: 100%;
  margin-top: 8px;
}

.chart-container-large {
  height: 120px;
}

.mini-chart {
  height: 100%;
  width: 100%;
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

.amount {
  font-weight: 600;
  color: #f56c6c;
}

.risk-score {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 4px;
}

.score-text {
  font-size: 12px;
  color: #666;
}

.no-score,
.no-suggestion {
  color: #c0c4cc;
}

.suggestion-text {
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
  text-overflow: ellipsis;
  cursor: help;
}

.pagination-wrapper {
  display: flex;
  justify-content: center;
  margin-top: 20px;
}

:deep(.el-table__row) {
  cursor: pointer;
}

:deep(.el-table__row:hover) {
  background-color: #f5f7fa;
}

:deep(.el-progress-bar__outer) {
  background-color: #e4e7ed;
}

:deep(.el-form-item) {
  margin-bottom: 16px;
}
</style> 