package config

import (
	"os"
	"time"
)

type Config struct {
	Server   ServerConfig
	DB       DBConfig
	MQ       MQConfig
	UseMySQL bool
}

type ServerConfig struct {
	Port int
}

type DBConfig struct {
	Type string
	DSN  string
}

type MQConfig struct {
	URL            string
	Topic          string
	ConsumerGroup  string
	ReconnectDelay time.Duration
}

func NewConfig() *Config {
	useMySQL := os.Getenv("USE_MYSQL") == "true"

	var dbConfig DBConfig
	if useMySQL {
		dbConfig = DBConfig{
			Type: "mysql",
			DSN:  os.Getenv("MYSQL_DSN"),
		}
		if dbConfig.DSN == "" {
			dbConfig.DSN = "root:password@tcp(localhost:3306)/notification?charset=utf8mb4&parseTime=True&loc=Local"
		}
	} else {
		dbConfig = DBConfig{
			Type: "memory",
			DSN:  "notification.db",
		}
	}

	return &Config{
		Server: ServerConfig{
			Port: 8080,
		},
		DB:       dbConfig,
		UseMySQL: useMySQL,
		MQ: MQConfig{
			URL:            "localhost:9092",
			Topic:          "notification_callbacks",
			ConsumerGroup:  "notification_system",
			ReconnectDelay: 5 * time.Second,
		},
	}
}
