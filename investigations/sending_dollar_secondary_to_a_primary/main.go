package main

import (
	"os"

	"context"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"kevinalbs.com/go-bootstrap/util"
)

func main() {
	var uri string
	if uri = os.Getenv("MONGODB_URI"); uri == "" {
		uri = "mongodb://localhost:27017"
	}
	rp := readpref.Secondary()
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(uri).SetMonitor(util.CreateMonitor()))
	if err != nil {
		panic(err)
	}
	err = client.Ping(context.TODO(), rp)
	if err != nil {
		panic(err)
	}
}

/* Outputs:
2021/11/12 18:22:58 => CommandStarted: ping {"ping":{"$numberInt":"1"},"lsid":{"id":{"$binary":{"base64":"dvo8SfTVQzWsEJN4/cBchg==","subType":"04"}}},"$clusterTime":{"clusterTime":{"$timestamp":{"t":1636759370,"i":1}},"signature":{"hash":{"$binary":{"base64":"AAAAAAAAAAAAAAAAAAAAAAAAAAA=","subType":"00"}},"keyId":{"$numberLong":"0"}}},"$db":"admin","$readPreference":{"mode":"secondary"}}
2021/11/12 18:22:58 <= CommandSucceeded: ping {"ok":{"$numberDouble":"1.0"},"$clusterTime":{"clusterTime":{"$timestamp":{"t":1636759370,"i":1}},"signature":{"hash":{"$binary":{"base64":"AAAAAAAAAAAAAAAAAAAAAAAAAAA=","subType":"00"}},"keyId":{"$numberLong":"0"}}},"operationTime":{"$timestamp":{"t":1636759370,"i":1}}}
*/
