<template>
  <div class="loan-detail">
    <!-- 顶部面包屑 -->
    <el-breadcrumb class="breadcrumb" separator="/">
      <el-breadcrumb-item to="/approval/loan">贷款审批</el-breadcrumb-item>
      <el-breadcrumb-item>申请详情</el-breadcrumb-item>
    </el-breadcrumb>

    <div v-if="loading" class="loading-container">
      <el-skeleton :rows="10" animated />
    </div>

    <div v-else-if="loanDetail" class="detail-content">
      <!-- 申请基本信息 -->
      <el-card class="info-card" shadow="never">
        <template #header>
          <div class="card-header">
            <span>申请基本信息</span>
            <div class="header-actions">
              <el-tag :type="getStatusColor(loanDetail.status)" size="large">
                {{ getStatusName(loanDetail.status) }}
              </el-tag>
            </div>
          </div>
        </template>

        <el-descriptions :column="3" border>
          <el-descriptions-item label="申请编号">{{ loanDetail.application_number }}</el-descriptions-item>
          <el-descriptions-item label="申请人">{{ loanDetail.applicant_name }}</el-descriptions-item>
          <el-descriptions-item label="身份证号">{{ loanDetail.id_card }}</el-descriptions-item>
          <el-descriptions-item label="手机号">{{ loanDetail.phone }}</el-descriptions-item>
          <el-descriptions-item label="邮箱">{{ loanDetail.email || '未填写' }}</el-descriptions-item>
          <el-descriptions-item label="贷款类型">
            <el-tag :type="getLoanTypeColor(loanDetail.loan_type)">
              {{ getLoanTypeName(loanDetail.loan_type) }}
            </el-tag>
          </el-descriptions-item>
          <el-descriptions-item label="申请金额">
            <span class="amount">{{ formatAmount(loanDetail.amount) }}</span>
          </el-descriptions-item>
          <el-descriptions-item label="贷款期限">{{ loanDetail.term }}个月</el-descriptions-item>
          <el-descriptions-item label="利率">{{ (loanDetail.interest_rate * 100).toFixed(2) }}%</el-descriptions-item>
          <el-descriptions-item label="申请时间" :span="2">{{ formatDateTime(loanDetail.created_at) }}</el-descriptions-item>
          <el-descriptions-item label="贷款用途" :span="3">{{ loanDetail.purpose }}</el-descriptions-item>
        </el-descriptions>
      </el-card>

      <!-- 财务信息 -->
      <el-card class="info-card" shadow="never">
        <template #header>
          <span>财务信息</span>
        </template>

        <el-descriptions :column="2" border>
          <el-descriptions-item label="月收入">{{ formatAmount(loanDetail.monthly_income) }}</el-descriptions-item>
          <el-descriptions-item label="年收入">{{ formatAmount(loanDetail.annual_income) }}</el-descriptions-item>
          <el-descriptions-item label="现有负债">{{ formatAmount(loanDetail.existing_debt) }}</el-descriptions-item>
          <el-descriptions-item label="负债收入比">{{ ((loanDetail.existing_debt / loanDetail.annual_income) * 100).toFixed(1) }}%</el-descriptions-item>
          <el-descriptions-item label="资产总值">{{ formatAmount(loanDetail.total_assets) }}</el-descriptions-item>
          <el-descriptions-item label="工作年限">{{ loanDetail.work_experience }}年</el-descriptions-item>
          <el-descriptions-item label="工作单位" :span="2">{{ loanDetail.employer }}</el-descriptions-item>
        </el-descriptions>
      </el-card>

      <!-- 风险评估 -->
      <el-card class="info-card" shadow="never">
        <template #header>
          <span>风险评估</span>
        </template>

        <div class="risk-assessment">
          <el-row :gutter="20">
            <el-col :span="8">
              <div class="score-card">
                <div class="score-title">信用评分</div>
                <div class="score-value">{{ loanDetail.credit_score }}</div>
                <el-progress
                  :percentage="(loanDetail.credit_score / 850) * 100"
                  :color="getScoreColor(loanDetail.credit_score)"
                  :show-text="false"
                />
              </div>
            </el-col>
            <el-col :span="8">
              <div class="score-card">
                <div class="score-title">还款能力</div>
                <div class="score-value">{{ loanDetail.repayment_ability }}</div>
                <el-progress
                  :percentage="loanDetail.repayment_ability"
                  :color="getScoreColor(loanDetail.repayment_ability * 8.5)"
                  :show-text="false"
                />
              </div>
            </el-col>
            <el-col :span="8">
              <div class="score-card">
                <div class="score-title">综合评级</div>
                <div class="score-value">{{ loanDetail.overall_rating }}</div>
                <el-tag :type="getRatingColor(loanDetail.overall_rating)" size="large">
                  {{ loanDetail.overall_rating }}
                </el-tag>
              </div>
            </el-col>
          </el-row>

          <div class="risk-factors">
            <h4>风险因素分析</h4>
            <el-table :data="loanDetail.risk_factors" size="small">
              <el-table-column prop="factor" label="风险因素" width="200" />
              <el-table-column prop="level" label="风险等级" width="120">
                <template #default="{ row }">
                  <el-tag :type="getRiskLevelColor(row.level)" size="small">
                    {{ row.level }}
                  </el-tag>
                </template>
              </el-table-column>
              <el-table-column prop="description" label="说明" />
              <el-table-column prop="impact" label="影响程度" width="100">
                <template #default="{ row }">
                  <el-rate v-model="row.impact" disabled size="small" />
                </template>
              </el-table-column>
            </el-table>
          </div>
        </div>
      </el-card>

      <!-- 审批历史 -->
      <el-card class="info-card" shadow="never">
        <template #header>
          <span>审批历史</span>
        </template>

        <el-timeline>
          <el-timeline-item
            v-for="record in loanDetail.approval_history"
            :key="record.id"
            :timestamp="formatDateTime(record.created_at)"
            :type="getTimelineType(record.action)"
          >
            <div class="timeline-content">
              <div class="timeline-header">
                <span class="action">{{ record.action }}</span>
                <span class="operator">{{ record.operator }}</span>
              </div>
              <div class="timeline-comment" v-if="record.comment">
                {{ record.comment }}
              </div>
            </div>
          </el-timeline-item>
        </el-timeline>
      </el-card>

      <!-- 附件材料 -->
      <el-card class="info-card" shadow="never">
        <template #header>
          <span>申请材料</span>
        </template>

        <div class="attachments">
          <el-row :gutter="16">
            <el-col :span="6" v-for="attachment in loanDetail.attachments" :key="attachment.id">
              <div class="attachment-item">
                <div class="attachment-icon">
                  <el-icon><Document /></el-icon>
                </div>
                <div class="attachment-info">
                  <div class="attachment-name">{{ attachment.name }}</div>
                  <div class="attachment-size">{{ formatFileSize(attachment.size) }}</div>
                </div>
                <div class="attachment-actions">
                  <el-button size="small" @click="previewFile(attachment)">预览</el-button>
                  <el-button size="small" @click="downloadFile(attachment)">下载</el-button>
                </div>
              </div>
            </el-col>
          </el-row>
        </div>
      </el-card>

      <!-- 审批操作 -->
      <el-card v-if="loanDetail.status === 'pending'" class="action-card" shadow="never">
        <template #header>
          <span>审批操作</span>
        </template>

        <div class="approval-actions">
          <el-form :model="approvalForm" label-width="100px">
            <el-form-item label="审批意见">
              <el-input
                v-model="approvalForm.comment"
                type="textarea"
                :rows="4"
                placeholder="请输入审批意见"
              />
            </el-form-item>
            <el-form-item label="批准金额" v-if="approvalForm.action === 'approve'">
              <el-input-number
                v-model="approvalForm.approved_amount"
                :min="0"
                :max="loanDetail.amount"
                :precision="2"
                style="width: 200px"
              />
            </el-form-item>
            <el-form-item label="调整利率" v-if="approvalForm.action === 'approve'">
              <el-input-number
                v-model="approvalForm.adjusted_rate"
                :min="0"
                :max="1"
                :precision="4"
                :step="0.0001"
                style="width: 200px"
              />
            </el-form-item>
          </el-form>

          <div class="action-buttons">
            <el-button 
              type="success" 
              size="large" 
              @click="handleApprove"
              :loading="submitting"
            >
              <el-icon><Check /></el-icon>
              通过申请
            </el-button>
            <el-button 
              type="danger" 
              size="large" 
              @click="handleReject"
              :loading="submitting"
            >
              <el-icon><Close /></el-icon>
              拒绝申请
            </el-button>
            <el-button 
              type="warning" 
              size="large" 
              @click="handleReturn"
              :loading="submitting"
            >
              <el-icon><Back /></el-icon>
              退回修改
            </el-button>
          </div>
        </div>
      </el-card>
    </div>

    <div v-else class="error-state">
      <el-result
        icon="warning"
        title="申请未找到"
        sub-title="请检查申请ID是否正确"
      >
        <template #extra>
          <el-button type="primary" @click="$router.back()">返回列表</el-button>
        </template>
      </el-result>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { Document, Check, Close, Back } from '@element-plus/icons-vue'
import { ElMessage, ElMessageBox } from 'element-plus'

const route = useRoute()
const router = useRouter()

const loading = ref(true)
const submitting = ref(false)
const loanDetail = ref<any>(null)

// 审批表单
const approvalForm = reactive({
  action: '',
  comment: '',
  approved_amount: 0,
  adjusted_rate: 0
})

// 模拟详情数据
const mockLoanDetail = {
  id: route.params.id,
  application_number: 'LN202412001',
  applicant_name: '张三',
  id_card: '320101199001011234',
  phone: '13800138001',
  email: 'zhangsan@example.com',
  loan_type: 'personal_credit',
  amount: 500000,
  term: 36,
  interest_rate: 0.065,
  purpose: '用于扩大农业种植规模，购买种子、化肥等农资',
  status: 'pending',
  created_at: new Date('2024-12-01 09:30:00'),
  monthly_income: 25000,
  annual_income: 300000,
  existing_debt: 80000,
  total_assets: 1200000,
  work_experience: 8,
  employer: '张家农业合作社',
  credit_score: 720,
  repayment_ability: 85,
  overall_rating: 'A',
  risk_factors: [
    {
      factor: '收入稳定性',
      level: '低风险',
      description: '农业收入较为稳定，有合作社保障',
      impact: 2
    },
    {
      factor: '负债比率',
      level: '中风险',
      description: '负债收入比26.7%，处于可接受范围',
      impact: 3
    },
    {
      factor: '担保情况',
      level: '低风险',
      description: '有农业资产作为担保',
      impact: 2
    }
  ],
  approval_history: [
    {
      id: 1,
      action: '申请提交',
      operator: '张三',
      comment: '提交贷款申请',
      created_at: new Date('2024-12-01 09:30:00')
    },
    {
      id: 2,
      action: '初审通过',
      operator: '李审核',
      comment: '基本资料齐全，进入风险评估',
      created_at: new Date('2024-12-01 14:20:00')
    }
  ],
  attachments: [
    {
      id: 1,
      name: '身份证正面.jpg',
      size: 2048576,
      type: 'image/jpeg'
    },
    {
      id: 2,
      name: '收入证明.pdf',
      size: 1536000,
      type: 'application/pdf'
    },
    {
      id: 3,
      name: '土地承包合同.pdf',
      size: 3072000,
      type: 'application/pdf'
    }
  ]
}

// 格式化金额
const formatAmount = (amount: number) => {
  if (amount >= 10000) {
    return `¥${(amount / 10000).toFixed(1)}万`
  }
  return `¥${amount.toLocaleString()}`
}

// 格式化日期时间
const formatDateTime = (date: Date) => {
  return date.toLocaleString('zh-CN', {
    year: 'numeric',
    month: '2-digit',
    day: '2-digit',
    hour: '2-digit',
    minute: '2-digit'
  })
}

// 格式化文件大小
const formatFileSize = (bytes: number) => {
  if (bytes === 0) return '0 B'
  const k = 1024
  const sizes = ['B', 'KB', 'MB', 'GB']
  const i = Math.floor(Math.log(bytes) / Math.log(k))
  return parseFloat((bytes / Math.pow(k, i)).toFixed(2)) + ' ' + sizes[i]
}

// 获取贷款类型名称
const getLoanTypeName = (type: string) => {
  const types: Record<string, string> = {
    personal_credit: '个人信用贷',
    business_loan: '企业经营贷',
    mortgage: '抵押贷款',
    vehicle: '车辆贷款'
  }
  return types[type] || type
}

// 获取贷款类型颜色
const getLoanTypeColor = (type: string) => {
  const colors: Record<string, string> = {
    personal_credit: 'primary',
    business_loan: 'success',
    mortgage: 'warning',
    vehicle: 'info'
  }
  return colors[type] || 'default'
}

// 获取状态名称
const getStatusName = (status: string) => {
  const statuses: Record<string, string> = {
    pending: '待审批',
    approved: '已通过',
    rejected: '已拒绝',
    returned: '已退回'
  }
  return statuses[status] || status
}

// 获取状态颜色
const getStatusColor = (status: string) => {
  const colors: Record<string, string> = {
    pending: 'warning',
    approved: 'success',
    rejected: 'danger',
    returned: 'info'
  }
  return colors[status] || 'default'
}

// 获取评分颜色
const getScoreColor = (score: number) => {
  if (score >= 700) return '#67C23A'
  if (score >= 600) return '#E6A23C'
  return '#F56C6C'
}

// 获取评级颜色
const getRatingColor = (rating: string) => {
  const colors: Record<string, string> = {
    'AAA': 'success',
    'AA': 'success',
    'A': 'primary',
    'BBB': 'warning',
    'BB': 'warning',
    'B': 'danger'
  }
  return colors[rating] || 'info'
}

// 获取风险等级颜色
const getRiskLevelColor = (level: string) => {
  const colors: Record<string, string> = {
    '低风险': 'success',
    '中风险': 'warning',
    '高风险': 'danger'
  }
  return colors[level] || 'info'
}

// 获取时间线类型
const getTimelineType = (action: string) => {
  if (action.includes('通过')) return 'success'
  if (action.includes('拒绝')) return 'danger'
  if (action.includes('退回')) return 'warning'
  return 'primary'
}

// 加载详情数据
const loadDetail = async () => {
  loading.value = true
  try {
    // 模拟API调用
    await new Promise(resolve => setTimeout(resolve, 1000))
    loanDetail.value = mockLoanDetail
    approvalForm.approved_amount = loanDetail.value.amount
    approvalForm.adjusted_rate = loanDetail.value.interest_rate
  } catch (error) {
    ElMessage.error('加载申请详情失败')
  } finally {
    loading.value = false
  }
}

// 审批操作
const handleApprove = () => {
  approvalForm.action = 'approve'
  ElMessageBox.confirm('确定要通过这个申请吗？', '确认审批', {
    confirmButtonText: '确定',
    cancelButtonText: '取消',
    type: 'success'
  }).then(() => {
    submitApproval()
  })
}

const handleReject = () => {
  approvalForm.action = 'reject'
  if (!approvalForm.comment.trim()) {
    ElMessage.warning('拒绝申请必须填写审批意见')
    return
  }
  ElMessageBox.confirm('确定要拒绝这个申请吗？', '确认审批', {
    confirmButtonText: '确定',
    cancelButtonText: '取消',
    type: 'warning'
  }).then(() => {
    submitApproval()
  })
}

const handleReturn = () => {
  approvalForm.action = 'return'
  if (!approvalForm.comment.trim()) {
    ElMessage.warning('退回申请必须填写退回原因')
    return
  }
  ElMessageBox.confirm('确定要退回这个申请吗？', '确认操作', {
    confirmButtonText: '确定',
    cancelButtonText: '取消',
    type: 'warning'
  }).then(() => {
    submitApproval()
  })
}

// 提交审批
const submitApproval = async () => {
  submitting.value = true
  try {
    await new Promise(resolve => setTimeout(resolve, 1000))
    ElMessage.success('审批操作成功')
    router.push('/approval/loan')
  } catch (error) {
    ElMessage.error('审批操作失败')
  } finally {
    submitting.value = false
  }
}

// 文件操作
const previewFile = (attachment: any) => {
  ElMessage.info(`预览文件：${attachment.name}`)
}

const downloadFile = (attachment: any) => {
  ElMessage.info(`下载文件：${attachment.name}`)
}

onMounted(() => {
  loadDetail()
})
</script>

<style scoped>
.loan-detail {
  padding: 20px;
  background: #f5f5f5;
  min-height: calc(100vh - 60px);
}

.breadcrumb {
  margin-bottom: 20px;
}

.loading-container {
  padding: 40px;
}

.detail-content {
  display: flex;
  flex-direction: column;
  gap: 20px;
}

.info-card, .action-card {
  border: none;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  font-weight: 500;
}

.header-actions {
  display: flex;
  gap: 8px;
}

.amount {
  font-weight: 600;
  color: #E6A23C;
  font-size: 16px;
}

.risk-assessment {
  padding: 20px 0;
}

.score-card {
  text-align: center;
  padding: 20px;
  background: #f8f9fa;
  border-radius: 8px;
}

.score-title {
  font-size: 14px;
  color: #909399;
  margin-bottom: 8px;
}

.score-value {
  font-size: 28px;
  font-weight: bold;
  color: #303133;
  margin-bottom: 12px;
}

.risk-factors {
  margin-top: 30px;
}

.risk-factors h4 {
  margin-bottom: 16px;
  color: #303133;
}

.timeline-content {
  padding: 0 0 12px 0;
}

.timeline-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 4px;
}

.action {
  font-weight: 500;
  color: #303133;
}

.operator {
  font-size: 12px;
  color: #909399;
}

.timeline-comment {
  color: #606266;
  font-size: 14px;
}

.attachments {
  padding: 10px 0;
}

.attachment-item {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 16px;
  border: 1px solid #EBEEF5;
  border-radius: 8px;
  margin-bottom: 12px;
}

.attachment-icon {
  width: 40px;
  height: 40px;
  display: flex;
  align-items: center;
  justify-content: center;
  background: #f0f9ff;
  border-radius: 6px;
  color: #409EFF;
  font-size: 20px;
}

.attachment-info {
  flex: 1;
}

.attachment-name {
  font-weight: 500;
  color: #303133;
  margin-bottom: 4px;
}

.attachment-size {
  font-size: 12px;
  color: #909399;
}

.attachment-actions {
  display: flex;
  gap: 8px;
}

.approval-actions {
  padding: 20px 0;
}

.action-buttons {
  display: flex;
  gap: 12px;
  justify-content: center;
  margin-top: 20px;
}

.error-state {
  display: flex;
  justify-content: center;
  align-items: center;
  min-height: 400px;
}

@media (max-width: 768px) {
  .loan-detail {
    padding: 10px;
  }
  
  .card-header {
    flex-direction: column;
    gap: 12px;
    align-items: stretch;
  }
  
  .action-buttons {
    flex-direction: column;
  }
  
  .attachment-item {
    flex-direction: column;
    text-align: center;
  }
}
</style> 