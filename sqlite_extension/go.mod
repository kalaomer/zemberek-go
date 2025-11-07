module github.com/kalaomer/zemberek-go/sqlite_extension

go 1.24.0

require (
	github.com/kalaomer/zemberek-go v0.1.0
	github.com/mattn/go-sqlite3 v1.14.24
)

require google.golang.org/protobuf v1.36.10 // indirect

replace github.com/kalaomer/zemberek-go => ../
