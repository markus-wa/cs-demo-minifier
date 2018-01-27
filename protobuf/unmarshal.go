package protobuf

import (
	fmt "fmt"
	"io"
	"io/ioutil"

	common "github.com/markus-wa/demoinfocs-golang/common"

	gen "github.com/markus-wa/cs-demo-minifier/protobuf/gen"
	rep "github.com/markus-wa/cs-demo-minifier/replay"
)

// UnmarshalReplay deserializes protobuf data from a io.Reader into a Replay.
func UnmarshalReplay(r io.Reader, replay *rep.Replay) error {
	var pbReplay gen.Replay
	b, err := ioutil.ReadAll(r)
	if err != nil {
		return err
	}

	err = pbReplay.Unmarshal(b)
	if err != nil {
		return err
	}

	replay.Header = mapFromHeader(pbReplay.Header)
	replay.Entities = mapFromEntities(pbReplay.Entities)
	replay.Snapshots = mapFromSnapshots(pbReplay.Snapshots)
	replay.Ticks = mapFromTicks(pbReplay.Ticks)

	return nil
}

func mapFromHeader(header *gen.Replay_Header) rep.Header {
	return rep.Header{
		MapName:      header.Map,
		SnapshotRate: int(header.SnapshotRate),
		TickRate:     header.TickRate,
	}
}

func mapFromEntities(entities []*gen.Replay_Entity) []rep.Entity {
	result := make([]rep.Entity, 0)
	for _, e := range entities {
		result = append(result, rep.Entity{
			ID:    int(e.Id),
			Team:  mapFromTeam(e.Team),
			Name:  e.Name,
			IsNpc: e.IsNpc,
		})
	}
	return result
}

func mapFromTeam(team gen.Team) int {
	var result common.Team
	switch team {
	case gen.Team_TERRORIST:
		result = common.TeamTerrorists
	case gen.Team_COUNTER_TERRORIST:
		result = common.TeamCounterTerrorists
	case gen.Team_SPECTATOR:
		result = common.TeamSpectators
	default:
		result = common.TeamUnassigned
	}
	return int(result)
}

func mapFromSnapshots(snaps []*gen.Replay_Snapshot) []rep.Snapshot {
	result := make([]rep.Snapshot, 0)
	for _, s := range snaps {
		result = append(result, rep.Snapshot{
			Tick:          int(s.Tick),
			EntityUpdates: mapFromEntityUpdates(s.EntityUpdates),
		})
	}
	return result
}

func mapFromEntityUpdates(entityUpdates []*gen.Replay_Snapshot_EntityUpdate) []rep.EntityUpdate {
	result := make([]rep.EntityUpdate, 0)
	for _, u := range entityUpdates {
		result = append(result, rep.EntityUpdate{
			Angle:         int(u.Angle),
			Armor:         int(u.Armor),
			EntityID:      int(u.EntityId),
			FlashDuration: u.FlashDuration,
			Hp:            int(u.Hp),
			Positions:     mapFromPositions(u.Positions),
			IsNpc:         u.IsNpc,
			Team:          mapFromTeam(u.Team),
		})
	}
	return result
}

func mapFromPositions(positions []*gen.Point) []rep.Point {
	result := make([]rep.Point, 0)
	for _, p := range positions {
		result = append(result, mapFromPosition(p))
	}
	return result
}

func mapFromPosition(p *gen.Point) rep.Point {
	return rep.Point{
		X: int(p.X),
		Y: int(p.Y),
	}
}

func mapFromTicks(ticks []*gen.Replay_Tick) []rep.Tick {
	result := make([]rep.Tick, 0)
	for _, t := range ticks {
		result = append(result, rep.Tick{
			Nr:     int(t.Nr),
			Events: mapFromEvents(t.Events),
		})
	}
	return result
}

func mapFromEvents(events []*gen.Replay_Tick_Event) []rep.Event {
	result := make([]rep.Event, 0)
	for _, e := range events {
		result = append(result, rep.Event{
			Name:       mapFromEventKind(e.Kind),
			Attributes: mapFromAttributes(e.Attributes),
		})
	}
	return result
}

func mapFromAttributes(attrs []*gen.Replay_Tick_Event_Attribute) []rep.EventAttribute {
	if attrs == nil {
		return nil
	}
	result := make([]rep.EventAttribute, 0)
	for _, a := range attrs {
		result = append(result, rep.EventAttribute{
			Key:    mapFromAttributeKind(a.Kind),
			NumVal: a.NumberValue,
			StrVal: a.StringValue,
		})
	}
	return result
}

func mapFromAttributeKind(key gen.Replay_Tick_Event_Attribute_Kind) string {
	kind, ok := attributeKindMap.GetInverse(key)
	if !ok {
		panic(fmt.Errorf("Unknown attribute kind %q", key))
	}
	return kind.(string)
}

func mapFromEventKind(name gen.Replay_Tick_Event_Kind) string {
	kind, ok := eventKindMap.GetInverse(name)
	if !ok {
		panic(fmt.Errorf("Unknown event kind %q", name))
	}
	return kind.(string)
}
