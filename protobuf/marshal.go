package protobuf

import (
	io "io"

	common "github.com/markus-wa/demoinfocs-golang/common"

	gen "github.com/markus-wa/cs-demo-minifier/protobuf/gen"
	rep "github.com/markus-wa/cs-demo-minifier/replay"
)

// MarshalReplay serializes a Replay as protobuf to an io.Writer
func MarshalReplay(r rep.Replay, w io.Writer) error {
	pbReplay := gen.Replay{
		Entities: mapToEntities(r.Entities),
		Header: &gen.Replay_Header{
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

func mapToEntities(entities []rep.Entity) []*gen.Replay_Entity {
	result := make([]*gen.Replay_Entity, 0)
	for _, e := range entities {
		result = append(result, &gen.Replay_Entity{
			Id:    int32(e.ID),
			Team:  mapToTeam(e.Team),
			Name:  e.Name,
			IsNpc: e.IsNpc,
		})
	}
	return result
}

func mapToTeam(team int) gen.Team {
	var result gen.Team
	switch common.Team(team) {
	case common.TeamTerrorists:
		result = gen.Team_TERRORIST
	case common.TeamCounterTerrorists:
		result = gen.Team_COUNTER_TERRORIST
	case common.TeamSpectators:
		result = gen.Team_SPECTATOR
	default:
		result = gen.Team_UNASSIGNED
	}
	return result
}

func mapToSnapshots(snaps []rep.Snapshot) []*gen.Replay_Snapshot {
	result := make([]*gen.Replay_Snapshot, 0)
	for _, s := range snaps {
		result = append(result, &gen.Replay_Snapshot{
			Tick:          int32(s.Tick),
			EntityUpdates: mapToEntityUpdates(s.EntityUpdates),
		})
	}
	return result
}

func mapToEntityUpdates(entityUpdates []rep.EntityUpdate) []*gen.Replay_Snapshot_EntityUpdate {
	result := make([]*gen.Replay_Snapshot_EntityUpdate, 0)
	for _, u := range entityUpdates {
		result = append(result, &gen.Replay_Snapshot_EntityUpdate{
			AngleX:        int32(u.AngleX),
			AngleY:        int32(u.AngleY),
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

func mapToPositions(positions []rep.Point) []*gen.Point {
	result := make([]*gen.Point, 0)
	for _, p := range positions {
		result = append(result, mapToPosition(p))
	}
	return result
}

func mapToPosition(p rep.Point) *gen.Point {
	return &gen.Point{
		X: int32(p.X),
		Y: int32(p.Y),
		Z: int32(p.Z),
	}
}

func mapToTicks(ticks []rep.Tick) []*gen.Replay_Tick {
	result := make([]*gen.Replay_Tick, 0)
	for _, t := range ticks {
		result = append(result, &gen.Replay_Tick{
			Nr:     int32(t.Nr),
			Events: mapToEvents(t.Events),
		})
	}
	return result
}

func mapToEvents(events []rep.Event) []*gen.Replay_Tick_Event {
	result := make([]*gen.Replay_Tick_Event, 0)
	for _, e := range events {
		kind := mapToEventKind(e.Name)
		// Custom events
		if kind == gen.Replay_Tick_Event_CUSTOM {
			e.Attributes = append(e.Attributes, rep.EventAttribute{
				Key:    attrKindEventName,
				StrVal: e.Name,
			})
		}
		result = append(result, &gen.Replay_Tick_Event{
			Kind:       kind,
			Attributes: mapToAttributes(e.Attributes),
		})
	}
	return result
}

func mapToAttributes(attrs []rep.EventAttribute) []*gen.Replay_Tick_Event_Attribute {
	if attrs == nil {
		return nil
	}
	result := make([]*gen.Replay_Tick_Event_Attribute, 0)
	for _, a := range attrs {
		kind := mapToAttributeKind(a.Key)
		// Custom attributes
		var customName string
		if kind == gen.Replay_Tick_Event_Attribute_CUSTOM {
			customName = a.Key
		}
		result = append(result, &gen.Replay_Tick_Event_Attribute{
			Kind:        kind,
			NumberValue: a.NumVal,
			StringValue: a.StrVal,
			CustomName:  customName,
		})
	}
	return result
}

func mapToAttributeKind(key string) gen.Replay_Tick_Event_Attribute_Kind {
	kind, ok := attributeKindMap.Get(key)
	if !ok {
		return gen.Replay_Tick_Event_Attribute_CUSTOM
	}
	return kind.(gen.Replay_Tick_Event_Attribute_Kind)
}

func mapToEventKind(name string) gen.Replay_Tick_Event_Kind {
	kind, ok := eventKindMap.Get(name)
	if !ok {
		return gen.Replay_Tick_Event_CUSTOM
	}
	return kind.(gen.Replay_Tick_Event_Kind)
}
