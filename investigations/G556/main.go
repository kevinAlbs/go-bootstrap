package main

import (
	"fmt"
	"os"
	"sync"

	"context"

	"go.mongodb.org/mongo-driver/event"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/description"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	var uri string
	if uri = os.Getenv("MONGODB_URI"); uri == "" {
		uri = "mongodb://localhost:27017"
	}

	var latestTD *description.Topology
	var lockTD sync.RWMutex

	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(uri).SetServerMonitor(
		&event.ServerMonitor{
			TopologyDescriptionChanged: func(tdce *event.TopologyDescriptionChangedEvent) {
				lockTD.Lock()
				defer lockTD.Unlock()
				latestTD = &tdce.NewDescription
			},
		},
	))
	if err != nil {
		panic(err)
	}

	err = client.Ping(context.TODO(), nil)
	if err != nil {
		panic(err)
	}

	{
		lockTD.Lock()
		defer lockTD.Unlock()
		fmt.Printf("latestTD.Kind=%v\n", latestTD.Kind)
	}
}
