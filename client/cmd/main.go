package main

import (
	"client/internal/admin"
	"client/internal/app"
	"client/internal/config"
	"context"
	"flag"
	"go.mongodb.org/mongo-driver/mongo"
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

	ctx := context.Background()
	uri := "mongodb://localhost:27017"
	dbName := "test"
	tripCollectionName := "trips"
	UsersToTripsCollectionName := "usersToTrips"
	mongoDB, err := admin.NewMongoDB(ctx, uri, dbName, tripCollectionName, UsersToTripsCollectionName)
	if err != nil {
		log.Fatal(err)
	}
	defer func(client *mongo.Client, ctx context.Context) {
		err := client.Disconnect(ctx)
		if err != nil {

		}
	}(&mongoDB.Client, ctx)

	a := app.NewApp(cfg, mongoDB)

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	a.Run()

	//TODO tests start

	//TODO tests end

	<-ctx.Done()
	ctx, _ = context.WithTimeout(ctx, 3*time.Second)
	a.Stop(ctx)
}
