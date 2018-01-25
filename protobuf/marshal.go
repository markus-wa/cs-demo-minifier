package protobuf

import (
	fmt "fmt"
	io "io"

	rep "github.com/markus-wa/cs-demo-minifier/replay"
	common "github.com/markus-wa/demoinfocs-golang/common"
)

// MarshalReplay serializes a Replay as protobuf to an io.Writer
func MarshalReplay(r rep.Replay, w io.Writer) error {
	pbReplay := Replay{
		Entities: mapToEntities(r.Entities),
		Header: &Replay_Header{
			Map:          r.Header.MapName,
			SnapshotRate: int32(r.Header.SnapshotRate),
			TickRate:     r.Header.TickRate,
		},
		Snapshots: mapToSnapshots(r.Snapshots),
		Ticks:     mapToTicks(r.Ticks),
	}

	data, err := pbReplay.Marshal()
	if err != nil {
		return err
	}
	_, err = w.Write(data)
	return err
}

func mapToEntities(entities []rep.Entity) []*Replay_Entity {
	result := make([]*Replay_Entity, 0)
	for _, e := range entities {
		result = append(result, &Replay_Entity{
			Id:    int32(e.ID),
			Team:  mapToTeam(e.Team),
			Name:  e.Name,
			IsNpc: e.IsNpc,
		})
	}
	return result
}

func mapToTeam(team int) Team {
	var result Team
	switch common.Team(team) {
	case common.TeamTerrorists:
		result = Team_TERRORIST
	case common.TeamCounterTerrorists:
		result = Team_COUNTER_TERRORIST
	case common.TeamSpectators:
		result = Team_SPECTATOR
	default:
		result = Team_UNASSIGNED
	}
	return result
}

func mapToSnapshots(snaps []rep.Snapshot) []*Replay_Snapshot {
	result := make([]*Replay_Snapshot, 0)
	for _, s := range snaps {
		result = append(result, &Replay_Snapshot{
			Tick:          int32(s.Tick),
			EntityUpdates: mapToEntityUpdates(s.EntityUpdates),
		})
	}
	return result
}

func mapToEntityUpdates(entityUpdates []rep.EntityUpdate) []*Replay_Snapshot_EntityUpdate {
	result := make([]*Replay_Snapshot_EntityUpdate, 0)
	for _, u := range entityUpdates {
		result = append(result, &Replay_Snapshot_EntityUpdate{
			Angle:         int32(u.Angle),
			Armor:         int32(u.Armor),
			EntityId:      int32(u.EntityID),
			FlashDuration: u.FlashDuration,
			Hp:            int32(u.Hp),
			Positions:     mapToPositions(u.Positions),
			IsNpc:         u.IsNpc,
			Team:          mapToTeam(u.Team),
		})
	}
	return result
}

func mapToPositions(positions []rep.Point) []*Point {
	result := make([]*Point, 0)
	for _, p := range positions {
		result = append(result, mapToPosition(p))
	}
	return result
}

func mapToPosition(p rep.Point) *Point {
	return &Point{
		X: int32(p.X),
		Y: int32(p.Y),
	}
}

func mapToTicks(ticks []rep.Tick) []*Replay_Tick {
	result := make([]*Replay_Tick, 0)
	for _, t := range ticks {
		result = append(result, &Replay_Tick{
			Nr:     int32(t.Nr),
			Events: mapToEvents(t.Events),
		})
	}
	return result
}

func mapToEvents(events []rep.Event) []*Replay_Tick_Event {
	result := make([]*Replay_Tick_Event, 0)
	for _, e := range events {
		result = append(result, &Replay_Tick_Event{
			Kind:       mapToEventKind(e.Name),
			Attributes: mapToAttributes(e.Attributes),
		})
	}
	return result
}

func mapToAttributes(attrs []rep.EventAttribute) []*Replay_Tick_Event_Attribute {
	if attrs == nil {
		return nil
	}
	result := make([]*Replay_Tick_Event_Attribute, 0)
	for _, a := range attrs {
		result = append(result, &Replay_Tick_Event_Attribute{
			Kind:        mapToAttributeKind(a.Key),
			NumberValue: a.NumVal,
			StringValue: a.StrVal,
		})
	}
	return result
}

func mapToAttributeKind(key string) Replay_Tick_Event_Attribute_Kind {
	kind, ok := attributeKindMap.Get(key)
	if !ok {
		panic(fmt.Errorf("Unknown attribute kind %q", key))
	}
	return kind.(Replay_Tick_Event_Attribute_Kind)
}

func mapToEventKind(name string) Replay_Tick_Event_Kind {
	kind, ok := eventKindMap.Get(name)
	if !ok {
		panic(fmt.Errorf("Unknown event kind %q", name))
	}
	return kind.(Replay_Tick_Event_Kind)
}
