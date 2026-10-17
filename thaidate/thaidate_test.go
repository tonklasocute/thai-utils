package thaidate

import (
	"errors"
	"testing"
	"time"
)

func TestToBEToCE(t *testing.T) {
	tests := []struct{ ce, be int }{
		{2026, 2569},
		{2000, 2543},
		{1941, 2484},
		{0, 543},
		{-543, 0},
	}
	for _, tt := range tests {
		if got := ToBE(tt.ce); got != tt.be {
			t.Errorf("ToBE(%d) = %d, want %d", tt.ce, got, tt.be)
		}
		if got := ToCE(tt.be); got != tt.ce {
			t.Errorf("ToCE(%d) = %d, want %d", tt.be, got, tt.ce)
		}
	}
}

func TestNames(t *testing.T) {
	months := []struct {
		m           time.Month
		full, short string
	}{
		{time.January, "มกราคม", "ม.ค."},
		{time.February, "กุมภาพันธ์", "ก.พ."},
		{time.March, "มีนาคม", "มี.ค."},
		{time.April, "เมษายน", "เม.ย."},
		{time.May, "พฤษภาคม", "พ.ค."},
		{time.June, "มิถุนายน", "มิ.ย."},
		{time.July, "กรกฎาคม", "ก.ค."},
		{time.August, "สิงหาคม", "ส.ค."},
		{time.September, "กันยายน", "ก.ย."},
		{time.October, "ตุลาคม", "ต.ค."},
		{time.November, "พฤศจิกายน", "พ.ย."},
		{time.December, "ธันวาคม", "ธ.ค."},
		{0, "", ""},
		{13, "", ""},
	}
	for _, tt := range months {
		if got := MonthName(tt.m); got != tt.full {
			t.Errorf("MonthName(%d) = %q, want %q", tt.m, got, tt.full)
		}
		if got := ShortMonthName(tt.m); got != tt.short {
			t.Errorf("ShortMonthName(%d) = %q, want %q", tt.m, got, tt.short)
		}
	}

	days := []struct {
		d           time.Weekday
		full, short string
	}{
		{time.Sunday, "อาทิตย์", "อา."},
		{time.Monday, "จันทร์", "จ."},
		{time.Tuesday, "อังคาร", "อ."},
		{time.Wednesday, "พุธ", "พ."},
		{time.Thursday, "พฤหัสบดี", "พฤ."},
		{time.Friday, "ศุกร์", "ศ."},
		{time.Saturday, "เสาร์", "ส."},
		{-1, "", ""},
		{7, "", ""},
	}
	for _, tt := range days {
		if got := WeekdayName(tt.d); got != tt.full {
			t.Errorf("WeekdayName(%d) = %q, want %q", tt.d, got, tt.full)
		}
		if got := ShortWeekdayName(tt.d); got != tt.short {
			t.Errorf("ShortWeekdayName(%d) = %q, want %q", tt.d, got, tt.short)
		}
	}
}

func TestFormat(t *testing.T) {
	sat := time.Date(2026, time.October, 17, 14, 30, 5, 123456789, time.UTC)
	leap := time.Date(2024, time.February, 29, 9, 5, 0, 0, time.UTC)
	y2k := time.Date(2000, time.January, 1, 0, 0, 0, 0, time.UTC)
	ancient := time.Date(-600, time.March, 1, 0, 0, 0, 0, time.UTC)

	tests := []struct {
		name   string
		t      time.Time
		layout string
		want   string
	}{
		{"long", sat, DateLong, "17 ตุลาคม 2569"},
		{"short", sat, DateShort, "17 ต.ค. 69"},
		{"full", sat, DateFull, "วันเสาร์ที่ 17 ตุลาคม พ.ศ. 2569"},
		{"numeric", sat, DateNumeric, "17/10/2569"},
		{"date time", sat, DateTimeLong, "17 ตุลาคม 2569 เวลา 14:30 น."},
		{"leap day", leap, "Mon 2 Jan 2006", "พฤ. 29 ก.พ. 2567"},
		{"y2k", y2k, "Monday 2 January 2006", "เสาร์ 1 มกราคม 2543"},
		{"y2k short year", y2k, "02/01/06", "01/01/43"},
		{"negative BE year", ancient, "2006|06", "-57|57"},
		{"zero-padded short year", time.Date(1966, 5, 1, 0, 0, 0, 0, time.UTC), "06", "09"},
		{"zero-padded long year", time.Date(100, 5, 1, 0, 0, 0, 0, time.UTC), "2006", "0643"},
		{"stdlib tokens", sat, "15:04:05.000 PM MST -07:00", "14:30:05.123 PM UTC +00:00"},
		{"underscore day", y2k, "_2 Jan", " 1 ม.ค."},
		{"underscore year", sat, "_2006", "_2569"},
		{"double underscore", sat, "__2", "290"},
		{"year day", sat, "002", "290"},
		{"comma fraction", sat, "05,00", "05,12"},
		{"dot not fraction", sat, "2006.01", "2569.10"},
		{"lowercase pm", sat, "3pm", "2pm"},
		{"hour 15", sat, "15h", "14h"},
		{"zone Z", sat, "Z07:00", "Z"},
		{"Jan followed by lowercase is literal", sat, "Janet", "Janet"},
		{"Mon followed by lowercase is literal", sat, "Monster", "Monster"},
		{"day then year without separator", sat, "0206", "1769"},
		{"literal text", sat, "ปี 2006", "ปี 2569"},
		{"empty", sat, "", ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Format(tt.t, tt.layout); got != tt.want {
				t.Errorf("Format(%v, %q) = %q, want %q", tt.t, tt.layout, got, tt.want)
			}
		})
	}
}

// TestFormatMatchesStdlib checks that layouts without Thai tokens format
// exactly as the time package does.
func TestFormatMatchesStdlib(t *testing.T) {
	tm := time.Date(2026, time.October, 7, 8, 9, 10, 987654321, time.FixedZone("ICT", 7*3600))
	layouts := []string{
		time.Kitchen, time.TimeOnly, time.DateTime[10:],
		"15:04:05.999999 -0700 Z0700 -07 Z07:00:00 -070000",
		"3:4:5 pm", "1/2 _2 __2 002", ".00,99 .9x", "MST",
	}
	for _, layout := range layouts {
		if got, want := Format(tm, layout), tm.Format(layout); got != want {
			t.Errorf("Format(%q) = %q, want %q (stdlib)", layout, got, want)
		}
	}
}

func TestParse(t *testing.T) {
	tests := []struct {
		in      string
		want    time.Time
		wantErr error
	}{
		{"17/10/2569", time.Date(2026, 10, 17, 0, 0, 0, 0, time.UTC), nil},
		{"1/1/2543", time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC), nil},
		{"29/02/2567", time.Date(2024, 2, 29, 0, 0, 0, 0, time.UTC), nil},
		{" 31/12/2568 ", time.Date(2025, 12, 31, 0, 0, 0, 0, time.UTC), nil},
		{"๑๗/๑๐/๒๕๖๙", time.Date(2026, 10, 17, 0, 0, 0, 0, time.UTC), nil},
		{"29/02/2568", time.Time{}, ErrInvalidDate},
		{"31/04/2569", time.Time{}, ErrInvalidDate},
		{"00/01/2569", time.Time{}, ErrInvalidDate},
		{"01/13/2569", time.Time{}, ErrInvalidDate},
		{"01/00/2569", time.Time{}, ErrInvalidDate},
		{"", time.Time{}, ErrInvalidFormat},
		{"17-10-2569", time.Time{}, ErrInvalidFormat},
		{"17/10/69", time.Time{}, ErrInvalidFormat},
		{"17/10/25690", time.Time{}, ErrInvalidFormat},
		{"117/10/2569", time.Time{}, ErrInvalidFormat},
		{"/10/2569", time.Time{}, ErrInvalidFormat},
		{"1a/10/2569", time.Time{}, ErrInvalidFormat},
		{"17/10/2569/1", time.Time{}, ErrInvalidFormat},
		{"+1/10/2569", time.Time{}, ErrInvalidFormat},
	}
	for _, tt := range tests {
		got, err := Parse(tt.in)
		if !got.Equal(tt.want) || !errors.Is(err, tt.wantErr) {
			t.Errorf("Parse(%q) = (%v, %v), want (%v, %v)", tt.in, got, err, tt.want, tt.wantErr)
		}
	}
}

func TestParseInLocation(t *testing.T) {
	ict := time.FixedZone("ICT", 7*3600)
	got, err := ParseInLocation("17/10/2569", ict)
	if err != nil {
		t.Fatal(err)
	}
	if want := time.Date(2026, 10, 17, 0, 0, 0, 0, ict); !got.Equal(want) || got.Location() != ict {
		t.Errorf("ParseInLocation = %v, want %v", got, want)
	}
}

func TestParseFormatRoundTrip(t *testing.T) {
	start := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	for d := 0; d < 366; d++ {
		tm := start.AddDate(0, 0, d)
		s := Format(tm, DateNumeric)
		got, err := Parse(s)
		if err != nil || !got.Equal(tm) {
			t.Fatalf("Parse(Format(%v)) = (%v, %v)", tm, got, err)
		}
	}
}

func TestDigits(t *testing.T) {
	tests := []struct{ arabic, thai string }{
		{"0123456789", "๐๑๒๓๔๕๖๗๘๙"},
		{"17 ตุลาคม 2569", "๑๗ ตุลาคม ๒๕๖๙"},
		{"฿1,250.50", "฿๑,๒๕๐.๕๐"},
		{"no digits", "no digits"},
		{"", ""},
	}
	for _, tt := range tests {
		if got := ToThaiDigits(tt.arabic); got != tt.thai {
			t.Errorf("ToThaiDigits(%q) = %q, want %q", tt.arabic, got, tt.thai)
		}
		if got := ToArabicDigits(tt.thai); got != tt.arabic {
			t.Errorf("ToArabicDigits(%q) = %q, want %q", tt.thai, got, tt.arabic)
		}
	}
	// Mixed input converts only the relevant digits.
	if got := ToArabicDigits("๑2๓"); got != "123" {
		t.Errorf("ToArabicDigits(mixed) = %q", got)
	}
	if got := ToThaiDigits("๑2๓"); got != "๑๒๓" {
		t.Errorf("ToThaiDigits(mixed) = %q", got)
	}
}

func BenchmarkFormat(b *testing.B) {
	tm := time.Date(2026, 10, 17, 0, 0, 0, 0, time.UTC)
	for i := 0; i < b.N; i++ {
		Format(tm, DateFull)
	}
}
