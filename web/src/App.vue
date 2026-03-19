<template>
  <v-app :theme="appStore.darkMode ? 'grapeDark' : 'grapeLight'" class="gv-noise">
    <router-view v-slot="{ Component }">
      <transition name="fade" mode="out-in">
        <component :is="Component" />
      </transition>
    </router-view>
  </v-app>
</template>

<script setup lang="ts">
import { useAppStore } from './stores/app'
const appStore = useAppStore()
</script>
