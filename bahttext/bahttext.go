// Package bahttext spells out amounts of money in Thai words, the way they
// are written on cheques, receipts and tax invoices. It is the Go equivalent
// of Excel's BAHTTEXT function:
//
//	bahttext.Format(1250.50) // "หนึ่งพันสองร้อยห้าสิบบาทห้าสิบสตางค์"
//
// Rules applied:
//
//   - A 1 in the units place is read "เอ็ด" when any higher digit is non-zero
//     (11 → สิบเอ็ด, 101 → หนึ่งร้อยเอ็ด, 1,000,001 → หนึ่งล้านเอ็ด), and
//     "หนึ่ง" when it stands alone.
//   - The tens place reads "สิบ" for 1 (not "หนึ่งสิบ") and "ยี่สิบ" for 2.
//   - "ล้าน" repeats for every further group of six digits, so 10^12 reads
//     หนึ่งล้านล้าน. There is no upper limit when using [FormatString].
//   - Amounts are rounded half-up to whole satang. Whole amounts end in
//     "บาทถ้วน"; amounts with satang end in "สตางค์", and amounts below one
//     baht are read in satang only (0.25 → ยี่สิบห้าสตางค์).
//   - Zero is ศูนย์บาทถ้วน and negative amounts are prefixed with "ลบ".
package bahttext

import (
	"errors"
	"math"
	"strconv"
	"strings"
)

// ErrInvalidNumber is returned by [FormatString] when its input is not a
// decimal number. Compare it with [errors.Is].
var ErrInvalidNumber = errors.New("bahttext: invalid number")

const (
	zeroBaht = "ศูนย์บาทถ้วน"
	minus    = "ลบ"
	baht     = "บาท"
	exact    = "ถ้วน"
	satang   = "สตางค์"
	million  = "ล้าน"
)

var (
	digitWords = [...]string{"", "หนึ่ง", "สอง", "สาม", "สี่", "ห้า", "หก", "เจ็ด", "แปด", "เก้า"}
	placeWords = [...]string{"", "สิบ", "ร้อย", "พัน", "หมื่น", "แสน"}
)

// Format returns amount written in Thai words. It rounds to the nearest
// satang, with halves rounded away from zero, based on the shortest decimal
// representation of amount, so Format(1.005) reads one baht and one satang.
//
// Prefer [FormatString] for values that come from user input or a decimal
// type, since float64 cannot represent every decimal amount exactly and loses
// precision above 2^53. Format returns an empty string for NaN and ±Inf.
func Format(amount float64) string {
	if math.IsNaN(amount) || math.IsInf(amount, 0) {
		return ""
	}
	// The shortest round-trip representation is always a valid decimal.
	s, _ := FormatString(strconv.FormatFloat(amount, 'f', -1, 64))
	return s
}

// FormatString returns the decimal number s written in Thai words. s may have
// a leading sign, thousands separators and any number of decimal places, e.g.
// "1250.5", "-42", "1,000,000.00" or ".75"; surrounding whitespace is
// ignored. Decimal places beyond the second are rounded half-up.
//
// It returns [ErrInvalidNumber] if s is not a valid decimal number.
func FormatString(s string) (string, error) {
	neg, intPart, fracPart, ok := parseDecimal(s)
	if !ok {
		return "", ErrInvalidNumber
	}

	// Round to two decimal places, half-up.
	frac := fracPart + "00"
	st := int(frac[0]-'0')*10 + int(frac[1]-'0')
	if len(fracPart) > 2 && fracPart[2] >= '5' {
		st++
	}
	if st == 100 {
		st = 0
		intPart = increment(intPart)
	}

	intPart = strings.TrimLeft(intPart, "0")
	if intPart == "" && st == 0 {
		return zeroBaht, nil
	}

	var b strings.Builder
	b.Grow(len(intPart)*12 + 64)
	if neg {
		b.WriteString(minus)
	}
	if intPart != "" {
		writeNumber(&b, intPart)
		b.WriteString(baht)
		if st == 0 {
			b.WriteString(exact)
			return b.String(), nil
		}
	}
	writeNumber(&b, strconv.Itoa(st))
	b.WriteString(satang)
	return b.String(), nil
}

// writeNumber writes the Thai reading of d, a string of ASCII digits with no
// leading zeros.
func writeNumber(b *strings.Builder, d string) {
	n := len(d)
	for i := 0; i < n; i++ {
		digit := d[i] - '0'
		place := (n - 1 - i) % 6
		if digit != 0 {
			switch {
			case place == 0 && digit == 1 && i > 0:
				// d has no leading zeros, so i > 0 means a higher
				// non-zero digit exists.
				b.WriteString("เอ็ด")
			case place == 1 && digit == 1:
				// "สิบ", not "หนึ่งสิบ".
			case place == 1 && digit == 2:
				b.WriteString("ยี่")
			default:
				b.WriteString(digitWords[digit])
			}
			b.WriteString(placeWords[place])
		}
		if place == 0 && i < n-1 {
			b.WriteString(million)
		}
	}
}

// parseDecimal splits s into its sign, integer digits and fraction digits.
// Thousands separators are removed from the integer part.
func parseDecimal(s string) (neg bool, intPart, fracPart string, ok bool) {
	s = strings.TrimSpace(s)
	if s != "" && (s[0] == '-' || s[0] == '+') {
		neg = s[0] == '-'
		s = s[1:]
	}
	intPart, fracPart, hasDot := strings.Cut(s, ".")
	if intPart == "" && fracPart == "" {
		return false, "", "", false
	}
	if hasDot && fracPart == "" {
		// Reject "1." to avoid silently accepting truncated input.
		return false, "", "", false
	}
	if !allDigits(fracPart) {
		return false, "", "", false
	}
	if strings.Contains(intPart, ",") {
		groups := strings.Split(intPart, ",")
		for i, g := range groups {
			if (i == 0 && (len(g) < 1 || len(g) > 3)) || (i > 0 && len(g) != 3) {
				return false, "", "", false
			}
		}
		intPart = strings.Join(groups, "")
	}
	if !allDigits(intPart) {
		return false, "", "", false
	}
	return neg, intPart, fracPart, true
}

func allDigits(s string) bool {
	for i := 0; i < len(s); i++ {
		if s[i] < '0' || s[i] > '9' {
			return false
		}
	}
	return true
}

// increment adds one to the decimal digit string d.
func increment(d string) string {
	b := []byte(d)
	for i := len(b) - 1; i >= 0; i-- {
		if b[i] < '9' {
			b[i]++
			return string(b)
		}
		b[i] = '0'
	}
	return "1" + string(b)
}
