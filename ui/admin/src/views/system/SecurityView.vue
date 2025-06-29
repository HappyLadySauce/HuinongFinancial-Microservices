<template>
  <div class="security-view">
    <!-- 页面头部 -->
    <PageHeader
      title="安全设置"
      subtitle="系统安全配置和权限管理"
      :icon="Lock"
      :breadcrumbs="breadcrumbs"
    >
      <template #actions>
        <el-button @click="handleSaveAll" type="primary">
          <el-icon><Check /></el-icon>
          保存所有设置
        </el-button>
      </template>
    </PageHeader>

    <el-row :gutter="24">
      <!-- 左侧设置 -->
      <el-col :span="16">
        <!-- 登录安全 -->
        <el-card class="security-card" shadow="never">
          <template #header>
            <div class="card-header">
              <span>
                <el-icon><Key /></el-icon>
                登录安全设置
              </span>
            </div>
          </template>

          <el-form :model="loginSecurity" label-width="140px">
            <el-form-item label="密码策略">
              <el-checkbox-group v-model="loginSecurity.passwordPolicy">
                <el-checkbox label="minLength">最少8位字符</el-checkbox>
                <el-checkbox label="requireUppercase">包含大写字母</el-checkbox>
                <el-checkbox label="requireLowercase">包含小写字母</el-checkbox>
                <el-checkbox label="requireNumbers">包含数字</el-checkbox>
                <el-checkbox label="requireSpecialChars">包含特殊字符</el-checkbox>
              </el-checkbox-group>
            </el-form-item>

            <el-form-item label="密码有效期">
              <el-select v-model="loginSecurity.passwordExpiry" style="width: 200px">
                <el-option label="永不过期" :value="0" />
                <el-option label="30天" :value="30" />
                <el-option label="60天" :value="60" />
                <el-option label="90天" :value="90" />
                <el-option label="180天" :value="180" />
              </el-select>
              <span class="setting-desc">密码过期后需要重新设置</span>
            </el-form-item>

            <el-form-item label="登录失败限制">
              <el-input-number
                v-model="loginSecurity.maxLoginAttempts"
                :min="3"
                :max="10"
                style="width: 200px"
              />
              <span class="setting-desc">次失败后锁定账户</span>
            </el-form-item>

            <el-form-item label="账户锁定时间">
              <el-select v-model="loginSecurity.lockoutDuration" style="width: 200px">
                <el-option label="15分钟" :value="15" />
                <el-option label="30分钟" :value="30" />
                <el-option label="1小时" :value="60" />
                <el-option label="24小时" :value="1440" />
                <el-option label="永久锁定" :value="-1" />
              </el-select>
            </el-form-item>

            <el-form-item label="强制双因子认证">
              <el-switch v-model="loginSecurity.forceTwoFactor" />
              <span class="setting-desc">要求所有用户启用双因子认证</span>
            </el-form-item>

            <el-form-item label="会话超时">
              <el-select v-model="loginSecurity.sessionTimeout" style="width: 200px">
                <el-option label="30分钟" :value="30" />
                <el-option label="1小时" :value="60" />
                <el-option label="2小时" :value="120" />
                <el-option label="4小时" :value="240" />
                <el-option label="8小时" :value="480" />
              </el-select>
              <span class="setting-desc">无操作自动退出</span>
            </el-form-item>
          </el-form>
        </el-card>

        <!-- 访问控制 -->
        <el-card class="security-card" shadow="never">
          <template #header>
            <div class="card-header">
              <span>
                <el-icon><Lock /></el-icon>
                访问控制
              </span>
            </div>
          </template>

          <el-form :model="accessControl" label-width="140px">
            <el-form-item label="IP白名单">
              <el-switch v-model="accessControl.enableIPWhitelist" />
              <span class="setting-desc">启用IP白名单访问控制</span>
            </el-form-item>

            <el-form-item v-if="accessControl.enableIPWhitelist" label="允许的IP地址">
              <div class="ip-list">
                <div v-for="(ip, index) in accessControl.allowedIPs" :key="index" class="ip-item">
                  <el-input v-model="accessControl.allowedIPs[index]" placeholder="192.168.1.0/24" />
                  <el-button @click="removeIP(index)" type="danger" size="small">
                    <el-icon><Delete /></el-icon>
                  </el-button>
                </div>
                <el-button @click="addIP" size="small">
                  <el-icon><Plus /></el-icon>
                  添加IP
                </el-button>
              </div>
            </el-form-item>

            <el-form-item label="API访问限制">
              <el-switch v-model="accessControl.enableAPIRateLimit" />
              <span class="setting-desc">启用API访问频率限制</span>
            </el-form-item>

            <el-form-item v-if="accessControl.enableAPIRateLimit" label="请求频率限制">
              <div style="display: flex; align-items: center; gap: 8px;">
                <el-input-number
                  v-model="accessControl.rateLimit.requests"
                  :min="1"
                  :max="1000"
                  style="width: 120px"
                />
                <span>次/</span>
                <el-select v-model="accessControl.rateLimit.period" style="width: 100px">
                  <el-option label="分钟" value="minute" />
                  <el-option label="小时" value="hour" />
                  <el-option label="天" value="day" />
                </el-select>
              </div>
            </el-form-item>

            <el-form-item label="跨域访问控制">
              <el-switch v-model="accessControl.enableCORS" />
              <span class="setting-desc">允许跨域请求</span>
            </el-form-item>

            <el-form-item v-if="accessControl.enableCORS" label="允许的域名">
              <el-input
                v-model="accessControl.allowedOrigins"
                type="textarea"
                :rows="3"
                placeholder="https://example.com&#10;https://app.example.com"
              />
              <span class="setting-desc">每行一个域名</span>
            </el-form-item>
          </el-form>
        </el-card>

        <!-- 数据安全 -->
        <el-card class="security-card" shadow="never">
          <template #header>
            <div class="card-header">
              <span>
                <el-icon><Document /></el-icon>
                数据安全
              </span>
            </div>
          </template>

          <el-form :model="dataSecurity" label-width="140px">
            <el-form-item label="数据加密">
              <el-switch v-model="dataSecurity.enableEncryption" />
              <span class="setting-desc">启用敏感数据加密存储</span>
            </el-form-item>

            <el-form-item label="备份加密">
              <el-switch v-model="dataSecurity.enableBackupEncryption" />
              <span class="setting-desc">备份文件使用加密</span>
            </el-form-item>

            <el-form-item label="操作日志">
              <el-switch v-model="dataSecurity.enableAuditLog" />
              <span class="setting-desc">记录所有操作日志</span>
            </el-form-item>

            <el-form-item label="日志保留期">
              <el-select v-model="dataSecurity.logRetention" style="width: 200px">
                <el-option label="30天" :value="30" />
                <el-option label="90天" :value="90" />
                <el-option label="180天" :value="180" />
                <el-option label="1年" :value="365" />
                <el-option label="永久保留" :value="-1" />
              </el-select>
            </el-form-item>

            <el-form-item label="敏感信息脱敏">
              <el-switch v-model="dataSecurity.enableDataMasking" />
              <span class="setting-desc">在日志中隐藏敏感信息</span>
            </el-form-item>
          </el-form>
        </el-card>
      </el-col>

      <!-- 右侧状态 -->
      <el-col :span="8">
        <!-- 安全状态 -->
        <el-card class="status-card" shadow="never">
          <template #header>
            <div class="card-header">
              <span>
                <el-icon><Monitor /></el-icon>
                安全状态
              </span>
              <el-button size="small" @click="loadSecurityStatus">
                <el-icon><Refresh /></el-icon>
              </el-button>
            </div>
          </template>

          <div class="security-status">
            <div class="status-item">
              <div class="status-label">
                <el-icon class="status-icon success"><SuccessFilled /></el-icon>
                系统安全评分
              </div>
              <div class="status-value">
                <span class="score">{{ securityStatus.score }}</span>/100
              </div>
            </div>

            <div class="status-item">
              <div class="status-label">
                <el-icon class="status-icon" :class="getVulnerabilityClass()"><WarningFilled /></el-icon>
                安全漏洞
              </div>
              <div class="status-value">{{ securityStatus.vulnerabilities }}个</div>
            </div>

            <div class="status-item">
              <div class="status-label">
                <el-icon class="status-icon info"><InfoFilled /></el-icon>
                最后检查
              </div>
              <div class="status-value">{{ formatDateTime(securityStatus.lastCheck) }}</div>
            </div>

            <div class="status-item">
              <div class="status-label">
                <el-icon class="status-icon warning"><Timer /></el-icon>
                证书过期
              </div>
              <div class="status-value">{{ securityStatus.certExpiry }}天</div>
            </div>
          </div>

          <el-button @click="handleSecurityScan" style="width: 100%; margin-top: 16px;">
            <el-icon><Search /></el-icon>
            安全扫描
          </el-button>
        </el-card>

        <!-- 威胁监控 -->
        <el-card class="threat-card" shadow="never">
          <template #header>
            <div class="card-header">
              <span>
                <el-icon><Warning /></el-icon>
                威胁监控
              </span>
            </div>
          </template>

          <div class="threat-list">
            <div
              v-for="threat in threatAlerts"
              :key="threat.id"
              class="threat-item"
              :class="threat.level"
            >
              <div class="threat-icon">
                <el-icon v-if="threat.level === 'high'"><WarningFilled /></el-icon>
                <el-icon v-else-if="threat.level === 'medium'"><Warning /></el-icon>
                <el-icon v-else><InfoFilled /></el-icon>
              </div>
              <div class="threat-content">
                <div class="threat-title">{{ threat.title }}</div>
                <div class="threat-time">{{ formatDateTime(threat.time) }}</div>
              </div>
              <div class="threat-action">
                <el-button size="small" @click="handleThreatAction(threat)">
                  处理
                </el-button>
              </div>
            </div>
          </div>
        </el-card>

        <!-- 安全建议 -->
        <el-card class="suggestions-card" shadow="never">
          <template #header>
            <div class="card-header">
              <span>
                <el-icon><InfoFilled /></el-icon>
                安全建议
              </span>
            </div>
          </template>

          <div class="suggestions-list">
            <div
              v-for="suggestion in securitySuggestions"
              :key="suggestion.id"
              class="suggestion-item"
            >
              <el-icon class="suggestion-icon"><CircleCheck /></el-icon>
              <span>{{ suggestion.text }}</span>
            </div>
          </div>
        </el-card>
      </el-col>
    </el-row>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, computed } from 'vue'
import { ElMessage } from 'element-plus'
import dayjs from 'dayjs'
import {
  Lock,
  Check,
  Key,
  Document,
  Monitor,
  Refresh,
  Search,
  Warning,
  WarningFilled,
  SuccessFilled,
  InfoFilled,
  Timer,
  CircleCheck,
  Delete,
  Plus
} from '@element-plus/icons-vue'

// 导入通用组件
import { PageHeader } from '@/components/common'

// 登录安全设置
const loginSecurity = reactive({
  passwordPolicy: ['minLength', 'requireNumbers'],
  passwordExpiry: 90,
  maxLoginAttempts: 5,
  lockoutDuration: 30,
  forceTwoFactor: false,
  sessionTimeout: 120
})

// 访问控制设置
const accessControl = reactive({
  enableIPWhitelist: false,
  allowedIPs: ['192.168.1.0/24'],
  enableAPIRateLimit: true,
  rateLimit: {
    requests: 100,
    period: 'minute'
  },
  enableCORS: true,
  allowedOrigins: 'https://huinong.com\nhttps://app.huinong.com'
})

// 数据安全设置
const dataSecurity = reactive({
  enableEncryption: true,
  enableBackupEncryption: true,
  enableAuditLog: true,
  logRetention: 90,
  enableDataMasking: true
})

// 安全状态
const securityStatus = ref({
  score: 85,
  vulnerabilities: 3,
  lastCheck: '2024-12-29 18:00:00',
  certExpiry: 45
})

// 威胁警报
const threatAlerts = ref([
  {
    id: 1,
    level: 'high',
    title: '检测到异常登录尝试',
    time: '2024-12-29 17:45:00'
  },
  {
    id: 2,
    level: 'medium',
    title: 'SQL注入攻击尝试',
    time: '2024-12-29 16:30:00'
  },
  {
    id: 3,
    level: 'low',
    title: '密码策略不符合要求',
    time: '2024-12-29 15:20:00'
  }
])

// 安全建议
const securitySuggestions = ref([
  { id: 1, text: '启用双因子认证以提高账户安全性' },
  { id: 2, text: '定期更新系统补丁和安全更新' },
  { id: 3, text: '配置强密码策略和密码过期提醒' },
  { id: 4, text: '启用IP白名单限制管理员访问' },
  { id: 5, text: '定期备份重要数据并测试恢复流程' }
])

// 计算属性
const breadcrumbs = computed(() => [
  { title: '首页', to: '/' },
  { title: '系统管理', to: '/system' },
  { title: '安全设置' }
])

// 方法
const handleSaveAll = () => {
  console.log('保存安全设置:', {
    loginSecurity,
    accessControl,
    dataSecurity
  })
  ElMessage.success('安全设置已保存')
}

const loadSecurityStatus = () => {
  console.log('刷新安全状态')
  ElMessage.success('安全状态已更新')
}

const handleSecurityScan = () => {
  console.log('开始安全扫描')
  ElMessage.info('安全扫描已启动，请稍候...')
}

const handleThreatAction = (threat: any) => {
  console.log('处理威胁:', threat)
  ElMessage.success(`威胁 "${threat.title}" 已标记为已处理`)
}

const addIP = () => {
  accessControl.allowedIPs.push('')
}

const removeIP = (index: number) => {
  accessControl.allowedIPs.splice(index, 1)
}

const getVulnerabilityClass = () => {
  const count = securityStatus.value.vulnerabilities
  if (count === 0) return 'success'
  if (count <= 2) return 'warning'
  return 'danger'
}

const formatDateTime = (datetime: string) => {
  return dayjs(datetime).format('MM-DD HH:mm')
}
</script>

<style scoped>
.security-view {
  height: 100vh;
  display: flex;
  flex-direction: column;
  gap: 16px;
  padding: 0 20px;
  overflow: hidden;
}

.security-card,
.status-card,
.threat-card,
.suggestions-card {
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

.ip-list {
  width: 100%;
}

.ip-item {
  display: flex;
  align-items: center;
  gap: 8px;
  margin-bottom: 8px;
}

.security-status {
  padding: 16px 0;
}

.status-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 16px;
  padding: 12px 0;
  border-bottom: 1px solid #f0f0f0;
}

.status-item:last-child {
  border-bottom: none;
}

.status-label {
  display: flex;
  align-items: center;
  gap: 8px;
  font-size: 14px;
  color: #666;
}

.status-icon.success {
  color: #67c23a;
}

.status-icon.warning {
  color: #e6a23c;
}

.status-icon.danger {
  color: #f56c6c;
}

.status-icon.info {
  color: #409eff;
}

.status-value {
  font-weight: 600;
  font-size: 16px;
}

.score {
  font-size: 20px;
  color: #67c23a;
}

.threat-list {
  padding: 16px 0;
}

.threat-item {
  display: flex;
  align-items: center;
  margin-bottom: 16px;
  padding: 12px;
  border-radius: 8px;
  border: 1px solid #f0f0f0;
}

.threat-item.high {
  border-color: #f56c6c;
  background-color: #fef0f0;
}

.threat-item.medium {
  border-color: #e6a23c;
  background-color: #fdf6ec;
}

.threat-item.low {
  border-color: #409eff;
  background-color: #ecf5ff;
}

.threat-icon {
  margin-right: 12px;
}

.threat-item.high .threat-icon {
  color: #f56c6c;
}

.threat-item.medium .threat-icon {
  color: #e6a23c;
}

.threat-item.low .threat-icon {
  color: #409eff;
}

.threat-content {
  flex: 1;
}

.threat-title {
  font-weight: 500;
  margin-bottom: 4px;
}

.threat-time {
  font-size: 12px;
  color: #666;
}

.suggestions-list {
  padding: 16px 0;
}

.suggestion-item {
  display: flex;
  align-items: center;
  gap: 8px;
  margin-bottom: 12px;
  font-size: 14px;
  color: #666;
}

.suggestion-icon {
  color: #67c23a;
}

:deep(.el-card__header) {
  padding: 16px 20px;
  border-bottom: 1px solid #f0f0f0;
}

:deep(.el-card__body) {
  padding: 20px;
}

:deep(.el-form-item__label) {
  font-weight: 500;
}
</style>
