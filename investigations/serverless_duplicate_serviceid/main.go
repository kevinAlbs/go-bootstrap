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

func toJson(v interface{}) string {
	bytes, err := bson.MarshalExtJSON(v, true, false)
	if err != nil {
		log.Fatal(err)
	}
	return string(bytes)
}

const verbose bool = false

var findServiceIDs []string

func logCommandStarted(context context.Context, event *event.CommandStartedEvent) {
	if event.CommandName == "find" {
		findServiceIDs = append(findServiceIDs, event.ServiceID.Hex())
	}
	if !verbose {
		return
	}
	fmt.Printf("=> CommandStarted: %v with serviceID: %v\n", event.CommandName, event.ServiceID)
}

func logPoolEvent(event *event.PoolEvent) {
	if !verbose {
		return
	}
	fmt.Printf("=> PoolEvent: %v with connectionID: %v\n", event.Type, event.ConnectionID)
}

func main() {
	var uri string
	if uri = os.Getenv("MONGODB_URI"); uri == "" {
		log.Fatal("You must set your 'MONGODB_URI' environmental variable.")
	}

	client, err := mongo.Connect(context.TODO(), options.Client().
		ApplyURI(uri).
		SetMonitor(&event.CommandMonitor{logCommandStarted, nil, nil}).
		SetPoolMonitor(&event.PoolMonitor{logPoolEvent}))
	if err != nil {
		panic(err)
	}
	defer func() {
		if err = client.Disconnect(context.TODO()); err != nil {
			panic(err)
		}
	}()

	coll := client.Database("test").Collection("test")
	coll.Drop(context.TODO())
	coll.InsertOne(context.TODO(), bson.D{{"x", 1}})
	coll.InsertOne(context.TODO(), bson.D{{"x", 1}})

	cursor1, err := coll.Find(context.TODO(), bson.D{}, options.Find().SetBatchSize(1))
	if err != nil {
		panic(err)
	}
	defer cursor1.Close(context.TODO())

	cursor2, err := coll.Find(context.TODO(), bson.D{}, options.Find().SetBatchSize(1))
	if err != nil {
		panic(err)
	}
	defer cursor2.Close(context.TODO())

	fmt.Println("captured the following serviceIDs on the 'find' commands")
	for _, serviceID := range findServiceIDs {
		fmt.Println(serviceID)
	}
	// Sample output:
	// captured the following serviceIDs on the 'find' commands
	// 6127bde11b3bb25a4ac842f8
	// 6127bde11b3bb25a4ac842f8
}
