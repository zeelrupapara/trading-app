package services

import (
	"context"
	"fmt"
	"log"
	"strconv"
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

	currentPrice := 0.0
	for _, price := range ticker {
		if price.Symbol == symbol {
			currentPrice, _ = strconv.ParseFloat(price.Price, 64)
			break
		}
	}

	if currentPrice == 0.0 {
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

func (s *OrderService) GetOrders(limit uint, offset uint, userId string) ([]models.Order, error) {
	return s.orderModel.GetOrders(limit, offset, userId)
}

func (s *OrderService) GetUserPositions(userID string) ([]models.PositionSummary, error) {
	var summaries []models.PositionSummary

	// Retrieve all orders for the user to process
	orders, err := s.orderModel.GetAllOrders(userID)
	if err != nil {
		return nil, err
	}

	// Map to track the positions
	positionMap := make(map[string]*models.PositionSummary)

	// Iterate through the orders in the order of creation
	for _, order := range orders {
		volume := float64(order.Volume)

		// Initialize PositionSummary if it doesn't exist
		if _, exists := positionMap[order.Symbol]; !exists {
			positionMap[order.Symbol] = &models.PositionSummary{
				Symbol:            order.Symbol,
				HoldingVolume:     0,
				HoldingInvestment: 0,
				ProfitLoss:        0,
			}
		}

		if order.OrderType == "buy" {
			// Increase holding volume and investment for buys
			positionMap[order.Symbol].HoldingVolume += volume
			positionMap[order.Symbol].HoldingInvestment += volume * order.Price
		} else if order.OrderType == "sell" {
			// Decrease holding volume for sells, tracking FIFO
			if positionMap[order.Symbol].HoldingVolume < volume {
				// Handle overselling scenario if needed
				continue // Or trigger an error as appropriate
			}
			positionMap[order.Symbol].HoldingVolume -= volume
		}
	}

	// Calculate current price and profit/loss for each position
	for _, position := range positionMap {
		// Fetch current price for the position's symbol
		priceData, err := s.binanceClient.Client.NewListPricesService().Symbol(position.Symbol).Do(context.Background())
		if err != nil {
			return nil, fmt.Errorf("could not get current price for symbol %s: %v", position.Symbol, err)
		}

		for _, price := range priceData {
			if price.Symbol == position.Symbol {
				currentPrice, err := strconv.ParseFloat(price.Price, 64)
				if err != nil {
					log.Printf("Error parsing current price for symbol %s: %v", position.Symbol, err)
					continue
				}
				position.CurrentPrice = currentPrice
				break
			}
		}

		// Calculate profit and loss
		position.ProfitLoss = (position.CurrentPrice * position.HoldingVolume) - position.HoldingInvestment

		// Add position to summaries
		summaries = append(summaries, *position)
	}

	return summaries, nil
}
