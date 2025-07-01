<template>
  <div class="system-view">
    <div class="page-header">
      <h2 class="page-title">系统设置</h2>
      <div class="header-actions">
        <el-button @click="refreshData" :loading="loading">
          <el-icon><Refresh /></el-icon>
          刷新
        </el-button>
      </div>
    </div>

    <el-row :gutter="20">
      <!-- 左侧配置列表 -->
      <el-col :span="16">
        <!-- AI审批设置 -->
        <el-card class="config-card" shadow="never">
          <template #header>
            <div class="card-header">
              <span>
                <el-icon><Cpu /></el-icon>
                AI审批设置
              </span>
            </div>
          </template>
          
          <div class="config-section">
            <div class="config-item">
              <div class="config-label">
                <h4>AI审批功能</h4>
                <p>启用或禁用AI自动审批功能</p>
              </div>
              <div class="config-control">
                <el-switch
                  v-model="aiApprovalEnabled"
                  @change="handleAIApprovalToggle"
                  :loading="toggleLoading"
                  size="large"
                />
              </div>
            </div>
            
            <div class="config-item">
              <div class="config-label">
                <h4>风险阈值设置</h4>
                <p>设置AI审批的风险评分阈值</p>
              </div>
              <div class="config-control">
                <el-input-number
                  v-model="riskThreshold"
                  :min="0"
                  :max="100"
                  :step="5"
                  @change="updateRiskThreshold"
                />
                <span class="unit">分</span>
              </div>
            </div>
            
            <div class="config-item">
              <div class="config-label">
                <h4>自动批准金额上限</h4>
                <p>AI可以自动批准的最大金额</p>
              </div>
              <div class="config-control">
                <el-input-number
                  v-model="autoApprovalLimit"
                  :min="1000"
                  :max="100000"
                  :step="1000"
                  @change="updateAutoApprovalLimit"
                />
                <span class="unit">元</span>
              </div>
            </div>
          </div>
        </el-card>

        <!-- 系统配置 -->
        <el-card class="config-card" shadow="never">
          <template #header>
            <div class="card-header">
              <span>
                <el-icon><Setting /></el-icon>
                系统配置
              </span>
              <el-button type="primary" size="small" @click="showAddConfigDialog">
                <el-icon><Plus /></el-icon>
                添加配置
              </el-button>
            </div>
          </template>
          
          <el-table
            v-loading="loading"
            :data="configurations"
            stripe
            style="width: 100%"
          >
            <el-table-column prop="config_key" label="配置键" width="200">
              <template #default="{ row }">
                <code class="config-key">{{ row.config_key }}</code>
              </template>
            </el-table-column>
            
            <el-table-column prop="config_value" label="配置值" min-width="200">
              <template #default="{ row }">
                <el-input
                  v-if="row.editing"
                  v-model="row.config_value"
                  size="small"
                  @blur="saveConfig(row)"
                  @keyup.enter="saveConfig(row)"
                />
                <span v-else class="config-value">{{ row.config_value }}</span>
              </template>
            </el-table-column>
            
            <el-table-column prop="description" label="描述" min-width="250">
              <template #default="{ row }">
                <span class="config-desc">{{ row.description }}</span>
              </template>
            </el-table-column>
            
            <el-table-column prop="updated_at" label="更新时间" width="160">
              <template #default="{ row }">
                {{ formatDateTime(row.updated_at) }}
              </template>
            </el-table-column>
            
            <el-table-column label="操作" width="120" fixed="right">
              <template #default="{ row }">
                <el-button
                  v-if="!row.editing"
                  type="primary"
                  size="small"
                  @click="editConfig(row)"
                >
                  编辑
                </el-button>
                <el-button
                  v-else
                  type="success"
                  size="small"
                  @click="saveConfig(row)"
                >
                  保存
                </el-button>
              </template>
            </el-table-column>
          </el-table>
        </el-card>
      </el-col>

      <!-- 右侧系统信息 -->
      <el-col :span="8">
        <!-- 系统状态 -->
        <el-card class="status-card" shadow="never">
          <template #header>
            <div class="card-header">
              <span>
                <el-icon><Monitor /></el-icon>
                系统状态
              </span>
            </div>
          </template>
          
          <div class="status-list">
            <div class="status-item">
              <div class="status-label">AI审批状态</div>
              <el-tag :type="aiApprovalEnabled ? 'success' : 'danger'" size="small">
                {{ aiApprovalEnabled ? '运行中' : '已停止' }}
              </el-tag>
            </div>
            
            <div class="status-item">
              <div class="status-label">系统运行时间</div>
              <span class="status-value">{{ systemUptime }}</span>
            </div>
            
            <div class="status-item">
              <div class="status-label">数据库连接</div>
              <el-tag type="success" size="small">正常</el-tag>
            </div>
            
            <div class="status-item">
              <div class="status-label">API响应时间</div>
              <span class="status-value">{{ apiResponseTime }}ms</span>
            </div>
          </div>
        </el-card>

        <!-- 系统统计 -->
        <el-card class="stats-card" shadow="never">
          <template #header>
            <div class="card-header">
              <span>
                <el-icon><DataAnalysis /></el-icon>
                系统统计
              </span>
            </div>
          </template>
          
          <div class="stats-list">
            <div class="stats-item">
              <div class="stats-label">总用户数</div>
              <div class="stats-value">{{ systemStats.totalUsers }}</div>
            </div>
            
            <div class="stats-item">
              <div class="stats-label">活跃用户</div>
              <div class="stats-value">{{ systemStats.activeUsers }}</div>
            </div>
            
            <div class="stats-item">
              <div class="stats-label">今日申请</div>
              <div class="stats-value">{{ systemStats.todayApplications }}</div>
            </div>
            
            <div class="stats-item">
              <div class="stats-label">AI处理率</div>
              <div class="stats-value">{{ systemStats.aiProcessRate }}%</div>
            </div>
          </div>
        </el-card>

        <!-- 快速操作 -->
        <el-card class="actions-card" shadow="never">
          <template #header>
            <div class="card-header">
              <span>
                <el-icon><Operation /></el-icon>
                快速操作
              </span>
            </div>
          </template>
          
          <div class="action-buttons">
            <el-button type="primary" @click="clearCache" :loading="clearingCache">
              <el-icon><Delete /></el-icon>
              清除缓存
            </el-button>
            
            <el-button type="warning" @click="backupDatabase" :loading="backingUp">
              <el-icon><Download /></el-icon>
              备份数据
            </el-button>
            
            <el-button type="info" @click="exportLogs">
              <el-icon><Document /></el-icon>
              导出日志
            </el-button>
            
            <el-button type="danger" @click="restartService">
              <el-icon><RefreshRight /></el-icon>
              重启服务
            </el-button>
          </div>
        </el-card>
      </el-col>
    </el-row>

    <!-- 添加配置对话框 -->
    <el-dialog
      v-model="configDialogVisible"
      title="添加系统配置"
      width="500px"
      @close="resetConfigForm"
    >
      <el-form
        ref="configFormRef"
        :model="configForm"
        :rules="configRules"
        label-width="100px"
      >
        <el-form-item label="配置键" prop="config_key">
          <el-input
            v-model="configForm.config_key"
            placeholder="请输入配置键"
          />
        </el-form-item>
        
        <el-form-item label="配置值" prop="config_value">
          <el-input
            v-model="configForm.config_value"
            placeholder="请输入配置值"
          />
        </el-form-item>
        
        <el-form-item label="描述" prop="description">
          <el-input
            v-model="configForm.description"
            type="textarea"
            :rows="3"
            placeholder="请输入配置描述"
          />
        </el-form-item>
      </el-form>
      
      <template #footer>
        <el-button @click="configDialogVisible = false">取消</el-button>
        <el-button
          type="primary"
          @click="submitConfigForm"
          :loading="submitting"
        >
          创建
        </el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { ElMessage, ElMessageBox, type FormInstance, type FormRules } from 'element-plus'
import {
  Refresh,
  Cpu,
  Setting,
  Plus,
  Monitor,
  DataAnalysis,
  Operation,
  Delete,
  Download,
  Document,
  RefreshRight
} from '@element-plus/icons-vue'
import { 
  getSystemConfigurations, 
  updateSystemConfiguration, 
  toggleAIApproval 
} from '@/api/admin'
import type { SystemConfig } from '@/types'
import dayjs from 'dayjs'

const loading = ref(false)
const toggleLoading = ref(false)
const submitting = ref(false)
const clearingCache = ref(false)
const backingUp = ref(false)
const configDialogVisible = ref(false)
const configFormRef = ref<FormInstance>()

// 系统配置
const aiApprovalEnabled = ref(true)
const riskThreshold = ref(70)
const autoApprovalLimit = ref(50000)
const configurations = ref<(SystemConfig & { editing?: boolean })[]>([])

// 系统状态
const systemUptime = ref('7天 12小时 30分钟')
const apiResponseTime = ref(85)

// 系统统计
const systemStats = reactive({
  totalUsers: 125,
  activeUsers: 89,
  todayApplications: 23,
  aiProcessRate: 85
})

// 配置表单
const configForm = reactive({
  config_key: '',
  config_value: '',
  description: ''
})

const configRules: FormRules = {
  config_key: [
    { required: true, message: '请输入配置键', trigger: 'blur' }
  ],
  config_value: [
    { required: true, message: '请输入配置值', trigger: 'blur' }
  ],
  description: [
    { required: true, message: '请输入描述', trigger: 'blur' }
  ]
}

// 方法
const fetchConfigurations = async () => {
  try {
    loading.value = true
    const data = await getSystemConfigurations()
    configurations.value = data.map((config: SystemConfig) => ({
      ...config,
      editing: false
    }))
  } catch (error) {
    ElMessage.error('获取系统配置失败')
  } finally {
    loading.value = false
  }
}

const refreshData = () => {
  fetchConfigurations()
}

const handleAIApprovalToggle = async (enabled: boolean) => {
  try {
    toggleLoading.value = true
    await toggleAIApproval(enabled)
    ElMessage.success(`AI审批功能已${enabled ? '开启' : '关闭'}`)
  } catch (error) {
    ElMessage.error('切换AI审批状态失败')
    // 回滚状态
    aiApprovalEnabled.value = !enabled
  } finally {
    toggleLoading.value = false
  }
}

const updateRiskThreshold = async (value: number) => {
  try {
    await updateSystemConfiguration('ai_risk_threshold', value.toString())
    ElMessage.success('风险阈值更新成功')
  } catch (error) {
    ElMessage.error('更新失败')
  }
}

const updateAutoApprovalLimit = async (value: number) => {
  try {
    await updateSystemConfiguration('auto_approval_limit', value.toString())
    ElMessage.success('自动批准金额上限更新成功')
  } catch (error) {
    ElMessage.error('更新失败')
  }
}

const editConfig = (config: SystemConfig & { editing?: boolean }) => {
  config.editing = true
}

const saveConfig = async (config: SystemConfig & { editing?: boolean }) => {
  try {
    await updateSystemConfiguration(config.config_key, config.config_value)
    config.editing = false
    ElMessage.success('配置更新成功')
    fetchConfigurations()
  } catch (error) {
    ElMessage.error('更新配置失败')
  }
}

const showAddConfigDialog = () => {
  configDialogVisible.value = true
}

const resetConfigForm = () => {
  Object.assign(configForm, {
    config_key: '',
    config_value: '',
    description: ''
  })
}

const submitConfigForm = async () => {
  if (!configFormRef.value) return
  
  try {
    await configFormRef.value.validate()
    submitting.value = true
    
    // 这里应该调用创建配置的API，暂时使用更新配置的API
    await updateSystemConfiguration(configForm.config_key, configForm.config_value)
    
    ElMessage.success('配置创建成功')
    configDialogVisible.value = false
    resetConfigForm()
    fetchConfigurations()
  } catch (error) {
    ElMessage.error('创建配置失败')
  } finally {
    submitting.value = false
  }
}

const clearCache = async () => {
  try {
    await ElMessageBox.confirm('确定要清除系统缓存吗？', '提示', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning'
    })
    
    clearingCache.value = true
    // 模拟清除缓存
    await new Promise(resolve => setTimeout(resolve, 2000))
    ElMessage.success('缓存清除成功')
  } catch (error: any) {
    if (error === 'cancel') return
    ElMessage.error('清除缓存失败')
  } finally {
    clearingCache.value = false
  }
}

const backupDatabase = async () => {
  try {
    await ElMessageBox.confirm('确定要备份数据库吗？', '提示', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'info'
    })
    
    backingUp.value = true
    // 模拟备份数据库
    await new Promise(resolve => setTimeout(resolve, 3000))
    ElMessage.success('数据库备份成功')
  } catch (error: any) {
    if (error === 'cancel') return
    ElMessage.error('备份数据库失败')
  } finally {
    backingUp.value = false
  }
}

const exportLogs = () => {
  ElMessage.info('导出日志功能开发中...')
}

const restartService = async () => {
  try {
    await ElMessageBox.confirm('确定要重启服务吗？这将导致系统短暂不可用。', '警告', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning'
    })
    
    ElMessage.info('服务重启功能开发中...')
  } catch (error: any) {
    if (error === 'cancel') return
  }
}

const formatDateTime = (datetime: string) => {
  return dayjs(datetime).format('YYYY-MM-DD HH:mm')
}

onMounted(() => {
  fetchConfigurations()
})
</script>

<style scoped>
.system-view {
  padding: 0;
}

.page-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 20px;
}

.page-title {
  margin: 0;
  color: #333;
  font-size: 24px;
  font-weight: 600;
}

.config-card,
.status-card,
.stats-card,
.actions-card {
  margin-bottom: 20px;
  border-radius: 8px;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  font-weight: 600;
  gap: 8px;
}

.config-section {
  display: flex;
  flex-direction: column;
  gap: 24px;
}

.config-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 20px;
  background: #f8f9fa;
  border-radius: 8px;
}

.config-label h4 {
  margin: 0 0 4px 0;
  font-size: 16px;
  font-weight: 600;
  color: #333;
}

.config-label p {
  margin: 0;
  font-size: 13px;
  color: #666;
}

.config-control {
  display: flex;
  align-items: center;
  gap: 8px;
}

.unit {
  color: #666;
  font-size: 14px;
}

.config-key {
  background: #f1f2f6;
  padding: 2px 6px;
  border-radius: 4px;
  font-size: 12px;
  color: #333;
}

.config-value {
  word-break: break-all;
}

.config-desc {
  color: #666;
  font-size: 13px;
}

.status-list,
.stats-list {
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.status-item,
.stats-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 12px 0;
  border-bottom: 1px solid #f0f2f5;
}

.status-item:last-child,
.stats-item:last-child {
  border-bottom: none;
}

.status-label,
.stats-label {
  color: #666;
  font-size: 14px;
}

.status-value,
.stats-value {
  font-weight: 600;
  color: #333;
}

.action-buttons {
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  gap: 12px;
}

:deep(.el-form-item) {
  margin-bottom: 16px;
}
</style> 