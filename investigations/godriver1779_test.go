package common_test

import (
	"context"
	"fmt"
	"log"
	"testing"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

func TestGodriver1779(t *testing.T) {
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI("mongodb://localhost:27017/"))
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("Pinging primary")
	var res *mongo.SingleResult = client.Database("test").RunCommand(context.TODO(), bson.D{{"ping", 1}}, &options.RunCmdOptions{readpref.Primary()})
	if res.Err() != nil {
		log.Fatal(res.Err())
	}
	var doc bson.M
	res.Decode(&doc)
	fmt.Println(doc["ok"])

	fmt.Println("Pinging secondary")
	res = client.Database("test").RunCommand(context.TODO(), bson.D{{"ping", 1}}, &options.RunCmdOptions{readpref.Secondary()})
	if res.Err() != nil {
		log.Fatal(res.Err())
	}
	res.Decode(&doc)
	fmt.Println(doc["ok"])
}
