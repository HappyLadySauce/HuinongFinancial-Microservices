<template>
  <div class="profile-view">
    <!-- 页面头部 -->
    <PageHeader
      title="个人信息"
      subtitle="查看和编辑个人资料"
      :icon="User"
      :breadcrumbs="breadcrumbs"
    >
      <template #actions>
        <el-button type="primary" @click="handleEdit">
          <el-icon><Edit /></el-icon>
          编辑资料
        </el-button>
      </template>
    </PageHeader>

    <el-row :gutter="24">
      <!-- 左侧个人信息 -->
      <el-col :span="16">
        <el-card class="info-card" shadow="never">
          <template #header>
            <div class="card-header">
              <span>
                <el-icon><User /></el-icon>
                基本信息
              </span>
            </div>
          </template>

          <el-descriptions :column="2" border>
            <el-descriptions-item label="用户名">
              {{ userInfo.username }}
            </el-descriptions-item>
            <el-descriptions-item label="姓名">
              {{ userInfo.realName }}
            </el-descriptions-item>
            <el-descriptions-item label="邮箱">
              {{ userInfo.email }}
            </el-descriptions-item>
            <el-descriptions-item label="手机号">
              {{ userInfo.phone }}
            </el-descriptions-item>
            <el-descriptions-item label="部门">
              {{ userInfo.department }}
            </el-descriptions-item>
            <el-descriptions-item label="职位">
              {{ userInfo.position }}
            </el-descriptions-item>
            <el-descriptions-item label="角色">
              <StatusTag :status="userInfo.role" />
            </el-descriptions-item>
            <el-descriptions-item label="状态">
              <StatusTag :status="userInfo.status" :show-icon="true" />
            </el-descriptions-item>
            <el-descriptions-item label="注册时间">
              {{ formatDateTime(userInfo.createdAt) }}
            </el-descriptions-item>
            <el-descriptions-item label="最后登录">
              {{ formatDateTime(userInfo.lastLoginAt) }}
            </el-descriptions-item>
          </el-descriptions>
        </el-card>

        <!-- 权限信息 -->
        <el-card class="permissions-card" shadow="never">
          <template #header>
            <div class="card-header">
              <span>
                <el-icon><Key /></el-icon>
                权限信息
              </span>
            </div>
          </template>

          <div class="permissions-list">
            <el-tag
              v-for="permission in userInfo.permissions"
              :key="permission"
              size="small"
              type="primary"
              class="permission-tag"
            >
              {{ getPermissionName(permission) }}
            </el-tag>
          </div>
        </el-card>
      </el-col>

      <!-- 右侧头像和操作 -->
      <el-col :span="8">
        <el-card class="avatar-card" shadow="never">
          <template #header>
            <div class="card-header">
              <span>
                <el-icon><Avatar /></el-icon>
                头像设置
              </span>
            </div>
          </template>

          <div class="avatar-section">
            <el-avatar :size="120" :src="userInfo.avatar">
              {{ userInfo.realName?.charAt(0) }}
            </el-avatar>
            <el-button @click="handleAvatarUpload" style="margin-top: 16px;">
              <el-icon><Upload /></el-icon>
              更换头像
            </el-button>
          </div>
        </el-card>

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
            <el-button @click="handleChangePassword" style="width: 100%; margin-bottom: 12px;">
              <el-icon><Lock /></el-icon>
              修改密码
            </el-button>
            <el-button @click="handleLogout" style="width: 100%;" type="danger">
              <el-icon><SwitchButton /></el-icon>
              退出登录
            </el-button>
          </div>
        </el-card>
      </el-col>
    </el-row>

    <!-- 编辑个人信息对话框 -->
    <FormDialog
      v-model="editDialog.visible"
      title="编辑个人信息"
      :mode="'edit'"
      :initial-data="editDialog.data"
      :rules="editRules"
      @submit="handleEditSubmit"
      @cancel="editDialog.visible = false"
    >
      <template #form="{ form }">
        <el-form-item label="姓名" prop="realName">
          <el-input v-model="form.realName" />
        </el-form-item>
        
        <el-form-item label="邮箱" prop="email">
          <el-input v-model="form.email" />
        </el-form-item>
        
        <el-form-item label="手机号" prop="phone">
          <el-input v-model="form.phone" />
        </el-form-item>
        
        <el-form-item label="部门" prop="department">
          <el-input v-model="form.department" />
        </el-form-item>
        
        <el-form-item label="职位" prop="position">
          <el-input v-model="form.position" />
        </el-form-item>
      </template>
    </FormDialog>

    <!-- 修改密码对话框 -->
    <FormDialog
      v-model="passwordDialog.visible"
      title="修改密码"
      :mode="'edit'"
      :initial-data="passwordDialog.data"
      :rules="passwordRules"
      @submit="handlePasswordSubmit"
      @cancel="passwordDialog.visible = false"
    >
      <template #form="{ form }">
        <el-form-item label="当前密码" prop="currentPassword">
          <el-input v-model="form.currentPassword" type="password" show-password />
        </el-form-item>
        
        <el-form-item label="新密码" prop="newPassword">
          <el-input v-model="form.newPassword" type="password" show-password />
        </el-form-item>
        
        <el-form-item label="确认密码" prop="confirmPassword">
          <el-input v-model="form.confirmPassword" type="password" show-password />
        </el-form-item>
      </template>
    </FormDialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, computed } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage, ElMessageBox, type FormRules } from 'element-plus'
import dayjs from 'dayjs'
import {
  User,
  Edit,
  Key,
  Avatar,
  Upload,
  Tools,
  Lock,
  SwitchButton
} from '@element-plus/icons-vue'

// 导入通用组件
import { PageHeader, StatusTag, FormDialog } from '@/components/common'

const router = useRouter()

const userInfo = ref({
  id: 1,
  username: 'admin',
  realName: '系统管理员',
  email: 'admin@huinong.com',
  phone: '13888888888',
  department: '信息技术部',
  position: '系统管理员',
  role: 'admin',
  status: 'active',
  avatar: '',
  permissions: ['user:view', 'user:edit', 'approval:view', 'approval:edit', 'admin'],
  createdAt: '2024-01-01T00:00:00',
  lastLoginAt: '2024-12-29T10:30:00'
})

const editDialog = reactive({
  visible: false,
  data: {}
})

const passwordDialog = reactive({
  visible: false,
  data: {}
})

// 计算属性
const breadcrumbs = computed(() => [
  { title: '首页', to: '/' },
  { title: '个人信息' }
])

// 表单验证规则
const editRules: FormRules = {
  realName: [
    { required: true, message: '请输入姓名', trigger: 'blur' }
  ],
  email: [
    { required: true, message: '请输入邮箱', trigger: 'blur' },
    { type: 'email' as const, message: '请输入有效的邮箱地址', trigger: 'blur' }
  ],
  phone: [
    { required: true, message: '请输入手机号', trigger: 'blur' },
    { pattern: /^1[3-9]\d{9}$/, message: '请输入有效的手机号', trigger: 'blur' }
  ]
}

const passwordRules = {
  currentPassword: [
    { required: true, message: '请输入当前密码', trigger: 'blur' }
  ],
  newPassword: [
    { required: true, message: '请输入新密码', trigger: 'blur' },
    { min: 6, message: '密码长度不能少于6位', trigger: 'blur' }
  ],
  confirmPassword: [
    { required: true, message: '请确认新密码', trigger: 'blur' },
    {
      validator: (rule: any, value: string, callback: (error?: Error) => void) => {
        if (value !== (passwordDialog.data as any).newPassword) {
          callback(new Error('两次输入的密码不一致'))
        } else {
          callback()
        }
      },
      trigger: 'blur'
    }
  ]
}

// 方法
const handleEdit = () => {
  editDialog.data = {
    realName: userInfo.value.realName,
    email: userInfo.value.email,
    phone: userInfo.value.phone,
    department: userInfo.value.department,
    position: userInfo.value.position
  }
  editDialog.visible = true
}

const handleEditSubmit = async (data: any) => {
  try {
    console.log('更新个人信息:', data)
    // 实现更新逻辑
    Object.assign(userInfo.value, data)
    editDialog.visible = false
    ElMessage.success('个人信息更新成功')
  } catch (error) {
    console.error('更新失败:', error)
    ElMessage.error('更新失败')
  }
}

const handleChangePassword = () => {
  passwordDialog.data = {}
  passwordDialog.visible = true
}

const handlePasswordSubmit = async (data: any) => {
  try {
    console.log('修改密码:', data)
    // 实现密码修改逻辑
    passwordDialog.visible = false
    ElMessage.success('密码修改成功')
  } catch (error) {
    console.error('密码修改失败:', error)
    ElMessage.error('密码修改失败')
  }
}

const handleAvatarUpload = () => {
  ElMessage.info('头像上传功能待实现')
}

const handleLogout = async () => {
  try {
    await ElMessageBox.confirm('确定要退出登录吗？', '确认', {
      type: 'warning'
    })
    
    // 清除登录状态
    localStorage.removeItem('token')
    router.push('/login')
    ElMessage.success('已安全退出')
  } catch {
    // 用户取消
  }
}

const getPermissionName = (permission: string) => {
  const permissionMap: Record<string, string> = {
    'user:view': '查看用户',
    'user:edit': '编辑用户',
    'approval:view': '查看审批',
    'approval:edit': '编辑审批',
    'admin': '系统管理'
  }
  return permissionMap[permission] || permission
}

const formatDateTime = (datetime: string) => {
  return dayjs(datetime).format('YYYY-MM-DD HH:mm:ss')
}
</script>

<style scoped>
.profile-view {
  height: 100vh;
  display: flex;
  flex-direction: column;
  gap: 16px;
  padding: 0 20px;
  overflow: hidden;
}

.info-card,
.permissions-card,
.avatar-card,
.actions-card {
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

.permissions-list {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
  padding: 16px 0;
}

.permission-tag {
  margin: 0;
}

.avatar-section {
  display: flex;
  flex-direction: column;
  align-items: center;
  padding: 20px 0;
}

.action-buttons {
  padding: 16px 0;
}

:deep(.el-descriptions__cell) {
  padding: 12px;
}

:deep(.el-card__header) {
  padding: 16px 20px;
  border-bottom: 1px solid #f0f0f0;
}

:deep(.el-card__body) {
  padding: 20px;
}
</style> 