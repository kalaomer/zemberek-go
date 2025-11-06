#ifndef ZEMBEREK_TOKENIZER_H
#define ZEMBEREK_TOKENIZER_H

// Forward declarations from sqlite3
typedef struct sqlite3 sqlite3;
typedef struct sqlite3_stmt sqlite3_stmt;

// Include our own FTS5 header (which needs sqlite3 types)
// We'll define fts5_tokenizer_v2 ourselves instead
#include <stdint.h>

// FTS5 API structures (simplified, from fts5.h)
typedef struct Fts5Tokenizer Fts5Tokenizer;

// fts5_tokenizer_v2 structure
typedef struct fts5_tokenizer_v2 fts5_tokenizer_v2;
struct fts5_tokenizer_v2 {
  int iVersion;        /* Currently always 2 */

  int (*xCreate)(void*, const char **azArg, int nArg, Fts5Tokenizer **ppOut);
  void (*xDelete)(Fts5Tokenizer*);
  int (*xTokenize)(Fts5Tokenizer*,
      void *pCtx,
      int flags,            /* Mask of FTS5_TOKENIZE_* flags */
      const char *pText, int nText,
      const char *pLocale, int nLocale,
      int (*xToken)(
        void *pCtx,         /* Copy of 2nd argument to xTokenize() */
        int tflags,         /* Mask of FTS5_TOKEN_* flags */
        const char *pToken, /* Pointer to buffer containing token */
        int nToken,         /* Size of token in bytes */
        int iStart,         /* Byte offset of token within input text */
        int iEnd            /* Byte offset of end of token within input text */
      )
  );
};

// SQLite constants
#define SQLITE_OK           0
#define SQLITE_ERROR        1
#define SQLITE_NOMEM        7

// SQLite functions we need
extern void* sqlite3_malloc(int);
extern void sqlite3_free(void*);

// Tokenizer instance structure
typedef struct ZemberekTokenizer {
    int dummy; // Placeholder - we use Go's global morphology instance
} ZemberekTokenizer;

// Token callback type for easier reference
typedef int (*fts5_token_callback)(
    void *pCtx,
    int tflags,
    const char *pToken,
    int nToken,
    int iStart,
    int iEnd
);

// FTS5 Tokenizer v2 callbacks
int zemberekCreate(
    void *pUnused,
    const char **azArg,
    int nArg,
    Fts5Tokenizer **ppOut
);

void zemberekDelete(Fts5Tokenizer *pTokenizer);

int zemberekTokenize(
    Fts5Tokenizer *pTokenizer,
    void *pCtx,
    int flags,
    const char *pText,
    int nText,
    const char *pLocale,
    int nLocale,
    int (*xToken)(void*, int, const char*, int, int, int)
);

// Helper function to invoke xToken callback from Go
void invokeTokenCallback(
    fts5_token_callback xToken,
    void *pCtx,
    const char *pToken,
    int nToken,
    int iStart,
    int iEnd
);

// Get the fts5_tokenizer_v2 struct
fts5_tokenizer_v2* getZemberekTokenizerStruct(void);

#endif // ZEMBEREK_TOKENIZER_H
