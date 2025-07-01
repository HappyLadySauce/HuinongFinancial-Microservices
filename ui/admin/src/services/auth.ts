import { useUserStore } from '@/stores/user'
import { adminAuthApi, adminUserApi } from './api'
import router from '@/router'
import type { LoginRequest, AdminUser } from '@/types'

// JWT Token解码接口
interface JWTPayload {
  user_id: number
  phone: string
  user_type: string
  role: string
  iss: string
  exp: number
  nbf: number
  iat: number
}

// JWT Token管理类
class TokenManager {
  private refreshTimer: number | null = null
  
  // 检查Token是否即将过期，如果是则提醒用户
  checkTokenExpiry() {
    const userStore = useUserStore()
    
    if (!userStore.token || !userStore.isTokenValid()) {
      this.logout()
      return
    }
    
    if (userStore.isTokenExpiringSoon()) {
      console.warn('Token即将过期，请及时重新登录')
      // 可以在这里添加用户提醒逻辑
      // ElMessage.warning('登录即将过期，请及时重新登录')
    }
  }
  
  // 启动Token检查定时器
  startTokenCheck() {
    // 每5分钟检查一次
    this.refreshTimer = setInterval(() => {
      this.checkTokenExpiry()
    }, 5 * 60 * 1000)
  }
  
  // 停止Token检查定时器
  stopTokenCheck() {
    if (this.refreshTimer) {
      clearInterval(this.refreshTimer)
      this.refreshTimer = null
    }
  }
  
  // 登出
  logout() {
    this.stopTokenCheck()
    const userStore = useUserStore()
    userStore.logout()
    
    // 跳转到登录页
    if (router.currentRoute.value.path !== '/login') {
      router.push('/login')
    }
  }
}

// 创建全局Token管理器实例
const tokenManager = new TokenManager()

// 认证服务类
export class AuthService {
  // 登录
  static async login(phone: string, password: string): Promise<void> {
    try {
      const response = await adminAuthApi.login(phone, password)
      const userStore = useUserStore()
      
      // 设置token
      userStore.login(response)
      
      // 启动token检查
      tokenManager.startTokenCheck()
      
      // 获取完整的用户信息
      await this.fetchUserInfo(phone)
      
    } catch (error: any) {
      console.error('登录失败:', error)
      throw new Error(error.message || '登录失败')
    }
  }
  
  // 注册（管理员注册）
  static async register(data: LoginRequest): Promise<void> {
    try {
      const response = await adminAuthApi.register(data)
      const userStore = useUserStore()
      
      // 设置token
      userStore.login(response)
      
      // 启动token检查
      tokenManager.startTokenCheck()
      
      // 获取完整的用户信息
      await this.fetchUserInfo(data.phone)
      
    } catch (error: any) {
      console.error('注册失败:', error)
      throw new Error(error.message || '注册失败')
    }
  }
  
  // 登出
  static async logout(): Promise<void> {
    try {
      // 调用后端登出接口
      await adminAuthApi.logout()
    } catch (error) {
      console.error('后端登出失败:', error)
      // 即使后端登出失败，也要清除本地状态
    } finally {
      // 清除本地状态
      tokenManager.logout()
    }
  }
  
  // 获取用户信息
  static async fetchUserInfo(phone: string): Promise<AdminUser | null> {
    try {
      const response = await adminUserApi.getUserInfo(phone)
      const userStore = useUserStore()
      
      userStore.setUserInfo(response.user_info)
      return response.user_info
    } catch (error: any) {
      console.error('获取用户信息失败:', error)
      return null
    }
  }
  
  // 修改密码
  static async changePassword(phone: string, oldPassword: string, newPassword: string): Promise<void> {
    try {
      await adminAuthApi.changePassword(phone, oldPassword, newPassword)
    } catch (error: any) {
      console.error('修改密码失败:', error)
      throw new Error(error.message || '修改密码失败')
    }
  }
  
  // 检查用户是否已登录
  static isLoggedIn(): boolean {
    const userStore = useUserStore()
    return userStore.isLoggedIn && userStore.isTokenValid()
  }
  
  // 检查用户权限
  static hasPermission(requiredRole: string | string[]): boolean {
    const userStore = useUserStore()
    
    if (!this.isLoggedIn() || !userStore.userInfo) {
      return false
    }
    
    const userRole = userStore.userInfo.role
    
    if (Array.isArray(requiredRole)) {
      return requiredRole.includes(userRole)
    }
    
    return userRole === requiredRole
  }
  
  // 检查是否为管理员
  static isAdmin(): boolean {
    return this.hasPermission('admin')
  }
  
  // 检查是否为操作员
  static isOperator(): boolean {
    return this.hasPermission(['admin', 'operator'])
  }
  
  // 检查是否为审核员
  static isAuditor(): boolean {
    return this.hasPermission(['admin', 'auditor'])
  }
  
  // 权限守卫 - 用于路由守卫
  static checkRoutePermission(to: any): boolean {
    // 如果没有登录，需要登录
    if (!this.isLoggedIn()) {
      return false
    }
    
    // 检查路由所需权限
    const requiredRoles = to.meta?.roles as string[] | undefined
    
    if (!requiredRoles || requiredRoles.length === 0) {
      // 没有权限要求，允许访问
      return true
    }
    
    return this.hasPermission(requiredRoles)
  }
  
  // 初始化认证服务
  static init(): void {
    const userStore = useUserStore()
    
    // 从localStorage恢复状态
    userStore.initialize()
    
    // 如果有有效token，启动检查
    if (userStore.isLoggedIn && userStore.isTokenValid()) {
      tokenManager.startTokenCheck()
    }
  }
  
  // 清理资源
  static cleanup(): void {
    tokenManager.stopTokenCheck()
  }
}

// 路由守卫函数
export const authGuard = (to: any, from: any, next: any) => {
  // 登录页面直接放行
  if (to.path === '/login') {
    next()
    return
  }
  
  // 检查登录状态
  if (!AuthService.isLoggedIn()) {
    next({
      path: '/login',
      query: { redirect: to.fullPath }
    })
    return
  }
  
  // 检查路由权限
  if (!AuthService.checkRoutePermission(to)) {
    console.warn('权限不足，无法访问该页面')
    next('/403') // 跳转到403页面或首页
    return
  }
  
  next()
}

// 权限指令 - 用于v-permission
export const permissionDirective = {
  mounted(el: HTMLElement, binding: any) {
    const { value } = binding
    
    if (value && !AuthService.hasPermission(value)) {
      // 隐藏或禁用元素
      el.style.display = 'none'
      // 或者移除元素
      // el.parentNode?.removeChild(el)
    }
  },
  
  updated(el: HTMLElement, binding: any) {
    const { value } = binding
    
    if (value && !AuthService.hasPermission(value)) {
      el.style.display = 'none'
    } else {
      el.style.display = ''
    }
  }
}

// 用户信息工具函数
export const userUtils = {
  // 获取当前用户信息
  getCurrentUser(): AdminUser | null {
    const userStore = useUserStore()
    return userStore.userInfo
  },
  
  // 获取用户角色
  getUserRole(): string {
    const userStore = useUserStore()
    return userStore.getUserRole
  },
  
  // 获取用户ID
  getUserId(): number {
    const userStore = useUserStore()
    return userStore.getUserId
  },
  
  // 获取遮盖后的手机号
  getMaskedPhone(): string {
    const userStore = useUserStore()
    return userStore.getMaskedPhone
  }
}

// 导出默认认证服务
export default AuthService 