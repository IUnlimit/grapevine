<template>
  <div class="policies-page">
    <!-- Page header -->
    <div class="page-header">
      <div>
        <h1 class="page-title">策略中心</h1>
        <p class="page-desc">配置限流、突发容量与熔断阈值</p>
      </div>
    </div>

    <!-- Empty state -->
    <v-card v-if="!services.length" class="gv-glass pa-12 text-center">
      <v-icon size="48" color="medium-emphasis" class="mb-3">mdi-shield-off-outline</v-icon>
      <p style="opacity: 0.4; font-size: 14px">暂无服务，请先在路由管理中创建服务</p>
    </v-card>

    <!-- Policy cards -->
    <div class="policies-grid">
      <div v-for="svc in services" :key="svc.id" class="policy-card gv-glass">
        <!-- Card header -->
        <div class="policy-header">
          <div class="policy-service">
            <div class="suffix-badge">
              <span class="suffix-slash">/</span>{{ svc.suffix }}
            </div>
            <span class="policy-name">{{ svc.name }}</span>
          </div>
          <v-btn
            variant="tonal"
            color="primary"
            size="small"
            :loading="savingMap[svc.id]"
            @click="savePolicy(svc.id)"
            elevation="0"
          >
            <v-icon start size="16">mdi-content-save-outline</v-icon>
            保存
          </v-btn>
        </div>

        <v-divider class="my-4" style="opacity: 0.06" />

        <!-- Rate limit -->
        <div class="policy-control">
          <div class="control-header">
            <div class="control-icon control-icon--purple">
              <v-icon size="16" color="white">mdi-speedometer</v-icon>
            </div>
            <div class="control-info">
              <span class="control-label">限流 QPS</span>
              <span class="control-desc">每秒最大请求数</span>
            </div>
            <div class="control-value control-value--purple">
              {{ getPolicy(svc.id).rate_limit }}
            </div>
          </div>
          <v-slider
            :model-value="getPolicy(svc.id).rate_limit"
            @update:model-value="(v: number) => updateField(svc.id, 'rate_limit', v)"
            :min="0" :max="1000" :step="10"
            color="#7C3AED"
            track-color="rgba(124, 58, 237, 0.1)"
            thumb-label
            hide-details
          />
        </div>

        <!-- Burst -->
        <div class="policy-control">
          <div class="control-header">
            <div class="control-icon control-icon--blue">
              <v-icon size="16" color="white">mdi-waves</v-icon>
            </div>
            <div class="control-info">
              <span class="control-label">突发容量</span>
              <span class="control-desc">允许的瞬时突发请求数</span>
            </div>
            <div class="control-value control-value--blue">
              {{ getPolicy(svc.id).burst }}
            </div>
          </div>
          <v-slider
            :model-value="getPolicy(svc.id).burst"
            @update:model-value="(v: number) => updateField(svc.id, 'burst', v)"
            :min="0" :max="2000" :step="10"
            color="#3B82F6"
            track-color="rgba(59, 130, 246, 0.1)"
            thumb-label
            hide-details
          />
        </div>

        <!-- Circuit breaker -->
        <div class="policy-control">
          <div class="control-header">
            <div class="control-icon control-icon--red">
              <v-icon size="16" color="white">mdi-flash-alert</v-icon>
            </div>
            <div class="control-info">
              <span class="control-label">熔断阈值</span>
              <span class="control-desc">错误率超过此值触发熔断</span>
            </div>
            <div class="control-value control-value--red">
              {{ getPolicy(svc.id).circuit_break_threshold }}%
            </div>
          </div>
          <v-slider
            :model-value="getPolicy(svc.id).circuit_break_threshold"
            @update:model-value="(v: number) => updateField(svc.id, 'circuit_break_threshold', v)"
            :min="0" :max="100" :step="5"
            color="#EF4444"
            track-color="rgba(239, 68, 68, 0.1)"
            thumb-label
            hide-details
          />
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import api from '../api'

interface Service { id: number; suffix: string; name: string }
interface Policy { rate_limit: number; burst: number; circuit_break_threshold: number }

const services = ref<Service[]>([])
const policies = reactive<Record<number, Policy>>({})
const savingMap = reactive<Record<number, boolean>>({})

const defaultPolicy: Policy = { rate_limit: 100, burst: 200, circuit_break_threshold: 50 }

function getPolicy(serviceId: number): Policy {
  return policies[serviceId] ?? { ...defaultPolicy }
}

function updateField(serviceId: number, field: keyof Policy, value: number) {
  if (!policies[serviceId]) {
    policies[serviceId] = { ...defaultPolicy }
  }
  policies[serviceId][field] = value
}

async function fetchServices() {
  const res = await api.get('/admin/services')
  services.value = res.data
}

async function fetchPolicies() {
  for (const svc of services.value) {
    try {
      const res = await api.get(`/admin/policies/${svc.id}`)
      policies[svc.id] = res.data
    } catch {
      policies[svc.id] = { ...defaultPolicy }
    }
  }
}

async function savePolicy(serviceId: number) {
  savingMap[serviceId] = true
  try {
    await api.put(`/admin/policies/${serviceId}`, getPolicy(serviceId))
  } finally {
    savingMap[serviceId] = false
  }
}

onMounted(async () => {
  await fetchServices()
  await fetchPolicies()
})
</script>

<style scoped>
.policies-page {
  display: flex;
  flex-direction: column;
  gap: 20px;
}

.page-header {
  display: flex;
  align-items: flex-start;
  justify-content: space-between;
}

.page-title {
  font-size: 24px;
  font-weight: 700;
  letter-spacing: -0.5px;
  margin: 0;
  line-height: 1.2;
}

.page-desc {
  font-size: 13px;
  opacity: 0.45;
  margin-top: 4px;
}

/* Policies grid */
.policies-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(400px, 1fr));
  gap: 16px;
}

@media (max-width: 480px) {
  .policies-grid {
    grid-template-columns: 1fr;
  }
}

.policy-card {
  padding: 24px;
  border-radius: 20px;
}

.policy-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
}

.policy-service {
  display: flex;
  align-items: center;
  gap: 10px;
}

.suffix-badge {
  display: inline-flex;
  align-items: center;
  padding: 3px 10px;
  border-radius: 6px;
  background: rgba(124, 58, 237, 0.08);
  font-size: 13px;
  font-weight: 500;
}

.suffix-slash {
  opacity: 0.4;
  margin-right: 1px;
}

.policy-name {
  font-weight: 600;
  font-size: 15px;
}

/* Policy controls */
.policy-control {
  margin-bottom: 20px;
}
.policy-control:last-child {
  margin-bottom: 0;
}

.control-header {
  display: flex;
  align-items: center;
  gap: 12px;
  margin-bottom: 10px;
}

.control-icon {
  width: 32px;
  height: 32px;
  border-radius: 8px;
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
}
.control-icon--purple { background: linear-gradient(135deg, #7C3AED, #A78BFA); }
.control-icon--blue { background: linear-gradient(135deg, #3B82F6, #60A5FA); }
.control-icon--red { background: linear-gradient(135deg, #EF4444, #F87171); }

.control-info {
  flex: 1;
  display: flex;
  flex-direction: column;
}

.control-label {
  font-size: 13px;
  font-weight: 600;
}

.control-desc {
  font-size: 11px;
  opacity: 0.4;
}

.control-value {
  font-family: 'JetBrains Mono', monospace;
  font-size: 18px;
  font-weight: 700;
  min-width: 60px;
  text-align: right;
}
.control-value--purple {
  color: #A78BFA;
}
.control-value--blue {
  color: #60A5FA;
}
.control-value--red {
  color: #F87171;
}
</style>
