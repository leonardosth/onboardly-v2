import { defineStore } from 'pinia';
import axios from 'axios';

const API_URL = 'http://localhost:8080/api';

export const useAuthStore = defineStore('auth', {
  state: () => ({
    token: localStorage.getItem('token') || null,
    role: localStorage.getItem('role') || null,
    email: localStorage.getItem('email') || null,
  }),
  getters: {
    isAuthenticated: (state) => !!state.token,
    isAdmin: (state) => state.role === 'Admin',
    isAnalista: (state) => state.role === 'Analista',
  },
  actions: {
    async login(email, password) {
      try {
        const response = await axios.post(`${API_URL}/auth/login`, { email, password });
        const { token, role } = response.data;
        
        this.token = token;
        this.role = role;
        this.email = email;
        
        localStorage.setItem('token', token);
        localStorage.setItem('role', role);
        localStorage.setItem('email', email);
        
        // Setup axios default authorization header
        axios.defaults.headers.common['Authorization'] = `Bearer ${token}`;
        return true;
      } catch (error) {
        console.error('Login error:', error);
        throw error;
      }
    },
    async register(email, password, role) {
      try {
        await axios.post(`${API_URL}/auth/register`, { email, password, role });
        return true;
      } catch (error) {
        console.error('Registration error:', error);
        throw error;
      }
    },
    logout() {
      this.token = null;
      this.role = null;
      this.email = null;
      
      localStorage.removeItem('token');
      localStorage.removeItem('role');
      localStorage.removeItem('email');
      
      delete axios.defaults.headers.common['Authorization'];
    },
    initializeAuth() {
      if (this.token) {
        axios.defaults.headers.common['Authorization'] = `Bearer ${this.token}`;
      }
    }
  }
});
