import { useUserStore } from '@/stores/user'
import type { RouteLocationNormalized, NavigationGuardNext } from 'vue-router'
import { userApi } from '@/services/api'

// 需要登录的路由列表
const authRequiredRoutes = [
  '/loan/apply',
  '/loan/application',
  '/loan/my-applications',
  '/finance',
  '/me',
  '/user'
]

// 检查路由是否需要登录
export const isAuthRequired = (path: string): boolean => {
  return authRequiredRoutes.some(route => path.startsWith(route))
}

// 验证会话是否有效
const validateSession = async (): Promise<boolean> => {
  try {
    // 尝试获取当前会话信息来验证会话是否有效
    await userApi.getUserInfo()
    return true
  } catch (error) {
    console.warn('会话验证失败:', error)
    return false
  }
}

// 检查是否需要验证服务器端会话
const shouldValidateSession = (userStore: any): boolean => {
  // 如果没有登录状态，不需要验证
  if (!userStore.isLoggedIn) return false
  
  // 如果登录时间很近（比如5分钟内），可能是刚登录或页面刷新，跳过验证
  const now = Date.now()
  const loginTime = userStore.loginTime
  const fiveMinutes = 5 * 60 * 1000
  
  if (loginTime && (now - loginTime) < fiveMinutes) {
    return false
  }
  
  // 其他情况需要验证
  return true
}

// 路由守卫
export const authGuard = async (
  to: RouteLocationNormalized,
  from: RouteLocationNormalized,
  next: NavigationGuardNext
) => {
  const userStore = useUserStore()
  
  // 检查是否需要登录
  if (isAuthRequired(to.path)) {
    // 首先检查本地登录状态
    if (!userStore.isLoggedIn || !userStore.isTokenValid()) {
      // 未登录，跳转到登录页
      next({
        path: '/login',
        query: { redirect: to.fullPath } // 保存原始路径，登录后可以回到原页面
      })
      return
    }
    
    // 检查是否需要验证服务器端会话
    if (shouldValidateSession(userStore)) {
      try {
        const isValid = await validateSession()
        if (!isValid) {
          // 会话无效，清除本地状态并跳转到登录页
          userStore.logout()
          next({
            path: '/login',
            query: { redirect: to.fullPath }
          })
          return
        }
      } catch (error) {
        // 验证过程中出错，但如果是网络问题等，不要立即清除登录状态
        console.error('会话验证出错:', error)
        // 只有在明确的401错误时才清除状态，网络错误等允许用户继续使用
        if ((error as any)?.message && (error as any).message.includes('401')) {
          userStore.logout()
          next({
            path: '/login',
            query: { redirect: to.fullPath }
          })
          return
        }
        // 网络错误等情况，允许继续访问
      }
    }
  }
  
  // 如果已登录用户访问登录页，重定向到首页或原来要访问的页面
  if (to.path === '/login' && userStore.isLoggedIn && userStore.isTokenValid()) {
    const redirectPath = to.query.redirect as string || '/home'
    next(redirectPath)
    return
  }
  
  next()
}

// 检查用户权限
export const hasPermission = (permission: string): boolean => {
  const userStore = useUserStore()
  
  if (!userStore.isLoggedIn || !userStore.userInfo) {
    return false
  }
  
  // 这里可以根据用户角色或权限进行判断
  // 目前简单返回true，因为普通用户都有基本权限
  return true
}

// 格式化错误信息
export const getErrorMessage = (error: any): string => {
  if (typeof error === 'string') {
    return error
  }
  
  if (error?.message) {
    return error.message
  }
  
  if (error?.response?.data?.message) {
    return error.response.data.message
  }
  
  return '操作失败，请稍后重试'
} 