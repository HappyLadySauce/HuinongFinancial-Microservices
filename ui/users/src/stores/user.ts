import { defineStore } from 'pinia'
import { ref } from 'vue'

// 用户信息接口
export interface UserInfo {
  user_id: string
  phone: string
  nickname: string
  avatar_url: string
  real_name: string
  id_card_number: string
  address: string
}

// 登录响应接口
export interface LoginResponse {
  session_id: string
}

export const useUserStore = defineStore('user', {
    // 状态
    state: () => ({
        sessionId: ref(''),
        userInfo: ref<UserInfo | null>(null),
        isLoggedIn: ref(false),
        loginTime: ref<number>(0)
    }),
    // 操作
    actions: {
        // 设置session
        setSession(sessionId: string) {
            this.sessionId = sessionId
            this.isLoggedIn = true
            this.loginTime = Date.now()
            // 存储到localStorage
            localStorage.setItem('sessionId', sessionId)
            localStorage.setItem('loginTime', this.loginTime.toString())
        },
        
        // 设置用户信息
        setUserInfo(userInfo: UserInfo) {
            this.userInfo = userInfo
            // 存储到localStorage
            localStorage.setItem('userInfo', JSON.stringify(userInfo))
        },
        
        // 登录
        login(loginData: LoginResponse) {
            this.setSession(loginData.session_id)
        },
        
        // 登出
        logout() {
            this.sessionId = ''
            this.userInfo = null
            this.isLoggedIn = false
            this.loginTime = 0
            // 清除localStorage
            localStorage.removeItem('sessionId')
            localStorage.removeItem('userInfo')
            localStorage.removeItem('loginTime')
        },
        
        // 从localStorage恢复状态
        restoreFromStorage() {
            const sessionId = localStorage.getItem('sessionId')
            const userInfo = localStorage.getItem('userInfo')
            const loginTime = localStorage.getItem('loginTime')
            
            if (sessionId && loginTime) {
                // 检查登录时间是否过期（比如7天）
                const now = Date.now()
                const loginTimestamp = parseInt(loginTime)
                const sevenDays = 7 * 24 * 60 * 60 * 1000
                
                if (now - loginTimestamp < sevenDays) {
                    // 未过期，恢复登录状态
                    this.sessionId = sessionId
                    this.isLoggedIn = true
                    this.loginTime = loginTimestamp
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
                    localStorage.removeItem('userInfo')
                }
            }
        },
        
        // 检查session是否可能有效
        isTokenValid(): boolean {
            // 对于Cookie认证，主要依赖服务器验证
            // 这里只做基本的状态检查
            return this.isLoggedIn && this.sessionId !== ''
        },
        
        // 初始化store - 从localStorage恢复状态
        initialize() {
            // 只恢复基本的用户信息，不自动设置登录状态
            // 登录状态的验证将由路由守卫处理
            const userInfo = localStorage.getItem('userInfo')
            const sessionId = localStorage.getItem('sessionId')
            const loginTime = localStorage.getItem('loginTime')
            
            if (userInfo) {
                try {
                    this.userInfo = JSON.parse(userInfo)
                } catch (e) {
                    console.error('解析用户信息失败:', e)
                    localStorage.removeItem('userInfo')
                }
            }
            
            // 只有当所有必要信息都存在时才恢复登录状态
            if (sessionId && loginTime && userInfo) {
                this.sessionId = sessionId
                this.loginTime = parseInt(loginTime)
                this.isLoggedIn = true
            }
        }
    },
    // 计算
    getters: {
        getSessionId: (state) => state.sessionId,
        getUserInfo: (state) => state.userInfo,
        getIsLoggedIn: (state) => state.isLoggedIn,
        // 获取遮盖后的手机号
        getMaskedPhone: (state) => {
            if (state.userInfo?.phone) {
                const phone = state.userInfo.phone
                return phone.replace(/(\d{3})\d{4}(\d{4})/, '$1****$2')
            }
            return ''
        }
    },
})
