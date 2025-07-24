package ast

import (
	"fmt"
	"strings"
)

// Renderer interface converts AST nodes to string output
type Renderer interface {
	Render(root Node) (string, error)
}

// RendererFunc allows function-based renderers
type RendererFunc func(Node) (string, error)

func (rf RendererFunc) Render(root Node) (string, error) {
	return rf(root)
}

// HandlerFunc defines the signature for node handlers
type HandlerFunc func(Node, *strings.Builder) error

// Traverser handles tree traversal with pluggable node handlers
type Traverser struct {
	handlers map[string]HandlerFunc
}

// NewTraverser creates a new traverser with empty handler map
func NewTraverser() *Traverser {
	return &Traverser{
		handlers: make(map[string]HandlerFunc),
	}
}

// Handle registers a handler function for a specific node type
// Returns the traverser to enable method chaining
func (t *Traverser) Handle(nodeType string, handler HandlerFunc) *Traverser {
	t.handlers[nodeType] = handler
	return t
}

// Render traverses the AST tree and applies handlers to generate output
// Uses stack-based depth-first traversal for proper ordering
func (t *Traverser) Render(root Node) (string, error) {
	sb := &strings.Builder{}
	stack := []Node{root}

	for len(stack) > 0 {
		// Pop current node from stack
		current := stack[len(stack)-1]
		stack = stack[:len(stack)-1]

		// Find and apply handler for current node type
		handler, exists := t.handlers[current.Type()]
		if !exists {
			return "", fmt.Errorf("no handler registered for node type: %s", current.Type())
		}

		// Apply the handler to the current node
		if err := handler(current, sb); err != nil {
			return "", fmt.Errorf("handler error for node type %s: %w", current.Type(), err)
		}

		// Add children to stack in reverse order for correct depth-first traversal
		children := current.Children()
		for i := len(children) - 1; i >= 0; i-- {
			stack = append(stack, children[i])
		}
	}

	return sb.String(), nil
}

// HasHandler checks if a handler is registered for the given node type
func (t *Traverser) HasHandler(nodeType string) bool {
	_, exists := t.handlers[nodeType]
	return exists
}

// RemoveHandler removes a handler for the given node type
func (t *Traverser) RemoveHandler(nodeType string) {
	delete(t.handlers, nodeType)
}

// HandlerTypes returns all registered node types
func (t *Traverser) HandlerTypes() []string {
	types := make([]string, 0, len(t.handlers))
	for nodeType := range t.handlers {
		types = append(types, nodeType)
	}
	return types
}
