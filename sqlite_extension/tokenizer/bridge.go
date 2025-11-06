package tokenizer

/*
#cgo CFLAGS: -DSQLITE_ENABLE_FTS5
#cgo LDFLAGS: -lsqlite3
#include <stdlib.h>
#include "zemberek_tokenizer.h"
*/
import "C"
import (
	"sync"
	"unsafe"

	"github.com/kalaomer/zemberek-go/morphology"
)

var (
	// Global morphology instance (created once, reused)
	globalMorph *morphology.TurkishMorphology
	morphOnce   sync.Once
	morphMutex  sync.RWMutex
)

// getMorphology returns the global morphology instance (thread-safe singleton)
func getMorphology() *morphology.TurkishMorphology {
	morphOnce.Do(func() {
		globalMorph = morphology.CreateWithDefaults()
	})
	return globalMorph
}

//export goTokenizeText
func goTokenizeText(text *C.char, nText C.int, pCtx unsafe.Pointer, xToken C.fts5_token_callback) {
	// Convert C string to Go string
	goText := C.GoStringN(text, nText)

	// Get morphology instance (thread-safe)
	morphMutex.RLock()
	morph := getMorphology()
	morphMutex.RUnlock()

	// Perform stemming and tokenization
	tokens := morphology.StemTextWithPositions(goText, morph)

	// Invoke xToken callback for each token
	for _, token := range tokens {
		// Convert stem to C string
		cStem := C.CString(token.Stem)

		// Call the xToken callback via C helper function
		C.invokeTokenCallback(
			xToken,
			pCtx,
			cStem,
			C.int(len(token.Stem)),
			C.int(token.StartByte),
			C.int(token.EndByte),
		)

		// Free C string
		C.free(unsafe.Pointer(cStem))
	}
}

// GetTokenizerStruct returns a pointer to the FTS5 tokenizer structure
func GetTokenizerStruct() unsafe.Pointer {
	return unsafe.Pointer(C.getZemberekTokenizerStruct())
}
