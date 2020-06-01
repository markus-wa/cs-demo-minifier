package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"

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

	formatPtr := fl.String("format", "json", "Format into which the demo should me minified [json, msgpack, protobuf]")
	freqPtr := fl.Float64("freq", 0.5, "Snapshot frequency - per second")
	demPathPtr := fl.String("demo", "", "Demo file `path` (default stdin)")
	outPathPtr := fl.String("out", "", "Output file `path` (default stdout)")
	webserver := fl.Bool("server", false, "When set, the app starts a webserver instead. "+
		"The WebServer accepts all arguments the CLI does as headers. E.g x-freq. For now JSON Response is supported")
	httpPort := fl.Int("port", 8080, "The HTTP Port of the Webserver. Only considered when 'server' is set to true")

	err := fl.Parse(os.Args[1:])
	if err != nil {
		// Some parsing problem, the flag.Parse() already prints the error to stderr
		return
	}

	format := *formatPtr
	freq := *freqPtr
	demPath := *demPathPtr
	outPath := *outPathPtr

	if *webserver == true {
		fmt.Println("orld")
		log.Println("Started a WebServer")
		http.HandleFunc("/", HttpHandler)
		log.Fatal(http.ListenAndServe(":"+strconv.Itoa(*httpPort), nil))
	} else {
		err = minify(demPath, freq, format, outPath)
		if err != nil {
			panic(err)
		}
	}
}

func minify(demPath string, freq float64, format string, outPath string) error {
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

func HttpHandler(w http.ResponseWriter, r *http.Request) {
	freqStr := r.Header.Get("x-freq")
	freq, _ := strconv.ParseFloat(freqStr, 64)

	var fileAsBytes = StreamToByte(r.Body)
	byteReader := bytes.NewReader(fileAsBytes)
	var marshaller min.ReplayMarshaller = func(replay rep.Replay, w io.Writer) error {
		return json.NewEncoder(w).Encode(replay)
	}
	min.MinifyTo(byteReader, freq, marshaller, w)
}

func StreamToByte(stream io.Reader) []byte {
	buf := new(bytes.Buffer)
	buf.ReadFrom(stream)
	return buf.Bytes()
}
