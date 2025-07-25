# Feature Spec: AST Traversal Functions

## Purpose & User Problem
The current `Traverser` struct mixes traversal logic with marshaling concerns. We need clean, functional traversal utilities that eliminate boilerplate for walking AST trees, similar to Go's `filepath.Walk`. Users should be able to choose traversal order and provide their own visitor functions.

## Success Criteria
- Provide clean, functional traversal functions for different traversal orders
- Separate traversal logic from output generation/marshaling
- Easy to use API similar to `filepath.Walk`
- Predictable left-to-right child processing
- Support for common traversal patterns (pre-order, post-order, breadth-first)

## Proposed API Design

### Core Visitor Function Type
```go
// VisitorFunc is called for each node during traversal
type VisitorFunc func(node Node) error
```

### Traversal Functions
```go
// WithPreOrderTraversal visits parent before children (depth-first)
func WithPreOrderTraversal(root Node, visit VisitorFunc) error

// WithPostOrderTraversal visits children before parent (depth-first)  
func WithPostOrderTraversal(root Node, visit VisitorFunc) error

// WithBreadthFirstTraversal visits nodes level by level (queue-based)
func WithBreadthFirstTraversal(root Node, visit VisitorFunc) error
```

### Usage Examples
```go
// Count nodes by type
nodeCount := make(map[string]int)
err := WithPreOrderTraversal(root, func(node Node) error {
    nodeCount[node.Type()]++
    return nil
})

// Build text output
var sb strings.Builder
err := WithBreadthFirstTraversal(root, func(node Node) error {
    sb.WriteString(node.Content().String())
    sb.WriteString(" ")
    return nil
})

// Find specific nodes
var textNodes []Node
err := WithPreOrderTraversal(root, func(node Node) error {
    if node.Type() == "text" {
        textNodes = append(textNodes, node)
    }
    return nil
})
```

## Implementation Approaches

### Breadth-First (Queue-Based) - Recommended
```go
func WithBreadthFirstTraversal(root Node, visit VisitorFunc) error {
    if root == nil {
        return nil
    }
    
    queue := []Node{root}
    
    for len(queue) > 0 {
        current := queue[0]
        queue = queue[1:]
        
        // Visit current node
        if err := visit(current); err != nil {
            return err
        }
        
        // Add children to queue (preserves left-to-right order)
        queue = append(queue, current.Children()...)
    }
    
    return nil
}
```

### Pre-Order (Recursive)
```go
func WithPreOrderTraversal(root Node, visit VisitorFunc) error {
    if root == nil {
        return nil
    }
    
    return preOrderTraversal(root, visit)
}

func preOrderTraversal(node Node, visit VisitorFunc) error {
    // Visit current node first
    if err := visit(node); err != nil {
        return err
    }
    
    // Then visit children left-to-right
    for _, child := range node.Children() {
        if err := preOrderTraversal(child, visit); err != nil {
            return err
        }
    }
    
    return nil
}
```

### Post-Order (Recursive)
```go
func WithPostOrderTraversal(root Node, visit VisitorFunc) error {
    if root == nil {
        return nil
    }
    
    return postOrderTraversal(root, visit)
}

func postOrderTraversal(node Node, visit VisitorFunc) error {
    // Visit children first
    for _, child := range node.Children() {
        if err := postOrderTraversal(child, visit); err != nil {
            return err
        }
    }
    
    // Then visit current node
    return visit(node)
}
```

## Benefits of This Approach

### Separation of Concerns
- **Traversal**: Pure tree walking logic
- **Processing**: User-defined visitor functions
- **Marshaling**: Separate concern handled by dedicated marshalers

### Flexibility
- Users can choose appropriate traversal order for their use case
- Easy to compose different behaviors
- No need to manage handler registrations or state

### Go-Idiomatic
- Similar to `filepath.Walk`, `json.Walk` patterns
- Functional programming style
- Simple, predictable API

## Technical Considerations

### Error Handling
- Return early on first visitor error (fail-fast)
- Preserve original error context
- Input validation for nil nodes/visitors

### Performance
- Queue-based breadth-first is memory efficient for balanced trees
- Recursive approaches handle deep trees well
- No allocations for handler maps or state

## Migration from Current Traverser

### Old Code:
```go
traverser := NewTraverser()
traverser.Handle("text", func(node Node, sb *strings.Builder) error {
    sb.WriteString(node.Content().String())
    return nil
})
result, err := traverser.Marshal(root)
```

### New Code:
```go
var sb strings.Builder
err := WithBreadthFirstTraversal(root, func(node Node) error {
    if node.Type() == "text" {
        sb.WriteString(node.Content().String())
    }
    return nil
})
result := sb.String()
```

## Implementation Plan

1. **Create new traversal functions** in `traverser.go`
2. **Keep old Traverser struct** for backward compatibility (deprecated)
3. **Add comprehensive tests** for all traversal orders
4. **Update documentation** with usage examples
5. **Gradual migration** of existing code

## Out of Scope

- Complex traversal options (filters, max depth, context passing)
- Thread safety (traversal functions are stateless)
- Performance optimizations beyond basic implementation