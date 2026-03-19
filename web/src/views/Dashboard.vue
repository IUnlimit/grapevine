<template>
  <div>
    <!-- QPS 卡片 -->
    <v-row class="mb-4">
      <v-col cols="12" md="4">
        <v-card class="pa-6">
          <div class="text-medium-emphasis text-body-2 mb-1">当前 QPS</div>
          <div class="text-h3 font-weight-bold text-primary">{{ dashboard.qps }}</div>
        </v-card>
      </v-col>
      <v-col cols="12" md="4">
        <v-card class="pa-6">
          <div class="text-medium-emphasis text-body-2 mb-1">在线节点</div>
          <div class="text-h3 font-weight-bold text-success">{{ onlineCount }}</div>
        </v-card>
      </v-col>
      <v-col cols="12" md="4">
        <v-card class="pa-6">
          <div class="text-medium-emphasis text-body-2 mb-1">离线节点</div>
          <div class="text-h3 font-weight-bold text-error">{{ offlineCount }}</div>
        </v-card>
      </v-col>
    </v-row>

    <!-- 流量分布 -->
    <v-row class="mb-4">
      <v-col cols="12" md="6">
        <v-card class="pa-6">
          <div class="text-subtitle-1 font-weight-medium mb-4">流量分布</div>
          <Doughnut v-if="chartData.labels.length" :data="chartData" :options="chartOptions" />
          <div v-else class="text-center text-medium-emphasis py-8">暂无流量数据</div>
        </v-card>
      </v-col>
      <v-col cols="12" md="6">
        <v-card class="pa-6">
          <div class="text-subtitle-1 font-weight-medium mb-4">后缀请求统计</div>
          <v-table density="comfortable">
            <thead>
              <tr>
                <th>后缀</th>
                <th class="text-right">请求数</th>
                <th class="text-right">错误数</th>
              </tr>
            </thead>
            <tbody>
              <tr v-for="s in dashboard.suffixes" :key="s.suffix">
                <td><v-chip size="small" color="primary" variant="tonal">{{ s.suffix }}</v-chip></td>
                <td class="text-right">{{ s.requests }}</td>
                <td class="text-right text-error">{{ s.errors }}</td>
              </tr>
              <tr v-if="!dashboard.suffixes?.length">
                <td colspan="3" class="text-center text-medium-emphasis">暂无数据</td>
              </tr>
            </tbody>
          </v-table>
        </v-card>
      </v-col>
    </v-row>

    <!-- 节点状态 -->
    <v-card class="pa-6">
      <div class="text-subtitle-1 font-weight-medium mb-4">节点存活状态</div>
      <v-table density="comfortable">
        <thead>
          <tr>
            <th>服务</th>
            <th>后缀</th>
            <th>地址</th>
            <th>优先级</th>
            <th>状态</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="(ep, i) in dashboard.endpoints" :key="i">
            <td>{{ ep.service_name }}</td>
            <td><v-chip size="small" color="primary" variant="tonal">{{ ep.suffix }}</v-chip></td>
            <td>{{ ep.host }}:{{ ep.port }}</td>
            <td>{{ ep.priority }}</td>
            <td>
              <v-chip :color="ep.status === 'online' ? 'success' : 'error'" size="small" variant="flat">
                {{ ep.status === 'online' ? '在线' : '离线' }}
              </v-chip>
            </td>
          </tr>
          <tr v-if="!dashboard.endpoints?.length">
            <td colspan="5" class="text-center text-medium-emphasis">暂无节点</td>
          </tr>
        </tbody>
      </v-table>
    </v-card>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, onUnmounted } from 'vue'
import { Doughnut } from 'vue-chartjs'
import { Chart as ChartJS, ArcElement, Tooltip, Legend } from 'chart.js'
import api from '../api'

ChartJS.register(ArcElement, Tooltip, Legend)

interface SuffixSummary { suffix: string; requests: number; errors: number }
interface EndpointStatus { service_name: string; suffix: string; host: string; port: number; priority: number; status: string }
interface DashboardData { qps: number; suffixes: SuffixSummary[]; endpoints: EndpointStatus[] }

const dashboard = ref<DashboardData>({ qps: 0, suffixes: [], endpoints: [] })
let timer: ReturnType<typeof setInterval>

const onlineCount = computed(() => dashboard.value.endpoints?.filter((e) => e.status === 'online').length ?? 0)
const offlineCount = computed(() => dashboard.value.endpoints?.filter((e) => e.status === 'offline').length ?? 0)

const colors = ['#6A1B9A', '#AB47BC', '#CE93D8', '#E1BEE7', '#F3E5F5', '#4A148C', '#7B1FA2', '#9C27B0']

const chartData = computed(() => ({
  labels: dashboard.value.suffixes?.map((s) => s.suffix) ?? [],
  datasets: [{
    data: dashboard.value.suffixes?.map((s) => s.requests) ?? [],
    backgroundColor: colors.slice(0, dashboard.value.suffixes?.length ?? 0),
    borderWidth: 0,
  }],
}))

const chartOptions = { responsive: true, plugins: { legend: { position: 'bottom' as const } } }

async function fetchDashboard() {
  try {
    const res = await api.get('/admin/dashboard')
    dashboard.value = res.data
  } catch { /* ignore */ }
}

onMounted(() => {
  fetchDashboard()
  timer = setInterval(fetchDashboard, 3000)
})

onUnmounted(() => clearInterval(timer))
</script>
