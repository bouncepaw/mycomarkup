// Package mycocontext provides a wrapper over context.Context and some operations on the wrapper.
package mycocontext

import (
	"bytes"
	"context"
	"git.sr.ht/~bouncepaw/mycomarkup/v5/options"
)

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

// WithBuffer returns a copy of the given context but with a different input buffer.
func WithBuffer(ctx Context, buf *bytes.Buffer) Context {
	return Context{&mycoContext{context.WithValue(ctx, keyInputBuffer, buf)}}
}

func WithOptions(ctx Context, opts options.Options) Context {
	return Context{&mycoContext{context.WithValue(ctx, keyOptions, opts)}}
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
