package output

import "fmt"

// StringValue is a simple implementation of fmt.Stringer for string values
type StringValue string

func (s StringValue) String() string {
	return string(s)
}

// Text represents a leaf element containing text content
type Text struct {
	value fmt.Stringer
}

// NewText creates a new Text element with the given content
func NewText(content string) *Text {
	return &Text{value: StringValue(content)}
}

// NewTextStringer creates a new Text element with a fmt.Stringer value
func NewTextStringer(value fmt.Stringer) *Text {
	return &Text{value: value}
}

// Value returns the text content of this element
func (t *Text) Value() fmt.Stringer {
	return t.value
}

// Children returns an empty slice as Text elements have no children
func (t *Text) Children() []Element {
	return []Element{}
}

// Type returns "text" as the element type identifier
func (t *Text) Type() string {
	return "text"
}
