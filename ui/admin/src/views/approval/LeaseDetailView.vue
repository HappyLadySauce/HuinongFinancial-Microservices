<template>
  <div class="lease-detail">
    <!-- 页面头部 -->
    <PageHeader
      title="租赁审批详情"
      :show-back="true"
      :breadcrumbs="breadcrumbs"
      @back="handleBack"
    >
      <template #actions>
        <el-button @click="handleExport">
          <el-icon><Download /></el-icon>
          导出详情
        </el-button>
      </template>
    </PageHeader>

    <div class="detail-content">
      <!-- 审批进度 -->
      <el-card class="progress-card" shadow="never">
        <template #header>
          <div class="card-header">
            <span>
              <el-icon><Operation /></el-icon>
              审批进度
            </span>
            <StatusTag :status="detail.status" :show-icon="true" size="large" />
          </div>
        </template>

        <el-steps :active="getProgressStep(detail.status)" align-center>
          <el-step title="提交申请" :icon="Document" />
          <el-step title="初步审核" :icon="View" />
          <el-step title="风险评估" :icon="Warning" />
          <el-step title="最终审批" :icon="Select" />
          <el-step 
            :title="detail.status === 'approved' ? '审批通过' : detail.status === 'rejected' ? '审批拒绝' : '待完成'" 
            :icon="detail.status === 'approved' ? CircleCheck : detail.status === 'rejected' ? CircleClose : Clock"
          />
        </el-steps>
      </el-card>

      <el-row :gutter="24">
        <!-- 左侧详情 -->
        <el-col :span="16">
          <!-- 基本信息 -->
          <el-card class="info-card" shadow="never">
            <template #header>
              <div class="card-header">
                <span>
                  <el-icon><User /></el-icon>
                  租赁申请信息
                </span>
              </div>
            </template>
            
            <el-descriptions :column="2" border>
              <el-descriptions-item label="申请编号">
                <el-text type="primary" tag="b">{{ detail.application_id }}</el-text>
              </el-descriptions-item>
              <el-descriptions-item label="申请人">
                {{ detail.applicant_details?.real_name }}
              </el-descriptions-item>
              <el-descriptions-item label="身份证号">
                {{ maskIdCard(detail.applicant_details?.id_card_number) }}
              </el-descriptions-item>
              <el-descriptions-item label="联系电话">
                {{ detail.applicant_details?.phone || '-' }}
              </el-descriptions-item>
              <el-descriptions-item label="联系地址" span="2">
                {{ detail.applicant_details?.address }}
              </el-descriptions-item>
              <el-descriptions-item label="租赁设备">
                {{ detail.equipment_type }}
              </el-descriptions-item>
              <el-descriptions-item label="设备数量">
                {{ detail.equipment_quantity }}台
              </el-descriptions-item>
              <el-descriptions-item label="租赁金额">
                <span class="amount">¥{{ formatAmount(detail.lease_amount) }}</span>
              </el-descriptions-item>
              <el-descriptions-item label="租赁期限">
                {{ detail.lease_term }}个月
              </el-descriptions-item>
              <el-descriptions-item label="月租金">
                <span class="amount">¥{{ formatAmount(detail.monthly_rent) }}</span>
              </el-descriptions-item>
              <el-descriptions-item label="保证金">
                <span class="amount">¥{{ formatAmount(detail.security_deposit) }}</span>
              </el-descriptions-item>
              <el-descriptions-item label="租赁用途" span="2">
                {{ detail.lease_purpose }}
              </el-descriptions-item>
              <el-descriptions-item label="提交时间">
                {{ formatDateTime(detail.submitted_at) }}
              </el-descriptions-item>
              <el-descriptions-item label="更新时间">
                {{ formatDateTime(detail.updated_at) }}
              </el-descriptions-item>
              <el-descriptions-item v-if="detail.approved_amount" label="批准金额" span="2">
                <span class="amount approved">¥{{ formatAmount(detail.approved_amount) }}</span>
              </el-descriptions-item>
            </el-descriptions>
          </el-card>

          <!-- AI分析报告 -->
          <el-card v-if="detail.ai_analysis_report" class="ai-card" shadow="never">
            <template #header>
              <div class="card-header">
                <span>
                  <el-icon><Cpu /></el-icon>
                  AI智能分析报告
                </span>
                <div class="header-tags">
                  <el-tag type="info" size="small">
                    风险评分: {{ detail.ai_analysis_report.overall_risk_score }}分
                  </el-tag>
                  <el-tag :type="getRiskTagType(detail.ai_analysis_report.overall_risk_score)" size="small">
                    {{ getRiskLevel(detail.ai_analysis_report.overall_risk_score) }}
                  </el-tag>
                </div>
              </div>
            </template>
            
            <div class="ai-analysis">
              <!-- 风险评分 -->
              <div class="risk-scores">
                <div class="score-item">
                  <div class="score-label">设备价值评估</div>
                  <div class="score-value">{{ detail.ai_analysis_report.equipment_value_score }}分</div>
                  <el-progress 
                    :percentage="detail.ai_analysis_report.equipment_value_score" 
                    :color="getScoreColor(detail.ai_analysis_report.equipment_value_score)"
                    :show-text="false"
                  />
                </div>
                <div class="score-item">
                  <div class="score-label">信用评分</div>
                  <div class="score-value">{{ detail.ai_analysis_report.credit_score }}分</div>
                  <el-progress 
                    :percentage="detail.ai_analysis_report.credit_score" 
                    :color="getScoreColor(detail.ai_analysis_report.credit_score)"
                    :show-text="false"
                  />
                </div>
                <div class="score-item">
                  <div class="score-label">还款能力</div>
                  <div class="score-value">{{ detail.ai_analysis_report.repayment_ability_score }}分</div>
                  <el-progress 
                    :percentage="detail.ai_analysis_report.repayment_ability_score" 
                    :color="getScoreColor(detail.ai_analysis_report.repayment_ability_score)"
                    :show-text="false"
                  />
                </div>
              </div>

              <!-- AI建议 -->
              <div class="ai-recommendation">
                <h4>AI分析建议</h4>
                <p>{{ detail.ai_analysis_report.recommendation }}</p>
              </div>
            </div>
          </el-card>

          <!-- 设备信息 -->
          <el-card class="equipment-card" shadow="never">
            <template #header>
              <div class="card-header">
                <span>
                  <el-icon><Box /></el-icon>
                  设备详情
                </span>
              </div>
            </template>

            <el-table :data="detail.equipment_list" border>
              <el-table-column prop="name" label="设备名称" />
              <el-table-column prop="model" label="设备型号" />
              <el-table-column prop="brand" label="品牌" />
              <el-table-column prop="quantity" label="数量" />
              <el-table-column prop="unit_price" label="单价">
                <template #default="{ row }">
                  ¥{{ formatAmount(row.unit_price) }}
                </template>
              </el-table-column>
              <el-table-column prop="total_value" label="总价值">
                <template #default="{ row }">
                  ¥{{ formatAmount(row.total_value) }}
                </template>
              </el-table-column>
              <el-table-column prop="depreciation_rate" label="折旧率">
                <template #default="{ row }">
                  {{ (row.depreciation_rate * 100).toFixed(1) }}%
                </template>
              </el-table-column>
            </el-table>
          </el-card>

          <!-- 审批历史 -->
          <el-card class="history-card" shadow="never">
            <template #header>
              <div class="card-header">
                <span>
                  <el-icon><Clock /></el-icon>
                  审批历史
                </span>
              </div>
            </template>

            <el-timeline>
              <el-timeline-item
                v-for="(record, index) in detail.approval_history"
                :key="index"
                :timestamp="formatDateTime(record.created_at)"
                :type="getTimelineType(record.action)"
                :icon="getTimelineIcon(record.action)"
                size="large"
              >
                <el-card>
                  <div class="timeline-content">
                    <div class="timeline-header">
                      <span class="action">{{ record.action }}</span>
                      <span class="operator">{{ record.operator_name }}</span>
                    </div>
                    <div v-if="record.comment" class="timeline-comment">
                      {{ record.comment }}
                    </div>
                    <div v-if="record.attachments?.length" class="timeline-attachments">
                      <el-tag 
                        v-for="attachment in record.attachments" 
                        :key="attachment.id"
                        size="small"
                        type="info"
                      >
                        {{ attachment.name }}
                      </el-tag>
                    </div>
                  </div>
                </el-card>
              </el-timeline-item>
            </el-timeline>
          </el-card>
        </el-col>

        <!-- 右侧操作面板 -->
        <el-col :span="8">
          <!-- 审批操作 -->
          <el-card v-if="detail.status === 'pending'" class="action-card" shadow="never">
            <template #header>
              <div class="card-header">
                <span>
                  <el-icon><Tools /></el-icon>
                  审批操作
                </span>
              </div>
            </template>

            <div class="action-buttons">
              <el-button 
                type="success" 
                size="large" 
                @click="handleApprove"
                style="width: 100%; margin-bottom: 12px;"
              >
                <el-icon><CircleCheck /></el-icon>
                通过申请
              </el-button>
              <el-button 
                type="danger" 
                size="large" 
                @click="handleReject"
                style="width: 100%;"
              >
                <el-icon><CircleClose /></el-icon>
                拒绝申请
              </el-button>
            </div>
          </el-card>

          <!-- 申请人信息 -->
          <el-card class="applicant-card" shadow="never">
            <template #header>
              <div class="card-header">
                <span>
                  <el-icon><Avatar /></el-icon>
                  申请人信息
                </span>
              </div>
            </template>

            <div class="applicant-info">
              <div class="info-item">
                <span class="label">姓名：</span>
                <span class="value">{{ detail.applicant_details?.real_name }}</span>
              </div>
              <div class="info-item">
                <span class="label">年龄：</span>
                <span class="value">{{ detail.applicant_details?.age }}岁</span>
              </div>
              <div class="info-item">
                <span class="label">学历：</span>
                <span class="value">{{ detail.applicant_details?.education }}</span>
              </div>
              <div class="info-item">
                <span class="label">职业：</span>
                <span class="value">{{ detail.applicant_details?.occupation }}</span>
              </div>
              <div class="info-item">
                <span class="label">月收入：</span>
                <span class="value amount">¥{{ formatAmount(detail.applicant_details?.monthly_income) }}</span>
              </div>
              <div class="info-item">
                <span class="label">信用评级：</span>
                <StatusTag :status="detail.applicant_details?.credit_rating" />
              </div>
            </div>
          </el-card>

          <!-- 附件文档 -->
          <el-card class="documents-card" shadow="never">
            <template #header>
              <div class="card-header">
                <span>
                  <el-icon><Document /></el-icon>
                  相关文档
                </span>
              </div>
            </template>

            <div class="document-list">
              <div 
                v-for="doc in detail.documents" 
                :key="doc.id" 
                class="document-item"
              >
                <div class="doc-info">
                  <div class="doc-name">{{ getDocTypeName(doc.doc_type) }}</div>
                  <div class="doc-size">{{ doc.file_size }}</div>
                </div>
                <div class="doc-actions">
                  <el-button size="small" @click="previewDocument(doc)">预览</el-button>
                  <el-button size="small" type="primary" @click="downloadDocument(doc)">下载</el-button>
                </div>
              </div>
            </div>
          </el-card>
        </el-col>
      </el-row>
    </div>

    <!-- 审批对话框 -->
    <FormDialog
      v-model="approvalDialog.visible"
      :title="approvalDialog.type === 'approve' ? '审批通过' : '审批拒绝'"
      :mode="'edit'"
      :initial-data="approvalDialog.data"
      :rules="approvalRules"
      @submit="handleApprovalSubmit"
      @cancel="approvalDialog.visible = false"
    >
      <template #form="{ form }">
        <el-form-item label="处理意见" prop="comment">
          <el-input
            v-model="form.comment"
            type="textarea"
            :rows="4"
            placeholder="请输入处理意见"
          />
        </el-form-item>
        
        <el-form-item v-if="approvalDialog.type === 'approve'" label="批准金额" prop="approvedAmount">
          <el-input-number
            v-model="form.approvedAmount"
            :min="0"
            :max="form.requestedAmount"
            :step="1000"
            style="width: 100%"
          />
        </el-form-item>
        
        <el-form-item v-if="approvalDialog.type === 'approve'" label="月租金" prop="monthlyRent">
          <el-input-number
            v-model="form.monthlyRent"
            :min="0"
            :step="100"
            style="width: 100%"
          />
        </el-form-item>
      </template>
    </FormDialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, computed, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import dayjs from 'dayjs'
import {
  Download,
  Operation, // 替换Process为Operation（操作图标）
  Document,
  View,
  Warning,
  Select,
  Clock,
  CircleCheck,
  CircleClose,
  User,
  Cpu,
  Box,
  Tools,
  Avatar
} from '@element-plus/icons-vue'

// 导入通用组件
import { PageHeader, StatusTag, FormDialog } from '@/components/common'

const route = useRoute()
const router = useRouter()

const loading = ref(false)

const detail = ref({
  id: route.params.id,
  application_id: 'LS202412290001',
  status: 'pending',
  equipment_type: '农业机械',
  equipment_quantity: 2,
  lease_amount: 150000,
  lease_term: 24,
  monthly_rent: 7500,
  security_deposit: 15000,
  lease_purpose: '农业生产作业',
  submitted_at: '2024-12-29T10:30:00',
  updated_at: '2024-12-29T15:20:00',
  approved_amount: null,
  applicant_details: {
    real_name: '张三',
    id_card_number: '123456789012345678',
    phone: '13888888888',
    address: '江苏省南京市玄武区某某街道123号',
    age: 35,
    education: '本科',
    occupation: '农民专业合作社负责人',
    monthly_income: 12000,
    credit_rating: 'good'
  },
  equipment_list: [
    {
      name: '拖拉机',
      model: 'TX-1504',
      brand: '东方红',
      quantity: 1,
      unit_price: 80000,
      total_value: 80000,
      depreciation_rate: 0.15
    },
    {
      name: '收割机',
      model: 'GF-2024',
      brand: '雷沃',
      quantity: 1,
      unit_price: 70000,
      total_value: 70000,
      depreciation_rate: 0.18
    }
  ],
  ai_analysis_report: {
    overall_risk_score: 75,
    equipment_value_score: 85,
    credit_score: 72,
    repayment_ability_score: 78,
    recommendation: '该申请人具备良好的设备运营能力和还款意愿，建议批准租赁申请。设备价值评估合理，月租金在申请人承受范围内。'
  },
  approval_history: [
    {
      action: '提交申请',
      operator_name: '张三',
      created_at: '2024-12-29T10:30:00',
      comment: '提交农业机械租赁申请'
    },
    {
      action: '初步审核',
      operator_name: '李审核员',
      created_at: '2024-12-29T11:15:00',
      comment: '材料齐全，通过初步审核'
    }
  ],
  documents: [
    {
      id: 1,
      doc_type: 'id_card_front',
      file_size: '2.1MB',
      file_url: '/api/files/id_card_front.jpg'
    },
    {
      id: 2,
      doc_type: 'business_license',
      file_size: '1.8MB',
      file_url: '/api/files/business_license.pdf'
    },
    {
      id: 3,
      doc_type: 'equipment_quote',
      file_size: '3.2MB',
      file_url: '/api/files/equipment_quote.pdf'
    }
  ]
})

const approvalDialog = reactive({
  visible: false,
  type: 'approve' as 'approve' | 'reject',
  data: {}
})

// 计算属性
const breadcrumbs = computed(() => [
  { title: '首页', to: '/' },
  { title: '核心审批', to: '/approval' },
  { title: '租赁审批', to: '/approval/lease' },
  { title: '审批详情' }
])

// 表单验证规则
const approvalRules = {
  comment: [
    { required: true, message: '请输入处理意见', trigger: 'blur' }
  ],
  approvedAmount: [
    { required: true, message: '请输入批准金额', trigger: 'blur' }
  ],
  monthlyRent: [
    { required: true, message: '请输入月租金', trigger: 'blur' }
  ]
}

// 方法
const handleBack = () => {
  router.push('/approval/lease')
}

const handleApprove = () => {
  approvalDialog.type = 'approve'
  approvalDialog.data = {
    id: detail.value.id,
    requestedAmount: detail.value.lease_amount,
    approvedAmount: detail.value.lease_amount,
    monthlyRent: detail.value.monthly_rent
  }
  approvalDialog.visible = true
}

const handleReject = () => {
  approvalDialog.type = 'reject'
  approvalDialog.data = { id: detail.value.id }
  approvalDialog.visible = true
}

const handleApprovalSubmit = async (data: any, mode: string) => {
  try {
    console.log('提交审批:', data)
    // 实现审批提交逻辑
    approvalDialog.visible = false
    // 更新详情数据
  } catch (error) {
    console.error('审批失败:', error)
  }
}

const handleExport = () => {
  console.log('导出详情')
}

const getProgressStep = (status: string) => {
  const stepMap: Record<string, number> = {
    'submitted': 0,
    'reviewing': 1,
    'risk_assessment': 2,
    'pending': 3,
    'approved': 4,
    'rejected': 4
  }
  return stepMap[status] || 0
}

const getRiskTagType = (score: number) => {
  if (score >= 80) return 'success'
  if (score >= 60) return 'warning'
  return 'danger'
}

const getRiskLevel = (score: number) => {
  if (score >= 80) return '低风险'
  if (score >= 60) return '中风险'
  return '高风险'
}

const getScoreColor = (score: number) => {
  if (score >= 80) return '#67c23a'
  if (score >= 60) return '#e6a23c'
  return '#f56c6c'
}

const getTimelineType = (action: string) => {
  const typeMap: Record<string, string> = {
    '提交申请': 'primary',
    '初步审核': 'success',
    '风险评估': 'warning',
    '审批通过': 'success',
    '审批拒绝': 'danger'
  }
  return typeMap[action] || 'info'
}

const getTimelineIcon = (action: string) => {
  const iconMap: Record<string, any> = {
    '提交申请': Document,
    '初步审核': View,
    '风险评估': Warning,
    '审批通过': CircleCheck,
    '审批拒绝': CircleClose
  }
  return iconMap[action] || Clock
}

const getDocTypeName = (docType: string) => {
  const typeMap: Record<string, string> = {
    'id_card_front': '身份证正面',
    'id_card_back': '身份证背面',
    'business_license': '营业执照',
    'equipment_quote': '设备报价单',
    'financial_statement': '财务报表',
    'bank_statement': '银行流水'
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

const formatDateTime = (datetime: string) => {
  return dayjs(datetime).format('YYYY-MM-DD HH:mm:ss')
}

const previewDocument = (doc: any) => {
  window.open(doc.file_url, '_blank')
}

const downloadDocument = (doc: any) => {
  const link = document.createElement('a')
  link.href = doc.file_url
  link.download = getDocTypeName(doc.doc_type)
  link.click()
}

onMounted(() => {
  // 加载详情数据
  console.log('加载租赁详情:', route.params.id)
})
</script>

<style scoped>
.lease-detail {
  height: 100vh;
  display: flex;
  flex-direction: column;
  gap: 16px;
  padding: 0 20px;
  overflow: hidden;
}

.detail-content {
  flex: 1;
  overflow-y: auto;
  padding-bottom: 20px;
}

.progress-card {
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

.header-tags {
  display: flex;
  gap: 8px;
}

.info-card,
.ai-card,
.equipment-card,
.history-card,
.action-card,
.applicant-card,
.documents-card {
  margin-bottom: 20px;
  border: none;
  border-radius: 12px;
}

.amount {
  font-weight: 600;
  color: #409eff;
}

.amount.approved {
  color: #67c23a;
}

.ai-analysis {
  padding: 16px 0;
}

.risk-scores {
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  gap: 20px;
  margin-bottom: 24px;
}

.score-item {
  text-align: center;
}

.score-label {
  font-size: 14px;
  color: #666;
  margin-bottom: 8px;
}

.score-value {
  font-size: 24px;
  font-weight: bold;
  margin-bottom: 8px;
}

.ai-recommendation {
  background: #f8f9fa;
  padding: 16px;
  border-radius: 8px;
}

.ai-recommendation h4 {
  margin: 0 0 8px 0;
  color: #333;
}

.ai-recommendation p {
  margin: 0;
  color: #666;
  line-height: 1.6;
}

.timeline-content {
  padding: 8px 0;
}

.timeline-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 8px;
}

.action {
  font-weight: 500;
  color: #333;
}

.operator {
  color: #666;
  font-size: 12px;
}

.timeline-comment {
  color: #666;
  margin-bottom: 8px;
}

.timeline-attachments {
  display: flex;
  gap: 8px;
}

.action-buttons {
  padding: 16px 0;
}

.applicant-info {
  padding: 16px 0;
}

.info-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 12px;
  padding: 8px 0;
  border-bottom: 1px solid #f0f0f0;
}

.info-item:last-child {
  border-bottom: none;
}

.label {
  color: #666;
  font-size: 14px;
}

.value {
  color: #333;
  font-weight: 500;
}

.document-list {
  padding: 16px 0;
}

.document-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 12px 0;
  border-bottom: 1px solid #f0f0f0;
}

.document-item:last-child {
  border-bottom: none;
}

.doc-info {
  flex: 1;
}

.doc-name {
  font-weight: 500;
  color: #333;
  margin-bottom: 4px;
}

.doc-size {
  font-size: 12px;
  color: #999;
}

.doc-actions {
  display: flex;
  gap: 8px;
}

:deep(.el-steps) {
  padding: 20px 0;
}

:deep(.el-descriptions__cell) {
  padding: 12px;
}

:deep(.el-card__header) {
  padding: 16px 20px;
  border-bottom: 1px solid #f0f0f0;
}

:deep(.el-card__body) {
  padding: 20px;
}
</style> 