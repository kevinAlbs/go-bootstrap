package main

// An example running all of the five functions of CSFLE.
// Run with: go run -tags cse ./csfle/five_functions

import (
	"context"
	"fmt"
	"log"
	"os"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/bsontype"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/x/bsonx/bsoncore"
	"kevinalbs.com/go-bootstrap/util"
)

func main() {
	var uri string
	if uri = os.Getenv("MONGODB_URI"); uri == "" {
		uri = "mongodb://localhost:27017"
	}

	keyvaultClient, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(uri))
	if err != nil {
		log.Fatalf("Connect error: %v\n", err)
	}
	defer keyvaultClient.Disconnect(context.TODO())

	kmsProviders := util.GetExampleKMSProviders()

	// A ClientEncryption struct provides admin helpers with three functions:
	// 1. create a data key
	// 2. explicit encrypt
	// 3. explicit decrypt
	ceopts := options.ClientEncryption().
		SetKmsProviders(kmsProviders).
		SetKeyVaultNamespace("keyvault.datakeys")

	ce, err := mongo.NewClientEncryption(keyvaultClient, ceopts)
	if err != nil {
		log.Fatalf("NewClientEncryption error: %v\n", err)
	}

	fmt.Printf("CreateDataKey begin...\n")
	keyid, err := ce.CreateDataKey(context.TODO(), "local", options.DataKey())
	if err != nil {
		log.Fatalf("CreateDataKey error: %v\n", err)
	}
	fmt.Printf("Created key with a UUID: %v\n", keyid)
	fmt.Printf("... CreateDataKey end\n")

	fmt.Printf("Encrypt begin...\n")
	plaintext := bson.RawValue{Type: bsontype.String, Value: bsoncore.AppendString(nil, "test")}
	eOpts := options.Encrypt().SetAlgorithm("AEAD_AES_256_CBC_HMAC_SHA_512-Deterministic").SetKeyID(keyid)
	ciphertext, err := ce.Encrypt(context.TODO(), plaintext, eOpts)
	if err != nil {
		log.Fatalf("Encrypt error: %v\n", err)
	}
	fmt.Printf("Explicitly encrypted to ciphertext: %v\n", ciphertext)
	fmt.Printf("...Encrypt end\n")

	fmt.Printf("Decrypt begin...\n")
	plaintext, err = ce.Decrypt(context.TODO(), ciphertext)
	if err != nil {
		log.Fatalf("Decrypt error: %v\n", err)
	}
	fmt.Printf("Explicitly decrypted to plaintext: %v\n", plaintext)
	fmt.Printf("...Decrypt end\n")

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
		log.Fatalf("Connect error: %v\n", err)
	}
	defer encryptedClient.Disconnect(context.TODO())

	coll := encryptedClient.Database("db").Collection("coll")
	err = coll.Drop(context.TODO())
	if err != nil {
		log.Fatalf("Drop error: %v\n", err)
	}

	fmt.Printf("Automatic encryption begin...\n")
	_, err = coll.InsertOne(context.TODO(), bson.D{{"encryptMe", "test"}})
	if err != nil {
		log.Fatalf("InsertOne error: %v\n", err)
	}
	fmt.Printf("... Automatic encryption end\n")

	fmt.Printf("Automatic decryption begin...\n")
	res := coll.FindOne(context.TODO(), bson.D{})
	if res.Err() != nil {
		log.Fatalf("FindOne error: %v\n", res.Err())
	}
	var decoded bson.D
	if err = res.Decode(&decoded); err != nil {
		log.Fatalf("Decode error: %v\n", err)
	}
	fmt.Printf("Decrypted document: %v\n", decoded)
	fmt.Printf("... Automatic decryption end\n")

}

/* Sample output

% go run -tags cse ./csfle/five_functions
CreateDataKey begin...
Created key with a UUID: {4 [196 197 200 201 210 71 65 234 139 157 212 161 13 102 36 8]}
... CreateDataKey end
Encrypt begin...
Explicitly encrypted to ciphertext: {6 [1 196 197 200 201 210 71 65 234 139 157 212 161 13 102 36 8 2 38 68 198 29 99 203 69 209 202 202 140 41 122 86 72 42 239 177 170 47 93 252 34 157 217 69 145 254 58 115 188 31 117 85 200 232 16 54 76 242 119 65 81 146 197 47 34 134 82 195 119 233 76 38 149 132 139 212 23 221 246 79 106 80]}
...Encrypt end
Decrypt begin...
Explicitly decrypted to plaintext: "test"
...Decrypt end
Automatic encryption begin...
... Automatic encryption end
Automatic decryption begin...
Decrypted document: [{_id ObjectID("61814ab1a5436ec2bdad84ab")} {encryptMe test}]
... Automatic decryption end
*/
