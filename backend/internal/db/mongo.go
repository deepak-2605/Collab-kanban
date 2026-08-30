package db

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

func Connect(uri, dbName string) *mongo.Database {
	client, err := mongo.Connect(options.Client().ApplyURI(uri))

	if err != nil {
		log.Fatal(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel() // always release the context's resources

	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("connecte to MongoDB")

	return client.Database(dbName)
}
