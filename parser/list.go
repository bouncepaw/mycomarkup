package parser

import (
	"bytes"
	"context"
	"sync"
)

type List struct {
	Items  []ListItem
	Marker ListMarker
}

// Call only if there is a list item on the line.
func nextList(ctx context.Context) (list List, eof bool) {
	var contents []interface{}
	list = List{
		Items: make([]ListItem, 0),
	}
	for !eof {
		marker, level, found := markerOnNextLine(ctx)
		if !found {
			break
		}

		eatUntilSpace(ctx)
		contents, eof = nextListItem(ctx)
		item := ListItem{
			Marker:   marker,
			Level:    level,
			Contents: contents,
		}

		list.Items = append(list.Items, item)
	}
	list.Marker = list.Items[0].Marker // There should be at least one item!
	return list, eof
}

func readNextListItemsContents(ctx context.Context) (text bytes.Buffer, eof bool) {
	var (
		onNewLine       = true
		escaping        = false
		curlyBracesOpen = 0
		b               byte
	)
walker: // Read all item's contents
	for !eof {
		b, eof = nextByte(ctx)
	stateMachine: // I'm extremely sorry
		switch {
		case onNewLine && b != ' ':
			onNewLine = false
			goto stateMachine
		case onNewLine: // We just ignore spaces on line beginnings
		case escaping:
			escaping = false
			if b == '\n' && curlyBracesOpen == 0 {
				break walker
			}
			text.WriteByte(b)
		case b == '\\':
			escaping = true
			text.WriteByte('\\')

		case b == '{':
			if curlyBracesOpen > 0 {
				text.WriteByte('{')
			}
			curlyBracesOpen++
		case b == '}':
			if curlyBracesOpen != 1 {
				text.WriteByte('}')
			}
			if curlyBracesOpen >= 0 {
				curlyBracesOpen--
			}
		case b == '\n' && curlyBracesOpen == 0:
			break walker
		case b == '\n':
			text.WriteByte(b)
			onNewLine = true
		default:
			text.WriteByte(b)
		}
	}
	return text, eof
}

func nextListItem(ctx context.Context) (contents []interface{}, eof bool) {
	// Parse the text as a separate mycodoc
	var (
		text     bytes.Buffer
		blocksCh = make(chan interface{})
		blocks   = make([]interface{}, 0)
		wg       sync.WaitGroup
	)
	text, eof = readNextListItemsContents(ctx)

	wg.Add(1)
	go func() {
		Parse(context.WithValue(ctx, KeyInputBuffer, &text), blocksCh)
		wg.Done()
	}()
	for block := range blocksCh {
		blocks = append(blocks, block)
	}
	wg.Wait()

	return blocks, eof
}

type ListItem struct {
	Marker   ListMarker
	Level    uint
	Contents []interface{}
}

func eatUntilSpace(ctx context.Context) {
	// We do not care what is read, therefore we drop the read line.
	// We know that there //is// a space beforehand, therefore we drop the error.
	_, _ = inputFrom(ctx).ReadString(' ')
}

// ListMarker is the type of marker the ListItem has.
type ListMarker int

const (
	MarkerUnordered ListMarker = iota
	MarkerOrdered
	MarkerTodoDone
	MarkerTodo
)

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
	return MarkerUnordered, 0, false
}

func (m1 ListMarker) sameAs(m2 ListMarker) bool {
	return (m1 == m2) ||
		((m1 == MarkerTodoDone) && (m2 == MarkerTodo)) ||
		((m1 == MarkerTodo) && (m2 == MarkerTodoDone))
}
