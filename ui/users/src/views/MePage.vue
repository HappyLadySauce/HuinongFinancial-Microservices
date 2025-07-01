<!-- vue3 element-plus -->
<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import AppFooter from './components/footer.vue'
import { useUserStore } from '../stores/user'
import { userApi } from '../services/api'

// 当前选中的导航栏
const activeTab = ref('me')
const router = useRouter()
const userStore = useUserStore()

import {
  User,
  Service,
  Help,
  Document,
  Medal,
  Setting,
  Star,
  Location,
  ChatDotRound,
  QuestionFilled,
  Ticket,
  Warning,
  InfoFilled,
  Share,
  Money,
  Wallet,
  Clock,
  Present,
  Discount,
  Bell,
  Check,
  Edit
} from '@element-plus/icons-vue'

const loading = ref(false)
const userPoints = ref(520)
const growthValue = ref(120)
const hasSignedToday = ref(false)

// 计算属性：用户信息
const userName = computed(() => {
  if (userStore.userInfo?.real_name) {
    // 隐藏姓名中间字符
    const name = userStore.userInfo.real_name
    if (name.length <= 2) return name
    return name.charAt(0) + '*'.repeat(name.length - 2) + name.charAt(name.length - 1)
  }
  return userStore.userInfo?.nickname || '未设置'
})

const userPhone = computed(() => {
  return userStore.getMaskedPhone || '未绑定'
})

const userLevel = computed(() => {
  return '普通会员' // 可以根据用户信息动态计算
})

// 检查登录状态
const checkLoginStatus = () => {
  if (!userStore.isLoggedIn) {
    ElMessage.warning('请先登录')
    router.push('/login')
    return false
  }
  return true
}

// 签到
const signIn = () => {
  if (!checkLoginStatus()) return
  
  hasSignedToday.value = true
  userPoints.value += 5
  ElMessage.success('签到成功，获得5积分！')
}

// 编辑个人信息
const editProfile = () => {
  if (!checkLoginStatus()) return
  router.push('/profile/edit')
}

// 查看我的申请
const viewMyApplications = () => {
  if (!checkLoginStatus()) return
  router.push('/loan/my-applications?from=/me')
}

// 设置页面
const goToSettings = () => {
  if (!checkLoginStatus()) return
  router.push('/settings')
}

// 登出
const logout = () => {
  userStore.logout()
  ElMessage.success('已退出登录')
  router.push('/home')
}

// 功能模块点击处理
const handleModuleClick = (module: string) => {
  if (!checkLoginStatus()) return
  
  switch (module) {
    case 'orders':
      ElMessage.info('订单功能开发中')
      break
    case 'favorites':
      ElMessage.info('收藏功能开发中')
      break
    case 'address':
      ElMessage.info('地址管理功能开发中')
      break
    case 'vip':
      ElMessage.info('会员特权功能开发中')
      break
    case 'promotions':
      ElMessage.info('优惠活动功能开发中')
      break
    case 'notifications':
      ElMessage.info('消息通知功能开发中')
      break
    case 'service':
      ElMessage.info('在线客服功能开发中')
      break
    case 'faq':
      ElMessage.info('常见问题功能开发中')
      break
    case 'feedback':
      ElMessage.info('投诉反馈功能开发中')
      break
    case 'risk':
      ElMessage.info('风险提示功能开发中')
      break
    case 'about':
      ElMessage.info('关于我们功能开发中')
      break
    case 'share':
      ElMessage.info('分享推荐功能开发中')
      break
    default:
      ElMessage.info('功能开发中')
  }
}

// 加载用户信息 - 适配新API
const loadUserInfo = async () => {
  if (!userStore.isLoggedIn || !userStore.userInfo?.phone) return
  
  try {
    loading.value = true
    // 新API需要传递phone参数
    const response = await userApi.getUserInfo(userStore.userInfo.phone)
    // 新API直接返回用户信息，不需要.data
    userStore.setUserInfo(response)
  } catch (error: any) {
    console.error('加载用户信息失败:', error)
    // 不显示错误消息，避免影响用户体验
  } finally {
    loading.value = false
  }
}

// 组件挂载时加载数据
onMounted(() => {
  if (userStore.isLoggedIn) {
    loadUserInfo()
  }
})
</script>

<template>
  <div class="me-page">
    <!-- Header with User Info -->
    <div class="user-header">
      <div class="user-details">
        <el-avatar :size="55" class="user-avatar">{{ userName.charAt(0) }}</el-avatar>
        <div class="user-info">
          <div class="user-name">{{ userName }}</div>
          <div class="user-phone">{{ userPhone }}</div>
          <div class="user-level">
            <el-tag size="small" type="warning">{{ userLevel }}</el-tag>
          </div>
        </div>
        <el-button class="settings-btn" text>
          <el-icon><Setting /></el-icon>
        </el-button>
      </div>
      
      <!-- 用户积分和签到 -->
      <div class="user-stats">
        <div class="stats-item">
          <div class="stats-value">{{ userPoints }}</div>
          <div class="stats-label">积分</div>
        </div>
        <div class="stats-item">
          <div class="stats-value">{{ growthValue }}</div>
          <div class="stats-label">成长值</div>
        </div>
        <div class="sign-in-container">
          <el-button 
            type="success" 
            :disabled="hasSignedToday" 
            size="small" 
            @click="signIn"
            class="sign-in-btn"
          >
            <el-icon v-if="hasSignedToday"><Check /></el-icon>
            {{ hasSignedToday ? '已签到' : '签到' }}
          </el-button>
          <div class="sign-in-tip">每日签到得5积分</div>
        </div>
      </div>
    </div>

    <!-- Module Sections -->
    <div class="modules-container">
      <!-- User Module (Combined with Assets) -->
      <div class="module-section">
        <div class="module-title">
          <el-icon><User /></el-icon>
          <span>个人中心</span>
        </div>
        
        <!-- User Assets -->
        <div class="asset-row">
          <el-row :gutter="20" class="mb-15">
            <el-col :span="8">
              <div class="asset-item">
                <div class="asset-value">0.00</div>
                <div class="asset-label">账户余额(元)</div>
              </div>
            </el-col>
            <el-col :span="8">
              <div class="asset-item">
                <div class="asset-value">0</div>
                <div class="asset-label">优惠券(张)</div>
              </div>
            </el-col>
            <el-col :span="8">
              <div class="asset-item">
                <div class="asset-value">0</div>
                <div class="asset-label">积分</div>
              </div>
            </el-col>
          </el-row>
        </div>
        
        <el-divider class="my-divider" />
        
        <!-- User Functions -->
        <el-row :gutter="0" class="module-grid">
          <el-col :span="8">
            <div class="grid-item">
              <div class="grid-icon"><el-icon><Document /></el-icon></div>
              <div class="grid-label">我的订单</div>
            </div>
          </el-col>
          <el-col :span="8">
            <div class="grid-item">
              <div class="grid-icon"><el-icon><Star /></el-icon></div>
              <div class="grid-label">我的收藏</div>
            </div>
          </el-col>
          <el-col :span="8">
            <div class="grid-item">
              <div class="grid-icon"><el-icon><Location /></el-icon></div>
              <div class="grid-label">地址管理</div>
            </div>
          </el-col>
          <el-col :span="8">
            <div class="grid-item">
              <div class="grid-icon"><el-icon><Present /></el-icon></div>
              <div class="grid-label">会员特权</div>
            </div>
          </el-col>
          <el-col :span="8">
            <div class="grid-item">
              <div class="grid-icon"><el-icon><Discount /></el-icon></div>
              <div class="grid-label">优惠活动</div>
            </div>
          </el-col>
          <el-col :span="8">
            <div class="grid-item">
              <div class="grid-icon"><el-icon><Bell /></el-icon></div>
              <div class="grid-label">消息通知</div>
            </div>
          </el-col>
        </el-row>
      </div>

      <!-- Service & Protection Module -->
      <div class="module-section">
        <div class="module-title">
          <el-icon><Service /></el-icon>
          <span>服务与保障</span>
        </div>
        <el-row :gutter="0" class="module-grid">
          <el-col :span="8">
            <div class="grid-item">
              <div class="grid-icon"><el-icon><ChatDotRound /></el-icon></div>
              <div class="grid-label">在线客服</div>
            </div>
          </el-col>
          <el-col :span="8">
            <div class="grid-item">
              <div class="grid-icon"><el-icon><QuestionFilled /></el-icon></div>
              <div class="grid-label">常见问题</div>
            </div>
          </el-col>
          <el-col :span="8">
            <div class="grid-item">
              <div class="grid-icon"><el-icon><Ticket /></el-icon></div>
              <div class="grid-label">投诉反馈</div>
            </div>
          </el-col>
          <el-col :span="8">
            <div class="grid-item">
              <div class="grid-icon"><el-icon><Warning /></el-icon></div>
              <div class="grid-label">风险提示</div>
            </div>
          </el-col>
          <el-col :span="8">
            <div class="grid-item">
              <div class="grid-icon"><el-icon><InfoFilled /></el-icon></div>
              <div class="grid-label">消费者权益</div>
            </div>
          </el-col>
          <el-col :span="8">
            <div class="grid-item">
              <div class="grid-icon"><el-icon><Share /></el-icon></div>
              <div class="grid-label">举报中心</div>
            </div>
          </el-col>
        </el-row>
      </div>

      <!-- Business Module -->
      <div class="module-section">
        <div class="module-title">
          <el-icon><Money /></el-icon>
          <span>业务办理</span>
        </div>
        <el-row :gutter="0" class="module-grid">
          <el-col :span="8">
            <div class="grid-item">
              <div class="grid-icon"><el-icon><Money /></el-icon></div>
              <div class="grid-label">贷款申请</div>
            </div>
          </el-col>
          <el-col :span="8">
            <div class="grid-item">
              <div class="grid-icon"><el-icon><Wallet /></el-icon></div>
              <div class="grid-label">账户管理</div>
            </div>
          </el-col>
          <el-col :span="8">
            <div class="grid-item">
              <div class="grid-icon"><el-icon><Clock /></el-icon></div>
              <div class="grid-label">还款计划</div>
            </div>
          </el-col>
        </el-row>
      </div>
    </div>
  </div>
  <app-footer v-model:active-tab="activeTab" />
</template>

<style scoped>
.me-page {
  background-color: #f5f5f5;
  min-height: 100vh;
  padding-bottom: 20px;
}

.user-header {
  background: linear-gradient(135deg, #2e8b57, #4CAF50);
  color: white;
  padding: 8px;
  margin: 8px;
  position: relative;
  display: flex;
  flex-direction: column;
  align-items: stretch;
  border-radius: 12px;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1);
}

.user-details {
  width: 100%;
  padding: 0;
  display: flex;
  align-items: center;
  position: relative;
}

.user-avatar {
  background-color: #fff;
  color: #2e8b57;
  font-weight: bold;
  flex-shrink: 0;
}

.user-info {
  margin-left: 12px;
  flex: 1;
}

.user-name {
  font-size: 16px;
  font-weight: bold;
  margin-bottom: 3px;
}

.user-phone {
  font-size: 12px;
  opacity: 0.8;
  margin-bottom: 3px;
}

.settings-btn {
  color: white;
  font-size: 18px;
  position: absolute;
  top: 0;
  right: 0;
}

.modules-container {
  padding: 0 10px;
}

.module-section {
  background-color: white;
  border-radius: 8px;
  margin-bottom: 12px;
  padding: 15px;
  box-shadow: 0 2px 4px rgba(0,0,0,0.08);
}

.module-title {
  display: flex;
  align-items: center;
  margin-bottom: 15px;
  font-size: 16px;
  font-weight: 500;
  color: #333;
}

.module-title .el-icon {
  margin-right: 8px;
  font-size: 18px;
  color: #2e8b57;
}

.module-grid {
  text-align: center;
}

.asset-item {
  text-align: center;
}

.asset-value {
  font-size: 18px;
  font-weight: bold;
  color: #333;
}

.asset-label {
  font-size: 12px;
  color: #999;
  margin-top: 5px;
}

.grid-item {
  padding: 10px 0;
  cursor: pointer;
}

.grid-icon {
  display: flex;
  justify-content: center;
  margin-bottom: 8px;
}

.grid-icon .el-icon {
  font-size: 24px;
  color: #2e8b57;
  background-color: rgba(46, 139, 87, 0.08);
  padding: 10px;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
}

.grid-label {
  font-size: 12px;
  color: #666;
}

.mb-15 {
  margin-bottom: 15px;
}

.my-divider {
  margin: 15px 0;
}

.user-stats {
  display: flex;
  width: 100%;
  justify-content: space-around;
  margin-top: 8px;
  padding-top: 8px;
  border-top: 1px solid rgba(255, 255, 255, 0.15);
}

.stats-item {
  text-align: center;
}

.stats-value {
  font-size: 16px;
  font-weight: bold;
}

.stats-label {
  font-size: 11px;
  opacity: 0.8;
  margin-top: 1px;
}

.sign-in-container {
  display: flex;
  flex-direction: column;
  align-items: center;
}

.sign-in-btn {
  background-color: rgba(255, 255, 255, 0.2);
  border: none;
  padding: 3px 12px;
  font-size: 12px;
}

.sign-in-tip {
  font-size: 10px;
  margin-top: 5px;
  opacity: 0.7;
}
</style>
