package taxid_test

import (
	"errors"
	"fmt"

	"github.com/tonklasocute/thai-utils/taxid"
)

func ExampleValidate() {
	fmt.Println(taxid.Validate("0105556000017")) // company
	fmt.Println(taxid.Validate("1103702071811")) // individual
	fmt.Println(taxid.Validate("0105556000018"))
	// Output:
	// true
	// true
	// false
}

func ExampleCheck() {
	if err := taxid.Check("0105556"); errors.Is(err, taxid.ErrInvalidLength) {
		fmt.Println("tax ID must have 13 digits")
	}
	// Output: tax ID must have 13 digits
}

func ExampleFormat() {
	s, err := taxid.Format("0105556000017")
	if err != nil {
		panic(err)
	}
	fmt.Println(s)
	// Output: 0-1055-56000-01-7
}

func ExampleNormalize() {
	fmt.Println(taxid.Normalize("0-1055-56000-01-7"))
	// Output: 0105556000017
}

func ExampleIsJuristic() {
	fmt.Println(taxid.IsJuristic("0105556000017"))
	fmt.Println(taxid.IsJuristic("1103702071811"))
	// Output:
	// true
	// false
}
