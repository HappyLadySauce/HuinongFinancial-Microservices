<script setup lang="ts">
import { ref, onMounted, onUnmounted } from 'vue'
import { useRouter } from 'vue-router'

const router = useRouter()
const showFirstImage = ref(true)
const showSecondImage = ref(false)
let imageTimer: number | null = null

// 使用import.meta.url加载图片
const image1 = new URL('../assets/images/数字人民币1.jpg', import.meta.url).href
const image2 = new URL('../assets/images/数字人民币2.png', import.meta.url).href

const goBack = () => {
  if (imageTimer) {
    clearTimeout(imageTimer)
  }
  window.history.back()
}

onMounted(() => {
  // 第一张图片显示2秒
  imageTimer = setTimeout(() => {
    showFirstImage.value = false
    showSecondImage.value = true
    
    // 第二张图片显示1分钟后返回
    imageTimer = setTimeout(() => {
      window.history.back()
    }, 60000) // 1分钟 = 60000毫秒
  }, 2000) // 2秒 = 2000毫秒
})

onUnmounted(() => {
  // 清除定时器，避免内存泄漏
  if (imageTimer) {
    clearTimeout(imageTimer)
  }
})
</script>

<template>
  <div class="digital-currency-page">
    <!-- 返回按钮 -->
    <div v-if="showSecondImage" class="back-button" @click="goBack">
      <svg viewBox="0 0 24 24" width="24" height="24">
        <path d="M20,11H7.83l5.59-5.59L12,4l-8,8l8,8l1.41-1.41L7.83,13H20V11z" fill="#ffffff"/>
      </svg>
    </div>

    <div class="image-container">
      <!-- 第一张图片 -->
      <div v-if="showFirstImage" class="image-wrapper">
        <img :src="image1" alt="数字人民币启动页" class="fill-image" />
      </div>

      <!-- 第二张图片 -->
      <div v-if="showSecondImage" class="image-wrapper">
        <img :src="image2" alt="数字人民币钱包" class="fill-image" />
      </div>
    </div>
  </div>
</template>

<style scoped>
.digital-currency-page {
  position: relative;
  width: 100%;
  height: 100vh;
  background-color: #fff;
  overflow: hidden;
}

.image-container {
  width: 100%;
  height: 100vh;
  position: relative;
}

.image-wrapper {
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  display: flex;
  align-items: center;
  justify-content: center;
}

.fill-image {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.back-button {
  position: fixed;
  top: 16px;
  left: 16px;
  width: 40px;
  height: 40px;
  background-color: rgba(0, 0, 0, 0.5);
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 100;
  cursor: pointer;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.2);
}

.back-button:active {
  background-color: rgba(0, 0, 0, 0.7);
  transform: scale(0.95);
}
</style> 