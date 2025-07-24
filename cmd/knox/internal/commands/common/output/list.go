package output

import "fmt"

// List represents a container element that holds list items
type List struct {
	children []Element
}

// NewList creates a new empty List element
func NewList() *List {
	return &List{children: make([]Element, 0)}
}

// Add appends a child element to this list
func (l *List) Add(element Element) {
	l.children = append(l.children, element)
}

// Value returns empty string as List is a container element with no direct content
func (l *List) Value() fmt.Stringer {
	return StringValue("")
}

// Children returns all child elements of this list
func (l *List) Children() []Element {
	return l.children
}

// Type returns "list" as the element type identifier
func (l *List) Type() string {
	return "list"
}

// ListItem represents a single item within a list
type ListItem struct {
	value fmt.Stringer
}

// NewListItem creates a new ListItem with the given content
func NewListItem(content string) *ListItem {
	return &ListItem{value: StringValue(content)}
}

// NewListItemStringer creates a new ListItem with a fmt.Stringer value
func NewListItemStringer(value fmt.Stringer) *ListItem {
	return &ListItem{value: value}
}

// Value returns the content of this list item
func (li *ListItem) Value() fmt.Stringer {
	return li.value
}

// Children returns an empty slice as ListItem elements have no children
func (li *ListItem) Children() []Element {
	return []Element{}
}

// Type returns "listitem" as the element type identifier
func (li *ListItem) Type() string {
	return "listitem"
}
