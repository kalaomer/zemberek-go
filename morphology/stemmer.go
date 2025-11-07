package morphology

import (
	"runtime"
	"strings"
	"sync"
	"unicode"
	"unicode/utf8"
)

// Global thread-safe cache for stemming results
// Key: normalized (lowercase) word -> Value: stem
var stemCache sync.Map

// StemToken represents a stemmed token with its byte position in the original text
type StemToken struct {
	Stem      string // Stemmed form: "kitap"
	Original  string // Original word: "kitapları"
	StartByte int    // UTF-8 byte offset start
	EndByte   int    // UTF-8 byte offset end
}

// StemTextWithPositions extracts stems from text with byte positions.
//
// This is the MAIN function for FTS5 integration. It:
// 1. Tokenizes text into words (tracking byte offsets)
// 2. Performs parallel morphological analysis with caching
// 3. Extracts stems using Item.Root (for correct Turkish voicing)
// 4. Returns list of stems with their byte positions
//
// Performance optimizations:
// - Global cache (sync.Map) for repeated words
// - Worker pool for parallel analysis (uses all CPU cores)
// - Thread-safe: morphology.Analyze() is concurrent-safe
//
// Example:
//
//	morph := CreateWithDefaults()
//	text := "Kitapları okuyorum"
//	tokens := StemTextWithPositions(text, morph)
//	// tokens[0] = {Stem: "kitap", Original: "Kitapları", StartByte: 0, EndByte: 11}
//	// tokens[1] = {Stem: "oku", Original: "okuyorum", StartByte: 12, EndByte: 20}
func StemTextWithPositions(text string, morphology *TurkishMorphology) []StemToken {
	if morphology == nil {
		// Fallback: no stemming, just tokenize
		return tokenizeWithoutStemming(text)
	}

	wordInfos := tokenizeWithByteOffsets(text)

	// Filter word tokens and build job list
	type job struct {
		index int
		info  wordInfo
	}

	jobs := make([]job, 0, len(wordInfos))
	for i, info := range wordInfos {
		if isWordToken(info.text) {
			jobs = append(jobs, job{index: i, info: info})
		}
	}

	if len(jobs) == 0 {
		return []StemToken{}
	}

	// Result type with index for ordering
	type indexedResult struct {
		index int
		token StemToken
	}

	// Worker pool setup
	numWorkers := runtime.NumCPU()
	if numWorkers > len(jobs) {
		numWorkers = len(jobs) // Don't spawn more workers than jobs
	}

	jobChan := make(chan job, len(jobs))
	resultChan := make(chan indexedResult, len(jobs))

	// Start workers
	var wg sync.WaitGroup
	for w := 0; w < numWorkers; w++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := range jobChan {
				stem := stemWord(j.info.text, morphology)
				resultChan <- indexedResult{
					index: j.index,
					token: StemToken{
						Stem:      stem,
						Original:  j.info.text,
						StartByte: j.info.startByte,
						EndByte:   j.info.endByte,
					},
				}
			}
		}()
	}

	// Send jobs
	for _, j := range jobs {
		jobChan <- j
	}
	close(jobChan)

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

	// Build final ordered result
	result := make([]StemToken, 0, len(resultMap))
	for i := 0; i < len(wordInfos); i++ {
		if token, ok := resultMap[i]; ok {
			result = append(result, token)
		}
	}

	return result
}

// stemWord performs cached morphological analysis on a single word
// Uses global cache for performance (thread-safe with sync.Map)
func stemWord(word string, morphology *TurkishMorphology) string {
	// Normalize for cache key (lowercase)
	normalizedWord := strings.ToLower(word)

	// 1. Cache lookup
	if cached, ok := stemCache.Load(normalizedWord); ok {
		return cached.(string)
	}

	// 2. Cache miss: perform analysis
	stem := word // Default: use original if no stem found

	analysis := morphology.Analyze(word)
	if len(analysis.AnalysisResults) > 0 {
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
		var bestRoot string
		fewestMorphemes := -1
		shortestLen := -1

		for _, a := range analysis.AnalysisResults {
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

	// 3. Store in cache for future use
	stemCache.Store(normalizedWord, stem)

	return stem
}

// wordInfo holds tokenization info with byte positions
type wordInfo struct {
	text      string
	startByte int
	endByte   int
}

// tokenizeWithByteOffsets tokenizes text tracking UTF-8 byte offsets
func tokenizeWithByteOffsets(text string) []wordInfo {
	words := make([]wordInfo, 0)

	currentWord := ""
	wordStartByte := 0
	byteOffset := 0
	inWord := false

	for _, r := range text {
		runeByteLen := utf8.RuneLen(r)

		if isWordChar(r) {
			// Start or continue word
			if !inWord {
				wordStartByte = byteOffset
				inWord = true
			}
			currentWord += string(r)
		} else {
			// End of word
			if inWord {
				words = append(words, wordInfo{
					text:      currentWord,
					startByte: wordStartByte,
					endByte:   byteOffset,
				})
				currentWord = ""
				inWord = false
			}
			// Note: We could also emit punctuation here if needed
		}

		byteOffset += runeByteLen
	}

	// Final word if text ends with word char
	if inWord && currentWord != "" {
		words = append(words, wordInfo{
			text:      currentWord,
			startByte: wordStartByte,
			endByte:   byteOffset,
		})
	}

	return words
}

// tokenizeWithoutStemming returns tokens without morphological analysis
func tokenizeWithoutStemming(text string) []StemToken {
	wordInfos := tokenizeWithByteOffsets(text)
	result := make([]StemToken, 0, len(wordInfos))

	for _, info := range wordInfos {
		if !isWordToken(info.text) {
			continue
		}
		result = append(result, StemToken{
			Stem:      info.text,
			Original:  info.text,
			StartByte: info.startByte,
			EndByte:   info.endByte,
		})
	}

	return result
}

// isWordChar checks if rune is part of a word
func isWordChar(r rune) bool {
	// Letters (including Turkish)
	if unicode.IsLetter(r) {
		return true
	}
	// Digits
	if unicode.IsDigit(r) {
		return true
	}
	// Apostrophe for possessives: "Ali'nin"
	if r == '\'' {
		return true
	}
	return false
}

// isWordToken checks if token is a word (not just punctuation/numbers)
func isWordToken(token string) bool {
	if token == "" {
		return false
	}

	// Must contain at least one letter
	for _, r := range token {
		if unicode.IsLetter(r) {
			return true
		}
	}

	return false
}
