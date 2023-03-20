package main

// An example running all of the five functions of CSFLE.
// Run with: go run -tags cse ./investigations/encrypt_without_decrypt
//
// Set the environment variable KMS_PROVIDERS_PATH to the path of a JSON file with KMS credentials.
// KMS_PROVIDERS_PATH defaults to ~/.csfle/kms_providers.json.
// Here is a sample kms_providers.json that can be used for testing the "local" KMS provider:
// {
//     "local": {
//         "key": "Mng0NCt4ZHVUYUJCa1kxNkVyNUR1QURhZ2h2UzR2d2RrZzh0cFBwM3R6NmdWMDFBMUN3YkQ5aXRRMkhGRGdQV09wOGVNYUMxT2k3NjZKelhaQmRCZGJkTXVyZG9uSjFk"
//     }
// }
//
// Set the environment variable MONGODB_URI to set a custom URI. MONGODB_URI defaults to
// mongodb://localhost:27017.

import (
	"context"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/bsontype"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/x/bsonx/bsoncore"
)

// readFile reads a file into a byte slice.
func readFile(path string) []byte {
	file, err := os.Open(path)
	defer file.Close()

	if err != nil {
		log.Panicf("error in Open: %v", err)
	}
	contents, err := ioutil.ReadAll(file)
	if err != nil {
		log.Panicf("error in ReadAll: %v", err)
	}

	return contents
}

// getKMSProvidersFromFile reads a JSON file for use as the KmsProviders option.
func getKMSProvidersFromFile(path string) map[string]map[string]interface{} {
	var kmsProviders map[string]map[string]interface{}
	contents := readFile(path)
	err := bson.UnmarshalExtJSON(contents, false, &kmsProviders)
	if err != nil {
		log.Panicf("error in UnmarshalExtJSON: %v", err)
	}

	return kmsProviders
}

func main() {
	var uri string
	if uri = os.Getenv("MONGODB_URI"); uri == "" {
		uri = "mongodb://localhost:27017"
	}

	var kmsProvidersPath string
	if kmsProvidersPath = os.Getenv("KMS_PROVIDERS_PATH"); kmsProvidersPath == "" {
		dirname, err := os.UserHomeDir()
		if err != nil {
			log.Fatal(err)
		}
		kmsProvidersPath = dirname + "/.csfle/kms_providers.json"
	}

	keyvaultClient, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(uri))
	if err != nil {
		log.Panicf("Connect error: %v\n", err)
	}
	defer keyvaultClient.Disconnect(context.TODO())

	kmsProviders := getKMSProvidersFromFile(kmsProvidersPath)

	// Create a ClientEncryption object.
	// A ClientEncryption struct provides admin helpers with three functions:
	// 1. create a data key
	// 2. explicit encrypt
	// 3. explicit decrypt
	var ce *mongo.ClientEncryption
	{
		ceopts := options.ClientEncryption().
			SetKmsProviders(kmsProviders).
			SetKeyVaultNamespace("keyvault.datakeys")

		ce, err = mongo.NewClientEncryption(keyvaultClient, ceopts)
		if err != nil {
			log.Panicf("NewClientEncryption error: %v\n", err)
		}
	}

	// Create a key.
	var keyid primitive.Binary
	{
		fmt.Printf("CreateDataKey... begin\n")
		keyid, err = ce.CreateDataKey(context.TODO(), "local", options.DataKey())
		if err != nil {
			log.Panicf("CreateDataKey error: %v\n", err)
		}
		fmt.Printf("Created key with a UUID: %v\n", hex.EncodeToString(keyid.Data))
		fmt.Printf("CreateDataKey... end\n")
	}

	// Insert a document with automatic encryption.
	{
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
	}

	// Query the document using explicit encryption.
	{
		fmt.Printf("Encrypt... begin\n")
		plaintext := bson.RawValue{Type: bsontype.String, Value: bsoncore.AppendString(nil, "test")}
		eOpts := options.Encrypt().SetAlgorithm("AEAD_AES_256_CBC_HMAC_SHA_512-Deterministic").SetKeyID(keyid)
		ciphertext, err := ce.Encrypt(context.TODO(), plaintext, eOpts)
		if err != nil {
			log.Panicf("Encrypt error: %v\n", err)
		}
		fmt.Printf("Explicitly encrypted to ciphertext: %x\n", ciphertext)
		fmt.Printf("Encrypt... end\n")

		// Use an unencrypted client to prevent automatic decryption.
		unencryptedClient, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(uri))
		if err != nil {
			log.Panicf("Connect error: %v\n", err)
		}
		defer unencryptedClient.Disconnect(context.TODO())

		coll := unencryptedClient.Database("db").Collection("coll")

		fmt.Printf("FindOne with explicit encryption... begin\n")
		got := coll.FindOne(context.TODO(), bson.D{{"encryptMe", ciphertext}})
		if got.Err() != nil {
			log.Panicf("FindOne error: %v\n", got.Err())
		}
		gotBson, err := got.DecodeBytes()
		if err != nil {
			log.Panicf("DecodeBytes error: %v\n", err)
		}
		fmt.Printf("FindOne with explicit encryption got... %v\n", gotBson)
		fmt.Printf("FindOne with explicit encryption... end\n")

		// Explicitly decrypt the value.
		val, err := gotBson.LookupErr("encryptMe")
		if err != nil {
			log.Panicf("Unable to find 'encryptMe' in result: %v", gotBson)
		}
		subtype, data, ok := val.BinaryOK()
		if !ok {
			log.Panicf("Unable to get 'encryptMe' as binary: %v", gotBson)
		}
		gotCiphertext := primitive.Binary{Data: data, Subtype: subtype}
		decrypted, err := ce.Decrypt(context.TODO(), gotCiphertext)
		if err != nil {
			log.Panicf("Decrypt error: %v\n", err)
		}
		if decrypted.StringValue() != "test" {
			log.Panicf("Expected to decrypt to 'test', got %v", decrypted)
		}
	}

}

/* Sample output
CreateDataKey... begin
Created key with a UUID: b9fc9e5a3a3d45a3981a8dac1b6bf48f
CreateDataKey... end
CreateDataKey... begin
Created key with a UUID: d7085b867c2142938fe58552b8c49e76
CreateDataKey... end
Automatic encryption begin...
Automatic encryption... end
Automatic decryption begin...
Decrypted document: {"_id": {"$oid":"64186f69d3bb9e50d5ee584d"},"encryptMe": "test"}
Automatic decryption... end
Encrypt... begin
Explicitly encrypted to ciphertext: {6 01d7085b867c2142938fe58552b8c49e7602875c322e75f1d67cf81066f8fa92fc3e244d2a8d127fa2ce5fdca7704ad39a55d646ec23179862d889a2a7593b20d204cd5cf222e6b034cf7967771a5fc46e99}
Encrypt... end
FindOne with explicit encryption... begin
FindOne with explicit encryption got... {"_id": {"$oid":"64186f69d3bb9e50d5ee584d"},"encryptMe": {"$binary":{"base64":"AdcIW4Z8IUKTj+WFUrjEnnYCh1wyLnXx1nz4EGb4+pL8PiRNKo0Sf6LOX9yncErTmlXWRuwjF5hi2Imip1k7INIEzVzyIuawNM95Z3caX8RumQ==","subType":"06"}}}
FindOne with explicit encryption... end
*/
