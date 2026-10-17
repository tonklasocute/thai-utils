// Command thaidate demonstrates the thaidate package.
package main

import (
	"fmt"
	"time"

	"github.com/tonklasocute/thai-utils/thaidate"
)

func main() {
	now := time.Now()

	fmt.Println("Long:   ", thaidate.Format(now, thaidate.DateLong))
	fmt.Println("Short:  ", thaidate.Format(now, thaidate.DateShort))
	fmt.Println("Full:   ", thaidate.Format(now, thaidate.DateFull))
	fmt.Println("Time:   ", thaidate.Format(now, thaidate.DateTimeLong))
	fmt.Println("Digits: ", thaidate.ToThaiDigits(thaidate.Format(now, thaidate.DateLong)))

	t, err := thaidate.Parse("13/04/2569")
	if err != nil {
		panic(err)
	}
	fmt.Printf("Parsed 13/04/2569 -> %s (Songkran, CE %d)\n", t.Format(time.DateOnly), t.Year())
}
