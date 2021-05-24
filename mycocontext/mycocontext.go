package mycocontext

import (
	"bytes"
	"context"
)

// Context is a wrapper around context.Context providing type-level safety on presence of several values.
type Context interface {
	context.Context
	// HyphaName returns the name of the processed hypha.
	HyphaName() string

	// Input returns the buffer which contains all characters of the hypha text.
	Input() *bytes.Buffer

	GetRecursionLevel() uint
	WithIncrementedRecursionLevel() Context

	WebSiteURL() string
}

type CancelFunc context.CancelFunc

// I'm very well aware that storing context.Context inside structs is discouraged in most cases. But it should be ok this time.
type mycoContext struct {
	context.Context
}

// ContextFromStringInput returns the context for the given inputFrom.
func ContextFromStringInput(hyphaName, input string) (Context, CancelFunc) {
	ctx, cancel := context.WithCancel(
		context.WithValue(
			context.WithValue(
				context.WithValue(
					context.Background(),
					keyHyphaName,
					hyphaName),
				keyInputBuffer,
				bytes.NewBufferString(input),
			),
			keyRecursionLevel,
			0,
		),
	)
	return &mycoContext{ctx}, CancelFunc(cancel)
}

func WithBuffer(ctx Context, buf *bytes.Buffer) Context {
	return &mycoContext{context.WithValue(ctx, keyInputBuffer, buf)}
}
