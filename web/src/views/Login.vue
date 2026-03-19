<template>
  <div class="gv-auth-bg">
    <div class="login-container">
      <div class="login-card gv-glass">
        <!-- Logo area -->
        <div class="logo-section">
          <div class="logo-icon">
            <svg width="48" height="48" viewBox="0 0 48 48" fill="none">
              <circle cx="24" cy="14" r="5" fill="url(#g1)" opacity="0.9"/>
              <circle cx="16" cy="22" r="4" fill="url(#g1)" opacity="0.7"/>
              <circle cx="32" cy="22" r="4" fill="url(#g1)" opacity="0.7"/>
              <circle cx="20" cy="30" r="3.5" fill="url(#g1)" opacity="0.5"/>
              <circle cx="28" cy="30" r="3.5" fill="url(#g1)" opacity="0.5"/>
              <circle cx="24" cy="37" r="3" fill="url(#g1)" opacity="0.4"/>
              <path d="M24 4 C24 4 26 8 24 14" stroke="url(#g2)" stroke-width="1.5" fill="none" stroke-linecap="round"/>
              <defs>
                <linearGradient id="g1" x1="0" y1="0" x2="1" y2="1">
                  <stop offset="0%" stop-color="#A78BFA"/>
                  <stop offset="100%" stop-color="#7C3AED"/>
                </linearGradient>
                <linearGradient id="g2" x1="0" y1="0" x2="0" y2="1">
                  <stop offset="0%" stop-color="#34D399"/>
                  <stop offset="100%" stop-color="#10B981"/>
                </linearGradient>
              </defs>
            </svg>
          </div>
          <h1 class="brand-name">Grapevine</h1>
          <p class="brand-subtitle">API Gateway Console</p>
        </div>

        <!-- Form -->
        <v-form @submit.prevent="handleLogin" class="login-form">
          <div class="field-group">
            <label class="field-label">用户名</label>
            <v-text-field
              v-model="form.username"
              placeholder="输入管理员用户名"
              prepend-inner-icon="mdi-account-outline"
              :error-messages="error ? ' ' : ''"
              hide-details="auto"
              bg-color="transparent"
            />
          </div>
          <div class="field-group">
            <label class="field-label">密码</label>
            <v-text-field
              v-model="form.password"
              placeholder="输入密码"
              type="password"
              prepend-inner-icon="mdi-lock-outline"
              :error-messages="error"
              hide-details="auto"
              bg-color="transparent"
            />
          </div>
          <v-btn
            type="submit"
            block
            size="x-large"
            color="primary"
            :loading="loading"
            rounded="lg"
            class="login-btn mt-2"
            elevation="0"
          >
            <span class="login-btn-text">登 录</span>
            <v-icon end>mdi-arrow-right</v-icon>
          </v-btn>
        </v-form>

        <div class="login-footer">
          <v-icon size="14" color="medium-emphasis">mdi-shield-check-outline</v-icon>
          <span>安全连接已建立</span>
        </div>
      </div>
    </div>
  </div>
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
    router.push('/admin/dashboard')
  } catch {
    error.value = '用户名或密码错误'
  } finally {
    loading.value = false
  }
}
</script>

<style scoped>
.login-container {
  width: 100%;
  max-width: 420px;
  padding: 24px;
  z-index: 1;
  animation: slideUp 0.5s cubic-bezier(0.4, 0, 0.2, 1);
}

@keyframes slideUp {
  from {
    opacity: 0;
    transform: translateY(20px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}

.login-card {
  padding: 48px 40px;
  border-radius: 24px;
  position: relative;
  overflow: hidden;
}

.logo-section {
  text-align: center;
  margin-bottom: 40px;
}

.logo-icon {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  width: 80px;
  height: 80px;
  border-radius: 20px;
  background: rgba(124, 58, 237, 0.08);
  border: 1px solid rgba(167, 139, 250, 0.15);
  margin-bottom: 16px;
}

.brand-name {
  font-size: 28px;
  font-weight: 700;
  letter-spacing: -0.5px;
  background: linear-gradient(135deg, #A78BFA 0%, #C4B5FD 100%);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  background-clip: text;
  margin: 0;
  line-height: 1.2;
}

.brand-subtitle {
  font-size: 13px;
  font-weight: 400;
  opacity: 0.45;
  margin-top: 4px;
  letter-spacing: 0.5px;
  text-transform: uppercase;
}

.login-form {
  display: flex;
  flex-direction: column;
  gap: 20px;
}

.field-group {
  display: flex;
  flex-direction: column;
  gap: 6px;
}

.field-label {
  font-size: 12px;
  font-weight: 600;
  text-transform: uppercase;
  letter-spacing: 0.06em;
  opacity: 0.5;
  padding-left: 4px;
}

.login-btn {
  margin-top: 8px;
  height: 52px !important;
  font-size: 15px !important;
  font-weight: 600 !important;
  letter-spacing: 0.02em !important;
  background: linear-gradient(135deg, #7C3AED 0%, #8B5CF6 100%) !important;
  box-shadow: 0 4px 24px rgba(124, 58, 237, 0.3) !important;
  transition: box-shadow 0.25s ease, transform 0.25s ease !important;
}
.login-btn:hover {
  box-shadow: 0 8px 32px rgba(124, 58, 237, 0.45) !important;
  transform: translateY(-1px);
}

.login-footer {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 6px;
  margin-top: 32px;
  font-size: 11px;
  opacity: 0.35;
  letter-spacing: 0.02em;
}
</style>
