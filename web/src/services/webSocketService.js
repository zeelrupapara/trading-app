export const useMarketDataWebSocket = (symbol, callback) => {
    const socket = new WebSocket(`ws://localhost:8000/ws/v1/marketdata?symbol=${symbol}`);
    console.log(socket.url);
    
  
    socket.onmessage = (event) => {
        const data = JSON.parse(event.data);
        console.log(event);
        
        callback({
            symbol: data.symbol,
            bidPrice: data.bidPrice,
            bidQty: data.bidQty,
            askPrice: data.askPrice,
            askQty: data.askQty
        });
    };
  
    return socket;
  };