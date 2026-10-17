package idcard_test

import (
	"errors"
	"fmt"

	"github.com/tonklasocute/thai-utils/idcard"
)

func ExampleValidate() {
	fmt.Println(idcard.Validate("1-1037-02071-81-1"))
	fmt.Println(idcard.Validate("1103702071812"))
	// Output:
	// true
	// false
}

func ExampleCheck() {
	err := idcard.Check("1103702071812")
	if errors.Is(err, idcard.ErrInvalidChecksum) {
		fmt.Println("typo in ID number")
	}
	// Output: typo in ID number
}

func ExampleFormat() {
	s, err := idcard.Format("1103702071811")
	if err != nil {
		panic(err)
	}
	fmt.Println(s)
	// Output: 1-1037-02071-81-1
}

func ExampleNormalize() {
	fmt.Println(idcard.Normalize("1-1037-02071-81-1"))
	// Output: 1103702071811
}

func ExampleGenerate() {
	id := idcard.Generate() // random, e.g. "3849105732112"
	fmt.Println(len(id), idcard.Validate(id))
	// Output: 13 true
}
