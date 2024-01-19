package http

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"http-metric/internal/server/http/handlers"
	"http-metric/internal/server/http/middleware"
	"http-metric/internal/service/metric"
	"log/slog"
	"net/http"
)

type ServerHttp struct {
	log     *slog.Logger
	server  *http.Server
	port    int
	metrics *metric.Manager
	router  *http.ServeMux
}

func New(log *slog.Logger, port int, metrics *metric.Manager) *ServerHttp {
	server := &http.Server{
		Addr: fmt.Sprintf(":%d", port),
	}

	app := &ServerHttp{
		log:     log,
		server:  server,
		port:    port,
		metrics: metrics,
	}

	server.Handler = app.Routes()

	return app
}

func (a *ServerHttp) Start() error {
	const op = "httpServer.Start"
	logger := a.log.With(slog.String("op", op), slog.Int("port", a.port))

	logger.Info("starting server")

	err := a.server.ListenAndServe()
	if err != nil {
		return err
	}

	logger.Info(fmt.Sprintf("http server is running on %s", slog.String("addr", a.server.Addr)))

	return nil
}

func (a *ServerHttp) Stop() error {
	const op = "httpServer.Stop"
	a.log.With(slog.String("op", op)).Info("stopping HTTP server", slog.Int("port", a.port))

	err := a.server.Shutdown(context.Background())
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (a *ServerHttp) Routes() http.Handler {
	const op = "httpServer.Routes"
	l := a.log.With(slog.String("op", op))
	l.Info("setting up routes")

	r := gin.New()
	r.Use(middleware.LogRequest(a.log), middleware.Metric(a.metrics), middleware.Recovery(a.log))

	m := handlers.NewHttpMetricHandler(a.log, a.metrics)
	r.GET("/ping", m.Ping)
	r.GET("/requests_counter", m.RequestCounter)

	l.Info("set up routes done")
	return r
}
