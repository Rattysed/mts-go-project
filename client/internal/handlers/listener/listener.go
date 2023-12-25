package listener

import (
	"client/internal/config"
	"context"
	"fmt"
	"github.com/segmentio/kafka-go"
	"go.uber.org/zap"
	"os"
	"strings"
	"time"
)

type Listener struct {
	logger *zap.Logger
	reader *kafka.Reader
	writer *kafka.Writer
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

func New(cfg *config.KafkaConfig, logger *zap.Logger) (*Listener, error) {
	url := os.Getenv("kafkaURL") // TODO добавить в docker Compose -env
	l := &Listener{
		logger: logger,
		reader: getKafkaReader(url, cfg.ListenerName, ""),
		writer: newKafkaWriter(url, cfg.WriterName),
	}
	return l, nil
}

func (l *Listener) Serve() error {
	l.logger.Info("Started Listener serving")
	done := make(chan os.Signal, 1)
	go func() {
		time.Sleep(5 * time.Second) // даём время на запуск кафки
		for {
			msg, err := l.reader.ReadMessage(context.Background())
			if err != nil {
				l.logger.Fatal(err.Error())
			}
			fmt.Println(msg)
		}
	}()
	<-done

	l.Shutdown()
	return nil
}

func (l *Listener) Shutdown() {
	_, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	l.reader.Close()
	l.writer.Close()
}
