// This is a simplified HTTP service with three endpoints:
// /chatPost/?msg=foo to post a chat message.
// /chatGet to get the full chat.
// /chatClear to clear the chat.
package main

import (
	"os"
	"time"

	"context"

	"fmt"
	"net/http"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var mongoClient *mongo.Client

type Msg struct {
	Msg string
}

// chatPost adds msg to the chat.
func chatPost(ctx context.Context, msg string) error {
	coll := mongoClient.Database("chat").Collection("msg")
	_, err := coll.InsertOne(ctx, Msg{Msg: msg})
	return err
}

// chatGet gets the full chat.
func chatGet(ctx context.Context) (string, error) {
	var msgs []Msg
	coll := mongoClient.Database("chat").Collection("msg")
	cursor, err := coll.Find(ctx, bson.D{})
	if err != nil {
		return "", err
	}
	err = cursor.All(ctx, &msgs)
	if err != nil {
		return "", err
	}
	all_msgs := ""
	for _, msg := range msgs {
		all_msgs += "\n" + msg.Msg
	}
	return all_msgs, nil
}

// chatClear clears the chat.
func chatClear(ctx context.Context) error {
	coll := mongoClient.Database("chat").Collection("msg")
	_, err := coll.DeleteMany(ctx, bson.D{})
	if err != nil {
		return err
	}
	return nil
}

func chatPostHandler(w http.ResponseWriter, req *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	msg, ok := req.URL.Query()["msg"]
	if !ok || len(msg) < 0 {
		fmt.Fprintf(w, "'msg' query parameter not passed. Try http://localhost:8090/chatPost?msg=foo")
		return
	}
	err := chatPost(ctx, msg[0])
	if err != nil {
		fmt.Fprintf(w, "error posting: %v", err)
		return
	}
	fmt.Fprintf(w, "posted %q", msg[0])
}

func chatGetHandler(w http.ResponseWriter, req *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	all_msgs, err := chatGet(ctx)
	if err != nil {
		fmt.Fprintf(w, "error getting: %v", err)
		return
	}
	fmt.Fprintf(w, all_msgs)
}

func chatClearHandler(w http.ResponseWriter, req *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	err := chatClear(ctx)
	if err != nil {
		fmt.Fprintf(w, "error clearing: %v", err)
		return
	}
	fmt.Fprintf(w, "cleared messages")
}

func main() {
	var uri string
	if uri = os.Getenv("MONGODB_URI"); uri == "" {
		uri = "mongodb://localhost:27017"
	}
	c, err := mongo.Connect(context.Background(), options.Client().ApplyURI(uri))
	mongoClient = c
	defer c.Disconnect(context.Background())
	if err != nil {
		panic(err)
	}

	http.HandleFunc("/chatPost", chatPostHandler)
	http.HandleFunc("/chatGet", chatGetHandler)
	http.HandleFunc("/chatClear", chatClearHandler)
	fmt.Println("Listening on localhost:8090")
	http.ListenAndServe(":8090", nil)
}
