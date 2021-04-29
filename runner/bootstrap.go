package main

import (
	"fmt"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"context"
	"time"
)

func doPing(ctx context.Context, client *mongo.Client) error {
	err := client.Ping(ctx, nil)
	if err != nil {
		return fmt.Errorf("db error: %w", err)
	}
	return nil
}

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	opts := options.Client()
	// Uncomment to enable logging of all commands and replies.
	// opts.SetMonitor(util.CreateMonitor())

	client, err := mongo.Connect(ctx, opts)
	if err != nil {
		log.Fatal(err)
	}

	coll := client.Database("test").Collection("coll")
	var result struct {
		Value float64
	}
	filter := bson.D{{"name", "pi"}}
	ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	err = coll.FindOne(ctx, filter).Decode(&result)
	if err == mongo.ErrNoDocuments {
		// Do something when no record was found
		fmt.Println("record does not exist")
	} else {
		log.Fatal(err)
	}
}
