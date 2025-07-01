import { defineStore } from 'pinia'
import { ref } from 'vue'
import type { AdminUser, LoginResponse } from '@/types'

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

export const useUserStore = defineStore('user', {
    // 状态
    state: () => ({
        token: ref(''),
        userInfo: ref<AdminUser | null>(null),
        isLoggedIn: ref(false),
        loginTime: ref<number>(0)
    }),
    // 操作
    actions: {
        // 设置JWT token
        setToken(token: string) {
            this.token = token
            this.isLoggedIn = true
            this.loginTime = Date.now()
            // 存储到localStorage
            localStorage.setItem('adminToken', token)
            localStorage.setItem('adminLoginTime', this.loginTime.toString())
            
            // 解析JWT获取用户基本信息
            this.parseTokenInfo(token)
        },
        
        // 解析JWT token获取用户信息
        parseTokenInfo(token: string) {
            try {
                const payload = this.parseJWT(token)
                if (payload) {
                    // 创建基本用户信息对象
                    const userInfo: Partial<AdminUser> = {
                        id: payload.user_id,
                        phone: payload.phone,
                        role: payload.role,
                        name: '', // 需要从API获取完整信息
                        nickname: '',
                        age: 0,
                        gender: 0,
                        status: 1,
                        created_at: 0,
                        updated_at: 0
                    }
                    this.userInfo = userInfo as AdminUser
                }
            } catch (error) {
                console.error('解析JWT失败:', error)
            }
        },
        
        // 解析JWT token
        parseJWT(token: string): JWTPayload | null {
            try {
                const parts = token.split('.')
                if (parts.length !== 3) {
                    return null
                }
                
                const payload = parts[1]
                const decoded = atob(payload.replace(/-/g, '+').replace(/_/g, '/'))
                return JSON.parse(decoded) as JWTPayload
            } catch (error) {
                console.error('JWT解析错误:', error)
                return null
            }
        },
        
        // 设置完整用户信息
        setUserInfo(userInfo: AdminUser) {
            this.userInfo = userInfo
            // 存储到localStorage
            localStorage.setItem('adminUserInfo', JSON.stringify(userInfo))
        },
        
        // 登录
        login(loginData: LoginResponse) {
            this.setToken(loginData.token)
        },
        
        // 登出
        logout() {
            this.token = ''
            this.userInfo = null
            this.isLoggedIn = false
            this.loginTime = 0
            // 清除localStorage
            localStorage.removeItem('adminToken')
            localStorage.removeItem('adminUserInfo')
            localStorage.removeItem('adminLoginTime')
        },
        
        // 从localStorage恢复状态
        restoreFromStorage() {
            const token = localStorage.getItem('adminToken')
            const userInfo = localStorage.getItem('adminUserInfo')
            const loginTime = localStorage.getItem('adminLoginTime')
            
            if (token && loginTime) {
                // 检查Token是否过期
                if (this.isTokenValid(token)) {
                    // 未过期，恢复登录状态
                    this.token = token
                    this.isLoggedIn = true
                    this.loginTime = parseInt(loginTime)
                    this.parseTokenInfo(token)
                } else {
                    // 已过期，清除过期数据
                    this.logout()
                    return
                }
            }
            
            if (userInfo && this.isLoggedIn) {
                try {
                    this.userInfo = JSON.parse(userInfo)
                } catch (e) {
                    console.error('解析用户信息失败:', e)
                    localStorage.removeItem('adminUserInfo')
                }
            }
        },
        
        // 检查JWT token是否有效
        isTokenValid(token?: string): boolean {
            const tokenToCheck = token || this.token
            if (!tokenToCheck) return false
            
            try {
                const payload = this.parseJWT(tokenToCheck)
                if (!payload) return false
                
                // 检查过期时间
                const now = Math.floor(Date.now() / 1000)
                return payload.exp > now
            } catch (error) {
                console.error('Token验证失败:', error)
                return false
            }
        },
        
        // 获取Token过期时间
        getTokenExpiry(): number {
            if (!this.token) return 0
            
            try {
                const payload = this.parseJWT(this.token)
                return payload ? payload.exp * 1000 : 0
            } catch (error) {
                return 0
            }
        },
        
        // 检查Token是否即将过期（15分钟内）
        isTokenExpiringSoon(): boolean {
            const expiry = this.getTokenExpiry()
            if (!expiry) return false
            
            const now = Date.now()
            const fifteenMinutes = 15 * 60 * 1000
            return (expiry - now) < fifteenMinutes
        },
        
        // 初始化store - 从localStorage恢复状态
        initialize() {
            this.restoreFromStorage()
        }
    },
    // 计算
    getters: {
        getToken: (state) => state.token,
        getUserInfo: (state) => state.userInfo,
        getIsLoggedIn: (state) => state.isLoggedIn,
        // 获取遮盖后的手机号
        getMaskedPhone: (state) => {
            if (state.userInfo?.phone) {
                const phone = state.userInfo.phone
                return phone.replace(/(\d{3})\d{4}(\d{4})/, '$1****$2')
            }
            return ''
        },
        // 获取用户角色
        getUserRole: (state) => {
            return state.userInfo?.role || ''
        },
        // 获取用户ID
        getUserId: (state) => {
            return state.userInfo?.id || 0
        },
        // 检查是否为管理员
        isAdmin: (state) => {
            return state.userInfo?.role === 'admin'
        },
        // 检查是否为操作员
        isOperator: (state) => {
            return state.userInfo?.role === 'operator'
        },
        // 检查是否为审核员
        isAuditor: (state) => {
            return state.userInfo?.role === 'auditor'
        }
    },
})
