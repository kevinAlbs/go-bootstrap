package investigations

import (
	"context"
	"fmt"
	"log"
	"testing"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/event"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"kevinalbs.com/go-bootstrap/util"
)

var uri = "foo"

func TestCloudp92708(t *testing.T) {
	var connsCheckedOut int
	clientOpts := options.Client().ApplyURI(uri)
	clientOpts.SetMonitor(util.CreateMonitor())
	clientOpts.SetPoolMonitor(&event.PoolMonitor{
		Event: func(evt *event.PoolEvent) {
			fmt.Println(evt)
			switch evt.Type {
			case event.GetSucceeded:
				connsCheckedOut++
			case event.ConnectionReturned:
				connsCheckedOut--
			}
		},
	})
	client, err := mongo.Connect(context.Background(), clientOpts)
	if err != nil {
		log.Println(err)
		return
	}
	coll := client.Database("test").Collection("collection")
	coll.Drop(context.Background())
	docs := []interface{}{
		bson.D{{"_id", 1}},
		bson.D{{"_id", 2}},
		bson.D{{"_id", 3}},
		bson.D{{"_id", 4}},
	}
	_, err = coll.InsertMany(context.Background(), docs)
	if err != nil {
		log.Println(err)
		return
	}

	err = client.UseSession(context.Background(), func(sctx mongo.SessionContext) error {
		err = sctx.StartTransaction()
		if err != nil {
			return err
		}

		_, err = coll.Find(context.Background(), bson.D{}, options.Find().SetBatchSize(3))
		if err != nil {
			return err
		}

		return sctx.CommitTransaction(context.Background())
	})
	if err != nil {
		log.Println(err)
		return
	}

	log.Printf("remaining connections checked out: %v\n", connsCheckedOut)
}
