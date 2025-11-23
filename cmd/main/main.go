package main

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/hello-larin/avito-2025-autumn/internal/base"
	database "github.com/hello-larin/avito-2025-autumn/internal/db"
	prHttp "github.com/hello-larin/avito-2025-autumn/internal/delivery/pr"
	teamHttp "github.com/hello-larin/avito-2025-autumn/internal/delivery/team"
	userHttp "github.com/hello-larin/avito-2025-autumn/internal/delivery/user"
	prRepository "github.com/hello-larin/avito-2025-autumn/internal/repository/pr"
	teamRepository "github.com/hello-larin/avito-2025-autumn/internal/repository/team"
	userRepository "github.com/hello-larin/avito-2025-autumn/internal/repository/user"
	prUsecase "github.com/hello-larin/avito-2025-autumn/internal/usecase/pr"
	teamUsecase "github.com/hello-larin/avito-2025-autumn/internal/usecase/team"
	userUsecase "github.com/hello-larin/avito-2025-autumn/internal/usecase/user"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"

	"github.com/go-playground/validator/v10"
)

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		AddSource: true, // ← вот эта опция включает вывод файла и строки
		Level:     slog.LevelDebug,
	}))

	// Устанавливаем как глобальный логгер (опционально)
	slog.SetDefault(logger)

	pg, err := database.NewPostgresPool()
	if err != nil {
		slog.Warn("Failed to initialize PostgreSQL pool", "error", err)
	}
	defer pg.Close()
	validate := validator.New(validator.WithRequiredStructEnabled())
	prRepository := prRepository.New(pg)
	teamRepository := teamRepository.New(pg)
	userRepository := userRepository.New(pg)
	prUsecase := prUsecase.New(prRepository, userRepository)
	teamUsecase := teamUsecase.New(teamRepository, userRepository)
	userUsecase := userUsecase.New(userRepository, prRepository)
	prHTTP := prHttp.New(prUsecase, validate)
	teamHTTP := teamHttp.New(teamUsecase, validate)
	userHTTP := userHttp.New(userUsecase, validate)

	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(60 * time.Second))
	r.Use(middleware.StripSlashes)

	r.Get("/health", base.HealthHandler(pg))

	prHTTP.RegisterRoutes(r)
	teamHTTP.RegisterRoutes(r)
	userHTTP.RegisterRoutes(r)

	slog.Info("server configured")

	port := ":" + os.Getenv("APP_PORT")
	if port == ":" {
		port = "8080"
	}

	srv := &http.Server{
		Addr:              port,
		Handler:           r,
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       10 * time.Second,
		WriteTimeout:      15 * time.Second,
		IdleTimeout:       60 * time.Second,
		MaxHeaderBytes:    1 << 20,
	}

	go func() {
		// service connections
		if err = srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			slog.Error("listen", "error", err)
		}
	}()

	quit := make(chan os.Signal, 1)

	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	slog.Info("Shutdown Server ...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err = srv.Shutdown(ctx); err != nil {
		slog.Error("Server Shutdown:", "error", err)
		os.Exit(1)
	}

	<-ctx.Done()
	slog.Info("timeout of 5 seconds.")
	slog.Info("Server exiting")
}
