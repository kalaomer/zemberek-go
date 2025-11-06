package main

import (
	"database/sql"
	"fmt"
	"log"

	// Import the custom driver with Turkish tokenizer
	_ "github.com/kalaomer/zemberek-go/sqlite_extension/driver"
)

func main() {
	// Open database with custom driver
	db, err := sql.Open("sqlite3_turkish", ":memory:")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Create FTS5 table with turkish_stem tokenizer
	_, err = db.Exec(`
		CREATE VIRTUAL TABLE documents USING fts5(
			title,
			content,
			tokenize='turkish_stem'
		)
	`)
	if err != nil {
		log.Fatalf("Failed to create FTS5 table: %v", err)
	}

	fmt.Println("✓ FTS5 table created successfully with turkish_stem tokenizer")

	// Insert some test data
	testData := []struct {
		title   string
		content string
	}{
		{"Kitaplar", "Kitapları okuyorum ve çok seviyorum"},
		{"Yazılım", "Yazılım geliştiriyorum"},
		{"Bilgisayar", "Bilgisayarlar çok gelişti"},
	}

	for _, data := range testData {
		_, err = db.Exec("INSERT INTO documents (title, content) VALUES (?, ?)",
			data.title, data.content)
		if err != nil {
			log.Fatalf("Failed to insert data: %v", err)
		}
	}

	fmt.Println("✓ Test data inserted successfully")

	// Search example - "kitap" should match "Kitaplar" and "kitapları"
	fmt.Println("\nSearching for 'kitap':")
	rows, err := db.Query("SELECT title, content FROM documents WHERE documents MATCH 'kitap'")
	if err != nil {
		log.Fatalf("Search failed: %v", err)
	}
	defer rows.Close()

	found := false
	for rows.Next() {
		var title, content string
		if err := rows.Scan(&title, &content); err != nil {
			log.Fatal(err)
		}
		fmt.Printf("  Found: %s - %s\n", title, content)
		found = true
	}

	if !found {
		fmt.Println("  No results found")
	}

	fmt.Println("\n✓ Basic example completed successfully!")
}
