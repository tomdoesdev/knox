package output

import "fmt"

// Root represents the root container element of the output tree
type Root struct {
	children []Element
}

// NewRoot creates a new empty Root element
func NewRoot() *Root {
	return &Root{children: make([]Element, 0)}
}

// Add appends a child element to this root container
func (r *Root) Add(element Element) {
	r.children = append(r.children, element)
}

// Value returns empty string as Root is a container element with no direct content
func (r *Root) Value() fmt.Stringer {
	return StringValue("")
}

// Children returns all child elements of this root container
func (r *Root) Children() []Element {
	return r.children
}

// Type returns "root" as the element type identifier
func (r *Root) Type() string {
	return "root"
}
