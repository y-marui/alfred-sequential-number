// Package seq generates sequential values (decimal, binary, octal, hex,
// alphabetic, or a custom format) over one or more ranges.
//
// This is a direct port of the original workflow/seq.py: a pattern string
// containing %b/%o/%d/%x/%X/%a/%A or #/#a/#A placeholders is expanded
// against the Cartesian product of one or more "N" or "S-E" range
// specifiers.
package seq

import (
	"fmt"
	"strconv"
	"strings"
)

const lowerLetters = "abcdefghijklmnopqrstuvwxyz"

// ConvertAlphabetLower converts a 1-indexed positive integer to a
// bijective base-26 lowercase string (1 -> "a", 26 -> "z", 27 -> "aa").
func ConvertAlphabetLower(n int) (string, error) {
	return convertAlphabet(n, lowerLetters)
}

// ConvertAlphabetUpper converts a 1-indexed positive integer to a
// bijective base-26 uppercase string (1 -> "A", 26 -> "Z", 27 -> "AA").
func ConvertAlphabetUpper(n int) (string, error) {
	return convertAlphabet(n, strings.ToUpper(lowerLetters))
}

// convertAlphabet mirrors convert_alphabet/convert_Alphabet in seq.py: the
// first digit position divides by 26, every subsequent position by 27
// (bijective base-26 has no "zero" digit, so later positions shift by one).
func convertAlphabet(n int, letters string) (string, error) {
	n--
	if n < 0 {
		return "", fmt.Errorf("n must not be negative")
	}
	if n == 0 {
		return string(letters[0]), nil
	}
	var s []byte
	for n > 0 {
		var m int
		if len(s) == 0 {
			n, m = n/26, n%26
			s = append([]byte{letters[m]}, s...)
		} else {
			n, m = n/27, n%27
			// Python indexes ascii_uppercase[m-1], where negative
			// indexing wraps m==0 to the last letter ('Z'/'z'); Go has
			// no such wraparound, so it must be done explicitly.
			idx := m - 1
			if idx < 0 {
				idx += len(letters)
			}
			s = append([]byte{letters[idx]}, s...)
		}
	}
	return string(s), nil
}

// LoadRange parses a range specifier ("N" or "S-E") into a Go-style
// half-open [start, stop) range, mirroring load_range in seq.py: "N" means
// 1..N inclusive, "S-E" means S..E inclusive. Extra "-"-separated parts
// beyond the second are ignored.
func LoadRange(text string) (start, stop int, err error) {
	if isAllDigits(text) {
		n, err := strconv.Atoi(text)
		if err != nil {
			return 0, 0, err
		}
		return 1, n + 1, nil
	}
	parts := strings.Split(text, "-")
	if len(parts) < 2 {
		return 0, 0, fmt.Errorf("invalid range %q", text)
	}
	s, err := strconv.Atoi(parts[0])
	if err != nil {
		return 0, 0, fmt.Errorf("invalid range %q: %w", text, err)
	}
	e, err := strconv.Atoi(parts[1])
	if err != nil {
		return 0, 0, fmt.Errorf("invalid range %q: %w", text, err)
	}
	return s, e + 1, nil
}

func isAllDigits(s string) bool {
	if s == "" {
		return false
	}
	for _, r := range s {
		if r < '0' || r > '9' {
			return false
		}
	}
	return true
}

// token is one piece of a parsed pattern: either a literal string or a
// format field (verb is one of 'b','o','d','x','X','a','A').
type token struct {
	literal string
	verb    byte
}

// parsePattern turns a pattern string into a sequence of tokens, mirroring
// the %/# state machine in seq.py's loader(). '%' followed by
// b/o/d/x/X/a/A becomes a format field; any other character after '%' is
// an error. '#' followed by a/A becomes a format field; '#' followed by
// anything else (or nothing) becomes a decimal field, with that other
// character emitted as a literal right after it.
func parsePattern(pattern string) ([]token, error) {
	var tokens []token
	var lit strings.Builder

	flush := func() {
		if lit.Len() > 0 {
			tokens = append(tokens, token{literal: lit.String()})
			lit.Reset()
		}
	}

	afterPercent := false
	afterHash := false

	for _, c := range pattern {
		switch {
		case afterPercent:
			switch c {
			case 'b', 'o', 'd', 'x', 'X', 'a', 'A':
				flush()
				tokens = append(tokens, token{verb: byte(c)})
			default:
				return nil, fmt.Errorf("%%b and %%o, %%d, %%x, %%X, %%a, %%A are only supported")
			}
			afterPercent = false
		case afterHash:
			switch c {
			case 'a', 'A':
				flush()
				tokens = append(tokens, token{verb: byte(c)})
			case '%':
				flush()
				tokens = append(tokens, token{verb: 'd'})
				afterPercent = true
			default:
				flush()
				tokens = append(tokens, token{verb: 'd'})
				lit.WriteRune(c)
			}
			afterHash = false
		case c == '%':
			afterPercent = true
		case c == '#':
			afterHash = true
		default:
			lit.WriteRune(c)
		}
	}
	if afterHash {
		flush()
		tokens = append(tokens, token{verb: 'd'})
	}
	flush()
	return tokens, nil
}

// render fills the format fields in tokens from args, consuming one arg
// per field in order. Extra unused trailing args are ignored (matching
// Python's str.format(), which silently ignores unused positional args) —
// this is what makes a range with more values than a pattern has fields
// repeat its output.
func render(tokens []token, args []int) (string, error) {
	var b strings.Builder
	i := 0
	for _, t := range tokens {
		if t.verb == 0 {
			b.WriteString(t.literal)
			continue
		}
		if i >= len(args) {
			return "", fmt.Errorf("not enough range values for pattern")
		}
		v := args[i]
		i++
		switch t.verb {
		case 'd':
			b.WriteString(strconv.Itoa(v))
		case 'b':
			b.WriteString(strconv.FormatInt(int64(v), 2))
		case 'o':
			b.WriteString(strconv.FormatInt(int64(v), 8))
		case 'x':
			b.WriteString(strconv.FormatInt(int64(v), 16))
		case 'X':
			b.WriteString(strings.ToUpper(strconv.FormatInt(int64(v), 16)))
		case 'a':
			s, err := ConvertAlphabetLower(v)
			if err != nil {
				return "", err
			}
			b.WriteString(s)
		case 'A':
			s, err := ConvertAlphabetUpper(v)
			if err != nil {
				return "", err
			}
			b.WriteString(s)
		}
	}
	return b.String(), nil
}

// cartesianProduct returns every combination of one value per range, with
// ranges[0] varying slowest and the last range varying fastest — matching
// multi_range in seq.py.
func cartesianProduct(ranges [][2]int) [][]int {
	result := [][]int{{}}
	for _, r := range ranges {
		var next [][]int
		for _, prefix := range result {
			for i := r[0]; i < r[1]; i++ {
				tuple := make([]int, len(prefix)+1)
				copy(tuple, prefix)
				tuple[len(prefix)] = i
				next = append(next, tuple)
			}
		}
		result = next
	}
	return result
}

// Generate expands pattern against the Cartesian product of the ranges in
// narr (each "N" or "S-E"), returning one formatted string per
// combination. narr must be non-empty.
func Generate(pattern string, narr []string) ([]string, error) {
	if len(narr) == 0 {
		return nil, fmt.Errorf("at least one range is required")
	}

	tokens, err := parsePattern(pattern)
	if err != nil {
		return nil, err
	}

	ranges := make([][2]int, len(narr))
	for i, s := range narr {
		start, stop, err := LoadRange(s)
		if err != nil {
			return nil, err
		}
		ranges[i] = [2]int{start, stop}
	}

	tuples := cartesianProduct(ranges)
	results := make([]string, 0, len(tuples))
	for _, tup := range tuples {
		s, err := render(tokens, tup)
		if err != nil {
			return nil, err
		}
		results = append(results, s)
	}
	return results, nil
}
