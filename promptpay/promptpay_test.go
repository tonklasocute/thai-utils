package promptpay

import (
	"errors"
	"math"
	"strings"
	"testing"

	"github.com/tonklasocute/thai-utils/taxid"
)

// Golden payloads generated with the promptpay-qr npm package v0.5.0
// (https://github.com/dtinth/promptpay-qr), the de facto reference
// implementation used by many Thai apps.
var golden = []struct {
	name   string
	target string
	amount *float64
	want   string
}{
	{"phone static", "0812345678", nil,
		"00020101021129370016A000000677010111011300668123456785802TH530376463045D82"},
	{"phone static 2", "0801234567", nil,
		"00020101021129370016A000000677010111011300668012345675802TH530376463046197"},
	{"phone 1 baht", "0899999999", Amount(1),
		"00020101021229370016A000000677010111011300668999999995802TH530376454041.0063041A88"},
	{"phone 100.50", "0812345678", Amount(100.5),
		"00020101021229370016A000000677010111011300668123456785802TH53037645406100.506304F88B"},
	{"phone 1234.56", "0812345678", Amount(1234.56),
		"00020101021229370016A000000677010111011300668123456785802TH530376454071234.566304D937"},
	{"phone 1 satang", "0812345678", Amount(0.01),
		"00020101021229370016A000000677010111011300668123456785802TH530376454040.01630437B6"},
	{"phone 06 prefix", "0612345678", Amount(99999.99),
		"00020101021229370016A000000677010111011300666123456785802TH5303764540899999.996304F10B"},
	{"national ID static", "1103702071811", nil,
		"00020101021129370016A000000677010111021311037020718115802TH530376463042166"},
	{"national ID amount", "1103702071811", Amount(420),
		"00020101021229370016A000000677010111021311037020718115802TH53037645406420.0063045C1A"},
	{"tax ID static", "0105556000017", nil,
		"00020101021129370016A000000677010111021301055560000175802TH53037646304F905"},
	{"tax ID amount", "0105556000017", Amount(2500.75),
		"00020101021229370016A000000677010111021301055560000175802TH530376454072500.75630470D5"},
	{"e-wallet static", "012345678901234", nil,
		"00020101021129390016A00000067701011103150123456789012345802TH530376463049781"},
	{"e-wallet amount", "012345678901234", Amount(50),
		"00020101021229390016A00000067701011103150123456789012345802TH5303764540550.0063048F3D"},
}

func TestPayloadGolden(t *testing.T) {
	for _, tt := range golden {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Payload(tt.target, tt.amount)
			if err != nil {
				t.Fatalf("Payload(%q) error: %v", tt.target, err)
			}
			if got != tt.want {
				t.Errorf("Payload(%q)\n got  %s\n want %s", tt.target, got, tt.want)
			}
			if !VerifyChecksum(got) {
				t.Errorf("VerifyChecksum(%s) = false", got)
			}
		})
	}
}

func TestPayloadTargetFormats(t *testing.T) {
	// Different spellings of the same target must produce identical payloads.
	tests := []struct {
		name      string
		inputs    []string
		canonical string
	}{
		{"phone", []string{"0812345678", "081-234-5678", "081 234 5678", "+66812345678", "+66 81 234 5678", "66812345678"}, "0812345678"},
		{"national ID", []string{"1103702071811", "1-1037-02071-81-1", "1 1037 02071 81 1"}, "1103702071811"},
		{"e-wallet", []string{"012345678901234", "012-345678901234", " 012345678901234 "}, "012345678901234"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			want, err := Payload(tt.canonical, nil)
			if err != nil {
				t.Fatal(err)
			}
			for _, in := range tt.inputs {
				got, err := Payload(in, nil)
				if err != nil || got != want {
					t.Errorf("Payload(%q) = (%q, %v), want (%q, nil)", in, got, err, want)
				}
			}
		})
	}
}

func TestPayloadErrors(t *testing.T) {
	tests := []struct {
		name    string
		target  string
		amount  *float64
		wantErr []error
	}{
		{"empty", "", nil, []error{ErrInvalidTarget}},
		{"letters", "promptpay", nil, []error{ErrInvalidTarget}},
		{"landline", "021234567", nil, []error{ErrInvalidTarget}},
		{"short number", "12345", nil, []error{ErrInvalidTarget}},
		{"11 digits non-phone", "12345678901", nil, []error{ErrInvalidTarget}},
		{"14 digits", "12345678901234", nil, []error{ErrInvalidTarget}},
		{"16 digits", "1234567890123456", nil, []error{ErrInvalidTarget}},
		{"bad ID checksum", "1103702071812", nil, []error{ErrInvalidTarget, taxid.ErrInvalidChecksum}},
		{"zero amount", "0812345678", Amount(0), []error{ErrInvalidAmount}},
		{"negative amount", "0812345678", Amount(-10), []error{ErrInvalidAmount}},
		{"rounds to zero", "0812345678", Amount(0.004), []error{ErrInvalidAmount}},
		{"NaN", "0812345678", Amount(math.NaN()), []error{ErrInvalidAmount}},
		{"+Inf", "0812345678", Amount(math.Inf(1)), []error{ErrInvalidAmount}},
		{"too large", "0812345678", Amount(1e13), []error{ErrInvalidAmount}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Payload(tt.target, tt.amount)
			if got != "" {
				t.Errorf("Payload(%q) = %q, want empty", tt.target, got)
			}
			for _, want := range tt.wantErr {
				if !errors.Is(err, want) {
					t.Errorf("Payload(%q) error = %v, want errors.Is %v", tt.target, err, want)
				}
			}
		})
	}
}

func TestPayloadAmount(t *testing.T) {
	tests := []struct {
		amount float64
		want   string
	}{
		{1, "54041.00"},
		{0.005, "54040.01"},
		{12.345678, "540512.35"},
		{9999999999.99, "54139999999999.99"}, // 13 chars: maximum allowed
	}
	for _, tt := range tests {
		got, err := Payload("0812345678", Amount(tt.amount))
		if err != nil {
			t.Fatalf("Payload(amount=%v) error: %v", tt.amount, err)
		}
		if !strings.Contains(got, "5303764"+tt.want+"6304") {
			t.Errorf("Payload(amount=%v) = %s, want amount field %s", tt.amount, got, tt.want)
		}
	}
}

func TestCRC16(t *testing.T) {
	tests := []struct {
		in   string
		want uint16
	}{
		{"", 0xFFFF},
		{"123456789", 0x29B1}, // standard CRC-16/CCITT-FALSE check value
		{"00020101021129370016A000000677010111011300668123456785802TH53037646304", 0x5D82},
	}
	for _, tt := range tests {
		if got := CRC16([]byte(tt.in)); got != tt.want {
			t.Errorf("CRC16(%q) = %04X, want %04X", tt.in, got, tt.want)
		}
	}
}

func TestVerifyChecksum(t *testing.T) {
	valid := golden[0].want
	tests := []struct {
		name string
		in   string
		want bool
	}{
		{"valid", valid, true},
		{"lowercase hex", valid[:len(valid)-4] + "5d82", true},
		{"tampered data", strings.Replace(valid, "006681", "006682", 1), false},
		{"tampered crc", valid[:len(valid)-1] + "3", false},
		{"non-hex crc", valid[:len(valid)-4] + "ZZZZ", false},
		{"missing crc tag", valid[:len(valid)-8] + "64045D82", false},
		{"too short", "6304", false},
		{"empty", "", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := VerifyChecksum(tt.in); got != tt.want {
				t.Errorf("VerifyChecksum(%q) = %v, want %v", tt.in, got, tt.want)
			}
		})
	}
}

func TestAmount(t *testing.T) {
	p := Amount(12.5)
	if p == nil || *p != 12.5 {
		t.Fatalf("Amount(12.5) = %v", p)
	}
	if Amount(1) == Amount(1) {
		t.Error("Amount must return a fresh pointer on each call")
	}
}

func BenchmarkPayload(b *testing.B) {
	amount := Amount(100)
	for i := 0; i < b.N; i++ {
		_, _ = Payload("0812345678", amount)
	}
}
