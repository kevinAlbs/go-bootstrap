package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"sync"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func printConnectionId(client *mongo.Client) {

}

func main() {
	var uri string
	if uri = os.Getenv("MONGODB_URI"); uri == "" {
		log.Fatal("You must set your 'MONGODB_URI' environmental variable.")
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

	// Run "hello" on five concurrent goroutines.
	// This makes it likely to use different connections from the connection pool.
	var wg sync.WaitGroup
	for i := 1; i <= 5; i++ {
		wg.Add(1)
		go func() {
			resp := client.Database("db").RunCommand(context.TODO(), bson.D{{"hello", 1}})
			var respRaw bson.Raw
			resp.Decode(&respRaw)
			connectionId := respRaw.Lookup("connectionId")
			fmt.Println(connectionId)
			wg.Done()
		}()
	}

	wg.Wait()
}

// Sample output
// {"$numberInt":"55"}
// {"$numberInt":"56"}
// {"$numberInt":"57"}
// {"$numberInt":"58"}
// {"$numberInt":"59"}
