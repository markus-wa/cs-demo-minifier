package main_test

import (
	"fmt"
	min "github.com/markus-wa/cs-demo-minifier"
	"os"
	"testing"
)

var demPath = "test/cs-demos/default.dem"
var outDir = "test/results"
var outPath = outDir + "/demo.min"

func init() {
	if _, err := os.Stat(demPath); err != nil {
		panic(fmt.Sprintf("Can't read test demo %q", demPath))
	}
	if err := os.MkdirAll(outDir, 0777); err != nil {
		panic(fmt.Sprintf("Couldn't create output dir %q", outDir))
	}
}

func TestStdInOut(t *testing.T) {
	var f *os.File
	var err error

	f, err = os.Open(demPath)
	os.Stdin = f
	if err != nil {
		panic(err.Error())
	}

	out := outPath + ".stdout"
	f, err = os.Create(out)
	if err != nil {
		panic(err.Error())
	}

	stdOut := os.Stdout
	os.Stdout = f
	min.Minify([]string{})
	os.Stdout = stdOut

	assertOutFileCreated(out, t)
}

func TestInOut(t *testing.T) {
	out := outPath + ".out"
	min.Minify([]string{"-demo", demPath, "-out", out})
	assertOutFileCreated(out, t)
}

func TestFreq(t *testing.T) {
	min.Minify([]string{"-demo", demPath, "-freq", "0.2", "-out", os.TempDir() + "/demo-freq.out"})
}

func TestJson(t *testing.T) {
	testProtocol("json", ".json", t)
}

func TestMsgpack(t *testing.T) {
	testProtocol("msgpack", ".mp", t)
}

func TestProtobuf(t *testing.T) {
	testProtocol("protobuf", ".pb", t)
}

func testProtocol(protocol string, suffix string, t *testing.T) {
	min.Minify([]string{"-demo", demPath, "-protocol", protocol, "-out", outPath + suffix})
	assertOutFileCreated(outPath+suffix, t)
}

func assertOutFileCreated(path string, t *testing.T) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		t.Fatalf("Output file %s not created", path)
	}
}
