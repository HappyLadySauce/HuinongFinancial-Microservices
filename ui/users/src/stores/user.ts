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

// 登录响应接口 - 更新为JWT token
export interface LoginResponse {
  token: string  // JWT token替代session_id
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
        // 设置token - 更新为JWT token处理
        setSession(token: string) {
            this.sessionId = token  // sessionId现在存储JWT token
            this.isLoggedIn = true
            this.loginTime = Date.now()
            // 存储到localStorage
            localStorage.setItem('sessionId', token)  // 保持key名称不变以兼容
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
            this.setSession(loginData.token)
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
        
        // 检查JWT token是否可能有效
        isTokenValid(): boolean {
            if (!this.isLoggedIn || !this.sessionId) {
                return false
            }
            
            // 简单检查token格式（JWT应该有3个部分用.分割）
            const tokenParts = this.sessionId.split('.')
            if (tokenParts.length !== 3) {
                return false
            }
            
            try {
                // 检查token是否过期
                const payload = JSON.parse(atob(tokenParts[1]))
                const currentTime = Math.floor(Date.now() / 1000)
                return payload.exp > currentTime
            } catch (error) {
                console.error('Token验证失败:', error)
                return false
            }
        },
        
        // 初始化store - 从localStorage恢复JWT token状态
        initialize() {
            // 恢复用户信息
            const userInfo = localStorage.getItem('userInfo')
            const token = localStorage.getItem('sessionId')  // sessionId现在存储JWT token
            const loginTime = localStorage.getItem('loginTime')
            
            if (userInfo) {
                try {
                    this.userInfo = JSON.parse(userInfo)
                } catch (e) {
                    console.error('解析用户信息失败:', e)
                    localStorage.removeItem('userInfo')
                }
            }
            
            // 恢复JWT token状态
            if (token && loginTime) {
                this.sessionId = token
                this.loginTime = parseInt(loginTime)
                this.isLoggedIn = true
                
                // 验证token有效性
                if (!this.isTokenValid()) {
                    console.log('JWT token已过期，清除登录状态')
                    this.logout()
                }
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
