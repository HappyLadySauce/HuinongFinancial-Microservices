<template>
  <div class="users-view">
    <div class="page-header">
      <div class="page-title-container">
        <h2 class="page-title">用户管理</h2>
      </div>
      <div class="header-actions">
        <el-button type="primary" @click="showAddUserDialog">
          <el-icon><Plus /></el-icon>
          添加用户
        </el-button>
      </div>
    </div>

    <!-- 筛选条件 -->
    <el-card class="filter-card" shadow="never">
      <el-form :model="filterForm" inline size="default">
        <el-form-item label="用户角色">
          <el-select v-model="filterForm.role" placeholder="请选择角色" clearable>
            <el-option label="全部" value="" />
            <el-option label="管理员" value="ADMIN" />
            <el-option label="审批员" value="审批员" />
          </el-select>
        </el-form-item>
        
        <el-form-item label="用户名">
          <el-input
            v-model="filterForm.username"
            placeholder="请输入用户名"
            clearable
          />
        </el-form-item>
        
        <el-form-item>
          <el-button type="primary" @click="handleSearch">
            <el-icon><Search /></el-icon>
            搜索
          </el-button>
          <el-button @click="handleReset">
            <el-icon><RefreshLeft /></el-icon>
            重置
          </el-button>
        </el-form-item>
      </el-form>
    </el-card>

    <!-- 用户列表 -->
    <el-card class="table-card" shadow="never">
      <template #header>
        <div class="card-header">
          <span>用户列表</span>
          <div class="header-extra">
            <el-tag type="info">
              共 {{ pagination.total }} 个用户
            </el-tag>
          </div>
        </div>
      </template>
      
      <el-table
        v-loading="loading"
        :data="users"
        stripe
        style="width: 100%"
      >
        <el-table-column prop="admin_user_id" label="用户ID" width="180" />
        
        <el-table-column prop="username" label="用户名" width="120" />
        
        <el-table-column prop="display_name" label="显示名称" width="120" />
        
        <el-table-column prop="role" label="角色" width="100">
          <template #default="{ row }">
            <el-tag :type="getRoleType(row.role)" size="small">
              {{ getRoleName(row.role) }}
            </el-tag>
          </template>
        </el-table-column>
        
        <el-table-column prop="email" label="邮箱" min-width="200" />
        
        <el-table-column prop="status" label="状态" width="100">
          <template #default="{ row }">
            <el-tag :type="row.status === 0 ? 'success' : 'danger'" size="small">
              {{ row.status === 0 ? '正常' : '禁用' }}
            </el-tag>
          </template>
        </el-table-column>
        
        <el-table-column prop="created_at" label="创建时间" width="160">
          <template #default="{ row }">
            {{ formatDateTime(row.created_at) }}
          </template>
        </el-table-column>
        
        <el-table-column label="操作" width="200" fixed="right">
          <template #default="{ row }">
            <el-button
              v-if="row.status === 0"
              type="warning"
              size="small"
              @click="updateUserStatus(row, 1)"
            >
              禁用
            </el-button>
            <el-button
              v-else
              type="success"
              size="small"
              @click="updateUserStatus(row, 0)"
            >
              启用
            </el-button>
            <el-button
              type="primary"
              size="small"
              @click="editUser(row)"
            >
              编辑
            </el-button>
          </template>
        </el-table-column>
      </el-table>
      
      <!-- 分页 -->
      <div class="pagination-wrapper">
        <el-pagination
          v-model:current-page="pagination.page"
          v-model:page-size="pagination.limit"
          :page-sizes="[10, 20, 50, 100]"
          :total="pagination.total"
          layout="total, sizes, prev, pager, next, jumper"
          @size-change="handleSizeChange"
          @current-change="handleCurrentChange"
        />
      </div>
    </el-card>

    <!-- 添加/编辑用户对话框 -->
    <el-dialog
      v-model="userDialogVisible"
      :title="isEditMode ? '编辑用户' : '添加用户'"
      width="500px"
      @close="resetUserForm"
    >
      <el-form
        ref="userFormRef"
        :model="userForm"
        :rules="userRules"
        label-width="100px"
      >
        <el-form-item label="用户名" prop="username">
          <el-input
            v-model="userForm.username"
            placeholder="请输入用户名"
            :disabled="isEditMode"
          />
        </el-form-item>
        
        <el-form-item v-if="!isEditMode" label="密码" prop="password">
          <el-input
            v-model="userForm.password"
            type="password"
            placeholder="请输入密码"
            show-password
          />
        </el-form-item>
        
        <el-form-item label="显示名称" prop="display_name">
          <el-input
            v-model="userForm.display_name"
            placeholder="请输入显示名称"
          />
        </el-form-item>
        
        <el-form-item label="角色" prop="role">
          <el-select v-model="userForm.role" placeholder="请选择角色" style="width: 100%">
            <el-option label="管理员" value="ADMIN" />
            <el-option label="审批员" value="审批员" />
          </el-select>
        </el-form-item>
        
        <el-form-item label="邮箱" prop="email">
          <el-input
            v-model="userForm.email"
            placeholder="请输入邮箱地址"
          />
        </el-form-item>
      </el-form>
      
      <template #footer>
        <el-button @click="userDialogVisible = false">取消</el-button>
        <el-button
          type="primary"
          @click="submitUserForm"
          :loading="submitting"
        >
          {{ isEditMode ? '更新' : '创建' }}
        </el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { ElMessage, ElMessageBox, type FormInstance, type FormRules } from 'element-plus'
import {
  Plus,
  Refresh,
  Search,
  RefreshLeft
} from '@element-plus/icons-vue'
import { adminUserApi } from '@/services/api'
import type { AdminUser, PaginationResponse } from '@/types'
import dayjs from 'dayjs'

const loading = ref(false)
const submitting = ref(false)
const userDialogVisible = ref(false)
const isEditMode = ref(false)
const userFormRef = ref<FormInstance>()

// 筛选表单
const filterForm = reactive({
  role: '',
  username: ''
})

// 分页信息
const pagination = reactive({
  page: 1,
  limit: 20,
  total: 0
})

// 用户列表
const users = ref<AdminUser[]>([])

// 用户表单
const userForm = reactive({
  username: '',
  password: '',
  display_name: '',
  role: '',
  email: ''
})

const userRules: FormRules = {
  username: [
    { required: true, message: '请输入用户名', trigger: 'blur' },
    { min: 3, max: 20, message: '用户名长度在 3 到 20 个字符', trigger: 'blur' }
  ],
  password: [
    { required: true, message: '请输入密码', trigger: 'blur' },
    { min: 6, max: 20, message: '密码长度在 6 到 20 个字符', trigger: 'blur' }
  ],
  display_name: [
    { required: true, message: '请输入显示名称', trigger: 'blur' }
  ],
  role: [
    { required: true, message: '请选择角色', trigger: 'change' }
  ],
  email: [
    { required: true, message: '请输入邮箱地址', trigger: 'blur' },
    { type: 'email', message: '请输入正确的邮箱格式', trigger: 'blur' }
  ]
}

// 方法
const fetchUsers = async () => {
  try {
    loading.value = true
    const params = {
      role: filterForm.role,
      page: pagination.page,
      size: pagination.limit
    }
    
    const data = await adminUserApi.getUsers(params)
    users.value = data.list || []
    pagination.total = data.total || 0
  } catch (error) {
    ElMessage.error('获取用户列表失败')
  } finally {
    loading.value = false
  }
}

const refreshData = () => {
  fetchUsers()
}

const handleSearch = () => {
  pagination.page = 1
  fetchUsers()
}

const handleReset = () => {
  Object.assign(filterForm, {
    role: '',
    username: ''
  })
  pagination.page = 1
  fetchUsers()
}

const handleSizeChange = () => {
  pagination.page = 1
  fetchUsers()
}

const handleCurrentChange = () => {
  fetchUsers()
}

const showAddUserDialog = () => {
  isEditMode.value = false
  userDialogVisible.value = true
}

const editUser = (user: AdminUser) => {
  isEditMode.value = true
  Object.assign(userForm, {
    username: user.username,
    password: '',
    display_name: user.display_name,
    role: user.role,
    email: user.email
  })
  userDialogVisible.value = true
}

const resetUserForm = () => {
  Object.assign(userForm, {
    username: '',
    password: '',
    display_name: '',
    role: '',
    email: ''
  })
  isEditMode.value = false
}

const submitUserForm = async () => {
  if (!userFormRef.value) return
  
  try {
    await userFormRef.value.validate()
    submitting.value = true
    
    if (!isEditMode.value) {
      await adminUserApi.createUser({
        username: userForm.username,
        password: userForm.password,
        display_name: userForm.display_name,
        role: userForm.role,
        email: userForm.email
      })
      ElMessage.success('用户创建成功')
    } else {
      // 编辑用户功能暂时不实现，可以在这里添加编辑用户的API调用
      ElMessage.info('编辑用户功能开发中...')
    }
    
    userDialogVisible.value = false
    resetUserForm()
    fetchUsers()
  } catch (error) {
    ElMessage.error(isEditMode.value ? '更新用户失败' : '创建用户失败')
  } finally {
    submitting.value = false
  }
}

const updateUserStatus = async (user: AdminUser, status: number) => {
  try {
    const action = status === 0 ? '启用' : '禁用'
    await ElMessageBox.confirm(`确定要${action}用户 "${user.username}" 吗？`, '提示', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning'
    })
    
    await adminUserApi.updateUserStatus(user.admin_user_id, status)
    ElMessage.success(`用户${action}成功`)
    fetchUsers()
  } catch (error: any) {
    if (error === 'cancel') return
    ElMessage.error('更新用户状态失败')
  }
}

const getRoleType = (role: string) => {
  return role === 'ADMIN' ? 'danger' : 'primary'
}

const getRoleName = (role: string) => {
  const roleMap: Record<string, string> = {
    'ADMIN': '管理员',
    '审批员': '审批员'
  }
  return roleMap[role] || role
}

const formatDateTime = (datetime: string) => {
  return dayjs(datetime).format('YYYY-MM-DD HH:mm')
}

onMounted(() => {
  fetchUsers()
})
</script>

<style scoped>
.users-view {
  padding: 0;
}

.page-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 20px;
}

.page-title-container {
  display: flex;
  align-items: center;
}

.page-title {
  margin: 0;
  color: #1976D2;
  font-size: 24px;
  font-weight: 600;
}

.header-actions {
  display: flex;
  gap: 12px;
}

.filter-card {
  margin-bottom: 20px;
  border-radius: 8px;
}

.table-card {
  border-radius: 8px;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  font-weight: 600;
}

.pagination-wrapper {
  display: flex;
  justify-content: center;
  margin-top: 20px;
}

:deep(.el-form-item) {
  margin-bottom: 16px;
}
</style> 