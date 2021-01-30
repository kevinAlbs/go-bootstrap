package investigations

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type test struct {
	ID  primitive.ObjectID `bson:"_id,omitempty"`
	Val int                `bson:"val,omitempty"`
}

func TestCommitTimeout(t *testing.T) {
	// Replace the uri string with your MongoDB deployment's connection string.
	uri := "mongodb://localhost:27017"
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		panic(err)
	}

	defer func() {
		if err = client.Disconnect(ctx); err != nil {
			panic(err)
		}
	}()

	// Db and collection
	database := client.Database("test")
	testCollection := database.Collection("test")

	// Ping the primary
	if err := client.Ping(ctx, readpref.Primary()); err != nil {
		panic(err)
	}

	session, _ := client.StartSession()
	defer session.EndSession(ctx)

	fmt.Println("Successfully connected and pinged.")

	err = mongo.WithSession(context.Background(), session, func(sessionContext mongo.SessionContext) error {
		// Start transaction
		if err = session.StartTransaction(); err != nil {
			panic(fmt.Sprintf("Panic starting transaction: %s\n", err))
		}

		time.Sleep(1 * time.Second)

		result, err := testCollection.InsertOne(sessionContext, test{
			ID:  primitive.NewObjectID(),
			Val: 17,
		})
		if err != nil {
			panic(fmt.Sprintf("Panic inserting one: %s\n", err))
		}

		fmt.Printf("Insert result %+v\n", result)

		commitTimeoutCtx, commitCancel := context.WithTimeout(sessionContext, 0)
		defer commitCancel()
		// Infinite loop is key here to reproduce bug in REALMC-7914
		for {
			// commitErr != nil after the below call but the session.clientSession.state == "Committed" so we cannot abort
			commitErr := session.CommitTransaction(commitTimeoutCtx)
			if commitErr == nil {
				break
			}

			fmt.Println(commitErr.(mongo.CommandError).HasErrorLabel("UnknownTransactionCommitResult"))

			if errors.Is(commitErr, context.DeadlineExceeded) {
				fmt.Println("context deadline exceeded error")
				return commitErr
			}

			if errors.Is(commitErr, context.Canceled) {
				fmt.Println("context canceled error")
				return commitErr
			}
		}

		return nil
	})
	if err != nil {
		if abortErr := session.AbortTransaction(context.Background()); abortErr != nil {
			// We fail to Abort with error "cannot call abortTransaction after calling commitTransaction
			panic(abortErr)
		}
		panic(err)
	}
}
