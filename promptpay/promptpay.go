// Package promptpay generates PromptPay QR code payloads following the
// EMVCo merchant-presented QR specification as profiled by the Bank of
// Thailand (Thai QR Payment).
//
// The package produces the payload string only. To get an image, encode the
// string with any QR code library, for example:
//
//	payload, _ := promptpay.Payload("0812345678", promptpay.Amount(100))
//	png, _ := qrcode.Encode(payload, qrcode.Medium, 256) // github.com/skip2/go-qrcode
//
// Payloads are byte-for-byte compatible with the widely used promptpay-qr
// JavaScript library (https://github.com/dtinth/promptpay-qr).
package promptpay

import (
	"errors"
	"fmt"
	"math"
	"strconv"
	"strings"
	"unicode"

	"github.com/tonklasocute/thai-utils/phone"
	"github.com/tonklasocute/thai-utils/taxid"
)

// Sentinel errors returned by [Payload]. Compare them with [errors.Is].
var (
	// ErrInvalidTarget means the target is not a Thai mobile number, a valid
	// 13-digit national ID / tax ID, or a 15-digit e-wallet ID. When the
	// target looks like an ID but fails validation, the returned error also
	// wraps the underlying taxid error.
	ErrInvalidTarget = errors.New("promptpay: invalid target")
	// ErrInvalidAmount means the amount is not a finite number greater than
	// zero (after rounding to satang), or is too large to encode.
	ErrInvalidAmount = errors.New("promptpay: invalid amount")
)

// EMV tag IDs and fixed values used in a PromptPay payload.
const (
	tagPayloadFormat = "00"
	tagPOIMethod     = "01"
	tagMerchantInfo  = "29"
	tagCurrency      = "53"
	tagAmount        = "54"
	tagCountry       = "58"
	tagCRC           = "63"

	payloadFormatEMV = "01"
	poiStatic        = "11" // reusable QR, payer enters the amount
	poiDynamic       = "12" // QR carries a fixed amount

	subTagAID     = "00"
	subTagPhone   = "01"
	subTagTaxID   = "02"
	subTagEWallet = "03"

	aidPromptPay = "A000000677010111"
	currencyTHB  = "764" // ISO 4217 numeric code for Thai baht
	countryTH    = "TH"

	eWalletLength   = 15
	maxAmountLength = 13 // EMV limit for tag 54
)

// Amount returns a pointer to v, for passing a literal amount to [Payload]:
//
//	promptpay.Payload("0812345678", promptpay.Amount(250))
func Amount(v float64) *float64 {
	return &v
}

// Payload builds the PromptPay QR payload for target.
//
// target may be:
//
//   - a Thai mobile number in any format accepted by the phone package
//     (e.g. "081-234-5678" or "+66812345678");
//   - a 13-digit national ID or tax ID, with or without dashes;
//   - a 15-digit e-wallet ID.
//
// If amount is nil the payload is a static (reusable) QR code and the payer
// enters the amount in their banking app. Otherwise the payload is dynamic
// and carries the amount, rounded to two decimal places.
func Payload(target string, amount *float64) (string, error) {
	subTag, value, err := formatTarget(target)
	if err != nil {
		return "", err
	}

	poi := poiStatic
	var amountStr string
	if amount != nil {
		if amountStr, err = formatAmount(*amount); err != nil {
			return "", err
		}
		poi = poiDynamic
	}

	var b strings.Builder
	b.Grow(96)
	writeTLV(&b, tagPayloadFormat, payloadFormatEMV)
	writeTLV(&b, tagPOIMethod, poi)
	writeTLV(&b, tagMerchantInfo, tlv(subTagAID, aidPromptPay)+tlv(subTag, value))
	writeTLV(&b, tagCountry, countryTH)
	writeTLV(&b, tagCurrency, currencyTHB)
	if amount != nil {
		writeTLV(&b, tagAmount, amountStr)
	}
	// The CRC covers everything up to and including its own tag and length.
	b.WriteString(tagCRC + "04")
	fmt.Fprintf(&b, "%04X", CRC16([]byte(b.String())))
	return b.String(), nil
}

// VerifyChecksum reports whether payload ends with a CRC field (tag 63) whose
// value matches the CRC of the preceding data. It is useful for checking a
// payload read back from a scanned QR code. It does not validate any other
// field.
func VerifyChecksum(payload string) bool {
	const crcField = 8 // "6304" + 4 hex digits
	if len(payload) < crcField || payload[len(payload)-crcField:len(payload)-4] != tagCRC+"04" {
		return false
	}
	want, err := strconv.ParseUint(payload[len(payload)-4:], 16, 16)
	if err != nil {
		return false
	}
	return uint64(CRC16([]byte(payload[:len(payload)-4]))) == want
}

// CRC16 computes the CRC-16/CCITT-FALSE checksum of data (polynomial 0x1021,
// initial value 0xFFFF, no reflection, no final XOR), as required by the
// EMVCo QR specification.
func CRC16(data []byte) uint16 {
	crc := uint16(0xFFFF)
	for _, c := range data {
		crc ^= uint16(c) << 8
		for i := 0; i < 8; i++ {
			if crc&0x8000 != 0 {
				crc = crc<<1 ^ 0x1021
			} else {
				crc <<= 1
			}
		}
	}
	return crc
}

// formatTarget returns the merchant-information sub-tag and value for target.
func formatTarget(target string) (subTag, value string, err error) {
	if local, perr := phone.Normalize(target); perr == nil {
		if !phone.IsMobile(local) {
			return "", "", fmt.Errorf("%w: landline numbers cannot receive PromptPay", ErrInvalidTarget)
		}
		// "0066" + number without trunk 0: "0812345678" -> "0066812345678".
		return subTagPhone, "00" + "66" + local[1:], nil
	}

	digits, ok := stripSeparators(target)
	if !ok {
		return "", "", fmt.Errorf("%w: %q", ErrInvalidTarget, target)
	}
	switch len(digits) {
	case 13:
		if err := taxid.Check(digits); err != nil {
			return "", "", fmt.Errorf("%w: %w", ErrInvalidTarget, err)
		}
		return subTagTaxID, digits, nil
	case eWalletLength:
		return subTagEWallet, digits, nil
	default:
		return "", "", fmt.Errorf("%w: %q is not a mobile number, 13-digit ID or 15-digit e-wallet ID", ErrInvalidTarget, target)
	}
}

func formatAmount(v float64) (string, error) {
	if math.IsNaN(v) || math.IsInf(v, 0) || v <= 0 {
		return "", fmt.Errorf("%w: %v", ErrInvalidAmount, v)
	}
	s := strconv.FormatFloat(v, 'f', 2, 64)
	if s == "0.00" {
		return "", fmt.Errorf("%w: %v rounds to zero", ErrInvalidAmount, v)
	}
	if len(s) > maxAmountLength {
		return "", fmt.Errorf("%w: %v exceeds %d characters", ErrInvalidAmount, v, maxAmountLength)
	}
	return s, nil
}

// stripSeparators removes dashes and whitespace from s. ok is false if s
// contains anything other than digits and those separators.
func stripSeparators(s string) (string, bool) {
	var b strings.Builder
	for _, r := range s {
		switch {
		case r >= '0' && r <= '9':
			b.WriteRune(r)
		case r == '-' || unicode.IsSpace(r):
		default:
			return "", false
		}
	}
	return b.String(), true
}

func writeTLV(b *strings.Builder, tag, value string) {
	b.WriteString(tag)
	fmt.Fprintf(b, "%02d", len(value))
	b.WriteString(value)
}

func tlv(tag, value string) string {
	var b strings.Builder
	writeTLV(&b, tag, value)
	return b.String()
}
