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
        <router-link to="/meetings" class="shortcut-card">
          <span class="shortcut-icon">📅</span>
          <div class="shortcut-info">
            <h3>Minhas Agendas</h3>
            <p>Acompanhe e conclua suas reuniões de implantação</p>
          </div>
        </router-link>
      </div>

      <!-- Loading / Error States -->
      <div v-if="loading" class="loading-container">
        Carregando métricas consolidadas...
      </div>
      
      <div v-else class="dashboard-content">
        <!-- KPI Metrics Row -->
        <div class="dashboard-grid">
          <!-- Activation Rate Widget -->
          <div class="metric-card glow-blue">
            <div class="card-title">Índice Geral de Ativação</div>
            <div class="metric-value">{{ dashboardData.metrics?.activation_rate }}%</div>
            <div class="metric-progress-bar">
              <div 
                class="progress" 
                :style="'width: ' + dashboardData.metrics?.activation_rate + '%'" 
                :class="{ success: dashboardData.metrics?.activation_rate >= 86 }"
              ></div>
            </div>
            <div class="metric-footer">
              Meta recomendada: <strong>~86%</strong>
              <span v-if="dashboardData.metrics?.activation_rate >= 86" class="meta-status success">Meta Atingida!</span>
              <span v-else class="meta-status warning">Abaixo da meta</span>
            </div>
          </div>

          <!-- No-Show Rate Widget -->
          <div class="metric-card glow-red">
            <div class="card-title">Taxa de No-Show</div>
            <div class="metric-value">{{ dashboardData.metrics?.no_show_rate }}%</div>
            <div class="metric-progress-bar">
              <div 
                class="progress progress-red" 
                :style="'width: ' + dashboardData.metrics?.no_show_rate + '%'" 
                :class="{ success: dashboardData.metrics?.no_show_rate < 10 }"
              ></div>
            </div>
            <div class="metric-footer">
              Meta recomendada: <strong>&lt; 10%</strong>
              <span v-if="dashboardData.metrics?.no_show_rate < 10" class="meta-status success">Meta Atingida!</span>
              <span v-else class="meta-status alert">Meta Violada</span>
            </div>
          </div>

          <!-- Abandonment Rate Widget -->
          <div class="metric-card glow-orange">
            <div class="card-title">Taxa de Abandono</div>
            <div class="metric-value">{{ dashboardData.metrics?.abandonment_rate }}%</div>
            <div class="metric-progress-bar">
              <div 
                class="progress progress-red" 
                :style="'width: ' + dashboardData.metrics?.abandonment_rate + '%'" 
                :class="{ success: dashboardData.metrics?.abandonment_rate <= 20 }"
              ></div>
            </div>
            <div class="metric-footer">
              Meta recomendada: <strong>&le; 20%</strong>
              <span v-if="dashboardData.metrics?.abandonment_rate <= 20" class="meta-status success">Meta Atingida!</span>
              <span v-else class="meta-status alert">Meta Violada</span>
            </div>
          </div>

          <!-- 30-Day Activation Rate Widget -->
          <div class="metric-card glow-green">
            <div class="card-title">Ativação em até 30 Dias</div>
            <div class="metric-value">{{ dashboardData.metrics?.activation_30d_rate }}%</div>
            <div class="metric-progress-bar">
              <div 
                class="progress" 
                :style="'width: ' + dashboardData.metrics?.activation_30d_rate + '%'" 
                :class="{ success: dashboardData.metrics?.activation_30d_rate >= 80 }"
              ></div>
            </div>
            <div class="metric-footer">
              Meta recomendada: <strong>&ge; 80%</strong>
              <span v-if="dashboardData.metrics?.activation_30d_rate >= 80" class="meta-status success">Meta Atingida!</span>
              <span v-else class="meta-status warning">Abaixo da meta</span>
            </div>
          </div>

          <!-- First-Meeting Activation Rate Widget -->
          <div class="metric-card glow-purple">
            <div class="card-title">Ativação na 1ª Reunião</div>
            <div class="metric-value">{{ dashboardData.metrics?.first_meeting_activation_rate }}%</div>
            <div class="metric-progress-bar">
              <div 
                class="progress progress-purple" 
                :style="'width: ' + dashboardData.metrics?.first_meeting_activation_rate + '%'"
              ></div>
            </div>
            <div class="metric-footer">
              Acompanhamento operacional
            </div>
          </div>
        </div>

        <!-- Funnel Section -->
        <div class="funnel-card">
          <h3>Funil de Implantação</h3>
          <div class="funnel-container">
            <div class="funnel-stage stage-registered">
              <div class="stage-name">Inscritos</div>
              <div class="stage-value">{{ dashboardData.funnel?.registered }}</div>
              <div class="stage-pct">100% dos Projetos</div>
            </div>
            
            <div class="funnel-connector">
              <div class="connector-arrow">➔</div>
              <div class="connector-value">{{ getPct(dashboardData.funnel?.participants, dashboardData.funnel?.registered) }}% de conv.</div>
            </div>
            
            <div class="funnel-stage stage-participants">
              <div class="stage-name">Participantes</div>
              <div class="stage-value">{{ dashboardData.funnel?.participants }}</div>
              <div class="stage-pct">{{ getPct(dashboardData.funnel?.participants, dashboardData.funnel?.registered) }}% do total</div>
            </div>
            
            <div class="funnel-connector">
              <div class="connector-arrow">➔</div>
              <div class="connector-value">{{ getPct(dashboardData.funnel?.active, dashboardData.funnel?.participants) }}% de conv.</div>
            </div>
            
            <div class="funnel-stage stage-active">
              <div class="stage-name">Ativos</div>
              <div class="stage-value">{{ dashboardData.funnel?.active }}</div>
              <div class="stage-pct">{{ getPct(dashboardData.funnel?.active, dashboardData.funnel?.registered) }}% total</div>
            </div>
          </div>
        </div>

        <!-- Charts & Feed Grid -->
        <div class="analytics-grid">
          <!-- Cohort (Safra) Widget -->
          <div class="analytics-card">
            <h3>Desempenho por Safra (Mês de Compra)</h3>
            <div class="table-container">
              <table class="cohort-table">
                <thead>
                  <tr>
                    <th>Safra</th>
                    <th>Total Clientes</th>
                    <th>Clientes Ativados</th>
                    <th>Taxa de Ativação</th>
                  </tr>
                </thead>
                <tbody>
                  <tr 
                    v-for="cohort in dashboardData.cohorts" 
                    :key="cohort.month"
                    :class="{ 'target-met': cohort.activation_rate >= 80 }"
                  >
                    <td class="cohort-month">{{ formatMonth(cohort.month) }}</td>
                    <td>{{ cohort.total }}</td>
                    <td>{{ cohort.activated }}</td>
                    <td>
                      <div class="cohort-rate-cell">
                        <span class="rate-value">{{ cohort.activation_rate }}%</span>
                        <span v-if="cohort.activation_rate >= 80" class="badge-target-met">Meta 80% ✓</span>
                      </div>
                    </td>
                  </tr>
                  <tr v-if="!dashboardData.cohorts || dashboardData.cohorts.length === 0">
                    <td colspan="4" class="empty-table-text">Nenhuma safra registrada.</td>
                  </tr>
                </tbody>
              </table>
            </div>
          </div>

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
        </div>

        <div class="analytics-grid mt-4">
          <!-- Recent Activities Feed -->
          <div class="analytics-card full-width-card">
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

const getPct = (val, total) => {
  if (!total || !val) return 0;
  return Math.round((val / total) * 100);
};

const formatMonth = (monthStr) => {
  // Converts YYYY-MM to Month/YY (e.g. 2026-06 to Jun/26)
  if (!monthStr) return '';
  const parts = monthStr.split('-');
  if (parts.length < 2) return monthStr;
  const [year, month] = parts;
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

.main-content {
  padding: 2.5rem;
  max-width: 1200px;
  margin: 0 auto;
}

.nav-shortcuts {
  display: flex;
  gap: 1.5rem;
  margin-bottom: 2.5rem;
  flex-wrap: wrap;
}

.shortcut-card {
  display: flex;
  align-items: center;
  gap: 1.5rem;
  background: rgba(30, 41, 59, 0.4);
  border: 1px solid rgba(255, 255, 255, 0.05);
  border-radius: 16px;
  padding: 1.25rem 1.75rem;
  color: #f8fafc;
  text-decoration: none;
  transition: transform 0.2s, border-color 0.2s;
  flex: 1;
  min-width: 280px;
  max-width: 380px;
  box-shadow: 0 4px 6px -1px rgba(0, 0, 0, 0.1);
}

.shortcut-card:hover {
  transform: translateY(-2px);
  border-color: #38bdf8;
}

.shortcut-icon {
  font-size: 2.2rem;
}

.shortcut-info h3 {
  margin: 0 0 0.25rem 0;
  color: #38bdf8;
  font-size: 1.1rem;
}

.shortcut-info p {
  margin: 0;
  font-size: 0.85rem;
  color: #94a3b8;
}

.dashboard-content {
  display: flex;
  flex-direction: column;
  gap: 2.5rem;
}

.dashboard-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(220px, 1fr));
  gap: 1.5rem;
}

.metric-card {
  background: rgba(30, 41, 59, 0.5);
  border: 1px solid rgba(255, 255, 255, 0.05);
  border-radius: 16px;
  padding: 1.5rem;
  box-shadow: 0 10px 15px -3px rgba(0, 0, 0, 0.3);
  display: flex;
  flex-direction: column;
  justify-content: space-between;
}

.glow-blue { border-left: 4px solid #38bdf8; }
.glow-red { border-left: 4px solid #ef4444; }
.glow-orange { border-left: 4px solid #f97316; }
.glow-green { border-left: 4px solid #10b981; }
.glow-purple { border-left: 4px solid #a855f7; }

.card-title {
  font-size: 0.8rem;
  color: #94a3b8;
  text-transform: uppercase;
  letter-spacing: 0.05em;
  margin-bottom: 0.75rem;
  font-weight: 600;
}

.metric-value {
  font-size: 2.5rem;
  font-weight: 800;
  color: #f8fafc;
  margin-bottom: 1rem;
}

.metric-progress-bar {
  height: 6px;
  background: #1e293b;
  border-radius: 9999px;
  overflow: hidden;
  margin-bottom: 1rem;
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
  background: #22c55e; /* default success for inverse metrics */
}

.progress-red.success {
  background: #ef4444; /* violation */
}

.progress-purple {
  background: #a855f7;
}

.metric-footer {
  display: flex;
  justify-content: space-between;
  align-items: center;
  font-size: 0.8rem;
  color: #94a3b8;
}

.meta-status {
  font-weight: 600;
  padding: 0.1rem 0.4rem;
  border-radius: 4px;
  font-size: 0.75rem;
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

/* Funnel Styling */
.funnel-card {
  background: rgba(30, 41, 59, 0.4);
  border: 1px solid rgba(255, 255, 255, 0.05);
  border-radius: 16px;
  padding: 2rem;
  box-shadow: 0 4px 6px -1px rgba(0, 0, 0, 0.1);
}

.funnel-card h3 {
  margin-top: 0;
  margin-bottom: 2rem;
  color: #38bdf8;
  font-size: 1.25rem;
}

.funnel-container {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 1rem;
  flex-wrap: wrap;
}

.funnel-stage {
  flex: 1;
  min-width: 160px;
  padding: 1.5rem;
  border-radius: 12px;
  text-align: center;
  box-shadow: 0 4px 10px rgba(0, 0, 0, 0.2);
}

.stage-registered {
  background: linear-gradient(135deg, rgba(56, 189, 248, 0.15), rgba(2, 132, 199, 0.15));
  border: 1px solid rgba(56, 189, 248, 0.3);
}

.stage-participants {
  background: linear-gradient(135deg, rgba(168, 85, 247, 0.15), rgba(126, 34, 206, 0.15));
  border: 1px solid rgba(168, 85, 247, 0.3);
}

.stage-active {
  background: linear-gradient(135deg, rgba(52, 211, 153, 0.15), rgba(5, 150, 105, 0.15));
  border: 1px solid rgba(52, 211, 153, 0.3);
}

.stage-name {
  font-size: 0.9rem;
  color: #cbd5e1;
  text-transform: uppercase;
  letter-spacing: 0.05em;
  margin-bottom: 0.5rem;
  font-weight: 500;
}

.stage-value {
  font-size: 2.2rem;
  font-weight: 800;
  color: #f8fafc;
  margin-bottom: 0.5rem;
}

.stage-pct {
  font-size: 0.8rem;
  color: #94a3b8;
}

.funnel-connector {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  color: #475569;
}

.connector-arrow {
  font-size: 1.5rem;
}

.connector-value {
  font-size: 0.8rem;
  color: #38bdf8;
  font-weight: 600;
  white-space: nowrap;
}

/* Analytics Grid and Table */
.analytics-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(480px, 1fr));
  gap: 2rem;
}

.mt-4 {
  margin-top: 2rem;
}

.full-width-card {
  grid-column: 1 / -1;
}

.analytics-card {
  background: rgba(30, 41, 59, 0.4);
  border: 1px solid rgba(255, 255, 255, 0.05);
  border-radius: 16px;
  padding: 2rem;
  box-shadow: 0 4px 6px -1px rgba(0, 0, 0, 0.1);
}

.analytics-card h3 {
  margin-top: 0;
  margin-bottom: 1.5rem;
  color: #38bdf8;
  font-size: 1.2rem;
}

.table-container {
  overflow-x: auto;
}

.cohort-table {
  width: 100%;
  border-collapse: collapse;
  text-align: left;
}

.cohort-table th, .cohort-table td {
  padding: 0.85rem 1rem;
  border-bottom: 1px solid rgba(255, 255, 255, 0.05);
}

.cohort-table th {
  background: rgba(15, 23, 42, 0.4);
  color: #94a3b8;
  font-weight: 600;
  font-size: 0.8rem;
  text-transform: uppercase;
  letter-spacing: 0.05em;
}

.cohort-table tbody tr:hover {
  background: rgba(255, 255, 255, 0.02);
}

.cohort-table tbody tr.target-met {
  background: rgba(16, 185, 129, 0.04);
}

.cohort-table tbody tr.target-met:hover {
  background: rgba(16, 185, 129, 0.07);
}

.cohort-month {
  font-weight: 600;
  color: #f8fafc;
}

.cohort-rate-cell {
  display: flex;
  align-items: center;
  gap: 0.75rem;
}

.rate-value {
  font-weight: 600;
}

.badge-target-met {
  background: rgba(16, 185, 129, 0.15);
  color: #34d399;
  padding: 0.1rem 0.4rem;
  border-radius: 4px;
  font-size: 0.7rem;
  font-weight: 600;
}

.empty-table-text {
  text-align: center;
  color: #64748b;
  padding: 2rem;
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
  padding: 3rem;
  background: rgba(30, 41, 59, 0.3);
  border-radius: 12px;
  border: 1px dashed rgba(255, 255, 255, 0.1);
  color: #94a3b8;
}
</style>
