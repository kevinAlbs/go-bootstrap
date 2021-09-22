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
	opts := options.Client().ApplyURI(uri).
		SetPoolMonitor(util.CreatePoolMonitor()).
		SetMonitor(util.CreateMonitor()).
		SetServerMonitor(util.CreateServerMonitor())

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
			ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			defer cancel()
			res := client.Database("test").RunCommand(ctx, bson.D{{"ping", 1}})
			if res.Err() != nil {
				log.Fatal(res.Err())
			}
			log.Printf("iteration %v", iteration)
			time.Sleep(1 * time.Second)
			iteration++
		}
		wg.Done()
	}()
	wg.Wait()
}
