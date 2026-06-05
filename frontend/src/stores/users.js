import { defineStore } from 'pinia';
import * as userService from '../services/user';

export const useUsersStore = defineStore('users', {
  state: () => ({
    users: [],
    loading: false,
    error: null,
  }),
  actions: {
    async fetchUsers() {
      this.loading = true;
      this.error = null;
      try {
        this.users = (await userService.getUsers()) || [];
      } catch (err) {
        this.error = err.response?.data?.error || 'Erro ao carregar usuários';
      } finally {
        this.loading = false;
      }
    },
    async addUser(userData) {
      this.loading = true;
      this.error = null;
      try {
        const newUser = await userService.createUser(userData);
        this.users.unshift(newUser); // add at beginning since ordered by created_at DESC
        return newUser;
      } catch (err) {
        this.error = err.response?.data?.error || 'Erro ao cadastrar usuário';
        throw err;
      } finally {
        this.loading = false;
      }
    },
    async removeUser(id) {
      this.loading = true;
      this.error = null;
      try {
        await userService.deleteUser(id);
        this.users = this.users.filter(u => u.id !== id);
      } catch (err) {
        this.error = err.response?.data?.error || 'Erro ao excluir usuário';
        throw err;
      } finally {
        this.loading = false;
      }
    }
  }
});
