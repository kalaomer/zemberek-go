package tokenization

import (
	"strings"
	"testing"
)

// Helper function to match sentence tokenization
// Compares tokenized output with expected string (space-separated tokens)
func matchSentence(t *testing.T, tokenizer *TurkishTokenizer, input, expected string) {
	tokens := tokenizer.TokenizeToStrings(input)
	got := strings.Join(tokens, " ")
	if got != expected {
		t.Errorf("Tokenize(%q)\n  got:      %q\n  expected: %q", input, got, expected)
	}
}

// Helper function to match single token
func matchToken(t *testing.T, tokenizer *TurkishTokenizer, input, expected string) {
	tokens := tokenizer.TokenizeToStrings(input)
	if len(tokens) != 1 {
		t.Errorf("Tokenize(%q) returned %d tokens, expected 1: %v", input, len(tokens), tokens)
		return
	}
	if tokens[0] != expected {
		t.Errorf("Tokenize(%q) = %q, expected %q", input, tokens[0], expected)
	}
}

// TestAbbreviationsInSentence tests abbreviation detection in context
// Ported from Java: TurkishTokenizerTest.testAbbreviations()
func TestAbbreviationsInSentence(t *testing.T) {
	tokenizer := DEFAULT

	// Test that abbreviations are kept together
	matchToken(t, tokenizer, "Prof.", "Prof.")
	matchToken(t, tokenizer, "Dr.", "Dr.")
	matchToken(t, tokenizer, "yy.", "yy.")

	// Test that normal words with dots are split
	matchSentence(t, tokenizer, "kedi.", "kedi .")

	// Complex sentence with multiple abbreviations
	matchSentence(t, tokenizer,
		"Prof. Dr. Ahmet'e git! dedi Av. Mehmet.",
		"Prof. Dr. Ahmet'e git ! dedi Av. Mehmet .")
}

// TestAbbreviations2 tests more abbreviation cases
// Ported from Java: TurkishTokenizerTest.testAbbreviations2()
func TestAbbreviations2(t *testing.T) {
	tokenizer := DEFAULT

	// Capital letters with dots
	matchToken(t, tokenizer, "I.B.M.", "I.B.M.")
	matchToken(t, tokenizer, "I.B.M.'nin", "I.B.M.'nin")

	// Mixed case
	matchSentence(t, tokenizer, "İ.Ö,Ğ.Ş", "İ.Ö , Ğ.Ş")
	matchSentence(t, tokenizer, "İ.Ö.,Ğ.Ş.", "İ.Ö. , Ğ.Ş.")
}

// TestCapitalWords tests capital word handling
// Ported from Java: TurkishTokenizerTest.testCapitalWords()
func TestCapitalWords(t *testing.T) {
	tokenizer := DEFAULT

	matchToken(t, tokenizer, "TCDD", "TCDD")
	matchToken(t, tokenizer, "TCDD'ye", "TCDD'ye")
}

// TestAlphaNumerical tests alphanumeric word detection
// Ported from Java: TurkishTokenizerTest.testAlphaNumerical()
func TestAlphaNumerical(t *testing.T) {
	tokenizer := DEFAULT

	matchSentence(t, tokenizer,
		"F-16'yı, (H1N1) H1N1'den.",
		"F-16'yı , ( H1N1 ) H1N1'den .")
}

// TestURLInSentence tests URL tokenization in sentences
// Ported from Java: TurkishTokenizerTest.testUrl()
func TestURLInSentence(t *testing.T) {
	tokenizer := DEFAULT

	urls := []string{
		"http://t.co/gn32szS9",
		"http://www.fo.bar",
		"www.fo.bar",
		"fo.com.tr",
		"foo.net'e",
	}

	for _, url := range urls {
		tokens := tokenizer.TokenizeToStrings(url)
		if len(tokens) != 1 {
			t.Errorf("URL %q was split into %d tokens: %v", url, len(tokens), tokens)
		}
	}
}

// TestEmailInSentence tests email tokenization
// Ported from Java: TurkishTokenizerTest.testEmail()
func TestEmailInSentence(t *testing.T) {
	tokenizer := DEFAULT

	emails := []string{
		"fo@bar.baz",
		"fo.bar@bar.baz",
		"fo_.bar@bar.baz",
		"ali@gmail.com'u",
	}

	for _, email := range emails {
		tokens := tokenizer.TokenizeToStrings(email)
		if len(tokens) != 1 {
			t.Errorf("Email %q was split into %d tokens: %v", email, len(tokens), tokens)
		}
	}
}

// TestMentionAndHashTag tests social media tokens
// Ported from Java: TurkishTokenizerTest
func TestMentionAndHashTag(t *testing.T) {
	tokenizer := DEFAULT

	mentions := []string{"@bar", "@foo_bar", "@kemal'in"}
	for _, mention := range mentions {
		matchToken(t, tokenizer, mention, mention)
	}

	hashTags := []string{"#foo", "#foo_bar", "#foo_bar'a"}
	for _, tag := range hashTags {
		matchToken(t, tokenizer, tag, tag)
	}
}

// TestTokenBoundaries tests that start/end positions are correct
// Ported from Java: TurkishTokenizerTest.testTokenBoundaries()
func TestTokenBoundaries(t *testing.T) {
	tokenizer := DEFAULT

	input := "bir av. geldi."
	tokens := tokenizer.Tokenize(input)

	if len(tokens) != 4 {
		t.Fatalf("Expected 4 tokens, got %d: %v", len(tokens), tokens)
	}

	// Check positions (rune-based, not byte-based)
	expected := []struct {
		content string
		start   int
		end     int
	}{
		{"bir", 0, 2},
		{"av.", 4, 6},
		{"geldi", 8, 12},
		{".", 13, 13},
	}

	for i, exp := range expected {
		if tokens[i].Content != exp.content {
			t.Errorf("Token[%d].Content = %q, expected %q", i, tokens[i].Content, exp.content)
		}
		if tokens[i].Start != exp.start {
			t.Errorf("Token[%d].Start = %d, expected %d", i, tokens[i].Start, exp.start)
		}
		if tokens[i].End != exp.end {
			t.Errorf("Token[%d].End = %d, expected %d", i, tokens[i].End, exp.end)
		}
	}
}

// TestApostrophes tests different apostrophe types
// Ported from Java: TurkishTokenizerTest.testApostrophes()
func TestApostrophes(t *testing.T) {
	tokenizer := DEFAULT

	// U+0027
	matchToken(t, tokenizer, "foo'f", "foo'f")
	// U+2019
	matchToken(t, tokenizer, "foo'f", "foo'f")

	// Apostrophes at start/end
	matchSentence(t, tokenizer, "'foo", "' foo")
	matchSentence(t, tokenizer, "''foo'", "' ' foo '")
}

// TestUnknownWords tests handling of non-Turkish characters
// Ported from Java: TurkishTokenizerTest.testUnknownWord*()
func TestUnknownWords(t *testing.T) {
	tokenizer := DEFAULT

	// Arabic characters
	matchSentence(t, tokenizer, "زنبورك", "زنبورك")

	// Norwegian characters
	matchSentence(t, tokenizer, "Bjørn", "Bjørn")
}

// TestDotsInMiddle tests dot handling
// Ported from Java: TurkishTokenizerTest
func TestDotsInMiddle(t *testing.T) {
	tokenizer := DEFAULT

	matchSentence(t, tokenizer, "Ali.gel.", "Ali . gel .")
}

// TestUnderscores tests underscore handling
// Ported from Java: TurkishTokenizerTest
func TestUnderscores(t *testing.T) {
	tokenizer := DEFAULT

	matchSentence(t, tokenizer, "__he_llo__", "__he_llo__")
}

// TestCapitalLettersAfterQuotes tests Issue #64 fix
// Ported from Java: TurkishTokenizerTest.testCapitalLettersAfterQuotes()
func TestCapitalLettersAfterQuotes(t *testing.T) {
	tokenizer := DEFAULT

	matchSentence(t, tokenizer, "ANKARA'ya.", "ANKARA'ya .")
	matchSentence(t, tokenizer, "ANKARA'YA.", "ANKARA'YA .")
	matchSentence(t, tokenizer, "Ankara'YA.", "Ankara'YA .")
}

// TestPunctuation tests punctuation tokenization
// Ported from Java: TurkishTokenizerTest.testPunctuation()
func TestPunctuationInSentence(t *testing.T) {
	tokenizer := DEFAULT

	matchSentence(t, tokenizer,
		".,!:;$%\"'()[]{}&@®™©℠",
		". , ! : ; $ % \" ' ( ) [ ] { } & @ ® ™ © ℠")

	// Special multi-char punctuation
	matchToken(t, tokenizer, "...", "...")
	matchToken(t, tokenizer, "(!)", "(!)")
}

// TestBuilderPattern tests custom tokenizer building
func TestBuilderPattern(t *testing.T) {
	// Tokenizer that ignores punctuation and whitespace
	tokenizer := NewBuilder().
		AcceptAll().
		IgnoreTypes(Punctuation, NewLine, SpaceTab).
		Build()

	input := "Merhaba, dünya!"
	tokens := tokenizer.TokenizeToStrings(input)

	if len(tokens) != 2 {
		t.Fatalf("Expected 2 tokens, got %d: %v", len(tokens), tokens)
	}

	if tokens[0] != "Merhaba" || tokens[1] != "dünya" {
		t.Errorf("Got tokens %v, expected [Merhaba dünya]", tokens)
	}
}

// TestPredefinedInstances tests ALL and DEFAULT tokenizers
func TestPredefinedInstances(t *testing.T) {
	input := "Merhaba dünya"

	// DEFAULT ignores whitespace
	defaultTokens := DEFAULT.TokenizeToStrings(input)
	if len(defaultTokens) != 2 {
		t.Errorf("DEFAULT should ignore whitespace, got %d tokens: %v", len(defaultTokens), defaultTokens)
	}

	// ALL includes whitespace
	allTokens := ALL.TokenizeToStrings(input)
	if len(allTokens) != 3 {
		t.Errorf("ALL should include whitespace, got %d tokens: %v", len(allTokens), allTokens)
	}
}

// TestComplexSentence tests a real-world complex sentence
func TestComplexSentence(t *testing.T) {
	tokenizer := DEFAULT

	input := "Prof. Dr. Ahmet Yılmaz, @twitter'da #türkçe hakkında www.example.com'u paylaştı."

	tokens := tokenizer.Tokenize(input)

	// Check that we got reasonable tokenization
	if len(tokens) < 10 {
		t.Errorf("Expected at least 10 tokens for complex sentence, got %d", len(tokens))
	}

	// Verify some specific tokens
	foundProf := false
	foundDr := false
	foundMention := false
	foundHashTag := false
	foundURL := false

	for _, token := range tokens {
		switch token.Content {
		case "Prof.":
			foundProf = true
			if token.Type != Abbreviation {
				t.Errorf("'Prof.' should be Abbreviation, got %v", TokenTypeName(token.Type))
			}
		case "Dr.":
			foundDr = true
			if token.Type != Abbreviation {
				t.Errorf("'Dr.' should be Abbreviation, got %v", TokenTypeName(token.Type))
			}
		}

		if strings.HasPrefix(token.Content, "@") {
			foundMention = true
			if token.Type != Mention {
				t.Errorf("Mention should have type Mention, got %v", TokenTypeName(token.Type))
			}
		}

		if strings.HasPrefix(token.Content, "#") {
			foundHashTag = true
			if token.Type != HashTag {
				t.Errorf("HashTag should have type HashTag, got %v", TokenTypeName(token.Type))
			}
		}

		if strings.Contains(token.Content, "www.") || strings.HasPrefix(token.Content, "http") {
			foundURL = true
			if token.Type != URL {
				t.Errorf("URL should have type URL, got %v", TokenTypeName(token.Type))
			}
		}
	}

	if !foundProf {
		t.Error("Should find 'Prof.' abbreviation")
	}
	if !foundDr {
		t.Error("Should find 'Dr.' abbreviation")
	}
	if !foundMention {
		t.Error("Should find @mention")
	}
	if !foundHashTag {
		t.Error("Should find #hashtag")
	}
	if !foundURL {
		t.Error("Should find URL")
	}
}

// TestLegalDocument tests tokenization of Turkish legal text
func TestLegalDocument(t *testing.T) {
	tokenizer := DEFAULT

	input := "T.C. Mahkemesi E. 123/456 sayılı kararında 15.08.2023 tarihinde karar verdi."

	tokens := tokenizer.Tokenize(input)

	// Check for important legal tokens
	foundTC := false
	foundE := false
	foundCaseNumber := false
	foundDate := false

	for _, token := range tokens {
		if token.Content == "T.C." {
			foundTC = true
			if token.Type != AbbreviationWithDots {
				t.Errorf("'T.C.' should be AbbreviationWithDots, got %v", TokenTypeName(token.Type))
			}
		}

		if token.Content == "E." {
			foundE = true
			if token.Type != Abbreviation {
				t.Errorf("'E.' should be Abbreviation, got %v", TokenTypeName(token.Type))
			}
		}

		if strings.Contains(token.Content, "/") && strings.Contains(token.Content, "123") {
			foundCaseNumber = true
			if token.Type != Number {
				t.Errorf("'123/456' should be Number (fraction), got %v", TokenTypeName(token.Type))
			}
		}

		if strings.Contains(token.Content, "15.08") {
			foundDate = true
			if token.Type != Date {
				t.Errorf("'15.08.2023' should be Date, got %v", TokenTypeName(token.Type))
			}
		}
	}

	if !foundTC {
		t.Error("Should find 'T.C.' abbreviation")
	}
	if !foundE {
		t.Error("Should find 'E.' abbreviation")
	}
	if !foundCaseNumber {
		t.Error("Should find case number '123/456'")
	}
	if !foundDate {
		t.Error("Should find date '15.08.2023'")
	}
}

// TestAbbreviationsLoading tests that abbreviations were loaded
func TestAbbreviationsLoading(t *testing.T) {
	// Check that some common abbreviations are loaded
	commonAbbr := []string{"Prof.", "Dr.", "Av.", "T.C.", "vb.", "vs."}

	for _, abbr := range commonAbbr {
		if !IsAbbreviation(abbr) {
			t.Errorf("Abbreviation %q should be loaded", abbr)
		}

		// Check lowercase version too
		if !IsAbbreviation(strings.ToLower(abbr)) {
			t.Errorf("Lowercase abbreviation %q should be loaded", strings.ToLower(abbr))
		}
	}

	// Check that non-abbreviations are not matched
	if IsAbbreviation("kedi.") {
		t.Error("'kedi.' should not be an abbreviation")
	}
}

// TestOrdinalNumbers tests that ordinal numbers don't include trailing whitespace
func TestOrdinalNumbers(t *testing.T) {
	tokenizer := DEFAULT

	tests := []struct {
		input         string
		expectedToken string
		description   string
	}{
		{"16. mahkeme", "16.", "ordinal number before word"},
		{"1. sıra", "1.", "single digit ordinal"},
		{"34. madde", "34.", "two digit ordinal"},
		{"123.", "123.", "standalone ordinal number"},
	}

	for _, tt := range tests {
		tokens := tokenizer.Tokenize(tt.input)
		if len(tokens) == 0 {
			t.Errorf("Test %q: no tokens returned", tt.description)
			continue
		}

		firstToken := tokens[0]

		// Check content matches expected (no trailing space)
		if firstToken.Content != tt.expectedToken {
			t.Errorf("Test %q: got content %q, expected %q",
				tt.description, firstToken.Content, tt.expectedToken)
		}

		// Check it's a Number type
		if firstToken.Type != Number {
			t.Errorf("Test %q: got type %v, expected Number",
				tt.description, TokenTypeName(firstToken.Type))
		}

		// Ensure no trailing whitespace in content
		if strings.HasSuffix(firstToken.Content, " ") {
			t.Errorf("Test %q: ordinal number %q has trailing whitespace",
				tt.description, firstToken.Content)
		}
	}
}
