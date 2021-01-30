package investigations

import (
	"context"
	"fmt"
	"log"
	"testing"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/event"

	"github.com/bxcodec/faker"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Model struct {
	A int64 `bson:"a,omitempty"`
}

func logCommandStarted(context context.Context, event *event.CommandStartedEvent) {
	if event.CommandName == "insert" {
		return
	}
	log.Println("=> CommandStarted:", event.CommandName, ToJson(event.Command))
}

func logCommandSucceeded(context context.Context, event *event.CommandSucceededEvent) {
	if event.CommandName == "insert" {
		return
	}
	fmt.Println("<= CommandSucceeded:", event.CommandName)
	if event.CommandName == "getMore" {
		fmt.Println(" getMore reply: ", ToJson(event.Reply))
	}
}

func logCommandFailed(context context.Context, event *event.CommandFailedEvent) {
	fmt.Println("<= CommandFailed:", event.CommandName, event.Failure)
}

func ToJson(v interface{}) string {
	bytes, err := bson.MarshalExtJSON(v, true, false)
	if err != nil {
		log.Fatal(err)
	}
	return string(bytes)
}

func createMonitor() *event.CommandMonitor {
	return &event.CommandMonitor{logCommandStarted, logCommandSucceeded, logCommandFailed}
}

func Test1831(t *testing.T) {
	opts := options.Client().ApplyURI("mongodb://127.0.0.1:27017")
	opts.SetMonitor(createMonitor())
	client, err := mongo.NewClient(opts)
	if err != nil {
		panic(err)
	}

	if err := client.Connect(context.Background()); err != nil {
		panic(err)
	}

	genData(client)
	listData(client)
}

func genData(cli *mongo.Client) {
	cli.Database("test").Collection("test").Drop(context.Background())
	var coll = cli.Database("test").Collection("test")
	for i := 0; i < 1001; i++ {
		m := Model{A: faker.RandomUnixTime()}
		if _, err := coll.InsertOne(context.Background(), m); err != nil {
			panic(err)
		}
	}
}

func listData(cli *mongo.Client) {
	var coll = cli.Database("test").Collection("test")
	var data = []Model{}
	var opts = options.Find().SetSkip(0).SetLimit(1000)
	if cursor, err := coll.Find(context.Background(), bson.M{}, opts); err != nil {
		panic(err)
	} else {
		if err = cursor.All(context.Background(), &data); err != nil {
			panic(err)
		}
	}
	println(len(data))
}
