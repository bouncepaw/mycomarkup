package mycocontext

// key is used for setting and getting from mycocontext.Context in this project.
type key int

// These are keys for the context that floats around.
const (
	// keyInputBuffer is for storing *bytes.Buffer with unread bytes of the source document.
	keyInputBuffer key = iota
	// KeyRecursionLevel stores current level of transclusion recursion.
	keyRecursionLevel
	keyOptions
)
