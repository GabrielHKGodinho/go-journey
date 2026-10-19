package base62

import (
	"errors"
	"math"
	"testing"
)

func TestEncodeDecode(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		n    uint64
		code string
	}{
		{name: "zero", n: 0, code: "0"},
		{name: "largest single digit", n: 61, code: "Z"},
		{name: "first two-digit number", n: 62, code: "10"},
		{name: "two digits", n: 125, code: "21"},
		{name: "largest six-digit code", n: 56_800_235_584 - 1, code: "ZZZZZZ"},
		{name: "first seven-digit code (62^6)", n: 56_800_235_584, code: "1000000"},
		{name: "largest seven-digit code (62^7 - 1)", n: 3_521_614_606_207, code: "ZZZZZZZ"},
		{name: "max uint64 fits in maxLen digits", n: math.MaxUint64, code: "lYGhA16ahyf"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			if got := Encode(tt.n); got != tt.code {
				t.Errorf("Encode(%d) = %q, want %q", tt.n, got, tt.code)
			}

			got, err := Decode(tt.code)
			if err != nil {
				t.Fatalf("Decode(%q) unexpected error: %v", tt.code, err)
			}
			if got != tt.n {
				t.Errorf("Decode(%q) = %d, want %d", tt.code, got, tt.n)
			}
		})
	}
}

func TestDecodeErrors(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		input   string
		wantErr error
	}{
		{name: "empty string", input: "", wantErr: ErrEmpty},
		{name: "symbol in the middle", input: "ab-c", wantErr: ErrInvalidChar},
		{name: "trailing space", input: "abc ", wantErr: ErrInvalidChar},
		{name: "non-ASCII letter", input: "ação", wantErr: ErrInvalidChar},
		{name: "one past max uint64", input: "lYGhA16ahyg", wantErr: ErrOverflow},
		{name: "too many digits", input: "100000000000", wantErr: ErrOverflow},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got, err := Decode(tt.input)
			if !errors.Is(err, tt.wantErr) {
				t.Fatalf("Decode(%q) error = %v, want %v", tt.input, err, tt.wantErr)
			}
			if got != 0 {
				t.Errorf("Decode(%q) = %d on error, want 0", tt.input, got)
			}
		})
	}
}
