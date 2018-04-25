// Package protobuf provides a replay marshaller for protobuf.
// The generated code is located in the gen sub-package.
package protobuf

import (
	bimap "github.com/vishalkuo/bimap"

	gen "github.com/markus-wa/cs-demo-minifier/protobuf/gen"
	rep "github.com/markus-wa/cs-demo-minifier/replay"
)

const (
	attrKindEventName = "INTERNAL_EVENT_NAME"
)

var (
	attributeKindMap = bimap.NewBiMap()
	eventKindMap     = bimap.NewBiMap()
)

func init() {
	attributeKindMap.Insert(rep.AttrKindEntityID, gen.Replay_Tick_Event_Attribute_ENTITY_ID)
	attributeKindMap.Insert(rep.AttrKindVictim, gen.Replay_Tick_Event_Attribute_VICTIM)
	attributeKindMap.Insert(rep.AttrKindKiller, gen.Replay_Tick_Event_Attribute_KILLER)
	attributeKindMap.Insert(rep.AttrKindAssister, gen.Replay_Tick_Event_Attribute_ASSISTER)
	attributeKindMap.Insert(rep.AttrKindText, gen.Replay_Tick_Event_Attribute_TEXT)
	attributeKindMap.Insert(attrKindEventName, gen.Replay_Tick_Event_Attribute_EVENT_NAME)

	eventKindMap.Insert(rep.EventJump, gen.Replay_Tick_Event_JUMP)
	eventKindMap.Insert(rep.EventFire, gen.Replay_Tick_Event_FIRE)
	eventKindMap.Insert(rep.EventHurt, gen.Replay_Tick_Event_HURT)
	eventKindMap.Insert(rep.EventKill, gen.Replay_Tick_Event_KILL)
	eventKindMap.Insert(rep.EventFlashed, gen.Replay_Tick_Event_FLASHED)
	eventKindMap.Insert(rep.EventRoundStarted, gen.Replay_Tick_Event_ROUND_STARTED)
	eventKindMap.Insert(rep.EventSwapTeam, gen.Replay_Tick_Event_SWAP_TEAM)
	eventKindMap.Insert(rep.EventDisconnect, gen.Replay_Tick_Event_DISCONNECT)
	eventKindMap.Insert(rep.EventChatMessage, gen.Replay_Tick_Event_CHAT_MESSAGE)
}
