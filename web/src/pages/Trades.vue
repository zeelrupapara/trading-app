<template>
  <div class="container">
    <h2>{{ marketData.symbol }} - Trade Market</h2>
    <div class="card mb-3">
      <div class="card-body">
        <h5 class="card-title">Real-Time Data</h5>
        <p class="card-text">Bid Price: {{ marketData.bidPrice }}</p>
        <p class="card-text">Ask Price: {{ marketData.askPrice }}</p>
        <p class="card-text">Bid Quantity: {{ marketData.bidQty }}</p>
        <p class="card-text">Ask Quantity: {{ marketData.askQty }}</p>

        <input type="number" v-model="volume" placeholder="Enter volume" min="0" />

        <button class="btn btn-success me-2" @click="placeOrder('buy')">Buy</button>
        <button class="btn btn-danger" @click="placeOrder('sell')">Sell</button>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted, onUnmounted } from 'vue';
import { useRoute } from 'vue-router';
import { useMarketDataWebSocket } from '../services/webSocketService';
import { placeOrder as placeOrderAPI } from '../services/apiService';

const route = useRoute(); // To access the route object
const marketData = ref({
  symbol: route.query.symbol,
  bidPrice: '0.00',
  askPrice: '0.00',
  bidQty: '0.00',
  askQty: '0.00',
});

const volume = ref(0); // User input for volume
let socket;

const initializeWebSocket = () => {
  socket = useMarketDataWebSocket(marketData.value.symbol, (data) => {
    marketData.value = data; // Update market data with real-time data
  });
};

const placeOrder = async (type) => {
  if (volume.value <= 0) {
    alert("Volume must be greater than 0.");
    return;
  }

  const orderDetails = {
    symbol: marketData.value.symbol,
    volume: volume.value,
    type,
  };

  try {
    const response = await placeOrderAPI(orderDetails);
    alert(`Order ${type} placed successfully!`);
    console.log(response); // You might want to handle further actions or state changes here
  } catch (error) {
    console.error('Error placing the order:', error);
  }
};

onMounted(() => {
  initializeWebSocket();
});

onUnmounted(() => {
  if (socket) {
    socket.close();
  }
});
</script>

<style scoped>
.container {
  max-width: 800px;
  margin: 0 auto;
}

.card {
  box-shadow: 0 0 10px rgba(0, 0, 0, 0.1);
  margin-top: 2rem;
}
</style>