import { defineStore } from 'pinia';
import * as clientService from '../services/client';

export const useClientsStore = defineStore('clients', {
  state: () => ({
    clients: [],
    currentClient: null,
    loading: false,
    error: null,
  }),
  actions: {
    async fetchClients() {
      this.loading = true;
      this.error = null;
      try {
        this.clients = (await clientService.getClients()) || [];
      } catch (err) {
        this.error = err.response?.data?.error || 'Erro ao carregar clientes';
      } finally {
        this.loading = false;
      }
    },
    async fetchClient(id) {
      this.loading = true;
      this.error = null;
      try {
        this.currentClient = await clientService.getClient(id);
      } catch (err) {
        this.error = err.response?.data?.error || 'Erro ao obter detalhes do cliente';
      } finally {
        this.loading = false;
      }
    },
    async addClient(clientData) {
      this.loading = true;
      this.error = null;
      try {
        const newClient = await clientService.createClient(clientData);
        this.clients.push(newClient);
        return newClient;
      } catch (err) {
        this.error = err.response?.data?.error || 'Erro ao cadastrar cliente';
        throw err;
      } finally {
        this.loading = false;
      }
    },
    async editClient(id, clientData) {
      this.loading = true;
      this.error = null;
      try {
        const updated = await clientService.updateClient(id, clientData);
        const index = this.clients.findIndex(c => c.id === id);
        if (index !== -1) {
          this.clients[index] = updated;
        }
        if (this.currentClient && this.currentClient.id === id) {
          this.currentClient = updated;
        }
        return updated;
      } catch (err) {
        this.error = err.response?.data?.error || 'Erro ao editar cliente';
        throw err;
      } finally {
        this.loading = false;
      }
    },
    async removeClient(id) {
      this.loading = true;
      this.error = null;
      try {
        await clientService.deleteClient(id);
        this.clients = this.clients.filter(c => c.id !== id);
        if (this.currentClient && this.currentClient.id === id) {
          this.currentClient = null;
        }
      } catch (err) {
        this.error = err.response?.data?.error || 'Erro ao excluir cliente';
        throw err;
      } finally {
        this.loading = false;
      }
    }
  }
});
