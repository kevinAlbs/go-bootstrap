package fun

import (
	"fmt"
)

type person struct {
	Age int
}

func newPerson() *person {
	return &person{123}
}

type personIface interface{}

func main() {
	// A nil pointer is nil.
	var p *person
	fmt.Println("p == nil:", p == nil)

	// A nil interface is nil.
	var piface personIface
	fmt.Println("piface == nil:", piface == nil)

	// But an interface with a set type and nil value is not nil.
	piface = p
	fmt.Println("piface = p; piface == nil:", piface == nil)

	// Consequently, this leads to a trap.
	// Returning a nil pointer to a struct as an error would
	// fail the nil check.
	if err := badFunc(123); err != nil {
		fmt.Println("error: ", err)
	} else {
		fmt.Println("no error")
	}
}

// A real life danger is returning a concrete nil type as an error interface.
type inputError struct{}

func (e *inputError) Error() string {
	return fmt.Sprintf("error! unexpected negative input")
}

func badFunc(x int) error {
	var ret *inputError = nil

	fmt.Println("badFunc called:", x)

	if x < 0 {
		// error on input
		ret = &inputError{}
	}

	return ret
}
