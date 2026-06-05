import { defineStore } from 'pinia';
import * as projectService from '../services/project';

export const useProjectsStore = defineStore('projects', {
  state: () => ({
    projects: [],
    currentProject: null,
    loading: false,
    error: null,
  }),
  actions: {
    async fetchProjects() {
      this.loading = true;
      this.error = null;
      try {
        this.projects = (await projectService.getProjects()) || [];
      } catch (err) {
        this.error = err.response?.data?.error || 'Erro ao carregar projetos';
      } finally {
        this.loading = false;
      }
    },
    async fetchProject(id) {
      this.loading = true;
      this.error = null;
      try {
        this.currentProject = await projectService.getProject(id);
      } catch (err) {
        this.error = err.response?.data?.error || 'Erro ao obter detalhes do projeto';
      } finally {
        this.loading = false;
      }
    },
    async addProject(projectData) {
      this.loading = true;
      this.error = null;
      try {
        const newProject = await projectService.createProject(projectData);
        this.projects.push(newProject);
        return newProject;
      } catch (err) {
        this.error = err.response?.data?.error || 'Erro ao cadastrar projeto';
        throw err;
      } finally {
        this.loading = false;
      }
    },
    async editProject(id, projectData) {
      this.loading = true;
      this.error = null;
      try {
        const updated = await projectService.updateProject(id, projectData);
        const index = this.projects.findIndex(p => p.id === id);
        if (index !== -1) {
          this.projects[index] = updated;
        }
        if (this.currentProject && this.currentProject.id === id) {
          this.currentProject = updated;
        }
        return updated;
      } catch (err) {
        this.error = err.response?.data?.error || 'Erro ao editar projeto';
        throw err;
      } finally {
        this.loading = false;
      }
    },
    async removeProject(id) {
      this.loading = true;
      this.error = null;
      try {
        await projectService.deleteProject(id);
        this.projects = this.projects.filter(p => p.id !== id);
        if (this.currentProject && this.currentProject.id === id) {
          this.currentProject = null;
        }
      } catch (err) {
        this.error = err.response?.data?.error || 'Erro ao excluir projeto';
        throw err;
      } finally {
        this.loading = false;
      }
    }
  }
});
