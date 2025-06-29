<template>
  <div class="loan-approval">
    <!-- 页面头部 -->
    <PageHeader
      title="贷款审批"
      subtitle="管理和审批贷款申请"
      :icon="CreditCard"
      :show-stats="true"
      :stats="pageStats"
      :breadcrumbs="breadcrumbs"
      @stat-click="handleStatClick"
    >
      <template #actions>
        <el-button type="primary" @click="handleExport">
          <el-icon><Download /></el-icon>
          导出数据
        </el-button>
      </template>
    </PageHeader>

    <!-- 搜索和操作栏 -->
    <TableActions
      :selected-items="selectedApprovals"
      :total="pagination.total"
      :batch-actions="batchActions"
      @search="handleSearch"
      @reset="handleReset"
      @refresh="handleRefresh"
      @batch-action="handleBatchAction"
    >
      <template #filters="{ form }">
        <el-form-item label="状态">
          <el-select v-model="form.status" placeholder="选择状态" clearable style="width: 120px">
            <el-option label="全部" value="" />
            <el-option label="待审批" value="pending" />
            <el-option label="已通过" value="approved" />
            <el-option label="已拒绝" value="rejected" />
          </el-select>
        </el-form-item>
        
        <el-form-item label="贷款类型">
          <el-select v-model="form.type" placeholder="贷款类型" clearable style="width: 140px">
            <el-option label="全部类型" value="" />
            <el-option label="个人信用贷" value="personal_credit" />
            <el-option label="企业经营贷" value="business_loan" />
            <el-option label="抵押贷款" value="mortgage" />
            <el-option label="车辆贷款" value="vehicle" />
          </el-select>
        </el-form-item>

        <el-form-item label="申请时间">
          <el-date-picker
            v-model="form.dateRange"
            type="daterange"
            range-separator="至"
            start-placeholder="开始日期"
            end-placeholder="结束日期"
            style="width: 240px"
          />
        </el-form-item>

        <el-form-item label="关键词">
          <el-input
            v-model="form.keyword"
            placeholder="搜索申请人或申请编号"
            style="width: 200px"
          >
            <template #prefix>
              <el-icon><Search /></el-icon>
            </template>
          </el-input>
        </el-form-item>
      </template>

      <template #stats="{ total, selected }">
        <span>共 {{ total }} 条记录</span>
        <span v-if="selected > 0">，已选择 {{ selected }} 条</span>
      </template>
    </TableActions>

    <!-- 统计卡片组 -->
    <div class="stats-row">
      <el-row :gutter="16">
        <el-col :span="6">
          <StatCard
            :value="stats.pending"
            label="待审批"
            :icon="Clock"
            type="warning"
            :clickable="true"
            @click="handleStatFilter('pending')"
          />
        </el-col>
        <el-col :span="6">
          <StatCard
            :value="stats.approved"
            label="已通过"
            :icon="Check"
            type="success"
            :clickable="true"
            @click="handleStatFilter('approved')"
          />
        </el-col>
        <el-col :span="6">
          <StatCard
            :value="stats.rejected"
            label="已拒绝"
            :icon="Close"
            type="danger"
            :clickable="true"
            @click="handleStatFilter('rejected')"
          />
        </el-col>
        <el-col :span="6">
          <StatCard
            :value="stats.totalAmount"
            label="总金额"
            :icon="Money"
            type="primary"
            format="currency"
          />
        </el-col>
      </el-row>
    </div>

    <!-- 审批列表 -->
    <el-card class="list-card" shadow="never">
      <template #header>
        <div class="card-header">
          <span>贷款审批列表</span>
        </div>
      </template>

      <el-table
        v-loading="loading"
        :data="approvalList"
        @selection-change="handleSelectionChange"
        stripe
      >
        <el-table-column type="selection" width="55" />
        
        <el-table-column prop="application_number" label="申请编号" width="140">
          <template #default="{ row }">
            <el-link type="primary" @click="handleViewDetail(row)">
              {{ row.application_number }}
            </el-link>
          </template>
        </el-table-column>
        
        <el-table-column prop="applicant_name" label="申请人" width="100" />
        
        <el-table-column prop="loan_type" label="贷款类型" width="120">
          <template #default="{ row }">
            <StatusTag :status="row.loan_type" :status-map="loanTypeStatusMap" />
          </template>
        </el-table-column>
        
        <el-table-column prop="amount" label="申请金额" width="120">
          <template #default="{ row }">
            <span class="amount">{{ formatAmount(row.amount) }}</span>
          </template>
        </el-table-column>
        
        <el-table-column prop="term" label="期限" width="80">
          <template #default="{ row }">
            {{ row.term }}个月
          </template>
        </el-table-column>
        
        <el-table-column prop="interest_rate" label="利率" width="80">
          <template #default="{ row }">
            {{ (row.interest_rate * 100).toFixed(2) }}%
          </template>
        </el-table-column>
        
        <el-table-column prop="status" label="状态" width="100">
          <template #default="{ row }">
            <StatusTag :status="row.status" :show-icon="true" />
          </template>
        </el-table-column>
        
        <el-table-column prop="created_at" label="申请时间" width="160">
          <template #default="{ row }">
            {{ formatDateTime(row.created_at) }}
          </template>
        </el-table-column>
        
        <el-table-column prop="reviewed_at" label="审批时间" width="160">
          <template #default="{ row }">
            {{ row.reviewed_at ? formatDateTime(row.reviewed_at) : '-' }}
          </template>
        </el-table-column>
        
        <el-table-column label="操作" width="180" fixed="right">
          <template #default="{ row }">
            <el-button size="small" @click="handleViewDetail(row)">详情</el-button>
            <el-button 
              v-if="row.status === 'pending'" 
              type="success" 
              size="small" 
              @click="handleApprove(row)"
            >
              通过
            </el-button>
            <el-button 
              v-if="row.status === 'pending'" 
              type="danger" 
              size="small" 
              @click="handleReject(row)"
            >
              拒绝
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
          @size-change="handleSizeChange"
          @current-change="handlePageChange"
        />
      </div>
    </el-card>

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
        
        <el-form-item v-if="approvalDialog.type === 'approve'" label="批准利率" prop="approvedRate">
          <el-input-number
            v-model="form.approvedRate"
            :min="0.01"
            :max="0.5"
            :step="0.001"
            :precision="3"
            style="width: 100%"
          />
          <span style="margin-left: 8px; color: #999;">%</span>
        </el-form-item>
      </template>
    </FormDialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import dayjs from 'dayjs'
import {
  CreditCard,
  Download,
  Search,
  Clock,
  Check,
  Close,
  Money
} from '@element-plus/icons-vue'

// 导入通用组件
import { PageHeader, StatCard, StatusTag, TableActions, FormDialog } from '@/components/common'

const router = useRouter()

// 响应式数据
const loading = ref(false)
const selectedApprovals = ref<any[]>([])

const stats = reactive({
  pending: 156,
  approved: 89,
  rejected: 23,
  totalAmount: 12580000
})

const pagination = reactive({
  page: 1,
  size: 20,
  total: 268
})

const approvalList = ref<any[]>([])

const approvalDialog = reactive({
  visible: false,
  type: 'approve' as 'approve' | 'reject',
  data: {}
})

// 计算属性
const breadcrumbs = computed(() => [
  { title: '首页', to: '/' },
  { title: '核心审批', to: '/approval' },
  { title: '贷款审批' }
])

const pageStats = computed(() => [
  { label: '待处理', value: stats.pending, color: '#e6a23c', clickable: true, key: 'pending' },
  { label: '今日审批', value: 25, color: '#67c23a', clickable: true, key: 'today' },
  { label: '审批率', value: '79.2%', color: '#409eff' }
])

const batchActions = computed(() => [
  { key: 'approve', label: '批量通过', icon: Check },
  { key: 'reject', label: '批量拒绝', icon: Close, danger: true },
  { key: 'export', label: '导出选中', icon: Download }
])

// 状态映射
const loanTypeStatusMap = {
  'personal_credit': { type: 'primary', text: '个人信用贷' },
  'business_loan': { type: 'success', text: '企业经营贷' },
  'mortgage': { type: 'warning', text: '抵押贷款' },
  'vehicle': { type: 'info', text: '车辆贷款' }
}

// 表单验证规则
const approvalRules = {
  comment: [
    { required: true, message: '请输入处理意见', trigger: 'blur' }
  ],
  approvedAmount: [
    { required: true, message: '请输入批准金额', trigger: 'blur' }
  ],
  approvedRate: [
    { required: true, message: '请输入批准利率', trigger: 'blur' }
  ]
}

// 方法
const handleSearch = (filters: Record<string, any>) => {
  console.log('搜索过滤条件:', filters)
  fetchApprovalList()
}

const handleReset = () => {
  console.log('重置过滤条件')
  fetchApprovalList()
}

const handleRefresh = () => {
  fetchApprovalList()
}

const handleBatchAction = (action: string, items: any[]) => {
  console.log('批量操作:', action, items)
  switch (action) {
    case 'approve':
      // 批量通过逻辑
      break
    case 'reject':
      // 批量拒绝逻辑
      break
    case 'export':
      // 导出逻辑
      break
  }
}

const handleStatClick = (stat: any) => {
  if (stat.key) {
    handleStatFilter(stat.key)
  }
}

const handleStatFilter = (status: string) => {
  console.log('按状态筛选:', status)
  // 实现状态筛选逻辑
}

const handleSelectionChange = (selection: any[]) => {
  selectedApprovals.value = selection
}

const handleViewDetail = (row: any) => {
  router.push(`/approval/loan/${row.id}`)
}

const handleApprove = (row: any) => {
  approvalDialog.type = 'approve'
  approvalDialog.data = {
    id: row.id,
    requestedAmount: row.amount,
    approvedAmount: row.amount,
    approvedRate: row.interest_rate
  }
  approvalDialog.visible = true
}

const handleReject = (row: any) => {
  approvalDialog.type = 'reject'
  approvalDialog.data = { id: row.id }
  approvalDialog.visible = true
}

const handleApprovalSubmit = async (data: any, mode: string) => {
  try {
    console.log('提交审批:', data)
    // 实现审批提交逻辑
    approvalDialog.visible = false
    await fetchApprovalList()
  } catch (error) {
    console.error('审批失败:', error)
  }
}

const handleExport = () => {
  console.log('导出数据')
}

const handleSizeChange = (size: number) => {
  pagination.size = size
  fetchApprovalList()
}

const handlePageChange = (page: number) => {
  pagination.page = page
  fetchApprovalList()
}

const fetchApprovalList = async () => {
  loading.value = true
  try {
    // 模拟API调用
    await new Promise(resolve => setTimeout(resolve, 500))
    
    approvalList.value = Array.from({ length: pagination.size }, (_, index) => ({
      id: index + 1,
      application_number: `LA${Date.now() + index}`,
      applicant_name: `申请人${index + 1}`,
      loan_type: ['personal_credit', 'business_loan', 'mortgage', 'vehicle'][index % 4],
      amount: Math.floor(Math.random() * 1000000) + 50000,
      term: [12, 24, 36, 48][index % 4],
      interest_rate: 0.05 + Math.random() * 0.1,
      status: ['pending', 'approved', 'rejected'][index % 3],
      created_at: new Date(Date.now() - Math.random() * 30 * 24 * 60 * 60 * 1000).toISOString(),
      reviewed_at: index % 3 !== 0 ? new Date().toISOString() : null
    }))
  } finally {
    loading.value = false
  }
}

const formatAmount = (amount: number) => {
  return `¥${amount.toLocaleString()}`
}

const formatDateTime = (datetime: string) => {
  return dayjs(datetime).format('YYYY-MM-DD HH:mm:ss')
}

onMounted(() => {
  fetchApprovalList()
})
</script>

<style scoped>
.loan-approval {
  height: 100vh;
  display: flex;
  flex-direction: column;
  gap: 16px;
  padding: 0 20px;
  overflow: hidden;
}

.stats-row {
  margin-bottom: 16px;
}

.list-card {
  flex: 1;
  display: flex;
  flex-direction: column;
  overflow: hidden;
  border: none;
  border-radius: 12px;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  font-weight: 500;
}

.amount {
  font-weight: 600;
  color: #409eff;
}

.pagination-wrapper {
  display: flex;
  justify-content: center;
  margin-top: 16px;
  padding: 16px 0;
  border-top: 1px solid #f0f0f0;
}

:deep(.el-card__body) {
  flex: 1;
  display: flex;
  flex-direction: column;
  overflow: hidden;
}

:deep(.el-table) {
  flex: 1;
}
</style> 