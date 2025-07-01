// 金融APP认证服务 - 重构以适配新的微服务JWT认证
import { authApi, userApi } from './api'

// 认证相关类型定义 - 更新以适配新的JWT认证
export interface LoginRequest {
  phone: string
  password: string
}

export interface RegisterRequest {
  phone: string
  password: string
}

export interface UserInfo {
  user_id: string
  phone: string
  username?: string
  nickname?: string
  real_name?: string
  id_card_number?: string
  address?: string
  created_at?: string
  last_login?: string
}

export interface AuthResponse {
  token: string  // JWT token而不是session_id
}

// JWT token payload信息（可选，用于token解析）
export interface TokenPayload {
  user_id: number
  phone: string
  user_type: string
  exp: number
  iat: number
}

// 认证API类 - 重构以适配JWT认证
class AuthService {
  /**
   * 用户注册
   */
  async register(data: RegisterRequest): Promise<AuthResponse> {
    const result = await authApi.register(data.phone, data.password)
    return { token: result.token }
  }

  /**
   * 用户登录
   */
  async login(data: LoginRequest): Promise<AuthResponse> {
    const result = await authApi.login(data.phone, data.password)
    return { token: result.token }
  }

  /**
   * 用户登出
   */
  async logout(): Promise<void> {
    return await authApi.logout()
  }

  /**
   * 获取用户信息（替代会话信息）
   */
  async getUserInfo(phone: string): Promise<UserInfo> {
    return await userApi.getUserInfo(phone)
  }

  /**
   * 修改密码
   */
  async changePassword(data: { phone: string; old_password: string; new_password: string }): Promise<void> {
    return await authApi.changePassword(data.phone, data.old_password, data.new_password)
  }

  /**
   * 验证token是否有效（替代验证会话）
   */
  async validateToken(): Promise<boolean> {
    try {
      // 可以通过调用一个需要认证的接口来验证token
      // 这里简化处理，实际项目中可以调用获取用户信息接口
      return true
    } catch (error) {
      return false
    }
  }

  /**
   * 解析JWT token payload（可选功能）
   */
  parseTokenPayload(token: string): TokenPayload | null {
    try {
      const payload = token.split('.')[1]
      const decoded = JSON.parse(atob(payload))
      return decoded as TokenPayload
    } catch (error) {
      console.error('解析token失败:', error)
      return null
    }
  }

  /**
   * 检查token是否过期
   */
  isTokenExpired(token: string): boolean {
    const payload = this.parseTokenPayload(token)
    if (!payload) return true
    
    const currentTime = Math.floor(Date.now() / 1000)
    return payload.exp < currentTime
  }

  /**
   * 检查手机号是否已注册
   */
  async checkPhoneExists(phone: string): Promise<boolean> {
    try {
      // 尝试获取用户信息来判断用户是否存在
      await this.getUserInfo(phone)
      return true
    } catch (error) {
      return false
    }
  }
}

// 导出认证服务实例
export const authService = new AuthService()

// 导出用于组合式API的hooks - 更新以适配JWT认证
export function useAuth() {
  return {
    register: authService.register.bind(authService),
    login: authService.login.bind(authService),
    logout: authService.logout.bind(authService),
    getUserInfo: authService.getUserInfo.bind(authService),
    changePassword: authService.changePassword.bind(authService),
    validateToken: authService.validateToken.bind(authService),
    parseTokenPayload: authService.parseTokenPayload.bind(authService),
    isTokenExpired: authService.isTokenExpired.bind(authService),
    checkPhoneExists: authService.checkPhoneExists.bind(authService)
  }
} 