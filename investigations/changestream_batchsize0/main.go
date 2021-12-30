package main

import (
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
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(uri).SetMonitor(util.CreateStartedMonitor()))
	if err != nil {
		panic(err)
	}

	coll := client.Database("db").Collection("coll")
	opts := options.ChangeStream().SetBatchSize(0)
	var cursorBatchSize int32 = 1
	opts.CursorBatchSize = &cursorBatchSize
	cs, err := coll.Watch(context.TODO(), []bson.D{}, opts)
	if err != nil {
		panic(err)
	}
	for cs.Next(context.TODO()) {
		var res bson.D
		err := cs.Decode(&res)
		if err != nil {
			panic(err)
		}
	}
}

/*
Observed events with changes in https://github.com/kevinAlbs/mongo-go-driver/commit/0fe797f9415936d603984d1922992c1794383ee2:

2021/11/17 15:00:48 => CommandStarted: aggregate {"aggregate":"coll","pipeline":[{"$changeStream":{"fullDocument":"default"}}],"cursor":{"batchSize":{"$numberInt":"0"}},"lsid":{"id":{"$binary":{"base64":"lcBs+IsOQOeQOALtJjlR8g==","subType":"04"}}},"$clusterTime":{"clusterTime":{"$timestamp":{"t":1637179239,"i":1}},"signature":{"hash":{"$binary":{"base64":"AAAAAAAAAAAAAAAAAAAAAAAAAAA=","subType":"00"}},"keyId":{"$numberLong":"0"}}},"$db":"db","$readPreference":{"mode":"primary"}}
2021/11/17 15:00:48 => CommandStarted: getMore {"getMore":{"$numberLong":"5638971789422902517"},"collection":"coll","batchSize":{"$numberInt":"1"},"lsid":{"id":{"$binary":{"base64":"lcBs+IsOQOeQOALtJjlR8g==","subType":"04"}}},"$clusterTime":{"clusterTime":{"$timestamp":{"t":1637179239,"i":1}},"signature":{"hash":{"$binary":{"base64":"AAAAAAAAAAAAAAAAAAAAAAAAAAA=","subType":"00"}},"keyId":{"$numberLong":"0"}}},"$db":"db","$readPreference":{"mode":"primaryPreferred"}}
*/
