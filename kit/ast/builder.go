package ast

import (
	"fmt"
)

// Builder provides a fluent API for constructing AST trees
type Builder interface {
	// Node creation and navigation
	Node(nodeType string) Builder // Add child node and navigate to it
	Up() Builder                  // Navigate back to parent node
	Root() Builder                // Navigate back to root node

	// Content and attributes
	Content(content string) Builder             // Set content for current node
	Attr(key string, value interface{}) Builder // Set attribute on current node

	// Tree completion
	Build() Node // Return the built tree
}

// mutableNode is a node implementation used during tree construction
type mutableNode struct {
	nodeType   string
	content    fmt.Stringer
	children   []Node
	attributes map[string]*Attribute
	parent     *mutableNode // For navigation during building
}

func (n *mutableNode) Type() string {
	return n.nodeType
}

func (n *mutableNode) Content() fmt.Stringer {
	return n.content
}

func (n *mutableNode) Children() []Node {
	return n.children
}

func (n *mutableNode) Attributes() map[string]*Attribute {
	return n.attributes
}

func (n *mutableNode) GetAttribute(key string) (*Attribute, bool) {
	value, exists := n.attributes[key]
	return value, exists
}

// nodeBuilder implements the Builder interface
type nodeBuilder struct {
	root    *mutableNode   // Root of the tree being built
	current *mutableNode   // Current node being modified
	stack   []*mutableNode // Parent stack for Up() navigation
}

// NewBuilder creates a new builder with the specified root node type
func NewBuilder(rootType string) Builder {
	root := &mutableNode{
		nodeType:   rootType,
		content:    EmptyValue{},
		children:   make([]Node, 0),
		attributes: make(map[string]*Attribute),
		parent:     nil,
	}

	return &nodeBuilder{
		root:    root,
		current: root,
		stack:   make([]*mutableNode, 0),
	}
}

// Node adds a child node of the specified type and navigates to it
func (b *nodeBuilder) Node(nodeType string) Builder {
	child := &mutableNode{
		nodeType:   nodeType,
		content:    EmptyValue{},
		children:   make([]Node, 0),
		attributes: make(map[string]*Attribute),
		parent:     b.current,
	}

	// Add child to current node
	b.current.children = append(b.current.children, child)

	// Push current to stack and navigate to child
	b.stack = append(b.stack, b.current)
	b.current = child

	return b
}

// Up navigates back to the parent node
func (b *nodeBuilder) Up() Builder {
	if len(b.stack) == 0 {
		return b // Already at root, no-op
	}

	// Pop parent from stack
	b.current = b.stack[len(b.stack)-1]
	b.stack = b.stack[:len(b.stack)-1]

	return b
}

// Root navigates back to the root node
func (b *nodeBuilder) Root() Builder {
	b.current = b.root
	b.stack = b.stack[:0] // Clear stack
	return b
}

// Content sets the content of the current node
func (b *nodeBuilder) Content(content string) Builder {
	b.current.content = StringValue(content)
	return b
}

// Attr sets an attribute on the current node
func (b *nodeBuilder) Attr(key string, value any) Builder {
	b.current.attributes[key] = NewAttribute(value)
	return b
}

// Build returns the constructed tree as a Node
func (b *nodeBuilder) Build() Node {
	return b.root
}
