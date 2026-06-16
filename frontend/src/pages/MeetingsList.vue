<template>
  <div class="dashboard-layout">
    <div class="main-content">
      <div class="actions-row">
        <h1 class="page-title">Minhas Agendas</h1>
      </div>

      <!-- Filters -->
      <div class="filters-row">
        <select v-model="statusFilter" class="status-select" @change="onStatusFilterChange">
          <option value="completed">Concluídas</option>
          <option value="scheduled">Agendadas</option>
          <option value="">Todas</option>
        </select>
      </div>

      <!-- List / Table -->
      <div v-if="meetingsStore.loading" class="loading-container">
        Carregando agendas...
      </div>
      <div v-else-if="meetingsStore.error" class="error-container">
        {{ meetingsStore.error }}
      </div>
      <div v-else-if="meetingsStore.myMeetings?.length === 0" class="empty-container">
        Nenhuma agenda encontrada.
      </div>
      <div v-else class="table-container">
        <table class="meetings-table">
          <thead>
            <tr>
              <th>Título</th>
              <th>Cliente</th>
              <th>Projeto</th>
              <th>Data Agendada</th>
              <th>Status</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="meeting in meetingsStore.myMeetings" :key="meeting.id">
              <td class="meeting-title">{{ meeting.title }}</td>
              <td>{{ meeting.client_name || 'N/A' }}</td>
              <td>
                <router-link :to="'/projects/' + meeting.project_id" class="project-link">
                  {{ meeting.project_name || 'N/A' }}
                </router-link>
              </td>
              <td>{{ formatDateTime(meeting.scheduled_at) }}</td>
              <td>
                <span class="status-badge" :class="meeting.status.toLowerCase()">
                  {{ meeting.status === 'completed' ? 'Concluída' : 'Agendada' }}
                </span>
              </td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue';
import { useMeetingsStore } from '../stores/meetings';

const meetingsStore = useMeetingsStore();
const statusFilter = ref('completed');

onMounted(() => {
  fetchMyMeetings();
});

const fetchMyMeetings = () => {
  meetingsStore.fetchMyMeetings(statusFilter.value);
};

const onStatusFilterChange = () => {
  fetchMyMeetings();
};

const formatDateTime = (dateStr) => {
  if (!dateStr) return '';
  const date = new Date(dateStr);
  return date.toLocaleString('pt-BR', {
    day: '2-digit',
    month: '2-digit',
    year: 'numeric',
    hour: '2-digit',
    minute: '2-digit'
  });
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

.page-title {
  margin: 0;
  font-size: 1.5rem;
  color: #f8fafc;
}

.actions-row {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 2rem;
}

.filters-row {
  display: flex;
  gap: 1rem;
  margin-bottom: 2rem;
}

.status-select {
  padding: 0.75rem 1rem;
  border-radius: 8px;
  border: 1px solid #334155;
  background: rgba(30, 41, 59, 0.5);
  color: #f8fafc;
  font-size: 1rem;
  outline: none;
  width: 200px;
}

.status-select:focus {
  border-color: #38bdf8;
}

.table-container {
  background: rgba(30, 41, 59, 0.5);
  border: 1px solid rgba(255, 255, 255, 0.05);
  border-radius: 12px;
  overflow: hidden;
  box-shadow: 0 4px 6px -1px rgba(0, 0, 0, 0.1);
}

.meetings-table {
  width: 100%;
  border-collapse: collapse;
  text-align: left;
}

.meetings-table th, .meetings-table td {
  padding: 1rem 1.5rem;
  border-bottom: 1px solid rgba(255, 255, 255, 0.05);
}

.meetings-table th {
  background: rgba(15, 23, 42, 0.4);
  color: #94a3b8;
  font-weight: 600;
  font-size: 0.875rem;
  text-transform: uppercase;
  letter-spacing: 0.05em;
}

.meetings-table tbody tr:last-child td {
  border-bottom: none;
}

.meetings-table tbody tr:hover {
  background: rgba(255, 255, 255, 0.02);
}

.meeting-title {
  font-weight: 500;
  color: #f8fafc;
}

.project-link {
  color: #38bdf8;
  text-decoration: none;
  font-weight: 500;
}

.project-link:hover {
  text-decoration: underline;
}

.status-badge {
  font-size: 0.75rem;
  padding: 0.2rem 0.5rem;
  border-radius: 9999px;
  font-weight: 600;
  white-space: nowrap;
}

.status-badge.completed {
  background: rgba(52, 211, 153, 0.2);
  color: #34d399;
}

.status-badge.scheduled {
  background: rgba(251, 191, 36, 0.2);
  color: #fbbf24;
}

.loading-container, .error-container, .empty-container {
  text-align: center;
  padding: 3rem;
  background: rgba(30, 41, 59, 0.3);
  border-radius: 12px;
  border: 1px dashed rgba(255, 255, 255, 0.1);
  color: #94a3b8;
}
</style>
