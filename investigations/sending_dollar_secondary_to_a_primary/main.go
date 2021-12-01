// Run with:
// go run ./investigations/sending_dollar_secondary_to_a_primary/main.go
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
		uri = "mongodb://localhost:27017/?directConnection=true"
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

This was tested against a single-node replica set of the following MongoDB Server versions:
(5.1.0-alpha-414-g69ca8b9, 5.0.2, 4.4.3, 4.2.12, 4.0.22, 3.6.21)
*/
