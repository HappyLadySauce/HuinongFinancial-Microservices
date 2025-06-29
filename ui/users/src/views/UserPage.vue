<script setup lang="ts">
import { ref } from 'vue'
import AppFooter from './components/footer.vue'
import '../assets/icons/agri-icons.css'

// 当前选中的导航栏
const activeTab = ref('user')

// 用户信息
const userInfo = ref({
  name: '张*',
  avatar: '',
  level: '普通会员',
  points: 1254,
  certification: false
})

// 我的资产数据
const assetData = ref({
  totalAmount: '15,280.25',
  savingAmount: '10,000.00',
  investAmount: '5,000.00',
  loanAmount: '280.25'
})

// 业务办理项目
const businessItems = [
  { id: 1, name: '农业贷款', icon: 'loan', color: '#4CAF50' },
  { id: 2, name: '农业补贴申请', icon: 'subsidy', color: '#FF9800' },
  { id: 3, name: '农业保险', icon: 'insurance', color: '#2196F3' },
  { id: 4, name: '养老金查询', icon: 'pension', color: '#9C27B0' },
  { id: 5, name: '低保申请', icon: 'welfare', color: '#F44336' },
  { id: 6, name: '医疗报销', icon: 'medical', color: '#00BCD4' },
  { id: 7, name: '水电费缴纳', icon: 'utility', color: '#795548' },
  { id: 8, name: '农机预约', icon: 'machinery', color: '#607D8B' }
]

// 更多服务项目
const serviceItems = [
  { id: 1, name: '在线客服', icon: 'customer-service' },
  { id: 2, name: '意见反馈', icon: 'feedback' },
  { id: 3, name: '设置', icon: 'settings' },
  { id: 4, name: '关于我们', icon: 'about' }
]

// 消息通知
const notificationCount = ref(3)

// 显示认证弹窗
const showCertModal = ref(false)
const startCertification = () => {
  showCertModal.value = true
}

// 关闭认证弹窗
const closeCertModal = () => {
  showCertModal.value = false
}

// 查看资产详情
const viewAssetDetail = (type: string) => {
  console.log('查看资产详情:', type)
}

// 办理业务
const handleBusiness = (item: any) => {
  console.log('办理业务:', item.name)
}

// 更多服务
const handleService = (item: any) => {
  console.log('使用服务:', item.name)
}
</script>

<template>
  <div class="user-container">
    <!-- 顶部导航 -->
    <div class="top-nav">
      <div class="nav-title">个人中心</div>
      <div class="notification" @click="handleService(serviceItems[0])">
        <i class="notification-icon"></i>
        <span v-if="notificationCount > 0" class="notification-badge">{{ notificationCount }}</span>
      </div>
    </div>
    
    <!-- 个人中心模块 -->
    <div class="user-profile">
      <div class="profile-header">
        <div class="profile-avatar-wrapper">
          <div class="profile-avatar" :style="userInfo.avatar ? `background-image: url(${userInfo.avatar})` : ''">
            <span v-if="!userInfo.avatar">{{ userInfo.name.substring(0, 1) }}</span>
          </div>
        </div>
        <div class="profile-info">
          <div class="profile-name-row">
            <h3 class="profile-name">{{ userInfo.name }}</h3>
            <div class="profile-level">{{ userInfo.level }}</div>
          </div>
          <div class="profile-points">积分: {{ userInfo.points }}</div>
          <div v-if="!userInfo.certification" class="profile-cert" @click="startCertification">
            <i class="cert-icon"></i>
            <span>未实名认证</span>
          </div>
          <div v-else class="profile-cert certified">
            <i class="cert-icon-done"></i>
            <span>已实名认证</span>
          </div>
        </div>
      </div>
    </div>

    <!-- 我的资产模块 -->
    <div class="module-card asset-module">
      <div class="module-header">
        <h3>我的资产</h3>
        <div class="module-more" @click="viewAssetDetail('all')">查看全部</div>
      </div>
      <div class="asset-total">
        <div class="asset-label">总资产(元)</div>
        <div class="asset-amount">{{ assetData.totalAmount }}</div>
      </div>
      <div class="asset-details">
        <div class="asset-item" @click="viewAssetDetail('saving')">
          <div class="asset-item-amount">{{ assetData.savingAmount }}</div>
          <div class="asset-item-name">存款</div>
        </div>
        <div class="asset-item" @click="viewAssetDetail('invest')">
          <div class="asset-item-amount">{{ assetData.investAmount }}</div>
          <div class="asset-item-name">理财</div>
        </div>
        <div class="asset-item" @click="viewAssetDetail('loan')">
          <div class="asset-item-amount">{{ assetData.loanAmount }}</div>
          <div class="asset-item-name">贷款</div>
        </div>
      </div>
    </div>

    <!-- 业务办理模块 -->
    <div class="module-card">
      <div class="module-header">
        <h3>业务办理</h3>
      </div>
      <div class="business-grid">
        <div v-for="item in businessItems" 
             :key="item.id" 
             class="business-item"
             @click="handleBusiness(item)">
          <div class="business-icon" :style="`background-color: ${item.color}`">
            <i :class="`agri-icon agri-icon-${item.icon}`"></i>
          </div>
          <div class="business-name">{{ item.name }}</div>
        </div>
      </div>
    </div>

    <!-- 更多服务模块 -->
    <div class="module-card">
      <div class="module-header">
        <h3>更多服务</h3>
      </div>
      <div class="service-list">
        <div v-for="item in serviceItems" 
             :key="item.id" 
             class="service-item"
             @click="handleService(item)">
          <div class="service-item-left">
            <i :class="`service-icon service-icon-${item.icon}`"></i>
            <span>{{ item.name }}</span>
          </div>
          <i class="arrow-right"></i>
        </div>
      </div>
    </div>

    <!-- 认证弹窗 -->
    <div v-if="showCertModal" class="cert-modal">
      <div class="cert-modal-content">
        <div class="cert-modal-header">
          <h3>实名认证</h3>
          <div class="cert-close" @click="closeCertModal">×</div>
        </div>
        <div class="cert-modal-body">
          <p>请完成实名认证，享受更多金融服务</p>
          <button class="cert-btn">立即认证</button>
        </div>
      </div>
    </div>

    <!-- 底部导航栏 -->
    <app-footer v-model:active-tab="activeTab" />
  </div>
</template>

<style scoped>
.user-container {
  padding: 0;
  text-align: left;
  min-height: 100vh;
  padding-bottom: 60px; /* 为底部导航栏留出空间 */
  background-color: #f5f5f5;
  color: #333;
  font-size: 14px;
}

/* 顶部导航栏 */
.top-nav {
  display: flex;
  justify-content: space-between;
  align-items: center;
  height: 44px;
  background-color: #4CAF50;
  color: #fff;
  padding: 0 15px;
  position: sticky;
  top: 0;
  z-index: 100;
}

.nav-title {
  font-size: 18px;
  font-weight: 500;
}

.notification {
  position: relative;
  width: 24px;
  height: 24px;
}

.notification-icon {
  width: 24px;
  height: 24px;
  background-image: url("data:image/svg+xml,%3Csvg xmlns='http://www.w3.org/2000/svg' viewBox='0 0 24 24'%3E%3Cpath fill='%23fff' d='M12,22c1.1,0,2-0.9,2-2h-4C10,21.1,10.9,22,12,22z M18,16v-5c0-3.07-1.63-5.64-4.5-6.32V4c0-0.83-0.67-1.5-1.5-1.5S10.5,3.17,10.5,4v0.68C7.64,5.36,6,7.92,6,11v5l-2,2v1h16v-1L18,16z'/%3E%3C/svg%3E");
  background-size: cover;
}

.notification-badge {
  position: absolute;
  top: -5px;
  right: -5px;
  background-color: #f44336;
  color: #fff;
  border-radius: 50%;
  width: 16px;
  height: 16px;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 10px;
}

/* 个人中心模块 */
.user-profile {
  background-color: #4CAF50;
  color: #fff;
  padding: 15px 15px 25px 15px;
  border-radius: 0 0 15px 15px;
}

.profile-header {
  display: flex;
  align-items: center;
}

.profile-avatar-wrapper {
  margin-right: 15px;
}

.profile-avatar {
  width: 64px;
  height: 64px;
  border-radius: 50%;
  background-color: #fff;
  color: #4CAF50;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 28px;
  font-weight: bold;
  background-size: cover;
  background-position: center;
  border: 2px solid rgba(255, 255, 255, 0.6);
}

.profile-info {
  flex: 1;
}

.profile-name-row {
  display: flex;
  align-items: center;
  margin-bottom: 5px;
}

.profile-name {
  font-size: 18px;
  font-weight: 500;
  margin: 0;
  margin-right: 10px;
}

.profile-level {
  background-color: rgba(255, 255, 255, 0.2);
  border-radius: 10px;
  font-size: 12px;
  padding: 2px 8px;
}

.profile-points {
  font-size: 12px;
  margin-bottom: 8px;
}

.profile-cert {
  display: flex;
  align-items: center;
  font-size: 12px;
  background-color: rgba(255, 255, 255, 0.2);
  padding: 3px 10px;
  border-radius: 15px;
  display: inline-flex;
}

.profile-cert.certified {
  background-color: rgba(255, 255, 255, 0.3);
}

.cert-icon {
  width: 14px;
  height: 14px;
  background-image: url("data:image/svg+xml,%3Csvg xmlns='http://www.w3.org/2000/svg' viewBox='0 0 24 24'%3E%3Cpath fill='%23fff' d='M12,17c-2.76,0-5-2.24-5-5s2.24-5,5-5s5,2.24,5,5S14.76,17,12,17z M12,9c-1.65,0-3,1.35-3,3s1.35,3,3,3s3-1.35,3-3S13.65,9,12,9z M18,20v1h-1v-1h-1v-1h1v-1h1v1h1v1H18z M18,4v1h-1V4h-1V3h1V2h1v1h1v1H18z M6,20v1H5v-1H4v-1h1v-1h1v1h1v1H6z M6,4v1H5V4H4V3h1V2h1v1h1v1H6z'/%3E%3C/svg%3E");
  background-size: cover;
  margin-right: 4px;
}

.cert-icon-done {
  width: 14px;
  height: 14px;
  background-image: url("data:image/svg+xml,%3Csvg xmlns='http://www.w3.org/2000/svg' viewBox='0 0 24 24'%3E%3Cpath fill='%23fff' d='M9,16.17L4.83,12l-1.42,1.41L9,19L21,7l-1.41-1.41L9,16.17z'/%3E%3C/svg%3E");
  background-size: cover;
  margin-right: 4px;
}

/* 模块卡片通用样式 */
.module-card {
  background-color: #fff;
  margin: 12px;
  border-radius: 10px;
  padding: 15px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.05);
}

.module-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 15px;
}

.module-header h3 {
  font-size: 16px;
  font-weight: 500;
  margin: 0;
  color: #333;
}

.module-more {
  color: #666;
  font-size: 12px;
}

/* 资产模块 */
.asset-module {
  margin-top: -15px;
  position: relative;
  z-index: 10;
}

.asset-total {
  text-align: center;
  margin-bottom: 15px;
}

.asset-label {
  font-size: 12px;
  color: #999;
  margin-bottom: 5px;
}

.asset-amount {
  font-size: 24px;
  font-weight: bold;
  color: #333;
}

.asset-details {
  display: flex;
  justify-content: space-between;
  border-top: 1px solid #f0f0f0;
  padding-top: 15px;
}

.asset-item {
  flex: 1;
  text-align: center;
  padding: 0 5px;
}

.asset-item-amount {
  font-size: 16px;
  font-weight: 500;
  color: #333;
  margin-bottom: 5px;
}

.asset-item-name {
  font-size: 12px;
  color: #666;
}

/* 业务办理模块 */
.business-grid {
  display: grid;
  grid-template-columns: repeat(4, 1fr);
  gap: 15px;
}

.business-item {
  display: flex;
  flex-direction: column;
  align-items: center;
}

.business-icon {
  width: 40px;
  height: 40px;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  margin-bottom: 5px;
}

.business-icon i {
  color: #fff;
  font-size: 20px;
}

.business-name {
  font-size: 12px;
  color: #333;
  text-align: center;
}

/* 农业图标样式占位符 */
.agri-icon {
  width: 24px;
  height: 24px;
  display: inline-block;
}

/* 更多服务模块 */
.service-list {
  margin-top: -5px;
}

.service-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 12px 0;
  border-bottom: 1px solid #f0f0f0;
}

.service-item:last-child {
  border-bottom: none;
}

.service-item-left {
  display: flex;
  align-items: center;
}

.service-icon {
  width: 24px;
  height: 24px;
  margin-right: 10px;
}

/* 各个服务图标的样式 */
.service-icon-customer-service {
  background-image: url("data:image/svg+xml,%3Csvg xmlns='http://www.w3.org/2000/svg' viewBox='0 0 24 24'%3E%3Cpath fill='%234CAF50' d='M12,1C7.03,1,3,5.03,3,10v4c0,1.1,0.9,2,2,2h2v-6H5v-0.5C5,6.13,8.13,3,12,3s7,3.13,7,7.5V10h-2v6h2c1.1,0,2-0.9,2-2v-4C21,5.03,16.97,1,12,1z M12,17c-1.66,0-3,1.34-3,3s1.34,3,3,3s3-1.34,3-3S13.66,17,12,17z'/%3E%3C/svg%3E");
}

.service-icon-feedback {
  background-image: url("data:image/svg+xml,%3Csvg xmlns='http://www.w3.org/2000/svg' viewBox='0 0 24 24'%3E%3Cpath fill='%234CAF50' d='M20,2H4C2.9,2,2,2.9,2,4v12c0,1.1,0.9,2,2,2h14l4,4L20,2z M18,14H6v-2h12V14z M18,11H6V9h12V11z M18,8H6V6h12V8z'/%3E%3C/svg%3E");
}

.service-icon-settings {
  background-image: url("data:image/svg+xml,%3Csvg xmlns='http://www.w3.org/2000/svg' viewBox='0 0 24 24'%3E%3Cpath fill='%234CAF50' d='M19.14,12.94c0.04-0.3,0.06-0.61,0.06-0.94c0-0.32-0.02-0.64-0.07-0.94l2.03-1.58c0.18-0.14,0.23-0.41,0.12-0.61 l-1.92-3.32c-0.12-0.22-0.37-0.29-0.59-0.22l-2.39,0.96c-0.5-0.38-1.03-0.7-1.62-0.94L14.4,2.81c-0.04-0.24-0.24-0.41-0.48-0.41 h-3.84c-0.24,0-0.43,0.17-0.47,0.41L9.25,5.35C8.66,5.59,8.12,5.92,7.63,6.29L5.24,5.33c-0.22-0.08-0.47,0-0.59,0.22L2.74,8.87 C2.62,9.08,2.66,9.34,2.86,9.48l2.03,1.58C4.84,11.36,4.8,11.69,4.8,12s0.02,0.64,0.07,0.94l-2.03,1.58 c-0.18,0.14-0.23,0.41-0.12,0.61l1.92,3.32c0.12,0.22,0.37,0.29,0.59,0.22l2.39-0.96c0.5,0.38,1.03,0.7,1.62,0.94l0.36,2.54 c0.05,0.24,0.24,0.41,0.48,0.41h3.84c0.24,0,0.44-0.17,0.47-0.41l0.36-2.54c0.59-0.24,1.13-0.56,1.62-0.94l2.39,0.96 c0.22,0.08,0.47,0,0.59-0.22l1.92-3.32c0.12-0.22,0.07-0.47-0.12-0.61L19.14,12.94z M12,15.6c-1.98,0-3.6-1.62-3.6-3.6 s1.62-3.6,3.6-3.6s3.6,1.62,3.6,3.6S13.98,15.6,12,15.6z'/%3E%3C/svg%3E");
}

.service-icon-about {
  background-image: url("data:image/svg+xml,%3Csvg xmlns='http://www.w3.org/2000/svg' viewBox='0 0 24 24'%3E%3Cpath fill='%234CAF50' d='M12,2C6.48,2,2,6.48,2,12s4.48,10,10,10s10-4.48,10-10S17.52,2,12,2z M13,17h-2v-6h2V17z M13,9h-2V7h2V9z'/%3E%3C/svg%3E");
}

.arrow-right {
  width: 16px;
  height: 16px;
  background-image: url("data:image/svg+xml,%3Csvg xmlns='http://www.w3.org/2000/svg' viewBox='0 0 24 24'%3E%3Cpath fill='%23999' d='M10,6L8.59,7.41L13.17,12l-4.58,4.59L10,18l6-6L10,6z'/%3E%3C/svg%3E");
  background-size: cover;
}

/* 认证弹窗 */
.cert-modal {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background-color: rgba(0, 0, 0, 0.5);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 1000;
}

.cert-modal-content {
  width: 80%;
  max-width: 320px;
  background-color: #fff;
  border-radius: 10px;
  overflow: hidden;
}

.cert-modal-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 15px;
  border-bottom: 1px solid #f0f0f0;
}

.cert-modal-header h3 {
  margin: 0;
  font-size: 16px;
}

.cert-close {
  font-size: 24px;
  color: #999;
  cursor: pointer;
}

.cert-modal-body {
  padding: 20px 15px;
  text-align: center;
}

.cert-modal-body p {
  margin-bottom: 20px;
  color: #666;
}

.cert-btn {
  background-color: #4CAF50;
  color: #fff;
  border: none;
  border-radius: 20px;
  padding: 8px 30px;
  font-size: 14px;
}

/* 媒体查询，适配不同尺寸 */
@media screen and (max-width: 320px) {
  .business-grid {
    gap: 10px;
  }
  
  .business-name {
    font-size: 11px;
  }
}

@media screen and (min-width: 480px) {
  .user-container {
    max-width: 500px;
    margin: 0 auto;
  }
}
</style> 