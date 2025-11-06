#include "zemberek_tokenizer.h"
#include <stdlib.h>
#include <string.h>

// Forward declaration of Go function (will be exported from bridge.go)
extern void goTokenizeText(
    const char *text,
    int nText,
    void *pCtx,
    fts5_token_callback xToken
);

// xCreate: Create tokenizer instance
int zemberekCreate(
    void *pUnused,
    const char **azArg,
    int nArg,
    Fts5Tokenizer **ppOut
) {
    (void)pUnused;
    (void)azArg;
    (void)nArg;

    ZemberekTokenizer *pTokenizer = (ZemberekTokenizer*)sqlite3_malloc(sizeof(ZemberekTokenizer));
    if (!pTokenizer) {
        return SQLITE_NOMEM;
    }

    pTokenizer->dummy = 0;
    *ppOut = (Fts5Tokenizer*)pTokenizer;

    return SQLITE_OK;
}

// xDelete: Free tokenizer instance
void zemberekDelete(Fts5Tokenizer *pTokenizer) {
    if (pTokenizer) {
        sqlite3_free(pTokenizer);
    }
}

// xTokenize: Perform tokenization by calling Go code
int zemberekTokenize(
    Fts5Tokenizer *pTokenizer,
    void *pCtx,
    int flags,
    const char *pText,
    int nText,
    const char *pLocale,
    int nLocale,
    int (*xToken)(void*, int, const char*, int, int, int)
) {
    (void)pTokenizer;
    (void)flags;
    (void)pLocale;
    (void)nLocale;

    if (!pText || nText < 0) {
        return SQLITE_OK;
    }

    // Call Go function to perform stemming and tokenization
    goTokenizeText(pText, nText, pCtx, xToken);

    return SQLITE_OK;
}

// Helper function: Invoke xToken callback from Go
void invokeTokenCallback(
    fts5_token_callback xToken,
    void *pCtx,
    const char *pToken,
    int nToken,
    int iStart,
    int iEnd
) {
    xToken(pCtx, 0, pToken, nToken, iStart, iEnd);
}

// Static fts5_tokenizer_v2 structure
static fts5_tokenizer_v2 zemberekTokenizerModule = {
    2,                  // iVersion
    zemberekCreate,     // xCreate
    zemberekDelete,     // xDelete
    zemberekTokenize    // xTokenize
};

// Getter for the tokenizer struct
fts5_tokenizer_v2* getZemberekTokenizerStruct(void) {
    return &zemberekTokenizerModule;
}
