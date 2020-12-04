package common_test

import (
	"encoding/json"
	"fmt"
	"log"
	"testing"

	"go.mongodb.org/mongo-driver/bson"
)

type EmbedMe struct {
	A int32 `bson:"a" json:"a"`
}
type Outer struct {
	*EmbedMe
	B int32 `bson:"a" json:"a"`
}

func Test1777(t *testing.T) {
	// By default, encoding/json will overwrite fields of embedded structs.
	val := Outer{&EmbedMe{123}, 456}

	// Note, there is no encoding/json analog to the `inline` bson struct tag.

	log.Println("JSON", jsonMarshal(val))
	log.Println("BSON", bsonMarshal(val))
}

func jsonMarshal(val interface{}) string {
	res, err := json.Marshal(val)
	if err != nil {
		fmt.Println("jsonMarshal error", err)
	}

	return string(res)
}

func bsonMarshal(val interface{}) bson.Raw {
	res, err := bson.Marshal(val)
	if err != nil {
		fmt.Println("bsonMarshal error", err)
	}

	return bson.Raw(res)
}
