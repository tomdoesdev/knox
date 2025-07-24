package plain

import (
	"fmt"
	"strings"

	"github.com/tomdoesdev/knox/cmd/knox/internal/commands/common/output"
)

// RootFormatter handles rendering of root elements (no-op)
func RootFormatter(element output.Element, sb *strings.Builder) {
	// Root elements don't add any output themselves, only their children
}

// TextFormatter handles rendering of text elements
func TextFormatter(element output.Element, sb *strings.Builder) {
	sb.WriteString(element.Value().String())
}

// ListFormatter handles rendering of list elements (no-op container)
func ListFormatter(element output.Element, sb *strings.Builder) {
	// List elements don't add any output themselves, only their children
}

// ListItemFormatter handles rendering of list item elements
func ListItemFormatter(element output.Element, sb *strings.Builder) {
	sb.WriteString(fmt.Sprintf("  - %s\n", element.Value().String()))
}

// HeadingFormatter handles rendering of heading elements
func HeadingFormatter(element output.Element, sb *strings.Builder) {
	sb.WriteString(fmt.Sprintf("%s:\n", element.Value().String()))
}

// TextFormatter handles rendering of text elements with newline
func TextWithNewlineFormatter(element output.Element, sb *strings.Builder) {
	sb.WriteString(fmt.Sprintf("%s\n\n", element.Value().String()))
}
