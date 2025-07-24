# DOM-Like Output System Refactor

## Overview
Refactor the current component-based output system to use a DOM-like tree structure with pluggable formatters. This will provide better separation of concerns between content structure and presentation, enabling multiple output formats and easier styling.

## Problem Statement
The current output system has formatting logic embedded in each component's `String()` method. This makes it difficult to:
- Support multiple output formats (plain text, colored, JSON, etc.)
- Modify styling without changing component logic
- Test formatting independently from structure
- Add new element types without duplicating formatting patterns

## Goals
1. **Separation of Concerns**: Decouple content structure from presentation
2. **Pluggable Formatters**: Enable multiple rendering strategies
3. **Extensibility**: Make it easy to add new element types and formatters
4. **Backward Compatibility**: Maintain existing handler usage patterns where possible
5. **Error Safety**: Explicit handling of missing formatters

## Architecture

### Core Interfaces

```go
type Element interface {
    Value() fmt.Stringer    // Raw content (empty for container elements)
    Children() []Element    // Child elements
    Type() string          // Element type identifier ("text", "list", "heading")
}

type OutputFormatter struct {
    formatters map[string]func(Element, *strings.Builder)
}
```

### Element Types
- **Container Elements**: No value, only children (List, Heading)
- **Leaf Elements**: Have value, no children (Text, ListItem)
- **Mixed Elements**: Have both value and children (Heading with text and sub-elements)

### Formatting System
- **OutputFormatter**: Manages formatter functions and tree traversal
- **Formatter Functions**: `func(Element, *strings.Builder)` - handle single element rendering
- **Registration**: Builder pattern for registering formatters
- **Traversal**: Stack-based depth-first traversal

## Implementation Plan

### Phase 1: Core System
1. Create new `Element` interface and `OutputFormatter` struct
2. Implement stack-based tree traversal in `Render()` method
3. Add builder pattern for formatter registration
4. Create error handling for missing formatters

### Phase 2: Element Implementations
1. Create concrete element types (Text, List, ListItem, Heading)
2. Implement builder pattern methods (`Add()`, constructors)
3. Ensure proper parent-child relationships

### Phase 3: Built-in Formatters
1. Create `formatters/plain` package with basic formatters
2. Implement formatters for each element type
3. Ensure output matches current system behavior

### Phase 4: Migration
1. Update existing handlers to use new system
2. Verify output remains unchanged
3. Remove old component system

### Phase 5: Testing & Documentation
1. Add comprehensive tests for tree construction and rendering
2. Test error cases (missing formatters)
3. Document formatter creation patterns

## Technical Details

### OutputFormatter Implementation
```go
func NewOutputFormatter() *OutputFormatter
func (f *OutputFormatter) Using(elementType string, formatter func(Element, *strings.Builder)) *OutputFormatter  
func (f *OutputFormatter) Render(root Element) (string, error)
```

### Tree Traversal Algorithm
- Use stack-based depth-first traversal
- Apply formatter to current element
- Add children to stack in reverse order for correct processing
- Return error if formatter not found for element type

### Element Construction
```go
// Maintain familiar API
root := output.NewRoot()
heading := output.NewHeading("Title")
list := output.NewList()
text := output.NewText("content")

// Builder pattern preserved
root.Add(heading)
list.Add(text)
```

### Formatter Registration
```go
formatter := NewOutputFormatter().
    Using("text", plain.TextFormatter).
    Using("list", plain.ListFormatter).
    Using("heading", plain.HeadingFormatter)

output, err := formatter.Render(root)
```

## Migration Strategy

1. **Parallel Implementation**: Build new system alongside existing one
2. **Handler Migration**: Update one handler at a time to use new system
3. **Output Verification**: Ensure identical output during migration
4. **Gradual Removal**: Remove old system after all handlers migrated

## Benefits

1. **Multiple Output Formats**: Easy to add JSON, colored, or other formatters
2. **Testability**: Format logic can be tested independently
3. **Maintainability**: Clear separation between structure and presentation
4. **Extensibility**: New element types require no changes to formatting pipeline
5. **Flexibility**: Runtime formatter selection and modification

## Success Criteria

- [ ] All existing handlers produce identical output using new system
- [ ] New element types can be added without modifying OutputFormatter
- [ ] Multiple formatters can be swapped at runtime
- [ ] Clear error messages for missing formatters
- [ ] Performance equivalent to current system
- [ ] Clean, documented API for creating custom formatters

## Out of Scope

- Advanced styling features (colors, etc.) - will be added later
- Complex layout algorithms - keeping simple tree rendering
- Performance optimizations beyond basic implementation
- Alternative output formats (JSON, etc.) - foundation only

## Dependencies

- No external dependencies required
- Uses existing `fmt.Stringer` interface
- Compatible with current Go standard library usage