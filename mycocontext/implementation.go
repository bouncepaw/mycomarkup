package mycocontext

import (
	"bytes"
	"context"
	"github.com/bouncepaw/mycomarkup/v5/options"
)

// I'm very well aware that storing context.Context inside structs is discouraged in most cases. But it should be ok this time.
type mycoContext struct {
	context.Context
}

// See interface.go for description of the methods.

func (ctx *mycoContext) HyphaName() string {
	return ctx.Value(keyOptions).(options.Options).HyphaName
}

func (ctx *mycoContext) Input() *bytes.Buffer {
	return ctx.Value(keyInputBuffer).(*bytes.Buffer)
}

func (ctx *mycoContext) RecursionLevel() uint {
	return ctx.Value(keyRecursionLevel).(uint)
}

func (ctx *mycoContext) WithIncrementedRecursionLevel() Context {
	return &mycoContext{context.WithValue(ctx, keyRecursionLevel, ctx.RecursionLevel()+1)}
}

func (ctx *mycoContext) WebSiteURL() string {
	return ctx.Value(keyOptions).(options.Options).WebSiteURL
}

func (ctx *mycoContext) TransclusionSupported() bool {
	return ctx.Value(keyOptions).(options.Options).TransclusionSupported
}
