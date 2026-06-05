<template>
  <div class="auth-container">
    <div class="auth-card">
      <h2>🚀 Onboardly Cadastro</h2>
      <form @submit.prevent="handleRegister">
        <div class="form-group">
          <label for="email">E-mail</label>
          <input
            id="email"
            v-model="email"
            type="email"
            placeholder="Digite seu e-mail"
            required
          />
        </div>
        <div class="form-group">
          <label for="password">Senha</label>
          <input
            id="password"
            v-model="password"
            type="password"
            placeholder="Escolha uma senha"
            required
          />
        </div>
        <div class="form-group">
          <label for="role">Cargo / Nível de Acesso</label>
          <select id="role" v-model="role" required>
            <option value="Analista">Analista</option>
            <option value="Admin">Admin</option>
          </select>
        </div>
        <div v-if="error" class="error-message">
          {{ error }}
        </div>
        <div v-if="successMsg" class="success-message">
          {{ successMsg }}
        </div>
        <button type="submit" :disabled="loading">
          {{ loading ? 'Carregando...' : 'Cadastrar' }}
        </button>
      </form>
      <p class="auth-footer">
        Já tem conta? <router-link to="/login">Faça Login</router-link>
      </p>
    </div>
  </div>
</template>

<script setup>
import { ref } from 'vue';
import { useRouter } from 'vue-router';
import { useAuthStore } from '../stores/auth';

const email = ref('');
const password = ref('');
const role = ref('Analista');
const error = ref('');
const successMsg = ref('');
const loading = ref(false);

const auth = useAuthStore();
const router = useRouter();

const handleRegister = async () => {
  error.value = '';
  successMsg.value = '';
  loading.value = true;
  try {
    const success = await auth.register(email.value, password.value, role.value);
    if (success) {
      successMsg.value = 'Cadastro concluído com sucesso! Redirecionando para login...';
      setTimeout(() => {
        router.push('/login');
      }, 1500);
    }
  } catch (err) {
    error.value = err.response?.data?.error || 'Erro ao realizar cadastro. Tente outro e-mail.';
  } finally {
    loading.value = false;
  }
};
</script>

<style scoped>
.auth-container {
  display: flex;
  justify-content: center;
  align-items: center;
  min-height: 100vh;
  background: linear-gradient(135deg, #1e293b, #0f172a);
  font-family: 'Inter', sans-serif;
  color: #f8fafc;
}

.auth-card {
  background: rgba(30, 41, 59, 0.7);
  backdrop-filter: blur(10px);
  border: 1px solid rgba(255, 255, 255, 0.1);
  padding: 2.5rem;
  border-radius: 16px;
  width: 100%;
  max-width: 400px;
  box-shadow: 0 10px 25px -5px rgba(0, 0, 0, 0.3), 0 8px 10px -6px rgba(0, 0, 0, 0.3);
}

h2 {
  text-align: center;
  margin-bottom: 2rem;
  font-weight: 700;
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

input, select {
  width: 100%;
  padding: 0.75rem 1rem;
  border-radius: 8px;
  border: 1px solid #334155;
  background: #0f172a;
  color: #f8fafc;
  box-sizing: border-box;
  font-size: 1rem;
  transition: border-color 0.2s;
}

input:focus, select:focus {
  outline: none;
  border-color: #38bdf8;
}

button {
  width: 100%;
  padding: 0.75rem;
  border-radius: 8px;
  border: none;
  background: #38bdf8;
  color: #0f172a;
  font-size: 1rem;
  font-weight: 600;
  cursor: pointer;
  transition: background-color 0.2s;
  margin-top: 1rem;
}

button:hover {
  background: #0ea5e9;
}

button:disabled {
  background: #475569;
  cursor: not-allowed;
}

.error-message {
  background: rgba(239, 68, 68, 0.2);
  border: 1px solid #ef4444;
  color: #fca5a5;
  padding: 0.75rem;
  border-radius: 8px;
  font-size: 0.875rem;
  margin-bottom: 1rem;
  text-align: center;
}

.success-message {
  background: rgba(34, 197, 94, 0.2);
  border: 1px solid #22c55e;
  color: #86efac;
  padding: 0.75rem;
  border-radius: 8px;
  font-size: 0.875rem;
  margin-bottom: 1rem;
  text-align: center;
}

.auth-footer {
  text-align: center;
  margin-top: 1.5rem;
  font-size: 0.875rem;
  color: #94a3b8;
}

.auth-footer a {
  color: #38bdf8;
  text-decoration: none;
  font-weight: 500;
}

.auth-footer a:hover {
  text-decoration: underline;
}
</style>
