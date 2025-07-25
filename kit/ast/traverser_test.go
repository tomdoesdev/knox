package ast

import (
	"errors"
	"fmt"
	"strings"
	"testing"
)

// Mock Node implementation for testing
type mockNode struct {
	nodeType   string
	content    fmt.Stringer
	children   []Node
	attributes map[string]interface{}
}

func (m *mockNode) Type() string {
	return m.nodeType
}

func (m *mockNode) Content() fmt.Stringer {
	return m.content
}

func (m *mockNode) Children() []Node {
	return m.children
}

func (m *mockNode) Attributes() map[string]interface{} {
	return m.attributes
}

func (m *mockNode) GetAttribute(key string) (interface{}, bool) {
	value, exists := m.attributes[key]
	return value, exists
}

// Helper to create a mock node
func newMockNode(nodeType string, content string, children ...Node) *mockNode {
	return &mockNode{
		nodeType:   nodeType,
		content:    StringValue(content),
		children:   children,
		attributes: make(map[string]interface{}),
	}
}

// Create a test tree:
//
//	    root
//	   /    \
//	 left   right
//	/   \      \
//
// leaf1 leaf2  leaf3
func createTestTree() Node {
	leaf1 := newMockNode("leaf", "leaf1")
	leaf2 := newMockNode("leaf", "leaf2")
	leaf3 := newMockNode("leaf", "leaf3")

	left := newMockNode("branch", "left", leaf1, leaf2)
	right := newMockNode("branch", "right", leaf3)

	root := newMockNode("root", "root", left, right)

	return root
}

func TestWithBreadthFirstTraversal(t *testing.T) {
	root := createTestTree()

	var visited []string
	err := WithBreadthFirstTraversal(root, func(node Node) error {
		visited = append(visited, node.Content().String())
		return nil
	})

	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	// Breadth-first should visit level by level: root, left/right, leaf1/leaf2/leaf3
	expected := []string{"root", "left", "right", "leaf1", "leaf2", "leaf3"}

	if len(visited) != len(expected) {
		t.Fatalf("Expected %d nodes, got %d", len(expected), len(visited))
	}

	for i, expectedNode := range expected {
		if visited[i] != expectedNode {
			t.Errorf("Expected node %d to be %q, got %q", i, expectedNode, visited[i])
		}
	}
}

func TestWithPreOrderTraversal(t *testing.T) {
	root := createTestTree()

	var visited []string
	err := WithPreOrderTraversal(root, func(node Node) error {
		visited = append(visited, node.Content().String())
		return nil
	})

	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	// Pre-order should visit parent before children: root, left, leaf1, leaf2, right, leaf3
	expected := []string{"root", "left", "leaf1", "leaf2", "right", "leaf3"}

	if len(visited) != len(expected) {
		t.Fatalf("Expected %d nodes, got %d", len(expected), len(visited))
	}

	for i, expectedNode := range expected {
		if visited[i] != expectedNode {
			t.Errorf("Expected node %d to be %q, got %q", i, expectedNode, visited[i])
		}
	}
}

func TestWithPostOrderTraversal(t *testing.T) {
	root := createTestTree()

	var visited []string
	err := WithPostOrderTraversal(root, func(node Node) error {
		visited = append(visited, node.Content().String())
		return nil
	})

	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	// Post-order should visit children before parent: leaf1, leaf2, left, leaf3, right, root
	expected := []string{"leaf1", "leaf2", "left", "leaf3", "right", "root"}

	if len(visited) != len(expected) {
		t.Fatalf("Expected %d nodes, got %d", len(expected), len(visited))
	}

	for i, expectedNode := range expected {
		if visited[i] != expectedNode {
			t.Errorf("Expected node %d to be %q, got %q", i, expectedNode, visited[i])
		}
	}
}

func TestTraversalWithNilRoot(t *testing.T) {
	visitCalled := false

	// Test all traversal functions with nil root
	err := WithBreadthFirstTraversal(nil, func(node Node) error {
		visitCalled = true
		return nil
	})
	if err != nil {
		t.Errorf("Expected no error for nil root, got: %v", err)
	}
	if visitCalled {
		t.Error("Visitor should not be called for nil root")
	}

	err = WithPreOrderTraversal(nil, func(node Node) error {
		visitCalled = true
		return nil
	})
	if err != nil {
		t.Errorf("Expected no error for nil root, got: %v", err)
	}
	if visitCalled {
		t.Error("Visitor should not be called for nil root")
	}

	err = WithPostOrderTraversal(nil, func(node Node) error {
		visitCalled = true
		return nil
	})
	if err != nil {
		t.Errorf("Expected no error for nil root, got: %v", err)
	}
	if visitCalled {
		t.Error("Visitor should not be called for nil root")
	}
}

func TestTraversalWithNilVisitor(t *testing.T) {
	root := createTestTree()

	// Test all traversal functions with nil visitor
	err := WithBreadthFirstTraversal(root, nil)
	if err == nil || !strings.Contains(err.Error(), "visitor function cannot be nil") {
		t.Errorf("Expected nil visitor error, got: %v", err)
	}

	err = WithPreOrderTraversal(root, nil)
	if err == nil || !strings.Contains(err.Error(), "visitor function cannot be nil") {
		t.Errorf("Expected nil visitor error, got: %v", err)
	}

	err = WithPostOrderTraversal(root, nil)
	if err == nil || !strings.Contains(err.Error(), "visitor function cannot be nil") {
		t.Errorf("Expected nil visitor error, got: %v", err)
	}
}

func TestTraversalErrorHandling(t *testing.T) {
	root := createTestTree()
	expectedError := errors.New("visitor error")

	// Test that errors from visitor are propagated
	err := WithBreadthFirstTraversal(root, func(node Node) error {
		if node.Content().String() == "left" {
			return expectedError
		}
		return nil
	})
	if err != expectedError {
		t.Errorf("Expected error to be propagated, got: %v", err)
	}

	err = WithPreOrderTraversal(root, func(node Node) error {
		if node.Content().String() == "left" {
			return expectedError
		}
		return nil
	})
	if err != expectedError {
		t.Errorf("Expected error to be propagated, got: %v", err)
	}

	err = WithPostOrderTraversal(root, func(node Node) error {
		if node.Content().String() == "left" {
			return expectedError
		}
		return nil
	})
	if err != expectedError {
		t.Errorf("Expected error to be propagated, got: %v", err)
	}
}

func TestSingleNodeTraversal(t *testing.T) {
	singleNode := newMockNode("single", "only")

	// Test all traversal methods with single node
	var visited []string

	err := WithBreadthFirstTraversal(singleNode, func(node Node) error {
		visited = append(visited, node.Content().String())
		return nil
	})
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	if len(visited) != 1 || visited[0] != "only" {
		t.Errorf("Expected single visit to 'only', got: %v", visited)
	}

	visited = nil
	err = WithPreOrderTraversal(singleNode, func(node Node) error {
		visited = append(visited, node.Content().String())
		return nil
	})
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	if len(visited) != 1 || visited[0] != "only" {
		t.Errorf("Expected single visit to 'only', got: %v", visited)
	}

	visited = nil
	err = WithPostOrderTraversal(singleNode, func(node Node) error {
		visited = append(visited, node.Content().String())
		return nil
	})
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	if len(visited) != 1 || visited[0] != "only" {
		t.Errorf("Expected single visit to 'only', got: %v", visited)
	}
}

func TestLeftToRightOrdering(t *testing.T) {
	// Create a tree with many children to test left-to-right ordering
	child1 := newMockNode("child", "1")
	child2 := newMockNode("child", "2")
	child3 := newMockNode("child", "3")
	child4 := newMockNode("child", "4")

	root := newMockNode("root", "root", child1, child2, child3, child4)

	var visited []string
	err := WithBreadthFirstTraversal(root, func(node Node) error {
		visited = append(visited, node.Content().String())
		return nil
	})

	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	// Should visit root first, then children in left-to-right order
	expected := []string{"root", "1", "2", "3", "4"}

	if len(visited) != len(expected) {
		t.Fatalf("Expected %d nodes, got %d", len(expected), len(visited))
	}

	for i, expectedNode := range expected {
		if visited[i] != expectedNode {
			t.Errorf("Expected node %d to be %q, got %q", i, expectedNode, visited[i])
		}
	}
}

// Benchmark tests
func BenchmarkBreadthFirstTraversal(b *testing.B) {
	root := createTestTree()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		err := WithBreadthFirstTraversal(root, func(node Node) error {
			_ = node.Content().String() // Simulate some work
			return nil
		})
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkPreOrderTraversal(b *testing.B) {
	root := createTestTree()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		err := WithPreOrderTraversal(root, func(node Node) error {
			_ = node.Content().String() // Simulate some work
			return nil
		})
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkPostOrderTraversal(b *testing.B) {
	root := createTestTree()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		err := WithPostOrderTraversal(root, func(node Node) error {
			_ = node.Content().String() // Simulate some work
			return nil
		})
		if err != nil {
			b.Fatal(err)
		}
	}
}
