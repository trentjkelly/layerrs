package config

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/pressly/goose/v3"
	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/lib/pq"
)

type DBConfig struct {
	Host string
	Port string
	Database string
	User string
	Password string
	psqlURL string
}

const (
	PSQL_HOST_DOCKER = "PSQL_HOST_DOCKER"
	PSQL_HOST_LOCAL = "PSQL_HOST_LOCAL"
	PSQL_PORT = "PSQL_PORT_%s"
	PSQL_DB = "PSQL_DB_%s"
	PSQL_USER = "PSQL_USER_%s"
	PSQL_PASSWORD = "PSQL_PASSWORD_%s"
)

// Creates a new database configuration
func NewDBConfig(env string, isDocker bool) (*DBConfig, error) {
	dbConfig := new(DBConfig)

	if isDocker {
		dbConfig.Host = os.Getenv(PSQL_HOST_DOCKER)
	} else {
		dbConfig.Host = os.Getenv(PSQL_HOST_LOCAL)
	}
	if dbConfig.Host == "" {
		return nil, fmt.Errorf("could not find the environment variable PSQL_HOST_%s", env)
	}
	
	dbConfig.Port = os.Getenv(fmt.Sprintf(PSQL_PORT, env))
	if dbConfig.Port == "" {
		return nil, fmt.Errorf("could not find the environment variable PSQL_PORT_%s", env)
	}

	dbConfig.Database = os.Getenv(fmt.Sprintf(PSQL_DB, env))
	if dbConfig.Database == "" {
		return nil, fmt.Errorf("could not find the environment variable PSQL_DB_%s", env)
	}

	dbConfig.User = os.Getenv(fmt.Sprintf(PSQL_USER, env))
	if dbConfig.User == "" {
		return nil, fmt.Errorf("could not find the environment variable PSQL_USER_%s", env)
	}

	dbConfig.Password = os.Getenv(fmt.Sprintf(PSQL_PASSWORD, env))
	if dbConfig.Password == "" {
		return nil, fmt.Errorf("could not find the environment variable PSQL_PASSWORD_%s", env)
	}

	dbConfig.psqlURL = fmt.Sprintf("postgresql://%s:%s@%s:%s/%s?sslmode=disable", dbConfig.User, dbConfig.Password, dbConfig.Host, dbConfig.Port, dbConfig.Database)

	return dbConfig, nil
}

// Initializes a database connection, applies migrations, and returns a connection pool
func InitDB(env string, isDocker bool) (*pgxpool.Pool, error) {
	dbConfig, err := NewDBConfig(env, isDocker)
	if err != nil {
		return nil, fmt.Errorf("failed to create database config: %w", err)
	}

	db, err := dbConfig.CreatePSQLConnection()
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

	pool, err := dbConfig.CreatePSQLPoolConnection()
	if err != nil {
		return nil, fmt.Errorf("failed to create connection pool: %w", err)
	}

	return pool, nil
}

// Creates a connection pool to the PostgreSQL database
func (dbConfig *DBConfig) CreatePSQLPoolConnection() (*pgxpool.Pool, error) {
	cfg, err := pgxpool.ParseConfig(dbConfig.psqlURL)
	if err != nil {
		return nil, fmt.Errorf("failed to parse config: %w", err)
	}

	cfg.MaxConns = 10
	cfg.MinConns = 2
	cfg.MaxConnLifetime = 1 * time.Hour
	cfg.MaxConnIdleTime = 30 * time.Minute

	pool, err := pgxpool.NewWithConfig(context.Background(), cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to create connection pool: %w", err)
	}

	err = pool.Ping(context.Background())
	if err != nil {
		return nil, fmt.Errorf("failed to ping connection pool using url: %s: %w", dbConfig.psqlURL, err)
	}

	return pool, nil
}

// Creates a connection to the PostgreSQL database -- only used for migrations
func (dbConfig *DBConfig) CreatePSQLConnection() (*sql.DB, error) {

	db, err := sql.Open("postgres", dbConfig.psqlURL)
	if err != nil {
		return nil, fmt.Errorf("failed to open connection: %w", err)
	}

	log.Println(dbConfig.psqlURL)

	err = db.Ping()
	if err != nil {
		return nil, fmt.Errorf("failed to ping database for migrations using url: %s: %w", dbConfig.psqlURL, err)
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
