<template>
  <div class="system-config">
    <!-- 配置导航 -->
    <el-card class="nav-card" shadow="never">
      <el-tabs v-model="activeTab" @tab-change="handleTabChange">
        <el-tab-pane label="基础配置" name="basic">
          <el-icon><Setting /></el-icon>
        </el-tab-pane>
        <el-tab-pane label="审批配置" name="approval">
          <el-icon><DocumentChecked /></el-icon>
        </el-tab-pane>
        <el-tab-pane label="通知配置" name="notification">
          <el-icon><Bell /></el-icon>
        </el-tab-pane>
        <el-tab-pane label="安全配置" name="security">
          <el-icon><Lock /></el-icon>
        </el-tab-pane>
        <el-tab-pane label="集成配置" name="integration">
          <el-icon><Connection /></el-icon>
        </el-tab-pane>
      </el-tabs>
    </el-card>

    <!-- 基础配置 -->
    <el-card v-show="activeTab === 'basic'" class="config-card" shadow="never">
      <template #header>
        <div class="card-header">
          <span>基础配置</span>
          <el-button type="primary" @click="saveBasicConfig" :loading="saving">保存配置</el-button>
        </div>
      </template>

      <el-form :model="basicConfig" label-width="150px">
        <el-row :gutter="20">
          <el-col :span="12">
            <el-form-item label="系统名称">
              <el-input v-model="basicConfig.system_name" placeholder="请输入系统名称" />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="系统版本">
              <el-input v-model="basicConfig.system_version" placeholder="请输入系统版本" />
            </el-form-item>
          </el-col>
        </el-row>

        <el-row :gutter="20">
          <el-col :span="12">
            <el-form-item label="公司名称">
              <el-input v-model="basicConfig.company_name" placeholder="请输入公司名称" />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="联系电话">
              <el-input v-model="basicConfig.contact_phone" placeholder="请输入联系电话" />
            </el-form-item>
          </el-col>
        </el-row>

        <el-form-item label="系统描述">
          <el-input
            v-model="basicConfig.system_description"
            type="textarea"
            :rows="3"
            placeholder="请输入系统描述"
          />
        </el-form-item>

        <el-row :gutter="20">
          <el-col :span="12">
            <el-form-item label="时区设置">
              <el-select v-model="basicConfig.timezone" placeholder="选择时区" style="width: 100%">
                <el-option label="Asia/Shanghai (UTC+8)" value="Asia/Shanghai" />
                <el-option label="Asia/Hong_Kong (UTC+8)" value="Asia/Hong_Kong" />
                <el-option label="UTC (UTC+0)" value="UTC" />
              </el-select>
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="默认语言">
              <el-select v-model="basicConfig.default_language" placeholder="选择默认语言" style="width: 100%">
                <el-option label="简体中文" value="zh-CN" />
                <el-option label="English" value="en-US" />
              </el-select>
            </el-form-item>
          </el-col>
        </el-row>

        <el-row :gutter="20">
          <el-col :span="12">
            <el-form-item label="维护模式">
              <el-switch
                v-model="basicConfig.maintenance_mode"
                active-text="开启"
                inactive-text="关闭"
              />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="用户注册">
              <el-switch
                v-model="basicConfig.allow_registration"
                active-text="允许"
                inactive-text="禁止"
              />
            </el-form-item>
          </el-col>
        </el-row>
      </el-form>
    </el-card>

    <!-- 审批配置 -->
    <el-card v-show="activeTab === 'approval'" class="config-card" shadow="never">
      <template #header>
        <div class="card-header">
          <span>审批配置</span>
          <el-button type="primary" @click="saveApprovalConfig" :loading="saving">保存配置</el-button>
        </div>
      </template>

      <el-form :model="approvalConfig" label-width="150px">
        <el-row :gutter="20">
          <el-col :span="12">
            <el-form-item label="自动审批金额上限">
              <el-input-number
                v-model="approvalConfig.auto_approval_limit"
                :min="0"
                :step="10000"
                style="width: 100%"
              />
              <div class="form-tip">超过此金额需要人工审批</div>
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="审批超时时间">
              <el-input-number
                v-model="approvalConfig.approval_timeout_hours"
                :min="1"
                :max="168"
                style="width: 100%"
              />
              <div class="form-tip">小时</div>
            </el-form-item>
          </el-col>
        </el-row>

        <el-row :gutter="20">
          <el-col :span="12">
            <el-form-item label="需要多级审批">
              <el-switch
                v-model="approvalConfig.multi_level_approval"
                active-text="是"
                inactive-text="否"
              />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="启用风险评估">
              <el-switch
                v-model="approvalConfig.enable_risk_assessment"
                active-text="启用"
                inactive-text="禁用"
              />
            </el-form-item>
          </el-col>
        </el-row>

        <el-form-item label="审批流程模板">
          <el-table :data="approvalConfig.workflow_templates" size="small">
            <el-table-column prop="name" label="模板名称" width="200" />
            <el-table-column prop="type" label="适用类型" width="120">
              <template #default="{ row }">
                <el-tag :type="row.type === 'loan' ? 'primary' : 'success'" size="small">
                  {{ row.type === 'loan' ? '贷款' : '租赁' }}
                </el-tag>
              </template>
            </el-table-column>
            <el-table-column prop="steps" label="审批步骤" />
            <el-table-column prop="enabled" label="状态" width="80">
              <template #default="{ row }">
                <el-tag :type="row.enabled ? 'success' : 'info'" size="small">
                  {{ row.enabled ? '启用' : '禁用' }}
                </el-tag>
              </template>
            </el-table-column>
            <el-table-column label="操作" width="120">
              <template #default="{ row }">
                <el-button size="small" @click="editWorkflowTemplate(row)">编辑</el-button>
              </template>
            </el-table-column>
          </el-table>
        </el-form-item>
      </el-form>
    </el-card>

    <!-- 通知配置 -->
    <el-card v-show="activeTab === 'notification'" class="config-card" shadow="never">
      <template #header>
        <div class="card-header">
          <span>通知配置</span>
          <el-button type="primary" @click="saveNotificationConfig" :loading="saving">保存配置</el-button>
        </div>
      </template>

      <el-form :model="notificationConfig" label-width="150px">
        <el-form-item label="邮件通知">
          <el-row :gutter="20">
            <el-col :span="8">
              <el-form-item label="SMTP服务器">
                <el-input v-model="notificationConfig.smtp_host" placeholder="smtp.example.com" />
              </el-form-item>
            </el-col>
            <el-col :span="8">
              <el-form-item label="SMTP端口">
                <el-input-number v-model="notificationConfig.smtp_port" :min="1" :max="65535" />
              </el-form-item>
            </el-col>
            <el-col :span="8">
              <el-form-item label="启用SSL">
                <el-switch v-model="notificationConfig.smtp_ssl" />
              </el-form-item>
            </el-col>
          </el-row>

          <el-row :gutter="20">
            <el-col :span="12">
              <el-form-item label="发送邮箱">
                <el-input v-model="notificationConfig.smtp_username" placeholder="noreply@example.com" />
              </el-form-item>
            </el-col>
            <el-col :span="12">
              <el-form-item label="邮箱密码">
                <el-input
                  v-model="notificationConfig.smtp_password"
                  type="password"
                  placeholder="请输入邮箱密码"
                  show-password
                />
              </el-form-item>
            </el-col>
          </el-row>
        </el-form-item>

        <el-form-item label="短信通知">
          <el-row :gutter="20">
            <el-col :span="12">
              <el-form-item label="短信服务商">
                <el-select v-model="notificationConfig.sms_provider" placeholder="选择短信服务商">
                  <el-option label="阿里云" value="aliyun" />
                  <el-option label="腾讯云" value="tencent" />
                  <el-option label="华为云" value="huawei" />
                </el-select>
              </el-form-item>
            </el-col>
            <el-col :span="12">
              <el-form-item label="启用短信">
                <el-switch v-model="notificationConfig.sms_enabled" />
              </el-form-item>
            </el-col>
          </el-row>
        </el-form-item>

        <el-form-item label="通知事件">
          <el-table :data="notificationConfig.notification_events" size="small">
            <el-table-column prop="event" label="事件类型" width="200" />
            <el-table-column prop="description" label="事件描述" />
            <el-table-column prop="email_enabled" label="邮件通知" width="100">
              <template #default="{ row }">
                <el-switch v-model="row.email_enabled" size="small" />
              </template>
            </el-table-column>
            <el-table-column prop="sms_enabled" label="短信通知" width="100">
              <template #default="{ row }">
                <el-switch v-model="row.sms_enabled" size="small" />
              </template>
            </el-table-column>
          </el-table>
        </el-form-item>
      </el-form>
    </el-card>

    <!-- 安全配置 -->
    <el-card v-show="activeTab === 'security'" class="config-card" shadow="never">
      <template #header>
        <div class="card-header">
          <span>安全配置</span>
          <el-button type="primary" @click="saveSecurityConfig" :loading="saving">保存配置</el-button>
        </div>
      </template>

      <el-form :model="securityConfig" label-width="150px">
        <el-row :gutter="20">
          <el-col :span="12">
            <el-form-item label="密码最小长度">
              <el-input-number
                v-model="securityConfig.min_password_length"
                :min="6"
                :max="20"
                style="width: 100%"
              />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="密码复杂度">
              <el-checkbox-group v-model="securityConfig.password_complexity">
                <el-checkbox value="uppercase">大写字母</el-checkbox>
                <el-checkbox value="lowercase">小写字母</el-checkbox>
                <el-checkbox value="numbers">数字</el-checkbox>
                <el-checkbox value="symbols">特殊字符</el-checkbox>
              </el-checkbox-group>
            </el-form-item>
          </el-col>
        </el-row>

        <el-row :gutter="20">
          <el-col :span="12">
            <el-form-item label="会话超时时间">
              <el-input-number
                v-model="securityConfig.session_timeout_minutes"
                :min="5"
                :max="1440"
                style="width: 100%"
              />
              <div class="form-tip">分钟</div>
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="登录失败限制">
              <el-input-number
                v-model="securityConfig.max_login_attempts"
                :min="3"
                :max="10"
                style="width: 100%"
              />
              <div class="form-tip">次</div>
            </el-form-item>
          </el-col>
        </el-row>

        <el-row :gutter="20">
          <el-col :span="12">
            <el-form-item label="启用双因子认证">
              <el-switch v-model="securityConfig.enable_2fa" />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="强制HTTPS">
              <el-switch v-model="securityConfig.force_https" />
            </el-form-item>
          </el-col>
        </el-row>

        <el-form-item label="IP白名单">
          <el-input
            v-model="securityConfig.ip_whitelist"
            type="textarea"
            :rows="3"
            placeholder="每行一个IP地址或IP段，如：192.168.1.1 或 192.168.1.0/24"
          />
        </el-form-item>
      </el-form>
    </el-card>

    <!-- 集成配置 -->
    <el-card v-show="activeTab === 'integration'" class="config-card" shadow="never">
      <template #header>
        <div class="card-header">
          <span>集成配置</span>
          <el-button type="primary" @click="saveIntegrationConfig" :loading="saving">保存配置</el-button>
        </div>
      </template>

      <el-form :model="integrationConfig" label-width="150px">
        <el-form-item label="数据库配置">
          <el-row :gutter="20">
            <el-col :span="12">
              <el-form-item label="数据库类型">
                <el-select v-model="integrationConfig.database_type" placeholder="选择数据库类型">
                  <el-option label="TiDB" value="tidb" />
                  <el-option label="MySQL" value="mysql" />
                  <el-option label="PostgreSQL" value="postgresql" />
                </el-select>
              </el-form-item>
            </el-col>
            <el-col :span="12">
              <el-form-item label="连接池大小">
                <el-input-number v-model="integrationConfig.db_pool_size" :min="1" :max="100" />
              </el-form-item>
            </el-col>
          </el-row>
        </el-form-item>

        <el-form-item label="API配置">
          <el-row :gutter="20">
            <el-col :span="12">
              <el-form-item label="API限流">
                <el-input-number
                  v-model="integrationConfig.api_rate_limit"
                  :min="10"
                  :max="10000"
                  style="width: 100%"
                />
                <div class="form-tip">每分钟请求数</div>
              </el-form-item>
            </el-col>
            <el-col :span="12">
              <el-form-item label="API版本">
                <el-input v-model="integrationConfig.api_version" placeholder="v1.0" />
              </el-form-item>
            </el-col>
          </el-row>
        </el-form-item>

        <el-form-item label="第三方服务">
          <el-table :data="integrationConfig.third_party_services" size="small">
            <el-table-column prop="name" label="服务名称" width="150" />
            <el-table-column prop="endpoint" label="服务地址" />
            <el-table-column prop="enabled" label="状态" width="80">
              <template #default="{ row }">
                <el-tag :type="row.enabled ? 'success' : 'info'" size="small">
                  {{ row.enabled ? '启用' : '禁用' }}
                </el-tag>
              </template>
            </el-table-column>
            <el-table-column label="操作" width="120">
              <template #default="{ row }">
                <el-button size="small" @click="testConnection(row)">测试连接</el-button>
              </template>
            </el-table-column>
          </el-table>
        </el-form-item>
      </el-form>
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive } from 'vue'
import { 
  Setting, 
  DocumentChecked, 
  Bell, 
  Lock, 
  Connection 
} from '@element-plus/icons-vue'
import { ElMessage } from 'element-plus'

// 当前标签页
const activeTab = ref('basic')
const saving = ref(false)

// 基础配置
const basicConfig = reactive({
  system_name: '惠农金融管理系统',
  system_version: '1.0.0',
  company_name: '惠农金融科技有限公司',
  contact_phone: '400-800-8888',
  system_description: '专业的农业金融管理平台，提供贷款、租赁等金融服务',
  timezone: 'Asia/Shanghai',
  default_language: 'zh-CN',
  maintenance_mode: false,
  allow_registration: true
})

// 审批配置
const approvalConfig = reactive({
  auto_approval_limit: 500000,
  approval_timeout_hours: 72,
  multi_level_approval: true,
  enable_risk_assessment: true,
  workflow_templates: [
    {
      name: '标准贷款审批',
      type: 'loan',
      steps: '初审 → 风险评估 → 终审',
      enabled: true
    },
    {
      name: '快速租赁审批',
      type: 'lease',
      steps: '资质审核 → 合同审核',
      enabled: true
    },
    {
      name: '大额贷款审批',
      type: 'loan',
      steps: '初审 → 风险评估 → 委员会审核 → 终审',
      enabled: true
    }
  ]
})

// 通知配置
const notificationConfig = reactive({
  smtp_host: 'smtp.example.com',
  smtp_port: 587,
  smtp_ssl: true,
  smtp_username: 'noreply@example.com',
  smtp_password: '',
  sms_provider: 'aliyun',
  sms_enabled: true,
  notification_events: [
    {
      event: 'application_submitted',
      description: '申请提交',
      email_enabled: true,
      sms_enabled: false
    },
    {
      event: 'application_approved',
      description: '申请通过',
      email_enabled: true,
      sms_enabled: true
    },
    {
      event: 'application_rejected',
      description: '申请拒绝',
      email_enabled: true,
      sms_enabled: true
    },
    {
      event: 'payment_due',
      description: '还款提醒',
      email_enabled: true,
      sms_enabled: true
    }
  ]
})

// 安全配置
const securityConfig = reactive({
  min_password_length: 8,
  password_complexity: ['lowercase', 'numbers'],
  session_timeout_minutes: 120,
  max_login_attempts: 5,
  enable_2fa: false,
  force_https: true,
  ip_whitelist: ''
})

// 集成配置
const integrationConfig = reactive({
  database_type: 'tidb',
  db_pool_size: 20,
  api_rate_limit: 1000,
  api_version: 'v1.0',
  third_party_services: [
    {
      name: '征信系统',
      endpoint: 'https://api.credit.com',
      enabled: true
    },
    {
      name: '支付网关',
      endpoint: 'https://api.payment.com',
      enabled: true
    },
    {
      name: '短信服务',
      endpoint: 'https://api.sms.com',
      enabled: true
    }
  ]
})

// 标签页切换
const handleTabChange = (tabName: string) => {
  console.log('切换到标签页:', tabName)
}

// 保存配置
const saveBasicConfig = async () => {
  saving.value = true
  try {
    await new Promise(resolve => setTimeout(resolve, 1000))
    ElMessage.success('基础配置保存成功')
  } finally {
    saving.value = false
  }
}

const saveApprovalConfig = async () => {
  saving.value = true
  try {
    await new Promise(resolve => setTimeout(resolve, 1000))
    ElMessage.success('审批配置保存成功')
  } finally {
    saving.value = false
  }
}

const saveNotificationConfig = async () => {
  saving.value = true
  try {
    await new Promise(resolve => setTimeout(resolve, 1000))
    ElMessage.success('通知配置保存成功')
  } finally {
    saving.value = false
  }
}

const saveSecurityConfig = async () => {
  saving.value = true
  try {
    await new Promise(resolve => setTimeout(resolve, 1000))
    ElMessage.success('安全配置保存成功')
  } finally {
    saving.value = false
  }
}

const saveIntegrationConfig = async () => {
  saving.value = true
  try {
    await new Promise(resolve => setTimeout(resolve, 1000))
    ElMessage.success('集成配置保存成功')
  } finally {
    saving.value = false
  }
}

// 编辑工作流模板
const editWorkflowTemplate = (template: any) => {
  ElMessage.info(`编辑模板：${template.name}`)
}

// 测试连接
const testConnection = (service: any) => {
  ElMessage.info(`测试${service.name}连接...`)
}
</script>

<style scoped>
.system-config {
  padding: 20px;
  background: #f5f5f5;
  min-height: calc(100vh - 60px);
}

.nav-card, .config-card {
  border: none;
  margin-bottom: 20px;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  font-weight: 500;
}

.form-tip {
  font-size: 12px;
  color: #909399;
  margin-top: 4px;
}

:deep(.el-tab-pane) {
  padding: 0;
}

:deep(.el-tabs__item) {
  display: flex;
  align-items: center;
  gap: 8px;
}

:deep(.el-form-item) {
  margin-bottom: 18px;
}

@media (max-width: 768px) {
  .system-config {
    padding: 10px;
  }
  
  .card-header {
    flex-direction: column;
    gap: 12px;
    align-items: stretch;
  }
}
</style>
