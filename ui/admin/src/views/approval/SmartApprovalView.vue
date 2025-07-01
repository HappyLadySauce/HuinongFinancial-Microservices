<template>
  <div class="smart-approval">
    <!-- 顶部功能区 -->
    <div class="header-section">
      <el-card class="stats-card" shadow="never">
        <div class="stats-grid">
          <div class="stat-item">
            <div class="stat-icon ai">
              <el-icon><MagicStick /></el-icon>
            </div>
            <div class="stat-content">
              <div class="stat-number">{{ aiStats.processed }}</div>
              <div class="stat-label">AI处理量</div>
            </div>
          </div>
          
          <div class="stat-item">
            <div class="stat-icon accuracy">
              <el-icon><TrendCharts /></el-icon>
            </div>
            <div class="stat-content">
              <div class="stat-number">{{ aiStats.accuracy }}%</div>
              <div class="stat-label">准确率</div>
            </div>
          </div>
          
          <div class="stat-item">
            <div class="stat-icon time-saved">
              <el-icon><Timer /></el-icon>
            </div>
            <div class="stat-content">
              <div class="stat-number">{{ aiStats.timeSaved }}h</div>
              <div class="stat-label">节省时间</div>
            </div>
          </div>
          
          <div class="stat-item">
            <div class="stat-icon auto-rate">
              <el-icon><Cpu /></el-icon>
            </div>
            <div class="stat-content">
              <div class="stat-number">{{ aiStats.autoRate }}%</div>
              <div class="stat-label">自动化率</div>
            </div>
          </div>
        </div>
      </el-card>
    </div>

    <!-- AI配置区域 -->
    <el-card class="config-card" shadow="never">
      <template #header>
        <div class="card-header">
          <span>AI审批配置</span>
          <el-switch
            v-model="aiEnabled"
            class="ai-switch"
            active-text="AI审批已启用"
            inactive-text="AI审批已停用"
            @change="handleAIToggle"
          />
        </div>
      </template>

      <div class="config-content">
        <el-row :gutter="20">
          <el-col :span="8">
            <el-card class="config-item" shadow="hover">
              <div class="config-header">
                <el-icon><Money /></el-icon>
                <span>自动审批阈值</span>
              </div>
              <div class="config-body">
                <el-form-item label="贷款金额上限">
                  <el-input-number
                    v-model="aiConfig.loanLimit"
                    :min="0"
                    :max="10000000"
                    :step="10000"
                    style="width: 100%"
                  />
                </el-form-item>
                <el-form-item label="租赁金额上限">
                  <el-input-number
                    v-model="aiConfig.leaseLimit"
                    :min="0"
                    :max="5000000"
                    :step="5000"
                    style="width: 100%"
                  />
                </el-form-item>
              </div>
            </el-card>
          </el-col>
          
          <el-col :span="8">
            <el-card class="config-item" shadow="hover">
              <div class="config-header">
                <el-icon><User /></el-icon>
                <span>风险评估权重</span>
              </div>
              <div class="config-body">
                <el-form-item label="信用评分权重">
                  <el-slider v-model="aiConfig.creditWeight" :max="100" show-input />
                </el-form-item>
                <el-form-item label="收入验证权重">
                  <el-slider v-model="aiConfig.incomeWeight" :max="100" show-input />
                </el-form-item>
                <el-form-item label="历史记录权重">
                  <el-slider v-model="aiConfig.historyWeight" :max="100" show-input />
                </el-form-item>
              </div>
            </el-card>
          </el-col>
          
          <el-col :span="8">
            <el-card class="config-item" shadow="hover">
              <div class="config-header">
                <el-icon><Setting /></el-icon>
                <span>审批策略</span>
              </div>
              <div class="config-body">
                <el-form-item label="自动通过阈值">
                  <el-input-number
                    v-model="aiConfig.autoApproveScore"
                    :min="0"
                    :max="100"
                    style="width: 100%"
                  />
                </el-form-item>
                <el-form-item label="自动拒绝阈值">
                  <el-input-number
                    v-model="aiConfig.autoRejectScore"
                    :min="0"
                    :max="100"
                    style="width: 100%"
                  />
                </el-form-item>
              </div>
            </el-card>
          </el-col>
        </el-row>
        
        <div class="config-actions">
          <el-button type="primary" @click="saveAIConfig">保存配置</el-button>
          <el-button @click="resetAIConfig">重置配置</el-button>
          <el-button type="warning" @click="trainAIModel">重新训练模型</el-button>
        </div>
      </div>
    </el-card>

    <!-- AI审批队列 -->
    <el-card class="queue-card" shadow="never">
      <template #header>
        <div class="card-header">
          <span>AI审批队列</span>
          <div class="header-actions">
            <el-button size="small" @click="processAll" :loading="processing">
              <el-icon><VideoPlay /></el-icon>
              批量处理
            </el-button>
            <el-button size="small" @click="loadQueue">
              <el-icon><Refresh /></el-icon>
              刷新
            </el-button>
          </div>
        </div>
      </template>

      <el-table
        v-loading="queueLoading"
        :data="queueData"
        stripe
      >
        <el-table-column prop="id" label="申请ID" width="120" />
        
        <el-table-column prop="type" label="类型" width="100">
          <template #default="{ row }">
            <el-tag :type="row.type === 'loan' ? 'primary' : 'success'" size="small">
              {{ row.type === 'loan' ? '贷款' : '租赁' }}
            </el-tag>
          </template>
        </el-table-column>
        
        <el-table-column prop="applicant" label="申请人" width="120" />
        
        <el-table-column prop="amount" label="金额" width="120">
          <template #default="{ row }">
            <span class="amount">{{ formatAmount(row.amount) }}</span>
          </template>
        </el-table-column>
        
        <el-table-column prop="ai_score" label="AI评分" width="120">
          <template #default="{ row }">
            <div class="score-display">
              <el-progress
                :percentage="row.ai_score"
                :color="getScoreColor(row.ai_score)"
                :show-text="false"
                style="margin-bottom: 4px"
              />
              <span class="score-number">{{ row.ai_score }}</span>
            </div>
          </template>
        </el-table-column>
        
        <el-table-column prop="risk_level" label="风险等级" width="100">
          <template #default="{ row }">
            <el-tag :type="getRiskColor(row.risk_level)" size="small">
              {{ getRiskText(row.risk_level) }}
            </el-tag>
          </template>
        </el-table-column>
        
        <el-table-column prop="ai_recommendation" label="AI建议" width="120">
          <template #default="{ row }">
            <el-tag :type="getRecommendationColor(row.ai_recommendation)" size="small">
              {{ getRecommendationText(row.ai_recommendation) }}
            </el-tag>
          </template>
        </el-table-column>
        
        <el-table-column prop="confidence" label="置信度" width="100">
          <template #default="{ row }">
            {{ (row.confidence * 100).toFixed(1) }}%
          </template>
        </el-table-column>
        
        <el-table-column prop="processing_time" label="处理时间" width="120">
          <template #default="{ row }">
            {{ row.processing_time }}ms
          </template>
        </el-table-column>
        
        <el-table-column label="操作" width="200" fixed="right">
          <template #default="{ row }">
            <el-button size="small" @click="viewAIAnalysis(row)">查看分析</el-button>
            <el-button
              v-if="row.ai_recommendation === 'approve'"
              type="success"
              size="small"
              @click="executeAIDecision(row, 'approve')"
            >
              执行通过
            </el-button>
            <el-button
              v-if="row.ai_recommendation === 'reject'"
              type="danger"
              size="small"
              @click="executeAIDecision(row, 'reject')"
            >
              执行拒绝
            </el-button>
            <el-button
              v-if="row.ai_recommendation === 'manual'"
              type="warning"
              size="small"
              @click="transferToManual(row)"
            >
              转人工
            </el-button>
          </template>
        </el-table-column>
      </el-table>

      <!-- 分页 -->
      <div class="pagination-wrapper">
        <el-pagination
          v-model:current-page="queuePagination.page"
          v-model:page-size="queuePagination.size"
          :total="queuePagination.total"
          :page-sizes="[10, 20, 50, 100]"
          layout="total, sizes, prev, pager, next, jumper"
          @size-change="handleQueueSizeChange"
          @current-change="handleQueuePageChange"
        />
      </div>
    </el-card>

    <!-- AI分析详情对话框 -->
    <el-dialog
      v-model="analysisDialog.visible"
      title="AI风险分析报告"
      width="800px"
    >
      <div v-if="analysisDialog.data" class="analysis-content">
        <el-descriptions :column="2" border>
          <el-descriptions-item label="申请ID">{{ analysisDialog.data.id }}</el-descriptions-item>
          <el-descriptions-item label="AI评分">{{ analysisDialog.data.ai_score }}</el-descriptions-item>
          <el-descriptions-item label="风险等级">
            <el-tag :type="getRiskColor(analysisDialog.data.risk_level)">
              {{ getRiskText(analysisDialog.data.risk_level) }}
            </el-tag>
          </el-descriptions-item>
          <el-descriptions-item label="置信度">{{ (analysisDialog.data.confidence * 100).toFixed(1) }}%</el-descriptions-item>
        </el-descriptions>

        <div class="analysis-sections">
          <h4>风险因子分析</h4>
          <div class="risk-factors">
            <div v-for="factor in analysisDialog.data.risk_factors" :key="factor.name" class="factor-item">
              <div class="factor-header">
                <span class="factor-name">{{ factor.name }}</span>
                <span class="factor-score">{{ factor.score }}</span>
              </div>
              <el-progress
                :percentage="factor.score"
                :color="getScoreColor(factor.score)"
                :show-text="false"
              />
              <div class="factor-reason">{{ factor.reason }}</div>
            </div>
          </div>

          <h4>审批建议理由</h4>
          <div class="ai-reasoning">
            {{ analysisDialog.data.reasoning }}
          </div>

          <h4>相似案例对比</h4>
          <el-table :data="analysisDialog.data.similar_cases" size="small">
            <el-table-column prop="case_id" label="案例ID" width="120" />
            <el-table-column prop="similarity" label="相似度" width="100">
              <template #default="{ row }">
                {{ (row.similarity * 100).toFixed(1) }}%
              </template>
            </el-table-column>
            <el-table-column prop="final_decision" label="最终决策" width="100">
              <template #default="{ row }">
                <el-tag :type="row.final_decision === 'approved' ? 'success' : 'danger'" size="small">
                  {{ row.final_decision === 'approved' ? '通过' : '拒绝' }}
                </el-tag>
              </template>
            </el-table-column>
            <el-table-column prop="outcome" label="结果" show-overflow-tooltip />
          </el-table>
        </div>
      </div>
      
      <template #footer>
        <el-button @click="analysisDialog.visible = false">关闭</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { 
  MagicStick, 
  TrendCharts, 
  Timer, 
  Cpu, 
  Money, 
  User, 
  Setting, 
  VideoPlay, 
  Refresh 
} from '@element-plus/icons-vue'
import { ElMessage, ElMessageBox } from 'element-plus'

// AI状态
const aiEnabled = ref(true)
const processing = ref(false)
const queueLoading = ref(false)

// AI统计数据
const aiStats = reactive({
  processed: 1256,
  accuracy: 94.8,
  timeSaved: 428,
  autoRate: 87.3
})

// AI配置
const aiConfig = reactive({
  loanLimit: 500000,
  leaseLimit: 200000,
  creditWeight: 40,
  incomeWeight: 35,
  historyWeight: 25,
  autoApproveScore: 85,
  autoRejectScore: 30
})

// 队列数据
const queueData = ref([])
const queuePagination = reactive({
  page: 1,
  size: 20,
  total: 0
})

// 分析对话框
const analysisDialog = reactive({
  visible: false,
  data: null as any
})

// 模拟队列数据
const mockQueueData = [
  {
    id: 'LN202412005',
    type: 'loan',
    applicant: '张三',
    amount: 300000,
    ai_score: 78,
    risk_level: 'medium',
    ai_recommendation: 'manual',
    confidence: 0.76,
    processing_time: 1240,
    risk_factors: [
      { name: '信用评分', score: 75, reason: '信用记录良好，但历史较短' },
      { name: '收入验证', score: 82, reason: '收入稳定，符合要求' },
      { name: '负债比率', score: 68, reason: '负债比率略高，需关注' }
    ],
    reasoning: '申请人信用记录良好，收入稳定，但负债比率偏高，建议人工审核确认还款能力。',
    similar_cases: [
      { case_id: 'LN202411089', similarity: 0.89, final_decision: 'approved', outcome: '按时还款' },
      { case_id: 'LN202411156', similarity: 0.82, final_decision: 'approved', outcome: '轻微逾期' }
    ]
  },
  {
    id: 'LS202412008',
    type: 'lease',
    applicant: '李四',
    amount: 150000,
    ai_score: 92,
    risk_level: 'low',
    ai_recommendation: 'approve',
    confidence: 0.94,
    processing_time: 890,
    risk_factors: [
      { name: '信用评分', score: 95, reason: '优秀的信用记录' },
      { name: '资产证明', score: 88, reason: '资产充足，风险较低' },
      { name: '行业风险', score: 90, reason: '所属行业稳定' }
    ],
    reasoning: '申请人信用优秀，资产充足，租赁金额合理，建议自动通过。',
    similar_cases: [
      { case_id: 'LS202411234', similarity: 0.91, final_decision: 'approved', outcome: '合同顺利执行' },
      { case_id: 'LS202411445', similarity: 0.87, final_decision: 'approved', outcome: '提前还款' }
    ]
  },
  {
    id: 'LN202412006',
    type: 'loan',
    applicant: '王五',
    amount: 800000,
    ai_score: 25,
    risk_level: 'high',
    ai_recommendation: 'reject',
    confidence: 0.88,
    processing_time: 1580,
    risk_factors: [
      { name: '信用评分', score: 35, reason: '存在多次逾期记录' },
      { name: '收入验证', score: 28, reason: '收入不稳定' },
      { name: '负债比率', score: 15, reason: '负债比率过高' }
    ],
    reasoning: '申请人信用记录较差，收入不稳定，负债比率过高，违约风险极大，建议拒绝。',
    similar_cases: [
      { case_id: 'LN202411267', similarity: 0.85, final_decision: 'rejected', outcome: '避免损失' },
      { case_id: 'LN202411389', similarity: 0.79, final_decision: 'rejected', outcome: '风险控制' }
    ]
  }
]

// 格式化金额
const formatAmount = (amount: number) => {
  if (amount >= 10000) {
    return `¥${(amount / 10000).toFixed(1)}万`
  }
  return `¥${amount.toLocaleString()}`
}

// 获取评分颜色
const getScoreColor = (score: number) => {
  if (score >= 80) return '#67C23A'
  if (score >= 60) return '#E6A23C'
  return '#F56C6C'
}

// 获取风险等级颜色
const getRiskColor = (level: string) => {
  const colors: Record<string, string> = {
    low: 'success',
    medium: 'warning',
    high: 'danger'
  }
  return colors[level] || 'info'
}

// 获取风险等级文本
const getRiskText = (level: string) => {
  const texts: Record<string, string> = {
    low: '低风险',
    medium: '中风险',
    high: '高风险'
  }
  return texts[level] || level
}

// 获取建议颜色
const getRecommendationColor = (recommendation: string) => {
  const colors: Record<string, string> = {
    approve: 'success',
    reject: 'danger',
    manual: 'warning'
  }
  return colors[recommendation] || 'info'
}

// 获取建议文本
const getRecommendationText = (recommendation: string) => {
  const texts: Record<string, string> = {
    approve: '建议通过',
    reject: '建议拒绝',
    manual: '转人工审核'
  }
  return texts[recommendation] || recommendation
}

// 加载队列数据
const loadQueue = () => {
  queueLoading.value = true
  
  setTimeout(() => {
    queueData.value = mockQueueData as any
    queuePagination.total = mockQueueData.length
    queueLoading.value = false
  }, 500)
}

// AI开关切换
const handleAIToggle = (value: boolean) => {
  ElMessage.success(`AI审批已${value ? '启用' : '停用'}`)
}

// 保存AI配置
const saveAIConfig = () => {
  ElMessage.success('AI配置已保存')
}

// 重置AI配置
const resetAIConfig = () => {
  ElMessageBox.confirm('确定要重置AI配置吗？', '确认操作', {
    confirmButtonText: '确定',
    cancelButtonText: '取消',
    type: 'warning'
  }).then(() => {
    ElMessage.success('AI配置已重置')
  })
}

// 训练AI模型
const trainAIModel = () => {
  ElMessageBox.confirm('重新训练模型将需要较长时间，确定继续吗？', '训练模型', {
    confirmButtonText: '开始训练',
    cancelButtonText: '取消',
    type: 'warning'
  }).then(() => {
    ElMessage.success('模型训练已开始，请稍后查看进度')
  })
}

// 批量处理
const processAll = async () => {
  processing.value = true
  
  try {
    await new Promise(resolve => setTimeout(resolve, 2000))
    ElMessage.success('批量处理完成')
    loadQueue()
  } finally {
    processing.value = false
  }
}

// 查看AI分析
const viewAIAnalysis = (row: any) => {
  analysisDialog.data = row
  analysisDialog.visible = true
}

// 执行AI决策
const executeAIDecision = (row: any, decision: string) => {
  const action = decision === 'approve' ? '通过' : '拒绝'
  ElMessageBox.confirm(
    `确定要执行AI建议，${action}申请 ${row.id} 吗？`,
    '执行决策',
    {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: decision === 'approve' ? 'success' : 'warning'
    }
  ).then(() => {
    ElMessage.success(`已执行${action}决策`)
    loadQueue()
  })
}

// 转人工审核
const transferToManual = (row: any) => {
  ElMessage.success(`申请 ${row.id} 已转入人工审核队列`)
}

// 分页处理
const handleQueuePageChange = (page: number) => {
  queuePagination.page = page
  loadQueue()
}

const handleQueueSizeChange = (size: number) => {
  queuePagination.size = size
  queuePagination.page = 1
  loadQueue()
}

onMounted(() => {
  loadQueue()
})
</script>

<style scoped>
.smart-approval {
  padding: 20px;
  background: #f5f5f5;
  min-height: calc(100vh - 60px);
}

.header-section {
  margin-bottom: 20px;
}

.stats-card {
  border: none;
}

.stats-grid {
  display: grid;
  grid-template-columns: repeat(4, 1fr);
  gap: 20px;
}

.stat-item {
  display: flex;
  align-items: center;
  gap: 16px;
  padding: 16px;
  background: #fff;
  border-radius: 8px;
  box-shadow: 0 2px 4px rgba(0,0,0,0.1);
}

.stat-icon {
  width: 48px;
  height: 48px;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 24px;
  color: white;
}

.stat-icon.ai { background: linear-gradient(135deg, #667eea 0%, #764ba2 100%); }
.stat-icon.accuracy { background: linear-gradient(135deg, #f093fb 0%, #f5576c 100%); }
.stat-icon.time-saved { background: linear-gradient(135deg, #4facfe 0%, #00f2fe 100%); }
.stat-icon.auto-rate { background: linear-gradient(135deg, #43e97b 0%, #38f9d7 100%); }

.stat-number {
  font-size: 28px;
  font-weight: bold;
  color: #303133;
  margin-bottom: 4px;
}

.stat-label {
  font-size: 14px;
  color: #909399;
}

.config-card, .queue-card {
  border: none;
  margin-bottom: 20px;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  font-weight: 500;
}

.ai-switch {
  margin-left: 16px;
}

.config-content {
  padding: 20px 0;
}

.config-item {
  height: 280px;
}

.config-header {
  display: flex;
  align-items: center;
  gap: 8px;
  margin-bottom: 16px;
  font-weight: 500;
  color: #303133;
}

.config-body {
  padding: 0 16px;
}

.config-actions {
  text-align: center;
  margin-top: 20px;
  padding-top: 20px;
  border-top: 1px solid #EBEEF5;
}

.header-actions {
  display: flex;
  gap: 8px;
}

.score-display {
  display: flex;
  flex-direction: column;
  align-items: center;
}

.score-number {
  font-size: 12px;
  font-weight: 500;
  color: #303133;
}

.amount {
  font-weight: 500;
  color: #E6A23C;
}

.pagination-wrapper {
  margin-top: 16px;
  text-align: right;
}

.analysis-content {
  padding: 20px 0;
}

.analysis-sections {
  margin-top: 20px;
}

.analysis-sections h4 {
  margin: 20px 0 12px 0;
  color: #303133;
  font-size: 16px;
}

.risk-factors {
  display: grid;
  gap: 16px;
}

.factor-item {
  padding: 16px;
  background: #f8f9fa;
  border-radius: 8px;
}

.factor-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 8px;
}

.factor-name {
  font-weight: 500;
  color: #303133;
}

.factor-score {
  font-weight: bold;
  color: #409EFF;
}

.factor-reason {
  margin-top: 8px;
  font-size: 12px;
  color: #909399;
}

.ai-reasoning {
  padding: 16px;
  background: #f0f9ff;
  border-radius: 8px;
  color: #303133;
  line-height: 1.6;
}

@media (max-width: 768px) {
  .stats-grid {
    grid-template-columns: repeat(2, 1fr);
  }
  
  .config-content .el-row {
    flex-direction: column;
  }
  
  .config-item {
    height: auto;
    margin-bottom: 16px;
  }
}
</style> 