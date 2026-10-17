<div align="center">

# thai-utils

**Thailand-specific utilities for Go: national IDs, tax IDs, phone numbers, PromptPay QR, Buddhist-Era dates and baht text.**

[![Go Reference](https://pkg.go.dev/badge/github.com/tonklasocute/thai-utils.svg)](https://pkg.go.dev/github.com/tonklasocute/thai-utils)
[![CI](https://github.com/tonklasocute/thai-utils/actions/workflows/ci.yml/badge.svg)](https://github.com/tonklasocute/thai-utils/actions/workflows/ci.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/tonklasocute/thai-utils)](https://goreportcard.com/report/github.com/tonklasocute/thai-utils)
[![License: MIT](https://img.shields.io/badge/license-MIT-blue.svg)](LICENSE)
![Go version](https://img.shields.io/github/go-mod/go-version/tonklasocute/thai-utils)
![Dependencies](https://img.shields.io/badge/dependencies-0-brightgreen)

English · [ภาษาไทย](README.th.md)

</div>

---

```go
idcard.Validate("1-1037-02071-81-1")  // true
phone.ToE164("081-234-5678")          // "+66812345678"
promptpay.Payload("0812345678", nil)  // "000201010211…63045D82"
thaidate.Format(t, thaidate.DateFull) // "วันเสาร์ที่ 17 ตุลาคม พ.ศ. 2569"
bahttext.Format(1250.50)              // "หนึ่งพันสองร้อยห้าสิบบาทห้าสิบสตางค์"
```

## Features

| Package | What it does | Key API |
|---|---|---|
| [`idcard`](https://pkg.go.dev/github.com/tonklasocute/thai-utils/idcard) | Thai national ID (เลขประจำตัวประชาชน): checksum validation, formatting, test-data generation | `Validate` `Check` `Format` `Normalize` `Generate` |
| [`taxid`](https://pkg.go.dev/github.com/tonklasocute/thai-utils/taxid) | Thai tax ID (เลขประจำตัวผู้เสียภาษี) for individuals and juristic persons | `Validate` `Check` `Format` `Normalize` `IsJuristic` |
| [`phone`](https://pkg.go.dev/github.com/tonklasocute/thai-utils/phone) | Mobile and landline numbers: validate, normalize, E.164, display format | `Validate` `Normalize` `ToE164` `Format` `IsMobile` `IsLandline` |
| [`promptpay`](https://pkg.go.dev/github.com/tonklasocute/thai-utils/promptpay) | EMVCo / Thai QR PromptPay payloads, static or with an amount | `Payload` `Amount` `VerifyChecksum` `CRC16` |
| [`thaidate`](https://pkg.go.dev/github.com/tonklasocute/thai-utils/thaidate) | Buddhist Era (พ.ศ.) conversion, Thai month and weekday names, DD/MM/YYYY parsing, Thai digits | `ToBE` `ToCE` `Format` `Parse` `ToThaiDigits` |
| [`bahttext`](https://pkg.go.dev/github.com/tonklasocute/thai-utils/bahttext) | Amounts in Thai words, like Excel's `BAHTTEXT` | `Format` `FormatString` |

**Design principles**

- **Zero dependencies.** Standard library only, so nothing transitive lands in your `go.sum`.
- **Forgiving input, strict validation.** Dashes, spaces and `+66` prefixes are accepted; wrong checksums and impossible dates are not.
- **Typed errors.** Every package exports sentinel errors such as `idcard.ErrInvalidChecksum`, so you can tell users exactly what is wrong using `errors.Is`.
- **Verified against references.** PromptPay payloads match the [`promptpay-qr`](https://github.com/dtinth/promptpay-qr) npm package byte for byte (golden tests).
- **100% test coverage**, with table-driven tests, runnable examples and benchmarks.

## Install

```sh
go get github.com/tonklasocute/thai-utils
```

Requires Go 1.22 or newer.

## Quick start

### `idcard`: national ID

```go
import "github.com/tonklasocute/thai-utils/idcard"

idcard.Validate("1-1037-02071-81-1")  // true
idcard.Validate("1103702071812")      // false: wrong check digit
idcard.Format("1103702071811")        // "1-1037-02071-81-1", nil
idcard.Normalize("1-1037-02071-81-1") // "1103702071811"
idcard.Generate()                     // e.g. "3849105732112"

switch err := idcard.Check(input); {
case errors.Is(err, idcard.ErrInvalidLength):
	// "เลขบัตรประชาชนต้องมี 13 หลัก"
case errors.Is(err, idcard.ErrInvalidChecksum):
	// "เลขบัตรประชาชนไม่ถูกต้อง"
}
```

> [!WARNING]
> `Generate` is for test fixtures only. A random valid ID may belong to a real person.

### `taxid`: tax ID

```go
import "github.com/tonklasocute/thai-utils/taxid"

taxid.Validate("0105556000017")   // true
taxid.IsJuristic("0105556000017") // true: company / juristic person
taxid.IsJuristic("1103702071811") // false: individual
taxid.Format("0105556000017")     // "0-1055-56000-01-7", nil
```

### `phone`: phone numbers

```go
import "github.com/tonklasocute/thai-utils/phone"

phone.Validate("+66 81 234 5678")  // true
phone.Normalize("+66 81 234 5678") // "0812345678", nil
phone.ToE164("081-234-5678")       // "+66812345678", nil
phone.Format("0812345678")         // "081-234-5678", nil
phone.Format("021234567")          // "02-123-4567", nil
phone.Format("053123456")          // "053-123-456", nil
phone.IsMobile("021234567")        // false
```

Mobile numbers are `06`/`08`/`09` plus 8 digits. Landlines are `02` (Bangkok) or `03`/`04`/`05`/`07` plus 7 digits.

### `promptpay`: PromptPay QR payload

```go
import "github.com/tonklasocute/thai-utils/promptpay"

// Static QR: the payer enters the amount.
payload, err := promptpay.Payload("081-234-5678", nil)

// Dynamic QR: fixed amount of ฿100.50 to a national ID.
payload, err = promptpay.Payload("1-1037-02071-81-1", promptpay.Amount(100.50))
// 00020101021229370016A000000677010111021311037020718115802TH53037645406100.5063041B0C
```

The target can be a mobile number, a 13-digit national ID or tax ID, or a 15-digit e-wallet ID.

<details>
<summary><b>Rendering the QR image</b></summary>

`thai-utils` stays dependency-free, so it returns the payload string and leaves the image to you. With [`skip2/go-qrcode`](https://github.com/skip2/go-qrcode):

```go
import (
	"github.com/skip2/go-qrcode"
	"github.com/tonklasocute/thai-utils/promptpay"
)

payload, err := promptpay.Payload("0812345678", promptpay.Amount(250))
if err != nil {
	return err
}
// Write a 256×256 PNG that any Thai banking app can scan.
err = qrcode.WriteFile(payload, qrcode.Medium, 256, "promptpay.png")
```

To serve it over HTTP, use `qrcode.Encode` to get the PNG bytes instead.

</details>

### `thaidate`: Buddhist Era dates

```go
import "github.com/tonklasocute/thai-utils/thaidate"

t := time.Date(2026, time.October, 17, 14, 30, 0, 0, time.UTC)

thaidate.ToBE(2026)                             // 2569
thaidate.ToCE(2569)                             // 2026
thaidate.Format(t, thaidate.DateLong)           // "17 ตุลาคม 2569"
thaidate.Format(t, thaidate.DateShort)          // "17 ต.ค. 69"
thaidate.Format(t, thaidate.DateFull)           // "วันเสาร์ที่ 17 ตุลาคม พ.ศ. 2569"
thaidate.Format(t, "Mon 2 Jan 2006 15:04 น.")   // "ส. 17 ต.ค. 2569 14:30 น."
thaidate.Parse("17/10/2569")                    // 2026-10-17 00:00:00 UTC, nil
thaidate.ToThaiDigits("17/10/2569")             // "๑๗/๑๐/๒๕๖๙"
```

`Format` takes ordinary Go layouts. `2006`/`06` become Buddhist Era years, `January`/`Jan` become Thai month names and `Monday`/`Mon` become Thai weekday names. Every other token behaves exactly like `time.Format`.

### `bahttext`: amounts in Thai words

```go
import "github.com/tonklasocute/thai-utils/bahttext"

bahttext.Format(21)                    // "ยี่สิบเอ็ดบาทถ้วน"
bahttext.Format(1250.50)               // "หนึ่งพันสองร้อยห้าสิบบาทห้าสิบสตางค์"
bahttext.Format(0.25)                  // "ยี่สิบห้าสตางค์"
bahttext.Format(-100)                  // "ลบหนึ่งร้อยบาทถ้วน"
bahttext.FormatString("1,000,001.00")  // "หนึ่งล้านเอ็ดบาทถ้วน", nil
bahttext.FormatString("1000000000000") // "หนึ่งล้านล้านบาทถ้วน", nil
```

Prefer `FormatString` for money that comes from user input or a decimal type. It works on decimal strings of any size with exact half-up rounding, so `float64` precision never comes into play.

## Examples

Each package has a runnable program in [`examples/`](examples):

```sh
go run ./examples/bahttext 1250.50 99.99
go run ./examples/promptpay 0812345678 100
go run ./examples/thaidate
```

## Development

```sh
make test      # go test -race ./...
make cover     # coverage report (coverage.out, coverage.html)
make lint      # golangci-lint
make fmt       # gofmt + goimports
make bench     # benchmarks
make examples  # run every example program
```

## Contributing

Contributions are welcome, especially new Thailand-specific utilities. Ideas include address and postal-code data, vehicle license plates, bank account formats and the Thai lunar calendar.

1. Open an issue to discuss larger changes before you start.
2. Keep it dependency-free: standard library only.
3. Add table-driven tests (including invalid input) and an `Example` function for every new exported identifier. Coverage should stay at or near 100%.
4. Run `make fmt lint test` before you open a pull request.
5. Use [Conventional Commits](https://www.conventionalcommits.org/) for commit messages, e.g. `feat(phone): support 1-800 numbers`.

If you find an incorrect result, such as a misspelled Thai word or a payload a bank app rejects, please open an issue with the input and the expected output.

## License

[MIT](LICENSE) © Khomkrit Daengnuan
