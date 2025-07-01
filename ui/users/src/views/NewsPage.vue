<script setup lang="ts">
import { ref } from 'vue'
import AppFooter from './components/footer.vue'

// 新闻页面待开发

// 当前选中的导航栏 - 新闻页面没有对应的底部导航，保持首页选中
const activeTab = ref('home')

// 新闻数据模拟
const newsItems = ref([
  {
    id: 1, 
    title: '2025年农业补贴政策解读',
    intro: '农业农村部最新发布的农业补贴政策，种粮农民有望获得更多支持，提高粮食生产积极性。',
    image: '/src/assets/images/news1.png',
    date: '2025-05-01',
    author: '惠农金融',
    reads: '2456',
    url: '/news/detail/1'
  },
  {
    id: 2, 
    title: '夏粮增产增收 农业"半年报"亮眼',
    intro: '今年夏粮再获丰收，全国夏粮总产量比上年增加60多亿斤，为保障国家粮食安全提供了有力支撑。',
    image: '/src/assets/images/news2.png',
    date: '2025-05-15',
    author: '农业日报',
    reads: '3782',
    url: '/news/detail/2'
  },
  {
    id: 3, 
    title: '数字技术助力乡村振兴',
    intro: '各地积极推动数字技术与农业生产深度融合，打造数字乡村，推进农村现代化，助力乡村产业振兴。',
    image: '/src/assets/images/news3.png',
    date: '2025-05-20',
    author: '科技农业网',
    reads: '1845',
    url: '/news/detail/3'
  },
  {
    id: 4, 
    title: '惠农金融新产品上线，助力农民增收',
    intro: '惠农金融推出专为农户设计的普惠金融产品，低息贷款支持农业生产，解决农民资金难题。',
    image: '/src/assets/images/news4.png',
    date: '2025-05-25',
    author: '惠农金融',
    reads: '2103',
    url: '/news/detail/4'
  },
  {
    id: 5, 
    title: '农村集体经济组织如何发展壮大',
    intro: '发展壮大农村集体经济是实施乡村振兴战略的重要内容，本文分析了集体经济发展的多种模式。',
    image: '/src/assets/images/news5.png',
    date: '2025-05-28',
    author: '农村经济研究',
    reads: '1678',
    url: '/news/detail/5'
  }
])

// 新闻分类
const categories = ref([
  { id: 1, name: '全部' },
  { id: 2, name: '政策解读' },
  { id: 3, name: '农业科技' },
  { id: 4, name: '金融服务' },
  { id: 5, name: '乡村振兴' },
  { id: 6, name: '市场行情' }
])

// 当前选中的分类
const currentCategory = ref(categories.value[0])

// 切换分类
const switchCategory = (category: { id: number; name: string }) => {
  currentCategory.value = category
  // 这里应该有根据分类过滤新闻的逻辑
}
</script>

<template>
  <div class="news-container">
    <!-- 顶部导航 -->
    <div class="top-nav">
      <div class="nav-title">农业资讯</div>
    </div>
    
    <!-- 分类导航 -->
    <div class="category-nav">
      <div class="category-scroll">
        <div
          v-for="category in categories"
          :key="category.id"
          :class="['category-item', { active: currentCategory.id === category.id }]"
          @click="switchCategory(category)"
        >
          {{ category.name }}
        </div>
      </div>
    </div>
    
    <!-- 新闻列表 -->
    <div class="news-list">
      <div class="news-item" v-for="news in newsItems" :key="news.id">
        <div class="news-content">
          <div class="news-title">{{ news.title }}</div>
          <div class="news-intro">{{ news.intro }}</div>
          <div class="news-meta">
            <span class="news-source">{{ news.author }}</span>
            <span class="news-date">{{ news.date }}</span>
            <span class="news-reads">{{ news.reads }} 阅读</span>
          </div>
        </div>
        <div class="news-image">
          <img :src="news.image" alt="资讯图片" class="news-img"/>
        </div>
      </div>
    </div>
    
    <!-- 底部导航栏 -->
    <app-footer v-model:active-tab="activeTab" />
  </div>
</template>

<style scoped>
.news-container {
  padding: 0;
  text-align: left;
  min-height: 100vh;
  padding-bottom: 60px; /* 为底部导航栏留出空间 */
  background-color: #f5f5f5;
}

/* 顶部导航 */
.top-nav {
  display: flex;
  align-items: center;
  justify-content: center;
  height: 44px;
  background-color: #fff;
  position: sticky;
  top: 0;
  z-index: 100;
  border-bottom: 1px solid #eee;
}

.nav-title {
  font-size: 18px;
  font-weight: 500;
}

/* 分类导航 */
.category-nav {
  background-color: #fff;
  overflow-x: auto;
  white-space: nowrap;
  -webkit-overflow-scrolling: touch;
  scrollbar-width: none; /* Firefox */
  position: sticky;
  top: 44px;
  z-index: 99;
}

.category-nav::-webkit-scrollbar {
  display: none; /* Chrome, Safari, Edge */
}

.category-scroll {
  display: inline-flex;
  padding: 0 10px;
}

.category-item {
  padding: 12px 15px;
  font-size: 14px;
  color: #666;
  position: relative;
  transition: color 0.3s;
}

.category-item.active {
  color: #4CAF50;
  font-weight: 500;
}

.category-item.active::after {
  content: '';
  position: absolute;
  width: 20px;
  height: 2px;
  background-color: #4CAF50;
  bottom: 8px;
  left: 50%;
  transform: translateX(-50%);
  border-radius: 1px;
}

/* 新闻列表 */
.news-list {
  padding: 10px 12px;
}

.news-item {
  background-color: #fff;
  border-radius: 8px;
  margin-bottom: 10px;
  padding: 12px;
  display: flex;
  cursor: pointer;
  box-shadow: 0 1px 3px rgba(0,0,0,0.05);
  transition: transform 0.2s, box-shadow 0.2s;
}

.news-item:active {
  transform: scale(0.98);
}

@media screen and (min-width: 768px) {
  .news-list {
    max-width: 800px;
    margin: 0 auto;
    padding: 15px;
  }
  
  .news-item {
    padding: 15px;
    margin-bottom: 15px;
  }
}

.news-content {
  flex: 1;
  margin-right: 12px;
  display: flex;
  flex-direction: column;
}

.news-title {
  font-size: 16px;
  font-weight: 500;
  color: #333;
  margin-bottom: 8px;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
  line-height: 1.4;
}

@media screen and (max-width: 375px) {
  .news-title {
    font-size: 15px;
    margin-bottom: 6px;
  }
}

.news-intro {
  font-size: 14px;
  color: #666;
  margin-bottom: 10px;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
  line-height: 1.5;
  flex: 1;
}

@media screen and (max-width: 375px) {
  .news-intro {
    font-size: 13px;
    margin-bottom: 8px;
    -webkit-line-clamp: 2;
  }
}

.news-meta {
  font-size: 12px;
  color: #999;
  display: flex;
  gap: 12px;
}

.news-image {
  width: 110px;
  height: 80px;
  border-radius: 4px;
  overflow: hidden;
  flex-shrink: 0;
}

@media screen and (max-width: 375px) {
  .news-image {
    width: 100px;
    height: 70px;
  }
}

@media screen and (min-width: 480px) {
  .news-image {
    width: 120px;
    height: 90px;
  }
}

.news-img {
  width: 100%;
  height: 100%;
  object-fit: cover;
}
</style> 