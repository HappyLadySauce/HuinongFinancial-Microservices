// 金融APP认证服务
import apiClient from './api'

// 认证相关类型定义
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
  username: string
  created_at: string
  last_login?: string
}

export interface AuthResponse {
  session_id: string
}

export interface SessionInfo {
  session_id: string
  user_id: string
  logged_in: boolean
  created_at: string
  last_access: string
  expires_at: string
  ip_address: string
  user_agent: string
  device_id: string
  refresh_count: number
}

// 认证API类
class AuthService {
  /**
   * 用户注册
   */
  async register(data: RegisterRequest): Promise<AuthResponse> {
    return apiClient.post('/register', data)
  }

  /**
   * 用户登录
   */
  async login(data: LoginRequest): Promise<AuthResponse> {
    return apiClient.post('/login', data)
  }

  /**
   * 用户登出
   */
  async logout(): Promise<void> {
    return apiClient.post('/logout')
  }

  /**
   * 获取当前会话信息
   */
  async getSessionInfo(): Promise<SessionInfo> {
    return apiClient.get('/session/info')
  }

  /**
   * 获取用户所有会话
   */
  async getUserSessions(): Promise<SessionInfo[]> {
    const response: any = await apiClient.get('/session/list')
    return response.sessions || []
  }

  /**
   * 删除指定会话
   */
  async deleteSession(sessionId: string): Promise<void> {
    return apiClient.delete(`/session/${sessionId}`)
  }

  /**
   * 删除所有其他会话（保留当前会话）
   */
  async deleteAllOtherSessions(): Promise<void> {
    return apiClient.delete('/session/all')
  }

  /**
   * 刷新当前会话
   */
  async refreshSession(): Promise<void> {
    return apiClient.post('/session/refresh')
  }

  /**
   * 修改密码
   */
  async changePassword(data: { old_password: string; new_password: string }): Promise<void> {
    return apiClient.post('/user/change-password', data)
  }

  /**
   * 验证会话是否有效
   */
  async validateSession(): Promise<boolean> {
    try {
      await this.getSessionInfo()
      return true
    } catch (error) {
      return false
    }
  }

  /**
   * 检查手机号是否已注册
   */
  async checkPhoneExists(phone: string): Promise<boolean> {
    try {
      // 这里可以添加专门的检查接口，或者通过登录接口来判断
      // 暂时返回false，实际项目中需要后端提供相应接口
      return false
    } catch (error) {
      return false
    }
  }
}

// 导出认证服务实例
export const authService = new AuthService()

// 导出用于组合式API的hooks
export function useAuth() {
  return {
    register: authService.register.bind(authService),
    login: authService.login.bind(authService),
    logout: authService.logout.bind(authService),
    getSessionInfo: authService.getSessionInfo.bind(authService),
    getUserSessions: authService.getUserSessions.bind(authService),
    deleteSession: authService.deleteSession.bind(authService),
    deleteAllOtherSessions: authService.deleteAllOtherSessions.bind(authService),
    refreshSession: authService.refreshSession.bind(authService),
    changePassword: authService.changePassword.bind(authService),
    validateSession: authService.validateSession.bind(authService),
    checkPhoneExists: authService.checkPhoneExists.bind(authService)
  }
} 