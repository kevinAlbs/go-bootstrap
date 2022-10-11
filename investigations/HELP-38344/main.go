package main

import (
	"context"
	"log"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readconcern"
)

func main() {

	uri := "mongodb://localhost:27017/?readConcernLevel=majority"
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(uri))
	if err != nil {
		panic(err)
	}

	defer func() {
		if err := client.Disconnect(context.TODO()); err != nil {
			panic(err)
		}
	}()

	coll := client.Database("test").Collection("poc", options.Collection().SetReadConcern(readconcern.Local()))

	cursor, err := coll.Aggregate(
		context.TODO(),
		mongo.Pipeline{})
	if err != nil {
		log.Fatal(err)
	}
	cursor.Close(context.TODO())
}
