package main

// An example reproducing the abort of MONGOCRYPT-380.
// Run with: go run -tags cse ./investigations/mongocrypt380

import (
	"context"
	"fmt"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	keyvaultClient, err := mongo.Connect(context.TODO())
	if err != nil {
		log.Fatalf("Connect error: %v\n", err)
	}
	defer keyvaultClient.Disconnect(context.TODO())

	kmsProvidersTmpl := `
	{
		"local": {
			"key": {
				"$binary": {
					"base64": "%s",
					"subType": "00"
				}
			}
		}
	}
`
	// Using an empty string for "base64" results in an abort.
	kmsProvidersStr := fmt.Sprintf(kmsProvidersTmpl, "")
	// Using a non-empty string for "base64" with an incorrect length results in an error.
	// kmsProvidersStr := fmt.Sprintf(kmsProvidersTmpl, "AAAA")

	var kmsProviders map[string]map[string]interface{}
	err = bson.UnmarshalExtJSON([]byte(kmsProvidersStr), false, &kmsProviders)
	if err != nil {
		log.Fatal(err)
	}

	ceopts := options.ClientEncryption().
		SetKmsProviders(kmsProviders).
		SetKeyVaultNamespace("keyvault.datakeys")

	ce, err := mongo.NewClientEncryption(keyvaultClient, ceopts)
	defer ce.Close(context.TODO())
	if err != nil {
		log.Fatalf("NewClientEncryption error: %v\n", err)
	}
}
