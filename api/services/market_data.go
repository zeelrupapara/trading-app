package services

import (
	binance_connector "github.com/zeelrupapara/trading-api/pkg/binance"
)

type BinanceService struct {
	connector *binance_connector.BinanceConnector
}

// NewBinanceService initializes the service
func NewBinanceService() *BinanceService {
	connector := binance_connector.NewBinanceConnector()
	return &BinanceService{
		connector: connector,
	}
}

// RegisterClient registers a new client for a specific symbol
func (s *BinanceService) RegisterClient(symbol string, clientChan chan []byte) {
	s.connector.Subscribe(symbol, clientChan)
}

// UnregisterClient unregisters a client for a specific symbol
func (s *BinanceService) UnregisterClient(symbol string, clientChan chan []byte) {
	s.connector.Unsubscribe(symbol, clientChan)
}
