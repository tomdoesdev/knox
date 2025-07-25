package cli

import (
	"fmt"
	"strings"

	"github.com/tomdoesdev/knox/kit/ast"
)

// NewRenderer creates a new CLI renderer with standard formatting handlers
func NewRenderer() *ast.Traverser {
	return ast.NewTraverser().
		Handle("root", RootHandler).
		Handle("heading", HeadingHandler).
		Handle("list", ListHandler).
		Handle("listitem", ListItemHandler).
		Handle("text", TextHandler).
		Handle("section", SectionHandler)
}

// RootHandler handles rendering of root nodes (no-op)
func RootHandler(node ast.Node, sb *strings.Builder) error {
	// Root nodes don't add any output themselves, only their children
	return nil
}

// HeadingHandler handles rendering of heading nodes
func HeadingHandler(node ast.Node, sb *strings.Builder) error {
	sb.WriteString(fmt.Sprintf("%s:\n", node.Content().String()))
	return nil
}

// ListHandler handles rendering of list nodes (no-op container)
func ListHandler(node ast.Node, sb *strings.Builder) error {
	// List nodes don't add any output themselves, only their children
	return nil
}

// ListItemHandler handles rendering of list item nodes
func ListItemHandler(node ast.Node, sb *strings.Builder) error {
	sb.WriteString(fmt.Sprintf("  - %s\n", node.Content().String()))
	return nil
}

// TextHandler handles rendering of text nodes with spacing
func TextHandler(node ast.Node, sb *strings.Builder) error {
	sb.WriteString(fmt.Sprintf("%s\n\n", node.Content().String()))
	return nil
}

// SectionHandler handles rendering of section nodes
func SectionHandler(node ast.Node, sb *strings.Builder) error {
	sb.WriteString(fmt.Sprintf("=== %s ===\n", node.Content().String()))
	return nil
}

// WithAttributes creates a renderer that can handle styling attributes
func NewAttributeRenderer() *ast.Traverser {
	return ast.NewTraverser().
		Handle("root", RootHandler).
		Handle("heading", AttributeHeadingHandler).
		Handle("list", AttributeListHandler).
		Handle("listitem", AttributeListItemHandler).
		Handle("text", AttributeTextHandler).
		Handle("section", AttributeSectionHandler)
}

// AttributeHeadingHandler handles headings with potential level attributes
func AttributeHeadingHandler(node ast.Node, sb *strings.Builder) error {
	// Check for level attribute
	if level, exists := node.GetAttribute("level"); exists {
		if levelInt, ok := level.(int); ok && levelInt > 1 {
			// Sub-heading with indentation
			indent := strings.Repeat("  ", levelInt-1)
			sb.WriteString(fmt.Sprintf("%s%s:\n", indent, node.Content().String()))
			return nil
		}
	}
	// Default heading formatting
	return HeadingHandler(node, sb)
}

// AttributeListHandler handles lists with potential ordered attribute
func AttributeListHandler(node ast.Node, sb *strings.Builder) error {
	// Lists are containers, no direct output
	// The ordered attribute would be used by list items
	return nil
}

// AttributeListItemHandler handles list items with ordered list support
func AttributeListItemHandler(node ast.Node, sb *strings.Builder) error {
	// Check parent for ordered attribute (simplified version)
	sb.WriteString(fmt.Sprintf("  - %s\n", node.Content().String()))
	return nil
}

// AttributeTextHandler handles text with potential styling attributes
func AttributeTextHandler(node ast.Node, sb *strings.Builder) error {
	content := node.Content().String()

	// Check for styling attributes
	if bold, exists := node.GetAttribute("bold"); exists && bold.(bool) {
		content = fmt.Sprintf("**%s**", content) // Markdown-style bold
	}

	if color, exists := node.GetAttribute("color"); exists {
		// For now, just add color annotation (could be enhanced with ANSI codes later)
		content = fmt.Sprintf("[%s]%s", color, content)
	}

	sb.WriteString(fmt.Sprintf("%s\n\n", content))
	return nil
}

// AttributeSectionHandler handles sections with potential styling
func AttributeSectionHandler(node ast.Node, sb *strings.Builder) error {
	content := node.Content().String()

	// Check for level or importance attributes
	if level, exists := node.GetAttribute("level"); exists {
		if levelInt, ok := level.(int); ok {
			switch levelInt {
			case 1:
				sb.WriteString(fmt.Sprintf("=== %s ===\n", content))
			case 2:
				sb.WriteString(fmt.Sprintf("--- %s ---\n", content))
			default:
				sb.WriteString(fmt.Sprintf("* %s *\n", content))
			}
			return nil
		}
	}

	// Default section formatting
	return SectionHandler(node, sb)
}
