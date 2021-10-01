package parser

import (
	"bytes"
	"github.com/bouncepaw/mycomarkup/v3/blocks"
	"github.com/bouncepaw/mycomarkup/v3/mycocontext"
)

// Call only if there is a list item on the line.
func nextList(ctx mycocontext.Context) (list blocks.List, eof bool) {
	var contents []blocks.Block
	rootMarker, rootLevel, _ := markerOnNextLine(ctx)
	list = blocks.List{
		Items:  make([]blocks.ListItem, 0),
		Marker: rootMarker,
	}
	for !eof {
		marker, level, found := markerOnNextLine(ctx)
		if !found || (!marker.SameAs(list.Marker) && rootLevel == level) {
			break
		}

		_ = mycocontext.EatUntilSpace(ctx)
		contents, eof = nextListItem(ctx, rootLevel)
		item := blocks.ListItem{
			Marker:   marker,
			Level:    level,
			Contents: contents,
		}

		list.Items = append(list.Items, item)
	}
	return list, eof
}

func readNextListItemsContents(ctx mycocontext.Context) (text bytes.Buffer, eof bool) {
	var (
		onNewLine       = true
		escaping        = false
		curlyBracesOpen = 0
		b               byte
	)
walker: // Read all item's contents
	for !eof {
		b, eof = mycocontext.NextByte(ctx)
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

func nextListItem(
	ctx mycocontext.Context,
	rootLevel uint, // They have to have a level higher than this, though
) (contents []blocks.Block, eof bool) {
	// Parse the text as a separate mycodoc
	var (
		text    bytes.Buffer
		ast     = make([]blocks.Block, 0)
		subText bytes.Buffer
	)
	text, eof = readNextListItemsContents(ctx)

	// Grab the sublist text, if there is one. Each bullet is decremented by one asterisk.
	for !eof {
		_, level, found := markerOnNextLine(ctx)
		// We are not interested in same level or less-nested list items. Screw them! Forget them!
		if !found || level <= rootLevel {
			break
		}

		// I am so sure there is an asterisk we can simply drop.
		// Add a newline for proper parsing later on.
		// The space is left by EatUntilSpace at the end of the string.
		disnestedBullet := "\n" + mycocontext.EatUntilSpace(ctx)[1:]
		text.WriteString(disnestedBullet)

		subText, eof = readNextListItemsContents(ctx)
		_, _ = subText.WriteTo(&text) // Let's just hope it never fails. We are confident people.
	}

	parseSubdocumentForEachBlock(ctx, &text, func(block blocks.Block) {
		ast = append(ast, block)
	})

	return ast, eof
}

func looksLikeList(ctx mycocontext.Context) bool {
	_, level, found := markerOnNextLine(ctx)
	return found && level == 1
}

func markerOnNextLine(ctx mycocontext.Context) (m blocks.ListMarker, level uint, found bool) {
	var (
		onStart            = true
		onAsterisk         = false
		onSpecialCharacter = false
	)
	for _, b := range ctx.Input().Bytes() {
		switch {
		case onStart && b != '*':
			return blocks.MarkerUnordered, 0, false
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
				m = blocks.MarkerTodoDone
			case 'x':
				m = blocks.MarkerTodo
			case '.':
				m = blocks.MarkerOrdered
			}
		case onAsterisk:
			return blocks.MarkerUnordered, 0, false

		case onSpecialCharacter && b != ' ':
			return blocks.MarkerUnordered, 0, false
		case onSpecialCharacter:
			return m, level, true
		}
	}
	return blocks.MarkerUnordered, 0, false
}
