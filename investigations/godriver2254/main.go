package main

import (
	"fmt"
	"os"

	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"kevinalbs.com/go-bootstrap/util"
)

func main() {
	var uri string
	if uri = os.Getenv("MONGODB_URI"); uri == "" {
		uri = "mongodb://localhost:27017"
	}
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(uri).SetMonitor(util.CreateMonitor()))
	if err != nil {
		panic(err)
	}
	collection := client.Database("foo").Collection("bar")
	ctx := context.TODO()
	for i := 0; i < 10; i++ {
		_, err = collection.InsertOne(ctx, bson.M{"type": uint8(i), "value": i})
		if err != nil {
			fmt.Println(err)
			return
		}
	}
	err = client.Ping(context.TODO(), nil)
	if err != nil {
		panic(err)
	}
}
