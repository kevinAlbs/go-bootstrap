package main

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Replace the uri string with your MongoDB deployment's connection string.
const uri = "mongodb://user:password@localhost:27017/?tls=true&tlsCAFile=/Users/kevin.albertson/code/mongo-c-driver/src/libmongoc/tests/x509gen/ca.pem"

func main() {
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(uri))

	if err != nil {
		panic(err)
	}
	defer func() {
		if err = client.Disconnect(context.TODO()); err != nil {
			panic(err)
		}
	}()

	// begin updatemany
	myCollection := client.Database("sample_mflix").Collection("theaters")
	filter := bson.D{{"theaterId", 1002}}
	update := bson.D{{"$set", bson.D{{"test1", "abc"}, {"test2", 5}}}}

	result, err := myCollection.UpdateMany(context.TODO(), filter, update)
	// end updatemany

	if err != nil {
		panic(err)
	}

	fmt.Printf("Documents updated: %v\n", result.ModifiedCount)

}
