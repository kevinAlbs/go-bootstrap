package main

import (
	"log"
	"os"

	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"kevinalbs.com/go-bootstrap/util"
)

func main() {
	var uri string
	if uri = os.Getenv("MONGODB_URI"); uri == "" {
		uri = "mongodb://localhost:27017"
	}
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(uri).SetMonitor(util.CreateMonitor()))
	if err != nil {
		panic(err)
	}

	// Drop db.cappedColl
	coll := client.Database("db").Collection("cappedColl")
	err = coll.Drop(context.TODO())
	if err != nil {
		panic(err)
	}
	// Create db.cappedColl with {capped: true}
	err = client.Database("db").CreateCollection(context.TODO(), "cappedColl", options.CreateCollection().SetCapped(true).SetSizeInBytes(4096))
	if err != nil {
		panic(err)
	}
	// Insert three documents
	{
		for i := 0; i < 1; i++ {
			_, err = coll.InsertOne(context.TODO(), bson.M{"_id": i})
			if err != nil {
				panic(err)
			}
		}
	}
	// Create a tailable cursor
	{
		cursor, err := coll.Find(context.TODO(), bson.M{}, options.Find().SetCursorType(options.Tailable))
		if err != nil {
			panic(err)
		}
		// cursor.TryNext() will only send one getMore. If no result, will stop sending getMore.
		// cursor.Next() will block until a result or error.
		ok := cursor.Next(context.TODO())
		if !ok {
			log.Fatalf("cursor.Next() failed")
		}
		// Expect this to hang, waiting for results, until a document is inserted into db.cappedColl
		ok = cursor.Next(context.TODO())
		if !ok {
			log.Fatalf("cursor.Next() failed")
		}
	}
}
