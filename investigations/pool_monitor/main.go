/*
This is an example of configuring a PoolMonitor.

Run this example with: go run ./investigations/pool_monitor

Sample output:
2022/04/08 16:43:03 PoolEvent: {"type":"ConnectionPoolCreated","address":"localhost:27017","connectionId":0,"options":{"maxPoolSize":100,"minPoolSize":0,"maxIdleTimeMS":0},"reason":"","serviceId":null}
2022/04/08 16:43:03 PoolEvent: {"type":"ConnectionPoolReady","address":"localhost:27017","connectionId":0,"options":null,"reason":"","serviceId":null}
2022/04/08 16:43:03 PoolEvent: {"type":"ConnectionCheckOutStarted","address":"localhost:27017","connectionId":0,"options":null,"reason":"","serviceId":null}
2022/04/08 16:43:03 PoolEvent: {"type":"ConnectionCreated","address":"localhost:27017","connectionId":1,"options":null,"reason":"","serviceId":null}
2022/04/08 16:43:03 PoolEvent: {"type":"ConnectionReady","address":"localhost:27017","connectionId":1,"options":null,"reason":"","serviceId":null}
2022/04/08 16:43:03 PoolEvent: {"type":"ConnectionCheckedOut","address":"localhost:27017","connectionId":1,"options":null,"reason":"","serviceId":null}
2022/04/08 16:43:03 PoolEvent: {"type":"ConnectionCheckedIn","address":"localhost:27017","connectionId":1,"options":null,"reason":"","serviceId":null}
2022/04/08 16:43:03 PoolEvent: {"type":"ConnectionCheckOutStarted","address":"localhost:27017","connectionId":0,"options":null,"reason":"","serviceId":null}
2022/04/08 16:43:03 PoolEvent: {"type":"ConnectionCheckedOut","address":"localhost:27017","connectionId":1,"options":null,"reason":"","serviceId":null}
2022/04/08 16:43:03 PoolEvent: {"type":"ConnectionCheckedIn","address":"localhost:27017","connectionId":1,"options":null,"reason":"","serviceId":null}
2022/04/08 16:43:03 PoolEvent: {"type":"ConnectionClosed","address":"localhost:27017","connectionId":1,"options":null,"reason":"poolClosed","serviceId":null}
2022/04/08 16:43:03 PoolEvent: {"type":"ConnectionPoolClosed","address":"localhost:27017","connectionId":0,"options":null,"reason":"","serviceId":null}
*/
package main

import (
	"encoding/json"
	"log"
	"os"

	"context"

	"go.mongodb.org/mongo-driver/event"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func logPoolEvent(event *event.PoolEvent) {
	asJson, err := json.Marshal(event)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("PoolEvent: %v\n", string(asJson))
}

func CreatePoolMonitor() *event.PoolMonitor {
	return &event.PoolMonitor{logPoolEvent}
}

func main() {
	var uri string
	if uri = os.Getenv("MONGODB_URI"); uri == "" {
		uri = "mongodb://localhost:27017"
	}
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(uri).SetPoolMonitor(CreatePoolMonitor()))
	defer client.Disconnect(context.Background())
	if err != nil {
		panic(err)
	}
	err = client.Ping(context.TODO(), nil)
	if err != nil {
		panic(err)
	}
}
