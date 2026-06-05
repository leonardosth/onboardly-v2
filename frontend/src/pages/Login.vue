<template>
  <div class="auth-container">
    <div class="auth-card">
      <h2>🔐 Onboardly Login</h2>
      <form @submit.prevent="handleLogin">
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
            placeholder="Digite sua senha"
            required
          />
        </div>
        <div v-if="error" class="error-message">
          {{ error }}
        </div>
        <button type="submit" :disabled="loading">
          {{ loading ? 'Carregando...' : 'Entrar' }}
        </button>
      </form>
      <p class="auth-footer">
        Não tem conta? <router-link to="/register">Registre-se</router-link>
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
const error = ref('');
const loading = ref(false);

const auth = useAuthStore();
const router = useRouter();

const handleLogin = async () => {
  error.value = '';
  loading.value = true;
  try {
    const success = await auth.login(email.value, password.value);
    if (success) {
      router.push('/dashboard');
    }
  } catch (err) {
    error.value = err.response?.data?.error || 'Erro ao realizar login. Verifique suas credenciais.';
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

input {
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

input:focus {
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
