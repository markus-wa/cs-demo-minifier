# cs-demo-minifier

	TODO

[![GoDoc](https://godoc.org/github.com/markus-wa/cs-demo-minifier?status.svg)](https://godoc.org/github.com/markus-wa/cs-demo-minifier)
[![License](https://img.shields.io/badge/license-MIT-blue.svg?style=flat)](LICENSE.md)

## Go Get

	# Library
	go get github.com/markus-wa/cs-demo-minifier

	# Command line tool
	go get github.com/markus-wa/cs-demo-minifier/cmd/csminify

## Usage

	TODO

## Development

### Running tests

To run tests [Git LFS](https://git-lfs.github.com) is required.

```sh
git submodule init
git submodule update
pushd test/cs-demos && git lfs pull && popd
go test ./...
```

### Generating protobuf code

Should you need to re-generate the protobuf generated code in the `protobuf` package, you will need the following tools:

- The latest protobuf generator (`protoc`) from your package manager or https://github.com/google/protobuf/releases

- And `protoc-gen-gogofaster` from [gogoprotobuf](https://github.com/gogo/protobuf) to generate code for go.

		go get github.com/gogo/protobuf/protoc-gen-gogofaster

Make sure both are inside your `PATH` variable.

After installing these use `go generate ./protobuf` to generate the protobuf code.
