package ast

import (
	"fmt"
)

// VisitorFunc is called for each node during traversal
type VisitorFunc func(node Node) error

// WithBreadthFirstTraversal visits nodes level by level (queue-based)
// Children are processed in left-to-right order as they appear in the Children() slice
func WithBreadthFirstTraversal(root Node, visit VisitorFunc) error {
	if root == nil {
		return nil
	}
	if visit == nil {
		return fmt.Errorf("visitor function cannot be nil")
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

// WithPreOrderTraversal visits parent before children (depth-first)
// Children are processed in left-to-right order
func WithPreOrderTraversal(root Node, visit VisitorFunc) error {
	if root == nil {
		return nil
	}
	if visit == nil {
		return fmt.Errorf("visitor function cannot be nil")
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

// WithPostOrderTraversal visits children before parent (depth-first)
// Children are processed in left-to-right order
func WithPostOrderTraversal(root Node, visit VisitorFunc) error {
	if root == nil {
		return nil
	}
	if visit == nil {
		return fmt.Errorf("visitor function cannot be nil")
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
