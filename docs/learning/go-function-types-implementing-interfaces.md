# Go: Function Types Implementing Interfaces

## Overview

In Go, you can create function types that implement interfaces. This is a powerful pattern that allows functions to be used anywhere an interface is expected, without needing to wrap them in structs.

## How It Works

```go
// 1. Define a function type
type MarshalerFunc func(Node) ([]byte, error)

// 2. Add a method to that function type
func (rf MarshalerFunc) Marshal(root Node) ([]byte, error) {
    return rf(root)  // Call the function itself
}
```

Now `MarshalerFunc` implements any interface that requires a `Marshal(Node) ([]byte, error)` method.

## Step-by-Step Breakdown

### 1. Function Type Definition
```go
type MarshalerFunc func(Node) ([]byte, error)
```
This creates a new type based on a function signature.

### 2. Method on Function Type
```go
func (rf MarshalerFunc) Marshal(root Node) ([]byte, error) {
    return rf(root)  // rf is the function itself
}
```
This adds a method to the function type. Inside the method, `rf` refers to the actual function.

### 3. Interface Implementation
```go
type Marshaler interface {
    Marshal(root Node) ([]byte, error)
}
```
Since `MarshalerFunc` has a `Marshal` method with the right signature, it automatically implements the `Marshaler` interface.

## Practical Examples

### Basic Usage
```go
// Regular function
func simpleTextMarshaler(node Node) ([]byte, error) {
    return []byte(node.Content().String()), nil
}

// Convert function to implement interface
var marshaler Marshaler = MarshalerFunc(simpleTextMarshaler)

// Use it
data, err := marshaler.Marshal(myNode)
```

### Inline Function
```go
// Create marshaler inline
jsonMarshaler := MarshalerFunc(func(node Node) ([]byte, error) {
    return json.Marshal(map[string]interface{}{
        "type":    node.Type(),
        "content": node.Content().String(),
    })
})

// Use anywhere a Marshaler is expected
data, err := jsonMarshaler.Marshal(myNode)
```

### Comparison: Function Type vs Struct

**Without function types (more verbose):**
```go
type SimpleMarshaler struct{}

func (s SimpleMarshaler) Marshal(node Node) ([]byte, error) {
    return []byte(node.Content().String()), nil
}

marshaler := SimpleMarshaler{}
```

**With function types (concise):**
```go
marshaler := MarshalerFunc(func(node Node) ([]byte, error) {
    return []byte(node.Content().String()), nil
})
```

## Standard Library Examples

This pattern is used extensively in Go's standard library:

### http.HandlerFunc
```go
type HandlerFunc func(ResponseWriter, *Request)

func (f HandlerFunc) ServeHTTP(w ResponseWriter, r *Request) {
    f(w, r)
}

// Usage:
http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, "Hello, World!")
})
```

### sort.Reverse
```go
type reverse struct{ Interface }

func Reverse(data Interface) Interface {
    return &reverse{data}
}

// But also:
type StringSlice []string
func (p StringSlice) Len() int           { return len(p) }
func (p StringSlice) Less(i, j int) bool { return p[i] < p[j] }
func (p StringSlice) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }
```

## Benefits

1. **Lightweight**: No need to create structs for simple cases
2. **Flexible**: Easy to create implementations inline
3. **Composable**: Functions can be easily wrapped and combined
4. **Familiar**: Reduces boilerplate for simple interface implementations

## When to Use

- **Simple interfaces** with one primary method
- **Functional programming** patterns
- **Adapter patterns** where you want to make a function implement an interface
- **Quick implementations** for testing or prototyping

## When NOT to Use

- **Complex interfaces** with multiple methods
- **Stateful implementations** that need to store data
- **When you need multiple different methods** on the same type

## Real-World Use Case: Multiple Marshalers

```go
// Different marshaling strategies
var (
    textMarshaler = MarshalerFunc(func(node Node) ([]byte, error) {
        return []byte(node.Content().String()), nil
    })
    
    jsonMarshaler = MarshalerFunc(func(node Node) ([]byte, error) {
        return json.Marshal(map[string]interface{}{
            "type":    node.Type(),
            "content": node.Content().String(),
            "attrs":   node.Attributes(),
        })
    })
    
    yamlMarshaler = MarshalerFunc(func(node Node) ([]byte, error) {
        return yaml.Marshal(map[string]interface{}{
            "type":    node.Type(),
            "content": node.Content().String(),
            "attrs":   node.Attributes(),
        })
    })
)

// All can be used interchangeably
func processWithMarshaler(m Marshaler, node Node) {
    data, err := m.Marshal(node)
    // ...
}
```

## Key Takeaway

Function types implementing interfaces allow you to treat functions as first-class objects that satisfy interface contracts. This creates a bridge between functional and object-oriented programming styles in Go, making the language more expressive and flexible.