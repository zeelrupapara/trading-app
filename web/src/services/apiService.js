import axios from 'axios';

const apiClient = axios.create({
  baseURL: 'http://localhost:8000/api/v1',
  withCredentials: true,
  headers: {
    'Content-Type': 'application/json',
  },
});

// User Authentication
export const loginUser = async (email, password) => {
  const response = await apiClient.post('/login', { email, password }, { withCredentials: true });
  return response.data;
};

export const registerUser = async (userDetails) => {
  const response = await apiClient.post('/users', userDetails);
  return response.data; // change according to your response structure
};

// Fetch Trade History
export const getTradeHistory = async () => {
  const response = await apiClient.get('/trade-history');
  return response.data; // assuming response includes trade history
};

// Fetch Positions
export const getPositions = async () => {
  const response = await apiClient.get('/position');
  return response.data; // assuming response includes positions
};

// Place Orders
export const placeOrder = async (orderDetails) => {
  const response = await apiClient.post('/orders', orderDetails);
  return response.data; // assuming response includes order confirmation
};