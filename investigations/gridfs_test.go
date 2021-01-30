package investigations

import (
	"context"
	"fmt"
	"log"
	"testing"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/gridfs"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Cat struct {
	x int
}

type Cat2 Cat

func TestGridFS(t *testing.T) {
	ctx := context.Background()
	opts := options.Client().ApplyURI("mongodb://localhost:27017")
	client, _ := mongo.Connect(ctx, opts)
	db := client.Database("test")
	data := []byte("abc.def.ghi")
	var chunkSize int32 = 4

	var c Cat2
	c.x = 123

	testcases := []struct {
		firstRead          int
		skip               int64
		secondRead         int
		expectedSecondRead string
	}{
		{
			0, 1, 3, "bc.",
		},
		{
			3, 1, 3, "def",
		},
	}

	for _, tc := range testcases {
		bucket, _ := gridfs.NewBucket(db, options.GridFSBucket().SetChunkSizeBytes(chunkSize))

		ustream, _ := bucket.OpenUploadStream("foo")
		id := ustream.FileID
		ustream.Write(data)
		ustream.Close()

		dstream, _ := bucket.OpenDownloadStream(id)
		dst := make([]byte, tc.firstRead)
		_, err := dstream.Read(dst)
		if err != nil {
			log.Fatal(err)
		}

		dst = make([]byte, tc.secondRead)
		_, err = dstream.Skip(tc.skip)
		if err != nil {
			log.Fatal(err)
		}
		dstream.Read(dst)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("Second read: %q Expected: %q\n", string(dst), tc.expectedSecondRead)
	}
}

/*
Currently prints:
Second read: "abc" Expected: "bc."
Second read: ".de" Expected: "def"
*/
