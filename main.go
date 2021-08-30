package main

import (
	"context"
	"flag"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"github.com/FreedomCentral/central/cache"
	"github.com/FreedomCentral/central/env"
	"github.com/FreedomCentral/central/queue"
	"github.com/FreedomCentral/central/secret"
	"github.com/FreedomCentral/central/zaplog"
	"github.com/gin-gonic/gin"

	"go.uber.org/zap"
)

const (
	serviceName = "notification-service"
	port        = "server.port"
	dbengName   = "mysql.uri"
	cacheName   = "redis.uri"
	qName       = "rabbitmq.uri"
)

var logger *zap.SugaredLogger

func main() {
	flag.Parse()

	// Setup zap logger.
	logger = zaplog.Setup(serviceName).Sugar()

	// TIP: now use zap.S().Infow(..) for sugared logger
	// or use zap.L().Infow(...) for regular logger
	// These are safe for concurrent use.

	startService()
}

// startService sets up logging, connects to external datatabases, starts http server.
func startService() {
	// IMPORTANT: default environment variable BR_ENV is set via Makefile with values: dev, stg, prod.
	env.Set("dev")
	logger.Infof("Starting %s in %q", serviceName, env.Get())

	// Open secrets, if pass 'UseYAMLPlainText' then will use /config file instead of Vault.
	sec, err := secret.Open(serviceName, secret.UseYAMLPlainText)
	if err != nil {
		logger.Fatalf("Failed to open secrets for %q: %v", serviceName, err)
	}

	//db connect
	dbURI, err := sec.Get(dbengName)
	if err != nil {
		logger.Fatalf("Failed to get secret for %q: %v", dbengName, err)
	}
	mysqlConn, err := connectMySQL(dbURI)
	if err != nil {
		logger.Fatalf("Failed to connect to mysql %q: %v", dbURI, err)
	}

	//redis connect : add more redis instances if required

	redisURI, err := sec.Get(cacheName)
	if err != nil {
		logger.Fatalf("Failed to get secret for %q: %v", cacheName, err)
	}
	redisKV := cache.NewRedisDict(redisURI, 1) // TOOD: instead of number here use const from dbnumbers.go in central.cache package.

	//rmq connect
	rmqURI, err := sec.Get(qName)
	if err != nil {
		logger.Fatalf("Failed to get secret for %q: %v", qName, err)
	}
	q := queue.Connect(rmqURI)

	logger.Info("Connected to DB,Cache, Queue")

	// Create server instance. Main instance shared by all http handlers.
	// HERE we setup any additional variables that will be accessible to all Handlers
	// IMPORTANT variables here should be safe for concurrent use or protected by channels or sync types.

	srv := &Service{
		db:    mysqlConn,
		users: redisKV,
		queue: q,
	}

	router := setupRouter(srv, sec)
	port, err := sec.Get(port)
	if err != nil {
		logger.Fatalf("Failed to get secret for %q: %v", port, err)
	}

	listenAndServe(router, port)
}

func listenAndServe(router *gin.Engine, port string) {
	// Create context that listens for the interrupt signal from the OS.
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	srv := &http.Server{
		Addr:    port,
		Handler: router,
	}

	// Initializing the server in a goroutine so that
	// it won't block the graceful shutdown handling below
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatalf("listen: %s\n", err)
		}
	}()

	// Listen for the interrupt signal.
	<-ctx.Done()

	// Restore default behavior on the interrupt signal and notify user of shutdown.
	stop()
	logger.Info("Shutting down gracefully, press Ctrl+C again to force")

	// The context is used to inform the server it has 5 seconds to finish
	// the request it is currently handling
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		logger.Fatalf("Server forced to shutdown: %v", err)
	}

	logger.Info("Server exiting")
}
