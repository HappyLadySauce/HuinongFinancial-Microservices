import { createRouter, createWebHashHistory } from 'vue-router'
import { authGuard } from '@/utils/auth'
import IndexPage from '../views/IndexPage.vue'

const router = createRouter({
  history: createWebHashHistory(),
  routes: [
    {
      path: '/',
      name: 'splash',
      component: () => import('../views/SplashScreen.vue'),
    },
    {
      path: '/home',
      name: 'home',
      component: IndexPage
    },
    {
      path: '/login',
      name: 'login',
      component: () => import('../views/login/LoginPage.vue'),
    },
    {
      path: '/register',
      name: 'register',
      component: () => import('../views/login/RegisterPage.vue'),
    },
    {
      path: '/finance',
      name: 'finance',
      component: () => import('../views/FinancePage.vue'),
    },
    {
      path: '/machinery',
      name: 'machinery',
      component: () => import('../views/MachineryPage.vue'),
    },
    {
      path: '/me',
      name: 'me',
      component: () => import('../views/MePage.vue'),
    },
    {
      path: '/news',
      name: 'news',
      component: () => import('../views/NewsPage.vue'),
    },
    {
      path: '/loan/apply/:productId',
      name: 'loanApplication',
      component: () => import('../views/LoanApplicationPage.vue'),
    },
    {
      path: '/loan/apply',
      name: 'loanApplicationNew',
      component: () => import('../views/LoanApplicationPage.vue'),
    },
    {
      path: '/loan/application/:applicationId',
      name: 'loanApplicationDetail',
      component: () => import('../views/LoanApplicationDetailPage.vue'),
    },
    {
      path: '/loan/my-applications',
      name: 'myLoanApplications',
      component: () => import('../views/MyLoanApplicationsPage.vue'),
    },
    {
      path: '/loan/products/:productId',
      name: 'loanProductDetail',
      component: () => import('../views/LoanProductDetailPage.vue'),
    },
    {
      path: '/digital-currency',
      name: 'digitalCurrency',
      component: () => import('../views/DigitalCurrencyPage.vue'),
    }
  ],
})

// 添加路由守卫
router.beforeEach(authGuard)

export default router
