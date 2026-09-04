// Package seqcmd dispatches an Alfred seq query to internal/seq and builds
// the Script Filter response, porting the original workflow/filter.py.
//
// Commands:
//
//	seq <length or range>            — decimal (default)
//	seq bin <length or range>        — binary
//	seq oct <length or range>        — octal
//	seq hex <length or range>        — hex lower
//	seq Hex <length or range>        — hex upper
//	seq alf <length or range>        — alphabetic lower
//	seq Alf <length or range>        — alphabetic upper
//	seq fmt <length or range> [fmt]  — custom format (default fmt: item-#)
//	seq  <length or range>           — (leading space) preview every format
package seqcmd

import (
	"strings"

	"github.com/y-marui/alfred-sequential-number/internal/scriptfilter"
	"github.com/y-marui/alfred-sequential-number/internal/seq"
)

const (
	defaultLength = "12"
	defaultFormat = "item-#"
)

// formatEntry is one row of ALL_FORMATS in filter.py: a format string, its
// display description, and the subcommand keyword that selects it ("" for
// the default decimal format).
type formatEntry struct {
	format string
	desc   string
	label  string
}

var allFormats = []formatEntry{
	{"%d", "Paste Sequential Numbers (decimal)", ""},
	{"%b", "Paste Sequential Numbers (binary)", "bin"},
	{"%o", "Paste Sequential Numbers (octal)", "oct"},
	{"%x", "Paste Sequential Numbers (hex lower)", "hex"},
	{"%X", "Paste Sequential Numbers (hex upper)", "Hex"},
	{"%a", "Paste Sequential Numbers (alpha lower)", "alf"},
	{"%A", "Paste Sequential Numbers (alpha upper)", "Alf"},
	{defaultFormat, "Paste Sequential Numbers (custom format)", "fmt"},
}

// subcommand looks up a non-"fmt" subcommand keyword (bin/oct/hex/Hex/
// alf/Alf) in allFormats.
func subcommand(label string) (formatEntry, bool) {
	for _, e := range allFormats {
		if e.label != "" && e.label == label {
			return e, true
		}
	}
	return formatEntry{}, false
}

// makePreview renders a short preview of seq for an item's title: the
// full list joined by ", " when it has 5 or fewer entries, otherwise the
// first 3 entries plus the last, joined with "...".
func makePreview(values []string) string {
	if len(values) <= 5 {
		return strings.Join(values, ", ")
	}
	return strings.Join(values[:3], ", ") + ", ..., " + values[len(values)-1]
}

// makeItem generates values for format/narr and builds a result item, or
// returns ok=false if generation failed or produced nothing.
func makeItem(format string, narr []string, desc string) (item scriptfilter.Item, ok bool) {
	values, err := seq.Generate(format, narr)
	if err != nil || len(values) == 0 {
		return scriptfilter.Item{}, false
	}
	return scriptfilter.Item{
		Title:    makePreview(values),
		Subtitle: desc,
		Arg:      strings.Join(values, "\r\n"),
		Valid:    scriptfilter.BoolPtr(true),
	}, true
}

func noResults(desc string) scriptfilter.Response {
	return scriptfilter.Response{
		Items: []scriptfilter.Item{
			{Title: "No results", Subtitle: desc, Valid: scriptfilter.BoolPtr(false)},
		},
	}
}

// Dispatch parses the raw Alfred query and returns the Script Filter
// response, mirroring workflow/filter.py exactly.
func Dispatch(rawQuery string) scriptfilter.Response {
	if strings.HasPrefix(rawQuery, " ") {
		return dispatchAllFormats(rawQuery)
	}

	args := strings.Fields(rawQuery)

	var format, desc string
	var narr []string

	switch {
	case len(args) == 0:
		format, narr, desc = "%d", []string{defaultLength}, "Paste Sequential Numbers (decimal)"

	case args[0] == "fmt":
		rest := args[1:]
		if len(rest) == 0 {
			return scriptfilter.Response{
				Items: []scriptfilter.Item{
					{
						Title:    "Enter a length or range",
						Subtitle: "Paste Sequential Numbers (custom format)",
						Valid:    scriptfilter.BoolPtr(false),
					},
				},
			}
		}
		narr = []string{rest[0]}
		if len(rest) > 1 {
			format = rest[1]
		} else {
			format = defaultFormat
		}
		desc = "Paste Sequential Numbers (custom format)"

	default:
		if e, ok := subcommand(args[0]); ok {
			format = e.format
			if len(args) > 1 {
				narr = args[1:]
			} else {
				narr = []string{defaultLength}
			}
			desc = e.desc
		} else {
			format, narr, desc = "%d", args, "Paste Sequential Numbers (decimal)"
		}
	}

	values, err := seq.Generate(format, narr)
	if err != nil {
		values = nil
	}
	if len(values) == 0 {
		return noResults(desc)
	}

	return scriptfilter.Response{
		Items: []scriptfilter.Item{
			{
				Title:    makePreview(values),
				Subtitle: desc,
				Arg:      strings.Join(values, "\r\n"),
				Valid:    scriptfilter.BoolPtr(true),
			},
		},
	}
}

// dispatchAllFormats handles the "seq" + two-spaces trick: preview every
// registered format for one length/range.
func dispatchAllFormats(rawQuery string) scriptfilter.Response {
	length := strings.TrimSpace(rawQuery)
	if length == "" {
		length = defaultLength
	}
	narr := []string{length}

	var items []scriptfilter.Item
	for _, e := range allFormats {
		displayDesc := e.desc
		if e.label != "" {
			displayDesc = e.label + ": " + e.desc
		}
		if item, ok := makeItem(e.format, narr, displayDesc); ok {
			items = append(items, item)
		}
	}
	if len(items) == 0 {
		items = []scriptfilter.Item{{Title: "No results", Valid: scriptfilter.BoolPtr(false)}}
	}
	return scriptfilter.Response{Items: items}
}
