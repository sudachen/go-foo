module github.com/sudachen/go-foo

replace github.com/sudachen/go-iokit => ./go-iokit

go 1.13

require (
	github.com/google/go-cmp v0.4.0 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	golang.org/x/xerrors v0.0.0-20191204190536-9bdfabe68543
	gotest.tools v2.2.0+incompatible
)
