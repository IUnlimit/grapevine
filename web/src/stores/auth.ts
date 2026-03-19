import { defineStore } from 'pinia'
import { ref } from 'vue'
import api from '../api'

export const useAuthStore = defineStore('auth', () => {
  const token = ref(localStorage.getItem('gv_token') || '')
  const username = ref(localStorage.getItem('gv_username') || '')

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

  return { token, username, login, logout }
})
