// Package mycocontext provides a wrapper over context.Context and some operations on the wrapper.
package mycocontext

import (
	"bytes"
	"context"
)

// Context is the wrapper around context.Context providing type-level safety on presence of several values.
type Context interface {
	context.Context

	// HyphaName returns the name of the processed hypha.
	HyphaName() string

	// Input returns the buffer which contains all characters of the hypha text.
	Input() *bytes.Buffer

	// RecursionLevel returns current recursive transclusion level.
	RecursionLevel() uint

	// WithIncrementedRecursionLevel returns a copy of the context but with the recursion level incremented.
	//
	//     lvl1 := ctx.RecursionLevel()
	//     lvl2 := ctx.WithIncrementedRecursionLevel().RecursionLevel()
	//     lvl2 - lvl1 == 1
	WithIncrementedRecursionLevel() Context

	// WebSiteURL returns the URL of the wiki, including the protocol (http or https). It is used for generating OpenGraph meta tags.
	WebSiteURL() string
}

type CancelFunc context.CancelFunc

// ContextFromStringInput returns the context for the given input.
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

// WithBuffer returns a copy of the given context but with a different input buffer.
func WithBuffer(ctx Context, buf *bytes.Buffer) Context {
	return &mycoContext{context.WithValue(ctx, keyInputBuffer, buf)}
}
