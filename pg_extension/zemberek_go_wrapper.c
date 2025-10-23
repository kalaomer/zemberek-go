#include "postgres.h"
#include "fmgr.h"
#include "utils/builtins.h"
#include "libzemberek_go.h"

PG_MODULE_MAGIC;

/*
 * Wrapper function for the Go hello function
 */
PG_FUNCTION_INFO_V1(zemberek_go_hello_wrapper);

Datum
zemberek_go_hello_wrapper(PG_FUNCTION_ARGS)
{
    text *input_text;
    char *input_str;
    char *result_str;
    text *result_text;
    
    /* Get the input argument */
    input_text = PG_GETARG_TEXT_PP(0);
    
    /* Convert PostgreSQL text to C string */
    input_str = text_to_cstring(input_text);
    
    /* Call the Go function */
    result_str = GoHello(input_str);
    
    /* Convert result back to PostgreSQL text */
    result_text = cstring_to_text(result_str);
    
    /* Free the Go-allocated string */
    free(result_str);
    
    PG_RETURN_TEXT_P(result_text);
}
