<template>
  <div class="app-container">
    <!-- Sidebar -->
    <aside class="sidebar">
      <div class="brand">
        <span class="logo-icon">🚀</span>
        <h2>Onboardly</h2>
      </div>

      <nav class="nav-links">
        <router-link to="/dashboard" class="nav-item">
          <span class="nav-icon">📊</span>
          Dashboard
        </router-link>
        <router-link to="/clients" class="nav-item" :class="{ active: $route.path.startsWith('/clients') }">
          <span class="nav-icon">🏢</span>
          Clientes
        </router-link>
        <router-link to="/projects" class="nav-item" :class="{ active: $route.path.startsWith('/projects') }">
          <span class="nav-icon">📁</span>
          Projetos
        </router-link>
        <router-link to="/meetings" class="nav-item" :class="{ active: $route.path.startsWith('/meetings') }">
          <span class="nav-icon">📅</span>
          Agendas
        </router-link>
        <router-link v-if="auth.isAdmin" to="/users" class="nav-item" :class="{ active: $route.path.startsWith('/users') }">
          <span class="nav-icon">👥</span>
          Usuários
        </router-link>
      </nav>

      <div class="sidebar-footer">
        <div class="user-meta">
          <div class="user-email">{{ auth.email }}</div>
          <span class="role-badge" :class="auth.role?.toLowerCase() || ''">{{ auth.role }}</span>
        </div>
        <button class="logout-btn" @click="handleLogout">Sair</button>
      </div>
    </aside>

    <!-- Main Content Area -->
    <main class="main-content-wrapper">
      <router-view />
    </main>
  </div>
</template>

<script setup>
import { useAuthStore } from '../stores/auth';
import { useRouter } from 'vue-router';

const auth = useAuthStore();
const router = useRouter();

const handleLogout = () => {
  auth.logout();
  router.push('/login');
};
</script>

<style scoped>
.app-container {
  display: flex;
  min-height: 100vh;
  background: #0f172a;
  color: #f8fafc;
  font-family: 'Inter', sans-serif;
}

.sidebar {
  width: 280px;
  background: rgba(30, 41, 59, 0.95);
  border-right: 1px solid rgba(255, 255, 255, 0.05);
  display: flex;
  flex-direction: column;
  position: fixed;
  height: 100vh;
  z-index: 50;
  backdrop-filter: blur(10px);
}

.brand {
  padding: 2rem 1.5rem;
  display: flex;
  align-items: center;
  gap: 0.75rem;
  border-bottom: 1px solid rgba(255, 255, 255, 0.05);
}

.brand h2 {
  margin: 0;
  font-size: 1.5rem;
  color: #38bdf8;
  font-weight: 700;
}

.logo-icon {
  font-size: 1.5rem;
}

.nav-links {
  padding: 1.5rem 1rem;
  flex: 1;
  display: flex;
  flex-direction: column;
  gap: 0.5rem;
}

.nav-item {
  display: flex;
  align-items: center;
  gap: 1rem;
  padding: 0.85rem 1rem;
  border-radius: 8px;
  color: #94a3b8;
  text-decoration: none;
  font-weight: 500;
  transition: all 0.2s;
}

.nav-item:hover {
  background: rgba(255, 255, 255, 0.05);
  color: #f8fafc;
}

.nav-item.router-link-active,
.nav-item.active {
  background: rgba(56, 189, 248, 0.15);
  color: #38bdf8;
}

.nav-icon {
  font-size: 1.25rem;
}

.sidebar-footer {
  padding: 1.5rem;
  border-top: 1px solid rgba(255, 255, 255, 0.05);
}

.user-meta {
  margin-bottom: 1rem;
}

.user-email {
  font-size: 0.9rem;
  font-weight: 600;
  margin-bottom: 0.25rem;
  word-break: break-all;
}

.role-badge {
  display: inline-block;
  font-size: 0.75rem;
  padding: 0.2rem 0.5rem;
  border-radius: 9999px;
  background: #334155;
}

.role-badge.admin {
  background: rgba(56, 189, 248, 0.2);
  color: #38bdf8;
}

.role-badge.analista {
  background: rgba(148, 163, 184, 0.2);
  color: #cbd5e1;
}

.logout-btn {
  width: 100%;
  background: transparent;
  border: 1px solid #ef4444;
  color: #ef4444;
  padding: 0.6rem;
  border-radius: 6px;
  cursor: pointer;
  font-weight: 600;
  transition: background 0.2s;
}

.logout-btn:hover {
  background: rgba(239, 68, 68, 0.1);
}

.main-content-wrapper {
  flex: 1;
  margin-left: 280px; /* offset for fixed sidebar */
  min-height: 100vh;
  overflow-x: hidden;
}
</style>
