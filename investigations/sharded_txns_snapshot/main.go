// Q: "the Go driver uses as Snapshot for sharded clusters when they use transactions, regardless of options the application sends. Is that true?"
// A: No.
// Run this with:
// go run ./investigations/sharded_txns_snapshot
// 2022/02/17 18:44:38 Transaction with default read concern... begin
// 2022/02/17 18:44:38 => CommandStarted: insert {"insert":"coll","ordered":true,"lsid":{"id":{"$binary":{"base64":"VB+ZM84QQwmpgzg45vtnXg==","subType":"04"}}},"txnNumber":{"$numberLong":"1"},"startTransaction":true,"autocommit":false,"$clusterTime":{"clusterTime":{"$timestamp":{"t":1645141474,"i":1}},"signature":{"hash":{"$binary":{"base64":"AAAAAAAAAAAAAAAAAAAAAAAAAAA=","subType":"00"}},"keyId":{"$numberLong":"0"}}},"$db":"db","documents":[{"_id":{"$oid":"620edde6d18dd064b0ca8f13"},"x":{"$numberInt":"1"}}]}
// 2022/02/17 18:44:38 => CommandStarted: commitTransaction {"commitTransaction":{"$numberInt":"1"},"recoveryToken":{"recoveryShardId":"shard01"},"lsid":{"id":{"$binary":{"base64":"VB+ZM84QQwmpgzg45vtnXg==","subType":"04"}}},"txnNumber":{"$numberLong":"1"},"autocommit":false,"$clusterTime":{"clusterTime":{"$timestamp":{"t":1645141474,"i":1}},"signature":{"hash":{"$binary":{"base64":"AAAAAAAAAAAAAAAAAAAAAAAAAAA=","subType":"00"}},"keyId":{"$numberLong":"0"}}},"$db":"admin"}
// 2022/02/17 18:44:38 Transaction with default read concern... end
// 2022/02/17 18:44:38 Transaction with snapshot read concern... begin
// 2022/02/17 18:44:38 => CommandStarted: insert {"insert":"coll","ordered":true,"readConcern":{"level":"snapshot","afterClusterTime":{"$timestamp":{"t":1645141478,"i":1}}},"lsid":{"id":{"$binary":{"base64":"VB+ZM84QQwmpgzg45vtnXg==","subType":"04"}}},"txnNumber":{"$numberLong":"2"},"startTransaction":true,"autocommit":false,"$clusterTime":{"clusterTime":{"$timestamp":{"t":1645141478,"i":1}},"signature":{"hash":{"$binary":{"base64":"AAAAAAAAAAAAAAAAAAAAAAAAAAA=","subType":"00"}},"keyId":{"$numberLong":"0"}}},"$db":"db","documents":[{"_id":{"$oid":"620edde6d18dd064b0ca8f14"},"x":{"$numberInt":"1"}}]}
// 2022/02/17 18:44:38 => CommandStarted: commitTransaction {"commitTransaction":{"$numberInt":"1"},"recoveryToken":{"recoveryShardId":"shard01"},"lsid":{"id":{"$binary":{"base64":"VB+ZM84QQwmpgzg45vtnXg==","subType":"04"}}},"txnNumber":{"$numberLong":"2"},"autocommit":false,"$clusterTime":{"clusterTime":{"$timestamp":{"t":1645141478,"i":1}},"signature":{"hash":{"$binary":{"base64":"AAAAAAAAAAAAAAAAAAAAAAAAAAA=","subType":"00"}},"keyId":{"$numberLong":"0"}}},"$db":"admin"}
// 2022/02/17 18:44:38 Transaction with snapshot read concern... end
package main

import (
	"log"
	"os"

	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readconcern"
	"kevinalbs.com/go-bootstrap/util"
)

func main() {
	var uri string
	if uri = os.Getenv("MONGODB_URI"); uri == "" {
		uri = "mongodb://localhost:27017"
	}
	ctx := context.TODO()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri).SetMonitor(util.CreateMonitor(&util.MonitorOpts{StartedOnly: true})))
	if err != nil {
		log.Fatalf("error in Connect: %v", err)
	}
	coll := client.Database("db").Collection("coll")

	sess, err := client.StartSession()

	log.Println("Transaction with default read concern... begin")
	if err != nil {
		log.Fatalf("error in StartSession: %v", err)
	}
	err = sess.StartTransaction()
	if err != nil {
		log.Fatalf("error in StartTransaction: %v", err)
	}
	sessctx := mongo.NewSessionContext(ctx, sess)
	_, err = coll.InsertOne(sessctx, bson.D{{"x", 1}})
	if err != nil {
		log.Fatalf("error in InsertOne: %v", err)
	}
	err = sess.CommitTransaction(ctx)
	if err != nil {
		log.Fatalf("error in CommitTransaction: %v", err)
	}
	log.Println("Transaction with default read concern... end")

	log.Println("Transaction with snapshot read concern... begin")
	txnopts := options.Transaction()
	txnopts.SetReadConcern(readconcern.Snapshot())
	err = sess.StartTransaction(txnopts)
	if err != nil {
		log.Fatalf("error in StartTransaction: %v", err)
	}
	_, err = coll.InsertOne(sessctx, bson.D{{"x", 1}})
	if err != nil {
		log.Fatalf("error in InsertOne: %v", err)
	}
	err = sess.CommitTransaction(ctx)
	if err != nil {
		log.Fatalf("error in CommitTransaction: %v", err)
	}
	log.Println("Transaction with snapshot read concern... end")
}
