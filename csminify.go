// Package csminify provides functions for parsing CS:GO demos and minifying them into various formats.
package csminify

import (
	"bufio"
	"bytes"
	"io"
	"math"

	r3 "github.com/golang/geo/r3"
	dem "github.com/markus-wa/demoinfocs-golang/v2/pkg/demoinfocs"
	common "github.com/markus-wa/demoinfocs-golang/v2/pkg/demoinfocs/common"
	events "github.com/markus-wa/demoinfocs-golang/v2/pkg/demoinfocs/events"

	rep "github.com/markus-wa/cs-demo-minifier/replay"
)

// ReplayMarshaller is the signature for functions that serialize replay.Replay structs to an io.Writer
type ReplayMarshaller func(rep.Replay, io.Writer) error

// Minify wraps MinifyTo with a bytes.Buffer and returns the written bytes.
func Minify(r io.Reader, snapFreq float64, marshal ReplayMarshaller) ([]byte, error) {
	return MinifyWithConfig(r, DefaultReplayConfig(snapFreq), marshal)
}

// MinifyWithConfig wraps MinifyToWithConfig with a bytes.Buffer and returns the written bytes.
func MinifyWithConfig(r io.Reader, cfg ReplayConfig, marshal ReplayMarshaller) ([]byte, error) {
	var buf bytes.Buffer
	err := MinifyToWithConfig(r, cfg, marshal, bufio.NewWriter(&buf))
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// MinifyTo reads a demo from r, creates a replay and marshals it to w.
// See also: ToReplay
func MinifyTo(r io.Reader, snapFreq float64, marshal ReplayMarshaller, w io.Writer) error {
	return MinifyToWithConfig(r, DefaultReplayConfig(snapFreq), marshal, w)
}

// MinifyToWithConfig reads a demo from r, creates a replay and marshals it to w.
// See also: ToReplayWithConfig
func MinifyToWithConfig(r io.Reader, cfg ReplayConfig, marshal ReplayMarshaller, w io.Writer) error {
	replay, err := ToReplayWithConfig(r, cfg)
	if err != nil {
		return err
	}

	err = marshal(replay, w)
	return err
}

// DefaultReplayConfig returns the default configuration with a given snapshot frequency.
// May be overridden.
var DefaultReplayConfig = func(snapFreq float64) ReplayConfig {
	ec := new(EventCollector)
	EventHandlers.Default.RegisterAll(ec)
	return ReplayConfig{
		SnapshotFrequency: snapFreq,
		EventCollector:    ec,
	}
}

// ReplayConfig contains the configuration for generating a replay.
type ReplayConfig struct {
	SnapshotFrequency float64
	EventCollector    *EventCollector
	// TODO: Smoothify flag?
}

// ToReplay reads a demo from r, takes snapshots (snapFreq/sec) and records events into a Replay.
func ToReplay(r io.Reader, snapFreq float64) (rep.Replay, error) {
	return ToReplayWithConfig(r, DefaultReplayConfig(snapFreq))
}

// ToReplayWithConfig reads a demo from r, takes snapshots and records events into a Replay with a custom configuration.
func ToReplayWithConfig(r io.Reader, cfg ReplayConfig) (rep.Replay, error) {
	// TODO: Provide a way to pass on warnings to the caller
	p := dem.NewParser(r)
	header, err := p.ParseHeader()
	if err != nil {
		return rep.Replay{}, err
	}

	// Make the parser accessible for the custom event handlers
	cfg.EventCollector.parser = p

	m := newMinifier(p, cfg.EventCollector, cfg.SnapshotFrequency)

	m.replay.Header.MapName = header.MapName
	m.tickRate(p.TickRate())

	p.RegisterEventHandler(func(events.ConVarsUpdated) {
		if tickRate := p.TickRate(); tickRate != 0 {
			m.tickRate(tickRate)
		}
	})

	// Register event handlers from collector
	for _, h := range cfg.EventCollector.handlers {
		m.parser.RegisterEventHandler(h)
	}

	m.parser.RegisterEventHandler(m.frameDone)

	err = p.ParseToEnd()
	if err != nil {
		return rep.Replay{}, err
	}

	return m.replay, nil
}

type minifier struct {
	parser            dem.Parser
	replay            rep.Replay
	eventCollector    *EventCollector
	snapshotFrequency float64

	knownPlayerEntityIDs map[int]struct{}
}

func newMinifier(parser dem.Parser, eventCollector *EventCollector, snapshotFrequency float64) minifier {
	return minifier{
		parser:               parser,
		eventCollector:       eventCollector,
		knownPlayerEntityIDs: make(map[int]struct{}),
		snapshotFrequency:    snapshotFrequency,
	}
}

func (m *minifier) frameDone(e events.FrameDone) {
	tick := m.parser.CurrentFrame()
	// Is it snapshot o'clock?
	if tick%m.replay.Header.SnapshotRate == 0 {
		// TODO: There might be a better way to do this than having updateKnownPlayers() here
		m.updateKnownPlayers()

		snap := m.snapshot()
		m.replay.Snapshots = append(m.replay.Snapshots, snap)
	}

	// Did we collect any events in this frame?
	if len(m.eventCollector.events) > 0 {
		tickEvents := make([]rep.Event, len(m.eventCollector.events))
		copy(tickEvents, m.eventCollector.events)
		m.replay.Ticks = append(m.replay.Ticks, rep.Tick{
			Nr:     tick,
			Events: tickEvents,
		})
		// Clear events for next frame
		m.eventCollector.events = m.eventCollector.events[:0]
	}
}

func (m *minifier) snapshot() rep.Snapshot {
	snap := rep.Snapshot{
		Tick: m.parser.CurrentFrame(),
	}

	for _, pl := range m.parser.GameState().Participants().Playing() {
		if pl.IsAlive() {
			e := rep.EntityUpdate{
				EntityID:      pl.EntityID,
				Hp:            pl.Health(),
				Armor:         pl.Armor(),
				FlashDuration: float32(roundTo(float64(pl.FlashDuration), 0.1)), // Round to nearest 0.1 sec - saves space in JSON
				Positions:     []rep.Point{r3VectorToPoint(pl.Position())},
				AngleX:        int(pl.ViewDirectionX()),
				AngleY:        int(pl.ViewDirectionY()),
				HasHelmet:     pl.HasHelmet(),
				HasDefuseKit:  pl.HasDefuseKit(),
				Equipment:     toEntityEquipment(pl.Weapons()),
				Team:          int(pl.Team),
			}

			// FIXME: Smoothify Positions

			snap.EntityUpdates = append(snap.EntityUpdates, e)
		}
	}

	return snap
}

func (m *minifier) updateKnownPlayers() {
	for _, pl := range m.parser.GameState().Participants().All() {
		if pl.EntityID != 0 {
			if _, alreadyKnown := m.knownPlayerEntityIDs[pl.EntityID]; !alreadyKnown {
				ent := rep.Entity{
					ID:    pl.EntityID,
					Team:  int(pl.Team),
					Name:  pl.Name,
					IsNpc: pl.IsBot,
				}

				m.replay.Entities = append(m.replay.Entities, ent)

				m.knownPlayerEntityIDs[pl.EntityID] = struct{}{}
			}
		}
	}
}

func (m *minifier) tickRate(rate float64) {
	m.replay.Header.TickRate = rate
	m.replay.Header.SnapshotRate = int(math.Round(rate / m.snapshotFrequency))
}

func r3VectorToPoint(v r3.Vector) rep.Point {
	return rep.Point{X: int(v.X), Y: int(v.Y), Z: int(v.Z)}
}

// roundTo wraps math.Round and allows specifying the rounding precision.
func roundTo(x, precision float64) float64 {
	return math.Round(x/precision) * precision
}

func toEntityEquipment(eq []*common.Equipment) []rep.EntityEquipment {
	var equipmentForPlayer = make([]rep.EntityEquipment, 0, len(eq))

	for _, equipment := range eq {
		equipmentForPlayer = append(equipmentForPlayer, rep.EntityEquipment{
			Type:           int(equipment.Type),
			AmmoInMagazine: equipment.AmmoInMagazine(),
			AmmoReserve:    equipment.AmmoReserve(),
		})
	}

	return equipmentForPlayer
}
