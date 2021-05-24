package mycocontext

import (
	"bytes"
	"context"
)

// Key is used for setting and getting from mycocontext.Context in this project.
type Key int

// These are keys for the context that floats around.
const (
	// keyHyphaName is for storing current hypha name as a string here.
	keyHyphaName Key = iota
	// keyInputBuffer is for storing *bytes.Buffer with unread bytes of the source document.
	keyInputBuffer
	// KeyRecursionLevel stores current level of transclusion recursion.
	keyRecursionLevel
	//
	keyWebSiteURL
)

// HyphaName retrieves current hypha name from the given context.
func (ctx *mycoContext) HyphaName() string {
	return ctx.Value(keyHyphaName).(string)
}

// Input retrieves the current bytes buffer from the given context.
func (ctx *mycoContext) Input() *bytes.Buffer {
	return ctx.Value(keyInputBuffer).(*bytes.Buffer)
}

func (ctx *mycoContext) GetRecursionLevel() uint {
	return ctx.Value(keyRecursionLevel).(uint)
}

func (ctx *mycoContext) WithIncrementedRecursionLevel() Context {
	return &mycoContext{context.WithValue(ctx, keyRecursionLevel, ctx.GetRecursionLevel()+1)}
}

func (ctx *mycoContext) WebSiteURL() string {
	return ctx.Value(keyWebSiteURL).(string)
}
