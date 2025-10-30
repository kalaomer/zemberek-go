package sqlite_extension

import (
	"database/sql"
	"fmt"
	"strings"
)

// FTS5Helper provides helper functions for working with SQLite FTS5
type FTS5Helper struct {
	db        *sql.DB
	tokenizer *ZemberekTokenizer
}

// NewFTS5Helper creates a new FTS5 helper
func NewFTS5Helper(db *sql.DB) *FTS5Helper {
	return &FTS5Helper{
		db:        db,
		tokenizer: NewZemberekTokenizer(),
	}
}

// NewFTS5HelperWithTokenizer creates a new FTS5 helper with a custom tokenizer
func NewFTS5HelperWithTokenizer(db *sql.DB, tokenizer *ZemberekTokenizer) *FTS5Helper {
	return &FTS5Helper{
		db:        db,
		tokenizer: tokenizer,
	}
}

// CreateFTS5Table creates an FTS5 virtual table with Zemberek-compatible configuration
// Note: Uses unicode61 tokenizer with Turkish-friendly settings
func (h *FTS5Helper) CreateFTS5Table(tableName string, columns ...string) error {
	if len(columns) == 0 {
		return fmt.Errorf("at least one column is required")
	}

	// Use unicode61 with remove_diacritics for better Turkish support
	// While not perfect, it's better than ascii tokenizer
	sql := fmt.Sprintf(`
		CREATE VIRTUAL TABLE IF NOT EXISTS %s USING fts5(
			%s,
			tokenize='unicode61 remove_diacritics 2'
		)
	`, tableName, strings.Join(columns, ", "))

	_, err := h.db.Exec(sql)
	return err
}

// InsertDocument inserts a document into the FTS5 table
func (h *FTS5Helper) InsertDocument(tableName string, values map[string]string) error {
	if len(values) == 0 {
		return fmt.Errorf("no values provided")
	}

	columns := make([]string, 0, len(values))
	placeholders := make([]string, 0, len(values))
	args := make([]interface{}, 0, len(values))

	for col, val := range values {
		columns = append(columns, col)
		placeholders = append(placeholders, "?")
		args = append(args, val)
	}

	sql := fmt.Sprintf("INSERT INTO %s(%s) VALUES (%s)",
		tableName,
		strings.Join(columns, ", "),
		strings.Join(placeholders, ", "))

	_, err := h.db.Exec(sql, args...)
	return err
}

// Search performs a full-text search and returns results
func (h *FTS5Helper) Search(tableName, query string) (*sql.Rows, error) {
	// Normalize the search query using Zemberek tokenizer
	normalizedQuery := h.tokenizer.Tokenize(query)
	if len(normalizedQuery) == 0 {
		return nil, fmt.Errorf("empty search query")
	}

	// Join tokens with OR for broader matching
	searchTerms := strings.Join(normalizedQuery, " OR ")

	sql := fmt.Sprintf("SELECT * FROM %s WHERE %s MATCH ?", tableName, tableName)
	return h.db.Query(sql, searchTerms)
}

// SearchWithRank performs a full-text search with ranking
func (h *FTS5Helper) SearchWithRank(tableName string, columns []string, query string) (*sql.Rows, error) {
	// Normalize the search query
	normalizedQuery := h.tokenizer.Tokenize(query)
	if len(normalizedQuery) == 0 {
		return nil, fmt.Errorf("empty search query")
	}

	searchTerms := strings.Join(normalizedQuery, " OR ")

	columnList := "*"
	if len(columns) > 0 {
		columnList = strings.Join(columns, ", ")
	}

	sql := fmt.Sprintf(`
		SELECT %s, rank
		FROM %s
		WHERE %s MATCH ?
		ORDER BY rank
	`, columnList, tableName, tableName)

	return h.db.Query(sql, searchTerms)
}

// HighlightMatches returns snippets with highlighted matches
func (h *FTS5Helper) HighlightMatches(tableName, column, query string, maxSnippets int) (*sql.Rows, error) {
	normalizedQuery := h.tokenizer.Tokenize(query)
	if len(normalizedQuery) == 0 {
		return nil, fmt.Errorf("empty search query")
	}

	searchTerms := strings.Join(normalizedQuery, " OR ")

	sql := fmt.Sprintf(`
		SELECT snippet(%s, -1, '[', ']', '...', %d) as snippet, rank
		FROM %s
		WHERE %s MATCH ?
		ORDER BY rank
	`, tableName, maxSnippets, tableName, tableName)

	return h.db.Query(sql, searchTerms)
}

// NormalizeQuery normalizes a query string using the Zemberek tokenizer
func (h *FTS5Helper) NormalizeQuery(query string) string {
	tokens := h.tokenizer.Tokenize(query)
	return strings.Join(tokens, " ")
}

// Example usage and helper documentation:
//
// Creating an FTS5 table:
//
//   helper := sqlite_extension.NewFTS5Helper(db)
//   err := helper.CreateFTS5Table("documents", "title", "content")
//
// Inserting documents:
//
//   err := helper.InsertDocument("documents", map[string]string{
//       "title": "Türkçe Başlık",
//       "content": "Bu bir Türkçe metin örneğidir.",
//   })
//
// Searching:
//
//   rows, err := helper.Search("documents", "türkçe metin")
//   defer rows.Close()
//   for rows.Next() {
//       // Process results
//   }
//
// Searching with rank:
//
//   rows, err := helper.SearchWithRank("documents",
//       []string{"title", "content"}, "türkçe")
//   defer rows.Close()
//   for rows.Next() {
//       var title, content string
//       var rank float64
//       rows.Scan(&title, &content, &rank)
//       fmt.Printf("[%.2f] %s\n", rank, title)
//   }
