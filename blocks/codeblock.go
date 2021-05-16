package blocks

import "strings"

// CodeBlock represents a block of preformatted text.
type CodeBlock struct {
	language string
	contents string
}

// MakeCodeBlock returns a code block with the given language and contents.
func MakeCodeBlock(language, contents string) CodeBlock {
	return CodeBlock{
		language: language,
		contents: contents,
	}
}

// Language returns what kind of formal language the code block is written in. It returns "plain" if the language is not specified.
func (cb *CodeBlock) Language() string {
	// TODO: some form of protection should be done?
	if cb.language == "" {
		return "plain"
	}
	return cb.language
}

// Contents returns the code block's contents.
func (cb *CodeBlock) Contents() string {
	return strings.TrimPrefix(cb.contents, "\n")
}

// AddLine adds a line to the code block's contents. The line should be without line breaks.
func (cb *CodeBlock) AddLine(line string) {
	cb.contents += "\n" + line
}
