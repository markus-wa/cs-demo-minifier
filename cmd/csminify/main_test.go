package main

import (
	"fmt"
	"os"
	"testing"
)

var (
	demPath = "../../test/cs-demos/default.dem"
	outDir  = "../../test/results"
	outPath = outDir + "/demo.min"
)

func TestMain(m *testing.M) {
	// Check if test demo exists
	if _, err := os.Stat(demPath); err != nil {
		panic(fmt.Sprintf("Can't read test demo %q: %s", demPath, err.Error()))
	}
	// Create test result folder if it doesn't exist, somehow 0666 isn't enough for travis??
	if err := os.MkdirAll(outDir, 0777); err != nil {
		panic(fmt.Sprintf("Couldn't create output dir %q: %s", outDir, err.Error()))
	}
	os.Exit(m.Run())
}

func TestStdInOut(t *testing.T) {
	oldStdin := os.Stdin
	defer func() { os.Stdin = oldStdin }()

	f, err := os.Open(demPath)
	os.Stdin = f
	if err != nil {
		t.Fatal(err)
	}

	out := outPath + ".stdout"
	f, err = os.Create(out)
	if err != nil {
		t.Fatal(err)
	}

	stdOut := os.Stdout
	os.Stdout = f
	runMainWithArgs([]string{})
	os.Stdout = stdOut

	assertOutFileCreated(out, t)
}

func TestInOut(t *testing.T) {
	out := outPath + ".out"
	runMainWithArgs([]string{"-demo", demPath, "-out", out})
	assertOutFileCreated(out, t)
}

func TestFreq(t *testing.T) {
	runMainWithArgs([]string{"-demo", demPath, "-freq", "0.2", "-out", os.TempDir() + "/demo-freq.out"})
}

func TestJSON(t *testing.T) {
	testFormat("json", ".json", t)
}

func TestMsgPack(t *testing.T) {
	testFormat("msgpack", ".mp", t)
}

func TestProtobuf(t *testing.T) {
	testFormat("protobuf", ".pb", t)
}

func testFormat(format string, suffix string, t *testing.T) {
	runMainWithArgs([]string{"-demo", demPath, "-format", format, "-out", outPath + suffix})
	assertOutFileCreated(outPath+suffix, t)
}

func runMainWithArgs(args []string) {
	// Store original arguments and restore them at the end
	oldArgs := os.Args
	defer func() { os.Args = oldArgs }()

	// Add dummy command name to args
	args = append([]string{"cmd"}, args...)
	os.Args = args
	main()
}

func assertOutFileCreated(path string, t *testing.T) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		t.Fatalf("Output file %s not created", path)
	}
}
