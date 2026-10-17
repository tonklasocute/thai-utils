// Command bahttext spells out amounts in Thai words.
//
// Usage:
//
//	go run ./examples/bahttext [amount ...]
package main

import (
	"fmt"
	"os"

	"github.com/tonklasocute/thai-utils/bahttext"
)

func main() {
	amounts := os.Args[1:]
	if len(amounts) == 0 {
		amounts = []string{"1", "21", "101", "1250.50", "0.25", "1,000,001", "-42.10", "1000000000000"}
	}

	for _, a := range amounts {
		words, err := bahttext.FormatString(a)
		if err != nil {
			fmt.Printf("%16s  error: %v\n", a, err)
			continue
		}
		fmt.Printf("%16s  %s\n", a, words)
	}
}
