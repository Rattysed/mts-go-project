package admin

import (
	"client/internal/models"
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type MongoDB struct {
	Client                mongo.Client
	database              mongo.Database
	TripCollection        mongo.Collection
	UsersToTripCollection mongo.Collection
}

type UsersToTrip struct {
	User string `json:"user"`
	Trip string `json:"trip"`
}

func NewMongoDB(ctx context.Context, uri, dbName, TripCollectionName string, UserCollectionName string) (*MongoDB, error) {
	clientOptions := options.Client().ApplyURI(uri)
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return nil, err
	}

	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		return nil, err
	}

	db := client.Database(dbName)
	collection := db.Collection(TripCollectionName)
	UsersCollection := db.Collection(UserCollectionName)

	return &MongoDB{Client: *client, database: *db, TripCollection: *collection, UsersToTripCollection: *UsersCollection}, nil
}

func (m *MongoDB) AddTrip(ctx context.Context, trip *models.Trip, userId string) error {
	_, err := m.TripCollection.InsertOne(ctx, trip)
	if err != nil {
		return err
	}

	_, err = m.UsersToTripCollection.InsertOne(ctx, UsersToTrip{
		Trip: trip.ID,
		User: userId,
	})
	if err != nil {
		return err
	}

	return nil
}

func (m *MongoDB) FindTrip(ctx context.Context, tripId string) (*models.Trip, error) {
	fmt.Println()

	filter := bson.D{{"id", tripId}}
	result := models.Trip{}
	err := m.TripCollection.FindOne(ctx, filter).Decode(&result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}
