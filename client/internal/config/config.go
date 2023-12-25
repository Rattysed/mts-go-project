package config

import (
	"encoding/json"
	"os"
)

const (
	DefaultServeAddress = "localhost:9626"
	DefaultListenerName = "client_listener"
	DefaultWriterName   = "client_writer"
)

type KafkaConfig struct {
	ListenerName string `yaml:"listener_name"`
	WriterName   string `yaml:"writer_name"`
}

type AppConfig struct {
	IP      string `json:"ip"`
	Port    string `json:"port"`
	Version string `json:"version"`
}

type Config struct {
	App   AppConfig
	Kafka KafkaConfig
}

func NewConfig(filePath string) (*Config, error) {
	// Открытие файла
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	var cfg AppConfig
	// Десериализация
	err = json.Unmarshal(data, &cfg)
	if err != nil {
		return nil, err
	}

	return &Config{
		App: cfg,
		Kafka: KafkaConfig{
			ListenerName: DefaultListenerName,
			WriterName:   DefaultWriterName,
		},
	}, nil
}
