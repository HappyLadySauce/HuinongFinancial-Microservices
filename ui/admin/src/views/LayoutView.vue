<template>
  <div class="oa-layout">
    <el-container class="layout-container">
      <!-- 顶部导航栏 -->
      <el-header class="layout-header" height="60px">
        <div class="header-left">
          <div class="logo">
            <div class="logo-icon">
              <el-icon><HomeFilled /></el-icon>
            </div>
            <span class="logo-text" v-if="!isCollapse">惠农金融</span>
          </div>
          <el-button 
            link 
            @click="toggleCollapse" 
            class="collapse-btn"
            :icon="isCollapse ? Expand : Fold"
          />
        </div>

        <div class="header-center">
          <el-input
            v-model="searchText"
            placeholder="请输入需要查询的内容"
            class="search-input"
            clearable
          >
            <template #prefix>
              <el-icon><Search /></el-icon>
            </template>
          </el-input>
        </div>

        <div class="header-right">
          <!-- 审批状态 -->
          <div class="approval-status">
            <el-tooltip content="待审批数量" placement="bottom">
              <div class="status-item">
                <el-icon class="status-icon"><DocumentChecked /></el-icon>
                <el-badge :value="pendingApprovals" :max="99" class="status-badge">
                  <span class="status-text">待审批</span>
                </el-badge>
              </div>
            </el-tooltip>
          </div>

          <!-- 通知 -->
          <div class="notification">
            <el-popover placement="bottom" width="320" trigger="click">
              <template #reference>
                <el-badge :value="notifications.length" :max="99" :hidden="notifications.length === 0">
                  <el-button link class="notification-btn">
                    <el-icon><Bell /></el-icon>
                  </el-button>
                </el-badge>
              </template>
              <div class="notification-panel">
                <div class="notification-header">
                  <span>消息通知</span>
                  <el-button link size="small" @click="clearAllNotifications">清空</el-button>
                </div>
                <div class="notification-list">
                  <div 
                    v-for="item in notifications" 
                    :key="item.id" 
                    class="notification-item"
                    @click="handleNotificationClick(item)"
                  >
                    <div class="notification-content">
                      <div class="notification-title">{{ item.title }}</div>
                      <div class="notification-desc">{{ item.content }}</div>
                    </div>
                    <div class="notification-time">{{ formatTime(item.createdAt) }}</div>
                  </div>
                  <div v-if="notifications.length === 0" class="notification-empty">
                    暂无新消息
                  </div>
                </div>
              </div>
            </el-popover>
          </div>

          <!-- 用户信息 -->
          <el-dropdown @command="handleUserCommand" trigger="hover" class="user-dropdown">
            <div class="user-info">
              <el-avatar :size="36" class="user-avatar">
                <el-icon><UserFilled /></el-icon>
              </el-avatar>
              <div class="user-text" v-if="!isCollapse">
                <div class="user-name">{{ currentUser?.username || '管理员' }}</div>
                <div class="user-role">{{ getRoleName(currentUser?.role) }}</div>
              </div>
              <el-icon class="user-arrow"><ArrowDown /></el-icon>
            </div>
            <template #dropdown>
              <el-dropdown-menu>
                <el-dropdown-item command="profile">
                  <el-icon><User /></el-icon>
                  个人中心
                </el-dropdown-item>
                <el-dropdown-item command="settings">
                  <el-icon><Setting /></el-icon>
                  系统设置
                </el-dropdown-item>
                <el-dropdown-item divided command="logout">
                  <el-icon><SwitchButton /></el-icon>
                  退出登录
                </el-dropdown-item>
              </el-dropdown-menu>
            </template>
          </el-dropdown>
        </div>
      </el-header>

      <el-container class="main-container">
        <!-- 侧边栏 -->
        <el-aside :width="isCollapse ? '64px' : '240px'" class="layout-sidebar">
          <div class="sidebar-content">
            <el-scrollbar class="sidebar-scrollbar">
              <el-menu
                :default-active="activeMenu"
                :collapse="isCollapse"
                :unique-opened="true"
                router
                class="sidebar-menu"
              >
                <!-- 个人中心 -->
                <el-menu-item index="/dashboard" class="menu-item">
                  <el-icon><House /></el-icon>
                  <template #title>个人中心</template>
                </el-menu-item>

                <!-- 核心审批导航 -->
                <el-sub-menu index="approval" class="menu-group">
                  <template #title>
                    <el-icon><DocumentChecked /></el-icon>
                    <span>核心审批</span>
                  </template>
                  
                  <el-menu-item index="/approval/dashboard" class="sub-menu-item">
                    <el-icon><DataBoard /></el-icon>
                    <template #title>审批看板</template>
                  </el-menu-item>
                  
                  <el-menu-item index="/approval/loan" class="sub-menu-item">
                    <el-icon><CreditCard /></el-icon>
                    <template #title>贷款审批</template>
                  </el-menu-item>
                  
                  <el-menu-item index="/approval/lease" class="sub-menu-item">
                    <el-icon><Van /></el-icon>
                    <template #title>租赁审批</template>
                  </el-menu-item>
                  
                  <el-menu-item index="/approval/smart" class="sub-menu-item">
                    <el-icon><Cpu /></el-icon>
                    <template #title>智能审批</template>
                  </el-menu-item>
                </el-sub-menu>

                <!-- 运营管理 -->
                <el-sub-menu index="operation" class="menu-group" v-if="hasPermission('admin')">
                  <template #title>
                    <el-icon><Operation /></el-icon>
                    <span>运营管理</span>
                  </template>
                  
                  <el-menu-item index="/operation/products" class="sub-menu-item">
                    <el-icon><Goods /></el-icon>
                    <template #title>产品管理</template>
                  </el-menu-item>
                  
                  <el-menu-item index="/operation/users" class="sub-menu-item">
                    <el-icon><User /></el-icon>
                    <template #title>用户管理</template>
                  </el-menu-item>
                  
                  <el-menu-item index="/operation/statistics" class="sub-menu-item">
                    <el-icon><DataLine /></el-icon>
                    <template #title>数据统计</template>
                  </el-menu-item>
                </el-sub-menu>

                <!-- 系统管理 -->
                <el-sub-menu index="system" class="menu-group" v-if="hasPermission('admin')">
                  <template #title>
                    <el-icon><Setting /></el-icon>
                    <span>系统管理</span>
                  </template>
                  
                  <el-menu-item index="/system/config" class="sub-menu-item">
                    <el-icon><Tools /></el-icon>
                    <template #title>系统配置</template>
                  </el-menu-item>
                  
                  <el-menu-item index="/system/logs" class="sub-menu-item">
                    <el-icon><Document /></el-icon>
                    <template #title>操作日志</template>
                  </el-menu-item>
                  
                  <el-menu-item index="/system/security" class="sub-menu-item">
                    <el-icon><Lock /></el-icon>
                    <template #title>安全设置</template>
                  </el-menu-item>
                </el-sub-menu>
              </el-menu>
            </el-scrollbar>
          </div>
        </el-aside>

        <!-- 主内容区 -->
        <el-main class="layout-main">
          <div class="main-content">
            <router-view v-slot="{ Component, route }">
              <transition name="fade-transform" mode="out-in">
                <keep-alive>
                  <component :is="Component" :key="route.path" />
                </keep-alive>
              </transition>
            </router-view>
          </div>
        </el-main>
      </el-container>
    </el-container>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useUserStore } from '@/stores/user'
import { 
  House, Search, Bell, UserFilled, ArrowDown, User, Setting, SwitchButton,
  DocumentChecked, Expand, Fold, DataBoard, CreditCard, Van, Cpu,
  Operation, Goods, DataLine, Tools, Document, Lock, HomeFilled
} from '@element-plus/icons-vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import dayjs from 'dayjs'
import relativeTime from 'dayjs/plugin/relativeTime'
import 'dayjs/locale/zh-cn'

dayjs.extend(relativeTime)
dayjs.locale('zh-cn')

const route = useRoute()
const router = useRouter()
const userStore = useUserStore()

// 响应式数据
const isCollapse = ref(false)
const searchText = ref('')
const pendingApprovals = ref(0)
const notifications = ref([])

// 计算属性
const activeMenu = computed(() => route.path)
const currentUser = computed(() => userStore.userInfo)

// 方法
const toggleCollapse = () => {
  isCollapse.value = !isCollapse.value
}

const hasPermission = (permission: string) => {
  // 暂时允许所有登录用户访问
  return userStore.isLoggedIn
}

const getRoleName = (role: string) => {
  const roleMap = {
    'admin': '系统管理员',
    'manager': '审批经理',
    'user': '普通用户'
  }
  return roleMap[role] || '普通用户'
}

const formatTime = (time: string | Date) => {
  return dayjs(time).fromNow()
}

const handleUserCommand = (command: string) => {
  switch (command) {
    case 'profile':
      router.push('/profile')
      break
    case 'settings':
      router.push('/settings')
      break
    case 'logout':
      handleLogout()
      break
  }
}

const handleLogout = async () => {
  try {
    await ElMessageBox.confirm(
      '确定要退出登录吗？',
      '提示',
      {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning',
      }
    )
    
    userStore.logout()
    ElMessage.success('退出登录成功')
    router.push('/login')
  } catch (error) {
    console.log('用户取消退出')
  }
}

const handleNotificationClick = (notification: any) => {
  // 处理通知点击事件
  if (notification.link) {
    router.push(notification.link)
  }
}

const clearAllNotifications = () => {
  notifications.value = []
}

// 生命周期
onMounted(() => {
  // 初始化数据
  loadPendingApprovals()
  loadNotifications()
})

const loadPendingApprovals = async () => {
  // 加载待审批数量
  try {
    // 这里应该调用实际的API
    pendingApprovals.value = 12
  } catch (error) {
    console.error('加载待审批数量失败:', error)
  }
}

const loadNotifications = async () => {
  // 加载通知列表
  try {
    // 这里应该调用实际的API
    notifications.value = [
      {
        id: 1,
        title: '新的贷款申请',
        content: '用户张三提交了农业贷款申请，金额50万元',
        createdAt: new Date(),
        link: '/approval/loan'
      },
      {
        id: 2,
        title: '租赁审批通过',
        content: '拖拉机租赁申请已通过审批',
        createdAt: new Date(Date.now() - 3600000),
        link: '/approval/lease'
      }
    ]
  } catch (error) {
    console.error('加载通知失败:', error)
  }
}
</script>

<style scoped>
.oa-layout {
  height: 100vh;
  background-color: #f0f2f5;
}

.layout-container {
  height: 100%;
}

/* 顶部导航栏 */
.layout-header {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 0 24px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
  position: relative;
  z-index: 1000;
}

.header-left {
  display: flex;
  align-items: center;
  gap: 16px;
}

.logo {
  display: flex;
  align-items: center;
  gap: 12px;
}

.logo-icon {
  width: 36px;
  height: 36px;
  display: flex;
  align-items: center;
  justify-content: center;
  background: rgba(255, 255, 255, 0.2);
  border-radius: 8px;
  color: white;
  font-size: 20px;
}

.logo-text {
  font-size: 20px;
  font-weight: 600;
  color: white;
  white-space: nowrap;
}

.collapse-btn {
  color: white !important;
  font-size: 18px;
}

.collapse-btn:hover {
  background-color: rgba(255, 255, 255, 0.1) !important;
}

.header-center {
  flex: 1;
  max-width: 400px;
  margin: 0 40px;
}

.search-input {
  width: 100%;
}

.search-input :deep(.el-input__wrapper) {
  background-color: rgba(255, 255, 255, 0.2);
  border: 1px solid rgba(255, 255, 255, 0.3);
  border-radius: 20px;
  backdrop-filter: blur(10px);
}

.search-input :deep(.el-input__inner) {
  color: white;
}

.search-input :deep(.el-input__inner::placeholder) {
  color: rgba(255, 255, 255, 0.7);
}

.header-right {
  display: flex;
  align-items: center;
  gap: 24px;
}

.approval-status {
  display: flex;
  align-items: center;
}

.status-item {
  display: flex;
  align-items: center;
  gap: 8px;
  color: white;
  cursor: pointer;
  padding: 8px 12px;
  border-radius: 8px;
  transition: background-color 0.3s;
}

.status-item:hover {
  background-color: rgba(255, 255, 255, 0.1);
}

.status-icon {
  font-size: 18px;
}

.status-text {
  font-size: 14px;
  font-weight: 500;
}

.notification-btn {
  color: white !important;
  font-size: 18px;
  padding: 8px;
}

.notification-panel {
  padding: 0;
}

.notification-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 16px;
  border-bottom: 1px solid #ebeef5;
  font-weight: 600;
}

.notification-list {
  max-height: 300px;
  overflow-y: auto;
}

.notification-item {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  padding: 16px;
  border-bottom: 1px solid #f5f5f5;
  cursor: pointer;
  transition: background-color 0.3s;
}

.notification-item:hover {
  background-color: #f5f5f5;
}

.notification-content {
  flex: 1;
}

.notification-title {
  font-weight: 500;
  color: #333;
  margin-bottom: 4px;
}

.notification-desc {
  font-size: 12px;
  color: #666;
  line-height: 1.4;
}

.notification-time {
  font-size: 12px;
  color: #999;
  white-space: nowrap;
  margin-left: 12px;
}

.notification-empty {
  padding: 40px 16px;
  text-align: center;
  color: #999;
}

.user-dropdown {
  cursor: pointer;
}

.user-info {
  display: flex;
  align-items: center;
  gap: 12px;
  color: white;
}

.user-avatar {
  border: 2px solid rgba(255, 255, 255, 0.3);
}

.user-text {
  display: flex;
  flex-direction: column;
  align-items: flex-start;
}

.user-name {
  font-size: 14px;
  font-weight: 500;
  line-height: 1.2;
}

.user-role {
  font-size: 12px;
  opacity: 0.8;
  line-height: 1.2;
}

.user-arrow {
  font-size: 12px;
  opacity: 0.8;
}

/* 侧边栏 */
.layout-sidebar {
  background: white;
  border-right: 1px solid #e8eaec;
  transition: width 0.3s ease;
  overflow: hidden;
}

.sidebar-content {
  height: 100%;
  display: flex;
  flex-direction: column;
}

.sidebar-scrollbar {
  flex: 1;
}

.sidebar-menu {
  border: none;
  padding: 16px 0;
}

.menu-item {
  margin: 4px 16px;
  border-radius: 8px;
  height: 48px;
  line-height: 48px;
}

.menu-item:hover {
  background-color: #f5f7fa;
}

.menu-item.is-active {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: white;
}

.menu-group {
  margin: 8px 16px;
  border-radius: 8px;
}

.menu-group :deep(.el-sub-menu__title) {
  height: 48px;
  line-height: 48px;
  border-radius: 8px;
  padding-left: 16px;
  margin: 0;
}

.menu-group :deep(.el-sub-menu__title):hover {
  background-color: #f5f7fa;
}

.sub-menu-item {
  margin: 4px 0;
  border-radius: 6px;
  height: 40px;
  line-height: 40px;
}

.sub-menu-item:hover {
  background-color: #f5f7fa;
}

.sub-menu-item.is-active {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: white;
}

/* 主内容区 */
.layout-main {
  padding: 0;
  background-color: #f0f2f5;
  overflow-x: hidden;
}

.main-content {
  min-height: 100%;
  padding: 24px;
}

/* 页面切换动画 */
.fade-transform-enter-active,
.fade-transform-leave-active {
  transition: all 0.3s ease;
}

.fade-transform-enter-from {
  opacity: 0;
  transform: translateX(30px);
}

.fade-transform-leave-to {
  opacity: 0;
  transform: translateX(-30px);
}

/* 响应式设计 */
@media (max-width: 768px) {
  .header-center {
    display: none;
  }
  
  .header-right {
    gap: 16px;
  }
  
  .user-text {
    display: none;
  }
  
  .user-arrow {
    display: none;
  }
}
</style> 