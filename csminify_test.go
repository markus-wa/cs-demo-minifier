package csminify_test

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/alecthomas/jsonschema"
	"github.com/stretchr/testify/assert"
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

	parsedReplay, err = csminify.ToReplay(f, 0.5)
	if err != nil {
		panic(err.Error())
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

var update = flag.Bool("updateDocs", false, "update schema documentation")

const jsonSchemaFile = "schema.json"

func TestDoc(t *testing.T) {
	schema, err := json.MarshalIndent(jsonschema.Reflect(&rep.Replay{}), "", "\t")
	assert.NoError(t, err)

	if *update {
		f, err := os.OpenFile(jsonSchemaFile, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0755)
		assert.NoError(t, err)

		_, err = f.Write(schema)
		assert.NoError(t, err)
	} else {
		f, err := os.Open(jsonSchemaFile)
		assert.NoError(t, err)

		b, err := ioutil.ReadAll(f)
		assert.NoError(t, err)

		assert.Equal(t, b, schema)
	}
}
