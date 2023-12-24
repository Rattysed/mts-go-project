package config

import (
	"encoding/json"
	"os"
)

type Config struct {
	IP      string `json:"ip"`
	Port    string `json:"port"`
	Version string `json:"version"`
	JWT     string `json:"jwt"`
}

func NewConfig(filePath string) (*Config, error) {
	// Открытие файла
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	var cfg Config
	// Десериализация
	err = json.Unmarshal(data, &cfg)
	if err != nil {
		return nil, err
	}

	return &cfg, nil
}
