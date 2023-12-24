package app

import (
	"client/internal/admin"
	"client/internal/config"
	"context"
	"fmt"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
	"log"
	"net/http"
)

type App struct {
	cfg    *config.Config
	server *http.Server
	Logger *zap.Logger
}

func initServer(cfg *config.Config, logger *zap.Logger, admin *admin.MongoDB) http.Handler {
	//handler := handlers.NewController(manager.NewManager(cfg, logger), logger)

	router := chi.NewRouter()
	//router.Post("/offers", handler.CreateOffer)
	//router.Get("/offers/{offerID}", handler.ParseOffer)

	return router
}

func NewApp(cfg *config.Config, admin *admin.MongoDB) *App {
	logger, err := zap.NewProduction()
	if err != nil {
		log.Fatalf("Logger init error. %v", err)
		return nil
	}

	newServer := &http.Server{
		Addr:    cfg.IP + ":" + cfg.Port,
		Handler: initServer(cfg, logger, admin),
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
			fmt.Println(err)
		}
	}()
}

func (a *App) Stop(ctx context.Context) {
	a.Logger.Info("Closing app")
	fmt.Println(a.server.Shutdown(ctx))
}
