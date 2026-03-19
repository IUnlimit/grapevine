<template>
  <v-layout>
    <!-- Sidebar -->
    <v-navigation-drawer
      v-model="appStore.sidebarOpen"
      :rail="!appStore.sidebarOpen"
      permanent
      :color="appStore.darkMode ? '#120E1F' : '#F8F6FC'"
      :width="260"
    >
      <!-- Brand -->
      <div class="sidebar-brand" :class="{ 'sidebar-brand--rail': !appStore.sidebarOpen }">
        <div class="brand-icon">
          <svg width="28" height="28" viewBox="0 0 48 48" fill="none">
            <circle cx="24" cy="14" r="5" fill="url(#nav1)" opacity="0.9"/>
            <circle cx="16" cy="22" r="4" fill="url(#nav1)" opacity="0.7"/>
            <circle cx="32" cy="22" r="4" fill="url(#nav1)" opacity="0.7"/>
            <circle cx="20" cy="30" r="3.5" fill="url(#nav1)" opacity="0.5"/>
            <circle cx="28" cy="30" r="3.5" fill="url(#nav1)" opacity="0.5"/>
            <defs>
              <linearGradient id="nav1" x1="0" y1="0" x2="1" y2="1">
                <stop offset="0%" stop-color="#A78BFA"/>
                <stop offset="100%" stop-color="#7C3AED"/>
              </linearGradient>
            </defs>
          </svg>
        </div>
        <transition name="fade">
          <div v-if="appStore.sidebarOpen" class="brand-text">
            <span class="brand-title">Grapevine</span>
            <span class="brand-version">v1.0</span>
          </div>
        </transition>
      </div>

      <v-divider class="mx-4 my-2" style="opacity: 0.08" />

      <!-- Navigation -->
      <div class="sidebar-nav">
        <div
          v-if="appStore.sidebarOpen"
          class="nav-section-label"
        >
          导航
        </div>
        <v-list nav density="comfortable" class="px-3">
          <v-list-item
            v-for="item in navItems"
            :key="item.path"
            :to="item.path"
            :prepend-icon="item.icon"
            :title="item.title"
            rounded="lg"
            class="nav-item mb-1"
            :class="{ 'nav-item--active': route.path === item.path }"
          >
            <template #prepend>
              <div class="nav-icon-wrap" :class="{ 'nav-icon-wrap--active': route.path === item.path }">
                <v-icon :icon="item.icon" size="20" />
              </div>
            </template>
          </v-list-item>
        </v-list>
      </div>

      <!-- Bottom -->
      <template #append>
        <v-divider class="mx-4 mb-2" style="opacity: 0.08" />
        <v-list nav density="comfortable" class="px-3 pb-4">
          <v-list-item
            prepend-icon="mdi-logout"
            title="退出登录"
            rounded="lg"
            class="nav-item nav-item--logout"
            @click="handleLogout"
          >
            <template #prepend>
              <div class="nav-icon-wrap">
                <v-icon icon="mdi-logout" size="20" />
              </div>
            </template>
          </v-list-item>
        </v-list>
      </template>
    </v-navigation-drawer>

    <!-- Top bar -->
    <v-app-bar flat :color="appStore.darkMode ? '#0F0B1A' : '#F4F2F7'" :height="64" class="topbar">
      <v-btn
        icon
        variant="text"
        size="small"
        @click="appStore.toggleSidebar"
        class="ml-2"
      >
        <v-icon>{{ appStore.sidebarOpen ? 'mdi-menu-open' : 'mdi-menu' }}</v-icon>
      </v-btn>

      <div class="topbar-title">
        <span class="topbar-title-text">{{ currentTitle }}</span>
      </div>

      <v-spacer />

      <!-- Theme toggle -->
      <v-btn
        icon
        variant="text"
        size="small"
        @click="appStore.toggleDark"
        class="mr-1"
      >
        <v-icon size="20">{{ appStore.darkMode ? 'mdi-white-balance-sunny' : 'mdi-moon-waning-crescent' }}</v-icon>
      </v-btn>

      <!-- User chip -->
      <div class="user-chip mr-4">
        <div class="user-avatar">
          <v-icon size="16" color="white">mdi-account</v-icon>
        </div>
        <span class="user-name">{{ authStore.username }}</span>
      </div>
    </v-app-bar>

    <!-- Main content -->
    <v-main>
      <div class="main-content">
        <router-view v-slot="{ Component }">
          <transition name="slide-up" mode="out-in">
            <component :is="Component" />
          </transition>
        </router-view>
      </div>
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
  { path: '/admin/dashboard', icon: 'mdi-view-dashboard-outline', title: '仪表盘' },
  { path: '/admin/services', icon: 'mdi-routes', title: '路由管理' },
  { path: '/admin/endpoints', icon: 'mdi-server-network', title: '节点调度' },
  { path: '/admin/policies', icon: 'mdi-shield-check-outline', title: '策略中心' },
]

const currentTitle = computed(() => {
  const item = navItems.find((n) => n.path === route.path)
  return item?.title ?? 'Grapevine'
})

function handleLogout() {
  authStore.logout()
  router.push('/admin/login')
}
</script>

<style scoped>
.sidebar-brand {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 20px 20px 8px;
  transition: padding 0.2s ease;
}
.sidebar-brand--rail {
  padding: 20px 12px 8px;
  justify-content: center;
}

.brand-icon {
  width: 40px;
  height: 40px;
  border-radius: 12px;
  background: rgba(124, 58, 237, 0.08);
  border: 1px solid rgba(167, 139, 250, 0.12);
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
}

.brand-text {
  display: flex;
  flex-direction: column;
  min-width: 0;
}

.brand-title {
  font-size: 17px;
  font-weight: 700;
  letter-spacing: -0.3px;
  line-height: 1.2;
}

.brand-version {
  font-size: 11px;
  opacity: 0.35;
  font-weight: 500;
}

.nav-section-label {
  font-size: 10px;
  font-weight: 700;
  text-transform: uppercase;
  letter-spacing: 0.1em;
  opacity: 0.3;
  padding: 16px 24px 8px;
}

.sidebar-nav {
  flex: 1;
}

.nav-item {
  transition: all 0.2s ease !important;
  margin-bottom: 2px;
}

.nav-icon-wrap {
  width: 36px;
  height: 36px;
  border-radius: 10px;
  display: flex;
  align-items: center;
  justify-content: center;
  transition: all 0.2s ease;
  margin-right: 4px;
}

.nav-icon-wrap--active {
  background: linear-gradient(135deg, rgba(124, 58, 237, 0.15), rgba(167, 139, 250, 0.1));
}

.nav-item--active {
  background: rgba(124, 58, 237, 0.06) !important;
}

.nav-item--logout:hover {
  color: #EF4444 !important;
}

/* Top bar */
.topbar {
  border-bottom: 1px solid rgba(255, 255, 255, 0.04) !important;
}

.topbar-title {
  margin-left: 8px;
}

.topbar-title-text {
  font-size: 15px;
  font-weight: 600;
  letter-spacing: -0.2px;
}

.user-chip {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 6px 14px 6px 6px;
  border-radius: 20px;
  background: rgba(124, 58, 237, 0.08);
  border: 1px solid rgba(167, 139, 250, 0.1);
}

.user-avatar {
  width: 28px;
  height: 28px;
  border-radius: 50%;
  background: linear-gradient(135deg, #7C3AED, #A78BFA);
  display: flex;
  align-items: center;
  justify-content: center;
}

.user-name {
  font-size: 13px;
  font-weight: 500;
}

/* Main content */
.main-content {
  padding: 28px 32px;
  max-width: 1400px;
}

@media (max-width: 960px) {
  .main-content {
    padding: 16px;
  }
}
</style>
