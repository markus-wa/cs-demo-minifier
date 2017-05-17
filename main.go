package main

import (
	"flag"
	"fmt"
	min "github.com/markus-wa/cs-demo-minifier/csminify"
	"github.com/markus-wa/cs-demo-minifier/csminify/json"
	"io"
	"os"
)

func main() {
	p := flag.String("protocol", "json", "Set protocol")
	demPath := flag.String("demo", "", "Demo file")
	outPath := flag.String("out", "", "Demo file")
	flag.Parse()

	var marshaller min.ReplayMarshaller
	switch *p {
	case "json":
		marshaller = json.MarshalReplay
	default:
		panic(fmt.Sprintf("Protocol '%s' unknown", p))
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
			panic(err)
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
			panic(err)
		}
	}

	min.MinifyTo(in, 0.5, marshaller, out)
}
