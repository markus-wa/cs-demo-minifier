package main

import (
	"encoding/json"
	"flag"
	"fmt"
	min "github.com/markus-wa/cs-demo-minifier/csminify"
	pb "github.com/markus-wa/cs-demo-minifier/csminify/protobuf"
	rep "github.com/markus-wa/cs-demo-minifier/csminify/replay"
	"gopkg.in/vmihailenco/msgpack.v2"
	"io"
	"os"
)

func main() {
	Minify(os.Args[1:])
}

func Minify(args []string) {
	fl := new(flag.FlagSet)
	prot := fl.String("protocol", "json", "Protocol to minify the demo to")
	freq := fl.Float64("freq", 0.5, "Snapshot frequency - per second")
	demPath := fl.String("demo", "", "Demo file path")
	outPath := fl.String("out", "", "Output file path")
	fl.Parse(args)

	var marshaller min.ReplayMarshaller
	switch *prot {
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
	switch *demPath {
	case "":
		in = os.Stdin
	default:
		f, err := os.Open(*demPath)
		defer f.Close()
		in = f
		if err != nil {
			panic(err.Error())
		}
	}

	var out io.Writer
	switch *outPath {
	case "":
		out = os.Stdout
	default:
		f, err := os.Create(*outPath)
		defer f.Close()
		out = f
		if err != nil {
			panic(err.Error())
		}
	}

	min.MinifyTo(in, float32(*freq), marshaller, out)
}
