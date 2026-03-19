<template>
  <div class="dashboard">
    <!-- Page header -->
    <div class="page-header">
      <div>
        <h1 class="page-title">仪表盘</h1>
        <p class="page-desc">实时监控网关运行状态</p>
      </div>
      <div class="header-actions">
        <v-chip size="small" variant="tonal" color="success" prepend-icon="mdi-circle-small">
          实时刷新中
        </v-chip>
      </div>
    </div>

    <!-- Stat cards -->
    <div class="stat-grid">
      <div class="gv-stat-card stat-card" :style="{ background: cardBg }">
        <div class="stat-card-inner">
          <div class="stat-icon stat-icon--purple">
            <v-icon size="22" color="white">mdi-lightning-bolt</v-icon>
          </div>
          <div class="stat-content">
            <span class="stat-label">当前 QPS</span>
            <span class="stat-value">{{ dashboard.qps }}</span>
          </div>
        </div>
        <div class="stat-sparkline stat-sparkline--purple"></div>
      </div>

      <div class="gv-stat-card stat-card" :style="{ background: cardBg }">
        <div class="stat-card-inner">
          <div class="stat-icon stat-icon--green">
            <v-icon size="22" color="white">mdi-server</v-icon>
          </div>
          <div class="stat-content">
            <span class="stat-label">在线节点</span>
            <span class="stat-value stat-value--green">{{ onlineCount }}</span>
          </div>
        </div>
      </div>

      <div class="gv-stat-card stat-card" :style="{ background: cardBg }">
        <div class="stat-card-inner">
          <div class="stat-icon stat-icon--red">
            <v-icon size="22" color="white">mdi-server-off</v-icon>
          </div>
          <div class="stat-content">
            <span class="stat-label">离线节点</span>
            <span class="stat-value stat-value--red">{{ offlineCount }}</span>
          </div>
        </div>
      </div>

      <div class="gv-stat-card stat-card" :style="{ background: cardBg }">
        <div class="stat-card-inner">
          <div class="stat-icon stat-icon--amber">
            <v-icon size="22" color="white">mdi-chart-timeline-variant</v-icon>
          </div>
          <div class="stat-content">
            <span class="stat-label">路由总数</span>
            <span class="stat-value stat-value--amber">{{ dashboard.suffixes?.length ?? 0 }}</span>
          </div>
        </div>
      </div>
    </div>

    <!-- Charts row -->
    <div class="charts-grid">
      <!-- Traffic distribution -->
      <v-card class="chart-card gv-glass">
        <div class="card-header">
          <div>
            <h3 class="card-title">流量分布</h3>
            <p class="card-desc">各后缀请求占比</p>
          </div>
        </div>
        <div class="chart-container">
          <Doughnut v-if="chartData.labels.length" :data="chartData" :options="chartOptions" />
          <div v-else class="empty-state">
            <v-icon size="40" color="medium-emphasis" class="mb-2">mdi-chart-donut</v-icon>
            <span>暂无流量数据</span>
          </div>
        </div>
      </v-card>

      <!-- Suffix stats table -->
      <v-card class="chart-card gv-glass">
        <div class="card-header">
          <div>
            <h3 class="card-title">后缀请求统计</h3>
            <p class="card-desc">各路由请求与错误计数</p>
          </div>
        </div>
        <div class="table-container">
          <v-table density="comfortable" class="stats-table">
            <thead>
              <tr>
                <th>后缀</th>
                <th class="text-right">请求数</th>
                <th class="text-right">错误数</th>
                <th class="text-right">错误率</th>
              </tr>
            </thead>
            <tbody>
              <tr v-for="s in dashboard.suffixes" :key="s.suffix" class="table-row">
                <td>
                  <div class="suffix-badge">
                    <span class="suffix-slash">/</span>{{ s.suffix }}
                  </div>
                </td>
                <td class="text-right font-mono">{{ s.requests.toLocaleString() }}</td>
                <td class="text-right font-mono text-error">{{ s.errors }}</td>
                <td class="text-right">
                  <v-chip
                    size="x-small"
                    :color="getErrorRate(s) > 10 ? 'error' : getErrorRate(s) > 0 ? 'warning' : 'success'"
                    variant="tonal"
                  >
                    {{ getErrorRate(s).toFixed(1) }}%
                  </v-chip>
                </td>
              </tr>
              <tr v-if="!dashboard.suffixes?.length">
                <td colspan="4" class="text-center py-8" style="opacity: 0.4">暂无数据</td>
              </tr>
            </tbody>
          </v-table>
        </div>
      </v-card>
    </div>

    <!-- Endpoints status -->
    <v-card class="endpoints-card gv-glass">
      <div class="card-header">
        <div>
          <h3 class="card-title">节点存活状态</h3>
          <p class="card-desc">所有注册节点的实时健康状态</p>
        </div>
        <div class="d-flex align-center ga-3">
          <div class="legend-item">
            <span class="gv-dot gv-dot--online"></span>
            <span class="legend-text">在线 {{ onlineCount }}</span>
          </div>
          <div class="legend-item">
            <span class="gv-dot gv-dot--offline"></span>
            <span class="legend-text">离线 {{ offlineCount }}</span>
          </div>
        </div>
      </div>

      <v-table density="comfortable" class="endpoints-table">
        <thead>
          <tr>
            <th>状态</th>
            <th>服务</th>
            <th>后缀</th>
            <th>地址</th>
            <th>优先级</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="(ep, i) in dashboard.endpoints" :key="i" class="table-row">
            <td>
              <div class="status-cell">
                <span class="gv-dot" :class="ep.status === 'online' ? 'gv-dot--online' : 'gv-dot--offline'"></span>
                <span :class="ep.status === 'online' ? 'text-success' : 'text-error'" style="font-size: 12px; font-weight: 500">
                  {{ ep.status === 'online' ? '在线' : '离线' }}
                </span>
              </div>
            </td>
            <td class="font-weight-medium">{{ ep.service_name }}</td>
            <td>
              <div class="suffix-badge">
                <span class="suffix-slash">/</span>{{ ep.suffix }}
              </div>
            </td>
            <td>
              <code class="address-code">{{ ep.host }}:{{ ep.port }}</code>
            </td>
            <td>
              <v-chip size="x-small" :color="ep.priority <= 1 ? 'primary' : 'default'" variant="tonal">
                P{{ ep.priority }}
              </v-chip>
            </td>
          </tr>
          <tr v-if="!dashboard.endpoints?.length">
            <td colspan="5" class="text-center py-8" style="opacity: 0.4">暂无节点</td>
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
import { useAppStore } from '../stores/app'
import api from '../api'

ChartJS.register(ArcElement, Tooltip, Legend)

const appStore = useAppStore()

interface SuffixSummary { suffix: string; requests: number; errors: number }
interface EndpointStatus { service_name: string; suffix: string; host: string; port: number; priority: number; status: string }
interface DashboardData { qps: number; suffixes: SuffixSummary[]; endpoints: EndpointStatus[] }

const dashboard = ref<DashboardData>({ qps: 0, suffixes: [], endpoints: [] })
let timer: ReturnType<typeof setInterval>

const onlineCount = computed(() => dashboard.value.endpoints?.filter((e) => e.status === 'online').length ?? 0)
const offlineCount = computed(() => dashboard.value.endpoints?.filter((e) => e.status === 'offline').length ?? 0)

const cardBg = computed(() => appStore.darkMode ? '#1A1528' : '#FFFFFF')

function getErrorRate(s: SuffixSummary): number {
  if (!s.requests) return 0
  return (s.errors / s.requests) * 100
}

const chartColors = ['#7C3AED', '#A78BFA', '#34D399', '#60A5FA', '#FBBF24', '#F87171', '#C4B5FD', '#6EE7B7']

const chartData = computed(() => ({
  labels: dashboard.value.suffixes?.map((s) => '/' + s.suffix) ?? [],
  datasets: [{
    data: dashboard.value.suffixes?.map((s) => s.requests) ?? [],
    backgroundColor: chartColors.slice(0, dashboard.value.suffixes?.length ?? 0),
    borderWidth: 0,
    hoverOffset: 6,
  }],
}))

const chartOptions = {
  responsive: true,
  maintainAspectRatio: false,
  cutout: '65%',
  plugins: {
    legend: {
      position: 'bottom' as const,
      labels: {
        padding: 16,
        usePointStyle: true,
        pointStyleWidth: 8,
        font: { family: 'DM Sans', size: 12 },
        color: appStore.darkMode ? '#E8E0F0' : '#1E1B2E',
      },
    },
  },
}

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

<style scoped>
.dashboard {
  display: flex;
  flex-direction: column;
  gap: 24px;
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

/* Stat grid */
.stat-grid {
  display: grid;
  grid-template-columns: repeat(4, 1fr);
  gap: 16px;
}

@media (max-width: 1200px) {
  .stat-grid {
    grid-template-columns: repeat(2, 1fr);
  }
}

@media (max-width: 600px) {
  .stat-grid {
    grid-template-columns: 1fr;
  }
}

.stat-card {
  border-radius: 16px;
  padding: 24px;
  position: relative;
  overflow: hidden;
}

.stat-card-inner {
  display: flex;
  align-items: center;
  gap: 16px;
  position: relative;
  z-index: 1;
}

.stat-icon {
  width: 44px;
  height: 44px;
  border-radius: 12px;
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
}
.stat-icon--purple { background: linear-gradient(135deg, #7C3AED, #A78BFA); }
.stat-icon--green { background: linear-gradient(135deg, #10B981, #34D399); }
.stat-icon--red { background: linear-gradient(135deg, #EF4444, #F87171); }
.stat-icon--amber { background: linear-gradient(135deg, #F59E0B, #FBBF24); }

.stat-content {
  display: flex;
  flex-direction: column;
}

.stat-label {
  font-size: 12px;
  font-weight: 500;
  opacity: 0.5;
  text-transform: uppercase;
  letter-spacing: 0.04em;
}

.stat-value {
  font-size: 28px;
  font-weight: 700;
  letter-spacing: -1px;
  line-height: 1.1;
  margin-top: 2px;
  background: linear-gradient(135deg, #A78BFA, #C4B5FD);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  background-clip: text;
}
.stat-value--green {
  background: linear-gradient(135deg, #10B981, #34D399);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  background-clip: text;
}
.stat-value--red {
  background: linear-gradient(135deg, #EF4444, #F87171);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  background-clip: text;
}
.stat-value--amber {
  background: linear-gradient(135deg, #F59E0B, #FBBF24);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  background-clip: text;
}

/* Charts grid */
.charts-grid {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 16px;
}

@media (max-width: 960px) {
  .charts-grid {
    grid-template-columns: 1fr;
  }
}

.chart-card {
  padding: 24px;
}

.card-header {
  display: flex;
  align-items: flex-start;
  justify-content: space-between;
  margin-bottom: 20px;
}

.card-title {
  font-size: 15px;
  font-weight: 600;
  letter-spacing: -0.2px;
  margin: 0;
}

.card-desc {
  font-size: 12px;
  opacity: 0.4;
  margin-top: 2px;
}

.chart-container {
  height: 260px;
  display: flex;
  align-items: center;
  justify-content: center;
}

.empty-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 4px;
  opacity: 0.4;
  font-size: 13px;
}

.table-container {
  max-height: 300px;
  overflow-y: auto;
}

/* Suffix badge */
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

.font-mono {
  font-family: 'JetBrains Mono', monospace !important;
  font-size: 13px;
}

/* Endpoints card */
.endpoints-card {
  padding: 24px;
}

.legend-item {
  display: flex;
  align-items: center;
  gap: 6px;
}

.legend-text {
  font-size: 12px;
  font-weight: 500;
  opacity: 0.6;
}

.status-cell {
  display: flex;
  align-items: center;
  gap: 8px;
}

.address-code {
  font-family: 'JetBrains Mono', monospace !important;
  font-size: 12px;
  padding: 2px 8px;
  border-radius: 4px;
  background: rgba(124, 58, 237, 0.06);
}

.table-row {
  transition: background 0.15s ease;
}
.table-row:hover {
  background: rgba(124, 58, 237, 0.03);
}
</style>
