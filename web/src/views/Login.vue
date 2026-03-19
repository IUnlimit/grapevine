<template>
  <v-container class="fill-height" fluid>
    <v-row align="center" justify="center">
      <v-col cols="12" sm="8" md="4">
        <v-card class="pa-8" elevation="8" rounded="xl">
          <div class="text-center mb-6">
            <v-icon icon="mdi-fruit-grapes" size="64" color="primary" />
            <h1 class="text-h5 font-weight-bold mt-2">Grapevine</h1>
            <p class="text-body-2 text-medium-emphasis">API Gateway 管理后台</p>
          </div>

          <v-form @submit.prevent="handleLogin">
            <v-text-field
              v-model="form.username"
              label="用户名"
              prepend-inner-icon="mdi-account"
              :error-messages="error ? ' ' : ''"
              class="mb-2"
            />
            <v-text-field
              v-model="form.password"
              label="密码"
              type="password"
              prepend-inner-icon="mdi-lock"
              :error-messages="error"
              class="mb-4"
            />
            <v-btn type="submit" block size="large" color="primary" :loading="loading" rounded="lg">
              登 录
            </v-btn>
          </v-form>
        </v-card>
      </v-col>
    </v-row>
  </v-container>
</template>

<script setup lang="ts">
import { reactive, ref } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '../stores/auth'

const authStore = useAuthStore()
const router = useRouter()

const form = reactive({ username: '', password: '' })
const loading = ref(false)
const error = ref('')

async function handleLogin() {
  loading.value = true
  error.value = ''
  try {
    await authStore.login(form.username, form.password)
    router.push('/dashboard')
  } catch {
    error.value = '用户名或密码错误'
  } finally {
    loading.value = false
  }
}
</script>
