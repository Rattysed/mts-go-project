package main

import (
	"flag"
	"fmt"
	logBase "log"
	"trip/internal/app"
	"trip/logger"
)

func getConfigPath() string {
	var configPath string

	flag.StringVar(&configPath, "c", ".config/local.config.yaml", "path to config file")
	flag.Parse()

	return configPath
}

func main() {
	log, err := logger.GetLogger(true)
	if err != nil {
		logBase.Fatal(err)
	}

	path := getConfigPath()
	config, err := app.NewConfig(path)

	if err != nil {
		log.Fatal(err.Error())
	}
	log.Info("Config had opened")

	a, err := app.New(config, log)
	if err != nil {
		log.Fatal(err.Error())
	}
	log.Info("Config had read")

	fmt.Print(a)
}
