package main

import (
	"context"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Replace the uri string with your MongoDB deployment's connection string.
const uri = "mongodb://localhost:27017"

func main() {
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(uri))
	if err != nil {
		panic(err)
	}
	coll := client.Database("db").Collection("coll")
	opts := options.Aggregate()
	//.SetBackground(true) // Proposed API.
	pipeline := mongo.Pipeline{
		{{"$out", bson.D{{"db", "foo"}, {"coll", "bar"}}}},
	}
	cursor, err := coll.Aggregate(context.TODO(), pipeline, opts)
	if err != nil {
		log.Fatal("Aggregate error: ", err)
	}
	var results []bson.D
	err = cursor.All(context.TODO(), &results)
	if err != nil {
		panic(err)
	}
}
