package csminify

import (
	"github.com/golang/geo/r3"
	rep "github.com/markus-wa/cs-demo-minifier/csminify/replay"
)

func r3VectorToPoint(v r3.Vector) rep.Point {
	return rep.Point{X: v.X, Y: v.Y}
}

func createGameEvent(eventName string, eventString string) rep.GameEvent {
	return rep.GameEvent{
		Event: rep.Event{
			Name: eventName,
		},
		EventString: eventString,
	}
}

func createEntityEvent(eventName string, entityID int) rep.EntityEvent {
	return rep.EntityEvent{
		Event: rep.Event{
			Name: eventName,
		},
		EntityID: entityID,
	}
}
