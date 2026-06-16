import axios from 'axios';

const API_URL = 'http://localhost:8080/api';

export const getMeetings = async (projectId) => {
  const response = await axios.get(`${API_URL}/meetings?project_id=${projectId}`);
  return response.data;
};

export const getMyMeetings = async (status) => {
  const response = await axios.get(`${API_URL}/meetings/mine?status=${status || ''}`);
  return response.data;
};

export const createMeeting = async (meetingData) => {
  const response = await axios.post(`${API_URL}/meetings`, meetingData);
  return response.data;
};

export const completeMeeting = async (meetingId, activateClient) => {
  const response = await axios.post(`${API_URL}/meetings/${meetingId}/complete`, {
    activate_client: activateClient,
  });
  return response.data;
};
