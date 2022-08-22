// Package mycocontext provides a context type for all things Mycomarkup. It also implements the context.Context interface, if that's your thing.
package mycocontext

import (
	"bytes"
	"context"
	"git.sr.ht/~bouncepaw/mycomarkup/v5/options"
)

// Context contains all context related to the current Mycomarkup document.
//
// It is cheap to copy it.
type Context struct {
	*mycoContext
}

// CancelFunc is a function you call to cancel the context.
//
// Why would you, though? It is currently unused. I doubt it works, actually.
type CancelFunc context.CancelFunc

// ContextFromBuffer returns the context for the given input.
func ContextFromBuffer(input *bytes.Buffer, opts options.Options) (Context, CancelFunc) {
	// Why is it like that?
	opts.TransclusionSupported = true

	ctx := context.Background()
	ctx = context.WithValue(ctx, keyOptions, opts)
	ctx = context.WithValue(ctx, keyInputBuffer, input)
	ctx = context.WithValue(ctx, keyRecursionLevel, uint(0))
	ctx, cancel := context.WithCancel(ctx)

	return Context{&mycoContext{ctx}}, CancelFunc(cancel)
}

// ContextFromBytes returns the context for the given input.
func ContextFromBytes(input []byte, opts options.Options) (Context, CancelFunc) {
	return ContextFromBuffer(bytes.NewBuffer(input), opts)
}

// ContextFromStringInput returns the context for the given input.
func ContextFromStringInput(input string, opts options.Options) (Context, CancelFunc) {
	return ContextFromBuffer(bytes.NewBufferString(input), opts)
}
