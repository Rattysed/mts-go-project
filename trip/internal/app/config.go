package app

import (
	"io/ioutil"
	"time"

	"gopkg.in/yaml.v3"
)

const (
	AppName                = "auth"
	DefaultServeAddress    = "localhost:9626"
	DefaultShutdownTimeout = 20 * time.Second
	DefaultListenerName    = "trip_listener"
	DefaultToDriverName    = "driver_listener"
	DefaultToClientName    = "client_listener"
	DefaultDSN             = "dsn://"
	DefaultMigrationsDir   = "file://migrations/auth"
)

type AppConfig struct {
	Debug           bool          `yaml:"debug"`
	ShutdownTimeout time.Duration `yaml:"shutdown_timeout"`
}

type DatabaseConfig struct {
	DSN           string `yaml:"dsn"`
	MigrationsDir string `yaml:"migrations_dir"`
}

type KafkaConfig struct {
	ListenerName string `yaml:"listener_name"`
	ToDriverName string `yaml:"to_driver_name"`
	ToClientName string `yaml:"to_client_name"`
}

type Config struct {
	App      AppConfig      `yaml:"app"`
	Database DatabaseConfig `yaml:"database"`
	Kafka    KafkaConfig    `yaml:"kafka"`
}

func NewConfig(fileName string) (*Config, error) {
	data, err := ioutil.ReadFile(fileName)
	if err != nil {
		return nil, err
	}

	cnf := Config{
		App: AppConfig{
			ShutdownTimeout: DefaultShutdownTimeout,
		},
		Database: DatabaseConfig{
			DSN:           DefaultDSN,
			MigrationsDir: DefaultMigrationsDir,
		},

		Kafka: KafkaConfig{
			ListenerName: DefaultListenerName,
			ToDriverName: DefaultToDriverName,
			ToClientName: DefaultToClientName,
		},
	}

	if err := yaml.Unmarshal(data, &cnf); err != nil {
		return nil, err
	}

	return &cnf, nil
}
