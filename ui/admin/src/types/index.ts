// 用户相关类型
export interface AdminUser {
  admin_user_id: string
  username: string
  role: 'ADMIN' | '审批员'
  display_name: string
  email: string
  status: number
  created_at: string
  updated_at: string
}

export interface LoginResponse {
  admin_user_id: string
  username: string
  role: string
  token: string
}

// 工作台数据类型
export interface DashboardData {
  stats: {
    total_applications: number
    pending_review: number
    approved_today: number
    ai_efficiency: number
  }
  pending_tasks: Array<{
    task_id: string
    task_type: string
    title: string
    priority: 'high' | 'medium' | 'low'
    created_at: string
    application_id?: string
  }>
  recent_activities: Array<{
    activity_id: string
    activity_type: string
    description: string
    timestamp: string
    operator: string
  }>
}

// 贷款申请相关类型
export interface LoanApplication {
  application_id: string
  applicant_name: string
  amount: number
  status: string
  submission_time: string
  ai_risk_score?: number
  ai_suggestion?: string
}

export interface ApplicationDetail {
  application_id: string
  amount: number
  term_months: number
  purpose: string
  status: string
  submitted_at: string
  updated_at: string
  approved_amount?: number
  applicant_details: {
    user_id: string
    real_name: string
    id_card_number: string
    phone: string
    address: string
    bank_account: string
  }
  uploaded_documents_details: Array<{
    file_id: string
    doc_type: string
    file_url: string
    ocr_result?: string
  }>
  ai_analysis_report?: {
    overall_risk_score: number
    risk_factors: string[]
    data_verification_results: Array<{
      item: string
      result: string
    }>
    suggestion: string
  }
  history: Array<{
    status: string
    timestamp: string
    operator: string
  }>
}

// 操作日志类型
export interface OperationLog {
  id: string
  operator_id: string
  operator_name: string
  action: string
  target: string
  result: string
  ip_address: string
  user_agent: string
  occurred_at: string
}

// 系统配置类型
export interface SystemConfig {
  config_key: string
  config_value: string
  description: string
  updated_at: string
}

// 分页相关类型
export interface PaginationResponse<T> {
  data: T[]
  total: number
  page: number
  limit: number
} 