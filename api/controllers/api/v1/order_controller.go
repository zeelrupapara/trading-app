package v1

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/doug-martin/goqu/v9"
	"github.com/gofiber/fiber/v2"
	"github.com/zeelrupapara/trading-api/config"
	"github.com/zeelrupapara/trading-api/constants"
	"github.com/zeelrupapara/trading-api/models"
	binance_connector "github.com/zeelrupapara/trading-api/pkg/binance"
	"github.com/zeelrupapara/trading-api/services"
	"github.com/zeelrupapara/trading-api/structs"
	"github.com/zeelrupapara/trading-api/utils"
	"go.uber.org/zap"
)

// OrderController handles order-related requests
type OrderController struct {
	service *services.OrderService
	logger  *zap.Logger
}

// NewOrderController initializes the order service
func NewOrderController(logger *zap.Logger, db *goqu.Database, cfg config.AppConfig) (*OrderController, error) {
	// Initialize the binance service
	binance_client := binance_connector.NewBinanceClient(cfg.Binance.APIKey, cfg.Binance.APISecret)

	// Initialize the models
	orderModel, err := models.InitOrderModel(db)
	if err != nil {
		logger.Error("error while initializing order model", zap.Error(err))
		return nil, err
	}

	// Initialize the order service
	service := services.NewOrderService(binance_client, &orderModel)
	return &OrderController{service: service, logger: logger}, nil
}

// PlaceOrder handles the POST request to place a new order
// swagger:route POST /api/v1/orders Orders RequestPlaceOrder
//
// Place a new order.
//
//		Consumes:
//		- application/json
//
//		Schemes: http, https
//
//		Responses:
//		  200: ResponsePlaceOrder
//	      400: GenericResFailBadRequest
//		  500: GenericResError
func (ctrl *OrderController) PlaceOrder(c *fiber.Ctx) error {
	uid := c.Locals(constants.ContextUid).(string)

	var orderReq structs.ReqPlaceOrder

	err := json.Unmarshal(c.Body(), &orderReq)
	if err != nil {
		ctrl.logger.Error("error while unmarshal", zap.Error(err))
		return utils.JSONFail(c, http.StatusBadRequest, err.Error())
	}

	if orderReq.Type != "buy" && orderReq.Type != "sell" {
		ctrl.logger.Error("Order type must be 'buy' or 'sell'")
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Order type must be 'buy' or 'sell'"})
	}

	order, err := ctrl.service.PlaceOrder(orderReq.Symbol, orderReq.Volume, orderReq.Type, uid)
	if err != nil {
		log.Println("Error placing order:", err)
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to place order"})
	}

	ctrl.logger.Debug("Order placed", zap.Any("order", order))
	return c.JSON(order)
}

// GetAllOrders handles the GET request to retrieve all orders
// swagger:route GET /api/v1/trade-history Orders RequestTradeHistory
//
// Retrieve all orders.
//
//	Schemes: http, https
//
//	Responses:
//	  200: ResponseTradeHistory
//	  500: GenericResError
func (ctrl *OrderController) GetOrders(c *fiber.Ctx) error {
	uid := c.Locals(constants.ContextUid).(string)

	page, err := strconv.Atoi(c.Query("page", "1"))
	if err != nil {
		return utils.JSONFail(c, http.StatusBadRequest, "Invalid page")
	}

	// get the limit
	limit := constants.OrderPageLimit
	offset := (page - 1) * int(limit)

	orders, err := ctrl.service.GetOrders(uint(limit), uint(offset), uid)
	if err != nil {
		return utils.JSONError(c, http.StatusInternalServerError, err.Error())
	}

	ctrl.logger.Debug("Orders retrieved", zap.Any("orders", orders))
	return utils.JSONSuccess(c, http.StatusOK, orders)
}

// GetPositions handles the GET request to retrieve all positions
// swagger:route GET /api/v1/position Orders nil
//
// Retrieve all positions.
//
//	Schemes: http, https
//
//	Responses:
//	  200: ResponsePosition
//	  500: GenericResError
func (ctrl *OrderController) GetPositions(c *fiber.Ctx) error {
	uid := c.Locals(constants.ContextUid).(string)

	positions, err := ctrl.service.GetUserPositions(uid)
	ctrl.logger.Info("positions", zap.Any("positions", positions))
	if err != nil {
		ctrl.logger.Error("error while get user positions", zap.Any("id", uid), zap.Error(err))
		return utils.JSONError(c, http.StatusInternalServerError, err.Error())
	}
	return utils.JSONSuccess(c, http.StatusOK, positions)
}
