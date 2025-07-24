package ast

// Common node type constructors for typical use cases

// NewRoot creates a root container node
func NewRoot() *BasicNode {
	return NewContainerNode("root")
}

// NewHeading creates a heading node with title content
func NewHeading(title string) *BasicNode {
	return NewNode("heading", StringValue(title))
}

// NewList creates a list container node
func NewList() *BasicNode {
	return NewContainerNode("list")
}

// NewListItem creates a list item leaf node
func NewListItem(content string) *BasicNode {
	return NewNode("listitem", StringValue(content))
}

// NewSection creates a section container node
func NewSection(title string) *BasicNode {
	return NewNode("section", StringValue(title))
}

// Text creates a simple text node (alias for NewTextNode for brevity)
func Text(content string) *BasicNode {
	return NewTextNode(content)
}

// Container creates a generic container node (alias for NewContainerNode)
func Container(nodeType string) *BasicNode {
	return NewContainerNode(nodeType)
}
