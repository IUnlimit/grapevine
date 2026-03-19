<template>
  <v-app :theme="appStore.darkMode ? 'grapeDark' : 'grapePurple'">
    <router-view />
  </v-app>
</template>

<script setup lang="ts">
import { useAppStore } from './stores/app'
const appStore = useAppStore()
</script>
