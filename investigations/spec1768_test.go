package common

import (
	"encoding/base64"
	"fmt"
	"log"
	"os"
	"strings"
	"testing"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/event"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"context"
	"io/ioutil"
)

func readJSONFile(filename string) []byte {
	file, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}
	contents, err := ioutil.ReadAll(file)
	if err != nil {
		log.Fatal(err)
	}

	return contents
}

type DeadlockTest struct {
	clientTest            *mongo.Client
	clientMetadataOpts    *options.ClientOptions
	clientMetadataEvents  []bson.Raw
	clientKeyVaultOpts    *options.ClientOptions
	clientKeyVaultEvents  []bson.Raw
	clientEncryption      *mongo.ClientEncryption
	clientEncryptedOpts   *options.ClientOptions
	clientEncryptedEvents []bson.Raw
	ciphertext            primitive.Binary
}

func getKMSProviders() map[string]map[string]interface{} {
	const localBase64 = "Mng0NCt4ZHVUYUJCa1kxNkVyNUR1QURhZ2h2UzR2d2RrZzh0cFBwM3R6NmdWMDFBMUN3YkQ5aXRRMkhGRGdQV09wOGVNYUMxT2k3NjZKelhaQmRCZGJkTXVyZG9uSjFk"
	raw, err := ioutil.ReadAll(base64.NewDecoder(base64.StdEncoding, strings.NewReader(localBase64)))
	if err != nil {
		log.Fatal(err)
	}

	kmsProviders := make(map[string]map[string]interface{})
	kmsProviders["local"] = make(map[string]interface{})
	kmsProviders["local"]["key"] = primitive.Binary{0, raw}
	return kmsProviders
}

func setupTest() DeadlockTest {
	d := DeadlockTest{}
	ctx := context.Background()
	var err error

	clientTestOpts := options.Client().ApplyURI("mongodb://localhost:27017").SetMaxPoolSize(1)
	if d.clientTest, err = mongo.Connect(ctx, clientTestOpts); err != nil {
		log.Fatal(err)
	}

	d.clientMetadataOpts = options.Client().ApplyURI("mongodb://localhost:27017").SetMaxPoolSize(1)
	d.clientMetadataOpts.SetMonitor(&event.CommandMonitor{
		func(ctx context.Context, event *event.CommandStartedEvent) {
			fmt.Println("clientmetadata command started")
			tmp := make(bson.Raw, len(event.Command))
			copy(tmp, event.Command)
			d.clientMetadataEvents = append(d.clientMetadataEvents, tmp)
		}, nil, nil,
	})

	// Go driver takes client options, not a client, to configure the key vault client.
	d.clientKeyVaultOpts = options.Client().ApplyURI("mongodb://localhost:27017").SetMaxPoolSize(1)
	d.clientKeyVaultOpts.SetMonitor(&event.CommandMonitor{
		func(ctx context.Context, event *event.CommandStartedEvent) {
			fmt.Println("clientkeyvault command started")
			tmp := make(bson.Raw, len(event.Command))
			copy(tmp, event.Command)
			d.clientKeyVaultEvents = append(d.clientKeyVaultEvents, tmp)
		}, nil, nil,
	})

	keyvaultColl := d.clientTest.Database("keyvault").Collection("datakeys")
	dataColl := d.clientTest.Database("db").Collection("coll")
	if err := dataColl.Drop(ctx); err != nil {
		log.Fatal(err)
	}

	if err := keyvaultColl.Drop(ctx); err != nil {
		log.Fatal(err)
	}

	var keyDoc bson.M
	keyDocJSON := readJSONFile("./spec1768/external-key.json")
	if err := bson.UnmarshalExtJSON(keyDocJSON, true, &keyDoc); err != nil {
		log.Fatal(err)
	}

	if _, err := keyvaultColl.InsertOne(ctx, keyDoc); err != nil {
		log.Fatal(err)
	}

	var schema bson.M
	schemaJSON := readJSONFile("./spec1768/external-schema.json")
	if err := bson.UnmarshalExtJSON(schemaJSON, true, &schema); err != nil {
		log.Fatal(err)
	}

	createOpts := options.CreateCollection().SetValidator(schema)
	if err := d.clientTest.Database("db").CreateCollection(ctx, "coll", createOpts); err != nil {
		log.Fatal(err)
	}

	kmsProviders := getKMSProviders()
	json, _ := bson.MarshalExtJSON(&kmsProviders, true, false)
	fmt.Println(string(json))
	ceOpts := options.ClientEncryption().SetKmsProviders(getKMSProviders()).SetKeyVaultNamespace("keyvault.datakeys")
	if d.clientEncryption, err = mongo.NewClientEncryption(d.clientTest, ceOpts); err != nil {
		log.Fatal(err)
	}

	var in bson.RawValue // TODO: there is probably a better way to create a bson.RawValue from a Go native type.
	if bytes, err := bson.Marshal(bson.M{"v": "string0"}); err != nil {
		log.Fatal(err)
	} else {
		asRaw := bson.Raw(bytes)
		in, err = asRaw.LookupErr("v")
		if err != nil {
			log.Fatal(err)
		}
	}

	d.ciphertext, err = d.clientEncryption.Encrypt(ctx, in, options.Encrypt().SetAlgorithm("AEAD_AES_256_CBC_HMAC_SHA_512-Deterministic").SetKeyAltName("local"))
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("encrypted: %x\n", d.ciphertext)

	d.clientEncryptedOpts = options.Client().ApplyURI("mongodb://localhost:27017")
	d.clientEncryptedOpts.SetMonitor(&event.CommandMonitor{
		func(ctx context.Context, event *event.CommandStartedEvent) {
			fmt.Println("clientencrypted command started")
			tmp := make(bson.Raw, len(event.Command))
			copy(tmp, event.Command)
			d.clientEncryptedEvents = append(d.clientEncryptedEvents, tmp)
		}, nil, nil,
	})
	return d
}

func TestSpec1768Case5(t *testing.T) {
	d := setupTest()
	ctx := context.Background()

	aeOpts := options.AutoEncryption().SetKeyVaultNamespace("keyvault.datakeys").SetKmsProviders(getKMSProviders())
	aeOpts.SetKeyVaultClientOptions(d.clientKeyVaultOpts)
	// TODO: aeOpts.SetMetadataClientOptions (d.clientMetadataOpts)
	d.clientEncryptedOpts.SetMaxPoolSize(1).SetAutoEncryptionOptions(aeOpts)

	clientEncrypted, err := mongo.Connect(ctx, d.clientEncryptedOpts)
	if err != nil {
		log.Fatal(err)
	}
	coll := clientEncrypted.Database("db").Collection("coll")
	fmt.Println("about to insert")
	coll.InsertOne(ctx, bson.M{"_id": 0, "encrypted": "string0"})
	fmt.Println("insert complete")

	// TODO
}
