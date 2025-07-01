import { createRouter, createWebHistory } from 'vue-router'
import { useUserStore } from '@/stores/user'
import AuthService from '@/services/auth'

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: '/login',
      name: 'login',
      component: () => import('@/views/LoginView.vue'),
      meta: { requiresAuth: false, title: '登录' }
    },
    {
      path: '/',
      name: 'layout',
      component: () => import('@/views/LayoutView.vue'),
      meta: { requiresAuth: true },
      redirect: '/dashboard',
      children: [
        // 个人中心
        {
          path: '/dashboard',
          name: 'dashboard',
          component: () => import('@/views/DashboardView.vue'),
          meta: { 
            requiresAuth: true,
            title: '个人中心',
            permission: 'dashboard:view'
          }
        },

        // 核心审批导航
        {
          path: '/approval',
          name: 'approval',
          redirect: '/approval/dashboard',
          meta: { 
            requiresAuth: true,
            title: '核心审批',
            permission: 'approval:view'
          },
          children: [
            {
              path: 'dashboard',
              name: 'approval-dashboard',
              component: () => import('@/views/approval/DashboardView.vue'),
              meta: { 
                requiresAuth: true,
                title: '审批看板',
                permission: 'approval:view'
              }
            },
            {
              path: 'loan',
              name: 'approval-loan',
              component: () => import('@/views/approval/LoanApprovalView.vue'),
              meta: { 
                requiresAuth: true,
                title: '贷款审批',
                permission: 'approval:view'
              }
            },
            {
              path: 'loan/:id',
              name: 'approval-loan-detail',
              component: () => import('@/views/approval/LoanDetailView.vue'),
              meta: { 
                requiresAuth: true,
                title: '贷款审批详情',
                permission: 'approval:view'
              }
            },
            {
              path: 'lease',
              name: 'approval-lease',
              component: () => import('@/views/approval/LeaseApprovalView.vue'),
              meta: { 
                requiresAuth: true,
                title: '租赁审批',
                permission: 'approval:view'
              }
            },
            {
              path: 'lease/:id',
              name: 'approval-lease-detail',
              component: () => import('@/views/approval/LeaseDetailView.vue'),
              meta: { 
                requiresAuth: true,
                title: '租赁审批详情',
                permission: 'approval:view'
              }
            },
            {
              path: 'smart',
              name: 'approval-smart',
              component: () => import('@/views/approval/SmartApprovalView.vue'),
              meta: { 
                requiresAuth: true,
                title: '智能审批',
                permission: 'approval:view'
              }
            }
          ]
        },

        // 运营管理（管理员功能）
        {
          path: '/operation',
          name: 'operation',
          redirect: '/operation/products',
          meta: { 
            requiresAuth: true,
            title: '运营管理',
            permission: 'admin'
          },
          children: [
            {
              path: 'products',
              name: 'operation-products',
              component: () => import('@/views/operation/ProductsView.vue'),
              meta: { 
                requiresAuth: true,
                title: '产品管理',
                permission: 'admin'
              }
            },
            {
              path: 'users',
              name: 'operation-users',
              component: () => import('@/views/operation/UsersView.vue'),
              meta: { 
                requiresAuth: true,
                title: '用户管理',
                permission: 'admin'
              }
            },
            {
              path: 'statistics',
              name: 'operation-statistics',
              component: () => import('@/views/operation/StatisticsView.vue'),
              meta: { 
                requiresAuth: true,
                title: '数据统计',
                permission: 'admin'
              }
            }
          ]
        },

        // 系统管理
        {
          path: '/system',
          name: 'system',
          redirect: '/system/config',
          meta: { 
            requiresAuth: true,
            title: '系统管理',
            permission: 'admin'
          },
          children: [
            {
              path: 'config',
              name: 'system-config',
              component: () => import('@/views/system/ConfigView.vue'),
              meta: { 
                requiresAuth: true,
                title: '系统配置',
                permission: 'admin'
              }
            },
            {
              path: 'logs',
              name: 'system-logs',
              component: () => import('@/views/system/LogsView.vue'),
              meta: { 
                requiresAuth: true,
                title: '操作日志',
                permission: 'admin'
              }
            },
            {
              path: 'security',
              name: 'system-security',
              component: () => import('@/views/system/SecurityView.vue'),
              meta: { 
                requiresAuth: true,
                title: '安全设置',
                permission: 'admin'
              }
            }
          ]
        },

        // 个人设置
        {
          path: '/profile',
          name: 'profile',
          component: () => import('@/views/ProfileView.vue'),
          meta: { 
            requiresAuth: true,
            title: '个人信息',
            permission: 'user'
          }
        },
        {
          path: '/settings',
          name: 'settings',
          component: () => import('@/views/SettingsView.vue'),
          meta: { 
            requiresAuth: true,
            title: '系统设置',
            permission: 'user'
          }
        }
      ]
    },
    {
      path: '/:pathMatch(.*)*',
      name: 'NotFound',
      component: () => import('@/views/NotFoundView.vue'),
      meta: { title: '页面不存在' }
    }
  ],
})

// 路由守卫
router.beforeEach(async (to, from, next) => {
  // 设置页面标题
  if (to.meta.title) {
    document.title = `${to.meta.title} - 惠农金融OA系统`
  }
  
  // 初始化认证服务
  AuthService.init()
  
  // 检查是否需要认证
  if (to.meta.requiresAuth) {
    if (!AuthService.isLoggedIn()) {
      console.log('🔒 需要登录，重定向到登录页')
      next({
        path: '/login',
        query: { redirect: to.fullPath }
      })
      return
    }
    
    // 权限检查
    const requiredPermission = to.meta.permission as string
    if (requiredPermission) {
      if (requiredPermission === 'admin' && !AuthService.isAdmin()) {
        console.log('🚫 权限不足，需要管理员权限')
        next('/dashboard') // 重定向到主页或403页面
        return
      }
      
      if (requiredPermission === 'approval:view' && !AuthService.isAuditor()) {
        console.log('🚫 权限不足，需要审核权限')
        next('/dashboard')
        return
      }
    }
  }
  
  // 如果已登录且访问登录页，重定向到首页
  if (to.name === 'login' && AuthService.isLoggedIn()) {
    console.log('✅ 已登录，重定向到首页')
    next('/dashboard')
    return
  }
  
  next()
})

export default router
