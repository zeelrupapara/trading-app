package services

import (
	"context"
	"fmt"
	"time"

	"github.com/zeelrupapara/trading-api/models"
	binance_connector "github.com/zeelrupapara/trading-api/pkg/binance"
	"github.com/zeelrupapara/trading-api/utils"
)

type OrderService struct {
	binanceClient *binance_connector.BinanceClient
	orderModel    *models.OrderModel
}

func NewOrderService(client *binance_connector.BinanceClient, orderModel *models.OrderModel) *OrderService {
	return &OrderService{
		binanceClient: client,
		orderModel:    orderModel,
	}
}

// PlaceOrder handles placing a new order
func (s *OrderService) PlaceOrder(symbol string, volume float64, orderType string, uid string) (models.Order, error) {
	ticker, err := s.binanceClient.Client.NewListPricesService().Symbol(symbol).Do(context.Background())
	if err != nil {
		return models.Order{}, err
	}

	currentPrice := "0.0"
	for _, price := range ticker {
		if price.Symbol == symbol {
			currentPrice = price.Price
			break
		}
	}

	if currentPrice == "0.0" {
		return models.Order{}, fmt.Errorf("could not find price for symbol: %s", symbol)
	}

	orderID := utils.GenerateUUID("order")
	order := models.Order{
		ID:        orderID,
		Symbol:    symbol,
		Volume:    volume,
		OrderType: orderType,
		Price:     currentPrice,
		UserId:    uid,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	order, err = s.orderModel.InsertOrder(order)
	if err != nil {
		return models.Order{}, err
	}

	return order, nil
}
