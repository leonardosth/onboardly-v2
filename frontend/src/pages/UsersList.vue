<template>
  <div class="dashboard-layout">
    <div class="main-content">
      <div class="actions-row">
        <h1 class="page-title">Gerenciamento de Usuários</h1>
        <button class="add-btn" @click="showAddModal = true">➕ Novo Usuário</button>
      </div>

      <!-- Add User Modal -->
      <div v-if="showAddModal" class="modal-overlay">
        <div class="modal">
          <h3>Cadastrar Novo Usuário</h3>
          <form @submit.prevent="submitAddUser">
            <div class="form-group">
              <label for="modal-email">E-mail</label>
              <input id="modal-email" v-model="newUserEmail" type="email" placeholder="usuario@empresa.com" required />
            </div>
            <div class="form-group">
              <label for="modal-password">Senha</label>
              <input id="modal-password" v-model="newUserPassword" type="password" placeholder="Mínimo 8 caracteres, 1 letra, 1 número" required />
            </div>
            <div class="form-group">
              <label for="modal-role">Papel</label>
              <select id="modal-role" v-model="newUserRole" class="role-select" required>
                <option value="Analista">Analista</option>
                <option value="Admin">Administrador</option>
              </select>
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
      <div v-if="usersStore.loading" class="loading-container">
        Carregando usuários...
      </div>
      <div v-else-if="usersStore.error" class="error-container">
        {{ usersStore.error }}
      </div>
      <div v-else-if="usersStore.users?.length === 0" class="empty-container">
        Nenhum usuário encontrado.
      </div>
      <div v-else class="clients-grid">
        <div v-for="user in usersStore.users" :key="user.id" class="client-card">
          <div class="card-header">
            <h3>{{ user.email }}</h3>
            <span class="status-badge" :class="user.role.toLowerCase()">
              {{ user.role }}
            </span>
          </div>
          <p class="cnpj-text">
            <strong>Criado em:</strong> {{ new Date(user.created_at).toLocaleDateString('pt-BR') }}
          </p>
          <div class="card-footer">
            <!-- Prevent deleting oneself visually if possible, but backend enforces it anyway -->
            <button 
              v-if="auth.email !== user.email" 
              class="delete-btn" 
              @click="handleDelete(user.id, user.email)"
            >Excluir</button>
            <span v-else class="you-badge">Você</span>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue';
import { useAuthStore } from '../stores/auth';
import { useUsersStore } from '../stores/users';

const auth = useAuthStore();
const usersStore = useUsersStore();

const showAddModal = ref(false);
const newUserEmail = ref('');
const newUserPassword = ref('');
const newUserRole = ref('Analista');
const modalError = ref('');
const submitting = ref(false);

onMounted(() => {
  usersStore.fetchUsers();
});

const closeModal = () => {
  showAddModal.value = false;
  newUserEmail.value = '';
  newUserPassword.value = '';
  newUserRole.value = 'Analista';
  modalError.value = '';
};

const submitAddUser = async () => {
  modalError.value = '';
  submitting.value = true;
  try {
    await usersStore.addUser({
      email: newUserEmail.value,
      password: newUserPassword.value,
      role: newUserRole.value
    });
    closeModal();
  } catch (err) {
    modalError.value = err.response?.data?.error || 'Erro ao cadastrar usuário.';
  } finally {
    submitting.value = false;
  }
};

const handleDelete = async (id, email) => {
  if (confirm(`Tem certeza que deseja excluir o usuário ${email}?`)) {
    try {
      await usersStore.removeUser(id);
    } catch (err) {
      alert(err.response?.data?.error || 'Erro ao excluir usuário.');
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
  font-size: 1.1rem;
  color: #f8fafc;
  word-break: break-all;
}

.status-badge {
  font-size: 0.75rem;
  padding: 0.2rem 0.5rem;
  border-radius: 9999px;
  font-weight: 600;
  white-space: nowrap;
}

.status-badge.admin { background: rgba(56, 189, 248, 0.2); color: #38bdf8; }
.status-badge.analista { background: rgba(148, 163, 184, 0.2); color: #cbd5e1; }

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

.you-badge {
  color: #94a3b8;
  font-size: 0.85rem;
  font-style: italic;
}

.loading-container, .error-container, .empty-container {
  text-align: center;
  padding: 3rem;
  background: rgba(30, 41, 59, 0.3);
  border-radius: 12px;
  border: 1px dashed rgba(255, 255, 255, 0.1);
  color: #94a3b8;
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

.modal input, .modal select {
  width: 100%;
  padding: 0.75rem 1rem;
  border-radius: 8px;
  border: 1px solid #334155;
  background: #0f172a !important;
  color: #f8fafc !important;
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

.submit-btn:disabled {
  opacity: 0.5;
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
</style>
