// Package idcard validates, formats and generates Thai national ID numbers
// (เลขประจำตัวประชาชน).
//
// A Thai national ID has 13 digits. The last digit is a check digit computed
// from the first 12:
//
//	sum   = Σ digit[i] × (13 − i)   for i = 0..11
//	check = (11 − sum mod 11) mod 10
//
// The first digit is the person-type category (1–8). Personal IDs never start
// with 0; IDs starting with 0 are juristic-person tax IDs, which are handled by
// the taxid package.
//
// Functions accept IDs written with or without the usual separators, so
// "1103702071811", "1-1037-02071-81-1" and "1 1037 02071 81 1" are all
// equivalent.
package idcard

import (
	"errors"
	"math/rand/v2"

	"github.com/tonklasocute/thai-utils/internal/thaiid"
)

// Sentinel errors returned by [Check] and [Format]. Compare them with
// [errors.Is].
var (
	// ErrInvalidCharacter means the input contains something other than
	// digits, dashes and whitespace.
	ErrInvalidCharacter = errors.New("idcard: invalid character")
	// ErrInvalidLength means the input does not contain exactly 13 digits.
	ErrInvalidLength = errors.New("idcard: must contain exactly 13 digits")
	// ErrLeadingZero means the ID starts with 0, which is reserved for
	// juristic persons and never used for personal IDs.
	ErrLeadingZero = errors.New("idcard: personal ID must not start with 0")
	// ErrInvalidChecksum means the check digit does not match the first 12
	// digits.
	ErrInvalidChecksum = errors.New("idcard: invalid check digit")
)

// Check validates id and returns nil if it is a valid personal Thai national
// ID, or one of the package's sentinel errors describing the first problem
// found.
func Check(id string) error {
	_, err := parse(id)
	return err
}

// Validate reports whether id is a valid personal Thai national ID. Dashes and
// whitespace are ignored; any other non-digit character makes the ID invalid.
func Validate(id string) bool {
	return Check(id) == nil
}

// Format validates id and returns it in the standard X-XXXX-XXXXX-XX-X layout.
// It returns a sentinel error if id is not valid.
func Format(id string) (string, error) {
	d, err := parse(id)
	if err != nil {
		return "", err
	}
	return thaiid.Format(d), nil
}

// Normalize returns id with every non-digit character removed. It does not
// validate the result; use [Validate] or [Check] for that.
func Normalize(id string) string {
	return thaiid.Digits(id)
}

// Generate returns a random, checksum-valid personal ID made of 13 digits
// without separators. The first digit is between 1 and 8.
//
// Generate is meant for test data and fixtures only. It uses a
// non-cryptographic random source, and because the ID space is small a
// generated number may belong to a real person. Never present generated IDs
// as genuine.
func Generate() string {
	var b [thaiid.Length]byte
	b[0] = byte('1' + rand.IntN(8)) //nolint:gosec // test data, not security sensitive
	for i := 1; i < thaiid.Length-1; i++ {
		b[i] = byte('0' + rand.IntN(10)) //nolint:gosec // test data, not security sensitive
	}
	b[thaiid.Length-1] = thaiid.CheckDigit(string(b[:thaiid.Length-1]))
	return string(b[:])
}

func parse(id string) (string, error) {
	d, ok := thaiid.Clean(id)
	switch {
	case !ok:
		return "", ErrInvalidCharacter
	case len(d) != thaiid.Length:
		return "", ErrInvalidLength
	case d[0] == '0':
		return "", ErrLeadingZero
	case !thaiid.Valid(d):
		return "", ErrInvalidChecksum
	}
	return d, nil
}
