// Package genhtml provides utilities for converting Mycomarkup blocks into HTML documents. As of now, some parts of HTML generation are in other parts of the library, WIP.
package genhtml

import (
	"fmt"
	"html"
	"sort"
	"strings"

	"github.com/bouncepaw/mycomarkup/v3/blocks"
	"github.com/bouncepaw/mycomarkup/v3/mycocontext"
)

// This package shall not depend on anything other than blocks, links, globals, mycocontext, util.

// We might change this later or provide a way to change by the user.
const indentation = "\t"

// TagKind is the kind of a Tag. The way the tag is rendered depends on the kind.
type TagKind int

const (
	// ClosedTag is a tag that looks like that: <t>children</t>.
	ClosedTag TagKind = iota
	// UnclosedTag is a tag that looks like that: <t/>
	UnclosedTag
	// WrapperTag is a tag that looks like that: children
	WrapperTag
)

// Tag represents an HTML tag/DOM node.
type Tag struct {
	Name       string
	Kind       TagKind
	Attributes map[string]string
	Contents   string
	// nil if infertile
	Children []Tag
}

// String returns an indented pretty-printed representation of the Tag.
func (t Tag) String() (res string) {
	switch t.Kind {
	case UnclosedTag:
		return fmt.Sprintf("<%s%s/>\n", t.Name, attrs(t.Attributes))
	case ClosedTag:
		res += fmt.Sprintf("<%s%s>\n", t.Name, attrs(t.Attributes))
		var tmp string
		tmp += t.Contents
		for i, child := range t.Children {
			if i > 0 || (i == 0 && t.Contents != "") {
				tmp += "\n"
			}
			tmp += child.String()
		}
		res += eachLineIndented(tmp)
		res += fmt.Sprintf("</%s>", t.Name)
		return res
	case WrapperTag:
		res += t.Contents
		for i, child := range t.Children {
			if i > 0 || (i == 0 && t.Contents != "") {
				res += "\n"
			}
			res += child.String()
		}
		return res
	default:
		return "ERROR"
	}
}

func attrs(m map[string]string) (res string) {
	if len(m) == 0 {
		return ""
	}
	var parts []string
	for k, v := range m {
		// TODO: perform some escaping?
		parts = append(parts, fmt.Sprintf(` %s="%s"`, k, v))
	}
	// Sort so the output is the same for the same input.
	sort.Strings(parts)
	return strings.Join(parts, "")
}

// BlockToTag turns the given Block into a Tag depending on the Context and IDCounter.
func BlockToTag(ctx mycocontext.Context, block blocks.Block, counter *blocks.IDCounter) Tag {
	var attrs = map[string]string{}
	if counter.ShouldUseResults() {
		attrs["id"] = block.ID(counter)
	}
	switch block := block.(type) {
	case blocks.Formatted:
		var (
			contents string
			tagState = blocks.CleanStyleState()
		)
		for i, line := range block.Lines {
			if i > 0 {
				contents += `<br>`
			}
			for _, span := range line {
				switch s := span.(type) {
				case blocks.SpanTableEntry:
					contents += blocks.TagFromState(s.Kind(), tagState)
				case blocks.InlineLink:
					contents += fmt.Sprintf(
						`<a href="%s" class="%s">%s</a>`,
						s.Href(),
						s.Classes(),
						html.EscapeString(s.Display()),
					)
				case blocks.InlineText:
					contents += html.EscapeString(s.Contents)
				default:
					panic("unknown span")
				}
			}
			for stt, open := range tagState { // Close the unclosed
				if open {
					contents += blocks.TagFromState(stt, tagState)
				}
			}
		}
		return Tag{
			Name:       "",
			Kind:       WrapperTag,
			Attributes: attrs,
			Contents:   contents,
			Children:   nil,
		}
	case blocks.HorizontalLine:
		return Tag{
			Name:       "hr",
			Kind:       UnclosedTag,
			Attributes: attrs,
			Contents:   "",
			Children:   nil,
		}
	default:
		return Tag{
			Name:       "error",
			Kind:       UnclosedTag,
			Attributes: nil,
			Contents:   "",
			Children:   nil,
		}
	}
}

func eachLineIndented(s string) (res string) {
	for _, line := range strings.Split(s, "\n") {
		res += indentation + line
	}
	return res
}
