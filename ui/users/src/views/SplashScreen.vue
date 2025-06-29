<template>
  <div class="splash-container" :class="{ 'fade-out': isFadingOut }">
    <div class="splash-content">
      <!-- Logo动画 -->
      <div 
        class="logo-container"
        :style="{ transform: `scale(${logoScale})`, opacity: opacity }"
      >
        <img src="../assets/images/logo.png" alt="数字惠农" class="logo-image" />
      </div>
      
      <!-- APP名称 -->
      <div 
        class="app-title"
        :style="{ opacity: titleOpacity }"
      >
        数字惠农
      </div>
      
      <!-- 标语 -->
      <div 
        class="app-slogan"
        :style="{ opacity: sloganOpacity }"
      >
        数字惠农，助力美丽新乡村！
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'

const router = useRouter()
const logoScale = ref(0.5)
const opacity = ref(0)
const titleOpacity = ref(0)
const sloganOpacity = ref(0)
const isFadingOut = ref(false)

onMounted(() => {
  // 加载动画
  setTimeout(() => {
    logoScale.value = 1
    opacity.value = 1
    
    // 名称淡入
    setTimeout(() => {
      titleOpacity.value = 1
      
      // 标语淡入
      setTimeout(() => {
        sloganOpacity.value = 1
        
        // 等待一段时间后开始淡出
        setTimeout(() => {
          // 开始淡出动画
          isFadingOut.value = true
          
          // 等待动画完成后跳转到登录页
          setTimeout(() => {
            router.push('/login')
          }, 1000)
        }, 900)
      }, 400)
    }, 400)
  }, 300)
})
</script>

<style scoped>
.splash-container {
  width: 100vw;
  height: 100vh;
  display: flex;
  justify-content: center;
  align-items: center;
  background-color: #fff;
  overflow: hidden;
  padding: 20px;
  box-sizing: border-box;
  transition: opacity 1s ease-out;
}

.fade-out {
  opacity: 0;
}

.splash-content {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  width: 100%;
  max-width: 800px;
  transform: scale(1);
}

.logo-container {
  margin-bottom: clamp(10px, 3vh, 20px);
  transition: all 1s ease;
}

.logo-image {
  width: clamp(160px, 25vw, 240px);
  height: clamp(160px, 25vw, 240px);
  object-fit: contain;
}

.app-title {
  font-size: clamp(36px, 8vw, 54px);
  font-weight: bold;
  margin-bottom: clamp(20px, 5vh, 30px);
  background-image: linear-gradient(to right, #FFD700, #C49000, #FFD700);
  -webkit-background-clip: text;
  background-clip: text;
  -webkit-text-fill-color: transparent;
  transition: opacity 0.8s ease;
  text-align: center;
  width: 100%;
  font-family: "STZhongsong", "华文中宋", "SimSun", "宋体", serif;
  letter-spacing: 4px;
  text-shadow: 0 4px 8px rgba(169, 130, 0, 0.2);
}

.app-slogan {
  /* 字体调小 */
  font-size: clamp(28px, 8vw, 36px);
  color: #666;
  text-align: center;
  transition: opacity 0.8s ease;
  width: 100%;
  line-height: 1.5;
  font-family: "STKaiti", "华文楷体", "YouYuan", "幼圆", sans-serif;
  letter-spacing: 2px;
}

/* 横屏模式适配 */
@media screen and (orientation: landscape) and (max-height: 500px) {
  .splash-content {
    flex-direction: row;
    align-items: center;
    max-width: 90%;
  }
  
  .logo-container {
    margin-bottom: 0;
    margin-right: 30px;
  }
  
  .logo-image {
    width: clamp(120px, 20vh, 180px);
    height: clamp(120px, 20vh, 180px);
  }
  
  .app-title {
    font-size: clamp(34px, 7vh, 46px);
    margin-bottom: 16px;
    text-align: left;
  }
  
  .app-slogan {
    font-size: clamp(32px, 6vh, 40px);
    text-align: left;
  }
}

/* 小屏幕手机适配 */
@media screen and (max-width: 320px) {
  .logo-image {
    width: 140px;
    height: 140px;
  }
  
  .app-title {
    font-size: 36px;
  }
  
  .app-slogan {
    font-size: 36px;
  }
}

/* 平板设备适配 */
@media screen and (min-width: 768px) and (max-width: 1024px) {
  .logo-container {
    margin-bottom: 40px;
  }
  
  .app-title {
    margin-bottom: 30px;
  }
}

/* 大屏设备适配 */
@media screen and (min-width: 1025px) {
  .logo-image {
    width: 280px;
    height: 280px;
  }
  
  .app-title {
    font-size: 60px;
    margin-bottom: 40px;
  }
  
  .app-slogan {
    font-size: 48px;
  }
}

/* 超窄屏幕设备适配 */
@media screen and (max-width: 280px) {
  .splash-content {
    transform: scale(0.8);
  }
  
  .logo-image {
    width: 120px;
    height: 120px;
  }
  
  .app-title {
    font-size: 34px;
  }
  
  .app-slogan {
    font-size: 32px;
  }
}

/* 适配奇怪屏幕比例 */
@media screen and (max-aspect-ratio: 3/4) and (max-width: 400px) {
  .splash-content {
    transform: scale(0.8);
  }
}

@media screen and (min-aspect-ratio: 18/9) {
  .splash-content {
    flex-direction: row;
    max-width: 900px;
  }
  
  .logo-container {
    margin-right: 50px;
    margin-bottom: 0;
  }
}
</style> 