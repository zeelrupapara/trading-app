import { defineStore } from 'pinia';
import { loginUser, registerUser } from '../services/apiService';

export const useAuthStore = defineStore('auth', {
  state: () => ({
    user: null,
  }),
  actions: {
    async login(email, password) {
      const userData = await loginUser(email, password);
      this.user = userData;
      localStorage.setItem('jwt', userData.token); // Store your JWT
    },
    async register(userDetails) {
      const userData = await registerUser(userDetails);
      this.user = userData;
      localStorage.setItem('jwt', userData.token); // Store your JWT
    },
    logout() {
      this.user = null;
      localStorage.removeItem('jwt');
    },
  },
});