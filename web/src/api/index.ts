import axios from 'axios'
import { useAuthStore } from '../stores/auth'

const api = axios.create({
  baseURL: '/api',
  timeout: 10000,
})

api.interceptors.request.use((config) => {
  const auth = useAuthStore()
  if (auth.token) {
    config.headers.Authorization = `Bearer ${auth.token}`
  }
  return config
})

api.interceptors.response.use(
  (res) => res,
  (err) => {
    if (err.response?.status === 401) {
      // 仅在非登录/setup请求时清除 token 并跳转
      const url = err.config?.url || ''
      if (!url.includes('/login') && !url.includes('/setup')) {
        const auth = useAuthStore()
        auth.logout()
        window.location.href = '/admin/login'
      }
    }
    return Promise.reject(err)
  }
)

export default api
