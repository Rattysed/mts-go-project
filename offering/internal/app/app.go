package app

import (
	"context"
	"fmt"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
	"log"
	"net/http"
	"offering/internal/config"
	"offering/internal/handlers"
	"offering/internal/manager"
)

type App struct {
	cfg    *config.Config
	server *http.Server
	Logger *zap.Logger
}

func initServer(cfg *config.Config, logger *zap.Logger) http.Handler {
	handler := handlers.NewController(manager.NewManager(cfg, logger), logger)

	router := chi.NewRouter()
	router.Post("/offers", handler.CreateOffer)
	router.Get("/offers/{offerID}", handler.ParseOffer)

	return router
}

func NewApp(cfg *config.Config) *App {
	logger, err := zap.NewProduction()
	if err != nil {
		log.Fatalf("Logger init error. %v", err)
		return nil
	}

	newServer := &http.Server{
		Addr:    cfg.IP + ":" + cfg.Port,
		Handler: initServer(cfg, logger),
	}

	return &App{
		cfg:    cfg,
		server: newServer,
		Logger: logger,
	}
}

func (a *App) Run() {
	a.Logger.Info("Starting app")
	go func() {
		err := a.server.ListenAndServe()
		if err != nil {
			a.Logger.Warn(err.Error())
		}
	}()
}

func (a *App) Stop(ctx context.Context) {
	a.Logger.Info("Closing app")
	fmt.Println(a.server.Shutdown(ctx))
}
