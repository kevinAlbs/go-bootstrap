// Q: Does the Go driver parse legacy UUID format (subtype 3) and preserve the subtype?
// A: Yes. When parsing a subtype 3, the subtype 3 is kept.
// To run:
// Generate a legacy UUID with:
//   python generate_legacy.py.
// Run this with:
//   go run ./main.go

package main

import (
	"fmt"
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
	err = client.Ping(context.TODO(), nil)
	if err != nil {
		panic(err)
	}

	coll := client.Database("uuid_db").Collection("uuid_coll")
	res := coll.FindOne(context.TODO(), bson.D{})
	if res.Err() != nil {
		panic(res.Err())
	}

	str, err := bson.MarshalExtJSON(res, true /* canonical */, false /* escapeHTML */)
	fmt.Printf("uuid_id.uuid_coll contains this document: %v\n", string(str))
}
