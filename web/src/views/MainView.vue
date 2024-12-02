<!-- src/views/MainView.vue -->
<template>
  <div>
    <header>
      <h1>Mock Trading Platform</h1>
      <nav>
        <RouterLink to="/">Home</RouterLink>
        <RouterLink to="/trades">Trades</RouterLink>
        <RouterLink to="/history">Trade-History</RouterLink>
        <RouterLink to="/positions">Positions</RouterLink>
        <RouterLink v-if="!isAuthenticated" to="/login">Login</RouterLink>
        <RouterLink v-if="!isAuthenticated" to="/register">Register</RouterLink>
        <button v-if="isAuthenticated" @click="logout">Logout</button>
      </nav>
    </header>
    
    <main>
      <RouterView />
    </main>
  </div>
</template>

<script setup>
import { RouterLink } from 'vue-router';
import { useAuthStore } from '../stores/auth';
import { computed } from 'vue';

const authStore = useAuthStore();

const isAuthenticated = computed(() => !!authStore.user);

const logout = () => {
  authStore.logout();
};
</script>

<style scoped>
header {
  display: flex;
  justify-content: space-between;
  padding: 1rem;
  background-color: #f8f9fa;
}

nav a {
  margin: 0 10px;
}
</style>