package config

import (
	"log"

	"github.com/spf13/viper"
)

type Config struct {
	PostgresDSN string
	RedisAddr   string
	RabbitMQURL string
	JWTSecret   string
	Port        string
}

func Load() *Config {
	viper.SetConfigFile(".env")
	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Error reading config: %v", err)
	}

	return &Config{
		PostgresDSN: viper.GetString("POSTGRES_DSN"),
		RedisAddr:   viper.GetString("REDIS_ADDR"),
		RabbitMQURL: viper.GetString("RABBITMQ_URL"),
		JWTSecret:   viper.GetString("JWT_SECRET"),
		Port:        viper.GetString("PORT"),
	}
}
