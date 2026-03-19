import { defineStore } from 'pinia'
import { ref } from 'vue'

export const useAppStore = defineStore('app', () => {
  const darkMode = ref(false)
  const sidebarOpen = ref(true)

  function toggleDark() {
    darkMode.value = !darkMode.value
  }

  function toggleSidebar() {
    sidebarOpen.value = !sidebarOpen.value
  }

  return { darkMode, sidebarOpen, toggleDark, toggleSidebar }
})
