package investigations

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func TestError(t *testing.T) {
	ctx := context.TODO()
	client, err := mongo.Connect(ctx)
	assert.Nil(t, err)
	_, err = client.Database("db").Collection("coll").InsertOne(ctx, bson.D{{"_id", "foo"}})
	if cmdErr, ok := err.(mongo.CommandError); ok {
		fmt.Println(cmdErr.Code)
	}
	if err != nil {
		fmt.Println("Insert failed", err)
		var serverError mongo.ServerError

		if errors.As(err, &serverError) {
			if serverError.HasErrorLabel("RetryableWrite") {
				fmt.Println("Going to retry...")
			}
		}
	}
}
