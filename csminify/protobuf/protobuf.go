package protobuf

//go:generate protoc -I=proto --gogofaster_out=Mgoogle/protobuf/any.proto=github.com/gogo/protobuf/types:. proto/*.proto

import (
	rep "github.com/markus-wa/cs-demo-minifier/csminify/replay"
	"github.com/markus-wa/demoinfocs-golang/common"
	"io"
)

// MarshalReplay serializes a Replay as protobuf to an io.Writer
func MarshalReplay(replay rep.Replay, w io.Writer) error {
	ticks, err := mapTicks(replay.Ticks)
	if err != nil {
		return err
	}
	rep := Replay{
		Entities: mapEntities(replay.Entities),
		Header: &Replay_Header{
			Map:          replay.Header.MapName,
			SnapshotRate: int32(replay.Header.SnapshotRate),
			TickRate:     replay.Header.TickRate,
		},
		Snapshots: mapSnapshots(replay.Snapshots),
		Ticks:     ticks,
	}
	data, err := rep.Marshal()
	if err != nil {
		return err
	}
	_, err = w.Write(data)
	return err
}

func mapEntities(entities []rep.Entity) []*Replay_Entity {
	result := make([]*Replay_Entity, 0)
	for _, e := range entities {
		var t Team
		switch common.Team(e.Team) {
		case common.Team_Terrorists:
			t = Team_TERRORIST
		case common.Team_CounterTerrorists:
			t = Team_COUNTER_TERRORIST
		case common.Team_Spectators:
			t = Team_SPECTATOR
		default:
			t = Team_UNASSIGNED
		}

		result = append(result, &Replay_Entity{
			Id:    int32(e.ID),
			Team:  t,
			Name:  e.Name,
			IsNpc: e.IsNpc,
		})
	}
	return result
}

func mapSnapshots(snaps []rep.Snapshot) []*Replay_Snapshot {
	result := make([]*Replay_Snapshot, 0)
	for _, s := range snaps {
		result = append(result, &Replay_Snapshot{
			Tick:          int32(s.Tick),
			EntityUpdates: mapEntityUpdates(s.EntityUpdates),
		})
	}
	return result
}

func mapEntityUpdates(entityUpdates []rep.EntityUpdate) []*Replay_Snapshot_EntityUpdate {
	result := make([]*Replay_Snapshot_EntityUpdate, 0)
	for _, u := range entityUpdates {
		result = append(result, &Replay_Snapshot_EntityUpdate{
			Angle:         u.Angle,
			Armor:         int32(u.Armor),
			EntityId:      int32(u.EntityID),
			FlashDuration: u.FlashDuration,
			Hp:            int32(u.Hp),
			Positions:     mapPositions(u.Positions),
		})
	}
	return result
}

func mapPositions(positions []rep.Point) []*Point {
	result := make([]*Point, 0)
	for _, p := range positions {
		result = append(result, mapPosition(p))
	}
	return result
}

func mapPosition(p rep.Point) *Point {
	return &Point{
		X: p.X,
		Y: p.Y,
	}
}

func mapTicks(ticks []rep.Tick) ([]*Replay_Tick, error) {
	result := make([]*Replay_Tick, 0)
	for _, t := range ticks {
		e, err := mapEvents(t.Events)
		if err != nil {
			return nil, err
		}
		result = append(result, &Replay_Tick{
			Nr:     int32(t.Nr),
			Events: e,
		})
	}
	return result, nil
}

func mapEvents(events []rep.Event) ([]*Replay_Tick_Event, error) {
	result := make([]*Replay_Tick_Event, 0)
	for _, e := range events {
		result = append(result, &Replay_Tick_Event{
			Event:      e.Name,
			Attributes: mapAttributes(e.Attributes),
			//Details:    details,
		})
	}
	return result, nil
}

func mapAttributes(attrs []rep.EventAttribute) []*Replay_Tick_Event_Attribute {
	result := make([]*Replay_Tick_Event_Attribute, 0)
	for _, a := range attrs {
		result = append(result, &Replay_Tick_Event_Attribute{
			Key:         a.Key,
			NumberValue: a.NumVal,
			StringValue: a.StrVal,
		})
	}
	return result
}
