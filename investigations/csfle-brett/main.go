package main

import (
	"C"
	"context"
	"fmt"
	"os"
	"time"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sts"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type SchemaSlice struct {
	deterministic [][]string
	random        [][]string
}

type AWSCredentials struct {
	AccessKeyId     *string
	SecretAccessKey *string
	SessionToken    *string
}

func getAWSToken() (*AWSCredentials, error) {
	cfg, err := config.LoadDefaultConfig(
		context.TODO(),
		config.WithSharedConfigProfile("Services.User-331472312345"),
	)
	if err != nil {
		return nil, err
	}

	// Create a STS client
	svc := sts.NewFromConfig(cfg)

	roleToAssumeArn := "arn:aws:iam::331472312345:role/ce-training-kms"
	sessionName := "test_session"
	var durationSeconds int32 = 3200
	result, err := svc.AssumeRole(context.TODO(), &sts.AssumeRoleInput{
		RoleArn:         &roleToAssumeArn,
		RoleSessionName: &sessionName,
		DurationSeconds: &durationSeconds,
	})

	if err != nil {
		return nil, err
	}

	return &AWSCredentials{
		AccessKeyId:     result.Credentials.AccessKeyId,
		SecretAccessKey: result.Credentials.SecretAccessKey,
		SessionToken:    result.Credentials.SessionToken,
	}, nil
}

func getAWSCredentialsFromEnv() (*AWSCredentials, error) {
	accessKeyId := os.Getenv("CSFLE_AWS_TEMP_ACCESS_KEY_ID")
	secretAccessKey := os.Getenv("CSFLE_AWS_TEMP_SECRET_ACCESS_KEY")
	sessionToken := os.Getenv("CSFLE_AWS_TEMP_SESSION_TOKEN")
	return &AWSCredentials{
		AccessKeyId:     &accessKeyId,
		SecretAccessKey: &secretAccessKey,
		SessionToken:    &sessionToken,
	}, nil
}

func createClient(c string) (*mongo.Client, error) {
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(c))

	if err != nil {
		return nil, err
	}

	return client, nil
}

func createManualEncryptionClient(c *mongo.Client, kp map[string]map[string]interface{}, kns string) (*mongo.ClientEncryption, error) {
	o := options.ClientEncryption().SetKeyVaultNamespace(kns).SetKmsProviders(kp)
	client, err := mongo.NewClientEncryption(c, o)
	if err != nil {
		return nil, err
	}

	return client, nil
}

func createEncryptedClient(cs string, kp map[string]map[string]interface{}, s bson.M, ks string) (*mongo.Client, error) {
	a := options.AutoEncryption().SetKeyVaultNamespace(ks).SetKmsProviders(kp).SetSchemaMap(s)

	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(cs).
		SetAutoEncryptionOptions(a))

	if err != nil {
		return nil, err
	}

	return client, nil
}

func encryptManual(ce *mongo.ClientEncryption, d primitive.Binary, a string, b interface{}) primitive.Binary {
	var out primitive.Binary
	rawValueType, rawValueData, err := bson.MarshalValue(b)
	if err != nil {
		panic(err)
	}

	rawValue := bson.RawValue{Type: rawValueType, Value: rawValueData}

	encryptionOpts := options.Encrypt().
		SetAlgorithm(a).
		SetKeyID(d)

	out, err = ce.Encrypt(context.TODO(), rawValue, encryptionOpts)
	if err != nil {
		panic(err)
	}

	return out
}

func main() {
	var (
		kmsProvider      map[string]map[string]interface{}
		keySpace         = "__secret.__keyvault"
		connectionString = "mongodb://localhost:27017"
		clientEncryption *mongo.ClientEncryption
		client           *mongo.Client
		exitCode         = 0
		result           *mongo.InsertOneResult
		findResult       bson.M
		dek              primitive.Binary
	)

	defer func() {
		os.Exit(exitCode)
	}()

	// role, err := getAWSToken()
	role, err := getAWSCredentialsFromEnv()
	if err != nil {
		fmt.Printf("Token error: %s\n", err)
		exitCode = 1
		return
	}

	kmsProvider = map[string]map[string]interface{}{
		"aws": {
			"accessKeyId":     role.AccessKeyId,
			"secretAccessKey": role.SecretAccessKey,
			"sessionToken":    role.SessionToken,
		},
	}

	client, err = createClient(connectionString)
	if err != nil {
		fmt.Printf("MDB client error: %s\n", err)
		exitCode = 1
		return
	}

	coll := client.Database("__secret").Collection("__keyvault")

	opts := options.FindOne().SetProjection(bson.D{{Key: "_id", Value: 1}})
	err = coll.FindOne(context.TODO(), bson.D{{Key: "keyAltNames", Value: "dataKey1"}}, opts).Decode(&findResult)
	if err != nil {
		fmt.Printf("DEK find error: %s\n", err)
		exitCode = 1
		return
	}
	dek = findResult["_id"].(primitive.Binary)

	clientEncryption, err = createManualEncryptionClient(client, kmsProvider, keySpace)
	if err != nil {
		fmt.Printf("ClientEncrypt error: %s\n", err)
		exitCode = 1
		return
	}
	defer func() { _ = clientEncryption.Close(context.TODO()) }()

	payload := bson.D{
		{Key: "_id", Value: 2315},
		{Key: "name", Value: bson.D{
			{Key: "firstname", Value: encryptManual(clientEncryption, dek, "AEAD_AES_256_CBC_HMAC_SHA_512-Deterministic", "Will")},
			{Key: "lastname", Value: encryptManual(clientEncryption, dek, "AEAD_AES_256_CBC_HMAC_SHA_512-Deterministic", "T")},
		}},
		{Key: "address", Value: bson.D{
			{Key: "streetAddress", Value: encryptManual(clientEncryption, dek, "AEAD_AES_256_CBC_HMAC_SHA_512-Random", "537 White Hills Rd")},
			{Key: "suburbCounty", Value: encryptManual(clientEncryption, dek, "AEAD_AES_256_CBC_HMAC_SHA_512-Random", "Evandale")},
			{Key: "zipPostcode", Value: "7258"},
			{Key: "stateProvince", Value: "Tasmania"},
			{Key: "country", Value: "Oz"},
		}},
		{Key: "dob", Value: encryptManual(clientEncryption, dek, "AEAD_AES_256_CBC_HMAC_SHA_512-Random", time.Date(1989, 1, 1, 0, 0, 0, 0, time.Local))},
		{Key: "phoneNumber", Value: encryptManual(clientEncryption, dek, "AEAD_AES_256_CBC_HMAC_SHA_512-Random", "+61 400 000 111")},
		{Key: "salary", Value: encryptManual(clientEncryption, dek, "AEAD_AES_256_CBC_HMAC_SHA_512-Random", bson.D{
			{Key: "current", Value: 99000.00},
			{Key: "startDate", Value: time.Date(2022, 6, 1, 0, 0, 0, 0, time.Local)},
			{Key: "history", Value: bson.D{
				{Key: "salary", Value: 89000.00},
				{Key: "startDate", Value: time.Date(2021, 8, 1, 0, 0, 0, 0, time.Local)},
			}},
		})},
		{Key: "taxIdentifier", Value: encryptManual(clientEncryption, dek, "AEAD_AES_256_CBC_HMAC_SHA_512-Random", "103-443-923")},
		{Key: "role", Value: []string{"IC"}},
	}
	fmt.Println(payload)

	coll = client.Database("companyData").Collection("employee")

	result, err = coll.InsertOne(context.TODO(), payload)
	if err != nil {
		fmt.Printf("Insert error: %s\n", err)
		exitCode = 1
		return
	}
	fmt.Print(result)

	exitCode = 0
}
