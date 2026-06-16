import { defineStore } from 'pinia';
import * as meetingService from '../services/meeting';

export const useMeetingsStore = defineStore('meetings', {
  state: () => ({
    meetings: [],
    myMeetings: [],
    loading: false,
    error: null,
  }),
  actions: {
    async fetchMeetings(projectId) {
      this.loading = true;
      this.error = null;
      try {
        this.meetings = await meetingService.getMeetings(projectId) || [];
      } catch (err) {
        this.error = err.response?.data?.error || 'Erro ao carregar reuniões';
      } finally {
        this.loading = false;
      }
    },
    async fetchMyMeetings(status) {
      this.loading = true;
      this.error = null;
      try {
        this.myMeetings = await meetingService.getMyMeetings(status) || [];
      } catch (err) {
        this.error = err.response?.data?.error || 'Erro ao carregar minhas agendas';
      } finally {
        this.loading = false;
      }
    },
    async scheduleMeeting(meetingData) {
      this.loading = true;
      this.error = null;
      try {
        const newMeeting = await meetingService.createMeeting(meetingData);
        this.meetings.push(newMeeting);
        return newMeeting;
      } catch (err) {
        this.error = err.response?.data?.error || 'Erro ao agendar reunião';
        throw err;
      } finally {
        this.loading = false;
      }
    },
    async completeMeeting(meetingId, activateClient) {
      this.loading = true;
      this.error = null;
      try {
        const result = await meetingService.completeMeeting(meetingId, activateClient);
        
        // Update meetings array if the completed meeting is in it
        const idx = this.meetings.findIndex(m => m.id === meetingId);
        if (idx !== -1) {
          this.meetings[idx] = result.meeting;
        }

        // Update myMeetings array if the completed meeting is in it
        const myIdx = this.myMeetings.findIndex(m => m.id === meetingId);
        if (myIdx !== -1) {
          this.myMeetings[myIdx] = {
            ...this.myMeetings[myIdx],
            ...result.meeting
          };
        }
        
        return result;
      } catch (err) {
        this.error = err.response?.data?.error || 'Erro ao concluir reunião';
        throw err;
      } finally {
        this.loading = false;
      }
    }
  }
});
