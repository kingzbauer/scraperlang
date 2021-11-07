package interpreter

// Accessor defines an interface for retrieving object attributes.
// Certain runtime objects implement this interface e.g arrays and maps to allow
// for looping across it's elements.
type Accessor interface {
	Get(string) interface{}
}

// Keyer interface allows for indexing into an array or map by index and key respectively
type Keyer interface {
	GetValue(key interface{}) interface{}
}
