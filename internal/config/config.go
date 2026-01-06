package config

import (
	"time"
)

type Config struct {
	Server ServerConfig
	DB     DBConfig
	MQ     MQConfig
}

type ServerConfig struct {
	Port int
}

type DBConfig struct {
	DSN string
}

type MQConfig struct {
	URL            string
	Topic          string
	ConsumerGroup  string
	ReconnectDelay time.Duration
}

func NewConfig() *Config {
	return &Config{
		Server: ServerConfig{
			Port: 8080,
		},
		DB: DBConfig{
			DSN: "notification.db",
		},
		MQ: MQConfig{
			URL:            "localhost:9092",
			Topic:          "notification_callbacks",
			ConsumerGroup:  "notification_system",
			ReconnectDelay: 5 * time.Second,
		},
	}
}
