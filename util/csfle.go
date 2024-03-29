package util

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func GetExampleKMSProviders() map[string]map[string]interface{} {
	var kmsProviders map[string]map[string]interface{}
	const filename = "/Users/kevin.albertson/.csfle/kms_providers.json"

	file, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}
	contents, err := ioutil.ReadAll(file)
	if err != nil {
		log.Fatal(err)
	}

	err = bson.UnmarshalExtJSON(contents, false, &kmsProviders)
	if err != nil {
		log.Fatal(err)
	}

	return kmsProviders
}

// GetExampleSchemaMap gets a schema map to encrypt "secret" on "db.coll". It uses a key with the altname "example" (or creates one if it does not exist).
func GetExampleSchemaMap() map[string]interface{} {
	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI("mongodb://localhost:27017/"))
	if err != nil {
		log.Fatal(err)
	}

	// Check if there is a key already with the alternate name "example".
	keyvaultColl := client.Database("keyvault").Collection("datakeys")
	res := keyvaultColl.FindOne(context.Background(), bson.M{"keyAltNames": "example"})
	var keyid primitive.Binary
	if res.Err() == mongo.ErrNoDocuments {
		fmt.Println("key with keyAltNames:example not found, creating!")
		ceopts := options.ClientEncryption().SetKmsProviders(GetExampleKMSProviders()).SetKeyVaultNamespace("keyvault.datakeys")
		ce, err := mongo.NewClientEncryption(client, ceopts)
		if err != nil {
			log.Fatal(err)
		}
		dkopts := options.DataKey().SetKeyAltNames([]string{"example"})
		keyid, err = ce.CreateDataKey(context.Background(), "local", dkopts)
		if err != nil {
			log.Fatal(err)
		}
	} else if res.Err() != nil {
		log.Fatal("failed to find", res.Err())
	} else {
		var doc bson.M
		if err = res.Decode(&doc); err != nil {
			log.Fatal(err)
		}
		keyid = doc["_id"].(primitive.Binary)
	}
	fmt.Printf("keyid is %T, %v\n", keyid, keyid)

	schemaMap := bson.M{
		"db.coll": bson.M{
			"bsonType": "object",
			"properties": bson.M{
				"secret": bson.M{
					"encrypt": bson.M{
						"keyId":     bson.A{keyid},
						"algorithm": "AEAD_AES_256_CBC_HMAC_SHA_512-Deterministic",
						"bsonType":  "string",
					},
				},
			},
		},
	}

	out, err := bson.MarshalExtJSON(schemaMap, true, false)
	fmt.Printf("Produced this thing: %s\n", string(out))

	var ret map[string]interface{}

	// Convert from bson.M to map[string]interface{}.
	bytes, err := bson.Marshal(schemaMap)
	if err != nil {
		log.Fatal(err)
	}
	err = bson.Unmarshal(bytes, &ret)
	if err != nil {
		log.Fatal(err)
	}
	return ret
}
