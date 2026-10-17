// Package thaidate converts between Buddhist Era (พุทธศักราช, พ.ศ.) and
// Common Era years, formats dates with Thai month and weekday names, parses
// Thai-style DD/MM/YYYY dates, and converts between Arabic and Thai digits.
//
// [Format] accepts the same reference-time layouts as [time.Time.Format].
// The tokens below are rendered the Thai way; every other token behaves
// exactly as in the standard library:
//
//	2006     Buddhist Era year, e.g. 2569
//	06       last two digits of the Buddhist Era year, e.g. 69
//	January  full Thai month name, e.g. ตุลาคม
//	Jan      abbreviated Thai month name, e.g. ต.ค.
//	Monday   full Thai weekday name, e.g. เสาร์
//	Mon      abbreviated Thai weekday name, e.g. ส.
//
// Thai text in a layout, such as "วันที่" or "พ.ศ.", is copied verbatim.
package thaidate

import (
	"errors"
	"strconv"
	"strings"
	"time"
)

// EraOffset is the difference between a Buddhist Era year and the
// corresponding Common Era (Gregorian) year.
const EraOffset = 543

// Common layouts for [Format]. Sample output is for Saturday 17 October 2026.
const (
	DateLong     = "2 January 2006"                   // 17 ตุลาคม 2569
	DateShort    = "2 Jan 06"                         // 17 ต.ค. 69
	DateFull     = "วันMondayที่ 2 January พ.ศ. 2006" // วันเสาร์ที่ 17 ตุลาคม พ.ศ. 2569
	DateNumeric  = "02/01/2006"                       // 17/10/2569
	DateTimeLong = "2 January 2006 เวลา 15:04 น."     // 17 ตุลาคม 2569 เวลา 14:30 น.
)

// Sentinel errors returned by [Parse] and [ParseInLocation]. Compare them
// with [errors.Is].
var (
	// ErrInvalidFormat means the input is not in DD/MM/YYYY form.
	ErrInvalidFormat = errors.New("thaidate: date must be in DD/MM/YYYY form")
	// ErrInvalidDate means the input is well-formed but names a day that
	// does not exist, such as 30/02/2567.
	ErrInvalidDate = errors.New("thaidate: no such date")
)

var (
	monthNames = [...]string{
		"มกราคม", "กุมภาพันธ์", "มีนาคม", "เมษายน", "พฤษภาคม", "มิถุนายน",
		"กรกฎาคม", "สิงหาคม", "กันยายน", "ตุลาคม", "พฤศจิกายน", "ธันวาคม",
	}
	shortMonthNames = [...]string{
		"ม.ค.", "ก.พ.", "มี.ค.", "เม.ย.", "พ.ค.", "มิ.ย.",
		"ก.ค.", "ส.ค.", "ก.ย.", "ต.ค.", "พ.ย.", "ธ.ค.",
	}
	weekdayNames = [...]string{
		"อาทิตย์", "จันทร์", "อังคาร", "พุธ", "พฤหัสบดี", "ศุกร์", "เสาร์",
	}
	shortWeekdayNames = [...]string{
		"อา.", "จ.", "อ.", "พ.", "พฤ.", "ศ.", "ส.",
	}
)

// ToBE converts a Common Era year to a Buddhist Era year: ToBE(2026) == 2569.
func ToBE(year int) int {
	return year + EraOffset
}

// ToCE converts a Buddhist Era year to a Common Era year: ToCE(2569) == 2026.
func ToCE(year int) int {
	return year - EraOffset
}

// MonthName returns the full Thai name of m, e.g. "มกราคม" for January. It
// returns an empty string if m is out of range.
func MonthName(m time.Month) string {
	if m < time.January || m > time.December {
		return ""
	}
	return monthNames[m-1]
}

// ShortMonthName returns the abbreviated Thai name of m, e.g. "ม.ค." for
// January. It returns an empty string if m is out of range.
func ShortMonthName(m time.Month) string {
	if m < time.January || m > time.December {
		return ""
	}
	return shortMonthNames[m-1]
}

// WeekdayName returns the full Thai name of d without the "วัน" prefix, e.g.
// "จันทร์" for Monday. It returns an empty string if d is out of range.
func WeekdayName(d time.Weekday) string {
	if d < time.Sunday || d > time.Saturday {
		return ""
	}
	return weekdayNames[d]
}

// ShortWeekdayName returns the abbreviated Thai name of d, e.g. "จ." for
// Monday. It returns an empty string if d is out of range.
func ShortWeekdayName(d time.Weekday) string {
	if d < time.Sunday || d > time.Saturday {
		return ""
	}
	return shortWeekdayNames[d]
}

// Format returns t formatted according to layout, using Buddhist Era years
// and Thai month and weekday names. See the package documentation for the
// tokens that are localized; all other layout tokens behave exactly as in
// [time.Time.Format].
func Format(t time.Time, layout string) string {
	var b strings.Builder
	b.Grow(len(layout) + 16)
	start := 0 // start of the pending run of standard-library tokens
	for i := 0; i < len(layout); {
		n, thai := nextToken(layout[i:])
		if !thai {
			i += n
			continue
		}
		if start < i {
			b.WriteString(t.Format(layout[start:i]))
		}
		b.WriteString(renderThai(t, layout[i:i+n]))
		i += n
		start = i
	}
	if start < len(layout) {
		b.WriteString(t.Format(layout[start:]))
	}
	return b.String()
}

// nextToken reports the byte length of the layout token at the start of s and
// whether Format renders it itself. It mirrors the tokenizer of the time
// package so that splitting a layout around Thai tokens never changes how the
// remaining parts are interpreted.
func nextToken(s string) (n int, thai bool) {
	has := func(p string) bool { return strings.HasPrefix(s, p) }
	lowerAt := func(i int) bool { return len(s) > i && s[i] >= 'a' && s[i] <= 'z' }

	switch s[0] {
	case 'J':
		if has("January") {
			return 7, true
		}
		if has("Jan") && !lowerAt(3) {
			return 3, true
		}
	case 'M':
		if has("Monday") {
			return 6, true
		}
		if has("Mon") && !lowerAt(3) {
			return 3, true
		}
		if has("MST") {
			return 3, false
		}
	case '0':
		if len(s) >= 2 && s[1] >= '1' && s[1] <= '6' {
			return 2, s[1] == '6'
		}
		if has("002") {
			return 3, false
		}
	case '1':
		if has("15") {
			return 2, false
		}
	case '2':
		if has("2006") {
			return 4, true
		}
	case '_':
		// "_2006" is a literal underscore followed by the year.
		if has("_2") && !has("_2006") {
			return 2, false
		}
		if has("__2") {
			return 3, false
		}
	case 'P', 'p':
		if has("PM") || has("pm") {
			return 2, false
		}
	case '-', 'Z':
		for _, z := range [...]string{"07:00:00", "070000", "07:00", "0700", "07"} {
			if strings.HasPrefix(s[1:], z) {
				return 1 + len(z), false
			}
		}
	case '.', ',':
		if len(s) >= 2 && (s[1] == '0' || s[1] == '9') {
			j := 1
			for j < len(s) && s[j] == s[1] {
				j++
			}
			if j == len(s) || s[j] < '0' || s[j] > '9' {
				return j, false
			}
		}
	}
	return 1, false
}

func renderThai(t time.Time, token string) string {
	switch token {
	case "January":
		return monthNames[t.Month()-1]
	case "Jan":
		return shortMonthNames[t.Month()-1]
	case "Monday":
		return weekdayNames[t.Weekday()]
	case "Mon":
		return shortWeekdayNames[t.Weekday()]
	case "2006":
		return pad(ToBE(t.Year()), 4)
	default: // "06"
		yy := ToBE(t.Year()) % 100
		if yy < 0 {
			yy = -yy
		}
		return pad(yy, 2)
	}
}

func pad(v, width int) string {
	s := strconv.Itoa(v)
	if v >= 0 && len(s) < width {
		s = strings.Repeat("0", width-len(s)) + s
	}
	return s
}

// Parse parses a date written as DD/MM/YYYY with a Buddhist Era year, such
// as "17/10/2569", and returns midnight UTC of that day. Day and month may
// have one or two digits; the year must have four. Thai digits (๐–๙) are
// accepted as well as Arabic digits.
func Parse(s string) (time.Time, error) {
	return ParseInLocation(s, time.UTC)
}

// ParseInLocation is like [Parse] but returns midnight in loc. loc must not
// be nil.
func ParseInLocation(s string, loc *time.Location) (time.Time, error) {
	parts := strings.Split(strings.TrimSpace(ToArabicDigits(s)), "/")
	if len(parts) != 3 {
		return time.Time{}, ErrInvalidFormat
	}
	day, ok1 := atoi(parts[0], 1, 2)
	month, ok2 := atoi(parts[1], 1, 2)
	year, ok3 := atoi(parts[2], 4, 4)
	if !ok1 || !ok2 || !ok3 {
		return time.Time{}, ErrInvalidFormat
	}

	t := time.Date(ToCE(year), time.Month(month), day, 0, 0, 0, 0, loc)
	if t.Day() != day || int(t.Month()) != month {
		return time.Time{}, ErrInvalidDate
	}
	return t, nil
}

// atoi parses s as an unsigned decimal of minLen to maxLen ASCII digits.
func atoi(s string, minLen, maxLen int) (int, bool) {
	if len(s) < minLen || len(s) > maxLen {
		return 0, false
	}
	n := 0
	for i := 0; i < len(s); i++ {
		if s[i] < '0' || s[i] > '9' {
			return 0, false
		}
		n = n*10 + int(s[i]-'0')
	}
	return n, true
}

// ToThaiDigits replaces the Arabic digits 0–9 in s with the Thai digits ๐–๙.
// All other characters are left unchanged.
func ToThaiDigits(s string) string {
	return strings.Map(func(r rune) rune {
		if r >= '0' && r <= '9' {
			return r - '0' + '๐'
		}
		return r
	}, s)
}

// ToArabicDigits replaces the Thai digits ๐–๙ in s with the Arabic digits
// 0–9. All other characters are left unchanged.
func ToArabicDigits(s string) string {
	return strings.Map(func(r rune) rune {
		if r >= '๐' && r <= '๙' {
			return r - '๐' + '0'
		}
		return r
	}, s)
}
