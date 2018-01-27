# cs-demo-minifier

This tool and library aims to provide a way of converting CS:GO demos into a more easily digestible format while decreasing the data size and retaining all important information.
	
It is still under development and the data formats may change in backwards-incompatible ways without notice.

[![GoDoc](https://godoc.org/github.com/markus-wa/cs-demo-minifier?status.svg)](https://godoc.org/github.com/markus-wa/cs-demo-minifier)
[![Build Status](https://travis-ci.org/markus-wa/cs-demo-minifier.svg?branch=master)](https://travis-ci.org/markus-wa/cs-demo-minifier)
[![codecov](https://codecov.io/gh/markus-wa/cs-demo-minifier/branch/master/graph/badge.svg)](https://codecov.io/gh/markus-wa/cs-demo-minifier)
[![License](https://img.shields.io/badge/license-MIT-blue.svg?style=flat)](LICENSE.md)

## Go Get

	# Library
	go get github.com/markus-wa/cs-demo-minifier

	# Command line tool
	go get github.com/markus-wa/cs-demo-minifier/cmd/csminify

## Usage

### Command Line

The following command takes one snapshot of a demo every two seconds (`-freq 0.5`) and saves the resulting, replay in the `MessagePack` format to `demo.mp`.

	csminify -demo /path/to/demo.dem -format msgpack -freq 0.5 -out demo.mp

### Library

This is an example on how to minify a demo to JSON and decode it to a `replay.Replay` again.

```go
import (
	"bytes"
	"encoding/json"
	"log"
	"os"

	csminify "github.com/markus-wa/cs-demo-minifier"
	rep "github.com/markus-wa/cs-demo-minifier/replay"
)

func main() {
	// Open the demo file
	f, err := os.Open("/path/to/demo.dem")
	defer f.Close()
	if err != nil {
		log.Fatal(err)
	}

	// Minify the replay to a byte buffer - or any other io.Writer (JSON)
	// Take 0.5 snapshots per second (=one every two seconds)
	freq := 0.5
	buf := new(bytes.Buffer)
	err = csminify.MinifyTo(f, freq, marshalJSON, buf)
	if err != nil {
		log.Fatal(err)
	}

	// Decoding it again is just as easy
	var r rep.Replay
	err = json.NewDecoder(buf).Decode(&r)
	if err != nil {
		log.Fatal(err)
	}
}

// JSON marshaller
func marshalJSON(r rep.Replay, w io.Writer) error {
	return json.NewEncoder(w).Encode(r)
}
```

MessagePack marshalling works pretty much the same way as JSON.<br>
For Protobuf use `protobuf.Unmarshal()` (in the sub-package).

## Supported Formats

Format | Command Line (`-format` Flag)
--- | --- | ---
JSON | `json`
MessagePack | `msgpack`, `mp`
Protocol Buffers | `protobuf`, `proto`, `pb`

More formats can be added programmatically by implementing the `ReplayMarshaller` interface.

If you would like to see additional formats supported please open a feature request (issue) or a pull request if you already have an implementation ready.

## Development

### Running Tests

To run tests [Git LFS](https://git-lfs.github.com) is required.

```sh
git submodule init
git submodule update
pushd test/cs-demos && git lfs pull && popd
go test ./...
```

### Generating Protobuf Code

Should you need to re-generate the protobuf generated code in the `protobuf` package, you will need the following tools:

- The latest protobuf generator (`protoc`) from your package manager or https://github.com/google/protobuf/releases

- And `protoc-gen-gogofaster` from [gogoprotobuf](https://github.com/gogo/protobuf) to generate code for go.

		go get github.com/gogo/protobuf/protoc-gen-gogofaster

Make sure both are inside your `PATH` variable.

After installing these use `go generate ./protobuf` to generate the protobuf code.
