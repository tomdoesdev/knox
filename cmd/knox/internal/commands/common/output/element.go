package output

import "fmt"

// Element represents a node in the DOM-like output tree
type Element interface {
	// Value returns the raw content of the element
	// Returns empty string for container elements that have no direct content
	Value() fmt.Stringer

	// Children returns all child elements
	// Returns empty slice for leaf elements
	Children() []Element

	// Type returns the element type identifier (e.g., "text", "list", "heading")
	Type() string
}
