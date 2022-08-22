// Package mycocontext provides a wrapper over context.Context and some operations on the wrapper.
package mycocontext

import (
	"bytes"
	"context"
	"git.sr.ht/~bouncepaw/mycomarkup/v5/options"
)

// Context is the wrapper around context.Context providing type-level safety on presence of several values.
//
// TODO: Get rid of the interface.
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

	// TransclusionSupported is true if Mycomarkup is invoked as a stand-alone program, false if used as a library.
	TransclusionSupported() bool
}

func Options(ctx Context) options.Options {
	return ctx.Value(keyOptions).(options.Options)
}

// CancelFunc is a function you call to cancel the context. Why would you, though? It is currently unused.
type CancelFunc context.CancelFunc

// ContextFromBuffer returns the context for the given input.
func ContextFromBuffer(input *bytes.Buffer, opts options.Options) (Context, CancelFunc) {
	opts.TransclusionSupported = true
	ctx, cancel := context.WithCancel(
		context.WithValue(
			context.WithValue(
				context.WithValue(
					context.Background(),
					keyOptions,
					opts),
				keyInputBuffer,
				input,
			),
			keyRecursionLevel,
			uint(0),
		))
	return &mycoContext{ctx}, CancelFunc(cancel)
}

// ContextFromBytes returns the context for the given input.
func ContextFromBytes(input []byte, opts options.Options) (Context, CancelFunc) {
	return ContextFromBuffer(bytes.NewBuffer(input), opts)
}

// ContextFromStringInput returns the context for the given input.
func ContextFromStringInput(input string, opts options.Options) (Context, CancelFunc) {
	return ContextFromBuffer(bytes.NewBufferString(input), opts)
}

// WithBuffer returns a copy of the given context but with a different input buffer.
func WithBuffer(ctx Context, buf *bytes.Buffer) Context {
	return &mycoContext{context.WithValue(ctx, keyInputBuffer, buf)}
}

func WithOptions(ctx Context, opts options.Options) Context {
	return &mycoContext{context.WithValue(ctx, keyOptions, opts)}
}

// TODO: get rid of these three below

func HyphaExists(ctx Context, hyphaName string) bool {
	return ctx.Value(keyOptions).(options.Options).HyphaExists(hyphaName)
}

func IterateHyphaNamesWith(ctx Context, f func(string)) {
	ctx.Value(keyOptions).(options.Options).IterateHyphaNamesWith(f)
}

func HyphaHTMLData(ctx Context, name string) (rawText, binaryHtml string, err error) {
	return ctx.Value(keyOptions).(options.Options).HyphaHTMLData(name)
}
