// To run this example:
// 	 go run -tags cse ./csfle/qe_updateOne
// libmongocrypt is required to run this example.
// These environment variables are understood:
//   MONGODB_URI to a MongoDB URI.
//   SHOW_COMMAND_STARTED=ON to show command started events.
package main

import (
	"encoding/base64"
	"fmt"
	"log"
	"os"

	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/bsontype"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/event"
	"go.mongodb.org/mongo-driver/mongo"

	"go.mongodb.org/mongo-driver/mongo/options"
)

func ToJson(v interface{}) string {
	bytes, err := bson.MarshalExtJSON(v, true, false)
	if err != nil {
		log.Fatal(err)
	}
	return string(bytes)
}
func logCommandStarted(context context.Context, event *event.CommandStartedEvent) {
	cmdJSON, err := bson.MarshalExtJSONIndent(event.Command, true, false, "", "  ")
	if err != nil {
		log.Panicf("error in MarshalExtJSONIndent: %v", err)
	}
	log.Printf("=> CommandStarted event: %v\n%v", event.CommandName, string(cmdJSON))
}

func main() {
	var uri string
	if uri = os.Getenv("MONGODB_URI"); uri == "" {
		uri = "mongodb://localhost:27017"
	}
	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(uri))
	defer client.Disconnect(context.Background())
	if err != nil {
		log.Panicf("error in Connect: %v", err)
	}

	// Drop data from previous runs.
	{
		err = client.Database("db").Drop(context.Background())
		if err != nil {
			log.Panicf("error dropping database db: %v", err)
		}
	}

	// Create KMS provider credentials.
	var kmsProviders map[string]map[string]interface{}
	{
		// Use a 96 byte local Key Encryption Key (KEK) for testing.
		localKEK, err := base64.StdEncoding.DecodeString("01MICl70sM95GKn9ilmhssTroPYHCjTSv750baQFJIPlMLHHy+ZVbQE6l//sb+ZEkA9nz4TM/ETLOxBiQKfYS2jgrKpDZ9jkrDmTWEIe5JfGZoNCrnvWTSjPJ4/PRmrb")
		if err != nil {
			log.Panicf("error in DecodeString: %v", err)
		}

		kmsProviders = map[string]map[string]interface{}{
			"local": {"key": localKEK},
		}
	}

	// Create a Data Encryption Key (DEK).
	var keyID primitive.Binary
	{
		ceOpts := options.ClientEncryption().SetKmsProviders(kmsProviders).SetKeyVaultNamespace("db.datakeys")
		ce, err := mongo.NewClientEncryption(client, ceOpts)
		if err != nil {
			log.Panicf("error in NewClientEncryption: %v", err)
		}
		keyID, err = ce.CreateDataKey(context.Background(), "local")
		if err != nil {
			log.Panicf("error in CreateDataKey: %v", err)
		}
	}

	// Create a client configured with Queryable Encryption.
	var encryptedClient *mongo.Client
	{
		aeOpts := options.AutoEncryption().
			SetKeyVaultNamespace("db.datakeys").
			SetKmsProviders(kmsProviders)

		co := options.Client().ApplyURI(uri).SetAutoEncryptionOptions(aeOpts)
		if showCmdStarted := os.Getenv("SHOW_COMMAND_STARTED"); showCmdStarted == "ON" {
			co.SetMonitor(&event.CommandMonitor{Started: logCommandStarted})
		}
		encryptedClient, err = mongo.Connect(context.Background(), co)
		defer encryptedClient.Disconnect(context.Background())
		if err != nil {
			log.Panicf("error in Connect: %v", err)
		}
	}

	// Create collection with encrypted fields.
	var encryptedColl *mongo.Collection
	{
		encryptedFields := bson.M{
			"fields": bson.A{
				bson.M{
					"keyId":    keyID,
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
		err = encryptedClient.Database("db").CreateCollection(context.Background(), "encrypted", cco)
		if err != nil {
			log.Panicf("error in CreateCollection: %v", err)
		}
		encryptedColl = encryptedClient.Database("db").Collection("encrypted")
	}

	// UpdateOne in encrypted collection with upsert.
	{
		filter := bson.M{"_id": 1, "encryptedIndexed": "foo"}
		update := bson.M{"$set": bson.M{"foo": "bar"}}
		opts := options.Update().SetUpsert(true)
		_, err = encryptedColl.UpdateOne(context.Background(),
			filter,
			update,
			opts)
		if err != nil {
			log.Panicf("error in UpdateOne: %v", err)
		}
	}

	// Use an unencrypted client to check data was encrypted.
	{
		coll := client.Database("db").Collection("encrypted")
		res := coll.FindOne(context.Background(), bson.M{"_id": 1})
		if res.Err() != nil {
			log.Panicf("error in FindOne: %v", res.Err())
		}
		raw, err := res.DecodeBytes()
		if err != nil {
			log.Panicf("error in DecodeBytes: %v", err)
		}
		asJson, err := bson.MarshalExtJSON(raw, true, false)
		if err != nil {
			log.Panicf("error in MarshalExtJSON: %v", err)
		}
		gotType := raw.Lookup("encryptedIndexed").Type
		if gotType != bsontype.Binary {
			log.Panicf("expected Binary field 'encryptedIndexed', got: %v", string(asJson))
		}
	}

	fmt.Printf("Inserted into db.encrypted\n")
}
