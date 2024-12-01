package v1

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2" // Importing the websocket package
	"github.com/zeelrupapara/trading-api/services"
	"go.uber.org/zap"
)

type MarketDataController struct {
	service *services.BinanceService
	logger  *zap.Logger
}

func NewMarketDataController(logger *zap.Logger) *MarketDataController {
	service := services.NewBinanceService()
	return &MarketDataController{service: service, logger: logger}
}
// ServeMarketData handles the WebSocket connection for market data
// swagger:route GET /ws/v1/marketdata MarketData RequestMarketData
//
// Handles the WebSocket connection for market data.
//
//  Consumes:
//		- application/json
//
//  Schemes: ws, wss
//
//  Responses:
//	  200: ResponseMarketData
//	  400: GenericResFailBadRequest
//	  500: GenericResError

func (ctrl *MarketDataController) ServeMarketData() fiber.Handler {
	ctrl.logger.Info("Starting market data controller")
	return websocket.New(func(conn *websocket.Conn) {
		// Getting the symbol from the query
		symbol := conn.Query("symbol")
		if symbol == "" {
			ctrl.logger.Error("Symbol query parameter is required.")
			conn.WriteMessage(websocket.TextMessage, []byte("Symbol query parameter is required."))
			conn.Close()
			return
		}

		// Create a channel for messages to be sent to the client
		clientChannel := make(chan []byte)
		ctrl.service.RegisterClient(symbol, clientChannel)
		ctrl.logger.Info("Connected to market data", zap.String("symbol", symbol))

		defer func() {
			ctrl.service.UnregisterClient(symbol, clientChannel)
			ctrl.logger.Info("Disconnected from market data", zap.String("symbol", symbol))
			conn.Close()
		}()

		// Goroutine to handle incoming messages from the service
		go func() {
			for message := range clientChannel {
				if err := conn.WriteMessage(websocket.TextMessage, message); err != nil {
					ctrl.logger.Error("Error writing message:", zap.Error(err))
					break
				}
			}
		}()

		// Keep the connection open by reading incoming messages
		for {
			if _, _, err := conn.ReadMessage(); err != nil {
				ctrl.logger.Error("Error reading message:", zap.Error(err))
				break
			}
		}
	})
}
