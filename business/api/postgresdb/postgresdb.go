package postgresdb

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/url"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v5/pgxpool"
)

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
	// perform db migration, fatal if not successful
	err = migration(config)
	if err != nil {
		log.Fatalf("Failed to perform db migration: %s", err)
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

func migration(config Config) error {
	migrationSourceUrl := "file://business/api/migration"

	dbUrl := buildConnectionString(config)
	m, err := migrate.New(migrationSourceUrl, dbUrl)
	if err != nil {
		return fmt.Errorf("error building db migration url: %s, %s, %s\n", migrationSourceUrl, dbUrl, err.Error())
	}
	if err = m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return fmt.Errorf("Error with DB migration: %s\n", err.Error())
	}
	return nil
}
