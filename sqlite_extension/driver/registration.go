package driver

/*
#cgo CFLAGS: -DSQLITE_ENABLE_FTS5
#include <sqlite3.h>
#include <stdlib.h>

// Forward type declarations (Fts5Tokenizer not in system sqlite3.h)
typedef struct Fts5Tokenizer Fts5Tokenizer;

// Import C functions from tokenizer package
extern int zemberekCreate(void*, const char**, int, Fts5Tokenizer**);
extern void zemberekDelete(Fts5Tokenizer*);
extern int zemberekTokenize(Fts5Tokenizer*, void*, int, const char*, int, const char*, int, void*);

// C helper to get fts5_api pointer
static fts5_api* getFts5Api(sqlite3 *db) {
    fts5_api *pApi = NULL;
    sqlite3_stmt *stmt = NULL;

    if (sqlite3_prepare_v2(db, "SELECT fts5(?1)", -1, &stmt, NULL) != SQLITE_OK) {
        return NULL;
    }

    sqlite3_bind_pointer(stmt, 1, (void*)&pApi, "fts5_api_ptr", NULL);
    sqlite3_step(stmt);
    sqlite3_finalize(stmt);

    return pApi;
}

// Wrapper for xTokenize that ignores locale (v1 compatible)
static int zemberekTokenizeWrapper(
    Fts5Tokenizer *pTokenizer,
    void *pCtx,
    int flags,
    const char *pText,
    int nText,
    int (*xToken)(void*, int, const char*, int, int, int)
) {
    // Call v2 tokenize with NULL locale
    return zemberekTokenize(pTokenizer, pCtx, flags, pText, nText, NULL, 0, (void*)xToken);
}

// Register the tokenizer
static int registerZemberekTokenizer(sqlite3 *db) {
    fts5_api *pApi = getFts5Api(db);
    if (!pApi) {
        return SQLITE_ERROR;
    }

    // Create v1 tokenizer struct on stack
    fts5_tokenizer tokenizer = {
        zemberekCreate,
        zemberekDelete,
        zemberekTokenizeWrapper
    };

    // Register using v1 API
    return pApi->xCreateTokenizer(pApi, "turkish_stem", NULL, &tokenizer, NULL);
}
*/
import "C"
import (
	"fmt"
	"reflect"

	sqlite3 "github.com/mattn/go-sqlite3"
)

// RegisterTurkishTokenizer registers the Zemberek Turkish stemmer tokenizer
func RegisterTurkishTokenizer(conn *sqlite3.SQLiteConn) error {
	// Extract sqlite3* handle using reflection
	db := extractSQLite3Handle(conn)
	if db == nil {
		return fmt.Errorf("failed to extract sqlite3 handle from connection")
	}

	// Register the tokenizer
	rc := C.registerZemberekTokenizer(db)
	if rc != C.SQLITE_OK {
		return fmt.Errorf("failed to register turkish_stem tokenizer: error %d", rc)
	}

	return nil
}

// extractSQLite3Handle uses reflection to get the unexported db field from SQLiteConn
func extractSQLite3Handle(conn *sqlite3.SQLiteConn) *C.sqlite3 {
	// SQLiteConn structure:
	// type SQLiteConn struct {
	//     mu          sync.Mutex    <- Field 0
	//     db          *C.sqlite3    <- Field 1 (what we want!)
	//     loc         *time.Location
	//     ...
	// }

	connValue := reflect.ValueOf(conn)
	if connValue.Kind() == reflect.Ptr {
		connValue = connValue.Elem()
	}

	// Get field by index (field 1 is 'db')
	dbField := connValue.Field(1)
	if !dbField.IsValid() {
		return nil
	}

	// Field is already a pointer type (*C.sqlite3), get its value
	// Use UnsafePointer() to get the pointer value itself, not its address
	dbPtr := (*C.sqlite3)(dbField.UnsafePointer())
	return dbPtr
}
