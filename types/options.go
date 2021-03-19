package main

// Options is an object holding all custom stuff.
type Options struct {
	// CheckForExistence is a function that gets a map, where keys are canonical hypha names and values are true if the given hypha exists. If the function is not passed, consider that all hyphae exist.
	CheckForExistence func(map[string]bool)
}

func DefaultCheckForExistence(m map[string]bool) {
	for k, _ := range m {
		m[k] = true
	}
}
