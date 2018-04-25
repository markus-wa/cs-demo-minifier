// Package csminify provides functions for parsing CS:GO demos and minifying them into various formats.
package csminify

import (
	"bufio"
	"bytes"
	"io"
	"math"

	r3 "github.com/golang/geo/r3"
	dem "github.com/markus-wa/demoinfocs-golang"
	events "github.com/markus-wa/demoinfocs-golang/events"

	rep "github.com/markus-wa/cs-demo-minifier/replay"
)

// ReplayMarshaller is the signature for functions that serialize replay.Replay structs to an io.Writer
type ReplayMarshaller func(rep.Replay, io.Writer) error

// Minify wraps MinifyTo with a bytes.Buffer and returns the written bytes.
func Minify(r io.Reader, snapFreq float32, marshal ReplayMarshaller) ([]byte, error) {
	var buf bytes.Buffer
	err := MinifyTo(r, snapFreq, marshal, bufio.NewWriter(&buf))
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// MinifyTo reads a demo from r, creates a replay and marshals it to w.
// See also: ToReplay
func MinifyTo(r io.Reader, snapFreq float32, marshal ReplayMarshaller, w io.Writer) error {
	replay, err := ToReplay(r, snapFreq)
	if err != nil {
		return err
	}

	err = marshal(replay, w)
	return err
}

// ToReplay reads a demo from r, takes snapshots (snapFreq/sec) and records events into a Replay.
func ToReplay(r io.Reader, snapFreq float32) (rep.Replay, error) {
	ec := new(EventCollector)
	EventHandlers.Default.RegisterAll(ec)
	return ToReplayWithCustomEvents(r, snapFreq, ec)
}

// ToReplayWithCustomEvents is like ToReplay but with a custom EventCollector.
func ToReplayWithCustomEvents(r io.Reader, snapFreq float32, eventCollector *EventCollector) (rep.Replay, error) {
	// FIXME: Smoothify flag
	// TODO: Maybe pass a WarnHandler along
	p := dem.NewParser(r, nil)
	_, err := p.ParseHeader()
	if err != nil {
		return rep.Replay{}, err
	}

	// Make the parser accessible for the custom event handlers
	eventCollector.parser = p

	m := minifier{parser: p, eventCollector: eventCollector}

	m.replay.Header.MapName = p.Map()
	m.replay.Header.TickRate = p.FrameRate()
	m.replay.Header.SnapshotRate = int(math.Round(float64(m.replay.Header.TickRate / snapFreq)))

	// Register event handlers from collector
	for _, h := range eventCollector.handlers {
		m.parser.RegisterEventHandler(h)
	}

	m.parser.RegisterEventHandler(m.tickDone)

	err = p.ParseToEnd()
	if err != nil {
		return rep.Replay{}, err
	}

	// TODO: There's probably a better place for this
	for _, pl := range m.parser.Participants() {
		ent := rep.Entity{
			ID:    pl.EntityID,
			Team:  int(pl.Team),
			Name:  pl.Name,
			IsNpc: pl.IsBot,
		}

		m.replay.Entities = append(m.replay.Entities, ent)
	}

	return m.replay, nil
}

type minifier struct {
	parser         *dem.Parser
	replay         rep.Replay
	eventCollector *EventCollector
}

func (m *minifier) tickDone(e events.TickDoneEvent) {
	tick := m.parser.CurrentFrame()
	// Is it snapshot o'clock?
	if tick%m.replay.Header.SnapshotRate == 0 {
		snap := rep.Snapshot{
			Tick: tick,
		}

		for _, pl := range m.parser.PlayingParticipants() {
			if pl.IsAlive() {
				e := rep.EntityUpdate{
					EntityID:      pl.EntityID,
					Hp:            pl.Hp,
					Armor:         pl.Armor,
					FlashDuration: float32(roundTo(float64(pl.FlashDuration), 0.1)), // Round to nearest 0.1 sec - saves space in JSON
					Positions:     []rep.Point{r3VectorToPoint(pl.Position)},
					Angle:         int(pl.ViewDirectionX),
				}
				// FIXME: Smoothify
				snap.EntityUpdates = append(snap.EntityUpdates, e)
			}
		}

		m.replay.Snapshots = append(m.replay.Snapshots, snap)
	}

	// Did we collect any events in this tick?
	if len(m.eventCollector.events) > 0 {
		tickEvents := make([]rep.Event, len(m.eventCollector.events))
		copy(tickEvents, m.eventCollector.events)
		m.replay.Ticks = append(m.replay.Ticks, rep.Tick{
			Nr:     tick,
			Events: tickEvents,
		})
		// Clear events for next tick
		m.eventCollector.events = m.eventCollector.events[:0]
	}
}

func r3VectorToPoint(v r3.Vector) rep.Point {
	return rep.Point{X: int(v.X), Y: int(v.Y)}
}

// roundTo wraps math.Round and allows specifying the rounding precision.
func roundTo(x, precision float64) float64 {
	return math.Round(x/precision) * precision
}
