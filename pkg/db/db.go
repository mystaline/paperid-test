// To connect and manage db lifecycle
package db

import (
	"context"
	"fmt"
	"log"
	"sync"

	"github.com/bwmarrin/snowflake"
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
		return nil, fmt.Errorf("parse db config: %w", err)
	}

	pool, err := pgxpool.NewWithConfig(context.Background(), cfg)
	if err != nil {
		return nil, fmt.Errorf("create pool: %w", err)
	}

	if err := pool.Ping(context.Background()); err != nil {
		pool.Close()
		return nil, fmt.Errorf("ping db: %w", err)
	}

	return pool, nil
}

var (
	idNode     *snowflake.Node
	idNodeOnce sync.Once
)

func InitIDGenerator() {
	idNodeOnce.Do(func() {
		var err error
		idNode, err = snowflake.NewNode(1)
		if err != nil {
			log.Fatalf("snowflake init failed: %v", err)
		}
	})
}

func IDNode() *snowflake.Node {
	if idNode == nil {
		log.Fatal("call db.InitIDGenerator() first")
	}
	return idNode
}
