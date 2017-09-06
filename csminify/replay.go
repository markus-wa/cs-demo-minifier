package csminify

import (
	"github.com/golang/geo/r3"
)

type replay struct {
	Header    header     `json:"header"`
	Entities  []entity   `json:"entities"`
	Snapshots []snapshot `json:"snapshots"`
	Ticks     []tick     `json:"ticks"`
}

type header struct {
	MapName      string  `json:"map"`
	TickRate     float32 `json:"tickRate"`     // How many ticks per second
	SnapshotRate int     `json:"snapshotRate"` // How many ticks per snapshot
}

// Players & NPCs
type entity struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Team  int    `json:"team"`
	IsNpc bool   `json:"isNpc,omitempty"`
}

type snapshot struct {
	Tick          int            `json:"tick"`
	EntityUpdates []entityUpdate `json:"entityUpdates"`
}

// Players & NPCs
type entityUpdate struct {
	EntityID      int     `json:"entityId"`
	Team          int     `json:"team,omitempty"`
	IsNpc         bool    `json:"isNpc,omitempty"`
	Positions     []point `json:"positions,omitempty"` // This allows us smoother replay with less overhead compared to higher snapshot rate
	Angle         float32 `json:"angle,omitempty"`
	Hp            int     `json:"hp,omitempty"`
	Armor         int     `json:"armor,omitempty"`
	FlashDuration float32 `json:"flashDuration,omitempty"`
}

type point struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
}

func r3VectorToPoint(v r3.Vector) point {
	return point{X: v.X, Y: v.Y}
}

type tick struct {
	Nr           int           `json:"nr"`
	GameEvents   []gameEvent   `json:"gameEvents"`
	MapEvents    []mapEvent    `json:"mapEvents"`
	EntityEvents []entityEvent `json:"entityEvents"`
}

type event struct {
	Tick int    `json:"tick"`
	Name string `json:"name"`
}

type gameEvent struct {
	event
	EventString string `json:"event"`
}

func createGameEvent(eventName string, eventString string) gameEvent {
	return gameEvent{
		event: event{
			Name: eventName,
		},
		EventString: eventString,
	}
}

type mapEvent struct {
	event
	Location    point `json:"location"`
	TriggeredBy int   `json:"triggeredByEntityId, omitempty"` // FIXME: Can entityId be 0? if so: dont omitempty
}

type entityEvent struct {
	event
	EntityID int `json:"entityId"`
}

func createEntityEvent(eventName string, entityID int) entityEvent {
	return entityEvent{
		event: event{
			Name: eventName,
		},
		EntityID: entityID,
	}
}
