package tokenization

import "fmt"

// TokenType represents the type of a token
type TokenType int

const (
	SpaceTab TokenType = iota + 1
	NewLine
	Word
	WordAlphanumerical
	WordWithSymbol
	Abbreviation
	AbbreviationWithDots
	Punctuation
	RomanNumeral
	Number
	PercentNumeral
	Time
	Date
	URL
	Email
	HashTag
	Mention
	MetaTag
	Emoji
	Emoticon
	UnknownWord
	Unknown
)

// Token represents a lexical token
type Token struct {
	Content    string
	Type       TokenType
	Start      int
	End        int
	Normalized string
}

// NewToken creates a new Token
func NewToken(content string, tokenType TokenType, start, end int, normalized ...string) *Token {
	norm := content
	if len(normalized) > 0 {
		norm = normalized[0]
	}
	return &Token{
		Content:    content,
		Type:       tokenType,
		Start:      start,
		End:        end,
		Normalized: norm,
	}
}

// String returns string representation of the token
func (t *Token) String() string {
	return fmt.Sprintf("[%s %d %d-%d]", t.Content, t.Type, t.Start, t.End)
}

// SimpleTokenize performs basic tokenization by spaces and punctuation
func SimpleTokenize(text string) []string {
	tokens := make([]string, 0)
	current := ""

	for _, r := range text {
		if r == ' ' || r == '\t' || r == '\n' {
			if current != "" {
				tokens = append(tokens, current)
				current = ""
			}
		} else if r == '.' || r == ',' || r == '!' || r == '?' || r == ':' || r == ';' {
			if current != "" {
				tokens = append(tokens, current)
				current = ""
			}
			tokens = append(tokens, string(r))
		} else {
			current += string(r)
		}
	}

	if current != "" {
		tokens = append(tokens, current)
	}

	return tokens
}

// TokenTypeName returns the name of a token type
func TokenTypeName(t TokenType) string {
	switch t {
	case SpaceTab:
		return "SpaceTab"
	case NewLine:
		return "NewLine"
	case Word:
		return "Word"
	case WordAlphanumerical:
		return "WordAlphanumerical"
	case WordWithSymbol:
		return "WordWithSymbol"
	case Abbreviation:
		return "Abbreviation"
	case AbbreviationWithDots:
		return "AbbreviationWithDots"
	case Punctuation:
		return "Punctuation"
	case RomanNumeral:
		return "RomanNumeral"
	case Number:
		return "Number"
	case PercentNumeral:
		return "PercentNumeral"
	case Time:
		return "Time"
	case Date:
		return "Date"
	case URL:
		return "URL"
	case Email:
		return "Email"
	case HashTag:
		return "HashTag"
	case Mention:
		return "Mention"
	case MetaTag:
		return "MetaTag"
	case Emoji:
		return "Emoji"
	case Emoticon:
		return "Emoticon"
	case UnknownWord:
		return "UnknownWord"
	case Unknown:
		return "Unknown"
	default:
		return fmt.Sprintf("Unknown(%d)", t)
	}
}
