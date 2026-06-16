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
              <span :class="'status-badge ' + projectsStore.currentProject.status.toLowerCase().replace(/\s+/g, '-')">
                {{ projectsStore.currentProject.status }}
              </span>
              <span class="active-badge" :class="projectsStore.currentProject.activated_at ? 'active' : 'inactive'">
                {{ projectsStore.currentProject.activated_at ? 'Ativo' : 'Inativo' }}
              </span>
            </div>
          </div>
          
          <div class="project-meta-details">
            <p v-if="projectsStore.currentProject.activated_at">
              <strong>Ativado em:</strong> {{ formatDateTime(projectsStore.currentProject.activated_at) }}
            </p>
            <p v-else class="text-warning">
              ⚠️ Este cliente ainda não foi ativado.
            </p>
          </div>

          <!-- Action Buttons -->
          <div class="action-buttons-row">
            <button 
              v-if="!projectsStore.currentProject.activated_at && projectsStore.currentProject.status !== 'Go-Live'"
              class="action-btn activate-btn"
              @click="openActivationModal"
            >
              🚀 Marcar como Ativo
            </button>
            <button 
              v-if="projectsStore.currentProject.status !== 'Go-Live'"
              class="action-btn finalize-btn"
              @click="confirmFinalizeProject"
            >
              🏁 Finalizar Projeto
            </button>
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
                <div class="meeting-title-row">
                  <h4>{{ meeting.title }}</h4>
                  <span class="status-badge-mini" :class="meeting.status">
                    {{ meeting.status === 'completed' ? 'Concluída' : 'Agendada' }}
                  </span>
                </div>
                <p class="meeting-time">
                  <strong>Horário:</strong> {{ formatDateTime(meeting.scheduled_at) }}
                </p>
                <p v-if="meeting.status === 'completed' && meeting.completed_at" class="meeting-time">
                  <strong>Concluída em:</strong> {{ formatDateTime(meeting.completed_at) }}
                </p>
              </div>
              <div class="meeting-actions">
                <button 
                  v-if="meeting.status === 'scheduled'" 
                  class="complete-meeting-inline-btn"
                  @click="openCompleteModalForMeeting(meeting)"
                >
                  Concluir
                </button>
                <span v-else class="no-show-badge" :class="{ alert: meeting.no_show }">
                  {{ meeting.no_show ? 'No-Show' : 'Compareceu' }}
                </span>
              </div>
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

    <!-- Activation & Completion Modal -->
    <div v-if="showActivationModal" class="modal-overlay">
      <div class="modal">
        <h3>Concluir Reunião e Ativar</h3>
        <form @submit.prevent="submitActivation">
          <div class="form-group">
            <label for="meeting-select">Selecione a Reunião</label>
            <select id="meeting-select" v-model="selectedMeetingId" required class="form-select">
              <option value="" disabled>Selecione uma reunião agendada...</option>
              <option 
                v-for="m in scheduledMeetings" 
                :key="m.id" 
                :value="m.id"
              >
                {{ m.title }} ({{ formatDateTime(m.scheduled_at) }})
              </option>
            </select>
            <p v-if="scheduledMeetings.length === 0" class="text-error mt-2">
              Nenhuma reunião agendada encontrada. Agende uma reunião primeiro.
            </p>
          </div>

          <div class="form-checkbox-group">
            <input 
              id="activate-client-checkbox" 
              v-model="activateClientCheckbox" 
              type="checkbox" 
            />
            <label for="activate-client-checkbox">
              Marcar cliente como ativo (executou tarefas durante a implantação)
            </label>
          </div>

          <div v-if="activationError" class="modal-error">{{ activationError }}</div>

          <div class="modal-actions">
            <button type="button" class="cancel-btn" @click="closeActivationModal">Cancelar</button>
            <button 
              type="submit" 
              class="submit-btn" 
              :disabled="submittingActivation || !selectedMeetingId"
            >
              Salvar
            </button>
          </div>
        </form>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue';
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
const showActivationModal = ref(false);
const selectedMeetingId = ref('');
const activateClientCheckbox = ref(true);
const submittingActivation = ref(false);
const activationError = ref('');

const scheduledMeetings = computed(() => {
  return meetingsStore.meetings.filter(m => m.status === 'scheduled');
});

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

const openActivationModal = () => {
  activationError.value = '';
  // Default to selecting the first scheduled meeting if available
  const scheduled = scheduledMeetings.value;
  if (scheduled.length > 0) {
    selectedMeetingId.value = scheduled[0].id;
  } else {
    selectedMeetingId.value = '';
  }
  activateClientCheckbox.value = true;
  showActivationModal.value = true;
};

const openCompleteModalForMeeting = (meeting) => {
  activationError.value = '';
  selectedMeetingId.value = meeting.id;
  activateClientCheckbox.value = !projectsStore.currentProject.activated_at; // auto-check if not already active
  showActivationModal.value = true;
};

const closeActivationModal = () => {
  showActivationModal.value = false;
  selectedMeetingId.value = '';
  activationError.value = '';
};

const submitActivation = async () => {
  if (!selectedMeetingId.value) return;
  activationError.value = '';
  submittingActivation.value = true;
  try {
    const result = await meetingsStore.completeMeeting(selectedMeetingId.value, activateClientCheckbox.value);
    
    // Refresh project details to get updated activated_at / status
    await projectsStore.fetchProject(route.params.id);
    // Refresh meetings
    await meetingsStore.fetchMeetings(route.params.id);
    
    closeActivationModal();
  } catch (err) {
    activationError.value = err.response?.data?.error || 'Erro ao concluir reunião.';
  } finally {
    submittingActivation.value = false;
  }
};

const confirmFinalizeProject = async () => {
  const confirmMsg = 'Deseja finalizar o projeto e colocá-lo em Go-Live?';
  if (confirm(confirmMsg)) {
    try {
      await projectsStore.finalizeProject(route.params.id);
    } catch (err) {
      alert(err.response?.data?.error || 'Erro ao finalizar projeto.');
    }
  }
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
  margin-bottom: 1rem;
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
  padding: 0.25rem 0.6rem;
  border-radius: 6px;
}

.active-badge.active {
  background: rgba(52, 211, 153, 0.15);
  color: #4ade80;
}

.active-badge.inactive {
  background: rgba(239, 68, 68, 0.15);
  color: #ef4444;
}

.project-meta-details {
  margin-bottom: 1.5rem;
  color: #94a3b8;
  font-size: 0.95rem;
}

.text-warning {
  color: #fbbf24;
}

.action-buttons-row {
  display: flex;
  gap: 1rem;
  margin-bottom: 2rem;
  border-bottom: 1px solid rgba(255, 255, 255, 0.05);
  padding-bottom: 1.5rem;
}

.action-btn {
  padding: 0.75rem 1.5rem;
  border-radius: 8px;
  font-weight: 600;
  font-size: 0.95rem;
  cursor: pointer;
  border: none;
  transition: all 0.2s;
}

.activate-btn {
  background: #38bdf8;
  color: #0f172a;
}

.activate-btn:hover {
  background: #0ea5e9;
}

.finalize-btn {
  background: #10b981;
  color: #ffffff;
}

.finalize-btn:hover {
  background: #059669;
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
  background: rgba(56, 189, 248, 0.1);
  color: #38bdf8;
  border: 1px solid rgba(56, 189, 248, 0.2);
  padding: 0.5rem 1.25rem;
  border-radius: 8px;
  font-weight: 600;
  cursor: pointer;
  transition: all 0.2s;
}

.add-meeting-btn:hover {
  background: rgba(56, 189, 248, 0.2);
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

.meeting-title-row {
  display: flex;
  align-items: center;
  gap: 1rem;
  margin-bottom: 0.5rem;
}

.meeting-title-row h4 {
  margin: 0;
  font-size: 1.05rem;
  color: #f8fafc;
}

.status-badge-mini {
  font-size: 0.7rem;
  font-weight: 600;
  padding: 0.1rem 0.4rem;
  border-radius: 4px;
  text-transform: uppercase;
}

.status-badge-mini.completed {
  background: rgba(52, 211, 153, 0.15);
  color: #34d399;
}

.status-badge-mini.scheduled {
  background: rgba(251, 191, 36, 0.15);
  color: #fbbf24;
}

.meeting-time {
  margin: 0;
  font-size: 0.85rem;
  color: #94a3b8;
}

.complete-meeting-inline-btn {
  background: rgba(52, 211, 153, 0.1);
  color: #34d399;
  border: 1px solid rgba(52, 211, 153, 0.2);
  padding: 0.4rem 1rem;
  border-radius: 6px;
  font-weight: 600;
  cursor: pointer;
  transition: all 0.2s;
}

.complete-meeting-inline-btn:hover {
  background: rgba(52, 211, 153, 0.2);
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

/* Modal Styling */
.modal-overlay {
  position: fixed;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  background: rgba(15, 23, 42, 0.8);
  display: flex;
  justify-content: center;
  align-items: center;
  z-index: 100;
}

.modal {
  background: #1e293b;
  border: 1px solid rgba(255, 255, 255, 0.1);
  padding: 2.5rem;
  border-radius: 16px;
  width: 100%;
  max-width: 480px;
  box-shadow: 0 20px 25px -5px rgba(0, 0, 0, 0.5);
}

.modal h3 {
  margin-top: 0;
  margin-bottom: 1.5rem;
  color: #38bdf8;
}

.form-group {
  margin-bottom: 1.5rem;
}

label {
  display: block;
  margin-bottom: 0.5rem;
  font-size: 0.875rem;
  color: #94a3b8;
}

.form-select {
  width: 100%;
  padding: 0.75rem 1rem;
  border-radius: 8px;
  border: 1px solid #334155;
  background: #0f172a;
  color: #f8fafc;
  font-size: 1rem;
  box-sizing: border-box;
}

.form-checkbox-group {
  display: flex;
  align-items: center;
  gap: 0.75rem;
  margin-bottom: 1.5rem;
}

.form-checkbox-group input {
  width: auto;
  cursor: pointer;
}

.form-checkbox-group label {
  margin-bottom: 0;
  cursor: pointer;
  user-select: none;
}

.modal-actions {
  display: flex;
  justify-content: flex-end;
  gap: 1rem;
  margin-top: 2rem;
}

.cancel-btn {
  background: transparent;
  border: 1px solid #475569;
  color: #94a3b8;
  padding: 0.5rem 1.25rem;
  border-radius: 6px;
  cursor: pointer;
}

.cancel-btn:hover {
  background: rgba(255, 255, 255, 0.05);
}

.submit-btn {
  background: #38bdf8;
  color: #0f172a;
  border: none;
  padding: 0.5rem 1.25rem;
  border-radius: 6px;
  font-weight: 600;
  cursor: pointer;
}

.submit-btn:hover {
  background: #0ea5e9;
}

.submit-btn:disabled {
  background: #1e293b;
  border: 1px solid #334155;
  color: #475569;
  cursor: not-allowed;
}

.modal-error {
  color: #fca5a5;
  background: rgba(239, 68, 68, 0.2);
  border: 1px solid #ef4444;
  padding: 0.5rem;
  border-radius: 6px;
  font-size: 0.875rem;
  margin-top: 1rem;
  text-align: center;
}

.text-error {
  color: #ef4444;
  font-size: 0.875rem;
}

.mt-2 {
  margin-top: 0.5rem;
}
</style>
