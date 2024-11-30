package main

import (
	"github.com/zeelrupapara/trading-api/cli"
	"github.com/zeelrupapara/trading-api/config"
	"github.com/zeelrupapara/trading-api/logger"

	"go.uber.org/zap"
)

// Golang API.
//
//	Schemes: https
//	Host: localhost
//	BasePath: /api/v1
//	Version: 0.0.1-alpha
//
//	Consumes:
//	- application/json
//
//	Produces:
//	- application/json
//
// swagger:meta
func main() {
	// Collecting config from env or file or flag
	cfg := config.GetConfig()

	logger, err := logger.NewRootLogger(cfg.Debug, cfg.IsDevelopment)
	if err != nil {
		panic(err)
	}
	zap.ReplaceGlobals(logger)

	err = cli.Init(cfg, logger)
	if err != nil {
		panic(err)
	}
}
