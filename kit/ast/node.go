package ast

import "fmt"

// Node represents a node in an abstract syntax tree
// This is a pure dat
type Node interface {
	// Content returns the primary content/value of this node
	// Returns empty stringer for container-only nodes
	Content() fmt.Stringer

	// Children returns all child nodes
	// Returns empty slice for leaf nodes
	Children() []Node

	// Attributes returns metadata associated with this node
	// Can contain any structured data as attributes
	Attributes() map[string]*Attribute

	// GetAttribute retrieves a specific metadata attribute
	GetAttribute(key string) (*Attribute, bool)

	// Type returns a string identifier for the node type
	Type() string
}

// MutableNode interface for nodes that can be modified after creation
type MutableNode interface {
	Node

	// SetContent updates the node's content
	SetContent(content fmt.Stringer)

	// AddChild appends a child node
	AddChild(child Node)

	// SetAttribute sets a metadata attribute
	SetAttribute(key string, value any)
}

// StringValue provides a simple implementation of fmt.Stringer for string content
type StringValue string

func (s StringValue) String() string {
	return string(s)
}

// EmptyValue represents empty content for container nodes
type EmptyValue struct{}

func (e EmptyValue) String() string {
	return ""
}
