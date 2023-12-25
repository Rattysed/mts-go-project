package main

import (
	"flag"
	"fmt"
	logBase "log"
	"time"
	"trip/internal/app"
	"trip/internal/config"
	"trip/logger"
)

func getConfigPath() string {
	var configPath string

	flag.StringVar(&configPath, "c", "../.config/local.config.yaml", "path to config file")
	flag.Parse()

	return configPath
}

func main() {
	log, err := logger.GetLogger(true)
	if err != nil {
		logBase.Fatal(err)
	}
	log.Info("Got logger")

	path := getConfigPath()
	config, err := config.NewConfig(path)

	if err != nil {
		log.Fatal(err.Error())
	}
	log.Info("Config had opened")

	a, err := app.New(config, log)
	if err != nil {
		log.Info("Пиздец")
		time.Sleep(5 * time.Second)
		log.Fatal(err.Error())
	}
	log.Info("Created app by config")

	if err = a.Serve(); err != nil {
		log.Fatal(err.Error())
	}

	fmt.Print(a)
}
