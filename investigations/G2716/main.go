package main

import (
	"log"
	"reflect"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/bsoncodec"
	"go.mongodb.org/mongo-driver/bson/bsonoptions"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func main() {
	type Foo struct {
		MyBytes []byte
	}

	f := Foo{MyBytes: nil}

	// Marshal without a Registry. An empty byte slice encodes to `nil`.
	{
		got, err := bson.Marshal(f)
		if err != nil {
			log.Fatalf("error in Marshal: %v", err)
		}
		var decoded bson.M
		bson.Unmarshal(got, &decoded)

		if decoded["mybytes"] != nil {
			log.Fatalf("expected nil, got: %v", decoded["mybytes"])
		}
	}

	// Marshal with a Registry with a codec to encode Nil as an empty bytes value.
	{
		r := bson.NewRegistryBuilder().RegisterCodec(reflect.TypeOf([]byte{}), bsoncodec.NewByteSliceCodec(bsonoptions.ByteSliceCodec().SetEncodeNilAsEmpty(true))).Build()
		got, err := bson.MarshalWithRegistry(r, f)
		if err != nil {
			log.Fatalf("error in Marshal: %v", err)
		}
		var decoded bson.M
		bson.Unmarshal(got, &decoded)

		if _, ok := decoded["mybytes"].(primitive.Binary); !ok {
			log.Fatalf("expected primitive.Binary, got: %T", decoded["mybytes"])
		}
	}

	// Marshal with EncodeContext
	{
		ectx := bsoncodec.EncodeContext{Registry: bson.DefaultRegistry, NilByteSliceAsEmpty: true}
		got, err := bson.MarshalWithContext(ectx, f)
		if err != nil {
			log.Fatalf("error in Marshal: %v", err)
		}
		var decoded bson.M
		bson.Unmarshal(got, &decoded)

		if _, ok := decoded["mybytes"].(primitive.Binary); !ok {
			log.Fatalf("expected primitive.Binary, got: %T", decoded["mybytes"]) // Results in 'expected primitive.Binary, got: <nil>'
		}
	}

}
