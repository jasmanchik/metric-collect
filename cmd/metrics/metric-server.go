package main

import (
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"http-metric/internal/app"
	"http-metric/internal/config"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

func main() {
	cfg := config.MustLoad()
	logger := SetupLogger(cfg.Env)

	go func() {
		http.Handle("/metrics", promhttp.Handler())
		err := http.ListenAndServe(":8081", nil)
		if err != nil {
			logger.Error("can't start http metric server", slog.String("error", err.Error())) //nolint:govet
		}
	}()

	//requestDuration := prometheus.NewHistogramVec(
	//	prometheus.HistogramOpts{
	//		Name:    "http_request_duration_seconds",
	//		Help:    "Duration of requests",
	//		Buckets: prometheus.DefBuckets,
	//	},
	//	[]string{"time"},
	//)
	//err := prometheus.Register(requestDuration)
	//if err != nil {
	//	return
	//}
	//
	//go func() {
	//	for {
	//		start := time.Now()
	//
	//		randTime := rand.Intn(3)
	//		time.Sleep(time.Duration(randTime) * time.Second)
	//
	//		requestDuration.WithLabelValues(fmt.Sprintf("%d", 200)).Observe(time.Since(start).Seconds())
	//	}
	//}()
	//
	//ch := make(chan int, 1)
	//<-ch
	application := app.New(logger, cfg.HTTP.Port, cfg.HTTP.Timeout)
	appErrors := application.HTTPServer.Run()

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, syscall.SIGTERM, syscall.SIGINT)
	l := logger.With(slog.String("op", "main"))
	select {
	case err := <-appErrors:
		l.Error("http server error", slog.String("error", err.Error()))
	case sig := <-shutdown:
		l.Info("stopping application", slog.String("signal", sig.String()))
		err := application.HTTPServer.Stop()
		if err != nil {
			l.Error("http server stop error", slog.String("error", err.Error()))
		}
		l.Info("application stopped")
	}
}

func SetupLogger(env string) *slog.Logger {

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
