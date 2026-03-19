<template>
  <div class="services-page">
    <!-- Page header -->
    <div class="page-header">
      <div>
        <h1 class="page-title">路由管理</h1>
        <p class="page-desc">管理网关的后缀路由规则</p>
      </div>
      <v-btn
        color="primary"
        prepend-icon="mdi-plus"
        @click="openDialog()"
        elevation="0"
        class="add-btn"
      >
        新增路由
      </v-btn>
    </div>

    <!-- Table card -->
    <v-card class="table-card gv-glass">
      <v-table density="comfortable">
        <thead>
          <tr>
            <th>ID</th>
            <th>后缀</th>
            <th>名称</th>
            <th>描述</th>
            <th>状态</th>
            <th class="text-right">操作</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="svc in services" :key="svc.id" class="table-row">
            <td>
              <span class="id-badge">#{{ svc.id }}</span>
            </td>
            <td>
              <div class="suffix-badge">
                <span class="suffix-slash">/</span>{{ svc.suffix }}
              </div>
            </td>
            <td class="font-weight-medium">{{ svc.name }}</td>
            <td style="opacity: 0.5; font-size: 13px">{{ svc.description || '—' }}</td>
            <td>
              <div class="status-pill" :class="svc.is_active ? 'status-pill--active' : 'status-pill--inactive'">
                <span class="status-dot"></span>
                {{ svc.is_active ? '启用' : '禁用' }}
              </div>
            </td>
            <td class="text-right">
              <v-btn
                icon="mdi-pencil-outline"
                size="small"
                variant="text"
                @click="openDialog(svc)"
                class="action-btn"
              />
              <v-btn
                icon="mdi-delete-outline"
                size="small"
                variant="text"
                color="error"
                @click="confirmDelete(svc)"
                class="action-btn"
              />
            </td>
          </tr>
          <tr v-if="!services.length">
            <td colspan="6" class="text-center py-12">
              <div class="empty-state">
                <v-icon size="48" color="medium-emphasis" class="mb-3">mdi-routes</v-icon>
                <p style="opacity: 0.4; font-size: 14px">暂无路由配置</p>
                <v-btn
                  variant="tonal"
                  color="primary"
                  size="small"
                  prepend-icon="mdi-plus"
                  @click="openDialog()"
                  class="mt-3"
                >
                  创建第一个路由
                </v-btn>
              </div>
            </td>
          </tr>
        </tbody>
      </v-table>
    </v-card>

    <!-- Create/Edit dialog -->
    <v-dialog v-model="dialog" max-width="480" persistent>
      <v-card class="dialog-card">
        <div class="dialog-header">
          <h3 class="dialog-title">{{ editing ? '编辑路由' : '新增路由' }}</h3>
          <v-btn icon="mdi-close" size="small" variant="text" @click="dialog = false" />
        </div>
        <v-divider style="opacity: 0.06" />
        <div class="dialog-body">
          <v-form @submit.prevent="handleSave">
            <div class="field-group">
              <label class="field-label">后缀 (Suffix)</label>
              <v-text-field
                v-model="form.suffix"
                placeholder="例如: api, ws, static"
                hide-details
                bg-color="transparent"
              />
            </div>
            <div class="field-group">
              <label class="field-label">服务名称</label>
              <v-text-field
                v-model="form.name"
                placeholder="输入服务名称"
                hide-details
                bg-color="transparent"
              />
            </div>
            <div class="field-group">
              <label class="field-label">描述</label>
              <v-text-field
                v-model="form.description"
                placeholder="可选描述信息"
                hide-details
                bg-color="transparent"
              />
            </div>
            <v-switch
              v-model="form.is_active"
              label="启用"
              color="success"
              hide-details
              class="mt-2"
            />
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
            确定要删除路由「{{ deleteTarget?.name }}」吗？此操作不可撤销。
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
import { ref, reactive, onMounted } from 'vue'
import api from '../api'

interface Service {
  id?: number
  suffix: string
  name: string
  description: string
  is_active: boolean
}

const services = ref<Service[]>([])
const dialog = ref(false)
const deleteDialog = ref(false)
const editing = ref(false)
const saving = ref(false)
const deleting = ref(false)
const editId = ref<number | null>(null)
const deleteTarget = ref<Service | null>(null)

const form = reactive<Service>({ suffix: '', name: '', description: '', is_active: true })

async function fetchServices() {
  const res = await api.get('/admin/services')
  services.value = res.data
}

function openDialog(svc?: Service) {
  if (svc) {
    editing.value = true
    editId.value = svc.id!
    Object.assign(form, { suffix: svc.suffix, name: svc.name, description: svc.description, is_active: svc.is_active })
  } else {
    editing.value = false
    editId.value = null
    Object.assign(form, { suffix: '', name: '', description: '', is_active: true })
  }
  dialog.value = true
}

async function handleSave() {
  saving.value = true
  try {
    if (editing.value && editId.value) {
      await api.put(`/admin/services/${editId.value}`, form)
    } else {
      await api.post('/admin/services', form)
    }
    dialog.value = false
    await fetchServices()
  } finally {
    saving.value = false
  }
}

function confirmDelete(svc: Service) {
  deleteTarget.value = svc
  deleteDialog.value = true
}

async function handleDelete() {
  if (!deleteTarget.value?.id) return
  deleting.value = true
  try {
    await api.delete(`/admin/services/${deleteTarget.value.id}`)
    deleteDialog.value = false
    await fetchServices()
  } finally {
    deleting.value = false
  }
}

onMounted(fetchServices)
</script>

<style scoped>
.services-page {
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

.add-btn {
  box-shadow: 0 2px 12px rgba(124, 58, 237, 0.25) !important;
}

.table-card {
  overflow: hidden;
}

.table-row {
  transition: background 0.15s ease;
}
.table-row:hover {
  background: rgba(124, 58, 237, 0.03);
}

.id-badge {
  font-family: 'JetBrains Mono', monospace;
  font-size: 12px;
  opacity: 0.4;
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

.status-pill {
  display: inline-flex;
  align-items: center;
  gap: 6px;
  padding: 3px 10px;
  border-radius: 20px;
  font-size: 12px;
  font-weight: 500;
}
.status-pill--active {
  background: rgba(16, 185, 129, 0.1);
  color: #34D399;
}
.status-pill--inactive {
  background: rgba(255, 255, 255, 0.05);
  opacity: 0.5;
}
.status-dot {
  width: 6px;
  height: 6px;
  border-radius: 50%;
  background: currentColor;
}
.status-pill--active .status-dot {
  box-shadow: 0 0 6px rgba(52, 211, 153, 0.5);
}

.action-btn {
  opacity: 0.5;
  transition: opacity 0.15s ease;
}
.action-btn:hover {
  opacity: 1;
}

.empty-state {
  display: flex;
  flex-direction: column;
  align-items: center;
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
