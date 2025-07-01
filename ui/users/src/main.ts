import { createApp } from 'vue'
import { createPinia } from 'pinia'

import App from './App.vue'
import router from './router'
import ElementPlus from 'element-plus'
import 'element-plus/dist/index.css'
import * as ElementPlusIconsVue from '@element-plus/icons-vue'


const app = createApp(App)

app.use(ElementPlus)
const pinia = createPinia()
app.use(pinia)

// 注册 Element Plus 图标组件
for (const [key, component] of Object.entries(ElementPlusIconsVue)) {
  app.component(key, component as any)
}

// 在挂载路由之前初始化用户状态
const { useUserStore } = await import('./stores/user')
const userStore = useUserStore()
// 恢复用户状态，这样页面刷新后不会丢失登录状态
userStore.restoreFromStorage()

app.use(router)

app.mount('#app')
