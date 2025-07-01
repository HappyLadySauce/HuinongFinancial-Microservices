<script setup lang="ts">
import { ref, onMounted, computed } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage, ElMessageBox } from 'element-plus'
import { useUserStore } from '../stores/user'
import { loanApprovalApi } from '../services/api'
import type { LoanApproval } from '../services/api'

const router = useRouter()
const userStore = useUserStore()
const loading = ref(false)
const refreshing = ref(false)

// 申请列表
const applications = ref<LoanApproval[]>([])
const total = ref(0)
const page = ref(1)
const limit = ref(10)

// 状态筛选
const selectedStatus = ref('')
const statusOptions = [
  { label: '全部', value: '' },
  { label: '待审批', value: 'pending' },
  { label: '已批准', value: 'approved' },
  { label: '已拒绝', value: 'rejected' }
]

// 状态映射
const statusMap = {
  'pending': { text: '待审批', color: '#E6A23C', icon: '⏳' },
  'approved': { text: '已批准', color: '#67C23A', icon: '✅' },
  'rejected': { text: '已拒绝', color: '#F56C6C', icon: '❌' }
}

// 是否有更多数据
const hasMore = computed(() => {
  return applications.value.length < total.value
})

// 格式化时间
const formatDate = (dateString: string) => {
  const date = new Date(dateString)
  return date.toLocaleDateString('zh-CN', {
    year: 'numeric',
    month: '2-digit',
    day: '2-digit'
  })
}

// 格式化详细时间
const formatDateTime = (dateString: string) => {
  const date = new Date(dateString)
  return date.toLocaleString('zh-CN', {
    year: 'numeric',
    month: '2-digit',
    day: '2-digit',
    hour: '2-digit',
    minute: '2-digit'
  })
}

// 获取状态信息
const getStatusInfo = (status: string) => {
  return statusMap[status as keyof typeof statusMap] || {
    text: status,
    color: '#909399',
    icon: '❓'
  }
}

// 加载申请列表
const loadApplications = async (reset = false) => {
  try {
    if (reset) {
      loading.value = true
      page.value = 1
      applications.value = []
    } else {
      refreshing.value = true
    }

    const params = {
      page: page.value,
      page_size: limit.value
    }

    const response = await loanApprovalApi.getMyApprovals(params)
    
    let newApplications = response.list
    
    // 如果有状态筛选，在前端过滤
    if (selectedStatus.value) {
      newApplications = newApplications.filter(app => app.status === selectedStatus.value)
    }
    
    if (reset) {
      applications.value = newApplications
    } else {
      applications.value.push(...newApplications)
    }
    
    total.value = response.total

  } catch (error: any) {
    console.error('加载申请列表失败:', error)
    ElMessage.error('加载申请列表失败')
  } finally {
    loading.value = false
    refreshing.value = false
  }
}

// 加载更多
const loadMore = () => {
  if (hasMore.value && !refreshing.value) {
    page.value++
    loadApplications()
  }
}

// 筛选状态变化
const handleStatusChange = () => {
  loadApplications(true)
}

// 查看申请详情
const viewDetail = (applicationId: number) => {
  router.push(`/loan/application/${applicationId}`)
}

// 编辑申请（仅对pending状态）
const editApplication = async (application: LoanApproval) => {
  if (application.status !== 'pending') {
    ElMessage.warning('只能编辑待审批的申请')
    return
  }
  
  // 这里可以跳转到编辑页面或打开编辑对话框
  ElMessage.info('编辑功能暂未实现')
}

// 删除申请
const deleteApplication = async (application: LoanApproval) => {
  if (application.status !== 'pending') {
    ElMessage.warning('只能删除待审批的申请')
    return
  }
  
  try {
    await ElMessageBox.confirm(
      `确认删除申请"${application.name}"吗？删除后无法恢复。`,
      '确认删除',
      {
        confirmButtonText: '确认删除',
        cancelButtonText: '取消',
        type: 'warning'
      }
    )
    
    await loanApprovalApi.delete(application.id)
    ElMessage.success('删除成功')
    
    // 重新加载列表
    loadApplications(true)
    
  } catch (error: any) {
    if (error.message !== 'cancel') {
      console.error('删除申请失败:', error)
      ElMessage.error(error.message || '删除申请失败')
    }
  }
}

// 返回逻辑
const goBack = () => {
  // 检查路由query参数中是否有来源页面信息
  const from = router.currentRoute.value.query.from as string
  
  if (from) {
    // 如果有明确的来源页面，跳转回去
    router.push(from)
  } else {
    // 尝试使用history.length判断是否可以安全返回
    if (window.history.length > 1) {
      // 有历史记录，尝试返回
      router.go(-1)
    } else {
      // 没有历史记录或者是直接访问，跳转到金融页面
      router.push('/finance')
    }
  }
}

// 申请新贷款
const applyNewLoan = () => {
  router.push('/finance')
}

// 组件挂载时加载数据
onMounted(() => {
  // 检查登录状态
  if (!userStore.isLoggedIn) {
    ElMessage.error('请先登录')
    router.push('/login')
    return
  }
  
  loadApplications(true)
})
</script>

<template>
  <div class="my-applications-page">
    <!-- 顶部导航 -->
    <div class="top-nav">
      <div class="nav-left" @click="goBack">
        <el-icon><ArrowLeft /></el-icon>
      </div>
      <div class="nav-title">我的申请</div>
      <div class="nav-right">
        <el-icon @click="loadApplications(true)" :class="{ 'is-loading': refreshing }">
          <Refresh />
        </el-icon>
      </div>
    </div>

    <div class="page-content">
      <!-- 顶部统计 -->
      <div class="stats-card">
        <div class="stat-item">
          <div class="stat-value">{{ total }}</div>
          <div class="stat-label">总申请数</div>
        </div>
        <div class="stat-item">
          <div class="stat-value">
            {{ applications.filter(app => app.status === 'pending').length }}
          </div>
          <div class="stat-label">待审批</div>
        </div>
        <div class="stat-item">
          <div class="stat-value">
            {{ applications.filter(app => app.status === 'approved').length }}
          </div>
          <div class="stat-label">已批准</div>
        </div>
        <div class="stat-item">
          <div class="stat-value">
            {{ (applications.filter(app => app.status === 'approved').reduce((sum, app) => sum + app.amount, 0) / 10000).toFixed(1) }}万
          </div>
          <div class="stat-label">批准总额</div>
        </div>
      </div>

      <!-- 筛选栏 -->
      <div class="filter-bar">
        <el-select 
          v-model="selectedStatus" 
          placeholder="筛选状态"
          @change="handleStatusChange"
          style="width: 150px"
        >
          <el-option
            v-for="option in statusOptions"
            :key="option.value"
            :label="option.label"
            :value="option.value"
          />
        </el-select>
      </div>

      <!-- 申请列表 -->
      <div v-if="!loading" class="applications-list">
        <div 
          v-for="app in applications" 
          :key="app.id"
          class="application-card"
          @click="viewDetail(app.id)"
        >
          <div class="card-header">
            <div class="app-title">
              <h3 class="app-name">{{ app.name }}</h3>
              <div class="app-type">{{ app.type }}</div>
            </div>
            <div 
              class="status-badge"
              :style="{ 
                backgroundColor: getStatusInfo(app.status).color,
                color: 'white'
              }"
            >
              <span class="status-icon">{{ getStatusInfo(app.status).icon }}</span>
              <span class="status-text">{{ getStatusInfo(app.status).text }}</span>
            </div>
          </div>
          
          <div class="card-body">
            <div class="amount-section">
            <div class="amount-info">
              <span class="amount-label">申请金额</span>
              <span class="amount-value">¥{{ app.amount.toLocaleString() }}</span>
              </div>
              <div class="duration-info">
                <span class="duration-label">期限</span>
                <span class="duration-value">{{ app.duration }}个月</span>
              </div>
            </div>
            
            <div class="description-section">
              <div class="description-label">申请用途:</div>
              <div class="description-text">{{ app.description }}</div>
            </div>

            <div v-if="app.status === 'approved' && app.suggestions" class="suggestions-section">
              <div class="suggestions-label">审批意见:</div>
              <div class="suggestions-text">{{ app.suggestions }}</div>
              <div v-if="app.auditor" class="auditor-info">审批人: {{ app.auditor }}</div>
            </div>

            <div v-if="app.status === 'rejected' && app.suggestions" class="rejection-section">
              <div class="rejection-label">拒绝原因:</div>
              <div class="rejection-text">{{ app.suggestions }}</div>
              <div v-if="app.auditor" class="auditor-info">审批人: {{ app.auditor }}</div>
            </div>
          </div>

          <div class="card-footer">
            <div class="time-info">
              <div class="created-time">
                <span class="time-label">申请时间:</span>
                <span class="time-value">{{ formatDateTime(app.created_at) }}</span>
              </div>
              <div v-if="app.updated_at !== app.created_at" class="updated-time">
                <span class="time-label">更新时间:</span>
                <span class="time-value">{{ formatDateTime(app.updated_at) }}</span>
              </div>
            </div>
            
            <!-- 操作按钮 -->
            <div class="actions" @click.stop>
              <el-button 
                v-if="app.status === 'pending'"
                type="warning" 
                size="small"
                @click="editApplication(app)"
              >
                编辑
              </el-button>
              <el-button 
                v-if="app.status === 'pending'"
                type="danger" 
                size="small"
                @click="deleteApplication(app)"
              >
                删除
              </el-button>
              <el-button 
                type="primary" 
                size="small"
                @click="viewDetail(app.id)"
              >
                查看详情
              </el-button>
            </div>
          </div>
        </div>

        <!-- 空状态 -->
        <div v-if="applications.length === 0" class="empty-state">
          <div class="empty-icon">📋</div>
          <div class="empty-title">暂无申请记录</div>
          <div class="empty-desc">
            {{ selectedStatus ? '当前筛选条件下暂无申请记录' : '您还没有提交过贷款申请' }}
          </div>
          <el-button type="primary" @click="applyNewLoan" class="empty-action">
            {{ selectedStatus ? '查看全部申请' : '立即申请' }}
            </el-button>
        </div>

        <!-- 加载更多 -->
        <div v-if="hasMore" class="load-more">
          <el-button 
            @click="loadMore" 
            :loading="refreshing"
            type="text"
            class="load-more-btn"
          >
            {{ refreshing ? '加载中...' : '加载更多' }}
          </el-button>
        </div>

        <!-- 已加载全部 -->
        <div v-else-if="applications.length > 0" class="load-complete">
          <p>已加载全部 {{ applications.length }} 条记录</p>
        </div>
      </div>

      <!-- 加载状态 -->
      <div v-if="loading" class="loading-container">
        <el-icon class="is-loading"><Loading /></el-icon>
        <p>加载中...</p>
      </div>
    </div>

    <!-- 底部操作栏 -->
    <div class="bottom-actions">
      <el-button 
        type="primary" 
        size="large" 
        @click="applyNewLoan"
        class="apply-btn"
      >
        申请新贷款
      </el-button>
    </div>
  </div>
</template>

<style scoped>
.my-applications-page {
  min-height: 100vh;
  background-color: #f5f5f5;
  display: flex;
  flex-direction: column;
}

.top-nav {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 12px 16px;
  background: white;
  box-shadow: 0 1px 2px rgba(0, 0, 0, 0.1);
  position: sticky;
  top: 0;
  z-index: 100;
}

.nav-left, .nav-right {
  cursor: pointer;
  padding: 8px;
  width: 32px;
  display: flex;
  align-items: center;
  justify-content: center;
}

.nav-title {
  font-size: 18px;
  font-weight: 600;
  color: #2c3e50;
}

.page-content {
  flex: 1;
  padding: 16px;
  max-width: 600px;
  margin: 0 auto;
  width: 100%;
  box-sizing: border-box;
}

.stats-card {
  background: white;
  border-radius: 12px;
  padding: 20px;
  margin-bottom: 16px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
  display: grid;
  grid-template-columns: repeat(4, 1fr);
  gap: 16px;
}

.stat-item {
  text-align: center;
}

.stat-value {
  font-size: 24px;
  font-weight: 700;
  color: #27ae60;
  margin-bottom: 4px;
}

.stat-label {
  font-size: 12px;
  color: #7f8c8d;
}

.filter-bar {
  margin-bottom: 16px;
  display: flex;
  justify-content: flex-end;
}

.applications-list {
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.application-card {
  background: white;
  border-radius: 12px;
  padding: 20px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
  cursor: pointer;
  transition: all 0.3s ease;
  border-left: 4px solid #27ae60;
}

.application-card:hover {
  transform: translateY(-2px);
  box-shadow: 0 8px 20px rgba(0, 0, 0, 0.15);
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  margin-bottom: 16px;
}

.app-title {
  flex: 1;
}

.app-name {
  font-size: 18px;
  font-weight: 600;
  color: #2c3e50;
  margin: 0 0 4px 0;
}

.app-type {
  font-size: 13px;
  color: #7f8c8d;
  background: #f8f9fa;
  padding: 2px 8px;
  border-radius: 12px;
  display: inline-block;
}

.status-badge {
  display: flex;
  align-items: center;
  gap: 4px;
  padding: 6px 12px;
  border-radius: 16px;
  font-size: 13px;
  font-weight: 500;
  white-space: nowrap;
}

.status-icon {
  font-size: 14px;
}

.card-body {
  margin-bottom: 16px;
}

.amount-section {
  display: grid;
  grid-template-columns: 2fr 1fr;
  gap: 16px;
  margin-bottom: 16px;
  padding: 12px;
  background: #f8f9fa;
  border-radius: 8px;
}

.amount-info, .duration-info {
  display: flex;
  flex-direction: column;
  align-items: center;
}

.amount-label, .duration-label {
  font-size: 12px;
  color: #7f8c8d;
  margin-bottom: 4px;
}

.amount-value {
  font-size: 20px;
  font-weight: 600;
  color: #27ae60;
}

.duration-value {
  font-size: 16px;
  font-weight: 500;
  color: #2c3e50;
}

.description-section {
  margin-bottom: 12px;
}

.description-label {
  font-size: 13px;
  color: #7f8c8d;
  margin-bottom: 4px;
}

.description-text {
  font-size: 14px;
  color: #2c3e50;
  line-height: 1.5;
}

.suggestions-section {
  background: #e8f5e8;
  padding: 12px;
  border-radius: 8px;
  margin-bottom: 12px;
}

.suggestions-label {
  font-size: 13px;
  color: #27ae60;
  font-weight: 500;
  margin-bottom: 4px;
}

.suggestions-text {
  font-size: 14px;
  color: #2c3e50;
  line-height: 1.5;
  margin-bottom: 4px;
}

.rejection-section {
  background: #ffebee;
  padding: 12px;
  border-radius: 8px;
  margin-bottom: 12px;
}

.rejection-label {
  font-size: 13px;
  color: #e74c3c;
  font-weight: 500;
  margin-bottom: 4px;
}

.rejection-text {
  font-size: 14px;
  color: #2c3e50;
  line-height: 1.5;
  margin-bottom: 4px;
}

.auditor-info {
  font-size: 12px;
  color: #7f8c8d;
  text-align: right;
}

.card-footer {
  display: flex;
  justify-content: space-between;
  align-items: flex-end;
  padding-top: 12px;
  border-top: 1px solid #f0f0f0;
}

.time-info {
  display: flex;
  flex-direction: column;
  gap: 4px;
  flex: 1;
}

.created-time, .updated-time {
  display: flex;
  align-items: center;
  gap: 4px;
}

.time-label {
  font-size: 11px;
  color: #7f8c8d;
}

.time-value {
  font-size: 11px;
  color: #2c3e50;
}

.actions {
  display: flex;
  gap: 8px;
  flex-shrink: 0;
}

.empty-state {
  text-align: center;
  padding: 60px 20px;
  background: white;
  border-radius: 12px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
}

.empty-icon {
  font-size: 48px;
  margin-bottom: 16px;
}

.empty-title {
  font-size: 18px;
  font-weight: 600;
  color: #2c3e50;
  margin-bottom: 8px;
}

.empty-desc {
  font-size: 14px;
  color: #7f8c8d;
  margin-bottom: 20px;
  line-height: 1.5;
}

.empty-action {
  padding: 12px 24px;
  border-radius: 20px;
}

.load-more {
  text-align: center;
  padding: 20px;
}

.load-more-btn {
  color: #27ae60;
  font-weight: 500;
}

.load-complete {
  text-align: center;
  padding: 20px;
  color: #7f8c8d;
  font-size: 14px;
}

.loading-container {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  height: 200px;
  color: #7f8c8d;
}

.loading-container .el-icon {
  font-size: 32px;
  margin-bottom: 12px;
}

.bottom-actions {
  padding: 16px;
  background: white;
  box-shadow: 0 -2px 8px rgba(0, 0, 0, 0.1);
}

.apply-btn {
  width: 100%;
  height: 50px;
  font-size: 16px;
  font-weight: 600;
  border-radius: 25px;
}

@media (max-width: 768px) {
  .page-content {
    padding: 12px;
  }
  
  .stats-card {
    grid-template-columns: repeat(2, 1fr);
    padding: 16px;
    gap: 12px;
  }
  
  .stat-value {
    font-size: 20px;
  }
  
  .application-card {
    padding: 16px;
  }
  
  .card-header {
    flex-direction: column;
    gap: 8px;
    align-items: flex-start;
  }
  
  .amount-section {
    grid-template-columns: 1fr;
    gap: 8px;
  }
  
  .amount-info, .duration-info {
    flex-direction: row;
    justify-content: space-between;
  }
  
  .card-footer {
    flex-direction: column;
    gap: 12px;
    align-items: flex-start;
  }
  
  .actions {
    width: 100%;
    justify-content: flex-end;
  }
}
</style> 