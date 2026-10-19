// Package base62 converts non-negative integers to and from base-62 strings,
// the format used for URL shortener codes.
package base62

import (
	"errors"
	"fmt"
	"math"
	"strings"
)

// alphabet maps each digit value (0-61) to its character.
const alphabet = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

const (
	base = uint64(len(alphabet))

	// maxLen is how many base-62 digits the largest uint64 needs:
	// 62^10 < math.MaxUint64 < 62^11.
	maxLen = 11
)

var (
	// ErrInvalidChar is returned when a string contains a character outside the alphabet.
	ErrInvalidChar = errors.New("base62: invalid character")
	// ErrEmpty is returned when decoding an empty string.
	ErrEmpty = errors.New("base62: empty string")
	// ErrOverflow is returned when a string represents a number larger than math.MaxUint64.
	ErrOverflow = errors.New("base62: value overflows uint64")
)

// Encode returns the base-62 representation of n.
func Encode(n uint64) string {
	if n == 0 {
		return alphabet[:1] // "0": the loop below would never run
	}

	// Successive division yields the digits from least to most significant,
	// so the buffer is filled from the end and comes out in reading order.
	var buf [maxLen]byte
	i := len(buf)
	for n > 0 {
		i--
		buf[i] = alphabet[n%base]
		n /= base
	}
	return string(buf[i:]) // only the digits written, not the unused prefix
}

// Decode returns the number represented by the base-62 string s.
func Decode(s string) (uint64, error) {
	if s == "" {
		return 0, ErrEmpty
	}

	var n uint64
	for i := 0; i < len(s); i++ {
		digit := strings.IndexByte(alphabet, s[i])
		if digit == -1 {
			return 0, fmt.Errorf("%w %q at position %d", ErrInvalidChar, s[i], i)
		}

		// Reject before multiplying: uint64 arithmetic wraps around silently.
		// n*base + digit <= MaxUint64  <=>  n <= (MaxUint64 - digit) / base
		if n > (math.MaxUint64-uint64(digit))/base {
			return 0, fmt.Errorf("%w: %q", ErrOverflow, s)
		}

		// Horner's method: left to right, each step shifts the number one
		// base-62 position and adds the new digit.
		n = n*base + uint64(digit)
	}
	return n, nil
}
