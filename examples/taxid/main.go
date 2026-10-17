// Command taxid demonstrates the taxid package.
package main

import (
	"fmt"

	"github.com/tonklasocute/thai-utils/taxid"
)

func main() {
	inputs := []string{
		"0105556000017",
		"1-1037-02071-81-1",
		"0105556000018",
	}

	for _, in := range inputs {
		formatted, err := taxid.Format(in)
		if err != nil {
			fmt.Printf("%-20s invalid: %v\n", in, err)
			continue
		}
		kind := "individual"
		if taxid.IsJuristic(in) {
			kind = "juristic person"
		}
		fmt.Printf("%-20s %s (%s)\n", in, formatted, kind)
	}
}
