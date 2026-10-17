package thaidate_test

import (
	"fmt"
	"time"

	"github.com/tonklasocute/thai-utils/thaidate"
)

func ExampleToBE() {
	fmt.Println(thaidate.ToBE(2026))
	// Output: 2569
}

func ExampleToCE() {
	fmt.Println(thaidate.ToCE(2569))
	// Output: 2026
}

func ExampleFormat() {
	t := time.Date(2026, time.October, 17, 14, 30, 0, 0, time.UTC)

	fmt.Println(thaidate.Format(t, thaidate.DateLong))
	fmt.Println(thaidate.Format(t, thaidate.DateShort))
	fmt.Println(thaidate.Format(t, thaidate.DateFull))
	fmt.Println(thaidate.Format(t, "Mon 2 Jan 2006 15:04 น."))
	// Output:
	// 17 ตุลาคม 2569
	// 17 ต.ค. 69
	// วันเสาร์ที่ 17 ตุลาคม พ.ศ. 2569
	// ส. 17 ต.ค. 2569 14:30 น.
}

func ExampleParse() {
	t, err := thaidate.Parse("17/10/2569")
	if err != nil {
		panic(err)
	}
	fmt.Println(t.Format(time.DateOnly))
	// Output: 2026-10-17
}

func ExampleMonthName() {
	fmt.Println(thaidate.MonthName(time.April), thaidate.ShortMonthName(time.April))
	// Output: เมษายน เม.ย.
}

func ExampleWeekdayName() {
	fmt.Println(thaidate.WeekdayName(time.Friday), thaidate.ShortWeekdayName(time.Friday))
	// Output: ศุกร์ ศ.
}

func ExampleToThaiDigits() {
	fmt.Println(thaidate.ToThaiDigits("17/10/2569"))
	// Output: ๑๗/๑๐/๒๕๖๙
}

func ExampleToArabicDigits() {
	fmt.Println(thaidate.ToArabicDigits("๑๗/๑๐/๒๕๖๙"))
	// Output: 17/10/2569
}
