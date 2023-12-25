package admin

import (
	"client/internal/config"
	"client/internal/models"
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"
	"time"
)

type DBController struct {
	cfg         *config.Config
	mongoClient *mongo.Client
	Logger      *zap.Logger
}

func NewDBController(cfg *config.Config, mongoClient *mongo.Client, logger *zap.Logger) *DBController {
	return &DBController{
		cfg:         cfg,
		mongoClient: mongoClient,
		Logger:      logger,
	}
}

func (dbc *DBController) AddTrip(trip models.Trip) {
	collection := dbc.mongoClient.Database("my_mongo").Collection("trips")
	ctx := context.Background()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_, err := collection.InsertOne(ctx, trip)
	if err != nil {
		dbc.Logger.Warn("Insert error")
	}
}

func (dbc *DBController) GetTrip(tripID string, userID string) *models.Trip {
	collection := dbc.mongoClient.Database("my_mongo").Collection("trips")
	filter := bson.M{"id": tripID, "clientid": userID}
	ctx := context.Background()

	var trip models.Trip
	err := collection.FindOne(ctx, filter).Decode(&trip)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			dbc.Logger.Warn("Trip not found")
			return nil
		}
		dbc.Logger.Warn("Internal server error")
		return nil
	}

	return &trip
}

func (dbc *DBController) ListTrips(userID string) []models.Trip {
	collection := dbc.mongoClient.Database("my_mongo").Collection("trips")
	filter := bson.M{"clientid": userID}
	ctx := context.Background()

	cursor, err := collection.Find(ctx, filter)

	if err != nil {
		dbc.Logger.Warn("Internal server error")
		return nil
	}
	defer func(cur *mongo.Cursor, ctx context.Context) {
		err := cur.Close(ctx)
		if err != nil {

		}
	}(cursor, ctx)

	var trips []models.Trip
	for cursor.Next(ctx) {
		var trip models.Trip
		err := cursor.Decode(&trip)
		if err != nil {
			dbc.Logger.Warn("Decoding error")
		}
		trips = append(trips, trip)
	}

	return trips
}

func (dbc *DBController) ChangeStatus(tripID string, userID string, status string) error {
	collection := dbc.mongoClient.Database("my_mongo").Collection("trips")
	filter := bson.M{"id": tripID, "clientid": userID}
	ctx := context.Background()
	update := bson.M{
		"$set": bson.M{
			"status": status,
		},
	}

	updateResult, err := collection.UpdateOne(ctx, filter, update)
	if err != nil {
		dbc.Logger.Warn("Failed to update status")
		return err
	}

	if updateResult.ModifiedCount == 0 {
		dbc.Logger.Warn("Trip not found or not authorized")
		return err
	}
	return nil
}
