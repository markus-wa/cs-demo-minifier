package csminify

import (
	"bufio"
	"bytes"
	"encoding/json"
	dem "github.com/markus-wa/demoinfocs-golang"
	"github.com/markus-wa/demoinfocs-golang/events"
	"io"
	"math"
)

func Minify(r io.Reader, snapsPerSec float32) []byte {
	var buf bytes.Buffer
	MinifyTo(r, bufio.NewWriter(&buf), snapsPerSec)
	return buf.Bytes()
}

func MinifyTo(r io.Reader, w io.Writer, snapsPerSec float32) {
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

	// FIXME: Don't ignore error & use Marshal instead of MarshalIndent
	//b, _ := json.MarshalIndent(m.replay, "", "\t")
	b, _ := json.Marshal(m.replay)

	w.Write(b)
}

type minifier struct {
	parser      *dem.Parser
	replay      replay
	currentTick tick
}

func (m *minifier) matchStarted(e events.MatchStartedEvent) {

	for _, pl := range m.parser.PlayingParticipants() {
		ent := entity{
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
		snap := snapshot{
			Tick: tick,
		}

		for _, pl := range m.parser.PlayingParticipants() {
			if pl.IsAlive() {
				e := entityUpdate{
					EntityID:      pl.EntityID,
					Hp:            pl.Hp,
					Armor:         pl.Armor,
					FlashDuration: pl.FlashDuration,
					Positions:     []point{r3VectorToPoint(pl.Position)}, // Maybe round this to save space
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
		m.currentTick = tick{}
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
