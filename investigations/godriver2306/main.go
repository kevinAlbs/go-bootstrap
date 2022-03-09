package main

import (
	"fmt"
	"log"
	"os"

	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"kevinalbs.com/go-bootstrap/util"
)

type User struct {
	Id    primitive.ObjectID
	Email string
}

func collectionExists(db *mongo.Database, collName string) bool {
	collNames, err := db.ListCollectionNames(context.TODO(), bson.D{})
	if err != nil {
		log.Fatalf("error in ListCollectionNames: %v", err)
	}
	for _, cn := range collNames {
		if cn == collName {
			return true
		}
	}
	return false
}

func main() {
	var uri string
	if uri = os.Getenv("MONGODB_URI"); uri == "" {
		uri = "mongodb://localhost:27017"
	}
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(uri).SetMonitor(util.CreateMonitor(&util.MonitorOpts{StartedOnly: true})))
	if err != nil {
		panic(err)
	}
	err = client.Ping(context.TODO(), nil)
	if err != nil {
		panic(err)
	}
	ctx := context.TODO()

	// Q: can users pass a struct to InsertOne? A: Yes.
	// coll := client.Database("db").Collection("coll")
	// u := &User{Email: "foo@bar.com"}
	// _, err = coll.InsertOne(ctx, &u)
	// if err != nil {
	// 	panic(err)
	// }

	// Q: Does creating an index also create the collection? A:
	db := client.Database("db")
	err = db.Drop(ctx)
	if err != nil {
		log.Fatalf("error in Drop: %v", err)
	}
	coll := db.Collection("foo")
	// Check that "foo" does not exist.
	if collectionExists(db, "foo") {
		log.Fatalf("unexpected 'foo' exists")
	}

	indexName, err := coll.Indexes().CreateOne(ctx, mongo.IndexModel{Keys: bson.D{{"x", 1}}})
	fmt.Printf("created index: %v\n", indexName)
	if !collectionExists(db, "foo") {
		log.Fatalf("unexpected 'foo' does not exist")
	}

}
