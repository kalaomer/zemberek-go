package sqlite_extension

/*
#cgo CFLAGS: -DSQLITE_CORE
#include <sqlite3ext.h>
#include <stdlib.h>
#include <string.h>

// Forward declarations
static int zemberekCreate(void *pCtx, const char **azArg, int nArg, Fts5Tokenizer **ppOut);
static void zemberekDelete(Fts5Tokenizer *pTokenizer);
static int zemberekTokenize(Fts5Tokenizer *pTokenizer, void *pCtx, int flags, const char *pText, int nText,
                           int (*xToken)(void*, int, const char*, int, int, int));

// Tokenizer structure
typedef struct ZemberekTokenizer {
    int dummy; // Placeholder, actual state managed in Go
} ZemberekTokenizer;

// FTS5 tokenizer module
static fts5_tokenizer zemberekTokenizerModule = {
    zemberekCreate,
    zemberekDelete,
    zemberekTokenize
};

// Create tokenizer instance
static int zemberekCreate(void *pCtx, const char **azArg, int nArg, Fts5Tokenizer **ppOut) {
    ZemberekTokenizer *pTok = (ZemberekTokenizer*)sqlite3_malloc(sizeof(ZemberekTokenizer));
    if (pTok == NULL) {
        return SQLITE_NOMEM;
    }
    memset(pTok, 0, sizeof(ZemberekTokenizer));
    *ppOut = (Fts5Tokenizer*)pTok;
    return SQLITE_OK;
}

// Delete tokenizer instance
static void zemberekDelete(Fts5Tokenizer *pTokenizer) {
    if (pTokenizer) {
        sqlite3_free(pTokenizer);
    }
}

// Tokenize function - calls into Go
extern int goTokenize(void *pCtx, int flags, const char *pText, int nText,
                     int (*xToken)(void*, int, const char*, int, int, int));

static int zemberekTokenize(Fts5Tokenizer *pTokenizer, void *pCtx, int flags,
                           const char *pText, int nText,
                           int (*xToken)(void*, int, const char*, int, int, int)) {
    return goTokenize(pCtx, flags, pText, nText, xToken);
}

// Registration function
static int registerZemberekTokenizer(sqlite3 *db, const char *zName) {
    int rc = SQLITE_OK;
    fts5_api *pApi = NULL;

    // Get FTS5 API
    sqlite3_stmt *pStmt = NULL;
    rc = sqlite3_prepare_v2(db, "SELECT fts5(?1)", -1, &pStmt, 0);
    if (rc != SQLITE_OK) {
        return rc;
    }

    sqlite3_bind_pointer(pStmt, 1, (void*)&pApi, "fts5_api_ptr", 0);
    sqlite3_step(pStmt);

    if (pApi == NULL) {
        sqlite3_finalize(pStmt);
        return SQLITE_ERROR;
    }

    sqlite3_finalize(pStmt);

    // Register tokenizer
    rc = pApi->xCreateTokenizer(pApi, zName, (void*)pApi, &zemberekTokenizerModule, NULL);

    return rc;
}
*/
import "C"
import (
	"strings"
	"unsafe"

	"github.com/kalaomer/zemberek-go/tokenization"
)

// TokenizeCallback is the type for the FTS5 token callback
type TokenizeCallback func(flags int, token string, start, end int) int

var (
	currentCallback TokenizeCallback
)

//export goTokenize
func goTokenize(pCtx unsafe.Pointer, flags C.int, pText *C.char, nText C.int,
	xToken unsafe.Pointer) C.int {

	// Convert C string to Go string
	text := C.GoStringN(pText, nText)

	// Create callback wrapper
	callback := func(tflags int, token string, start, end int) int {
		cToken := C.CString(token)
		defer C.free(unsafe.Pointer(cToken))

		// Call the FTS5 callback
		xTokenFunc := *(*func(unsafe.Pointer, C.int, *C.char, C.int, C.int, C.int) C.int)(unsafe.Pointer(&xToken))
		rc := xTokenFunc(pCtx, C.int(tflags), cToken, C.int(len(token)), C.int(start), C.int(end))
		return int(rc)
	}

	// Use zemberek tokenizer
	rc := tokenizeWithZemberek(text, callback)

	return C.int(rc)
}

// tokenizeWithZemberek performs tokenization using Zemberek
func tokenizeWithZemberek(text string, callback func(int, string, int, int) int) int {
	// Use simple tokenization for now
	tokens := tokenization.SimpleTokenize(text)

	offset := 0
	for _, token := range tokens {
		// Find token position in original text
		idx := strings.Index(text[offset:], token)
		if idx == -1 {
			continue
		}

		start := offset + idx
		end := start + len(token)

		// Skip whitespace-only tokens
		if strings.TrimSpace(token) == "" {
			offset = end
			continue
		}

		// Convert to lowercase for matching
		normalizedToken := strings.ToLower(token)

		// FTS5_TOKEN_COLOCATED = 0x0001
		rc := callback(0, normalizedToken, start, end)
		if rc != 0 {
			return rc
		}

		offset = end
	}

	return 0 // SQLITE_OK
}

// RegisterTokenizer registers the Zemberek tokenizer with SQLite
func RegisterTokenizer(db unsafe.Pointer, name string) error {
	cName := C.CString(name)
	defer C.free(unsafe.Pointer(cName))

	rc := C.registerZemberekTokenizer((*C.sqlite3)(db), cName)
	if rc != C.SQLITE_OK {
		return ErrTokenizerRegistration
	}

	return nil
}

// ErrTokenizerRegistration is returned when tokenizer registration fails
var ErrTokenizerRegistration = &TokenizerError{msg: "failed to register tokenizer"}

// TokenizerError represents a tokenizer error
type TokenizerError struct {
	msg string
}

func (e *TokenizerError) Error() string {
	return e.msg
}
