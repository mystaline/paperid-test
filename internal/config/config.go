package config

import (
	"cmp"
	"os"
)

type AppConfig struct {
	Port string

	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string
	DBSSLMode  string
}

func LoadConfig() AppConfig {
	appConfig := AppConfig{}

	appConfig.Port = cmp.Or(os.Getenv("PORT"), "8080")

	appConfig.DBHost = cmp.Or(os.Getenv("POSTGRES_HOST"), "localhost")
	appConfig.DBPort = cmp.Or(os.Getenv("POSTGRES_PORT"), "5432")
	appConfig.DBUser = cmp.Or(os.Getenv("POSTGRES_USER"), "postgres")
	appConfig.DBPassword = cmp.Or(os.Getenv("POSTGRES_PASSWORD"), "postgres")
	appConfig.DBName = cmp.Or(os.Getenv("POSTGRES_DB"), "paperid_test")
	appConfig.DBSSLMode = cmp.Or(os.Getenv("POSTGRES_SSLMODE"), "disable")

	return appConfig
}
