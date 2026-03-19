import { createRouter, createWebHistory } from 'vue-router'
import { useAuthStore } from '../stores/auth'

const routes = [
  {
    path: '/admin/login',
    name: 'Login',
    component: () => import('../views/Login.vue'),
    meta: { guest: true },
  },
  {
    path: '/admin/setup',
    name: 'Setup',
    component: () => import('../views/Setup.vue'),
    meta: { guest: true },
  },
  {
    path: '/admin',
    component: () => import('../layouts/DefaultLayout.vue'),
    meta: { requiresAuth: true },
    children: [
      { path: '', redirect: '/admin/dashboard' },
      { path: 'dashboard', name: 'Dashboard', component: () => import('../views/Dashboard.vue') },
      { path: 'services', name: 'Services', component: () => import('../views/Services.vue') },
      { path: 'endpoints', name: 'Endpoints', component: () => import('../views/Endpoints.vue') },
      { path: 'policies', name: 'Policies', component: () => import('../views/Policies.vue') },
    ],
  },
  {
    path: '/',
    redirect: '/admin',
  },
]

const router = createRouter({
  history: createWebHistory(),
  routes,
})

router.beforeEach(async (to, _from, next) => {
  const auth = useAuthStore()

  // 每次会话首次导航时检查是否需要初始化
  if (!auth.setupChecked) {
    await auth.checkSetup()
  }

  // 需要初始化 → 强制跳转 Setup（除非已经在 Setup 页）
  if (auth.needSetup) {
    if (to.name === 'Setup') {
      next()
    } else {
      next('/admin/setup')
    }
    return
  }

  // 已初始化后，不允许再访问 Setup
  if (to.name === 'Setup') {
    next('/admin/login')
    return
  }

  if (to.meta.requiresAuth && !auth.token) {
    next('/admin/login')
  } else if (to.meta.guest && auth.token) {
    next('/admin/dashboard')
  } else {
    next()
  }
})

export default router
