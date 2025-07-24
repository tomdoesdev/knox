package output

import (
	"fmt"
	"strings"
)

type Printable interface {
	Add(v fmt.Stringer)
}
type Root struct {
	children []fmt.Stringer
}

func NewRoot() *Root {
	return &Root{children: make([]fmt.Stringer, 0)}
}

func (c *Root) Add(v fmt.Stringer) {
	c.children = append(c.children, v)
}

func (c *Root) String() string {
	sb := &strings.Builder{}

	for _, child := range c.children {
		sb.WriteString(child.String())
	}

	return sb.String()
}
