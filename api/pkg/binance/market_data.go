package binance_connector

import (
	"encoding/json"
	"log"
	"sync"

	binance "github.com/adshao/go-binance/v2"
)

type BinanceConnector struct {
	clients map[string]map[chan []byte]bool 
	mu      sync.Mutex
}

func NewBinanceConnector() *BinanceConnector {
	return &BinanceConnector{
		clients: make(map[string]map[chan []byte]bool),
	}
}

func (bc *BinanceConnector) Subscribe(symbol string, clientChan chan []byte) {
	bc.mu.Lock()
	defer bc.mu.Unlock()

	if _, exists := bc.clients[symbol]; !exists {
		bc.clients[symbol] = make(map[chan []byte]bool)
		go bc.startWebSocket(symbol)
	}
	bc.clients[symbol][clientChan] = true
}

func (bc *BinanceConnector) Unsubscribe(symbol string, clientChan chan []byte) {
	bc.mu.Lock()
	defer bc.mu.Unlock()

	if clients, exists := bc.clients[symbol]; exists {
		delete(clients, clientChan)
		if len(clients) == 0 {
			delete(bc.clients, symbol)
		}
		close(clientChan)
	}
}

func (bc *BinanceConnector) Broadcast(symbol string, message []byte) {
	bc.mu.Lock()
	defer bc.mu.Unlock()

	if clients, exists := bc.clients[symbol]; exists {
		for clientChan := range clients {
			select {
			case clientChan <- message:
			default:
				close(clientChan)
				delete(clients, clientChan)
			}
		}
	}
}

func (bc *BinanceConnector) startWebSocket(symbol string) {
	wsDepthHandler := func(event *binance.WsDepthEvent) {
		message := marshalEvent(event)
		bc.Broadcast(symbol, message)
	}

	errHandler := func(err error) {
		log.Println("WebSocket Error:", err)
	}

	doneC, _, err := binance.WsDepthServe(symbol, wsDepthHandler, errHandler)
	if err != nil {
		log.Fatalf("Failed to start WebSocket for symbol %s: %v", symbol, err)
	}
	<-doneC
}

func marshalEvent(event *binance.WsDepthEvent) []byte {
	depthData := struct {
		Symbol   string `json:"symbol"`
		BidPrice string `json:"bidPrice"`
		BidQty   string `json:"bidQty"`
		AskPrice string `json:"askPrice"`
		AskQty   string `json:"askQty"`
	}{
		Symbol:   event.Symbol,
		BidPrice: event.Bids[0].Price,
		BidQty:   event.Bids[0].Quantity,
		AskPrice: event.Asks[0].Price,
		AskQty:   event.Asks[0].Quantity,
	}

	// Marshal struct to JSON
	jsonData, err := json.Marshal(depthData)
	if err != nil {
		log.Println("Error marshaling depth event to JSON:", err)
		return nil
	}

	return jsonData
}
