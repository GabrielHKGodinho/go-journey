package encodedecode

import (
	"errors"
	"slices"
	"strconv"
	"strings"
	"testing"
)

type codec struct {
	encode func([]string) string
	decode func(string) ([]string, error)
}

var codecs = map[string]codec{
	"length-prefix": {encode, decode},
	"escape":        {encodeEscape, decodeEscape},
}

func TestRoundTrip(t *testing.T) {
	tests := []struct {
		name  string
		words []string
	}{
		{"example", []string{"neet", "code", "love", "you"}},
		{"empty list", []string{}},
		{"one empty string", []string{""}},
		{"only empty strings", []string{"", ""}},
		{"delimiter inside words", []string{"4;", "", "ab;12"}},
		{"word ending in backslash", []string{`a\`, "b"}},
		{"only special characters", []string{`\;`, ";;", `\\`}},
		{"word starting with digits", []string{"12abc", "9"}},
		{"length with two digits", []string{strings.Repeat("x", 15), "y"}},
		{"multi-byte characters", []string{"olá", "日本"}},
	}

	for codecName, c := range codecs {
		for _, tt := range tests {
			t.Run(codecName+"/"+tt.name, func(t *testing.T) {
				encoded := c.encode(tt.words)
				got, err := c.decode(encoded)
				if err != nil {
					t.Fatalf("decode(%q) returned error: %v", encoded, err)
				}
				if !slices.Equal(got, tt.words) {
					t.Errorf("decode(encode(%q)) = %q, want %q", tt.words, got, tt.words)
				}
			})
		}
	}
}

// The round trip alone would pass for a pair of functions that agree on a
// wrong format, so the wire format is pinned down here too.
func TestEncodeFormat(t *testing.T) {
	tests := []struct {
		name  string
		words []string
		want  string
	}{
		{"example", []string{"neet", "code"}, "4;neet4;code"},
		{"delimiter inside words", []string{"4;", "", "ab;12"}, "2;4;0;5;ab;12"},
		{"length counts bytes, not runes", []string{"olá"}, "4;olá"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := encode(tt.words); got != tt.want {
				t.Errorf("encode(%q) = %q, want %q", tt.words, got, tt.want)
			}
		})
	}
}

func TestDecodeMalformed(t *testing.T) {
	tests := []struct {
		name    string
		decode  func(string) ([]string, error)
		encoded string
	}{
		{"missing delimiter", decode, "4neet"},
		{"non-numeric length", decode, "x;abc"},
		{"empty length", decode, ";abc"},
		{"negative length", decode, "-1;a"},
		{"plus sign", decode, "+1;a"},
		{"length too long", decode, "9;abc"},
		{"dangling escape", decodeEscape, `abc\`},
		{"unterminated word", decodeEscape, "abc;def"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.decode(tt.encoded)
			if !errors.Is(err, ErrMalformed) {
				t.Errorf("decode(%q) = %q, %v; want ErrMalformed", tt.encoded, got, err)
			}
		})
	}
}

// encodeConcat is encode written with += instead of strings.Builder.
// It exists only to be measured against encode.
func encodeConcat(words []string) string {
	encoded := ""
	for _, word := range words {
		encoded += strconv.Itoa(len(word)) + ";" + word
	}
	return encoded
}

func benchmarkInput() []string {
	words := make([]string, 1000)
	for i := range words {
		words[i] = strings.Repeat("a;b\\", 16) // 64 bytes, half of them need escaping
	}
	return words
}

var (
	sinkString string
	sinkWords  []string
)

func BenchmarkEncode(b *testing.B) {
	words := benchmarkInput()
	encoders := map[string]func([]string) string{
		"builder": encode,
		"concat":  encodeConcat,
		"escape":  encodeEscape,
	}
	for name, enc := range encoders {
		b.Run(name, func(b *testing.B) {
			b.ReportAllocs()
			for b.Loop() {
				sinkString = enc(words)
			}
		})
	}
}

func BenchmarkDecode(b *testing.B) {
	words := benchmarkInput()
	for name, c := range codecs {
		encoded := c.encode(words)
		b.Run(name, func(b *testing.B) {
			b.ReportAllocs()
			for b.Loop() {
				sinkWords, _ = c.decode(encoded)
			}
		})
	}
}
