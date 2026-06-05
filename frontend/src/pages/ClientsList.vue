<template>
  <div class="dashboard-layout">


    <div class="main-content">
      <div class="actions-row">
        <router-link to="/dashboard" class="back-link">← Voltar para Dashboard</router-link>
        <button class="add-btn" @click="showAddModal = true">➕ Novo Cliente</button>
      </div>

      <!-- Add Client Modal -->
      <div v-if="showAddModal" class="modal-overlay">
        <div class="modal">
          <h3>Cadastrar Novo Cliente</h3>
          <form @submit.prevent="submitAddClient">
            <div class="form-group">
              <label for="modal-name">Nome do Cliente</label>
              <input id="modal-name" v-model="newClientName" type="text" placeholder="Razão Social / Nome Fantasia" required />
            </div>
            <div class="form-group">
              <label for="modal-cnpj">CNPJ</label>
              <input id="modal-cnpj" v-model="newClientCNPJ" type="text" placeholder="00.000.000/0000-00" required />
            </div>
            <div v-if="modalError" class="modal-error">{{ modalError }}</div>
            <div class="modal-actions">
              <button type="button" class="cancel-btn" @click="closeModal">Cancelar</button>
              <button type="submit" class="submit-btn" :disabled="submitting">Salvar</button>
            </div>
          </form>
        </div>
      </div>

      <!-- List -->
      <div v-if="clientsStore.loading" class="loading-container">
        Carregando clientes...
      </div>
      <div v-else-if="clientsStore.error" class="error-container">
        {{ clientsStore.error }}
      </div>
      <div v-else-if="clientsStore.clients?.length === 0" class="empty-container">
        Nenhum cliente cadastrado.
      </div>
      <div v-else class="clients-grid">
        <div v-for="client in clientsStore.clients" :key="client.id" class="client-card">
          <h3>{{ client.name }}</h3>
          <p class="cnpj-text"><strong>CNPJ:</strong> {{ client.cnpj }}</p>
          <div class="card-footer">
            <router-link :to="'/clients/' + client.id" class="details-btn">Ver Detalhes</router-link>
            <button v-if="auth.isAdmin" class="delete-btn" @click="handleDelete(client.id)">Excluir</button>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue';
import { useRouter } from 'vue-router';
import { useAuthStore } from '../stores/auth';
import { useClientsStore } from '../stores/clients';

const auth = useAuthStore();
const clientsStore = useClientsStore();
const router = useRouter();

const showAddModal = ref(false);
const newClientName = ref('');
const newClientCNPJ = ref('');
const modalError = ref('');
const submitting = ref(false);

onMounted(() => {
  clientsStore.fetchClients();
});

const handleLogout = () => {
  auth.logout();
  router.push('/login');
};

const closeModal = () => {
  showAddModal.value = false;
  newClientName.value = '';
  newClientCNPJ.value = '';
  modalError.value = '';
};

const submitAddClient = async () => {
  modalError.value = '';
  submitting.value = true;
  try {
    await clientsStore.addClient({
      name: newClientName.value,
      cnpj: newClientCNPJ.value
    });
    closeModal();
  } catch (err) {
    modalError.value = err.response?.data?.error || 'Erro ao salvar cliente. Verifique se o CNPJ é único.';
  } finally {
    submitting.value = false;
  }
};

const handleDelete = async (id) => {
  if (confirm('Tem certeza que deseja excluir este cliente? Todos os projetos associados serão removidos.')) {
    try {
      await clientsStore.removeClient(id);
    } catch (err) {
      alert('Erro ao excluir cliente.');
    }
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
  max-width: 1200px;
  margin: 0 auto;
}

.actions-row {
  display: flex;
  justify-content: space-between;
  align-items: center;
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

.modal .form-group {
  margin-bottom: 1.5rem;
  text-align: left;
}

.modal label {
  display: block;
  margin-bottom: 0.5rem;
  font-size: 0.875rem;
  color: #94a3b8;
}

.modal input {
  width: 100%;
  padding: 0.75rem 1rem;
  border-radius: 8px;
  border: 1px solid #334155;
  background: #0f172a !important;
  color: #f8fafc !important;
  box-sizing: border-box;
  font-size: 1rem;
}

.add-btn {
  background: #38bdf8;
  color: #0f172a;
  border: none;
  padding: 0.75rem 1.25rem;
  border-radius: 8px;
  font-weight: 600;
  cursor: pointer;
  transition: background-color 0.2s;
}

.add-btn:hover {
  background: #0ea5e9;
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
}

.client-card h3 {
  margin-top: 0;
  margin-bottom: 1rem;
  font-size: 1.25rem;
  color: #f8fafc;
}

.cnpj-text {
  color: #94a3b8;
  font-size: 0.9rem;
  margin-bottom: 1.5rem;
}

.card-footer {
  display: flex;
  justify-content: space-between;
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

.delete-btn {
  background: transparent;
  border: none;
  color: #ef4444;
  cursor: pointer;
  font-weight: 600;
  padding: 0;
}

.delete-btn:hover {
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
