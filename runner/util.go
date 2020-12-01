package main

import (
	"context"
	"fmt"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/event"
)

func logCommandStarted(context context.Context, event *event.CommandStartedEvent) {
	log.Println("=> CommandStarted:", event.CommandName, ToJson(event.Command))
}

func logCommandSucceeded(context context.Context, event *event.CommandSucceededEvent) {
	fmt.Println("<= CommandSucceeded:", event.CommandName, ToJson(event.Reply))
}

func logCommandFailed(context context.Context, event *event.CommandFailedEvent) {
	fmt.Println("<= CommandFailed:", event.CommandName, event.Failure)
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
