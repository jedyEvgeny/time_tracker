package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DatabaseScheme          string
	HTTPServerHost          string
	HTTPServerPort          string
	HTTPServerPath          string
	StoragePostgresHost     string
	StoragePostgresPort     string
	StoragePostgresUser     string
	StoragePostgresPassword string
	StoragePostgresDBName   string
	StoragePostgresSSLMode  string
	HTTPClientHost          string
}

func NewConfig() (Config, error) {
	err := godotenv.Load(".env")
	if err != nil {
		log.Printf("ошибка загрузки переменных окружения: %v\n", err)
		return Config{}, err
	}
	cfg := Config{
		HTTPServerHost:          os.Getenv("APP_HTTP_SERVER_HOST"),
		HTTPServerPort:          os.Getenv("APP_HTTP_SERVER_PORT"),
		HTTPServerPath:          os.Getenv("APP_HTTP_SERVER_PATH"),
		DatabaseScheme:          os.Getenv("APP_DATABASE_SCHEME"),
		StoragePostgresHost:     os.Getenv("APP_STORAGE_POSTGRES_HOST"),
		StoragePostgresPort:     os.Getenv("APP_STORAGE_POSTGRES_PORT"),
		StoragePostgresUser:     os.Getenv("APP_STORAGE_POSTGRES_USER"),
		StoragePostgresPassword: os.Getenv("APP_STORAGE_POSTGRES_PASSWORD"),
		StoragePostgresDBName:   os.Getenv("APP_STORAGE_POSTGRES_DBNAME"),
		StoragePostgresSSLMode:  os.Getenv("APP_STORAGE_POSTGRES_SSLMODE"),
		HTTPClientHost:          os.Getenv("APP_HTTP_CLIENT_HOST"),
	}
	return cfg, nil
}
