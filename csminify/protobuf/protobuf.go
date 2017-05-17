package protobuf

//go:generate protoc -I=proto --gogofaster_out=Mgoogle/protobuf/descriptor.proto=github.com/gogo/protobuf/protoc-gen-gogo/descriptor:. proto/*.proto

import (
	rep "github.com/markus-wa/cs-demo-minifier/csminify/replay"
	"github.com/markus-wa/demoinfocs-golang/common"
	"io"
)

func MarshalReplay(replay rep.Replay, w io.Writer) error {
	rep := Replay{
		Entities: mapEntities(replay.Entities),
		Header: &Replay_Header{
			Map:          replay.Header.MapName,
			SnapshotRate: int32(replay.Header.SnapshotRate),
			TickRate:     replay.Header.TickRate,
		},
		Snapshots: mapSnapshots(replay.Snapshots),
		Ticks:     mapTicks(replay.Ticks),
	}
	data, err := rep.Marshal()
	w.Write(data)
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
	return nil
}

func mapTicks(ticks []rep.Tick) []*Replay_Tick {
	return nil
}
