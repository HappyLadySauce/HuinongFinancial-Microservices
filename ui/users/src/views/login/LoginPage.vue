<script setup lang="ts">
// 导入组合式API
import { reactive, ref } from 'vue'
// 导入 element-plus 表单实例和表单规则类型
import type { FormInstance, FormRules } from 'element-plus'
import { ElMessage } from 'element-plus'
import { useUserStore } from '../../stores/user'
import { useRouter } from 'vue-router'
import { authApi, userApi } from '../../services/api'

// 定义ref对象表单实例
const FormRef = ref<FormInstance>()
const router = useRouter()
const userStore = useUserStore()
const loading = ref(false)

// 定义表单数据
const FormData = reactive({
  phone: '',
  password: '',
  remember: false
})

// 定义表单验证
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

// 定义表单规则
const rules = reactive<FormRules<typeof FormData>>({
    phone: [{ validator: validatePhone, trigger: 'blur' }],
    password: [{ validator: validatePassword, trigger: 'blur' }]
})

// 定义表单提交
const submitForm = async (formEl: FormInstance | undefined) => {
  if (!formEl) return
  
  try {
    const valid = await formEl.validate()
    if (!valid) return
    
    loading.value = true
    
    // 调用新的JWT登录API
    const response = await authApi.login(FormData.phone, FormData.password)
    
    // 保存JWT token
    userStore.login(response)
    
    // 获取用户详细信息
    try {
      const userInfoResponse = await userApi.getUserInfo(FormData.phone)
      userStore.setUserInfo(userInfoResponse)
    } catch (error) {
      console.warn('获取用户信息失败:', error)
      // 即使获取用户信息失败，也保持登录状态
    }
    
    ElMessage.success('登录成功')
    
    // 获取重定向路径，优先使用用户原来要访问的页面
    const redirectPath = (router.currentRoute.value.query.redirect as string) || '/home'
    router.push(redirectPath)
    
  } catch (error: any) {
    console.error('登录失败:', error)
    ElMessage.error(error.message || '登录失败，请检查手机号和密码')
  } finally {
    loading.value = false
  }
}

// 忘记密码
const forgotPassword = () => {
  ElMessage.info('忘记密码功能暂未开放，请联系客服')
}

// 立即注册
const register = () => {
  router.push('/register')
}
</script>

<template>
  <div class="login-container">
    <div class="login-box">
      <div class="logo">
        <img src="../../assets/images/logo.png" alt="数字惠农" />
      </div>
      
      <h1 class="system-title">数字惠农</h1>
      <p class="system-subtitle">金融助力，共富乡村</p>

      <!-- 测试账号提示 -->
      <div class="test-account-tip">
        <el-alert 
          title="测试账号" 
          type="info" 
          :closable="false"
          show-icon
        >
          <template #default>
            <p>手机号：13452552490</p>
            <p>密码：13452552490</p>
          </template>
        </el-alert>
      </div>

      <el-form 
        ref="FormRef" 
        :model="FormData" 
        :rules="rules" 
        class="form-container"
      >
        <div class="form-item">
          <div class="form-label">手机号码</div>
          <div class="input-wrapper">
            <el-form-item prop="phone">
              <el-input 
                v-model="FormData.phone" 
                placeholder="请输入手机号码" 
                prefix-icon="el-icon-phone"
              />
            </el-form-item>
          </div>
        </div>

        <div class="form-item">
          <div class="form-label">密码</div>
          <div class="input-wrapper">
            <el-form-item prop="password">
              <el-input
                v-model="FormData.password"
                type="password"
                placeholder="请输入密码"
                prefix-icon="el-icon-lock"
                @keyup.enter="submitForm(FormRef)"
              />
            </el-form-item>
          </div>
        </div>

        <div class="remember-forgot">
          <div class="remember">
            <el-checkbox v-model="FormData.remember">记住我</el-checkbox>
          </div>
          <div class="forgot">
            <a @click="forgotPassword" class="forgot-link">忘记密码？</a>
          </div>
        </div>

        <div class="submit-btn">
          <el-button 
            type="primary" 
            @click="submitForm(FormRef)" 
            class="login-button"
            :loading="loading"
            :disabled="loading"
          >
            {{ loading ? '登录中...' : '登录' }}
          </el-button>
        </div>
        
        <div class="register-btn">
          <el-button @click="register" class="register-button">
            立即注册
          </el-button>
        </div>
      </el-form>
    </div>
  </div>
</template>

<style scoped>


/* 登录页面 */
.login-container {
  width: 100vw;
  height: 100vh;
  /* 居中 */
  display: flex;
  /* 水平居中 */
  justify-content: center;
  /* 垂直居中 */
  align-items: center;
  /* 白色背景 */
  background: #ffffff;
  position: relative;
  overflow: hidden;
}

/* 删除不需要的背景动画 */
@keyframes bgShimmer {
  0%, 50%, 100% {
    background-color: #ffffff;
  }
}

.login-box {
  /* 宽度 */
  width: 85%;
  max-width: 450px;
  /* 背景颜色 - 半透明 */
  background-color: rgba(255, 255, 255, 0.92);
  /* 圆角 */
  border-radius: 10px;
  /* 内边距 */
  padding: 40px;
  /* 文本居中 */
  text-align: center;
  /* 高级阴影效果 */
  box-shadow: 0 10px 25px rgba(0, 0, 0, 0.1), 0 2px 10px rgba(0, 0, 0, 0.05);
  /* 边框效果 */
  border: 1px solid rgba(255, 255, 255, 0.3);
  /* 玻璃态模糊效果 */
  backdrop-filter: blur(5px);
  /* 动画效果 */
  animation: float 5s infinite ease-in-out;
  /* 相对定位 */
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
  /* 下边距 */
  margin-bottom: 20px;
  /* 放大 */
  transform: scale(1.5);
}

.logo img {
  /* 宽度 */
  width: 100px;
  /* 高度 */
  height: 100px;
  /* 移除背景色以显示原始图片 */
  background-color: transparent;
  /* 保留原始图片形状 */
  border-radius: 0;
  /* 移除内边距 */
  padding: 0;
}

.system-title {
  /* 金色渐变 */
  background-image: linear-gradient(to right, #D4AF37, #FBD341, #D4AF37);
  -webkit-background-clip: text;
  background-clip: text;
  -webkit-text-fill-color: transparent;
  /* 字体大小 */
  font-size: clamp(24px, 5vw, 28px);
  /* 下边距 */
  margin-bottom: 10px;
  /* 加粗显示更明显 */
  font-weight: bold;
}

.system-subtitle {
  color: #666;
  font-size: clamp(16px, 4vw, 20px);
  margin-bottom: 30px;
}

.test-account-tip {
  margin-bottom: 20px;
}

.test-account-tip :deep(.el-alert) {
  background: #f0f9ff;
  border: 1px solid #bae6fd;
}

.test-account-tip :deep(.el-alert__content) {
  text-align: left;
}

.test-account-tip p {
  margin: 2px 0;
  font-size: 14px;
  color: #0369a1;
}

.form-container {
  margin: 30px 0;
}

.form-item {
  margin-bottom: 25px;
  text-align: left;
}

.form-label {
  color: #333;
  font-size: clamp(16px, 4vw, 20px);
  margin-bottom: 5px;
}

.input-wrapper {
  width: 100%;
}

/* 表单验证样式调整 */
:deep(.el-form-item) {
  margin-bottom: 0;
}

:deep(.el-form-item__content) {
  line-height: normal;
}

:deep(.el-form-item__error) {
  padding-top: 2px;
  font-size: 12px;
}

.remember-forgot {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 25px;
}

.remember :deep(.el-checkbox__label) {
  color: #333;
  font-size: clamp(16px, 4vw, 20px);
}

.forgot-link {
  background-image: linear-gradient(to right, #e6a23c, #f5c069, #e6a23c);
  -webkit-background-clip: text;
  background-clip: text;
  -webkit-text-fill-color: transparent;
  cursor: pointer;
  text-decoration: none;
  font-weight: bold;
  font-size: clamp(16px, 4vw, 20px);
}

.login-button {
  width: 100%;
  height: 45px;
  background-image: linear-gradient(to right, #e6a23c, #f5c069, #e6a23c);
  color: white;
  border: none;
  border-radius: 4px;
  font-size: 18px;
  font-weight: bold;
  transition: all 0.3s ease;
}

.login-button:hover {
  background-image: linear-gradient(to right, #d39531, #e6a23c, #d39531);
  border: none;
  box-shadow: 0 4px 8px rgba(214, 146, 32, 0.3);
}

.register-button {
  width: 100%;
  height: 45px;
  background: linear-gradient(to right, #ffffff, #fdf6ec, #ffffff);
  color: #e6a23c;
  border: 1px solid #e6a23c;
  border-radius: 4px;
  font-size: 18px;
  font-weight: bold;
  margin-top: 15px;
  transition: all 0.3s ease;
}

.register-button:hover {
  background: linear-gradient(to right, #fdf6ec, #ffeed8, #fdf6ec);
  border: 1px solid #d39531;  
}

:deep(.el-input__wrapper) {
  /* 增大 */
  height: 60px;
}

/* 媒体查询 - 响应式设计 */
@media screen and (max-width: 768px) {
  .login-box {
    padding: 30px 20px;
  }
  
  .logo {
    transform: scale(1.2);
  }
  
  .logo img {
    width: 80px;
    height: 80px;
  }
  
  .form-container {
    margin: 20px 0;
  }
  
  .form-item {
    margin-bottom: 20px;
  }
  
  :deep(.el-input__wrapper) {
    height: 50px;
  }
  
  .login-button, .register-button {
    height: 40px;
  }
}

@media screen and (max-width: 480px) {
  .login-box {
    padding: 25px 15px;
  }
  
  .logo {
    transform: scale(1);
    margin-bottom: 15px;
  }
  
  .logo img {
    width: 70px;
    height: 70px;
  }
  
  .form-container {
    margin: 15px 0;
  }
  
  .form-item {
    margin-bottom: 15px;
  }
  
  .remember-forgot {
    margin-bottom: 20px;
  }
  
  :deep(.el-input__wrapper) {
    height: 45px;
  }
}
</style>