package main

// An example rewrapping local => AWS => local with a different Key Encryption Key (KEK).
// Run with: go run -tags cse ./csfle/rewrap_local
//
// Set the environment variables:
// LOCAL_KEK_1 to a base64 string representing 96 bytes.
// AWS_ACCESS_KEY_ID to the AWS credential.
// AWS_SECRET_ACCESS_KEY to the AWS credential.
// LOCAL_KEK_2 to a base64 string representing 96 bytes.
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
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	var uri string
	if uri = os.Getenv("MONGODB_URI"); uri == "" {
		uri = "mongodb://localhost:27017"
	}

	LOCAL_KEK_1 := os.Getenv("LOCAL_KEK_1")
	if LOCAL_KEK_1 == "" {
		log.Panicf("Expected environment variable LOCAL_KEK_1 to be set to a base64 string representing 96 bytes.")
	}

	AWS_ACCESS_KEY_ID := os.Getenv("AWS_ACCESS_KEY_ID")
	if AWS_ACCESS_KEY_ID == "" {
		log.Panicf("Expected environment variable AWS_ACCESS_KEY_ID to be set.")
	}

	AWS_SECRET_ACCESS_KEY := os.Getenv("AWS_SECRET_ACCESS_KEY")
	if AWS_SECRET_ACCESS_KEY == "" {
		log.Panicf("Expected environment variable AWS_SECRET_ACCESS_KEY to be set.")
	}

	LOCAL_KEK_2 := os.Getenv("LOCAL_KEK_2")
	if LOCAL_KEK_2 == "" {
		log.Panicf("Expected environment variable LOCAL_KEK_2 to be set to a base64 string representing 96 bytes.")
	}

	keyvaultClient, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(uri))
	if err != nil {
		log.Panicf("Connect error: %v\n", err)
	}
	defer keyvaultClient.Disconnect(context.TODO())

	// Create data key with local and LOCAL_KEK_1.
	var keyid primitive.Binary
	{
		// Configure ClientEncryption using LOCAL_KEK_1.
		kmsProviders := map[string]map[string]interface{}{
			"local": {
				"key": LOCAL_KEK_1,
			},
			"aws": {
				"accessKeyId":     AWS_ACCESS_KEY_ID,
				"secretAccessKey": AWS_SECRET_ACCESS_KEY,
			},
		}

		ceopts := options.ClientEncryption().
			SetKmsProviders(kmsProviders).
			SetKeyVaultNamespace("keyvault.datakeys")

		ce, err := mongo.NewClientEncryption(keyvaultClient, ceopts)
		if err != nil {
			log.Panicf("NewClientEncryption error: %v\n", err)
		}

		fmt.Printf("CreateDataKey: 'local with LOCAL_KEK_1'... begin\n")
		keyid, err = ce.CreateDataKey(context.TODO(), "local", options.DataKey())
		if err != nil {
			log.Panicf("CreateDataKey error: %v\n", err)
		}
		fmt.Printf("Created key with a UUID: %v\n", hex.EncodeToString(keyid.Data))
		fmt.Printf("CreateDataKey: 'local with LOCAL_KEK_1'... end\n")
	}

	// Rewrap from local to aws.
	{
		// Configure ClientEncryption using LOCAL_KEK_1.
		kmsProviders := map[string]map[string]interface{}{
			"local": {
				"key": LOCAL_KEK_1,
			},
			"aws": {
				"accessKeyId":     AWS_ACCESS_KEY_ID,
				"secretAccessKey": AWS_SECRET_ACCESS_KEY,
			},
		}

		ceopts := options.ClientEncryption().
			SetKmsProviders(kmsProviders).
			SetKeyVaultNamespace("keyvault.datakeys")

		ce, err := mongo.NewClientEncryption(keyvaultClient, ceopts)
		if err != nil {
			log.Panicf("NewClientEncryption error: %v\n", err)
		}

		fmt.Printf("RewrapManyDataKey from 'local with LOCAL_KEK_1' to 'aws' ... begin\n")
		_, err = ce.RewrapManyDataKey(context.TODO(), bson.M{"_id": keyid}, options.RewrapManyDataKey().SetProvider("aws").SetMasterKey(bson.M{
			"region": "us-east-1",
			"key":    "arn:aws:kms:us-east-1:579766882180:key/89fcc2c4-08b0-4bd9-9f25-e30687b580d0",
		}))
		if err != nil {
			log.Panicf("RewrapManyDataKey error: %v\n", err)
		}
		fmt.Printf("RewrapManyDataKey from 'local with LOCAL_KEK_1' to 'aws' ... end\n")
	}

	// Rewrap from aws to local.
	{
		// Configure ClientEncryption using LOCAL_KEK_2.
		kmsProviders := map[string]map[string]interface{}{
			"local": {
				"key": LOCAL_KEK_2,
			},
			"aws": {
				"accessKeyId":     AWS_ACCESS_KEY_ID,
				"secretAccessKey": AWS_SECRET_ACCESS_KEY,
			},
		}

		ceopts := options.ClientEncryption().
			SetKmsProviders(kmsProviders).
			SetKeyVaultNamespace("keyvault.datakeys")

		ce, err := mongo.NewClientEncryption(keyvaultClient, ceopts)
		if err != nil {
			log.Panicf("NewClientEncryption error: %v\n", err)
		}

		fmt.Printf("RewrapManyDataKey from 'aws' to 'local with LOCAL_KEK_2' ... begin\n")
		_, err = ce.RewrapManyDataKey(context.TODO(), bson.M{"_id": keyid}, options.RewrapManyDataKey().SetProvider("aws").SetMasterKey(bson.M{
			"region": "us-east-1",
			"key":    "arn:aws:kms:us-east-1:579766882180:key/89fcc2c4-08b0-4bd9-9f25-e30687b580d0",
		}))
		if err != nil {
			log.Panicf("RewrapManyDataKey error: %v\n", err)
		}
		fmt.Printf("RewrapManyDataKey from 'aws' to 'local with LOCAL_KEK_2' ... end\n")
	}

}
