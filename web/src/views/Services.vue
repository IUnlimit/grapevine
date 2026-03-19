<template>
  <div>
    <div class="d-flex align-center mb-4">
      <div class="text-h6 font-weight-medium">路由管理</div>
      <v-spacer />
      <v-btn color="primary" prepend-icon="mdi-plus" @click="openDialog()">新增路由</v-btn>
    </div>

    <v-card>
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
          <tr v-for="svc in services" :key="svc.id">
            <td>{{ svc.id }}</td>
            <td><v-chip size="small" color="primary" variant="tonal">/{{ svc.suffix }}</v-chip></td>
            <td>{{ svc.name }}</td>
            <td class="text-medium-emphasis">{{ svc.description || '-' }}</td>
            <td>
              <v-chip :color="svc.is_active ? 'success' : 'default'" size="small" variant="flat">
                {{ svc.is_active ? '启用' : '禁用' }}
              </v-chip>
            </td>
            <td class="text-right">
              <v-btn icon="mdi-pencil" size="small" variant="text" @click="openDialog(svc)" />
              <v-btn icon="mdi-delete" size="small" variant="text" color="error" @click="confirmDelete(svc)" />
            </td>
          </tr>
          <tr v-if="!services.length">
            <td colspan="6" class="text-center text-medium-emphasis py-6">暂无路由配置</td>
          </tr>
        </tbody>
      </v-table>
    </v-card>

    <!-- 新增/编辑对话框 -->
    <v-dialog v-model="dialog" max-width="500" persistent>
      <v-card class="pa-6">
        <div class="text-h6 mb-4">{{ editing ? '编辑路由' : '新增路由' }}</div>
        <v-form @submit.prevent="handleSave">
          <v-text-field v-model="form.suffix" label="后缀 (Suffix)" class="mb-2" />
          <v-text-field v-model="form.name" label="服务名称" class="mb-2" />
          <v-text-field v-model="form.description" label="描述" class="mb-2" />
          <v-switch v-model="form.is_active" label="启用" color="primary" class="mb-4" />
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
        <p class="text-body-2 text-medium-emphasis mb-4">确定要删除路由 <strong>{{ deleteTarget?.name }}</strong> 吗？此操作不可撤销。</p>
        <div class="d-flex justify-end ga-2">
          <v-btn variant="text" @click="deleteDialog = false">取消</v-btn>
          <v-btn color="error" :loading="deleting" @click="handleDelete">删除</v-btn>
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
