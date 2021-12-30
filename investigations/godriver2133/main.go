package main

import (
	"context"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/event"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	uri = "mongodb://localhost:27017/?serverSelectionTimeoutMS=1000"
)

func main() {
	client := setupClient()
	defer client.Disconnect(context.TODO())

	coll := client.Database("foo").Collection("bar")
	if err := coll.Drop(context.TODO()); err != nil {
		panic(err)
	}

	docs := []interface{}{
		bson.M{"x": 1},
		bson.M{"x": 2},
	}
	if _, err := coll.InsertMany(context.TODO(), docs); err != nil {
		panic(err)
	}

	cursor, err := coll.Aggregate(context.TODO(), bson.D{})
	if err != nil {
		panic(err)
	}
	defer cursor.Close(context.TODO())

	var res []bson.M
	if err := cursor.All(context.TODO(), &res); err != nil {
		panic(err)
	}
	log.Println(res)
}

func setupClient() *mongo.Client {
	monitor := &event.CommandMonitor{
		Started: func(ctx context.Context, evt *event.CommandStartedEvent) {
			log.Printf("%s\n\n", evt.Command)
		},
	}
	opts := options.Client().ApplyURI(uri).SetMonitor(monitor)
	client, err := mongo.Connect(context.TODO(), opts)
	if err != nil {
		panic(err)
	}
	return client
}
