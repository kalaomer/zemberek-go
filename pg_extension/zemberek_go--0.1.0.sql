-- complain if script is sourced in psql, rather than via CREATE EXTENSION
\echo Use "CREATE EXTENSION zemberek_go" to load this file. \quit

-- Dummy UDF that demonstrates goroutine usage
CREATE OR REPLACE FUNCTION zemberek_go_hello(name TEXT)
RETURNS TEXT
AS 'MODULE_PATHNAME', 'zemberek_go_hello_wrapper'
LANGUAGE C STRICT;

COMMENT ON FUNCTION zemberek_go_hello(TEXT) IS 'A demo function that uses Go goroutines to process the input';
