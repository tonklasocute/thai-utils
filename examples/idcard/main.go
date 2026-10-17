// Command idcard demonstrates the idcard package.
package main

import (
	"errors"
	"fmt"

	"github.com/tonklasocute/thai-utils/idcard"
)

func main() {
	inputs := []string{
		"1-1037-02071-81-1",
		"1103702071812",
		"0105556000017",
		"12345",
	}

	for _, in := range inputs {
		formatted, err := idcard.Format(in)
		switch {
		case err == nil:
			fmt.Printf("%-20s valid   -> %s\n", in, formatted)
		case errors.Is(err, idcard.ErrInvalidChecksum):
			fmt.Printf("%-20s invalid -> check digit mismatch\n", in)
		default:
			fmt.Printf("%-20s invalid -> %v\n", in, err)
		}
	}

	fmt.Println("random test ID:", idcard.Generate())
}
