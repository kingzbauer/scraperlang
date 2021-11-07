package interpreter

import "fmt"

// Map is a runtime hash map implementation
type Map struct {
	instance map[string]interface{}
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
