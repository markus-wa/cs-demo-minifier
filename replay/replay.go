// Package replay contains all types that make up a replay.
package replay

// Possible attribute kinds
const (
	AttrKindEntityID        = "entityId"
	AttrKindVictim          = "victim"
	AttrKindKiller          = "killer"
	AttrKindAssister        = "assister"
	AttrKindText            = "text"
	AttrKindSender          = "sender"
	AttrKindWeapon          = "weapon"
	AttrKindAttacker        = "attacker"
	AttrKindPlayer          = "player"
	AttrKindFlashDur        = "flashDuration"
	AttrKindTrajectoryStart = "trajectory_start"
	AttrKindTrajectoryEnd   = "trajectory_end"
	AttrKindTrajectory      = "trajectory"
)

// Possible event types
const (
	EventJump                     = "jump"
	EventFire                     = "fire"
	EventHurt                     = "hurt"
	EventKill                     = "kill"
	EventFlashed                  = "flashed"
	EventRoundStarted             = "round_started"
	EventRoundEnded               = "round_ended"
	EventSwapTeam                 = "swap_team"
	EventDisconnect               = "disconnect"
	EventChatMessage              = "chat_message"
	EventFootstep                 = "footstep"
	EventGrenadeProjectileDestroy = "grenade_destroy"
)

// Replay contains a minified demo
type Replay struct {
	Header    Header     `json:"header" msgpack:"header"`
	Entities  []Entity   `json:"entities" msgpack:"entities"`
	Snapshots []Snapshot `json:"snapshots" msgpack:"snapshots"`
	Ticks     []Tick     `json:"ticks" msgpack:"ticks"`
}

// Header holds the replay's general information
type Header struct {
	MapName      string  `json:"map" msgpack:"map"`
	TickRate     float64 `json:"tickRate" msgpack:"tickRate"`         // How many ticks per second
	SnapshotRate int     `json:"snapshotRate" msgpack:"snapshotRate"` // How many ticks per snapshot
}

// Entity holds players & NPCs
type Entity struct {
	ID    int    `json:"id" msgpack:"id"`
	Name  string `json:"name" msgpack:"name"`
	Team  int    `json:"team" msgpack:"team"`
	IsNpc bool   `json:"isNpc,omitempty" msgpack:"isNpc,omitempty"`
}

// Snapshot contains state changes since the last snapshot
type Snapshot struct {
	Tick          int            `json:"tick" msgpack:"tick"`
	EntityUpdates []EntityUpdate `json:"entityUpdates" msgpack:"entityUpdates"`
}

// EntityUpdate contains changes of player & NPCs attributes
type EntityUpdate struct {
	EntityID      int     `json:"entityId" msgpack:"entityId"`
	Team          int     `json:"team,omitempty" msgpack:"team,omitempty"`
	IsNpc         bool    `json:"isNpc,omitempty" msgpack:"isNpc,omitempty"`
	Positions     []Point `json:"positions,omitempty" msgpack:"positions,omitempty"` // This allows us smoother replay with less overhead compared to higher snapshot rate
	Angle         int     `json:"angle,omitempty" msgpack:"angle,omitempty"`
	Hp            int     `json:"hp,omitempty" msgpack:"hp,omitempty"`
	Armor         int     `json:"armor,omitempty" msgpack:"armor,omitempty"`
	FlashDuration float32 `json:"flashDuration,omitempty" msgpack:"flashDuration,omitempty"`
}

// Point is a position on the map
type Point struct {
	X int `json:"x" msgpack:"x"`
	Y int `json:"y" msgpack:"y"`
	Z int `json:"z" msgpack:"z"`
}

// Tick contains all events occurring at a specific tick
type Tick struct {
	Nr     int     `json:"nr" msgpack:"nr"`
	Events []Event `json:"events" msgpack:"events"`
}

// Event contains a game event
type Event struct {
	Name       string           `json:"name" msgpack:"name"`
	Attributes []EventAttribute `json:"attrs,omitempty" msgpack:"attrs,omitempty"`
}

// EventAttribute stores an additional attribute to an event
type EventAttribute struct {
	Key           string  `json:"key" msgpack:"key"`
	StrVal        string  `json:"strVal,omitempty" msgpack:"strVal,omitempty"`
	NumVal        float64 `json:"numVal,omitempty" msgpack:"numVal,omitempty"`
	TrajectoryVal []Point `json:"trajectory,omitempty" msgpack:"trajectory,omitempty"`
}
