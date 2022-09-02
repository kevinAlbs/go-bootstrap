package main

// Are server error codes accessible from authentication errors?

import (
	"fmt"
	"os"

	"context"

	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/x/mongo/driver"
)

func GetErrorCode(err error) int {
	switch e := errors.Cause(err).(type) {
	case mongo.CommandError:
		return int(e.Code)
	case driver.Error:
		return int(e.Code)
	case driver.WriteCommandError:
		for _, we := range e.WriteErrors {
			return int(we.Code)
		}
		if e.WriteConcernError != nil {
			return int(e.WriteConcernError.Code)
		}
		return 0
	case driver.QueryFailureError:
		codeVal, err := e.Response.LookupErr("code")
		if err == nil {
			code, _ := codeVal.Int32OK()
			return int(code)
		}
		return 0 // this shouldn't happen
	case mongo.WriteError:
		return e.Code
	case mongo.BulkWriteError:
		return e.Code
	case mongo.WriteConcernError:
		return e.Code
	case mongo.WriteException:
		for _, we := range e.WriteErrors {
			return GetErrorCode(we)
		}
		if e.WriteConcernError != nil {
			return e.WriteConcernError.Code
		}
		return 0
	case mongo.BulkWriteException:
		// Return the first error code.
		for _, ecase := range e.WriteErrors {
			return GetErrorCode(ecase)
		}
		if e.WriteConcernError != nil {
			return e.WriteConcernError.Code
		}
		return 0
	default:
		return 0
	}
}

// GetErrorCodeV2 returns the MongoDB server error code on an error.
// Wrapped errors are also checked.
// Returns 0 if no MongoDB server error code is found.
func GetErrorCodeV2(err error) int {
	e := errors.Cause(err)

	var commandError mongo.CommandError
	if errors.As(e, &commandError) {
		return int(commandError.Code)
	}
	var driverError driver.Error
	if errors.As(e, &driverError) {
		return int(driverError.Code)
	}
	// TODO: add other error types.
	return 0
}

// unwrap returns the inner error if err implements Unwrap(), otherwise it returns nil.
func unwrap(err error) error {
	u, ok := err.(interface {
		Unwrap() error
	})
	if !ok {
		return nil
	}
	return u.Unwrap()
}

func main() {
	var uri string
	if uri = os.Getenv("MONGODB_URI"); uri == "" {
		uri = "mongodb://user:bad-pwd@localhost:27017"
	}
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(uri))
	if err != nil {
		panic(err)
	}
	// Use an invalid command.
	// If the username/password is incorrect, an authentication error is returned.
	// If the username/password is correct, an invalid command error is returned.
	res := client.Database("db").RunCommand(context.TODO(), bson.D{{"invalid", 1}})
	err = res.Err()
	if err != nil {
		got := GetErrorCode(err)
		fmt.Printf("GetErrorCode returned: %v\n", got)

		got = GetErrorCodeV2(err)
		fmt.Printf("GetErrorCodeV2 returned: %v\n", got)

		i := 0
		fmt.Printf("unwrapping errors ... begin\n")
		for ; err != nil; err = unwrap(err) {
			fmt.Printf("error %v is type %T: %v\n", i, err, err)
			i += 1
		}
		fmt.Printf("unwrapping errors ... end\n")
	}

}

/* Sample output:

GetErrorCode returned: 0
GetErrorCodeV2 returned: 18
unwrapping errors ... begin
error 0 is type topology.ConnectionError: connection() error occurred during connection handshake: auth error: unable to authenticate using mechanism "SCRAM-SHA-256": (AuthenticationFailed) Authentication failed.
error 1 is type *auth.Error: auth error: unable to authenticate using mechanism "SCRAM-SHA-256": (AuthenticationFailed) Authentication failed.
error 2 is type *auth.Error: unable to authenticate using mechanism "SCRAM-SHA-256": (AuthenticationFailed) Authentication failed.
error 3 is type driver.Error: (AuthenticationFailed) Authentication failed.
unwrapping errors ... end

*/
