package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	uri := "foo"
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(uri))
	if err != nil {
		log.Fatalf("connect error: %v", err)
	}

	for {
		db := client.Database("test")
		res := db.RunCommand(context.Background(), bson.D{{"ping", 1}})
		var reply bson.D
		res.Decode(&reply)
		replyStr, err := bson.MarshalExtJSON(reply, false, false)
		if err != nil {
			fmt.Printf("err : %v\n", err)
		}
		fmt.Printf("res : %v\n", string(replyStr))
		time.Sleep(10 * time.Second)
	}

}
