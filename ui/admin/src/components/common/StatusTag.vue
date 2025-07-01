<template>
  <el-tag 
    :type="tagType" 
    :size="size"
    :effect="effect"
    :class="tagClass"
  >
    <el-icon v-if="showIcon" class="tag-icon">
      <component :is="statusIcon" />
    </el-icon>
    {{ displayText }}
  </el-tag>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import {
  Clock,
  Check,
  Close,
  Warning,
  Document,
  Loading,
  CircleCheck,
  CircleClose,
  QuestionFilled
} from '@element-plus/icons-vue'

interface Props {
  status: string
  statusMap?: Record<string, { type: string; text?: string; icon?: any }>
  showIcon?: boolean
  size?: 'large' | 'default' | 'small'
  effect?: 'dark' | 'light' | 'plain'
  customText?: string
}

const props = withDefaults(defineProps<Props>(), {
  showIcon: false,
  size: 'default',
  effect: 'light'
})

// 默认状态映射
const defaultStatusMap = {
  // 通用状态
  'active': { type: 'success', text: '活跃', icon: Check },
  'inactive': { type: 'info', text: '非活跃', icon: Close },
  'disabled': { type: 'danger', text: '禁用', icon: Close },
  'enabled': { type: 'success', text: '启用', icon: Check },
  
  // 审批状态
  'pending': { type: 'warning', text: '待审批', icon: Clock },
  'approved': { type: 'success', text: '已通过', icon: Check },
  'rejected': { type: 'danger', text: '已拒绝', icon: Close },
  'processing': { type: 'primary', text: '处理中', icon: Loading },
  'submitted': { type: 'info', text: '已提交', icon: Document },
  
  // 贷款/租赁状态
  'draft': { type: 'info', text: '草稿', icon: Document },
  'reviewing': { type: 'warning', text: '审核中', icon: Clock },
  'risk_assessment': { type: 'warning', text: '风险评估', icon: Warning },
  'manual_review': { type: 'warning', text: '人工审核', icon: Clock },
  'approved_pending': { type: 'primary', text: '批准待放款', icon: Clock },
  'disbursed': { type: 'success', text: '已放款', icon: CircleCheck },
  'repaying': { type: 'primary', text: '还款中', icon: Loading },
  'completed': { type: 'success', text: '已完成', icon: CircleCheck },
  'overdue': { type: 'danger', text: '逾期', icon: Warning },
  'defaulted': { type: 'danger', text: '违约', icon: CircleClose },
  
  // AI状态
  'ai_processing': { type: 'primary', text: 'AI处理中', icon: Loading },
  'ai_approved': { type: 'success', text: 'AI通过', icon: Check },
  'ai_rejected': { type: 'danger', text: 'AI拒绝', icon: Close },
  'manual_required': { type: 'warning', text: '需人工审核', icon: Clock },
  
  // 智能审批状态
  'smart_approved': { type: 'success', text: '智能通过', icon: Check },
  'smart_rejected': { type: 'danger', text: '智能拒绝', icon: Close },
  'manual': { type: 'warning', text: '转人工审批', icon: Clock },
  
  // 风险等级
  'low_risk': { type: 'success', text: '低风险', icon: Check },
  'medium_risk': { type: 'warning', text: '中风险', icon: Warning },
  'high_risk': { type: 'danger', text: '高风险', icon: Warning },
  
  // 产品状态
  'online': { type: 'success', text: '上架', icon: Check },
  'offline': { type: 'info', text: '下架', icon: Close },
  'maintaining': { type: 'warning', text: '维护中', icon: Warning },
  
  // 用户状态
  'verified': { type: 'success', text: '已认证', icon: Check },
  'unverified': { type: 'warning', text: '未认证', icon: Warning },
  'locked': { type: 'danger', text: '已锁定', icon: Close },
  
  // 系统状态
  'normal': { type: 'success', text: '正常', icon: Check },
  'warning': { type: 'warning', text: '警告', icon: Warning },
  'error': { type: 'danger', text: '错误', icon: Close },
  'maintenance': { type: 'info', text: '维护', icon: Warning }
}

const statusConfig = computed(() => {
  const statusMap = props.statusMap || defaultStatusMap
  return statusMap[props.status] || { 
    type: 'info', 
    text: props.status, 
    icon: QuestionFilled 
  }
})

const tagType = computed(() => statusConfig.value.type as any)

const displayText = computed(() => {
  return props.customText || statusConfig.value.text || props.status
})

const statusIcon = computed(() => statusConfig.value.icon || QuestionFilled)

const tagClass = computed(() => [
  'status-tag',
  `status-${props.status}`,
  { 'has-icon': props.showIcon }
])
</script>

<style scoped>
.status-tag {
  display: inline-flex;
  align-items: center;
  gap: 4px;
}

.tag-icon {
  font-size: 12px;
}

.has-icon {
  padding-left: 6px;
  padding-right: 8px;
}
</style> 