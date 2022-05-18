package main

// An example running all of the five functions of CSFLE.
// Run with: go run -tags cse ./csfle/five_functions
//
// Set the environment variable KMS_PROVIDERS_PATH to the path of a JSON file with KMS credentials.
// KMS_PROVIDERS_PATH defaults to ~/.csfle/kms_providers.json.
//
// Set the environment variable MONGODB_URI to set a custom URI. MONGODB_URI defaults to
// mongodb://localhost:27017.

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/bsontype"
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

	// A ClientEncryption struct provides admin helpers with three functions:
	// 1. create a data key
	// 2. explicit encrypt
	// 3. explicit decrypt
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
	fmt.Printf("Created key with a UUID: %v\n", keyid)
	fmt.Printf("CreateDataKey... end\n")

	fmt.Printf("Encrypt... begin\n")
	plaintext := bson.RawValue{Type: bsontype.String, Value: bsoncore.AppendString(nil, "test")}
	eOpts := options.Encrypt().SetAlgorithm("AEAD_AES_256_CBC_HMAC_SHA_512-Deterministic").SetKeyID(keyid)
	ciphertext, err := ce.Encrypt(context.TODO(), plaintext, eOpts)
	if err != nil {
		log.Panicf("Encrypt error: %v\n", err)
	}
	fmt.Printf("Explicitly encrypted to ciphertext: %v\n", ciphertext)
	fmt.Printf("Encrypt... end\n")

	fmt.Printf("Decrypt... begin\n")
	plaintext, err = ce.Decrypt(context.TODO(), ciphertext)
	if err != nil {
		log.Panicf("Decrypt error: %v\n", err)
	}
	fmt.Printf("Explicitly decrypted to plaintext: %v\n", plaintext)
	fmt.Printf("Decrypt... end\n")

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

	fmt.Printf("Automatic decryption begin...\n")
	res := coll.FindOne(context.TODO(), bson.D{})
	if res.Err() != nil {
		log.Panicf("FindOne error: %v\n", res.Err())
	}
	var decoded bson.D
	if err = res.Decode(&decoded); err != nil {
		log.Panicf("Decode error: %v\n", err)
	}
	fmt.Printf("Decrypted document: %v\n", decoded)
	fmt.Printf("Automatic decryption... end\n")

}

/* Sample output

% go run -tags cse ./csfle/five_functions
CreateDataKey... begin
Created key with a UUID: {4 [105 88 198 58 69 248 71 92 157 207 178 81 113 252 177 2]}
CreateDataKey... end
Encrypt... begin
Explicitly encrypted to ciphertext: {6 [1 105 88 198 58 69 248 71 92 157 207 178 81 113 252 177 2 2 13 122 179 77 14 163 123 162 194 62 19 16 223 162 247 51 233 15 182 32 25 28 210 107 225 52 140 158 129 40 27 130 35 42 93 117 176 102 153 98 168 11 180 39 1 200 196 79 101 194 230 60 104 126 113 218 112 250 99 196 7 81 40 37]}
Encrypt... end
Decrypt... begin
Explicitly decrypted to plaintext: "test"
Decrypt... end
Automatic encryption begin...
Automatic encryption... end
Automatic decryption begin...
Decrypted document: [{_id ObjectID("628501673dad7c11a8a4fbaf")} {encryptMe test}]
Automatic decryption... end
*/
