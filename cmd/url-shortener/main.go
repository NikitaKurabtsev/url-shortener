package main

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/NikitaKurabtsev/url-shortener/internal/config"
	"github.com/NikitaKurabtsev/url-shortener/internal/lib/logger/sl"
	"github.com/NikitaKurabtsev/url-shortener/internal/storage/sqlite"
	"github.com/go-chi/chi/v5"
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

	// TODO: init router: chi, chi render
	router := chi.NewRouter()

	// TODO: init run server: net/http
}

func setupLogger(env string) *slog.Logger {
	var log *slog.Logger

	switch env {
	case envLocal:
		log = slog.New(
			slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
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
