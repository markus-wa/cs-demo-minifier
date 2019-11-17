package csminify_test

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/alecthomas/jsonschema"
	"github.com/markus-wa/cs-demo-minifier/protobuf"
	"github.com/stretchr/testify/assert"
	"gopkg.in/vmihailenco/msgpack.v2"
	"io/ioutil"
	"os"
	"strings"
	"testing"

	csminify "github.com/markus-wa/cs-demo-minifier"
	rep "github.com/markus-wa/cs-demo-minifier/replay"
	nondefaultrep "github.com/markus-wa/cs-demo-minifier/replay/nondefault"
)

var (
	demPath      = "test/cs-demos/default.dem"
	chatDemoPath = "test/cs-demos/set/2017-05-17-ECSSeason3NA-liquid-vs-renegades-cobblestone.dem"
)

var nonDefaultReplay, parsedReplay, minimalExampleReplay rep.Replay

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

	parsedReplay, err = csminify.ToReplay(f, 0.5)
	if err != nil {
		panic(err.Error())
	}

	minimalExampleReplay = exampleReplay()
}

func exampleReplay() rep.Replay {
	var distinctEventTicks []rep.Tick

	knownEvents := make(map[string]struct{})
	for i := range parsedReplay.Ticks {
		for j := range parsedReplay.Ticks[i].Events {
			eventName := parsedReplay.Ticks[i].Events[j].Name
			if _, alreadyKnown := knownEvents[eventName]; !alreadyKnown {
				knownEvents[eventName] = struct{}{}
				distinctEventTicks = append(distinctEventTicks, parsedReplay.Ticks[i])
				break
			}
		}
	}

	return rep.Replay{
		Header:    parsedReplay.Header,
		Entities:  parsedReplay.Entities,
		Snapshots: parsedReplay.Snapshots[0:2],
		Ticks:     distinctEventTicks,
	}
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

func TestChat(t *testing.T) {
	f, err := os.Open(chatDemoPath)
	defer f.Close()
	if err != nil {
		t.Fatal(err)
	}

	r, err := csminify.ToReplay(f, 0.2)
	if err != nil {
		t.Fatal(err)
	}

	ok := false
	for _, t := range r.Ticks {
		for _, e := range t.Events {
			if e.Name == rep.EventChatMessage {
				for _, a := range e.Attributes {
					if a.Key == rep.AttrKindText && len(a.StrVal) > 0 {
						ok = true
					}
				}
			}
		}
	}
	if !ok {
		t.Fatal("No chat events recorded when there should have been some")
	}
}

func TestExtraHandlers(t *testing.T) {
	f, err := os.Open(demPath)
	defer f.Close()
	if err != nil {
		t.Fatal(err)
	}

	// Register extra handlers
	ec := new(csminify.EventCollector)
	csminify.EventHandlers.Extra.RegisterAll(ec)

	r, err := csminify.ToReplayWithConfig(f, csminify.ReplayConfig{SnapshotFrequency: 0.2, EventCollector: ec})
	if err != nil {
		t.Fatal(err)
	}

	// Check for extra events
	// Currently only footsteps
	ok := false
	for _, t := range r.Ticks {
		for _, e := range t.Events {
			if e.Name == rep.EventFootstep {
				for _, a := range e.Attributes {
					if a.Key == rep.AttrKindEntityID && a.NumVal > 0 {
						ok = true
					}
				}
			}
		}
	}
	if !ok {
		t.Fatal("No footstep events recorded when there should have been some")
	}
}

func TestDemoSet(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test due to -short flag")
	}

	demSetPath := "test/cs-demos/set"
	dems, err := ioutil.ReadDir(demSetPath)
	if err != nil {
		t.Fatal(err)
	}

	for _, d := range dems {
		name := d.Name()
		if strings.HasSuffix(name, ".dem") {
			fmt.Printf("Parsing '%s/%s'\n", demSetPath, name)
			func() {
				var f *os.File
				f, err = os.Open(demSetPath + "/" + name)
				if err != nil {
					t.Error(err)
				}
				defer f.Close()

				defer func() {
					pErr := recover()
					if pErr != nil {
						t.Errorf("Parsing of '%s/%s' paniced: %s\n", demSetPath, name, pErr)
					}
				}()

				b, err := csminify.Minify(f, 0.2, marshalJSON)
				if err != nil {
					t.Fatal(err)
				}
				if len(b) == 0 {
					t.Fatal("Resulting []byte of minification is empty")
				}
			}()
		}
	}
}

var updateGolden = flag.Bool("updateGolden", false, "update schema documentation")

const jsonSchemaFile = "schema.json"

func TestDoc(t *testing.T) {
	schema, err := json.MarshalIndent(jsonschema.Reflect(&rep.Replay{}), "", "\t")
	assert.NoError(t, err)

	if *updateGolden {
		updateGoldenFile(t, jsonSchemaFile, schema)
	} else {
		b := readFile(t, jsonSchemaFile)

		assert.Equal(t, b, schema)
	}
}

func TestExample_Json_Minimal(t *testing.T) {
	data, err := json.MarshalIndent(&minimalExampleReplay, "", "\t")
	assert.NoError(t, err)

	testGoldenOrUpdate(t, "examples/minimal.json", data, unmarshalJSON)
}

func TestExample_MsgPack_Minimal(t *testing.T) {
	data, err := msgpack.Marshal(&minimalExampleReplay)
	assert.NoError(t, err)

	testGoldenOrUpdate(t, "examples/minimal.mp", data, unmarshalMsgPack)
}

func TestExample_Protobuf_Minimal(t *testing.T) {
	data := bytes.Buffer{}
	err := protobuf.MarshalReplay(minimalExampleReplay, &data)
	assert.NoError(t, err)

	testGoldenOrUpdate(t, "examples/minimal.pb", data.Bytes(), protobuf.UnmarshalReplay)
}

func testGoldenOrUpdate(t *testing.T, fileName string, data []byte, unmarshaller replayUnmarshaller) {
	if *updateGolden {
		updateGoldenFile(t, fileName, data)
	} else {
		goldenBytes := readFile(t, fileName)

		// just check if it's the same length as the order of the elements might be different
		assert.Len(t, data, len(goldenBytes))
	}
}

func readFile(t *testing.T, fileName string) []byte {
	f, err := os.Open(fileName)
	assert.NoError(t, err)

	b, err := ioutil.ReadAll(f)
	assert.NoError(t, err)

	return b
}

func updateGoldenFile(t *testing.T, fileName string, bytes []byte) {
	f, err := os.OpenFile(fileName, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0755)
	assert.NoError(t, err)

	_, err = f.Write(bytes)
	assert.NoError(t, err)
}
