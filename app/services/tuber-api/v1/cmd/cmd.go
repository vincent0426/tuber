package cmd

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"runtime"
	"syscall"

	"github.com/TSMC-Uber/server/app/services/tuber-api/v1/config"
	"github.com/TSMC-Uber/server/business/sys/database"
	v1 "github.com/TSMC-Uber/server/business/web/v1"
	"github.com/TSMC-Uber/server/business/web/v1/debug"

	"github.com/TSMC-Uber/server/foundation/logger"
	"go.uber.org/zap"
)

func Main(build string, routeAdder v1.RouteAdder) error {
	log, err := logger.New("TUBER-API")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer log.Sync()

	ctx := context.Background()
	if err := run(ctx, log, build, routeAdder); err != nil {
		log.Errorw("startup", "ERROR", err)
		log.Sync()
		os.Exit(1)
	}

	return nil
}

func run(ctx context.Context, log *zap.SugaredLogger, build string, routeAdder v1.RouteAdder) error {

	// -------------------------------------------------------------------------
	// GOMAXPROCS

	log.Infow("startup", "GOMAXPROCS", runtime.GOMAXPROCS(0), "BUILD-", build)

	// -------------------------------------------------------------------------
	// Configuration

	cfg, err := config.New()
	if err != nil {
		return fmt.Errorf("loading config: %w", err)
	}

	log.Infow("startup", "status", "loaded configuration", "config", cfg)
	log.Infow("startup", "status", "debug v1 router started", "host", cfg.Web.DebugHost)

	go func() {
		if err := http.ListenAndServe(cfg.Web.DebugHost, debug.Mux(build, log)); err != nil {
			log.Errorw("shutdown", "status", "debug v1 router closed", "host", cfg.Web.DebugHost, "ERROR", err)
		}
	}()

	// -------------------------------------------------------------------------
	// App Starting

	log.Infow("starting service", "version", build)
	defer log.Infow("shutdown complete")

	// -------------------------------------------------------------------------
	// Database Support

	log.Infow("startup", "status", "initializing database support", "host", cfg.DB.Host)

	db, err := database.Open(database.Config{
		User:         cfg.DB.User,
		Password:     cfg.DB.Password,
		Host:         cfg.DB.Host,
		Name:         cfg.DB.Name,
		MaxIdleConns: cfg.DB.MaxIdleConns,
		MaxOpenConns: cfg.DB.MaxOpenConns,
		DisableTLS:   cfg.DB.DisableTLS,
	})
	if err != nil {
		return fmt.Errorf("connecting to db: %w", err)
	}
	defer func() {
		log.Infow("shutdown", "status", "stopping database support", "host", cfg.DB.Host)
		db.Close()
	}()

	// -------------------------------------------------------------------------
	// Start API Service

	log.Infow("startup", "status", "initializing V1 API support")

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, syscall.SIGINT, syscall.SIGTERM)

	cfgMux := v1.APIMuxConfig{
		Shutdown: shutdown,
		Log:      log,
		// Auth:     auth,
		DB: db,
	}

	apiMux := v1.APIMux(cfgMux, routeAdder)

	api := http.Server{
		Addr:         cfg.Web.APIHost,
		Handler:      apiMux,
		ReadTimeout:  cfg.Web.ReadTimeout,
		WriteTimeout: cfg.Web.WriteTimeout,
		IdleTimeout:  cfg.Web.IdleTimeout,
		ErrorLog:     zap.NewStdLog(log.Desugar()),
	}

	serverErrors := make(chan error, 1)

	go func() {
		log.Infow("startup", "status", "api router started", "host", api.Addr)
		serverErrors <- api.ListenAndServe()
	}()

	// -------------------------------------------------------------------------
	// Shutdown

	select {
	case err := <-serverErrors:
		return fmt.Errorf("server error: %w", err)

	case sig := <-shutdown:
		log.Infow("shutdown", "status", "shutdown started", "signal", sig)
		defer log.Infow("shutdown", "status", "shutdown complete", "signal", sig)

		ctx, cancel := context.WithTimeout(context.Background(), cfg.Web.ShutdownTimeout)
		defer cancel()

		if err := api.Shutdown(ctx); err != nil {
			api.Close()
			return fmt.Errorf("could not stop server gracefully: %w", err)
		}
	}

	return nil
}
