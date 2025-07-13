// Package db contains database access logic.
package db

import (
	"context"
	"embed"
	"errors"
	"fmt"
	"log"
	"net/url"
	"time"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

//go:embed migration/*.sql
var migrationFiles embed.FS

type Config struct {
	User       string
	Password   string
	HostPort   string
	Database   string
	DisableTLS bool
}

func New(config Config) *pgxpool.Pool {
	pool, err := pgxpool.New(context.Background(), buildConnectionString(config))
	if err != nil {
		log.Fatalf("Unable to initialise db connection %v\n", err)
	}
	return pool
}

func buildConnectionString(config Config) string {
	format := "postgresql://%s:%s@%s/%s"

	// auto-configure for local dockerized instance
	if config.DisableTLS {
		format = "postgresql://%s:%s@%s/%s?sslmode=disable"
	}

	return fmt.Sprintf(
		format,
		config.User,
		url.QueryEscape(config.Password),
		config.HostPort,
		config.Database,
	)
}

func Migrate(config Config) error {
	sourceDriver, err := iofs.New(migrationFiles, "migration")
	if err != nil {
		return fmt.Errorf("error creating migration source driver: %w", err)
	}

	dbURL := buildConnectionString(config)
	m, err := migrate.NewWithSourceInstance("iofs", sourceDriver, dbURL)
	if err != nil {
		return fmt.Errorf("error building db migration: %w", err)
	}
	if err = m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return fmt.Errorf("error with DB migration: %w", err)
	}
	return nil
}

func InitDatabase(ctx context.Context, config Config, newDBName string) error {
	connString := buildConnectionString(config)
	// 2. Establish a single, temporary connection
	conn, err := connectWithRetry(ctx, connString)
	if err != nil {
		return fmt.Errorf("failed to connect to maintenance database: %w", err)
	}
	defer conn.Close(ctx)

	// 3. Check if the database already exists to avoid errors
	var exists bool
	err = conn.QueryRow(ctx, "SELECT EXISTS(SELECT 1 FROM pg_database WHERE datname = $1)", newDBName).Scan(&exists)
	if err != nil {
		return fmt.Errorf("failed to check if database exists: %w", err)
	}
	if exists {
		return nil
	}

	// 4. Execute CREATE DATABASE using pgx.Identifier for safety against SQL injection
	// Using fmt.Sprintf with %q is a simpler but less robust alternative.
	createCommand := fmt.Sprintf("CREATE DATABASE %s", pgx.Identifier{newDBName}.Sanitize())

	_, err = conn.Exec(ctx, createCommand)
	if err != nil {
		return fmt.Errorf("failed to execute CREATE DATABASE command: %w", err)
	}
	return nil
}

func connectWithRetry(ctx context.Context, connStr string) (*pgx.Conn, error) {
	var conn *pgx.Conn
	var err error

	// Define retry parameters
	maxRetries := 10
	retryDelay := 1 * time.Second

	for range maxRetries {
		// Try to connect
		conn, err = pgx.Connect(ctx, connStr)
		if err == nil {
			return conn, nil
		}
		time.Sleep(retryDelay)
	}

	// If all retries fail, return the last error
	return nil, fmt.Errorf("failed to connect to database after %d attempts: %w", maxRetries, err)
}
