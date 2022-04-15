package main

import (
	"fmt"
	"io"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/x/bsonx/bsoncore"
)

type MyReader struct {
	data  []byte
	index int
}

func (mr *MyReader) Read(p []byte) (n int, err error) {
	if mr.index >= len(mr.data) {
		return 0, io.EOF
	}
	want := len(p)
	written := 0
	for i := 0; i < want; i++ {
		p[i] = mr.data[mr.index]
		written += 1
		mr.index += 1
		if mr.index >= len(mr.data) {
			return written, io.EOF
		}
	}
	return written, nil
}

func newMyReader() *MyReader {
	var mr MyReader
	mr.data = make([]byte, 0)
	// Append two BSON documents side by side.
	b1, err := bson.Marshal(bson.D{{"x", 1}})
	if err != nil {
		log.Fatalf("error in bson.Marshal: %v", err)
	}
	for _, b := range b1 {
		mr.data = append(mr.data, b)
	}
	b2, err := bson.Marshal(bson.D{{"y", 1}})
	if err != nil {
		log.Fatalf("error in bson.Marshal: %v", err)
	}
	for _, b := range b2 {
		mr.data = append(mr.data, b)
	}
	return &mr
}

func main() {
	var raw bson.Raw
	raw, err := bson.Marshal(bson.D{
		{"x", 1},
		{"y", 2},
		{"z", "foo"},
		{"w", bson.D{
			{"x", 123},
		}},
	})
	if err != nil {
		log.Fatalf("error in Marshal: %v", err)
	}

	// NewFromIOReader - Success.
	{
		mr := newMyReader()
		raw, err := bson.NewFromIOReader(mr)
		if err != nil {
			log.Fatalf("error in NewFromIOReader: %v", err)
		}
		_, err = raw.LookupErr("x")
		if err != nil {
			log.Fatal("expected to find 'x', but could not")
		}
	}

	// Validate - Success.
	{
		err = raw.Validate()
		if err != nil {
			log.Fatalf("error in Validate: %v", err)
		}
	}

	// Validate - Failure.
	{
		first_byte := raw[0]
		raw[0] = 0xFF
		err = raw.Validate()
		if err == nil {
			log.Fatalf("expected error in Validate, got none")
		}
		raw[0] = first_byte
	}

	// Lookup – Success.
	{
		rawValue := raw.Lookup("x")
		if rawValue.Type != bson.TypeInt32 {
			log.Fatalf("Expected 'x' to be TypeInt32, got %v", rawValue.Type)
		}
		// How do I get the RawValue.Value represented in a Go type?
		_, ok := rawValue.Int32OK()
		if !ok {
			log.Fatalf("Error in Int32OK: %v", err)
		}
	}

	// Lookup - Failure.
	{
		rawValue := raw.Lookup("notfound")
		if rawValue.Type != 0 {
			log.Fatalf("Expected empty rawValue, got: %v", rawValue)
		}
	}

	// Lookup - Recursing.
	{
		rawValue := raw.Lookup("w", "x")
		if rawValue.Type == 0 {
			log.Fatalf("Expected Lookup('w', 'x') to find value, but got empty")
		}
		_, ok := rawValue.Int32OK()
		if !ok {
			log.Fatalf("Expected 'w.x' to be Int32, but not OK")
		}
	}

	// LookupErr - Success.
	{
		rawValue, err := raw.LookupErr("y")
		if err != nil {
			log.Fatalf("Error in LookupErr for 'y': %v", err)
		}
		_, ok := rawValue.AsInt32OK()
		if !ok {
			log.Fatalf("AsInt32OK for 'y' not OK")
		}
	}

	// LookupErr - Failure.
	{
		_, err := raw.LookupErr("notfound")
		// Errors are returned from 'x' package.
		if err != bsoncore.ErrElementNotFound {
			log.Fatalf("expected LookupErr to return an error, but got no error")
		}
	}

	// Elements – Success.
	{
		elements, err := raw.Elements()
		if err != nil {
			log.Fatalf("error in Elements: %v", err)
		}

		fmt.Println("begin iterating elements ... ")
		for _, element := range elements {
			/* Print Int32 elements. */
			if element.Value().Type == bson.TypeInt32 {
				fmt.Printf("... key=%v, int32 value=%v\n", element.Key(), element.Value().Int32())
			} else {
				fmt.Printf("... key=%v, %v\n", element.Key(), element.Value().Type)
			}
		}
		fmt.Println("end iterating elements ... ")
	}

	// Values - Success.
	{
		fmt.Println("begin iterating values ...")
		values, err := raw.Values()
		if err != nil {
			log.Fatalf("error in Values: %v", err)
		}
		for _, value := range values {
			fmt.Printf("... value=%v\n", value)
		}
		fmt.Println("end iterating values ...")
	}

	// Index - Success.
	{
		rawElement := raw.Index(0)
		if rawElement.Key() != "x" {
			log.Fatalf("expected first element to be 'x', got: %v", rawElement.Key())
		}
	}

	// IndexErr – Failure;
	{
		_, err := raw.IndexErr(100)
		if err == nil {
			log.Fatalf("expected error from IndexErr, got none")
		}
	}

	// String
	{
		fmt.Printf("String() of document: %v", raw.String())
	}
}
