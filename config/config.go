package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/spf13/cast"
)

type Config struct {
	ORDER_SERVICE_PORT string
	AUTH_SERVICE_PORT  string
	DB_HOST            string
	DB_PORT            string
	DB_NAME            string
	DB_USER            string
	DB_PASSWORD        string
	REDIS_HOST         string
	REDIS_PORT         string
	REDIS_PASSWORD     string
	LOG_PATH           string
	APP_PASSWORD       string
}

func Load() *Config {
	if err := godotenv.Load(".env"); err != nil {
		log.Fatal("Error while loading .env file")
	}

	config := Config{}

	config.ORDER_SERVICE_PORT = cast.ToString(coalesce("ORDER_SERVICE_PORT", "localhost"))
	config.AUTH_SERVICE_PORT = cast.ToString(coalesce("AUTH_SERVICE_PORT", "localhost"))
	config.DB_HOST = cast.ToString(coalesce("DB_HOST", "localhost"))
	config.DB_PORT = cast.ToString(coalesce("DB_PORT", "5432"))
	config.DB_USER = cast.ToString(coalesce("DB_USER", "postgres"))
	config.DB_NAME = cast.ToString(coalesce("DB_NAME", "name"))
	config.DB_PASSWORD = cast.ToString(coalesce("DB_PASSWORD", "root"))
	config.REDIS_HOST = cast.ToString(coalesce("REDIS_HOST", "root"))
	config.REDIS_PORT = cast.ToString(coalesce("REDIS_PORT", "root"))
	config.REDIS_PASSWORD = cast.ToString(coalesce("REDIS_PASSWORD", "root"))
	config.LOG_PATH = cast.ToString(coalesce("LOG_PATH", "areyouinterested.log"))
	config.APP_PASSWORD = cast.ToString(coalesce("APP_PASSWORD", "COMMONMAN"))

	return &config
}

func coalesce(key string, value interface{}) interface{} {
	val, exist := os.LookupEnv(key)
	if exist {
		return val
	}
	return value
}
