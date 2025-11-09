package tokenization

import (
	"testing"
)

// TestNumbers tests number detection (ported from Java TurkishTokenizerTest)
func TestNumbers(t *testing.T) {
	tests := []struct {
		input    string
		expected TokenType
	}{
		// Integer
		{"1", Number},
		{"123", Number},
		{"-3", Number},
		{"45", Number},
		{"+100", Number},

		// Decimal
		{"3.14", Number},
		{"-1.34", Number},
		{"-3,14", Number},
		{"1,35", Number},

		// With Turkish suffix
		{"100'e", Number},
		{"3.14'ten", Number},
		{"45'in", Number},

		// Scientific notation
		{"1e10", Number},
		{"-3e4", Number},
		{"1.35E-9", Number},
		{"1e10'dur", Number},

		// Fraction
		{"1/2", Number},
		{"-3/4", Number},
		{"123/456", Number},

		// Thousand separator
		{"1.000.000", Number},
		{"2.345.531", Number},
		{"1,000,000", Number},

		// Ordinal
		{"2.", Number},
		{"34.", Number},

		// Percent
		{"%2.5", PercentNumeral},
		{"%100", PercentNumeral},
		{"%2.5'ten", PercentNumeral},
		{"%100'e", PercentNumeral},
	}

	for _, tt := range tests {
		got := DetermineTokenType(tt.input)
		if got != tt.expected {
			t.Errorf("DetermineTokenType(%q) = %v, want %v", tt.input, TokenTypeName(got), TokenTypeName(tt.expected))
		}
	}
}

// TestTime tests time detection
func TestTime(t *testing.T) {
	tests := []string{
		"10:20",
		"10:20:53",
		"00:00",
		"23:59",
		"10.20",
		"10.20.00",
		"10:20'de",
		"10.20.00'da",
		"23:59'a",
	}

	for _, input := range tests {
		got := DetermineTokenType(input)
		if got != Time {
			t.Errorf("DetermineTokenType(%q) = %v, want Time", input, TokenTypeName(got))
		}
	}
}

// TestDate tests date detection
func TestDate(t *testing.T) {
	tests := []string{
		"1/1/2011",
		"02/12/1998",
		"31/12/99",
		"1.1.2011",
		"02.12.1998",
		"31.12.99",
		"02.12.1998'de",
		"1/1/2011'e",
		"15.08.2023'ten",
	}

	for _, input := range tests {
		got := DetermineTokenType(input)
		if got != Date {
			t.Errorf("DetermineTokenType(%q) = %v, want Date", input, TokenTypeName(got))
		}
	}
}

// TestURL tests URL detection
func TestURL(t *testing.T) {
	tests := []string{
		"http://t.co/gn32szS9",
		"https://www.google.com",
		"http://www.fo.bar",
		"www.fo.bar",
		"www.google.com",
		"foo.com",
		"bar.org",
		"baz.edu",
		"test.gov",
		"example.net",
		"info.info",
		"foo.com.tr",
		"www.fo.bar'da",
		"foo.net'e",
		"example.com/path",
		"https://www.google.com.tr/search?q=test",
	}

	for _, input := range tests {
		got := DetermineTokenType(input)
		if got != URL {
			t.Errorf("DetermineTokenType(%q) = %v, want URL", input, TokenTypeName(got))
		}
	}
}

// TestEmail tests email detection
func TestEmail(t *testing.T) {
	tests := []string{
		"fo@bar.baz",
		"foo@bar.baz",
		"fo.bar@bar.baz",
		"fo_bar@bar.baz",
		"ali@gmail.com",
		"test@example.org",
		"ali@gmail.com'u",
		"test@domain.com.tr",
		"user_name@example.co.uk",
	}

	for _, input := range tests {
		got := DetermineTokenType(input)
		if got != Email {
			t.Errorf("DetermineTokenType(%q) = %v, want Email", input, TokenTypeName(got))
		}
	}
}

// TestMention tests mention detection
func TestMention(t *testing.T) {
	tests := []string{
		"@bar",
		"@foo_bar",
		"@kemal",
		"@user123",
		"@kemal'in",
		"@foo_bar'a",
		"@test_user'dan",
	}

	for _, input := range tests {
		got := DetermineTokenType(input)
		if got != Mention {
			t.Errorf("DetermineTokenType(%q) = %v, want Mention", input, TokenTypeName(got))
		}
	}
}

// TestHashTag tests hashtag detection
func TestHashTag(t *testing.T) {
	tests := []string{
		"#foo",
		"#foo_bar",
		"#tag",
		"#türkçe",
		"#test123",
		"#tag'a",
		"#foo_bar'dan",
		"#hashtag'ı",
	}

	for _, input := range tests {
		got := DetermineTokenType(input)
		if got != HashTag {
			t.Errorf("DetermineTokenType(%q) = %v, want HashTag", input, TokenTypeName(got))
		}
	}
}

// TestMetaTag tests meta tag detection
func TestMetaTag(t *testing.T) {
	tests := []string{
		"<tag>",
		"<meta>",
		"<html>",
		"<div>",
	}

	for _, input := range tests {
		got := DetermineTokenType(input)
		if got != MetaTag {
			t.Errorf("DetermineTokenType(%q) = %v, want MetaTag", input, TokenTypeName(got))
		}
	}
}

// TestEmoticons tests emoticon detection
func TestEmoticons(t *testing.T) {
	tests := []string{
		":)",
		":-)",
		":-]",
		":D",
		":-D",
		"8-)",
		";)",
		";‑)",
		":(",
		":-(",
		":'(",
		":')",
		":P",
		":p",
		":|",
		"=|",
		"=)",
		"=(",
		":‑/",
		":/",
		":^)",
		"¯\\_(ツ)_/¯",
		"O_o",
		"o_O",
		"O_O",
		"\\o/",
		"<3",
	}

	for _, input := range tests {
		got := DetermineTokenType(input)
		if got != Emoticon {
			t.Errorf("DetermineTokenType(%q) = %v, want Emoticon", input, TokenTypeName(got))
		}
	}
}

// TestRomanNumerals tests Roman numeral detection
func TestRomanNumerals(t *testing.T) {
	tests := []string{
		"I",
		"II",
		"III",
		"IV",
		"V",
		"IX",
		"X",
		"XII",
		"XX",
		"L",
		"C",
		"D",
		"M",
		"MCMXC",
		"IV.",
		"IX'u",
		"XII'ye",
	}

	for _, input := range tests {
		got := DetermineTokenType(input)
		if got != RomanNumeral {
			t.Errorf("DetermineTokenType(%q) = %v, want RomanNumeral", input, TokenTypeName(got))
		}
	}
}

// TestAbbreviationWithDots tests abbreviation with dots detection
func TestAbbreviationWithDots(t *testing.T) {
	tests := []string{
		"I.B.M.",
		"T.C.K.",
		"A.B.C.",
		"T.C.",
		"A.Ş.",
		"I.B.M.'nin",
		"T.C.K.'yı",
	}

	for _, input := range tests {
		got := DetermineTokenType(input)
		if got != AbbreviationWithDots {
			t.Errorf("DetermineTokenType(%q) = %v, want AbbreviationWithDots", input, TokenTypeName(got))
		}
	}
}

// TestWordWithSymbol tests word with symbol detection
func TestWordWithSymbol(t *testing.T) {
	tests := []string{
		"F-16",
		"H1N1-A",
		"covid-19",
		"F-16'yı",
	}

	for _, input := range tests {
		got := DetermineTokenType(input)
		if got != WordWithSymbol {
			t.Errorf("DetermineTokenType(%q) = %v, want WordWithSymbol", input, TokenTypeName(got))
		}
	}
}

// TestWordAlphanumerical tests alphanumerical word detection
func TestWordAlphanumerical(t *testing.T) {
	tests := []string{
		"F16",
		"H1N1",
		"covid19",
		"test123",
	}

	for _, input := range tests {
		got := DetermineTokenType(input)
		if got != WordAlphanumerical {
			t.Errorf("DetermineTokenType(%q) = %v, want WordAlphanumerical", input, TokenTypeName(got))
		}
	}
}

// TestWords tests pure word detection
func TestWords(t *testing.T) {
	tests := []string{
		"merhaba",
		"İstanbul",
		"Ankara",
		"kedi",
		"köpek",
		"Ahmet",
		"Mehmet",
		"kitap",
		"Ahmet'e",
		"İstanbul'a",
		"kitapları",
		"TCDD",
		"ANKARA",
	}

	for _, input := range tests {
		got := DetermineTokenType(input)
		if got != Word && got != WordAlphanumerical {
			t.Errorf("DetermineTokenType(%q) = %v, want Word or WordAlphanumerical", input, TokenTypeName(got))
		}
	}
}

// TestPunctuation tests punctuation detection
func TestPunctuation(t *testing.T) {
	tests := []string{
		".",
		",",
		"!",
		"?",
		":",
		";",
		"...",
		"(!)",
		"(?)",
		"(",
		")",
		"[",
		"]",
		"{",
		"}",
		"-",
		"/",
		"'",
		"\"",
	}

	for _, input := range tests {
		got := DetermineTokenType(input)
		if got != Punctuation {
			t.Errorf("DetermineTokenType(%q) = %v, want Punctuation", input, TokenTypeName(got))
		}
	}
}

// TestWhitespace tests whitespace detection
func TestWhitespace(t *testing.T) {
	tests := []struct {
		input    string
		expected TokenType
	}{
		{" ", SpaceTab},
		{"  ", SpaceTab},
		{"\t", SpaceTab},
		{" \t ", SpaceTab},
		{"\n", NewLine},
		{"\r", NewLine},
		{"\r\n", NewLine},
	}

	for _, tt := range tests {
		got := DetermineTokenType(tt.input)
		if got != tt.expected {
			t.Errorf("DetermineTokenType(%q) = %v, want %v", tt.input, TokenTypeName(got), TokenTypeName(tt.expected))
		}
	}
}

// TestHelperMethods tests token helper methods
func TestHelperMethods(t *testing.T) {
	// IsNumeral
	numberToken := &Token{Type: Number}
	if !numberToken.IsNumeral() {
		t.Error("Number token should return true for IsNumeral()")
	}

	percentToken := &Token{Type: PercentNumeral}
	if !percentToken.IsNumeral() {
		t.Error("PercentNumeral token should return true for IsNumeral()")
	}

	romanToken := &Token{Type: RomanNumeral}
	if !romanToken.IsNumeral() {
		t.Error("RomanNumeral token should return true for IsNumeral()")
	}

	// IsWhiteSpace
	spaceToken := &Token{Type: SpaceTab}
	if !spaceToken.IsWhiteSpace() {
		t.Error("SpaceTab token should return true for IsWhiteSpace()")
	}

	newlineToken := &Token{Type: NewLine}
	if !newlineToken.IsWhiteSpace() {
		t.Error("NewLine token should return true for IsWhiteSpace()")
	}

	// IsWebRelated
	urlToken := &Token{Type: URL}
	if !urlToken.IsWebRelated() {
		t.Error("URL token should return true for IsWebRelated()")
	}

	emailToken := &Token{Type: Email}
	if !emailToken.IsWebRelated() {
		t.Error("Email token should return true for IsWebRelated()")
	}

	mentionToken := &Token{Type: Mention}
	if !mentionToken.IsWebRelated() {
		t.Error("Mention token should return true for IsWebRelated()")
	}

	hashtagToken := &Token{Type: HashTag}
	if !hashtagToken.IsWebRelated() {
		t.Error("HashTag token should return true for IsWebRelated()")
	}

	// IsEmoji
	emoticonToken := &Token{Type: Emoticon}
	if !emoticonToken.IsEmoji() {
		t.Error("Emoticon token should return true for IsEmoji()")
	}

	// IsWord
	wordToken := &Token{Type: Word}
	if !wordToken.IsWord() {
		t.Error("Word token should return true for IsWord()")
	}

	abbrToken := &Token{Type: Abbreviation}
	if !abbrToken.IsWord() {
		t.Error("Abbreviation token should return true for IsWord()")
	}

	// IsUnidentified
	unknownToken := &Token{Type: Unknown}
	if !unknownToken.IsUnidentified() {
		t.Error("Unknown token should return true for IsUnidentified()")
	}

	unknownWordToken := &Token{Type: UnknownWord}
	if !unknownWordToken.IsUnidentified() {
		t.Error("UnknownWord token should return true for IsUnidentified()")
	}
}

// TestApostropheNormalization tests apostrophe normalization
func TestApostropheNormalization(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"foo'bar", "foo'bar"},   // U+0027 unchanged
		{"foo'bar", "foo'bar"},   // U+2019 -> U+0027
		{"test's", "test's"},     // U+2019 -> U+0027
	}

	for _, tt := range tests {
		got := NormalizeApostrophe(tt.input)
		if got != tt.expected {
			t.Errorf("NormalizeApostrophe(%q) = %q, want %q", tt.input, got, tt.expected)
		}
	}
}

// TestEdgeCases tests edge cases
func TestEdgeCases(t *testing.T) {
	// Empty string
	if got := DetermineTokenType(""); got != Unknown {
		t.Errorf("DetermineTokenType(\"\") = %v, want Unknown", TokenTypeName(got))
	}

	// Single character
	if got := DetermineTokenType("a"); got != Word {
		t.Errorf("DetermineTokenType(\"a\") = %v, want Word", TokenTypeName(got))
	}

	// Arabic characters (UnknownWord)
	arabicInput := "زنبورك"
	got := DetermineTokenType(arabicInput)
	if got != UnknownWord && got != Word {
		t.Errorf("DetermineTokenType(%q) = %v, want UnknownWord or Word", arabicInput, TokenTypeName(got))
	}
}
