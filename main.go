package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"runtime"
	"syscall"
	"time"

	"github.com/ardanlabs/conf/v3"
	"github.com/danipurwadi/internal-transfer-system/app/debug"
	"github.com/danipurwadi/internal-transfer-system/app/middleware"
	"github.com/danipurwadi/internal-transfer-system/app/transferapp"
	"github.com/danipurwadi/internal-transfer-system/business/api/db"
	"github.com/danipurwadi/internal-transfer-system/business/transferbus"
	"github.com/danipurwadi/internal-transfer-system/business/transferbus/stores/transferdb"
	"github.com/danipurwadi/internal-transfer-system/foundation/logger"
	"github.com/danipurwadi/internal-transfer-system/foundation/web"
)

var build = "develop"

func main() {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("Global recovery from panic. Error: %s", r)
		}
	}()

	var log *logger.Logger

	events := logger.Events{
		Error: func(ctx context.Context, r logger.Record) {
			log.Info(ctx, "******* SEND ALERT *******")
		},
	}

	traceIDFn := func(ctx context.Context) string {
		return web.GetTraceID(ctx)
	}

	log = logger.NewWithEvents(os.Stdout, logger.LevelInfo, "TRANSFER", traceIDFn, events)

	ctx := context.Background()

	if err := run(ctx, log); err != nil {
		log.Error(ctx, "startup", "msg", err)
		os.Exit(1)
	}
}

func run(ctx context.Context, log *logger.Logger) error {
	// -------------------------------------------------------------------------
	// GOMAXPROCS
	log.Info(ctx, "startup", "GOMAXPROCS", runtime.GOMAXPROCS(0))

	// -------------------------------------------------------------------------
	// CONFIGS
	cfg := struct {
		conf.Version
		Web struct {
			ReadTimeout        time.Duration `conf:"default:5s"`
			WriteTimeout       time.Duration `conf:"default:10s"`
			IdleTimeout        time.Duration `conf:"default:120s"`
			ShutdownTimeout    time.Duration `conf:"default:20s"`
			APIHost            string        `conf:"default:0.0.0.0:8080"`
			DebugHost          string        `conf:"default:0.0.0.0:8090"`
			CORSAllowedOrigins []string      `conf:"default:*"`
		}
		DB struct {
			User       string `conf:"default:postgres"`
			Password   string `conf:"default:password,mask"`
			Host       string `conf:"default:host.docker.internal"`
			Port       int    `conf:"default:5432"`
			Name       string `conf:"default:transfer"`
			DisableTLS bool   `conf:"default:true"`
		}
	}{
		Version: conf.Version{
			Build: build,
			Desc:  "Transfer",
		},
	}

	const prefix = "TRANSFER"
	_, err := conf.Parse(prefix, &cfg)
	if err != nil {
		return fmt.Errorf("parsing config: %w", err)
	}

	// -------------------------------------------------------------------------
	// App Starting

	log.Info(ctx, "starting service", "version", cfg.Build)
	defer log.Info(ctx, "shutdown complete")

	out, err := conf.String(&cfg)
	if err != nil {
		return fmt.Errorf("generating config for output: %w", err)
	}
	log.Info(ctx, "startup", "config", out)

	// -------------------------------------------------------------------------
	// Database Support

	log.Info(ctx, "startup", "status", "initializing database support", "hostport", cfg.DB.Host)

	dbConfig := db.Config{
		User:       cfg.DB.User,
		Password:   cfg.DB.Password,
		HostPort:   fmt.Sprintf("%s:%d", cfg.DB.Host, cfg.DB.Port),
		Database:   cfg.DB.Name,
		DisableTLS: cfg.DB.DisableTLS,
	}

	dbConn := db.New(dbConfig)
	db.Migrate(dbConfig)

	dbClient := transferdb.NewTxQueries(dbConn)

	// -------------------------------------------------------------------------
	// Start API Service

	log.Info(ctx, "startup", "status", "initializing V1 API support")

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, syscall.SIGINT, syscall.SIGTERM)

	// initialise business layer
	transferBus := transferbus.New(dbClient)

	// initialise app layer
	transferApp := transferapp.NewApp(transferBus)

	// intitialise and register routes to the client
	webClient := web.NewClient(shutdown, middleware.Logger(log), middleware.Errors(log))
	transferApp.Routes(webClient)

	api := http.Server{
		Addr:         cfg.Web.APIHost,
		Handler:      webClient,
		ReadTimeout:  cfg.Web.ReadTimeout,
		WriteTimeout: cfg.Web.WriteTimeout,
		IdleTimeout:  cfg.Web.IdleTimeout,
		ErrorLog:     logger.NewStdLogger(log, logger.LevelError),
	}

	serverErrors := make(chan error, 1)

	go func() {
		log.Info(ctx, "startup", "status", "api router started", "host", api.Addr)
		serverErrors <- api.ListenAndServe()
	}()

	// -------------------------------------------------------------------------
	// Start Debug Service

	go func() {
		log.Info(ctx, "startup", "status", "debug v1 router started", "host", cfg.Web.DebugHost)

		if err := http.ListenAndServe(cfg.Web.DebugHost, debug.Mux()); err != nil {
			log.Error(ctx, "shutdown", "status", "debug v1 router closed", "host", cfg.Web.DebugHost, "msg", err)
		}
	}()

	// -------------------------------------------------------------------------
	// Handle shutdown

	select {
	case err := <-serverErrors:
		return fmt.Errorf("server error: %w", err)

	case sig := <-shutdown:
		log.Info(ctx, "shutdown", "status", "shutdown started", "signal", sig)
		defer log.Info(ctx, "shutdown", "status", "shutdown complete", "signal", sig)

		ctx, cancel := context.WithTimeout(ctx, cfg.Web.ShutdownTimeout)
		defer cancel()

		if err := api.Shutdown(ctx); err != nil {
			api.Close()
			return fmt.Errorf("could not stop server gracefully: %w", err)
		}
	}

	return nil
}
