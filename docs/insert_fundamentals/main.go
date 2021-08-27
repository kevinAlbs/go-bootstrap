package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	var uri string
	if uri = os.Getenv("MONGODB_URI"); uri == "" {
		log.Fatal("You must set your 'MONGODB_URI' environmental variable. See\n\t https://docs.mongodb.com/drivers/go/current/usage-examples/")
	}

	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(uri))
	if err != nil {
		panic(err)
	}
	defer func() {
		if err = client.Disconnect(context.TODO()); err != nil {
			panic(err)
		}
	}()

	coll := client.Database("foo").Collection("bar")
	docs := []interface{}{
		bson.D{{"_id", 1}, {"country", "Tanzania"}},
		bson.D{{"_id", 2}, {"country", "Lithuania"}},
		bson.D{{"_id", 1}, {"country", "Vietnam"}},
		bson.D{{"_id", 3}, {"country", "Argentina"}},
	}

	opts := options.InsertMany().SetOrdered(false)
	result, err := coll.InsertMany(context.TODO(), docs, opts)
	list_ids := result.InsertedIDs
	if err != nil {
		fmt.Printf("A bulk write error occurred, but %v documents were still inserted.\n", len(list_ids))
	}
}
