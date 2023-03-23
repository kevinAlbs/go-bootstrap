/*
Create a document with two `_id` fields using bson.D.

Run with: go run ./investigations/double_id/main.go

Sample output:
% go run ./investigations/double_id/main.go
asBSON=[23 0 0 0 16 95 105 100 0 1 0 0 0 16 95 105 100 0 2 0 0 0 0]
asJSON={"_id":{"$numberInt":"1"},"_id":{"$numberInt":"2"}}
*/
package main

import (
	"fmt"
	"os"

	"go.mongodb.org/mongo-driver/bson"
)

func main() {
	var uri string
	if uri = os.Getenv("MONGODB_URI"); uri == "" {
		uri = "mongodb://localhost:27017"
	}
	got := bson.D{{"_id", 1}, {"_id", 2}}
	asBsonBytes, err := bson.Marshal(got)
	if err != nil {
		panic(err)
	}
	fmt.Printf("asBSON=%v\n", asBsonBytes)

	asJSON, err := bson.MarshalExtJSON(got, true /* canonical */, false /* escapeHTML */)
	if err != nil {
		panic(err)
	}
	fmt.Printf("asJSON=%s\n", asJSON)
}
