<template>
  <div class="container">
    <h2>Trade History</h2>
    <table class="table">
      <thead>
        <tr>
          <th>Order ID</th>
          <th>Symbol</th>
          <th>Type</th>
          <th>Volume</th>
          <th>Price</th>
          <th>Date</th>
        </tr>
      </thead>
      <tbody>
        <tr v-if="tradeHistory.length === 0">
          <td colspan="6" class="text-center">No trade history available</td>
        </tr>
        <tr v-for="trade in tradeHistory" :key="trade.id">
          <td>{{ trade.id }}</td>
          <td>{{ trade.symbol }}</td>
          <td>{{ trade.order_type }}</td>
          <td>{{ trade.volume }}</td>
          <td>{{ trade.price }}</td>
          <td>{{ trade.created_at }}</td>
        </tr>
      </tbody>
    </table>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue';
import { getTradeHistory } from '../services/apiService';

const tradeHistory = ref([]);

const fetchTradeHistory = async () => {
  try {
    const response = await getTradeHistory();
    tradeHistory.value = response; // Ensure you access the correct structure
  } catch (error) {
    console.error('Error fetching trade history:', error.response?.data || error.message);
  }
};

onMounted(() => {
  fetchTradeHistory();
});
</script>