package investigations

import (
	"context"
	"fmt"
	"testing"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"kevinalbs.com/go-bootstrap/util"
)

var (
	collection *mongo.Collection
)

type ChangeEvent struct {
	FullDocument map[string]interface{}
}

type Temp struct {
	Field1 string
	Field2 string
	Field3 int
	Field4 map[string]int
	Field5 map[string]string
}

func Create() {
	t := Temp{
		Field1: "val1",
		Field3: 2,
		Field5: map[string]string{},
	}
	res, err := collection.InsertOne(context.Background(), t)
	if err != nil {
		panic(err)
	}
	fmt.Println(res.InsertedID)
}

func Update() {
	filter := bson.M{"field1": "val1"}
	t := Temp{
		Field1: "new hello!",
		Field2: "This shouldn't show up...",
		Field3: 4,
		Field5: map[string]string{},
	}
	_, err := collection.UpdateOne(context.Background(), filter,
		mongo.Pipeline{bson.D{{Key: "$set", Value: t}}}) // this works if we supply a bson.D, instead of a mongo.Pipeline
	if err != nil {
		panic(err)
	}
}

func Test1967(t *testing.T) {
	opts := options.Client().ApplyURI("mongodb://localhost:27017")
	opts.SetMonitor(util.CreateMonitor())
	client, err := mongo.Connect(context.Background(), opts)
	if err != nil {
		panic(err)
	}
	collection = client.Database("temp").Collection("temp")
	Create()
	Update()
	fmt.Println("success!")
}
