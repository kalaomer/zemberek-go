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

### Demo Function: zemberek_go_hello

This is a demonstration function that shows how Go goroutines work within a PostgreSQL UDF.

```sql
SELECT zemberek_go_hello('Merhaba');
```

Expected output:
```
Hello from Zemberek-Go! Results: [Goroutine-1 processed 'Merhaba'; Goroutine-2 processed 'Merhaba'; Goroutine-3 processed 'Merhaba']
```

## Adding New Functions

To add new functionality from the zemberek-go modules:

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

## Example: Adding a Morphology Analysis Function

Here's an example of how you might add a morphology analysis function:

### In `zemberek_go.go`:
```go
//export AnalyzeTurkish
func AnalyzeTurkish(word *C.char) *C.char {
    goWord := C.GoString(word)
    // Use zemberek-go morphology package
    // ... implementation ...
    result := "analysis result"
    return C.CString(result)
}
```

### In `zemberek_go_wrapper.c`:
```c
PG_FUNCTION_INFO_V1(analyze_turkish_wrapper);

Datum
analyze_turkish_wrapper(PG_FUNCTION_ARGS)
{
    text *input_text = PG_GETARG_TEXT_PP(0);
    char *input_str = text_to_cstring(input_text);
    char *result_str = AnalyzeTurkish(input_str);
    text *result_text = cstring_to_text(result_str);
    free(result_str);
    PG_RETURN_TEXT_P(result_text);
}
```

### In `zemberek_go--0.1.0.sql`:
```sql
CREATE OR REPLACE FUNCTION analyze_turkish(word TEXT)
RETURNS TEXT
AS 'MODULE_PATHNAME', 'analyze_turkish_wrapper'
LANGUAGE C STRICT;
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
