<template>
  <div class="smart-approval-container">
    <el-card class="header-card">
      <template #header>
        <div class="card-header">
          <span>智能审批</span>
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
              <div class="stat-number smart">{{ stats.smartApproved }}</div>
              <div class="stat-label">智能审批通过</div>
            </div>
          </el-col>
          <el-col :span="6">
            <div class="stat-item">
              <div class="stat-number pending">{{ stats.pending }}</div>
              <div class="stat-label">待人工审批</div>
            </div>
          </el-col>
          <el-col :span="6">
            <div class="stat-item">
              <div class="stat-number rejected">{{ stats.rejected }}</div>
              <div class="stat-label">智能拒绝</div>
            </div>
          </el-col>
        </el-row>
      </div>
    </el-card>

    <el-card class="content-card">
      <template #header>
        <div class="card-header">
          <span>智能审批列表</span>
          <div class="header-actions">
            <el-select v-model="selectedStatus" placeholder="请选择状态" style="width: 160px" @change="handleStatusChange">
              <el-option label="全部" value="" />
              <el-option label="智能通过" value="approved" />
              <el-option label="智能拒绝" value="rejected" />
              <el-option label="转人工审批" value="manual" />
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
        style="width: 100%"
        @sort-change="handleSortChange"
      >
        <el-table-column prop="id" label="申请ID" width="100" sortable="custom" />
        <el-table-column prop="applicant_name" label="申请人" width="120" />
        <el-table-column prop="phone" label="手机号" width="140" />
        <el-table-column prop="loan_type" label="贷款类型" width="120" />
        <el-table-column prop="amount" label="申请金额" width="120" sortable="custom">
          <template #default="{ row }">
            <span class="amount">{{ formatAmount(row.amount) }}</span>
          </template>
        </el-table-column>
        <el-table-column prop="ai_confidence" label="AI置信度" width="100" sortable="custom">
          <template #default="{ row }">
            <el-progress
              :percentage="row.ai_confidence"
              :color="getConfidenceColor(row.ai_confidence)"
              :stroke-width="8"
              :show-text="false"
            />
            <span class="confidence-text">{{ row.ai_confidence }}%</span>
          </template>
        </el-table-column>
        <el-table-column prop="ai_result" label="AI审批结果" width="120">
          <template #default="{ row }">
            <el-tag
              :type="getResultTagType(row.ai_result)"
              effect="dark"
            >
              {{ getResultText(row.ai_result) }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="ai_reason" label="AI审批原因" min-width="200" show-overflow-tooltip />
        <el-table-column prop="created_at" label="申请时间" width="160" sortable="custom">
          <template #default="{ row }">
            {{ formatTime(row.created_at) }}
          </template>
        </el-table-column>
        <el-table-column label="操作" width="200" fixed="right">
          <template #default="{ row }">
            <el-button
              type="primary"
              size="small"
              @click="viewDetail(row)"
            >
              查看详情
            </el-button>
            <el-button
              v-if="row.ai_result === 'manual'"
              type="warning"
              size="small"
              @click="manualReview(row)"
            >
              人工审批
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

    <!-- AI审批详情对话框 -->
    <el-dialog
      v-model="detailDialogVisible"
      title="智能审批详情"
      width="800px"
      destroy-on-close
    >
      <div v-if="selectedRecord" class="detail-content">
        <el-descriptions :column="2" border>
          <el-descriptions-item label="申请ID">{{ selectedRecord.id }}</el-descriptions-item>
          <el-descriptions-item label="申请人">{{ selectedRecord.applicant_name }}</el-descriptions-item>
          <el-descriptions-item label="手机号">{{ selectedRecord.phone }}</el-descriptions-item>
          <el-descriptions-item label="身份证号">{{ selectedRecord.id_card }}</el-descriptions-item>
          <el-descriptions-item label="贷款类型">{{ selectedRecord.loan_type }}</el-descriptions-item>
          <el-descriptions-item label="申请金额">{{ formatAmount(selectedRecord.amount) }}</el-descriptions-item>
          <el-descriptions-item label="AI置信度">{{ selectedRecord.ai_confidence }}%</el-descriptions-item>
          <el-descriptions-item label="AI审批结果">
            <el-tag :type="getResultTagType(selectedRecord.ai_result)">
              {{ getResultText(selectedRecord.ai_result) }}
            </el-tag>
          </el-descriptions-item>
          <el-descriptions-item label="AI审批原因" :span="2">{{ selectedRecord.ai_reason }}</el-descriptions-item>
          <el-descriptions-item label="申请时间" :span="2">{{ formatTime(selectedRecord.created_at) }}</el-descriptions-item>
        </el-descriptions>

        <div class="ai-analysis" v-if="selectedRecord.ai_analysis">
          <h4>AI分析详情</h4>
          <div class="analysis-charts">
            <el-row :gutter="20">
              <el-col :span="8">
                <div class="chart-card">
                  <div class="chart-title">个人信用值</div>
                  <div ref="creditValueChartRef" class="analysis-chart"></div>
                </div>
              </el-col>
              <el-col :span="8">
                <div class="chart-card">
                  <div class="chart-title">资产信息</div>
                  <div ref="assetInfoChartRef" class="analysis-chart"></div>
                </div>
              </el-col>
              <el-col :span="8">
                <div class="chart-card">
                  <div class="chart-title">守约记录</div>
                  <div ref="complianceRecordChartRef" class="analysis-chart"></div>
                </div>
              </el-col>
            </el-row>
          </div>
        </div>
      </div>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted, nextTick } from 'vue'
import { ElMessage } from 'element-plus'
import { Refresh, Search } from '@element-plus/icons-vue'
import dayjs from 'dayjs'
import * as echarts from 'echarts'

// 接口定义
interface SmartApprovalRecord {
  id: string
  applicant_name: string
  phone: string
  id_card: string
  loan_type: string
  amount: number
  ai_confidence: number
  ai_result: string
  ai_reason: string
  ai_analysis: string
  created_at: string
}

// 响应式数据
const loading = ref(false)
const tableData = ref<SmartApprovalRecord[]>([])
const currentPage = ref(1)
const pageSize = ref(20)
const total = ref(0)
const selectedStatus = ref('')
const searchKeyword = ref('')
const detailDialogVisible = ref(false)
const selectedRecord = ref<SmartApprovalRecord | null>(null)

// 图表DOM引用
const assetInfoChartRef = ref<HTMLElement | null>(null)
const creditValueChartRef = ref<HTMLElement | null>(null)
const complianceRecordChartRef = ref<HTMLElement | null>(null)

// 图表实例
let assetInfoChart: echarts.ECharts | null = null
let creditValueChart: echarts.ECharts | null = null
let complianceRecordChart: echarts.ECharts | null = null

// 统计数据
const stats = reactive({
  totalApplications: 0,
  smartApproved: 0,
  pending: 0,
  rejected: 0
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
      // 先随机生成审批结果
      const resultTypes = ['approved', 'rejected', 'manual']
      const aiResult = resultTypes[Math.floor(Math.random() * 3)]
      
      // 根据审批结果选择对应的原因和置信度
      let aiReason = ''
      let aiConfidence = 0
      
      switch (aiResult) {
        case 'approved':
          aiReason = [
            '申请人信用记录良好，符合贷款条件',
            '申请材料完整，风险评估通过',
            '收入稳定，还款能力充足',
            '信用评分优秀，自动通过审批'
          ][Math.floor(Math.random() * 4)]
          aiConfidence = Math.floor(Math.random() * 20) + 80 // 80-100%
          break
        case 'rejected':
          aiReason = [
            '申请人负债率过高，不符合贷款条件',
            '信用记录存在不良记录，拒绝申请',
            '收入证明不充分，风险过高',
            '申请金额超出风险承受范围'
          ][Math.floor(Math.random() * 4)]
          aiConfidence = Math.floor(Math.random() * 25) + 75 // 75-100%
          break
        case 'manual':
          aiReason = [
            '申请人收入不稳定，建议人工审核',
            '申请材料存在疑问，需人工核实',
            'AI置信度不足，转人工审批',
            '申请金额较大，建议人工复审'
          ][Math.floor(Math.random() * 4)]
          aiConfidence = Math.floor(Math.random() * 20) + 60 // 60-80%
          break
      }
      
      // 申请人姓名数组
      const surnames = ['张', '李', '王', '赵', '孙']
      const applicantName = `${surnames[index % surnames.length]}*`
      
      return {
        id: `APP${String(currentPage.value * 100 + index + 1).padStart(6, '0')}`,
        applicant_name: applicantName,
        phone: `138${String(Math.random() * 100000000).substring(0, 8)}`,
        id_card: `420${String(Math.random() * 100000000000000).substring(0, 15)}`,
        loan_type: ['惠农贷', '农机贷', '种植贷', '养殖贷'][Math.floor(Math.random() * 4)],
        amount: Math.floor(Math.random() * 500000) + 10000,
        ai_confidence: aiConfidence,
        ai_result: aiResult,
        ai_reason: aiReason,
        ai_analysis: `AI风险评估详细分析报告...
        
基于以下因素进行综合评估：
- 申请人基本信息完整性: ${Math.floor(Math.random() * 40) + 60}%
- 信用历史记录评分: ${Math.floor(Math.random() * 40) + 60}%
- 收入稳定性评估: ${Math.floor(Math.random() * 40) + 60}%
- 负债比率分析: ${Math.floor(Math.random() * 40) + 60}%
- 行业风险评估: ${Math.floor(Math.random() * 40) + 60}%

最终AI置信度: ${aiConfidence}%
推荐结果: ${aiResult === 'approved' ? '通过' : aiResult === 'rejected' ? '拒绝' : '转人工审核'}`,
        created_at: new Date(Date.now() - Math.random() * 30 * 24 * 60 * 60 * 1000).toISOString()
      }
    })
    
    tableData.value = mockData
    total.value = 500
  } catch (error) {
    ElMessage.error('获取数据失败')
  } finally {
    loading.value = false
  }
}

const fetchStats = async () => {
  try {
    // 模拟统计数据
    stats.totalApplications = 1234
    stats.smartApproved = 856
    stats.pending = 278
    stats.rejected = 100
  } catch (error) {
    console.error('获取统计数据失败:', error)
  }
}

const handleStatusChange = () => {
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

const viewDetail = (row: any) => {
  selectedRecord.value = row
  detailDialogVisible.value = true
  // 等待对话框渲染完成后初始化图表
  nextTick(() => {
    initAnalysisCharts()
  })
}

const manualReview = (row: any) => {
  ElMessage.info(`转入人工审批流程: ${row.id}`)
  // 这里可以跳转到人工审批页面
}

// 初始化分析图表
const initAnalysisCharts = () => {
  if (!selectedRecord.value) return
  
  // 根据申请人信息生成个性化数据
  const record = selectedRecord.value
  const seed = parseInt(record.id.replace(/\D/g, '')) || 1000 // 使用ID中的数字作为种子
  
  // 生成基于种子的随机数函数
  const seededRandom = (min: number, max: number, offset: number = 0) => {
    const x = Math.sin((seed + offset) * 12.9898) * 43758.5453
    return Math.floor((x - Math.floor(x)) * (max - min + 1)) + min
  }
  
  // 仪表盘 - 个人信用值（基于申请人信息生成）
  if (creditValueChartRef.value) {
    creditValueChart = echarts.init(creditValueChartRef.value)
    
    // 根据AI审批结果生成个人信用值
    const baseScore = record.ai_result === 'approved' ? 750 : record.ai_result === 'rejected' ? 550 : 650
    const creditValue = Math.max(350, Math.min(850, baseScore + seededRandom(-50, 50, 10)))
    
    const creditOption = {
      tooltip: {
        formatter: '个人信用值: {c}'
      },
      series: [
        {
          name: '个人信用值',
          type: 'gauge',
          min: 350,
          max: 850,
          splitNumber: 5,
          radius: '80%',
          axisLine: {
            lineStyle: {
              width: 15,
              color: [
                [0.3, '#F44336'],
                [0.7, '#FF9800'],
                [1, '#4CAF50']
              ]
            }
          },
          pointer: {
            itemStyle: {
              color: 'auto'
            }
          },
          axisTick: {
            distance: -15,
            length: 8,
            lineStyle: {
              color: '#fff',
              width: 2
            }
          },
          splitLine: {
            distance: -15,
            length: 15,
            lineStyle: {
              color: '#fff',
              width: 4
            }
          },
          axisLabel: {
            color: 'auto',
            distance: 25,
            fontSize: 10
          },
          detail: {
            valueAnimation: true,
            formatter: '{value}',
            color: 'auto',
            fontSize: 20,
            offsetCenter: [0, '70%']
          },
          data: [
            {
              value: creditValue,
              name: '信用值'
            }
          ]
        }
      ]
    }
    creditValueChart.setOption(creditOption)
  }
  
  // 柱状图 - 资产信息（房产、理财产品、社保、公积金）
  if (assetInfoChartRef.value) {
    assetInfoChart = echarts.init(assetInfoChartRef.value)
    
    // 生成资产信息数据
    const assetCategories = ['房产', '理财产品', '社保', '公积金']
    const assetValues = assetCategories.map((_, index) => {
      const baseValue = record.ai_result === 'approved' ? 50000 : record.ai_result === 'rejected' ? 20000 : 35000
      const value = Math.max(0, baseValue + seededRandom(-20000, 30000, index + 50))
      return Math.floor(value)
    })
    
    const assetOption = {
      tooltip: {
        trigger: 'axis',
        axisPointer: {
          type: 'shadow'
        },
        formatter: function(params: any) {
          return `${params[0].name}: ¥${params[0].value.toLocaleString()}`
        }
      },
      grid: {
        left: '3%',
        right: '4%',
        bottom: '3%',
        top: '3%',
        containLabel: true
      },
      xAxis: {
        type: 'category',
        data: assetCategories,
        axisLabel: {
          interval: 0,
          fontSize: 11,
          textStyle: {
            textAlign: 'center'
          }
        },
        axisLine: {
          lineStyle: {
            color: '#E0E0E0'
          }
        }
      },
      yAxis: {
        type: 'value',
        min: 0,
        axisLabel: {
          textStyle: {
            textAlign: 'center'
          },
          formatter: function(value: number) {
            return (value / 10000).toFixed(0) + '万'
          }
        },
        axisLine: {
          show: false
        },
        axisTick: {
          show: false
        },
        splitLine: {
          lineStyle: {
            color: '#EFEFEF'
          }
        }
      },
      series: [
        {
          name: '资产价值',
          type: 'bar',
          data: assetValues,
          itemStyle: {
            color: new echarts.graphic.LinearGradient(0, 0, 0, 1, [
              { offset: 0, color: '#4CAF50' },
              { offset: 1, color: '#81C784' }
            ])
          },
          label: {
            show: true,
            position: 'top',
            textStyle: {
              textAlign: 'center',
              fontSize: 10
            },
            formatter: function(params: any) {
              return (params.value / 10000).toFixed(1) + '万'
            }
          }
        }
      ]
    }
    assetInfoChart.setOption(assetOption)
  }
  
  // 雷达图 - 守约记录（先用后付、免押租车、免押租物、先乘后付、先借后还）
  if (complianceRecordChartRef.value) {
    complianceRecordChart = echarts.init(complianceRecordChartRef.value)
    
    const complianceCategories = ['先用后付', '免押租车', '免押租物', '先乘后付', '先借后还']
    
    // 根据AI审批结果生成守约记录评分
    const baseScore = record.ai_result === 'approved' ? 85 : record.ai_result === 'rejected' ? 60 : 75
    const complianceScores = complianceCategories.map((_, index) => {
      const variation = seededRandom(-15, 15, index + 100)
      return Math.max(40, Math.min(100, baseScore + variation))
    })
    
    const complianceOption = {
      tooltip: {
        trigger: 'item'
      },
      radar: {
        indicator: complianceCategories.map(name => ({
          name: name,
          max: 100,
          min: 0
        })),
        radius: '70%',
        splitNumber: 4,
        axisName: {
          fontSize: 11,
          color: '#333'
        },
        splitLine: {
          lineStyle: {
            color: '#E0E0E0'
          }
        },
        splitArea: {
          show: true,
          areaStyle: {
            color: ['rgba(76, 175, 80, 0.1)', 'rgba(76, 175, 80, 0.05)']
          }
        }
      },
      series: [
        {
          name: '守约记录',
          type: 'radar',
          data: [
            {
              value: complianceScores,
              name: '守约评分',
              itemStyle: {
                color: '#4CAF50'
              },
              areaStyle: {
                color: 'rgba(76, 175, 80, 0.3)'
              },
              lineStyle: {
                color: '#4CAF50',
                width: 2
              }
            }
          ],
          label: {
            show: true,
            fontSize: 10,
            color: '#333',
            formatter: function(params: any) {
              return params.value
            }
          }
        }
      ]
    }
    complianceRecordChart.setOption(complianceOption)
  }
}

// 工具方法
const formatAmount = (amount: number) => {
  return `¥${amount.toLocaleString()}`
}

const formatTime = (time: string) => {
  return dayjs(time).format('YYYY-MM-DD HH:mm:ss')
}

const getConfidenceColor = (confidence: number) => {
  if (confidence >= 80) return '#67c23a'
  if (confidence >= 60) return '#e6a23c'
  return '#f56c6c'
}

const getResultTagType = (result: string) => {
  switch (result) {
    case 'approved': return 'success'
    case 'rejected': return 'danger'
    case 'manual': return 'warning'
    default: return 'info'
  }
}

const getResultText = (result: string) => {
  switch (result) {
    case 'approved': return '智能通过'
    case 'rejected': return '智能拒绝'
    case 'manual': return '转人工审批'
    default: return '未知'
  }
}

onMounted(() => {
  refreshData()
})

// 组件卸载时销毁图表
const destroyCharts = () => {
  assetInfoChart?.dispose()
  creditValueChart?.dispose()
  complianceRecordChart?.dispose()
}
</script>

<style scoped>
.smart-approval-container {
  height: 100vh;
  display: flex;
  flex-direction: column;
  gap: 20px;
  padding: 0 20px;
  overflow: hidden;
}

.header-card {
  flex-shrink: 0;
}

.content-card {
  flex: 1;
  display: flex;
  flex-direction: column;
  overflow: hidden;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
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
  padding: 30px;
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

.stat-number.smart {
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

.confidence-text {
  margin-left: 8px;
  font-size: 12px;
  color: #666;
}

.pagination-wrapper {
  margin-top: 20px;
  display: flex;
  justify-content: center;
  flex-shrink: 0;
}

.detail-content {
  max-height: 600px;
  overflow-y: auto;
  padding-right: 10px;
  /* 详情对话框滚动条样式 */
  scrollbar-width: thin;
  scrollbar-color: #409eff #f1f1f1;
}

.detail-content::-webkit-scrollbar {
  width: 12px;
}

.detail-content::-webkit-scrollbar-track {
  background: #f1f1f1;
  border-radius: 6px;
}

.detail-content::-webkit-scrollbar-thumb {
  background: #409eff;
  border-radius: 6px;
  border: 2px solid #f1f1f1;
}

.detail-content::-webkit-scrollbar-thumb:hover {
  background: #337ecc;
}

.ai-analysis {
  margin-top: 20px;
}

.ai-analysis h4 {
  margin-bottom: 12px;
  color: #333;
}

.analysis-charts {
  margin-top: 16px;
}

.chart-card {
  background: #fff;
  border: 1px solid #ebeef5;
  border-radius: 8px;
  padding: 16px;
  transition: box-shadow 0.3s;
}

.chart-card:hover {
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1);
}

.chart-title {
  font-size: 14px;
  font-weight: 600;
  color: #333;
  margin-bottom: 12px;
  text-align: center;
}

.analysis-chart {
  height: 200px;
  width: 100%;
}

/* Element Plus卡片内容区域样式 */
:deep(.el-card__body) {
  flex: 1;
  display: flex;
  flex-direction: column;
  overflow: hidden;
  padding: 20px;
}

/* 表格容器样式 */
:deep(.el-table) {
  flex: 1;
  overflow: hidden;
}

/* 表格主体滚动区域 */
:deep(.el-table__body-wrapper) {
  max-height: calc(100vh - 400px);
  overflow-y: auto !important;
  overflow-x: auto !important;
  /* 表格滚动条样式 */
  scrollbar-width: thin;
  scrollbar-color: #409eff #f1f1f1;
}

:deep(.el-table__body-wrapper::-webkit-scrollbar) {
  width: 12px;
  height: 12px;
}

:deep(.el-table__body-wrapper::-webkit-scrollbar-track) {
  background: #f1f1f1;
  border-radius: 6px;
}

:deep(.el-table__body-wrapper::-webkit-scrollbar-thumb) {
  background: #409eff;
  border-radius: 6px;
  border: 2px solid #f1f1f1;
}

:deep(.el-table__body-wrapper::-webkit-scrollbar-thumb:hover) {
  background: #337ecc;
}

:deep(.el-table__body-wrapper::-webkit-scrollbar-corner) {
  background: #f1f1f1;
}

/* 表格头部固定 */
:deep(.el-table__header-wrapper) {
  overflow: visible;
}

/* 分页组件位置调整 */
:deep(.el-pagination) {
  margin-top: 20px;
  justify-content: center;
}
</style> 