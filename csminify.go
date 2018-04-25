// Package csminify provides functions for parsing CS:GO demos and minifying them into various formats.
package csminify

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os"

	r3 "github.com/golang/geo/r3"
	dem "github.com/markus-wa/demoinfocs-golang"
	common "github.com/markus-wa/demoinfocs-golang/common"
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
	if err != nil {
		return err
	}
	return nil
}

// ToReplay reads a demo from r, takes snapshots (snapFreq/sec) and records events into a Replay.
func ToReplay(r io.Reader, snapFreq float32) (rep.Replay, error) {
	// FIXME: Smoothify flag
	// TODO: Maybe pass a WarnHandler along
	p := dem.NewParser(r, nil)
	_, err := p.ParseHeader()
	if err != nil {
		return rep.Replay{}, err
	}

	m := minifier{parser: p}

	m.replay.Header.MapName = p.Map()
	m.replay.Header.TickRate = p.FrameRate()
	m.replay.Header.SnapshotRate = int(round(float64(m.replay.Header.TickRate / snapFreq)))

	m.parser.RegisterEventHandler(m.tickDone)
	m.parser.RegisterEventHandler(m.roundStarted)
	m.parser.RegisterEventHandler(m.playerKilled)
	m.parser.RegisterEventHandler(m.playerHurt)
	m.parser.RegisterEventHandler(m.playerFlashed)
	m.parser.RegisterEventHandler(m.playerJump)
	m.parser.RegisterEventHandler(m.chatMessage)
	// TODO: Maybe hook up SayText (admin / console messages)
	//m.parser.RegisterEventHandler(m.sayText)

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
	parser      *dem.Parser
	replay      rep.Replay
	currentTick rep.Tick
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
					FlashDuration: float32(roundTo(float64(pl.FlashDuration), 0.1)), // Round to nearest 0.1
					Positions:     []rep.Point{r3VectorToPoint(pl.Position)},        // Maybe round the coordinates to save space
					Angle:         int(pl.ViewDirectionX),
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
	m.currentTick.Events = append(m.currentTick.Events, createEvent(rep.EventRoundStarted))
}

func (m *minifier) playerKilled(e events.PlayerKilledEvent) {
	if e.Victim == nil {
		return
	}

	eb := buildEvent(rep.EventKill).intAttr(rep.AttrKindVictim, e.Victim.EntityID)

	if e.Killer != nil && e.Killer != e.Victim {
		eb.intAttr(rep.AttrKindKiller, e.Killer.EntityID)
	}

	if e.Assister != nil {
		eb.intAttr(rep.AttrKindAssister, e.Assister.EntityID)
	}

	m.currentTick.Events = append(m.currentTick.Events, eb.build())
}

func (m *minifier) playerHurt(e events.PlayerHurtEvent) {
	m.addEntityEvent(rep.EventHurt, e.Player)
}

func (m *minifier) addEntityEvent(eventName string, pl *common.Player) {
	if pl != nil {
		m.currentTick.Events = append(m.currentTick.Events, createEntityEvent(eventName, pl.EntityID))
	} else {
		fmt.Fprintf(os.Stderr, "WARNING: Received %q event without player info\n", eventName)
	}
}

func (m *minifier) playerFlashed(e events.PlayerFlashedEvent) {
	m.addEntityEvent(rep.EventHurt, e.Player)
}

func (m *minifier) playerJump(e events.PlayerJumpEvent) {
	m.addEntityEvent(rep.EventJump, e.Player)
}

func (m *minifier) playerTeamChange(e events.PlayerTeamChangeEvent) {
	m.addEntityEvent(rep.EventSwapTeam, e.Player)
}

func (m *minifier) playerDisconnect(e events.PlayerDisconnectEvent) {
	m.addEntityEvent(rep.EventDisconnect, e.Player)
}

func (m *minifier) weaponFired(e events.WeaponFiredEvent) {
	m.addEntityEvent(rep.EventFire, e.Shooter)
}

func r3VectorToPoint(v r3.Vector) rep.Point {
	return rep.Point{X: int(v.X), Y: int(v.Y)}
}

func (m *minifier) chatMessage(e events.ChatMessageEvent) {
	eb := buildEvent(rep.EventChatMessage)
	eb = eb.stringAttr(rep.AttrKindText, e.Text)

	// Skip for now, probably always true anyway
	//eb = eb.boolAttr("isChatAll", e.IsChatAll)

	if e.Sender != nil {
		eb = eb.intAttr(rep.AttrKindSender, e.Sender.EntityID)
	} else {
	}

	m.currentTick.Events = append(m.currentTick.Events, eb.build())
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

func (b eventBuilder) intAttr(key string, value int) eventBuilder {
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

func createEntityEvent(eventName string, entityID int) rep.Event {
	return buildEvent(eventName).intAttr(rep.AttrKindEntityID, entityID).build()
}
