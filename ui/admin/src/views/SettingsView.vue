<template>
  <div class="settings-view">
    <!-- 页面头部 -->
    <PageHeader
      title="系统设置"
      subtitle="个人偏好和系统配置"
      :icon="Setting"
      :breadcrumbs="breadcrumbs"
    />

    <el-row :gutter="24">
      <!-- 左侧设置项 -->
      <el-col :span="16">
        <!-- 个人偏好 -->
        <el-card class="settings-card" shadow="never">
          <template #header>
            <div class="card-header">
              <span>
                <el-icon><User /></el-icon>
                个人偏好设置
              </span>
            </div>
          </template>

          <el-form :model="personalSettings" label-width="120px">
            <el-form-item label="界面主题">
              <el-radio-group v-model="personalSettings.theme">
                <el-radio label="light">浅色主题</el-radio>
                <el-radio label="dark">深色主题</el-radio>
                <el-radio label="auto">跟随系统</el-radio>
              </el-radio-group>
            </el-form-item>

            <el-form-item label="界面语言">
              <el-select v-model="personalSettings.language" style="width: 200px">
                <el-option label="简体中文" value="zh-CN" />
                <el-option label="English" value="en-US" />
              </el-select>
            </el-form-item>

            <el-form-item label="侧边栏折叠">
              <el-switch v-model="personalSettings.sidebarCollapsed" />
            </el-form-item>

            <el-form-item label="面包屑导航">
              <el-switch v-model="personalSettings.showBreadcrumb" />
            </el-form-item>

            <el-form-item label="页面动画">
              <el-switch v-model="personalSettings.enableAnimation" />
            </el-form-item>
          </el-form>
        </el-card>

        <!-- 通知设置 -->
        <el-card class="settings-card" shadow="never">
          <template #header>
            <div class="card-header">
              <span>
                <el-icon><Bell /></el-icon>
                通知设置
              </span>
            </div>
          </template>

          <el-form :model="notificationSettings" label-width="120px">
            <el-form-item label="邮件通知">
              <el-switch v-model="notificationSettings.emailNotification" />
              <span class="setting-desc">接收重要通知的邮件提醒</span>
            </el-form-item>

            <el-form-item label="短信通知">
              <el-switch v-model="notificationSettings.smsNotification" />
              <span class="setting-desc">接收审批相关的短信通知</span>
            </el-form-item>

            <el-form-item label="浏览器通知">
              <el-switch v-model="notificationSettings.browserNotification" />
              <span class="setting-desc">浏览器桌面通知提醒</span>
            </el-form-item>

            <el-form-item label="声音提醒">
              <el-switch v-model="notificationSettings.soundNotification" />
              <span class="setting-desc">新消息声音提醒</span>
            </el-form-item>

            <el-form-item label="工作时间">
              <el-time-picker
                v-model="notificationSettings.workingHours"
                is-range
                range-separator="至"
                start-placeholder="开始时间"
                end-placeholder="结束时间"
                format="HH:mm"
                value-format="HH:mm"
                style="width: 200px"
              />
              <span class="setting-desc">仅在工作时间内发送通知</span>
            </el-form-item>
          </el-form>
        </el-card>

        <!-- 安全设置 -->
        <el-card class="settings-card" shadow="never">
          <template #header>
            <div class="card-header">
              <span>
                <el-icon><Lock /></el-icon>
                安全设置
              </span>
            </div>
          </template>

          <el-form :model="securitySettings" label-width="120px">
            <el-form-item label="自动锁屏">
              <el-select v-model="securitySettings.autoLockTime" style="width: 200px">
                <el-option label="不自动锁屏" :value="0" />
                <el-option label="5分钟" :value="5" />
                <el-option label="10分钟" :value="10" />
                <el-option label="30分钟" :value="30" />
                <el-option label="1小时" :value="60" />
              </el-select>
              <span class="setting-desc">无操作自动锁屏时间</span>
            </el-form-item>

            <el-form-item label="双因子认证">
              <el-switch v-model="securitySettings.twoFactorAuth" />
              <span class="setting-desc">启用双因子身份验证</span>
            </el-form-item>

            <el-form-item label="登录日志">
              <el-switch v-model="securitySettings.loginLog" />
              <span class="setting-desc">记录登录活动日志</span>
            </el-form-item>

            <el-form-item label="会话管理">
              <el-button @click="showSessionDialog = true">管理活跃会话</el-button>
              <span class="setting-desc">查看和管理当前登录会话</span>
            </el-form-item>
          </el-form>
        </el-card>
      </el-col>

      <!-- 右侧操作 -->
      <el-col :span="8">
        <!-- 快速操作 -->
        <el-card class="actions-card" shadow="never">
          <template #header>
            <div class="card-header">
              <span>
                <el-icon><Tools /></el-icon>
                快速操作
              </span>
            </div>
          </template>

          <div class="action-buttons">
            <el-button @click="handleSaveSettings" type="primary" style="width: 100%; margin-bottom: 12px;">
              <el-icon><Check /></el-icon>
              保存设置
            </el-button>
            <el-button @click="handleResetSettings" style="width: 100%; margin-bottom: 12px;">
              <el-icon><RefreshLeft /></el-icon>
              重置为默认
            </el-button>
            <el-button @click="handleExportSettings" style="width: 100%; margin-bottom: 12px;">
              <el-icon><Download /></el-icon>
              导出设置
            </el-button>
            <el-button @click="handleImportSettings" style="width: 100%;">
              <el-icon><Upload /></el-icon>
              导入设置
            </el-button>
          </div>
        </el-card>

        <!-- 系统信息 -->
        <el-card class="info-card" shadow="never">
          <template #header>
            <div class="card-header">
              <span>
                <el-icon><InfoFilled /></el-icon>
                系统信息
              </span>
            </div>
          </template>

          <div class="system-info">
            <div class="info-item">
              <span class="label">系统版本：</span>
              <span class="value">v1.0.0</span>
            </div>
            <div class="info-item">
              <span class="label">构建时间：</span>
              <span class="value">2024-12-29</span>
            </div>
            <div class="info-item">
              <span class="label">浏览器：</span>
              <span class="value">{{ getBrowserInfo() }}</span>
            </div>
            <div class="info-item">
              <span class="label">屏幕分辨率：</span>
              <span class="value">{{ getScreenResolution() }}</span>
            </div>
          </div>
        </el-card>
      </el-col>
    </el-row>

    <!-- 会话管理对话框 -->
    <el-dialog v-model="showSessionDialog" title="会话管理" width="800px">
      <el-table :data="activeSessions" border>
        <el-table-column prop="id" label="会话ID" width="120" />
        <el-table-column prop="device" label="设备信息" />
        <el-table-column prop="location" label="登录地点" />
        <el-table-column prop="loginTime" label="登录时间" />
        <el-table-column prop="lastActivity" label="最后活动" />
        <el-table-column label="状态" width="80">
          <template #default="{ row }">
            <StatusTag :status="row.status" />
          </template>
        </el-table-column>
        <el-table-column label="操作" width="100">
          <template #default="{ row }">
            <el-button 
              v-if="row.id !== currentSessionId" 
              size="small" 
              type="danger" 
              @click="handleKillSession(row.id)"
            >
              终止
            </el-button>
            <el-tag v-else size="small" type="success">当前</el-tag>
          </template>
        </el-table-column>
      </el-table>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, computed } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import {
  Setting,
  User,
  Bell,
  Lock,
  Tools,
  Check,
  RefreshLeft,
  Download,
  Upload,
  InfoFilled
} from '@element-plus/icons-vue'

// 导入通用组件
import { PageHeader, StatusTag } from '@/components/common'

const showSessionDialog = ref(false)
const currentSessionId = ref('session-123')

const personalSettings = reactive({
  theme: 'light',
  language: 'zh-CN',
  sidebarCollapsed: false,
  showBreadcrumb: true,
  enableAnimation: true
})

const notificationSettings = reactive({
  emailNotification: true,
  smsNotification: false,
  browserNotification: true,
  soundNotification: false,
  workingHours: ['09:00', '18:00']
})

const securitySettings = reactive({
  autoLockTime: 30,
  twoFactorAuth: false,
  loginLog: true
})

const activeSessions = ref([
  {
    id: 'session-123',
    device: 'Chrome 120.0 / Windows 10',
    location: '江苏省南京市',
    loginTime: '2024-12-29 10:30:00',
    lastActivity: '2024-12-29 18:15:00',
    status: 'active'
  },
  {
    id: 'session-124',
    device: 'Safari 17.0 / macOS',
    location: '江苏省苏州市',
    loginTime: '2024-12-28 14:20:00',
    lastActivity: '2024-12-28 17:45:00',
    status: 'inactive'
  }
])

// 计算属性
const breadcrumbs = computed(() => [
  { title: '首页', to: '/' },
  { title: '系统设置' }
])

// 方法
const handleSaveSettings = async () => {
  try {
    console.log('保存设置:', {
      personal: personalSettings,
      notification: notificationSettings,
      security: securitySettings
    })
    // 实现保存逻辑
    ElMessage.success('设置保存成功')
  } catch (error) {
    console.error('保存失败:', error)
    ElMessage.error('保存失败')
  }
}

const handleResetSettings = async () => {
  try {
    await ElMessageBox.confirm('确定要重置为默认设置吗？', '确认', {
      type: 'warning'
    })
    
    // 重置为默认值
    Object.assign(personalSettings, {
      theme: 'light',
      language: 'zh-CN',
      sidebarCollapsed: false,
      showBreadcrumb: true,
      enableAnimation: true
    })
    
    Object.assign(notificationSettings, {
      emailNotification: true,
      smsNotification: false,
      browserNotification: true,
      soundNotification: false,
      workingHours: ['09:00', '18:00']
    })
    
    Object.assign(securitySettings, {
      autoLockTime: 30,
      twoFactorAuth: false,
      loginLog: true
    })
    
    ElMessage.success('设置已重置为默认值')
  } catch {
    // 用户取消
  }
}

const handleExportSettings = () => {
  const settings = {
    personal: personalSettings,
    notification: notificationSettings,
    security: securitySettings
  }
  
  const blob = new Blob([JSON.stringify(settings, null, 2)], { type: 'application/json' })
  const url = URL.createObjectURL(blob)
  const link = document.createElement('a')
  link.href = url
  link.download = 'settings.json'
  link.click()
  URL.revokeObjectURL(url)
  
  ElMessage.success('设置已导出')
}

const handleImportSettings = () => {
  const input = document.createElement('input')
  input.type = 'file'
  input.accept = '.json'
  input.onchange = (e: any) => {
    const file = e.target.files[0]
    if (file) {
      const reader = new FileReader()
      reader.onload = (e: any) => {
        try {
          const settings = JSON.parse(e.target.result)
          Object.assign(personalSettings, settings.personal || {})
          Object.assign(notificationSettings, settings.notification || {})
          Object.assign(securitySettings, settings.security || {})
          ElMessage.success('设置导入成功')
        } catch (error) {
          ElMessage.error('设置文件格式错误')
        }
      }
      reader.readAsText(file)
    }
  }
  input.click()
}

const handleKillSession = async (sessionId: string) => {
  try {
    await ElMessageBox.confirm('确定要终止该会话吗？', '确认', {
      type: 'warning'
    })
    
    // 实现终止会话逻辑
    const index = activeSessions.value.findIndex(s => s.id === sessionId)
    if (index > -1) {
      activeSessions.value.splice(index, 1)
    }
    
    ElMessage.success('会话已终止')
  } catch {
    // 用户取消
  }
}

const getBrowserInfo = () => {
  const ua = navigator.userAgent
  if (ua.includes('Chrome')) return 'Chrome'
  if (ua.includes('Firefox')) return 'Firefox'
  if (ua.includes('Safari')) return 'Safari'
  if (ua.includes('Edge')) return 'Edge'
  return '未知浏览器'
}

const getScreenResolution = () => {
  return `${screen.width} × ${screen.height}`
}
</script>

<style scoped>
.settings-view {
  height: 100vh;
  display: flex;
  flex-direction: column;
  gap: 16px;
  padding: 0 20px;
  overflow: hidden;
}

.settings-card,
.actions-card,
.info-card {
  margin-bottom: 20px;
  border: none;
  border-radius: 12px;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  font-weight: 500;
}

.setting-desc {
  margin-left: 12px;
  font-size: 12px;
  color: #999;
}

.action-buttons {
  padding: 16px 0;
}

.system-info {
  padding: 16px 0;
}

.info-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 12px;
  padding: 8px 0;
  border-bottom: 1px solid #f0f0f0;
}

.info-item:last-child {
  border-bottom: none;
}

.label {
  color: #666;
  font-size: 14px;
}

.value {
  color: #333;
  font-weight: 500;
}

:deep(.el-form-item__label) {
  font-weight: 500;
}

:deep(.el-card__header) {
  padding: 16px 20px;
  border-bottom: 1px solid #f0f0f0;
}

:deep(.el-card__body) {
  padding: 20px;
}
</style> 