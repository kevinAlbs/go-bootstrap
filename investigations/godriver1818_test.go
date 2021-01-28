package common_test

import (
	"log"
	"testing"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"

	"context"
)

func TestGodriver1818(t *testing.T) {
	client, err := mongo.Connect(context.TODO())
	if err != nil {
		log.Fatal(err)
	}

	// test.coll contains one document: { "_id" : 0, "outsideIp" : [ "a", "b" ] }
	coll := client.Database("test").Collection("coll")
	update := bson.D{
		{"$addToSet", bson.D{{"outsideIp2", "string"}}},
	}
	filter := bson.D{{"_id", 0}}
	coll.UpdateOne(context.TODO(), filter, update)
	// the document should be updated to: { "_id" : 0, "outsideIp" : [ "a", "b", "string" ] }
}
