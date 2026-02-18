package main

import (
	"context"
	"database/sql"
	"flag"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "github.com/lib/pq"

	httpserver "github.com/akonovalovdev/DDD_example/internal/adapters/http"
	"github.com/akonovalovdev/DDD_example/internal/adapters/http/handlers"
	"github.com/akonovalovdev/DDD_example/internal/adapters/repository/postgres"
	"github.com/akonovalovdev/DDD_example/internal/adapters/skinport"
	"github.com/akonovalovdev/DDD_example/internal/application"
	"github.com/akonovalovdev/DDD_example/internal/config"
	"github.com/akonovalovdev/DDD_example/internal/pkg/cache"
)

func main() {
	configPath := flag.String("config", "config/config.yaml", "path to config file")
	flag.Parse()

	logger := setupLogger("info")
	logger.Info("starting application")

	cfg, err := config.Load(*configPath)
	if err != nil {
		logger.Error("failed to load config", slog.Any("error", err))
		os.Exit(1)
	}

	logger = setupLogger(cfg.Log.Level)

	db, err := setupDatabase(cfg.Database)
	if err != nil {
		logger.Error("failed to connect to database", slog.Any("error", err))
		os.Exit(1)
	}
	defer db.Close()
	logger.Info("connected to database")

	itemCache := cache.NewInMemoryCache(time.Minute)
	defer itemCache.Close()

	skinportClient := skinport.NewClient(cfg.Skinport.APIURL, cfg.Skinport.Timeout)
	userRepo := postgres.NewUserRepository(db)
	transactionRepo := postgres.NewTransactionRepository(db)

	itemService := application.NewItemService(skinportClient, itemCache, cfg.Cache.TTL)
	balanceService := application.NewBalanceService(userRepo, transactionRepo)

	// Прогрев кеша при запуске (опционально, не блокирует старт при ошибке)
	logger.Info("warming up items cache...")
	if err := itemService.WarmUp(context.Background()); err != nil {
		logger.Warn("cache warm-up failed, will retry on first request", slog.Any("error", err))
	} else {
		logger.Info("items cache warmed up successfully")
	}

	itemHandler := handlers.NewItemHandler(itemService, logger)
	balanceHandler := handlers.NewBalanceHandler(balanceService, logger)

	server := httpserver.NewServer(
		cfg.Server.Port,
		cfg.Server.ReadTimeout,
		cfg.Server.WriteTimeout,
		itemHandler,
		balanceHandler,
		logger,
	)

	go func() {
		if err := server.Start(); err != nil && err != http.ErrServerClosed {
			logger.Error("server error", slog.Any("error", err))
			os.Exit(1)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), cfg.Server.ShutdownTimeout)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		logger.Error("server forced to shutdown", slog.Any("error", err))
	}

	logger.Info("server stopped")
}

func setupLogger(level string) *slog.Logger {
	var logLevel slog.Level
	switch level {
	case "debug":
		logLevel = slog.LevelDebug
	case "info":
		logLevel = slog.LevelInfo
	case "warn":
		logLevel = slog.LevelWarn
	case "error":
		logLevel = slog.LevelError
	default:
		logLevel = slog.LevelInfo
	}

	handler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: logLevel,
	})

	return slog.New(handler)
}

func setupDatabase(cfg config.DatabaseConfig) (*sql.DB, error) {
	db, err := sql.Open("postgres", cfg.URL)
	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(cfg.MaxOpenConns)
	db.SetMaxIdleConns(cfg.MaxIdleConns)
	db.SetConnMaxLifetime(cfg.ConnMaxLifetime)

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}
