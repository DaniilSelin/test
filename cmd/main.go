package main

import (
	"fmt"
	"os"
	"os/signal"
	"context"
	"time"
	"io"
	"net/http"

	"github.com/jackc/pgx/v5/pgxpool"
    "go.uber.org/zap"

	"quotebook/config"
	"quotebook/internal/logger"
	"quotebook/internal/database"
	"quotebook/internal/repository"
	"quotebook/internal/service"
	"quotebook/internal/transport/http/api"
)

func main() {
	ctx := context.Background()
	ctx, cancel := signal.NotifyContext(ctx, os.Interrupt)
	defer cancel()

	srv, dbPool, logBase, err := run(ctx, os.Stdout, os.Args);
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
	defer dbPool.Close()

	<-ctx.Done()
	logBase.Info(ctx, "Shutdown signal received")

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer shutdownCancel()

    if err := srv.Shutdown(shutdownCtx); err != nil {
        logBase.Error(ctx, "Server shutdown failed", zap.Error(err))
    }
    logBase.Info(ctx, "Server exited gracefully")
}

func run(ctx context.Context, w io.Writer, args []string) (*http.Server, *pgxpool.Pool, *logger.Logger, error) {
    // Конфиг и логгер
    cfg, err := config.LoadConfig("config/config.yml")
    if err != nil {
        return nil, nil, nil, err
    }
    logBase, err := logger.New(cfg)
    if err != nil {
        return nil, nil, nil, err
    }
    ctx = logger.CtxWWithLogger(ctx, logBase)

    // Подключение к БД и миграции
    dbPool, err := database.Connect(ctx, cfg)
    if err != nil {
        return nil, nil, nil, err
    }
    if err := database.RunMigrations(ctx, cfg, dbPool); err != nil {
        return nil, nil, nil, err
    }

    // Репозиторий и quote-сервис
    repo := repository.NewQuoteRepository(dbPool, cfg)
    qSrv := service.NewQuoteService(cfg, repo)

    // роутер
    handler := api.NewHandler(logBase, cfg, qSrv)
    router := api.NewRouter(handler)

    // HTTP-сервер
    addr := fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port)
    srv := &http.Server{
        Addr:    addr,
        Handler: router,
    }

    logBase.Info(ctx, "starting server", zap.String("addr", addr))
    go func() {
        if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
            logBase.Error(ctx, "ListenAndServe failed", zap.Error(err))
        }
    }()

    return srv, dbPool, logBase, nil
}