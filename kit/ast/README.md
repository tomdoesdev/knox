# Knox AST (Abstract Syntax Tree)

A generic tree data structure for representing and manipulating hierarchical data in Knox. The AST module provides a clean separation between content structure and presentation, enabling multiple output formats and rich metadata support.

## Overview

The AST module consists of three main components:

- **Nodes**: Represent tree structure with content, children, and attributes
- **Traverser**: Handles tree traversal with pluggable rendering logic
- **Renderers**: Convert AST trees to specific output formats

## Core Interfaces

### Node Interface

```go
type Node interface {
    Content() fmt.Stringer              // Primary content/value
    Children() []Node                   // Child nodes
    Attributes() map[string]interface{} // Metadata
    GetAttribute(key string) (interface{}, bool) // Get specific attribute
    Type() string                       // Node type identifier
}
```

### MutableNode Interface

```go
type MutableNode interface {
    Node
    SetContent(content fmt.Stringer)         // Update content
    AddChild(child Node)                     // Add child node
    SetAttribute(key string, value interface{}) // Set metadata
    GetAttribute(key string) (interface{}, bool) // Get metadata
}
```

## Creating Nodes

### Basic Constructors

```go
// Text node (leaf)
text := ast.NewTextNode("Hello, World!")

// Container node
container := ast.NewContainerNode("section")

// Node with content and type
node := ast.NewNode("custom", ast.StringValue("content"))
```

### Convenience Constructors

```go
// Common node types
root := ast.NewRoot()
heading := ast.NewHeading("My Title")
list := ast.NewList()
item := ast.NewListItem("Item 1")
text := ast.Text("Simple text")  // Alias for NewTextNode
```

## Building Trees

```go
// Create a document structure
root := ast.NewRoot()

// Add a heading with content
heading := ast.NewHeading("Projects")
heading.AddChild(ast.Text("Current workspace projects:"))

// Create a list of items
list := ast.NewList()
list.AddChild(ast.NewListItem("project-alpha"))
list.AddChild(ast.NewListItem("project-beta"))

// Add metadata to nodes
activeItem := ast.NewListItem("project-gamma")
activeItem.SetAttribute("active", true)
activeItem.SetAttribute("priority", "high")
list.AddChild(activeItem)

// Build the tree
heading.AddChild(list)
root.AddChild(heading)
```

## Using Attributes

Attributes provide rich metadata support for styling, semantic information, and custom data:

```go
// Styling attributes
node.SetAttribute("color", "red")
node.SetAttribute("bold", true)
node.SetAttribute("level", 2)

// Semantic attributes
node.SetAttribute("active", true)
node.SetAttribute("type", "warning")
node.SetAttribute("count", 42)

// Custom application data
node.SetAttribute("project_id", "abc-123")
node.SetAttribute("last_modified", time.Now())
```

## Rendering with Traverser

### Basic Traverser Usage

```go
// Create a traverser with handlers
traverser := ast.NewTraverser().
    Handle("root", func(node ast.Node, sb *strings.Builder) error {
        // Root nodes typically don't output anything
        return nil
    }).
    Handle("heading", func(node ast.Node, sb *strings.Builder) error {
        sb.WriteString(fmt.Sprintf("# %s\n", node.Content().String()))
        return nil
    }).
    Handle("text", func(node ast.Node, sb *strings.Builder) error {
        sb.WriteString(node.Content().String())
        return nil
    })

// Render the tree
output, err := traverser.Render(root)
if err != nil {
    log.Fatal(err)
}
fmt.Print(output)
```

### Attribute-Aware Handlers

```go
traverser.Handle("listitem", func(node ast.Node, sb *strings.Builder) error {
    content := node.Content().String()
    
    // Check for active attribute
    if active, exists := node.GetAttribute("active"); exists && active.(bool) {
        content = fmt.Sprintf("%s [ACTIVE]", content)
    }
    
    // Check for priority styling
    if priority, exists := node.GetAttribute("priority"); exists {
        switch priority.(string) {
        case "high":
            content = fmt.Sprintf("!! %s", content)
        case "low":
            content = fmt.Sprintf("- %s", content)
        }
    }
    
    sb.WriteString(fmt.Sprintf("  • %s\n", content))
    return nil
})
```

## Pre-built Renderers

### CLI Renderer

Knox includes a CLI renderer for terminal output:

```go
import "github.com/tomdoesdev/knox/cmd/knox/internal/commands/renderers/cli"

renderer := cli.NewRenderer()
output, err := renderer.Render(root)
```

### Creating Custom Renderers

```go
// Markdown renderer example
func NewMarkdownRenderer() *ast.Traverser {
    return ast.NewTraverser().
        Handle("root", func(node ast.Node, sb *strings.Builder) error {
            return nil // Root is invisible
        }).
        Handle("heading", func(node ast.Node, sb *strings.Builder) error {
            level := 1
            if l, exists := node.GetAttribute("level"); exists {
                level = l.(int)
            }
            prefix := strings.Repeat("#", level)
            sb.WriteString(fmt.Sprintf("%s %s\n\n", prefix, node.Content()))
            return nil
        }).
        Handle("listitem", func(node ast.Node, sb *strings.Builder) error {
            sb.WriteString(fmt.Sprintf("- %s\n", node.Content()))
            return nil
        })
}
```

## Common Patterns

### Status Display

```go
func renderStatus(projects []string, activeProject string) string {
    root := ast.NewRoot()
    
    heading := ast.NewHeading("Projects")
    list := ast.NewList()
    
    for _, project := range projects {
        item := ast.NewListItem(project)
        if project == activeProject {
            item.SetAttribute("active", true)
        }
        list.AddChild(item)
    }
    
    heading.AddChild(list)
    root.AddChild(heading)
    
    renderer := cli.NewRenderer()
    output, _ := renderer.Render(root)
    return output
}
```

### Configuration Trees

```go
func buildConfigTree(config map[string]interface{}) ast.Node {
    root := ast.NewContainerNode("config")
    
    for key, value := range config {
        node := ast.NewNode("setting", ast.StringValue(fmt.Sprintf("%v", value)))
        node.SetAttribute("key", key)
        node.SetAttribute("type", reflect.TypeOf(value).String())
        root.AddChild(node)
    }
    
    return root
}
```

### Error Handling

```go
// Traverser returns errors for missing handlers
traverser := ast.NewTraverser().
    Handle("known-type", handler)

_, err := traverser.Render(nodeWithUnknownType)
if err != nil {
    // Error: "no handler registered for node type: unknown-type"
    log.Printf("Rendering failed: %v", err)
}

// Check for handler existence
if !traverser.HasHandler("some-type") {
    log.Println("Handler missing for some-type")
}
```

## Best Practices

### Node Types

Use consistent, descriptive node type names:

```go
// Good
"heading", "list", "listitem", "text", "section", "code", "table"

// Less clear
"h1", "ul", "li", "p", "div"
```

### Attributes

- Use lowercase keys with underscores: `"is_active"`, `"line_number"`
- Store semantic data, not presentation: `"priority": "high"` not `"color": "red"`
- Use consistent value types: booleans for flags, strings for enums

### Error Handling

Always handle traverser errors:

```go
output, err := renderer.Render(root)
if err != nil {
    return fmt.Errorf("failed to render AST: %w", err)
}
```

### Performance

- Reuse traversers when possible - they're stateless
- Avoid deeply nested trees (>10 levels) for better performance
- Pre-allocate string builders for large outputs

## Testing

```go
func TestAST(t *testing.T) {
    // Create test tree
    root := ast.NewRoot()
    root.AddChild(ast.NewHeading("Test"))
    
    // Create simple renderer
    renderer := ast.NewTraverser().
        Handle("root", func(node ast.Node, sb *strings.Builder) error {
            return nil
        }).
        Handle("heading", func(node ast.Node, sb *strings.Builder) error {
            sb.WriteString(node.Content().String())
            return nil
        })
    
    // Test rendering
    output, err := renderer.Render(root)
    assert.NoError(t, err)
    assert.Equal(t, "Test", output)
}
```

## Future Extensions

The AST module is designed for extensibility:

- **JSON/YAML renderers** for data serialization
- **Color/styling renderers** for enhanced CLI output  
- **Template renderers** for document generation
- **Validation** for tree structure constraints
- **Transformation utilities** for tree manipulation

## Package Structure

```
kit/ast/
├── README.md           # This file
├── node.go            # Core interfaces
├── basic_node.go      # BasicNode implementation  
├── traverser.go       # Tree traversal and rendering
└── constructors.go    # Convenience constructors
```