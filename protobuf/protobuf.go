// Package protobuf provides a replay marshaller for protobuf.
// The generated code is located in the gen sub-package.
package protobuf

import (
	bimap "github.com/vishalkuo/bimap"

	gen "github.com/markus-wa/cs-demo-minifier/protobuf/gen"
)

var (
	attributeKindMap = bimap.NewBiMap()
	eventKindMap     = bimap.NewBiMap()
)

func init() {
	attributeKindMap.Insert("entityId", gen.Replay_Tick_Event_Attribute_ENTITY_ID)
	attributeKindMap.Insert("victim", gen.Replay_Tick_Event_Attribute_VICTIM)
	attributeKindMap.Insert("killer", gen.Replay_Tick_Event_Attribute_KILLER)
	attributeKindMap.Insert("assister", gen.Replay_Tick_Event_Attribute_ASSISTER)

	eventKindMap.Insert("jump", gen.Replay_Tick_Event_JUMP)
	eventKindMap.Insert("fire", gen.Replay_Tick_Event_FIRE)
	eventKindMap.Insert("hurt", gen.Replay_Tick_Event_HURT)
	eventKindMap.Insert("kill", gen.Replay_Tick_Event_KILL)
	eventKindMap.Insert("flashed", gen.Replay_Tick_Event_FLASHED)
	eventKindMap.Insert("round_started", gen.Replay_Tick_Event_ROUND_STARTED)
	eventKindMap.Insert("swap_team", gen.Replay_Tick_Event_SWAP_TEAM)
	eventKindMap.Insert("disconnect", gen.Replay_Tick_Event_DISCONNECT)
}
