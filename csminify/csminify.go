package csminify

import (
	"bufio"
	"bytes"
	dem "github.com/markus-wa/demoinfocs-golang"
	"github.com/markus-wa/demoinfocs-golang/events"
	"io"
	"reflect"
	"strconv"
)

type minifier struct {
	writer *bufio.Writer
	parser *dem.Parser
}

func (m *minifier) matchStarted(interface{}) {
	m.parser.EventDispatcher().RegisterHandler(reflect.TypeOf((*events.TickDoneEvent)(nil)).Elem(), m.tickDone)
}

func (m *minifier) tickDone(interface{}) {
	if m.parser.CurrentTick()&15 == 0 {
		m.writer.WriteString(strconv.Itoa(m.parser.CurrentTick()) + " - ")
	}
}

func Minify(r io.Reader) []byte {
	buf := bytes.Buffer{}
	MinifyTo(r, bufio.NewWriter(&buf))
	return buf.Bytes()
}

func MinifyTo(r io.Reader, w io.Writer) {
	p := dem.NewParser(r)
	p.ParseHeader()

	m := minifier{writer: bufio.NewWriter(w), parser: p}

	p.EventDispatcher().RegisterHandler(reflect.TypeOf((*events.MatchStartedEvent)(nil)).Elem(), m.matchStarted)

	p.ParseToEnd(nil)
}
