package common_test

import (
	"context"
	"fmt"
	"log"
	"testing"

	"go.mongodb.org/mongo-driver/bson"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func TestUnorderedError(t *testing.T) {
	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		log.Fatalf("connect: %v", err)
	}

	coll := client.Database("test").Collection("test")
	err = coll.Drop(context.Background())
	coll.InsertOne(context.Background(), bson.D{{"u", 1}})
	coll.InsertOne(context.Background(), bson.D{{"d", 1}})

	res := client.Database("admin").RunCommand(context.Background(), bson.D{
		{"configureFailPoint", "failCommand"},
		{"mode", "off"},
	})
	if res.Err() != nil {
		log.Fatalf("disable failpoint: %v", res.Err())
	}

	// Configure different failpoints on update and delete.
	res = client.Database("admin").RunCommand(context.Background(), bson.D{
		{"configureFailPoint", "failCommand"},
		{"mode", bson.D{{"times", 2}}},
		{"data", bson.D{{"failCommands", bson.A{"update"}}, {"errorCode", 1}}},
	})
	if res.Err() != nil {
		log.Fatalf("enable failpoint: %v", res.Err())
	}

	res = client.Database("admin").RunCommand(context.Background(), bson.D{
		{"configureFailPoint", "failCommand"},
		{"mode", bson.D{{"times", 2}}},
		{"data", bson.D{{"failCommands", bson.A{"delete"}}, {"errorCode", 2}}},
	})
	if res.Err() != nil {
		log.Fatalf("enable failpoint: %v", res.Err())
	}

	u1 := mongo.UpdateOneModel{Filter: bson.D{{"u", 1}}, Update: bson.D{{"$set", bson.D{{"u", 2}}}}}
	d1 := mongo.DeleteOneModel{Filter: bson.D{{"d", 1}}}
	i1 := mongo.InsertOneModel{Document: bson.D{{"x", 1}}}
	bulk := []mongo.WriteModel{&u1, &d1, &i1}

	if err != nil {
		log.Fatalf("drop: %v", err)
	}
	bulkRes, err := coll.BulkWrite(context.Background(), bulk, options.BulkWrite().SetOrdered(false))
	fmt.Println("modified", bulkRes.ModifiedCount)
	fmt.Println("deleted", bulkRes.DeletedCount)
	fmt.Println("inserted", bulkRes.InsertedCount)
	fmt.Println("error", err)

}
