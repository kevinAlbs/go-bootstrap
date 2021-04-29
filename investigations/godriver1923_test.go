package investigations

import (
	"context"
	"fmt"
	"testing"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"kevinalbs.com/go-bootstrap/util"
)

type user struct {
	Token string
	Query map[string][]string `bson:",inline"`
}

func TestEmbeddedNull(t *testing.T) {
	opts := options.Client().ApplyURI("mongodb://127.0.0.1:27017")
	opts.SetMonitor(util.CreateMonitor())
	client, err := mongo.NewClient(opts)
	if err != nil {
		panic(err)
	}

	if err := client.Connect(context.Background()); err != nil {
		panic(err)
	}

	coll := client.Database("test").Collection("test")
	// doc := bson.D{{"a%00%0d%00%00%00%02%30%00%01%00%00%00%00%00%04%62", "b"}}
	// doc2 := bson.M{"a\u0000": "b"}
	var u user
	u.Token = "abc"
	u.Query = make(map[string][]string)
	u.Query["blah\u0000\u0000foo"] = []string{"test", "this"}
	_, err = coll.InsertOne(context.TODO(), u)
	if err != nil {
		fmt.Println(err)
	}
}
