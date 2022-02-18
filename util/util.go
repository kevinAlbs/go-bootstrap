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

type MonitorOpts struct {
	StartedOnly bool
}

func CreateMonitor(opts ...*MonitorOpts) *event.CommandMonitor {
	var mopts MonitorOpts
	for _, mopt := range opts {
		mopts.StartedOnly = mopt.StartedOnly
	}
	if mopts.StartedOnly {
		return &event.CommandMonitor{Started: logCommandStarted}
	}
	return &event.CommandMonitor{Started: logCommandStarted, Succeeded: logCommandSucceeded, Failed: logCommandFailed}
}

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

func serverDescriptionChanged(event *event.ServerDescriptionChangedEvent) {
	log.Printf("ServerDescriptionChangedEvent: %v, %v => %v\n", event.Address, event.PreviousDescription.Kind, event.NewDescription.Kind)
}
func serverOpening(event *event.ServerOpeningEvent) {
	log.Println("ServerOpeningEvent", event)
}
func serverClosed(event *event.ServerClosedEvent) {
	log.Printf("ServerClosedEvent", event)
}
func topologyDescriptionChanged(event *event.TopologyDescriptionChangedEvent) {
	log.Printf("TopologyDescriptionChangedEvent: %v => %v\n", event.PreviousDescription.Kind, event.NewDescription.Kind)
}
func topologyOpening(event *event.TopologyOpeningEvent) {
	log.Println("TopologyOpeningEvent")
}
func topologyClosed(event *event.TopologyClosedEvent) {
	log.Println("TopologyClosedEvent")
}
func serverHeartbeatStarted(event *event.ServerHeartbeatStartedEvent) {
	log.Println("ServerHeartbeatStartedEvent", event)
}
func serverHeartbeatSucceeded(event *event.ServerHeartbeatSucceededEvent) {
	log.Println("ServerHeartbeatSucceededEvent", event)
}
func serverHeartbeatFailed(event *event.ServerHeartbeatFailedEvent) {
	log.Println("ServerHeartbeatFailedEvent", event)
}

func CreateServerMonitor() *event.ServerMonitor {
	return &event.ServerMonitor{serverDescriptionChanged, serverOpening, serverClosed, topologyDescriptionChanged, topologyOpening, topologyClosed, serverHeartbeatStarted, serverHeartbeatSucceeded, serverHeartbeatFailed}
}
