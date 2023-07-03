package main

// Q: Is an error returned if automatic decryption fails to find a key?
// A: Yes.
//
// Run with: go run -tags cse ./investigations/decrypt_with_missing_key
//
// Set the environment variable MONGODB_URI to set a custom URI. MONGODB_URI defaults to
// mongodb://localhost:27017.

import (
	"context"
	"encoding/hex"
	"fmt"
	"log"
	"os"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	var uri string
	if uri = os.Getenv("MONGODB_URI"); uri == "" {
		uri = "mongodb://localhost:27017"
	}

	keyvaultClient, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(uri))
	if err != nil {
		log.Panicf("Connect error: %v\n", err)
	}
	defer keyvaultClient.Disconnect(context.TODO())

	// Create KMS Providers. Use 96 bytes of zeroes local masterkey for testing.
	kmsProviders := map[string]map[string]interface{}{
		"local": {"key": make([]byte, 96)},
	}

	// Create a ClientEncryption object.
	ceopts := options.ClientEncryption().
		SetKmsProviders(kmsProviders).
		SetKeyVaultNamespace("keyvault.datakeys")

	ce, err := mongo.NewClientEncryption(keyvaultClient, ceopts)
	if err != nil {
		log.Panicf("NewClientEncryption error: %v\n", err)
	}

	fmt.Printf("CreateDataKey... begin\n")
	keyid, err := ce.CreateDataKey(context.TODO(), "local", options.DataKey())
	if err != nil {
		log.Panicf("CreateDataKey error: %v\n", err)
	}
	fmt.Printf("Created key with a UUID: %v\n", hex.EncodeToString(keyid.Data))
	fmt.Printf("CreateDataKey... end\n")

	schema := bson.D{
		{"bsonType", "object"},
		{"properties", bson.D{
			{"encryptMe", bson.D{
				{"encrypt", bson.D{
					{"keyId", bson.A{keyid}},
					{"bsonType", "string"},
					{"algorithm", "AEAD_AES_256_CBC_HMAC_SHA_512-Deterministic"},
				}},
			}},
		}},
	}
	schemaMap := map[string]interface{}{"db.coll": schema}

	aeOpts := options.AutoEncryption().
		SetKmsProviders(kmsProviders).
		SetKeyVaultNamespace("keyvault.datakeys").
		SetSchemaMap(schemaMap)

	encryptedClient, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(uri).SetAutoEncryptionOptions(aeOpts))
	if err != nil {
		log.Panicf("Connect error: %v\n", err)
	}
	defer encryptedClient.Disconnect(context.TODO())

	coll := encryptedClient.Database("db").Collection("coll")
	err = coll.Drop(context.TODO())
	if err != nil {
		log.Panicf("Drop error: %v\n", err)
	}

	fmt.Printf("Automatic encryption begin...\n")
	_, err = coll.InsertOne(context.TODO(), bson.D{{"encryptMe", "test"}})
	if err != nil {
		log.Panicf("InsertOne error: %v\n", err)
	}
	fmt.Printf("Automatic encryption... end\n")

	// Delete the key.
	_, err = ce.DeleteKey(context.TODO(), keyid)
	if err != nil {
		log.Panicf("DeleteKey error: %v\n", err)
	}

	// encryptedClient has cached the key. Recreate the encryptedClient.
	encryptedClient.Disconnect(context.TODO())
	encryptedClient, err = mongo.Connect(context.TODO(), options.Client().ApplyURI(uri).SetAutoEncryptionOptions(aeOpts))
	if err != nil {
		log.Panicf("Connect error: %v\n", err)
	}
	defer encryptedClient.Disconnect(context.TODO())

	coll = encryptedClient.Database("db").Collection("coll")

	fmt.Printf("Automatic decryption begin...\n")
	res := coll.FindOne(context.TODO(), bson.D{})
	if res.Err() != nil {
		log.Panicf("FindOne error: %v\n", res.Err())
	}
	var decoded bson.Raw
	if err = res.Decode(&decoded); err != nil {
		log.Panicf("Decode error: %v\n", err)
	}
	fmt.Printf("Decrypted document: %v\n", decoded)
	fmt.Printf("Automatic decryption... end\n")

}

/* Sample output:

CreateDataKey... begin
Created key with a UUID: afa7db95914f4fb4a92fff133ca4a6d4
CreateDataKey... end
Automatic encryption begin...
Automatic encryption... end
Automatic decryption begin...
2023/07/03 13:32:45 FindOne error: mongocrypt error 1: not all keys requested were satisfied. Verify that key vault DB/collection name was correctly specified.
panic: FindOne error: mongocrypt error 1: not all keys requested were satisfied. Verify that key vault DB/collection name was correctly specified.

*/
