package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/kalaomer/zemberek-go/sqlite_extension/driver"
)

func main() {
	// Open database
	db, err := sql.Open("sqlite3_turkish", ":memory:")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Create FTS5 table
	_, err = db.Exec(`
		CREATE VIRTUAL TABLE articles USING fts5(
			title,
			content,
			author,
			tokenize='turkish_stem'
		)
	`)
	if err != nil {
		log.Fatalf("Failed to create table: %v", err)
	}

	// Insert sample Turkish articles
	articles := []struct {
		title   string
		content string
		author  string
	}{
		{
			"Kitap İncelemesi",
			"Bu kitapları okuyorken çok keyif aldım. Yazarın kitapları gerçekten etkileyici.",
			"Ahmet Yılmaz",
		},
		{
			"Yazılım Geliştirme",
			"Yazılım geliştirirken en önemli şey test yazmaktır. Testler yazılımın kalitesini artırır.",
			"Ayşe Demir",
		},
		{
			"Bilgisayar Tarihi",
			"İlk bilgisayarlar çok büyüktü. Günümüz bilgisayarları ise çok küçük ve güçlü.",
			"Mehmet Kaya",
		},
		{
			"Okuma Alışkanlığı",
			"Düzenli okumak insanın hayal gücünü geliştirir. Her gün en az bir saat okumaya çalışıyorum.",
			"Fatma Şahin",
		},
		{
			"Teknoloji Haberleri",
			"Yeni çıkan yazılımlar çok hızlı. Geliştiriciler sürekli yenilik yapıyor.",
			"Ali Çelik",
		},
	}

	for _, article := range articles {
		_, err = db.Exec("INSERT INTO articles (title, content, author) VALUES (?, ?, ?)",
			article.title, article.content, article.author)
		if err != nil {
			log.Fatalf("Failed to insert: %v", err)
		}
	}

	fmt.Println("Sample Turkish articles database created.\n")

	// Test queries - demonstrating stemming
	queries := []string{
		"kitap",      // Should match: kitap, kitapları, kitaplar
		"yaz",        // Should match: yazarın, yazılım, yazmak, yazıyorum
		"bilgisayar", // Should match: bilgisayar, bilgisayarlar, bilgisayarları
		"oku",        // Should match: okuyorken, okumak, okumaya, okuyorum
		"geliştir",   // Should match: geliştirme, geliştirirken, geliştiriyor, geliştirir
	}

	for _, query := range queries {
		fmt.Printf("🔍 Search query: '%s'\n", query)
		fmt.Println(string(make([]rune, 50)))

		rows, err := db.Query(`
			SELECT title, content, author
			FROM articles
			WHERE articles MATCH ?
			ORDER BY rank
		`, query)

		if err != nil {
			log.Printf("Query failed: %v", err)
			continue
		}

		count := 0
		for rows.Next() {
			var title, content, author string
			if err := rows.Scan(&title, &content, &author); err != nil {
				log.Fatal(err)
			}
			count++
			fmt.Printf("  %d. %s\n", count, title)
			fmt.Printf("     Author: %s\n", author)
			fmt.Printf("     Content: %s\n", content)
			fmt.Println()
		}
		rows.Close()

		if count == 0 {
			fmt.Println("  No results found.\n")
		}

		fmt.Printf("  Total results: %d\n\n", count)
	}

	// Complex query example
	fmt.Println("🔍 Complex query: 'kitap OR yazılım'")
	fmt.Println(string(make([]rune, 50)))

	rows, err := db.Query(`
		SELECT title, content
		FROM articles
		WHERE articles MATCH 'kitap OR yazılım'
		ORDER BY rank
	`)

	if err != nil {
		log.Printf("Query failed: %v", err)
	} else {
		count := 0
		for rows.Next() {
			var title, content string
			if err := rows.Scan(&title, &content); err != nil {
				log.Fatal(err)
			}
			count++
			fmt.Printf("  %d. %s\n", count, title)
		}
		rows.Close()
		fmt.Printf("\n  Total results: %d\n", count)
	}

	fmt.Println("\n✅ Search example completed!")
}
