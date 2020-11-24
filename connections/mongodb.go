package connections

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

var client *mongo.Client
var database *mongo.Database

func init() {
	client, err := mongo.NewClient(options.Client().ApplyURI("mongodb+srv://root:MYdXPr7PEDwwf9hG@code-and-test.hug2u.mongodb.net/code-and-test"))

	if err != nil {
		log.Fatal(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

	defer cancel()

	err = client.Connect(ctx)
	err = client.Ping(ctx, readpref.Primary())

	if err != nil {
		log.Fatal("Couldn't connect to the database", err)
	} else {
		log.Println("Connected!")
	}

	database = client.Database("code-and-test")
}

func GetModel(collection string) *mongo.Collection {
	return database.Collection(collection)
}
