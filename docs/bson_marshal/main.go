package main

import (
	"fmt"
	"log"

	"go.mongodb.org/mongo-driver/bson"
)

type myStruct struct {
	Foo int32
}

func main() {
	doc, err := bson.Marshal(bson.D{{"foo", 1}})
	if err != nil {
		log.Fatalf("Marshal error: %v", err)
	}
	var dst myStruct
	err = bson.Unmarshal(doc, &dst)
	if err != nil {
		log.Fatalf("Unmarshal error: %v", err)
	}
	fmt.Printf("Unmarshalled field Foo: %v\n", dst.Foo)
}
