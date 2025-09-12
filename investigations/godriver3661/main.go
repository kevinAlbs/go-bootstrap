package main

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

func main() {
	opts := options.Client().ApplyURI("mongodb://localhost:27017")
	kmsProviders := map[string]map[string]interface{}{
		"aws": {
			"accessKeyId":     "foo",
			"secretAccessKey": "bar",
		},
	}

	opts.SetAutoEncryptionOptions(options.AutoEncryption().SetKmsProviders(kmsProviders).SetKeyVaultNamespace("db.keyvault").SetExtraOptions(map[string]any{
		"mongocryptdPath": "C:\\bin\\mongodb-win32-x86_64-enterprise-windows-8.3.0-alpha0-1179-gb3a8643\\bin\\mongocryptd.exe",
	}))
	client, err := mongo.Connect(opts)
	if err != nil {
		panic(err)
	}

	defer func() {
		if err := client.Disconnect(context.TODO()); err != nil {
			panic(err)
		}
	}()

	_, err = client.Database("db").Collection("coll").InsertOne(context.TODO(), bson.M{"foo": "bar"})
	if err != nil {
		panic(err)
	}
	fmt.Printf("Test passed!\n")
}
