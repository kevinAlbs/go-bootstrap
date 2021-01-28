package main

import (
	"fmt"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/x/bsonx/bsoncore"
	"go.mongodb.org/mongo-driver/mongo"

	"context"
	"time"
)

// DoPing pings
func DoPing(ctx context.Context, client *mongo.Client) error {
	err := client.Ping(ctx, nil)
	if err != nil {
		return fmt.Errorf("db error: %w", err)
	}
	return nil
}

type Test {
	arr bsoncore.Array
}

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Set a command monitor to log all commands and replies.
	// opts.SetMonitor(CreateMonitor())

	client, err := mongo.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}

	coll := client.Database("test").Collection("coll")
	update := bson.D{
		{"$addToSet", bson.D{{"outsideIp", "string"}}},
	}
	filter := bson.D{{"_id", 0}}
	updateResult, err := coll.UpdateOne(context.TODO(), filter, update)
	fmt.Println(updateResult)

	var t Test
	

	// I have a single document: { "_id" : 0, "outsideIp" : [ "a", "b" ] }


}
