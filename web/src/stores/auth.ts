import { defineStore } from 'pinia'
import { ref } from 'vue'
import api from '../api'

export const useAuthStore = defineStore('auth', () => {
  const token = ref(localStorage.getItem('gv_token') || '')
  const username = ref(localStorage.getItem('gv_username') || '')
  const needSetup = ref(false)
  const setupChecked = ref(false)

  async function checkSetup() {
    try {
      const res = await api.get('/admin/setup/status')
      needSetup.value = res.data.need_setup
      // 需要初始化时，清除可能残留的过期 token
      if (needSetup.value) {
        token.value = ''
        username.value = ''
        localStorage.removeItem('gv_token')
        localStorage.removeItem('gv_username')
      }
    } catch {
      needSetup.value = false
    }
    setupChecked.value = true
  }

  async function setup(user: string, password: string) {
    const res = await api.post('/admin/setup', { username: user, password })
    token.value = res.data.token
    username.value = res.data.username
    needSetup.value = false
    localStorage.setItem('gv_token', token.value)
    localStorage.setItem('gv_username', username.value)
  }

  async function login(user: string, password: string) {
    const res = await api.post('/admin/login', { username: user, password })
    token.value = res.data.token
    username.value = res.data.username
    localStorage.setItem('gv_token', token.value)
    localStorage.setItem('gv_username', username.value)
  }

  function logout() {
    token.value = ''
    username.value = ''
    localStorage.removeItem('gv_token')
    localStorage.removeItem('gv_username')
  }

  return { token, username, needSetup, setupChecked, checkSetup, setup, login, logout }
})
