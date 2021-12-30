// Run with:
// export MONGODB_URI="<URI to a serverless instance>"
// go run ./investigations/serverless_listcollections_batchsize

package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/event"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func logCommandStarted(context context.Context, event *event.CommandStartedEvent) {
	fmt.Println("=> CommandStarted:", event.CommandName, ToJson(event.Command))
}

func logCommandSucceeded(context context.Context, event *event.CommandSucceededEvent) {
	fmt.Println("<= CommandSucceeded:", event.CommandName, ToJson(event.Reply))
}

func logCommandFailed(context context.Context, event *event.CommandFailedEvent) {
	fmt.Println("<= CommandFailed:", event.CommandName, event.Failure)
}

func ToJson(v interface{}) string {
	bytes, err := bson.MarshalExtJSON(v, true, false)
	if err != nil {
		log.Fatal(err)
	}
	return string(bytes)
}

func CreateMonitor() *event.CommandMonitor {
	return &event.CommandMonitor{logCommandStarted, logCommandSucceeded, logCommandFailed}
}

func createCollection(db *mongo.Database, name string) {
	_, err := db.Collection(name).InsertOne(context.TODO(), bson.D{{"x", 1}})
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	uri := os.Getenv("MONGODB_URI")
	opts := options.Client().ApplyURI(uri)
	opts.SetMonitor(CreateMonitor())
	client, err := mongo.Connect(context.TODO(), opts)
	if err != nil {
		log.Fatal(err)
	}

	defer func() {
		client.Disconnect(context.TODO())
	}()

	db := client.Database("foo")
	createCollection(db, "coll0")
	createCollection(db, "coll1")
	createCollection(db, "coll2")

	lcopts := options.ListCollections().SetBatchSize(1)
	cursor, err := db.ListCollections(context.TODO(), bson.D{}, lcopts)
	if err != nil {
		log.Fatal(err)
	}

	var results []bson.D
	err = cursor.All(context.TODO(), &results)
	if err != nil {
		log.Fatal(err)
	}
}

// The expectation is to receive one result on the initial listCollections command.
