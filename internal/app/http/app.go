package http

import (
	"context"
	"fmt"
	"http-metric/internal/app/http/handlers"
	"http-metric/internal/metrics"
	"http-metric/internal/platform/web"
	"log/slog"
	"net/http"
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

type Handler func(w http.ResponseWriter, r *http.Request) error

func (a *App) Handle(pattern string, handler Handler) {
	const op = "http.Handle"
	l := a.log.With(slog.String("op", op), slog.String("endpoint", pattern))
	fn := func(w http.ResponseWriter, r *http.Request) {
		l.Info("handling request")
		if err := handler(w, r); err != nil {
			l.Error("request error", slog.String("error", err.Error()))
			if err := web.RespondError(w); err != nil {
				l.Error("can't respond with error", slog.String("error", err.Error()))
			}
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

	httpMetric := handlers.NewHttpMetric(a.log, a.metrics.HttpMetric)
	a.Handle("/ping", httpMetric.Ping)
	a.Handle("/requests_counter", httpMetric.RequestCounter)
	a.Handle("/metrics", httpMetric.Metrics)

	return a.mux
}
