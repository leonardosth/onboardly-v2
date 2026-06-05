<template>
  <div class="dashboard-layout">
    <div class="main-content">
      <div class="actions-row">
        <h1 class="page-title">Projetos de Implantação</h1>
      </div>

      <!-- Filters -->
      <div class="filters-row">
        <input 
          v-model="searchQuery" 
          type="text" 
          placeholder="Buscar projeto ou cliente..." 
          class="search-input"
        />
        <select v-model="statusFilter" class="status-select">
          <option value="">Todos os Status</option>
          <option value="Backlog">Backlog</option>
          <option value="Em andamento">Em andamento</option>
          <option value="Go-Live">Go-Live</option>
        </select>
      </div>

      <!-- List -->
      <div v-if="projectsStore.loading" class="loading-container">
        Carregando projetos...
      </div>
      <div v-else-if="projectsStore.error" class="error-container">
        {{ projectsStore.error }}
      </div>
      <div v-else-if="projectsStore.projects?.length === 0" class="empty-container">
        Nenhum projeto cadastrado.
      </div>
      <div v-else-if="filteredProjects.length === 0" class="empty-container">
        Nenhum projeto encontrado para os filtros aplicados.
      </div>
      <div v-else class="clients-grid">
        <div v-for="project in filteredProjects" :key="project.id" class="client-card">
          <div class="card-header">
            <h3>{{ project.name }}</h3>
            <span class="status-badge" :class="project.status.replace(/\s+/g, '-').toLowerCase()">
              {{ project.status }}
            </span>
          </div>
          <p class="cnpj-text">
            <strong>Cliente:</strong> {{ project.client?.name || 'Cliente não encontrado' }}
          </p>
          <div class="card-footer">
            <router-link :to="'/projects/' + project.id" class="details-btn">Ver Detalhes</router-link>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue';
import { useProjectsStore } from '../stores/projects';

const projectsStore = useProjectsStore();

const searchQuery = ref('');
const statusFilter = ref('');

onMounted(() => {
  projectsStore.fetchProjects();
});

const filteredProjects = computed(() => {
  let list = projectsStore.projects || [];
  
  if (statusFilter.value) {
    list = list.filter(p => p.status === statusFilter.value);
  }
  
  if (searchQuery.value) {
    const q = searchQuery.value.toLowerCase();
    list = list.filter(p => 
      p.name.toLowerCase().includes(q) || 
      (p.client && p.client.name.toLowerCase().includes(q))
    );
  }
  
  return list;
});
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

.search-input, .status-select {
  padding: 0.75rem 1rem;
  border-radius: 8px;
  border: 1px solid #334155;
  background: rgba(30, 41, 59, 0.5);
  color: #f8fafc;
  font-size: 1rem;
  outline: none;
}

.search-input {
  flex: 1;
}

.search-input:focus, .status-select:focus {
  border-color: #38bdf8;
}

.clients-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(300px, 1fr));
  gap: 1.5rem;
}

.client-card {
  background: rgba(30, 41, 59, 0.5);
  border: 1px solid rgba(255, 255, 255, 0.05);
  border-radius: 12px;
  padding: 1.5rem;
  box-shadow: 0 4px 6px -1px rgba(0, 0, 0, 0.1);
  display: flex;
  flex-direction: column;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  margin-bottom: 1rem;
}

.client-card h3 {
  margin: 0;
  font-size: 1.25rem;
  color: #f8fafc;
}

.status-badge {
  font-size: 0.75rem;
  padding: 0.2rem 0.5rem;
  border-radius: 9999px;
  font-weight: 600;
  white-space: nowrap;
}

.status-badge.backlog { background: #334155; color: #cbd5e1; }
.status-badge.em-andamento { background: rgba(56, 189, 248, 0.2); color: #38bdf8; }
.status-badge.go-live { background: rgba(52, 211, 153, 0.2); color: #34d399; }

.cnpj-text {
  color: #94a3b8;
  font-size: 0.9rem;
  margin-bottom: 1.5rem;
  flex: 1;
}

.card-footer {
  display: flex;
  justify-content: flex-end;
  align-items: center;
}

.details-btn {
  color: #38bdf8;
  text-decoration: none;
  font-weight: 600;
  font-size: 0.95rem;
}

.details-btn:hover {
  text-decoration: underline;
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
