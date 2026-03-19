<template>
  <v-layout>
    <v-navigation-drawer v-model="appStore.sidebarOpen" :rail="!appStore.sidebarOpen" permanent color="primary-darken-1">
      <v-list-item
        prepend-icon="mdi-fruit-grapes"
        title="Grapevine"
        subtitle="API Gateway"
        class="my-2"
      />
      <v-divider />
      <v-list nav density="comfortable">
        <v-list-item
          v-for="item in navItems"
          :key="item.path"
          :to="item.path"
          :prepend-icon="item.icon"
          :title="item.title"
          rounded="lg"
          class="mx-2 my-1"
        />
      </v-list>
      <template #append>
        <v-list nav density="comfortable">
          <v-list-item
            prepend-icon="mdi-logout"
            title="退出登录"
            rounded="lg"
            class="mx-2 my-1"
            @click="handleLogout"
          />
        </v-list>
      </template>
    </v-navigation-drawer>

    <v-app-bar flat density="compact" color="surface">
      <v-app-bar-nav-icon @click="appStore.toggleSidebar" />
      <v-toolbar-title class="text-body-1 font-weight-medium">
        {{ currentTitle }}
      </v-toolbar-title>
      <v-spacer />
      <v-btn :icon="appStore.darkMode ? 'mdi-weather-sunny' : 'mdi-weather-night'" variant="text" @click="appStore.toggleDark" />
      <v-chip class="mr-4" prepend-icon="mdi-account" variant="tonal" color="primary">
        {{ authStore.username }}
      </v-chip>
    </v-app-bar>

    <v-main>
      <v-container fluid class="pa-6">
        <router-view />
      </v-container>
    </v-main>
  </v-layout>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useAppStore } from '../stores/app'
import { useAuthStore } from '../stores/auth'

const appStore = useAppStore()
const authStore = useAuthStore()
const route = useRoute()
const router = useRouter()

const navItems = [
  { path: '/dashboard', icon: 'mdi-view-dashboard', title: '仪表盘' },
  { path: '/services', icon: 'mdi-routes', title: '路由管理' },
  { path: '/endpoints', icon: 'mdi-server-network', title: '节点调度' },
  { path: '/policies', icon: 'mdi-shield-check', title: '策略中心' },
]

const currentTitle = computed(() => {
  const item = navItems.find((n) => n.path === route.path)
  return item?.title ?? 'Grapevine'
})

function handleLogout() {
  authStore.logout()
  router.push('/login')
}
</script>
