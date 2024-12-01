package routes

import (
	"fmt"
	"sync"

	"github.com/zeelrupapara/trading-api/config"
	"github.com/zeelrupapara/trading-api/constants"
	"github.com/zeelrupapara/trading-api/middlewares"

	"github.com/doug-martin/goqu/v9"
	"github.com/gofiber/contrib/swagger"
	"github.com/gofiber/fiber/v2"
	controller "github.com/zeelrupapara/trading-api/controllers/api/v1"
	pMetrics "github.com/zeelrupapara/trading-api/pkg/prometheus"
	"go.uber.org/zap"
)

var mu sync.Mutex

// Setup func
func Setup(app *fiber.App, goqu *goqu.Database, logger *zap.Logger, config config.AppConfig, pMetrics *pMetrics.PrometheusMetrics) error {
	mu.Lock()

	app.Use(middlewares.LogHandler(logger, pMetrics))

	app.Use(swagger.New(swagger.Config{
		BasePath: "/api/v1/",
		FilePath: "./assets/swagger.json",
		Path:     "docs",
		Title:    "Swagger API Docs",
	}))

	// For API group
	api := app.Group("/api")
	api_v1 := api.Group("/v1")

	// For WS group
	ws := app.Group("/ws")
	ws_v1 := ws.Group("/v1")

	middlewares := middlewares.NewMiddleware(config, logger)

	err := setupAuthController(api_v1, goqu, logger, config)
	if err != nil {
		return err
	}

	err = setupUserController(api_v1, goqu, logger, middlewares)
	if err != nil {
		return err
	}

	err = healthCheckController(app, goqu, logger)
	if err != nil {
		return err
	}

	err = metricsController(app, goqu, logger, pMetrics)
	if err != nil {
		return err
	}

	err = marketDataController(ws_v1, logger)
	if err != nil {
		return err
	}

	err = ordersController(api_v1, goqu, logger, config, middlewares)
	if err != nil {
		return err
	}

	mu.Unlock()
	return nil
}

func setupAuthController(v1 fiber.Router, goqu *goqu.Database, logger *zap.Logger, config config.AppConfig) error {
	authController, err := controller.NewAuthController(goqu, logger, config)
	if err != nil {
		return err
	}
	v1.Post("/login", authController.Login)
	v1.Post("/logout", authController.Logout)
	return nil
}

func setupUserController(v1 fiber.Router, goqu *goqu.Database, logger *zap.Logger, middlewares middlewares.Middleware) error {
	userController, err := controller.NewUserController(goqu, logger)
	if err != nil {
		return err
	}

	userRouter := v1.Group("/users")
	userRouter.Post("/", userController.CreateUser)
	userRouter.Get(fmt.Sprintf("/:%s", constants.ParamUid), middlewares.Authenticated, userController.GetUser)
	return nil
}

func healthCheckController(app *fiber.App, goqu *goqu.Database, logger *zap.Logger) error {
	healthController, err := controller.NewHealthController(goqu, logger)
	if err != nil {
		return err
	}

	healthz := app.Group("/healthz")
	healthz.Get("/", healthController.Overall)
	healthz.Get("/db", healthController.Db)
	return nil
}

func metricsController(app *fiber.App, db *goqu.Database, logger *zap.Logger, pMetrics *pMetrics.PrometheusMetrics) error {
	metricsController, err := controller.InitMetricsController(db, logger, pMetrics)
	if err != nil {
		return nil
	}

	app.Get("/metrics", metricsController.Metrics)
	return nil
}

func marketDataController(ws fiber.Router, logger *zap.Logger) error {
	marketDataController := controller.NewMarketDataController(logger)
	ws.Get("/marketdata", marketDataController.ServeMarketData())
	return nil
}

func ordersController(v1 fiber.Router, goqu *goqu.Database, logger *zap.Logger, cfg config.AppConfig, middleware middlewares.Middleware) error {
	orderController, err := controller.NewOrderController(logger, goqu, cfg)
	if err != nil {
		return err
	}

	v1.Post("/orders", middleware.Authenticated, orderController.PlaceOrder)
	v1.Get("/trade-history", middleware.Authenticated, orderController.GetOrders)
	v1.Get("/position", middleware.Authenticated, orderController.GetPositions)
	return nil
}
