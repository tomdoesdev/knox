package ast

import (
	"io"
)

type Marshaler interface {
	Marshal(root Node) ([]byte, error)
}
type StreamMarshaler interface {
	Marshal(w io.Writer, root Node) error
}

type MarshalerFunc func(Node) ([]byte, error)

func (rf MarshalerFunc) Marshal(root Node) ([]byte, error) {
	return rf(root)
}
