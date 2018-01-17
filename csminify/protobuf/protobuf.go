// Package protobuf provides a replay marshaller for protobuf.
// Use 'go generate' to generate the code from the .proto files inside the proto sub directory.
package protobuf

// -I=proto is required, otherwise the generated .pb.go file will be put inside the proto directory.
// No idea what that is about to be honest . . .
//go:generate protoc -I=proto --gogofaster_out=. proto/replay.proto

import (
	bimap "github.com/vishalkuo/bimap"
)

var (
	attributeKindMap = bimap.NewBiMap()
	eventKindMap     = bimap.NewBiMap()
)

func init() {
	attributeKindMap.Insert("entityId", Replay_Tick_Event_Attribute_ENTITY_ID)
	attributeKindMap.Insert("victim", Replay_Tick_Event_Attribute_VICTIM)
	attributeKindMap.Insert("killer", Replay_Tick_Event_Attribute_KILLER)
	attributeKindMap.Insert("assister", Replay_Tick_Event_Attribute_ASSISTER)

	eventKindMap.Insert("jump", Replay_Tick_Event_JUMP)
	eventKindMap.Insert("fire", Replay_Tick_Event_FIRE)
	eventKindMap.Insert("hurt", Replay_Tick_Event_HURT)
	eventKindMap.Insert("kill", Replay_Tick_Event_KILL)
	eventKindMap.Insert("flashed", Replay_Tick_Event_FLASHED)
	eventKindMap.Insert("round_started", Replay_Tick_Event_ROUND_STARTED)
	eventKindMap.Insert("swap_team", Replay_Tick_Event_SWAP_TEAM)
	eventKindMap.Insert("disconnect", Replay_Tick_Event_DISCONNECT)
}
