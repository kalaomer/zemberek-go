package morphology

import (
	"runtime"
	"sync"
)

// deduplicateAndStem processes unique words and returns stem map.
// This is critical optimization: 1000 tokens might have only 100 unique words (10x dedup).
func deduplicateAndStem(jobs []stemmingJob, morphology *TurkishMorphology) map[string]string {
	// Build unique word set
	uniqueWords := make(map[string]bool)
	for _, j := range jobs {
		content := preprocessTokenForStemming(j.token)
		uniqueWords[content] = true
	}

	// Stem only unique words
	stemResults := make(map[string]string, len(uniqueWords))
	for word := range uniqueWords {
		stemResults[word] = stemWord(word, morphology)
	}

	return stemResults
}

// processStemsParallel processes unique words in parallel using worker pool.
// Returns map of word -> stem for all unique words.
func processStemsParallel(jobs []stemmingJob, morphology *TurkishMorphology) map[string]string {
	// Build unique word set
	uniqueWords := make(map[string]bool)
	for _, j := range jobs {
		content := preprocessTokenForStemming(j.token)
		uniqueWords[content] = true
	}

	// Convert to slice for worker pool
	wordList := make([]string, 0, len(uniqueWords))
	for word := range uniqueWords {
		wordList = append(wordList, word)
	}

	// Worker pool setup
	numWorkers := runtime.NumCPU()
	if numWorkers > len(wordList) {
		numWorkers = len(wordList)
	}

	type stemResult struct {
		word string
		stem string
	}

	jobChan := make(chan string, len(wordList))
	resultChan := make(chan stemResult, len(wordList))

	// Start workers
	var wg sync.WaitGroup
	for w := 0; w < numWorkers; w++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for word := range jobChan {
				stem := stemWord(word, morphology)
				resultChan <- stemResult{word: word, stem: stem}
			}
		}()
	}

	// Send jobs
	for _, word := range wordList {
		jobChan <- word
	}
	close(jobChan)

	// Wait and close results
	go func() {
		wg.Wait()
		close(resultChan)
	}()

	// Collect results
	stemResults := make(map[string]string, len(wordList))
	for sr := range resultChan {
		stemResults[sr.word] = sr.stem
	}

	return stemResults
}

// buildResultMap creates final StemToken map from jobs and stem results.
// Handles both stemmed and non-stemmed tokens.
func buildResultMap(jobs []stemmingJob, stemResults map[string]string) map[int]StemToken {
	resultMap := make(map[int]StemToken, len(jobs))

	for _, j := range jobs {
		var stem string
		if j.needsStem {
			content := preprocessTokenForStemming(j.token)
			stem = stemResults[content]
		} else {
			stem = j.token.Content // No stemming
		}

		resultMap[j.index] = StemToken{
			Stem:      stem,
			Original:  j.token.Content,
			Type:      j.token.Type,
			StartByte: j.startByte,
			EndByte:   j.endByte,
		}
	}

	return resultMap
}

// orderedResults converts result map to ordered slice.
func orderedResults(resultMap map[int]StemToken) []StemToken {
	// Find max index
	maxIndex := -1
	for idx := range resultMap {
		if idx > maxIndex {
			maxIndex = idx
		}
	}

	// Build ordered result
	result := make([]StemToken, 0, len(resultMap))
	for i := 0; i <= maxIndex; i++ {
		if token, ok := resultMap[i]; ok {
			result = append(result, token)
		}
	}

	return result
}
