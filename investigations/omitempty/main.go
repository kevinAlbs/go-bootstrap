// Run with:
// export MONGODB_URI="<URI to a serverless instance>"
// go run ./investigations/serverless_listcollections_batchsize

package main

import (
	"fmt"
	"log"

	"go.mongodb.org/mongo-driver/bson"
)

type IndexConfig struct {
	Options map[string]interface{} `json:"options,omitempty"`
}

func main() {
	fromjson := "{ \"options\": { \"$eq\": null } }"
	var cfg IndexConfig
	err := bson.UnmarshalExtJSON([]byte(fromjson), true /* canonical */, &cfg)
	if err != nil {
		log.Fatal(err)
	}
	tojson, err := bson.MarshalExtJSON(cfg, true /* canonical */, false /* escapeHTML */)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Result: %v\n", string(tojson))
	// Result: {"options":{"$eq":null}}
}
