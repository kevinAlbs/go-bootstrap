package investigations

import (
	"context"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/jamiealquiza/tachymeter"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/event"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readconcern"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

var connectionsReady int32

func monitorConnections(opts *options.ClientOptions) {
	poolMonitor := &event.PoolMonitor{
		Event: func(evt *event.PoolEvent) {
			switch evt.Type {
			case event.ConnectionReady:
				val := atomic.AddInt32(&connectionsReady, 1)
				if val%100 == 0 {
					fmt.Printf("ConnectionsReady: %v", val)
				}
			case event.ConnectionCreated:

			}
		},
	}
	opts.SetPoolMonitor(poolMonitor)
}

func TestHelp26173(t2 *testing.T) {
	num := flag.Int("n", 1000, "concurrency")
	zzz := flag.Bool("sleep", true, "pause execution")
	flag.Parse()
	if *zzz {
		fmt.Println("Running with 50 microsecond delay")
	} else {
		fmt.Println("Running without 50 microsecond delay")
	}
	ctx := context.Background()
	opts := options.
		Client().
		SetReadConcern(readconcern.Local()).
		SetReadPreference(readpref.Nearest()).
		ApplyURI("mongodb://localhost:27017,localhost:27018,localhost:27019/?replicaSet=replset").SetMinPoolSize(100)
	monitorConnections(opts)
	mc, err := mongo.Connect(ctx, opts)
	if err != nil {
		log.Fatal("error connecting to mongo", err)
	}
	defer mc.Disconnect(ctx)
	time.Sleep(10 * time.Second)
	fmt.Println("about to start")
	_ = mc.Ping(ctx, readpref.Nearest())
	tlColl := mc.Database("test").Collection("tax_lot")
	users, err := tlColl.Distinct(ctx, "owner_uuid", bson.D{}, options.Distinct())
	if err != nil {
		log.Fatal("error finding all owner_uuids: ", err)
	}
	n := *num
	t := tachymeter.New(&tachymeter.Config{
		Size: n,
	})
	var wg sync.WaitGroup
	wallStart := time.Now()
	for i := 0; i < n; i++ {
		wg.Add(1)
		ownerUUID := users[rand.Intn(len(users))].(string)
		if *zzz {
			time.Sleep(50 * time.Microsecond)
		}
		go func(uuid string) {
			defer wg.Done()
			var p mongo.Pipeline
			p = mongo.Pipeline{
				{{
					Key: "$match",
					Value: bson.D{
						{Key: "owner_uuid", Value: uuid},
						{Key: "current_amount", Value: bson.D{{
							Key: "$gt", Value: 0.0,
						}}},
					},
				}},
				{{
					Key: "$group",
					Value: bson.D{
						{Key: "_id", Value: "$asset"},
						{Key: "position", Value: bson.D{
							{Key: "$sum", Value: "$current_amount"},
						}},
					},
				}},
			}
			start := time.Now()
			c, err := tlColl.Aggregate(ctx, p)
			if err != nil {
				log.Fatal("error aggregating: ", err)
			}
			var results []bson.M
			if err := c.All(ctx, &results); err != nil {
				log.Fatal(err)
			}
			t.AddTime(time.Since(start))
		}(ownerUUID)
	}
	wg.Wait()
	t.SetWallTime(time.Since(wallStart))
	fmt.Println(t.Calc().String())
	fmt.Printf("total connections opened: %v\n", connectionsCreated)
}
