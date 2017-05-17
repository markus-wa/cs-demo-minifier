package csminify

import (
	"bufio"
	"bytes"
	rep "github.com/markus-wa/cs-demo-minifier/csminify/replay"
	dem "github.com/markus-wa/demoinfocs-golang"
	"github.com/markus-wa/demoinfocs-golang/events"
	"io"
	"math"
)

type ReplayMarshaller func(rep rep.Replay, w io.Writer) error

func Minify(r io.Reader, marshal ReplayMarshaller, snapsPerSec float32) []byte {
	var buf bytes.Buffer
	MinifyTo(r, snapsPerSec, marshal, bufio.NewWriter(&buf))
	return buf.Bytes()
}

func MinifyTo(r io.Reader, snapsPerSec float32, marshal ReplayMarshaller, w io.Writer) {
	p := dem.NewParser(r)
	p.ParseHeader()

	m := minifier{parser: p}

	m.replay.Header.TickRate = p.TickRate()

	f := float64(m.replay.Header.TickRate / snapsPerSec)

	// How on earth is there still no math.Round()?! https://github.com/golang/go/issues/4594
	if math.Abs(f) >= 0.5 {
		m.replay.Header.SnapshotRate = int(f + math.Copysign(0.5, f))
	}

	m.replay.Header.MapName = p.Map()

	p.RegisterEventHandler(m.matchStarted)

	p.ParseToEnd(nil)

	// FIXME: Don't ignore error
	marshal(m.replay, w)
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
	if tick := m.parser.CurrentTick(); tick%m.replay.Header.SnapshotRate == 0 {
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

	if len(m.currentTick.GameEvents) > 0 || len(m.currentTick.MapEvents) > 0 || len(m.currentTick.EntityEvents) > 0 {
		m.currentTick.Nr = m.parser.CurrentTick()
		m.replay.Ticks = append(m.replay.Ticks, m.currentTick)
		m.currentTick = rep.Tick{}
	}
}

func (m *minifier) roundStarted(e events.RoundStartedEvent) {
	m.currentTick.GameEvents = append(m.currentTick.GameEvents, createGameEvent("round_started", "New round started"))
}

func (m *minifier) playerKilled(e events.PlayerKilledEvent) {
	if e.Victim == nil {
		return
	}

	var msg string

	if e.Killer == nil || e.Killer == e.Victim {
		msg = e.Victim.Name + " killed himself"
	} else {
		msg = e.Killer.Name + " killed " + e.Victim.Name
	}

	if e.Assister != nil {
		msg += " with the help of " + e.Assister.Name
	}

	m.currentTick.GameEvents = append(m.currentTick.GameEvents, createGameEvent("kill", msg))
	m.currentTick.EntityEvents = append(m.currentTick.EntityEvents, createEntityEvent("die", e.Victim.EntityID))
}

func (m *minifier) playerHurt(e events.PlayerHurtEvent) {
	m.currentTick.EntityEvents = append(m.currentTick.EntityEvents, createEntityEvent("hurt", e.Player.EntityID))
}

func (m *minifier) playerFlashed(e events.PlayerFlashedEvent) {
	m.currentTick.EntityEvents = append(m.currentTick.EntityEvents, createEntityEvent("flashed", e.Player.EntityID))
}

func (m *minifier) playerJump(e events.PlayerJumpEvent) {
	m.currentTick.EntityEvents = append(m.currentTick.EntityEvents, createEntityEvent("jump", e.Player.EntityID))
}

func (m *minifier) playerTeamChange(e events.PlayerTeamChangeEvent) {
	m.currentTick.GameEvents = append(m.currentTick.GameEvents, createGameEvent("player_swap_team", e.Player.Name+" switched teams"))
	m.currentTick.EntityEvents = append(m.currentTick.EntityEvents, createEntityEvent("swap_team", e.Player.EntityID))
}

func (m *minifier) playerDisconnect(e events.PlayerDisconnectEvent) {
	m.currentTick.GameEvents = append(m.currentTick.GameEvents, createGameEvent("player_disconnect", e.Player.Name+" disconnected"))
	m.currentTick.EntityEvents = append(m.currentTick.EntityEvents, createEntityEvent("disconnect", e.Player.EntityID))
}

func (m *minifier) weaponFired(e events.WeaponFiredEvent) {
	m.currentTick.EntityEvents = append(m.currentTick.EntityEvents, createEntityEvent("fire", e.Shooter.EntityID))
}
