package output

import (
	"fmt"
	"strings"
)

type List struct {
	items []fmt.Stringer
}

func NewList() *List {
	return &List{items: make([]fmt.Stringer, 0)}
}

func (l *List) Add(v fmt.Stringer) {
	l.items = append(l.items, v)
}

func (l *List) String() string {
	sb := &strings.Builder{}
	sb.WriteString("List<\n")

	for _, child := range l.items {
		sb.WriteString(fmt.Sprintf("\tListItem<%s>\n", child))
	}

	sb.WriteString(">")
	return sb.String()
}
