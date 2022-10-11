package main

import (
	"context"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
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

	coll := client.Database("test").Collection("poc")

	opts := options.AggregateOptions{
		Custom: bson.M{"readConcern": bson.M{"level": "local"}},
	}

	cursor, err := coll.Aggregate(
		context.TODO(),
		mongo.Pipeline{}, &opts)
	if err != nil {
		log.Fatal(err)
	}
	cursor.Close(context.TODO())
}
