package csminify

import (
	rep "github.com/markus-wa/cs-demo-minifier/replay" 
	dem "github.com/markus-wa/demoinfocs-golang/v2/pkg/demoinfocs"
	events "github.com/markus-wa/demoinfocs-golang/v2/pkg/demoinfocs/events"
)

// EventCollector provides the possibility of adding custom events to replays.
// First all demo-event handlers must be registered via AddHandler().
// The registered handlers can add replay-events to the collector via AddEvent().
// The handlers can access game-state information via Parser().
// After a tick ends all events that were added to the collector during the tick will be stored into the replay.
type EventCollector struct {
	handlers []interface{}
	events   []rep.Event
	parser   dem.Parser
}

// AddHandler adds a handler which will be registered on the Parser to the collector.
// The handler should use EventCollector.AddEvent() and be of the type
// func(<EventType>) where EventType is the type of the event to be handled.
// The handler parameter is of type interface because lolnogenerics.
// See: github.com/markus-wa/demoinfocs-golang/v2/pkg/demoinfocs.Parser.RegisterEventHandler()
// GoDoc: https://pkg.go.dev/github.com/markus-wa/demoinfocs-golang/v2/pkg/demoinfocs?tab=doc#Parser
func (ec *EventCollector) AddHandler(handler interface{}) {
	ec.handlers = append(ec.handlers, handler)
}

// AddEvent adds an event to the collector.
// The event will be added to the replay after the current tick ends.
func (ec *EventCollector) AddEvent(event rep.Event) {
	ec.events = append(ec.events, event)
}

// Parser returns the demo-parser through which custom handlers can access game-state information.
// Returns nil before minification has started - so don't call this before you need it.
func (ec *EventCollector) Parser() dem.Parser {
	return ec.parser
}

// EventHandlers provides functions for registering the out-of-the-box provided handlers on EventCollectors.
// The handlers are divided into two groups: Default and Extra.
// Default contains the handlers that are used if no custom EventCollector is specified.
// Extra contains other handlers that are usually not required (i.e. a 'footstep' handler).
var EventHandlers eventHandlers

type eventHandlers struct {
	Default defaultEventHandlers
	Extra   extraEventHandlers
}

type defaultEventHandlers struct{}

func (defaultEventHandlers) RegisterAll(ec *EventCollector) {
	EventHandlers.Default.RegisterMatchStarted(ec)
	EventHandlers.Default.RegisterGamePhaseChanged(ec)
	EventHandlers.Default.RegisterRoundStarted(ec)
	EventHandlers.Default.RegisterRoundEnded(ec)
	EventHandlers.Default.RegisterPlayerKilled(ec)
	EventHandlers.Default.RegisterPlayerHurt(ec)
	EventHandlers.Default.RegisterPlayerFlashed(ec)
	EventHandlers.Default.RegisterPlayerJump(ec)
	EventHandlers.Default.RegisterPlayerTeamChange(ec)
	EventHandlers.Default.RegisterPlayerDisconnect(ec)
	EventHandlers.Default.RegisterWeaponFired(ec)
	EventHandlers.Default.RegisterChatMessage(ec)
	EventHandlers.Default.RegisterGrenadeEvents(ec)
}

func (defaultEventHandlers) RegisterMatchStarted(ec *EventCollector) {
	ec.AddHandler(func(e events.MatchStart) {
		ec.AddEvent(createEvent(rep.EventMatchStarted))
	})
}

func (defaultEventHandlers) RegisterGamePhaseChanged(ec *EventCollector) {
	ec.AddHandler(func(e events.GamePhaseChanged) {
		eb := buildEvent(rep.EventGamePhaseChanged)
		eb.intAttr("oldGamePhase", int(e.OldGamePhase))
		eb.intAttr("newGamePhase", int(e.NewGamePhase))
		ec.AddEvent(eb.build())
	})
}

func (defaultEventHandlers) RegisterRoundStarted(ec *EventCollector) {
	ec.AddHandler(func(e events.RoundStart) {
		ec.AddEvent(createEvent(rep.EventRoundStarted))
	})
}

func (defaultEventHandlers) RegisterRoundEnded(ec *EventCollector) {
	ec.AddHandler(func(e events.RoundEnd) {
		eb := buildEvent(rep.EventRoundEnded)
		eb.intAttr("winner", int(e.Winner))
		eb.intAttr("reason", int(e.Reason))
		ec.AddEvent(eb.build())
	})
}

func (defaultEventHandlers) RegisterPlayerKilled(ec *EventCollector) {
	ec.AddHandler(func(e events.Kill) {
		eb := buildEvent(rep.EventKill)
		eb.intAttr(rep.AttrKindVictim, e.Victim.EntityID)
		eb.intAttr(rep.AttrKindWeapon, int(e.Weapon.Type))

		if e.Killer != nil {
			eb.intAttr(rep.AttrKindKiller, e.Killer.EntityID)
		}

		if e.Assister != nil {
			eb.intAttr(rep.AttrKindAssister, e.Assister.EntityID)
		}

		ec.AddEvent(eb.build())
	})
}

func (defaultEventHandlers) RegisterPlayerHurt(ec *EventCollector) {
	ec.AddHandler(func(e events.PlayerHurt) {
		ec.AddEvent(createEntityEvent(rep.EventHurt, e.Player.EntityID))
	})
}

func (defaultEventHandlers) RegisterPlayerFlashed(ec *EventCollector) {
	ec.AddHandler(func(e events.PlayerFlashed) {
		ec.AddEvent(createEntityEvent(rep.EventFlashed, e.Player.EntityID))
	})
}

func (defaultEventHandlers) RegisterPlayerJump(ec *EventCollector) {
	ec.AddHandler(func(e events.PlayerJump) {
		if e.Player == nil {
			return
		}

		ec.AddEvent(createEntityEvent(rep.EventJump, e.Player.EntityID))
	})
}

func (defaultEventHandlers) RegisterPlayerTeamChange(ec *EventCollector) {
	ec.AddHandler(func(e events.PlayerTeamChange) {
		if e.Player == nil {
			return
		}

		ec.AddEvent(createEntityEvent(rep.EventSwapTeam, e.Player.EntityID))
	})
}

func (defaultEventHandlers) RegisterPlayerDisconnect(ec *EventCollector) {
	ec.AddHandler(func(e events.PlayerDisconnected) {
		ec.AddEvent(createEntityEvent(rep.EventDisconnect, e.Player.EntityID))
	})
}

func (defaultEventHandlers) RegisterWeaponFired(ec *EventCollector) {
	ec.AddHandler(func(e events.WeaponFire) {
		ec.AddEvent(createEntityEvent(rep.EventFire, e.Shooter.EntityID))
	})
}

func (defaultEventHandlers) RegisterChatMessage(ec *EventCollector) {
	ec.AddHandler(func(e events.ChatMessage) {
		eb := buildEvent(rep.EventChatMessage)
		eb.stringAttr(rep.AttrKindText, e.Text)

		// Skip for now, probably always true anyway
		//eb.boolAttr("isChatAll", e.IsChatAll)

		if e.Sender != nil {
			eb.intAttr(rep.AttrKindSender, e.Sender.EntityID)
		}

		ec.AddEvent(eb.build())
	})
}

func (defaultEventHandlers) RegisterGrenadeEvents(ec *EventCollector) {
  ec.AddHandler(func(e events.SmokeStart) {
    eb := buildEvent(rep.EventSmokeStart)
    eb.floatAttr("x", e.Position.X)
    eb.floatAttr("y", e.Position.Y)
    eb.floatAttr("z", e.Position.Z)
    eb.intAttr(rep.AttrKindEntityID, e.Thrower.EntityID);
    ec.AddEvent(eb.build())
  })
  ec.AddHandler(func(e events.SmokeExpired) {
    eb := buildEvent(rep.EventSmokeExpired)
    eb.floatAttr("x", e.Position.X)
    eb.floatAttr("y", e.Position.Y)
    eb.floatAttr("z", e.Position.Z)
    eb.intAttr(rep.AttrKindEntityID, e.Thrower.EntityID);
    ec.AddEvent(eb.build())
  })
  ec.AddHandler(func(e events.DecoyStart) {
    eb := buildEvent(rep.EventDecoyStart)
    eb.floatAttr("x", e.Position.X)
    eb.floatAttr("y", e.Position.Y)
    eb.floatAttr("z", e.Position.Z)
    eb.intAttr(rep.AttrKindEntityID, e.Thrower.EntityID);
    ec.AddEvent(eb.build())
  })
  ec.AddHandler(func(e events.DecoyExpired) {
    eb := buildEvent(rep.EventDecoyExpired)
    eb.floatAttr("x", e.Position.X)
    eb.floatAttr("y", e.Position.Y)
    eb.floatAttr("z", e.Position.Z)
    eb.intAttr(rep.AttrKindEntityID, e.Thrower.EntityID);
    ec.AddEvent(eb.build())
  })
  ec.AddHandler(func(e events.FireGrenadeStart) {
    eb := buildEvent(rep.EventFireGrenadeStart)
    eb.floatAttr("x", e.Position.X)
    eb.floatAttr("y", e.Position.Y)
    eb.floatAttr("z", e.Position.Z)
    //eb.intAttr(rep.AttrKindEntityID, e.Thrower.EntityID);
    ec.AddEvent(eb.build())
  })
  ec.AddHandler(func(e events.FireGrenadeExpired) {
    eb := buildEvent(rep.EventFireGrenadeExpired)
    eb.floatAttr("x", e.Position.X)
    eb.floatAttr("y", e.Position.Y)
    eb.floatAttr("z", e.Position.Z)
    //eb.intAttr(rep.AttrKindEntityID, e.Thrower.EntityID);
    ec.AddEvent(eb.build())
  })
  ec.AddHandler(func(e events.HeExplode) {
    eb := buildEvent(rep.EventHEGrenadeExplosion)
    eb.floatAttr("x", e.Position.X)
    eb.floatAttr("y", e.Position.Y)
    eb.floatAttr("z", e.Position.Z)
    eb.intAttr(rep.AttrKindEntityID, e.Thrower.EntityID);
    ec.AddEvent(eb.build())
  })
  ec.AddHandler(func(e events.FlashExplode) {
    eb := buildEvent(rep.EventFlashExplosion)
    eb.floatAttr("x", e.Position.X)
    eb.floatAttr("y", e.Position.Y)
    eb.floatAttr("z", e.Position.Z)
    eb.intAttr(rep.AttrKindEntityID, e.Thrower.EntityID);
    ec.AddEvent(eb.build())
  })
}

type extraEventHandlers struct{}

func (extraEventHandlers) RegisterAll(ec *EventCollector) {
	EventHandlers.Extra.RegisterFootstep(ec)
}

func (extraEventHandlers) RegisterFootstep(ec *EventCollector) {
	ec.AddHandler(func(e events.Footstep) {
		ec.AddEvent(createEntityEvent(rep.EventFootstep, e.Player.EntityID))
	})
}

type eventBuilder struct {
	event rep.Event
}

func (b *eventBuilder) stringAttr(key string, value string) *eventBuilder {
	b.event.Attributes = append(b.event.Attributes, rep.EventAttribute{
		Key:    key,
		StrVal: value,
	})
	return b
}

func (b *eventBuilder) intAttr(key string, value int) *eventBuilder {
	b.event.Attributes = append(b.event.Attributes, rep.EventAttribute{
		Key:    key,
		NumVal: float64(value),
	})
	return b
}

func (b *eventBuilder) floatAttr(key string, value float64) *eventBuilder {
  b.event.Attributes = append(b.event.Attributes, rep.EventAttribute{
    Key:    key,
    NumVal: value,
  })
  return b
}

func (b eventBuilder) build() rep.Event {
	return b.event
}

func buildEvent(eventName string) *eventBuilder {
	return &eventBuilder{
		event: createEvent(eventName),
	}
}

func createEvent(eventName string) rep.Event {
	return rep.Event{
		Name: eventName,
	}
}

func createEntityEvent(eventName string, entityID int) rep.Event {
	return buildEvent(eventName).intAttr(rep.AttrKindEntityID, entityID).build()
}
