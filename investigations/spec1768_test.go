package common

// import (
// 	"bytes"
// 	"encoding/base64"
// 	"fmt"
// 	"log"
// 	"os"
// 	"strings"
// 	"sync"
// 	"testing"

// 	"go.mongodb.org/mongo-driver/bson"
// 	"go.mongodb.org/mongo-driver/bson/primitive"
// 	"go.mongodb.org/mongo-driver/event"
// 	"go.mongodb.org/mongo-driver/mongo"
// 	"go.mongodb.org/mongo-driver/mongo/options"

// 	"context"
// 	"io/ioutil"
// )

// const DEBUG = false

// func readJSONFile(filename string) []byte {
// 	file, err := os.Open(filename)
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	contents, err := ioutil.ReadAll(file)
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	return contents
// }

// func toJson(v interface{}) string {
// 	bytes, err := bson.MarshalExtJSON(v, true, false)
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	return string(bytes)
// }

// type DeadlockTest struct {
// 	clientTest           *mongo.Client
// 	clientMetadataOpts   *options.ClientOptions
// 	clientMetadataEvents []bson.Raw
// 	clientKeyVaultOpts   *options.ClientOptions
// 	clientKeyVaultEvents []bson.Raw
// 	clientEncryption     *mongo.ClientEncryption
// 	ciphertext           primitive.Binary
// }

// func getKMSProviders() map[string]map[string]interface{} {
// 	const localBase64 = "Mng0NCt4ZHVUYUJCa1kxNkVyNUR1QURhZ2h2UzR2d2RrZzh0cFBwM3R6NmdWMDFBMUN3YkQ5aXRRMkhGRGdQV09wOGVNYUMxT2k3NjZKelhaQmRCZGJkTXVyZG9uSjFk"
// 	raw, err := ioutil.ReadAll(base64.NewDecoder(base64.StdEncoding, strings.NewReader(localBase64)))
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	kmsProviders := make(map[string]map[string]interface{})
// 	kmsProviders["local"] = make(map[string]interface{})
// 	kmsProviders["local"]["key"] = primitive.Binary{0, raw}
// 	return kmsProviders
// }

// func makeCaptureMonitor(name string, dest *[]bson.Raw) *event.CommandMonitor {
// 	return &event.CommandMonitor{func(ctx context.Context, event *event.CommandStartedEvent) {
// 		fmt.Printf("[%v] command started (%v, %v)\n", name, event.CommandName, event.DatabaseName)
// 		if DEBUG {
// 			fmt.Println(toJson(event.Command))
// 		}
// 		tmp := make(bson.Raw, len(event.Command))
// 		copy(tmp, event.Command)
// 		*dest = append(*dest, tmp)
// 	}, nil, nil}
// }

// func setupTest() *DeadlockTest {
// 	d := DeadlockTest{}
// 	ctx := context.Background()
// 	var err error

// 	// TODO: configure with read/write concern majority.

// 	clientTestOpts := options.Client().ApplyURI("mongodb://localhost:27017").SetMaxPoolSize(1)
// 	if d.clientTest, err = mongo.Connect(ctx, clientTestOpts); err != nil {
// 		log.Fatal(err)
// 	}

// 	d.clientMetadataOpts = options.Client().ApplyURI("mongodb://localhost:27017").SetMaxPoolSize(1)
// 	d.clientMetadataOpts.SetMonitor(makeCaptureMonitor("clientMetadata", &d.clientMetadataEvents))

// 	// Go driver takes client options, not a client, to configure the key vault client.
// 	d.clientKeyVaultOpts = options.Client().ApplyURI("mongodb://localhost:27017").SetMaxPoolSize(1)
// 	d.clientKeyVaultOpts.SetMonitor(makeCaptureMonitor("clientKeyVault", &d.clientKeyVaultEvents))

// 	keyvaultColl := d.clientTest.Database("keyvault").Collection("datakeys")
// 	dataColl := d.clientTest.Database("db").Collection("coll")
// 	if err := dataColl.Drop(ctx); err != nil {
// 		log.Fatal(err)
// 	}

// 	if err := keyvaultColl.Drop(ctx); err != nil {
// 		log.Fatal(err)
// 	}

// 	var keyDoc bson.M
// 	keyDocJSON := readJSONFile("./spec1768/external-key.json")
// 	if err := bson.UnmarshalExtJSON(keyDocJSON, true, &keyDoc); err != nil {
// 		log.Fatal(err)
// 	}

// 	if _, err := keyvaultColl.InsertOne(ctx, keyDoc); err != nil {
// 		log.Fatal(err)
// 	}

// 	var schema bson.M
// 	schemaJSON := readJSONFile("./spec1768/external-schema.json")
// 	if err := bson.UnmarshalExtJSON(schemaJSON, true, &schema); err != nil {
// 		log.Fatal(err)
// 	}

// 	createOpts := options.CreateCollection().SetValidator(bson.M{"$jsonSchema": schema})
// 	if err := d.clientTest.Database("db").CreateCollection(ctx, "coll", createOpts); err != nil {
// 		log.Fatal(err)
// 	}

// 	kmsProviders := getKMSProviders()
// 	json, _ := bson.MarshalExtJSON(&kmsProviders, true, false)
// 	fmt.Println(string(json))
// 	ceOpts := options.ClientEncryption().SetKmsProviders(getKMSProviders()).SetKeyVaultNamespace("keyvault.datakeys")
// 	if d.clientEncryption, err = mongo.NewClientEncryption(d.clientTest, ceOpts); err != nil {
// 		log.Fatal(err)
// 	}

// 	var in bson.RawValue // TODO: there is probably a better way to create a bson.RawValue from a Go native type.
// 	if bytes, err := bson.Marshal(bson.M{"v": "string0"}); err != nil {
// 		log.Fatal(err)
// 	} else {
// 		asRaw := bson.Raw(bytes)
// 		in, err = asRaw.LookupErr("v")
// 		if err != nil {
// 			log.Fatal(err)
// 		}
// 	}

// 	d.ciphertext, err = d.clientEncryption.Encrypt(ctx, in, options.Encrypt().SetAlgorithm("AEAD_AES_256_CBC_HMAC_SHA_512-Deterministic").SetKeyAltName("local"))
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	return &d
// }

// func TestSpec1768Case4(t *testing.T) {
// 	d := setupTest()
// 	ctx := context.Background()

// 	var clientEncryptedEvents []bson.Raw

// 	aeOpts := options.AutoEncryption()
// 	aeOpts.SetKeyVaultNamespace("keyvault.datakeys")
// 	aeOpts.SetKmsProviders(getKMSProviders())
// 	aeOpts.SetKeyVaultClientOptions(d.clientKeyVaultOpts)
// 	aeOpts.SetBypassAutoEncryption(true)

// 	ceOpts := options.Client().ApplyURI("mongodb://localhost:27017")
// 	ceOpts.SetMonitor(makeCaptureMonitor("clientEncrypted", &clientEncryptedEvents))
// 	ceOpts.SetMaxPoolSize(1)
// 	ceOpts.SetAutoEncryptionOptions(aeOpts)

// 	clientEncrypted, err := mongo.Connect(ctx, ceOpts)
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	coll := clientEncrypted.Database("db").Collection("coll")
// 	_, err = coll.InsertOne(ctx, bson.M{"_id": 0, "encrypted": d.ciphertext})
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	res := coll.FindOne(ctx, bson.M{"_id": 0})
// 	if res.Err() != nil {
// 		log.Fatal(res.Err())
// 	}
// 	// TODO: what is the right way to compare BSON.
// 	if raw, err := res.DecodeBytes(); err != nil {
// 		log.Fatal(res.Err())
// 	} else {
// 		expected, _ := bson.Marshal(bson.D{{"_id", 0}, {"encrypted", "string0"}})
// 		if bytes.Compare(expected, raw) != 0 {
// 			log.Fatal("not equal")
// 		}
// 	}

// 	// TODO: actually compare the command started event commands.
// 	if len(clientEncryptedEvents) != 2 {
// 		log.Fatal("expected 2 events in clientEncryptedEvents")
// 	}

// 	if len(d.clientKeyVaultEvents) != 1 {
// 		log.Fatal("expected 1 event in clientKeyVaultEvents")
// 	}

// 	if len(d.clientMetadataEvents) != 0 {
// 		log.Fatal("expected 0 event in clientMetadataEvents")
// 	}
// }

// func TestSpec1768Case5(t *testing.T) {
// 	d := setupTest()
// 	ctx := context.Background()

// 	aeOpts := options.AutoEncryption().SetKeyVaultNamespace("keyvault.datakeys").SetKmsProviders(getKMSProviders())
// 	aeOpts.SetKeyVaultClientOptions(d.clientKeyVaultOpts)
// 	aeOpts.SetMetadataClientOptions(d.clientMetadataOpts)

// 	var clientEncryptedEvents []bson.Raw
// 	ceOpts := options.Client().ApplyURI("mongodb://localhost:27017")
// 	ceOpts.SetMonitor(makeCaptureMonitor("clientEncrypted", &clientEncryptedEvents))
// 	ceOpts.SetMaxPoolSize(1)
// 	ceOpts.SetAutoEncryptionOptions(aeOpts)

// 	clientEncrypted, err := mongo.Connect(ctx, ceOpts)
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	coll := clientEncrypted.Database("db").Collection("coll")
// 	_, err = coll.InsertOne(ctx, bson.M{"_id": 0, "encrypted": "string0"})
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	res := coll.FindOne(ctx, bson.M{"_id": 0})
// 	if res.Err() != nil {
// 		log.Fatal(res.Err())
// 	}
// 	// TODO: what is the right way to compare BSON.
// 	if raw, err := res.DecodeBytes(); err != nil {
// 		log.Fatal(res.Err())
// 	} else {
// 		expected, _ := bson.Marshal(bson.D{{"_id", 0}, {"encrypted", "string0"}})
// 		if bytes.Compare(expected, raw) != 0 {
// 			log.Fatal("not equal")
// 		}
// 	}

// 	// TODO: actually compare the command started event commands.
// 	if len(clientEncryptedEvents) != 2 {
// 		log.Fatal("expected 2 events in clientEncryptedEvents")
// 	}

// 	if len(d.clientKeyVaultEvents) != 1 {
// 		log.Fatal("expected 1 event in clientKeyVaultEvents")
// 	}

// 	if len(d.clientMetadataEvents) != 1 {
// 		log.Fatal("expected 1 event in clientMetadataEvents")
// 	}
// }

// func TestSpec1768Case6(t *testing.T) {
// 	d := setupTest()
// 	ctx := context.Background()

// 	aeOpts := options.AutoEncryption().SetKeyVaultNamespace("keyvault.datakeys").SetKmsProviders(getKMSProviders())
// 	aeOpts.SetKeyVaultClientOptions(d.clientKeyVaultOpts)
// 	aeOpts.SetMetadataClientOptions(d.clientMetadataOpts)

// 	ceOpts := options.Client().ApplyURI("mongodb://localhost:27017")
// 	ceOpts.SetMaxPoolSize(1)
// 	ceOpts.SetAutoEncryptionOptions(aeOpts)

// 	clientEncrypted, err := mongo.Connect(ctx, ceOpts)
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	var wg sync.WaitGroup
// 	for i := 0; i < 10; i++ {
// 		wg.Add(1)
// 		go func(i int) {
// 			coll := clientEncrypted.Database("db").Collection("coll")
// 			_, err = coll.InsertOne(ctx, bson.M{"_id": i, "encrypted": "string0"})
// 			if err != nil {
// 				log.Fatal(err)
// 			}

// 			res := coll.FindOne(ctx, bson.M{"_id": i})
// 			if res.Err() != nil {
// 				log.Fatal(res.Err())
// 			}
// 			// TODO: what is the right way to compare BSON.
// 			if raw, err := res.DecodeBytes(); err != nil {
// 				log.Fatal(res.Err())
// 			} else {
// 				expected, _ := bson.Marshal(bson.D{{"_id", i}, {"encrypted", "string0"}})
// 				if bytes.Compare(expected, raw) != 0 {
// 					log.Fatal("not equal")
// 				}
// 			}
// 			wg.Done()
// 		}(i)
// 	}

// 	wg.Wait()
// }
