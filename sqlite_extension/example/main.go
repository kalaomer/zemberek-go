package main

import (
	"database/sql"
	"fmt"
	"log"
	"strings"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	// Open database
	db, err := sql.Open("sqlite3", "./turkish_fts5.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Note: The Zemberek tokenizer must be registered through CGO
	// This example assumes the tokenizer has been compiled as an SQLite extension

	// Create FTS5 table with Zemberek tokenizer
	_, err = db.Exec(`
		CREATE VIRTUAL TABLE IF NOT EXISTS documents USING fts5(
			title,
			content,
			tokenize='zemberek'
		)
	`)
	if err != nil {
		log.Printf("Warning: Could not create FTS5 table with zemberek tokenizer: %v", err)
		log.Printf("Falling back to unicode61 tokenizer")

		// Fallback to standard tokenizer
		_, err = db.Exec(`
			CREATE VIRTUAL TABLE IF NOT EXISTS documents USING fts5(
				title,
				content,
				tokenize='unicode61 remove_diacritics 2'
			)
		`)
		if err != nil {
			log.Fatal(err)
		}
	}

	// Insert sample Turkish documents
	documents := []struct {
		title   string
		content string
	}{
		{
			title:   "Türkçe Metin",
			content: "Bu bir Türkçe metin örneğidir. Zemberek, Türkçe doğal dil işleme kütüphanesidir.",
		},
		{
			title:   "İstanbul",
			content: "İstanbul, Türkiye'nin en büyük şehridir. Boğaziçi köprüsü çok güzeldir.",
		},
		{
			title:   "Yazılım Geliştirme",
			content: "Go programlama dili hızlı ve verimlidir. SQLite hafif bir veritabanıdır.",
		},
		{
			title:   "Doğal Dil İşleme",
			content: "Doğal dil işleme, bilgisayarların insan dilini anlaması için kullanılır.",
		},
	}

	// Clear existing data
	_, err = db.Exec("DELETE FROM documents")
	if err != nil {
		log.Fatal(err)
	}

	// Insert documents
	stmt, err := db.Prepare("INSERT INTO documents(title, content) VALUES (?, ?)")
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()

	for _, doc := range documents {
		_, err := stmt.Exec(doc.title, doc.content)
		if err != nil {
			log.Fatal(err)
		}
	}

	fmt.Println("Inserted", len(documents), "documents")
	fmt.Println()

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

		rows, err := db.Query(`
			SELECT title, content, rank
			FROM documents
			WHERE documents MATCH ?
			ORDER BY rank
		`, query)
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
			fmt.Printf("      %s\n", truncate(content, 80))
			fmt.Println()
		}
		rows.Close()

		if count == 0 {
			fmt.Println("  No results found")
			fmt.Println()
		}
	}
}

func truncate(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen] + "..."
}
