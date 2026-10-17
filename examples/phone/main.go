// Command phone demonstrates the phone package.
package main

import (
	"fmt"

	"github.com/tonklasocute/thai-utils/phone"
)

func main() {
	inputs := []string{
		"0812345678",
		"+66 81 234 5678",
		"(02) 123-4567",
		"053123456",
		"+1 415 555 0100",
		"1669",
	}

	fmt.Printf("%-18s %-14s %-14s %s\n", "input", "display", "E.164", "type")
	for _, in := range inputs {
		display, err := phone.Format(in)
		if err != nil {
			fmt.Printf("%-18s error: %v\n", in, err)
			continue
		}
		e164, _ := phone.ToE164(in)
		kind := "landline"
		if phone.IsMobile(in) {
			kind = "mobile"
		}
		fmt.Printf("%-18s %-14s %-14s %s\n", in, display, e164, kind)
	}
}
