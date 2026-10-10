package encodedecode

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

// ErrMalformed is returned when a decoder receives input that its encoder
// could never have produced.
var ErrMalformed = errors.New("malformed encoded string")

// encode writes every word as "<length>;<word>", where length is in bytes.
// The decoder only searches for ';' inside the length field, which holds
// digits only, so the words themselves may contain any character.
func encode(words []string) string {
	var sb strings.Builder
	for _, word := range words {
		sb.WriteString(strconv.Itoa(len(word)))
		sb.WriteByte(';')
		sb.WriteString(word)
	}
	return sb.String()
}

// decode reverses encode. It reads the length up to the first ';' and then
// skips exactly that many bytes, so the content of a word is never scanned.
// Each returned word is a substring of encoded: no bytes are copied.
func decode(encoded string) ([]string, error) {
	words := []string{}
	i := 0
	for i < len(encoded) {
		offset := strings.IndexByte(encoded[i:], ';')
		if offset == -1 {
			return nil, fmt.Errorf("%w: missing ';' after length at byte %d", ErrMalformed, i)
		}
		sep := i + offset

		field := encoded[i:sep]
		length, err := strconv.Atoi(field)
		// Atoi accepts "-3" and "+3"; encode only ever writes plain digits.
		if err != nil || field[0] == '-' || field[0] == '+' {
			return nil, fmt.Errorf("%w: invalid length %q at byte %d", ErrMalformed, field, i)
		}

		start := sep + 1
		if length > len(encoded)-start {
			return nil, fmt.Errorf("%w: length %d at byte %d exceeds remaining input", ErrMalformed, length, i)
		}
		words = append(words, encoded[start:start+length])
		i = start + length
	}
	return words, nil
}

// encodeEscape is the escaping alternative: ';' terminates every word, and
// both ';' and the escape character '\' are escaped inside a word.
// A terminator after every word (not a separator between words) keeps
// []string{} ("") distinct from []string{""} (";").
func encodeEscape(words []string) string {
	var sb strings.Builder
	for _, word := range words {
		for i := 0; i < len(word); i++ {
			if word[i] == '\\' || word[i] == ';' {
				sb.WriteByte('\\')
			}
			sb.WriteByte(word[i])
		}
		sb.WriteByte(';')
	}
	return sb.String()
}

// decodeEscape reverses encodeEscape. Unlike decode, it must look at every
// byte and copy each word, because the output differs from the input.
func decodeEscape(encoded string) ([]string, error) {
	words := []string{}
	var sb strings.Builder
	for i := 0; i < len(encoded); i++ {
		switch encoded[i] {
		case '\\':
			i++
			if i == len(encoded) {
				return nil, fmt.Errorf("%w: dangling escape at the end of input", ErrMalformed)
			}
			sb.WriteByte(encoded[i])
		case ';':
			words = append(words, sb.String())
			sb.Reset()
		default:
			sb.WriteByte(encoded[i])
		}
	}
	if sb.Len() > 0 {
		return nil, fmt.Errorf("%w: last word has no terminating ';'", ErrMalformed)
	}
	return words, nil
}
