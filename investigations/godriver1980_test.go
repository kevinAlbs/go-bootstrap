package investigations

import (
	"context"
	"fmt"
	"log"
	"testing"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/bsontype"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// TEST DATA:
// > db.doctest.insertOne({"title": "test Document", "subdoc": {"key1": "value1", "key2": 123, "key3": ObjectId(), "key4": NumberLong("64"), "key5": NumberInt("123")}})

var coll *mongo.Collection

type MyMap struct{}

func (i *MyMap) UnmarshalBSON(data []byte) error {
	fmt.Println("calling UnmarshalBSON")
	return nil
}

func (i *MyMap) UnmarshalBSONValue(t bsontype.Type, data []byte) error {
	if i == nil {
		return fmt.Errorf("cannot unmarshal to nil")
	}
	var err error
	var raw bson.Raw = data
	// This results in an invalid document type per
	// https://github.com/mongodb/mongo-go-driver/blob/v1.5.1/bson/bsontype/bsontype.go#L12
	// https://github.com/mongodb/mongo-go-driver/blob/v1.5.1/bson/bsontype/bsontype.go#L92
	log.Printf("%10s -> %02d:%s", "document", t, t)

	elements, err := raw.Elements()
	if err != nil {
		log.Fatalf("Raw Element issue: %v", err)
	}

	for _, e := range elements {
		log.Printf("%10s -> %02d:%s", e.Key(), e.Value().Type, e.Value().Type)
	}
	return err
}

func Test1980(t *testing.T) {
	uri := "mongodb://localhost:27017"

	client, err := mongo.Connect(context.Background(),
		options.Client().ApplyURI(uri))

	if err != nil {
		log.Fatalf("failed connection: %s, %v", uri, err)
	}
	coll = client.Database("test").Collection("test")

	var myMap MyMap

	cursor, err := coll.Aggregate(context.Background(), bson.A{})

	if err != nil {
		log.Fatalf("failed to aggregate: %v", err)
	}

	log.Printf("-------- test 1 --------")
	for cursor.Next(context.Background()) {
		if err = cursor.Decode(&myMap); err != nil {
			log.Fatalf("failed to decode: %v", err)
		}
	}

	log.Printf("-------- test 2 --------")
	// The following aggregation will also show top level document to be
	// invalid
	cursor, err = coll.Aggregate(context.Background(),
		bson.A{bson.M{"$replaceRoot": bson.M{
			"newRoot": "$subdoc"}}},
	)

	if err != nil {
		log.Fatalf("failed to aggregate: %v", err)
	}

	for cursor.Next(context.Background()) {
		if err = cursor.Decode(&myMap); err != nil {
			log.Fatalf("failed to decode: %v", err)
		}
	}

}
