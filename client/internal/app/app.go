package app

import (
	"client/internal/admin"
	"client/internal/config"
	"client/internal/handlers"
	"client/internal/handlers/listener"
	"context"
	"fmt"
	"github.com/go-chi/chi/v5"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"
	"log"
	"net/http"
	"time"
)

type App struct {
	cfg      *config.Config
	server   *http.Server
	Logger   *zap.Logger
	client   *mongo.Client
	listener *listener.Listener
}

func initServer(cfg *config.Config, logger *zap.Logger, client *mongo.Client) http.Handler {
	handler := handlers.NewController(admin.NewDBController(cfg, client, logger))

	router := chi.NewRouter()
	router.Post("/trips", handler.AddTrip)
	router.Get("/trips/{trip_id}", handler.GetTrip)
	router.Get("/trips", handler.ListTrips)
	router.Post("/trips/{trip_id}/cancel", handler.CancelTrip)

	return router
}

func NewApp(cfg *config.Config) *App {
	logger, err := zap.NewProduction()
	if err != nil {
		log.Fatalf("Logger init error. %v", err)
		return nil
	}

	client, err := initMongo()
	if err != nil {
		logger.Warn("Init mongo error")
	}

	l, err := listener.New(&cfg.Kafka, logger, admin.NewDBController(cfg, client, logger))

	if err != nil {
		log.Fatal("Failed to create listener")
	}

	newServer := &http.Server{
		Addr:    cfg.App.IP + ":" + cfg.App.Port,
		Handler: initServer(cfg, logger, client),
	}

	return &App{
		cfg:      cfg,
		server:   newServer,
		Logger:   logger,
		listener: l,
		client:   client,
	}
}

func initMongo() (*mongo.Client, error) {
	log.Println("Connecting to mongo")

	client, err := mongo.NewClient(options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		log.Fatal(err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Connected to MongoDB")
	return client, nil
}

func (a *App) Run() {
	a.Logger.Info("Starting app client")
	go func() {
		err := a.server.ListenAndServe()
		if err != nil {
			a.Logger.Fatal(err.Error())
		}
	}()

	go func() {
		err := a.listener.Serve()
		if err != nil {
			a.Logger.Fatal(err.Error())
		}
	}()
}

func (a *App) Stop(ctx context.Context) {
	a.Logger.Info("Closing app")
	fmt.Println(a.server.Shutdown(ctx))
}
