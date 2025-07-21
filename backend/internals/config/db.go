package config

import (
	"fmt"
	"path/filepath"
	"time"
	"os"
	"context"
	"database/sql"
	"github.com/pressly/goose/v3"
	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/lib/pq"
)

// Initializes a database connection, applies migrations, and returns a connection pool
func InitDB() (*pgxpool.Pool, error) {
	db, err := CreatePSQLConnection()
	if err != nil {
		return nil, fmt.Errorf("failed to create connection: %w", err)
	}

	err = ApplyMigrations(db)
	if err != nil {
		return nil, fmt.Errorf("failed to apply migrations: %w", err)
	}

	err = db.Close()
	if err != nil {
		return nil, fmt.Errorf("failed to close connection: %w", err)
	}

	pool, err := CreatePSQLPoolConnection()
	if err != nil {
		return nil, fmt.Errorf("failed to create connection pool: %w", err)
	}

	return pool, nil
}

// Creates a connection pool to the PostgreSQL database
func CreatePSQLPoolConnection() (*pgxpool.Pool, error) {
	// Create connection pool configuration	
	host := os.Getenv("PSQL_HOST")
	port := os.Getenv("PSQL_PORT")
	database := os.Getenv("PSQL_DB")
	user := os.Getenv("PSQL_USER")
	password := os.Getenv("PSQL_PASSWORD")
	
	psqlURL := fmt.Sprintf("postgresql://%s:%s/%s?user=%s&password=%s&sslmode=disable", host, port, database, user, password)

	cfg, err := pgxpool.ParseConfig(psqlURL)
	if err != nil {
		return nil, fmt.Errorf("failed to parse config: %w", err)
	}

	// Settings for DB connection pool
	cfg.MaxConns = 10
	cfg.MinConns = 2
	cfg.MaxConnLifetime = 1 * time.Hour
	cfg.MaxConnIdleTime = 30 * time.Minute

	// Create connection pool
	pool, err := pgxpool.NewWithConfig(context.Background(), cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to create connection pool: %w", err)
	}

	// Check that connection works properly
	err = pool.Ping(context.Background())
	if err != nil {
		return nil, fmt.Errorf("failed to ping connection pool: %w", err)
	}

	return pool, nil
}

// Creates a connection to the PostgreSQL database -- only used for migrations
func CreatePSQLConnection() (*sql.DB, error) {
	host := os.Getenv("PSQL_HOST")
	port := os.Getenv("PSQL_PORT")
	database := os.Getenv("PSQL_DB")
	user := os.Getenv("PSQL_USER")
	password := os.Getenv("PSQL_PASSWORD")

	psqlURL := fmt.Sprintf("postgresql://%s:%s/%s?user=%s&password=%s&sslmode=disable", host, port, database, user, password)

	db, err := sql.Open("postgres", psqlURL)
	if err != nil {
		return nil, fmt.Errorf("failed to open connection: %w", err)
	}

	err = db.Ping()
	if err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return db, nil
}

// Applies the migrations to the PostgreSQL database
func ApplyMigrations(db *sql.DB) (error) {
	if err := goose.SetDialect("postgres"); err != nil {
		return fmt.Errorf("failed to set dialect: %w", err)
	}

	dir, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get current directory: %w", err)
	}

	migrationsDir := filepath.Join(dir, "migrations", "db")

	if err := goose.Up(db, migrationsDir); err != nil {
		return fmt.Errorf("failed to apply migrations: %w", err)
	}

	return nil
}
