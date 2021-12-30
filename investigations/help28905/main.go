package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readconcern"
	"go.mongodb.org/mongo-driver/mongo/writeconcern"
)

// go.mongodb.org/mongo-driver v1.7.1

func main() {
	ctx := context.Background()

	err := json.Unmarshal(ActionsData, &actDef)
	if err != nil {
		log.Fatalf("%v", err)
		panic(err)
	}

	mdbOpts := options.Client()
	mdbOpts.ApplyURI(os.Getenv("MONGODB_URI"))

	client, err := mongo.Connect(ctx, mdbOpts)
	if err != nil {
		panic(err)
	}

	dbName := "sxodb8"
	// create DB
	db := client.Database(dbName)
	defer db.Client().Disconnect(ctx)

	// populate collection options
	ops := options.Collection()
	ops.SetReadConcern(readconcern.Majority())
	ops.SetWriteConcern(writeconcern.New(writeconcern.W(2), writeconcern.J(false)))
	actColl := db.Collection("definitionactivities", ops)

	// run for number of tenants
	for tenant := 90; tenant < 100; tenant++ {
		fmt.Printf("Start operation for tenant: %v \n", tenant)

		var prevActionUN string
		// create actions in given tenant
		for i, act := range actDef.Actions {

			fmt.Printf("Start new session for tenant: %v, for action: %s \n", tenant, act.UniqueName)
			// Step1: start new session
			sess, err := client.StartSession()
			if err != nil {
				fmt.Printf("StartSession %v", err)
				return
			}

			// update ctx mongo session context
			sessCtx := mongo.NewSessionContext(context.Background(), sess)

			// Start a transaction and sessCtx as the Context parameter to InsertOne and FindOne so both operations will be
			// run in the transaction.
			if err = sess.StartTransaction(options.Transaction().SetReadConcern(readconcern.Majority()).SetWriteConcern(writeconcern.New(writeconcern.W(2), writeconcern.J(false)))); err != nil {
				fmt.Printf("StartTransaction: tenant:%v err %v", tenant, err)
				return
			}
			// unmarshal result
			rslt := primitive.M{}
			if i > 0 {
				fmt.Printf("Read previous created action:%s from tenant: %v \n", prevActionUN, tenant)
				filter := map[string]interface{}{
					"tenant_id":   strconv.Itoa(tenant),
					"unique_name": prevActionUN,
				}
				if err := actColl.FindOne(sessCtx, filter).Decode(&rslt); err != nil {
					fmt.Printf("FindOne: tenant: %v, filter: %v -> err %v", tenant, filter, err)
					_ = abortTransaction(sessCtx)
					return
				}
			}
			act.TenantID = strconv.Itoa(tenant)
			act.CreatedOn = time.Now().UTC()
			_, err = actColl.InsertOne(sessCtx, act)
			if err != nil {
				fmt.Printf("CreateDocument: create error: %v, %T", err, err)
				_ = abortTransaction(sessCtx)
				return
			}

			prevActionUN = act.UniqueName
			commitTransaction(sessCtx, tenant)
		}
		fmt.Printf("Ended operation for tenant: %v", tenant)
	}
}

// AbortTransaction aborts the current transaction
func abortTransaction(ctx context.Context) error {
	fmt.Println("abortTransaction ...")
	sess := mongo.SessionFromContext(ctx)
	if sess == nil {
		fmt.Print("abortTransaction: transaction not active - abort ignored")
		return nil
	}

	defer sess.EndSession(ctx)

	err := sess.AbortTransaction(context.Background())
	if err != nil {
		return fmt.Errorf("abortTransaction error: %w", err)
	}
	return nil
}

// CommitTransaction commits the current transaction
func commitTransaction(ctx context.Context, tenant int) error {
	fmt.Printf("commitTransaction for tenant: %v \n", tenant)
	sess := mongo.SessionFromContext(ctx)
	if sess == nil {
		fmt.Print("commitTransaction: transaction not active - abort ignored")
		return nil
	}

	defer sess.EndSession(ctx)

	err := sess.CommitTransaction(context.Background())
	if err != nil {
		return err
	}

	fmt.Printf("commitTransaction: sess transaction state => %v \n", sess.(mongo.XSession).ClientSession().TransactionState)
	fmt.Printf("commitTransaction: session ID: %v, transaction number: %v \n", sess.(mongo.XSession).ClientSession().SessionID, sess.(mongo.XSession).ClientSession().TxnNumber)

	return nil
}

type Definition struct {
	Rev        int64     `json:"_rev" bson:"_rev"`
	TenantID   string    `json:"tenant_id" bson:"tenant_id"`
	Name       string    `json:"name,omitempty" bson:"name,omitempty"`
	Title      string    `json:"title,omitempty" bson:"title,omitempty"`
	Type       string    `json:"type" bson:"type"`
	BaseType   string    `json:"base_type" bson:"base_type"`
	CreatedOn  time.Time `json:"created_on" bson:"created_on"`
	UniqueName string    `json:"unique_name" bson:"unique_name"`
	ObjectType string    `json:"object_type" bson:"object_type"`
}

type DataDef struct {
	Actions []*Definition `json:"actions"`
}

//var actDef Definition
var actDef *DataDef

var ActionsData = []byte(` {  "actions": [
	{
	  "unique_name": "definition_activity_01HZWTH71OQ4X2FSX0cDaNcGqEl8WskfXIw",
	  "name": "JSONPath Query",
	  "title": "Extract Workflow ID",
	  "type": "corejava.jsonpathquery",
	  "base_type": "activity",
	  "object_type": "definition_activity"
	},
	{
	  "unique_name": "definition_activity_01HZWTJ0PVP120DXuhNptANKui2N6Hnhz8r",
	  "name": "HTTP Request",
	  "title": "Adding Upper Action  to Generic Worklfow",
	  "type": "web-service.http_request",
	  "base_type": "activity",
	  "object_type": "definition_activity"
	},
	{
	  "unique_name": "definition_activity_01HZWTKX7FBBP58lxVe2XOJSMpucTw7lazx",
	  "name": "HTTP Request",
	  "title": "Validate Workflow",
	  "type": "web-service.http_request",
	  "base_type": "activity",
	  "object_type": "definition_activity"
	},
	{
	  "unique_name": "definition_activity_01HZWTN737LPU3ne8AnzkjJqALBU84xUKqn",
	  "name": "HTTP Request",
	  "title": "Executing Workflow",
	  "type": "web-service.http_request",
	  "base_type": "activity",
	  "object_type": "definition_activity"
	},
	{
	  "unique_name": "definition_activity_01HZWTOMMC9IB2sTDOt9ZFqeZJpQP5466G9",
	  "name": "JSONPath Query",
	  "title": "Extract instance ID",
	  "type": "corejava.jsonpathquery",
	  "base_type": "activity",
	  "object_type": "definition_activity"
	},
	{
	  "unique_name": "definition_activity_01HZWTPZRSPIJ2EMyXTrjOVG6V8ylqHE5N4",
	  "name": "Sleep",
	  "title": "Sleep",
	  "type": "core.sleep",
	  "base_type": "activity",
	  "object_type": "definition_activity"
	},
	{
	  "unique_name": "definition_activity_01HZWTRHGOL371OCHRFBVK08XOhusfgwIOJ",
	  "name": "HTTP Request",
	  "title": "Getting Instance Information indetails",
	  "type": "web-service.http_request",
	  "base_type": "activity",
	  "object_type": "definition_activity"
	},
	{
	  "unique_name": "definition_activity_01HZWTT9X0Q3Y4c4Vy94Hc4U9xFSO6mHLtK",
	  "name": "JSONPath Query",
	  "title": "Extracting Instance and actions status",
	  "type": "corejava.jsonpathquery",
	  "base_type": "activity",
	  "object_type": "definition_activity"
	},
	{
	  "unique_name": "definition_activity_01HZWTV0HCDSF5rbQhkmk8LszsDKBHJQscE",
	  "name": "HTTP Request",
	  "title": "Delete Instance",
	  "type": "web-service.http_request",
	  "base_type": "activity",
	  "object_type": "definition_activity"
	},
	{
	  "unique_name": "definition_activity_01HZWTX4V7CAK0y0of11Le3BRuB0xuvQDIl",
	  "name": "HTTP Request",
	  "title": "Delete Atomic Workflow",
	  "type": "web-service.http_request",
	  "base_type": "activity",
	  "object_type": "definition_activity"
	}
  ]
  }`)
