package replay

type Replay struct {
	Header    Header     `json:"header"`
	Entities  []Entity   `json:"entities"`
	Snapshots []Snapshot `json:"snapshots"`
	Ticks     []Tick     `json:"ticks"`
}

type Header struct {
	MapName      string  `json:"map"`
	TickRate     float32 `json:"tickRate"`     // How many ticks per second
	SnapshotRate int     `json:"snapshotRate"` // How many ticks per snapshot
}

// Players & NPCs
type Entity struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Team  int    `json:"team"`
	IsNpc bool   `json:"isNpc,omitempty"`
}

type Snapshot struct {
	Tick          int            `json:"tick"`
	EntityUpdates []EntityUpdate `json:"entityUpdates"`
}

// Players & NPCs
type EntityUpdate struct {
	EntityID      int     `json:"entityId"`
	Team          int     `json:"team,omitempty"`
	IsNpc         bool    `json:"isNpc,omitempty"`
	Positions     []Point `json:"positions,omitempty"` // This allows us smoother replay with less overhead compared to higher snapshot rate
	Angle         float32 `json:"angle,omitempty"`
	Hp            int     `json:"hp,omitempty"`
	Armor         int     `json:"armor,omitempty"`
	FlashDuration float32 `json:"flashDuration,omitempty"`
}

type Point struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
}

type Tick struct {
	Nr           int           `json:"nr"`
	GameEvents   []GameEvent   `json:"gameEvents"`
	MapEvents    []MapEvent    `json:"mapEvents"`
	EntityEvents []EntityEvent `json:"entityEvents"`
}

type Event struct {
	Tick int    `json:"tick"`
	Name string `json:"name"`
}

type GameEvent struct {
	Event
	EventString string `json:"event"`
}

type MapEvent struct {
	Event
	Location    Point `json:"location"`
	TriggeredBy int   `json:"triggeredByEntityId, omitempty"` // FIXME: Can entityId be 0? if so: dont omitempty
}

type EntityEvent struct {
	Event
	EntityID int `json:"entityId"`
}
