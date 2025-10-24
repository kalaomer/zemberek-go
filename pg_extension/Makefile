EXTENSION = zemberek_go
DATA = zemberek_go--0.1.0.sql
MODULE_big = zemberek_go
HEADERS = libzemberek_go.h

# Only C object files here (no .a archives)
OBJS = zemberek_go_wrapper.o 

# Link the generated archive into the PostgreSQL extension
SHLIB_LINK += $(CURDIR)/libzemberek_go.a

# Cleaning
EXTRA_CLEAN = libzemberek_go.a libzemberek_go.h

PG_CONFIG = pg_config
PGXS := $(shell $(PG_CONFIG) --pgxs)
include $(PGXS)

# Build the Go static archive (and header) via c-archive
libzemberek_go.a: zemberek_go.go
	go build -buildmode=c-archive -o libzemberek_go.a zemberek_go.go

# Ensure the wrapper is compiled after the Go archive & uses its header
zemberek_go_wrapper.o: zemberek_go_wrapper.c libzemberek_go.a

# Make the final shared library re-link if the Go archive changed
$(MODULE_big).$(DLSUFFIX): libzemberek_go.a
