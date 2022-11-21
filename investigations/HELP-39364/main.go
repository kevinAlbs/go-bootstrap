// Run with `go run investigations/HELP-39364`
// If an error occurs on UpdateOne, this will panic.
package main

import (
	"fmt"
	"log"
	"os"
	"sync"

	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/event"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Set quiet to false to print command monitoring events.
const quiet bool = true

// toJson marshals a BSON type into a JSON string.
func toJson(v interface{}) string {
	bytes, err := bson.MarshalExtJSON(v, true, false)
	if err != nil {
		log.Fatal(err)
	}
	return string(bytes)
}

func logCommandStarted(context context.Context, event *event.CommandStartedEvent) {
	if !quiet {
		log.Println("=> CommandStarted:", event.CommandName, toJson(event.Command))
	}
}

func logCommandSucceeded(context context.Context, event *event.CommandSucceededEvent) {
	if !quiet {
		log.Println("<= CommandSucceeded:", event.CommandName, toJson(event.Reply))
	}
}

func logCommandFailed(context context.Context, event *event.CommandFailedEvent) {
	if !quiet {
		log.Println("<= CommandFailed:", event.CommandName, event.Failure)
	}
}

func AtomicIncrement(ctx context.Context, wg *sync.WaitGroup, coll *mongo.Collection, key string, fieldName string, count int) {
	doc := bson.M{"$inc": bson.M{fieldName: count}}

	// err object contains the reported error
	_, err := coll.UpdateOne(ctx, bson.M{"_id": key}, doc, options.Update().SetUpsert(true))
	if err != nil {
		log.Panicf("Got error on UpdateOne: %v", err)
	}

	wg.Done()
}

// When running this many times, sometimes, the above error occurs
// If we do the first update synchronously (by removing "go" statement), issue never appears again for the given limited use case
func repro(coll *mongo.Collection, wg *sync.WaitGroup) {
	wg.Add(1)
	wg.Add(1)
	go AtomicIncrement(context.TODO(), wg, coll, "document_key", "field1", 5)
	go AtomicIncrement(context.TODO(), wg, coll, "document_key", "field2", 6)
}

func main() {
	var uri string
	if uri = os.Getenv("MONGODB_URI"); uri == "" {
		uri = "mongodb://localhost:27017"
	}
	monitor := &event.CommandMonitor{Started: logCommandStarted, Succeeded: logCommandSucceeded, Failed: logCommandFailed}
	clientOpts := options.Client().ApplyURI(uri).SetMonitor(monitor)
	client, err := mongo.Connect(context.TODO(), clientOpts)
	if err != nil {
		panic(err)
	}

	coll := client.Database("db").Collection("coll")

	// Drop data from prior runs.
	err = coll.Drop(context.TODO())
	if err != nil {
		panic(err)
	}

	trials := 1000
	fmt.Printf("Testing with %v trials ... begin\n", trials)
	for i := 0; i < trials; i++ {
		// Remove data from prior to ensure one upsert succeeds.
		_, err := coll.DeleteMany(context.Background(), bson.M{})
		if err != nil {
			panic(err)
		}
		var wg sync.WaitGroup
		repro(coll, &wg)
		wg.Wait()
	}
	fmt.Printf("Testing with %v trials ... end\n", trials)
}
