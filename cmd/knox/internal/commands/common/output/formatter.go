package output

import (
	"fmt"
	"strings"
)

// FormatterFunc is a function that renders an element to a string builder
type FormatterFunc func(Element, *strings.Builder)

// OutputFormatter manages formatter functions and handles tree traversal
type OutputFormatter struct {
	formatters map[string]FormatterFunc
}

// NewOutputFormatter creates a new OutputFormatter with an empty formatter map
func NewOutputFormatter() *OutputFormatter {
	return &OutputFormatter{
		formatters: make(map[string]FormatterFunc),
	}
}

// Using registers a formatter function for a specific element type
// Returns the formatter to enable method chaining
func (f *OutputFormatter) Using(elementType string, formatter FormatterFunc) *OutputFormatter {
	f.formatters[elementType] = formatter
	return f
}

// Render traverses the element tree and applies formatters to generate output
// Uses stack-based depth-first traversal
// Returns error if any element type doesn't have a registered formatter
func (f *OutputFormatter) Render(root Element) (string, error) {
	sb := &strings.Builder{}
	stack := []Element{root}

	for len(stack) > 0 {
		// Pop current element from stack
		current := stack[len(stack)-1]
		stack = stack[:len(stack)-1]

		// Find and apply formatter for current element
		formatter, exists := f.formatters[current.Type()]
		if !exists {
			return "", fmt.Errorf("no formatter registered for element type: %s", current.Type())
		}

		formatter(current, sb)

		// Add children to stack in reverse order for correct depth-first traversal
		children := current.Children()
		for i := len(children) - 1; i >= 0; i-- {
			stack = append(stack, children[i])
		}
	}

	return sb.String(), nil
}
