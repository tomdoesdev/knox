# Feature Spec: AST Attribute Type System

## Purpose & User Problem
The current attribute system uses `map[string]any` which requires manual type assertions throughout the codebase. This leads to verbose, error-prone code when accessing node attributes. Users need a type-safe, ergonomic way to work with attributes that provides clear error messages and handles type conversions gracefully.

## Success Criteria
- Type-safe attribute access with clear error messages
- Eliminate manual type assertions in consuming code
- Support common type conversions (string to bool, string to int, etc.)
- Maintain backward compatibility where possible
- Provide discoverable API through helper methods
- Centralized conversion logic for consistency

## Proposed API Design

### Core Attribute Type
```go
// Attribute wraps a value with type-safe accessor methods
type Attribute struct {
    value any
}

// Constructor functions
func NewAttribute(value any) Attribute
func StringAttr(value string) Attribute
func BoolAttr(value bool) Attribute
func IntAttr(value int) Attribute
func FloatAttr(value float64) Attribute

// Type-safe accessors
func (a Attribute) AsString() (string, error)
func (a Attribute) AsBool() (bool, error) 
func (a Attribute) AsInt() (int, error)
func (a Attribute) AsFloat() (float64, error)

// Convenience methods
func (a Attribute) AsStringOr(defaultValue string) string
func (a Attribute) AsBoolOr(defaultValue bool) bool
func (a Attribute) AsIntOr(defaultValue int) int

// Type checking
func (a Attribute) IsString() bool
func (a Attribute) IsBool() bool
func (a Attribute) IsInt() bool
func (a Attribute) IsFloat() bool

// Raw access for advanced cases
func (a Attribute) Value() any
```

### Updated Node Interface
```go
type Node interface {
    // ... existing methods ...
    
    // Updated attribute methods
    Attributes() map[string]Attribute
    GetAttribute(key string) (Attribute, bool)
}

type MutableNode interface {
    Node
    
    // ... existing methods ...
    
    // Updated setter - accepts any value, wraps in Attribute
    SetAttribute(key string, value any)
    // Alternative explicit setter
    SetAttr(key string, attr Attribute)
}
```

### Usage Examples

#### Basic Type-Safe Access
```go
// Current (verbose, error-prone):
if val, exists := node.GetAttribute("enabled"); exists {
    if enabled, ok := val.(bool); ok {
        fmt.Println("Enabled:", enabled)
    } else {
        fmt.Println("Error: not a boolean")
    }
}

// New (clean, type-safe):
if attr, exists := node.GetAttribute("enabled"); exists {
    if enabled, err := attr.AsBool(); err == nil {
        fmt.Println("Enabled:", enabled)
    } else {
        fmt.Printf("Error: %v\n", err)
    }
}
```

#### With Default Values
```go
// Get boolean with fallback
enabled := node.GetAttribute("enabled").AsBoolOr(false)

// Get string with fallback
name := node.GetAttribute("name").AsStringOr("untitled")
```

#### Type Conversion Support
```go
// String "true" -> bool true
attr := NewAttribute("true")
enabled, err := attr.AsBool() // Returns true, nil

// String "42" -> int 42
attr := NewAttribute("42") 
count, err := attr.AsInt() // Returns 42, nil

// Int 1 -> bool true
attr := NewAttribute(1)
flag, err := attr.AsBool() // Returns true, nil
```

## Implementation Approach

### Attribute Type Implementation
```go
type Attribute struct {
    value any
}

func NewAttribute(value any) Attribute {
    return Attribute{value: value}
}

func (a Attribute) AsString() (string, error) {
    switch v := a.value.(type) {
    case string:
        return v, nil
    case fmt.Stringer:
        return v.String(), nil
    case bool:
        return strconv.FormatBool(v), nil
    case int, int8, int16, int32, int64:
        return fmt.Sprintf("%d", v), nil
    case float32, float64:
        return fmt.Sprintf("%g", v), nil
    default:
        return "", fmt.Errorf("cannot convert %T to string", v)
    }
}

func (a Attribute) AsBool() (bool, error) {
    switch v := a.value.(type) {
    case bool:
        return v, nil
    case string:
        return strconv.ParseBool(v)
    case int, int8, int16, int32, int64:
        // Convert numeric: 0 = false, non-zero = true
        return reflect.ValueOf(v).Int() != 0, nil
    default:
        return false, fmt.Errorf("cannot convert %T to bool", v)
    }
}

func (a Attribute) AsInt() (int, error) {
    switch v := a.value.(type) {
    case int:
        return v, nil
    case int8, int16, int32, int64:
        return int(reflect.ValueOf(v).Int()), nil
    case float32, float64:
        f := reflect.ValueOf(v).Float()
        return int(f), nil
    case string:
        return strconv.Atoi(v)
    case bool:
        if v {
            return 1, nil
        }
        return 0, nil
    default:
        return 0, fmt.Errorf("cannot convert %T to int", v)
    }
}
```

### Node Interface Updates
```go
type mutableNode struct {
    // ... existing fields ...
    attributes map[string]Attribute // Changed from map[string]any
}

func (n *mutableNode) GetAttribute(key string) (Attribute, bool) {
    attr, exists := n.attributes[key]
    return attr, exists
}

func (n *mutableNode) SetAttribute(key string, value any) {
    n.attributes[key] = NewAttribute(value)
}

func (n *mutableNode) SetAttr(key string, attr Attribute) {
    n.attributes[key] = attr
}
```

### Builder Integration
```go
// Builder methods remain the same - auto-wrap values
func (b *nodeBuilder) Attr(key string, value any) Builder {
    b.current.attributes[key] = NewAttribute(value)
    return b
}

// Usage stays clean:
config := NewBuilder("config").
    Attr("enabled", true).          // bool
    Attr("port", 8080).            // int  
    Attr("name", "myapp").         // string
    Build()
```

## Migration Strategy

### Phase 1: Add Attribute Type (Non-Breaking)
- Add `Attribute` type alongside existing `any` system
- Provide conversion utilities
- Update internal implementations to use `Attribute`

### Phase 2: Update Interfaces (Breaking Change)
- Change return types to use `Attribute`
- Update all consuming code
- Provide migration guide

### Phase 3: Cleanup
- Remove legacy type assertion helpers
- Update documentation and examples

## Error Handling Design

### Clear Error Messages
```go
attr := NewAttribute([]string{"not", "convertible"})
_, err := attr.AsBool()
// Error: "cannot convert []string to bool"

attr := NewAttribute("not-a-number")
_, err := attr.AsInt() 
// Error: "strconv.Atoi: parsing \"not-a-number\": invalid syntax"
```

### Fallback Methods
```go
// Never errors - always returns a value
enabled := attr.AsBoolOr(false)
port := attr.AsIntOr(8080)
name := attr.AsStringOr("unknown")
```

## Benefits

### Developer Experience
- **Discoverable**: IDE shows available conversion methods
- **Type Safe**: Eliminates runtime type assertion panics
- **Consistent**: Same conversion logic everywhere
- **Clear Errors**: Helpful error messages for debugging

### Code Quality
- **Less Verbose**: No manual type assertions needed
- **More Readable**: Intent is clear from method names
- **Safer**: Handles edge cases centrally
- **Maintainable**: Conversion logic in one place

## Testing Strategy

### Unit Tests for Attribute Type
- Test all type conversions (string->bool, int->string, etc.)
- Test error cases with invalid conversions
- Test fallback methods with default values
- Test edge cases (nil values, empty strings, etc.)

### Integration Tests
- Test with Node implementations
- Test with Builder API
- Test with existing traversal functions
- Test serialization/deserialization if needed

## Backward Compatibility

### Migration Helpers
```go
// Helper for gradual migration
func GetLegacyAttribute(node Node, key string) (any, bool) {
    if attr, exists := node.GetAttribute(key); exists {
        return attr.Value(), true
    }
    return nil, false
}
```

### Documentation
- Clear migration guide with before/after examples
- Highlight benefits of new approach
- Provide automated refactoring suggestions where possible

## Out of Scope

- Complex type conversions (struct mapping, JSON parsing)
- Custom type conversion registration
- Attribute validation/schema enforcement
- Performance optimizations (can be added later)
- Serialization format changes (maintain JSON compatibility)

## Questions for Clarification

1. Should we support custom type converters for specialized types?
2. How should we handle nil values in attributes?
3. Should attribute keys be case-sensitive or case-insensitive?
4. Do you want strict type conversion or permissive (e.g., "1" -> true)?
5. Should we add attribute metadata (descriptions, validation rules)?