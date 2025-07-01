import { useUserStore } from '@/stores/user'
import type { UserInfo, LoginResponse } from '@/stores/user'
import router from '@/router'

// API基础配置
const API_BASE_URL = '/api/v1'

// 统一响应格式
interface ApiResponse<T = any> {
  code: number
  message: string
  data: T
}

// 分页响应格式
export interface PaginatedResponse<T> {
  list: T[]
  total: number
  page: number
  page_size: number
}

// 通用请求函数
const apiRequest = async (url: string, options: RequestInit = {}) => {
  const defaultOptions: RequestInit = {
    credentials: 'include',
    headers: {
      'Content-Type': 'application/json',
      ...options.headers
    }
  }

  const finalOptions = { ...defaultOptions, ...options }
  
  try {
    const response = await fetch(`${API_BASE_URL}${url}`, finalOptions)
    const result = await response.json()
    
    if (result.code === 200) {
      return result.data
    } else {
      // 根据错误码进行不同处理
      switch (result.code) {
        case 401:
          // 未登录或会话过期，清除本地状态并跳转到登录页
          const userStore = useUserStore()
          userStore.logout()
          // 使用路由跳转而不是直接操作location
          if (router.currentRoute.value.path !== '/login') {
            router.push({
              path: '/login',
              query: { redirect: router.currentRoute.value.fullPath }
            })
          }
          throw new Error(result.message || '登录已过期，请重新登录')
        case 403:
          throw new Error('权限不足，无法访问该资源')
        case 404:
          throw new Error('请求的资源不存在')
        case 500:
          throw new Error('服务器内部错误，请稍后重试')
        default:
          throw new Error(result.message || '请求失败')
      }
    }
  } catch (error: any) {
    if (error.message) {
      throw error
    } else {
      throw new Error('网络请求失败，请检查网络连接')
    }
  }
}

// 贷款申请相关接口
export interface LoanApproval {
  id: number
  user_id: number
  name: string
  type: string
  amount: number
  duration: number
  description: string
  status: 'pending' | 'approved' | 'rejected'
  suggestions: string
  auditor: string
  created_at: string
  updated_at: string
}

export interface LoanApprovalRequest {
  name: string
  type: string
  amount: number
  duration: number
  description: string
}

// 租赁申请相关接口
export interface LeaseApproval {
  id: number
  user_id: number
  name: string
  type: string
  start_at: string
  end_at: string
  description: string
  status: 'pending' | 'approved' | 'rejected'
  suggestions: string
  auditor: string
  created_at: string
  updated_at: string
}

export interface LeaseApprovalRequest {
  name: string
  type: string
  start_at: string
  end_at: string
  description: string
}

// 产品类型接口
export interface ProductTypes {
  types: string[]
}

// 审批历史接口
export interface ReviewRequest {
  status: 'approved' | 'rejected'
  suggestions: string
  auditor: string
}

// 文件上传结果接口
export interface FileUploadResult {
  file_id: string
  file_url: string
  file_name: string
  file_size: number
}

// 用户认证API
export const authApi = {
  login: async (phone: string, password: string) => {
    return await apiRequest('/login', {
      method: 'POST',
      body: JSON.stringify({ phone, password })
    })
  },
  
  logout: async () => {
    return await apiRequest('/logout', {
      method: 'POST'
    })
  }
}

// 贷款申请API
export const loanApprovalApi = {
  // 创建贷款申请
  create: async (data: LoanApprovalRequest): Promise<LoanApproval> => {
    return await apiRequest('/loan/approvals', {
      method: 'POST',
      body: JSON.stringify(data)
    })
  },

  // 获取用户贷款申请列表
  getMyApprovals: async (params: { page?: number; page_size?: number } = {}): Promise<PaginatedResponse<LoanApproval>> => {
    const queryParams = new URLSearchParams()
    if (params.page) queryParams.append('page', params.page.toString())
    if (params.page_size) queryParams.append('page_size', params.page_size.toString())
    
    const url = `/loan/approvals${queryParams.toString() ? '?' + queryParams.toString() : ''}`
    return await apiRequest(url)
  },

  // 获取单个贷款申请详情
  getDetail: async (id: number): Promise<LoanApproval> => {
    return await apiRequest(`/loan/approvals/${id}`)
  },

  // 更新贷款申请
  update: async (id: number, data: Partial<LoanApprovalRequest>): Promise<LoanApproval> => {
    return await apiRequest(`/loan/approvals/${id}`, {
      method: 'PUT',
      body: JSON.stringify(data)
    })
  },

  // 删除贷款申请
  delete: async (id: number) => {
    return await apiRequest(`/loan/approvals/${id}`, {
      method: 'DELETE'
    })
  },

  // 获取贷款产品类型
  getTypes: async (): Promise<ProductTypes> => {
    return await apiRequest('/loan/types')
  }
}

// 租赁申请API
export const leaseApprovalApi = {
  // 创建租赁申请
  create: async (data: LeaseApprovalRequest): Promise<LeaseApproval> => {
    return await apiRequest('/lease/approvals', {
      method: 'POST',
      body: JSON.stringify(data)
    })
  },

  // 获取用户租赁申请列表
  getMyApprovals: async (params: { page?: number; page_size?: number } = {}): Promise<PaginatedResponse<LeaseApproval>> => {
    const queryParams = new URLSearchParams()
    if (params.page) queryParams.append('page', params.page.toString())
    if (params.page_size) queryParams.append('page_size', params.page_size.toString())
    
    const url = `/lease/approvals${queryParams.toString() ? '?' + queryParams.toString() : ''}`
    return await apiRequest(url)
  },

  // 获取单个租赁申请详情
  getDetail: async (id: number): Promise<LeaseApproval> => {
    return await apiRequest(`/lease/approvals/${id}`)
  },

  // 更新租赁申请
  update: async (id: number, data: Partial<LeaseApprovalRequest>): Promise<LeaseApproval> => {
    return await apiRequest(`/lease/approvals/${id}`, {
      method: 'PUT',
      body: JSON.stringify(data)
    })
  },

  // 删除租赁申请
  delete: async (id: number) => {
    return await apiRequest(`/lease/approvals/${id}`, {
      method: 'DELETE'
    })
  },

  // 获取租赁产品类型
  getTypes: async (): Promise<ProductTypes> => {
    return await apiRequest('/lease/types')
  }
}

// 管理员审批API（需要管理员权限）
export const adminApi = {
  // 获取所有贷款申请
  getAllLoanApprovals: async (params: { 
    page?: number; 
    page_size?: number; 
    status?: string 
  } = {}): Promise<PaginatedResponse<LoanApproval>> => {
    const queryParams = new URLSearchParams()
    if (params.page) queryParams.append('page', params.page.toString())
    if (params.page_size) queryParams.append('page_size', params.page_size.toString())
    if (params.status) queryParams.append('status', params.status)
    
    const url = `/admin/loan/approvals${queryParams.toString() ? '?' + queryParams.toString() : ''}`
    return await apiRequest(url)
  },

  // 审批贷款申请
  reviewLoanApproval: async (id: number, data: ReviewRequest) => {
    return await apiRequest(`/admin/loan/approvals/${id}/review`, {
      method: 'POST',
      body: JSON.stringify(data)
    })
  },

  // 获取所有租赁申请
  getAllLeaseApprovals: async (params: { 
    page?: number; 
    page_size?: number; 
    status?: string 
  } = {}): Promise<PaginatedResponse<LeaseApproval>> => {
    const queryParams = new URLSearchParams()
    if (params.page) queryParams.append('page', params.page.toString())
    if (params.page_size) queryParams.append('page_size', params.page_size.toString())
    if (params.status) queryParams.append('status', params.status)
    
    const url = `/admin/lease/approvals${queryParams.toString() ? '?' + queryParams.toString() : ''}`
    return await apiRequest(url)
  },

  // 审批租赁申请
  reviewLeaseApproval: async (id: number, data: ReviewRequest) => {
    return await apiRequest(`/admin/lease/approvals/${id}/review`, {
      method: 'POST',
      body: JSON.stringify(data)
    })
  }
}

// 创建API客户端实例
const apiClient = {
  // 通用请求方法
  async request<T>(
    endpoint: string,
    options: RequestInit = {}
  ): Promise<T> {
    const userStore = useUserStore()
    const url = `${API_BASE_URL}${endpoint}`
    
    // 设置默认headers
    const headers: Record<string, string> = {
      'Content-Type': 'application/json',
      ...((options.headers as Record<string, string>) || {})
    }

    // 如果有session且需要认证，添加Authorization header
    if (userStore.sessionId && userStore.isTokenValid()) {
      headers.Authorization = `Bearer ${userStore.sessionId}`
    }

    const config: RequestInit = {
      ...options,
      headers
    }

    try {
      const response = await fetch(url, config)
      
      if (!response.ok) {
        throw new Error(`HTTP error! status: ${response.status}`)
      }

      const result = await response.json()
      
      // 检查API响应码 (后端成功状态码是200)
      if (result.code !== 200) {
        throw new Error(result.message || '请求失败')
      }

      return result.data || result
    } catch (error) {
      console.error('API请求失败:', error)
      throw error
    }
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
  },

  // 文件上传
  async uploadFile(file: File, purpose?: string): Promise<ApiResponse<FileUploadResult>> {
    const userStore = useUserStore()
    const url = `${API_BASE_URL}/files/upload`
    
    const formData = new FormData()
    formData.append('file', file)
    if (purpose) {
      formData.append('purpose', purpose)
    }

    const headers: Record<string, string> = {}
    if (userStore.sessionId && userStore.isTokenValid()) {
      headers.Authorization = `Bearer ${userStore.sessionId}`
    }

    const response = await fetch(url, {
      method: 'POST',
      headers,
      body: formData
    })

    if (!response.ok) {
      throw new Error(`上传失败: ${response.status}`)
    }

    return response.json()
  }
}

// 兼容性接口定义（用于与旧代码兼容）
export interface RegisterRequest {
  phone: string
  password: string
}

export interface LoginRequest {
  phone: string
  password: string
}

export interface UpdateUserRequest {
  nickname?: string
  avatar_url?: string
  real_name?: string
  address?: string
}

// 贷款产品接口（用于产品展示页面）
export interface LoanProduct {
  product_id: string
  name: string
  description: string
  category: string
  min_amount: number
  max_amount: number
  min_term_months: number
  max_term_months: number
  interest_rate_yearly: string
  status: number
  repayment_methods?: string[]
  application_conditions?: string
  required_documents?: Array<{
    type: string
    desc: string
  }>
}

// 兼容性类型别名
export type LoanApplication = LoanApproval
export type LeaseApplication = LeaseApproval
export type LoanApplicationRequest = LoanApprovalRequest
export type LeaseApplicationRequest = LeaseApprovalRequest

// 模拟用户数据存储
const mockUsers = [
  {
    user_id: '1',
    phone: '13800138000',
    password: '123456',
    nickname: '张三',
    avatar_url: 'https://cube.elemecdn.com/3/7c/3ea6beec64369c2642b92c6726f1epng.png',
    real_name: '张三',
    id_card_number: '110101199001011234',
    address: '北京市朝阳区'
  },
  {
    user_id: '2', 
    phone: '13800138001',
    password: '123456',
    nickname: '李四',
    avatar_url: 'https://cube.elemecdn.com/0/88/03b0d39583f48206768a7534e55bcpng.png',
    real_name: '李四',
    id_card_number: '110101199001011235',
    address: '上海市浦东新区'
  }
]

// 用户相关API
export const userApi = {
  // 注册 - 模拟实现
  async register(data: RegisterRequest): Promise<ApiResponse<{ user_id: string }>> {
    return new Promise((resolve, reject) => {
      setTimeout(() => {
        // 检查手机号是否已存在
        const existingUser = mockUsers.find(user => user.phone === data.phone)
        if (existingUser) {
          reject(new Error('手机号已注册'))
          return
        }
        
        // 模拟创建新用户
        const newUser = {
          user_id: String(mockUsers.length + 1),
          phone: data.phone,
          password: data.password,
          nickname: '新用户',
          avatar_url: 'https://cube.elemecdn.com/3/7c/3ea6beec64369c2642b92c6726f1epng.png',
          real_name: '',
          id_card_number: '',
          address: ''
        }
        
        mockUsers.push(newUser)
        
        resolve({
          code: 0,
          message: '注册成功',
          data: { user_id: newUser.user_id }
        })
      }, 1000) // 模拟网络延迟
    })
  },

  // 登录 - 调用真实后端API
  async login(data: LoginRequest): Promise<ApiResponse<LoginResponse>> {
    try {
      const response = await apiRequest('/login', {
        method: 'POST',
        body: JSON.stringify(data)
      })
      
      return {
        code: 200,
          message: '登录成功',
        data: response
      }
    } catch (error: any) {
      throw new Error(error.message || '登录失败')
    }
  },

  // 获取用户信息 - 调用真实后端API
  async getUserInfo(): Promise<ApiResponse<UserInfo>> {
    try {
      const response = await apiRequest('/session/info')
      
      // 从session信息中构造用户信息
      const mockUserInfo = {
        user_id: response.session.user_id,
        phone: '13800138000', // 临时使用测试手机号
        nickname: '测试用户',
        avatar_url: 'https://cube.elemecdn.com/3/7c/3ea6beec64369c2642b92c6726f1epng.png',
        real_name: '张三',
        id_card_number: '110101199001011234',
        address: '北京市朝阳区'
      }
      
      return {
        code: 200,
          message: '获取用户信息成功',
        data: mockUserInfo
      }
    } catch (error: any) {
      throw new Error(error.message || '获取用户信息失败')
    }
  },

  // 更新用户信息
  updateUserInfo(data: UpdateUserRequest): Promise<ApiResponse> {
    return apiClient.put('/users/me', data)
  }
}

// 兼容性贷款API（保持旧接口可用）
export const loanApi = {
  // 获取贷款产品列表 - 兼容旧接口
  getProducts(category?: string): Promise<ApiResponse<LoanProduct[]>> {
    const params = category ? { category } : undefined
    return apiClient.get('/loans/products', params)
  },

  // 获取贷款产品详情 - 兼容旧接口
  getProductDetail(productId: string): Promise<ApiResponse<LoanProduct>> {
    return apiClient.get(`/loans/products/${productId}`)
  },

  // 提交贷款申请 - 兼容旧接口
  submitApplication(data: LoanApplicationRequest): Promise<ApiResponse<{ application_id: string }>> {
    return apiClient.post('/loans/applications', data)
  },

  // 获取贷款申请详情 - 兼容旧接口
  getApplicationDetail(applicationId: string): Promise<ApiResponse<LoanApplication>> {
    return apiClient.get(`/loans/applications/${applicationId}`)
  },

  // 获取我的贷款申请列表 - 兼容旧接口
  getMyApplications(params?: {
    status?: string
    page?: number
    limit?: number
  }): Promise<PaginatedResponse<LoanApplication>> {
    return apiClient.get('/loans/applications/my', params)
  },

  // 新接口的别名，保持向后兼容
  ...loanApprovalApi
}

// 租赁API别名
export const leaseApi = leaseApprovalApi

// 文件服务API
export const fileApi = {
  // 文件上传
  upload(file: File, purpose?: string): Promise<ApiResponse<FileUploadResult>> {
    return apiClient.uploadFile(file, purpose)
  }
}

// 健康检查API
export const healthApi = {
  check(): Promise<ApiResponse<{ status: string; service: string; version: string }>> {
    return apiClient.get('/health')
  }
}

export default apiClient 