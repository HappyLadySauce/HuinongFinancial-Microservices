<script setup lang="ts">
import { reactive, ref } from 'vue'
import type { FormInstance, FormRules } from 'element-plus'
import { ElMessage } from 'element-plus'
import { useRouter } from 'vue-router'
import { authApi } from '../../services/api'

const FormRef = ref<FormInstance>()
const router = useRouter()
const loading = ref(false)

// 表单数据
const FormData = reactive({
  phone: '',
  password: '',
  confirmPassword: ''
})

// 表单验证规则
const validatePhone = (rule: any, value: any, callback: any) => {
  if (value === '') {
    callback(new Error('请输入手机号码'))
  } else if (!/^1[3-9]\d{9}$/.test(value)) {
    callback(new Error('请输入正确的手机号码'))
  } else {
    callback()
  }
}

const validatePassword = (rule: any, value: any, callback: any) => {
  if (value === '') {
    callback(new Error('请输入密码'))
  } else if (value.length < 6) {
    callback(new Error('密码不能少于6位'))
  } else {
    callback()
  }
}

const validateConfirmPassword = (rule: any, value: any, callback: any) => {
  if (value === '') {
    callback(new Error('请确认密码'))
  } else if (value !== FormData.password) {
    callback(new Error('两次输入的密码不一致'))
  } else {
    callback()
  }
}

const rules = reactive<FormRules<typeof FormData>>({
  phone: [{ validator: validatePhone, trigger: 'blur' }],
  password: [{ validator: validatePassword, trigger: 'blur' }],
  confirmPassword: [{ validator: validateConfirmPassword, trigger: 'blur' }]
})

// 注册
const submitForm = async (formEl: FormInstance | undefined) => {
  if (!formEl) return
  
  try {
    const valid = await formEl.validate()
    if (!valid) return
    
    loading.value = true
    
    // 使用新的JWT注册API
    await authApi.register(FormData.phone, FormData.password)
    
    ElMessage.success('注册成功，请登录')
    router.push('/login')
    
  } catch (error: any) {
    console.error('注册失败:', error)
    ElMessage.error(error.message || '注册失败')
  } finally {
    loading.value = false
  }
}

// 返回登录
const goToLogin = () => {
  router.push('/login')
}
</script>

<template>
  <div class="register-container">
    <div class="register-box">
      <div class="logo">
        <img src="../../assets/images/logo.png" alt="数字惠农" />
      </div>
      
      <h1 class="system-title">创建账户</h1>
      <p class="system-subtitle">加入数字惠农，享受金融服务</p>

      <el-form 
        ref="FormRef" 
        :model="FormData" 
        :rules="rules" 
        class="register-form"
      >
        <div class="form-item">
          <div class="form-label">手机号码</div>
          <el-input 
            v-model="FormData.phone" 
            placeholder="请输入手机号码" 
            size="large"
            maxlength="11"
          />
        </div>

        <div class="form-item">
          <div class="form-label">密码</div>
          <el-input
            v-model="FormData.password"
            type="password"
            placeholder="请输入密码"
            size="large"
            show-password
          />
        </div>

        <div class="form-item">
          <div class="form-label">确认密码</div>
          <el-input
            v-model="FormData.confirmPassword"
            type="password"
            placeholder="请再次输入密码"
            size="large"
            show-password
          />
        </div>

        <div class="submit-btn">
          <el-button 
            type="primary" 
            @click="submitForm(FormRef)" 
            class="register-button"
            size="large"
            :loading="loading"
            :disabled="loading"
          >
            {{ loading ? '注册中...' : '立即注册' }}
          </el-button>
        </div>
        
        <div class="login-link">
          <span>已有账户？</span>
          <a @click="goToLogin" class="link">立即登录</a>
        </div>
      </el-form>
    </div>
  </div>
</template>

<style scoped>
.register-container {
  width: 100vw;
  height: 100vh;
  display: flex;
  justify-content: center;
  align-items: center;
  background: #ffffff;
  position: relative;
  overflow: hidden;
}

.register-box {
  width: 85%;
  max-width: 450px;
  background-color: rgba(255, 255, 255, 0.92);
  border-radius: 10px;
  padding: 40px;
  text-align: center;
  box-shadow: 0 10px 25px rgba(0, 0, 0, 0.1), 0 2px 10px rgba(0, 0, 0, 0.05);
  border: 1px solid rgba(255, 255, 255, 0.3);
  backdrop-filter: blur(5px);
  animation: float 5s infinite ease-in-out;
  position: relative;
  z-index: 1;
}

@keyframes float {
  0%, 100% {
    transform: translateY(0);
  }
  50% {
    transform: translateY(-10px);
  }
}

.logo {
  margin-bottom: 20px;
  transform: scale(1.5);
}

.logo img {
  width: 100px;
  height: 100px;
}

.system-title {
  font-size: 28px;
  font-weight: bold;
  color: #2c3e50;
  margin-bottom: 10px;
  text-shadow: 2px 2px 4px rgba(0, 0, 0, 0.1);
}

.system-subtitle {
  font-size: 16px;
  color: #7f8c8d;
  margin-bottom: 30px;
}

.register-form {
  text-align: left;
}

.form-item {
  margin-bottom: 20px;
}

.form-label {
  font-size: 14px;
  color: #2c3e50;
  margin-bottom: 8px;
  font-weight: 500;
}

.submit-btn {
  margin: 30px 0 20px 0;
}

.register-button {
  width: 100%;
  height: 50px;
  font-size: 16px;
  font-weight: bold;
  border-radius: 25px;
  background: linear-gradient(135deg, #27ae60, #2ecc71);
  border: none;
  box-shadow: 0 4px 15px rgba(46, 204, 113, 0.3);
  transition: all 0.3s ease;
}

.register-button:hover {
  transform: translateY(-2px);
  box-shadow: 0 6px 20px rgba(46, 204, 113, 0.4);
}

.login-link {
  text-align: center;
  color: #7f8c8d;
  font-size: 14px;
}

.link {
  color: #27ae60;
  cursor: pointer;
  text-decoration: none;
  font-weight: bold;
  margin-left: 5px;
}

.link:hover {
  text-decoration: underline;
}

:deep(.el-input__wrapper) {
  border-radius: 8px;
  box-shadow: 0 0 0 1px #dcdfe6 inset;
}

:deep(.el-input__wrapper:hover) {
  box-shadow: 0 0 0 1px #c0c4cc inset;
}

:deep(.el-input__wrapper.is-focus) {
  box-shadow: 0 0 0 1px #27ae60 inset;
}
</style> 