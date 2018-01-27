package csminify_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"reflect"
	"testing"

	diff "github.com/d4l3k/messagediff"
	msgpack "gopkg.in/vmihailenco/msgpack.v2"

	csminify "github.com/markus-wa/cs-demo-minifier"
	protobuf "github.com/markus-wa/cs-demo-minifier/protobuf"
	rep "github.com/markus-wa/cs-demo-minifier/replay"
	nondefaultrep "github.com/markus-wa/cs-demo-minifier/replay/nondefault"
)

var demPath = "test/cs-demos/default.dem"

var nonDefaultReplay, parsedReplay rep.Replay

func TestMain(m *testing.M) {
	nonDefaultReplay = nondefaultrep.GetNonDefaultReplay()

	if _, err := os.Stat(demPath); err != nil {
		panic(fmt.Sprintf("Can't read test demo %q", demPath))
	}

	initParsedReplay()

	os.Exit(m.Run())
}

func initParsedReplay() {
	f, err := os.Open(demPath)
	defer f.Close()
	if err != nil {
		panic(err.Error())
	}

	csminify.ToReplay(f, 0.5, &parsedReplay)
}

func TestMinify(t *testing.T) {
	f, err := os.Open(demPath)
	defer f.Close()
	if err != nil {
		t.Fatal(err)
	}

	b, err := csminify.Minify(f, 0.2, marshalJSON)
	if err != nil {
		t.Fatal(err)
	}
	if len(b) == 0 {
		t.Fatal("Resulting []byte of minification is empty")
	}
}

// Test data preservation of JSON marshalling & unmarshalling with a 'non-default' replay.
func TestJSONNonDefault(t *testing.T) {
	testDataPreservation(nonDefaultReplay, marshalJSON, unmarshalJSON, t)
}

// Test data preservation of JSON marshalling & unmarshalling with a real, parsed replay.
func TestJSONDemo(t *testing.T) {
	testDataPreservation(parsedReplay, marshalJSON, unmarshalJSON, t)
}

func marshalJSON(r rep.Replay, w io.Writer) error {
	return json.NewEncoder(w).Encode(r)
}

func unmarshalJSON(r io.Reader, rp *rep.Replay) error {
	return json.NewDecoder(r).Decode(rp)
}

// Test data preservation of MessagePack marshalling & unmarshalling with a 'non-default' replay.
func TestMsgPackNonDefault(t *testing.T) {
	testDataPreservation(nonDefaultReplay, marshalMsgPack, unmarshalMsgPack, t)
}

// Test data preservation of MessagePack marshalling & unmarshalling with a real, parsed replay.
func TestMsgPackDemo(t *testing.T) {
	testDataPreservation(parsedReplay, marshalMsgPack, unmarshalMsgPack, t)
}

func marshalMsgPack(r rep.Replay, w io.Writer) error {
	return msgpack.NewEncoder(w).Encode(r)
}

func unmarshalMsgPack(r io.Reader, rp *rep.Replay) error {
	return msgpack.NewDecoder(r).Decode(rp)
}

// Test data preservation of Protobuf marshalling & unmarshalling with a 'non-default' replay.
func TestProtobufNonDefault(t *testing.T) {
	testDataPreservation(nonDefaultReplay, protobuf.MarshalReplay, protobuf.UnmarshalReplay, t)
}

// Test data preservation of Protobuf marshalling & unmarshalling with a real, parsed replay.
func TestProtobufDemo(t *testing.T) {
	testDataPreservation(parsedReplay, protobuf.MarshalReplay, protobuf.UnmarshalReplay, t)
}

type replayUnmarshaller func(io.Reader, *rep.Replay) error

func testDataPreservation(replay rep.Replay, marshal csminify.ReplayMarshaller, unmarshal replayUnmarshaller, t *testing.T) {
	buf := new(bytes.Buffer)
	err := marshal(replay, buf)
	if err != nil {
		t.Fatal(err)
	}

	var r rep.Replay
	err = unmarshal(buf, &r)
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(replay, r) {
		t.Error("Original replay and unmarshalled replay aren't equal")
		diff, eq := diff.PrettyDiff(replay, r)
		if !eq {
			t.Logf("Diff:\n%s", diff)
		} else {
			t.Log("DeepEqual reported inequality but couldn't find differences in the structs")
		}
	}
}
