<template>
  <div class="gv-auth-bg">
    <div class="setup-container">
      <div class="setup-card gv-glass">
        <!-- Logo area -->
        <div class="logo-section">
          <div class="logo-icon">
            <svg width="48" height="48" viewBox="0 0 48 48" fill="none">
              <circle cx="24" cy="14" r="5" fill="url(#gs1)" opacity="0.9"/>
              <circle cx="16" cy="22" r="4" fill="url(#gs1)" opacity="0.7"/>
              <circle cx="32" cy="22" r="4" fill="url(#gs1)" opacity="0.7"/>
              <circle cx="20" cy="30" r="3.5" fill="url(#gs1)" opacity="0.5"/>
              <circle cx="28" cy="30" r="3.5" fill="url(#gs1)" opacity="0.5"/>
              <circle cx="24" cy="37" r="3" fill="url(#gs1)" opacity="0.4"/>
              <path d="M24 4 C24 4 26 8 24 14" stroke="url(#gs2)" stroke-width="1.5" fill="none" stroke-linecap="round"/>
              <defs>
                <linearGradient id="gs1" x1="0" y1="0" x2="1" y2="1">
                  <stop offset="0%" stop-color="#34D399"/>
                  <stop offset="100%" stop-color="#10B981"/>
                </linearGradient>
                <linearGradient id="gs2" x1="0" y1="0" x2="0" y2="1">
                  <stop offset="0%" stop-color="#A78BFA"/>
                  <stop offset="100%" stop-color="#7C3AED"/>
                </linearGradient>
              </defs>
            </svg>
          </div>
          <h1 class="brand-name">初始化设置</h1>
          <p class="brand-subtitle">创建管理员账号以开始使用 Grapevine</p>
        </div>

        <!-- Steps indicator -->
        <div class="steps-indicator">
          <div class="step active">
            <div class="step-dot"></div>
            <span>创建账号</span>
          </div>
          <div class="step-line"></div>
          <div class="step">
            <div class="step-dot"></div>
            <span>开始使用</span>
          </div>
        </div>

        <!-- Form -->
        <v-form @submit.prevent="handleSetup" class="setup-form">
          <div class="field-group">
            <label class="field-label">管理员用户名</label>
            <v-text-field
              v-model="form.username"
              placeholder="设置管理员用户名"
              prepend-inner-icon="mdi-account-outline"
              :error-messages="usernameError"
              hide-details="auto"
              bg-color="transparent"
            />
          </div>
          <div class="field-group">
            <label class="field-label">密码</label>
            <v-text-field
              v-model="form.password"
              placeholder="设置密码（至少 6 位）"
              type="password"
              prepend-inner-icon="mdi-lock-outline"
              hide-details="auto"
              bg-color="transparent"
            />
          </div>
          <div class="field-group">
            <label class="field-label">确认密码</label>
            <v-text-field
              v-model="form.confirmPassword"
              placeholder="再次输入密码"
              type="password"
              prepend-inner-icon="mdi-lock-check-outline"
              :error-messages="error"
              hide-details="auto"
              bg-color="transparent"
            />
          </div>
          <v-btn
            type="submit"
            block
            size="x-large"
            :loading="loading"
            rounded="lg"
            class="setup-btn mt-2"
            elevation="0"
          >
            <v-icon start>mdi-rocket-launch-outline</v-icon>
            <span>创建并开始</span>
          </v-btn>
        </v-form>

        <div class="setup-footer">
          <v-icon size="14" color="medium-emphasis">mdi-information-outline</v-icon>
          <span>此页面仅在首次使用时出现</span>
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

const form = reactive({ username: '', password: '', confirmPassword: '' })
const loading = ref(false)
const error = ref('')
const usernameError = ref('')

async function handleSetup() {
  usernameError.value = ''
  error.value = ''

  if (!form.username.trim()) {
    usernameError.value = '请输入用户名'
    return
  }
  if (form.password.length < 6) {
    error.value = '密码至少 6 位'
    return
  }
  if (form.password !== form.confirmPassword) {
    error.value = '两次密码不一致'
    return
  }

  loading.value = true
  try {
    await authStore.setup(form.username, form.password)
    router.push('/admin/dashboard')
  } catch {
    error.value = '创建失败，请重试'
  } finally {
    loading.value = false
  }
}
</script>

<style scoped>
.setup-container {
  width: 100%;
  max-width: 440px;
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

.setup-card {
  padding: 48px 40px;
  border-radius: 24px;
  position: relative;
  overflow: hidden;
}

.logo-section {
  text-align: center;
  margin-bottom: 32px;
}

.logo-icon {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  width: 80px;
  height: 80px;
  border-radius: 20px;
  background: rgba(16, 185, 129, 0.08);
  border: 1px solid rgba(52, 211, 153, 0.15);
  margin-bottom: 16px;
}

.brand-name {
  font-size: 26px;
  font-weight: 700;
  letter-spacing: -0.5px;
  background: linear-gradient(135deg, #34D399 0%, #6EE7B7 100%);
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
  margin-top: 6px;
}

.steps-indicator {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 12px;
  margin-bottom: 32px;
}

.step {
  display: flex;
  align-items: center;
  gap: 8px;
  font-size: 12px;
  font-weight: 500;
  opacity: 0.35;
}
.step.active {
  opacity: 1;
}
.step.active .step-dot {
  background: #34D399;
  box-shadow: 0 0 8px rgba(52, 211, 153, 0.5);
}

.step-dot {
  width: 8px;
  height: 8px;
  border-radius: 50%;
  background: rgba(255, 255, 255, 0.2);
}

.step-line {
  width: 40px;
  height: 1px;
  background: rgba(255, 255, 255, 0.1);
}

.setup-form {
  display: flex;
  flex-direction: column;
  gap: 18px;
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

.setup-btn {
  margin-top: 8px;
  height: 52px !important;
  font-size: 15px !important;
  font-weight: 600 !important;
  background: linear-gradient(135deg, #10B981 0%, #34D399 100%) !important;
  color: #0F0B1A !important;
  box-shadow: 0 4px 24px rgba(16, 185, 129, 0.3) !important;
  transition: box-shadow 0.25s ease, transform 0.25s ease !important;
}
.setup-btn:hover {
  box-shadow: 0 8px 32px rgba(16, 185, 129, 0.45) !important;
  transform: translateY(-1px);
}

.setup-footer {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 6px;
  margin-top: 28px;
  font-size: 11px;
  opacity: 0.35;
  letter-spacing: 0.02em;
}
</style>
