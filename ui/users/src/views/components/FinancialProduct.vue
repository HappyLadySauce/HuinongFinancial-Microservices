<script setup lang="ts">
import { defineProps } from 'vue';

const props = defineProps({
  // 产品类型: 股票型、债券型等
  type: {
    type: String,
    default: '',
  },
  // 产品名称
  name: {
    type: String,
    required: true,
  },
  // 年化收益率
  annualReturn: {
    type: String,
    default: '0.00%',
  },
  // 最短投资期限
  minPeriod: {
    type: String,
    default: '无固定期限',
  },
  // 风险等级: 低、中、高
  riskLevel: {
    type: String,
    default: '中',
  },
  // 农业相关标签
  agriTag: {
    type: String,
    default: '',
  },
  // 是否为推荐产品
  isRecommended: {
    type: Boolean,
    default: false,
  }
});
</script>

<template>
  <div class="financial-product-card" :class="{ 'is-recommended': isRecommended }">
    <div class="product-header">
      <span class="product-type">{{ type }}</span>
      <span v-if="agriTag" class="product-tag">{{ agriTag }}</span>
    </div>
    
    <div class="product-name">{{ name }}</div>
    
    <div class="product-info">
      <div class="return-info">
        <div class="return-rate">{{ annualReturn }}</div>
        <div class="return-label">预期年化</div>
      </div>
      
      <div class="period-info">
        <div class="period">{{ minPeriod }}</div>
        <div class="period-label">投资期限</div>
      </div>
    </div>
    
    <div class="risk-level" :class="`risk-${riskLevel}`">
      风险等级：{{ riskLevel }}
    </div>
    
    <div class="product-button">
      立即购买
    </div>
  </div>
</template>

<style scoped>
.financial-product-card {
  width: 100%;
  background-color: #fff;
  border-radius: 10px;
  padding: 15px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.05);
  margin-bottom: 15px;
  transition: transform 0.2s ease, box-shadow 0.2s ease;
}

/* 响应式布局 */
@media screen and (min-width: 480px) {
  .financial-product-card {
    width: calc(50% - 10px);
  }
}

@media screen and (min-width: 768px) {
  .financial-product-card {
    width: calc(33.333% - 10px);
  }
}

@media screen and (max-width: 375px) {
  .financial-product-card {
    padding: 12px;
    margin-bottom: 12px;
  }
}

.financial-product-card:active {
  transform: scale(0.98);
}

.is-recommended {
  border: 1px solid #4CAF50;
  position: relative;
}

.is-recommended::after {
  content: '推荐';
  position: absolute;
  top: 0;
  right: 0;
  background-color: #4CAF50;
  color: white;
  font-size: 12px;
  padding: 2px 8px;
  border-radius: 0 10px 0 10px;
}

.product-header {
  display: flex;
  justify-content: space-between;
  margin-bottom: 10px;
}

.product-type {
  background-color: #f5f5f5;
  color: #666;
  padding: 2px 8px;
  border-radius: 4px;
  font-size: 12px;
}

.product-tag {
  background-color: #e8f5e9;
  color: #4CAF50;
  padding: 2px 8px;
  border-radius: 4px;
  font-size: 12px;
}

.product-name {
  font-size: 16px;
  font-weight: 500;
  margin-bottom: 15px;
  line-height: 1.4;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
}

@media screen and (max-width: 375px) {
  .product-name {
    font-size: 15px;
    margin-bottom: 12px;
  }
}

.product-info {
  display: flex;
  justify-content: space-between;
  margin-bottom: 15px;
}

@media screen and (max-width: 375px) {
  .product-info {
    margin-bottom: 12px;
  }
}

.return-info, .period-info {
  display: flex;
  flex-direction: column;
}

.return-rate {
  font-size: 18px;
  font-weight: bold;
  color: #f44336;
}

@media screen and (max-width: 375px) {
  .return-rate {
    font-size: 16px;
  }
}

.period {
  font-size: 16px;
  font-weight: 500;
}

@media screen and (max-width: 375px) {
  .period {
    font-size: 14px;
  }
}

.return-label, .period-label {
  font-size: 12px;
  color: #999;
  margin-top: 5px;
}

.risk-level {
  margin-bottom: 15px;
  font-size: 12px;
  padding: 3px 0;
}

.risk-低 {
  color: #4CAF50;
}

.risk-中 {
  color: #ff9800;
}

.risk-高 {
  color: #f44336;
}

.product-button {
  background-color: #4CAF50;
  color: white;
  text-align: center;
  padding: 8px 0;
  border-radius: 4px;
  font-size: 14px;
  font-weight: 500;
  cursor: pointer;
  transition: background-color 0.2s ease;
}

.product-button:active {
  background-color: #3e8e41;
}
</style>
