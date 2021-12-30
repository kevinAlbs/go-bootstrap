package main

import (
	"context"
	"log"
	"os"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/event"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readconcern"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"go.mongodb.org/mongo-driver/mongo/writeconcern"
)

func logCommandStarted(context context.Context, event *event.CommandStartedEvent) {
	bytes, err := bson.MarshalExtJSON(event.Command, true, false)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("=> CommandStarted:", event.CommandName, string(bytes))
}

func createMonitor() *event.CommandMonitor {
	return &event.CommandMonitor{Started: logCommandStarted}
}

func main() {
	var uri string
	if uri = os.Getenv("MONGODB_URI"); uri == "" {
		log.Fatal("You must set your 'MONGODB_URI' environmental variable.")
	}
	opts := options.Client().ApplyURI(uri).SetMonitor(createMonitor())

	client, err := mongo.Connect(context.TODO(), opts)
	if err != nil {
		log.Fatal(err)
	}

	defer func() {
		client.Disconnect(context.TODO())
	}()

	// Use a session with a majority read concern and majority write concern for "read your own writes" semantics.
	// See https://docs.mongodb.com/manual/core/causal-consistency-read-write-concerns/
	session, err := client.StartSession(options.Session())
	if err != nil {
		log.Fatal(err)
	}

	ctx := mongo.NewSessionContext(context.TODO(), session)

	rp := readpref.Secondary()
	rc := readconcern.Majority()
	wc := writeconcern.New(writeconcern.WMajority())
	coll := client.Database("db").Collection("coll", options.Collection().SetReadPreference(rp).SetReadConcern(rc).SetWriteConcern(wc))

	err = coll.Drop(ctx)
	if err != nil {
		log.Fatal(err)
	}

	_, err = coll.InsertMany(ctx, []interface{}{bson.D{}, bson.D{}, bson.D{}})
	if err != nil {
		log.Fatal(err)
	}

	cursor, err := coll.Find(ctx, bson.D{}, options.Find().SetBatchSize(1))
	if err != nil {
		log.Fatal(err)
	}

	// var results []bson.D
	// err = cursor.All(ctx, &results)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	_ = cursor.Close(ctx)
}

// Sample run using a replica set:
// % MONGODB_URI="mongodb://localhost:27017,localhost:27018" go run ./investigations/help28242/main.go
// 2021/10/29 10:47:46 => CommandStarted: drop {"drop":"coll","writeConcern":{"w":"majority"},"lsid":{"id":{"$binary":{"base64":"elte+cbQT02NSZsPPkugdA==","subType":"04"}}},"$clusterTime":{"clusterTime":{"$timestamp":{"t":1635518860,"i":4}},"signature":{"hash":{"$binary":{"base64":"AAAAAAAAAAAAAAAAAAAAAAAAAAA=","subType":"00"}},"keyId":{"$numberLong":"0"}}},"$db":"db"}
// 2021/10/29 10:47:47 => CommandStarted: insert {"insert":"coll","ordered":true,"writeConcern":{"w":"majority"},"lsid":{"id":{"$binary":{"base64":"elte+cbQT02NSZsPPkugdA==","subType":"04"}}},"txnNumber":{"$numberLong":"1"},"$clusterTime":{"clusterTime":{"$timestamp":{"t":1635518866,"i":1}},"signature":{"hash":{"$binary":{"base64":"AAAAAAAAAAAAAAAAAAAAAAAAAAA=","subType":"00"}},"keyId":{"$numberLong":"0"}}},"$db":"db","documents":[{"_id":{"$oid":"617c099310a1d186c02cf1ec"}},{"_id":{"$oid":"617c099310a1d186c02cf1ed"}}]}
// 2021/10/29 10:47:47 => CommandStarted: find {"find":"coll","batchSize":{"$numberInt":"1"},"filter":{},"readConcern":{"level":"majority","afterClusterTime":{"$timestamp":{"t":1635518867,"i":3}}},"lsid":{"id":{"$binary":{"base64":"elte+cbQT02NSZsPPkugdA==","subType":"04"}}},"$clusterTime":{"clusterTime":{"$timestamp":{"t":1635518867,"i":3}},"signature":{"hash":{"$binary":{"base64":"AAAAAAAAAAAAAAAAAAAAAAAAAAA=","subType":"00"}},"keyId":{"$numberLong":"0"}}},"$db":"db","$readPreference":{"mode":"secondary"}}
// 2021/10/29 10:47:48 => CommandStarted: getMore {"getMore":{"$numberLong":"1160624589967466247"},"collection":"coll","batchSize":{"$numberInt":"1"},"lsid":{"id":{"$binary":{"base64":"elte+cbQT02NSZsPPkugdA==","subType":"04"}}},"$clusterTime":{"clusterTime":{"$timestamp":{"t":1635518867,"i":3}},"signature":{"hash":{"$binary":{"base64":"AAAAAAAAAAAAAAAAAAAAAAAAAAA=","subType":"00"}},"keyId":{"$numberLong":"0"}}},"$db":"db","$readPreference":{"mode":"primaryPreferred"}}
// 2021/10/29 10:47:48 => CommandStarted: getMore {"getMore":{"$numberLong":"1160624589967466247"},"collection":"coll","batchSize":{"$numberInt":"1"},"lsid":{"id":{"$binary":{"base64":"elte+cbQT02NSZsPPkugdA==","subType":"04"}}},"$clusterTime":{"clusterTime":{"$timestamp":{"t":1635518867,"i":3}},"signature":{"hash":{"$binary":{"base64":"AAAAAAAAAAAAAAAAAAAAAAAAAAA=","subType":"00"}},"keyId":{"$numberLong":"0"}}},"$db":"db","$readPreference":{"mode":"primaryPreferred"}}
