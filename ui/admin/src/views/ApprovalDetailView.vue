<template>
  <div class="approval-detail">
    <div class="page-header">
      <div class="header-left">
        <el-button @click="goBack" size="default">
          <el-icon><ArrowLeft /></el-icon>
          返回
        </el-button>
        <h2 class="page-title">审批详情</h2>
      </div>
      <div class="header-actions">
        <el-button @click="refreshData" :loading="loading">
          <el-icon><Refresh /></el-icon>
          刷新
        </el-button>
        <el-button type="primary" @click="exportDetail" v-if="detail.application_info?.application_id">
          <el-icon><Download /></el-icon>
          导出详情
        </el-button>
      </div>
    </div>

    <!-- 统计概览 -->
    <div class="stats-row" v-if="detail.application_info?.application_id">
      <el-card class="stat-item" shadow="hover">
        <div class="stat-content">
          <div class="stat-value">{{ detail.application_info.application_id }}</div>
          <div class="stat-label">申请编号</div>
          <div class="chart-container">
            <div ref="statusChartRef" class="mini-chart"></div>
          </div>
        </div>
      </el-card>
      <el-card class="stat-item" shadow="hover">
        <div class="stat-content">
          <div class="stat-value amount">¥{{ formatAmount(detail.application_info?.amount || 0) }}</div>
          <div class="stat-label">申请金额</div>
          <div class="chart-container">
            <div ref="amountChartRef" class="mini-chart"></div>
          </div>
        </div>
      </el-card>
      <el-card class="stat-item" shadow="hover">
        <div class="stat-content">
          <div class="stat-value" :class="getRiskClass(50)">
            50分
          </div>
          <div class="stat-label">风险评分</div>
          <div class="chart-container">
            <div ref="riskChartRef" class="mini-chart"></div>
          </div>
        </div>
      </el-card>
      <el-card class="stat-item" shadow="hover">
        <div class="stat-content">
          <div class="stat-value" :class="getStatusClass(detail.application_info?.status)">
            {{ detail.application_info?.status || '未知' }}
          </div>
          <div class="stat-label">当前状态</div>
          <div class="chart-container">
            <div ref="progressChartRef" class="mini-chart"></div>
          </div>
        </div>
      </el-card>
    </div>

    <div v-loading="loading" class="detail-content">
      <el-row :gutter="20">
        <!-- 左侧主要信息 -->
        <el-col :span="16">
          <!-- 基本信息 -->
          <el-card class="info-card" shadow="never">
            <template #header>
              <div class="card-header">
                <span>
                  <el-icon><Document /></el-icon>
                  申请基本信息
                </span>
                <div class="header-tags">
                  <el-tag :type="getStatusType(detail.application_info?.status || '')" size="large">
                    {{ detail.application_info?.status || '未知' }}
                  </el-tag>
                  <el-tag v-if="detail.application_info?.duration" type="info" size="small">
                    {{ detail.application_info.duration }}个月
                  </el-tag>
                </div>
              </div>
            </template>
            
            <el-descriptions :column="2" border>
              <el-descriptions-item label="申请编号">
                <el-text type="primary" tag="b">{{ detail.application_info?.application_id }}</el-text>
              </el-descriptions-item>
              <el-descriptions-item label="申请人">
                {{ detail.application_info?.applicant_name || '-' }}
              </el-descriptions-item>
              <el-descriptions-item label="用户ID">
                {{ detail.application_info?.user_id || '-' }}
              </el-descriptions-item>
              <el-descriptions-item label="产品名称">
                {{ detail.application_info?.name || '-' }}
              </el-descriptions-item>
              <el-descriptions-item label="产品类型" span="2">
                {{ detail.application_info?.type || '-' }}
              </el-descriptions-item>
              <el-descriptions-item label="申请金额">
                <span class="amount">¥{{ formatAmount(detail.application_info?.amount || 0) }}</span>
              </el-descriptions-item>
              <el-descriptions-item label="申请期限">
                {{ detail.application_info?.duration }}个月
              </el-descriptions-item>
              <el-descriptions-item label="贷款用途" span="2">
                {{ detail.application_info?.purpose || '-' }}
              </el-descriptions-item>
              <el-descriptions-item label="提交时间">
                {{ formatDateTime(detail.application_info?.created_at) }}
              </el-descriptions-item>
              <el-descriptions-item label="更新时间">
                {{ formatDateTime(detail.application_info?.updated_at) }}
              </el-descriptions-item>
            </el-descriptions>
          </el-card>

          <!-- 状态信息 -->
          <el-card class="status-card" shadow="never">
            <template #header>
              <div class="card-header">
                <span>
                  <el-icon><Cpu /></el-icon>
                  申请状态信息
                </span>
                <div class="header-tags">
                  <el-tag :type="getStatusType(detail.application_info?.status || '')" size="small">
                    {{ detail.application_info?.status || '未知' }}
                  </el-tag>
                </div>
              </div>
            </template>
            
            <div class="status-info">
              <el-descriptions :column="2" border>
                <el-descriptions-item label="当前状态">
                  <el-tag :type="getStatusType(detail.application_info?.status || '')" size="small">
                    {{ detail.application_info?.status || '未知' }}
                  </el-tag>
                </el-descriptions-item>
                <el-descriptions-item label="产品ID">
                  {{ detail.application_info?.product_id || '-' }}
                </el-descriptions-item>
                <el-descriptions-item label="创建时间">
                  {{ formatDateTime(detail.application_info?.created_at) }}
                </el-descriptions-item>
                <el-descriptions-item label="更新时间">
                  {{ formatDateTime(detail.application_info?.updated_at) }}
                </el-descriptions-item>
              </el-descriptions>
            </div>
          </el-card>

          <!-- 申请信息 -->
          <el-card class="info-card" shadow="never">
            <template #header>
              <div class="card-header">
                <span>
                  <el-icon><Folder /></el-icon>
                  申请信息汇总
                </span>
              </div>
            </template>
            
            <div class="summary-info">
              <el-alert
                title="申请处理提示"
                description="该申请正在等待审核处理，请及时查看并作出审批决定。"
                type="info"
                :closable="false"
                show-icon
              />
            </div>
          </el-card>
        </el-col>

        <!-- 右侧操作区域 -->
        <el-col :span="8">
          <!-- 审批操作 -->
          <el-card v-if="canReview(detail.application_info?.status)" class="review-card" shadow="never">
            <template #header>
              <div class="card-header">
                <span>
                  <el-icon><Check /></el-icon>
                  审批操作
                </span>
                <el-tag type="warning" size="small">待处理</el-tag>
              </div>
            </template>
            
            <el-form
              ref="reviewFormRef"
              :model="reviewForm"
              :rules="reviewRules"
              label-width="80px"
            >
              <el-form-item label="审批决策" prop="decision" required>
                <el-radio-group v-model="reviewForm.decision">
                  <el-radio value="approved" class="decision-radio">
                    <el-icon><CircleCheck /></el-icon>
                    批准
                  </el-radio>
                  <el-radio value="rejected" class="decision-radio">
                    <el-icon><CircleClose /></el-icon>
                    拒绝
                  </el-radio>
                  <el-radio value="require_more_info" class="decision-radio">
                    <el-icon><Warning /></el-icon>
                    补充材料
                  </el-radio>
                </el-radio-group>
              </el-form-item>
              
              <el-form-item
                v-if="reviewForm.decision === 'approved'"
                label="批准金额"
                prop="approved_amount"
              >
                <el-input-number
                  v-model="reviewForm.approved_amount"
                  :min="1"
                  :max="detail.application_info?.amount || 0"
                  :step="1000"
                  style="width: 100%"
                  :formatter="(value: number) => `¥ ${value}`.replace(/\B(?=(\d{3})+(?!\d))/g, ',')"
                  :parser="(value: string) => value.replace(/¥\s?|(,*)/g, '')"
                />
                <div class="amount-hint">
                  申请金额: ¥{{ formatAmount(detail.application_info?.amount || 0) }}
                </div>
              </el-form-item>
              
              <el-form-item label="审批意见" prop="comments" required>
                <el-input
                  v-model="reviewForm.comments"
                  type="textarea"
                  :rows="4"
                  placeholder="请输入审批意见"
                  show-word-limit
                  maxlength="500"
                />
              </el-form-item>
              
              <el-form-item
                v-if="reviewForm.decision === 'require_more_info'"
                label="补充说明"
                prop="required_info_details"
              >
                <el-input
                  v-model="reviewForm.required_info_details"
                  type="textarea"
                  :rows="3"
                  placeholder="请说明需要补充的材料或信息"
                  show-word-limit
                  maxlength="300"
                />
              </el-form-item>
              
              <el-form-item>
                <el-button
                  type="primary"
                  @click="submitApprovalReview"
                  :loading="submitting"
                  style="width: 100%"
                  size="large"
                >
                  <el-icon><Check /></el-icon>
                  提交审批
                </el-button>
              </el-form-item>
            </el-form>
          </el-card>

          <!-- 申请进度 -->
          <el-card class="progress-card" shadow="never">
            <template #header>
              <div class="card-header">
                <span>
                  <el-icon><Operation /></el-icon>
                  申请进度
                </span>
              </div>
            </template>
            
            <div class="progress-steps">
              <el-steps direction="vertical" :active="getProgressStep(detail.application_info?.status)" finish-status="success">
                <el-step title="申请提交" :description="formatDateTime(detail.application_info?.created_at)" />
                <el-step title="初步审核" description="系统初步验证" />
                <el-step title="人工复核" description="专业审批员审核" />
                <el-step title="审批完成" description="最终审批结果" />
              </el-steps>
            </div>
          </el-card>

          <!-- 审批历史 -->
          <el-card class="history-card" shadow="never">
            <template #header>
              <div class="card-header">
                <span>
                  <el-icon><Clock /></el-icon>
                  审批历史
                </span>
                <el-tag type="info" size="small">
                  2 条记录
                </el-tag>
              </div>
            </template>
            
            <el-timeline>
              <el-timeline-item
                :timestamp="formatDateTime(detail.application_info?.created_at)"
                placement="top"
                type="primary"
              >
                <div class="timeline-content">
                  <div class="timeline-status">申请提交</div>
                  <div class="timeline-operator">申请人: {{ detail.application_info?.applicant_name }}</div>
                  <div class="timeline-comments">提交贷款申请</div>
                </div>
              </el-timeline-item>
              <el-timeline-item
                v-if="detail.application_info?.updated_at !== detail.application_info?.created_at"
                :timestamp="formatDateTime(detail.application_info?.updated_at)"
                placement="top"
                type="warning"
              >
                <div class="timeline-content">
                  <div class="timeline-status">状态更新</div>
                  <div class="timeline-operator">系统</div>
                  <div class="timeline-comments">申请状态已更新</div>
                </div>
              </el-timeline-item>
            </el-timeline>
          </el-card>
        </el-col>
      </el-row>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted, onBeforeUnmount } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { ElMessage, type FormInstance, type FormRules } from 'element-plus'
import {
  ArrowLeft,
  Refresh,
  Document,
  Cpu,
  Folder,
  Check,
  CircleCheck,
  CircleClose,
  Warning,
  Clock,
  Download,
  View,
  Operation,
  TrendCharts
} from '@element-plus/icons-vue'
import { adminLoanApprovalApi } from '@/services/api'
import type { LoanApplication } from '@/types'
import type { ApplicationDetail } from '@/types'
import dayjs from 'dayjs'
import * as echarts from 'echarts'

const route = useRoute()
const router = useRouter()

const loading = ref(false)
const submitting = ref(false)
const reviewFormRef = ref<FormInstance>()

// 图表DOM引用
const statusChartRef = ref<HTMLElement | null>(null)
const amountChartRef = ref<HTMLElement | null>(null)
const riskChartRef = ref<HTMLElement | null>(null)
const progressChartRef = ref<HTMLElement | null>(null)

// 图表实例
let statusChart: echarts.ECharts | null = null
let amountChart: echarts.ECharts | null = null
let riskChart: echarts.ECharts | null = null
let progressChart: echarts.ECharts | null = null

const detail = ref<{ application_info: LoanApplication }>({} as { application_info: LoanApplication })

// 审批表单
const reviewForm = reactive({
  decision: '',
  approved_amount: 0,
  comments: '',
  required_info_details: ''
})

const reviewRules: FormRules = {
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
const fetchDetail = async () => {
  try {
    loading.value = true
    const applicationId = route.params.id as string
    const response = await adminLoanApprovalApi.getDetail(applicationId)
    detail.value = { application_info: response.application_info }
    
    // 初始化审批表单
    if (canReview(response.application_info.status)) {
      reviewForm.approved_amount = response.application_info.amount
    }
  } catch (error) {
    ElMessage.error('获取申请详情失败')
  } finally {
    loading.value = false
  }
}

const refreshData = () => {
  fetchDetail()
}

const goBack = () => {
  router.back()
}

const canReview = (status: string) => {
  return status === '待人工复核'
}

const submitApprovalReview = async () => {
  if (!reviewFormRef.value) return
  
  try {
    await reviewFormRef.value.validate()
    submitting.value = true
    
    const submitData = {
      decision: reviewForm.decision as 'approved' | 'rejected' | 'require_more_info',
      approved_amount: reviewForm.decision === 'approved' ? reviewForm.approved_amount : undefined,
      comments: reviewForm.comments,
      required_info_details: reviewForm.decision === 'require_more_info' ? reviewForm.required_info_details : undefined
    }
    
    await adminLoanApprovalApi.reviewApproval(detail.value.application_info.application_id, submitData)
    
    ElMessage.success('审批提交成功')
    fetchDetail() // 重新获取详情
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

const getTimelineType = (status: string) => {
  const typeMap: Record<string, string> = {
    '已提交': 'primary',
    'AI_审批中': 'warning',
    '待人工复核': 'info',
    '已批准': 'success',
    '已拒绝': 'danger'
  }
  return typeMap[status] || 'primary'
}

const getRiskColor = (score: number) => {
  if (score <= 30) return '#67c23a'
  if (score <= 70) return '#e6a23c'
  return '#f56c6c'
}

const getRiskLevel = (score: number) => {
  if (score <= 30) return '低风险'
  if (score <= 70) return '中风险'
  return '高风险'
}

const getRiskDescription = (score: number) => {
  if (score <= 30) return '风险较低，可以考虑批准'
  if (score <= 70) return '风险中等，需要谨慎考虑'
  return '风险较高，建议拒绝'
}

const getRiskTagType = (score: number) => {
  if (score <= 30) return 'success'
  if (score <= 70) return 'warning'
  return 'danger'
}

const getStatusClass = (status: string) => {
  const statusMap: Record<string, string> = {
    'AI_审批中': 'warning',
    '待人工复核': 'info',
    '已批准': 'success',
    '已拒绝': 'danger'
  }
  return statusMap[status] || 'info'
}

const getProgressStep = (status: string) => {
  const stepMap: Record<string, number> = {
    '已提交': 0,
    'AI_审批中': 1,
    '待人工复核': 2,
    '已批准': 3,
    '已拒绝': 3
  }
  return stepMap[status] || 0
}

const getTimelineIcon = (status: string) => {
  const iconMap: Record<string, string> = {
    '已提交': 'el-icon-document',
    'AI_审批中': 'el-icon-cpu',
    '待人工复核': 'el-icon-user',
    '已批准': 'el-icon-check',
    '已拒绝': 'el-icon-close'
  }
  return iconMap[status] || 'el-icon-document'
}

const getDocTypeName = (docType: string) => {
  const typeMap: Record<string, string> = {
    'id_card_front': '身份证正面',
    'id_card_back': '身份证背面',
    'land_contract': '土地承包合同',
    'income_proof': '收入证明',
    'asset_proof': '资产证明'
  }
  return typeMap[docType] || docType
}

const maskIdCard = (idCard?: string) => {
  if (!idCard) return '-'
  return idCard.replace(/(.{4}).*(.{4})/, '$1****$2')
}

const formatAmount = (amount: number) => {
  return amount.toLocaleString()
}

const formatDateTime = (datetime: string | number | undefined) => {
  if (!datetime) return '-'
  // 处理Unix时间戳
  if (typeof datetime === 'number') {
    return dayjs(datetime * 1000).format('YYYY-MM-DD HH:mm:ss')
  }
  // 处理ISO字符串
  return dayjs(datetime).format('YYYY-MM-DD HH:mm:ss')
}

const previewDocument = (doc: any) => {
  // 打开文档预览
  window.open(doc.file_url, '_blank')
}

const downloadDocument = (doc: any) => {
  // 下载文档
  const link = document.createElement('a')
  link.href = doc.file_url
  link.download = getDocTypeName(doc.doc_type)
  link.click()
}

const exportDetail = () => {
  // 实现导出详情功能
  console.log('导出详情')
}

const getRiskClass = (score?: number) => {
  if (!score) return 'info'
  if (score <= 30) return 'success'
  if (score <= 70) return 'warning'
  return 'danger'
}

// 初始化图表
const initCharts = () => {
  // 状态图表
  if (statusChartRef.value) {
    statusChart = echarts.init(statusChartRef.value)
    const option = {
      series: [{
        type: 'pie',
        radius: ['50%', '80%'],
        data: [
          { value: 1, name: '当前状态', itemStyle: { color: '#409eff' } },
          { value: 0, name: '其他', itemStyle: { color: '#f0f2f5' } }
        ],
        label: { show: false },
        emphasis: { label: { show: false } }
      }]
    }
    statusChart.setOption(option)
  }

  // 金额图表
  if (amountChartRef.value) {
    amountChart = echarts.init(amountChartRef.value)
    const option = {
      series: [{
        type: 'bar',
        data: [detail.value.application_info?.amount || 0],
        itemStyle: {
          color: new echarts.graphic.LinearGradient(0, 0, 0, 1, [
            { offset: 0, color: '#67c23a' },
            { offset: 1, color: '#85ce61' }
          ])
        },
        barWidth: '60%'
      }],
      xAxis: { show: false },
      yAxis: { show: false },
      grid: { left: 0, right: 0, top: 0, bottom: 0 }
    }
    amountChart.setOption(option)
  }

  // 风险图表
  if (riskChartRef.value) {
    riskChart = echarts.init(riskChartRef.value)
    const riskScore = 50 // 固定风险评分
    const option = {
      series: [{
        type: 'gauge',
        radius: '80%',
        data: [{ value: riskScore, name: '风险评分' }],
        detail: { show: false },
        title: { show: false },
        axisLine: {
          lineStyle: {
            width: 8,
            color: [[0.3, '#67c23a'], [0.7, '#e6a23c'], [1, '#f56c6c']]
          }
        },
        pointer: { show: false },
        axisTick: { show: false },
        axisLabel: { show: false },
        splitLine: { show: false }
      }]
    }
    riskChart.setOption(option)
  }

  // 进度图表
  if (progressChartRef.value) {
    progressChart = echarts.init(progressChartRef.value)
    const progress = getProgressStep(detail.value.application_info?.status) * 25
    const option = {
      series: [{
        type: 'pie',
        radius: ['40%', '70%'],
        data: [
          { value: progress, name: '已完成', itemStyle: { color: '#409eff' } },
          { value: 100 - progress, name: '未完成', itemStyle: { color: '#f0f2f5' } }
        ],
        label: { show: false },
        emphasis: { label: { show: false } }
      }]
    }
    progressChart.setOption(option)
  }
}

const resizeCharts = () => {
  statusChart?.resize()
  amountChart?.resize()
  riskChart?.resize()
  progressChart?.resize()
}

onMounted(() => {
  fetchDetail()
  // 等待DOM更新后初始化图表
  setTimeout(() => {
    initCharts()
    window.addEventListener('resize', resizeCharts)
  }, 200)
})

onBeforeUnmount(() => {
  window.removeEventListener('resize', resizeCharts)
  statusChart?.dispose()
  amountChart?.dispose()
  riskChart?.dispose()
  progressChart?.dispose()
})
</script>

<style scoped>
.approval-detail {
  padding: 0;
}

.page-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 20px;
}

.header-left {
  display: flex;
  align-items: center;
  gap: 16px;
}

.page-title {
  margin: 0;
  color: #333;
  font-size: 24px;
  font-weight: 600;
}

.info-card,
.ai-card,
.documents-card,
.review-card,
.history-card {
  margin-bottom: 20px;
  border-radius: 8px;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  font-weight: 600;
  gap: 8px;
}

.amount {
  font-weight: 600;
  color: #f56c6c;
}

.amount.approved {
  color: #67c23a;
}

.ai-analysis {
  display: flex;
  flex-direction: column;
  gap: 24px;
}

.analysis-section h4 {
  margin: 0 0 12px 0;
  color: #333;
  font-size: 16px;
  font-weight: 600;
}

.risk-score-display {
  display: flex;
  align-items: center;
  gap: 20px;
}

.risk-score-text {
  font-size: 18px;
  font-weight: 600;
  color: #333;
}

.risk-info {
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.risk-level {
  font-size: 16px;
  font-weight: 500;
  color: #666;
}

.risk-description {
  font-size: 14px;
  color: #999;
}

.risk-factors {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
}

.risk-factor-tag {
  margin: 0;
}

.verification-results {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.verification-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 8px 12px;
  background: #f8f9fa;
  border-radius: 6px;
}

.verification-label {
  color: #333;
  font-size: 14px;
}

.documents-list {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.document-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 12px;
  border: 1px solid #ebeef5;
  border-radius: 6px;
  background: #fafafa;
}

.doc-info {
  display: flex;
  align-items: center;
  gap: 12px;
}

.doc-icon {
  font-size: 20px;
  color: #409eff;
}

.doc-details {
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.doc-type {
  font-weight: 500;
  color: #333;
}

.doc-ocr {
  font-size: 12px;
  color: #666;
}

.doc-meta {
  font-size: 12px;
  color: #999;
}

.doc-actions {
  display: flex;
  gap: 8px;
}

.decision-radio {
  display: flex;
  align-items: center;
  gap: 4px;
  width: 100%;
  margin-bottom: 12px;
}

.amount-hint {
  font-size: 12px;
  color: #999;
  margin-top: 4px;
}

.timeline-content {
  background: #fff;
  padding: 8px 12px;
  border-radius: 6px;
  border: 1px solid #ebeef5;
}

.timeline-status {
  font-weight: 500;
  color: #333;
  margin-bottom: 4px;
}

.timeline-operator {
  font-size: 12px;
  color: #666;
}

.timeline-comments {
  font-size: 12px;
  color: #999;
}

:deep(.el-descriptions__label) {
  font-weight: 500;
}

:deep(.el-radio) {
  margin-right: 0;
  margin-bottom: 12px;
}

:deep(.el-timeline-item__content) {
  padding-bottom: 12px;
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

.stat-value.amount {
  color: #f56c6c;
}

.stat-value.success {
  color: #67c23a;
}

.stat-value.warning {
  color: #e6a23c;
}

.stat-value.danger {
  color: #f56c6c;
}

.stat-value.info {
  color: #409eff;
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

.mini-chart {
  height: 100%;
  width: 100%;
}

.header-tags {
  display: flex;
  gap: 8px;
  align-items: center;
}

.progress-steps {
  padding: 20px 0;
}

.progress-card {
  margin-bottom: 20px;
  border-radius: 8px;
}

.risk-info {
  display: flex;
  flex-direction: column;
  gap: 8px;
  margin-left: 20px;
}

.risk-description {
  font-size: 14px;
  color: #999;
  line-height: 1.4;
}

.doc-meta {
  font-size: 12px;
  color: #999;
  margin-top: 4px;
}

.timeline-comments {
  font-size: 12px;
  color: #999;
  margin-top: 4px;
  padding: 4px 8px;
  background: #f5f7fa;
  border-radius: 4px;
}

.analysis-section h4 {
  margin: 0 0 16px 0;
  color: #333;
  font-size: 16px;
  font-weight: 600;
  display: flex;
  align-items: center;
  gap: 8px;
}

.verification-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 12px 16px;
  background: #f8f9fa;
  border-radius: 8px;
  border: 1px solid #ebeef5;
  transition: all 0.3s ease;
}

.verification-item:hover {
  background: #f0f2f5;
  border-color: #d3d6db;
}

.document-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 16px;
  border: 1px solid #ebeef5;
  border-radius: 8px;
  background: #fafafa;
  transition: all 0.3s ease;
}

.document-item:hover {
  background: #f5f7fa;
  border-color: #c0c4cc;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
}

.doc-details {
  display: flex;
  flex-direction: column;
  gap: 6px;
}

.doc-type {
  font-weight: 500;
  color: #333;
  font-size: 14px;
}

.doc-ocr {
  font-size: 12px;
  color: #666;
  line-height: 1.3;
}

.decision-radio {
  display: flex;
  align-items: center;
  gap: 8px;
  width: 100%;
  margin-bottom: 16px;
  padding: 12px;
  border: 1px solid #ebeef5;
  border-radius: 8px;
  transition: all 0.3s ease;
}

.decision-radio:hover {
  background: #f5f7fa;
  border-color: #c0c4cc;
}

.amount-hint {
  font-size: 12px;
  color: #999;
  margin-top: 8px;
  text-align: left;
}

/* 响应式设计 */
@media (max-width: 1200px) {
  .stats-row {
    grid-template-columns: repeat(2, 1fr);
  }
}

@media (max-width: 768px) {
  .stats-row {
    grid-template-columns: 1fr;
  }
  
  .stat-content {
    padding: 20px;
  }
  
  .stat-value {
    font-size: 24px;
  }
}
</style> 