package seqcmd

import (
	"strings"
	"testing"
)

func values(t *testing.T, arg string) []string {
	t.Helper()
	if arg == "" {
		return nil
	}
	return strings.Split(arg, "\r\n")
}

func firstItem(t *testing.T, query string) (title, subtitle, arg string, valid bool) {
	t.Helper()
	resp := Dispatch(query)
	if len(resp.Items) != 1 {
		t.Fatalf("Dispatch(%q): got %d items, want 1: %+v", query, len(resp.Items), resp.Items)
	}
	item := resp.Items[0]
	v := item.Valid != nil && *item.Valid
	return item.Title, item.Subtitle, item.Arg, v
}

func assertValues(t *testing.T, query string, want []string) {
	t.Helper()
	_, _, arg, valid := firstItem(t, query)
	if !valid {
		t.Fatalf("Dispatch(%q): item is not valid", query)
	}
	got := values(t, arg)
	if len(got) != len(want) {
		t.Fatalf("Dispatch(%q) values = %v, want %v", query, got, want)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Errorf("Dispatch(%q) values[%d] = %q, want %q", query, i, got[i], want[i])
		}
	}
}

func TestDecimal(t *testing.T) {
	assertValues(t, "5", []string{"1", "2", "3", "4", "5"})
	assertValues(t, "3-5", []string{"3", "4", "5"})
}

func TestBinary(t *testing.T) {
	assertValues(t, "bin 4", []string{"1", "10", "11", "100"})
	assertValues(t, "bin 3-4", []string{"11", "100"})
}

func TestOctal(t *testing.T) {
	assertValues(t, "oct 8", []string{"1", "2", "3", "4", "5", "6", "7", "10"})
}

func TestHexLower(t *testing.T) {
	assertValues(t, "hex 3", []string{"1", "2", "3"})
	_, _, arg, _ := firstItem(t, "hex 16")
	if !strings.Contains(arg, "a") || !strings.Contains(arg, "f") {
		t.Errorf("Dispatch(%q) arg = %q, want to contain both %q and %q", "hex 16", arg, "a", "f")
	}
}

func TestHexUpper(t *testing.T) {
	_, _, arg, _ := firstItem(t, "Hex 16")
	if !strings.Contains(arg, "A") || !strings.Contains(arg, "F") {
		t.Errorf("Dispatch(%q) arg = %q, want to contain both %q and %q", "Hex 16", arg, "A", "F")
	}
}

func TestAlphaLower(t *testing.T) {
	assertValues(t, "alf 3", []string{"a", "b", "c"})
	got := values(t, itemArg(t, "alf 27"))
	if got[len(got)-1] != "aa" {
		t.Errorf("Dispatch(%q) last value = %q, want %q", "alf 27", got[len(got)-1], "aa")
	}
}

// TestAlphaLowerBeyondZZ regression-tests the m==0 boundary in
// internal/seq.convertAlphabet (703 == "aza"): the naive Go port of
// Python's ascii_lowercase[m-1] panicked here because Python's negative
// indexing silently wraps m==0 to the last letter.
func TestAlphaLowerBeyondZZ(t *testing.T) {
	title, _, arg, valid := firstItem(t, "alf 703")
	if !valid {
		t.Fatalf("Dispatch(%q): item is not valid (title=%q)", "alf 703", title)
	}
	got := values(t, arg)
	if got[len(got)-1] != "aza" {
		t.Errorf("Dispatch(%q) last value = %q, want %q", "alf 703", got[len(got)-1], "aza")
	}
}

func TestAlphaUpper(t *testing.T) {
	assertValues(t, "Alf 3", []string{"A", "B", "C"})
	got := values(t, itemArg(t, "Alf 27"))
	if got[len(got)-1] != "AA" {
		t.Errorf("Dispatch(%q) last value = %q, want %q", "Alf 27", got[len(got)-1], "AA")
	}
}

func itemArg(t *testing.T, query string) string {
	t.Helper()
	_, _, arg, _ := firstItem(t, query)
	return arg
}

func TestFormat(t *testing.T) {
	assertValues(t, "fmt 3", []string{"item-1", "item-2", "item-3"})
	assertValues(t, "fmt 3 item-%d", []string{"item-1", "item-2", "item-3"})
	assertValues(t, "fmt 3 %a", []string{"a", "b", "c"})
}

func TestFormatMissingRangeShowsPlaceholder(t *testing.T) {
	title, subtitle, _, valid := firstItem(t, "fmt")
	if title != "Enter a length or range" || subtitle != "Paste Sequential Numbers (custom format)" || valid {
		t.Errorf("Dispatch(%q) = (%q, %q, valid=%v), want placeholder", "fmt", title, subtitle, valid)
	}
}

func TestEmptyQueryUsesDefaultLength(t *testing.T) {
	assertValues(t, "", []string{"1", "2", "3", "4", "5", "6", "7", "8", "9", "10", "11", "12"})
}

func TestUnrecognizedFirstTokenFallsBackToDecimal(t *testing.T) {
	// A bare range like "3-5" isn't a registered subcommand, so the whole
	// query is treated as narr with the default decimal format.
	assertValues(t, "3-5", []string{"3", "4", "5"})
}

func TestMultipleRangeArgsRepeatOuterDimension(t *testing.T) {
	// "seq bin 3 5": two ranges but the %b pattern has one field, so the
	// outer range's value repeats once per inner range value — matches
	// the original Python implementation's str.format() behavior.
	assertValues(t, "bin 3 5", []string{
		"1", "1", "1", "1", "1",
		"10", "10", "10", "10", "10",
		"11", "11", "11", "11", "11",
	})
}

func TestDoubleSpacePreviewsAllFormats(t *testing.T) {
	resp := Dispatch(" 3")
	if len(resp.Items) != len(allFormats) {
		t.Fatalf("Dispatch(%q): got %d items, want %d", " 3", len(resp.Items), len(allFormats))
	}
	wantSubtitles := []string{
		"Paste Sequential Numbers (decimal)",
		"bin: Paste Sequential Numbers (binary)",
		"oct: Paste Sequential Numbers (octal)",
		"hex: Paste Sequential Numbers (hex lower)",
		"Hex: Paste Sequential Numbers (hex upper)",
		"alf: Paste Sequential Numbers (alpha lower)",
		"Alf: Paste Sequential Numbers (alpha upper)",
		"fmt: Paste Sequential Numbers (custom format)",
	}
	for i, want := range wantSubtitles {
		if resp.Items[i].Subtitle != want {
			t.Errorf("Dispatch(%q).Items[%d].Subtitle = %q, want %q", " 3", i, resp.Items[i].Subtitle, want)
		}
	}
}

func TestDoubleSpaceEmptyUsesDefaultLength(t *testing.T) {
	resp := Dispatch(" ")
	if len(resp.Items) == 0 {
		t.Fatal("Dispatch(\" \"): got no items")
	}
	got := values(t, resp.Items[0].Arg)
	if len(got) != 12 {
		t.Errorf("Dispatch(\" \").Items[0] has %d values, want 12 (default length)", len(got))
	}
}
