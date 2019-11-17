# cs-demo-minifier

This tool and library aims to provide a way of converting CS:GO demos into a more easily digestible format while decreasing the data size and retaining all important information. It is based on the demo parser [demoinfocs-golang](https://github.com/markus-wa/demoinfocs-golang).
	
The project is still under development and the data formats may change in backwards-incompatible ways without notice.

[![GoDoc](https://godoc.org/github.com/markus-wa/cs-demo-minifier?status.svg)](https://godoc.org/github.com/markus-wa/cs-demo-minifier)
[![Build Status](https://travis-ci.org/markus-wa/cs-demo-minifier.svg?branch=master)](https://travis-ci.org/markus-wa/cs-demo-minifier)
[![codecov](https://codecov.io/gh/markus-wa/cs-demo-minifier/branch/master/graph/badge.svg)](https://codecov.io/gh/markus-wa/cs-demo-minifier)
[![Go Report](https://goreportcard.com/badge/github.com/markus-wa/cs-demo-minifier)](https://goreportcard.com/report/github.com/markus-wa/cs-demo-minifier)
[![License](https://img.shields.io/badge/license-MIT-blue.svg?style=flat)](LICENSE.md)
[![FOSSA Status](https://app.fossa.io/api/projects/git%2Bgithub.com%2Fmarkus-wa%2Fcs-demo-minifier.svg?type=shield)](https://app.fossa.io/projects/git%2Bgithub.com%2Fmarkus-wa%2Fcs-demo-minifier?ref=badge_shield)


## Discussions / Chat

Use gitter to ask questions and discuss ideas about this project.<br>
There are also [some other rooms](https://gitter.im/csgodemos) available around the topic of CS:GO demos.

[![Gitter chat](https://badges.gitter.im/csgodemos/demo-minifier.png)](https://gitter.im/csgodemos/demo-minifier)


## Download

### Binaries

See the [releases](https://github.com/markus-wa/cs-demo-minifier/releases) page for pre-compiled executables for Linux, MacOS & Windows.

### Go Get

	# Library
	go get -u github.com/markus-wa/cs-demo-minifier

	# Command line tool
	go get -u github.com/markus-wa/cs-demo-minifier/cmd/csminify


## Usage

### Command Line

The following command takes one snapshot of a demo every two seconds (`-freq 0.5`) and saves the resulting replay in the `MessagePack` format to `demo.mp`.

	csminify -demo /path/to/demo.dem -format msgpack -freq 0.5 -out demo.mp

#### Options

```
$ csminify -help
Usage of csminify:
  -demo path
        Demo file path (default stdin)
  -format string
        Format into which the demo should me minified [json, msgpack, protobuf] (default "json")
  -freq float
        Snapshot frequency - per second (default 0.5)
  -out path
        Output file path (default stdout)
```


### Supported Formats

| Format | Command Line (`-format` Flag) | Structure | Minimal Example | Full Example |
| --- | --- | --- | --- | --- |
| JSON | `json` | [schema.json](schema.json) | [minimal.json](examples/minimal.json) | see [releases](https://github.com/markus-wa/cs-demo-minifier/releases) page |
| MessagePack | `msgpack`, `mp` | [schema.json](schema.json) | [minimal.mp](examples/minimal.mp) | see [releases](https://github.com/markus-wa/cs-demo-minifier/releases) page |
| Protocol Buffers | `protobuf`, `proto`, `pb` | [replay.proto](protobuf/gen/proto/replay.proto) | [minimal.pb](examples/minimal.pb) | see [releases](https://github.com/markus-wa/cs-demo-minifier/releases) page |

The minimal examples contain an extract of a demo with each event being included at least once.
Events and attributes are also are documented in [events.md](events.md).

More formats can be added programmatically by implementing the `ReplayMarshaller` interface.

If you would like to see additional formats supported please open a feature request (issue) or a pull request if you already have an implementation ready.

### Unix pipes and `jq`

As the CLI supports Unix pipes, you can combine it with other tools such as [`jq`](https://stedolan.github.io/jq/).

#### Examples

Get the map name of a demo:
```
$ csminfy < demo.dem | jq -r '.header.map'
de_cache
```

Select the first three kills:
```
$ csminify < test/cs-demos/default.dem | jq -r '[ .ticks[] as $parent |
  $parent.events[] | select(.name=="kill") as $kill |
  $kill.attrs[] | select(.key=="victim") as $victim |
  $kill.attrs[] | select(.key=="killer") as $killer |
  $kill.attrs[] | select(.key=="weapon") as $weapon |
  {
    tick: $parent.nr,
    kill: { victim: $victim.numVal, killer: $killer.numVal, weapon: $weapon.numVal }
  }] | .[0:3]'

[
  {
    "tick": 43,
    "kill": {
      "victim": 9,
      "killer": 2,
      "weapon": 303
    }
  },
  {
    "tick": 1029,
    "kill": {
      "victim": 7,
      "killer": 4,
      "weapon": 9
    }
  },
  {
    "tick": 1057,
    "kill": {
      "victim": 11,
      "killer": 4,
      "weapon": 9
    }
  }
]
```


### Library

This is an example on how to minify a demo to JSON and decode it to a `replay.Replay` again.

```go
package main

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

		go get -u github.com/gogo/protobuf/protoc-gen-gogofaster

Make sure both are inside your `PATH` variable.

After installing these use `go generate ./protobuf` to generate the protobuf code.


## License

This project is licensed under the [MIT license](LICENSE.md).

[![FOSSA Status](https://app.fossa.io/api/projects/git%2Bgithub.com%2Fmarkus-wa%2Fcs-demo-minifier.svg?type=large)](https://app.fossa.io/projects/git%2Bgithub.com%2Fmarkus-wa%2Fcs-demo-minifier?ref=badge_large)
