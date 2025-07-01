import { useUserStore } from '@/stores/user'
import type { UserInfo, LoginResponse } from '@/stores/user'
import router from '@/router'

// API基础配置 - 更新为新的微服务地址
const API_BASE_URL = 'http://127.0.0.1:8080/api/v1'

// 统一响应格式 - 根据新后端调整
interface ApiResponse<T = any> {
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

// 通用请求函数 - 重构以支持JWT认证
const apiRequest = async (url: string, options: RequestInit = {}): Promise<any> => {
  const userStore = useUserStore()
  
  const defaultHeaders: Record<string, string> = {
    'Content-Type': 'application/json'
  }

  // 添加JWT认证头
  if (userStore.sessionId && userStore.isTokenValid()) {
    defaultHeaders.Authorization = `Bearer ${userStore.sessionId}`
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

// 租赁申请相关接口 - 根据新API重构
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

export interface LeaseApplicationRequest {
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
}

export interface LeaseApplicationUpdateRequest {
  purpose: string
  delivery_address: string
  contact_phone: string
}

// 贷款申请相关接口 - 根据新API重构
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

export interface LoanApplicationRequest {
  product_id: number
  name: string
  type: string
  amount: number
  duration: number
  purpose: string
}

export interface LoanApplicationUpdateRequest {
  amount: number
  duration: number
  purpose: string
}

// 租赁产品接口
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

// 贷款产品接口
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

// 用户认证API - 重构以适配新的JWT认证
export const authApi = {
  // 用户注册
  register: async (phone: string, password: string): Promise<{ token: string }> => {
    return await apiRequest('/auth/register', {
      method: 'POST',
      body: JSON.stringify({ phone, password })
    })
  },

  // 用户登录 - 返回JWT token
  login: async (phone: string, password: string): Promise<{ token: string }> => {
    return await apiRequest('/auth/login', {
      method: 'POST',
      body: JSON.stringify({ phone, password })
    })
  },
  
  // 用户登出
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

// 用户信息API
export const userApi = {
  // 获取用户信息
  getUserInfo: async (phone: string): Promise<UserInfo> => {
    return await apiRequest('/user/info', {
      method: 'GET',
      body: JSON.stringify({ phone })
    })
  },

  // 更新用户信息
  updateUserInfo: async (user_info: Partial<UserInfo>): Promise<UserInfo> => {
    return await apiRequest('/user/info', {
      method: 'PUT',
      body: JSON.stringify({ user_info })
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

// 租赁申请API - 重构为新的微服务接口
export const leaseApprovalApi = {
  // 创建租赁申请
  create: async (data: LeaseApplicationRequest): Promise<{ application_id: string }> => {
    return await apiRequest('/lease/applications', {
      method: 'POST',
      body: JSON.stringify(data)
    })
  },

  // 获取我的租赁申请列表
  getMyApprovals: async (params: { page?: number; size?: number; status?: string } = {}): Promise<PaginatedResponse<LeaseApplication>> => {
    const queryParams = new URLSearchParams()
    if (params.page) queryParams.append('page', params.page.toString())
    if (params.size) queryParams.append('size', params.size.toString())
    if (params.status) queryParams.append('status', params.status)
    
    const url = `/lease/applications${queryParams.toString() ? '?' + queryParams.toString() : ''}`
    return await apiRequest(url)
  },

  // 获取单个租赁申请详情
  getDetail: async (id: string): Promise<{ application_info: LeaseApplication }> => {
    return await apiRequest(`/lease/applications/${id}`)
  },

  // 更新租赁申请
  update: async (id: string, data: LeaseApplicationUpdateRequest): Promise<{ application_info: LeaseApplication }> => {
    return await apiRequest(`/lease/applications/${id}`, {
      method: 'PUT',
      body: JSON.stringify(data)
    })
  },

  // 取消租赁申请
  cancel: async (id: string, reason: string): Promise<void> => {
    return await apiRequest(`/lease/applications/${id}/cancel`, {
      method: 'POST',
      body: JSON.stringify({ reason })
    })
  }
}

// 贷款申请API - 重构为新的微服务接口
export const loanApprovalApi = {
  // 创建贷款申请
  create: async (data: LoanApplicationRequest): Promise<{ application_id: string }> => {
    return await apiRequest('/loan/applications', {
      method: 'POST',
      body: JSON.stringify(data)
    })
  },

  // 获取我的贷款申请列表
  getMyApprovals: async (params: { page?: number; size?: number; status?: string } = {}): Promise<PaginatedResponse<LoanApplication>> => {
    const queryParams = new URLSearchParams()
    if (params.page) queryParams.append('page', params.page.toString())
    if (params.size) queryParams.append('size', params.size.toString())
    if (params.status) queryParams.append('status', params.status)
    
    const url = `/loan/applications${queryParams.toString() ? '?' + queryParams.toString() : ''}`
    return await apiRequest(url)
  },

  // 获取单个贷款申请详情
  getDetail: async (id: string): Promise<{ application_info: LoanApplication }> => {
    return await apiRequest(`/loan/applications/${id}`)
  },

  // 更新贷款申请
  update: async (id: string, data: LoanApplicationUpdateRequest): Promise<{ application_info: LoanApplication }> => {
    return await apiRequest(`/loan/applications/${id}`, {
      method: 'PUT',
      body: JSON.stringify(data)
    })
  },

  // 取消贷款申请
  cancel: async (id: string, reason: string): Promise<void> => {
    return await apiRequest(`/loan/applications/${id}/cancel`, {
      method: 'POST',
      body: JSON.stringify({ reason })
    })
  }
}

// 租赁产品API
export const leaseProductApi = {
  // 获取租赁产品列表
  getProducts: async (params: { 
    page?: number; 
    size?: number; 
    type?: string; 
    brand?: string; 
    status?: number;
    keyword?: string;
  } = {}): Promise<PaginatedResponse<LeaseProduct>> => {
    const queryParams = new URLSearchParams()
    queryParams.append('page', (params.page || 1).toString())
    queryParams.append('size', (params.size || 10).toString())
    if (params.type) queryParams.append('type', params.type)
    if (params.brand) queryParams.append('brand', params.brand)
    if (params.status !== undefined) queryParams.append('status', params.status.toString())
    if (params.keyword) queryParams.append('keyword', params.keyword)
    
    return await apiRequest(`/leaseproduct/products?${queryParams.toString()}`)
  },

  // 获取租赁产品详情
  getProductDetail: async (productCode: string): Promise<{ data: LeaseProduct }> => {
    return await apiRequest(`/leaseproduct/products/${productCode}`)
  },

  // 检查库存可用性
  checkInventoryAvailability: async (data: {
    product_code: string;
    quantity: number;
    start_date: string;
    end_date: string;
  }): Promise<{ available: boolean; available_count: number }> => {
    return await apiRequest('/leaseproduct/products/check-inventory', {
      method: 'POST',
      body: JSON.stringify(data)
    })
  }
}

// 贷款产品API
export const loanProductApi = {
  // 获取贷款产品列表
  getProducts: async (params: { 
    page?: number; 
    size?: number; 
    type?: string; 
    status?: number;
    keyword?: string;
  } = {}): Promise<PaginatedResponse<LoanProduct>> => {
    const queryParams = new URLSearchParams()
    queryParams.append('page', (params.page || 1).toString())
    queryParams.append('size', (params.size || 10).toString())
    if (params.type) queryParams.append('type', params.type)
    if (params.status !== undefined) queryParams.append('status', params.status.toString())
    if (params.keyword) queryParams.append('keyword', params.keyword)
    
    return await apiRequest(`/loanproduct/products?${queryParams.toString()}`)
  },

  // 获取贷款产品详情
  getProductDetail: async (id: number): Promise<{ data: LoanProduct }> => {
    return await apiRequest(`/loanproduct/products/${id}`)
  }
}

// 兼容性类型别名
export type LoanApproval = LoanApplication
export type LeaseApproval = LeaseApplication  
export type LoanApprovalRequest = LoanApplicationRequest
export type LeaseApprovalRequest = LeaseApplicationRequest

// 兼容性API别名
export const loanApi = loanApprovalApi
export const leaseApi = leaseApprovalApi

// 创建API客户端实例 - 简化版本
const apiClient = {
  // 通用请求方法
  async request<T>(endpoint: string, options: RequestInit = {}): Promise<T> {
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

export default apiClient