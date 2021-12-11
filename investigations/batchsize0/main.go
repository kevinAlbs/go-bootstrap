// Run with: go run ./investigations/batchsize0
// After running, use db.adminCommand({currentOp: 1}) to see command running server-side.
package main

import (
	"fmt"
	"os"
	"time"

	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"kevinalbs.com/go-bootstrap/util"
)

func runAggregate(coll *mongo.Collection) error {
	// Simulate an expensive pipeline by using sleep(30 * 1000)
	pipeline := bson.D{
		{
			"$addFields",
			bson.D{
				{
					"bar",
					bson.D{
						{
							"$function",
							bson.D{
								{"body", primitive.JavaScript("sleep(30 * 1000); return 'bar';")},
								{"args", bson.A{}},
								{"lang", "js"},
							},
						},
					},
				},
			},
		},
	}

	// Set a context deadline of one second.
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*1)
	defer cancel()
	cursor, err := coll.Aggregate(ctx, []bson.D{pipeline})
	if err != nil {
		return err
	}
	defer cursor.Close(ctx)
	for cursor.Next(ctx) {
		var out bson.M
		cursor.Decode(&out)
		fmt.Printf("Got back: %v\n", out)
	}
	return nil
}

func main() {
	var uri string
	if uri = os.Getenv("MONGODB_URI"); uri == "" {
		uri = "mongodb://localhost:27017"
	}
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(uri).SetMonitor(util.CreateMonitor()))
	if err != nil {
		panic(err)
	}
	coll := client.Database("db").Collection("coll")
	err = coll.Drop(context.TODO())
	if err != nil {
		panic(err)
	}

	docs := []interface{}{bson.D{{"_id", 1}}, bson.D{{"_id", 2}}, bson.D{{"_id", 3}}}
	_, err = coll.InsertMany(context.TODO(), docs)
	if err != nil {
		panic(err)
	}

	// Run an expensive aggregate with a one second timeout.
	err = runAggregate(coll)
	if err != nil {
		fmt.Printf("aggregate got error: %v\n", err)
	}

}
