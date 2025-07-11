package postgresdb

import (
	"context"
	"fmt"
	"log"
	"net/url"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Config struct {
	User     string
	Password string
	Host     string
	Port     int
	Database string
}

func New(config Config) *pgxpool.Pool {
	pool, err := pgxpool.New(context.Background(), buildConnectionString(config))
	if err != nil {
		log.Fatalf("Unable to initialise db connection %v\n", err)
	}
	return pool
}

func buildConnectionString(config Config) string {
	format := "postgresql://%s:%s@%s:%d/%s"

	// auto-configure for local dockerized instance
	if config.Password == "password" {
		format = "postgresql://%s:%s@%s:%d/%s?sslmode=disable"
	}

	return fmt.Sprintf(
		format,
		config.User,
		url.QueryEscape(config.Password),
		config.Host,
		config.Port,
		config.Database,
	)
}
