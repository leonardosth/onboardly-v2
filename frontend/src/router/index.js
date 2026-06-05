import { createRouter, createWebHistory } from 'vue-router';
import { useAuthStore } from '../stores/auth';

import Login from '../pages/Login.vue';
import Register from '../pages/Register.vue';
import MainLayout from '../layouts/MainLayout.vue';

// Lazy load secure views
const Dashboard = () => import('../pages/Dashboard.vue');
const ClientsList = () => import('../pages/ClientsList.vue');
const ClientDetail = () => import('../pages/ClientDetail.vue');
const ProjectDetail = () => import('../pages/ProjectDetail.vue');
const ProjectsList = () => import('../pages/ProjectsList.vue');
const UsersList = () => import('../pages/UsersList.vue');

const routes = [
  {
    path: '/login',
    name: 'Login',
    component: Login,
    meta: { guest: true }
  },
  {
    path: '/register',
    name: 'Register',
    component: Register,
    meta: { guest: true }
  },
  {
    path: '/',
    component: MainLayout,
    meta: { requiresAuth: true },
    children: [
      { path: '', redirect: '/dashboard' },
      { path: 'dashboard', name: 'Dashboard', component: Dashboard },
      { path: 'clients', name: 'ClientsList', component: ClientsList },
      { path: 'clients/:id', name: 'ClientDetail', component: ClientDetail },
      { path: 'projects', name: 'ProjectsList', component: ProjectsList },
      { path: 'projects/:id', name: 'ProjectDetail', component: ProjectDetail },
      { path: 'users', name: 'UsersList', component: UsersList, meta: { adminOnly: true } }
    ]
  },
  {
    path: '/:pathMatch(.*)*',
    redirect: '/dashboard'
  }
];

const router = createRouter({
  history: createWebHistory(),
  routes
});

// Guard routes to check authentication requirements
router.beforeEach((to, from, next) => {
  const authStore = useAuthStore();
  const authenticated = authStore.isAuthenticated;

  if (to.meta.requiresAuth && !authenticated) {
    next('/login');
  } else if (to.meta.adminOnly && !authStore.isAdmin) {
    next('/dashboard');
  } else if (to.meta.guest && authenticated) {
    next('/dashboard');
  } else {
    next();
  }
});

export default router;
