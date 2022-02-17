package main

import (
	"fmt"
	"os"

	"context"
	"encoding/hex"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"kevinalbs.com/go-bootstrap/util"
)

type foo struct {
	Bar *Int32type
}

type Int32type int32

func printJSON(raw interface{}) {
	asJson, err := bson.MarshalExtJSON(raw, true, false)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(asJson))
}

func (b *Int32type) UnmarshalBSON(in []byte) error {
	return nil
}

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

	test1 := foo{}
	raw, err := bson.Marshal(&test1)
	fmt.Println(hex.EncodeToString(raw))
	// Prints: "0a0000000a6261720000"
	printJSON(&test1)
	// Prints: { "bar": Null }
	var unmarshalled foo
	bson.Unmarshal(raw, &unmarshalled)
	fmt.Println(unmarshalled.Bar)
	// Prints: <nil>

	test2 := foo{}
	var int1 Int32type = 0
	test2.Bar = &int1
	raw, err = bson.Marshal(&test2)
	fmt.Println(hex.EncodeToString(raw))
	// Prints: ""
	printJSON(&test2)
	// Prints: {  }
	unmarshalled = foo{}
	bson.Unmarshal(raw, &unmarshalled)
	fmt.Println(*unmarshalled.Bar)
	// Prints: 0
}

/* I suspect there is little risk of backwards breaking behavior.

Currently

nil => BSON Null
BSON Null =>
*/
