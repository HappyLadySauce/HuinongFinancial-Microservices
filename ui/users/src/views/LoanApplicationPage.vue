<script setup lang="ts">
import { ref, reactive, onMounted, computed } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { ElMessage, ElMessageBox } from 'element-plus'
import type { FormInstance, FormRules } from 'element-plus'
import { useUserStore } from '../stores/user'
import { loanApprovalApi } from '../services/api'
import type { LoanApplicationRequest } from '../services/api'

const router = useRouter()
const route = useRoute()
const userStore = useUserStore()
const formRef = ref<FormInstance>()
const loading = ref(false)

// 可用的贷款类型
const availableTypes = ref<string[]>([])

// 申请表单数据 - 更新以匹配新API
const formData = reactive<LoanApplicationRequest>({
  product_id: 1, // 默认产品ID，实际应该从产品选择中获取
  name: '',
  type: '',
  amount: 10000,
  duration: 12,
  purpose: '' // 新API使用purpose而不是description
})

// 贷款类型配置
const typeConfigs = {
  '农业贷': {
    icon: '🌾',
    minAmount: 1000,
    maxAmount: 20000,
    minDuration: 6,
    maxDuration: 36,
    placeholder: '请详细说明农业生产用途...',
    defaultName: '农业生产贷款申请'
  },
  '创业贷': {
    icon: '🚀',
    minAmount: 1000,
    maxAmount: 20000,
    minDuration: 12,
    maxDuration: 60,
    placeholder: '请详细描述创业项目计划...',
    defaultName: '创业项目贷款申请'
  },
  '消费贷': {
    icon: '🛒',
    minAmount: 1000,
    maxAmount: 20000,
    minDuration: 3,
    maxDuration: 24,
    placeholder: '请说明消费用途...',
    defaultName: '消费贷款申请'
  },
  '经营贷': {
    icon: '🏢',
    minAmount: 1000,
    maxAmount: 20000,
    minDuration: 12,
    maxDuration: 60,
    placeholder: '请描述经营项目和资金用途...',
    defaultName: '经营贷款申请'
  },
  '助学贷': {
    icon: '🎓',
    minAmount: 1000,
    maxAmount: 20000,
    minDuration: 12,
    maxDuration: 120,
    placeholder: '请说明学习计划和资金用途...',
    defaultName: '助学贷款申请'
  }
}

// 当前选择类型的配置
const currentTypeConfig = computed(() => {
  if (!formData.type) return null
  return typeConfigs[formData.type as keyof typeof typeConfigs]
})

// 期限选项
const durationOptions = computed(() => {
  if (!currentTypeConfig.value) return []
  const options = []
  const { minDuration, maxDuration } = currentTypeConfig.value
  for (let i = minDuration; i <= maxDuration; i += 6) {
    options.push({
      value: i,
      label: `${i}个月`
    })
  }
  return options
})

// 表单验证规则
const rules = reactive<FormRules>({
  amount: [
    { required: true, message: '请输入贷款金额', trigger: 'blur' },
    { 
      validator: (rule: any, value: number, callback: any) => {
        if (!currentTypeConfig.value) {
          callback()
          return
        }
        const { minAmount, maxAmount } = currentTypeConfig.value
        if (value < minAmount) {
          callback(new Error(`贷款金额不能少于${minAmount.toLocaleString()}元`))
        } else if (value > maxAmount) {
          callback(new Error(`贷款金额不能超过${maxAmount.toLocaleString()}元`))
        } else {
          callback()
        }
      }, 
      trigger: 'blur' 
    }
  ],
  duration: [
    { required: true, message: '请选择贷款期限', trigger: 'change' }
  ],
  purpose: [
    { required: true, message: '请输入申请用途', trigger: 'blur' },
    { min: 20, message: '申请用途不能少于20个字符', trigger: 'blur' },
    { max: 500, message: '申请用途不能超过500个字符', trigger: 'blur' }
  ]
})

// 加载贷款类型 - 使用静态数据替代API调用
const loadLoanTypes = async () => {
  try {
    // 新的API暂时没有getTypes方法，使用静态数据
    availableTypes.value = ['农业贷', '创业贷', '消费贷', '经营贷', '助学贷']
    
    // 如果URL中有type参数，设置默认类型
    const typeFromQuery = route.query.type as string
    if (typeFromQuery) {
      // 检查类型是否有效
      if (availableTypes.value.includes(typeFromQuery)) {
        formData.type = typeFromQuery
        updateFormDataForType(typeFromQuery)
      } else {
        ElMessage.error('无效的贷款类型')
        router.push('/finance')
        return
      }
    } else {
      // 如果没有指定类型，返回理财页面
      ElMessage.warning('请先选择贷款类型')
      router.push('/finance')
      return
    }
  } catch (error: any) {
    console.error('加载贷款类型失败:', error)
    ElMessage.error('加载贷款类型失败')
  }
}

// 根据类型更新表单数据
const updateFormDataForType = (type: string) => {
  const config = typeConfigs[type as keyof typeof typeConfigs]
  if (config) {
    formData.amount = config.minAmount
    formData.duration = config.minDuration
    // 自动设置申请名称
    formData.name = config.defaultName
  }
}

// 提交申请
const submitApplication = async () => {
  if (!formRef.value) return
  
  try {
    // 验证表单
    const valid = await formRef.value.validate()
    if (!valid) return
    
    // 确认提交
    await ElMessageBox.confirm(
      `确认提交${formData.type}申请吗？申请金额：${formData.amount.toLocaleString()}元`,
      '确认提交',
      {
        confirmButtonText: '确认提交',
        cancelButtonText: '再检查一下',
        type: 'warning'
      }
    )
    
    loading.value = true
    
    // 提交申请
    const response = await loanApprovalApi.create(formData)
    
    ElMessage.success('申请提交成功')
    
    // 跳转到我的申请列表页面，传递来源信息
    router.push('/loan/my-applications?from=/loan/apply')
    
  } catch (error: any) {
    if (error.message !== 'cancel') {
      console.error('提交申请失败:', error)
      ElMessage.error(error.message || '提交申请失败')
    }
  } finally {
    loading.value = false
  }
}

// 返回上一页
const goBack = () => {
  router.go(-1)
}

// 格式化金额
const formatAmount = (amount: number) => {
  if (amount >= 10000) {
    return `${(amount / 10000).toFixed(1)}万`
  }
  return amount.toLocaleString()
}

// 组件挂载时加载数据
onMounted(() => {
  // 检查登录状态
  if (!userStore.isLoggedIn) {
    ElMessage.error('请先登录')
    router.push('/login')
    return
  }
  
  loadLoanTypes()
})
</script>

<template>
  <div class="loan-application-page">
    <!-- 顶部导航 -->
    <div class="top-nav">
      <div class="nav-left" @click="goBack">
        <el-icon><ArrowLeft /></el-icon>
      </div>
      <div class="nav-title">贷款申请</div>
      <div class="nav-right"></div>
    </div>

    <div class="page-content">
      <!-- 申请表单 -->
      <div class="form-container">
        <div class="form-header">
          <h3>填写申请信息</h3>
          <p>请仔细填写以下信息，确保信息真实有效</p>
        </div>

        <el-form
          ref="formRef"
          :model="formData"
          :rules="rules"
          label-width="100px"
          class="application-form"
        >
          <div class="form-section" v-if="currentTypeConfig">
            <div class="section-title">申请信息</div>

            <!-- 显示当前申请类型和名称 -->
            <div class="current-type-display">
              <div class="type-header">
                <span class="type-icon">{{ currentTypeConfig.icon }}</span>
                <div class="type-content">
                  <div class="type-name">{{ formData.type }}</div>
                  <div class="application-name">{{ formData.name }}</div>
                </div>
              </div>
              <div class="type-params">
                <div class="param-item">
                  <span class="param-label">金额范围:</span>
                  <span class="param-value">{{ formatAmount(currentTypeConfig.minAmount) }} - {{ formatAmount(currentTypeConfig.maxAmount) }}元</span>
                </div>
                <div class="param-item">
                  <span class="param-label">期限范围:</span>
                  <span class="param-value">{{ currentTypeConfig.minDuration }} - {{ currentTypeConfig.maxDuration }}个月</span>
                </div>
              </div>
            </div>
          </div>

          <div class="form-section" v-if="currentTypeConfig">
            <div class="section-title">申请金额和期限</div>

          <el-form-item label="贷款金额" prop="amount">
            <el-input-number
              v-model="formData.amount"
                :min="currentTypeConfig.minAmount"
                :max="currentTypeConfig.maxAmount"
              :step="1000"
              controls-position="right"
              style="width: 100%"
                :formatter="(value: number) => `¥ ${value.toLocaleString()}`"
                :parser="(value: string) => parseInt(value.replace(/\D/g, ''), 10)"
            />
              <div class="form-tip">
                可申请金额：{{ formatAmount(currentTypeConfig.minAmount) }} - {{ formatAmount(currentTypeConfig.maxAmount) }}元
              </div>
          </el-form-item>

            <el-form-item label="贷款期限" prop="duration">
              <el-select v-model="formData.duration" style="width: 100%">
              <el-option
                  v-for="option in durationOptions"
                :key="option.value"
                :label="option.label"
                :value="option.value"
              />
            </el-select>
              <div class="form-tip">
                可选期限：{{ currentTypeConfig.minDuration }} - {{ currentTypeConfig.maxDuration }}个月
              </div>
          </el-form-item>
          </div>

          <div class="form-section">
            <div class="section-title">用途说明</div>

            <el-form-item label="申请用途" prop="purpose">
            <el-input
                v-model="formData.purpose"
              type="textarea"
                :rows="5"
                :placeholder="currentTypeConfig?.placeholder || '请详细说明贷款用途，不少于20个字符'"
                maxlength="500"
              show-word-limit
            />
              <div class="form-tip">
                请详细描述贷款的具体用途，有助于加快审核速度
              </div>
          </el-form-item>
          </div>
        </el-form>
      </div>

      <!-- 申请预览 -->
      <div class="preview-container" v-if="formData.type">
        <div class="preview-header">
          <h3>申请预览</h3>
          <p>请确认以下信息无误</p>
            </div>
            
        <div class="preview-card">
          <div class="preview-item">
            <span class="preview-label">申请名称:</span>
            <span class="preview-value">{{ formData.name }}</span>
          </div>
          <div class="preview-item">
            <span class="preview-label">贷款类型:</span>
            <span class="preview-value">
              <span class="type-badge">
                {{ currentTypeConfig?.icon }}
                {{ formData.type }}
              </span>
            </span>
              </div>
          <div class="preview-item">
            <span class="preview-label">申请金额:</span>
            <span class="preview-value amount">¥{{ formData.amount.toLocaleString() }}</span>
            </div>
          <div class="preview-item">
            <span class="preview-label">贷款期限:</span>
            <span class="preview-value">{{ formData.duration }}个月</span>
            </div>
          <div class="preview-item">
            <span class="preview-label">申请用途:</span>
            <span class="preview-value">{{ formData.purpose || '未填写' }}</span>
          </div>
        </div>
      </div>

      <!-- 提交按钮 -->
      <div class="submit-container">
        <div class="submit-tips">
          <div class="tip-item">
            <el-icon class="tip-icon"><InfoFilled /></el-icon>
            <span>提交后将进入AI智能审核，通常1-3个工作日内完成审核</span>
          </div>
          <div class="tip-item">
            <el-icon class="tip-icon"><WarningFilled /></el-icon>
            <span>请确保信息真实有效，虚假信息将影响审核结果</span>
          </div>
        </div>

        <el-button 
          type="primary" 
          size="large" 
          @click="submitApplication"
          :loading="loading"
          :disabled="!formData.type"
          class="submit-btn"
        >
          {{ loading ? '提交中...' : '提交申请' }}
        </el-button>
      </div>
    </div>
  </div>
</template>

<style scoped>
.loan-application-page {
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

.form-container {
  background: white;
  border-radius: 12px;
  padding: 24px;
  margin-bottom: 16px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
}

.form-header {
  text-align: center;
  margin-bottom: 24px;
  padding-bottom: 16px;
  border-bottom: 1px solid #f0f0f0;
}

.form-header h3 {
  margin: 0 0 8px;
  font-size: 20px;
  font-weight: 600;
  color: #2c3e50;
}

.form-header p {
  margin: 0;
  font-size: 14px;
  color: #7f8c8d;
}

.form-section {
  margin-bottom: 24px;
}

.section-title {
  font-size: 16px;
  font-weight: 600;
  color: #2c3e50;
  margin-bottom: 16px;
  padding-bottom: 8px;
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

.form-tip {
  font-size: 12px;
  color: #7f8c8d;
  margin-top: 4px;
  line-height: 1.4;
}

/* 当前申请类型显示 */
.current-type-display {
  background: linear-gradient(135deg, #27ae60 0%, #2ecc71 100%);
  border-radius: 12px;
  padding: 20px;
  color: white;
  margin-bottom: 20px;
  box-shadow: 0 4px 12px rgba(39, 174, 96, 0.3);
}

.type-header {
  display: flex;
  align-items: center;
  gap: 16px;
  margin-bottom: 16px;
}

.type-header .type-icon {
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

.type-content {
  flex: 1;
}

.type-name {
  font-size: 20px;
  font-weight: 600;
  margin-bottom: 4px;
}

.application-name {
  font-size: 14px;
  opacity: 0.9;
}

.type-params {
  display: flex;
  gap: 24px;
  flex-wrap: wrap;
}

.param-item {
  display: flex;
  align-items: center;
  gap: 8px;
}

.param-label {
  font-size: 14px;
  opacity: 0.8;
}

.param-value {
  font-size: 14px;
  font-weight: 500;
}

.type-option {
  display: flex;
  align-items: center;
  gap: 8px;
}

.type-icon {
  font-size: 16px;
}

.type-text {
  font-size: 14px;
}

.type-info {
  margin-top: 12px;
}

.info-card {
  background: linear-gradient(135deg, #f8f9fa 0%, #e9ecef 100%);
  border-radius: 8px;
  padding: 16px;
  border-left: 4px solid #27ae60;
}

.info-header {
  display: flex;
  align-items: center;
  gap: 8px;
  margin-bottom: 12px;
}

.info-icon {
  font-size: 18px;
}

.info-title {
  font-size: 14px;
  font-weight: 600;
  color: #2c3e50;
}

.info-content {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.info-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.info-label {
  font-size: 13px;
  color: #7f8c8d;
}

.info-value {
  font-size: 13px;
  color: #2c3e50;
  font-weight: 500;
}

.preview-container {
  background: white;
  border-radius: 12px;
  padding: 20px;
  margin-bottom: 16px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
}

.preview-header {
  text-align: center;
  margin-bottom: 16px;
  padding-bottom: 12px;
  border-bottom: 1px solid #f0f0f0;
}

.preview-header h3 {
  margin: 0 0 4px;
  font-size: 16px;
  font-weight: 600;
  color: #2c3e50;
}

.preview-header p {
  margin: 0;
  font-size: 12px;
  color: #7f8c8d;
}

.preview-card {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.preview-item {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  padding: 8px 0;
  border-bottom: 1px solid #f8f9fa;
}

.preview-item:last-child {
  border-bottom: none;
}

.preview-label {
  font-size: 14px;
  color: #7f8c8d;
  width: 80px;
  flex-shrink: 0;
}

.preview-value {
  font-size: 14px;
  color: #2c3e50;
  font-weight: 500;
  text-align: right;
  flex: 1;
  word-break: break-word;
}

.preview-value.amount {
  color: #27ae60;
  font-weight: 600;
  font-size: 16px;
}

.type-badge {
  display: inline-flex;
  align-items: center;
  gap: 4px;
  padding: 4px 8px;
  background: #e8f5e8;
  color: #27ae60;
  border-radius: 12px;
  font-size: 12px;
  font-weight: 500;
}

.submit-container {
  background: white;
  border-radius: 12px;
  padding: 20px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
}

.submit-tips {
  margin-bottom: 16px;
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.tip-item {
  display: flex;
  align-items: flex-start;
  gap: 8px;
  padding: 8px;
  background: #f8f9fa;
  border-radius: 6px;
  font-size: 12px;
  color: #666;
  line-height: 1.4;
}

.tip-icon {
  flex-shrink: 0;
  margin-top: 1px;
}

.submit-btn {
  width: 100%;
  height: 50px;
  font-size: 16px;
  font-weight: 600;
  border-radius: 25px;
}

:deep(.el-form-item__label) {
  font-weight: 500;
  color: #2c3e50;
}

:deep(.el-input-number) {
  width: 100%;
}

:deep(.el-input-number .el-input__inner) {
  text-align: left;
}

:deep(.el-select-dropdown__item) {
  padding: 8px 20px;
}

@media (max-width: 768px) {
  .page-content {
    padding: 12px;
  }
  
  .form-container {
    padding: 20px;
  }
  
  .preview-container {
    padding: 16px;
  }
  
  .submit-container {
    padding: 16px;
  }
  
  .info-content {
    gap: 6px;
  }
  
  .info-item {
    flex-direction: column;
    align-items: flex-start;
    gap: 4px;
  }
  
  .preview-item {
    flex-direction: column;
    align-items: flex-start;
    gap: 4px;
  }
  
  .preview-label {
    width: auto;
  }
  
  .preview-value {
    text-align: left;
  }
}
</style> 