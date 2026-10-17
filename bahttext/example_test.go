package bahttext_test

import (
	"fmt"

	"github.com/tonklasocute/thai-utils/bahttext"
)

func ExampleFormat() {
	fmt.Println(bahttext.Format(1))
	fmt.Println(bahttext.Format(21))
	fmt.Println(bahttext.Format(1250.50))
	fmt.Println(bahttext.Format(0.25))
	fmt.Println(bahttext.Format(-100))
	// Output:
	// หนึ่งบาทถ้วน
	// ยี่สิบเอ็ดบาทถ้วน
	// หนึ่งพันสองร้อยห้าสิบบาทห้าสิบสตางค์
	// ยี่สิบห้าสตางค์
	// ลบหนึ่งร้อยบาทถ้วน
}

func ExampleFormatString() {
	s, err := bahttext.FormatString("1,000,001.00")
	if err != nil {
		panic(err)
	}
	fmt.Println(s)

	// No float64 limits: amounts of any size are exact.
	s, _ = bahttext.FormatString("1000000000000")
	fmt.Println(s)
	// Output:
	// หนึ่งล้านเอ็ดบาทถ้วน
	// หนึ่งล้านล้านบาทถ้วน
}
