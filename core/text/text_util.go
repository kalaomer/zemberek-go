package text

import (
	"regexp"
)

var (
	apostropheRegex      = regexp.MustCompile(`[′´` + "`" + `'']`)
	quotesRegex          = regexp.MustCompile(`[""»«″]|''`)
	quotesApostropheReg  = regexp.MustCompile(`[′´` + "`" + `'']`)
	hyphensRegex         = regexp.MustCompile(`[–]`)
)

// NormalizeApostrophes converts different apostrophe symbols to a unified form
func NormalizeApostrophes(inp string) string {
	return apostropheRegex.ReplaceAllString(inp, "'")
}

// NormalizeQuotesHyphens converts different quote and hyphen symbols to unified forms
func NormalizeQuotesHyphens(inp string) string {
	inp = quotesRegex.ReplaceAllString(inp, "\"")
	inp = quotesApostropheReg.ReplaceAllString(inp, "'")
	inp = hyphensRegex.ReplaceAllString(inp, "-")
	return inp
}
