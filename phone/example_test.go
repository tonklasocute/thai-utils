package phone_test

import (
	"errors"
	"fmt"

	"github.com/tonklasocute/thai-utils/phone"
)

func ExampleValidate() {
	fmt.Println(phone.Validate("081-234-5678"))
	fmt.Println(phone.Validate("+66 2 123 4567"))
	fmt.Println(phone.Validate("1234"))
	// Output:
	// true
	// true
	// false
}

func ExampleNormalize() {
	s, err := phone.Normalize("+66 81 234 5678")
	if err != nil {
		panic(err)
	}
	fmt.Println(s)
	// Output: 0812345678
}

func ExampleNormalize_error() {
	_, err := phone.Normalize("+1 415 555 0100")
	fmt.Println(errors.Is(err, phone.ErrInvalidCountryCode))
	// Output: true
}

func ExampleToE164() {
	s, err := phone.ToE164("081-234-5678")
	if err != nil {
		panic(err)
	}
	fmt.Println(s)
	// Output: +66812345678
}

func ExampleFormat() {
	for _, in := range []string{"0812345678", "021234567", "053123456"} {
		s, err := phone.Format(in)
		if err != nil {
			panic(err)
		}
		fmt.Println(s)
	}
	// Output:
	// 081-234-5678
	// 02-123-4567
	// 053-123-456
}

func ExampleIsMobile() {
	fmt.Println(phone.IsMobile("0812345678"))
	fmt.Println(phone.IsMobile("021234567"))
	// Output:
	// true
	// false
}

func ExampleIsLandline() {
	fmt.Println(phone.IsLandline("02-123-4567"))
	// Output: true
}
