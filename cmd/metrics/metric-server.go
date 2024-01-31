package main

import (
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"http-metric/internal/app"
	"http-metric/internal/config"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
)

func main() {
	cfg := config.MustLoad()
	logger := SetupLogger(cfg.LogLevel)

	go func() {
		http.Handle("/metrics", promhttp.Handler())
		err := http.ListenAndServe(":"+strconv.Itoa(cfg.DebugPort), nil)
		if err != nil {
			logger.Error("can't start http metric server", slog.String("error", err.Error())) //nolint:govet
		}
	}()

	application := app.New(logger, cfg.HTTP.Port)
	go application.Run()

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, syscall.SIGTERM, syscall.SIGINT)
	l := logger.With(slog.String("op", "main"))
	select {
	case sig := <-shutdown:
		l.Info("stopping application", slog.String("signal", sig.String()))
		err := application.Stop()
		if err != nil {
			l.Error("http server stop error", slog.String("error", err.Error()))
		}
		l.Info("application stopped")
	}
}

func SetupLogger(level slog.Level) *slog.Logger {

	var log *slog.Logger

	log = slog.New(
		slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: level}),
	)

	return log
}
