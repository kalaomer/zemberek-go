package morphology

import (
	"runtime"
	"strings"
	"sync"
	"unicode/utf8"

	analysispkg "github.com/kalaomer/zemberek-go/morphology/analysis"
	"github.com/kalaomer/zemberek-go/core/turkish"
	"github.com/kalaomer/zemberek-go/tokenization"
)

// Global thread-safe cache for stemming results
// Key: normalized (lowercase) word -> Value: stem
var stemCache sync.Map

// StemToken represents a stemmed token with its byte position in the original text
type StemToken struct {
	Stem      string                  // Stemmed form: "kitap"
	Original  string                  // Original word: "kitapları"
	Type      tokenization.TokenType  // Token type from TurkishTokenizer
	StartByte int                     // UTF-8 byte offset start
	EndByte   int                     // UTF-8 byte offset end
}

// stemmingJob holds information for a token to be stemmed
type stemmingJob struct {
	index     int
	token     *tokenization.Token
	startByte int
	endByte   int
	needsStem bool
}

// StemTextWithPositions extracts stems from text with byte positions.
//
// This is the MAIN function for FTS5 integration. It:
// 1. Tokenizes text using TurkishTokenizer (28 token types)
// 2. Filters out punctuation, whitespace, emoticons, unknown tokens
// 3. Preserves URL/Email/Number/Date/Time without stemming
// 4. Performs parallel morphological analysis with caching on words
// 5. Extracts stems using Item.Root (for correct Turkish voicing)
// 6. Returns list of stems with their byte positions
//
// Performance optimizations:
// - Global cache (sync.Map) for repeated words
// - Worker pool for parallel analysis (uses all CPU cores)
// - Thread-safe: morphology.Analyze() is concurrent-safe
//
// Token Type Handling:
// - Punctuation, Whitespace, Emoticon, Unknown → FILTERED (not included in results)
// - URL, Email, Number, Date, Time → PRESERVED (kept as-is, not stemmed)
// - Word, Abbreviation, WordAlphanumerical → STEMMED
// - Mention, HashTag → PREPROCESSED (strip prefix) then STEMMED
//
// Example:
//
//	morph := CreateWithDefaults()
//	text := "Kitapları www.google.com okuyorum."
//	tokens := StemTextWithPositions(text, morph)
//	// tokens[0] = {Stem: "kitap", Original: "Kitapları", Type: WordAlphanumerical, StartByte: 0, EndByte: 11}
//	// tokens[1] = {Stem: "www.google.com", Original: "www.google.com", Type: URL, StartByte: 12, EndByte: 26}
//	// tokens[2] = {Stem: "oku", Original: "okuyorum", Type: WordAlphanumerical, StartByte: 27, EndByte: 35}
//	// Note: "." at end is filtered out (Punctuation type)
func StemTextWithPositions(text string, morphology *TurkishMorphology) []StemToken {
	// Use TurkishTokenizer (ignores whitespace by default)
	tokenizer := tokenization.DEFAULT
	tokens := tokenizer.Tokenize(text)

	if len(tokens) == 0 {
		return []StemToken{}
	}

	// Convert rune positions to byte positions
	// Token.Start and Token.End are RUNE positions, we need BYTE positions
	bytePositions := calculateBytePositions(text, tokens)

	// Filter tokens and build job list
	jobs := make([]stemmingJob, 0, len(tokens))
	for i, token := range tokens {
		// Skip tokens that should be completely filtered out (Punctuation, etc.)
		if shouldFilterToken(token.Type) {
			continue
		}

		needsStem := !shouldSkipStemming(token.Type)

		jobs = append(jobs, stemmingJob{
			index:     i,
			token:     token,
			startByte: bytePositions[i].start,
			endByte:   bytePositions[i].end,
			needsStem: needsStem,
		})
	}

	// If morphology is nil or no tokens need stemming, return original tokens
	if morphology == nil {
		return tokensWithoutStemming(jobs)
	}

	// Result type with index for ordering
	type indexedResult struct {
		index int
		token StemToken
	}

	// Filter jobs that need stemming
	stemmingJobs := make([]stemmingJob, 0, len(jobs))
	for _, j := range jobs {
		if j.needsStem {
			stemmingJobs = append(stemmingJobs, j)
		}
	}

	if len(stemmingJobs) == 0 {
		return tokensWithoutStemming(jobs)
	}

	// Worker pool setup
	numWorkers := runtime.NumCPU()
	if numWorkers > len(stemmingJobs) {
		numWorkers = len(stemmingJobs)
	}

	jobChan := make(chan stemmingJob, len(stemmingJobs))
	resultChan := make(chan indexedResult, len(jobs))

	// Start workers
	var wg sync.WaitGroup
	for w := 0; w < numWorkers; w++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := range jobChan {
				// Preprocess token content (e.g., strip # from hashtag)
				content := preprocessTokenForStemming(j.token)
				stem := stemWord(content, morphology)

				resultChan <- indexedResult{
					index: j.index,
					token: StemToken{
						Stem:      stem,
						Original:  j.token.Content,
						Type:      j.token.Type,
						StartByte: j.startByte,
						EndByte:   j.endByte,
					},
				}
			}
		}()
	}

	// Send jobs that need stemming
	for _, j := range stemmingJobs {
		jobChan <- j
	}
	close(jobChan)

	// Also add jobs that don't need stemming (as-is)
	for _, j := range jobs {
		if !j.needsStem {
			resultChan <- indexedResult{
				index: j.index,
				token: StemToken{
					Stem:      j.token.Content, // No stemming, use original
					Original:  j.token.Content,
					Type:      j.token.Type,
					StartByte: j.startByte,
					EndByte:   j.endByte,
				},
			}
		}
	}

	// Wait for workers and close results
	go func() {
		wg.Wait()
		close(resultChan)
	}()

	// Collect results (preserving order)
	resultMap := make(map[int]StemToken, len(jobs))
	for r := range resultChan {
		resultMap[r.index] = r.token
	}

	// Build final ordered result using all indices from resultMap
	result := make([]StemToken, 0, len(resultMap))

	// Find max index to ensure we process all tokens
	maxIndex := -1
	for idx := range resultMap {
		if idx > maxIndex {
			maxIndex = idx
		}
	}

	// Iterate through all possible indices in order
	for i := 0; i <= maxIndex; i++ {
		if token, ok := resultMap[i]; ok {
			result = append(result, token)
		}
	}

	return result
}

// shouldFilterToken returns true if token should be completely filtered out from results
// These tokens are not useful for FTS5 searching
func shouldFilterToken(tokenType tokenization.TokenType) bool {
	switch tokenType {
	case tokenization.Punctuation,  // . , ! ? : ; etc.
		tokenization.SpaceTab,      // space, tab (already filtered by DEFAULT tokenizer)
		tokenization.NewLine,       // \n, \r (already filtered by DEFAULT tokenizer)
		tokenization.Emoticon,      // :) :( ^_^
		tokenization.MetaTag,       // <xml>
		tokenization.Unknown,       // unidentified characters
		tokenization.UnknownWord:   // unidentified words
		return true
	default:
		return false
	}
}

// shouldSkipStemming returns true if token type should not be stemmed but should be preserved
// These tokens are kept as-is for searching (URL, Email, Number, Date, Time)
func shouldSkipStemming(tokenType tokenization.TokenType) bool {
	switch tokenType {
	case tokenization.URL,           // www.example.com
		tokenization.Email,          // user@example.com
		tokenization.Number,         // 123, 123/456, 12.5
		tokenization.PercentNumeral, // %50
		tokenization.Date,           // 15.08.2023
		tokenization.Time:           // 14:30
		return true
	default:
		return false
	}
}

// preprocessTokenForStemming prepares token content for stemming
// Removes prefixes like # for hashtags, @ for mentions
func preprocessTokenForStemming(token *tokenization.Token) string {
	content := token.Content

	switch token.Type {
	case tokenization.HashTag:
		// Remove # prefix: "#gündem" -> "gündem"
		if len(content) > 0 && content[0] == '#' {
			return content[1:]
		}
	case tokenization.Mention:
		// Remove @ prefix: "@kullanici" -> "kullanici"
		if len(content) > 0 && content[0] == '@' {
			return content[1:]
		}
	}

	return content
}

// tokensWithoutStemming returns tokens without morphological analysis
func tokensWithoutStemming(jobs []stemmingJob) []StemToken {
	result := make([]StemToken, len(jobs))

	for i, j := range jobs {
		result[i] = StemToken{
			Stem:      j.token.Content,
			Original:  j.token.Content,
			Type:      j.token.Type,
			StartByte: j.startByte,
			EndByte:   j.endByte,
		}
	}

	return result
}

// bytePosition holds byte position info
type bytePosition struct {
	start int
	end   int
}

// calculateBytePositions converts rune positions to byte positions
// Token.Start and Token.End are RUNE indices, we need BYTE offsets
func calculateBytePositions(text string, tokens []*tokenization.Token) []bytePosition {
	positions := make([]bytePosition, len(tokens))

	// Build rune->byte mapping
	runeToByteOffset := make([]int, 0, len(text))
	byteOffset := 0

	for range text {
		runeToByteOffset = append(runeToByteOffset, byteOffset)
		_, size := utf8.DecodeRuneInString(text[byteOffset:])
		byteOffset += size
	}

	// Convert token positions
	for i, token := range tokens {
		startByte := 0
		endByte := byteOffset // Default to end of text

		if token.Start < len(runeToByteOffset) {
			startByte = runeToByteOffset[token.Start]
		}

		// token.End is inclusive (last rune index)
		// We want exclusive byte offset (one past last byte)
		if token.End+1 < len(runeToByteOffset) {
			endByte = runeToByteOffset[token.End+1]
		} else if token.End < len(runeToByteOffset) {
			// Last token - calculate end byte
			endByte = runeToByteOffset[token.End]
			_, size := utf8.DecodeRuneInString(text[endByte:])
			endByte += size
		}

		positions[i] = bytePosition{
			start: startByte,
			end:   endByte,
		}
	}

	return positions
}

// stemWord performs cached morphological analysis on a single word
// Uses global cache for performance (thread-safe with sync.Map)
func stemWord(word string, morphology *TurkishMorphology) string {
	// Normalize for cache key and default stem
	// This removes dots and converts to lowercase (same as NormalizeForAnalysis)
	normalizedWord := turkish.Instance.ToLower(word)
	normalizedWord = strings.ReplaceAll(normalizedWord, ".", "")
	if normalizedWord == "" {
		normalizedWord = turkish.Instance.ToLower(word)
	}

	// 1. Cache lookup
	if cached, ok := stemCache.Load(normalizedWord); ok {
		return cached.(string)
	}

	// 2. Cache miss: perform analysis
	stem := normalizedWord // Default: use normalized if no stem found

	analysis := morphology.Analyze(word)
	if len(analysis.AnalysisResults) > 0 {
		// Prefer non-proper nouns to avoid matching proper names in dictionary
		// Example: "mahkemesi" should match "mahkeme" (3 morphemes),
		//          not "Mahkemesi" (proper noun, 2 morphemes)
		// Try non-proper nouns first, fallback to all results if none found
		candidates := make([]*analysispkg.SingleAnalysis, 0, len(analysis.AnalysisResults))
		for _, a := range analysis.AnalysisResults {
			if a.Item != nil && a.Item.SecondaryPos.GetStringForm() != "Prop" {
				candidates = append(candidates, a)
			}
		}
		// If no non-proper results, use all (proper nouns included)
		if len(candidates) == 0 {
			candidates = analysis.AnalysisResults
		}

		// Select best root using two-stage strategy:
		// 1. FEWEST morphemes (prefer base forms with fewer affixes)
		// 2. SHORTEST root (tie-breaker when morpheme counts are equal)
		//
		// Examples:
		//   "araba":
		//     [0] araba (1 morpheme) ← SELECTED (fewest morphemes)
		//     [1] araba (2 morphemes)
		//     [2] arap (3 morphemes, shorter but more complex parse)
		//
		//   "arabanın":
		//     [0] araban (6 chars, 3 morphemes)
		//     [2] araba (5 chars, 3 morphemes) ← SELECTED (same morphemes, shorter)
		//
		//   "insandan":
		//     [0] insan (5 chars, 3 morphemes) ← SELECTED (fewer morphemes)
		//     [1] insa (4 chars, 4 morphemes) - shorter but more complex
		//
		//   "mahkemesi":
		//     [0] mahkeme (3 morphemes, non-Prop) ← SELECTED (proper nouns excluded)
		//     [1] Mahkemesi (2 morphemes, Prop) - fewer morphemes but excluded
		var bestRoot string
		fewestMorphemes := -1
		shortestLen := -1

		for _, a := range candidates {
			var candidateRoot string
			morphemeCount := len(a.MorphemeDataList)

			// Use dictionary item root for correct voicing handling
			// Example: "kitabı" -> item.Root = "kitap" (not surface "kitab")
			if a.Item != nil && a.Item.Root != "" {
				candidateRoot = a.Item.Root
			} else {
				// Fallback to surface stem if no dictionary item
				candidateRoot = a.GetStem()
			}

			if candidateRoot != "" {
				// Select if: 1) fewer morphemes OR 2) same morphemes but shorter
				if fewestMorphemes == -1 ||
					morphemeCount < fewestMorphemes ||
					(morphemeCount == fewestMorphemes && len(candidateRoot) < shortestLen) {
					bestRoot = candidateRoot
					fewestMorphemes = morphemeCount
					shortestLen = len(candidateRoot)
				}
			}
		}

		if bestRoot != "" {
			stem = bestRoot
		}
	}

	// Always lowercase stems for case-insensitive FTS5 search
	// This ensures "Ankara", "ANKARA", and "ankara" all produce "ankara"
	// This ensures "CMUK", "Cmuk", and "cmuk" all produce "cmuk"
	stem = turkish.Instance.ToLower(stem)

	// 3. Store in cache for future use
	stemCache.Store(normalizedWord, stem)

	return stem
}
