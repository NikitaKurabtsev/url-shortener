package main

import (
	"fmt"
	"github.com/NikitaKurabtsev/url-shortener/internal/config"
	"log/slog"
	"os"
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
	log.Debug("debug message")

	// TODO: init storage: sqlite

	// TODO: init router: chi, chi render

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
