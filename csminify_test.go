package csminify_test

import (
	"fmt"
	"os"
	"testing"

	csminify "github.com/markus-wa/cs-demo-minifier"
	rep "github.com/markus-wa/cs-demo-minifier/replay"
	nondefaultrep "github.com/markus-wa/cs-demo-minifier/replay/nondefault"
)

var (
	demPath      = "test/cs-demos/default.dem"
	chatDemoPath = "test/cs-demos/2017-05-17-ECSSeason3NA-liquid-vs-renegades-cobblestone.dem"
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
