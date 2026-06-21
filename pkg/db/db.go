// To connect and manage db lifecycle
package db

import (
	"context"
	"fmt"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/mystaline/paperid-test/internal/config"
)

func NewPool(appConfig config.AppConfig) (*pgxpool.Pool, error) {
	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		appConfig.DBHost,
		appConfig.DBPort,
		appConfig.DBUser,
		appConfig.DBPassword,
		appConfig.DBName,
		appConfig.DBSSLMode,
	)

	cfg, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		log.Fatalf("Invalid PostgreSQL config: %v", err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	return pgxpool.NewWithConfig(ctx, cfg)
}
