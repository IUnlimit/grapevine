import { defineStore } from 'pinia'
import { ref } from 'vue'

export const useAppStore = defineStore('app', () => {
  const darkMode = ref(localStorage.getItem('gv_dark') !== 'false')
  const sidebarOpen = ref(true)

  function toggleDark() {
    darkMode.value = !darkMode.value
    localStorage.setItem('gv_dark', String(darkMode.value))
  }

  function toggleSidebar() {
    sidebarOpen.value = !sidebarOpen.value
  }

  return { darkMode, sidebarOpen, toggleDark, toggleSidebar }
})
