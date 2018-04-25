package csminify_test

import (
	"bytes"
	"encoding/json"
	"io"
	"os"
	"reflect"
	"testing"

	diff "github.com/d4l3k/messagediff"
	events "github.com/markus-wa/demoinfocs-golang/events"
	msgpack "gopkg.in/vmihailenco/msgpack.v2"

	min "github.com/markus-wa/cs-demo-minifier"
	protobuf "github.com/markus-wa/cs-demo-minifier/protobuf"
	rep "github.com/markus-wa/cs-demo-minifier/replay"
)

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

// Test data preservation of Protobuf marshalling & unmarshalling with a custom events & attributes.
func TestProtobufCustomEvents(t *testing.T) {
	f, err := os.Open(demPath)
	defer f.Close()
	if err != nil {
		panic(err.Error())
	}

	ec := new(min.EventCollector)
	ec.AddHandler(func(events.BombPlantedEvent) {
		// Test access to parser
		ec.Parser().IngameTick()

		ec.AddEvent(rep.Event{
			Name:       "custom_plant_event",
			Attributes: []rep.EventAttribute{rep.EventAttribute{Key: "custom_attribute", StrVal: "test"}},
		})
	})

	customReplay, err := min.ToReplayWithCustomEvents(f, 0.5, ec)
	if err != nil {
		panic(err.Error())
	}
	testDataPreservation(customReplay, protobuf.MarshalReplay, protobuf.UnmarshalReplay, t)
}

type replayUnmarshaller func(io.Reader, *rep.Replay) error

func testDataPreservation(replay rep.Replay, marshal min.ReplayMarshaller, unmarshal replayUnmarshaller, t *testing.T) {
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
