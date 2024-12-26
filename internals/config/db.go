package config

import (
	"time"
	"os"
	"log"
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
)

func CreatePSQLPoolConnection() *pgxpool.Pool {

	// Create connection pool configuration
	cfg, err := pgxpool.ParseConfig(os.Getenv("PSQL_DATABASE_URL"))

	if err != nil {
		log.Fatal(err)
	}

	// Settings for DB connection pool
	cfg.MaxConns = 10
	cfg.MinConns = 2
	cfg.MaxConnLifetime = 1 * time.Hour
	cfg.MaxConnIdleTime = 30 * time.Minute

	// Create connection pool
	pool, err := pgxpool.NewWithConfig(context.Background(), cfg)

	if err != nil {
		log.Fatal(err)
	}

	// Check that connection works properly
	err = pool.Ping(context.Background())

	if err != nil {
		log.Fatal(err)
	}

	return pool
}