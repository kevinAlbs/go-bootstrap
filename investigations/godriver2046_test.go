package investigations

import (
	"fmt"
	"log"
	"testing"

	"go.mongodb.org/mongo-driver/bson"
)

type myStruct struct {
	F float32 `bson:"f,truncate"`
}

func Test2046(t *testing.T) {
	doc, err := bson.Marshal(bson.D{{"f", float64(.1)}})
	if err != nil {
		log.Fatalf("Marshal error: %v", err)
	}
	var dst myStruct
	err = bson.Unmarshal(doc, &dst)
	if err != nil {
		log.Fatalf("Unmarshal error: %v", err)
	}
	expectedF64 := float64(0.1)
	expectedF32 := float32(expectedF64)
	if dst.F != expectedF32 {
		fmt.Printf("Error: value of myStruct.F: %#+v != expected %#+v\n", dst.F, expectedF32)
	}
}
