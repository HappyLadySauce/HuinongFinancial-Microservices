<template>
  <div class="login-container">
    <div class="login-content">
      <div class="login-header">
        <h1>数字惠农OA管理系统</h1>
        <p>Digital Agriculture OA Management System</p>
      </div>
      
      <el-card class="login-card" shadow="always">
        <template #header>
          <div class="card-header">
            <span>管理员登录</span>
          </div>
        </template>
        
        <el-form
          ref="loginFormRef"
          :model="loginForm"
          :rules="loginRules"
          size="large"
          @submit.prevent="handleLogin"
        >
          <el-form-item prop="phone">
            <el-input
              v-model="loginForm.phone"
              placeholder="请输入手机号"
              :prefix-icon="User"
              clearable
            />
          </el-form-item>
          
          <el-form-item prop="password">
            <el-input
              v-model="loginForm.password"
              type="password"
              placeholder="请输入密码"
              :prefix-icon="Lock"
              show-password
              clearable
              @keyup.enter="handleLogin"
            />
          </el-form-item>
          
          <el-form-item>
            <el-checkbox v-model="rememberMe">记住我</el-checkbox>
          </el-form-item>
          
          <el-form-item>
            <el-button
              type="primary"
              size="large"
              style="width: 100%"
              :loading="loading"
              @click="handleLogin"
            >
              登录
            </el-button>
          </el-form-item>
        </el-form>
        
        <div class="login-tips">
          <h4>测试账号</h4>
          <p><strong>管理员：</strong>13452552349 / admin123</p>
          <p><strong>备用账号：</strong>13800138000 / 123456</p>
        </div>
      </el-card>
    </div>
  </div>
</template>

<script setup lang="ts">
import { reactive, ref } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage, type FormInstance, type FormRules } from 'element-plus'
import { User, Lock } from '@element-plus/icons-vue'
import { useUserStore } from '@/stores/user'
import { userApi } from '@/services/api'

const router = useRouter()
const userStore = useUserStore()

const loginFormRef = ref<FormInstance>()
const loading = ref(false)
const rememberMe = ref(false)

const loginForm = reactive({
  phone: '',
  password: ''
})

const loginRules: FormRules = {
  phone: [
    { required: true, message: '请输入手机号', trigger: 'blur' },
    { pattern: /^1[3-9]\d{9}$/, message: '请输入正确的手机号格式', trigger: 'blur' }
  ],
  password: [
    { required: true, message: '请输入密码', trigger: 'blur' },
    { min: 6, max: 20, message: '密码长度在 6 到 20 个字符', trigger: 'blur' }
  ]
}

const handleLogin = async () => {
  if (!loginFormRef.value) return
  
  try {
    await loginFormRef.value.validate()
    loading.value = true
    
    // 调用登录API
    const loginResponse = await userApi.login({
      phone: loginForm.phone,
      password: loginForm.password
    })
    
    // 设置用户登录状态
    userStore.login(loginResponse.data)
    
    // 获取用户信息
    try {
      const userInfoResponse = await userApi.getUserInfo()
      userStore.setUserInfo(userInfoResponse.data)
    } catch (error) {
      console.warn('获取用户信息失败:', error)
    }
    
    ElMessage.success('登录成功')
    router.push('/dashboard')
  } catch (error: any) {
    ElMessage.error(error.message || '登录失败，请检查用户名和密码')
  } finally {
    loading.value = false
  }
}
</script>

<style scoped>
.login-container {
  min-height: 100vh;
  background-image: url('@/assets/背景.jpg');
  background-size: cover;
  background-position: center;
  background-repeat: no-repeat;
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 20px;
  position: relative;
}

.login-container::before {
  content: '';
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background-color: rgba(0, 0, 0, 0.1);
  z-index: 1;
}

.login-content {
  width: 100%;
  max-width: 400px;
  position: relative;
  z-index: 2;
}

.login-header {
  text-align: center;
  margin-bottom: 40px;
  color: white;
}

.login-header h1 {
  font-size: 32px;
  font-weight: 600;
  margin-bottom: 10px;
  text-shadow: 0 2px 8px rgba(0, 0, 0, 0.5);
}

.login-header p {
  font-size: 16px;
  opacity: 0.95;
  margin: 0;
  text-shadow: 0 1px 4px rgba(0, 0, 0, 0.3);
}

.login-card {
  border-radius: 16px;
  border: none;
  box-shadow: 0 15px 50px rgba(0, 0, 0, 0.3);
  backdrop-filter: blur(8px);
  background-color: rgba(255, 255, 255, 0.85);
  overflow: hidden;
  transition: all 0.3s ease;
  width: 100%;
}

.login-card:hover {
  transform: translateY(-5px);
  box-shadow: 0 20px 60px rgba(0, 0, 0, 0.4);
}

.card-header {
  text-align: center;
  font-size: 22px;
  font-weight: 600;
  color: #333;
  padding: 15px 0;
  border-bottom: 1px solid rgba(0, 0, 0, 0.05);
}

.login-tips {
  margin-top: 15px;
  padding: 12px;
  background: #f1f8e9;
  border-radius: 8px;
  border-left: 4px solid #4CAF50;
}

.login-tips h4 {
  margin: 0 0 8px 0;
  color: #2E7D32;
  font-size: 14px;
  font-weight: 600;
}

.login-tips p {
  margin: 4px 0;
  color: #555;
  font-size: 13px;
  line-height: 1.4;
}

.login-tips strong {
  color: #2E7D32;
  font-weight: 600;
}

:deep(.el-form-item) {
  margin-bottom: 18px;
}

:deep(.el-input__inner) {
  border-radius: 8px;
  height: 42px;
  font-size: 14px;
  border: none;
  transition: all 0.3s ease;
  background-color: #f9f9f9;
}

:deep(.el-input__wrapper) {
  box-shadow: none;
  border: none;
  background-color: #f9f9f9;
  border-radius: 8px;
}

:deep(.el-input__wrapper:hover) {
  box-shadow: 0 0 0 1px rgba(76, 175, 80, 0.3);
}

:deep(.el-input__wrapper.is-focus) {
  box-shadow: 0 0 0 1px rgba(76, 175, 80, 0.5);
}

:deep(.el-input__wrapper:hover .el-input__inner) {
  border-color: transparent;
}

:deep(.el-input__wrapper.is-focus .el-input__inner) {
  border-color: transparent;
}

:deep(.el-button) {
  border-radius: 8px;
  font-weight: 500;
  background-color: #4CAF50;
  border-color: #4CAF50;
  height: 42px;
  font-size: 15px;
  letter-spacing: 2px;
  transition: all 0.3s ease;
}

:deep(.el-button:hover) {
  background-color: #388E3C;
  border-color: #388E3C;
  transform: translateY(-2px);
  box-shadow: 0 8px 20px rgba(76, 175, 80, 0.4);
}

:deep(.el-checkbox__input.is-checked .el-checkbox__inner) {
  background-color: #4CAF50;
  border-color: #4CAF50;
}

:deep(.el-input__prefix-inner) {
  display: flex;
  align-items: center;
  font-size: 16px;
  color: #4CAF50;
}

:deep(.el-form) {
  padding: 5px 25px 20px;
}

:deep(.el-checkbox) {
  height: 30px;
}
</style> 