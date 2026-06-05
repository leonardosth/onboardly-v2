<template>
  <div class="dashboard-layout">


    <div class="main-content">
      <div class="actions-row">
        <router-link to="/clients" class="back-link">← Voltar para Lista de Clientes</router-link>
      </div>

      <div v-if="clientsStore.loading" class="loading-container">
        Carregando detalhes do cliente...
      </div>
      <div v-else-if="clientsStore.error" class="error-container">
        {{ clientsStore.error }}
      </div>
      <div v-else-if="clientsStore.currentClient" class="client-details-container">
        <div class="client-info-card">
          <h2>{{ clientsStore.currentClient.name }}</h2>
          <p><strong>CNPJ:</strong> {{ clientsStore.currentClient.cnpj }}</p>
          <p><strong>Cadastrado em:</strong> {{ new Date(clientsStore.currentClient.created_at).toLocaleDateString() }}</p>
        </div>

        <!-- Projects section -->
        <div class="projects-section">
          <div class="section-header">
            <h3>🚀 Projetos de Implantação</h3>
            <button class="add-project-btn" @click="showAddProjectModal = true">➕ Novo Projeto</button>
          </div>

          <div v-if="projectsLoading" class="loading-container">
            Carregando projetos...
          </div>
          <div v-else-if="projects.length === 0" class="empty-container">
            Nenhum projeto cadastrado para este cliente.
          </div>
          <div v-else class="projects-list">
            <div v-for="project in projects" :key="project.id" class="project-item">
              <div class="project-meta">
                <h4>{{ project.name }}</h4>
                <span :class="'status-badge ' + project.status.toLowerCase().replace(' ', '-')">
                  {{ project.status }}
                </span>
                <span class="active-badge" :class="project.is_active ? 'active' : 'inactive'">
                  {{ project.is_active ? 'Ativo' : 'Inativo' }}
                </span>
              </div>
              <router-link :to="'/projects/' + project.id" class="manage-project-btn">Gerenciar</router-link>
            </div>
          </div>
        </div>
      </div>

      <!-- Add Project Modal -->
      <div v-if="showAddProjectModal" class="modal-overlay">
        <div class="modal">
          <h3>Criar Novo Projeto</h3>
          <form @submit.prevent="submitAddProject">
            <div class="form-group">
              <label for="project-name">Nome do Projeto</label>
              <input id="project-name" v-model="newProjectName" type="text" placeholder="Ex: Implantação Onboardly ERP" required />
            </div>
            <div v-if="projectModalError" class="modal-error">{{ projectModalError }}</div>
            <div class="modal-actions">
              <button type="button" class="cancel-btn" @click="closeProjectModal">Cancelar</button>
              <button type="submit" class="submit-btn" :disabled="submittingProject">Salvar</button>
            </div>
          </form>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue';
import { useRoute, useRouter } from 'vue-router';
import axios from 'axios';
import { useAuthStore } from '../stores/auth';
import { useClientsStore } from '../stores/clients';

const auth = useAuthStore();
const clientsStore = useClientsStore();
const route = useRoute();
const router = useRouter();

const projects = ref([]);
const projectsLoading = ref(false);

const showAddProjectModal = ref(false);
const newProjectName = ref('');
const projectModalError = ref('');
const submittingProject = ref(false);

onMounted(async () => {
  const clientId = route.params.id;
  await clientsStore.fetchClient(clientId);
  fetchProjects();
});

const fetchProjects = async () => {
  projectsLoading.value = true;
  try {
    const response = await axios.get(`http://localhost:8080/api/projects?client_id=${route.params.id}`);
    projects.value = response.data || [];
  } catch (err) {
    console.error('Error fetching projects:', err);
  } finally {
    projectsLoading.value = false;
  }
};

const handleLogout = () => {
  auth.logout();
  router.push('/login');
};

const closeProjectModal = () => {
  showAddProjectModal.value = false;
  newProjectName.value = '';
  projectModalError.value = '';
};

const submitAddProject = async () => {
  projectModalError.value = '';
  submittingProject.value = true;
  try {
    await axios.post('http://localhost:8080/api/projects', {
      client_id: route.params.id,
      name: newProjectName.value
    });
    closeProjectModal();
    fetchProjects();
  } catch (err) {
    projectModalError.value = err.response?.data?.error || 'Erro ao criar projeto.';
  } finally {
    submittingProject.value = false;
  }
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
  font-weight: 500;
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
  color: #94a3b8;
  text-decoration: none;
  font-weight: 500;
}

.back-link:hover {
  color: #38bdf8;
}

.client-info-card {
  background: rgba(30, 41, 59, 0.5);
  border: 1px solid rgba(255, 255, 255, 0.05);
  border-radius: 12px;
  padding: 2rem;
  margin-bottom: 2.5rem;
}

.client-info-card h2 {
  margin-top: 0;
  margin-bottom: 1rem;
  color: #f8fafc;
}

.client-info-card p {
  color: #94a3b8;
  margin: 0.5rem 0;
}

.projects-section {
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

.add-project-btn {
  background: #38bdf8;
  color: #0f172a;
  border: none;
  padding: 0.5rem 1rem;
  border-radius: 6px;
  font-weight: 600;
  cursor: pointer;
}

.add-project-btn:hover {
  background: #0ea5e9;
}

.projects-list {
  display: flex;
  flex-direction: column;
  gap: 1rem;
}

.project-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
  background: rgba(30, 41, 59, 0.5);
  padding: 1rem 1.5rem;
  border-radius: 8px;
  border: 1px solid rgba(255, 255, 255, 0.02);
}

.project-meta {
  display: flex;
  align-items: center;
  gap: 1.5rem;
}

.project-meta h4 {
  margin: 0;
  font-size: 1.1rem;
  color: #f8fafc;
}

.status-badge {
  font-size: 0.8rem;
  font-weight: 600;
  padding: 0.25rem 0.6rem;
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
  font-size: 0.75rem;
  font-weight: 500;
}

.active-badge.active {
  color: #4ade80;
}

.active-badge.inactive {
  color: #ef4444;
}

.manage-project-btn {
  color: #38bdf8;
  text-decoration: none;
  font-weight: 600;
  font-size: 0.9rem;
}

.manage-project-btn:hover {
  text-decoration: underline;
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
  max-width: 450px;
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

input {
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
