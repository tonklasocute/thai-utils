<div align="center">

# thai-utils

**ไลบรารี Go สำหรับข้อมูลเฉพาะของไทย: เลขบัตรประชาชน เลขผู้เสียภาษี เบอร์โทรศัพท์ QR พร้อมเพย์ วันที่แบบ พ.ศ. และการอ่านจำนวนเงินเป็นตัวอักษร**

[![Go Reference](https://pkg.go.dev/badge/github.com/tonklasocute/thai-utils.svg)](https://pkg.go.dev/github.com/tonklasocute/thai-utils)
[![CI](https://github.com/tonklasocute/thai-utils/actions/workflows/ci.yml/badge.svg)](https://github.com/tonklasocute/thai-utils/actions/workflows/ci.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/tonklasocute/thai-utils)](https://goreportcard.com/report/github.com/tonklasocute/thai-utils)
[![License: MIT](https://img.shields.io/badge/license-MIT-blue.svg)](LICENSE)
![Go version](https://img.shields.io/github/go-mod/go-version/tonklasocute/thai-utils)
![Dependencies](https://img.shields.io/badge/dependencies-0-brightgreen)

[English](README.md) · ภาษาไทย

</div>

---

```go
idcard.Validate("1-1037-02071-81-1")  // true
phone.ToE164("081-234-5678")          // "+66812345678"
promptpay.Payload("0812345678", nil)  // "000201010211…63045D82"
thaidate.Format(t, thaidate.DateFull) // "วันเสาร์ที่ 17 ตุลาคม พ.ศ. 2569"
bahttext.Format(1250.50)              // "หนึ่งพันสองร้อยห้าสิบบาทห้าสิบสตางค์"
```

## ความสามารถ

| แพ็กเกจ | ทำอะไรได้ | ฟังก์ชันหลัก |
|---|---|---|
| [`idcard`](https://pkg.go.dev/github.com/tonklasocute/thai-utils/idcard) | เลขประจำตัวประชาชน: ตรวจ check digit, จัดรูปแบบ, สุ่มเลขสำหรับข้อมูลทดสอบ | `Validate` `Check` `Format` `Normalize` `Generate` |
| [`taxid`](https://pkg.go.dev/github.com/tonklasocute/thai-utils/taxid) | เลขประจำตัวผู้เสียภาษีอากร ทั้งบุคคลธรรมดาและนิติบุคคล | `Validate` `Check` `Format` `Normalize` `IsJuristic` |
| [`phone`](https://pkg.go.dev/github.com/tonklasocute/thai-utils/phone) | เบอร์มือถือและเบอร์บ้าน: ตรวจสอบ, แปลงเป็นรูปแบบมาตรฐาน, E.164, รูปแบบสำหรับแสดงผล | `Validate` `Normalize` `ToE164` `Format` `IsMobile` `IsLandline` |
| [`promptpay`](https://pkg.go.dev/github.com/tonklasocute/thai-utils/promptpay) | สร้าง payload QR พร้อมเพย์ตามมาตรฐาน EMVCo / Thai QR Payment ทั้งแบบไม่ระบุและระบุจำนวนเงิน | `Payload` `Amount` `VerifyChecksum` `CRC16` |
| [`thaidate`](https://pkg.go.dev/github.com/tonklasocute/thai-utils/thaidate) | แปลง พ.ศ./ค.ศ., ชื่อเดือนและวันภาษาไทย, แปลงข้อความ วว/ดด/ปปปป, เลขไทย | `ToBE` `ToCE` `Format` `Parse` `ToThaiDigits` |
| [`bahttext`](https://pkg.go.dev/github.com/tonklasocute/thai-utils/bahttext) | อ่านจำนวนเงินเป็นตัวอักษรภาษาไทย เหมือนฟังก์ชัน `BAHTTEXT` ของ Excel | `Format` `FormatString` |

**หลักการออกแบบ**

- **ไม่มี dependency ภายนอก** ใช้แค่ standard library จึงไม่มีแพ็กเกจอื่นติดเข้ามาใน `go.sum` ของคุณ
- **รับ input ได้หลายแบบ แต่ตรวจสอบเข้มงวด** มีขีด เว้นวรรค หรือขึ้นต้นด้วย `+66` ก็ใช้ได้ แต่ check digit ผิดหรือวันที่ที่ไม่มีอยู่จริงจะไม่ผ่าน
- **error แยกชนิดชัดเจน** ทุกแพ็กเกจมี sentinel error เช่น `idcard.ErrInvalidChecksum` ใช้ `errors.Is` เพื่อบอกผู้ใช้ได้ตรงจุดว่าผิดตรงไหน
- **ตรวจกับ implementation อ้างอิง** payload พร้อมเพย์ตรงกับแพ็กเกจ npm [`promptpay-qr`](https://github.com/dtinth/promptpay-qr) ทุกไบต์ (golden tests)
- **test coverage 100%** มีทั้ง table-driven tests, example ที่รันได้ และ benchmark

## ติดตั้ง

```sh
go get github.com/tonklasocute/thai-utils
```

ต้องใช้ Go 1.22 ขึ้นไป

## เริ่มต้นใช้งาน

### `idcard`: เลขบัตรประชาชน

```go
import "github.com/tonklasocute/thai-utils/idcard"

idcard.Validate("1-1037-02071-81-1")  // true
idcard.Validate("1103702071812")      // false: check digit ผิด
idcard.Format("1103702071811")        // "1-1037-02071-81-1", nil
idcard.Normalize("1-1037-02071-81-1") // "1103702071811"
idcard.Generate()                     // เช่น "3849105732112"

switch err := idcard.Check(input); {
case errors.Is(err, idcard.ErrInvalidLength):
	// "เลขบัตรประชาชนต้องมี 13 หลัก"
case errors.Is(err, idcard.ErrInvalidChecksum):
	// "เลขบัตรประชาชนไม่ถูกต้อง"
}
```

> [!WARNING]
> `Generate` มีไว้สำหรับข้อมูลทดสอบเท่านั้น เลขที่สุ่มได้อาจตรงกับเลขของคนจริง

### `taxid`: เลขประจำตัวผู้เสียภาษี

```go
import "github.com/tonklasocute/thai-utils/taxid"

taxid.Validate("0105556000017")   // true
taxid.IsJuristic("0105556000017") // true: นิติบุคคล
taxid.IsJuristic("1103702071811") // false: บุคคลธรรมดา
taxid.Format("0105556000017")     // "0-1055-56000-01-7", nil
```

### `phone`: เบอร์โทรศัพท์

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

เบอร์มือถือขึ้นต้นด้วย `06`/`08`/`09` ตามด้วยตัวเลข 8 หลัก เบอร์บ้านขึ้นต้นด้วย `02` (กรุงเทพฯ และปริมณฑล) หรือ `03`/`04`/`05`/`07` ตามด้วยตัวเลข 7 หลัก

### `promptpay`: payload QR พร้อมเพย์

```go
import "github.com/tonklasocute/thai-utils/promptpay"

// QR แบบไม่ระบุจำนวนเงิน (ผู้จ่ายกรอกเอง)
payload, err := promptpay.Payload("081-234-5678", nil)

// QR แบบระบุจำนวนเงิน 100.50 บาท ไปยังเลขบัตรประชาชน
payload, err = promptpay.Payload("1-1037-02071-81-1", promptpay.Amount(100.50))
// 00020101021229370016A000000677010111021311037020718115802TH53037645406100.5063041B0C
```

ปลายทางเป็นเบอร์มือถือ เลขบัตรประชาชนหรือเลขผู้เสียภาษี 13 หลัก หรือเลข e-Wallet 15 หลักก็ได้

<details>
<summary><b>วิธีสร้างภาพ QR</b></summary>

`thai-utils` ไม่มี dependency จึงคืนค่าเป็นข้อความ payload แล้วให้คุณเลือกไลบรารีสร้างภาพเอง ตัวอย่างนี้ใช้ [`skip2/go-qrcode`](https://github.com/skip2/go-qrcode):

```go
import (
	"github.com/skip2/go-qrcode"
	"github.com/tonklasocute/thai-utils/promptpay"
)

payload, err := promptpay.Payload("0812345678", promptpay.Amount(250))
if err != nil {
	return err
}
// เขียนไฟล์ PNG ขนาด 256×256 ที่สแกนได้ด้วยแอปธนาคารทุกแอป
err = qrcode.WriteFile(payload, qrcode.Medium, 256, "promptpay.png")
```

ถ้าจะส่งภาพผ่าน HTTP ให้ใช้ `qrcode.Encode` เพื่อรับเป็น PNG bytes แทน

</details>

### `thaidate`: วันที่แบบ พ.ศ.

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

`Format` ใช้ layout แบบเดียวกับ Go ปกติ โดย `2006`/`06` จะเป็นปี พ.ศ., `January`/`Jan` เป็นชื่อเดือนภาษาไทย และ `Monday`/`Mon` เป็นชื่อวันภาษาไทย token อื่น ๆ ทำงานเหมือน `time.Format` ทุกอย่าง

### `bahttext`: อ่านจำนวนเงินเป็นตัวอักษร

```go
import "github.com/tonklasocute/thai-utils/bahttext"

bahttext.Format(21)                    // "ยี่สิบเอ็ดบาทถ้วน"
bahttext.Format(1250.50)               // "หนึ่งพันสองร้อยห้าสิบบาทห้าสิบสตางค์"
bahttext.Format(0.25)                  // "ยี่สิบห้าสตางค์"
bahttext.Format(-100)                  // "ลบหนึ่งร้อยบาทถ้วน"
bahttext.FormatString("1,000,001.00")  // "หนึ่งล้านเอ็ดบาทถ้วน", nil
bahttext.FormatString("1000000000000") // "หนึ่งล้านล้านบาทถ้วน", nil
```

ถ้าจำนวนเงินมาจาก input ของผู้ใช้หรือชนิดข้อมูลทศนิยม แนะนำให้ใช้ `FormatString` เพราะทำงานกับตัวเลขทศนิยมแบบข้อความได้ไม่จำกัดขนาด และปัดเศษสตางค์แบบปัดครึ่งขึ้นได้แม่นยำ ไม่ติดปัญหาความละเอียดของ `float64`

## ตัวอย่างโปรแกรม

ทุกแพ็กเกจมีโปรแกรมตัวอย่างที่รันได้ใน [`examples/`](examples)

```sh
go run ./examples/bahttext 1250.50 99.99
go run ./examples/promptpay 0812345678 100
go run ./examples/thaidate
```

## การพัฒนา

```sh
make test      # go test -race ./...
make cover     # รายงาน coverage (coverage.out, coverage.html)
make lint      # golangci-lint
make fmt       # gofmt + goimports
make bench     # benchmark
make examples  # รันโปรแกรมตัวอย่างทั้งหมด
```

## ร่วมพัฒนา

ยินดีรับ contribution ทุกรูปแบบ โดยเฉพาะยูทิลิตีใหม่ ๆ ที่เกี่ยวกับข้อมูลของไทย เช่น ที่อยู่และรหัสไปรษณีย์ ป้ายทะเบียนรถ รูปแบบเลขบัญชีธนาคาร หรือปฏิทินจันทรคติไทย

1. ถ้าเป็นการเปลี่ยนแปลงใหญ่ เปิด issue พูดคุยกันก่อนเริ่มทำ
2. ไม่เพิ่ม dependency ใช้แค่ standard library
3. เขียน table-driven tests (รวมกรณี input ไม่ถูกต้อง) และฟังก์ชัน `Example` สำหรับทุก identifier ที่ export ใหม่ coverage ควรอยู่ที่หรือใกล้ 100%
4. รัน `make fmt lint test` ก่อนเปิด pull request
5. เขียน commit message ตาม [Conventional Commits](https://www.conventionalcommits.org/) เช่น `feat(phone): support 1-800 numbers`

ถ้าเจอผลลัพธ์ที่ไม่ถูกต้อง เช่น สะกดคำภาษาไทยผิด หรือ payload ที่แอปธนาคารไม่รับ ช่วยเปิด issue พร้อม input และผลลัพธ์ที่ควรได้

## สัญญาอนุญาต

[MIT](LICENSE) © Khomkrit Daengnuan
