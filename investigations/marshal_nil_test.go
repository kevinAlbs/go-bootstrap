package investigations

import (
	"fmt"
	"testing"

	"go.mongodb.org/mongo-driver/bson"
)

type Test struct {
	X []int
}

func TestMarshalNilSlice(t *testing.T) {
	var ts Test
	ts.X = []int{1, 2, 3}
	ts.X = []int{}

	bytes, err := bson.MarshalExtJSON(&ts, true, false)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(string(bytes))
}
