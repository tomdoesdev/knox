# Feature Spec: AST Node Builder

## Purpose & User Problem
Building complex AST trees manually requires verbose code with lots of temporary variables and manual parent-child linking. Users need a fluent, chainable API that makes tree construction intuitive and readable, similar to HTML builders or configuration DSLs.

## Success Criteria
- Fluent API with method chaining for intuitive tree building
- Support for setting node content, attributes, and hierarchical structure
- Navigation methods to move up/down the tree during construction
- Type-safe builder that prevents invalid tree states
- Clean, readable code that mirrors the tree structure visually
- Integration with existing Node interface

## Proposed API Design

### Core Builder Interface
```go
type Builder interface {
    // Node creation and navigation
    Node(nodeType string) Builder     // Add child node and navigate to it
    Up() Builder                      // Navigate back to parent node
    Root() Builder                    // Navigate back to root node
    
    // Content and attributes
    Content(content string) Builder   // Set content for current node
    Attr(key string, value interface{}) Builder // Set attribute on current node
    
    // Tree completion
    Build() Node                      // Return the built tree
}

// Entry point
func NewBuilder(rootType string) Builder
```

### Usage Examples

#### Simple Tree
```go
// Build: root -> child -> grandchild
tree := NewBuilder("root").
    Content("root content").
    Node("child").
        Content("child content").
        Node("grandchild").
            Content("leaf content").
        Up().
    Up().
    Build()
```

#### Complex Configuration Tree
```go
config := NewBuilder("config").
    Attr("version", "2.1").
    Attr("environment", "production").
    
    Node("database").
        Attr("type", "postgresql").
        Node("host").Content("localhost").Up().
        Node("port").Content("5432").Up().
        Node("name").Content("myapp_db").Up().
        Node("ssl").Content("true").Up().
    Up().
    
    Node("server").
        Node("host").Content("0.0.0.0").Up().
        Node("port").Content("8080").Up().
        Node("timeout").Content("30s").Up().
    Up().
    
    Node("logging").
        Attr("level", "info").
        Node("file").Content("/var/log/app.log").Up().
        Node("max_size").Content("100MB").Up().
    Up().
    
    Build()
```

#### Array-like Structures
```go
list := NewBuilder("list").
    Node("item").Content("first").Up().
    Node("item").Content("second").Up().
    Node("item").Content("third").Up().
    Build()
```

## Implementation Approach

### Builder State Management
```go
type nodeBuilder struct {
    root    *mutableNode    // Root of the tree being built
    current *mutableNode    // Current node being modified
    stack   []*mutableNode  // Parent stack for Up() navigation
}

type mutableNode struct {
    nodeType   string
    content    fmt.Stringer
    children   []Node
    attributes map[string]interface{}
    parent     *mutableNode // For navigation during building
}
```

### Core Methods Implementation
```go
func NewBuilder(rootType string) Builder {
    root := &mutableNode{
        nodeType:   rootType,
        content:    EmptyValue{},
        children:   make([]Node, 0),
        attributes: make(map[string]interface{}),
        parent:     nil,
    }
    
    return &nodeBuilder{
        root:    root,
        current: root,
        stack:   make([]*mutableNode, 0),
    }
}

func (b *nodeBuilder) Node(nodeType string) Builder {
    child := &mutableNode{
        nodeType:   nodeType,
        content:    EmptyValue{},
        children:   make([]Node, 0),
        attributes: make(map[string]interface{}),
        parent:     b.current,
    }
    
    // Add child to current node
    b.current.children = append(b.current.children, child)
    
    // Push current to stack and navigate to child
    b.stack = append(b.stack, b.current)
    b.current = child
    
    return b
}

func (b *nodeBuilder) Up() Builder {
    if len(b.stack) == 0 {
        return b // Already at root
    }
    
    // Pop parent from stack
    b.current = b.stack[len(b.stack)-1]
    b.stack = b.stack[:len(b.stack)-1]
    
    return b
}

func (b *nodeBuilder) Content(content string) Builder {
    b.current.content = StringValue(content)
    return b
}

func (b *nodeBuilder) Attr(key string, value interface{}) Builder {
    b.current.attributes[key] = value
    return b
}

func (b *nodeBuilder) Root() Builder {
    b.current = b.root
    b.stack = b.stack[:0] // Clear stack
    return b
}

func (b *nodeBuilder) Build() Node {
    return b.root
}
```

### Node Interface Implementation
```go
func (n *mutableNode) Type() string {
    return n.nodeType
}

func (n *mutableNode) Content() fmt.Stringer {
    return n.content
}

func (n *mutableNode) Children() []Node {
    return n.children
}

func (n *mutableNode) Attributes() map[string]interface{} {
    return n.attributes
}

func (n *mutableNode) GetAttribute(key string) (interface{}, bool) {
    value, exists := n.attributes[key]
    return value, exists
}
```

## Technical Considerations

### Navigation Stack
- Use a stack to track parent nodes for `Up()` operations
- Stack enables unlimited nesting depth
- `Root()` provides quick navigation back to start

### Memory Efficiency
- Builder only exists during construction
- Final `Build()` returns standard Node interface
- No circular references (parent links only during building)

### Error Handling
- `Up()` on root node is safe (no-op)
- Invalid operations don't panic, they're ignored gracefully
- Builder state is always valid

### Type Safety
- All methods return Builder interface for chaining
- Final Build() returns Node interface
- No way to create invalid tree states

## Benefits

### Readability
Code structure visually mirrors the tree structure through indentation and chaining.

### Efficiency
- Single pass construction
- No temporary variables needed
- Minimal memory overhead during building

### Flexibility
- Can build any tree structure
- Mix content and attributes freely
- Navigate freely during construction

## Integration with Existing Code

### Compatibility
Built nodes implement the existing `Node` interface, so they work with:
- All traversal functions (`WithBreadthFirstTraversal`, etc.)
- Marshaler implementations
- Any existing code that accepts `Node`

### Migration
Existing manual node construction can be replaced gradually:

**Before:**
```go
root := newMockNode("config", "")
db := newMockNode("database", "")
db.SetAttribute("type", "postgresql")
host := newMockNode("host", "localhost")
db.AddChild(host)
root.AddChild(db)
```

**After:**
```go
root := NewBuilder("config").
    Node("database").
        Attr("type", "postgresql").
        Node("host").Content("localhost").Up().
    Up().
    Build()
```

## Testing Strategy

### Unit Tests
- Test each builder method individually
- Verify navigation works correctly
- Test edge cases (multiple Up() calls, etc.)

### Integration Tests  
- Build complex trees and verify structure
- Test with traversal functions
- Test with marshalers

### Example Tests
```go
func TestSimpleBuilder(t *testing.T) {
    tree := NewBuilder("root").
        Content("test").
        Node("child").
            Content("child content").
        Up().
        Build()
    
    // Verify structure matches expectations
}

func TestComplexNavigation(t *testing.T) {
    tree := NewBuilder("root").
        Node("a").
            Node("b").
                Node("c").Up().
            Up().
        Up().
        Node("d").Up().
        Build()
    
    // Verify complex navigation worked correctly
}
```

## Out of Scope

- Builder validation (ensuring required fields are set)
- Builder serialization/deserialization  
- Thread safety (builders are single-use, single-threaded)
- Advanced navigation (sibling traversal, path queries)
- Conditional building (if/else logic)

## Implementation Plan

1. **Create mutableNode struct** that implements Node interface
2. **Implement nodeBuilder struct** with navigation stack
3. **Add all fluent API methods** with proper chaining
4. **Create comprehensive tests** for all functionality
5. **Add documentation** with usage examples
6. **Integration testing** with existing traversal/marshal code