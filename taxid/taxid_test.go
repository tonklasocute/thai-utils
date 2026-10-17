package taxid

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
		{"juristic", "0105556000017", nil},
		{"juristic dashes", "0-1055-36092-64-1", nil},
		{"government", "0994000158785", nil},
		{"individual", "1103702071811", nil},
		{"individual spaces", "3 1010 01234 56 5", nil},
		{"empty", "", ErrInvalidLength},
		{"too short", "010555600001", ErrInvalidLength},
		{"too long", "01055560000170", ErrInvalidLength},
		{"letter", "O105556000017", ErrInvalidCharacter},
		{"slash", "0105556/000017", ErrInvalidCharacter},
		{"bad checksum", "0105556000018", ErrInvalidChecksum},
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
		{"0105556000017", true},
		{"1103702071811", true},
		{"0105556000010", false},
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
		{"0105556000017", "0-1055-56000-01-7", nil},
		{"0 1055 36092 64 1", "0-1055-36092-64-1", nil},
		{"1103702071811", "1-1037-02071-81-1", nil},
		{"0105556000018", "", ErrInvalidChecksum},
		{"123", "", ErrInvalidLength},
		{"abc", "", ErrInvalidCharacter},
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
		{"0-1055-56000-01-7", "0105556000017"},
		{"Tax ID: 0105556000017", "0105556000017"},
		{"", ""},
	}
	for _, tt := range tests {
		if got := Normalize(tt.in); got != tt.want {
			t.Errorf("Normalize(%q) = %q, want %q", tt.in, got, tt.want)
		}
	}
}

func TestIsJuristic(t *testing.T) {
	tests := []struct {
		in   string
		want bool
	}{
		{"0105556000017", true},
		{"0-9940-00158-78-5", true},
		{"1103702071811", false},
		{"0105556000018", false}, // starts with 0 but invalid checksum
		{"0", false},
		{"", false},
	}
	for _, tt := range tests {
		if got := IsJuristic(tt.in); got != tt.want {
			t.Errorf("IsJuristic(%q) = %v, want %v", tt.in, got, tt.want)
		}
	}
}
