package morphology

import (
	"fmt"
	"strings"
	"sync"
	"testing"
)

// Sample Turkish court decision text (realistic workload)
const benchmarkText = `
Mahkeme kararı incelendi. Davacı tarafın iddialarına göre, davalı taraf sözleşmeye aykırı
davranmıştır. Mahkeme, davacının sunduğu belgeleri incelemiş ve davalının savunmasını
değerlendirmiştir. Tanık ifadeleri alınmış ve bilirkişi raporu hazırlanmıştır.
Yargılama sürecinde tarafların beyanları dinlenmiş ve kanıtlar değerlendirilmiştir.
Mahkeme heyeti, dosya kapsamındaki tüm belgeleri inceledikten sonra kararını vermiştir.
Davacının taleplerinin bir kısmı kabul edilmiş, bir kısmı ise reddedilmiştir.
Karar, taraflara tebliğ edilecek ve kesinleşme sürecine girecektir.
`

// BenchmarkStemTextWithPositions_Parallel tests current parallel + cache implementation
func BenchmarkStemTextWithPositions_Parallel(b *testing.B) {
	morph := CreateWithDefaults()

	// Clear cache to measure full impact
	stemCache = sync.Map{}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = StemTextWithPositions(benchmarkText, morph)
	}
}

// BenchmarkStemTextWithPositions_ParallelCached tests with warm cache
func BenchmarkStemTextWithPositions_ParallelCached(b *testing.B) {
	morph := CreateWithDefaults()

	// Warm up cache
	_ = StemTextWithPositions(benchmarkText, morph)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = StemTextWithPositions(benchmarkText, morph)
	}
}

// BenchmarkStemTextWithPositions_LargeDocument tests with larger document
func BenchmarkStemTextWithPositions_LargeDocument(b *testing.B) {
	morph := CreateWithDefaults()

	// Create large document (10x repeated text, ~1000 words)
	largeText := strings.Repeat(benchmarkText, 10)

	// Clear cache
	stemCache = sync.Map{}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = StemTextWithPositions(largeText, morph)
	}
}

// BenchmarkStemTextWithPositions_VeryLargeDocument tests with very large document
func BenchmarkStemTextWithPositions_VeryLargeDocument(b *testing.B) {
	morph := CreateWithDefaults()

	// Create very large document (100x repeated text, ~10000 words)
	veryLargeText := strings.Repeat(benchmarkText, 100)

	// Clear cache
	stemCache = sync.Map{}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = StemTextWithPositions(veryLargeText, morph)
	}
}

// BenchmarkCacheLookup tests cache lookup performance
func BenchmarkCacheLookup(b *testing.B) {
	morph := CreateWithDefaults()

	// Warm up cache with common words
	commonWords := []string{
		"mahkeme", "karar", "davacı", "davalı", "taraf",
		"sözleşme", "belge", "tanık", "bilirkişi", "rapor",
	}

	for _, word := range commonWords {
		stemWord(word, morph)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		word := commonWords[i%len(commonWords)]
		_ = stemWord(word, morph)
	}
}

// Example output to show usage
func ExampleStemTextWithPositions_parallel() {
	morph := CreateWithDefaults()
	text := "Mahkeme davacının taleplerini inceledi."

	tokens := StemTextWithPositions(text, morph)
	for _, token := range tokens {
		fmt.Printf("%s -> %s\n", token.Original, token.Stem)
	}
	// Output will show parallel + cached stemming results
}
