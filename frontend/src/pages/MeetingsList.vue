<template>
  <div class="dashboard-layout">
    <div class="main-content">
      <div class="actions-row">
        <h1 class="page-title">Minhas Agendas</h1>
      </div>

      <!-- Filters -->
      <div class="filters-row">
        <select v-model="statusFilter" class="status-select" @change="onStatusFilterChange">
          <option value="completed">Concluídas</option>
          <option value="scheduled">Agendadas</option>
          <option value="">Todas</option>
        </select>
        <button class="new-meeting-btn" @click="showScheduleModal = true">📅 Nova Agenda</button>
      </div>

      <!-- List / Table -->
      <div v-if="meetingsStore.loading" class="loading-container">
        Carregando agendas...
      </div>
      <div v-else-if="meetingsStore.error" class="error-container">
        {{ meetingsStore.error }}
      </div>
      <div v-else-if="meetingsStore.myMeetings?.length === 0" class="empty-container">
        Nenhuma agenda encontrada.
      </div>
      <div v-else class="calendar-container">
        <FullCalendar :options="calendarOptions" />
      </div>

      <MeetingScheduler 
        v-if="showScheduleModal" 
        @close="showScheduleModal = false"
        @scheduled="onMeetingScheduled"
      />
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted, computed } from 'vue';
import { useRouter } from 'vue-router';
import { useMeetingsStore } from '../stores/meetings';
import MeetingScheduler from '../components/MeetingScheduler.vue';

import FullCalendar from '@fullcalendar/vue3';
import dayGridPlugin from '@fullcalendar/daygrid';
import timeGridPlugin from '@fullcalendar/timegrid';
import interactionPlugin from '@fullcalendar/interaction';
import ptBrLocale from '@fullcalendar/core/locales/pt-br';

const router = useRouter();
const meetingsStore = useMeetingsStore();
const statusFilter = ref('');
const showScheduleModal = ref(false);

onMounted(() => {
  fetchMyMeetings();
});

const fetchMyMeetings = () => {
  meetingsStore.fetchMyMeetings(statusFilter.value);
};

const onStatusFilterChange = () => {
  fetchMyMeetings();
};

const onMeetingScheduled = () => {
  fetchMyMeetings();
};

const handleEventClick = (info) => {
  const projectId = info.event.extendedProps.projectId;
  if (projectId) {
    router.push('/projects/' + projectId);
  }
};

const calendarOptions = computed(() => {
  return {
    plugins: [dayGridPlugin, timeGridPlugin, interactionPlugin],
    initialView: 'dayGridMonth',
    headerToolbar: {
      left: 'prev,next today',
      center: 'title',
      right: 'dayGridMonth,timeGridWeek'
    },
    locale: ptBrLocale,
    events: meetingsStore.myMeetings.map(m => ({
      id: m.id,
      title: `${m.client_name || 'N/A'} - ${m.title}`,
      start: m.scheduled_at,
      backgroundColor: m.status === 'completed' ? '#059669' : '#d97706',
      borderColor: m.status === 'completed' ? '#047857' : '#b45309',
      extendedProps: {
        projectId: m.project_id
      }
    })),
    eventClick: handleEventClick,
    height: 'auto',
    buttonText: {
      today: 'Hoje',
      month: 'Mês',
      week: 'Semana'
    }
  };
});
</script>



<style scoped>
.dashboard-layout {
  min-height: 100vh;
  background: #0f172a;
  color: #f8fafc;
  font-family: 'Inter', sans-serif;
}

.main-content {
  padding: 2.5rem;
  max-width: 1200px;
  margin: 0 auto;
}

.page-title {
  margin: 0;
  font-size: 1.5rem;
  color: #f8fafc;
}

.actions-row {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 2rem;
}

.filters-row {
  display: flex;
  gap: 1rem;
  margin-bottom: 2rem;
}

.status-select {
  padding: 0.75rem 1rem;
  border-radius: 8px;
  border: 1px solid #334155;
  background: rgba(30, 41, 59, 0.5);
  color: #f8fafc;
  font-size: 1rem;
  outline: none;
  width: 200px;
}

.status-select:focus {
  border-color: #38bdf8;
}

.new-meeting-btn {
  background: #38bdf8;
  color: #0f172a;
  border: none;
  padding: 0.75rem 1.25rem;
  border-radius: 8px;
  font-weight: 600;
  font-size: 0.875rem;
  cursor: pointer;
  transition: background-color 0.2s;
  display: flex;
  align-items: center;
  gap: 0.5rem;
}

.new-meeting-btn:hover {
  background: #0ea5e9;
}

.calendar-container {
  background: rgba(30, 41, 59, 0.5);
  border: 1px solid rgba(255, 255, 255, 0.05);
  border-radius: 12px;
  overflow: hidden;
  box-shadow: 0 4px 6px -1px rgba(0, 0, 0, 0.1);
  padding: 1.5rem;
}

/* FullCalendar customizations for dark theme */
:deep(.fc) {
  color: #f8fafc;
}

:deep(.fc-theme-standard td),
:deep(.fc-theme-standard th) {
  border-color: rgba(255, 255, 255, 0.05);
}

:deep(.fc-col-header-cell) {
  background: rgba(15, 23, 42, 0.4);
  color: #94a3b8;
  padding: 0.5rem 0;
  font-weight: 600;
  text-transform: uppercase;
  font-size: 0.875rem;
}

:deep(.fc-daygrid-day) {
  background: rgba(30, 41, 59, 0.2);
}

:deep(.fc-daygrid-day-number) {
  color: #cbd5e1;
  padding: 0.5rem;
}

:deep(.fc-button-primary) {
  background-color: #38bdf8 !important;
  border-color: #38bdf8 !important;
  color: #0f172a !important;
  font-weight: 500 !important;
  text-transform: capitalize !important;
}

:deep(.fc-button-primary:hover) {
  background-color: #0284c7 !important;
  border-color: #0284c7 !important;
}

:deep(.fc-button-primary:disabled) {
  background-color: #0f172a !important;
  border-color: #334155 !important;
  color: #94a3b8 !important;
}

:deep(.fc-button-active) {
  background-color: #0ea5e9 !important;
  border-color: #0ea5e9 !important;
}

:deep(.fc-event) {
  cursor: pointer;
  padding: 2px 4px;
  border-radius: 4px;
  font-size: 0.8rem;
  font-weight: 500;
  transition: opacity 0.2s;
}

:deep(.fc-event:hover) {
  opacity: 0.8;
}

.loading-container, .error-container, .empty-container {
  text-align: center;
  padding: 3rem;
  background: rgba(30, 41, 59, 0.3);
  border-radius: 12px;
  border: 1px dashed rgba(255, 255, 255, 0.1);
  color: #94a3b8;
}
</style>
