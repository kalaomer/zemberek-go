# Zemberek-Go PostgreSQL Extension

This directory contains a PostgreSQL extension that exposes Zemberek Turkish NLP functionality as User-Defined Functions (UDFs) in PostgreSQL.

## Architecture

The extension uses the following architecture:
- **Go code** (`zemberek_go.go`): Contains the actual implementation using Go and goroutines
- **C wrapper** (`zemberek_go_wrapper.c`): Bridges PostgreSQL's C API with Go code via CGO
- **SQL definitions** (`zemberek_go--0.1.0.sql`): Defines the SQL functions available in PostgreSQL
- **Control file** (`zemberek_go.control`): Extension metadata for PostgreSQL

## Prerequisites

- PostgreSQL development headers (postgresql-server-dev)
- Go 1.16 or higher
- GCC or compatible C compiler
- Make

### Installing Prerequisites

**macOS:**
```bash
brew install postgresql
```

**Ubuntu/Debian:**
```bash
sudo apt-get install postgresql-server-dev-all build-essential
```

**RHEL/CentOS:**
```bash
sudo yum install postgresql-devel gcc make
```

## Building the Extension

1. Ensure `pg_config` is in your PATH:
```bash
which pg_config
```

2. Build and install the extension:
```bash
cd pg_extension
make
sudo make install
```

## Installing the Extension in PostgreSQL

Connect to your PostgreSQL database and run:

```sql
CREATE EXTENSION zemberek_go;
```

## Usage

### Normalize Turkish Text

Normalizes informal Turkish text to formal Turkish:

```sql
SELECT zemberek_normalize('mrhba nasilsin');
```

### Morphological Analysis

Analyzes the morphological structure of Turkish words:

```sql
SELECT zemberek_analyze('kitaplarımızdan');
```

Returns detailed morphological analysis showing stems, suffixes, and grammatical features.

### Extract Word Stem

Extracts the root/stem from a Turkish word:

```sql
SELECT zemberek_stem('kitaplarımızdan');
-- Returns: kitap
```

### Check if Word Has Valid Analysis

Returns true if the word has valid morphological analysis (useful for spell checking):

```sql
SELECT zemberek_has_analysis('kitap');
-- Returns: true

SELECT zemberek_has_analysis('xyzabc');
-- Returns: false
```

### Batch Processing

Process multiple words or sentences:

```sql
-- Normalize a column of Turkish text
SELECT id, zemberek_normalize(text_column) 
FROM your_table;

-- Find words without valid analysis (potential typos)
SELECT word 
FROM turkish_words 
WHERE NOT zemberek_has_analysis(word);

-- Extract stems for search indexing
SELECT DISTINCT zemberek_stem(word) 
FROM turkish_vocabulary;
```

## Adding New Functions

To add more functionality from the zemberek-go modules:

1. **Add Go function** in `zemberek_go.go`:
   - Use `//export FunctionName` comment before the function
   - Use C types in the signature (or convert Go types)
   - Return C types using `C.CString()` for strings

2. **Add C wrapper** in `zemberek_go_wrapper.c`:
   - Create a new `PG_FUNCTION_INFO_V1()` declaration
   - Implement the wrapper function to convert PostgreSQL types to C types
   - Call your Go function
   - Convert the result back to PostgreSQL types

3. **Add SQL function** in `zemberek_go--0.1.0.sql`:
   - Define the SQL function signature
   - Reference the C wrapper function name

4. **Rebuild and reinstall**:
```bash
make clean
make
sudo make install
```

5. **Update the extension** in your database:
```sql
DROP EXTENSION zemberek_go CASCADE;
CREATE EXTENSION zemberek_go;
```

## Example: Adding Additional Functions

The current implementation provides four core functions. Here's how the pattern works if you want to add more:

### Example Structure

**In `zemberek_go.go`:**
```go
//export YourFunction
func YourFunction(input *C.char) *C.char {
    goInput := C.GoString(input)
    // Process using zemberek-go libraries
    result := processInput(goInput)
    return C.CString(result)
}
```

**In `zemberek_go_wrapper.c`:**
```c
PG_FUNCTION_INFO_V1(your_function);

Datum
your_function(PG_FUNCTION_ARGS)
{
    text *input_text = PG_GETARG_TEXT_PP(0);
    char *input_str = text_to_cstring(input_text);
    char *result_str = YourFunction(input_str);
    text *result_text = cstring_to_text(result_str);
    free(result_str);
    PG_RETURN_TEXT_P(result_text);
}
```

**In `zemberek_go--0.1.0.sql`:**
```sql
CREATE OR REPLACE FUNCTION your_function(input TEXT)
RETURNS TEXT
AS 'MODULE_PATHNAME', 'your_function'
LANGUAGE C STRICT;
```

**In `libzemberek_go.h`:**
```c
extern char* YourFunction(char* input);
```

## Troubleshooting

### Error: "pg_config: command not found"
Ensure PostgreSQL is installed and `pg_config` is in your PATH.

### Error: "could not access file"
Run `sudo make install` to copy the extension files to PostgreSQL's library directory.

### Extension not loading
Check PostgreSQL logs for detailed error messages:
```bash
tail -f /var/log/postgresql/postgresql-*.log
```

### Memory leaks
Always free C strings allocated by Go using `free()` in the C wrapper after converting them to PostgreSQL types.

## Notes

- The current implementation uses `buildmode=c-archive` which creates a static library
- All Go strings returned via CGO must be freed by the C code using `free()`
- Goroutines work normally within PostgreSQL, but be mindful of resource usage
- For production use, consider adding proper error handling and memory management
- The extension is marked as `relocatable = true` in the control file

## License

This extension follows the same license as the parent zemberek-go project.
