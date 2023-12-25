package listener

import (
	"client/internal/admin"
	"client/internal/config"
	"client/internal/models"
	"context"
	"encoding/json"
	"github.com/segmentio/kafka-go"
	"go.uber.org/zap"
	"os"
	"strings"
	"time"
)

type Listener struct {
	Logger *zap.Logger
	Reader *kafka.Reader
	Writer *kafka.Writer
	dbc    *admin.DBController
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

func New(cfg *config.KafkaConfig, logger *zap.Logger, dbc *admin.DBController) (*Listener, error) {
	url := os.Getenv("kafkaURL") // TODO добавить в docker Compose -env
	l := &Listener{
		Logger: logger,
		Reader: getKafkaReader(url, cfg.ListenerName, ""),
		Writer: newKafkaWriter(url, cfg.WriterName),
		dbc:    dbc,
	}
	return l, nil
}

func (l *Listener) Serve() error {
	l.Logger.Info("Started Listener serving")
	done := make(chan os.Signal, 1)
	go func() {
		time.Sleep(5 * time.Second) // даём время на запуск кафки
		l.Logger.Info("Kafka listener awakened")
		for {
			msg, err := l.Reader.ReadMessage(context.Background())
			if err != nil {
				l.Logger.Error("Error occurred after reading message " + err.Error())
			}
			var values models.Event
			if string(msg.Value) != "" {
				if err := json.Unmarshal(msg.Value, &values); err != nil {
					l.Logger.Error("Failed to unmarshal json " + err.Error())
				}
			}
			l.Logger.Info(values.Type)
			switch values.Type {
			case "trip.event.accepted":
				l.OnAccept(values)
			case "trip.event.canceled":
				l.OnCancel(values)
			case "trip.event.created":
				l.OnCreate(values)
			case "trip.event.ended":
				l.OnEnd(values)
			case "trip.event.started":
				l.OnStart(values)
			}
			time.Sleep(100 * time.Millisecond)
		}
	}()
	<-done

	l.Shutdown()
	return nil
}

func (l *Listener) OnAccept(values models.Event) {

}

func (l *Listener) OnCancel(values models.Event) {

}

func (l *Listener) OnCreate(values models.Event) {

}

func (l *Listener) OnEnd(values models.Event) {

}

func (l *Listener) OnStart(values models.Event) {

}

func (l *Listener) Shutdown() {
	_, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	l.Reader.Close()
	l.Writer.Close()
}
