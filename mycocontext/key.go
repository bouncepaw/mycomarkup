package mycocontext

// key is used for setting and getting from mycocontext.Context in this project.
type key int

// These are keys for the context that floats around.
const (
	// keyHyphaName is for storing current hypha name as a string here.
	keyHyphaName key = iota
	// keyInputBuffer is for storing *bytes.Buffer with unread bytes of the source document.
	keyInputBuffer
	// KeyRecursionLevel stores current level of transclusion recursion.
	keyRecursionLevel
	//
	keyWebSiteURL
)
