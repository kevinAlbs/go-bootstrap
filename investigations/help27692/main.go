package main

import (
	"context"
	"log"
	"os"
	"sync"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"kevinalbs.com/go-bootstrap/util"
)

func main() {
	var uri string
	if uri = os.Getenv("MONGODB_URI"); uri == "" {
		log.Fatal("You must set your 'MONGODB_URI' environmental variable.")
	}
	opts := options.Client().ApplyURI(uri).SetPoolMonitor(util.CreatePoolMonitor()).SetMonitor(util.CreateMonitor())

	client, err := mongo.Connect(context.TODO(), opts)
	if err != nil {
		log.Fatal(err)
	}

	defer func() {
		client.Disconnect(context.TODO())
	}()

	iteration := 1
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		for {
			client.Database("test").RunCommand(context.TODO(), bson.D{{"ping", 1}})
			log.Printf("iteration %v", iteration)
			time.Sleep(5 * time.Second)
			iteration++
		}
		wg.Done()
	}()
	wg.Wait()
}
