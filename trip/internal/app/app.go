package app

import (
	"context"
	"fmt"
	"github.com/golang-migrate/migrate/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"os"
	"trip/internal/handlers/listener"

	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"go.uber.org/zap"

	"trip/internal/config"
)

type App struct {
	logger *zap.Logger

	listener *listener.Listener
}

func New(cfg *config.Config, logger *zap.Logger) (*App, error) {
	//db, err := initDB(context.Background(), &cfg.Database)
	//if err != nil {
	//	return nil, err
	//}

	logger.Info("Db init successfully finished")

	l, err := listener.New(&cfg.Kafka, logger)
	if err != nil {
		return nil, err
	}
	logger.Info("Logger successfully created")
	a := &App{
		logger:   logger,
		listener: l,
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

func initDB(ctx context.Context, cfg *config.DatabaseConfig) (*pgxpool.Pool, error) {
	pgxConfig, err := pgxpool.ParseConfig(cfg.DSN)
	if err != nil {
		return nil, err
	}

	pool, err := pgxpool.ConnectConfig(ctx, pgxConfig)
	if err != nil {
		return nil, fmt.Errorf("unable to connect to database: %w", err)
	}

	// migrations

	m, err := migrate.New(cfg.MigrationsDir, cfg.DSN)
	if err != nil {
		return nil, err
	}

	if err := m.Down(); err != nil && err != migrate.ErrNoChange {
		return nil, err
	}

	if err := m.Up(); err != nil {
		return nil, err
	}

	return pool, nil
}
