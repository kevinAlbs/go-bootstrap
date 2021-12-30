// Test using AdvanceClusterTime and AdvanceOperationTime.
// Run with:
// export MONGODB_URI="<URI>"
// go run ./investigations/synchronizing_sessions/main.go

package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readconcern"
	"go.mongodb.org/mongo-driver/mongo/writeconcern"
)

func main() {
	uri := os.Getenv("MONGODB_URI")
	opts := options.Client().ApplyURI(uri)
	client, err := mongo.Connect(context.TODO(), opts)
	if err != nil {
		log.Fatal(err)
	}

	defer client.Disconnect(context.TODO())

	sess1, err := client.StartSession()
	if err != nil {
		log.Fatal("StartSession error", err)
	}

	txnOpts := options.Transaction().
		SetReadConcern(readconcern.Majority()).
		SetWriteConcern(writeconcern.New(writeconcern.WMajority()))

	if err := sess1.StartTransaction(txnOpts); err != nil {
		log.Fatal("StartTransaction error", err)
	}

	ctx := mongo.NewSessionContext(context.TODO(), sess1)
	_, err = client.Database("db").Collection("coll").InsertOne(ctx, bson.D{{"x", 1}})
	if err != nil {
		log.Fatal("InsertOne error", err)
	}

	if err = sess1.CommitTransaction(context.TODO()); err != nil {
		log.Fatal("CommitTransaction error", err)
	}

	sess2, err := client.StartSession()
	sess2.AdvanceClusterTime(sess1.ClusterTime())
	sess2.AdvanceOperationTime(sess1.OperationTime())
	if err != nil {
		log.Fatal("StartSession error: ", err)
	}
	if err := sess2.StartTransaction(txnOpts); err != nil {
		log.Fatal("StartTransaction error", err)
	}

	ctx = mongo.NewSessionContext(context.TODO(), sess2)
	res := client.Database("db").Collection("coll").FindOne(ctx, bson.D{{"x", 1}})
	if res.Err() != nil {
		log.Fatal("FindOne error", res.Err())
	}
	var resBson bson.M
	if err = res.Decode(&resBson); err != nil {
		log.Fatal("Decode error", err)
	}

	fmt.Println("Got result", resBson)

	if err = sess2.CommitTransaction(context.TODO()); err != nil {
		log.Fatal("CommitTransaction error", err)
	}

}
