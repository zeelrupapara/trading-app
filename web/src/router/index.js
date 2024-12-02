// src/router/index.js
import { createRouter, createWebHistory } from 'vue-router';
import Home from '../pages/Home.vue';
import Login from '../pages/Login.vue';
import Register from '../pages/Register.vue';
import Trades from '../pages/Trades.vue';
import TradeHistory from '../pages/TradeHistory.vue';
import Positions from '../pages/Positions.vue';

const routes = [
  { path: '/', component: Home },
  { path: '/login', component: Login },
  { path: '/register', component: Register },
  { path: '/trades', component: Trades, meta: { requiresAuth: true } },
  { path: '/history', component: TradeHistory, meta: { requiresAuth: true } },
  { path: '/positions', component: Positions, meta: { requiresAuth: true } },
];

const router = createRouter({
  history: createWebHistory(),
  routes,
});

// Navigation guard for protecting routes
router.beforeEach((to, from, next) => {
  const requiresAuth = to.meta.requiresAuth;
  const isAuthenticated = !!localStorage.getItem('jwt');

  if (requiresAuth && !isAuthenticated) {
    next('/login'); // Redirect to login if not authenticated
  } else {
    next();
  }
});

export default router;