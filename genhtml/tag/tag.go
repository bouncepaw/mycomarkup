// Package tag provides the data type for (X)HTML tags/DOM nodes.
package tag

import (
	"fmt"
	"sort"
	"strings"
)

// Kind is the kind of a Tag. The way the tag is rendered depends on the kind.
type Kind int

const (
	// Closed is a tag that looks like that: <t>children</t>.
	Closed Kind = iota
	// Unclosed is a tag that looks like that: <t/>
	Unclosed
	// Wrapper is a tag that looks like that: children
	Wrapper
)

// Tag represents an HTML tag/DOM node.
type Tag struct {
	Name       string
	Kind       Kind
	Attributes map[string]string
	Contents   string
	// nil if infertile
	Children []Tag
}

// New returns a new Tag with the given data.
func New(name string, kind Kind, attributes map[string]string, contents string, children []Tag) Tag {
	return Tag{
		Name:       name,
		Kind:       kind,
		Attributes: attributes,
		Contents:   contents,
		Children:   children,
	}
}

// String returns an indented pretty-printed representation of the Tag.
func (t Tag) String() (res string) {
	switch t.Kind {
	case Unclosed:
		return fmt.Sprintf("<%s%s/>\n", t.Name, attrs(t.Attributes))
	case Closed:
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
	case Wrapper:
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

func eachLineIndented(s string) (res string) {
	for _, line := range strings.Split(s, "\n") {
		res += "\t" + line
	}
	return res
}
