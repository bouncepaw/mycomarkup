// Package tag provides the data type for (X)HTML tags/DOM nodes.
package tag

import (
	"fmt"
	"github.com/bouncepaw/mycomarkup/v3/util/lines"
	"sort"
	"strings"
)

// tagKind is the kind of a Tag. The way the tag is rendered depends on the kind.
type tagKind int

const (
	// closed is a tag that looks like that: <t>children</t>.
	closed tagKind = iota
	// unclosed is a tag that looks like that: <t/>
	unclosed
	// wrapper is a tag that looks like that: children
	wrapper
)

// Tag represents an HTML tag/DOM node.
type Tag struct {
	name       string
	kind       tagKind
	attributes map[string]string
	contents   []lines.Line
	children   []Tag
}

// NewUnclosed returns a new unclosed tag.
func NewUnclosed(name string) Tag {
	return Tag{
		name:       name,
		kind:       unclosed,
		attributes: nil,
		contents:   nil,
		children:   nil,
	}
}

// NewClosed returns a new closed tag.
func NewClosed(name string) Tag {
	return Tag{
		name:       name,
		kind:       closed,
		attributes: nil,
		contents:   nil,
		children:   nil,
	}
}

// NewWrapper returns a new wrapper tag.
func NewWrapper() Tag {
	return Tag{
		name:       "",
		kind:       wrapper,
		attributes: nil,
		contents:   nil,
		children:   nil,
	}
}

// WithChildren returns the tag but with the given children. Previous children of the tag are discarded.
//
// This is a no-op for unclosed tags.
func (t Tag) WithChildren(children ...Tag) Tag {
	if t.kind == unclosed {
		return t
	}
	t.children = children
	return t
}

// WithAttrs return the tag but with the given attributes. Previous attributes of the tag are discarded.
//
// This is a no-op for wrapper tags.
func (t Tag) WithAttrs(attributes map[string]string) Tag {
	if t.kind == wrapper {
		return t
	}
	t.attributes = attributes
	return t
}

// WithContentsLines returns the tag but with the given lines of contents.
//
// Contents is like children, but just text, not tags.
//
// This is a no-op for unclosed tags.
func (t Tag) WithContentsLines(lines ...lines.Line) Tag {
	if t.kind == unclosed {
		return t
	}
	t.contents = lines
	return t
}

// WithContentsStrings is like WithContentsLines but it wraps strings into indented lines for you.
func (t Tag) WithContentsStrings(strs ...string) Tag {
	var contentsLines []lines.Line
	for _, str := range strs {
		contentsLines = append(contentsLines, lines.IndentableFrom(str))
	}
	return t.WithContentsLines(contentsLines...)
}

// String returns a string representation of the tag.
func (t Tag) String() string {
	var res string
	for _, line := range t.Lines() {
		res += line.String()
	}
	return res
}

// Lines returns rendered lines of the tag.
func (t Tag) Lines() (res []lines.Line) {
	switch t.kind {
	case unclosed:
		return []lines.Line{
			lines.IndentableFrom(fmt.Sprintf("<%s%s/>", t.name, attrs(t.attributes))),
		}

	case wrapper:
		res = t.contents
		for _, child := range t.children {
			res = append(res, child.Lines()...)
		}
		return res

	case closed:
		res = []lines.Line{
			lines.IndentableFrom(fmt.Sprintf("<%s%s>", t.name, attrs(t.attributes))),
		}
		res = append(res, t.contents...)
		for _, child := range t.children {
			for _, line := range child.Lines() {
				res = append(res, line.Indented())
			}
		}
		res = append(res, lines.IndentableFrom(fmt.Sprintf("</%s>\n", t.name)))
		return res

	default:
		res = append(res, lines.UnindentableFrom("ERROR"))
	}
	return res
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
