package main

import (
	"client/internal/app"
	"client/internal/config"
	"context"
	"flag"
	"log"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	var cfgPath string
	flag.StringVar(&cfgPath, "config", "configs/.config.json", "set config path")
	flag.Parse()

	cfg, err := config.NewConfig(cfgPath)
	if err != nil {
		log.Fatal(err)
	}

	a := app.NewApp(cfg)

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	a.Run()

	//TODO tests start

	//TODO tests end

	<-ctx.Done()
	ctx, _ = context.WithTimeout(ctx, 3*time.Second)
	a.Stop(ctx)
}
