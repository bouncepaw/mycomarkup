package mycocontext

import (
	"bytes"
	"context"
	"git.sr.ht/~bouncepaw/mycomarkup/v5/options"
)

// I'm very well aware that storing context.Context inside structs is discouraged in most cases. But it should be ok this time.
type mycoContext struct {
	context.Context
}

// Options returns Options for how to parse and render the document.
func (ctx *mycoContext) Options() options.Options {
	return ctx.Value(keyOptions).(options.Options)
}

// HyphaName returns the name of the Mycomarkup document being parsed.
//
// This method is the same as:. It is just used so much.
//
//     ctx.Options().HyphaName
func (ctx *mycoContext) HyphaName() string {
	// TODO: Rename it to DocumentName when the time comes, because Mycomarkup is not Mycorrhiza only.
	return ctx.Options().HyphaName
}

// Input returns the pointer to the input buffer.
func (ctx *mycoContext) Input() *bytes.Buffer {
	return ctx.Value(keyInputBuffer).(*bytes.Buffer)
}

// RecursionLevel returns the current recursion level. The recursion level can be increased by WithIncrementedRecursionLevel.
func (ctx *mycoContext) RecursionLevel() uint {
	return ctx.Value(keyRecursionLevel).(uint)
}

// WithIncrementedRecursionLevel returns a copy of the context, except it has an incremented recursion level. Use that in translcusion, and check so that it's not really high.
func (ctx *mycoContext) WithIncrementedRecursionLevel() Context {
	return Context{&mycoContext{context.WithValue(ctx, keyRecursionLevel, ctx.RecursionLevel()+1)}}
}

// WithBuffer returns a copy of the given context but with a different input buffer.
func (ctx *mycoContext) WithBuffer(buf *bytes.Buffer) Context {
	return Context{&mycoContext{context.WithValue(ctx, keyInputBuffer, buf)}}
}

// WithOptions returns a copy of the given copy of the given context but with different options.
func (ctx *mycoContext) WithOptions(opts options.Options) Context {
	return Context{&mycoContext{context.WithValue(ctx, keyOptions, opts)}}
}

type key int

const (
	keyInputBuffer key = iota
	keyRecursionLevel
	keyOptions
)
