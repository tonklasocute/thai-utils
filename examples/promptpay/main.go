// Command promptpay prints PromptPay QR payloads.
//
// Usage:
//
//	go run ./examples/promptpay [target] [amount]
//
// Paste the output into any QR code generator, or render it in Go with a QR
// library such as github.com/skip2/go-qrcode, then scan it with a Thai
// banking app.
package main

import (
	"fmt"
	"os"
	"strconv"

	"github.com/tonklasocute/thai-utils/promptpay"
)

func main() {
	target := "0812345678"
	var amount *float64

	if len(os.Args) > 1 {
		target = os.Args[1]
	}
	if len(os.Args) > 2 {
		v, err := strconv.ParseFloat(os.Args[2], 64)
		if err != nil {
			fmt.Fprintln(os.Stderr, "invalid amount:", err)
			os.Exit(2)
		}
		amount = &v
	}

	payload, err := promptpay.Payload(target, amount)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	fmt.Println(payload)
}
