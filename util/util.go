package util

import (
	"context"
	"encoding/json"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/event"
)

func logCommandStarted(context context.Context, event *event.CommandStartedEvent) {
	log.Println("=> CommandStarted:", event.CommandName, ToJson(event.Command))
}

func logCommandSucceeded(context context.Context, event *event.CommandSucceededEvent) {
	log.Println("<= CommandSucceeded:", event.CommandName, ToJson(event.Reply))
}

func logCommandFailed(context context.Context, event *event.CommandFailedEvent) {
	log.Println("<= CommandFailed:", event.CommandName, event.Failure)
}

func ToJson(v interface{}) string {
	bytes, err := bson.MarshalExtJSON(v, true, false)
	if err != nil {
		log.Fatal(err)
	}
	return string(bytes)
}

func CreateMonitor() *event.CommandMonitor {
	return &event.CommandMonitor{logCommandStarted, logCommandSucceeded, logCommandFailed}
}

func logPoolEvent(event *event.PoolEvent) {
	asJson, err := json.Marshal(event)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("PoolEvent: %v ", string(asJson))
}

func CreatePoolMonitor() *event.PoolMonitor {
	return &event.PoolMonitor{logPoolEvent}
}
