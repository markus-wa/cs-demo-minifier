package csminify

import (
	"bufio"
	"bytes"
	"github.com/golang/geo/r3"
	rep "gitlab.com/markus-wa/cs-demo-minifier/csminify/replay"
	dem "github.com/markus-wa/demoinfocs-golang"
	"github.com/markus-wa/demoinfocs-golang/events"
	"io"
	"math"
)

// ReplayMarshaller is the signature for functions that serialize replay.Replay structs to an io.Writer
type ReplayMarshaller func(rep rep.Replay, w io.Writer) error

// Minify wraps MinifyTo with a bytes.Buffer and returns the written bytes
func Minify(r io.Reader, marshal ReplayMarshaller, snapsPerSec float32) (error, []byte) {
	var buf bytes.Buffer
	err := MinifyTo(r, snapsPerSec, marshal, bufio.NewWriter(&buf))
	if err != nil {
		return err, nil
	}
	return nil, buf.Bytes()
}

// MinifyTo reads a demo from r, creates snapshots with a frequency of snapFreq/sec and writers the result of marshal to w
func MinifyTo(r io.Reader, snapFreq float32, marshal ReplayMarshaller, w io.Writer) error {
	p := dem.NewParser(r)
	err := p.ParseHeader()
	if err != nil {
		return err;
	}

	m := minifier{parser: p}

	m.replay.Header.TickRate = p.FrameRate()

	snapRate := float64(m.replay.Header.TickRate / snapFreq)

	// How on earth is there still no math.Round()?! https://github.com/golang/go/issues/4594
	if math.Abs(snapRate) >= 0.5 {
		m.replay.Header.SnapshotRate = int(snapRate + math.Copysign(0.5, snapRate))
	}

	m.replay.Header.MapName = p.Map()

	p.RegisterEventHandler(m.matchStarted)

	err = p.ParseToEnd()
	if err != nil {
		return err
	}

	// FIXME: Don't ignore error
	marshal(m.replay, w)
	return nil
}

type minifier struct {
	parser      *dem.Parser
	replay      rep.Replay
	currentTick rep.Tick
}

func (m *minifier) matchStarted(e events.MatchStartedEvent) {

	for _, pl := range m.parser.PlayingParticipants() {
		ent := rep.Entity{
			ID:    pl.EntityID,
			Team:  int(pl.Team),
			Name:  pl.Name,
			IsNpc: pl.IsBot,
		}
		// FIXME: Smoothify flag

		m.replay.Entities = append(m.replay.Entities, ent)
	}

	m.parser.RegisterEventHandler(m.tickDone)
	m.parser.RegisterEventHandler(m.roundStarted)
	m.parser.RegisterEventHandler(m.playerKilled)
	m.parser.RegisterEventHandler(m.playerHurt)
	m.parser.RegisterEventHandler(m.playerFlashed)
	m.parser.RegisterEventHandler(m.playerJump)
}

func (m *minifier) tickDone(e events.TickDoneEvent) {
	if tick := m.parser.CurrentFrame(); tick%m.replay.Header.SnapshotRate == 0 {
		snap := rep.Snapshot{
			Tick: tick,
		}

		for _, pl := range m.parser.PlayingParticipants() {
			if pl.IsAlive() {
				e := rep.EntityUpdate{
					EntityID:      pl.EntityID,
					Hp:            pl.Hp,
					Armor:         pl.Armor,
					FlashDuration: pl.FlashDuration,
					Positions:     []rep.Point{r3VectorToPoint(pl.Position)}, // Maybe round the coordinates to save space
					Angle:         pl.ViewDirectionX,
				}
				// FIXME: Smoothify
				snap.EntityUpdates = append(snap.EntityUpdates, e)
			}
		}

		m.replay.Snapshots = append(m.replay.Snapshots, snap)
	}

	if len(m.currentTick.Events) > 0 {
		m.currentTick.Nr = m.parser.CurrentFrame()
		m.replay.Ticks = append(m.replay.Ticks, m.currentTick)
		m.currentTick = rep.Tick{}
	}
}

func (m *minifier) roundStarted(e events.RoundStartedEvent) {
	m.currentTick.Events = append(m.currentTick.Events, createEvent("round_started"))
}

func (m *minifier) playerKilled(e events.PlayerKilledEvent) {
	if e.Victim == nil {
		return
	}

	eb := buildEvent("kill").numAttr("victim", e.Victim.EntityID)

	if e.Killer != nil && e.Killer != e.Victim {
		eb.numAttr("killer", e.Killer.EntityID)
	}

	if e.Assister != nil {
		eb.numAttr("assister", e.Assister.EntityID)
	}

	m.currentTick.Events = append(m.currentTick.Events, eb.build())
}

func (m *minifier) playerHurt(e events.PlayerHurtEvent) {
	m.currentTick.Events = append(m.currentTick.Events, createEntityEvent("hurt", e.Player.EntityID))
}

func (m *minifier) playerFlashed(e events.PlayerFlashedEvent) {
	m.currentTick.Events = append(m.currentTick.Events, createEntityEvent("flashed", e.Player.EntityID))
}

func (m *minifier) playerJump(e events.PlayerJumpEvent) {
	m.currentTick.Events = append(m.currentTick.Events, createEntityEvent("jump", e.Player.EntityID))
}

func (m *minifier) playerTeamChange(e events.PlayerTeamChangeEvent) {
	m.currentTick.Events = append(m.currentTick.Events, createEntityEvent("swap_team", e.Player.EntityID))
}

func (m *minifier) playerDisconnect(e events.PlayerDisconnectEvent) {
	m.currentTick.Events = append(m.currentTick.Events, createEntityEvent("disconnect", e.Player.EntityID))
}

func (m *minifier) weaponFired(e events.WeaponFiredEvent) {
	m.currentTick.Events = append(m.currentTick.Events, createEntityEvent("fire", e.Shooter.EntityID))
}

func r3VectorToPoint(v r3.Vector) rep.Point {
	return rep.Point{X: v.X, Y: v.Y}
}

type eventBuilder struct {
	event rep.Event
}

func (b eventBuilder) stringAttr(key string, value string) eventBuilder {
	b.event.Attributes = append(b.event.Attributes, rep.EventAttribute{
		Key:    key,
		StrVal: value,
	})
	return b
}

func (b eventBuilder) numAttr(key string, value int) eventBuilder {
	b.event.Attributes = append(b.event.Attributes, rep.EventAttribute{
		Key:    key,
		NumVal: float64(value),
	})
	return b
}

func (b eventBuilder) build() rep.Event {
	return b.event
}

func buildEvent(eventName string) eventBuilder {
	return eventBuilder{
		event: createEvent(eventName),
	}
}

func createEvent(eventName string) rep.Event {
	return rep.Event{
		Name: eventName,
	}
}

func createEntityEvent(eventName string, entityId int) rep.Event {
	return buildEvent(eventName).numAttr("entityId", entityId).build()
}
