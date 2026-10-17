// Package phone validates, normalizes and formats Thai telephone numbers.
//
// Two kinds of numbers are recognized, written in their local (trunk-prefix)
// form:
//
//   - Mobile: 10 digits starting with 06, 08 or 09, e.g. 081-234-5678.
//   - Landline: 9 digits starting with 02 (Bangkok metropolitan area) or with
//     03, 04, 05 or 07 (provincial area codes), e.g. 02-123-4567 or
//     053-123-456.
//
// Input may use the international prefix "+66" or "66" in place of the
// leading 0, and may contain spaces, dashes, dots and parentheses, so
// "+66 81 234 5678", "66812345678", "(02) 123-4567" and "+66 (0)2 123 4567"
// are all accepted. Short service codes (1xxx) and toll-free 1-800 numbers are
// not supported.
package phone

import (
	"errors"
	"strings"
	"unicode"
)

// Sentinel errors returned by [Normalize], [ToE164] and [Format]. Compare them
// with [errors.Is].
var (
	// ErrInvalidCharacter means the input contains a character other than
	// digits, a leading '+', whitespace, dashes, dots or parentheses.
	ErrInvalidCharacter = errors.New("phone: invalid character")
	// ErrInvalidCountryCode means the number has an international prefix
	// other than +66.
	ErrInvalidCountryCode = errors.New("phone: country code must be +66")
	// ErrInvalidPrefix means the number does not start with a known Thai
	// mobile or landline prefix.
	ErrInvalidPrefix = errors.New("phone: unknown number prefix")
	// ErrInvalidLength means the number has the wrong number of digits for
	// its prefix.
	ErrInvalidLength = errors.New("phone: invalid number length")
)

const countryCode = "66"

// Validate reports whether s is a valid Thai mobile or landline number.
func Validate(s string) bool {
	_, _, err := parse(s)
	return err == nil
}

// IsMobile reports whether s is a valid Thai mobile number.
func IsMobile(s string) bool {
	_, mobile, err := parse(s)
	return err == nil && mobile
}

// IsLandline reports whether s is a valid Thai landline number.
func IsLandline(s string) bool {
	_, mobile, err := parse(s)
	return err == nil && !mobile
}

// Normalize validates s and returns it in local form with digits only, e.g.
// "0812345678" or "021234567".
func Normalize(s string) (string, error) {
	local, _, err := parse(s)
	return local, err
}

// ToE164 validates s and returns it in E.164 international format, e.g.
// "+66812345678".
func ToE164(s string) (string, error) {
	local, _, err := parse(s)
	if err != nil {
		return "", err
	}
	return "+" + countryCode + local[1:], nil
}

// Format validates s and returns it in the conventional Thai display format:
//
//   - mobile:              081-234-5678
//   - Bangkok landline:    02-123-4567
//   - provincial landline: 053-123-456
func Format(s string) (string, error) {
	local, mobile, err := parse(s)
	if err != nil {
		return "", err
	}
	switch {
	case mobile:
		return local[:3] + "-" + local[3:6] + "-" + local[6:], nil
	case local[1] == '2':
		return local[:2] + "-" + local[2:5] + "-" + local[5:], nil
	default:
		return local[:3] + "-" + local[3:6] + "-" + local[6:], nil
	}
}

// parse returns the local digits-only form of s and whether it is a mobile
// number.
func parse(s string) (local string, mobile bool, err error) {
	s = strings.TrimSpace(s)
	plus := strings.HasPrefix(s, "+")
	if plus {
		s = s[1:]
	}

	var b strings.Builder
	b.Grow(len(s))
	for _, r := range s {
		switch {
		case r >= '0' && r <= '9':
			b.WriteRune(r)
		case r == '-' || r == '.' || r == '(' || r == ')' || unicode.IsSpace(r):
			// Separator: skip.
		default:
			return "", false, ErrInvalidCharacter
		}
	}
	d := b.String()

	switch {
	case plus && !strings.HasPrefix(d, countryCode):
		return "", false, ErrInvalidCountryCode
	case strings.HasPrefix(d, countryCode):
		// Drop the country code. Some people keep the trunk 0 after it,
		// as in "+66 (0)81 ...", so drop that too.
		d = "0" + strings.TrimPrefix(d[len(countryCode):], "0")
	case d == "":
		return "", false, ErrInvalidLength
	case d[0] != '0':
		return "", false, ErrInvalidPrefix
	}

	if len(d) < 2 {
		return "", false, ErrInvalidLength
	}
	switch d[1] {
	case '6', '8', '9':
		if len(d) != 10 {
			return "", false, ErrInvalidLength
		}
		return d, true, nil
	case '2', '3', '4', '5', '7':
		if len(d) != 9 {
			return "", false, ErrInvalidLength
		}
		return d, false, nil
	default:
		return "", false, ErrInvalidPrefix
	}
}
