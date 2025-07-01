import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import { ElMessage } from 'element-plus'

export const useAIStatusStore = defineStore('aiStatus', () => {
  // 全局AI审批开关
  const globalAIStatus = ref(false)
  
  // AI工作状态
  const aiWorkStatus = ref({
    totalTasks: 0,
    processingTasks: 0,
    completedTasks: 0,
    failedTasks: 0,
    averageProcessingTime: 0,
    successRate: 0
  })
  
  // 各模块AI状态
  const moduleStatus = ref({
    loan: {
      enabled: false,
      confidenceThreshold: 0.8,
      totalProcessed: 0,
      successRate: 0
    },
    lease: {
      enabled: false,
      confidenceThreshold: 0.8,
      totalProcessed: 0,
      successRate: 0
    },
    insurance: {
      enabled: false,
      confidenceThreshold: 0.8,
      totalProcessed: 0,
      successRate: 0
    }
  })
  
  // 计算属性
  const isAIEnabled = computed(() => globalAIStatus.value)
  
  const totalProcessedTasks = computed(() => 
    Object.values(moduleStatus.value).reduce((sum, module) => sum + module.totalProcessed, 0)
  )
  
  const averageSuccessRate = computed(() => {
    const modules = Object.values(moduleStatus.value)
    const enabledModules = modules.filter(module => module.enabled)
    if (enabledModules.length === 0) return 0
    
    const totalRate = enabledModules.reduce((sum, module) => sum + module.successRate, 0)
    return Math.round(totalRate / enabledModules.length)
  })
  
  // 动作
  const toggleAIStatus = async (silent = false) => {
    try {
      globalAIStatus.value = !globalAIStatus.value
      
      // 更新各模块状态
      Object.keys(moduleStatus.value).forEach(key => {
        moduleStatus.value[key as keyof typeof moduleStatus.value].enabled = globalAIStatus.value
      })
      
      if (!silent) {
        ElMessage.success(
          globalAIStatus.value ? '全局AI审批已启动' : '全局AI审批已暂停'
        )
      }
      
      // 这里可以调用API更新后端状态
      // await updateGlobalAIStatus(globalAIStatus.value)
      
    } catch (error) {
      ElMessage.error('AI状态切换失败')
      // 回滚状态
      globalAIStatus.value = !globalAIStatus.value
    }
  }
  
  const updateModuleStatus = (moduleKey: string, config: any) => {
    if (moduleStatus.value[moduleKey as keyof typeof moduleStatus.value]) {
      Object.assign(moduleStatus.value[moduleKey as keyof typeof moduleStatus.value], config)
    }
  }
  
  const updateWorkStatus = (status: any) => {
    Object.assign(aiWorkStatus.value, status)
  }
  
  const resetStatistics = () => {
    aiWorkStatus.value = {
      totalTasks: 0,
      processingTasks: 0,
      completedTasks: 0,
      failedTasks: 0,
      averageProcessingTime: 0,
      successRate: 0
    }
    
    Object.keys(moduleStatus.value).forEach(key => {
      moduleStatus.value[key as keyof typeof moduleStatus.value].totalProcessed = 0
      moduleStatus.value[key as keyof typeof moduleStatus.value].successRate = 0
    })
  }
  
  const loadAIStatus = async () => {
    try {
      // 这里调用API获取当前AI状态
      // const response = await getAIStatus()
      // globalAIStatus.value = response.globalEnabled
      // moduleStatus.value = response.modules
      // aiWorkStatus.value = response.workStatus
      
      // 模拟数据
      globalAIStatus.value = false
      updateWorkStatus({
        totalTasks: 1250,
        processingTasks: 23,
        completedTasks: 1180,
        failedTasks: 47,
        averageProcessingTime: 45,
        successRate: 94.2
      })
      
    } catch (error) {
      console.error('加载AI状态失败:', error)
    }
  }
  
  return {
    // 状态
    globalAIStatus,
    aiWorkStatus,
    moduleStatus,
    
    // 计算属性
    isAIEnabled,
    totalProcessedTasks,
    averageSuccessRate,
    
    // 方法
    toggleAIStatus,
    updateModuleStatus,
    updateWorkStatus,
    resetStatistics,
    loadAIStatus
  }
}) 