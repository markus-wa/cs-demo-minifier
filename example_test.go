package csminify_test

import (
	"bytes"
	"encoding/json"
	"os"
	"testing"

	csminify "github.com/markus-wa/cs-demo-minifier"
	rep "github.com/markus-wa/cs-demo-minifier/replay"
)

// Make sure the example from the README.md compiles and runs
func TestExample(t *testing.T) {
	// Open the demo file
	f, err := os.Open(demPath)
	defer f.Close()
	if err != nil {
		t.Fatal(err)
	}

	// Minify the replay to a byte buffer (JSON)
	buf := new(bytes.Buffer)
	err = csminify.MinifyTo(f, 0.5, marshalJSON, buf)
	if err != nil {
		t.Fatal(err)
	}

	// Decoding the it again is just as easy
	var r rep.Replay
	err = json.NewDecoder(buf).Decode(&r)
	if err != nil {
		t.Fatal(err)
	}
}
