package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"sync"
	"sync/atomic"
	"time"

	"github.com/jamiealquiza/tachymeter"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/event"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readconcern"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

/*
Hypothesis #1: The 50 microsecond delay helps with connection reuse.

I tested this by setting MinPoolSize=100 (the default MaxPoolSize) and waiting
for 3 * 100 ConnectionReady events before proceeding with the test.
However, enabling the 50 microsecond sleep still consistently improves the p95.

To reproduce:
./run-repeatedly.sh -minPoolSize 100
Running test 10 times with flags: -minPoolSize 100
70.771382ms
65.685711ms
60.857761ms
69.17513ms
71.254732ms
70.480517ms
56.138508ms
71.075244ms
78.001828ms
51.642829ms

./run-repeatedly.sh -minPoolSize 100 -sleep
Running test 10 times with flags: -minPoolSize 100 -sleep
4.421227ms
4.952742ms
4.733304ms
3.681267ms
5.495439ms
3.562335ms
19.728049ms
698.098µs
9.154847ms
11.795527ms

*/

var connectionsReady int32
var connectionsCreated int32

func monitorConnections(opts *options.ClientOptions) {
	poolMonitor := &event.PoolMonitor{
		Event: func(evt *event.PoolEvent) {
			switch evt.Type {
			case event.ConnectionReady:
				atomic.AddInt32(&connectionsReady, 1)
			case event.ConnectionCreated:
				atomic.AddInt32(&connectionsCreated, 1)
			}
		},
	}
	opts.SetPoolMonitor(poolMonitor)
}

func main() {
	num := flag.Int("n", 1000, "concurrency")
	zzz := flag.Bool("sleep", false, "pause execution")
	minPoolSize := flag.Int("minPoolSize", 0, "set min pool size")
	quiet := flag.Bool("quiet", false, "only print P95 metric")
	flag.Parse()
	if !*quiet {
		if *zzz {
			fmt.Println("Running with 50 microsecond delay")
		} else {
			fmt.Println("Running without 50 microsecond delay")
		}
	}
	ctx := context.Background()
	opts := options.
		Client().
		SetReadConcern(readconcern.Local()).
		SetReadPreference(readpref.Nearest()).
		ApplyURI("mongodb://localhost:27017,localhost:27018,localhost:27019/?replicaSet=replset")

	if *minPoolSize != 0 {
		if !*quiet {
			fmt.Printf("Setting MinPoolSize to %v\n", *minPoolSize)
		}
		opts.SetMinPoolSize(uint64(*minPoolSize))
	} else {
		if !*quiet {
			fmt.Println("Not setting MinPoolSize")
		}
	}

	monitorConnections(opts)
	mc, err := mongo.Connect(ctx, opts)
	if err != nil {
		log.Fatal("error connecting to mongo", err)
	}
	defer mc.Disconnect(ctx)
	if *minPoolSize != 0 {
		// Wait until minPoolSize * numServers is connected.
		numToWaitFor := int32(3 * *minPoolSize)
		if !*quiet {
			fmt.Printf("Waiting for %v connections to become ready before proceeding\n", numToWaitFor)
		}
		for atomic.LoadInt32(&connectionsReady) < numToWaitFor {
			time.Sleep(1 * time.Second)
		}
	}
	if !*quiet {
		fmt.Println("Proceeding with test")
	}
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
	metrics := t.Calc()
	if !*quiet {
		fmt.Println(metrics.String())
		fmt.Printf("Total connectionsCreated: %v\n", connectionsCreated)
		fmt.Printf("Total connectionsReady: %v\n", connectionsReady)
	}
	fmt.Printf("%v\n", metrics.Time.P95)
}
