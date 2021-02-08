package investigations

import (
	"context"
	"testing"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"kevinalbs.com/go-bootstrap/util"
)

func TestReadPrefCommandCursor(t *testing.T) {
	goctx := context.Background()
	client, err := mongo.Connect(goctx, options.Client().ApplyURI("mongodb://localhost:27017").SetMonitor(util.CreateMonitor()))
	if err != nil {
		panic(err)
	}

	db := client.Database("test")
	opts := options.RunCmd().SetReadPreference(readpref.Secondary())

	cur, err := db.RunCommandCursor(goctx, bson.D{{"find", "foo"}, {"batchSize", 10}}, opts)
	if err != nil {
		panic(err)
	}
	for cur.Next(goctx) {

	}
}
