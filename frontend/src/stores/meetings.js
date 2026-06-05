import { defineStore } from 'pinia';
import axios from 'axios';

const API_URL = 'http://localhost:8080/api';

export const useMeetingsStore = defineStore('meetings', {
  state: () => ({
    meetings: [],
    loading: false,
    error: null,
  }),
  actions: {
    async fetchMeetings(projectId) {
      this.loading = true;
      this.error = null;
      try {
        const response = await axios.get(`${API_URL}/meetings?project_id=${projectId}`);
        this.meetings = response.data || [];
      } catch (err) {
        this.error = err.response?.data?.error || 'Erro ao carregar reuniões';
      } finally {
        this.loading = false;
      }
    },
    async scheduleMeeting(meetingData) {
      this.loading = true;
      this.error = null;
      try {
        const response = await axios.post(`${API_URL}/meetings`, meetingData);
        this.meetings.push(response.data);
        return response.data;
      } catch (err) {
        this.error = err.response?.data?.error || 'Erro ao agendar reunião';
        throw err;
      } finally {
        this.loading = false;
      }
    }
  }
});
