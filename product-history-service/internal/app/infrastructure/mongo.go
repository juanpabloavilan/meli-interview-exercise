package infrastructure

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoDB struct {
	mongo.Database
}

func NewMongoDB(ctx context.Context, mongoConnectionString string, databaseName string) (*mongo.Database, error) {
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoConnectionString))
	if err != nil {
		return nil, err
	}

	return client.Database(databaseName), nil
}
