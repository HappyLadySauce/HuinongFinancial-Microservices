<template>
  <div class="loan-approval">
    <!-- 顶部筛选栏 -->
    <el-card class="filter-card" shadow="never">
      <div class="filter-row">
        <div class="filter-left">
          <el-select v-model="filterStatus" placeholder="选择状态" clearable style="width: 120px">
            <el-option label="全部" value="" />
            <el-option label="待审批" value="pending" />
            <el-option label="已通过" value="approved" />
            <el-option label="已拒绝" value="rejected" />
          </el-select>
          
          <el-select v-model="filterType" placeholder="贷款类型" clearable style="width: 140px">
            <el-option label="全部类型" value="" />
            <el-option label="个人信用贷" value="personal_credit" />
            <el-option label="企业经营贷" value="business_loan" />
            <el-option label="抵押贷款" value="mortgage" />
            <el-option label="车辆贷款" value="vehicle" />
          </el-select>

          <el-date-picker
            v-model="dateRange"
            type="daterange"
            range-separator="至"
            start-placeholder="开始日期"
            end-placeholder="结束日期"
            style="width: 240px"
          />
        </div>
        
        <div class="filter-right">
          <el-input
            v-model="searchKeyword"
            placeholder="搜索申请人或申请编号"
            style="width: 200px"
            @keyup.enter="handleSearch"
          >
            <template #prefix>
              <el-icon><Search /></el-icon>
            </template>
          </el-input>
          <el-button type="primary" @click="handleSearch">搜索</el-button>
          <el-button @click="handleReset">重置</el-button>
        </div>
      </div>
    </el-card>

    <!-- 数据统计卡片 -->
    <div class="stats-row">
      <el-card class="stat-card" shadow="hover">
        <div class="stat-content">
          <div class="stat-number">{{ stats.pending }}</div>
          <div class="stat-label">待审批</div>
        </div>
      </el-card>
      <el-card class="stat-card" shadow="hover">
        <div class="stat-content">
          <div class="stat-number">{{ stats.approved }}</div>
          <div class="stat-label">已通过</div>
        </div>
      </el-card>
      <el-card class="stat-card" shadow="hover">
        <div class="stat-content">
          <div class="stat-number">{{ stats.rejected }}</div>
          <div class="stat-label">已拒绝</div>
        </div>
      </el-card>
      <el-card class="stat-card" shadow="hover">
        <div class="stat-content">
          <div class="stat-number">{{ formatAmount(stats.totalAmount) }}</div>
          <div class="stat-label">总金额</div>
        </div>
      </el-card>
    </div>

    <!-- 审批列表 -->
    <el-card class="list-card" shadow="never">
      <template #header>
        <div class="card-header">
          <span>贷款审批列表</span>
          <div class="header-actions">
            <el-button size="small" @click="handleRefresh">
              <el-icon><Refresh /></el-icon>
              刷新
            </el-button>
            <el-button type="primary" size="small" @click="handleBatchApprove" :disabled="selectedIds.length === 0">
              批量审批
            </el-button>
          </div>
        </div>
      </template>

      <el-table
        v-loading="loading"
        :data="approvalList"
        @selection-change="handleSelectionChange"
        stripe
      >
        <el-table-column type="selection" width="55" />
        
        <el-table-column prop="id" label="申请编号" width="140">
          <template #default="{ row }">
            <el-link type="primary" @click="handleViewDetail(row)">
              LN{{ String(row.id).padStart(6, '0') }}
            </el-link>
          </template>
        </el-table-column>
        
        <el-table-column prop="name" label="贷款名称" width="120" />
        
        <el-table-column prop="type" label="贷款类型" width="120">
          <template #default="{ row }">
            <el-tag :type="getLoanTypeColor(row.type)" size="small">
              {{ getLoanTypeName(row.type) }}
            </el-tag>
          </template>
        </el-table-column>
        
        <el-table-column prop="amount" label="申请金额" width="120">
          <template #default="{ row }">
            <span class="amount">{{ formatAmount(row.amount) }}</span>
          </template>
        </el-table-column>
        
        <el-table-column prop="duration" label="期限" width="80">
          <template #default="{ row }">
            {{ row.duration }}个月
          </template>
        </el-table-column>
        
        <el-table-column prop="description" label="描述" width="200">
          <template #default="{ row }">
            {{ row.description || '-' }}
          </template>
        </el-table-column>
        
        <el-table-column prop="status" label="状态" width="100">
          <template #default="{ row }">
            <el-tag :type="getStatusColor(row.status)" size="small">
              {{ getStatusName(row.status) }}
            </el-tag>
          </template>
        </el-table-column>
        
        <el-table-column prop="created_at" label="申请时间" width="160">
          <template #default="{ row }">
            {{ formatDateTime(row.created_at) }}
          </template>
        </el-table-column>
        
        <el-table-column prop="updated_at" label="更新时间" width="160">
          <template #default="{ row }">
            {{ formatDateTime(row.updated_at) }}
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
    <el-dialog
      v-model="approvalDialog.visible"
      :title="approvalDialog.type === 'approve' ? '审批通过' : '审批拒绝'"
      width="500px"
    >
      <el-form :model="approvalDialog.form" label-width="80px">
        <el-form-item label="申请编号">
          <span>LN{{ String(approvalDialog.currentRow?.id).padStart(6, '0') }}</span>
        </el-form-item>
        <el-form-item label="贷款名称">
          <span>{{ approvalDialog.currentRow?.name }}</span>
        </el-form-item>
        <el-form-item label="申请金额">
          <span>{{ formatAmount(approvalDialog.currentRow?.amount || 0) }}</span>
        </el-form-item>
        <el-form-item label="审批意见" required>
          <el-input
            v-model="approvalDialog.form.comment"
            type="textarea"
            :rows="4"
            placeholder="请输入审批意见"
          />
        </el-form-item>
        <el-form-item v-if="approvalDialog.type === 'approve'" label="批准金额">
          <el-input-number
            v-model="approvalDialog.form.approved_amount"
            :min="0"
            :max="approvalDialog.currentRow?.amount || 0"
            :precision="2"
            style="width: 100%"
          />
        </el-form-item>
      </el-form>
      
      <template #footer>
        <el-button @click="approvalDialog.visible = false">取消</el-button>
        <el-button
          type="primary"
          :loading="approvalDialog.loading"
          @click="handleConfirmApproval"
        >
          确认
        </el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { Search, Refresh } from '@element-plus/icons-vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { adminApi, type LoanApproval } from '@/services/api'

const router = useRouter()

// 筛选条件
const filterStatus = ref('')
const filterType = ref('')
const dateRange = ref<[Date, Date] | null>(null)
const searchKeyword = ref('')

// 列表数据
const loading = ref(false)
const approvalList = ref<LoanApproval[]>([])
const selectedIds = ref<number[]>([])

// 统计数据
const stats = reactive({
  pending: 0,
  approved: 0,
  rejected: 0,
  totalAmount: 0
})

// 分页
const pagination = reactive({
  page: 1,
  size: 20,
  total: 0
})

// 审批对话框
const approvalDialog = reactive({
  visible: false,
  type: '' as 'approve' | 'reject',
  loading: false,
  currentRow: null as any,
  form: {
    comment: '',
    approved_amount: 0
  }
})

// 格式化金额
const formatAmount = (amount: number) => {
  if (amount >= 10000) {
    return `¥${(amount / 10000).toFixed(1)}万`
  }
  return `¥${amount.toLocaleString()}`
}

// 格式化日期时间
const formatDateTime = (dateString: string) => {
  if (!dateString) return '-'
  const date = new Date(dateString)
  return date.toLocaleString('zh-CN', {
    year: 'numeric',
    month: '2-digit',
    day: '2-digit',
    hour: '2-digit',
    minute: '2-digit'
  })
}

// 获取贷款类型名称
const getLoanTypeName = (type: string) => {
  const types: Record<string, string> = {
    personal_credit: '个人信用贷',
    business_loan: '企业经营贷',
    mortgage: '抵押贷款',
    vehicle: '车辆贷款',
    '个人信用贷': '个人信用贷',
    '企业经营贷': '企业经营贷',
    '抵押贷款': '抵押贷款',
    '车辆贷款': '车辆贷款'
  }
  return types[type] || type
}

// 获取贷款类型颜色
const getLoanTypeColor = (type: string) => {
  const colors: Record<string, string> = {
    personal_credit: 'primary',
    business_loan: 'success',
    mortgage: 'warning',
    vehicle: 'info',
    '个人信用贷': 'primary',
    '企业经营贷': 'success',
    '抵押贷款': 'warning',
    '车辆贷款': 'info'
  }
  return colors[type] || 'default'
}

// 获取状态名称
const getStatusName = (status: string) => {
  const statuses: Record<string, string> = {
    pending: '待审批',
    approved: '已通过',
    rejected: '已拒绝'
  }
  return statuses[status] || status
}

// 获取状态颜色
const getStatusColor = (status: string) => {
  const colors: Record<string, string> = {
    pending: 'warning',
    approved: 'success',
    rejected: 'danger'
  }
  return colors[status] || 'default'
}

// 加载数据
const loadData = async () => {
  loading.value = true
  
  try {
    const params: any = {
      page: pagination.page,
      page_size: pagination.size
    }
    
    // 添加筛选条件
    if (filterStatus.value) {
      params.status = filterStatus.value
    }
    
    const response = await adminApi.getAllLoanApprovals(params)
    
    approvalList.value = response.list
    pagination.total = response.total
    
    // 更新统计数据
    await loadStats()
    
  } catch (error) {
    console.error('加载数据失败:', error)
    ElMessage.error('加载数据失败')
  } finally {
    loading.value = false
  }
}

// 加载统计数据
const loadStats = async () => {
  try {
    // 并行获取各状态的统计数据
    const [pendingRes, approvedRes, rejectedRes, allRes] = await Promise.all([
      adminApi.getAllLoanApprovals({ page: 1, page_size: 1, status: 'pending' }),
      adminApi.getAllLoanApprovals({ page: 1, page_size: 1, status: 'approved' }),
      adminApi.getAllLoanApprovals({ page: 1, page_size: 1, status: 'rejected' }),
      adminApi.getAllLoanApprovals({ page: 1, page_size: 1 })
    ])
    
    stats.pending = pendingRes.total
    stats.approved = approvedRes.total
    stats.rejected = rejectedRes.total
    
    // 计算总金额（仅已通过的申请）
    if (approvedRes.total > 0) {
      const approvedList = await adminApi.getAllLoanApprovals({ 
        page: 1, 
        page_size: approvedRes.total, 
        status: 'approved' 
      })
      stats.totalAmount = approvedList.list.reduce((total, item) => total + item.amount, 0)
    }
    
  } catch (error) {
    console.error('加载统计数据失败:', error)
  }
}

// 搜索
const handleSearch = () => {
  loadData()
}

// 重置筛选
const handleReset = () => {
  filterStatus.value = ''
  filterType.value = ''
  dateRange.value = null
  searchKeyword.value = ''
  loadData()
}

// 刷新
const handleRefresh = () => {
  loadData()
}

// 选择变更
const handleSelectionChange = (selection: any[]) => {
  selectedIds.value = selection.map(item => item.id)
}

// 分页变更
const handlePageChange = (page: number) => {
  pagination.page = page
  loadData()
}

const handleSizeChange = (size: number) => {
  pagination.size = size
  pagination.page = 1
  loadData()
}

// 查看详情
const handleViewDetail = (row: LoanApproval) => {
  router.push(`/approval/loan/${row.id}`)
}

// 审批通过
const handleApprove = (row: any) => {
  approvalDialog.visible = true
  approvalDialog.type = 'approve'
  approvalDialog.currentRow = row
  approvalDialog.form.approved_amount = row.amount
  approvalDialog.form.comment = ''
}

// 审批拒绝
const handleReject = (row: any) => {
  approvalDialog.visible = true
  approvalDialog.type = 'reject'
  approvalDialog.currentRow = row
  approvalDialog.form.comment = ''
}

// 确认审批
const handleConfirmApproval = async () => {
  if (!approvalDialog.form.comment.trim()) {
    ElMessage.warning('请输入审批意见')
    return
  }

  approvalDialog.loading = true
  
  try {
    const reviewData = {
      status: approvalDialog.type === 'approve' ? 'approved' as const : 'rejected' as const,
      suggestions: approvalDialog.form.comment,
      auditor: '当前管理员' // 这里应该从登录用户信息获取
    }
    
    await adminApi.reviewLoanApproval(approvalDialog.currentRow.id, reviewData)
    
    const action = approvalDialog.type === 'approve' ? '通过' : '拒绝'
    ElMessage.success(`审批${action}成功`)
    
    approvalDialog.visible = false
    loadData()
  } catch (error) {
    console.error('审批操作失败:', error)
    ElMessage.error('审批操作失败')
  } finally {
    approvalDialog.loading = false
  }
}

// 批量审批
const handleBatchApprove = () => {
  ElMessageBox.confirm(
    `确定要批量审批选中的 ${selectedIds.value.length} 条记录吗？`,
    '批量审批',
    {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning'
    }
  ).then(() => {
    ElMessage.success('批量审批操作已提交')
    loadData()
  }).catch(() => {
    // 用户取消
  })
}

onMounted(() => {
  loadData()
})
</script>

<style scoped>
.loan-approval {
  padding: 20px;
  background: #f5f5f5;
  min-height: calc(100vh - 60px);
}

.filter-card {
  margin-bottom: 16px;
  border: none;
}

.filter-row {
  display: flex;
  justify-content: space-between;
  align-items: center;
  gap: 16px;
}

.filter-left {
  display: flex;
  gap: 16px;
  align-items: center;
}

.filter-right {
  display: flex;
  gap: 8px;
  align-items: center;
}

.stats-row {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(200px, 1fr));
  gap: 16px;
  margin-bottom: 16px;
}

.stat-card {
  border: none;
  text-align: center;
}

.stat-content {
  padding: 8px 0;
}

.stat-number {
  font-size: 28px;
  font-weight: bold;
  color: #409EFF;
  margin-bottom: 4px;
}

.stat-label {
  font-size: 14px;
  color: #909399;
}

.list-card {
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
  font-weight: 500;
  color: #E6A23C;
}

.pagination-wrapper {
  margin-top: 16px;
  text-align: right;
}

@media (max-width: 768px) {
  .filter-row {
    flex-direction: column;
    align-items: stretch;
  }
  
  .filter-left,
  .filter-right {
    flex-wrap: wrap;
  }
  
  .stats-row {
    grid-template-columns: repeat(2, 1fr);
  }
}
</style> 