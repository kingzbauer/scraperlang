package interpreter

import "fmt"

// Map is a runtime hash map implementation
type Map struct {
	instance map[string]interface{}
}

// Call treats Map as a callable and allows looping through the entries
func (m *Map) Call(args ...interface{}) interface{} {
	fun := args[0]
	// fun needs to be a callable
	if call, ok := fun.(Callable); ok {
		// If the arity of the callable is one, we only pass the value
		if call.Arity() == 1 {
			for _, value := range m.instance {
				call.Call(value)
			}
		} else if call.Arity() == 2 {
			for key, value := range m.instance {
				call.Call(key, value)
			}
		} else {
			panic(Error{
				msg: fmt.Sprintf("Map `loop` accepts a callable with arity 1 or 2, got %d", call.Arity()),
			})
		}
		return nil
	}

	panic(Error{
		msg: "Map 'loop' expects a callable as it's only argument",
	})
}

// Arity implements the callable interface for map. Supports arity of 1
func (m *Map) Arity() int {
	return 1
}

// Get allow map to implement the Accessor interface. Map provides a set of predefined attributes
//
//	1. `loop callable` which allows to loop across the entries of the map by key value
func (m *Map) Get(attr string) interface{} {
	switch attr {
	case "loop":
		return m
	default:
		panic(Error{
			msg: fmt.Sprintf("Map does not have an attribute %q", attr),
		})
	}
}

// GetValue returns the value pointed to by the key. Expects the key to be a string.
// Returns nil incase of a missing key
func (m *Map) GetValue(key interface{}) interface{} {
	keyString, ok := key.(string)
	if !ok {
		panic(Error{
			msg: "Expected a string as a map key",
		})
	}
	return m.instance[keyString]
}

func (m *Map) String() string {
	return fmt.Sprintf("#Map %v", m.instance)
}

// Array is a runtime list implementation
type Array struct {
	entries []interface{}
}

// Get implements the Accessor interface for array.
func (a *Array) Get(attr string) interface{} {
	switch attr {
	case "loop":
		return a
	case "size":
		return len(a.entries)
	case "last":
		if len(a.entries) == 0 {
			return nil
		}
		return a.entries[len(a.entries)-1]
	case "first":
		if len(a.entries) == 0 {
			return nil
		}
		return a.entries[0]
	default:
		panic(Error{
			msg: fmt.Sprintf("Array does not have an attribute %q", attr),
		})
	}
}

// Call allow a array to be callable for looping purposes
func (a *Array) Call(args ...interface{}) interface{} {
	fun := args[0]
	// fun needs to be a callable
	if call, ok := fun.(Callable); ok {
		// If the arity of the callable is one, we only pass the value
		if call.Arity() == 1 {
			for _, value := range a.entries {
				call.Call(value)
			}
		} else if call.Arity() == 2 {
			for key, value := range a.entries {
				call.Call(key, value)
			}
		} else {
			panic(Error{
				msg: fmt.Sprintf("Array `loop` accepts a callable with arity 1 or 2, got %d", call.Arity()),
			})
		}
		return nil
	}

	panic(Error{
		msg: "Array 'loop' expects a callable as it's only argument",
	})
}

// Arity returns the number of arguments array requires when treated as a callable
func (a *Array) Arity() int {
	return 1
}

// GetValue indexes into the list
func (a *Array) GetValue(key interface{}) interface{} {
	index, ok := key.(float64)
	if !ok {
		panic(Error{
			msg: "Expected an int as a list index",
		})
	}
	return a.entries[int(index)]
}

func (a *Array) String() string {
	return fmt.Sprintf("#Array %v", a.entries)
}
