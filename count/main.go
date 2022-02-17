package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"kevinalbs.com/go-bootstrap/util"
)

func main() {
	uri := "mongodb://localhost:27017"
	if os.Getenv("MONGODB_URI") == "" {
		fmt.Printf("Environment variable MONGODB_URI unset, using default: %v\n", uri)
	} else {
		uri = os.Getenv("MONGODB_URI")
	}

	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(uri).SetMonitor(util.CreateMonitor()))
	if err != nil {
		log.Fatalf("connect error: %v", err)
	}

	db := client.Database("test")
	cnt, err := db.Collection("coll").CountDocuments(context.TODO(), bson.D{})
	if err != nil {
		log.Fatalf("error in CountDocuments: %v", err)
	}
	fmt.Printf("CountDocuments on 'test.coll' returned: %v", cnt)

}
