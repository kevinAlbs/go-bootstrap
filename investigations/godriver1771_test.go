package investigations

import (
	"context"
	"fmt"
	"log"
	"testing"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MyStruct struct {
	X primitive.ObjectID `bson:"_id"`
}

func TestGodriver1771(t *testing.T) {
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		fmt.Println(err)
		return
	}

	matchStage := bson.D{{"$match", bson.D{}}}
	cursor, err := client.Database("test").Collection("test").Aggregate(context.TODO(), mongo.Pipeline{matchStage})
	if err != nil {
		log.Fatalf("aggregate error: %v", err)
	}

	if !cursor.Next(context.TODO()) {
		log.Fatalf("no results")
	}

	res := &MyStruct{}
	err = cursor.Decode(&res)
	if err != nil {
		log.Fatalf("decode error: %v", err)
	}
}
