# AST Module Migration

## Overview
Migrate from the tightly-coupled Element-based output system to a generic Abstract Syntax Tree (AST) module in the `kit/` package. This provides a reusable tree data structure that can be used throughout Knox for output formatting, configuration trees, data serialization, and more.

## Problem Statement
The current Element-based system has several limitations:
- **Tightly coupled to CLI output** - can't be reused for other tree structures
- **Limited semantic meaning** - `Value()` doesn't clearly express intent
- **Missing metadata support** - no way to attach attributes to nodes
- **Package location** - buried in commands/output instead of being a core utility
- **Single purpose** - designed only for output formatting

## Goals
1. **Generic Tree Structure**: Create a true AST that can represent any tree-based data
2. **Reusability**: Available throughout Knox ecosystem via `kit/ast` package
3. **Semantic Clarity**: Clear interfaces that express intent (`Content`, `Children`, `Attributes`)
4. **Metadata Support**: Ability to attach arbitrary attributes to nodes
5. **Multiple Use Cases**: Support CLI output, config trees, data serialization, etc.
6. **Backward Compatibility**: Maintain existing CLI output behavior during migration

## Architecture

### Core AST Interface

```go
package ast

import "fmt"

// Node represents a node in an abstract syntax tree
type Node interface {
    // Content returns the primary content/value of this node
    // Returns empty stringer for container-only nodes
    Content() fmt.Stringer
    
    // Children returns all child nodes
    // Returns empty slice for leaf nodes
    Children() []Node
    
    // Attributes returns metadata associated with this node
    // Can contain type information, styling hints, etc.
    Attributes() map[string]interface{}
    
    // Type returns a string identifier for the node type
    Type() string
}

// Mutable interface for nodes that can be modified after creation
type MutableNode interface {
    Node
    
    // SetContent updates the node's content
    SetContent(content fmt.Stringer)
    
    // AddChild appends a child node
    AddChild(child Node)
    
    // SetAttribute sets a metadata attribute
    SetAttribute(key string, value interface{})
}
```

### Concrete Node Implementation

```go
// BasicNode provides a simple implementation of MutableNode
type BasicNode struct {
    content    fmt.Stringer
    children   []Node
    attributes map[string]interface{}
    nodeType   string
}

// Constructor functions for common node types
func NewTextNode(content string) *BasicNode
func NewContainerNode(nodeType string) *BasicNode
func NewNode(nodeType string, content fmt.Stringer) *BasicNode
```

### Renderer Interface

```go
// Renderer converts AST nodes to string output
type Renderer interface {
    Render(root Node) (string, error)
}

// RendererFunc allows function-based renderers
type RendererFunc func(Node) (string, error)

// Traverser handles tree traversal with pluggable node handlers
type Traverser struct {
    handlers map[string]func(Node, *strings.Builder) error
}

func NewTraverser() *Traverser
func (t *Traverser) Handle(nodeType string, handler func(Node, *strings.Builder) error) *Traverser
func (t *Traverser) Render(root Node) (string, error)
```

## Migration Strategy

### Phase 1: Create AST Module
1. Implement `kit/ast` package with core interfaces
2. Create `BasicNode` implementation with common constructors
3. Implement `Traverser` with stack-based tree traversal
4. Add comprehensive tests

### Phase 2: Create Output Renderers
1. Create `cmd/knox/internal/commands/renderers/cli` package
2. Implement CLI-specific renderers using AST `Traverser`
3. Port existing plain formatters to AST node handlers
4. Ensure output matches existing behavior

### Phase 3: Update Element System
1. Create adapter layer from Element interface to AST Node
2. Update existing Element implementations to wrap AST nodes
3. Maintain backward compatibility for existing code

### Phase 4: Migrate Handlers
1. Update status handler to use AST directly instead of Elements
2. Verify output remains identical
3. Add attribute support for future enhancements (colors, etc.)

### Phase 5: Cleanup
1. Remove Element-based system once migration complete
2. Update documentation and examples
3. Remove old v1 and v2 output packages

### Phase 6: Expand Usage
1. Use AST for configuration file representation
2. Implement JSON/YAML renderers
3. Add support for other tree-structured data in Knox

## Implementation Details

### Node Type Semantics
- **"root"**: Container node representing document root
- **"heading"**: Mixed node with title content and child sections
- **"list"**: Container node for list items
- **"listitem"**: Leaf node containing list item content
- **"text"**: Leaf node containing plain text content
- **"section"**: Container node for grouping related content

### Attribute Usage Examples
```go
// Styling attributes
node.SetAttribute("color", "red")
node.SetAttribute("bold", true)

// Semantic attributes  
node.SetAttribute("level", 2) // for headings
node.SetAttribute("ordered", true) // for lists

// Custom data
node.SetAttribute("project_active", true)
node.SetAttribute("vault_count", 5)
```

### Renderer Implementation Example
```go
cliRenderer := ast.NewTraverser().
    Handle("root", func(node ast.Node, sb *strings.Builder) error {
        // Root nodes don't output anything themselves
        return nil
    }).
    Handle("heading", func(node ast.Node, sb *strings.Builder) error {
        sb.WriteString(fmt.Sprintf("%s:\n", node.Content().String()))
        return nil
    }).
    Handle("text", func(node ast.Node, sb *strings.Builder) error {
        sb.WriteString(fmt.Sprintf("%s\n\n", node.Content().String()))
        return nil
    })
```

## Benefits

### Immediate Benefits
1. **Decoupled Architecture**: AST separated from rendering concerns
2. **Reusable Foundation**: Available throughout Knox for any tree structure
3. **Enhanced Semantics**: Clear intent with `Content()`, `Children()`, `Attributes()`
4. **Metadata Support**: Rich attribute system for extensibility

### Future Capabilities
1. **Multiple Output Formats**: JSON, YAML, XML renderers
2. **Enhanced CLI Output**: Colors, styling, formatting options
3. **Configuration Trees**: Represent Knox config as AST
4. **Template Systems**: Use AST for template processing
5. **Data Serialization**: Generic tree serialization/deserialization
6. **Documentation Generation**: Generate docs from AST structures

## Success Criteria

- [ ] AST module created in `kit/ast` with comprehensive interface
- [ ] CLI renderer produces identical output to current system
- [ ] Status handler successfully migrated to use AST directly
- [ ] All existing tests pass with new AST system
- [ ] Performance equivalent or better than current Element system
- [ ] Clear documentation and examples for AST usage
- [ ] Attribute system working and tested
- [ ] Foundation ready for multiple renderer implementations

## Out of Scope

- **Advanced rendering features** (colors, etc.) - foundation only
- **Alternative output formats** - CLI renderer focus initially  
- **Configuration file parsing** - AST foundation only
- **Performance optimizations** - functional implementation priority
- **Breaking API changes** - maintain backward compatibility

## Dependencies

- No external dependencies required
- Uses existing `fmt.Stringer` interface
- Compatible with current Go standard library patterns
- Builds on existing Knox architectural patterns