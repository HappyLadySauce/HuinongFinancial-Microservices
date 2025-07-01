<script setup lang="ts">
import { ref, onMounted, computed } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { ElMessage, ElMessageBox } from 'element-plus'
import { useUserStore } from '../stores/user'
import { loanApprovalApi } from '../services/api'
import type { LoanApplication } from '../services/api'

const router = useRouter()
const route = useRoute()
const userStore = useUserStore()
const loading = ref(false)

// 申请详情
const application = ref<LoanApplication | null>(null)

// 状态映射
const statusMap = {
  'pending': { 
    text: '待审批', 
    color: '#E6A23C', 
    icon: '⏳',
    bgColor: '#fef6e7',
    description: '您的申请正在审核中，请耐心等待...'
  },
  'approved': { 
    text: '已批准', 
    color: '#67C23A', 
    icon: '✅',
    bgColor: '#e8f5e8',
    description: '恭喜！您的申请已获得批准'
  },
  'rejected': { 
    text: '已拒绝', 
    color: '#F56C6C', 
    icon: '❌',
    bgColor: '#ffebee',
    description: '很抱歉，您的申请未能通过审核'
  }
}

// 获取状态信息
const getStatusInfo = (status: string) => {
  return statusMap[status as keyof typeof statusMap] || {
    text: status,
    color: '#909399',
    icon: '❓',
    bgColor: '#f5f5f5',
    description: '申请状态未知'
  }
}

// 格式化时间 - 支持Unix时间戳
const formatDateTime = (timestamp: number) => {
  const date = new Date(timestamp * 1000) // Unix时间戳转换为毫秒
  return date.toLocaleString('zh-CN', {
    year: 'numeric',
    month: '2-digit',
    day: '2-digit',
    hour: '2-digit',
    minute: '2-digit'
  })
}

// 获取申请详情
const loadApplicationDetail = async () => {
  const applicationId = route.params.id as string
  if (!applicationId) {
    // 静默跳转到申请列表页面，不显示错误提示
    router.push('/loan/my-applications?from=/loan/application')
    return
  }

  try {
    loading.value = true
    const response = await loanApprovalApi.getDetail(applicationId)
    application.value = response.application_info // 新API返回格式为 { application_info: LoanApplication }
  } catch (error: any) {
    console.error('加载申请详情失败:', error)
    ElMessage.error('加载申请详情失败')
    router.go(-1)
  } finally {
    loading.value = false
  }
}

// 删除申请 - 使用cancel方法
const deleteApplication = async () => {
  if (!application.value || application.value.status !== 'pending') {
    ElMessage.warning('只能删除待审批的申请')
    return
  }

  try {
    await ElMessageBox.confirm(
      `确认删除申请"${application.value.name}"吗？删除后无法恢复。`,
      '确认删除',
      {
        confirmButtonText: '确认删除',
        cancelButtonText: '取消',
        type: 'warning'
      }
    )

    // 新API使用cancel方法取消申请
    await loanApprovalApi.cancel(application.value.application_id, '用户主动删除申请')
    ElMessage.success('删除成功')
    
    // 返回列表页，传递来源信息
    router.push('/loan/my-applications?from=/loan/application')

  } catch (error: any) {
    if (error.message !== 'cancel') {
      console.error('删除申请失败:', error)
      ElMessage.error(error.message || '删除申请失败')
    }
  }
}

// 编辑申请
const editApplication = () => {
  if (!application.value || application.value.status !== 'pending') {
    ElMessage.warning('只能编辑待审批的申请')
    return
  }
  
  ElMessage.info('编辑功能暂未实现')
}

// 返回上一页
const goBack = () => {
  router.go(-1)
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
  
  loadApplicationDetail()
})
</script>

<template>
  <div class="application-detail-page">
    <!-- 顶部导航 -->
    <div class="top-nav">
      <div class="nav-left" @click="goBack">
        <el-icon><ArrowLeft /></el-icon>
      </div>
      <div class="nav-title">申请详情</div>
      <div class="nav-right">
        <el-icon @click="loadApplicationDetail" :class="{ 'is-loading': loading }">
          <Refresh />
        </el-icon>
      </div>
    </div>

    <!-- 加载状态 -->
    <div v-if="loading" class="loading-container">
      <el-icon class="is-loading"><Loading /></el-icon>
      <p>加载中...</p>
    </div>

    <!-- 申请详情内容 -->
    <div v-if="!loading && application" class="page-content">
      <!-- 状态卡片 -->
      <div 
        class="status-card"
        :style="{ backgroundColor: getStatusInfo(application.status).bgColor }"
      >
        <div class="status-header">
          <div class="status-icon">{{ getStatusInfo(application.status).icon }}</div>
          <div class="status-content">
            <div 
              class="status-text"
              :style="{ color: getStatusInfo(application.status).color }"
            >
              {{ getStatusInfo(application.status).text }}
            </div>
            <div class="status-desc">
              {{ getStatusInfo(application.status).description }}
            </div>
          </div>
        </div>
        
        <!-- 新API暂时没有auditor字段，先隐藏 -->
        <!--
        <div v-if="application.auditor" class="auditor-section">
          <div class="auditor-info">
            <span class="auditor-label">审批人:</span>
            <span class="auditor-name">{{ application.auditor }}</span>
          </div>
          <div v-if="application.updated_at !== application.created_at" class="audit-time">
            <span class="time-label">审批时间:</span>
            <span class="time-value">{{ formatDateTime(application.updated_at) }}</span>
        </div>
        </div>
        -->
      </div>

      <!-- 基本信息 -->
      <div class="info-card">
        <div class="card-header">
          <h3 class="card-title">申请信息</h3>
          <div class="app-id">ID: {{ application.application_id }}</div>
        </div>
        
        <div class="info-content">
          <div class="info-section">
            <div class="section-title">基本信息</div>
            <div class="info-grid">
          <div class="info-item">
                <span class="info-label">申请名称</span>
                <span class="info-value">{{ application.name }}</span>
          </div>
          <div class="info-item">
                <span class="info-label">贷款类型</span>
                <span class="info-value">
                  <span class="type-badge">{{ application.type }}</span>
                </span>
          </div>
          <div class="info-item">
                <span class="info-label">申请金额</span>
                <span class="info-value amount">¥{{ application.amount.toLocaleString() }}</span>
          </div>
          <div class="info-item">
                <span class="info-label">贷款期限</span>
                <span class="info-value">{{ application.duration }}个月</span>
              </div>
            </div>
          </div>

          <div class="info-section">
            <div class="section-title">申请用途</div>
            <div class="description-content">
              {{ application.purpose }}
          </div>
          </div>

          <!-- 新API暂时没有suggestions字段，先隐藏 -->
          <!--
          <div v-if="application.suggestions" class="info-section">
            <div 
              class="section-title"
              :style="{ color: application.status === 'approved' ? '#27ae60' : '#e74c3c' }"
            >
              {{ application.status === 'approved' ? '审批意见' : '拒绝原因' }}
          </div>
            <div 
              class="suggestions-content"
              :class="{ 
                'approved-suggestion': application.status === 'approved',
                'rejected-suggestion': application.status === 'rejected' 
              }"
            >
              {{ application.suggestions }}
        </div>
      </div>
      -->

          <div class="info-section">
            <div class="section-title">时间记录</div>
            <div class="time-grid">
              <div class="time-item">
                <span class="time-label">申请时间</span>
                <span class="time-value">{{ formatDateTime(application.created_at) }}</span>
              </div>
              <div v-if="application.updated_at !== application.created_at" class="time-item">
                <span class="time-label">更新时间</span>
                <span class="time-value">{{ formatDateTime(application.updated_at) }}</span>
              </div>
            </div>
          </div>
        </div>
      </div>

      <!-- 操作按钮 -->
      <div class="actions-card">
        <div class="actions-header">
          <h3>可用操作</h3>
        </div>
        
        <div class="actions-content">
          <el-button 
            v-if="application.status === 'pending'"
            type="warning" 
            size="large"
            @click="editApplication"
            class="action-btn"
          >
            <el-icon><Edit /></el-icon>
            编辑申请
          </el-button>
          
        <el-button 
            v-if="application.status === 'pending'"
            type="danger" 
          size="large"
            @click="deleteApplication"
          class="action-btn"
        >
            <el-icon><Delete /></el-icon>
            删除申请
        </el-button>
        
        <el-button 
            type="primary" 
          size="large"
            @click="applyNewLoan"
          class="action-btn"
        >
            <el-icon><Plus /></el-icon>
            申请新贷款
        </el-button>
        </div>
      </div>

      <!-- 申请流程 -->
      <div class="timeline-card">
        <div class="card-header">
          <h3 class="card-title">申请流程</h3>
        </div>
        
        <div class="timeline-content">
          <div class="timeline-item completed">
            <div class="timeline-dot"></div>
            <div class="timeline-content-item">
              <div class="timeline-title">提交申请</div>
              <div class="timeline-time">{{ formatDateTime(application.created_at) }}</div>
              <div class="timeline-desc">申请已成功提交，等待审核</div>
            </div>
          </div>
          
          <div 
            class="timeline-item"
            :class="{ 
              'completed': application.status !== 'pending',
              'rejected': application.status === 'rejected'
            }"
          >
            <div class="timeline-dot"></div>
            <div class="timeline-content-item">
              <div class="timeline-title">
                {{ application.status === 'approved' ? '审核通过' : 
                   application.status === 'rejected' ? '审核拒绝' : '审核中' }}
              </div>
              <div v-if="application.status !== 'pending'" class="timeline-time">
                {{ formatDateTime(application.updated_at) }}
              </div>
              <div class="timeline-desc">
                {{ application.status === 'approved' ? '申请已通过审核，贷款即将发放' :
                   application.status === 'rejected' ? '申请未通过审核，请查看拒绝原因' :
                   '正在进行审核，请耐心等待' }}
              </div>
            </div>
          </div>
          
          <div 
            v-if="application.status === 'approved'"
            class="timeline-item pending"
          >
            <div class="timeline-dot"></div>
            <div class="timeline-content-item">
              <div class="timeline-title">放款</div>
              <div class="timeline-desc">等待银行放款</div>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- 错误状态 -->
    <div v-if="!loading && !application" class="error-state">
      <div class="error-icon">❌</div>
      <div class="error-title">申请不存在</div>
      <div class="error-desc">未找到相关申请信息，请检查申请ID是否正确</div>
      <el-button type="primary" @click="goBack">返回上一页</el-button>
    </div>
  </div>
</template>

<style scoped>
.application-detail-page {
  min-height: 100vh;
  background-color: #f5f5f5;
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
  padding: 16px;
  max-width: 600px;
  margin: 0 auto;
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.status-card {
  border-radius: 12px;
  padding: 20px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
}

.status-header {
  display: flex;
  align-items: center;
  gap: 16px;
  margin-bottom: 16px;
}

.status-icon {
  font-size: 32px;
  width: 50px;
  height: 50px;
  display: flex;
  align-items: center;
  justify-content: center;
  background: rgba(255, 255, 255, 0.9);
  border-radius: 50%;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
}

.status-content {
  flex: 1;
}

.status-text {
  font-size: 20px;
  font-weight: 600;
  margin-bottom: 4px;
}

.status-desc {
  font-size: 14px;
  color: #666;
  line-height: 1.5;
}

.auditor-section {
  padding-top: 16px;
  border-top: 1px solid rgba(255, 255, 255, 0.5);
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.auditor-info, .audit-time {
  display: flex;
  align-items: center;
  gap: 4px;
}

.auditor-label, .time-label {
  font-size: 12px;
  color: #666;
}

.auditor-name, .time-value {
  font-size: 13px;
  color: #2c3e50;
  font-weight: 500;
}

.info-card, .actions-card, .timeline-card {
  background: white;
  border-radius: 12px;
  padding: 20px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 20px;
  padding-bottom: 12px;
  border-bottom: 1px solid #f0f0f0;
}

.card-title {
  font-size: 18px;
  font-weight: 600;
  color: #2c3e50;
  margin: 0;
}

.app-id {
  font-size: 12px;
  color: #7f8c8d;
  background: #f8f9fa;
  padding: 4px 8px;
  border-radius: 12px;
}

.info-section {
  margin-bottom: 24px;
}

.info-section:last-child {
  margin-bottom: 0;
}

.section-title {
  font-size: 16px;
  font-weight: 600;
  color: #2c3e50;
  margin-bottom: 12px;
  padding-bottom: 6px;
  border-bottom: 2px solid #27ae60;
  position: relative;
}

.section-title::before {
  content: '';
  position: absolute;
  left: 0;
  bottom: -2px;
  width: 30px;
  height: 2px;
  background: #27ae60;
}

.info-grid {
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  gap: 16px;
}

.info-item {
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.info-label {
  font-size: 13px;
  color: #7f8c8d;
}

.info-value {
  font-size: 15px;
  color: #2c3e50;
  font-weight: 500;
}

.info-value.amount {
  color: #27ae60;
  font-size: 18px;
  font-weight: 600;
}

.type-badge {
  display: inline-flex;
  align-items: center;
  padding: 4px 12px;
  background: #e8f5e8;
  color: #27ae60;
  border-radius: 16px;
  font-size: 13px;
  font-weight: 500;
}

.description-content {
  background: #f8f9fa;
  padding: 16px;
  border-radius: 8px;
  font-size: 14px;
  color: #2c3e50;
  line-height: 1.6;
  border-left: 4px solid #27ae60;
}

.suggestions-content {
  padding: 16px;
  border-radius: 8px;
  font-size: 14px;
  color: #2c3e50;
  line-height: 1.6;
  border-left: 4px solid;
}

.suggestions-content.approved-suggestion {
  background: #e8f5e8;
  border-left-color: #27ae60;
}

.suggestions-content.rejected-suggestion {
  background: #ffebee;
  border-left-color: #e74c3c;
}

.time-grid {
  display: grid;
  grid-template-columns: 1fr;
  gap: 12px;
}

.time-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 12px 16px;
  background: #f8f9fa;
  border-radius: 8px;
}

.actions-header {
  margin-bottom: 16px;
}

.actions-header h3 {
  font-size: 16px;
  font-weight: 600;
  color: #2c3e50;
  margin: 0;
}

.actions-content {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.action-btn {
  width: 100%;
  height: 50px;
  font-size: 16px;
  font-weight: 500;
  border-radius: 25px;
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 8px;
}

.timeline-content {
  position: relative;
}

.timeline-content::before {
  content: '';
  position: absolute;
  left: 15px;
  top: 24px;
  bottom: 24px;
  width: 2px;
  background: #e1e1e1;
}

.timeline-item {
  position: relative;
  padding: 16px 0 16px 48px;
}

.timeline-dot {
  position: absolute;
  left: 0;
  top: 20px;
  width: 32px;
  height: 32px;
  border-radius: 50%;
  background: #e1e1e1;
  border: 4px solid white;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
}

.timeline-item.completed .timeline-dot {
  background: #27ae60;
}

.timeline-item.rejected .timeline-dot {
  background: #e74c3c;
}

.timeline-item.pending .timeline-dot {
  background: #E6A23C;
}

.timeline-content-item {
  padding-left: 12px;
}

.timeline-title {
  font-size: 16px;
  font-weight: 600;
  color: #2c3e50;
  margin-bottom: 4px;
}

.timeline-time {
  font-size: 12px;
  color: #7f8c8d;
  margin-bottom: 4px;
}

.timeline-desc {
  font-size: 14px;
  color: #666;
  line-height: 1.5;
}

.loading-container {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  height: 300px;
  color: #7f8c8d;
}

.loading-container .el-icon {
  font-size: 32px;
  margin-bottom: 12px;
}

.error-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  height: 400px;
  padding: 20px;
  text-align: center;
}

.error-icon {
  font-size: 48px;
  margin-bottom: 16px;
}

.error-title {
  font-size: 20px;
  font-weight: 600;
  color: #2c3e50;
  margin-bottom: 8px;
}

.error-desc {
  font-size: 14px;
  color: #7f8c8d;
  margin-bottom: 20px;
  line-height: 1.5;
}

@media (max-width: 768px) {
  .page-content {
    padding: 12px;
}

  .status-card, .info-card, .actions-card, .timeline-card {
    padding: 16px;
  }
  
  .status-header {
    flex-direction: column;
    text-align: center;
    gap: 12px;
  }
  
  .auditor-section {
    flex-direction: column;
    gap: 8px;
    align-items: flex-start;
  }
  
  .info-grid {
    grid-template-columns: 1fr;
    gap: 12px;
  }
  
  .time-item {
    flex-direction: column;
    align-items: flex-start;
    gap: 4px;
  }
  
  .actions-content {
    gap: 8px;
  }
  
  .action-btn {
    height: 45px;
    font-size: 15px;
  }
}
</style> 