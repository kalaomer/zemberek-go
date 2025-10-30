package main

import (
	"database/sql"
	"fmt"
	"log"
	"strings"

	_ "github.com/mattn/go-sqlite3"
	"github.com/kalaomer/zemberek-go/sqlite_extension"
)

func main() {
	// Open database
	db, err := sql.Open("sqlite3", "./turkish_fts5.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Create FTS5 helper
	helper := sqlite_extension.NewFTS5Helper(db)

	// Create FTS5 table
	err = helper.CreateFTS5Table("documents", "title", "content")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("✓ Created FTS5 table")

	// Clear existing data
	_, err = db.Exec("DELETE FROM documents")
	if err != nil {
		log.Fatal(err)
	}

	// Insert sample Turkish documents
	documents := []map[string]string{
		{
			"title":   "Türkçe Metin",
			"content": "Bu bir Türkçe metin örneğidir. Zemberek, Türkçe doğal dil işleme kütüphanesidir.",
		},
		{
			"title":   "İstanbul",
			"content": "İstanbul, Türkiye'nin en büyük şehridir. Boğaziçi köprüsü çok güzeldir.",
		},
		{
			"title":   "Yazılım Geliştirme",
			"content": "Go programlama dili hızlı ve verimlidir. SQLite hafif bir veritabanıdır.",
		},
		{
			"title":   "Doğal Dil İşleme",
			"content": "Doğal dil işleme, bilgisayarların insan dilini anlaması için kullanılır.",
		},
	}

	// Insert documents
	for _, doc := range documents {
		err := helper.InsertDocument("documents", doc)
		if err != nil {
			log.Fatal(err)
		}
	}

	fmt.Printf("✓ Inserted %d documents\n\n", len(documents))

	// Perform searches
	searches := []string{
		"türkçe",
		"istanbul",
		"yazılım",
		"doğal dil",
		"güzel",
	}

	for _, query := range searches {
		fmt.Printf("Searching for: %s\n", query)
		fmt.Println(strings.Repeat("-", 50))

		rows, err := helper.SearchWithRank("documents", []string{"title", "content"}, query)
		if err != nil {
			log.Printf("Search error: %v", err)
			continue
		}

		count := 0
		for rows.Next() {
			var title, content string
			var rank float64
			err := rows.Scan(&title, &content, &rank)
			if err != nil {
				log.Fatal(err)
			}

			count++
			fmt.Printf("  [%.2f] %s\n", rank, title)
			fmt.Printf("      %s\n", truncate(content, 60))
		}
		rows.Close()

		if count == 0 {
			fmt.Println("  No results found")
		}
		fmt.Println()
	}

	// Demonstrate tokenization
	fmt.Println("\n" + strings.Repeat("=", 50))
	fmt.Println("TOKENIZATION EXAMPLES")
	fmt.Println(strings.Repeat("=", 50))

	tokenizer := sqlite_extension.NewZemberekTokenizer()

	testTexts := []string{
		"İstanbul çok güzel bir şehir",
		"TÜRKÇE METIN",
		"Zemberek ile doğal dil işleme",
	}

	for _, text := range testTexts {
		tokens := tokenizer.Tokenize(text)
		fmt.Printf("\nText: %s\n", text)
		fmt.Printf("Tokens: %v\n", tokens)
	}

	// Demonstrate Turkish case conversion
	fmt.Println("\n" + strings.Repeat("=", 50))
	fmt.Println("TURKISH CASE CONVERSION")
	fmt.Println(strings.Repeat("=", 50))

	caseExamples := []struct {
		lowercase string
		uppercase string
	}{
		{"istanbul", "İSTANBUL"},
		{"ıstanbul", "ISTANBUL"},
		{"çalışma", "ÇALIŞMA"},
	}

	for _, ex := range caseExamples {
		result := sqlite_extension.TurkishUpperCase(ex.lowercase)
		fmt.Printf("  %s -> %s\n", ex.lowercase, result)
	}
}

func truncate(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen] + "..."
}
