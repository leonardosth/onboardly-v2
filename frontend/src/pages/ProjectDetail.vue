<template>
  <div class="dashboard-layout">


    <div class="main-content">
      <div class="actions-row">
        <button class="back-link" @click="goBack">← Voltar para o Cliente</button>
      </div>

      <div v-if="projectsStore.loading || meetingsStore.loading" class="loading-container">
        Carregando detalhes do projeto...
      </div>
      <div v-else-if="projectsStore.error" class="error-container">
        {{ projectsStore.error }}
      </div>
      <div v-else-if="projectsStore.currentProject" class="project-details-container">
        <!-- Project Info Card -->
        <div class="project-info-card">
          <div class="card-header">
            <h2>{{ projectsStore.currentProject.name }}</h2>
            <div class="badges">
              <span :class="'status-badge ' + projectsStore.currentProject.status.toLowerCase().replace(' ', '-')">
                {{ projectsStore.currentProject.status }}
              </span>
              <span class="active-badge" :class="projectsStore.currentProject.is_active ? 'active' : 'inactive'">
                {{ projectsStore.currentProject.is_active ? 'Ativo' : 'Inativo' }}
              </span>
            </div>
          </div>
          
          <!-- Status Progression Controller -->
          <div class="status-progression">
            <h3>Progresso da Implantação</h3>
            <div class="status-buttons">
              <button
                class="status-btn"
                :class="{ active: projectsStore.currentProject.status === 'Backlog' }"
                @click="updateProjectStatus('Backlog')"
              >
                Backlog
              </button>
              <button
                class="status-btn"
                :class="{ active: projectsStore.currentProject.status === 'Em andamento' }"
                @click="updateProjectStatus('Em andamento')"
              >
                Em andamento
              </button>
              <button
                class="status-btn"
                :class="{ active: projectsStore.currentProject.status === 'Go-Live' }"
                @click="updateProjectStatus('Go-Live')"
              >
                Go-Live
              </button>
            </div>
          </div>
        </div>

        <!-- Meetings Section -->
        <div class="meetings-section">
          <div class="section-header">
            <h3>📅 Reuniões & Interações</h3>
            <button class="add-meeting-btn" @click="showScheduleModal = true">📅 Agendar Reunião</button>
          </div>

          <div v-if="meetingsStore.meetings.length === 0" class="empty-container">
            Nenhuma reunião agendada para este projeto.
          </div>
          <div v-else class="meetings-list">
            <div v-for="meeting in meetingsStore.meetings" :key="meeting.id" class="meeting-item">
              <div class="meeting-details">
                <h4>{{ meeting.title }}</h4>
                <p class="meeting-time">
                  <strong>Horário:</strong> {{ new Date(meeting.scheduled_at).toLocaleString() }}
                </p>
              </div>
              <span class="no-show-badge" :class="{ alert: meeting.no_show }">
                {{ meeting.no_show ? 'No-Show' : 'Compareceu' }}
              </span>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- Meeting Scheduler Modal Component -->
    <MeetingScheduler
      v-if="showScheduleModal"
      :project-id="route.params.id"
      @close="showScheduleModal = false"
      @scheduled="onMeetingScheduled"
    />
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue';
import { useRoute, useRouter } from 'vue-router';
import { useAuthStore } from '../stores/auth';
import { useProjectsStore } from '../stores/projects';
import { useMeetingsStore } from '../stores/meetings';
import MeetingScheduler from '../components/MeetingScheduler.vue';

const auth = useAuthStore();
const projectsStore = useProjectsStore();
const meetingsStore = useMeetingsStore();
const route = useRoute();
const router = useRouter();

const showScheduleModal = ref(false);

onMounted(async () => {
  const projectId = route.params.id;
  await projectsStore.fetchProject(projectId);
  meetingsStore.fetchMeetings(projectId);
});

const goBack = () => {
  if (projectsStore.currentProject) {
    router.push(`/clients/${projectsStore.currentProject.client_id}`);
  } else {
    router.push('/clients');
  }
};

const handleLogout = () => {
  auth.logout();
  router.push('/login');
};

const updateProjectStatus = async (newStatus) => {
  try {
    await projectsStore.updateStatus(route.params.id, newStatus);
  } catch (err) {
    alert('Erro ao atualizar status do projeto.');
  }
};

const onMeetingScheduled = () => {
  meetingsStore.fetchMeetings(route.params.id);
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

.logout-btn {
  background: transparent;
  border: 1px solid #ef4444;
  color: #ef4444;
  padding: 0.5rem 1rem;
  border-radius: 6px;
  cursor: pointer;
}

.logout-btn:hover {
  background: rgba(239, 68, 68, 0.1);
}

.main-content {
  padding: 2.5rem;
  max-width: 1000px;
  margin: 0 auto;
}

.actions-row {
  margin-bottom: 2rem;
}

.back-link {
  background: transparent;
  border: none;
  color: #94a3b8;
  cursor: pointer;
  font-weight: 500;
  font-size: 1rem;
}

.back-link:hover {
  color: #38bdf8;
}

.project-info-card {
  background: rgba(30, 41, 59, 0.5);
  border: 1px solid rgba(255, 255, 255, 0.05);
  border-radius: 12px;
  padding: 2rem;
  margin-bottom: 2.5rem;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 2rem;
}

.card-header h2 {
  margin: 0;
  color: #f8fafc;
}

.badges {
  display: flex;
  gap: 1rem;
  align-items: center;
}

.status-badge {
  font-size: 0.85rem;
  font-weight: 600;
  padding: 0.35rem 0.75rem;
  border-radius: 9999px;
}

.status-badge.backlog {
  background: rgba(148, 163, 184, 0.2);
  color: #94a3b8;
}

.status-badge.em-andamento {
  background: rgba(56, 189, 248, 0.2);
  color: #38bdf8;
}

.status-badge.go-live {
  background: rgba(34, 197, 94, 0.2);
  color: #4ade80;
}

.active-badge {
  font-size: 0.8rem;
  font-weight: 600;
}

.active-badge.active {
  color: #4ade80;
}

.active-badge.inactive {
  color: #ef4444;
}

.status-progression h3 {
  margin-top: 0;
  margin-bottom: 1rem;
  font-size: 1.1rem;
  color: #94a3b8;
}

.status-buttons {
  display: flex;
  gap: 1rem;
}

.status-btn {
  background: transparent;
  border: 1px solid #475569;
  color: #94a3b8;
  padding: 0.6rem 1.5rem;
  border-radius: 8px;
  cursor: pointer;
  font-weight: 600;
  transition: all 0.2s;
}

.status-btn.active {
  background: #38bdf8;
  border-color: #38bdf8;
  color: #0f172a;
}

.meetings-section {
  background: rgba(30, 41, 59, 0.3);
  border: 1px solid rgba(255, 255, 255, 0.05);
  border-radius: 12px;
  padding: 2rem;
}

.section-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 2rem;
}

.section-header h3 {
  margin: 0;
  font-size: 1.35rem;
  color: #38bdf8;
}

.add-meeting-btn {
  background: #38bdf8;
  color: #0f172a;
  border: none;
  padding: 0.5rem 1.25rem;
  border-radius: 8px;
  font-weight: 600;
  cursor: pointer;
}

.add-meeting-btn:hover {
  background: #0ea5e9;
}

.meetings-list {
  display: flex;
  flex-direction: column;
  gap: 1rem;
}

.meeting-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
  background: rgba(30, 41, 59, 0.5);
  padding: 1.25rem 1.5rem;
  border-radius: 8px;
  border: 1px solid rgba(255, 255, 255, 0.02);
}

.meeting-details h4 {
  margin: 0 0 0.5rem 0;
  font-size: 1.05rem;
  color: #f8fafc;
}

.meeting-time {
  margin: 0;
  font-size: 0.85rem;
  color: #94a3b8;
}

.no-show-badge {
  font-size: 0.8rem;
  font-weight: 600;
  color: #4ade80;
}

.no-show-badge.alert {
  color: #ef4444;
}

.loading-container, .empty-container {
  text-align: center;
  padding: 2.5rem;
  color: #94a3b8;
}
</style>
