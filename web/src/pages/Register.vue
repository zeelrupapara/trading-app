<template>
    <div class="container">
      <h2>Register</h2>
      <form @submit.prevent="handleRegister">
        <div class="mb-3">
          <label for="email" class="form-label">Email</label>
          <input type="email" v-model="email" class="form-control" required />
        </div>
        <div class="mb-3">
          <label for="password" class="form-label">Password</label>
          <input type="password" v-model="password" class="form-control" required />
        </div>
        <div class="mb-3">
          <label for="firstName" class="form-label">First Name</label>
          <input type="text" v-model="firstName" class="form-control" required />
        </div>
        <div class="mb-3">
          <label for="lastName" class="form-label">Last Name</label>
          <input type="text" v-model="lastName" class="form-control" required />
        </div>
        <button type="submit" class="btn btn-primary">Register</button>
      </form>
    </div>
  </template>
  
  <script setup>
  import { ref } from 'vue';
  import { useAuthStore } from '../stores/auth';
  
  const email = ref('');
  const password = ref('');
  const firstName = ref('');
  const lastName = ref('');
  const authStore = useAuthStore();
  
  const handleRegister = async () => {
    await authStore.register({ email: email.value, password: password.value, first_name: firstName.value, last_name: lastName.value, roles: "user" });
    window.location.href = '/';
  };
  </script>