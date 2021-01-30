package investigations

import (
	"fmt"
	"log"
	"strings"
	"testing"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/bsonrw"
)

type StructAB struct {
	A string
	B string
}

type StructABC struct {
	A string
	B string
	C StructAB
}

func Test1235_Problem(t *testing.T) {
	// Let's do some basic things. Let's unmarshall an extended JSON string into a struct.
	var d1, d2 StructABC
	full := ([]byte)(`{"a": "a", "b": "b", "c": {"a": "c", "b": "c"}}`)
	err := bson.UnmarshalExtJSON(full, true, &d1)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("UnmarshalExtJSON of StructABC", d1)

	// The problem, when the reader skips over an extra document/array field in a nested document, an error occurs.
	extraneous := ([]byte)(`{"a": "a", "b": "b", "extra": {"x": {"y": 1}}, "c": {"a": "c", "b": "c"}}`)

	err = bson.UnmarshalExtJSON(extraneous, true, &d2)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("UnmarshalExtJSON of StructABC", d2)
}

func Test1235_OriginalProblem(t *testing.T) {
	reader, err := bsonrw.NewExtJSONValueReader(strings.NewReader(`{"a": 123, "b": { "c": 2, "d": 3, "e": {"x": 1} }, "c": 456}`), true)
	if err != nil {
		log.Fatal(err)
	}
	dr, err := reader.ReadDocument()
	if err != nil {
		log.Fatal(err)
	}
	// Now we have a document reader.
	key, vr, err := dr.ReadElement()
	fmt.Println("read", key)
	if err != nil {
		log.Fatal(err)
	}
	vr.Skip()

	key, vr, err = dr.ReadElement()
	fmt.Println("read", key)
	if err != nil {
		log.Fatal(err)
	}

	if err := vr.Skip(); err != nil {
		log.Fatal(err)
	}

	fmt.Println(vr.Type())
}
