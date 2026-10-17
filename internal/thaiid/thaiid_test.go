package thaiid

import "testing"

func TestClean(t *testing.T) {
	tests := []struct {
		in     string
		want   string
		wantOK bool
	}{
		{"1103702071811", "1103702071811", true},
		{"1-1037-02071-81-1", "1103702071811", true},
		{" 1 1037 02071 81 1\t", "1103702071811", true},
		{"", "", true},
		{"1103702071811x", "", false},
		{"1.1037.02071.81.1", "", false},
		{"๑๑๐๓๗๐๒๐๗๑๘๑๑", "", false},
	}
	for _, tt := range tests {
		got, ok := Clean(tt.in)
		if got != tt.want || ok != tt.wantOK {
			t.Errorf("Clean(%q) = (%q, %v), want (%q, %v)", tt.in, got, ok, tt.want, tt.wantOK)
		}
	}
}

func TestDigits(t *testing.T) {
	tests := []struct{ in, want string }{
		{"1-1037-02071-81-1", "1103702071811"},
		{"abc", ""},
		{"", ""},
		{"a1b2c3", "123"},
	}
	for _, tt := range tests {
		if got := Digits(tt.in); got != tt.want {
			t.Errorf("Digits(%q) = %q, want %q", tt.in, got, tt.want)
		}
	}
}

func TestCheckDigit(t *testing.T) {
	tests := []struct {
		in   string
		want byte
	}{
		{"110370207181", '1'},
		{"000000000000", '1'}, // (11 - 0) % 10
		{"010555600001", '7'},
	}
	for _, tt := range tests {
		if got := CheckDigit(tt.in); got != tt.want {
			t.Errorf("CheckDigit(%q) = %q, want %q", tt.in, got, tt.want)
		}
	}
}

func TestValid(t *testing.T) {
	tests := []struct {
		in   string
		want bool
	}{
		{"1103702071811", true},
		{"0000000000001", true},
		{"1103702071812", false},
		{"110370207181", false},
		{"11037020718111", false},
		{"110370207181a", false},
		{"", false},
	}
	for _, tt := range tests {
		if got := Valid(tt.in); got != tt.want {
			t.Errorf("Valid(%q) = %v, want %v", tt.in, got, tt.want)
		}
	}
}

func TestFormat(t *testing.T) {
	if got, want := Format("1103702071811"), "1-1037-02071-81-1"; got != want {
		t.Errorf("Format = %q, want %q", got, want)
	}
}
