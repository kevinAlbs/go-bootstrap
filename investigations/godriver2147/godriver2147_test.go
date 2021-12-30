package main

import (
	"context"
	"encoding/base64"
	"fmt"
	"testing"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	fakeKeyB64 = "dZFUVO2Cniuy1W/1Gj/aoL9ORTYm4uaftf9vbdG0ndgbvFV9p2MEJm8POy7m2MWncuyLl+KuPyNEmXbC5QPwYGuWJXYU6PoT5rxxz7pSP2i58y6ZMqWolJTKHdxz+rlu"
)

// code placeholder

func Test_CanYouClearASession(t *testing.T) {
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		t.Fatal(err)
	}
	ctx := context.TODO()
	session, err := client.StartSession()
	if err != nil {
		t.Fatal(err)
	}
	sessionCtx := mongo.NewSessionContext(ctx, session)
	retrieved := mongo.SessionFromContext(sessionCtx)
	fmt.Println(retrieved)
	sessionCtxCleared := mongo.NewSessionContext(sessionCtx, nil)
	retrieved = mongo.SessionFromContext(sessionCtxCleared)
	fmt.Println(retrieved)
	res := client.Database("admin").RunCommand(sessionCtxCleared, bson.D{{"ping", 1}})
	if res.Err() != nil {
		t.Fatal(res.Err())
	}
}

//taken from mongo_driver 1.5.1
func Test_Example_explictEncryptionWithAutomaticDecryption(t *testing.T) {
	// Automatic encryption requires MongoDB 4.2 enterprise, but automatic decryption is supported for all users.
	fakeKey, err := base64.StdEncoding.DecodeString(fakeKeyB64)
	if err != nil {
		t.Fatal(err)
	}
	//var localMasterKey []byte // This must be the same master key that was used to create the encryption key.
	kmsProviders := map[string]map[string]interface{}{
		"local": {
			"key": fakeKey,
		},
	}

	// The MongoDB namespace (db.collection) used to store the encryption data keys.
	keyVaultDBName, keyVaultCollName := "encryption", "testKeyVault"
	keyVaultNamespace := keyVaultDBName + "." + keyVaultCollName

	// Create the Client for reading/writing application data. Configure it with BypassAutoEncryption=true to disable
	// automatic encryption but keep automatic decryption. Setting BypassAutoEncryption will also bypass spawning
	// mongocryptd in the driver.
	autoEncryptionOpts := options.AutoEncryption().
		SetKmsProviders(kmsProviders).
		SetKeyVaultNamespace(keyVaultNamespace).
		SetBypassAutoEncryption(true)
	clientOpts := options.Client().
		ApplyURI("mongodb://localhost:27017").
		SetAutoEncryptionOptions(autoEncryptionOpts).
		SetMaxPoolSize(0)
	client, err := mongo.Connect(context.TODO(), clientOpts)
	if err != nil {
		panic(err)
	}
	defer func() { _ = client.Disconnect(context.TODO()) }()

	// Get a handle to the application collection and clear existing data.
	coll := client.Database("test").Collection("coll")
	_ = coll.Drop(context.TODO())

	// Set up the key vault for this example.
	keyVaultColl := client.Database(keyVaultDBName).Collection(keyVaultCollName)
	_ = keyVaultColl.Drop(context.TODO())
	// Ensure that two data keys cannot share the same keyAltName.
	keyVaultIndex := mongo.IndexModel{
		Keys: bson.D{{"keyAltNames", 1}},
		Options: options.Index().
			SetUnique(true).
			SetPartialFilterExpression(bson.D{
				{"keyAltNames", bson.D{
					{"$exists", true},
				}},
			}),
	}
	if _, err = keyVaultColl.Indexes().CreateOne(context.TODO(), keyVaultIndex); err != nil {
		panic(err)
	}

	// Create the ClientEncryption object to use for explicit encryption/decryption. The Client passed to
	// NewClientEncryption is used to read/write to the key vault. This can be the same Client used by the main
	// application.
	clientEncryptionOpts := options.ClientEncryption().
		SetKmsProviders(kmsProviders).
		SetKeyVaultNamespace(keyVaultNamespace)
	clientEncryption, err := mongo.NewClientEncryption(client, clientEncryptionOpts)
	if err != nil {
		panic(err)
	}
	defer func() { _ = clientEncryption.Close(context.TODO()) }()

	// Create a new data key for the encrypted field.
	dataKeyOpts := options.DataKey().SetKeyAltNames([]string{"go_encryption_example"})
	dataKeyID, err := clientEncryption.CreateDataKey(context.TODO(), "local", dataKeyOpts)
	if err != nil {
		panic(err)
	}

	// Create a bson.RawValue to encrypt and encrypt it using the key that was just created.
	rawValueType, rawValueData, err := bson.MarshalValue("123456789")
	if err != nil {
		panic(err)
	}
	rawValue := bson.RawValue{Type: rawValueType, Value: rawValueData}
	encryptionOpts := options.Encrypt().
		SetAlgorithm("AEAD_AES_256_CBC_HMAC_SHA_512-Deterministic").
		SetKeyID(dataKeyID)
	encryptedField, err := clientEncryption.Encrypt(context.TODO(), rawValue, encryptionOpts)
	if err != nil {
		panic(err)
	}

	// Insert a document with the encrypted field and then find it. The FindOne call will automatically decrypt the
	// field in the document.
	if _, err = coll.InsertOne(context.TODO(), bson.D{{"encryptedField", encryptedField}}); err != nil {
		panic(err)
	}
	var foundDoc bson.M
	//if err = coll.FindOne(context.TODO(), bson.D{}).Decode(&foundDoc); err != nil {
	// panic(err)
	//}
	//fmt.Printf("Decrypted document: %v\n", foundDoc)
	ctx := context.Background()
	err = client.UseSessionWithOptions(ctx, nil, func(sessionContext mongo.SessionContext) error { //sess
		if err = coll.FindOne(sessionContext, bson.D{}).Decode(&foundDoc); err != nil {
			panic(err)
		}
		fmt.Printf("Decrypted document: %v\n", foundDoc)
		return nil
	})
}
