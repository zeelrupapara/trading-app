package config

type BinanceConfig struct {
	APIKey    string `envconfig:"BINANCE_API_KEY" validate:"required"`
	APISecret string `envconfig:"BINANCE_API_SECRET" validate:"required"`
}
