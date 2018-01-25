package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os"

	msgpack "gopkg.in/vmihailenco/msgpack.v2"

	min "github.com/markus-wa/cs-demo-minifier"
	pb "github.com/markus-wa/cs-demo-minifier/protobuf"
	rep "github.com/markus-wa/cs-demo-minifier/replay"
)

func main() {
	fl := new(flag.FlagSet)
	protPtr := fl.String("protocol", "json", "Protocol to minify the demo to")
	freqPtr := fl.Float64("freq", 0.5, "Snapshot frequency - per second")
	demPathPtr := fl.String("demo", "", "Demo file path")
	outPathPtr := fl.String("out", "", "Output file path")
	fl.Parse(os.Args[1:])

	prot := *protPtr
	freq := float32(*freqPtr)
	demPath := *demPathPtr
	outPath := *outPathPtr

	var marshaller min.ReplayMarshaller
	switch prot {
	case "json":
		marshaller = func(replay rep.Replay, w io.Writer) error {
			return json.NewEncoder(w).Encode(replay)
		}

	case "protobuf":
		fallthrough
	case "proto":
		fallthrough
	case "pb":
		marshaller = pb.MarshalReplay

	case "msgpack":
		fallthrough
	case "mp":
		marshaller = func(rep rep.Replay, w io.Writer) error {
			return msgpack.NewEncoder(w).Encode(rep)
		}

	default:
		panic(fmt.Sprintf("Protocol '%s' unknown", prot))
	}

	var in io.Reader
	switch demPath {
	case "":
		in = os.Stdin
	default:
		f, err := os.Open(demPath)
		defer f.Close()
		in = f
		if err != nil {
			panic(err.Error())
		}
	}

	var out io.Writer
	switch outPath {
	case "":
		out = os.Stdout
	default:
		f, err := os.Create(outPath)
		defer f.Close()
		out = f
		if err != nil {
			panic(err.Error())
		}
	}

	err := min.MinifyTo(in, freq, marshaller, out)
	if err != nil {
		panic(err.Error())
	}
}
