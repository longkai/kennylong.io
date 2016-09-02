package render

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
	// Revalidate the given entries.
	Revalidate(adds, mods, dels []string) error
}
