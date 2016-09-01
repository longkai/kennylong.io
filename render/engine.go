package render

import "os"

// Engine render engine.
type Engine interface {
	// Ls list entries.
	Ls(key string, size int) (interface{}, error)
	// Get an entry.
	Get(key string) (interface{}, error)
	// Post the given key.
	Post(key string) (interface{}, error)
	// Put an entry.
	Put(key string) (interface{}, error)
	// Del an entry.
	Del(key string) (interface{}, error)
}

// Traveller a traveller travel somewhere to meet sth. interesting.
type Traveller interface {
	// Travel a place.
	Travel(place string)
	// Meet we only meet funny things.
	Meet(sth string)
	// Fun is it?
	Fun(place string, sth os.FileInfo) bool
}
