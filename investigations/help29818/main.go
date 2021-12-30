package main

import (
	"fmt"
	"os"

	"context"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"kevinalbs.com/go-bootstrap/util"
)

func main() {
	var uri string
	if uri = os.Getenv("MONGODB_URI"); uri == "" {
		uri = "mongodb://localhost:27017"
	}
	clientopts := options.Client().
		ApplyURI(uri).
		SetServerMonitor(util.CreateServerMonitor())

	client, err := mongo.Connect(context.TODO(), clientopts)
	defer client.Disconnect(context.TODO())
	if err != nil {
		panic(err)
	}

	fmt.Println("ping")
	err = client.Ping(context.TODO(), nil)
	if err != nil {
		panic(err)
	}
}
