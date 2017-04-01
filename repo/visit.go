package repo

const (
	// Adds indicates the addition files in this visit.
	Adds = iota
	// Mods indicates the modification files in this visit.
	Mods
)

// Visitor post process the found documents.
type Visitor interface {
	Visit(docs Docs, cookie map[int]interface{})
}
