<template>
  <div>
    <div class="text-h6 font-weight-medium mb-4">策略中心</div>

    <v-row v-if="!services.length">
      <v-col>
        <v-card class="pa-8 text-center">
          <v-icon icon="mdi-shield-off" size="48" color="medium-emphasis" class="mb-2" />
          <p class="text-medium-emphasis">暂无服务，请先在路由管理中创建服务</p>
        </v-card>
      </v-col>
    </v-row>

    <v-row>
      <v-col v-for="svc in services" :key="svc.id" cols="12" md="6">
        <v-card class="pa-6">
          <div class="d-flex align-center mb-4">
            <v-chip color="primary" variant="tonal" class="mr-2">/{{ svc.suffix }}</v-chip>
            <span class="text-subtitle-1 font-weight-medium">{{ svc.name }}</span>
          </div>

          <div class="mb-4">
            <div class="d-flex justify-space-between mb-1">
              <span class="text-body-2">限流 QPS</span>
              <span class="text-body-2 font-weight-medium">{{ getPolicy(svc.id).rate_limit }}</span>
            </div>
            <v-slider
              :model-value="getPolicy(svc.id).rate_limit"
              @update:model-value="(v: number) => updateField(svc.id, 'rate_limit', v)"
              :min="0" :max="1000" :step="10"
              color="primary" thumb-label hide-details
            />
          </div>

          <div class="mb-4">
            <div class="d-flex justify-space-between mb-1">
              <span class="text-body-2">突发容量 (Burst)</span>
              <span class="text-body-2 font-weight-medium">{{ getPolicy(svc.id).burst }}</span>
            </div>
            <v-slider
              :model-value="getPolicy(svc.id).burst"
              @update:model-value="(v: number) => updateField(svc.id, 'burst', v)"
              :min="0" :max="2000" :step="10"
              color="secondary" thumb-label hide-details
            />
          </div>

          <div class="mb-4">
            <div class="d-flex justify-space-between mb-1">
              <span class="text-body-2">熔断阈值 (错误率 %)</span>
              <span class="text-body-2 font-weight-medium">{{ getPolicy(svc.id).circuit_break_threshold }}%</span>
            </div>
            <v-slider
              :model-value="getPolicy(svc.id).circuit_break_threshold"
              @update:model-value="(v: number) => updateField(svc.id, 'circuit_break_threshold', v)"
              :min="0" :max="100" :step="5"
              color="error" thumb-label hide-details
            />
          </div>

          <div class="d-flex justify-end">
            <v-btn color="primary" variant="tonal" :loading="savingMap[svc.id]" @click="savePolicy(svc.id)">
              保存策略
            </v-btn>
          </div>
        </v-card>
      </v-col>
    </v-row>
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
