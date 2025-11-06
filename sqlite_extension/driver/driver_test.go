package driver

import (
	"database/sql"
	"testing"
)

func TestDriverRegistration(t *testing.T) {
	// Test that the driver is registered
	db, err := sql.Open("sqlite3_turkish", ":memory:")
	if err != nil {
		t.Fatalf("Failed to open database: %v", err)
	}
	defer db.Close()

	// Ping to ensure connection works
	if err := db.Ping(); err != nil {
		t.Fatalf("Failed to ping database: %v", err)
	}
}

func TestFTS5TableCreation(t *testing.T) {
	db, err := sql.Open("sqlite3_turkish", ":memory:")
	if err != nil {
		t.Fatalf("Failed to open database: %v", err)
	}
	defer db.Close()

	// Create FTS5 table with turkish_stem tokenizer
	_, err = db.Exec(`
		CREATE VIRTUAL TABLE test_docs USING fts5(
			content,
			tokenize='turkish_stem'
		)
	`)
	if err != nil {
		t.Fatalf("Failed to create FTS5 table: %v", err)
	}

	// Verify table was created
	var tableName string
	err = db.QueryRow("SELECT name FROM sqlite_master WHERE type='table' AND name='test_docs'").Scan(&tableName)
	if err != nil {
		t.Fatalf("Failed to query table: %v", err)
	}

	if tableName != "test_docs" {
		t.Errorf("Expected table name 'test_docs', got '%s'", tableName)
	}
}

func TestBasicInsertAndSearch(t *testing.T) {
	db, err := sql.Open("sqlite3_turkish", ":memory:")
	if err != nil {
		t.Fatalf("Failed to open database: %v", err)
	}
	defer db.Close()

	// Create table
	_, err = db.Exec(`
		CREATE VIRTUAL TABLE documents USING fts5(
			title,
			content,
			tokenize='turkish_stem'
		)
	`)
	if err != nil {
		t.Fatalf("Failed to create table: %v", err)
	}

	// Insert test data
	testData := []struct {
		title   string
		content string
	}{
		{"Kitaplar", "Kitapları okuyorum"},
		{"Yazılım", "Yazılım geliştiriyorum"},
		{"Okul", "Okulda ders çalışıyorum"},
	}

	for _, data := range testData {
		_, err = db.Exec("INSERT INTO documents (title, content) VALUES (?, ?)",
			data.title, data.content)
		if err != nil {
			t.Fatalf("Failed to insert data: %v", err)
		}
	}

	// Test search - "kitap" should match "Kitaplar" and "Kitapları"
	t.Run("Search for kitap", func(t *testing.T) {
		rows, err := db.Query("SELECT title FROM documents WHERE documents MATCH 'kitap'")
		if err != nil {
			t.Fatalf("Search failed: %v", err)
		}
		defer rows.Close()

		count := 0
		for rows.Next() {
			var title string
			if err := rows.Scan(&title); err != nil {
				t.Fatal(err)
			}
			count++
			if title != "Kitaplar" {
				t.Errorf("Expected title 'Kitaplar', got '%s'", title)
			}
		}

		if count == 0 {
			t.Error("Search returned no results")
		}
	})

	// Test search - "yazılım" should match "Yazılım" (exact stem match)
	t.Run("Search for yazılım", func(t *testing.T) {
		rows, err := db.Query("SELECT title FROM documents WHERE documents MATCH 'yazılım'")
		if err != nil {
			t.Fatalf("Search failed: %v", err)
		}
		defer rows.Close()

		found := false
		for rows.Next() {
			var title string
			if err := rows.Scan(&title); err != nil {
				t.Fatal(err)
			}
			found = true
			if title != "Yazılım" {
				t.Errorf("Expected title 'Yazılım', got '%s'", title)
			}
		}

		if !found {
			t.Error("Search for 'yazılım' returned no results")
		}
	})

	// Test search - "okuyorum" should match "Kitapları okuyorum"
	// (both stem to "oku" when in context)
	t.Run("Search for okuyorum", func(t *testing.T) {
		rows, err := db.Query("SELECT title FROM documents WHERE documents MATCH 'okuyorum'")
		if err != nil {
			t.Fatalf("Search failed: %v", err)
		}
		defer rows.Close()

		count := 0
		for rows.Next() {
			var title string
			if err := rows.Scan(&title); err != nil {
				t.Fatal(err)
			}
			count++
		}

		if count == 0 {
			t.Error("Search for 'okuyorum' returned no results")
		}
	})
}

func TestComplexQueries(t *testing.T) {
	db, err := sql.Open("sqlite3_turkish", ":memory:")
	if err != nil {
		t.Fatalf("Failed to open database: %v", err)
	}
	defer db.Close()

	// Create and populate table
	_, err = db.Exec(`
		CREATE VIRTUAL TABLE docs USING fts5(content, tokenize='turkish_stem')
	`)
	if err != nil {
		t.Fatalf("Failed to create table: %v", err)
	}

	docs := []string{
		"Kitapları okuyorum",
		"Yazılım geliştiriyorum",
		"Kitap ve yazılım",
	}

	for _, doc := range docs {
		_, err = db.Exec("INSERT INTO docs (content) VALUES (?)", doc)
		if err != nil {
			t.Fatalf("Failed to insert: %v", err)
		}
	}

	// Test OR query
	t.Run("OR query", func(t *testing.T) {
		rows, err := db.Query("SELECT content FROM docs WHERE docs MATCH 'kitap OR yazılım'")
		if err != nil {
			t.Fatalf("Query failed: %v", err)
		}
		defer rows.Close()

		count := 0
		for rows.Next() {
			count++
			var content string
			rows.Scan(&content)
		}

		if count != 3 {
			t.Errorf("Expected 3 results, got %d", count)
		}
	})

	// Test AND query
	t.Run("AND query", func(t *testing.T) {
		rows, err := db.Query("SELECT content FROM docs WHERE docs MATCH 'kitap AND yazılım'")
		if err != nil {
			t.Fatalf("Query failed: %v", err)
		}
		defer rows.Close()

		count := 0
		expectedContent := "Kitap ve yazılım"
		for rows.Next() {
			count++
			var content string
			rows.Scan(&content)
			if content != expectedContent {
				t.Errorf("Expected '%s', got '%s'", expectedContent, content)
			}
		}

		if count != 1 {
			t.Errorf("Expected 1 result, got %d", count)
		}
	})
}

func TestMultipleConnections(t *testing.T) {
	// Test that the tokenizer works across multiple connections
	for i := 0; i < 3; i++ {
		db, err := sql.Open("sqlite3_turkish", ":memory:")
		if err != nil {
			t.Fatalf("Connection %d: Failed to open database: %v", i, err)
		}

		_, err = db.Exec(`
			CREATE VIRTUAL TABLE test USING fts5(content, tokenize='turkish_stem')
		`)
		if err != nil {
			db.Close()
			t.Fatalf("Connection %d: Failed to create table: %v", i, err)
		}

		db.Close()
	}
}

func BenchmarkFTS5Insert(b *testing.B) {
	db, err := sql.Open("sqlite3_turkish", ":memory:")
	if err != nil {
		b.Fatal(err)
	}
	defer db.Close()

	_, err = db.Exec(`
		CREATE VIRTUAL TABLE bench USING fts5(content, tokenize='turkish_stem')
	`)
	if err != nil {
		b.Fatal(err)
	}

	text := "Kitapları okuyorum ve çok seviyorum. Yazılım geliştiriyorum."

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err = db.Exec("INSERT INTO bench (content) VALUES (?)", text)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkFTS5Search(b *testing.B) {
	db, err := sql.Open("sqlite3_turkish", ":memory:")
	if err != nil {
		b.Fatal(err)
	}
	defer db.Close()

	_, err = db.Exec(`
		CREATE VIRTUAL TABLE bench USING fts5(content, tokenize='turkish_stem')
	`)
	if err != nil {
		b.Fatal(err)
	}

	// Insert test data
	for i := 0; i < 100; i++ {
		_, err = db.Exec("INSERT INTO bench (content) VALUES (?)",
			"Kitapları okuyorum ve çok seviyorum")
		if err != nil {
			b.Fatal(err)
		}
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		rows, err := db.Query("SELECT content FROM bench WHERE bench MATCH 'kitap'")
		if err != nil {
			b.Fatal(err)
		}
		for rows.Next() {
			var content string
			rows.Scan(&content)
		}
		rows.Close()
	}
}
