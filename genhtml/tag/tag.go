// Package tag provides the data type for (X)HTML tags/DOM nodes.
package tag

import (
	"fmt"
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
	contents   string
	children   []Tag
}

// NewUnclosed returns a new unclosed tag.
func NewUnclosed(name string, attributes map[string]string) Tag {
	return Tag{
		name:       name,
		kind:       unclosed,
		attributes: attributes,
		contents:   "",
		children:   nil,
	}
}

// NewClosed returns a new closed tag.
func NewClosed(name string, attributes map[string]string, contents string, children ...Tag) Tag {
	return Tag{
		name:       name,
		kind:       closed,
		attributes: attributes,
		contents:   contents,
		children:   children,
	}
}

// NewWrapper returns a new wrapper tag.
func NewWrapper(contents string, children ...Tag) Tag {
	return Tag{
		name:       "",
		kind:       wrapper,
		attributes: map[string]string{},
		contents:   contents,
		children:   children,
	}
}

// String returns an indented pretty-printed representation of the Tag.
func (t Tag) String() (res string) {
	switch t.kind {
	case unclosed:
		return fmt.Sprintf("<%s%s/>\n", t.name, attrs(t.attributes))
	case closed:
		res += fmt.Sprintf("<%s%s>\n", t.name, attrs(t.attributes))
		var tmp string
		tmp += t.contents
		for i, child := range t.children {
			if i > 0 || (i == 0 && t.contents != "") {
				tmp += "\n"
			}
			tmp += child.String()
		}
		res += eachLineIndented(tmp)
		res += fmt.Sprintf("</%s>\n", t.name)
		return res
	case wrapper:
		res += t.contents
		for i, child := range t.children {
			if i > 0 || (i == 0 && t.contents != "") {
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

func eachLineIndented(s string) (res string) {
	for i, line := range strings.Split(s, "\n") {
		if i > 0 {
			res += "\n"
		}
		res += "\t" + line
	}
	return res
}
