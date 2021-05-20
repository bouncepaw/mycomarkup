package parser

import "context"

/*
## List examples
### Basic unordered flat list
* Item 1
* Item 2
* Item 3

### Basic ordered flat list
*. Item 1
*. Item 2
*. Item 3

### Basic to do list
*v Done
*x Not done
*v Done

### Multiline contents
* {
  //Any mycomarkup supported//

  * Other lists too
}
* This is single-line again

### Nesting lists
* As a shortcut for multiline contents with other lists only,
** you can use this syntax.
* You just have to
** increase
*** the amount of asterisks

### Mixing different types of lists
* This is from one list
*. But this is from a different one.
*v You can't mix them on one level

*v You can only mix
*x to do items

*. You can nest
** Items of different type
** But they have to be the same on one level

### Parsing approach
We read all list items of the same type and their contents (single-line or multi-line). The contents are parsed as if they were separate Mycomarkup documents.
*/

type List struct {
	Marker   ListMarker
	Level    uint
	Contents []interface{}
}

// ListMarker is the type of marker the ListEntry has.
type ListMarker int

const (
	MarkerUnordered ListMarker = iota
	MarkerOrdered
	MarkerTodoDone
	MarkerTodo
)

func eatUntilSpace(ctx context.Context) {
	// We do not care what is read, therefore we drop the read line.
	// We know that there //is// a space beforehand, therefore we drop the error.
	_, _ = inputFrom(ctx).ReadString(' ')
}

func looksLikeList(ctx context.Context) bool {
	_, level, found := markerOnNextLine(ctx)
	return found && level == 1
}

func markerOnNextLine(ctx context.Context) (m ListMarker, level uint, found bool) {
	var (
		onStart            = true
		onAsterisk         = false
		onSpecialCharacter = false
	)
	for _, b := range inputFrom(ctx).Bytes() {
		switch {
		case onStart && b != '*':
			return MarkerUnordered, 0, false
		case onStart:
			level = 1
			onStart = false
			onAsterisk = true

		case onAsterisk && b == '*':
			level++
		case onAsterisk && b == ' ':
			return m, level, true
		case onAsterisk && (b == 'v' || b == 'x' || b == '.'):
			onAsterisk = false
			onSpecialCharacter = true
			switch b {
			case 'v':
				m = MarkerTodoDone
			case 'x':
				m = MarkerTodo
			case '.':
				m = MarkerOrdered
			}
		case onAsterisk:
			return MarkerUnordered, 0, false

		case onSpecialCharacter && b != ' ':
			return MarkerUnordered, 0, false
		case onSpecialCharacter:
			return m, level, true
		}
	}
	panic("unreachable")
}

func (m1 ListMarker) sameAs(m2 ListMarker) bool {
	return (m1 == m2) ||
		((m1 == MarkerTodoDone) && (m2 == MarkerTodo)) ||
		((m1 == MarkerTodo) && (m2 == MarkerTodoDone))
}

func (m ListMarker) HTMLTemplate() string {
	switch m {
	case MarkerUnordered:
		return `<li class="item_unordered">%s</li>`
	case MarkerOrdered:
		return `<li class="item_ordered">%s</li>`
	case MarkerTodoDone:
		return `<li class="item_todo item_todo-done"><input type="checkbox" disabled checked>%s</li>`
	case MarkerTodo:
		return `<li class="item_todo"><input type="checkbox" disabled>%s</li>`
	}
	panic("unreachable")
}
