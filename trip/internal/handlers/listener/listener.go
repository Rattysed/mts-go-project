package listener

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/segmentio/kafka-go"
	"go.uber.org/zap"

	"trip/internal/config"
	"trip/models"
)

type Listener struct {
	Logger       *zap.Logger
	Reader       *kafka.Reader
	DriverWriter *kafka.Writer
	ClientWriter *kafka.Writer
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
		Logger:       logger,
		Reader:       getKafkaReader(url, cfg.ListenerName, ""),
		DriverWriter: newKafkaWriter(url, cfg.ToDriverName),
		ClientWriter: newKafkaWriter(url, cfg.ToClientName),
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
			case "trip.command.accept":
				l.OnAccept(values)
			case "trip.command.cancel":
				l.OnCancel(values)
			case "trip.command.create":
				l.OnCreate(values)
			case "trip.command.end":
				l.OnEnd(values)
			case "trip.command.start":
				l.OnStart(values)
			}
			time.Sleep(100 * time.Millisecond)
		}
	}()
	<-done

	l.Shutdown()
	return nil
}

func (l *Listener) Shutdown() {
	_, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	l.Reader.Close()
	l.ClientWriter.Close()
	l.DriverWriter.Close()
}

func (l *Listener) OnAccept(values models.Event) {

}

func (l *Listener) OnCancel(values models.Event) {

}

func (l *Listener) OnCreate(values models.Event) {
	offerId := values.Data["offer_id"]
	resp, err := http.Get("http://offering:8080/offers/" + offerId)
	if err != nil {
		l.Logger.Error("Error occurred while requesting offer: " + err.Error())
		return
	}
	bytes, err := io.ReadAll(resp.Body)
	if err != nil {
		l.Logger.Error("Error occurred while requesting offer: " + err.Error())
		return
	}
	var offer models.Offer
	err = json.Unmarshal(bytes, &offer)
	if err != nil {
		l.Logger.Error("Something wrong wih requested offer: " + err.Error())
		return
	}
	l.Logger.Info("Parsed offer with id: " + offer.Id)

}

func (l *Listener) OnEnd(values models.Event) {

}

func (l *Listener) OnStart(values models.Event) {

}
