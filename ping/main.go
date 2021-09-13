package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	uri := "mongodb://localhost:27017"
	if os.Getenv("MONGODB_URI") == "" {
		fmt.Printf("Environment variable MONGODB_URI unset, using default: %v\n", uri)
	} else {
		uri = os.Getenv("MONGODB_URI")
	}

	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(uri))
	if err != nil {
		log.Fatalf("connect error: %v", err)
	}

	db := client.Database("test")
	res := db.RunCommand(context.Background(), bson.D{{"ping", 1}})
	if nil != res.Err() {
		log.Fatalf("error sending 'ping': %v", res.Err())
	}
	var reply bson.D
	res.Decode(&reply)
	replyStr, err := bson.MarshalExtJSON(reply, false, false)
	if err != nil {
		log.Fatalf("error decoding reply: %v\n", err)
	}
	fmt.Printf("'ping' command replied with: %v\n", string(replyStr))
}
