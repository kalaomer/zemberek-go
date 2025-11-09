package tokenization

import (
	"regexp"

	"github.com/kalaomer/zemberek-go/core/turkish"
)

// Pattern fragments (Turkish characters and common patterns)
// Using TurkishAlphabet.Instance to avoid character duplication
var (
	turkishLetters        = turkish.Instance.Lowercase        // "abcçdefgğhıijklmnoöprsştuüvyzxwqâîû"
	turkishLettersCapital = turkish.Instance.Uppercase        // "ABCÇDEFGĞHIIJKLMNOÖPRSŞTUÜVYZXWQÂÎÛ"
	turkishLettersAll     = turkish.Instance.AllLetters       // lowercase + uppercase
	allTurkishAlphanumeric = turkish.Instance.AllLettersAndDigits // "0123456789" + AllLetters
	allTurkishAlphanumericUnderscore = turkish.Instance.AllLettersDigitsUnderscore // "0123456789" + AllLetters + "_"

	// Apostrophe: U+0027 (') or U+2019 (')
	apostrophe = `['']`

	// Turkish suffix pattern: 'e, 'a, 'den, 'te, etc.
	aposAndSuffix = apostrophe + `[` + turkishLettersAll + `]+`
)

// Compiled regex patterns for token detection
// Pattern priority order matters: URL before Email, Time before Number, etc.
var (
	// Time: 10:20, 10:20:53, 10.20.00'da
	timePattern = regexp.MustCompile(`^[0-2][0-9][:\.][0-5][0-9]([:\.][0-5][0-9])?(` + aposAndSuffix + `)?`)

	// Date: 1/1/2011, 02.12.1998'de, 1.1.11, 02/12/1998
	datePattern = regexp.MustCompile(`^[0-3]?[0-9][\./][0-1]?[0-9][\./]([1][7-9][0-9][0-9]|[2][0][0-9][0-9]|[0-9][0-9])(` + aposAndSuffix + `)?`)

	// Percent: %2.5, %100'e
	percentPattern = regexp.MustCompile(`^%[+\-]?[0-9]+([.,][0-9]+)?(` + aposAndSuffix + `)?`)

	// Number patterns (order matters - most specific first):
	// 1. Float with exponent: 1.35E-9, 1e10'dur
	numberExpPattern = regexp.MustCompile(`^[+\-]?[0-9]+([.,][0-9]+)?[Ee][+\-]?[0-9]+(` + aposAndSuffix + `)?`)

	// 2. Fraction: 1/2, -3/4, 123/456
	numberFractionPattern = regexp.MustCompile(`^[+\-]?[0-9]+/[0-9]+(` + aposAndSuffix + `)?`)

	// 3. Thousand separator (dot): 1.000.000, 2.500
	numberThousandDotPattern = regexp.MustCompile(`^([0-9]+\.)+[0-9]+(` + aposAndSuffix + `)?`)

	// 4. Thousand separator (comma): 2,345,531
	numberThousandCommaPattern = regexp.MustCompile(`^([0-9]+,)+[0-9]+(` + aposAndSuffix + `)?`)

	// 5. Decimal: -1.35, 3,1'e
	numberDecimalPattern = regexp.MustCompile(`^[+\-]?[0-9]+[.,][0-9]+(` + aposAndSuffix + `)?`)

	// 6. Ordinal: 2., 34.'ncü
	numberOrdinalPattern = regexp.MustCompile(`^[0-9]+\.(` + aposAndSuffix + `)?`)

	// 7. Integer: -3, 45, 100'e
	numberIntegerPattern = regexp.MustCompile(`^[+\-]?[0-9]+(` + aposAndSuffix + `)?`)

	// URL: http://foo.bar, www.foo.bar, foo.com'da, foo.com.tr/path
	urlPattern = regexp.MustCompile(`^(https?://|www\.)[` + allTurkishAlphanumeric + `\-_/?&+;=\[\].:]+(` + aposAndSuffix + `)?|^[` + allTurkishAlphanumericUnderscore + `]+\.(com|org|edu|gov|net|info)(\.tr)?(\/[` + allTurkishAlphanumeric + `\-_/?&+;=\[\].]+)?(` + aposAndSuffix + `)?`)

	// Email: ali@gmail.com, foo.bar@domain.com.tr'ye
	emailPattern = regexp.MustCompile(`^[` + allTurkishAlphanumericUnderscore + `]+\.?[` + allTurkishAlphanumericUnderscore + `]+@([` + allTurkishAlphanumericUnderscore + `]+\.)+[` + allTurkishAlphanumericUnderscore + `]+(` + aposAndSuffix + `)?`)

	// Mention: @kemal, @user_name'in
	mentionPattern = regexp.MustCompile(`^@[` + allTurkishAlphanumericUnderscore + `]+(` + aposAndSuffix + `)?`)

	// HashTag: #tag, #turkish_tag'a
	hashTagPattern = regexp.MustCompile(`^#[` + allTurkishAlphanumericUnderscore + `]+(` + aposAndSuffix + `)?`)

	// MetaTag: <tag>, <meta>
	metaTagPattern = regexp.MustCompile(`^<[` + allTurkishAlphanumericUnderscore + `]+>`)

	// Emoticons (from ANTLR grammar)
	emoticonPattern = regexp.MustCompile(`^(:\)|:-\)|:-\]|:D|:-D|8-\)|;\)|;‑\)|:\(|:-\(|:'\(|:'\)|:P|:p|:\||=\||=\)|=\(|:‑/|:/|:\^\)|¯\\_\(ツ\)_/¯|O_o|o_O|O_O|\\o/|<3)`)

	// Roman Numeral: IX, IV'ü, XII. (only valid Roman numeral characters)
	// Must contain ONLY I, V, X, L, C, D, M - not other letters
	romanNumeralPattern = regexp.MustCompile(`^[IVXLCDM]+\.?(` + aposAndSuffix + `)?$`)

	// Abbreviation with dots: I.B.M., T.C.K.'yı, İ.Ö (but not single letter like E.)
	// Matches either: (1) 2+ letter+dot pairs with optional final letter (T.C., I.B.M.), OR (2) exactly 1 letter+dot with REQUIRED final letter (İ.Ö)
	// Order matters: try longer pattern first
	abbreviationWithDotsPattern = regexp.MustCompile(`^([` + turkishLettersCapital + `]\.){2,}[` + turkishLettersCapital + `]?(` + aposAndSuffix + `)?|^[` + turkishLettersCapital + `]\.[` + turkishLettersCapital + `](` + aposAndSuffix + `)?`)

	// Word with symbol: F-16'yı, H1N1-A, covid-19
	wordWithSymbolPattern = regexp.MustCompile(`^[` + allTurkishAlphanumeric + `]+-[` + allTurkishAlphanumeric + `]+(` + aposAndSuffix + `)?`)

	// Alphanumeric word: F16, H1N1, covid19
	wordAlphanumericalPattern = regexp.MustCompile(`^[` + allTurkishAlphanumeric + `]+(` + aposAndSuffix + `)?`)

	// Pure Turkish word: merhaba, İstanbul, Ahmet'e
	wordPattern = regexp.MustCompile(`^[` + turkishLettersAll + `]+(` + aposAndSuffix + `)?`)

	// Punctuation: ., ,, !, ?, :, ;, ..., (!), (?)
	punctuationPattern = regexp.MustCompile(`^(\.\.\.|\(!\)|\(\?\)|[.,!?%$&*+@:;®™©℠>…=\\/\(\)\[\]\{\}\^'\"\-])`)

	// Whitespace patterns
	spaceTabPattern = regexp.MustCompile(`^[ \t]+`)
	newLinePattern  = regexp.MustCompile(`^[\r\n]+`)

	// Unknown word: characters not matching any pattern
	unknownWordPattern = regexp.MustCompile(`^[^ \n\r\t.,!?%$&*+@:;…®™©℠=>'"'"»«\\\-/\(\)\[\]\{\}\^]+`)
)
