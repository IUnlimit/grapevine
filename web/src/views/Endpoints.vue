<template>
  <div>
    <div class="d-flex align-center mb-4">
      <div class="text-h6 font-weight-medium">节点调度</div>
      <v-spacer />
      <v-select
        v-model="selectedService"
        :items="serviceOptions"
        item-title="text"
        item-value="value"
        label="筛选服务"
        density="compact"
        variant="outlined"
        hide-details
        style="max-width: 240px"
        class="mr-4"
        clearable
      />
      <v-btn color="primary" prepend-icon="mdi-plus" @click="openDialog()">新增节点</v-btn>
    </div>

    <v-card>
      <v-table density="comfortable">
        <thead>
          <tr>
            <th>ID</th>
            <th>所属服务</th>
            <th>地址</th>
            <th>优先级</th>
            <th>状态</th>
            <th class="text-right">操作</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="ep in endpoints" :key="ep.id">
            <td>{{ ep.id }}</td>
            <td>{{ getServiceName(ep.service_id) }}</td>
            <td class="font-weight-medium">{{ ep.host }}:{{ ep.port }}</td>
            <td>
              <v-chip size="small" :color="ep.priority <= 1 ? 'primary' : 'default'" variant="tonal">
                P{{ ep.priority }}
              </v-chip>
            </td>
            <td>
              <v-chip :color="ep.status === 'online' ? 'success' : 'error'" size="small" variant="flat">
                {{ ep.status === 'online' ? '在线' : '离线' }}
              </v-chip>
            </td>
            <td class="text-right">
              <v-btn
                v-if="ep.status === 'online'"
                size="small" variant="tonal" color="warning" class="mr-1"
                @click="toggleStatus(ep, 'offline')"
              >强制下线</v-btn>
              <v-btn
                v-else
                size="small" variant="tonal" color="success" class="mr-1"
                @click="toggleStatus(ep, 'online')"
              >恢复上线</v-btn>
              <v-btn icon="mdi-pencil" size="small" variant="text" @click="openDialog(ep)" />
              <v-btn icon="mdi-delete" size="small" variant="text" color="error" @click="confirmDelete(ep)" />
            </td>
          </tr>
          <tr v-if="!endpoints.length">
            <td colspan="6" class="text-center text-medium-emphasis py-6">暂无节点</td>
          </tr>
        </tbody>
      </v-table>
    </v-card>

    <!-- 新增/编辑对话框 -->
    <v-dialog v-model="dialog" max-width="500" persistent>
      <v-card class="pa-6">
        <div class="text-h6 mb-4">{{ editing ? '编辑节点' : '新增节点' }}</div>
        <v-form @submit.prevent="handleSave">
          <v-select
            v-model="form.service_id"
            :items="serviceOptions"
            item-title="text"
            item-value="value"
            label="所属服务"
            class="mb-2"
          />
          <v-text-field v-model="form.host" label="主机地址" class="mb-2" />
          <v-text-field v-model.number="form.port" label="端口" type="number" class="mb-2" />
          <v-text-field v-model.number="form.priority" label="优先级 (值越小越优先)" type="number" class="mb-4" />
          <div class="d-flex justify-end ga-2">
            <v-btn variant="text" @click="dialog = false">取消</v-btn>
            <v-btn type="submit" color="primary" :loading="saving">保存</v-btn>
          </div>
        </v-form>
      </v-card>
    </v-dialog>

    <!-- 删除确认 -->
    <v-dialog v-model="deleteDialog" max-width="400">
      <v-card class="pa-6">
        <div class="text-h6 mb-2">确认删除</div>
        <p class="text-body-2 text-medium-emphasis mb-4">确定要删除节点 <strong>{{ deleteTarget?.host }}:{{ deleteTarget?.port }}</strong> 吗？</p>
        <div class="d-flex justify-end ga-2">
          <v-btn variant="text" @click="deleteDialog = false">取消</v-btn>
          <v-btn color="error" :loading="deleting" @click="handleDelete">删除</v-btn>
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
