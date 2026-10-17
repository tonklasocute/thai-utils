package phone

import (
	"errors"
	"testing"
)

var normalizeTests = []struct {
	name    string
	in      string
	want    string
	wantErr error
}{
	// Mobile.
	{"mobile plain", "0812345678", "0812345678", nil},
	{"mobile 06", "0612345678", "0612345678", nil},
	{"mobile 09", "0991234567", "0991234567", nil},
	{"mobile dashes", "081-234-5678", "0812345678", nil},
	{"mobile spaces", "081 234 5678", "0812345678", nil},
	{"mobile dots", "081.234.5678", "0812345678", nil},
	{"mobile +66", "+66812345678", "0812345678", nil},
	{"mobile +66 spaced", "+66 81 234 5678", "0812345678", nil},
	{"mobile +66 dashed", "+66-81-234-5678", "0812345678", nil},
	{"mobile 66", "66812345678", "0812345678", nil},
	{"mobile +66 (0)", "+66 (0)81 234 5678", "0812345678", nil},
	{"mobile 66 keeps trunk 0", "660812345678", "0812345678", nil},
	{"surrounding whitespace", "  0812345678\n", "0812345678", nil},

	// Landline.
	{"bangkok", "021234567", "021234567", nil},
	{"bangkok formatted", "02-123-4567", "021234567", nil},
	{"bangkok parens", "(02) 123-4567", "021234567", nil},
	{"bangkok +66", "+6621234567", "021234567", nil},
	{"chiang mai", "053123456", "053123456", nil},
	{"chonburi", "038-123-456", "038123456", nil},
	{"khon kaen", "043123456", "043123456", nil},
	{"phuket", "+66 76 123 456", "076123456", nil},

	// Errors.
	{"empty", "", "", ErrInvalidLength},
	{"blank", "   ", "", ErrInvalidLength},
	{"only zero", "0", "", ErrInvalidLength},
	{"only country code", "+66", "", ErrInvalidLength},
	{"letters", "08l2345678", "", ErrInvalidCharacter},
	{"plus in middle", "08+12345678", "", ErrInvalidCharacter},
	{"slash", "081/2345678", "", ErrInvalidCharacter},
	{"foreign country code", "+1 415 555 0100", "", ErrInvalidCountryCode},
	{"plus 0", "+0812345678", "", ErrInvalidCountryCode},
	{"no trunk prefix", "812345678", "", ErrInvalidPrefix},
	{"prefix 01", "0112345678", "", ErrInvalidPrefix},
	{"prefix 00", "0012345678", "", ErrInvalidPrefix},
	{"short code", "1669", "", ErrInvalidPrefix},
	{"mobile too short", "081234567", "", ErrInvalidLength},
	{"mobile too long", "08123456789", "", ErrInvalidLength},
	{"landline too short", "02123456", "", ErrInvalidLength},
	{"landline too long", "0212345678", "", ErrInvalidLength},
	{"provincial 10 digits", "0531234567", "", ErrInvalidLength},
	{"+66 mobile too short", "+6681234567", "", ErrInvalidLength},
}

func TestNormalize(t *testing.T) {
	for _, tt := range normalizeTests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Normalize(tt.in)
			if got != tt.want || !errors.Is(err, tt.wantErr) {
				t.Errorf("Normalize(%q) = (%q, %v), want (%q, %v)", tt.in, got, err, tt.want, tt.wantErr)
			}
		})
	}
}

func TestValidate(t *testing.T) {
	for _, tt := range normalizeTests {
		t.Run(tt.name, func(t *testing.T) {
			if got, want := Validate(tt.in), tt.wantErr == nil; got != want {
				t.Errorf("Validate(%q) = %v, want %v", tt.in, got, want)
			}
		})
	}
}

func TestIsMobileIsLandline(t *testing.T) {
	tests := []struct {
		in           string
		wantMobile   bool
		wantLandline bool
	}{
		{"0812345678", true, false},
		{"+66 61 234 5678", true, false},
		{"0991234567", true, false},
		{"021234567", false, true},
		{"053123456", false, true},
		{"0712345678", false, false}, // landline prefix, mobile length
		{"081234567", false, false},  // mobile prefix, landline length
		{"", false, false},
		{"abc", false, false},
	}
	for _, tt := range tests {
		if got := IsMobile(tt.in); got != tt.wantMobile {
			t.Errorf("IsMobile(%q) = %v, want %v", tt.in, got, tt.wantMobile)
		}
		if got := IsLandline(tt.in); got != tt.wantLandline {
			t.Errorf("IsLandline(%q) = %v, want %v", tt.in, got, tt.wantLandline)
		}
	}
}

func TestToE164(t *testing.T) {
	tests := []struct {
		in      string
		want    string
		wantErr error
	}{
		{"0812345678", "+66812345678", nil},
		{"081-234-5678", "+66812345678", nil},
		{"+66812345678", "+66812345678", nil},
		{"02-123-4567", "+6621234567", nil},
		{"053 123 456", "+6653123456", nil},
		{"12345", "", ErrInvalidPrefix},
		{"+44 20 7946 0958", "", ErrInvalidCountryCode},
	}
	for _, tt := range tests {
		got, err := ToE164(tt.in)
		if got != tt.want || !errors.Is(err, tt.wantErr) {
			t.Errorf("ToE164(%q) = (%q, %v), want (%q, %v)", tt.in, got, err, tt.want, tt.wantErr)
		}
	}
}

func TestFormat(t *testing.T) {
	tests := []struct {
		in      string
		want    string
		wantErr error
	}{
		{"0812345678", "081-234-5678", nil},
		{"+66 81 234 5678", "081-234-5678", nil},
		{"0612345678", "061-234-5678", nil},
		{"021234567", "02-123-4567", nil},
		{"+6621234567", "02-123-4567", nil},
		{"053123456", "053-123-456", nil},
		{"(038) 123456", "038-123-456", nil},
		{"0812", "", ErrInvalidLength},
		{"x", "", ErrInvalidCharacter},
	}
	for _, tt := range tests {
		got, err := Format(tt.in)
		if got != tt.want || !errors.Is(err, tt.wantErr) {
			t.Errorf("Format(%q) = (%q, %v), want (%q, %v)", tt.in, got, err, tt.want, tt.wantErr)
		}
	}
}

// TestRoundTrip checks that every output format parses back to the same
// local number.
func TestRoundTrip(t *testing.T) {
	for _, in := range []string{"0812345678", "021234567", "053123456"} {
		e164, err := ToE164(in)
		if err != nil {
			t.Fatal(err)
		}
		formatted, err := Format(in)
		if err != nil {
			t.Fatal(err)
		}
		for _, s := range []string{e164, formatted} {
			if got, err := Normalize(s); got != in || err != nil {
				t.Errorf("Normalize(%q) = (%q, %v), want (%q, nil)", s, got, err, in)
			}
		}
	}
}
