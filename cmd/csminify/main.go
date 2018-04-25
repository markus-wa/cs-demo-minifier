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
	fl.Usage = func() {
		fmt.Fprintln(os.Stderr, "Usage of csminify:")
		fl.PrintDefaults()
		fmt.Fprintln(os.Stderr)
		fmt.Fprintln(os.Stderr, "Direct bug reports and feature requests to https://github.com/markus-wa/cs-demo-minifier")
	}

	formatPtr := fl.String("format", "json", "Format into which the demo should me minified")
	freqPtr := fl.Float64("freq", 0.5, "Snapshot frequency - per second")
	demPathPtr := fl.String("demo", "", "Demo file `path` (default stdin)")
	outPathPtr := fl.String("out", "", "Output file `path` (default stdout)")

	err := fl.Parse(os.Args[1:])
	if err != nil {
		// Some parsing problem, the flag.Parse() already prints the error to stderr
		return
	}

	format := *formatPtr
	freq := float32(*freqPtr)
	demPath := *demPathPtr
	outPath := *outPathPtr

	err = minify(demPath, freq, format, outPath)
	if err != nil {
		panic(err)
	}
}

func minify(demPath string, freq float32, format string, outPath string) error {
	var marshaller min.ReplayMarshaller
	switch format {
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
		fmt.Fprintf(os.Stderr, "Format '%s' unknown, known formats are 'json', 'msgpack' & 'protobuf'\n", format)
		os.Exit(1)
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
			return err
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
			return err
		}
	}

	return min.MinifyTo(in, freq, marshaller, out)
}
