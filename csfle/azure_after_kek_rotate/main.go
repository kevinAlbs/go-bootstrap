package main

// An example creating a DEK with Azure and encrypting a field.
// Run with: go run -tags cse ./csfle/azure_after_kek_rotate
//
// Set the environment variable KMS_PROVIDERS_PATH to the path of a JSON file with KMS credentials.
// KMS_PROVIDERS_PATH defaults to ~/.csfle/azure_after_kek_rotate.json.
//
// Set the environment variable MONGODB_URI to set a custom URI. MONGODB_URI defaults to
// mongodb://localhost:27017.
//
// Pass USE_KEYID_HEX to use an existing data key.

import (
	"context"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"go.mongodb.org/mongo-driver/bson"
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
		kmsProvidersPath = dirname + "/.csfle/azure_after_kek_rotate.json"
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

	var keyid primitive.Binary
	if USE_KEYID_HEX := os.Getenv("USE_KEYID_HEX"); USE_KEYID_HEX != "" {
		// Use supplied key ID.
		got, err := hex.DecodeString(USE_KEYID_HEX)
		if err != nil {
			log.Panicf("Failed to decode USE_KEYID_HEX: %v", err)
		}
		keyid = primitive.Binary{Data: got, Subtype: 4}
	} else {
		fmt.Printf("CreateDataKey... begin\n")
		keyid, err = ce.CreateDataKey(context.TODO(), "azure", options.DataKey().SetMasterKey(
			bson.D{
				{"keyVaultEndpoint", "key-vault-kevinalbs.vault.azure.net/"},
				{"keyName", "KeyName-KevinAlbs"},
			},
		))
		if err != nil {
			log.Panicf("CreateDataKey error: %v\n", err)
		}
		fmt.Printf("Created key with a UUID: %v\n", hex.EncodeToString(keyid.Data))
		fmt.Printf("CreateDataKey... end\n")
	}

	fmt.Printf("Encrypt... begin\n")
	plaintext := bson.RawValue{Type: bson.TypeString, Value: bsoncore.AppendString(nil, "test")}
	eOpts := options.Encrypt().SetAlgorithm("AEAD_AES_256_CBC_HMAC_SHA_512-Deterministic").SetKeyID(keyid)
	ciphertext, err := ce.Encrypt(context.TODO(), plaintext, eOpts)
	if err != nil {
		log.Panicf("Encrypt error: %v\n", err)
	}
	fmt.Printf("Explicitly encrypted to ciphertext: %x\n", ciphertext)
	fmt.Printf("Encrypt... end\n")

}

/*
## Results

Created a DEK with an Azure KEK without setting `keyVersion`.
Encrypted a value with the DEK. This succeeded.
Rotated the KEK in the Azure portal.
Encrypted a value with the DEK. This resulted in an error:
    Error in KMS response. HTTP status=400. Response body=
    {"error":{"code":"BadParameter","message":"The parameter is incorrect.\r\n"}}
Updated the DEK document in Compass to add `"keyVersion": "(original KEK version)"`
Encrypted a value with the DEK. This succeeded.
*/
