// Run with: go run ./investigations/apiVersionWithExplain
package main

import (
	"os"

	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"kevinalbs.com/go-bootstrap/util"
)

func main() {
	var uri string
	if uri = os.Getenv("MONGODB_URI"); uri == "" {
		uri = "mongodb://localhost:27017"
	}

	copts := options.Client().
		ApplyURI(uri).
		SetServerAPIOptions(options.ServerAPI(options.ServerAPIVersion1)).
		SetMonitor(util.CreateMonitor())

	client, err := mongo.Connect(context.TODO(), copts)
	if err != nil {
		panic(err)
	}
	res := client.Database("db").RunCommand(context.TODO(), bson.D{
		{"explain", bson.D{{"find", "foo"}}},
	})
	if res.Err() != nil {
		panic(res.Err())
	}
}

/*
The outgoing command is:

{
    "explain": {
        "find": "foo"
    },
    "lsid": {
        "id": {
            "$binary": {
                "base64": "yNsoEvj4TlOVMqrJWCgdig==",
                "subType": "04"
            }
        }
    },
    "$clusterTime": {
        "clusterTime": {
            "$timestamp": {
                "t": 1671728629,
                "i": 1
            }
        },
        "signature": {
            "hash": {
                "$binary": {
                    "base64": "AAAAAAAAAAAAAAAAAAAAAAAAAAA=",
                    "subType": "00"
                }
            },
            "keyId": {
                "$numberLong": "0"
            }
        }
    },
    "apiVersion": "1",
    "$db": "db",
    "$readPreference": {
        "mode": "primary"
    }
}

*/
