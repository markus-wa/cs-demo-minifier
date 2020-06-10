package protobuf

import (
	"fmt"
	"io"
	"io/ioutil"

	common "github.com/markus-wa/demoinfocs-golang/v2/pkg/demoinfocs/common"

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
	if entities == nil {
		return nil
	}

	result := make([]rep.Entity, len(entities))
	for i, e := range entities {
		result[i] = rep.Entity{
			ID:    int(e.Id),
			Team:  mapFromTeam(e.Team),
			Name:  e.Name,
			IsNpc: e.IsNpc,
		}
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
	if snaps == nil {
		return nil
	}

	result := make([]rep.Snapshot, len(snaps))
	for i, s := range snaps {
		result[i] = rep.Snapshot{
			Tick:          int(s.Tick),
			EntityUpdates: mapFromEntityUpdates(s.EntityUpdates),
		}
	}

	return result
}

func mapFromEntityUpdates(entityUpdates []*gen.Replay_Snapshot_EntityUpdate) []rep.EntityUpdate {
	if entityUpdates == nil {
		return nil
	}

	result := make([]rep.EntityUpdate, len(entityUpdates))
	for i, u := range entityUpdates {
		result[i] = rep.EntityUpdate{
			AngleX:        int(u.AngleX),
			AngleY:        int(u.AngleY),
			Armor:         int(u.Armor),
			EntityID:      int(u.EntityId),
			FlashDuration: u.FlashDuration,
			Hp:            int(u.Hp),
			Positions:     mapFromPositions(u.Positions),
			IsNpc:         u.IsNpc,
			Team:          mapFromTeam(u.Team),
			HasHelmet:     u.HasHelmet,
			HasDefuseKit:  u.HasDefuseKit,
			Equipment:     mapFromEquipment(u.Equipment),
		}
	}

	return result
}

func mapFromEquipment(equipment []*gen.Replay_Snapshot_EntityEquipment) []rep.EntityEquipment {
	if equipment == nil {
		return nil
	}

	result := make([]rep.EntityEquipment, len(equipment))
	for i, eq := range equipment {
		result[i] = rep.EntityEquipment{
			Type:           int(eq.Type),
			AmmoInMagazine: int(eq.AmmoInMagazine),
			AmmoReserve:    int(eq.AmmoReserve),
		}
	}

	return result
}

func mapFromPositions(positions []*gen.Point) []rep.Point {
	if positions == nil {
		return nil
	}

	result := make([]rep.Point, len(positions))
	for i, p := range positions {
		result[i] = mapFromPosition(p)
	}

	return result
}

func mapFromPosition(p *gen.Point) rep.Point {
	return rep.Point{
		X: int(p.X),
		Y: int(p.Y),
		Z: int(p.Z),
	}
}

func mapFromTicks(ticks []*gen.Replay_Tick) []rep.Tick {
	if ticks == nil {
		return nil
	}

	result := make([]rep.Tick, len(ticks))
	for i, t := range ticks {
		result[i] = rep.Tick{
			Nr:     int(t.Nr),
			Events: mapFromEvents(t.Events),
		}
	}

	return result
}

func mapFromEvents(events []*gen.Replay_Tick_Event) []rep.Event {
	if events == nil {
		return nil
	}

	result := make([]rep.Event, len(events))
	for i, e := range events {
		// Custom events
		var name string
		if e.Kind == gen.Replay_Tick_Event_CUSTOM {
			for _, attr := range e.Attributes {
				if attr.Kind == gen.Replay_Tick_Event_Attribute_EVENT_NAME {
					name = attr.StringValue
				}
			}
		} else {
			name = mapFromEventKind(e.Kind)
		}
		result[i] = rep.Event{
			Name:       name,
			Attributes: mapFromAttributes(e.Attributes),
		}
	}

	return result
}

func mapFromAttributes(attrs []*gen.Replay_Tick_Event_Attribute) []rep.EventAttribute {
	if attrs == nil {
		return nil
	}

	result := make([]rep.EventAttribute, 0, len(attrs))
	for _, a := range attrs {
		// Skip the internal 'event name' attributes
		if a.Kind != gen.Replay_Tick_Event_Attribute_EVENT_NAME {
			// Custom attributes
			var key string
			if a.Kind == gen.Replay_Tick_Event_Attribute_CUSTOM {
				key = a.CustomName
			} else {
				key = mapFromAttributeKind(a.Kind)
			}
			result = append(result, rep.EventAttribute{
				Key:    key,
				NumVal: a.NumberValue,
				StrVal: a.StringValue,
			})
		}
	}

	return result
}

func mapFromAttributeKind(key gen.Replay_Tick_Event_Attribute_Kind) string {
	kind, ok := attributeKindMap.GetInverse(key)
	if !ok {
		panic(fmt.Errorf("unknown attribute kind %q", key))
	}
	return kind.(string)
}

func mapFromEventKind(name gen.Replay_Tick_Event_Kind) string {
	kind, ok := eventKindMap.GetInverse(name)
	if !ok {
		panic(fmt.Errorf("unknown event kind %q", name))
	}
	return kind.(string)
}
