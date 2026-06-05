import axios from 'axios';

const API_URL = 'http://localhost:8080/api';

export const getClients = async () => {
  const response = await axios.get(`${API_URL}/clients`);
  return response.data;
};

export const getClient = async (id) => {
  const response = await axios.get(`${API_URL}/clients/${id}`);
  return response.data;
};

export const createClient = async (clientData) => {
  const response = await axios.post(`${API_URL}/clients`, clientData);
  return response.data;
};

export const updateClient = async (id, clientData) => {
  const response = await axios.put(`${API_URL}/clients/${id}`, clientData);
  return response.data;
};

export const deleteClient = async (id) => {
  const response = await axios.delete(`${API_URL}/clients/${id}`);
  return response.data;
};
