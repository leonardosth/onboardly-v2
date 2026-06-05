<template>
  <div class="dashboard-layout">

    <div class="main-content">
      <!-- Quick Navigation -->
      <div class="nav-shortcuts">
        <router-link to="/clients" class="shortcut-card">
          <span class="shortcut-icon">🏢</span>
          <div class="shortcut-info">
            <h3>Gerenciar Clientes</h3>
            <p>Cadastre, visualize e edite a base de clientes</p>
          </div>
        </router-link>
      </div>

      <!-- KPI Metrics Row -->
      <div v-if="loading" class="loading-container">
        Carregando métricas consolidadas...
      </div>
      <div v-else class="dashboard-grid">
        <!-- Activation Rate Widget -->
        <div class="metric-card glow-blue">
          <div class="card-title">Índice de Ativação</div>
          <div class="metric-value">{{ dashboardData.metrics?.activation_rate }}%</div>
          <div class="metric-progress-bar">
            <div class="progress" :style="'width: ' + dashboardData.metrics?.activation_rate + '%'" :class="{ success: dashboardData.metrics?.activation_rate >= 86 }"></div>
          </div>
          <div class="metric-footer">
            Meta recomendada: <strong>~86%</strong>
            <span v-if="dashboardData.metrics?.activation_rate >= 86" class="meta-status success">Meta Atingida!</span>
            <span v-else class="meta-status warning">Abaixo da meta</span>
          </div>
        </div>

        <!-- No-Show Rate Widget -->
        <div class="metric-card glow-red">
          <div class="card-title">Taxa de No-Show (Reuniões)</div>
          <div class="metric-value">{{ dashboardData.metrics?.no_show_rate }}%</div>
          <div class="metric-progress-bar">
            <div class="progress progress-red" :style="'width: ' + dashboardData.metrics?.no_show_rate + '%'" :class="{ success: dashboardData.metrics?.no_show_rate < 10 }"></div>
          </div>
          <div class="metric-footer">
            Meta recomendada: <strong>&lt; 10%</strong>
            <span v-if="dashboardData.metrics?.no_show_rate < 10" class="meta-status success">Meta Atingida!</span>
            <span v-else class="meta-status alert">Meta Violada</span>
          </div>
        </div>
      </div>

      <!-- Charts & Feed Grid -->
      <div v-if="!loading" class="analytics-grid">
        <!-- Monthly History SVG Bar Chart -->
        <div class="analytics-card">
          <h3>Evolução de Projetos Concluídos (6 Meses)</h3>
          <div class="chart-container">
            <svg viewBox="0 0 500 220" class="bar-chart">
              <!-- Grid lines -->
              <line x1="40" y1="30" x2="480" y2="30" stroke="#334155" stroke-dasharray="4" />
              <line x1="40" y1="90" x2="480" y2="90" stroke="#334155" stroke-dasharray="4" />
              <line x1="40" y1="150" x2="480" y2="150" stroke="#334155" stroke-dasharray="4" />
              <line x1="40" y1="180" x2="480" y2="180" stroke="#475569" />

              <!-- Y-axis labels -->
              <text x="15" y="35" fill="#94a3b8" font-size="10">Max</text>
              <text x="15" y="105" fill="#94a3b8" font-size="10">Mid</text>
              <text x="15" y="185" fill="#94a3b8" font-size="10">0</text>

              <!-- Bars -->
              <g v-for="(item, index) in dashboardData.history" :key="item.month">
                <!-- Calc dynamic height based on deployments -->
                <rect
                  :x="60 + index * 70"
                  :y="180 - Math.min(item.deployments * 15, 140)"
                  width="35"
                  :height="Math.min(item.deployments * 15, 140)"
                  fill="url(#blue-grad)"
                  rx="4"
                />
                <!-- Value tag -->
                <text
                  :x="77 + index * 70"
                  :y="175 - Math.min(item.deployments * 15, 140)"
                  fill="#f8fafc"
                  font-size="10"
                  text-anchor="middle"
                >
                  {{ item.deployments }}
                </text>
                <!-- Month name label -->
                <text
                  :x="77 + index * 70"
                  y="200"
                  fill="#94a3b8"
                  font-size="9"
                  text-anchor="middle"
                >
                  {{ formatMonth(item.month) }}
                </text>
              </g>

              <!-- Gradient Definition -->
              <defs>
                <linearGradient id="blue-grad" x1="0%" y1="0%" x2="0%" y2="100%">
                  <stop offset="0%" stop-color="#38bdf8" />
                  <stop offset="100%" stop-color="#0284c7" />
                </linearGradient>
              </defs>
            </svg>
          </div>
        </div>

        <!-- Recent Activities Feed -->
        <div class="analytics-card">
          <h3>Feed de Atividades Recentes</h3>
          <div v-if="dashboardData.recent_activities?.length === 0" class="empty-feed">
            Nenhuma atividade registrada recentemente.
          </div>
          <div v-else class="activity-feed">
            <div v-for="act in dashboardData.recent_activities" :key="act.id" class="feed-item">
              <span class="feed-icon">
                {{ act.entity_type === 'Client' ? '🏢' : act.entity_type === 'Project' ? '🚀' : '📅' }}
              </span>
              <div class="feed-details">
                <p class="feed-desc">{{ act.description }}</p>
                <span class="feed-time">{{ formatTime(act.created_at) }}</span>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue';
import { useRouter } from 'vue-router';
import axios from 'axios';
import { useAuthStore } from '../stores/auth';

const auth = useAuthStore();
const router = useRouter();

const loading = ref(true);
const dashboardData = ref({});

onMounted(() => {
  fetchDashboard();
});

const fetchDashboard = async () => {
  loading.value = true;
  try {
    const response = await axios.get('http://localhost:8080/api/dashboard');
    dashboardData.value = response.data;
  } catch (err) {
    console.error('Error fetching dashboard metrics:', err);
  } finally {
    loading.value = false;
  }
};

const handleLogout = () => {
  auth.logout();
  router.push('/login');
};

const formatMonth = (monthStr) => {
  // Converts YYYY-MM to Month/YY (e.g. 2026-06 to Jun/26)
  if (!monthStr) return '';
  const [year, month] = monthStr.split('-');
  const months = ['Jan', 'Fev', 'Mar', 'Abr', 'Mai', 'Jun', 'Jul', 'Ago', 'Set', 'Out', 'Nov', 'Dez'];
  const monthIdx = parseInt(month, 10) - 1;
  return `${months[monthIdx]}/${year.substring(2)}`;
};

const formatTime = (timeStr) => {
  if (!timeStr) return '';
  return new Date(timeStr).toLocaleString();
};
</script>

<style scoped>
.dashboard-layout {
  min-height: 100vh;
  background: #0f172a;
  color: #f8fafc;
  font-family: 'Inter', sans-serif;
}

.header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 1.5rem 2.5rem;
  background: rgba(30, 41, 59, 0.8);
  border-bottom: 1px solid rgba(255, 255, 255, 0.1);
  backdrop-filter: blur(10px);
}

.brand {
  display: flex;
  align-items: center;
  gap: 0.75rem;
}

.logo-icon {
  font-size: 1.5rem;
}

.header h1 {
  font-size: 1.5rem;
  font-weight: 700;
  color: #38bdf8;
  margin: 0;
}

.user-meta {
  display: flex;
  align-items: center;
  gap: 1.5rem;
}

.role-badge {
  font-size: 0.875rem;
  padding: 0.35rem 0.75rem;
  border-radius: 9999px;
  background: #334155;
}

.role-badge.admin {
  background: rgba(56, 189, 248, 0.2);
  color: #38bdf8;
}

.role-badge.analista {
  background: rgba(148, 163, 184, 0.2);
  color: #cbd5e1;
}

.logout-btn {
  background: transparent;
  border: 1px solid #ef4444;
  color: #ef4444;
  padding: 0.5rem 1rem;
  border-radius: 6px;
  cursor: pointer;
  font-weight: 500;
}

.logout-btn:hover {
  background: rgba(239, 68, 68, 0.1);
}

.main-content {
  padding: 2.5rem;
  max-width: 1200px;
  margin: 0 auto;
}

.nav-shortcuts {
  margin-bottom: 2.5rem;
}

.shortcut-card {
  display: flex;
  align-items: center;
  gap: 1.5rem;
  background: rgba(30, 41, 59, 0.4);
  border: 1px solid rgba(255, 255, 255, 0.05);
  border-radius: 16px;
  padding: 1.5rem 2rem;
  color: #f8fafc;
  text-decoration: none;
  transition: transform 0.2s, border-color 0.2s;
  max-width: 350px;
}

.shortcut-card:hover {
  transform: translateY(-2px);
  border-color: #38bdf8;
}

.shortcut-icon {
  font-size: 2.5rem;
}

.shortcut-info h3 {
  margin: 0 0 0.25rem 0;
  color: #38bdf8;
}

.shortcut-info p {
  margin: 0;
  font-size: 0.875rem;
  color: #94a3b8;
}

.dashboard-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(280px, 1fr));
  gap: 2rem;
  margin-bottom: 2.5rem;
}

.metric-card {
  background: rgba(30, 41, 59, 0.5);
  border: 1px solid rgba(255, 255, 255, 0.05);
  border-radius: 16px;
  padding: 2rem;
  box-shadow: 0 10px 15px -3px rgba(0, 0, 0, 0.3);
}

.glow-blue {
  border-left: 4px solid #38bdf8;
}

.glow-red {
  border-left: 4px solid #ef4444;
}

.card-title {
  font-size: 0.9rem;
  color: #94a3b8;
  text-transform: uppercase;
  letter-spacing: 0.05em;
  margin-bottom: 1rem;
}

.metric-value {
  font-size: 3rem;
  font-weight: 800;
  color: #f8fafc;
  margin-bottom: 1.5rem;
}

.metric-progress-bar {
  height: 8px;
  background: #1e293b;
  border-radius: 9999px;
  overflow: hidden;
  margin-bottom: 1.5rem;
}

.progress {
  height: 100%;
  border-radius: 9999px;
  background: #f59e0b; /* default warning */
}

.progress.success {
  background: #22c55e;
}

.progress-red {
  background: #22c55e; /* default success */
}

.progress-red.success {
  background: #ef4444; /* violation */
}

.metric-footer {
  display: flex;
  justify-content: space-between;
  align-items: center;
  font-size: 0.85rem;
  color: #94a3b8;
}

.meta-status {
  font-weight: 600;
  padding: 0.2rem 0.5rem;
  border-radius: 4px;
}

.meta-status.success {
  background: rgba(34, 197, 94, 0.15);
  color: #4ade80;
}

.meta-status.warning {
  background: rgba(245, 158, 11, 0.15);
  color: #fbbf24;
}

.meta-status.alert {
  background: rgba(239, 68, 68, 0.15);
  color: #fca5a5;
}

.analytics-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(450px, 1fr));
  gap: 2.5rem;
}

.analytics-card {
  background: rgba(30, 41, 59, 0.4);
  border: 1px solid rgba(255, 255, 255, 0.05);
  border-radius: 16px;
  padding: 2rem;
}

.analytics-card h3 {
  margin-top: 0;
  margin-bottom: 1.5rem;
  color: #38bdf8;
  font-size: 1.2rem;
}

.chart-container {
  display: flex;
  justify-content: center;
  align-items: center;
}

.bar-chart {
  width: 100%;
  max-width: 450px;
  height: 200px;
}

.activity-feed {
  max-height: 250px;
  overflow-y: auto;
  display: flex;
  flex-direction: column;
  gap: 1rem;
}

.feed-item {
  display: flex;
  align-items: flex-start;
  gap: 1rem;
  background: rgba(15, 23, 42, 0.3);
  padding: 0.75rem 1rem;
  border-radius: 8px;
  border: 1px solid rgba(255, 255, 255, 0.02);
}

.feed-icon {
  font-size: 1.25rem;
}

.feed-details {
  flex: 1;
}

.feed-desc {
  margin: 0 0 0.25rem 0;
  font-size: 0.9rem;
}

.feed-time {
  font-size: 0.75rem;
  color: #64748b;
}

.loading-container, .empty-feed {
  text-align: center;
  padding: 2rem;
  color: #94a3b8;
}
</style>
