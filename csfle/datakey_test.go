package csfle

import (
	"context"
	"fmt"
	"log"
	"testing"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"kevinalbs.com/go-bootstrap/util"
)

// Drop and create a key vault collection.
// Create a key, and print out the UUID.

const (
	uri = "mongodb://localhost:27017"
)

func TestCreateDataKey(t *testing.T) {
	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(uri))
	if err != nil {
		t.Error(err)
	}

	kmsProviders := util.GetExampleKMSProviders()

	// A ClientEncryption struct provides admin helpers:
	// - create a data key
	// - explicit encrypt / decrypt
	ceopts := options.ClientEncryption().
		SetKmsProviders(kmsProviders).
		SetKeyVaultNamespace("keyvault.datakeys")

	ce, err := mongo.NewClientEncryption(client, ceopts)
	if err != nil {
		t.Error(err)
	}

	dkopts := options.DataKey().SetKeyAltNames([]string{"example"})
	keyid, err := ce.CreateDataKey(context.Background(), "local", dkopts)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("created key with a UUID: %v", keyid)
}
