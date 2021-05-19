package parser

import (
	"bytes"
	"context"
	"errors"
)

type Key int

// These are keys for the context that floats around.
const (
	KeyHyphaName Key = iota
	KeyInputBuffer
)

// ParsingDone is returned by Context when the parsing is done because there is no more input.
var ParsingDone = errors.New("parsing done")

// ContextFromStringInput returns the context for the given input.
func ContextFromStringInput(hyphaName, input string) (context.Context, context.CancelFunc) {
	ctx, cancel := context.WithCancel(
		context.WithValue(
			context.WithValue(
				context.Background(),
				KeyHyphaName,
				hyphaName),
			KeyInputBuffer,
			bytes.NewBufferString(input),
		),
	)
	return ctx, cancel
}

func hyphaName(ctx context.Context) string {
	return ctx.Value(KeyHyphaName).(string)
}

func input(ctx context.Context) *bytes.Buffer {
	return ctx.Value(KeyInputBuffer).(*bytes.Buffer)
}
