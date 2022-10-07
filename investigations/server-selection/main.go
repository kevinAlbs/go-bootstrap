// Run this example with: go run investigations/server-selection/main.go

package main

import (
	"context"
	"fmt"
	"log"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/version"
)

func main() {
	fmt.Printf("Driver version: %v\n", version.Driver)
	uri := "mongodb://example.com:27017"
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(uri))
	if err != nil {
		log.Fatalf("error in Connect: %v\n", err)
	}
	// Expect a server selection error after 30 seconds.
	err = client.Ping(context.TODO(), nil)
	if err != nil {
		log.Fatalf("error in Ping: %v\n", err)
	}
}
