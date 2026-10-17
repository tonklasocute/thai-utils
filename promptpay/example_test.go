package promptpay_test

import (
	"errors"
	"fmt"

	"github.com/tonklasocute/thai-utils/promptpay"
)

func ExamplePayload() {
	// Static QR: the payer types in the amount.
	payload, err := promptpay.Payload("081-234-5678", nil)
	if err != nil {
		panic(err)
	}
	fmt.Println(payload)
	// Output: 00020101021129370016A000000677010111011300668123456785802TH530376463045D82
}

func ExamplePayload_withAmount() {
	// Dynamic QR for exactly 100.50 baht to a national ID.
	payload, err := promptpay.Payload("1-1037-02071-81-1", promptpay.Amount(100.50))
	if err != nil {
		panic(err)
	}
	fmt.Println(payload)
	// Output: 00020101021229370016A000000677010111021311037020718115802TH53037645406100.5063041B0C
}

func ExamplePayload_error() {
	_, err := promptpay.Payload("02-123-4567", nil)
	fmt.Println(errors.Is(err, promptpay.ErrInvalidTarget))
	// Output: true
}

func ExampleVerifyChecksum() {
	scanned := "00020101021129370016A000000677010111011300668123456785802TH530376463045D82"
	fmt.Println(promptpay.VerifyChecksum(scanned))
	// Output: true
}

func ExampleCRC16() {
	fmt.Printf("%04X\n", promptpay.CRC16([]byte("123456789")))
	// Output: 29B1
}
