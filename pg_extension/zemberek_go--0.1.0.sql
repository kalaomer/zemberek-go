-- complain if script is sourced in psql, rather than via CREATE EXTENSION
\echo Use "CREATE EXTENSION zemberek_go" to load this file. \quit

-- Normalize informal Turkish text to formal Turkish
CREATE OR REPLACE FUNCTION zemberek_normalize(text TEXT)
RETURNS TEXT
AS 'MODULE_PATHNAME', 'zemberek_normalize'
LANGUAGE C STRICT;

COMMENT ON FUNCTION zemberek_normalize(TEXT) IS 'Normalizes informal Turkish text to formal Turkish (e.g., "gitmişim" -> "gitmişim")';

-- Perform morphological analysis on Turkish words
CREATE OR REPLACE FUNCTION zemberek_analyze(text TEXT)
RETURNS TEXT
AS 'MODULE_PATHNAME', 'zemberek_analyze'

LANGUAGE C STRICT;

COMMENT ON FUNCTION zemberek_analyze(TEXT) IS 'Performs morphological analysis on a Turkish word and returns possible analyses';

-- Extract the stem from a Turkish word
CREATE OR REPLACE FUNCTION zemberek_stem(text TEXT)
RETURNS TEXT
AS 'MODULE_PATHNAME', 'zemberek_stem'
LANGUAGE C STRICT;

COMMENT ON FUNCTION zemberek_stem(TEXT) IS 'Extracts the stem (root) from a Turkish word';

-- Check if a Turkish word has valid morphological analysis
CREATE OR REPLACE FUNCTION zemberek_has_analysis(text TEXT)
RETURNS BOOLEAN
AS 'MODULE_PATHNAME', 'zemberek_has_analysis'
LANGUAGE C STRICT;

COMMENT ON FUNCTION zemberek_has_analysis(TEXT) IS 'Returns true if the Turkish word has valid morphological analysis';

