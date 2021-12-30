package main

import (
	"os"

	"context"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"kevinalbs.com/go-bootstrap/util"
)

func main() {
	var uri string
	if uri = os.Getenv("MONGODB_URI"); uri == "" {
		uri = "mongodb://localhost:27017"
	}
	awsCreds := options.Credential{
		Username:      os.Getenv("TEST_AWS_TEMP_ACCESS_KEY_ID"),
		Password:      os.Getenv("TEST_AWS_TEMP_SECRET_ACCESS_KEY"),
		AuthMechanism: "MONGODB-AWS",
		AuthMechanismProperties: map[string]string{
			"AWS_SESSION_TOKEN": os.Getenv("TEST_AWS_TEMP_SESSION_TOKEN"),
		},
		AuthSource: "$external",
	}
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(uri).SetAuth(awsCreds).SetMonitor(util.CreateMonitor()))
	if err != nil {
		panic(err)
	}
	err = client.Ping(context.TODO(), nil)
	if err != nil {
		panic(err)
	}
}
