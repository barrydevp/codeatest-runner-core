package main

import (
   "context"
   "fmt"
   "log"
   "time"

   "go.mongodb.org/mongo-driver/bson"
   "go.mongodb.org/mongo-driver/mongo"
   "go.mongodb.org/mongo-driver/mongo/options"
   // "go.mongodb.org/mongo-driver/mongo/readpref"
)

type User struct {
    ExternalId string `json:"external_id"`
    Name string `json:"name"`
    Group string `json:"group"`
    CreatedAt string `json:"created_at"`
}

func main() {
    client, err := mongo.NewClient(options.Client().ApplyURI("mongodb+srv://root:MYdXPr7PEDwwf9hG@code-and-test.hug2u.mongodb.net/code-and-test"))

    if err != nil {
        log.Fatal(err)
    }

    fmt.Println("STARTING")

    ctx, _ := context.WithTimeout(context.Background(), 1 * time.Second)

    err = client.Connect(ctx)
    if err != nil {
            log.Fatal(err)
    }

    defer client.Disconnect(ctx)

    UserCollection := client.Database("code-and-test").Collection("users")

    cur, err := UserCollection.Find(ctx, bson.D{})

    if err != nil {
        log.Fatal(err)
    }

    defer cur.Close(ctx)

    for cur.Next(ctx) {
        var user User
        err := cur.Decode(&user)

        if err != nil {
            log.Fatal(err)
        }

        fmt.Println(user.Name)
    }

}
