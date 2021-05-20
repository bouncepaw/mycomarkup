package blocks

import (
	"errors"
	"regexp"
	"strings"
)

// List is a nesting list, either ordered or unordered.
type List struct {
	curr      *ListItem
	hyphaName string
	ordered   bool
	finalized bool
}

func MakeList(line, hyphaName string) (*List, bool) {
	list := &List{
		hyphaName: hyphaName,
		curr:      newListItem(nil),
	}
	return list, list.ProcessLine(line)
}

func (list *List) pushItem() {
	item := newListItem(list.curr)
	list.curr.children = append(list.curr.children, item)
	list.curr = item
}

func (list *List) popItem() {
	if list.curr == nil {
		return
	}
	list.curr = list.curr.parent
}

func (list *List) balance(level int) {
	for level > list.curr.depth {
		list.pushItem()
	}

	for level < list.curr.depth {
		list.popItem()
	}
}

var listRe = regexp.MustCompile(`^\*+\.? `)

func okForList(line string) bool {
	return listRe.MatchString(line)
}

func (list *List) ProcessLine(line string) (done bool) {
	if !okForList(line) {
		list.Finalize()
		return true
	}
	level, offset, ordered, err := parseListItem(line)
	if err != nil {
		list.Finalize()
		return true
	}

	// update ordered flag if the current node is the root one
	// (i.e. no parsing has been done yet)
	if list.curr.parent == nil {
		list.ordered = ordered
	}

	// if list type has suddenly changed (ill-formatted list), quit
	if ordered != list.ordered {
		list.Finalize()
		return true
	}

	list.balance(level)

	// if the current node already has content, create a new one
	// to prevent overwriting existing content (effectively creating
	// a new sibling node)
	if len(list.curr.content) > 0 {
		list.popItem()
		list.pushItem()
	}

	list.curr.content = line[offset:]

	return false
}

func (list *List) Finalize() {
	if !list.finalized {
		// close all opened nodes, effectively going up to the root node
		list.balance(0)
		list.finalized = true
	}
}

func (list *List) RenderAsHtml() (html string) {
	// for a good measure
	list.Finalize()

	b := &strings.Builder{}

	// fire up recursive render process
	list.curr.renderAsHtmlTo(b, list.hyphaName, list.ordered)

	return "\n" + b.String()
}

func parseListItem(line string) (level int, offset int, ordered bool, err error) {
	for line[level] == '*' {
		level++
	}

	if line[level] == '.' {
		ordered = true
		offset = level + 2
	} else {
		ordered = false
		offset = level + 1
	}

	if line[offset-1] != ' ' || len(line) < offset+2 || level < 1 || level > 6 {
		err = errors.New("ill-formatted list item")
	}
	return
}

func MatchesList(line string) bool {
	level, _, _, err := parseListItem(line)
	return err == nil && level == 1
}

// ListItem is an entry in a List.
type ListItem struct {
	content  string
	parent   *ListItem
	children []*ListItem
	depth    int
}

func newListItem(parent *ListItem) *ListItem {
	depth := 0
	if parent != nil {
		depth = parent.depth + 1
	}
	return &ListItem{
		parent:   parent,
		children: make([]*ListItem, 0),
		depth:    depth,
	}
}

func (item *ListItem) renderAsHtmlTo(b *strings.Builder, hyphaName string, ordered bool) {
	if len(item.content) > 0 {
		b.WriteString("<li>")
		b.WriteString(ParagraphToHtml(hyphaName, item.content))
	}

	if len(item.children) > 0 {
		if ordered {
			b.WriteString("<ol>")
		} else {
			b.WriteString("<ul>")
		}

		for _, child := range item.children {
			child.renderAsHtmlTo(b, hyphaName, ordered)
		}

		if ordered {
			b.WriteString("</ol>")
		} else {
			b.WriteString("</ul>")
		}
	}

	if len(item.content) > 0 {
		b.WriteString("</li>")
	}
}
