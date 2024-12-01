package config

import (
	"log"

	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
)

var AllConfig AppConfig

type AppConfig struct {
	IsDevelopment bool   `envconfig:"IS_DEVELOPMENT"`
	Debug         bool   `envconfig:"DEBUG"`
	Port          string `envconfig:"APP_PORT"`
	Secret        string `envconfig:"JWT_SECRET"`
	DB            DBConfig
	Binance       BinanceConfig
}

func GetConfig() AppConfig {
	err := godotenv.Load()
	if err != nil {
		log.Println("warning .env file not found, scanning from OS ENV")
	}

	AllConfig = AppConfig{}

	err = envconfig.Process("APP_PORT", &AllConfig)
	if err != nil {
		log.Fatal(err)
	}

	return AllConfig
}
