<template>
  <el-card class="stat-card" :class="cardClass" shadow="hover" @click="handleClick">
    <div class="stat-content">
      <div v-if="showIcon" class="stat-icon" :class="iconClass">
        <el-icon>
          <component :is="icon" />
        </el-icon>
      </div>
      <div class="stat-info">
        <div class="stat-number" :class="valueClass">{{ formattedValue }}</div>
        <div class="stat-label">{{ label }}</div>
        <div v-if="change" class="stat-change" :class="changeClass">
          <el-icon>
            <ArrowUp v-if="changeType === 'up'" />
            <ArrowDown v-if="changeType === 'down'" />
            <Minus v-if="changeType === 'same'" />
          </el-icon>
          {{ change }}
        </div>
      </div>
    </div>
  </el-card>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { ArrowUp, ArrowDown, Minus } from '@element-plus/icons-vue'

interface Props {
  value: string | number
  label: string
  icon?: any
  type?: 'primary' | 'success' | 'warning' | 'danger' | 'info'
  change?: string
  changeType?: 'up' | 'down' | 'same'
  showIcon?: boolean
  clickable?: boolean
  format?: 'number' | 'currency' | 'percent'
}

const props = withDefaults(defineProps<Props>(), {
  type: 'primary',
  showIcon: true,
  clickable: false,
  format: 'number'
})

const emit = defineEmits<{
  click: []
}>()

const formattedValue = computed(() => {
  if (typeof props.value === 'string') return props.value
  
  switch (props.format) {
    case 'currency':
      return `¥${props.value.toLocaleString()}`
    case 'percent':
      return `${props.value}%`
    default:
      return props.value.toLocaleString()
  }
})

const cardClass = computed(() => [
  `stat-card-${props.type}`,
  { 'stat-card-clickable': props.clickable }
])

const iconClass = computed(() => `stat-icon-${props.type}`)

const valueClass = computed(() => `stat-number-${props.type}`)

const changeClass = computed(() => props.changeType)

const handleClick = () => {
  if (props.clickable) {
    emit('click')
  }
}
</script>

<style scoped>
.stat-card {
  border: none;
  border-radius: 12px;
  overflow: hidden;
  height: 120px;
  position: relative;
  transition: all 0.3s ease;
}

.stat-card:hover {
  transform: translateY(-2px);
  box-shadow: 0 8px 25px rgba(0, 0, 0, 0.1);
}

.stat-card-clickable {
  cursor: pointer;
}

.stat-content {
  display: flex;
  align-items: center;
  height: 100%;
  padding: 20px;
}

.stat-icon {
  width: 60px;
  height: 60px;
  border-radius: 12px;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 24px;
  margin-right: 16px;
}

.stat-icon-primary {
  background: linear-gradient(135deg, #4facfe 0%, #00f2fe 100%);
  color: white;
}

.stat-icon-success {
  background: linear-gradient(135deg, #43e97b 0%, #38f9d7 100%);
  color: white;
}

.stat-icon-warning {
  background: linear-gradient(135deg, #fa709a 0%, #fee140 100%);
  color: white;
}

.stat-icon-danger {
  background: linear-gradient(135deg, #ff6b6b 0%, #ee5a24 100%);
  color: white;
}

.stat-icon-info {
  background: linear-gradient(135deg, #a8edea 0%, #fed6e3 100%);
  color: #333;
}

.stat-info {
  flex: 1;
}

.stat-number {
  font-size: 32px;
  font-weight: 700;
  line-height: 1;
  margin-bottom: 4px;
}

.stat-number-primary {
  color: #409eff;
}

.stat-number-success {
  color: #67c23a;
}

.stat-number-warning {
  color: #e6a23c;
}

.stat-number-danger {
  color: #f56c6c;
}

.stat-number-info {
  color: #333;
}

.stat-label {
  font-size: 14px;
  color: #666;
  margin-bottom: 8px;
}

.stat-change {
  display: flex;
  align-items: center;
  gap: 4px;
  font-size: 12px;
  font-weight: 500;
}

.stat-change.up {
  color: #67c23a;
}

.stat-change.down {
  color: #f56c6c;
}

.stat-change.same {
  color: #909399;
}
</style> 