<script setup lang="ts">
defineProps({
  type: {
    type: String,
    required: true
  },
  size: {
    type: Number,
    default: 24
  }
})
</script>

<template>
  <div class="service-icon">
    <div class="inner-icon">
      <!-- 根据type选择对应图标 -->
      <div v-if="type === 'loan'" class="icon-loan">
        <!-- 惠农贷 - 钱币图标 -->
        <svg viewBox="0 0 40 40" xmlns="http://www.w3.org/2000/svg">
          <circle cx="20" cy="20" r="16" fill="#60a917"/>
          <path d="M24 16V14C24 12.9 23.1 12 22 12H16V20H22C23.1 20 24 19.1 24 18V16ZM22 16H20V14H22V16Z" fill="white"/>
          <path d="M16 22V24C16 25.1 16.9 26 18 26H24V18H18C16.9 18 16 18.9 16 20V22ZM18 22H20V24H18V22Z" fill="white"/>
        </svg>
      </div>

      <div v-else-if="type === 'insurance'" class="icon-insurance">
        <!-- 惠农医保 - 医疗十字图标 -->
        <svg viewBox="0 0 40 40" xmlns="http://www.w3.org/2000/svg">
          <circle cx="20" cy="20" r="16" fill="#4FC3F7"/>
          <path d="M28 18H22V12H18V18H12V22H18V28H22V22H28V18Z" fill="white"/>
        </svg>
      </div>

      <div v-else-if="type === 'subsidy'" class="icon-subsidy">
        <!-- 农业补贴 - 小麦图标 -->
        <svg viewBox="0 0 40 40" xmlns="http://www.w3.org/2000/svg">
          <circle cx="20" cy="20" r="16" fill="#FF9800"/>
          <path d="M20 10C17 10 13 11.5 13 15C13 16.7 14.5 18.4 16 20C17.4 21.4 19 22.5 19 25H21C21 22.5 22.6 21.4 24 20C25.5 18.4 27 16.7 27 15C27 11.5 23 10 20 10Z" fill="white"/>
          <path d="M18 28H22V31H18V28Z" fill="white"/>
        </svg>
      </div>

      <div v-else-if="type === 'pension'" class="icon-pension">
        <!-- 个人养老金 - 银行卡图标 -->
        <svg viewBox="0 0 40 40" xmlns="http://www.w3.org/2000/svg">
          <circle cx="20" cy="20" r="16" fill="#8BC34A"/>
          <path d="M28 14H12V16H28V14ZM28 18H12V26H28V18ZM14 22H20V24H14V22Z" fill="white"/>
        </svg>
      </div>

      <div v-else-if="type === 'expert'" class="icon-expert">
        <!-- 农技专家 - 农民图标 -->
        <svg viewBox="0 0 40 40" xmlns="http://www.w3.org/2000/svg">
          <circle cx="20" cy="20" r="16" fill="#795548"/>
          <path d="M20 12C21.66 12 23 13.34 23 15C23 16.66 21.66 18 20 18C18.34 18 17 16.66 17 15C17 13.34 18.34 12 20 12Z" fill="#FFD54F"/>
          <path d="M24 19H16C14.9 19 14 19.9 14 21V27H16V21H24V27H26V21C26 19.9 25.1 19 24 19Z" fill="#FFD54F"/>
          <path d="M22 16V13L25 15L22 16Z" fill="#FFD54F"/>
        </svg>
      </div>

      <div v-else-if="type === 'machinery'" class="icon-machinery">
        <!-- 农机租赁 - 拖拉机图标 -->
        <svg viewBox="0 0 40 40" xmlns="http://www.w3.org/2000/svg">
          <circle cx="20" cy="20" r="16" fill="#4CAF50"/>
          <path d="M27 14H24V10H14C12.9 10 12 10.9 12 12V22H14C14 23.66 15.34 25 17 25C18.66 25 20 23.66 20 22H24C24 23.66 25.34 25 27 25C28.66 25 30 23.66 30 22H31V18L27 14ZM27 20.5C26.17 20.5 25.5 21.17 25.5 22C25.5 22.83 26.17 23.5 27 23.5C27.83 23.5 28.5 22.83 28.5 22C28.5 21.17 27.83 20.5 27 20.5ZM26 18L28.59 15H24V18H26ZM17 23.5C16.17 23.5 15.5 22.83 15.5 22C15.5 21.17 16.17 20.5 17 20.5C17.83 20.5 18.5 21.17 18.5 22C18.5 22.83 17.83 23.5 17 23.5Z" fill="white"/>
        </svg>
      </div>

      <div v-else-if="type === 'digital-currency'" class="icon-digital-currency">
        <!-- 数字人民币 - 人民币图标 -->
        <svg viewBox="0 0 40 40" xmlns="http://www.w3.org/2000/svg">
          <circle cx="20" cy="20" r="16" fill="#69B342"/>
          <path d="M20 10C14.48 10 10 14.48 10 20C10 25.52 14.48 30 20 30C25.52 30 30 25.52 30 20C30 14.48 25.52 10 20 10ZM23 27H21V24H19V27H17V24H14V22H17V20H14V18H16.25L13 13H15.75L19 18H20L23.25 13H26L22.75 18H25V20H22V22H25V24H23V27Z" fill="white"/>
        </svg>
      </div>

      <div v-else-if="type === 'welfare'" class="icon-welfare">
        <!-- 低保服务 - 盾牌图标 -->
        <svg viewBox="0 0 40 40" xmlns="http://www.w3.org/2000/svg">
          <circle cx="20" cy="20" r="16" fill="#5D9CEC"/>
          <path d="M20 10L12 14V19C12 23.4 15.4 27.2 20 28C24.6 27.2 28 23.4 28 19V14L20 10ZM26 19C26 22.31 23.66 25.18 20 25.91C16.34 25.18 14 22.31 14 19V15.34L20 12.32L26 15.34V19Z" fill="white"/>
          <path d="M18 17H22V23H18V17Z" fill="white"/>
        </svg>
      </div>

      <!-- 默认图标 -->
      <div v-else class="icon-default">
        <svg viewBox="0 0 40 40" xmlns="http://www.w3.org/2000/svg">
          <circle cx="20" cy="20" r="16" fill="#9E9E9E"/>
          <path d="M19 15H21V21H19V15ZM19 23H21V25H19V23Z" fill="white"/>
        </svg>
      </div>
    </div>
  </div>
</template>

<style scoped>
.service-icon {
  width: 50px;
  height: 50px;
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 0;
}

.inner-icon {
  width: 40px;
  height: 40px;
  display: flex;
  align-items: center;
  justify-content: center;
}

svg {
  width: 100%;
  height: 100%;
}
</style> 