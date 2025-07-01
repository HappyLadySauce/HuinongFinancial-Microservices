<template>
  <el-card class="chart-card" shadow="never">
    <template #header>
      <div class="chart-header">
        <div class="header-left">
          <el-icon v-if="icon" class="header-icon">
            <component :is="icon" />
          </el-icon>
          <span class="header-title">{{ title }}</span>
          <el-tag v-if="subtitle" type="info" size="small" class="header-subtitle">
            {{ subtitle }}
          </el-tag>
        </div>
        
        <div class="header-actions">
          <slot name="header-actions" />
          
          <!-- 时间范围选择 -->
          <el-select 
            v-if="showTimeRange"
            v-model="selectedTimeRange"
            size="small"
            style="width: 120px"
            @change="handleTimeRangeChange"
          >
            <el-option
              v-for="option in timeRangeOptions"
              :key="option.value"
              :label="option.label"
              :value="option.value"
            />
          </el-select>
          
          <!-- 刷新按钮 -->
          <el-tooltip v-if="showRefresh" content="刷新数据" placement="top">
            <el-button 
              size="small" 
              circle 
              @click="handleRefresh"
              :loading="refreshLoading"
            >
              <el-icon><Refresh /></el-icon>
            </el-button>
          </el-tooltip>
          
          <!-- 全屏按钮 -->
          <el-tooltip v-if="showFullscreen" content="全屏显示" placement="top">
            <el-button size="small" circle @click="toggleFullscreen">
              <el-icon><FullScreen /></el-icon>
            </el-button>
          </el-tooltip>
        </div>
      </div>
    </template>

    <!-- 图表内容 -->
    <div 
      ref="chartContainer" 
      class="chart-container"
      :style="{ height: chartHeight }"
      v-loading="loading"
      element-loading-text="加载中..."
    >
      <!-- 无数据状态 -->
      <div v-if="!loading && !hasData" class="no-data">
        <el-empty 
          :description="noDataText"
          :image-size="100"
        >
          <el-button type="primary" @click="handleRefresh">
            重新加载
          </el-button>
        </el-empty>
      </div>
    </div>

    <!-- 图表说明 -->
    <div v-if="showLegend && legendData.length > 0" class="chart-legend">
      <div 
        v-for="item in legendData"
        :key="item.name"
        class="legend-item"
        :style="{ color: item.color }"
      >
        <span class="legend-color" :style="{ backgroundColor: item.color }"></span>
        <span class="legend-name">{{ item.name }}</span>
        <span v-if="item.value !== undefined" class="legend-value">{{ item.value }}</span>
      </div>
    </div>

    <!-- 全屏对话框 -->
    <el-dialog
      v-model="fullscreenVisible"
      :title="title"
      width="90%"
      fullscreen
      @close="fullscreenVisible = false"
    >
      <div 
        ref="fullscreenChartContainer"
        class="fullscreen-chart"
        v-loading="loading"
      ></div>
    </el-dialog>
  </el-card>
</template>

<script setup lang="ts">
import { ref, onMounted, onUnmounted, watch, computed, nextTick } from 'vue'
import { Refresh, FullScreen } from '@element-plus/icons-vue'
import * as echarts from 'echarts'

interface TimeRangeOption {
  label: string
  value: string
}

interface LegendItem {
  name: string
  color: string
  value?: string | number
}

interface Props {
  title: string
  subtitle?: string
  icon?: any
  chartHeight?: string
  loading?: boolean
  showRefresh?: boolean
  showFullscreen?: boolean
  showTimeRange?: boolean
  showLegend?: boolean
  timeRangeOptions?: TimeRangeOption[]
  defaultTimeRange?: string
  noDataText?: string
  chartOptions?: any
  legendData?: LegendItem[]
}

const props = withDefaults(defineProps<Props>(), {
  chartHeight: '300px',
  loading: false,
  showRefresh: true,
  showFullscreen: true,
  showTimeRange: false,
  showLegend: false,
  timeRangeOptions: () => [
    { label: '今日', value: 'today' },
    { label: '本周', value: 'week' },
    { label: '本月', value: 'month' },
    { label: '本年', value: 'year' }
  ],
  defaultTimeRange: 'month',
  noDataText: '暂无数据',
  chartOptions: () => ({}),
  legendData: () => []
})

const emit = defineEmits<{
  refresh: []
  timeRangeChange: [range: string]
}>()

const chartContainer = ref<HTMLElement>()
const fullscreenChartContainer = ref<HTMLElement>()
const selectedTimeRange = ref(props.defaultTimeRange)
const refreshLoading = ref(false)
const fullscreenVisible = ref(false)

let chartInstance: echarts.ECharts | null = null
let fullscreenChartInstance: echarts.ECharts | null = null

const hasData = computed(() => {
  if (!props.chartOptions || !props.chartOptions.series) return false
  return props.chartOptions.series.some((series: any) => 
    series.data && series.data.length > 0
  )
})

// 初始化图表
const initChart = () => {
  if (!chartContainer.value) return

  chartInstance = echarts.init(chartContainer.value)
  updateChart()
  
  // 监听窗口大小变化
  const resizeObserver = new ResizeObserver(() => {
    chartInstance?.resize()
  })
  resizeObserver.observe(chartContainer.value)
}

// 更新图表
const updateChart = () => {
  if (!chartInstance || !hasData.value) return

  const defaultOptions = {
    tooltip: {
      trigger: 'axis',
      backgroundColor: 'rgba(50, 50, 50, 0.9)',
      borderColor: 'transparent',
      textStyle: {
        color: '#fff'
      }
    },
    grid: {
      left: '3%',
      right: '4%',
      bottom: '3%',
      containLabel: true
    },
    responsive: true
  }

  const options = {
    ...defaultOptions,
    ...props.chartOptions
  }

  chartInstance.setOption(options, true)
}

// 初始化全屏图表
const initFullscreenChart = async () => {
  await nextTick()
  if (!fullscreenChartContainer.value) return

  fullscreenChartInstance = echarts.init(fullscreenChartContainer.value)
  fullscreenChartInstance.setOption(props.chartOptions, true)
}

// 处理刷新
const handleRefresh = async () => {
  refreshLoading.value = true
  try {
    emit('refresh')
  } finally {
    setTimeout(() => {
      refreshLoading.value = false
    }, 500)
  }
}

// 处理时间范围变化
const handleTimeRangeChange = (range: string) => {
  emit('timeRangeChange', range)
}

// 切换全屏
const toggleFullscreen = async () => {
  fullscreenVisible.value = true
  await initFullscreenChart()
}

// 监听图表选项变化
watch(() => props.chartOptions, () => {
  updateChart()
  if (fullscreenChartInstance) {
    fullscreenChartInstance.setOption(props.chartOptions, true)
  }
}, { deep: true })

// 监听全屏对话框关闭
watch(fullscreenVisible, (visible) => {
  if (!visible && fullscreenChartInstance) {
    fullscreenChartInstance.dispose()
    fullscreenChartInstance = null
  }
})

onMounted(() => {
  initChart()
})

onUnmounted(() => {
  if (chartInstance) {
    chartInstance.dispose()
  }
  if (fullscreenChartInstance) {
    fullscreenChartInstance.dispose()
  }
})

// 暴露给父组件的方法
const resize = () => {
  chartInstance?.resize()
  fullscreenChartInstance?.resize()
}

const getChartInstance = () => chartInstance

defineExpose({
  resize,
  getChartInstance
})
</script>

<style scoped>
.chart-card {
  border: none;
  border-radius: 12px;
}

.chart-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.header-left {
  display: flex;
  align-items: center;
  gap: 8px;
}

.header-icon {
  font-size: 18px;
  color: #409eff;
}

.header-title {
  font-weight: 500;
  font-size: 16px;
}

.header-subtitle {
  margin-left: 8px;
}

.header-actions {
  display: flex;
  align-items: center;
  gap: 8px;
}

.chart-container {
  width: 100%;
  position: relative;
}

.no-data {
  display: flex;
  align-items: center;
  justify-content: center;
  height: 100%;
  min-height: 200px;
}

.chart-legend {
  display: flex;
  flex-wrap: wrap;
  gap: 16px;
  margin-top: 16px;
  padding-top: 16px;
  border-top: 1px solid #f0f0f0;
}

.legend-item {
  display: flex;
  align-items: center;
  gap: 6px;
  font-size: 12px;
}

.legend-color {
  width: 12px;
  height: 12px;
  border-radius: 2px;
}

.legend-name {
  color: #666;
}

.legend-value {
  font-weight: 500;
  color: #333;
}

.fullscreen-chart {
  width: 100%;
  height: 70vh;
}

:deep(.el-card__header) {
  padding: 16px 20px;
  border-bottom: 1px solid #f0f0f0;
}

:deep(.el-card__body) {
  padding: 20px;
}
</style> 