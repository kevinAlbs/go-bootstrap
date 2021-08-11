package main

// Prerequisites:
// Get 2.4 server by using the download_mongodb.sh script in drivers-evergreen-tools
// Run this file with:
// go run ./investigations/connect_to_24/connect_to_24.go

import (
	"context"
	"fmt"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	uri := "mongodb://localhost:27017"
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(uri))
	if err != nil {
		log.Fatalf("connect error: %v", err)
	}

	db := client.Database("test")
	res := db.RunCommand(context.Background(), bson.D{{"ping", 1}})
	if nil != res.Err() {
		panic(res.Err())
		// On a 2.4.14 server, this panics with the error:
		// panic: server at localhost:27017 reports wire version 0, but this version of the Go driver requires at least 2 (MongoDB 2.6)
	}
	var reply bson.D
	res.Decode(&reply)
	replyStr, err := bson.MarshalExtJSON(reply, false, false)
	if err != nil {
		fmt.Printf("err : %v\n", err)
	}
	fmt.Printf("res : %v\n", string(replyStr))

}
