package common_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"testing"

	"go.mongodb.org/mongo-driver/bson/bsonrw"

	"github.com/kr/pretty"
	"go.mongodb.org/mongo-driver/bson"
)

/*
=== RUN   TestUnmarshalBson
recv:unmarshal_test.receiver{
    Total: primitive.A{
        primitive.D{
            {
                Key:   "count",
                Value: int64(100),
            },
        },
    },
    Data: primitive.A{
        primitive.D{
            {
                Key:   "name",
                Value: "Alice",
            },
        },
    },
}
a:[]unmarshal_test.Total{
}
b:[]unmarshal_test.Data{
}
--- PASS: TestUnmarshalBson (0.00s)
=== RUN   TestUnmarshalJson
recv:unmarshal_test.receiver{
    Total: &[]unmarshal_test.Total{
        {Count:0x64},
    },
    Data: &[]unmarshal_test.Data{
        {Name:"Alice"},
    },
}
a:[]unmarshal_test.Total{
    {Count:0x64},
}
b:[]unmarshal_test.Data{
    {Name:"Alice"},
}
--- PASS: TestUnmarshalJson (0.00s)
PASS
*/

type Total struct {
	Count uint64 `bson:"count"`
}

type Data struct {
	Name string `bson:"name"`
}

type Demo struct {
	Total []Total `bson:"total"`
	Data  []Data  `bson:"data"`
}

var new = Demo{
	Total: []Total{
		{
			Count: 100,
		},
	},
	Data: []Data{
		{
			"Alice",
		},
	},
}

type receiver struct {
	Total interface{}
	Data  interface{}
}

type src struct {
	X int `bson:"x"`
}
type dstWithPtr struct {
	X *int32
}
type dstWithEmptyInterface struct {
	X interface{}
}

func TestMine(t *testing.T) {
	// Can you unmarshal anything into a pointer type?
	s := src{1}
	bytes, _ := bson.Marshal(s)
	j, _ := bson.MarshalExtJSON(bson.Raw(bytes), false, false)
	fmt.Printf("marshaled to: %v\n", string(j))
	var storage int32
	var d dstWithPtr
	d.X = &storage
	bson.Unmarshal(bytes, &d)
	fmt.Println("unmarshaled to dstWithPtr:")
	fmt.Printf("*X=%v\n", *d.X)
	fmt.Printf("storage=%v\n", storage)

	// reset
	storage = 0

	var d2 dstWithEmptyInterface
	d2.X = "foo"
	bson.Unmarshal(bytes, &d2)
	fmt.Println("unmarshaled to dstWithEmptyInterface:")
	// Note, we don't look at the value type of X.
	// The interace{}'s type of d2.X is int32 instead of *int32.
	fmt.Printf("X=%v\n", d2.X.(int32))
	fmt.Printf("storage=%v\n", storage)
}

type To struct {
	X interface{}
}

func TestDecodingIntoEmptyInterface(t *testing.T) {
	val := struct{ X int32 }{123}
	bytes, err := bson.Marshal(val)
	if err != nil {
		t.Fatal(err)
	}

	// Now, try to unmarshal into an empty interface.
	var t1 To
	t1.X = "blah blah blah"

	err = bson.Unmarshal(bytes, &t1)
	if err != nil {
		t.Fatal(err)
	}

	fmt.Println(t1)
}

func TestEncoder(t *testing.T) {
	var buf bytes.Buffer
	vw, _ := bsonrw.NewBSONValueWriter(&buf)
	enc, _ := bson.NewEncoder(vw)
	enc.Encode(bson.D{{"x", 123}})
	buf.WriteTo(os.Stdout)
}

func TestUnmarshalBson(t *testing.T) {
	a := make([]Total, 0)
	b := make([]Data, 0)
	recv := receiver{
		Total: &a,
		Data:  &b,
	}
	buf, err := bson.Marshal(new)
	if err != nil {
		log.Fatal(err)
	}
	j, _ := bson.MarshalExtJSON(bson.Raw(buf), false, false)
	fmt.Printf("as extended JSON: %v\n", string(j))
	err = bson.Unmarshal(buf, &recv)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("recv:%# v\n", pretty.Formatter(recv))
	fmt.Printf("a:%# v\n", pretty.Formatter(a))
	fmt.Printf("b:%# v\n", pretty.Formatter(b))
}

func TestUnmarshalJson(t *testing.T) {
	a := make([]Total, 0)
	b := make([]Data, 0)
	recv := receiver{
		Total: &a,
		Data:  &b,
	}
	buf, err := json.Marshal(new)
	if err != nil {
		log.Fatal(err)
	}
	err = json.Unmarshal(buf, &recv)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("recv:%# v\n", pretty.Formatter(recv))
	fmt.Printf("a:%# v\n", pretty.Formatter(a))
	fmt.Printf("b:%# v\n", pretty.Formatter(b))
}
