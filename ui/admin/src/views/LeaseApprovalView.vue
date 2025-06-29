<template>
  <div class="lease-approval-container">
    <el-card class="header-card">
      <template #header>
        <div class="card-header">
          <span>农机租赁审批</span>
          <el-button type="primary" @click="refreshData">
            <el-icon><Refresh /></el-icon>
            刷新
          </el-button>
        </div>
      </template>
      
      <div class="summary-stats">
        <el-row :gutter="20">
          <el-col :span="6">
            <div class="stat-item">
              <div class="stat-number">{{ stats.totalApplications }}</div>
              <div class="stat-label">总申请数</div>
            </div>
          </el-col>
          <el-col :span="6">
            <div class="stat-item">
              <div class="stat-number approved">{{ stats.approved }}</div>
              <div class="stat-label">已通过</div>
            </div>
          </el-col>
          <el-col :span="6">
            <div class="stat-item">
              <div class="stat-number pending">{{ stats.pending }}</div>
              <div class="stat-label">待审批</div>
            </div>
          </el-col>
          <el-col :span="6">
            <div class="stat-item">
              <div class="stat-number rejected">{{ stats.rejected }}</div>
              <div class="stat-label">已拒绝</div>
            </div>
          </el-col>
        </el-row>
      </div>
    </el-card>

    <el-card class="content-card">
      <template #header>
        <div class="card-header">
          <span>租赁申请列表</span>
          <div class="header-actions">
            <el-select v-model="selectedStatus" placeholder="审批状态" style="width: 140px" @change="handleStatusChange">
              <el-option label="全部" value="" />
              <el-option label="待审批" value="pending" />
              <el-option label="已通过" value="approved" />
              <el-option label="已拒绝" value="rejected" />
            </el-select>
            <el-select v-model="selectedMachineType" placeholder="农机类型" style="width: 140px" @change="handleMachineTypeChange">
              <el-option label="全部类型" value="" />
              <el-option label="拖拉机" value="tractor" />
              <el-option label="收割机" value="harvester" />
              <el-option label="播种机" value="seeder" />
              <el-option label="施肥机" value="fertilizer" />
              <el-option label="其他" value="other" />
            </el-select>
            <el-input
              v-model="searchKeyword"
              placeholder="请输入申请人姓名或手机号"
              style="width: 250px"
              @input="handleSearch"
            >
              <template #prefix>
                <el-icon><Search /></el-icon>
              </template>
            </el-input>
          </div>
        </div>
      </template>

      <el-table
        :data="tableData"
        v-loading="loading"
        stripe
        style="width: 100%;"
        @sort-change="handleSortChange"
        flexible
      >
        <el-table-column prop="id" label="申请ID" min-width="120" sortable="custom" />
        <el-table-column prop="applicant_name" label="申请人" min-width="100" />
        <el-table-column prop="phone" label="手机号" min-width="140" />
        <el-table-column prop="machine_type" label="农机类型" min-width="100">
          <template #default="{ row }">
            <el-tag :type="getMachineTypeTagType(row.machine_type)">
              {{ getMachineTypeText(row.machine_type) }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="machine_model" label="设备型号" min-width="140" show-overflow-tooltip />
        <el-table-column prop="lease_duration" label="租期" min-width="80">
          <template #default="{ row }">
            <span>{{ row.lease_duration }}天</span>
          </template>
        </el-table-column>
        <el-table-column prop="daily_rate" label="日租金" min-width="100" sortable="custom">
          <template #default="{ row }">
            <span class="amount">¥{{ row.daily_rate }}</span>
          </template>
        </el-table-column>
        <el-table-column prop="total_amount" label="租金总额" min-width="120" sortable="custom">
          <template #default="{ row }">
            <span class="total-amount">{{ formatAmount(row.total_amount) }}</span>
          </template>
        </el-table-column>
        <el-table-column prop="lease_area" label="作业面积" min-width="100">
          <template #default="{ row }">
            <span>{{ row.lease_area }}亩</span>
          </template>
        </el-table-column>
        <el-table-column prop="status" label="审批状态" min-width="100">
          <template #default="{ row }">
            <el-tag
              :type="getStatusTagType(row.status)"
              effect="dark"
            >
              {{ getStatusText(row.status) }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="created_at" label="申请时间" min-width="160" sortable="custom">
          <template #default="{ row }">
            {{ formatTime(row.created_at) }}
          </template>
        </el-table-column>
        <el-table-column label="操作" min-width="220" fixed="right">
          <template #default="{ row }">
            <el-button
              type="primary"
              size="small"
              @click="viewDetail(row)"
            >
              查看详情
            </el-button>
            <el-button
              v-if="row.status === 'pending'"
              type="success"
              size="small"
              @click="approveApplication(row)"
            >
              通过
            </el-button>
            <el-button
              v-if="row.status === 'pending'"
              type="danger"
              size="small"
              @click="rejectApplication(row)"
            >
              拒绝
            </el-button>
          </template>
        </el-table-column>
      </el-table>

      <div class="pagination-wrapper">
        <el-pagination
          v-model:current-page="currentPage"
          v-model:page-size="pageSize"
          :page-sizes="[10, 20, 50, 100]"
          :total="total"
          layout="total, sizes, prev, pager, next, jumper"
          @size-change="handleSizeChange"
          @current-change="handleCurrentChange"
        />
      </div>
    </el-card>

    <!-- 租赁详情对话框 -->
    <el-dialog
      v-model="detailDialogVisible"
      title="农机租赁申请详情"
      width="900px"
      destroy-on-close
    >
      <div v-if="selectedRecord" class="detail-content">
        <el-descriptions :column="2" border style="margin-bottom: 20px">
          <el-descriptions-item label="申请ID">{{ selectedRecord.id }}</el-descriptions-item>
          <el-descriptions-item label="申请人">{{ selectedRecord.applicant_name }}</el-descriptions-item>
          <el-descriptions-item label="手机号">{{ selectedRecord.phone }}</el-descriptions-item>
          <el-descriptions-item label="身份证号">{{ selectedRecord.id_card }}</el-descriptions-item>
          <el-descriptions-item label="农机类型">
            <el-tag :type="getMachineTypeTagType(selectedRecord.machine_type)">
              {{ getMachineTypeText(selectedRecord.machine_type) }}
            </el-tag>
          </el-descriptions-item>
          <el-descriptions-item label="设备型号">{{ selectedRecord.machine_model }}</el-descriptions-item>
          <el-descriptions-item label="租期">{{ selectedRecord.lease_duration }}天</el-descriptions-item>
          <el-descriptions-item label="日租金">¥{{ selectedRecord.daily_rate }}</el-descriptions-item>
          <el-descriptions-item label="租金总额">{{ formatAmount(selectedRecord.total_amount) }}</el-descriptions-item>
          <el-descriptions-item label="作业面积">{{ selectedRecord.lease_area }}亩</el-descriptions-item>
          <el-descriptions-item label="作业地址" :span="2">{{ selectedRecord.work_address }}</el-descriptions-item>
          <el-descriptions-item label="审批状态">
            <el-tag :type="getStatusTagType(selectedRecord.status)">
              {{ getStatusText(selectedRecord.status) }}
            </el-tag>
          </el-descriptions-item>
          <el-descriptions-item label="申请时间">{{ formatTime(selectedRecord.created_at) }}</el-descriptions-item>
        </el-descriptions>

        <div class="application-details">
          <h4>申请详情</h4>
          <el-card>
            <p><strong>作业需求：</strong>{{ selectedRecord.work_requirement }}</p>
            <p><strong>预期开始时间：</strong>{{ formatTime(selectedRecord.expected_start_time) }}</p>
            <p><strong>特殊要求：</strong>{{ selectedRecord.special_requirements || '无' }}</p>
          </el-card>
        </div>

        <div class="approval-actions" v-if="selectedRecord.status === 'pending'">
          <el-divider content-position="left">审批操作</el-divider>
          <el-form :model="approvalForm" label-width="100px">
            <el-form-item label="审批意见">
              <el-input
                v-model="approvalForm.comment"
                type="textarea"
                :rows="3"
                placeholder="请输入审批意见"
              />
            </el-form-item>
            <el-form-item>
              <el-button type="success" @click="submitApproval('approved')">
                通过申请
              </el-button>
              <el-button type="danger" @click="submitApproval('rejected')">
                拒绝申请
              </el-button>
            </el-form-item>
          </el-form>
        </div>
      </div>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Refresh, Search } from '@element-plus/icons-vue'
import dayjs from 'dayjs'

// 接口定义
interface LeaseApprovalRecord {
  id: string
  applicant_name: string
  phone: string
  id_card: string
  machine_type: string
  machine_model: string
  lease_duration: number
  daily_rate: number
  total_amount: number
  lease_area: number
  work_address: string
  work_requirement: string
  expected_start_time: string
  special_requirements: string
  status: string
  created_at: string
  approved_at?: string
  rejected_at?: string
  approval_comment?: string
}

// 响应式数据
const loading = ref(false)
const tableData = ref<LeaseApprovalRecord[]>([])
const currentPage = ref(1)
const pageSize = ref(20)
const total = ref(0)
const selectedStatus = ref('')
const selectedMachineType = ref('')
const searchKeyword = ref('')
const detailDialogVisible = ref(false)
const selectedRecord = ref<LeaseApprovalRecord | null>(null)

// 统计数据
const stats = reactive({
  totalApplications: 0,
  approved: 0,
  pending: 0,
  rejected: 0
})

// 审批表单
const approvalForm = reactive({
  comment: ''
})

// 方法
const refreshData = () => {
  fetchData()
  fetchStats()
}

const fetchData = async () => {
  loading.value = true
  try {
    // 模拟API调用
    await new Promise(resolve => setTimeout(resolve, 500))
    
    // 模拟数据
    const mockData = Array.from({ length: pageSize.value }, (_, index) => {
      const machineTypes = ['tractor', 'harvester', 'seeder', 'fertilizer', 'other']
      const machineType = machineTypes[Math.floor(Math.random() * machineTypes.length)]
      const leaseDuration = [3, 5, 7, 10, 15, 30][Math.floor(Math.random() * 6)]
      const dailyRate = Math.floor(Math.random() * 500) + 200
      const totalAmount = leaseDuration * dailyRate
      const statusTypes = ['pending', 'approved', 'rejected']
      const status = statusTypes[Math.floor(Math.random() * 3)]
      
      // 申请人姓名数组
      const surnames = ['张', '李', '王', '赵', '孙']
      const applicantName = `${surnames[index % surnames.length]}*`
      
      return {
        id: `LA${String(currentPage.value * 100 + index + 1).padStart(6, '0')}`,
        applicant_name: applicantName,
        phone: `138${String(Math.random() * 100000000).substring(0, 8)}`,
        id_card: `420${String(Math.random() * 100000000000000).substring(0, 15)}`,
        machine_type: machineType,
        machine_model: getMachineModel(machineType),
        lease_duration: leaseDuration,
        daily_rate: dailyRate,
        total_amount: totalAmount,
        lease_area: Math.floor(Math.random() * 100) + 20,
        work_address: `${['湖北省', '河南省', '安徽省', '江苏省'][Math.floor(Math.random() * 4)]}${['武汉市', '黄冈市', '孝感市', '荆州市'][Math.floor(Math.random() * 4)]}${['红安县', '麻城市', '英山县', '罗田县'][Math.floor(Math.random() * 4)]}某村`,
        work_requirement: getWorkRequirement(machineType),
        expected_start_time: new Date(Date.now() + Math.random() * 30 * 24 * 60 * 60 * 1000).toISOString(),
        special_requirements: ['无特殊要求', '需要专业操作员', '夜间作业', '紧急需求'][Math.floor(Math.random() * 4)],
        status: status,
        created_at: new Date(Date.now() - Math.random() * 30 * 24 * 60 * 60 * 1000).toISOString(),
        approved_at: status === 'approved' ? new Date(Date.now() - Math.random() * 5 * 24 * 60 * 60 * 1000).toISOString() : undefined,
        rejected_at: status === 'rejected' ? new Date(Date.now() - Math.random() * 5 * 24 * 60 * 60 * 1000).toISOString() : undefined,
        approval_comment: status !== 'pending' ? getApprovalComment(status) : undefined
      }
    })
    
    tableData.value = mockData
    total.value = 680
  } catch (error) {
    ElMessage.error('获取数据失败')
  } finally {
    loading.value = false
  }
}

const fetchStats = async () => {
  try {
    // 模拟统计数据
    stats.totalApplications = 680
    stats.approved = 445
    stats.pending = 156
    stats.rejected = 79
  } catch (error) {
    console.error('获取统计数据失败:', error)
  }
}

const handleStatusChange = () => {
  currentPage.value = 1
  fetchData()
}

const handleMachineTypeChange = () => {
  currentPage.value = 1
  fetchData()
}

const handleSearch = () => {
  currentPage.value = 1
  fetchData()
}

const handleSortChange = (sortInfo: any) => {
  console.log('排序变化:', sortInfo)
  fetchData()
}

const handleSizeChange = (size: number) => {
  pageSize.value = size
  currentPage.value = 1
  fetchData()
}

const handleCurrentChange = (page: number) => {
  currentPage.value = page
  fetchData()
}

const viewDetail = (row: LeaseApprovalRecord) => {
  selectedRecord.value = row
  approvalForm.comment = row.approval_comment || ''
  detailDialogVisible.value = true
}

const approveApplication = async (row: LeaseApprovalRecord) => {
  try {
    await ElMessageBox.confirm('确定要通过此租赁申请吗？', '确认审批', {
      confirmButtonText: '确定通过',
      cancelButtonText: '取消',
      type: 'success'
    })
    ElMessage.success('审批通过！')
    refreshData()
  } catch {
    // 用户取消
  }
}

const rejectApplication = async (row: LeaseApprovalRecord) => {
  try {
    await ElMessageBox.confirm('确定要拒绝此租赁申请吗？', '确认审批', {
      confirmButtonText: '确定拒绝',
      cancelButtonText: '取消',
      type: 'error'
    })
    ElMessage.warning('申请已拒绝！')
    refreshData()
  } catch {
    // 用户取消
  }
}

const submitApproval = async (decision: string) => {
  if (!approvalForm.comment.trim()) {
    ElMessage.warning('请填写审批意见')
    return
  }
  
  try {
    const actionText = decision === 'approved' ? '通过' : '拒绝'
    await ElMessageBox.confirm(`确定要${actionText}此申请吗？`, '确认审批', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: decision === 'approved' ? 'success' : 'error'
    })
    
    ElMessage.success(`申请已${actionText}！`)
    detailDialogVisible.value = false
    refreshData()
  } catch {
    // 用户取消
  }
}

// 工具方法
const formatAmount = (amount: number) => {
  return `¥${amount.toLocaleString()}`
}

const formatTime = (time: string) => {
  return dayjs(time).format('YYYY-MM-DD HH:mm:ss')
}

const getMachineTypeTagType = (type: string) => {
  switch (type) {
    case 'tractor': return 'primary'
    case 'harvester': return 'success'
    case 'seeder': return 'warning'
    case 'fertilizer': return 'info'
    default: return 'default'
  }
}

const getMachineTypeText = (type: string) => {
  switch (type) {
    case 'tractor': return '拖拉机'
    case 'harvester': return '收割机'
    case 'seeder': return '播种机'
    case 'fertilizer': return '施肥机'
    case 'other': return '其他'
    default: return '未知'
  }
}

const getMachineModel = (type: string) => {
  const models = {
    tractor: ['东方红LX904', '雷沃M504-B', '约翰迪尔5E-754'],
    harvester: ['雷沃谷神RG50', '久保田PRO688Q', '中联收获4LZ-8A'],
    seeder: ['雷沃阿波斯2BFX-12', '大华宝来2BF-16', '农哈哈2BFG-9'],
    fertilizer: ['雷沃施肥机FS150', '东风农机FS200', '久保田撒肥机SF180'],
    other: ['多功能农机A1', '通用型农机B2', '专用农机C3']
  }
  const typeModels = models[type as keyof typeof models] || models.other
  return typeModels[Math.floor(Math.random() * typeModels.length)]
}

const getWorkRequirement = (type: string) => {
  const requirements = {
    tractor: '耕地翻土，准备春耕播种',
    harvester: '小麦收割，秋季收获作业',
    seeder: '玉米播种，春季种植作业',
    fertilizer: '大田施肥，提高作物产量',
    other: '综合农田作业，多功能需求'
  }
  return requirements[type as keyof typeof requirements] || requirements.other
}

const getApprovalComment = (status: string) => {
  if (status === 'approved') {
    return ['申请材料完整，符合租赁条件', '申请人资质良好，予以通过', '租赁需求合理，同意出租'][Math.floor(Math.random() * 3)]
  } else {
    return ['申请材料不完整，需补充证明', '租赁条件不符合要求', '当前设备不可用'][Math.floor(Math.random() * 3)]
  }
}

const getStatusTagType = (status: string) => {
  switch (status) {
    case 'approved': return 'success'
    case 'rejected': return 'danger'
    case 'pending': return 'warning'
    default: return 'info'
  }
}

const getStatusText = (status: string) => {
  switch (status) {
    case 'approved': return '已通过'
    case 'rejected': return '已拒绝'
    case 'pending': return '待审批'
    default: return '未知'
  }
}

onMounted(() => {
  refreshData()
})
</script>

<style scoped>
.lease-approval-container {
  height: calc(100vh - 70px); /* 减去header高度 */
  display: flex;
  flex-direction: column;
  gap: 20px;
  padding: 20px;
  overflow: hidden;
  box-sizing: border-box;
  width: 100%; /* 确保容器占满宽度 */
  max-width: 100%; /* 防止容器超出视口 */
}

.header-card {
  flex-shrink: 0;
  width: 100%; /* 确保头部卡片占满宽度 */
}

.content-card {
  flex: 1;
  display: flex;
  flex-direction: column;
  overflow: hidden;
  min-height: 0; /* 确保flex子项能够正确收缩 */
  height: 100%; /* 确保占满可用高度 */
  width: 100%; /* 确保内容卡片占满宽度 */
  max-width: 100%; /* 防止卡片超出容器 */
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  flex-shrink: 0; /* 防止头部被压缩 */
}

.header-actions {
  display: flex;
  gap: 16px;
  align-items: center;
}

.summary-stats {
  margin-top: 16px;
}

.stat-item {
  text-align: center;
  padding: 20px;
  background: linear-gradient(135deg, #f5f7fa 0%, #c3cfe2 100%);
  border-radius: 8px;
  transition: all 0.3s ease;
}

.stat-item:hover {
  transform: translateY(-2px);
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1);
}

.stat-number {
  font-size: 28px;
  font-weight: bold;
  color: #409eff;
  margin-bottom: 8px;
}

.stat-number.approved {
  color: #67c23a;
}

.stat-number.pending {
  color: #e6a23c;
}

.stat-number.rejected {
  color: #f56c6c;
}

.stat-label {
  font-size: 14px;
  color: #666;
}

.amount {
  font-weight: 600;
  color: #e6a23c;
}

.total-amount {
  font-weight: 600;
  color: #67c23a;
}

.pagination-wrapper {
  margin-top: 20px;
  display: flex;
  justify-content: center;
  flex-shrink: 0;
  padding: 10px 0;
  border-top: 1px solid #ebeef5; /* 添加分隔线 */
  background-color: #fff; /* 确保背景色 */
}

.detail-content {
  max-height: 700px;
  overflow-y: auto;
  padding-right: 10px;
  /* 详情对话框滚动条样式 */
  scrollbar-width: thin;
  scrollbar-color: #67c23a #f1f1f1;
}

.detail-content::-webkit-scrollbar {
  width: 12px;
}

.detail-content::-webkit-scrollbar-track {
  background: #f1f1f1;
  border-radius: 6px;
}

.detail-content::-webkit-scrollbar-thumb {
  background: #67c23a;
  border-radius: 6px;
  border: 2px solid #f1f1f1;
}

.detail-content::-webkit-scrollbar-thumb:hover {
  background: #529b2e;
}

.application-details {
  margin: 20px 0;
}

.application-details h4 {
  margin-bottom: 12px;
  color: #333;
}

.application-details p {
  margin: 8px 0;
  line-height: 1.6;
}

.approval-actions {
  margin-top: 20px;
}

/* Element Plus卡片内容区域样式 */
:deep(.el-card__body) {
  flex: 1;
  display: flex;
  flex-direction: column;
  overflow: hidden;
  padding: 0 20px 20px 20px;
  min-height: 0;
  height: 100%;
  width: 100%; /* 确保卡片主体占满宽度 */
}

/* 卡片头部样式调整 */
:deep(.el-card__header) {
  padding: 20px;
  border-bottom: 1px solid #ebeef5;
  flex-shrink: 0;
  width: 100%; /* 确保卡片头部占满宽度 */
}

/* 表格容器样式 - 修复水平布局问题 */
:deep(.el-table) {
  flex: 1 !important;
  overflow: hidden;
  height: 100% !important;
  display: flex !important;
  flex-direction: column !important;
  width: 100% !important; /* 确保表格宽度占满容器 */
  max-width: 100% !important; /* 防止表格超出容器 */
}

/* 表格内部包装器 */
:deep(.el-table__inner-wrapper) {
  width: 100% !important;
  overflow: hidden;
  flex: 1;
  display: flex;
  flex-direction: column;
}

/* 表格头部 */
:deep(.el-table__header-wrapper) {
  overflow: visible;
  background-color: #fafafa;
  flex-shrink: 0;
  width: 100% !important;
}

/* 表格头部表格 */
:deep(.el-table__header) {
  width: 100% !important;
  table-layout: fixed !important; /* 使用固定布局确保列宽一致 */
  min-width: 100% !important;
}

/* 表格主体滚动区域 - 修复水平滚动 */
:deep(.el-table__body-wrapper) {
  flex: 1 !important;
  overflow: auto !important; /* 允许水平和垂直滚动 */
  height: auto !important;
  max-height: none !important;
  width: 100% !important;
  scrollbar-width: thin;
  scrollbar-color: #67c23a #f1f1f1;
}

/* 表格主体表格 */
:deep(.el-table__body) {
  width: 100% !important;
  table-layout: fixed !important; /* 使用固定布局确保列宽一致 */
  min-width: 100% !important;
}

/* 滚动条容器 */
:deep(.el-scrollbar) {
  height: 100%;
  width: 100%;
  flex: 1;
}

:deep(.el-scrollbar__wrap) {
  overflow: auto !important; /* 允许水平滚动 */
  width: 100%;
  height: 100%;
}

:deep(.el-scrollbar__view) {
  display: block !important;
  width: 100% !important;
  min-width: 100% !important; /* 确保视图占满宽度 */
}

/* 修复固定列的显示问题 */
:deep(.el-table-fixed-column--right) {
  background-color: #fff;
  position: sticky !important;
  right: 0 !important;
  z-index: 10;
}

/* 确保列宽度自适应 */
:deep(.el-table__cell) {
  padding: 8px 12px;
  word-break: break-word;
  overflow: hidden;
  text-overflow: ellipsis;
  box-sizing: border-box;
}

/* 表格列组 - 确保列宽度正确分配 */
:deep(colgroup col) {
  min-width: auto !important;
}

/* 修复表格行的显示 */
:deep(.el-table__row) {
  background-color: #fff;
  width: 100%;
}

:deep(.el-table__row--striped) {
  background-color: #fafafa;
}

/* 表格滚动条样式 */
:deep(.el-table__body-wrapper::-webkit-scrollbar) {
  width: 12px;
  height: 12px;
}

:deep(.el-table__body-wrapper::-webkit-scrollbar-track) {
  background: #f1f1f1;
  border-radius: 6px;
}

:deep(.el-table__body-wrapper::-webkit-scrollbar-thumb) {
  background: #67c23a;
  border-radius: 6px;
  border: 2px solid #f1f1f1;
}

:deep(.el-table__body-wrapper::-webkit-scrollbar-thumb:hover) {
  background: #529b2e;
}

:deep(.el-table__body-wrapper::-webkit-scrollbar-corner) {
  background: #f1f1f1;
}

/* 水平滚动条样式 */
:deep(.el-scrollbar__bar.is-horizontal) {
  height: 12px !important;
  bottom: 0 !important;
}

:deep(.el-scrollbar__bar.is-vertical) {
  width: 12px !important;
  right: 0 !important;
}

/* 响应式调整 */
@media (max-width: 1200px) {
  .lease-approval-container {
    padding: 15px;
    gap: 15px;
  }
}

@media (max-width: 768px) {
  .lease-approval-container {
    padding: 10px;
    gap: 10px;
  }
  
  .header-actions {
    flex-direction: column;
    gap: 10px;
    align-items: stretch;
  }
  
  .header-actions .el-select,
  .header-actions .el-input {
    width: 100% !important;
  }
}
</style> 