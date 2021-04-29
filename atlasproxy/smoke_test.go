package atlasproxy_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"kevinalbs.com/go-bootstrap/util"
)

func TestTransactions(t *testing.T) {
	uri := "mongodb://user:pencil@host.local.10gen.cc:9900,host.local.10gen.cc:9910,host.local.10gen.cc:9920/admin?tls=true&tlsCAFile=/Users/kevin.albertson/code/atlasproxy/main/ca.pem"
	opts := options.Client().ApplyURI(uri)
	opts.SetMonitor(util.CreateMonitor())
	client, err := mongo.Connect(context.TODO(), opts)
	assert.Nil(t, err, "connect error: %v", err)

	session, err := client.StartSession()
	assert.Nil(t, err, "startsession err: %v", session)

	ctx := mongo.NewSessionContext(context.Background(), session)

	err = session.StartTransaction()
	assert.Nil(t, err, "starttransaction err: %v", session)
	client.Database("db").Collection("coll").InsertOne(ctx, bson.D{{"x", 1}})
	err = session.CommitTransaction(ctx)
	assert.Nil(t, err, "committransaction err: %v", session)

}
