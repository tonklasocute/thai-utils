// Package thaiutils is a collection of small, dependency-free utilities for
// working with Thailand-specific data in Go.
//
// The root package contains no code; import the sub-package you need:
//
//   - idcard    validates, formats and generates 13-digit Thai national ID numbers.
//   - taxid     validates and formats Thai tax IDs, including juristic-person IDs.
//   - phone     validates, normalizes and formats Thai mobile and landline numbers.
//   - promptpay builds EMVCo-compliant PromptPay QR payloads (Thai QR Payment).
//   - thaidate  converts between Buddhist Era and Common Era years, formats dates
//     with Thai month and weekday names, and converts Thai digits.
//   - bahttext  spells out amounts of money in Thai words, like Excel's BAHTTEXT.
//
// Every sub-package uses the standard library only, so adding thai-utils to a
// project never pulls in transitive dependencies.
//
// Example:
//
//	import (
//		"github.com/tonklasocute/thai-utils/bahttext"
//		"github.com/tonklasocute/thai-utils/idcard"
//	)
//
//	idcard.Validate("1-1037-02071-81-1") // true
//	bahttext.Format(1250.50)            // "หนึ่งพันสองร้อยห้าสิบบาทห้าสิบสตางค์"
package thaiutils
