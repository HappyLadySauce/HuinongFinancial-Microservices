<script setup lang="ts">
import { ref, onMounted, onBeforeUnmount } from 'vue';
import type { CarouselImage } from '../assets/images/machinery-images';

const props = defineProps({
  images: {
    type: Array as () => CarouselImage[],
    required: true
  },
  autoplayDelay: {
    type: Number,
    default: 5000
  },
  height: {
    type: String,
    default: '300px'
  }
});

const currentSlide = ref(0);
const totalSlides = ref(props.images.length);
let autoplayInterval: number | null = null;

// 切换到下一张幻灯片
const nextSlide = () => {
  currentSlide.value = (currentSlide.value + 1) % totalSlides.value;
};

// 切换到上一张幻灯片
const prevSlide = () => {
  currentSlide.value = (currentSlide.value - 1 + totalSlides.value) % totalSlides.value;
};

// 切换到指定幻灯片
const goToSlide = (index: number) => {
  currentSlide.value = index;
};

// 开始自动播放
const startAutoplay = () => {
  stopAutoplay();
  autoplayInterval = window.setInterval(() => {
    nextSlide();
  }, props.autoplayDelay);
};

// 停止自动播放
const stopAutoplay = () => {
  if (autoplayInterval !== null) {
    clearInterval(autoplayInterval);
    autoplayInterval = null;
  }
};

// 组件挂载时开始自动播放
onMounted(() => {
  startAutoplay();
});

// 组件卸载前停止自动播放
onBeforeUnmount(() => {
  stopAutoplay();
});

// 暂停和恢复自动播放 (鼠标悬停和离开时)
const handleMouseEnter = () => {
  stopAutoplay();
};

const handleMouseLeave = () => {
  startAutoplay();
};
</script>

<template>
  <div 
    class="carousel-container" 
    :style="{ height: props.height }"
    @mouseenter="handleMouseEnter"
    @mouseleave="handleMouseLeave"
  >
    <div class="carousel-slides" :style="{ transform: `translateX(-${currentSlide * 100}%)` }">
      <div 
        v-for="(image, index) in images" 
        :key="index" 
        class="carousel-slide"
      >
        <img :src="image.url" :alt="image.alt || '农机图片'" class="carousel-image">
      </div>
    </div>
    
    <!-- 轮播图导航按钮 -->
    <button class="carousel-btn prev-btn" @click="prevSlide">
      <span class="carousel-btn-icon">&lt;</span>
    </button>
    <button class="carousel-btn next-btn" @click="nextSlide">
      <span class="carousel-btn-icon">&gt;</span>
    </button>
    
    <!-- 轮播图指示器 -->
    <div class="carousel-indicators">
      <button 
        v-for="(_, index) in images" 
        :key="index" 
        class="indicator-dot"
        :class="{ active: currentSlide === index }"
        @click="goToSlide(index)"
      ></button>
    </div>
  </div>
</template>

<style scoped>
.carousel-container {
  position: relative;
  width: 100%;
  overflow: hidden;
  border-radius: 4px;
  box-shadow: 0 2px 8px rgba(0,0,0,0.1);
}

.carousel-slides {
  display: flex;
  height: 100%;
  transition: transform 0.5s ease;
}

.carousel-slide {
  min-width: 100%;
  height: 100%;
}

.carousel-image {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.carousel-btn {
  position: absolute;
  top: 50%;
  transform: translateY(-50%);
  background-color: rgba(255, 255, 255, 0.7);
  border: none;
  border-radius: 50%;
  width: 40px;
  height: 40px;
  display: flex;
  align-items: center;
  justify-content: center;
  cursor: pointer;
  opacity: 0.7;
  transition: opacity 0.3s ease;
  z-index: 10;
}

.carousel-btn:hover {
  opacity: 1;
}

.prev-btn {
  left: 10px;
}

.next-btn {
  right: 10px;
}

.carousel-btn-icon {
  font-size: 18px;
  font-weight: bold;
}

.carousel-indicators {
  position: absolute;
  bottom: 10px;
  left: 50%;
  transform: translateX(-50%);
  display: flex;
  gap: 8px;
  z-index: 10;
}

.indicator-dot {
  width: 10px;
  height: 10px;
  border-radius: 50%;
  background-color: rgba(255, 255, 255, 0.5);
  border: none;
  cursor: pointer;
  transition: background-color 0.3s ease;
}

.indicator-dot.active {
  background-color: #4CAF50;
}
</style> 