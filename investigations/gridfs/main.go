package main

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/gridfs"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	ctx := context.Background()
	opts := options.
		Client().
		ApplyURI("mongodb://localhost:27017")

	client, err := mongo.Connect(ctx, opts)
	if err != nil {
		panic(err)
	}
	db := client.Database("db")
	bucket, err := gridfs.NewBucket(db, options.GridFSBucket())
	if err != nil {
		panic(err)
	}

	ustream, err := bucket.OpenUploadStream("foo")
	if err != nil {
		panic(err)
	}

	// _ = ustream.FileID
	_, err = ustream.Write([]byte{1, 2, 3})
	if err != nil {
		panic(err)
	}
	err = ustream.Close()
	if err != nil {
		panic(err)
	}
}
