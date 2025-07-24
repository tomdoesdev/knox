package output

import "fmt"

// Heading represents a mixed element that has both title content and child elements
type Heading struct {
	title    fmt.Stringer
	children []Element
}

// NewHeading creates a new Heading element with the given title
func NewHeading(title string) *Heading {
	return &Heading{
		title:    StringValue(title),
		children: make([]Element, 0),
	}
}

// NewHeadingStringer creates a new Heading element with a fmt.Stringer title
func NewHeadingStringer(title fmt.Stringer) *Heading {
	return &Heading{
		title:    title,
		children: make([]Element, 0),
	}
}

// Add appends a child element to this heading
func (h *Heading) Add(element Element) {
	h.children = append(h.children, element)
}

// Value returns the title content of this heading
func (h *Heading) Value() fmt.Stringer {
	return h.title
}

// Children returns all child elements of this heading
func (h *Heading) Children() []Element {
	return h.children
}

// Type returns "heading" as the element type identifier
func (h *Heading) Type() string {
	return "heading"
}
