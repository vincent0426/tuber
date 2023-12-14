package cmd

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"runtime"
	"syscall"
	"time"

	"github.com/TSMC-Uber/server/app/services/tuber-api/v1/config"
	"github.com/TSMC-Uber/server/business/core/locationws"
	"github.com/TSMC-Uber/server/business/sys/cachedb"
	"github.com/TSMC-Uber/server/business/sys/database"
	"github.com/TSMC-Uber/server/business/sys/mail"
	"github.com/TSMC-Uber/server/business/sys/mq"
	v1 "github.com/TSMC-Uber/server/business/web/v1"
	"github.com/TSMC-Uber/server/business/web/v1/auth"
	"github.com/TSMC-Uber/server/business/web/v1/debug"

	"github.com/TSMC-Uber/server/foundation/logger"
	"github.com/TSMC-Uber/server/foundation/web"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
)

func Main(routes, build string, routeAdder v1.RouteAdder) error {
	var log *logger.Logger

	events := logger.Events{
		Error: func(ctx context.Context, r logger.Record) {
			log.Info(ctx, "******* SEND ALERT ******")
		},
	}

	traceIDFunc := func(ctx context.Context) string {
		return web.GetTraceID(ctx)
	}

	log = logger.NewWithEvents(os.Stdout, logger.LevelInfo, "TUBER-API", traceIDFunc, events)

	ctx := context.Background()
	if err := run(ctx, log, build, routes, routeAdder); err != nil {
		log.Error(ctx, "startup", "status", "shutdown with error", "ERROR", err)
		os.Exit(1)
	}

	return nil
}

func run(ctx context.Context, log *logger.Logger, build string, routes string, routeAdder v1.RouteAdder) error {

	// -------------------------------------------------------------------------
	// GOMAXPROCS

	log.Info(ctx, "startup", "GOMAXPROCS", runtime.GOMAXPROCS(0))

	// -------------------------------------------------------------------------
	// Configuration

	cfg, err := config.New()
	if err != nil {
		return fmt.Errorf("loading config: %w", err)
	}

	log.Info(ctx, "startup", "status", "loaded configuration", "config", cfg)

	// -------------------------------------------------------------------------
	// App Starting

	log.Info(ctx, "starting service", "version", build)
	defer log.Info(ctx, "shutdown complete")

	// -------------------------------------------------------------------------
	// RabbitMQ Support

	log.Info(ctx, "startup", "status", "initializing rabbitmq support", "host", cfg.RabbitMQ.Host)

	_, err = mq.ConnectToRabbitMQ(mq.Config{
		Host:              cfg.RabbitMQ.Host,
		DelayQueueName:    cfg.RabbitMQ.DelayQueueName,
		DelayExchangeName: cfg.RabbitMQ.DelayExchangeName,
		DelayRoutingKey:   cfg.RabbitMQ.DelayRoutingKey,
	})
	if err != nil {
		return fmt.Errorf("connecting to rabbitmq: %w", err)
	}

	if routes == "locationserver" {
		go mail.StartSendEmailWorker()
	}

	defer func() {
		log.Info(ctx, "shutdown", "status", "stopping rabbitmq support", "host", cfg.RabbitMQ.Host)
		mq.Close()
	}()

	// -------------------------------------------------------------------------
	// Room Dispatcher
	locationws.NewRoomsDispatcher()

	// -------------------------------------------------------------------------
	// Database Support

	log.Info(ctx, "startup", "status", "initializing database support", "host", cfg.DB.Host)

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
		log.Info(ctx, "shutdown", "status", "stopping database support", "host", cfg.DB.Host)
		db.Close()
	}()

	// check status of database
	if err := database.StatusCheck(ctx, db); err != nil {
		return fmt.Errorf("database status check: %w", err)
	}

	cachedb, err := cachedb.Open(cachedb.Config{
		MasterHost:      cfg.Redis.Host.Master,
		MasterPassword:  cfg.Redis.Password,
		MasterDB:        cfg.Redis.DB,
		ReplicaHost:     cfg.Redis.Host.Replica,
		ReplicaPassword: cfg.Redis.Password,
		ReplicaDB:       cfg.Redis.DB,
	})
	if err != nil {
		return fmt.Errorf("connecting to master redis: %w", err)
	}
	defer func() {
		log.Info(ctx, "shutdown", "status", "stopping redis support", "host", cfg.Redis.Host)
		cachedb.Master.Close()
		cachedb.Replica.Close()
	}()

	authCfg := auth.Config{
		Log: log,
		DB:  db,
		// KeyLookup: vault,
	}

	auth, err := auth.New(authCfg)
	if err != nil {
		return fmt.Errorf("constructing auth: %w", err)
	}

	// -------------------------------------------------------------------------
	// Start Tracing Support

	log.Info(ctx, "startup", "status", "initializing OT/Tempo tracing support")

	traceProvider, err := startTracing(
		cfg.Tempo.ReporterURI,
		cfg.Tempo.ServiceName,
		cfg.Tempo.Probability,
	)
	if err != nil {
		return fmt.Errorf("starting tracing: %w", err)
	}
	defer traceProvider.Shutdown(context.Background())

	tracer := traceProvider.Tracer("service")

	// -------------------------------------------------------------------------
	// Start Debug Service

	go func() {
		log.Info(ctx, "startup", "status", "debug v1 router started", "host", cfg.Web.DebugHost)

		if err := http.ListenAndServe(cfg.Web.DebugHost, debug.Mux(build, log)); err != nil {
			log.Error(ctx, "shutdown", "status", "debug v1 router closed", "host", cfg.Web.DebugHost, "msg", err)
		}
	}()

	// -------------------------------------------------------------------------
	// Start API Service

	log.Info(ctx, "startup", "status", "initializing API support", "host", cfg.Web.APIHost)

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, syscall.SIGINT, syscall.SIGTERM)

	cfgMux := v1.APIMuxConfig{
		Shutdown: shutdown,
		Log:      log,
		Auth:     auth,
		DB:       db,
		Tracer:   tracer,
	}

	apiMux := v1.APIMux(cfgMux, routeAdder)

	api := http.Server{
		Addr:         cfg.Web.APIHost,
		Handler:      apiMux,
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
	// Shutdown

	select {
	case err := <-serverErrors:
		return fmt.Errorf("server error: %w", err)

	case sig := <-shutdown:
		log.Info(ctx, "shutdown", "status", "shutdown started", "signal", sig)
		defer log.Info(ctx, "shutdown", "status", "shutdown complete", "signal", sig)

		ctx, cancel := context.WithTimeout(context.Background(), cfg.Web.ShutdownTimeout)
		defer cancel()

		if err := api.Shutdown(ctx); err != nil {
			api.Close()
			return fmt.Errorf("could not stop server gracefully: %w", err)
		}
	}

	return nil
}

// =============================================================================

// startTracing configure open telemetry to be used with Grafana Tempo.
func startTracing(reporterURI string, serviceName string, probability float64) (*trace.TracerProvider, error) {
	exporter, err := otlptrace.New(
		context.Background(),
		otlptracegrpc.NewClient(
			otlptracegrpc.WithInsecure(), // This should be configurable
			otlptracegrpc.WithEndpoint(reporterURI),
		),
	)
	if err != nil {
		return nil, fmt.Errorf("creating new exporter: %w", err)
	}

	traceProvider := trace.NewTracerProvider(
		trace.WithSampler(trace.TraceIDRatioBased(probability)),
		trace.WithBatcher(exporter,
			trace.WithMaxExportBatchSize(trace.DefaultMaxExportBatchSize),
			trace.WithBatchTimeout(trace.DefaultScheduleDelay*time.Millisecond),
			trace.WithMaxExportBatchSize(trace.DefaultMaxExportBatchSize),
		),
		trace.WithResource(
			resource.NewWithAttributes(
				semconv.SchemaURL,
				semconv.ServiceNameKey.String(serviceName),
			),
		),
	)

	// We must set this provider as the global provider for things to work,
	// but we pass this provider around the program where needed to collect
	// our traces.
	otel.SetTracerProvider(traceProvider)

	// Chooses the HTTP header formats we extract incoming trace contexts from,
	// and the headers we set in outgoing requests.
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{},
		propagation.Baggage{},
	))

	return traceProvider, nil
}
