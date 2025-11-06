// Package driver provides a custom SQLite3 driver with Zemberek Turkish stemmer
// tokenizer integrated into FTS5.
//
// Usage:
//
//	import (
//	    "database/sql"
//	    _ "github.com/kalaomer/zemberek-go/sqlite_extension/driver"
//	)
//
//	func main() {
//	    db, err := sql.Open("sqlite3_turkish", "mydb.db")
//	    if err != nil {
//	        log.Fatal(err)
//	    }
//	    defer db.Close()
//
//	    // Create FTS5 table with Turkish stemmer
//	    _, err = db.Exec(`
//	        CREATE VIRTUAL TABLE documents USING fts5(
//	            title,
//	            content,
//	            tokenize='turkish_stem'
//	        )
//	    `)
//	}
package driver

import (
	"database/sql"
	"log"

	sqlite3 "github.com/mattn/go-sqlite3"
	// Import tokenizer to ensure its C code is linked
	_ "github.com/kalaomer/zemberek-go/sqlite_extension/tokenizer"
)

func init() {
	// Register custom SQLite3 driver with Zemberek tokenizer
	sql.Register("sqlite3_turkish",
		&sqlite3.SQLiteDriver{
			ConnectHook: func(conn *sqlite3.SQLiteConn) error {
				// Register the turkish_stem tokenizer on each new connection
				if err := RegisterTurkishTokenizer(conn); err != nil {
					log.Printf("Warning: Failed to register Turkish tokenizer: %v", err)
					return err
				}
				return nil
			},
		},
	)
}
