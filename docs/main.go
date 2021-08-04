package main

import (
	"context"
	"fmt"

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

	if err = client.Ping(context.TODO(), nil); err != nil {
		panic(err)
	}

	fmt.Println("Ping 1")

	if err = client.Disconnect(context.TODO()); err != nil {
		panic(err)
	}

	if err = client.Connect(context.TODO()); err != nil {
		panic(err) // panics here: server is closed
	}

	err = client.Ping(context.TODO(), nil)
	if err != nil {
		panic(err)
	}

	fmt.Println("Ping 2")

}
