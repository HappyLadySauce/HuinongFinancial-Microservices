<script setup lang="ts">
import { ref, defineEmits } from 'vue'
import { useRouter } from 'vue-router'
import { HomeFilled, Money, Van, User } from '@element-plus/icons-vue'

const router = useRouter()
const emit = defineEmits(['update:activeTab'])

interface NavItem {
  id: string;
  icon: string;
  name: string;
  path: string;
}

// 导航项目
const navItems: NavItem[] = [
  { id: 'home', icon: 'HomeFilled', name: '首页', path: '/home' },
  { id: 'finance', icon: 'Money', name: '理财', path: '/finance' },
  { id: 'machinery', icon: 'Van', name: '农机租赁', path: '/machinery' },
  { id: 'me', icon: 'User', name: '我的', path: '/me' }
]

// 当前活跃的导航标签
const props = defineProps({
  activeTab: {
    type: String,
    default: 'home'
  }
})

// 处理导航点击
const handleNavClick = (item: NavItem) => {
  emit('update:activeTab', item.id)
  router.push(item.path)
}
</script>

<template>
  <!-- 底部导航栏 -->
  <div class="footer-nav">
    <div 
      v-for="item in navItems" 
      :key="item.id" 
      class="nav-item" 
      :class="{ active: props.activeTab === item.id }" 
      @click="handleNavClick(item)"
    >
      <el-icon><component :is="item.icon" /></el-icon>
      <div>{{ item.name }}</div>
    </div>
  </div>
</template>

<style scoped>
/* 底部导航栏 */
.footer-nav {
  position: fixed;
  bottom: 0;
  left: 0;
  right: 0;
  height: 60px;
  background-color: white;
  display: flex;
  justify-content: space-around;
  align-items: center;
  box-shadow: 0 -2px 5px rgba(0, 0, 0, 0.05);
  z-index: 99;
}

.nav-item {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  font-size: 12px;
  color: #666;
  width: 25%;
  height: 100%;
}

.nav-item .el-icon {
  font-size: 22px;
  margin-bottom: 3px;
}

.nav-item.active {
  color: #4CAF50; /* 主题色 */
  font-weight: 500;
}
</style>
