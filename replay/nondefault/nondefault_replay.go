// Package nondefault is used for testing purposes only.
// It allows making sure the marshalling & unmarshalling of every field is
// tested by checking if the replay contains any default-values (e.g. 0 for int or "" for string etc.)
package nondefault

import (
	"reflect"
	"unsafe"

	rep "github.com/markus-wa/cs-demo-minifier/replay"
)

var nonDefaultReplay rep.Replay

// GetNonDefaultReplay returns a Replay that doesn't contain any default values.
// Used to test data preservation during marshalling & unmarshalling
func GetNonDefaultReplay() rep.Replay {
	return nonDefaultReplay
}

func init() {
	var ent []rep.Entity
	ent = append(ent, rep.Entity{
		ID:    5,
		IsNpc: true,
		Name:  "Batman",
		Team:  2,
	})

	var pos []rep.Point
	pos = append(pos, rep.Point{
		X: 1124,
		Y: -321,
		Z: 24,
	})

	var entUpd []rep.EntityUpdate
	entUpd = append(entUpd, rep.EntityUpdate{
		AngleX:        90,
		AngleY:        45,
		Armor:         80,
		EntityID:      5,
		FlashDuration: 2.35,
		Hp:            100,
		IsNpc:         true,
		Positions:     pos,
		Team:          1,
		HasDefuseKit:  true,
		HasHelmet:     true,
		Equipment: []rep.EntityEquipment{
			{
				Type: 1,
			},
		}})

	var snaps []rep.Snapshot
	snaps = append(snaps, rep.Snapshot{
		Tick:          1,
		EntityUpdates: entUpd,
	})

	var attrs []rep.EventAttribute
	attrs = append(attrs, rep.EventAttribute{
		Key:    rep.AttrKindEntityID,
		NumVal: 5,
		StrVal: "test",
	})
	attrs = append(attrs, rep.EventAttribute{
		Key:    "custom_attr",
		NumVal: 10,
		StrVal: "custom_val",
	})

	var events []rep.Event
	events = append(events, rep.Event{
		Name:       rep.EventJump,
		Attributes: attrs,
	})
	events = append(events, rep.Event{
		Name:       "custom_event",
		Attributes: attrs,
	})

	var ticks []rep.Tick
	ticks = append(ticks, rep.Tick{
		Nr:     5,
		Events: events,
	})

	replay := rep.Replay{
		Header: rep.Header{
			MapName:      "de_test",
			SnapshotRate: 64,
			TickRate:     128,
		},
		Entities:  ent,
		Snapshots: snaps,
		Ticks:     ticks,
	}

	// Check for nested default values in the testdata.
	// Default values could lead to false positive marshalling / unmarshalling tests.
	if !deepNonDefault(replay) {
		panic("Marshalling / unmarshalling test data (replay) contains default values")
	}
	nonDefaultReplay = replay
}

// TODO: deepNonDefault should build up a tree and allow default values
// as long as the field isn't set to default everywhere in the tree.

// deepNonDefault reports whether x and y are ``deeply un-equal,''
// basically the two values must have the same structure but different values.
// Code was taken form reflect/deepequal.go and modified.
func deepNonDefault(x interface{}) bool {
	if x == nil {
		return false
	}
	v := reflect.ValueOf(x)
	return deepValueNonDefault(v)
}

// During deepValueUnEqual, must keep track of checks that are
// in progress. The comparison algorithm assumes that all
// checks in progress are true when it reencounters them.
// Visited comparisons are stored in a map indexed by visit.
type visit struct {
	a1  unsafe.Pointer
	a2  unsafe.Pointer
	typ reflect.Type
}

// Tests for deep unequality using reflected types. The map argument tracks
// comparisons that have already been seen, which allows short circuiting on
// recursive types.
func deepValueNonDefault(v reflect.Value) bool {
	if !v.IsValid() {
		return false
	}

	switch v.Kind() {
	case reflect.Array:
		for i := 0; i < v.Len(); i++ {
			if !deepValueNonDefault(v.Index(i)) {
				return false
			}
		}
		return true

	case reflect.Slice:
		fallthrough
	case reflect.Map:
		if v.IsNil() {
			return false
		}
		if v.Len() == 0 {
			return false
		}
		for i := 0; i < v.Len(); i++ {
			if !deepValueNonDefault(v.Index(i)) {
				return false
			}
		}
		return true

	case reflect.Interface:
		if v.IsNil() {
			return false
		}
		return deepValueNonDefault(v.Elem())

	case reflect.Ptr:
		return deepValueNonDefault(v.Elem())

	case reflect.Struct:
		for i, n := 0, v.NumField(); i < n; i++ {
			if !deepValueNonDefault(v.Field(i)) {
				return false
			}
		}

		// Shouldn't be the default value
		return deepValueUnEqual(v, reflect.New(v.Type()).Elem(), make(map[visit]bool), 0)

	case reflect.Func:
		if v.IsNil() {
			return false
		}
		// Can't do better than this:
		return false

	default:
		// Normal inequality suffices
		return v.Interface() != nil
	}
}

// Tests for deep unequality using reflected types. The map argument tracks
// comparisons that have already been seen, which allows short circuiting on
// recursive types.
func deepValueUnEqual(v1, v2 reflect.Value, visited map[visit]bool, depth int) bool {
	if !v1.IsValid() || !v2.IsValid() {
		return false
	}
	if v1.Type() != v2.Type() {
		return false
	}

	// if depth > 10 { panic("deepValueUnEqual") }	// for debugging

	// We want to avoid putting more in the visited map than we need to.
	// For any possible reference cycle that might be encountered,
	// hard(t) needs to return true for at least one of the types in the cycle.
	hard := func(k reflect.Kind) bool {
		switch k {
		case reflect.Map, reflect.Slice, reflect.Ptr, reflect.Interface:
			return true
		}
		return false
	}

	if v1.CanAddr() && v2.CanAddr() && hard(v1.Kind()) {
		addr1 := unsafe.Pointer(v1.UnsafeAddr())
		addr2 := unsafe.Pointer(v2.UnsafeAddr())
		if uintptr(addr1) > uintptr(addr2) {
			// Canonicalize order to reduce number of entries in visited.
			// Assumes non-moving garbage collector.
			addr1, addr2 = addr2, addr1
		}

		// Short circuit if references are already seen.
		typ := v1.Type()
		v := visit{addr1, addr2, typ}
		if visited[v] {
			return true
		}

		// Remember for later.
		visited[v] = true
	}

	switch v1.Kind() {
	case reflect.Array:
		for i := 0; i < v1.Len(); i++ {
			if !deepValueUnEqual(v1.Index(i), v2.Index(i), visited, depth+1) {
				return false
			}
		}
		return true

	case reflect.Slice:
		fallthrough
	case reflect.Map:
		if v1.IsNil() || v2.IsNil() {
			return v1.IsNil() != v2.IsNil()
		}
		if v1.Pointer() == v2.Pointer() {
			return false
		}
		if v1.Len() == v2.Len() {
			return false
		}
		return true

	case reflect.Interface:
		if v1.IsNil() || v2.IsNil() {
			return v1.IsNil() != v2.IsNil()
		}
		return deepValueUnEqual(v1.Elem(), v2.Elem(), visited, depth+1)

	case reflect.Ptr:
		if v1.Pointer() == v2.Pointer() {
			return false
		}
		return deepValueUnEqual(v1.Elem(), v2.Elem(), visited, depth+1)

	case reflect.Struct:
		for i, n := 0, v1.NumField(); i < n; i++ {
			if !deepValueUnEqual(v1.Field(i), v2.Field(i), visited, depth+1) {
				return false
			}
		}
		return true

	case reflect.Func:
		if v1.IsNil() || v2.IsNil() {
			return v1.IsNil() != v2.IsNil()
		}
		// Can't do better than this:
		return false

	default:
		// Normal inequality suffices
		return v1.Interface() != v2.Interface()
	}
}
