package http

import (
	"context"
	"errors"
	"fmt"
	"github.com/prometheus/client_golang/prometheus"
	"http-metric/internal/app/http/handlers"
	"http-metric/internal/metrics"
	"http-metric/internal/platform/web"
	"log/slog"
	"net/http"
	"strconv"
	"time"
)

type App struct {
	log        *slog.Logger
	httpServer *http.Server
	port       int
	metrics    *metrics.MetricRecorder
	mux        *http.ServeMux
}

func New(log *slog.Logger, port int, timeout time.Duration, metrics *metrics.MetricRecorder) *App {
	server := &http.Server{
		Addr:        fmt.Sprintf(":%d", port),
		ReadTimeout: timeout,
	}

	app := &App{
		log:        log,
		httpServer: server,
		mux:        http.NewServeMux(),
		port:       port,
		metrics:    metrics,
	}

	server.Handler = app.Routes()

	return app
}

type Handler func(w http.ResponseWriter, r *http.Request) (any, error)

func (a *App) Handle(pattern string, handler Handler) {
	const op = "http.Handle"
	metric := a.metrics.HttpMetric
	l := a.log.With(slog.String("op", op), slog.String("endpoint", pattern))

	fn := func(w http.ResponseWriter, r *http.Request) {
		metric.TotalRequests++
		metric.PromTotalRequests.WithLabelValues(r.Method, pattern).Inc()
		defer func() {
			if err := recover(); err != nil {
				metric.TotalServerErrors++
				metric.PromTotalErrors.WithLabelValues("GET", "/ping", strconv.Itoa(http.StatusInternalServerError)).Inc()
				l.Error("panic", "error", err) //nolint:govet
				err := web.RespondError(w)
				if err != nil {
					l.Error("defer response error:", slog.String("error", err.Error())) //nolint:govet
				}
			}
		}()

		l.Info("handling request")
		timer := prometheus.NewTimer(metric.PromRequestDuration.WithLabelValues(r.Method, pattern))
		defer timer.ObserveDuration()
		data, err := handler(w, r)

		if err != nil {
			metric.TotalErrors++
			metric.PromTotalErrors.WithLabelValues(r.Method, pattern, strconv.Itoa(http.StatusBadRequest)).Inc()

			var httpError *handlers.ErrorHTTP
			if err != nil && errors.As(err, &httpError) {
				slog.Info("error code", "code", httpError.Code)
				err = web.Response(w, nil, httpError.Code)
				if err != nil {
					l.Warn("response error:", slog.String("error", err.Error())) //nolint:govet
					panic(fmt.Sprintf("response error: %v", err))
				}
			}
		}
		err = web.Response(w, data, http.StatusOK)
		if err != nil {
			l.Error("response error:", slog.String("error", err.Error())) //nolint:govet
			panic(fmt.Sprintf("response error: %v", err))
		}
	}
	a.mux.HandleFunc(pattern, fn)
}

func (a *App) Run() chan error {
	const op = "http_app.Run"
	logger := a.log.With(slog.String("op", op), slog.Int("port", a.port))

	appErrors := make(chan error)
	logger.Info("starting server")
	logger.Info(fmt.Sprintf("Server is running on http://localhost:%d", a.port))
	go func() {
		appErrors <- a.httpServer.ListenAndServe()
	}()

	logger.Info("http server is running", slog.String("addr", a.httpServer.Addr))

	return appErrors
}

func (a *App) Stop() error {
	const op = "http_app.Stop"
	a.log.With(slog.String("op", op)).Info("stopping HTTP server", slog.Int("port", a.port))

	err := a.httpServer.Shutdown(context.Background())
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (a *App) Routes() http.Handler {
	const op = "http_app.Routes"
	a.log.With(slog.String("op", op)).Info("setting up routes")

	metric := handlers.NewHttpMetric(a.log, a.metrics.HttpMetric)
	a.Handle("/ping", metric.Ping)
	a.Handle("/requests_counter", metric.RequestCounter)

	return a.mux
}
