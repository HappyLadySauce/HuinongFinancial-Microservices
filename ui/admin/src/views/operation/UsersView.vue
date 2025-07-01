<template>
  <div class="users-management">
    <!-- 顶部操作栏 -->
    <el-card class="operation-card" shadow="never">
      <div class="operation-row">
        <div class="operation-left">
          <el-button type="primary" @click="handleAddUser">
            <el-icon><Plus /></el-icon>
            新增用户
          </el-button>
          <el-button @click="handleImport">
            <el-icon><Upload /></el-icon>
            批量导入
          </el-button>
          <el-button @click="handleExport">
            <el-icon><Download /></el-icon>
            导出数据
          </el-button>
          <el-button type="warning" @click="handleBatchDisable" :disabled="selectedUsers.length === 0">
            <el-icon><Lock /></el-icon>
            批量禁用
          </el-button>
        </div>
        
        <div class="operation-right">
          <el-select v-model="filterRole" placeholder="用户角色" clearable style="width: 120px">
            <el-option label="全部角色" value="" />
            <el-option label="管理员" value="admin" />
            <el-option label="审批员" value="approver" />
            <el-option label="普通用户" value="user" />
          </el-select>
          <el-select v-model="filterStatus" placeholder="用户状态" clearable style="width: 120px">
            <el-option label="全部状态" value="" />
            <el-option label="正常" value="active" />
            <el-option label="禁用" value="disabled" />
            <el-option label="待审核" value="pending" />
          </el-select>
          <el-input
            v-model="searchKeyword"
            placeholder="搜索用户名或手机号"
            style="width: 200px"
            @keyup.enter="handleSearch"
          >
            <template #prefix>
              <el-icon><Search /></el-icon>
            </template>
          </el-input>
          <el-button type="primary" @click="handleSearch">搜索</el-button>
        </div>
      </div>
    </el-card>

    <!-- 用户统计卡片 -->
    <div class="stats-row">
      <el-card class="stat-card" shadow="hover">
        <div class="stat-content">
          <div class="stat-number">{{ stats.total }}</div>
          <div class="stat-label">总用户数</div>
        </div>
      </el-card>
      <el-card class="stat-card" shadow="hover">
        <div class="stat-content">
          <div class="stat-number">{{ stats.active }}</div>
          <div class="stat-label">活跃用户</div>
        </div>
      </el-card>
      <el-card class="stat-card" shadow="hover">
        <div class="stat-content">
          <div class="stat-number">{{ stats.disabled }}</div>
          <div class="stat-label">禁用用户</div>
        </div>
      </el-card>
      <el-card class="stat-card" shadow="hover">
        <div class="stat-content">
          <div class="stat-number">{{ stats.todayLogin }}</div>
          <div class="stat-label">今日登录</div>
        </div>
      </el-card>
    </div>

    <!-- 用户列表 -->
    <el-card class="list-card" shadow="never">
      <template #header>
        <div class="card-header">
          <span>用户列表</span>
          <div class="header-actions">
            <el-button size="small" @click="handleRefresh">
              <el-icon><Refresh /></el-icon>
              刷新
            </el-button>
          </div>
        </div>
      </template>

      <el-table
        v-loading="loading"
        :data="userList"
        stripe
        @selection-change="handleSelectionChange"
      >
        <el-table-column type="selection" width="55" />
        
        <el-table-column prop="id" label="用户ID" width="100" />
        
        <el-table-column prop="username" label="用户名" width="120" />
        
        <el-table-column prop="real_name" label="真实姓名" width="120">
          <template #default="{ row }">
            <span v-if="row.real_name">{{ row.real_name }}</span>
            <span v-else class="text-placeholder">未设置</span>
          </template>
        </el-table-column>
        
        <el-table-column prop="phone" label="手机号" width="140">
          <template #default="{ row }">
            <span v-if="row.phone">{{ row.phone }}</span>
            <span v-else class="text-placeholder">未设置</span>
          </template>
        </el-table-column>
        
        <el-table-column prop="email" label="邮箱" width="180" show-overflow-tooltip>
          <template #default="{ row }">
            <span v-if="row.email">{{ row.email }}</span>
            <span v-else class="text-placeholder">未设置</span>
          </template>
        </el-table-column>
        
        <el-table-column prop="role" label="角色" width="100">
          <template #default="{ row }">
            <el-tag :type="getRoleColor(row.role)" size="small">
              {{ getRoleName(row.role) }}
            </el-tag>
          </template>
        </el-table-column>
        
        <el-table-column prop="status" label="状态" width="100">
          <template #default="{ row }">
            <el-tag :type="getStatusColor(row.status)" size="small">
              {{ getStatusName(row.status) }}
            </el-tag>
          </template>
        </el-table-column>
        
        <el-table-column prop="last_login_at" label="最后登录" width="160">
          <template #default="{ row }">
            <span v-if="row.last_login_at">{{ formatDateTime(row.last_login_at) }}</span>
            <span v-else class="text-placeholder">从未登录</span>
          </template>
        </el-table-column>
        
        <el-table-column prop="created_at" label="注册时间" width="160">
          <template #default="{ row }">
            {{ formatDateTime(row.created_at) }}
          </template>
        </el-table-column>
        
        <el-table-column label="操作" width="200" fixed="right">
          <template #default="{ row }">
            <el-button size="small" @click="handleView(row)">详情</el-button>
            <el-button size="small" @click="handleEdit(row)">编辑</el-button>
            <el-button 
              v-if="row.status === 'disabled'"
              type="success"
              size="small"
              @click="handleToggleStatus(row, 'active')"
            >
              启用
            </el-button>
            <el-button 
              v-else
              type="warning"
              size="small"
              @click="handleToggleStatus(row, 'disabled')"
            >
              禁用
            </el-button>
            <el-button
              type="danger"
              size="small"
              @click="handleDelete(row)"
            >
              删除
            </el-button>
          </template>
        </el-table-column>
      </el-table>

      <!-- 分页 -->
      <div class="pagination-wrapper">
        <el-pagination
          v-model:current-page="pagination.page"
          v-model:page-size="pagination.size"
          :total="pagination.total"
          :page-sizes="[10, 20, 50, 100]"
          layout="total, sizes, prev, pager, next, jumper"
          @size-change="handleSizeChange"
          @current-change="handlePageChange"
        />
      </div>
    </el-card>

    <!-- 用户编辑对话框 -->
    <el-dialog
      v-model="userDialog.visible"
      :title="userDialog.isEdit ? '编辑用户' : '新增用户'"
      width="600px"
    >
      <el-form
        ref="userFormRef"
        :model="userDialog.form"
        :rules="userRules"
        label-width="100px"
      >
        <el-row :gutter="20">
          <el-col :span="12">
            <el-form-item label="用户名" prop="username">
              <el-input 
                v-model="userDialog.form.username" 
                placeholder="请输入用户名"
                :disabled="userDialog.isEdit"
              />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="真实姓名" prop="real_name">
              <el-input v-model="userDialog.form.real_name" placeholder="请输入真实姓名" />
            </el-form-item>
          </el-col>
        </el-row>

        <el-row :gutter="20">
          <el-col :span="12">
            <el-form-item label="手机号" prop="phone">
              <el-input v-model="userDialog.form.phone" placeholder="请输入手机号" />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="邮箱" prop="email">
              <el-input v-model="userDialog.form.email" placeholder="请输入邮箱" />
            </el-form-item>
          </el-col>
        </el-row>

        <el-row :gutter="20">
          <el-col :span="12">
            <el-form-item label="用户角色" prop="role">
              <el-select v-model="userDialog.form.role" placeholder="选择用户角色" style="width: 100%">
                <el-option label="管理员" value="admin" />
                <el-option label="审批员" value="approver" />
                <el-option label="普通用户" value="user" />
              </el-select>
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="用户状态" prop="status">
              <el-radio-group v-model="userDialog.form.status">
                <el-radio value="active">正常</el-radio>
                <el-radio value="disabled">禁用</el-radio>
                <el-radio value="pending">待审核</el-radio>
              </el-radio-group>
            </el-form-item>
          </el-col>
        </el-row>

        <el-form-item v-if="!userDialog.isEdit" label="初始密码" prop="password">
          <el-input 
            v-model="userDialog.form.password" 
            type="password" 
            placeholder="请输入初始密码"
            show-password
          />
        </el-form-item>

        <el-form-item label="备注信息">
          <el-input
            v-model="userDialog.form.notes"
            type="textarea"
            :rows="3"
            placeholder="请输入备注信息"
          />
        </el-form-item>
      </el-form>
      
      <template #footer>
        <el-button @click="userDialog.visible = false">取消</el-button>
        <el-button
          type="primary"
          :loading="userDialog.loading"
          @click="handleSaveUser"
        >
          保存
        </el-button>
      </template>
    </el-dialog>

    <!-- 用户详情对话框 -->
    <el-dialog
      v-model="detailDialog.visible"
      title="用户详情"
      width="700px"
    >
      <div v-if="detailDialog.user" class="user-detail">
        <el-descriptions :column="2" border>
          <el-descriptions-item label="用户ID">{{ detailDialog.user.id }}</el-descriptions-item>
          <el-descriptions-item label="用户名">{{ detailDialog.user.username }}</el-descriptions-item>
          <el-descriptions-item label="真实姓名">
            {{ detailDialog.user.real_name || '未设置' }}
          </el-descriptions-item>
          <el-descriptions-item label="手机号">
            {{ detailDialog.user.phone || '未设置' }}
          </el-descriptions-item>
          <el-descriptions-item label="邮箱">
            {{ detailDialog.user.email || '未设置' }}
          </el-descriptions-item>
          <el-descriptions-item label="用户角色">
            <el-tag :type="getRoleColor(detailDialog.user.role)">
              {{ getRoleName(detailDialog.user.role) }}
            </el-tag>
          </el-descriptions-item>
          <el-descriptions-item label="用户状态">
            <el-tag :type="getStatusColor(detailDialog.user.status)">
              {{ getStatusName(detailDialog.user.status) }}
            </el-tag>
          </el-descriptions-item>
          <el-descriptions-item label="注册时间">
            {{ formatDateTime(detailDialog.user.created_at) }}
          </el-descriptions-item>
          <el-descriptions-item label="最后登录">
            {{ detailDialog.user.last_login_at ? formatDateTime(detailDialog.user.last_login_at) : '从未登录' }}
          </el-descriptions-item>
          <el-descriptions-item label="登录次数">
            {{ detailDialog.user.login_count || 0 }}次
          </el-descriptions-item>
        </el-descriptions>

        <div class="detail-sections">
          <h4>近期活动记录</h4>
          <el-table :data="detailDialog.activities" size="small">
            <el-table-column prop="action" label="操作" width="120" />
            <el-table-column prop="description" label="描述" />
            <el-table-column prop="ip_address" label="IP地址" width="140" />
            <el-table-column prop="created_at" label="时间" width="160">
              <template #default="{ row }">
                {{ formatDateTime(row.created_at) }}
              </template>
            </el-table-column>
          </el-table>
        </div>
      </div>
      
      <template #footer>
        <el-button @click="detailDialog.visible = false">关闭</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { Plus, Upload, Download, Search, Refresh, Lock } from '@element-plus/icons-vue'
import { ElMessage, ElMessageBox } from 'element-plus'

// 筛选条件
const filterRole = ref('')
const filterStatus = ref('')
const searchKeyword = ref('')

// 列表数据
const loading = ref(false)
const userList = ref([])
const selectedUsers = ref([])

// 统计数据
const stats = reactive({
  total: 1245,
  active: 1189,
  disabled: 56,
  todayLogin: 89
})

// 分页
const pagination = reactive({
  page: 1,
  size: 20,
  total: 0
})

// 用户对话框
const userDialog = reactive({
  visible: false,
  loading: false,
  isEdit: false,
  form: {
    id: null,
    username: '',
    real_name: '',
    phone: '',
    email: '',
    role: 'user',
    status: 'active',
    password: '',
    notes: ''
  }
})

// 详情对话框
const detailDialog = reactive({
  visible: false,
  user: null as any,
  activities: []
})

// 表单验证规则
const userRules = {
  username: [
    { required: true, message: '请输入用户名', trigger: 'blur' },
    { min: 3, max: 20, message: '用户名长度为3-20个字符', trigger: 'blur' }
  ],
  phone: [
    { pattern: /^1[3-9]\d{9}$/, message: '请输入正确的手机号', trigger: 'blur' }
  ],
  email: [
    { type: 'email', message: '请输入正确的邮箱地址', trigger: 'blur' }
  ],
  role: [{ required: true, message: '请选择用户角色', trigger: 'change' }],
  password: [
    { required: true, message: '请输入初始密码', trigger: 'blur' },
    { min: 6, message: '密码长度不能少于6位', trigger: 'blur' }
  ]
}

// 模拟数据
const mockData = [
  {
    id: 1,
    username: 'admin001',
    real_name: '张管理',
    phone: '13800138001',
    email: 'admin@example.com',
    role: 'admin',
    status: 'active',
    created_at: new Date('2024-01-15'),
    last_login_at: new Date('2024-12-01 10:30:00'),
    login_count: 156
  },
  {
    id: 2,
    username: 'approver001',
    real_name: '李审批',
    phone: '13800138002',
    email: 'approver@example.com',
    role: 'approver',
    status: 'active',
    created_at: new Date('2024-02-20'),
    last_login_at: new Date('2024-12-01 09:15:00'),
    login_count: 89
  },
  {
    id: 3,
    username: 'user001',
    real_name: '王用户',
    phone: '13800138003',
    email: 'user@example.com',
    role: 'user',
    status: 'active',
    created_at: new Date('2024-03-10'),
    last_login_at: new Date('2024-11-28 16:45:00'),
    login_count: 23
  },
  {
    id: 4,
    username: 'disabled001',
    real_name: '禁用用户',
    phone: '13800138004',
    email: null,
    role: 'user',
    status: 'disabled',
    created_at: new Date('2024-05-15'),
    last_login_at: null,
    login_count: 0
  }
]

// 模拟活动记录
const mockActivities = [
  {
    action: '登录',
    description: '用户登录系统',
    ip_address: '192.168.1.100',
    created_at: new Date('2024-12-01 10:30:00')
  },
  {
    action: '申请审批',
    description: '提交贷款申请',
    ip_address: '192.168.1.100',
    created_at: new Date('2024-12-01 09:15:00')
  },
  {
    action: '修改密码',
    description: '用户修改登录密码',
    ip_address: '192.168.1.100',
    created_at: new Date('2024-11-30 14:20:00')
  }
]

// 格式化日期时间
const formatDateTime = (date: Date) => {
  return date.toLocaleString('zh-CN', {
    year: 'numeric',
    month: '2-digit',
    day: '2-digit',
    hour: '2-digit',
    minute: '2-digit'
  })
}

// 获取角色名称
const getRoleName = (role: string) => {
  const roles: Record<string, string> = {
    admin: '管理员',
    approver: '审批员',
    user: '普通用户'
  }
  return roles[role] || role
}

// 获取角色颜色
const getRoleColor = (role: string) => {
  const colors: Record<string, string> = {
    admin: 'danger',
    approver: 'warning',
    user: 'primary'
  }
  return colors[role] || 'default'
}

// 获取状态名称
const getStatusName = (status: string) => {
  const statuses: Record<string, string> = {
    active: '正常',
    disabled: '禁用',
    pending: '待审核'
  }
  return statuses[status] || status
}

// 获取状态颜色
const getStatusColor = (status: string) => {
  const colors: Record<string, string> = {
    active: 'success',
    disabled: 'danger',
    pending: 'warning'
  }
  return colors[status] || 'default'
}

// 加载数据
const loadData = () => {
  loading.value = true
  
  setTimeout(() => {
    userList.value = mockData as any
    pagination.total = mockData.length
    loading.value = false
  }, 500)
}

// 搜索
const handleSearch = () => {
  loadData()
}

// 刷新
const handleRefresh = () => {
  loadData()
}

// 选择变更
const handleSelectionChange = (selection: any[]) => {
  selectedUsers.value = selection
}

// 分页变更
const handlePageChange = (page: number) => {
  pagination.page = page
  loadData()
}

const handleSizeChange = (size: number) => {
  pagination.size = size
  pagination.page = 1
  loadData()
}

// 新增用户
const handleAddUser = () => {
  userDialog.visible = true
  userDialog.isEdit = false
  userDialog.form = {
    id: null,
    username: '',
    real_name: '',
    phone: '',
    email: '',
    role: 'user',
    status: 'active',
    password: '',
    notes: ''
  }
}

// 编辑用户
const handleEdit = (row: any) => {
  userDialog.visible = true
  userDialog.isEdit = true
  userDialog.form = { ...row, password: '' }
}

// 查看用户详情
const handleView = (row: any) => {
  detailDialog.user = row
  detailDialog.activities = mockActivities
  detailDialog.visible = true
}

// 保存用户
const handleSaveUser = async () => {
  userDialog.loading = true
  
  try {
    await new Promise(resolve => setTimeout(resolve, 1000))
    
    const action = userDialog.isEdit ? '更新' : '创建'
    ElMessage.success(`用户${action}成功`)
    
    userDialog.visible = false
    loadData()
  } catch (error) {
    ElMessage.error('操作失败')
  } finally {
    userDialog.loading = false
  }
}

// 切换用户状态
const handleToggleStatus = (row: any, status: string) => {
  const action = status === 'active' ? '启用' : '禁用'
  ElMessageBox.confirm(`确定要${action}用户 ${row.username} 吗？`, '确认操作', {
    confirmButtonText: '确定',
    cancelButtonText: '取消',
    type: 'warning'
  }).then(() => {
    ElMessage.success(`用户已${action}`)
    loadData()
  })
}

// 删除用户
const handleDelete = (row: any) => {
  ElMessageBox.confirm(`确定要删除用户 ${row.username} 吗？`, '确认删除', {
    confirmButtonText: '删除',
    cancelButtonText: '取消',
    type: 'danger'
  }).then(() => {
    ElMessage.success('用户已删除')
    loadData()
  })
}

// 批量禁用
const handleBatchDisable = () => {
  ElMessageBox.confirm(
    `确定要禁用选中的 ${selectedUsers.value.length} 个用户吗？`,
    '批量禁用',
    {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning'
    }
  ).then(() => {
    ElMessage.success('批量禁用操作已执行')
    loadData()
  })
}

// 导入
const handleImport = () => {
  ElMessage.info('导入功能开发中')
}

// 导出
const handleExport = () => {
  ElMessage.info('正在导出数据...')
}

onMounted(() => {
  loadData()
})
</script>

<style scoped>
.users-management {
  padding: 20px;
  background: #f5f5f5;
  min-height: calc(100vh - 60px);
}

.operation-card {
  margin-bottom: 16px;
  border: none;
}

.operation-row {
  display: flex;
  justify-content: space-between;
  align-items: center;
  gap: 16px;
}

.operation-left {
  display: flex;
  gap: 8px;
}

.operation-right {
  display: flex;
  gap: 8px;
  align-items: center;
}

.stats-row {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(200px, 1fr));
  gap: 16px;
  margin-bottom: 16px;
}

.stat-card {
  border: none;
  text-align: center;
}

.stat-content {
  padding: 8px 0;
}

.stat-number {
  font-size: 28px;
  font-weight: bold;
  color: #409EFF;
  margin-bottom: 4px;
}

.stat-label {
  font-size: 14px;
  color: #909399;
}

.list-card {
  border: none;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  font-weight: 500;
}

.header-actions {
  display: flex;
  gap: 8px;
}

.text-placeholder {
  color: #c0c4cc;
  font-style: italic;
}

.pagination-wrapper {
  margin-top: 16px;
  text-align: right;
}

.user-detail {
  padding: 20px 0;
}

.detail-sections {
  margin-top: 20px;
}

.detail-sections h4 {
  margin: 20px 0 12px 0;
  color: #303133;
  font-size: 16px;
}

@media (max-width: 768px) {
  .operation-row {
    flex-direction: column;
    align-items: stretch;
  }
  
  .operation-left,
  .operation-right {
    flex-wrap: wrap;
  }
  
  .stats-row {
    grid-template-columns: repeat(2, 1fr);
  }
}
</style> 