<template>
  <div class="products-management">
    <!-- 顶部操作栏 -->
    <el-card class="operation-card" shadow="never">
      <div class="operation-row">
        <div class="operation-left">
          <el-button type="primary" @click="handleAddProduct">
            <el-icon><Plus /></el-icon>
            新增产品
          </el-button>
          <el-button @click="handleImport">
            <el-icon><Upload /></el-icon>
            批量导入
          </el-button>
          <el-button @click="handleExport">
            <el-icon><Download /></el-icon>
            导出数据
          </el-button>
        </div>
        
        <div class="operation-right">
          <el-select v-model="filterCategory" placeholder="产品类别" clearable style="width: 140px">
            <el-option label="全部类别" value="" />
            <el-option label="信贷产品" value="loan" />
            <el-option label="租赁产品" value="lease" />
            <el-option label="保险产品" value="insurance" />
          </el-select>
          <el-select v-model="filterStatus" placeholder="产品状态" clearable style="width: 120px">
            <el-option label="全部状态" value="" />
            <el-option label="上架中" value="active" />
            <el-option label="已下架" value="inactive" />
            <el-option label="待审核" value="pending" />
          </el-select>
          <el-input
            v-model="searchKeyword"
            placeholder="搜索产品名称"
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

    <!-- 产品列表 -->
    <el-card class="list-card" shadow="never">
      <template #header>
        <div class="card-header">
          <span>产品列表</span>
          <div class="header-stats">
            <el-tag class="stat-tag">总计: {{ pagination.total }}</el-tag>
            <el-tag type="success" class="stat-tag">上架: {{ stats.active }}</el-tag>
            <el-tag type="warning" class="stat-tag">下架: {{ stats.inactive }}</el-tag>
          </div>
        </div>
      </template>

      <el-table
        v-loading="loading"
        :data="productList"
        stripe
        @selection-change="handleSelectionChange"
      >
        <el-table-column type="selection" width="55" />
        
        <el-table-column prop="id" label="产品ID" width="100" />
        
        <el-table-column prop="name" label="产品名称" width="200" show-overflow-tooltip />
        
        <el-table-column prop="category" label="产品类别" width="120">
          <template #default="{ row }">
            <el-tag :type="getCategoryColor(row.category)" size="small">
              {{ getCategoryName(row.category) }}
            </el-tag>
          </template>
        </el-table-column>
        
        <el-table-column prop="interest_rate" label="利率范围" width="120">
          <template #default="{ row }">
            {{ (row.min_rate * 100).toFixed(2) }}% - {{ (row.max_rate * 100).toFixed(2) }}%
          </template>
        </el-table-column>
        
        <el-table-column prop="amount_range" label="金额范围" width="180">
          <template #default="{ row }">
            {{ formatAmount(row.min_amount) }} - {{ formatAmount(row.max_amount) }}
          </template>
        </el-table-column>
        
        <el-table-column prop="term_range" label="期限范围" width="120">
          <template #default="{ row }">
            {{ row.min_term }} - {{ row.max_term }}个月
          </template>
        </el-table-column>
        
        <el-table-column prop="status" label="状态" width="100">
          <template #default="{ row }">
            <el-tag :type="getStatusColor(row.status)" size="small">
              {{ getStatusName(row.status) }}
            </el-tag>
          </template>
        </el-table-column>
        
        <el-table-column prop="applications_count" label="申请量" width="100" sortable />
        
        <el-table-column prop="approval_rate" label="通过率" width="100">
          <template #default="{ row }">
            {{ (row.approval_rate * 100).toFixed(1) }}%
          </template>
        </el-table-column>
        
        <el-table-column prop="created_at" label="创建时间" width="160">
          <template #default="{ row }">
            {{ formatDateTime(row.created_at) }}
          </template>
        </el-table-column>
        
        <el-table-column label="操作" width="220" fixed="right">
          <template #default="{ row }">
            <el-button size="small" @click="handleView(row)">查看</el-button>
            <el-button size="small" @click="handleEdit(row)">编辑</el-button>
            <el-button 
              v-if="row.status === 'inactive'"
              type="success"
              size="small"
              @click="handleToggleStatus(row, 'active')"
            >
              上架
            </el-button>
            <el-button 
              v-else
              type="warning"
              size="small"
              @click="handleToggleStatus(row, 'inactive')"
            >
              下架
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

    <!-- 产品编辑对话框 -->
    <el-dialog
      v-model="productDialog.visible"
      :title="productDialog.isEdit ? '编辑产品' : '新增产品'"
      width="800px"
    >
      <el-form
        ref="productFormRef"
        :model="productDialog.form"
        :rules="productRules"
        label-width="120px"
      >
        <el-row :gutter="20">
          <el-col :span="12">
            <el-form-item label="产品名称" prop="name">
              <el-input v-model="productDialog.form.name" placeholder="请输入产品名称" />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="产品类别" prop="category">
              <el-select v-model="productDialog.form.category" placeholder="选择产品类别" style="width: 100%">
                <el-option label="信贷产品" value="loan" />
                <el-option label="租赁产品" value="lease" />
                <el-option label="保险产品" value="insurance" />
              </el-select>
            </el-form-item>
          </el-col>
        </el-row>

        <el-form-item label="产品描述" prop="description">
          <el-input
            v-model="productDialog.form.description"
            type="textarea"
            :rows="3"
            placeholder="请输入产品描述"
          />
        </el-form-item>

        <el-row :gutter="20">
          <el-col :span="12">
            <el-form-item label="最低利率" prop="min_rate">
              <el-input-number
                v-model="productDialog.form.min_rate"
                :precision="4"
                :step="0.0001"
                :max="1"
                style="width: 100%"
              />
              <span class="input-suffix">%</span>
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="最高利率" prop="max_rate">
              <el-input-number
                v-model="productDialog.form.max_rate"
                :precision="4"
                :step="0.0001"
                :max="1"
                style="width: 100%"
              />
              <span class="input-suffix">%</span>
            </el-form-item>
          </el-col>
        </el-row>

        <el-row :gutter="20">
          <el-col :span="12">
            <el-form-item label="最低金额" prop="min_amount">
              <el-input-number
                v-model="productDialog.form.min_amount"
                :min="0"
                :step="1000"
                style="width: 100%"
              />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="最高金额" prop="max_amount">
              <el-input-number
                v-model="productDialog.form.max_amount"
                :min="0"
                :step="10000"
                style="width: 100%"
              />
            </el-form-item>
          </el-col>
        </el-row>

        <el-row :gutter="20">
          <el-col :span="12">
            <el-form-item label="最短期限" prop="min_term">
              <el-input-number
                v-model="productDialog.form.min_term"
                :min="1"
                style="width: 100%"
              />
              <span class="input-suffix">个月</span>
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="最长期限" prop="max_term">
              <el-input-number
                v-model="productDialog.form.max_term"
                :min="1"
                style="width: 100%"
              />
              <span class="input-suffix">个月</span>
            </el-form-item>
          </el-col>
        </el-row>

        <el-form-item label="申请条件" prop="requirements">
          <el-input
            v-model="productDialog.form.requirements"
            type="textarea"
            :rows="3"
            placeholder="请输入申请条件"
          />
        </el-form-item>

        <el-row :gutter="20">
          <el-col :span="12">
            <el-form-item label="产品状态" prop="status">
              <el-radio-group v-model="productDialog.form.status">
                <el-radio value="active">上架</el-radio>
                <el-radio value="inactive">下架</el-radio>
                <el-radio value="pending">待审核</el-radio>
              </el-radio-group>
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="排序权重" prop="sort_order">
              <el-input-number
                v-model="productDialog.form.sort_order"
                :min="0"
                style="width: 100%"
              />
            </el-form-item>
          </el-col>
        </el-row>
      </el-form>
      
      <template #footer>
        <el-button @click="productDialog.visible = false">取消</el-button>
        <el-button
          type="primary"
          :loading="productDialog.loading"
          @click="handleSaveProduct"
        >
          保存
        </el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { Plus, Upload, Download, Search } from '@element-plus/icons-vue'
import { ElMessage, ElMessageBox } from 'element-plus'

// 筛选条件
const filterCategory = ref('')
const filterStatus = ref('')
const searchKeyword = ref('')

// 列表数据
const loading = ref(false)
const productList = ref([])
const selectedProducts = ref([])

// 统计数据
const stats = reactive({
  active: 12,
  inactive: 3,
  pending: 2
})

// 分页
const pagination = reactive({
  page: 1,
  size: 20,
  total: 0
})

// 产品对话框
const productDialog = reactive({
  visible: false,
  loading: false,
  isEdit: false,
  form: {
    id: null,
    name: '',
    category: '',
    description: '',
    min_rate: 0.05,
    max_rate: 0.15,
    min_amount: 10000,
    max_amount: 1000000,
    min_term: 3,
    max_term: 36,
    requirements: '',
    status: 'active',
    sort_order: 0
  }
})

// 表单验证规则
const productRules = {
  name: [{ required: true, message: '请输入产品名称', trigger: 'blur' }],
  category: [{ required: true, message: '请选择产品类别', trigger: 'change' }],
  description: [{ required: true, message: '请输入产品描述', trigger: 'blur' }]
}

// 模拟数据
const mockData = [
  {
    id: 1,
    name: '惠农快贷',
    category: 'loan',
    description: '专为农户设计的快速信贷产品',
    min_rate: 0.055,
    max_rate: 0.085,
    min_amount: 50000,
    max_amount: 500000,
    min_term: 6,
    max_term: 24,
    status: 'active',
    applications_count: 256,
    approval_rate: 0.847,
    created_at: new Date('2024-10-15')
  },
  {
    id: 2,
    name: '农机租赁',
    category: 'lease',
    description: '农业机械设备融资租赁服务',
    min_rate: 0.065,
    max_rate: 0.095,
    min_amount: 100000,
    max_amount: 2000000,
    min_term: 12,
    max_term: 60,
    status: 'active',
    applications_count: 189,
    approval_rate: 0.912,
    created_at: new Date('2024-09-20')
  },
  {
    id: 3,
    name: '农业保险',
    category: 'insurance',
    description: '农作物种植保险产品',
    min_rate: 0.02,
    max_rate: 0.05,
    min_amount: 5000,
    max_amount: 100000,
    min_term: 12,
    max_term: 12,
    status: 'inactive',
    applications_count: 67,
    approval_rate: 0.956,
    created_at: new Date('2024-08-10')
  }
]

// 格式化金额
const formatAmount = (amount: number) => {
  if (amount >= 10000) {
    return `¥${(amount / 10000).toFixed(1)}万`
  }
  return `¥${amount.toLocaleString()}`
}

// 格式化日期时间
const formatDateTime = (date: Date) => {
  return date.toLocaleDateString('zh-CN')
}

// 获取类别名称
const getCategoryName = (category: string) => {
  const categories: Record<string, string> = {
    loan: '信贷产品',
    lease: '租赁产品',
    insurance: '保险产品'
  }
  return categories[category] || category
}

// 获取类别颜色
const getCategoryColor = (category: string) => {
  const colors: Record<string, string> = {
    loan: 'primary',
    lease: 'success',
    insurance: 'warning'
  }
  return colors[category] || 'default'
}

// 获取状态名称
const getStatusName = (status: string) => {
  const statuses: Record<string, string> = {
    active: '上架中',
    inactive: '已下架',
    pending: '待审核'
  }
  return statuses[status] || status
}

// 获取状态颜色
const getStatusColor = (status: string) => {
  const colors: Record<string, string> = {
    active: 'success',
    inactive: 'info',
    pending: 'warning'
  }
  return colors[status] || 'default'
}

// 加载数据
const loadData = () => {
  loading.value = true
  
  setTimeout(() => {
    productList.value = mockData as any
    pagination.total = mockData.length
    loading.value = false
  }, 500)
}

// 搜索
const handleSearch = () => {
  loadData()
}

// 选择变更
const handleSelectionChange = (selection: any[]) => {
  selectedProducts.value = selection
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

// 新增产品
const handleAddProduct = () => {
  productDialog.visible = true
  productDialog.isEdit = false
  productDialog.form = {
    id: null,
    name: '',
    category: '',
    description: '',
    min_rate: 0.05,
    max_rate: 0.15,
    min_amount: 10000,
    max_amount: 1000000,
    min_term: 3,
    max_term: 36,
    requirements: '',
    status: 'active',
    sort_order: 0
  }
}

// 编辑产品
const handleEdit = (row: any) => {
  productDialog.visible = true
  productDialog.isEdit = true
  productDialog.form = { ...row }
}

// 查看产品
const handleView = (row: any) => {
  ElMessage.info(`查看产品：${row.name}`)
}

// 保存产品
const handleSaveProduct = async () => {
  productDialog.loading = true
  
  try {
    await new Promise(resolve => setTimeout(resolve, 1000))
    
    const action = productDialog.isEdit ? '更新' : '创建'
    ElMessage.success(`产品${action}成功`)
    
    productDialog.visible = false
    loadData()
  } catch (error) {
    ElMessage.error('操作失败')
  } finally {
    productDialog.loading = false
  }
}

// 切换状态
const handleToggleStatus = (row: any, status: string) => {
  const action = status === 'active' ? '上架' : '下架'
  ElMessageBox.confirm(`确定要${action}产品 ${row.name} 吗？`, '确认操作', {
    confirmButtonText: '确定',
    cancelButtonText: '取消',
    type: 'warning'
  }).then(() => {
    ElMessage.success(`产品已${action}`)
    loadData()
  })
}

// 删除产品
const handleDelete = (row: any) => {
  ElMessageBox.confirm(`确定要删除产品 ${row.name} 吗？`, '确认删除', {
    confirmButtonText: '删除',
    cancelButtonText: '取消',
    type: 'danger'
  }).then(() => {
    ElMessage.success('产品已删除')
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
.products-management {
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

.list-card {
  border: none;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  font-weight: 500;
}

.header-stats {
  display: flex;
  gap: 8px;
}

.stat-tag {
  margin: 0;
}

.pagination-wrapper {
  margin-top: 16px;
  text-align: right;
}

.input-suffix {
  margin-left: 8px;
  color: #909399;
  font-size: 14px;
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
}
</style> 