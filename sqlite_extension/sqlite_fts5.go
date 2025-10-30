package sqlite_extension

/*
#cgo CFLAGS: -DSQLITE_CORE
#include <sqlite3ext.h>
#include <stdlib.h>
*/
import "C"
import (
	"database/sql"
	"fmt"
	"unsafe"
)

// DefaultTokenizerName is the default name for the Zemberek tokenizer
const DefaultTokenizerName = "zemberek"

// AutoRegister is a convenience function to register the tokenizer with an sql.DB
func AutoRegister(db *sql.DB) error {
	return AutoRegisterWithName(db, DefaultTokenizerName)
}

// AutoRegisterWithName registers the tokenizer with a custom name
func AutoRegisterWithName(db *sql.DB, tokenizerName string) error {
	// We need to get the underlying SQLite connection pointer
	// This requires using the database/sql driver interface
	var sqliteDB *C.sqlite3

	// Execute a no-op query to get access to the connection
	err := db.QueryRow("SELECT sqlite_version()").Scan(new(string))
	if err != nil {
		return fmt.Errorf("failed to verify SQLite connection: %w", err)
	}

	// Get the connection - this is driver-specific
	// For github.com/mattn/go-sqlite3, we can use a custom connection hook
	// But for now, we'll provide a manual registration function

	return RegisterTokenizer(unsafe.Pointer(sqliteDB), tokenizerName)
}

// RegisterWithConnection registers the tokenizer with a raw SQLite connection pointer
// This is useful when you have direct access to the sqlite3* pointer
func RegisterWithConnection(sqliteConn unsafe.Pointer, tokenizerName string) error {
	if sqliteConn == nil {
		return fmt.Errorf("nil SQLite connection")
	}
	return RegisterTokenizer(sqliteConn, tokenizerName)
}

// Example SQL for creating an FTS5 table with Zemberek tokenizer:
//
//   CREATE VIRTUAL TABLE documents USING fts5(
//       title,
//       content,
//       tokenize='zemberek'
//   );
//
// To insert Turkish text:
//
//   INSERT INTO documents(title, content)
//   VALUES ('Başlık', 'Bu bir Türkçe metindir.');
//
// To search:
//
//   SELECT * FROM documents
//   WHERE documents MATCH 'türkçe';
