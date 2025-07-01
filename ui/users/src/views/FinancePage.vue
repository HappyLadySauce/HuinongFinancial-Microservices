<script setup lang="ts">
import { ref, onMounted, computed } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import AppFooter from './components/footer.vue'
import { useUserStore } from '../stores/user'
import { loanApprovalApi } from '../services/api'
// 移除不存在的ProductTypes接口

const router = useRouter()
const userStore = useUserStore()
const activeTab = ref('finance')
const loading = ref(false)
const refreshing = ref(false)

// 贷款类型列表
const loanTypes = ref<string[]>([])
const selectedType = ref('')

// 我的申请统计数据
const myStats = ref({
  totalApplications: 0,
  pendingCount: 0,
  approvedAmount: 0
})

// 贷款类型配置
const loanTypeConfigs = {
  '农业贷': {
    description: '专为农业生产提供资金支持',
    icon: '🌾',
    color: '#27ae60',
    features: ['利率优惠', '期限灵活', '快速审批'],
    amountRange: '1万-50万',
    termRange: '6-36个月'
  },
  '创业贷': {
    description: '支持创新创业项目发展',
    icon: '🚀',
    color: '#3498db',
    features: ['低门槛', '高额度', '政策扶持'],
    amountRange: '5万-100万',
    termRange: '12-60个月'
  },
  '消费贷': {
    description: '满足日常消费需求',
    icon: '🛒',
    color: '#e74c3c',
    features: ['免抵押', '快放款', '随借随还'],
    amountRange: '1万-20万',
    termRange: '3-24个月'
  },
  '经营贷': {
    description: '助力企业经营发展',
    icon: '🏢',
    color: '#f39c12',
    features: ['额度高', '期限长', '循环使用'],
    amountRange: '10万-500万',
    termRange: '12-60个月'
  },
  '助学贷': {
    description: '教育投资，成就未来',
    icon: '🎓',
    color: '#9b59b6',
    features: ['利率低', '还款宽松', '政府贴息'],
    amountRange: '1万-10万',
    termRange: '12-120个月'
  }
}

// 筛选后的贷款类型
const filteredTypes = computed(() => {
  if (!selectedType.value) return loanTypes.value
  return loanTypes.value.filter(type => type === selectedType.value)
})

// 获取类型配置
const getTypeConfig = (type: string) => {
  return loanTypeConfigs[type as keyof typeof loanTypeConfigs] || {
    description: '专业金融服务',
    icon: '💰',
    color: '#7f8c8d',
    features: ['专业服务', '安全可靠'],
    amountRange: '面议',
    termRange: '面议'
  }
}

// 加载贷款类型 - 使用静态数据替代API调用
const loadLoanTypes = async () => {
  try {
    loading.value = true
    // 新的API暂时没有getTypes方法，使用静态数据
    loanTypes.value = ['农业贷', '创业贷', '消费贷', '经营贷', '助学贷']
  } catch (error: any) {
    console.error('加载贷款类型失败:', error)
    ElMessage.error('加载产品信息失败')
  } finally {
    loading.value = false
  }
}

// 加载我的申请统计
const loadMyStats = async () => {
  if (!userStore.isLoggedIn) return
  
  try {
    const response = await loanApprovalApi.getMyApprovals({ page: 1, size: 100 })
    myStats.value = {
      totalApplications: response.total,
      pendingCount: response.list.filter(app => app.status === 'pending').length,
      approvedAmount: response.list
        .filter(app => app.status === 'approved')
        .reduce((sum, app) => sum + app.amount, 0)
    }
  } catch (error: any) {
    console.error('加载统计数据失败:', error)
  }
}

// 刷新数据
const refreshData = async () => {
  try {
    refreshing.value = true
    await loadLoanTypes()
    ElMessage.success('刷新成功')
  } catch (error: any) {
    ElMessage.error('刷新失败')
  } finally {
    refreshing.value = false
  }
}

// 切换分类
const switchType = (type: string) => {
  selectedType.value = type
}

// 申请贷款
const applyLoan = (loanType: string) => {
  // 检查登录状态
  if (!userStore.isLoggedIn || !userStore.isTokenValid()) {
    ElMessage.warning('请先登录')
    router.push('/login')
    return
  }
  
  // 跳转到申请页面，传递贷款类型
  router.push(`/loan/apply?type=${encodeURIComponent(loanType)}`)
}

// 查看我的申请
const viewMyApplications = () => {
  if (!userStore.isLoggedIn) {
    ElMessage.warning('请先登录')
    router.push('/login')
    return
  }
  router.push('/loan/my-applications?from=/finance')
}

// 组件挂载时加载数据
onMounted(() => {
  loadLoanTypes()
})
</script>

<template>
  <div class="finance-page">
    <!-- 顶部导航 -->
    <div class="top-nav">
      <div class="nav-left">
        <el-icon @click="router.go(-1)"><ArrowLeft /></el-icon>
      </div>
      <div class="nav-title">惠农金融</div>
      <div class="nav-right">
        <el-icon @click="refreshData" :class="{ 'is-loading': refreshing }">
          <Refresh />
        </el-icon>
      </div>
    </div>

    <div class="page-content">
      <!-- 用户快捷操作 -->
      <div class="quick-actions" v-if="userStore.isLoggedIn">
        <div class="action-card primary" @click="viewMyApplications">
          <div class="card-icon">
            <svg viewBox="0 0 24 24" width="22" height="22" fill="white">
              <path d="M14,2H6C4.9,2,4,2.9,4,4v16c0,1.1,0.9,2,2,2h12c1.1,0,2-0.9,2-2V8L14,2z M16,18H8v-2h8V18z M16,14H8v-2h8V14z M13,9V3.5L18.5,9H13z"/>
            </svg>
          </div>
          <div class="card-content">
            <h3>我的申请</h3>
            <p>查看贷款申请进度</p>
          </div>
          <div class="card-arrow">
            <svg viewBox="0 0 24 24" width="18" height="18" fill="white" opacity="0.8">
              <path d="M8.59,16.59L13.17,12L8.59,7.41L10,6l6,6l-6,6L8.59,16.59z"/>
            </svg>
          </div>
        </div>
      </div>

      <!-- 登录提示 -->
      <div class="login-prompt" v-else>
        <div class="prompt-content">
          <el-icon class="prompt-icon"><User /></el-icon>
          <p>登录后享受更多金融服务</p>
          <el-button type="primary" @click="router.push('/login')">
            立即登录
          </el-button>
        </div>
      </div>

      <!-- 产品分类筛选 -->
      <div class="category-filter">
        <div class="filter-header">
          <span class="filter-icon">
            <svg viewBox="0 0 24 24" width="20" height="20" fill="#27ae60">
              <path d="M10,18V16H8V14H10V12H12V14H14V16H12V18H10M3,4H21V8H19.5C19.5,8 19,8.5 19,9C19,9.5 19.5,10 19.5,10H21V14H19.5C19.5,14 19,14.5 19,15C19,15.5 19.5,16 19.5,16H21V20H3V16H4.5C4.5,16 5,15.5 5,15C5,14.5 4.5,14 4.5,14H3V10H4.5C4.5,10 5,9.5 5,9C5,8.5 4.5,8 4.5,8H3V4Z"/>
            </svg>
          </span>
          <span class="filter-title">贷款类型</span>
        </div>
        <div class="category-tabs">
          <div 
            class="category-tab"
            :class="{ 'active': selectedType === '' }"
            @click="switchType('')"
          >
            <span class="tab-icon">📋</span>
            <span class="tab-text">全部类型</span>
          </div>
          <div 
            v-for="type in loanTypes" 
            :key="type"
            class="category-tab"
            :class="{ 'active': selectedType === type }"
            @click="switchType(type)"
          >
            <span class="tab-icon">{{ getTypeConfig(type).icon }}</span>
            <span class="tab-text">{{ type }}</span>
          </div>
        </div>
      </div>

      <!-- 贷款类型列表 -->
      <div class="products-section">
        <div class="section-header">
          <h3>可申请类型</h3>
          <span class="product-count">{{ filteredTypes.length }}种类型</span>
        </div>

        <!-- 加载状态 -->
        <div v-if="loading" class="loading-container">
          <el-icon class="is-loading"><Loading /></el-icon>
          <p>加载中...</p>
        </div>

        <!-- 贷款类型网格 -->
        <div v-else class="loan-types-grid">
          <div 
            v-for="type in filteredTypes" 
            :key="type"
            class="loan-type-card"
            @click="applyLoan(type)"
          >
            <div class="loan-icon">
              <span>{{ getTypeConfig(type).icon }}</span>
                </div>
            <div class="loan-content">
              <div class="loan-title">{{ type }}</div>
              <div class="loan-desc">{{ getTypeConfig(type).description }}</div>
              <div class="loan-features">
                <span 
                  v-for="feature in getTypeConfig(type).features.slice(0, 3)" 
                  :key="feature"
                  class="feature-tag"
                >
                  {{ feature }}
                </span>
              </div>
              <div class="loan-info">
                <div class="info-row">
                  <span class="info-label">贷款金额</span>
                  <span class="info-value">1千-2万</span>
            </div>
                <div class="info-row">
                  <span class="info-label">贷款期限</span>
                  <span class="info-value">{{ getTypeConfig(type).termRange }}</span>
                </div>
                </div>
              <div class="apply-button">
                立即申请{{ type }}
                </div>
            </div>
          </div>

          <!-- 空状态 -->
          <div v-if="filteredTypes.length === 0" class="empty-state">
            <el-empty description="暂无相关产品">
              <el-button type="primary" @click="switchType('')">
                查看全部类型
              </el-button>
            </el-empty>
          </div>
        </div>
      </div>

      <!-- 申请须知 -->
      <div class="notice-section">
        <div class="notice-header">
          <span class="notice-icon">📋</span>
          <span class="notice-title">申请须知</span>
        </div>
        <div class="notice-content">
          <div class="notice-item">
            <span class="notice-number">1</span>
            <span class="notice-text">填写真实有效的申请信息</span>
          </div>
          <div class="notice-item">
            <span class="notice-number">2</span>
            <span class="notice-text">申请提交后将进入AI智能审核</span>
          </div>
          <div class="notice-item">
            <span class="notice-number">3</span>
            <span class="notice-text">审核结果通常1-3个工作日内完成</span>
          </div>
          <div class="notice-item">
            <span class="notice-number">4</span>
            <span class="notice-text">如有疑问请联系客服咨询</span>
          </div>
        </div>
      </div>
    </div>

    <!-- 底部导航 -->
    <app-footer v-model:active-tab="activeTab" />
  </div>
</template>

<style scoped>
.finance-page {
  min-height: 100vh;
  background-color: #f5f5f5;
  padding-bottom: 80px;
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
  width: 32px;
  height: 32px;
  display: flex;
  align-items: center;
  justify-content: center;
  cursor: pointer;
  color: #2c3e50;
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
}

.quick-actions {
  margin-bottom: 16px;
}

.action-card {
  background: linear-gradient(135deg, #27ae60 0%, #2ecc71 100%);
  border-radius: 16px;
  padding: 22px;
  color: white;
  display: flex;
  align-items: center;
  cursor: pointer;
  transition: all 0.3s ease;
  margin-bottom: 16px;
  box-shadow: 0 8px 16px rgba(39, 174, 96, 0.2);
  position: relative;
  overflow: hidden;
}

.action-card:before {
  content: '';
  position: absolute;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  background: linear-gradient(45deg, rgba(255,255,255,0.1) 0%, rgba(255,255,255,0) 100%);
  z-index: 1;
}

.action-card:hover {
  transform: translateY(-3px);
  box-shadow: 0 12px 28px rgba(39, 174, 96, 0.3);
}

.card-icon {
  width: 48px;
  height: 48px;
  background: rgba(255, 255, 255, 0.2);
  border-radius: 12px;
  display: flex;
  align-items: center;
  justify-content: center;
  margin-right: 16px;
  position: relative;
  z-index: 2;
  box-shadow: 0 4px 8px rgba(0, 0, 0, 0.1);
}

.card-content {
  flex: 1;
  position: relative;
  z-index: 2;
}

.card-content h3 {
  margin: 0 0 6px;
  font-size: 20px;
  font-weight: 600;
  letter-spacing: 0.5px;
}

.card-content p {
  margin: 0;
  font-size: 14px;
  opacity: 0.9;
  font-weight: 400;
}

.card-arrow {
  opacity: 0.8;
  position: relative;
  z-index: 2;
  display: flex;
  align-items: center;
  justify-content: center;
}

.stats-row {
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  gap: 16px;
  background: white;
  border-radius: 16px;
  padding: 20px;
  box-shadow: 0 8px 16px rgba(0,0,0,0.08);
  position: relative;
  overflow: hidden;
}

.stats-row:before {
  content: '';
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  height: 4px;
  background: linear-gradient(90deg, #27ae60, #2ecc71);
}

.stat-item {
  text-align: center;
  position: relative;
  padding: 8px 0;
}

.stat-item:not(:last-child):after {
  content: '';
  position: absolute;
  right: -8px;
  top: 20%;
  height: 60%;
  width: 1px;
  background: #f0f0f0;
}

.stat-value {
  font-size: 26px;
  font-weight: 700;
  color: #27ae60;
  margin-bottom: 8px;
  position: relative;
}

.stat-label {
  font-size: 13px;
  color: #666;
  font-weight: 500;
}

.stat-icon {
  position: absolute;
  top: 8px;
  right: 16px;
  opacity: 0.2;
}

.login-prompt {
  background: white;
  border-radius: 12px;
  padding: 24px;
  text-align: center;
  margin-bottom: 16px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
}

.prompt-content {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 12px;
}

.prompt-icon {
  font-size: 32px;
  color: #27ae60;
}

.category-filter {
  background: white;
  border-radius: 16px;
  padding: 20px;
  margin-bottom: 20px;
  box-shadow: 0 4px 16px rgba(0,0,0,0.08);
  position: relative;
  overflow: hidden;
}

.category-filter:before {
  content: '';
  position: absolute;
  top: 0;
  left: 0;
  width: 4px;
  height: 100%;
  background: linear-gradient(to bottom, #27ae60, #2ecc71);
}

.filter-header {
  display: flex;
  align-items: center;
  margin-bottom: 16px;
}

.filter-icon {
  display: flex;
  align-items: center;
  justify-content: center;
  margin-right: 10px;
}

.filter-title {
  font-size: 16px;
  font-weight: 600;
  color: #2c3e50;
}

.category-tabs {
  display: flex;
  gap: 12px;
  overflow-x: auto;
  scrollbar-width: none;
  padding-bottom: 4px;
  -webkit-overflow-scrolling: touch;
}

.category-tabs::-webkit-scrollbar {
  display: none;
}

.category-tab {
  padding: 8px 16px;
  border-radius: 12px;
  border: 2px solid #e1e1e1;
  font-size: 14px;
  font-weight: 500;
  cursor: pointer;
  white-space: nowrap;
  transition: all 0.3s;
  background: white;
  color: #666;
  display: flex;
  align-items: center;
  gap: 6px;
}

.category-tab.active {
  background: #27ae60;
  color: white;
  border-color: #27ae60;
  transform: translateY(-1px);
  box-shadow: 0 4px 8px rgba(39, 174, 96, 0.2);
}

.category-tab:hover:not(.active) {
  border-color: #27ae60;
  color: #27ae60;
  background-color: rgba(39, 174, 96, 0.05);
}

.tab-icon {
  display: flex;
  align-items: center;
  justify-content: center;
}

.products-section {
  background: white;
  border-radius: 12px;
  padding: 16px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
  margin-bottom: 16px;
}

.section-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 16px;
  padding-bottom: 12px;
  border-bottom: 1px solid #f0f0f0;
}

.section-header h3 {
  margin: 0;
  font-size: 16px;
  color: #2c3e50;
}

.product-count {
  font-size: 12px;
  color: #7f8c8d;
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

/* 贷款类型网格布局，参考首页风格 */
.loan-types-grid {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 15px;
  margin-top: 10px;
}

.loan-type-card {
  background: white;
  border-radius: 12px;
  padding: 16px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
  transition: all 0.3s ease;
  cursor: pointer;
  display: flex;
  flex-direction: column;
  height: auto;
}

.loan-type-card:hover {
  transform: translateY(-2px);
  box-shadow: 0 4px 16px rgba(0, 0, 0, 0.15);
}

.loan-icon {
  text-align: center;
  margin-bottom: 12px;
}

.loan-icon span {
  font-size: 32px;
  width: 50px;
  height: 50px;
  display: inline-flex;
  align-items: center;
  justify-content: center;
  background: linear-gradient(135deg, #4CAF50 0%, #27ae60 100%);
  border-radius: 50%;
  box-shadow: 0 2px 8px rgba(76, 175, 80, 0.3);
}

.loan-content {
  flex: 1;
  display: flex;
  flex-direction: column;
}

.loan-title {
  font-size: 16px;
  font-weight: 600;
  color: #333;
  text-align: center;
  margin-bottom: 8px;
}

.loan-desc {
  font-size: 12px;
  color: #666;
  text-align: center;
  margin-bottom: 12px;
  line-height: 1.4;
}

.loan-features {
  display: flex;
  justify-content: center;
  gap: 4px;
  flex-wrap: wrap;
  margin-bottom: 12px;
}

.feature-tag {
  font-size: 10px;
  color: #4CAF50;
  background: #e8f5e8;
  padding: 2px 6px;
  border-radius: 10px;
  border: 1px solid #4CAF50;
}

.loan-info {
  margin-bottom: 12px;
}

.info-row {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 4px;
}

.info-label {
  font-size: 12px;
  color: #999;
}

.info-value {
  font-size: 12px;
  color: #333;
  font-weight: 500;
}

.apply-button {
  background: linear-gradient(135deg, #4CAF50 0%, #27ae60 100%);
  color: white;
  text-align: center;
  padding: 8px 16px;
  border-radius: 20px;
  font-size: 14px;
  font-weight: 500;
  margin-top: auto;
  transition: all 0.3s ease;
}

.loan-type-card:hover .apply-button {
  transform: translateY(-1px);
  box-shadow: 0 4px 12px rgba(76, 175, 80, 0.4);
}

.empty-state {
  text-align: center;
  padding: 40px 20px;
}

.notice-section {
  background: white;
  border-radius: 12px;
  padding: 20px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
}

.notice-header {
  display: flex;
  align-items: center;
  margin-bottom: 16px;
  gap: 8px;
}

.notice-icon {
  font-size: 20px;
}

.notice-title {
  font-size: 16px;
  font-weight: 600;
  color: #2c3e50;
}

.notice-content {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.notice-item {
  display: flex;
  align-items: center;
  gap: 12px;
}

.notice-number {
  width: 24px;
  height: 24px;
  background: #27ae60;
  color: white;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 12px;
  font-weight: 600;
  flex-shrink: 0;
}

.notice-text {
  font-size: 14px;
  color: #2c3e50;
  line-height: 1.5;
}

@media (max-width: 768px) {
  .page-content {
    padding: 12px;
  }
  
  .action-card {
    padding: 18px;
  }
  
  .card-icon {
    width: 42px;
    height: 42px;
  }
  
  .card-content h3 {
    font-size: 18px;
  }
  
  .stats-row {
    padding: 16px 10px;
  }
  
  .stat-value {
    font-size: 22px;
  }
  
  .stat-label {
    font-size: 12px;
  }
  
  .stat-icon {
    display: none;
  }
  
  .category-filter {
    padding: 16px;
    margin-bottom: 16px;
  }
  
  .category-tab {
    padding: 6px 14px;
    font-size: 13px;
  }

  .loan-types-grid {
    gap: 10px;
  }
  
  .loan-type-card {
    padding: 12px;
  }
  
  .loan-icon span {
    font-size: 28px;
    width: 45px;
    height: 45px;
  }
  
  .loan-title {
    font-size: 15px;
  }
  
  .loan-desc {
    font-size: 11px;
  }
  
  .feature-tag {
    font-size: 9px;
  }
}
</style> 