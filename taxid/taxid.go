// Package taxid validates and formats Thai taxpayer identification numbers
// (เลขประจำตัวผู้เสียภาษีอากร).
//
// Since 2013 every Thai tax ID has 13 digits and uses the same check-digit
// algorithm as the national ID card:
//
//	sum   = Σ digit[i] × (13 − i)   for i = 0..11
//	check = (11 − sum mod 11) mod 10
//
// For individuals the tax ID is the national ID number. Juristic persons
// (companies, partnerships, foundations, government agencies, ...) receive a
// number that starts with 0, which [IsJuristic] detects.
//
// Functions accept IDs with or without dashes and whitespace.
package taxid

import (
	"errors"

	"github.com/tonklasocute/thai-utils/internal/thaiid"
)

// Sentinel errors returned by [Check] and [Format]. Compare them with
// [errors.Is].
var (
	// ErrInvalidCharacter means the input contains something other than
	// digits, dashes and whitespace.
	ErrInvalidCharacter = errors.New("taxid: invalid character")
	// ErrInvalidLength means the input does not contain exactly 13 digits.
	ErrInvalidLength = errors.New("taxid: must contain exactly 13 digits")
	// ErrInvalidChecksum means the check digit does not match the first 12
	// digits.
	ErrInvalidChecksum = errors.New("taxid: invalid check digit")
)

// Check validates id and returns nil if it is a valid Thai tax ID, or one of
// the package's sentinel errors describing the first problem found.
func Check(id string) error {
	_, err := parse(id)
	return err
}

// Validate reports whether id is a valid Thai tax ID, for either an
// individual or a juristic person.
func Validate(id string) bool {
	return Check(id) == nil
}

// Format validates id and returns it in the X-XXXX-XXXXX-XX-X layout used on
// tax invoices and in Revenue Department systems.
func Format(id string) (string, error) {
	d, err := parse(id)
	if err != nil {
		return "", err
	}
	return thaiid.Format(d), nil
}

// Normalize returns id with every non-digit character removed. It does not
// validate the result.
func Normalize(id string) string {
	return thaiid.Digits(id)
}

// IsJuristic reports whether id is a valid tax ID belonging to a juristic
// person, that is, a valid ID whose first digit is 0. It returns false for
// invalid IDs.
func IsJuristic(id string) bool {
	d, err := parse(id)
	return err == nil && d[0] == '0'
}

func parse(id string) (string, error) {
	d, ok := thaiid.Clean(id)
	switch {
	case !ok:
		return "", ErrInvalidCharacter
	case len(d) != thaiid.Length:
		return "", ErrInvalidLength
	case !thaiid.Valid(d):
		return "", ErrInvalidChecksum
	}
	return d, nil
}
