<template>
  <div class="container">
    <h2>Welcome to Mock Trading Platform</h2>
    <p>Explore live trades, market data, and manage your portfolio efficiently.</p>

    <div class="row">
      <div class="col-md-4" v-for="symbol in symbols" :key="symbol.name">
        <div class="card mb-3" @click="navigateToTrades(symbol.symbol)" style="cursor: pointer;">
          <div class="card-body">
            <h5 class="card-title">{{ symbol.name }}</h5>
            <p class="card-text">{{ symbol.symbol }}</p>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref } from 'vue';
import { useRouter } from 'vue-router';

const symbols = ref([
  { name: 'BitcoinEth', symbol: 'ETHBTC' },
  { name: 'Bitcoin', symbol: 'BTCUSD' },
  { name: 'Ethereum', symbol: 'ETHUSD' },
  { name: 'Litecoin', symbol: 'LTCUSD' },
]);

const router = useRouter();

const navigateToTrades = (symbol) => {
  router.push({ path: '/trades', query: { symbol: symbol } });
};
</script>

<style scoped>
.card {
  transition: transform 0.2s;
}

.card:hover {
  transform: scale(1.02);
}
</style>