<template>
  <div class="endpoints-page">
    <!-- Page header -->
    <div class="page-header">
      <div>
        <h1 class="page-title">节点调度</h1>
        <p class="page-desc">管理后端服务节点与优先级</p>
      </div>
      <div class="header-actions">
        <v-select
          v-model="selectedService"
          :items="serviceOptions"
          item-title="text"
          item-value="value"
          label="筛选服务"
          density="compact"
          variant="outlined"
          hide-details
          clearable
          bg-color="transparent"
          style="max-width: 220px"
          class="filter-select"
        />
        <v-btn
          color="primary"
          prepend-icon="mdi-plus"
          @click="openDialog()"
          elevation="0"
          class="add-btn"
        >
          新增节点
        </v-btn>
      </div>
    </div>

    <!-- Endpoints grid -->
    <div v-if="endpoints.length" class="endpoints-grid">
      <div
        v-for="ep in endpoints"
        :key="ep.id"
        class="endpoint-card gv-stat-card"
        :class="{ 'endpoint-card--offline': ep.status !== 'online' }"
      >
        <div class="ep-header">
          <div class="ep-status">
            <span class="gv-dot" :class="ep.status === 'online' ? 'gv-dot--online' : 'gv-dot--offline'"></span>
            <span class="ep-status-text" :class="ep.status === 'online' ? 'text-success' : 'text-error'">
              {{ ep.status === 'online' ? '在线' : '离线' }}
            </span>
          </div>
          <v-chip size="x-small" :color="ep.priority <= 1 ? 'primary' : 'default'" variant="tonal">
            P{{ ep.priority }}
          </v-chip>
        </div>

        <div class="ep-address">
          <code>{{ ep.host }}:{{ ep.port }}</code>
        </div>

        <div class="ep-service">
          <v-icon size="14" class="mr-1" style="opacity: 0.4">mdi-tag-outline</v-icon>
          <span>{{ getServiceName(ep.service_id) }}</span>
        </div>

        <div class="ep-actions">
          <v-btn
            v-if="ep.status === 'online'"
            size="small"
            variant="tonal"
            color="warning"
            @click="toggleStatus(ep, 'offline')"
            class="ep-action-btn"
          >
            <v-icon start size="16">mdi-power-plug-off</v-icon>
            下线
          </v-btn>
          <v-btn
            v-else
            size="small"
            variant="tonal"
            color="success"
            @click="toggleStatus(ep, 'online')"
            class="ep-action-btn"
          >
            <v-icon start size="16">mdi-power-plug</v-icon>
            上线
          </v-btn>
          <v-btn icon="mdi-pencil-outline" size="x-small" variant="text" @click="openDialog(ep)" class="action-btn" />
          <v-btn icon="mdi-delete-outline" size="x-small" variant="text" color="error" @click="confirmDelete(ep)" class="action-btn" />
        </div>
      </div>
    </div>

    <!-- Empty state -->
    <v-card v-else class="gv-glass pa-12 text-center">
      <v-icon size="48" color="medium-emphasis" class="mb-3">mdi-server-network</v-icon>
      <p style="opacity: 0.4; font-size: 14px">暂无节点</p>
      <v-btn
        variant="tonal"
        color="primary"
        size="small"
        prepend-icon="mdi-plus"
        @click="openDialog()"
        class="mt-3"
      >
        添加第一个节点
      </v-btn>
    </v-card>

    <!-- Create/Edit dialog -->
    <v-dialog v-model="dialog" max-width="480" persistent>
      <v-card class="dialog-card">
        <div class="dialog-header">
          <h3 class="dialog-title">{{ editing ? '编辑节点' : '新增节点' }}</h3>
          <v-btn icon="mdi-close" size="small" variant="text" @click="dialog = false" />
        </div>
        <v-divider style="opacity: 0.06" />
        <div class="dialog-body">
          <v-form @submit.prevent="handleSave">
            <div class="field-group">
              <label class="field-label">所属服务</label>
              <v-select
                v-model="form.service_id"
                :items="serviceOptions"
                item-title="text"
                item-value="value"
                hide-details
                bg-color="transparent"
              />
            </div>
            <div class="field-row">
              <div class="field-group" style="flex: 2">
                <label class="field-label">主机地址</label>
                <v-text-field v-model="form.host" placeholder="127.0.0.1" hide-details bg-color="transparent" />
              </div>
              <div class="field-group" style="flex: 1">
                <label class="field-label">端口</label>
                <v-text-field v-model.number="form.port" type="number" placeholder="8080" hide-details bg-color="transparent" />
              </div>
            </div>
            <div class="field-group">
              <label class="field-label">优先级（值越小越优先）</label>
              <v-text-field v-model.number="form.priority" type="number" placeholder="0" hide-details bg-color="transparent" />
            </div>
            <div class="dialog-actions">
              <v-btn variant="text" @click="dialog = false">取消</v-btn>
              <v-btn type="submit" color="primary" :loading="saving" elevation="0">
                {{ editing ? '更新' : '创建' }}
              </v-btn>
            </div>
          </v-form>
        </div>
      </v-card>
    </v-dialog>

    <!-- Delete confirm -->
    <v-dialog v-model="deleteDialog" max-width="400">
      <v-card class="dialog-card">
        <div class="dialog-body" style="padding-top: 32px">
          <div class="text-center mb-4">
            <div class="delete-icon-wrap">
              <v-icon size="28" color="error">mdi-delete-alert-outline</v-icon>
            </div>
          </div>
          <h3 class="text-center mb-2" style="font-size: 16px; font-weight: 600">确认删除</h3>
          <p class="text-center mb-6" style="font-size: 13px; opacity: 0.5">
            确定要删除节点「{{ deleteTarget?.host }}:{{ deleteTarget?.port }}」吗？
          </p>
          <div class="dialog-actions">
            <v-btn variant="text" @click="deleteDialog = false">取消</v-btn>
            <v-btn color="error" :loading="deleting" @click="handleDelete" elevation="0">删除</v-btn>
          </div>
        </div>
      </v-card>
    </v-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, computed, onMounted, watch } from 'vue'
import api from '../api'

interface Endpoint {
  id?: number
  service_id: number
  host: string
  port: number
  priority: number
  status: string
}

interface Service { id: number; suffix: string; name: string }

const endpoints = ref<Endpoint[]>([])
const services = ref<Service[]>([])
const selectedService = ref<number | null>(null)
const dialog = ref(false)
const deleteDialog = ref(false)
const editing = ref(false)
const saving = ref(false)
const deleting = ref(false)
const editId = ref<number | null>(null)
const deleteTarget = ref<Endpoint | null>(null)

const form = reactive<Endpoint>({ service_id: 0, host: '127.0.0.1', port: 8080, priority: 0, status: 'online' })

const serviceOptions = computed(() =>
  services.value.map((s) => ({ text: `${s.name} (/${s.suffix})`, value: s.id }))
)

function getServiceName(serviceId: number) {
  const svc = services.value.find((s) => s.id === serviceId)
  return svc ? svc.name : `#${serviceId}`
}

async function fetchServices() {
  const res = await api.get('/admin/services')
  services.value = res.data
}

async function fetchEndpoints() {
  const params: Record<string, string> = {}
  if (selectedService.value) params.service_id = String(selectedService.value)
  const res = await api.get('/admin/endpoints', { params })
  endpoints.value = res.data
}

function openDialog(ep?: Endpoint) {
  if (ep) {
    editing.value = true
    editId.value = ep.id!
    Object.assign(form, { service_id: ep.service_id, host: ep.host, port: ep.port, priority: ep.priority, status: ep.status })
  } else {
    editing.value = false
    editId.value = null
    Object.assign(form, { service_id: services.value[0]?.id ?? 0, host: '127.0.0.1', port: 8080, priority: 0, status: 'online' })
  }
  dialog.value = true
}

async function handleSave() {
  saving.value = true
  try {
    if (editing.value && editId.value) {
      await api.put(`/admin/endpoints/${editId.value}`, form)
    } else {
      await api.post('/admin/endpoints', form)
    }
    dialog.value = false
    await fetchEndpoints()
  } finally {
    saving.value = false
  }
}

async function toggleStatus(ep: Endpoint, status: string) {
  await api.put(`/admin/endpoints/${ep.id}`, { ...ep, status })
  await fetchEndpoints()
}

function confirmDelete(ep: Endpoint) {
  deleteTarget.value = ep
  deleteDialog.value = true
}

async function handleDelete() {
  if (!deleteTarget.value?.id) return
  deleting.value = true
  try {
    await api.delete(`/admin/endpoints/${deleteTarget.value.id}`)
    deleteDialog.value = false
    await fetchEndpoints()
  } finally {
    deleting.value = false
  }
}

watch(selectedService, fetchEndpoints)

onMounted(async () => {
  await fetchServices()
  await fetchEndpoints()
})
</script>

<style scoped>
.endpoints-page {
  display: flex;
  flex-direction: column;
  gap: 20px;
}

.page-header {
  display: flex;
  align-items: flex-start;
  justify-content: space-between;
  flex-wrap: wrap;
  gap: 12px;
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

.header-actions {
  display: flex;
  align-items: center;
  gap: 12px;
}

.add-btn {
  box-shadow: 0 2px 12px rgba(124, 58, 237, 0.25) !important;
}

/* Endpoints grid */
.endpoints-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(280px, 1fr));
  gap: 14px;
}

.endpoint-card {
  padding: 20px;
  border-radius: 16px;
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.endpoint-card--offline {
  opacity: 0.65;
}

.ep-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
}

.ep-status {
  display: flex;
  align-items: center;
  gap: 8px;
}

.ep-status-text {
  font-size: 12px;
  font-weight: 600;
  text-transform: uppercase;
  letter-spacing: 0.04em;
}

.ep-address code {
  font-family: 'JetBrains Mono', monospace !important;
  font-size: 16px;
  font-weight: 500;
  letter-spacing: -0.3px;
}

.ep-service {
  display: flex;
  align-items: center;
  font-size: 12px;
  opacity: 0.5;
}

.ep-actions {
  display: flex;
  align-items: center;
  gap: 6px;
  margin-top: 4px;
  padding-top: 12px;
  border-top: 1px solid rgba(255, 255, 255, 0.04);
}

.ep-action-btn {
  flex: 1;
}

.action-btn {
  opacity: 0.4;
  transition: opacity 0.15s ease;
}
.action-btn:hover {
  opacity: 1;
}

/* Dialog styles */
.dialog-card {
  border-radius: 20px !important;
  overflow: hidden;
}

.dialog-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 20px 24px;
}

.dialog-title {
  font-size: 16px;
  font-weight: 600;
  margin: 0;
}

.dialog-body {
  padding: 24px;
}

.field-group {
  margin-bottom: 16px;
}

.field-row {
  display: flex;
  gap: 12px;
}

.field-label {
  display: block;
  font-size: 12px;
  font-weight: 600;
  text-transform: uppercase;
  letter-spacing: 0.06em;
  opacity: 0.5;
  margin-bottom: 6px;
  padding-left: 4px;
}

.dialog-actions {
  display: flex;
  justify-content: flex-end;
  gap: 8px;
  margin-top: 24px;
}

.delete-icon-wrap {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  width: 56px;
  height: 56px;
  border-radius: 16px;
  background: rgba(239, 68, 68, 0.1);
}
</style>
