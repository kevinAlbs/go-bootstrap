package main

import (
	"encoding/hex"
	"log"
	"os"

	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"kevinalbs.com/go-bootstrap/util"
)

func cleanup(uri string) {
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(uri))
	if err != nil {
		panic(err)
	}
	defer client.Disconnect(context.TODO())
	client.Database("foo").Drop(context.TODO())
	client.Database("foo").Drop(context.TODO())
}

func main() {
	var uri string
	if uri = os.Getenv("MONGODB_URI"); uri == "" {
		uri = "mongodb://localhost:27017"
	}
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(uri).SetMonitor(util.CreateMonitor()))
	if err != nil {
		panic(err)
	}
	defer client.Disconnect(context.TODO())

	// Q: Can you create a collection in a transaction?
	// A: Yes.
	// See https://www.mongodb.com/docs/manual/core/transactions/
	// > Starting in MongoDB 4.4, you can create collections in transactions implicitly or explicitly.
	// {
	// 	cleanup(uri)
	// 	sess, err := client.StartSession()
	// 	if err != nil {
	// 		log.Fatalf("failed to StartSession: %v", err)
	// 	}
	// 	err = sess.StartTransaction()
	// 	if err != nil {
	// 		log.Fatalf("failed to StartTransaction: %v", err)
	// 	}

	// 	sctx := mongo.NewSessionContext(context.Background(), sess)
	// 	err = client.Database("foo").CreateCollection(sctx, "bar")
	// 	if err != nil {
	// 		log.Fatalf("failed to CreateCollection: %v", err)
	// 	}
	// 	err = sess.CommitTransaction(context.Background())
	// 	if err != nil {
	// 		log.Fatalf("failed to CommitTransaction: %v", err)
	// 	}
	// }

	// Q: Does the server error if a collection is created with encryptedFields on a standalone?
	// A: Yes.
	{
		cleanup(uri)

		keyId, err := hex.DecodeString("0368a3851137418e904e486f70d8cb78")
		if err != nil {
			log.Fatalf("error in DecodeString: %v", err)
		}
		encryptedFields := bson.M{
			"fields": bson.A{
				bson.M{
					"keyId":    primitive.Binary{Subtype: 4, Data: keyId},
					"path":     "encryptedIndexed",
					"bsonType": "string",
					"queries": bson.M{
						"queryType":  "equality",
						"contention": 0,
					},
				},
			},
		}
		cco := options.CreateCollection().SetEncryptedFields(encryptedFields)
		err = client.Database("foo").CreateCollection(context.Background(), "encrypted", cco)
		if err != nil {
			log.Fatalf("error in CreateCollection: %v", err)
		}
	}

	// Q: Do data keys need to be created before the collection with encryptedFields?
	// A: No.
	{
		cleanup(uri)

		keyId, err := hex.DecodeString("0368a3851137418e904e486f70d8cb78")
		if err != nil {
			log.Fatalf("error in DecodeString: %v", err)
		}
		encryptedFields := bson.M{
			"fields": bson.A{
				bson.M{
					"keyId":    primitive.Binary{Subtype: 4, Data: keyId},
					"path":     "encryptedIndexed",
					"bsonType": "string",
					"queries": bson.M{
						"queryType":  "equality",
						"contention": 0,
					},
				},
			},
		}
		cco := options.CreateCollection().SetEncryptedFields(encryptedFields)
		err = client.Database("foo").CreateCollection(context.Background(), "encrypted", cco)
		if err != nil {
			log.Fatalf("error in CreateCollection: %v", err)
		}

	}
}
