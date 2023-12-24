package app

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/jmoiron/sqlx"
	"github.com/pressly/goose/v3"
	kafka "github.com/segmentio/kafka-go"
	"go.uber.org/zap"
)

const driverName = "postgres"

type App struct {
	db *sqlx.DB

	logger *zap.Logger

	listener     *kafka.Reader
	clientWriter *kafka.Writer
	driverWriter *kafka.Writer
}

func getKafkaReader(kafkaURL, topic, groupID string) *kafka.Reader {
	brokers := strings.Split(kafkaURL, ",")
	return kafka.NewReader(kafka.ReaderConfig{
		Brokers:  brokers,
		GroupID:  groupID,
		Topic:    topic,
		MinBytes: 10e3, // 10KB
		MaxBytes: 10e6, // 10MB
	})
}

func newKafkaWriter(kafkaURL, topic string) *kafka.Writer {
	return &kafka.Writer{
		Addr:     kafka.TCP(kafkaURL),
		Topic:    topic,
		Balancer: &kafka.LeastBytes{},
	}
}

func New(config *Config, logger *zap.Logger) (*App, error) {
	db, err := initDB(context.Background(), &config.Database)
	if err != nil {
		return nil, err
	}

	a := &App{
		db:     db,
		logger: logger,
	}
	return a, nil
}

func initDB(ctx context.Context, config *DatabaseConfig) (*sqlx.DB, error) {
	db, err := sqlx.Open(driverName, config.DSN)
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

	fs := os.DirFS(config.MigrationsDir)
	goose.SetBaseFS(fs)

	if err = goose.SetDialect(driverName); err != nil {
		panic(err)
	}

	if err = goose.UpContext(ctx, db.DB, "."); err != nil {
		panic(err)
	}

	return db, nil
}
