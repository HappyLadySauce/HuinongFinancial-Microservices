<script setup lang="ts">
import { ref, onMounted, computed } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { ElMessage } from 'element-plus'
import { useUserStore } from '../stores/user'
import { loanApi } from '../services/api'
import type { LoanProduct } from '../services/api'

const router = useRouter()
const route = useRoute()
const userStore = useUserStore()
const loading = ref(false)

// 产品详情
const product = ref<LoanProduct | null>(null)

// 格式化金额
const formatAmount = (amount: number) => {
  if (amount >= 10000) {
    return `${(amount / 10000).toFixed(1)}万`
  }
  return amount.toLocaleString()
}

// 格式化期限选项
const termOptions = computed(() => {
  if (!product.value) return []
  const options = []
  for (let i = product.value.min_term_months; i <= product.value.max_term_months; i += 6) {
    options.push(`${i}个月`)
  }
  return options
})

// 加载产品详情
const loadProductDetail = async () => {
  const productId = route.params.productId as string
  if (!productId) {
    ElMessage.error('缺少产品信息')
    router.go(-1)
    return
  }

  try {
    loading.value = true
    const response = await loanApi.getProductDetail(productId)
    product.value = response.data
  } catch (error: any) {
    console.error('加载产品详情失败:', error)
    ElMessage.error('加载产品详情失败')
    router.go(-1)
  } finally {
    loading.value = false
  }
}

// 申请贷款
const applyLoan = () => {
  if (!product.value) return
  
  // 检查登录状态
  if (!userStore.isLoggedIn || !userStore.isTokenValid()) {
    ElMessage.warning('请先登录')
    router.push('/login')
    return
  }
  
  // 跳转到申请页面
  router.push(`/loan/apply/${product.value.product_id}`)
}

// 返回上一页
const goBack = () => {
  router.go(-1)
}

// 组件挂载时加载数据
onMounted(() => {
  loadProductDetail()
})
</script>

<template>
  <div class="product-detail-page">
    <!-- 顶部导航 -->
    <div class="top-nav">
      <div class="nav-left" @click="goBack">
        <el-icon><ArrowLeft /></el-icon>
      </div>
      <div class="nav-title">产品详情</div>
      <div class="nav-right"></div>
    </div>

    <div class="page-content" v-if="product && !loading">
      <!-- 产品基本信息 -->
      <div class="product-header">
        <div class="product-info">
          <h1 class="product-name">{{ product.name }}</h1>
          <p class="product-desc">{{ product.description }}</p>
          <div class="product-tags">
            <el-tag type="success" size="small">{{ product.category }}</el-tag>
            <el-tag v-if="product.status === 0" type="success" size="small">可申请</el-tag>
            <el-tag v-else type="danger" size="small">暂停申请</el-tag>
          </div>
        </div>
      </div>

      <!-- 核心信息卡片 -->
      <div class="highlight-card">
        <div class="highlight-item">
          <div class="highlight-label">年利率</div>
          <div class="highlight-value">{{ product.interest_rate_yearly }}</div>
        </div>
        <div class="divider"></div>
        <div class="highlight-item">
          <div class="highlight-label">可贷金额</div>
          <div class="highlight-value">{{ formatAmount(product.min_amount) }}-{{ formatAmount(product.max_amount) }}元</div>
        </div>
      </div>

      <!-- 产品详细信息 -->
      <div class="detail-sections">
        <!-- 基本信息 -->
        <div class="detail-section">
          <h3 class="section-title">基本信息</h3>
          <div class="detail-list">
            <div class="detail-item">
              <span class="label">产品类型</span>
              <span class="value">{{ product.category }}</span>
            </div>
            <div class="detail-item">
              <span class="label">贷款金额</span>
              <span class="value">{{ formatAmount(product.min_amount) }} - {{ formatAmount(product.max_amount) }}元</span>
            </div>
            <div class="detail-item">
              <span class="label">贷款期限</span>
              <span class="value">{{ product.min_term_months }} - {{ product.max_term_months }}个月</span>
            </div>
            <div class="detail-item">
              <span class="label">年利率</span>
              <span class="value highlight">{{ product.interest_rate_yearly }}</span>
            </div>
            <div v-if="product.repayment_methods" class="detail-item">
              <span class="label">还款方式</span>
              <span class="value">{{ product.repayment_methods.join('、') }}</span>
            </div>
          </div>
        </div>

        <!-- 申请条件 -->
        <div v-if="product.application_conditions" class="detail-section">
          <h3 class="section-title">申请条件</h3>
          <div class="conditions-content">
            <p>{{ product.application_conditions }}</p>
          </div>
        </div>

        <!-- 所需材料 -->
        <div v-if="product.required_documents && product.required_documents.length > 0" class="detail-section">
          <h3 class="section-title">所需材料</h3>
          <div class="documents-list">
            <div 
              v-for="(doc, index) in product.required_documents" 
              :key="index"
              class="document-item"
            >
              <el-icon class="doc-icon"><Document /></el-icon>
              <span class="doc-desc">{{ doc.desc }}</span>
            </div>
          </div>
        </div>

        <!-- 期限选项 -->
        <div class="detail-section">
          <h3 class="section-title">可选期限</h3>
          <div class="term-options">
            <div 
              v-for="term in termOptions" 
              :key="term"
              class="term-option"
            >
              {{ term }}
            </div>
          </div>
        </div>
      </div>

      <!-- 申请按钮 -->
      <div class="apply-section">
        <el-button 
          type="primary" 
          size="large" 
          @click="applyLoan"
          :disabled="product.status !== 0"
          class="apply-btn"
        >
          {{ product.status === 0 ? '立即申请' : '暂停申请' }}
        </el-button>
        <p class="apply-tip">
          <el-icon><InfoFilled /></el-icon>
          申请提交后将进入AI智能审核，通常1-3个工作日内完成审核
        </p>
      </div>
    </div>

    <!-- 加载状态 -->
    <div v-if="loading" class="loading-container">
      <el-icon class="is-loading"><Loading /></el-icon>
      <p>加载中...</p>
    </div>
  </div>
</template>

<style scoped>
.product-detail-page {
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

.nav-left {
  cursor: pointer;
  padding: 8px;
}

.nav-title {
  font-size: 18px;
  font-weight: 600;
  color: #2c3e50;
}

.nav-right {
  width: 32px;
}

.page-content {
  padding: 16px;
  max-width: 600px;
  margin: 0 auto;
}

.product-header {
  background: white;
  border-radius: 12px;
  padding: 20px;
  margin-bottom: 16px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
}

.product-name {
  font-size: 24px;
  font-weight: 600;
  color: #2c3e50;
  margin: 0 0 12px 0;
}

.product-desc {
  font-size: 16px;
  color: #7f8c8d;
  line-height: 1.6;
  margin: 0 0 16px 0;
}

.product-tags {
  display: flex;
  gap: 8px;
}

.highlight-card {
  background: linear-gradient(135deg, #27ae60, #2ecc71);
  border-radius: 12px;
  padding: 24px;
  margin-bottom: 16px;
  color: white;
  display: flex;
  align-items: center;
  justify-content: space-around;
}

.highlight-item {
  text-align: center;
  flex: 1;
}

.highlight-label {
  font-size: 14px;
  opacity: 0.9;
  margin-bottom: 8px;
}

.highlight-value {
  font-size: 20px;
  font-weight: 600;
}

.divider {
  width: 1px;
  height: 40px;
  background: rgba(255, 255, 255, 0.3);
  margin: 0 20px;
}

.detail-sections {
  margin-bottom: 16px;
}

.detail-section {
  background: white;
  border-radius: 12px;
  padding: 20px;
  margin-bottom: 16px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
}

.section-title {
  font-size: 16px;
  font-weight: 600;
  color: #2c3e50;
  margin: 0 0 16px 0;
  padding-bottom: 8px;
  border-bottom: 2px solid #27ae60;
}

.detail-list {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.detail-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 8px 0;
}

.label {
  color: #7f8c8d;
  font-size: 14px;
}

.value {
  color: #2c3e50;
  font-weight: 500;
  text-align: right;
}

.value.highlight {
  color: #27ae60;
  font-weight: 600;
  font-size: 16px;
}

.conditions-content {
  color: #2c3e50;
  line-height: 1.6;
}

.conditions-content p {
  margin: 0;
}

.documents-list {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.document-item {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 8px 0;
}

.doc-icon {
  color: #27ae60;
  font-size: 18px;
}

.doc-desc {
  color: #2c3e50;
  font-size: 14px;
}

.term-options {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
}

.term-option {
  padding: 8px 16px;
  background: #f8f9fa;
  border: 1px solid #e1e1e1;
  border-radius: 20px;
  font-size: 14px;
  color: #2c3e50;
}

.apply-section {
  background: white;
  border-radius: 12px;
  padding: 20px;
  text-align: center;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
}

.apply-btn {
  width: 100%;
  height: 50px;
  font-size: 16px;
  font-weight: 600;
  border-radius: 25px;
  margin-bottom: 12px;
}

.apply-tip {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 4px;
  color: #7f8c8d;
  font-size: 12px;
  margin: 0;
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

@media (max-width: 768px) {
  .page-content {
    padding: 12px;
  }
  
  .product-header {
    padding: 16px;
  }
  
  .highlight-card {
    padding: 20px 16px;
    flex-direction: column;
    gap: 16px;
  }
  
  .divider {
    width: 100%;
    height: 1px;
    margin: 0;
  }
  
  .detail-section {
    padding: 16px;
  }
  
  .apply-section {
    padding: 16px;
  }
}
</style> 