package connections

import (
	"context"
	"log"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

var MongoURI string
var MongoClient *mongo.Client
var MongoDatabase *mongo.Database

func init() {
	MongoURI = os.Getenv("MONGODB_URI")

	if MongoURI == "" {
		log.Fatal("MONGODB_URI is required.")
	}

	client, err := mongo.NewClient(options.Client().ApplyURI(MongoURI))

	if err != nil {
		log.Fatal(err)
	}

	MongoClient = client

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

	defer cancel()

	err = client.Connect(ctx)
	err = client.Ping(ctx, readpref.Primary())

	err = client.Ping(context.TODO(), readpref.Primary())
	log.Println("PING AGAIN")

	if err != nil {
		log.Fatal("Couldn't connect to the database", err)
	} else {
		log.Println("Connected!")
	}

	MongoDatabase = client.Database("code-and-test")
}

func GetModel(collection string) *mongo.Collection {
	return MongoDatabase.Collection(collection)
}

func CloseMongoClient() {
	err := MongoClient.Disconnect(context.TODO())

	if err != nil {
		log.Fatal(err)
	}
	log.Println("Connection to MongoDB closed.")
}

func Ping() (err error) {
	err = MongoClient.Ping(context.TODO(), readpref.Primary())

	if err != nil {
		return
	}

	return
}
