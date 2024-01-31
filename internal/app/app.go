package app

import (
	"errors"
	"http-metric/internal/server/grpc"
	"http-metric/internal/service/metric"
	"log/slog"
	defHttp "net/http"
)

type WebServer interface {
	Start() error
	Stop() error
}

type App struct {
	server WebServer
	metric *metric.Manager
}

func (a App) Run() {
	err := a.server.Start()
	if err != nil {
		if !errors.Is(err, defHttp.ErrServerClosed) {
			panic(err)
		}
	}
}

func (a App) Stop() error {
	if err := a.server.Stop(); err != nil {
		return err
	}

	return nil
}

func New(logger *slog.Logger, port int) *App {
	metric := metric.New(logger)

	//server := http.New(logger, port, metric)
	server := grpc.New(logger, port, metric)

	return &App{
		server: server,
	}
}
