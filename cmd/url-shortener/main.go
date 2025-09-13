package main

import (
	"fmt"
	"log/slog"
	"net/http"
	"os"

	"github.com/NikitaKurabtsev/url-shortener/internal/config"
	"github.com/NikitaKurabtsev/url-shortener/internal/http-server/handlers/urls/save"
	mwLogger "github.com/NikitaKurabtsev/url-shortener/internal/http-server/middleware/logger"
	"github.com/NikitaKurabtsev/url-shortener/internal/lib/logger/handlers/slogpretty"
	"github.com/NikitaKurabtsev/url-shortener/internal/lib/logger/sl"
	"github.com/NikitaKurabtsev/url-shortener/internal/storage/sqlite"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

func main() {
	cfg := config.InitConfig()
	fmt.Println(cfg)

	log := setupLogger(envDev)
	log.Info("starting url-shortener", slog.String("env", cfg.Env))

	// TODO: init storage: sqlite
	storage, err := sqlite.NewStorage(cfg.StoragePath)
	if err != nil {
		log.Error("error occurred while initializing storage", sl.Err(err))
		os.Exit(1)
	}

	_ = storage

	// TODO: init router: chi, chi render
	router := chi.NewRouter()
	router.Use(middleware.RequestID)
	router.Use(mwLogger.New(log))
	router.Use(middleware.Recoverer)
	router.Use(middleware.URLFormat)

	router.Post("/url", save.New(log, storage))

	// TODO: init run server: net/http
	log.Info("start server", slog.String("address", cfg.Address))

	srv := http.Server{
		Addr:         cfg.Address,
		Handler:      router,
		ReadTimeout:  cfg.HTTPServer.Timeout,
		WriteTimeout: cfg.HTTPServer.Timeout,
		IdleTimeout:  cfg.HTTPServer.IdleTimeout,
	}

	if err := srv.ListenAndServe(); err != nil {
		log.Error("failed to start server")
	}

	log.Error("server stopped")

}

func setupLogger(env string) *slog.Logger {
	var log *slog.Logger

	switch env {
	case envLocal:
		log = setupPrettySlog()
	case envDev:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case envProd:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
		)
	}

	return log
}

func setupPrettySlog() *slog.Logger {
	opts := slogpretty.PrettyHandlerOptions{
		SlogOpts: &slog.HandlerOptions{
			Level: slog.LevelDebug,
		},
	}

	handler := opts.NewPrettyHandler(os.Stdout)

	return slog.New(handler)
}
