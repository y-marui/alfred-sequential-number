package seq

import (
	"reflect"
	"testing"
)

func TestConvertAlphabetLower(t *testing.T) {
	cases := []struct {
		n    int
		want string
	}{
		{1, "a"},
		{26, "z"},
		{27, "aa"},
		{703, "aza"},
	}
	for _, c := range cases {
		got, err := ConvertAlphabetLower(c.n)
		if err != nil {
			t.Fatalf("ConvertAlphabetLower(%d): unexpected error: %v", c.n, err)
		}
		if got != c.want {
			t.Errorf("ConvertAlphabetLower(%d) = %q, want %q", c.n, got, c.want)
		}
	}
}

func TestConvertAlphabetLowerNegativeErrors(t *testing.T) {
	if _, err := ConvertAlphabetLower(0); err == nil {
		t.Fatal("ConvertAlphabetLower(0): expected error, got nil")
	}
}

func TestConvertAlphabetUpper(t *testing.T) {
	cases := []struct {
		n    int
		want string
	}{
		{1, "A"},
		{26, "Z"},
		{27, "AA"},
		// Boundary around the m==0 case, where the original Python
		// implementation relies on negative-index wraparound
		// (ascii_uppercase[-1] == "Z"); verified against the original
		// convert_Alphabet(n) before it was removed.
		{700, "ZX"},
		{701, "ZY"},
		{702, "ZZ"},
		{703, "AZA"},
		{704, "AZB"},
		{730, "AAB"},
		{731, "AAC"},
		{18278, "ZZZ"},
	}
	for _, c := range cases {
		got, err := ConvertAlphabetUpper(c.n)
		if err != nil {
			t.Fatalf("ConvertAlphabetUpper(%d): unexpected error: %v", c.n, err)
		}
		if got != c.want {
			t.Errorf("ConvertAlphabetUpper(%d) = %q, want %q", c.n, got, c.want)
		}
	}
}

func TestLoadRange(t *testing.T) {
	cases := []struct {
		text      string
		start     int
		stop      int
		wantError bool
	}{
		{"5", 1, 6, false},
		{"3-5", 3, 6, false},
		{"abc", 0, 0, true},
	}
	for _, c := range cases {
		start, stop, err := LoadRange(c.text)
		if c.wantError {
			if err == nil {
				t.Errorf("LoadRange(%q): expected error, got nil", c.text)
			}
			continue
		}
		if err != nil {
			t.Fatalf("LoadRange(%q): unexpected error: %v", c.text, err)
		}
		if start != c.start || stop != c.stop {
			t.Errorf("LoadRange(%q) = (%d, %d), want (%d, %d)", c.text, start, stop, c.start, c.stop)
		}
	}
}

func TestGenerate(t *testing.T) {
	cases := []struct {
		name    string
		pattern string
		narr    []string
		want    []string
	}{
		{"decimal", "%d", []string{"3"}, []string{"1", "2", "3"}},
		{"binary", "%b", []string{"3"}, []string{"1", "10", "11"}},
		{"hash decimal", "#", []string{"3"}, []string{"1", "2", "3"}},
		{"hash alpha lower", "#a", []string{"3"}, []string{"a", "b", "c"}},
		{"hash alpha upper", "#A", []string{"3"}, []string{"A", "B", "C"}},
		{"custom format", "item-%d", []string{"3"}, []string{"item-1", "item-2", "item-3"}},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got, err := Generate(c.pattern, c.narr)
			if err != nil {
				t.Fatalf("Generate(%q, %v): unexpected error: %v", c.pattern, c.narr, err)
			}
			if !reflect.DeepEqual(got, c.want) {
				t.Errorf("Generate(%q, %v) = %v, want %v", c.pattern, c.narr, got, c.want)
			}
		})
	}
}

// TestGenerateExtraRangeIgnored preserves the original Python
// implementation's quirk: when a pattern has fewer format fields than
// narr has ranges, str.format() (and this port) silently ignores the
// extra positional args, so the outer range's value repeats once per
// inner range value. This is reachable in practice via e.g. "seq bin 3 5".
func TestGenerateExtraRangeIgnored(t *testing.T) {
	got, err := Generate("%b", []string{"3", "5"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	want := []string{
		"1", "1", "1", "1", "1",
		"10", "10", "10", "10", "10",
		"11", "11", "11", "11", "11",
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("Generate(%%b, [3 5]) = %v, want %v", got, want)
	}
}

func TestGenerateInvalidPercentVerb(t *testing.T) {
	if _, err := Generate("%z", []string{"3"}); err == nil {
		t.Fatal("expected error for invalid %z verb, got nil")
	}
}

func TestGenerateEmptyNarr(t *testing.T) {
	if _, err := Generate("%d", nil); err == nil {
		t.Fatal("expected error for empty narr, got nil")
	}
}
