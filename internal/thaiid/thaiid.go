// Package thaiid implements the 13-digit identifier scheme shared by Thai
// national ID cards and Thai tax IDs.
package thaiid

import (
	"strings"
	"unicode"
)

// Length is the number of digits in a Thai national ID or tax ID.
const Length = 13

// Clean removes the separators people commonly use when writing an ID
// (dashes and whitespace) and returns the remaining ASCII digits. ok is false
// if s contains any other character.
func Clean(s string) (digits string, ok bool) {
	var b strings.Builder
	b.Grow(Length)
	for _, r := range s {
		switch {
		case r >= '0' && r <= '9':
			b.WriteRune(r)
		case r == '-' || unicode.IsSpace(r):
			// Separator: skip.
		default:
			return "", false
		}
	}
	return b.String(), true
}

// Digits returns only the ASCII digits of s, discarding everything else.
func Digits(s string) string {
	var b strings.Builder
	b.Grow(Length)
	for _, r := range s {
		if r >= '0' && r <= '9' {
			b.WriteRune(r)
		}
	}
	return b.String()
}

// CheckDigit returns the ASCII check digit for the first 12 digits of d:
// (11 - Σ d[i]·(13-i) mod 11) mod 10. d must hold at least 12 ASCII digits.
func CheckDigit(d string) byte {
	sum := 0
	for i := 0; i < Length-1; i++ {
		sum += int(d[i]-'0') * (Length - i)
	}
	return byte((11-sum%11)%10) + '0'
}

// Valid reports whether d is exactly 13 ASCII digits with a correct check
// digit.
func Valid(d string) bool {
	if len(d) != Length {
		return false
	}
	for i := 0; i < Length; i++ {
		if d[i] < '0' || d[i] > '9' {
			return false
		}
	}
	return CheckDigit(d) == d[Length-1]
}

// Format lays out a 13-digit string as X-XXXX-XXXXX-XX-X. d must hold exactly
// 13 characters.
func Format(d string) string {
	return d[0:1] + "-" + d[1:5] + "-" + d[5:10] + "-" + d[10:12] + "-" + d[12:13]
}
