package app

import (
	"context"
	"fmt"
	"os"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/pressly/goose/v3"
	"go.uber.org/zap"

	"trip/internal/config"
	"trip/internal/handlers/listener"
)

const driverName = "postgres"

type App struct {
	logger *zap.Logger

	listener *listener.Listener

	db *sqlx.DB
}

func New(cfg *config.Config, logger *zap.Logger) (*App, error) {
	db, err := initDB(context.Background(), &cfg.Database)
	if err != nil {
		return nil, err
	}

	logger.Info("Db init successfully finished")

	l, err := listener.New(&cfg.Kafka, logger)
	if err != nil {
		return nil, err
	}
	logger.Info("Logger successfully created")
	a := &App{
		logger:   logger,
		listener: l,
		db:       db,
	}
	return a, nil
}

func (a *App) Serve() error {
	a.logger.Info("Starting app serving")
	done := make(chan os.Signal, 1)

	go func() {
		err := a.listener.Serve()
		if err != nil {
			a.logger.Fatal(err.Error())
		}
	}()
	<-done
	return nil
}

func initDB(ctx context.Context, cfg *config.DatabaseConfig) (*sqlx.DB, error) {
	db, err := sqlx.Open(driverName, cfg.DSN)
	if err != nil {
		return nil, fmt.Errorf("unable to connect to database: %w", err)
	}

	db.DB.SetMaxOpenConns(100)  // The default is 0 (unlimited)
	db.DB.SetMaxIdleConns(10)   // defaultMaxIdleConns = 2
	db.DB.SetConnMaxLifetime(0) // 0, connections are reused forever.

	if err = db.PingContext(ctx); err != nil {
		return nil, err
	}

	// migrations

	fs := os.DirFS(cfg.MigrationsDir)
	goose.SetBaseFS(fs)

	if err = goose.SetDialect(driverName); err != nil {
		return nil, err
	}

	if err = goose.UpContext(ctx, db.DB, "."); err != nil {
		return nil, err
	}

	return db, nil
}
