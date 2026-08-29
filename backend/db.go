package main

import (
	"context"
	"log"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

func connectMongo() *mongo.Client {
	connUri := os.Getenv("MONGO_URI")
	client, err := mongo.Connect(options.Client().ApplyURI(connUri))

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

	return client

}
