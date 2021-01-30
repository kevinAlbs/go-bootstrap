package main

import (
	"fmt"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"kevinalbs.com/go-bootstrap/util"

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
	doc := bson.D{{"x", 1}}
	coll.InsertOne(ctx, doc)
	fmt.Println("inserted", util.ToJson(doc), "into", coll.Name())
}
