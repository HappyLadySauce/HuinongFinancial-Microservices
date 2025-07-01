// 统一响应格式
export interface ApiResponse<T = any> {
  code?: number
  message?: string
  data?: T
  // JWT认证响应
  token?: string
  // 直接数据响应
  [key: string]: any
}

// 分页响应格式
export interface PaginatedResponse<T> {
  list: T[]
  total: number
  page?: number
  size?: number
}

// B端管理员用户相关类型
export interface AdminUser {
  id: number
  admin_user_id: string
  phone: string
  name: string
  nickname: string
  username: string
  display_name: string
  email: string
  age: number
  gender: number
  role: string
  status: number
  created_at: number
  updated_at: number
}

export interface LoginRequest {
  phone: string
  password: string
  role: string
}

export interface LoginResponse {
  token: string
}

// C端用户相关类型（管理员查看）
export interface AppUser {
  id: number
  phone: string
  name: string
  nickname: string
  age: number
  gender: number
  occupation: string
  address: string
  income: number
  created_at: number
  updated_at: number
}

export interface AppUserInfo {
  user_info: AppUser
}

// 租赁申请相关类型
export interface LeaseApplication {
  id: number
  application_id: string
  user_id: number
  applicant_name: string
  product_id: number
  product_code: string
  name: string
  type: string
  machinery: string
  start_date: string
  end_date: string
  duration: number
  daily_rate: number
  total_amount: number
  deposit: number
  delivery_address: string
  contact_phone: string
  purpose: string
  status: 'pending' | 'approved' | 'rejected' | 'cancelled'
  created_at: number
  updated_at: number
}

export interface LeaseApplicationDetail {
  application_info: LeaseApplication
}

// 租赁审批记录
export interface LeaseApproval {
  id: number
  application_id: number
  auditor_id: number
  auditor_name: string
  action: string
  suggestions: string
  approved_duration: number
  approved_amount: number
  approved_deposit: number
  created_at: number
}

export interface LeaseApprovalRequest {
  action: string
  suggestions: string
  approved_duration: number
  approved_amount: number
  approved_deposit: number
}

// 贷款申请相关类型
export interface LoanApplication {
  id: number
  application_id: string
  user_id: number
  applicant_name: string
  product_id: number
  name: string
  type: string
  amount: number
  duration: number
  purpose: string
  status: 'pending' | 'approved' | 'rejected' | 'cancelled'
  created_at: number
  updated_at: number
}

export interface LoanApplicationDetail {
  application_info: LoanApplication
}

// 贷款审批记录
export interface LoanApproval {
  id: number
  application_id: number
  auditor_id: number
  auditor_name: string
  action: string
  suggestions: string
  approved_amount: number
  approved_duration: number
  interest_rate: number
  created_at: number
}

export interface LoanApprovalRequest {
  action: string
  suggestions: string
  approved_amount: number
  approved_duration: number
  interest_rate: number
}

// 租赁产品相关类型
export interface LeaseProduct {
  id: number
  product_code: string
  name: string
  type: string
  machinery: string
  brand: string
  model: string
  daily_rate: number
  deposit: number
  max_duration: number
  min_duration: number
  description: string
  inventory_count: number
  available_count: number
  status: number
  created_at: number
  updated_at: number
}

export interface LeaseProductRequest {
  product_code: string
  name: string
  type: string
  machinery: string
  brand: string
  model: string
  daily_rate: number
  deposit: number
  max_duration: number
  min_duration: number
  description: string
  inventory_count: number
}

export interface LeaseProductUpdateRequest {
  name: string
  type: string
  machinery: string
  brand: string
  model: string
  daily_rate: number
  deposit: number
  max_duration: number
  min_duration: number
  description: string
  status: number
}

export interface LeaseProductDetail {
  data: LeaseProduct
}

// 贷款产品相关类型
export interface LoanProduct {
  id: number
  product_code: string
  name: string
  type: string
  max_amount: number
  min_amount: number
  max_duration: number
  min_duration: number
  interest_rate: number
  description: string
  status: number
  created_at: number
  updated_at: number
}

export interface LoanProductRequest {
  product_code: string
  name: string
  type: string
  max_amount: number
  min_amount: number
  max_duration: number
  min_duration: number
  interest_rate: number
  description: string
}

export interface LoanProductUpdateRequest {
  name: string
  type: string
  max_amount: number
  min_amount: number
  max_duration: number
  min_duration: number
  interest_rate: number
  description: string
}

export interface LoanProductDetail {
  data: LoanProduct
}

// 库存检查相关
export interface InventoryCheckRequest {
  product_code: string
  quantity: number
  start_date: string
  end_date: string
}

export interface InventoryCheckResponse {
  available: boolean
  available_count: number
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

// 查询参数接口
export interface ApplicationListParams {
  page?: number
  size?: number
  page_size?: number  // 兼容旧版本参数名
  user_id?: number
  product_code?: string
  status?: string
}

export interface ProductListParams {
  page?: number
  size?: number
  type?: string
  brand?: string
  status?: number
  keyword?: string
}

export interface UserListParams {
  page?: number
  size?: number
  role?: string
  status?: number
  keyword?: string
}

// 状态枚举
export enum ApplicationStatus {
  PENDING = 'pending',
  APPROVED = 'approved', 
  REJECTED = 'rejected',
  CANCELLED = 'cancelled'
}

export enum ProductStatus {
  INACTIVE = 0,
  ACTIVE = 1,
  DISABLED = 2
}

export enum UserStatus {
  INACTIVE = 0,
  ACTIVE = 1,
  SUSPENDED = 2
}

export enum UserRole {
  ADMIN = 'admin',
  OPERATOR = 'operator',
  AUDITOR = 'auditor'
}

// 通用时间格式化函数类型
export type DateFormatter = (timestamp: number) => string

// 兼容旧版本的类型别名
export type PaginationResponse<T> = PaginatedResponse<T>
export type ApplicationDetail = LoanApplicationDetail 