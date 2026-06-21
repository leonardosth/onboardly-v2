<template>
  <div class="modal-overlay">
    <div class="modal">
      <h3>📅 Agendar Reunião</h3>
      <form @submit.prevent="submitSchedule">
        <div class="form-group" v-if="!props.projectId">
          <label for="project-select">Projeto</label>
          <select id="project-select" v-model="selectedProjectId" required :disabled="loadingProjects">
            <option value="" disabled>Selecione um projeto...</option>
            <option v-for="p in projects" :key="p.id" :value="p.id">{{ p.name }}</option>
          </select>
        </div>
        <div class="form-group">
          <label for="meeting-title">Assunto / Título</label>
          <input
            id="meeting-title"
            v-model="title"
            type="text"
            placeholder="Ex: Reunião de Alinhamento Técnico"
            required
          />
        </div>
        <div class="form-row">
          <div class="form-group">
            <label for="meeting-date">Data</label>
            <input
              id="meeting-date"
              v-model="date"
              type="date"
              required
            />
          </div>
          <div class="form-group">
            <label for="meeting-time">Horário</label>
            <input
              id="meeting-time"
              v-model="timeVal"
              type="time"
              required
            />
          </div>
        </div>
        
        <div v-if="error" class="modal-error">
          {{ error }}
        </div>
        
        <div class="modal-actions">
          <button type="button" class="cancel-btn" @click="$emit('close')">Cancelar</button>
          <button type="submit" class="submit-btn" :disabled="loading">
            {{ loading ? 'Carregando...' : 'Agendar' }}
          </button>
        </div>
      </form>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue';
import { useMeetingsStore } from '../stores/meetings';
import { useAuthStore } from '../stores/auth';
import * as projectService from '../services/project';

const props = defineProps({
  projectId: {
    type: String,
    required: false,
    default: ''
  }
});

const emit = defineEmits(['close', 'scheduled']);

const selectedProjectId = ref(props.projectId);
const title = ref('');
const date = ref('');
const timeVal = ref('');
const error = ref('');
const loading = ref(false);

const projects = ref([]);
const loadingProjects = ref(false);

const meetingsStore = useMeetingsStore();
const authStore = useAuthStore();

onMounted(async () => {
  if (!props.projectId) {
    loadingProjects.value = true;
    try {
      projects.value = await projectService.getProjects() || [];
    } catch (err) {
      error.value = 'Erro ao carregar projetos.';
    } finally {
      loadingProjects.value = false;
    }
  }
});

const submitSchedule = async () => {
  error.value = '';
  
  if (!selectedProjectId.value) {
    error.value = 'Selecione um projeto.';
    return;
  }

  loading.value = true;
  
  try {
    // Combine date and time inputs into ISO 8601 string
    const localDateTimeStr = `${date.value}T${timeVal.value}:00`;
    const scheduledAt = new Date(localDateTimeStr).toISOString();

    await meetingsStore.scheduleMeeting({
      project_id: selectedProjectId.value,
      analyst_id: authStore.email, // Passing username/email or default ID
      title: title.value,
      scheduled_at: scheduledAt
    });

    emit('scheduled');
    emit('close');
  } catch (err) {
    error.value = err.response?.data?.error || 'Erro ao agendar reunião.';
  } finally {
    loading.value = false;
  }
};
</script>

<style scoped>
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
  max-width: 450px;
  box-shadow: 0 20px 25px -5px rgba(0, 0, 0, 0.5);
  box-sizing: border-box;
}

.modal h3 {
  margin-top: 0;
  margin-bottom: 1.5rem;
  color: #38bdf8;
}

.form-group {
  margin-bottom: 1.5rem;
}

.form-row {
  display: flex;
  gap: 1rem;
}

.form-row .form-group {
  flex: 1;
}

label {
  display: block;
  margin-bottom: 0.5rem;
  font-size: 0.875rem;
  color: #94a3b8;
}

input, select {
  width: 100%;
  padding: 0.75rem 1rem;
  border-radius: 8px;
  border: 1px solid #334155;
  background: #0f172a;
  color: #f8fafc;
  box-sizing: border-box;
  font-size: 1rem;
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
</style>
