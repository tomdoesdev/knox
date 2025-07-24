package ast

import "fmt"

// BasicNode provides a simple implementation of MutableNode
type BasicNode struct {
	content    fmt.Stringer
	children   []Node
	attributes map[string]interface{}
	nodeType   string
}

// NewTextNode creates a new leaf node with text content
func NewTextNode(content string) *BasicNode {
	return &BasicNode{
		content:    StringValue(content),
		children:   make([]Node, 0),
		attributes: make(map[string]interface{}),
		nodeType:   "text",
	}
}

// NewContainerNode creates a new container node with no content
func NewContainerNode(nodeType string) *BasicNode {
	return &BasicNode{
		content:    EmptyValue{},
		children:   make([]Node, 0),
		attributes: make(map[string]interface{}),
		nodeType:   nodeType,
	}
}

// NewNode creates a new node with specified type and content
func NewNode(nodeType string, content fmt.Stringer) *BasicNode {
	return &BasicNode{
		content:    content,
		children:   make([]Node, 0),
		attributes: make(map[string]interface{}),
		nodeType:   nodeType,
	}
}

// Content returns the primary content of this node
func (n *BasicNode) Content() fmt.Stringer {
	return n.content
}

// Children returns all child nodes
func (n *BasicNode) Children() []Node {
	return n.children
}

// Attributes returns metadata associated with this node
func (n *BasicNode) Attributes() map[string]interface{} {
	return n.attributes
}

// Type returns the node type identifier
func (n *BasicNode) Type() string {
	return n.nodeType
}

// SetContent updates the node's content
func (n *BasicNode) SetContent(content fmt.Stringer) {
	n.content = content
}

// AddChild appends a child node
func (n *BasicNode) AddChild(child Node) {
	n.children = append(n.children, child)
}

// SetAttribute sets a metadata attribute
func (n *BasicNode) SetAttribute(key string, value interface{}) {
	n.attributes[key] = value
}

// GetAttribute retrieves a metadata attribute
func (n *BasicNode) GetAttribute(key string) (interface{}, bool) {
	value, exists := n.attributes[key]
	return value, exists
}
