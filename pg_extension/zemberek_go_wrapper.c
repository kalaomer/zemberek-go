#include "postgres.h"
#include "fmgr.h"
#include "utils/builtins.h"
#include "libzemberek_go.h"

PG_MODULE_MAGIC;

/*
 * Normalize Turkish text
 */
PG_FUNCTION_INFO_V1(zemberek_normalize);

Datum
zemberek_normalize(PG_FUNCTION_ARGS)
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
    result_str = NormalizeTurkish(input_str);
    
    /* Convert result back to PostgreSQL text */
    result_text = cstring_to_text(result_str);
    
    /* Free the Go-allocated string */
    free(result_str);
    
    PG_RETURN_TEXT_P(result_text);
}

/*
 * Analyze Turkish word morphologically
 */
PG_FUNCTION_INFO_V1(zemberek_analyze);

Datum
zemberek_analyze(PG_FUNCTION_ARGS)
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
    result_str = AnalyzeTurkish(input_str);
    
    /* Convert result back to PostgreSQL text */
    result_text = cstring_to_text(result_str);
    
    /* Free the Go-allocated string */
    free(result_str);
    
    PG_RETURN_TEXT_P(result_text);
}

/*
 * Extract stem from Turkish word
 */
PG_FUNCTION_INFO_V1(zemberek_stem);

Datum
zemberek_stem(PG_FUNCTION_ARGS)
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
    result_str = StemTurkish(input_str);
    
    /* Convert result back to PostgreSQL text */
    result_text = cstring_to_text(result_str);
    
    /* Free the Go-allocated string */
    free(result_str);
    
    PG_RETURN_TEXT_P(result_text);
}

/*
 * Check if Turkish word has morphological analysis
 */
PG_FUNCTION_INFO_V1(zemberek_has_analysis);

Datum
zemberek_has_analysis(PG_FUNCTION_ARGS)
{
    text *input_text;
    char *input_str;
    int result;
    
    /* Get the input argument */
    input_text = PG_GETARG_TEXT_PP(0);
    
    /* Convert PostgreSQL text to C string */
    input_str = text_to_cstring(input_text);
    
    /* Call the Go function */
    result = HasTurkishAnalysis(input_str);
    
    PG_RETURN_BOOL(result != 0);
}

