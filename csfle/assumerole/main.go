package main

import (
	"fmt"
	"log"
	"os"

	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"kevinalbs.com/go-bootstrap/util"
)

func main() {
	var uri string
	if uri = os.Getenv("MONGODB_URI"); uri == "" {
		uri = "mongodb://localhost:27017"
	}
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(uri).SetMonitor(util.CreateMonitor()))
	if err != nil {
		panic(err)
	}

	kmsProviders := map[string]map[string]interface{}{
		"aws": {
			"accessKeyId":     os.Getenv("TEMP_ACCESS_KEY_ID"),
			"secretAccessKey": os.Getenv("TEMP_SECRET_ACCESS_KEY"),
			"sessionToken":    os.Getenv("TEMP_SESSION_TOKEN"),
		},
	}

	ceopts := options.ClientEncryption().
		SetKmsProviders(kmsProviders).
		SetKeyVaultNamespace("keyvault.datakeys")

	ce, err := mongo.NewClientEncryption(client, ceopts)
	if err != nil {
		log.Fatal(err)
	}

	dkopts := options.DataKey().SetKeyAltNames([]string{"example"}).SetMasterKey(bson.M{
		"region": os.Getenv("CMK_REGION"),
		"key":    os.Getenv("CMK_ARN"),
	})

	keyid, err := ce.CreateDataKey(context.Background(), "aws", dkopts)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("created key with a UUID: %v", keyid)
}
