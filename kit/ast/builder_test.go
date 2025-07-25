package ast

import (
	"testing"
)

func TestNewBuilder(t *testing.T) {
	builder := NewBuilder("root")
	tree := builder.Build()

	if tree.Type() != "root" {
		t.Errorf("Expected root type 'root', got '%s'", tree.Type())
	}

	if tree.Content().String() != "" {
		t.Errorf("Expected empty content, got '%s'", tree.Content().String())
	}

	if len(tree.Children()) != 0 {
		t.Errorf("Expected no children, got %d", len(tree.Children()))
	}

	if len(tree.Attributes()) != 0 {
		t.Errorf("Expected no attributes, got %d", len(tree.Attributes()))
	}
}

func TestBuilderContent(t *testing.T) {
	tree := NewBuilder("root").
		Content("test content").
		Build()

	if tree.Content().String() != "test content" {
		t.Errorf("Expected content 'test content', got '%s'", tree.Content().String())
	}
}

func TestBuilderAttributes(t *testing.T) {
	tree := NewBuilder("root").
		Attr("key1", "value1").
		Attr("key2", 42).
		Attr("key3", true).
		Build()

	attrs := tree.Attributes()

	if len(attrs) != 3 {
		t.Errorf("Expected 3 attributes, got %d", len(attrs))
	}

	if attrs["key1"] != "value1" {
		t.Errorf("Expected key1='value1', got '%v'", attrs["key1"])
	}

	if attrs["key2"] != 42 {
		t.Errorf("Expected key2=42, got '%v'", attrs["key2"])
	}

	if attrs["key3"] != true {
		t.Errorf("Expected key3=true, got '%v'", attrs["key3"])
	}

	// Test GetAttribute method
	value, exists := tree.GetAttribute("key1")
	if !exists || value != "value1" {
		t.Errorf("GetAttribute failed for key1")
	}

	_, exists = tree.GetAttribute("nonexistent")
	if exists {
		t.Error("GetAttribute should return false for nonexistent key")
	}
}

func TestSimpleTree(t *testing.T) {
	tree := NewBuilder("root").
		Content("root content").
		Node("child").
		Content("child content").
		Up().
		Build()

	// Check root
	if tree.Type() != "root" {
		t.Errorf("Expected root type 'root', got '%s'", tree.Type())
	}

	if tree.Content().String() != "root content" {
		t.Errorf("Expected root content 'root content', got '%s'", tree.Content().String())
	}

	children := tree.Children()
	if len(children) != 1 {
		t.Fatalf("Expected 1 child, got %d", len(children))
	}

	// Check child
	child := children[0]
	if child.Type() != "child" {
		t.Errorf("Expected child type 'child', got '%s'", child.Type())
	}

	if child.Content().String() != "child content" {
		t.Errorf("Expected child content 'child content', got '%s'", child.Content().String())
	}

	if len(child.Children()) != 0 {
		t.Errorf("Expected child to have no children, got %d", len(child.Children()))
	}
}

func TestDeepNesting(t *testing.T) {
	tree := NewBuilder("root").
		Node("level1").
		Node("level2").
		Node("level3").
		Content("deep content").
		Up().
		Up().
		Up().
		Build()

	// Navigate to deep node
	level1 := tree.Children()[0]
	level2 := level1.Children()[0]
	level3 := level2.Children()[0]

	if level3.Content().String() != "deep content" {
		t.Errorf("Expected deep content, got '%s'", level3.Content().String())
	}
}

func TestMultipleChildren(t *testing.T) {
	tree := NewBuilder("root").
		Node("child1").Content("first").Up().
		Node("child2").Content("second").Up().
		Node("child3").Content("third").Up().
		Build()

	children := tree.Children()
	if len(children) != 3 {
		t.Fatalf("Expected 3 children, got %d", len(children))
	}

	expectedContent := []string{"first", "second", "third"}
	expectedTypes := []string{"child1", "child2", "child3"}

	for i, child := range children {
		if child.Type() != expectedTypes[i] {
			t.Errorf("Child %d: expected type '%s', got '%s'", i, expectedTypes[i], child.Type())
		}

		if child.Content().String() != expectedContent[i] {
			t.Errorf("Child %d: expected content '%s', got '%s'", i, expectedContent[i], child.Content().String())
		}
	}
}

func TestRootNavigation(t *testing.T) {
	tree := NewBuilder("root").
		Node("child1").
		Node("grandchild").
		Up().
		Up().
		Root(). // Navigate back to root from anywhere
		Node("child2").
		Content("second child").
		Up().
		Build()

	children := tree.Children()
	if len(children) != 2 {
		t.Fatalf("Expected 2 children, got %d", len(children))
	}

	if children[1].Content().String() != "second child" {
		t.Errorf("Expected second child content, got '%s'", children[1].Content().String())
	}
}

func TestUpFromRoot(t *testing.T) {
	// Test that Up() from root is a no-op
	tree := NewBuilder("root").
		Up().Up().Up(). // Multiple Up() calls from root
		Content("still root").
		Build()

	if tree.Content().String() != "still root" {
		t.Errorf("Up() from root should be no-op, got '%s'", tree.Content().String())
	}
}

func TestComplexConfiguration(t *testing.T) {
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

	// Verify root attributes
	version, exists := config.GetAttribute("version")
	if !exists || version != "2.1" {
		t.Error("Config version attribute not set correctly")
	}

	environment, exists := config.GetAttribute("environment")
	if !exists || environment != "production" {
		t.Error("Config environment attribute not set correctly")
	}

	// Verify structure
	children := config.Children()
	if len(children) != 3 {
		t.Fatalf("Expected 3 main sections, got %d", len(children))
	}

	// Check database section
	database := children[0]
	if database.Type() != "database" {
		t.Errorf("Expected database section, got '%s'", database.Type())
	}

	dbChildren := database.Children()
	if len(dbChildren) != 4 {
		t.Errorf("Expected 4 database config items, got %d", len(dbChildren))
	}

	// Check server section
	server := children[1]
	if server.Type() != "server" {
		t.Errorf("Expected server section, got '%s'", server.Type())
	}

	// Check logging section
	logging := children[2]
	if logging.Type() != "logging" {
		t.Errorf("Expected logging section, got '%s'", logging.Type())
	}

	loggingLevel, exists := logging.GetAttribute("level")
	if !exists || loggingLevel != "info" {
		t.Error("Logging level attribute not set correctly")
	}
}

func TestBuilderWithTraversal(t *testing.T) {
	tree := NewBuilder("root").
		Node("child1").Content("first").Up().
		Node("child2").Content("second").Up().
		Node("child3").Content("third").Up().
		Build()

	// Test that builder-created trees work with traversal functions
	var visited []string
	err := WithBreadthFirstTraversal(tree, func(node Node) error {
		if node.Content().String() != "" {
			visited = append(visited, node.Content().String())
		}
		return nil
	})

	if err != nil {
		t.Fatalf("Traversal failed: %v", err)
	}

	expected := []string{"first", "second", "third"}
	if len(visited) != len(expected) {
		t.Fatalf("Expected %d nodes, got %d", len(expected), len(visited))
	}

	for i, expectedContent := range expected {
		if visited[i] != expectedContent {
			t.Errorf("Expected content '%s', got '%s'", expectedContent, visited[i])
		}
	}
}

func TestBuilderChaining(t *testing.T) {
	// Test that all methods return Builder for chaining
	var b Builder = NewBuilder("root")

	b = b.Content("test")
	b = b.Attr("key", "value")
	b = b.Node("child")
	b = b.Up()
	b = b.Root()

	tree := b.Build()

	if tree.Type() != "root" {
		t.Error("Builder chaining broke")
	}
}

func TestBuilderReuseAfterBuild(t *testing.T) {
	builder := NewBuilder("root").Content("original")
	tree1 := builder.Build()

	// Continue building after Build() - should still work
	tree2 := builder.Node("child").Up().Build()

	// Original tree should be unchanged
	if tree1.Content().String() != "original" {
		t.Error("Original tree was modified after Build()")
	}

	// New tree should have the child
	if len(tree2.Children()) != 1 {
		t.Error("Builder reuse after Build() failed")
	}
}

// Benchmark tests
func BenchmarkSimpleBuilder(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		tree := NewBuilder("root").
			Content("content").
			Node("child").
			Content("child content").
			Up().
			Build()
		_ = tree
	}
}

func BenchmarkComplexBuilder(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		tree := NewBuilder("config").
			Attr("version", "1.0").
			Node("database").
			Node("host").Content("localhost").Up().
			Node("port").Content("5432").Up().
			Up().
			Node("server").
			Node("port").Content("8080").Up().
			Up().
			Build()
		_ = tree
	}
}
