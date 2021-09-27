package blocks

import (
	"fmt"
	"html"
)

// CodeBlock represents a block of preformatted text.
type CodeBlock struct {
	language string
	contents string
}

// ID returns the code block's id which is "codeblock-" and a number.
func (cb CodeBlock) ID(counter *IDCounter) string {
	counter.codeblocks++
	return fmt.Sprintf("codeblock-%d", counter.codeblocks)
}

// NewCodeBlock returns a code block with the given language and contents.
func NewCodeBlock(language, contents string) CodeBlock {
	return CodeBlock{
		language: language,
		contents: contents,
	}
}

// Language returns what kind of formal language the code block is written in. It returns "plain" if the language is not specified. Returns escaped text otherwise.
func (cb CodeBlock) Language() string {
	if cb.language == "" {
		return "plain"
	}
	return html.EscapeString(cb.language)
}

// Contents returns the code block's contents.
func (cb CodeBlock) Contents() string {
	return cb.contents
}
