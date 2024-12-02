<template>
  <div class="container">
    <h2>Your Positions</h2>
    <table class="table">
      <thead>
        <tr>
          <th>Symbol</th>
          <th>Holding Volume</th>
          <th>Holding Investment</th>
          <th>Current Price</th>
          <th>Profit/Loss</th>
        </tr>
      </thead>
      <tbody>
        <tr v-for="position in positions" :key="position.symbol">
          <td>{{ position.symbol }}</td>
          <td>{{ position.holding_volume }}</td>
          <td>{{ position.holding_investment }}</td>
          <td>{{ position.current_price }}</td>
          <td>{{ position.profit_loss }}</td>
        </tr>
      </tbody>
    </table>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue';
import { getPositions } from '../services/apiService';

const positions = ref([]);

const fetchPositions = async () => {
  try {
    const response = await getPositions();
    positions.value = response;
  } catch (error) {
    console.error('Error fetching positions:', error);
  }
};

onMounted(() => {
  fetchPositions();
});
</script>