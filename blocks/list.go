package blocks

import "fmt"

// List is the block representing a set of related elements. It must be the same as all ListItem.Marker.
type List struct {
	// Items are the entries of the List. There should be at least one.
	Items []ListItem
	// Marker is the type of the list. All entries have the same type. See SameAs for information about same types.
	Marker ListMarker
}

// ID returns this list's id. It is list- and the list's number.
func (l List) ID(counter *IDCounter) string {
	counter.lists++
	return fmt.Sprintf("list-%d", counter.lists)
}

func (l List) isBlock() {}

// ListItem is an entry in a List.
type ListItem struct {
	Marker ListMarker
	// Level is equal to amount of asterisks.
	//     *    -> Level = 1
	//     **.  -> Level = 2
	Level uint
	// Contents are Mycomarkup blocks contained in this list item.
	Contents []Block
}

func (l ListItem) isBlock() {}

// ID returns an empty string because list items don't have ids.
func (l ListItem) ID(_ *IDCounter) string {
	return ""
}

// ListMarker is the type of a ListItem or a List.
type ListMarker int

const (
	// MarkerUnordered is for bullets like * (no point).
	MarkerUnordered ListMarker = iota
	// MarkerOrdered is for bullets like *. (with point).
	MarkerOrdered
	// MarkerTodoDone is for bullets like *v (with tick).
	MarkerTodoDone
	// MarkerTodo is for bullets like *x (with cross).
	MarkerTodo
)

// SameAs is true if both list markers are of the same type. MarkerTodoDone and MarkerTodo are considered same to each other. All other markers are different from each other.
func (m1 ListMarker) SameAs(m2 ListMarker) bool {
	return (m1 == m2) ||
		((m1 == MarkerTodoDone) && (m2 == MarkerTodo)) ||
		((m1 == MarkerTodo) && (m2 == MarkerTodoDone))
}
