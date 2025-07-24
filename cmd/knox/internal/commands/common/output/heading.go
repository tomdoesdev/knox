package output

import (
	"fmt"
	"strings"
)

type Heading struct {
	heading     string
	subHeadings []fmt.Stringer
}

func NewHeading(title string) *Heading {
	return &Heading{heading: title, subHeadings: make([]fmt.Stringer, 0)}
}
func (h *Heading) Add(v fmt.Stringer) {
	h.subHeadings = append(h.subHeadings, v)
}

func (h *Heading) String() string {
	sb := strings.Builder{}
	sb.WriteString(fmt.Sprintf("Heading<%s>\n", h.heading))

	for _, sub := range h.subHeadings {
		sb.WriteString(fmt.Sprintf("SubHeading<%s>\n", sub.String()))
	}

	return sb.String()
}
