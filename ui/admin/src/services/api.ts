import { useUserStore } from '@/stores/user'
import router from '@/router'
import type {
  ApiResponse,
  PaginatedResponse,
  AdminUser,
  LoginRequest,
  LoginResponse,
  AppUser,
  AppUserInfo,
  LeaseApplication,
  LeaseApplicationDetail,
  LeaseApproval,
  LeaseApprovalRequest,
  LoanApplication,
  LoanApplicationDetail,
  LoanApproval,
  LoanApprovalRequest,
  LeaseProduct,
  LeaseProductRequest,
  LeaseProductUpdateRequest,
  LeaseProductDetail,
  LoanProduct,
  LoanProductRequest,
  LoanProductUpdateRequest,
  LoanProductDetail,
  InventoryCheckRequest,
  InventoryCheckResponse,
  ApplicationListParams,
  ProductListParams,
  UserListParams,
  OperationLog,
  SystemConfig
} from '@/types'

// API基础配置
const API_BASE_URL = 'http://127.0.0.1:8080/api/v1'

// 统一响应格式处理
interface RequestOptions extends RequestInit {
  skipAuth?: boolean
}

// 通用请求函数
const apiRequest = async <T = any>(url: string, options: RequestOptions = {}): Promise<T> => {
  const userStore = useUserStore()
  
  const defaultHeaders: Record<string, string> = {
    'Content-Type': 'application/json'
  }

  // 添加JWT认证头（除非明确跳过）
  if (!options.skipAuth && userStore.token && userStore.isTokenValid()) {
    defaultHeaders.Authorization = `Bearer ${userStore.token}`
  }

  const finalOptions: RequestInit = {
    ...options,
    headers: {
      ...defaultHeaders,
      ...options.headers
    }
  }
  
  try {
    const response = await fetch(`${API_BASE_URL}${url}`, finalOptions)
    
    if (!response.ok) {
      throw new Error(`HTTP error! status: ${response.status}`)
    }
    
    const result = await response.json()
    
    // 处理JWT登录响应 - 直接返回token
    if (result.token) {
      return result
    }
    
    // 处理标准响应格式
    if (result.code !== undefined) {
      if (result.code === 200) {
        return result.data || result
      } else {
        throw new Error(result.message || '请求失败')
      }
    }
    
    // 直接返回数据（用于新的API格式）
    return result
    
  } catch (error: any) {
    // 处理401未授权错误
    if (error.message.includes('401') || error.message.includes('unauthorized')) {
      const userStore = useUserStore()
      userStore.logout()
      if (router.currentRoute.value.path !== '/login') {
        router.push({
          path: '/login',
          query: { redirect: router.currentRoute.value.fullPath }
        })
      }
      throw new Error('登录已过期，请重新登录')
    }
    
    throw error.message ? error : new Error('网络请求失败，请检查网络连接')
  }
}

// ================================
// B端用户认证API (oauser服务)
// ================================
export const adminAuthApi = {
  // 管理员注册
  register: async (data: LoginRequest): Promise<LoginResponse> => {
    return await apiRequest('/auth/register', {
      method: 'POST',
      body: JSON.stringify(data),
      skipAuth: true
    })
  },

  // 管理员登录
  login: async (phone: string, password: string): Promise<LoginResponse> => {
    return await apiRequest('/auth/login', {
      method: 'POST',
      body: JSON.stringify({ phone, password }),
      skipAuth: true
    })
  },
  
  // 管理员登出
  logout: async (): Promise<void> => {
    return await apiRequest('/auth/logout', {
      method: 'POST'
    })
  },

  // 修改密码
  changePassword: async (phone: string, old_password: string, new_password: string): Promise<void> => {
    return await apiRequest('/auth/password', {
      method: 'POST',
      body: JSON.stringify({ phone, old_password, new_password })
    })
  }
}

// ================================
// B端用户信息API (oauser服务)
// ================================
export const adminUserApi = {
  // 获取用户列表
  getUsers: async (params: UserListParams = {}): Promise<PaginatedResponse<AdminUser>> => {
    const queryParams = new URLSearchParams()
    if (params.page) queryParams.append('page', params.page.toString())
    if (params.size) queryParams.append('size', params.size.toString())
    if (params.role) queryParams.append('role', params.role)
    if (params.status !== undefined) queryParams.append('status', params.status.toString())
    if (params.keyword) queryParams.append('keyword', params.keyword)
    
    const url = `/admin/users${queryParams.toString() ? '?' + queryParams.toString() : ''}`
    return await apiRequest(url)
  },

  // 创建用户
  createUser: async (data: {
    username: string;
    password: string;
    display_name: string;
    role: string;
    email: string;
  }): Promise<{ user_info: AdminUser }> => {
    return await apiRequest('/admin/users', {
      method: 'POST',
      body: JSON.stringify(data)
    })
  },

  // 获取管理员信息
  getUserInfo: async (phone: string): Promise<{ user_info: AdminUser }> => {
    return await apiRequest('/user/info', {
      method: 'GET',
      body: JSON.stringify({ phone })
    })
  },

  // 更新管理员信息
  updateUserInfo: async (user_info: Partial<AdminUser>): Promise<{ user_info: AdminUser }> => {
    return await apiRequest('/user/info', {
      method: 'PUT',
      body: JSON.stringify({ user_info })
    })
  },

  // 更新用户状态
  updateUserStatus: async (user_id: string, status: number): Promise<{ status: number }> => {
    return await apiRequest('/admin/users/status', {
      method: 'PUT',
      body: JSON.stringify({ user_id, status })
    })
  },

  // 删除用户
  deleteUser: async (phone: string): Promise<void> => {
    return await apiRequest('/user/delete', {
      method: 'POST',
      body: JSON.stringify({ phone })
    })
  }
}

// ================================
// C端用户管理API (appuser服务)
// ================================
export const appUserApi = {
  // 获取C端用户信息
  getUserInfo: async (phone: string): Promise<AppUserInfo> => {
    return await apiRequest('/user/info', {
      method: 'GET',
      body: JSON.stringify({ phone })
    })
  },

  // 更新C端用户信息
  updateUserInfo: async (user_info: Partial<AppUser>): Promise<AppUserInfo> => {
    return await apiRequest('/user/info', {
      method: 'PUT',
      body: JSON.stringify({ user_info })
    })
  },

  // 删除C端用户
  deleteUser: async (phone: string): Promise<void> => {
    return await apiRequest('/user/delete', {
      method: 'POST',
      body: JSON.stringify({ phone })
    })
  }
}

// ================================
// 租赁申请管理API (lease服务)
// ================================
export const adminLeaseApi = {
  // 获取所有租赁申请列表
  getAllApplications: async (params: ApplicationListParams = {}): Promise<PaginatedResponse<LeaseApplication>> => {
    const queryParams = new URLSearchParams()
    if (params.page) queryParams.append('page', params.page.toString())
    if (params.size) queryParams.append('size', params.size.toString())
    if (params.user_id) queryParams.append('user_id', params.user_id.toString())
    if (params.product_code) queryParams.append('product_code', params.product_code)
    if (params.status) queryParams.append('status', params.status)
    
    const url = `/admin/lease/applications${queryParams.toString() ? '?' + queryParams.toString() : ''}`
    return await apiRequest(url)
  },

  // 获取租赁申请详情
  getApplicationDetail: async (id: string): Promise<LeaseApplicationDetail> => {
    return await apiRequest(`/admin/lease/applications/${id}`)
  },

  // 获取租赁审批记录
  getApprovals: async (application_id: string): Promise<{ list: LeaseApproval[] }> => {
    return await apiRequest(`/admin/lease/applications/${application_id}/approvals?application_id=${application_id}`)
  },

  // 审批租赁申请
  approveApplication: async (id: string, data: LeaseApprovalRequest): Promise<void> => {
    return await apiRequest(`/admin/lease/applications/${id}/approve`, {
      method: 'POST',
      body: JSON.stringify(data)
    })
  },

  // 审核申请 (别名方法，兼容性)
  reviewApproval: async (application_id: string, data: any): Promise<void> => {
    return await apiRequest(`/admin/lease/applications/${application_id}/review`, {
      method: 'POST',
      body: JSON.stringify(data)
    })
  },

  // 获取详情 (别名方法，兼容性)
  getDetail: async (id: string): Promise<{ application_info: LeaseApplication }> => {
    const result = await apiRequest(`/admin/lease/applications/${id}`)
    return { application_info: result }
  }
}

// ================================
// 贷款申请管理API (loan服务)
// ================================
export const adminLoanApi = {
  // 获取所有贷款申请列表
  getAllApplications: async (params: ApplicationListParams = {}): Promise<PaginatedResponse<LoanApplication>> => {
    const queryParams = new URLSearchParams()
    if (params.page) queryParams.append('page', params.page.toString())
    if (params.size) queryParams.append('size', params.size.toString())
    if (params.user_id) queryParams.append('user_id', params.user_id.toString())
    if (params.status) queryParams.append('status', params.status)
    
    const url = `/admin/loan/applications${queryParams.toString() ? '?' + queryParams.toString() : ''}`
    return await apiRequest(url)
  },

  // 获取贷款申请详情
  getApplicationDetail: async (id: string): Promise<LoanApplicationDetail> => {
    return await apiRequest(`/admin/loan/applications/${id}`)
  },

  // 获取贷款审批记录
  getApprovals: async (params: ApplicationListParams = {}): Promise<PaginatedResponse<LoanApplication>> => {
    const queryParams = new URLSearchParams()
    if (params.page) queryParams.append('page', params.page.toString())
    if (params.size) queryParams.append('size', params.size.toString())
    if (params.status) queryParams.append('status', params.status)
    
    const url = `/admin/loan/applications${queryParams.toString() ? '?' + queryParams.toString() : ''}`
    return await apiRequest(url)
  },

  // 审批贷款申请
  approveApplication: async (id: string, data: LoanApprovalRequest): Promise<void> => {
    return await apiRequest(`/admin/loan/applications/${id}/approve`, {
      method: 'POST',
      body: JSON.stringify(data)
    })
  },

  // 审核申请 (别名方法，兼容性)
  reviewApproval: async (application_id: string, data: any): Promise<void> => {
    return await apiRequest(`/admin/loan/applications/${application_id}/review`, {
      method: 'POST',
      body: JSON.stringify(data)
    })
  },

  // 获取详情 (别名方法，兼容性)
  getDetail: async (id: string): Promise<{ application_info: LoanApplication }> => {
    const result = await apiRequest(`/admin/loan/applications/${id}`)
    return { application_info: result }
  }
}

// ================================
// 租赁产品管理API (leaseproduct服务)
// ================================
export const adminLeaseProductApi = {
  // 获取所有租赁产品列表
  getAllProducts: async (params: ProductListParams = {}): Promise<PaginatedResponse<LeaseProduct>> => {
    const queryParams = new URLSearchParams()
    queryParams.append('page', (params.page || 1).toString())
    queryParams.append('size', (params.size || 10).toString())
    if (params.type) queryParams.append('type', params.type)
    if (params.brand) queryParams.append('brand', params.brand)
    if (params.status !== undefined) queryParams.append('status', params.status.toString())
    if (params.keyword) queryParams.append('keyword', params.keyword)
    
    return await apiRequest(`/admin/leaseproduct/products?${queryParams.toString()}`)
  },

  // 创建租赁产品
  createProduct: async (data: LeaseProductRequest): Promise<LeaseProductDetail> => {
    return await apiRequest('/admin/leaseproduct/products', {
      method: 'POST',
      body: JSON.stringify(data)
    })
  },

  // 获取租赁产品详情
  getProductDetail: async (productCode: string): Promise<LeaseProductDetail> => {
    return await apiRequest(`/admin/leaseproduct/products/${productCode}`)
  },

  // 更新租赁产品
  updateProduct: async (productCode: string, data: LeaseProductUpdateRequest): Promise<LeaseProductDetail> => {
    return await apiRequest(`/admin/leaseproduct/products/${productCode}`, {
      method: 'PUT',
      body: JSON.stringify(data)
    })
  },

  // 删除租赁产品
  deleteProduct: async (productCode: string): Promise<void> => {
    return await apiRequest(`/admin/leaseproduct/products/${productCode}`, {
      method: 'DELETE'
    })
  }
}

// C端租赁产品查看API
export const leaseProductApi = {
  // 获取租赁产品列表（C端视图）
  getProducts: async (params: ProductListParams = {}): Promise<PaginatedResponse<LeaseProduct>> => {
    const queryParams = new URLSearchParams()
    queryParams.append('page', (params.page || 1).toString())
    queryParams.append('size', (params.size || 10).toString())
    if (params.type) queryParams.append('type', params.type)
    if (params.brand) queryParams.append('brand', params.brand)
    if (params.status !== undefined) queryParams.append('status', params.status.toString())
    if (params.keyword) queryParams.append('keyword', params.keyword)
    
    return await apiRequest(`/leaseproduct/products?${queryParams.toString()}`)
  },

  // 获取租赁产品详情（C端视图）
  getProductDetail: async (productCode: string): Promise<LeaseProductDetail> => {
    return await apiRequest(`/leaseproduct/products/${productCode}`)
  },

  // 检查库存可用性
  checkInventoryAvailability: async (data: InventoryCheckRequest): Promise<InventoryCheckResponse> => {
    return await apiRequest('/leaseproduct/products/check-inventory', {
      method: 'POST',
      body: JSON.stringify(data)
    })
  }
}

// ================================
// 贷款产品管理API (loanproduct服务)
// ================================
export const adminLoanProductApi = {
  // 获取所有贷款产品列表
  getAllProducts: async (params: ProductListParams = {}): Promise<PaginatedResponse<LoanProduct>> => {
    const queryParams = new URLSearchParams()
    queryParams.append('page', (params.page || 1).toString())
    queryParams.append('size', (params.size || 10).toString())
    if (params.type) queryParams.append('type', params.type)
    if (params.status !== undefined) queryParams.append('status', params.status.toString())
    if (params.keyword) queryParams.append('keyword', params.keyword)
    
    return await apiRequest(`/admin/loanproduct/products?${queryParams.toString()}`)
  },

  // 创建贷款产品
  createProduct: async (data: LoanProductRequest): Promise<LoanProductDetail> => {
    return await apiRequest('/admin/loanproduct/products', {
      method: 'POST',
      body: JSON.stringify(data)
    })
  },

  // 获取贷款产品详情
  getProductDetail: async (id: number): Promise<LoanProductDetail> => {
    return await apiRequest(`/admin/loanproduct/products/${id}`)
  },

  // 更新贷款产品
  updateProduct: async (id: number, data: LoanProductUpdateRequest): Promise<LoanProductDetail> => {
    return await apiRequest(`/admin/loanproduct/products/${id}`, {
      method: 'PUT',
      body: JSON.stringify(data)
    })
  },

  // 更新产品状态
  updateProductStatus: async (id: number, status: number): Promise<void> => {
    return await apiRequest(`/admin/loanproduct/products/${id}/status`, {
      method: 'PUT',
      body: JSON.stringify({ status })
    })
  },

  // 删除贷款产品
  deleteProduct: async (id: number): Promise<void> => {
    return await apiRequest(`/admin/loanproduct/products/${id}`, {
      method: 'DELETE'
    })
  }
}

// C端贷款产品查看API
export const loanProductApi = {
  // 获取贷款产品列表（C端视图）
  getProducts: async (params: ProductListParams = {}): Promise<PaginatedResponse<LoanProduct>> => {
    const queryParams = new URLSearchParams()
    queryParams.append('page', (params.page || 1).toString())
    queryParams.append('size', (params.size || 10).toString())
    if (params.type) queryParams.append('type', params.type)
    if (params.status !== undefined) queryParams.append('status', params.status.toString())
    if (params.keyword) queryParams.append('keyword', params.keyword)
    
    return await apiRequest(`/loanproduct/products?${queryParams.toString()}`)
  },

  // 获取贷款产品详情（C端视图）
  getProductDetail: async (id: number): Promise<LoanProductDetail> => {
    return await apiRequest(`/loanproduct/products/${id}`)
  }
}

// ================================
// 创建API客户端实例
// ================================
const apiClient = {
  // 通用请求方法
  async request<T>(endpoint: string, options: RequestOptions = {}): Promise<T> {
    return await apiRequest(endpoint, options)
  },

  // GET请求
  get<T>(endpoint: string, params?: Record<string, any>): Promise<T> {
    const url = params ? `${endpoint}?${new URLSearchParams(params)}` : endpoint
    return this.request<T>(url, { method: 'GET' })
  },

  // POST请求
  post<T>(endpoint: string, data?: any): Promise<T> {
    return this.request<T>(endpoint, {
      method: 'POST',
      body: data ? JSON.stringify(data) : undefined
    })
  },

  // PUT请求
  put<T>(endpoint: string, data?: any): Promise<T> {
    return this.request<T>(endpoint, {
      method: 'PUT',
      body: data ? JSON.stringify(data) : undefined
    })
  },

  // DELETE请求
  delete<T>(endpoint: string): Promise<T> {
    return this.request<T>(endpoint, { method: 'DELETE' })
  }
}

// ================================
// 工具函数
// ================================

// 格式化Unix时间戳
export const formatTimestamp = (timestamp: number): string => {
  return new Date(timestamp * 1000).toLocaleString('zh-CN')
}

// 格式化金额
export const formatCurrency = (amount: number): string => {
  return new Intl.NumberFormat('zh-CN', {
    style: 'currency',
    currency: 'CNY'
  }).format(amount)
}

// 状态标签映射
export const statusLabels = {
  pending: '待审核',
  approved: '已通过',
  rejected: '已拒绝',
  cancelled: '已取消'
}

// 产品状态标签
export const productStatusLabels = {
  0: '停用',
  1: '启用',
  2: '禁用'
}

// 用户状态标签
export const userStatusLabels = {
  0: '停用',
  1: '启用',
  2: '暂停'
}

// 角色标签
export const roleLabels = {
  admin: '管理员',
  operator: '操作员',
  auditor: '审核员'
}

// ================================
// 别名导出 - 兼容旧版本命名
// ================================
export const adminLoanApprovalApi = adminLoanApi
export const adminLeaseApprovalApi = adminLeaseApi
// ================================
// 操作日志API
// ================================
export const adminLogApi = {
  // 获取操作日志列表
  getOperationLogs: async (params: {
    operator_id?: string;
    action?: string;
    start_date?: string;
    end_date?: string;
    page?: number;
    limit?: number;
  } = {}): Promise<PaginatedResponse<OperationLog>> => {
    const queryParams = new URLSearchParams()
    if (params.operator_id) queryParams.append('operator_id', params.operator_id)
    if (params.action) queryParams.append('action', params.action)
    if (params.start_date) queryParams.append('start_date', params.start_date)
    if (params.end_date) queryParams.append('end_date', params.end_date)
    if (params.page) queryParams.append('page', params.page.toString())
    if (params.limit) queryParams.append('size', params.limit.toString())
    
    const url = `/admin/logs${queryParams.toString() ? '?' + queryParams.toString() : ''}`
    return await apiRequest(url)
  }
}

// ================================
// 系统配置API
// ================================
export const adminSystemApi = {
  // 获取系统配置列表
  getSystemConfigurations: async (): Promise<SystemConfig[]> => {
    return await apiRequest('/admin/system/configurations')
  },

  // 更新系统配置
  updateSystemConfiguration: async (key: string, value: string): Promise<void> => {
    return await apiRequest('/admin/system/configurations', {
      method: 'PUT',
      body: JSON.stringify({ key, value })
    })
  },

  // 切换AI审批状态
  toggleAIApproval: async (enabled: boolean): Promise<void> => {
    return await apiRequest('/admin/system/ai-approval', {
      method: 'PUT',
      body: JSON.stringify({ enabled })
    })
  }
}

export const adminApi = {
  loan: adminLoanApi,
  lease: adminLeaseApi,
  user: adminUserApi,
  auth: adminAuthApi,
  leaseProduct: adminLeaseProductApi,
  loanProduct: adminLoanProductApi,
  log: adminLogApi,
  system: adminSystemApi,
  
  // 别名方法 - 兼容旧版本调用
  getAllLoanApprovals: adminLoanApi.getAllApplications,
  getAllLeaseApprovals: adminLeaseApi.getAllApplications,
  reviewLoanApproval: adminLoanApi.reviewApproval,
  reviewLeaseApproval: adminLeaseApi.approveApplication
}

// 类型导出
export type { 
  LeaseApproval,
  LoanApproval,
  ApplicationDetail
} from '@/types'

export default apiClient 