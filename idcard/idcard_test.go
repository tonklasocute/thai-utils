package idcard

import (
	"errors"
	"testing"
)

func TestCheck(t *testing.T) {
	tests := []struct {
		name string
		in   string
		want error
	}{
		{"plain", "1103702071811", nil},
		{"dashes", "1-1037-02071-81-1", nil},
		{"spaces", "1 1037 02071 81 1", nil},
		{"surrounding whitespace", "  3101001234565\n", nil},
		{"category 5", "5101234567890", nil},
		{"mixed separators", "1-2345 67890-12-1", nil},
		{"empty", "", ErrInvalidLength},
		{"only separators", "- -", ErrInvalidLength},
		{"too short", "110370207181", ErrInvalidLength},
		{"too long", "11037020718111", ErrInvalidLength},
		{"letter", "110370207181A", ErrInvalidCharacter},
		{"dots", "1.1037.02071.81.1", ErrInvalidCharacter},
		{"thai digits", "๑๑๐๓๗๐๒๐๗๑๘๑๑", ErrInvalidCharacter},
		{"leading zero", "0105556000017", ErrLeadingZero},
		{"bad checksum", "1103702071812", ErrInvalidChecksum},
		{"transposed digits", "1130702071811", ErrInvalidChecksum},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := Check(tt.in); !errors.Is(err, tt.want) {
				t.Errorf("Check(%q) = %v, want %v", tt.in, err, tt.want)
			}
		})
	}
}

func TestValidate(t *testing.T) {
	tests := []struct {
		in   string
		want bool
	}{
		{"1103702071811", true},
		{"1-1037-02071-81-1", true},
		{"3309900123451", true},
		{"1103702071810", false},
		{"0105556000017", false},
		{"abc", false},
		{"", false},
	}
	for _, tt := range tests {
		if got := Validate(tt.in); got != tt.want {
			t.Errorf("Validate(%q) = %v, want %v", tt.in, got, tt.want)
		}
	}
}

func TestFormat(t *testing.T) {
	tests := []struct {
		in      string
		want    string
		wantErr error
	}{
		{"1103702071811", "1-1037-02071-81-1", nil},
		{"1-1037-02071-81-1", "1-1037-02071-81-1", nil},
		{"3 1010 01234 56 5", "3-1010-01234-56-5", nil},
		{"1103702071812", "", ErrInvalidChecksum},
		{"12345", "", ErrInvalidLength},
		{"x", "", ErrInvalidCharacter},
	}
	for _, tt := range tests {
		got, err := Format(tt.in)
		if got != tt.want || !errors.Is(err, tt.wantErr) {
			t.Errorf("Format(%q) = (%q, %v), want (%q, %v)", tt.in, got, err, tt.want, tt.wantErr)
		}
	}
}

func TestNormalize(t *testing.T) {
	tests := []struct{ in, want string }{
		{"1-1037-02071-81-1", "1103702071811"},
		{"1 1037 02071 81 1", "1103702071811"},
		{"ID: 1103702071811.", "1103702071811"},
		{"no digits", ""},
		{"", ""},
	}
	for _, tt := range tests {
		if got := Normalize(tt.in); got != tt.want {
			t.Errorf("Normalize(%q) = %q, want %q", tt.in, got, tt.want)
		}
	}
}

func TestGenerate(t *testing.T) {
	for i := 0; i < 1000; i++ {
		id := Generate()
		if len(id) != 13 {
			t.Fatalf("Generate() = %q, want 13 digits", id)
		}
		if id[0] < '1' || id[0] > '8' {
			t.Fatalf("Generate() = %q, first digit must be 1-8", id)
		}
		if err := Check(id); err != nil {
			t.Fatalf("Generate() = %q, Check: %v", id, err)
		}
	}
}

func BenchmarkValidate(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Validate("1-1037-02071-81-1")
	}
}
